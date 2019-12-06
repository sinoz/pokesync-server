package status

import (
	"gitlab.com/pokesync/game-service/internal/game-service/game"
)

// Parameters holds information about this game server's notifiable status.
type Parameters struct {
	PlayerCount int
}

// Provider provides the list of status parameters to notify external
// services of.
type Provider interface {
	Provide() (Parameters, error)
}

// ProviderImpl is an implementation of a Provider.
type ProviderImpl struct {
	gameService *game.Service
}

// NewProvider constructs a new status Provider.
func NewProvider(gameService *game.Service) *ProviderImpl {
	return &ProviderImpl{gameService}
}

// Provide provides the list of status Parameters from other internal services.
func (provider *ProviderImpl) Provide() (Parameters, error) {
	parameters := &Parameters{}

	return *parameters, nil
}
