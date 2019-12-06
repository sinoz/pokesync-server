package discord

import "go.uber.org/zap"

// Config holds configurations specifically for the discord Service.
type Config struct {
	// TODO
}

// TODO
type Service struct {
	config Config
	logger *zap.SugaredLogger
}

// NewService constructs a new instance of a Service.
func NewService(config Config, logger *zap.SugaredLogger) *Service {
	return &Service{config: config, logger: logger}
}

// Stop stops this Service and cleans up resources.
func (service *Service) Stop() {
	// TODO
}
