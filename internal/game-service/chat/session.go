package chat

import (
	"fmt"

	"gitlab.com/pokesync/game-service/internal/game-service/client"
)

// SessionConfig TODO
type SessionConfig struct {
	BufferLimit int
}

// Session TODO
type Session struct {
	client  *client.Client
	config  SessionConfig
	mailbox chan client.Message
}

// SessionRegistry keeps track of chat Session's.
type SessionRegistry struct {
	sessions map[client.ID]*Session
}

// NewSession constructs a new instance of a chat Session.
func NewSession(cli *client.Client, config SessionConfig) *Session {
	session := &Session{
		client:  cli,
		config:  config,
		mailbox: make(chan client.Message, config.BufferLimit),
	}

	go session.receive()

	return session
}

// NewSessionRegistry constructs a new instance of a SessionRegistry.
func NewSessionRegistry() *SessionRegistry {
	return &SessionRegistry{sessions: make(map[client.ID]*Session)}
}

// receive receives client Message's and handles them.
func (session *Session) receive() {
	for message := range session.mailbox {
		switch msg := message.(type) {
		case *SubmitChatMessage:
			fmt.Println(msg)
		case *SelectChatChannel:
			fmt.Println(msg)
		}
	}
}

// Enqueue enqueues the given client Message into this Session's mailbox.
// Returns whether the message is consumed or to be dropped.
func (session *Session) Enqueue(message client.Message) bool {
	select {
	case session.mailbox <- message:
		return true

	default:
		return false
	}
}

// Stop stops the Session, closing its underlying mailbox of messages.
func (session *Session) Stop() {
	close(session.mailbox)
}

// Put puts the given Session into the registry under the specified client ID.
func (registry *SessionRegistry) Put(id client.ID, session *Session) {
	registry.sessions[id] = session
}

// Remove removes any Session that is associated with the specified client ID,
// from this registry.
func (registry *SessionRegistry) Remove(id client.ID) *Session {
	session, exists := registry.sessions[id]
	if !exists {
		return nil
	}

	delete(registry.sessions, id)
	return session
}

// Get looks up a Session instance by the given client ID. May return nil.
func (registry *SessionRegistry) Get(id client.ID) *Session {
	return registry.sessions[id]
}
