package discord

// Config holds configurations specifically for the discord Service.
type Config struct {
	// TODO
}

// TODO
type Service struct {
	config Config
}

// NewService constructs a new instance of a Service.
func NewService(config Config) *Service {
	return &Service{config: config}
}

// TearDown terminates this Service and cleans up resources.
func (service *Service) TearDown() {
	// TODO
}
