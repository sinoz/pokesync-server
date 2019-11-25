package character

import "gitlab.com/pokesync/game-service/internal/pokesync/client"

type Config struct {
	WorkerCount int
}

type Service struct {
	config  Config
	routing *client.Router
}

// NewService constructs a new instance of a Service.
func NewService(config Config, routing *client.Router) *Service {
	service := &Service{
		config:  config,
		routing: routing,
	}

	for i := 0; i < config.WorkerCount; i++ {
		service.spawnWorker()
	}

	return service
}

// TODO
func (service *Service) spawnWorker() {
	go func() {
		// TODO
	}()
}
