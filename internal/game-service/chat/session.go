package chat

import "gitlab.com/pokesync/game-service/internal/game-service/client"

// TODO
type SessionConfig struct {
	// TODO
}

// TODO
type Session struct {
	client *client.Client
	config SessionConfig
}

// NewSession constructs a new instance of a public chat Session.
func NewSession(cli *client.Client, config SessionConfig) *Session {
	return &Session{
		client: cli,
		config: config,
	}
}
