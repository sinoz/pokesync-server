package login

import (
	"fmt"
	"reflect"
	"time"

	"gitlab.com/pokesync/game-service/internal/game-service/account"
	"gitlab.com/pokesync/game-service/internal/game-service/client"
	"gitlab.com/pokesync/game-service/internal/game-service/game"
)

// Unbounded is for parameters such as the job limit.
const Unbounded = -1

// Config holds configurations for the login service.
type Config struct {
	JobLimit          int
	JobConsumeTimeout time.Duration

	WorkerCount int
}

// Job is a login job picked up and processed by a worker.
type Job struct {
	Client  *client.Client
	Request Request
}

// Service is an implementation of a login service providing
// login and authentication capabilities.
type Service struct {
	config        Config
	jobQueue      chan Job
	authenticator Authenticator
	routing       *client.Router
}

// NewService constructs a new login Service.
func NewService(config Config, authenticator Authenticator, routing *client.Router) *Service {
	var jobQueue chan Job
	if config.JobLimit == Unbounded {
		jobQueue = make(chan Job)
	} else {
		jobQueue = make(chan Job, config.JobLimit)
	}

	service := &Service{
		config:        config,
		jobQueue:      jobQueue,
		authenticator: authenticator,
		routing:       routing,
	}

	mailbox := routing.Subscribe("login_request")
	service.receiver(mailbox)

	for i := 0; i < config.WorkerCount; i++ {
		service.spawnWorker()
	}

	return service
}

// receiver receives and handles client messages from the specified mailbox.
func (service *Service) receiver(mailbox client.Mailbox) {
	go func() {
		for mail := range mailbox {
			switch message := mail.Payload.(type) {
			case *Request:
				service.handleRequest(mail.Client, message)

			default:
				fmt.Println(reflect.TypeOf(message))
			}
		}
	}()
}

// handleRequest handles the given login Request for the given Client. If no
// worker was able to pick up the job within a specific time frame, a timeout
// occurs and the client's request is denied.
func (service *Service) handleRequest(client *client.Client, request *Request) {
	job := Job{Client: client, Request: *request}

	select {
	case service.jobQueue <- job:
		// job has been picked up by a worker.

	case <-time.After(service.config.JobConsumeTimeout):
		// no worker was able to pick up the job within the time frame.
		// we have to notify the client of this failure.
		client.SendNow(&RequestTimedOut{})
		client.Terminate()
	}
}

// spawnWorker spawns a worker goroutine that reads from the
// service's job queue.
func (service *Service) spawnWorker() {
	go func() {
		for job := range service.jobQueue {
			email := account.Email(job.Request.Email)
			password := account.Password(job.Request.Password)

			if !email.Validate() || !password.Validate() {
				job.Client.SendNow(&InvalidCredentials{})
				job.Client.Terminate()

				continue
			}

			result, err := service.authenticator.Authenticate(email, password)
			if err != nil {
				job.Client.SendNow(&ErrorDuringAccountFetch{})
				job.Client.Terminate()

				continue
			}

			switch res := result.(type) {
			case AuthSuccess:
				service.routing.Publish(game.AuthenticationEventTopic, client.Mail{
					Client:  job.Client,
					Payload: game.Authenticated{Account: res.Account},
				})

			case CouldNotFindAccount, PasswordMismatch:
				job.Client.SendNow(&InvalidCredentials{})
				job.Client.Terminate()

				continue
			}
		}
	}()
}

// TearDown tears down this service, closing its job queue and
// cleaning up any resources it is holding.
func (service *Service) TearDown() {
	close(service.jobQueue)
}
