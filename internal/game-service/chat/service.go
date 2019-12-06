package chat

import (
	"gitlab.com/pokesync/game-service/internal/game-service/client"
	"go.uber.org/zap"
)

// Config holds configurations specifically for the chat service.
type Config struct {
}

// Service is an implementation of a public chat service and provides
// chatting capabilities for users across different channels.
type Service struct {
	config  Config
	logger  *zap.SugaredLogger
	routing *client.Router
}

// NewService constructs a new chat Service.
func NewService(config Config, logger *zap.SugaredLogger, routing *client.Router) *Service {
	return &Service{
		config:  config,
		logger:  logger,
		routing: routing,
	}
}

// Stop stops this Service and cleans up resources.
func (service *Service) Stop() {
	// TODO
}
