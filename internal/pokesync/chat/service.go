package chat

import (
	"gitlab.com/pokesync/game-service/internal/pokesync/client"
	"go.uber.org/zap"
)

// config holds configurations specifically for the chat service.
type Config struct {
	Logger *zap.SugaredLogger
}

// Service is an implementation of a public chat service and provides
// chatting capabilities for users across different channels.
type Service struct {
	config  Config
	routing *client.Router
}

// NewService constructs a new chat Service.
func NewService(config Config, routing *client.Router) *Service {
	return &Service{
		config:  config,
		routing: routing,
	}
}
