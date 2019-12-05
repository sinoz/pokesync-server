package login

import (
	"context"
	"reflect"

	"gitlab.com/pokesync/game-service/internal/game-service/account"
	"gitlab.com/pokesync/game-service/internal/game-service/client"
	"gitlab.com/pokesync/game-service/internal/game-service/game"
	"go.uber.org/zap"
)

// Config holds configurations for the login service.
type Config struct {
	Logger      *zap.SugaredLogger
	WorkerCount int
}

// Job is a login job picked up and processed by a worker.
type Job struct {
	Context context.Context
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
	jobQueue := make(chan Job)

	service := &Service{
		config:        config,
		jobQueue:      jobQueue,
		authenticator: authenticator,
		routing:       routing,
	}

	mailbox := routing.Subscribe("login_request")
	go service.receiver(mailbox)

	for i := 0; i < config.WorkerCount; i++ {
		go service.spawnWorker()
	}

	return service
}

// receiver receives and handles client messages from the specified mailbox.
func (service *Service) receiver(mailbox client.Mailbox) {
	for mail := range mailbox {
		switch message := mail.Payload.(type) {
		case *Request:
			service.queueRequest(mail.Context, mail.Client, message)
			break

		default:
			service.config.Logger.Errorf("unexpected message received of type %v", reflect.TypeOf(message))
		}
	}
}

// queueRequest buffers the given login Request of the given Client for processing.
func (service *Service) queueRequest(ctx context.Context, client *client.Client, request *Request) {
	service.jobQueue <- Job{Context: ctx, Client: client, Request: *request}
}

// spawnWorker spawns a worker goroutine that reads from the
// service's job queue.
func (service *Service) spawnWorker() {
	for job := range service.jobQueue {
		email := account.Email(job.Request.Email)
		password := account.Password(job.Request.Password)

		if !email.Validate() || !password.Validate() {
			job.Client.SendNow(&InvalidCredentials{})
			job.Client.Terminate()

			continue
		}

		result, err := service.authenticator.Authenticate(job.Context, email, password)
		if err != nil {
			job.Client.SendNow(&ErrorDuringAccountFetch{})
			job.Client.Terminate()

			continue
		}

		switch res := result.(type) {
		case AuthSuccess:
			service.routing.Publish(game.AuthenticationEventTopic, client.Mail{
				Context: job.Context,
				Client:  job.Client,
				Payload: game.Authenticated{Account: res.Account},
			})

		case TimedOut:
			job.Client.SendNow(&RequestTimedOut{})
			job.Client.Terminate()

			continue

		case CouldNotFindAccount, PasswordMismatch:
			job.Client.SendNow(&InvalidCredentials{})
			job.Client.Terminate()

			continue

		default:
			service.config.Logger.Errorf("unexpected authentication result type of %v", reflect.TypeOf(res))
		}
	}
}

// TearDown tears down this service, closing its job queue and
// cleaning up any resources it is holding.
func (service *Service) TearDown() {
	close(service.jobQueue)
}
