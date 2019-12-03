package account

import (
	"reflect"

	"go.uber.org/zap"
)

// Config holds configurations specific to the account Service.
type Config struct {
	WorkerCount int
	Logger      *zap.SugaredLogger
}

// LoadResult is the result from attempting to load an account.
type LoadResult struct {
	Account *Account
	Error   error
}

// Service is in charge of delegating of account-related jobs to
// its workers.
type Service struct {
	config     Config
	Repository Repository
	jobQueue   chan Job
}

// loadAccount is a type of job to load an account from the storage.
type loadAccount struct {
	email    Email
	password Password
	channel  chan<- LoadResult
}

// saveAccount is a type of job to save an account in the storage.
type saveAccount struct {
	email   Email
	account Account
}

// Job represents an account-related job.
type Job interface{}

// NewService constructs a new Service.
func NewService(config Config, repository Repository) *Service {
	service := &Service{
		config:     config,
		Repository: repository,
		jobQueue:   make(chan Job),
	}

	for i := 0; i < config.WorkerCount; i++ {
		go service.spawnWorker()
	}

	return service
}

// LoadAccount loads an Account.
func (service *Service) LoadAccount(email Email, password Password) <-chan LoadResult {
	result := make(chan LoadResult, 1)
	service.jobQueue <- loadAccount{email: email, password: password, channel: result}
	return result
}

// SaveAccount saves the given Account.
func (service *Service) SaveAccount(email Email, account Account) {
	service.jobQueue <- saveAccount{email: email, account: account}
}

// spawnWorker spawns an Account Job worker that processes jobs for the Service.
func (service *Service) spawnWorker() {
	for job := range service.jobQueue {
		switch j := job.(type) {
		case loadAccount:
			account, err := service.Repository.Get(j.email, j.password)
			if err != nil {
				j.channel <- LoadResult{Error: err}
				continue
			}

			j.channel <- LoadResult{Account: account}
			break

		case saveAccount:
			err := service.Repository.Put(j.email, j.account)
			if err != nil {
				service.config.Logger.Error(err)

				// TODO what to do with this account?
			}

			break

		default:
			service.config.Logger.Errorf("Unexpected job of type %v", reflect.TypeOf(j))
		}
	}
}

// TearDown terminates this Service and cleans up resources.
func (service *Service) TearDown() {
	close(service.jobQueue)
}
