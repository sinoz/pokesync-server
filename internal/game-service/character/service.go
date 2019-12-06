package character

import (
	"reflect"

	"gitlab.com/pokesync/game-service/internal/game-service/account"
	"go.uber.org/zap"
)

// Config holds configurations specific to the character Service.
type Config struct {
	WorkerCount int
}

// LoadResult is the result from attempting to load a character profile.
type LoadResult struct {
	Profile *Profile
	Error   error
}

// Service is in charge of delegating of character-related jobs to
// its workers.
type Service struct {
	config     Config
	logger     *zap.SugaredLogger
	Repository Repository
	jobQueue   chan Job
}

// loadProfile is a type of job to load a character profile from the storage.
type loadProfile struct {
	email   account.Email
	channel chan<- LoadResult
}

// saveProfile is a type of job to save an account in the storage.
type saveProfile struct {
	email   account.Email
	profile *Profile
}

// Job represents an account-related job.
type Job interface{}

// NewService constructs a new Service.
func NewService(config Config, logger *zap.SugaredLogger, repository Repository) *Service {
	service := &Service{
		config:     config,
		logger:     logger,
		Repository: repository,
		jobQueue:   make(chan Job),
	}

	for i := 0; i < config.WorkerCount; i++ {
		go service.worker()
	}

	return service
}

// LoadProfile loads an Account.
func (service *Service) LoadProfile(email account.Email) <-chan LoadResult {
	result := make(chan LoadResult, 1)
	service.jobQueue <- loadProfile{email: email, channel: result}
	return result
}

// SaveProfile saves the given Profile.
func (service *Service) SaveProfile(email account.Email, profile *Profile) {
	service.jobQueue <- saveProfile{email: email, profile: profile}
}

// worker continuously reads from the service's job queue until the
// queue is closed.
func (service *Service) worker() {
	for job := range service.jobQueue {
		switch j := job.(type) {
		case loadProfile:
			profile, err := service.Repository.Get(j.email)
			if err != nil {
				j.channel <- LoadResult{Error: err}
				continue
			}

			j.channel <- LoadResult{Profile: profile}
			break

		case saveProfile:
			err := service.Repository.Put(j.email, j.profile)
			if err != nil {
				service.logger.Error(err)

				// TODO what to do with this account?
			}

			break

		default:
			service.logger.Errorf("Unexpected job of type %v", reflect.TypeOf(j))
		}
	}
}

// Stop stops this Service and cleans up resources.
func (service *Service) Stop() {
	close(service.jobQueue)
}
