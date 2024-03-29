package status

import (
	"time"

	"go.uber.org/zap"
)

// Config holds configurations specific to this Service.
type Config struct {
	RefreshRate time.Duration
}

// Service is in charge of notifying external services of the status
// of this game server application.
type Service struct {
	config Config

	logger *zap.SugaredLogger

	notifier Notifier
	provider Provider

	quit chan bool
}

// NewService constructs a new instance of a Service.
func NewService(config Config, logger *zap.SugaredLogger, notifier Notifier, provider Provider) *Service {
	service := &Service{
		config: config,
		logger: logger,

		notifier: notifier,
		provider: provider,

		quit: make(chan bool, 1),
	}

	go service.worker()

	return service
}

// worker creates a new worker that continuously notifies external services
// of this game server's status at a preconfigured interval rate.
func (service *Service) worker() {
	for {
		select {
		case <-service.quit:
			return

		case <-time.After(service.config.RefreshRate):
			if err := service.notify(); err != nil {
				service.logger.Error(err.Error())
			}

			continue
		}
	}
}

// notify notifies external services of this game server's status.
func (service *Service) notify() error {
	status, err := service.provider.Provide()
	if err != nil {
		return err
	}

	return service.notifier.Notify(status)
}

// Stop stop this Service, no longer notifying external services.
// It also cleans up resources used by this Service.
func (service *Service) Stop() {
	service.quit <- true
}
