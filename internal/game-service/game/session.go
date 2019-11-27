package game

import (
	"sync"

	ecs "gitlab.com/pokesync/ecs/src"
	"gitlab.com/pokesync/game-service/internal/game-service/account"
	"gitlab.com/pokesync/game-service/internal/game-service/client"
)

// SessionConfig holds configurations specific to Session's.
type SessionConfig struct {
	CommandLimit int
	EventLimit   int
}

// Session acts as an interface between the client networking layer
// and the game simulation layer.
type Session struct {
	client *client.Client

	config SessionConfig

	Account account.Account
	Entity  *ecs.Entity

	commands chan client.Message
	events   chan client.Message
}

// SessionRegistry keeps track of Session's.
type SessionRegistry struct {
	sessions map[client.ID]*Session
	mutex    *sync.RWMutex
}

// NewSession constructs a new instance of a Session.
func NewSession(cl *client.Client, config SessionConfig, account account.Account, entity *ecs.Entity) *Session {
	return &Session{
		client: cl,

		config: config,

		Account: account,
		Entity:  entity,

		commands: make(chan client.Message, config.CommandLimit),
		events:   make(chan client.Message, config.EventLimit),
	}
}

// NewSessionRegistry constructs a new instance of a SessionRegistry.
func NewSessionRegistry() *SessionRegistry {
	return &SessionRegistry{
		sessions: make(map[client.ID]*Session),
		mutex:    &sync.RWMutex{},
	}
}

// DequeueCommand polls a command from the Session's buffer. May return nil
// if the queue is empty.
func (session *Session) DequeueCommand() client.Message {
	select {
	case command := <-session.commands:
		return command

	default:
		return nil
	}
}

// QueueCommand stores the given command into a queue of commands. The command
// may be dropped if this Session's command buffer has reached its capacity.
func (session *Session) QueueCommand(command client.Message) {
	select {
	case session.commands <- command:
		// successfully stored into the queue

	default:
		// message wasn't stored the buffer is full.
	}
}

// DequeueEvent polls an event from the Session's buffer. May return nil
// if the queue is empty.
func (session *Session) DequeueEvent() client.Message {
	select {
	case event := <-session.events:
		return event

	default:
		return nil
	}
}

// QueueEvent stores the given event into a queue of events. The event may be
// dropped if this Session's event buffer has reached its capacity.
func (session *Session) QueueEvent(event client.Message) {
	select {
	case session.events <- event:
		// successfully stored into the queue

	default:
		// message wasn't stored the buffer is full.
	}
}

// Send sends the given client Message directly to the underlying Client
// instance without any queueing.
func (session *Session) Send(message client.Message) {
	session.client.Send(message)
}

// Flush performs a flush call to write-and flush all of the queued
// up events to the socket connection.
func (session *Session) Flush() {
	session.client.Flush()
}

// Terminate terminates this Session and the underlying Client instance.
func (session *Session) Terminate() {
	session.client.Terminate()
}

// Put puts the given Session into the registry under the specified client ID.
func (registry *SessionRegistry) Put(id client.ID, session *Session) {
	registry.mutex.Lock()
	defer registry.mutex.Unlock()

	registry.sessions[id] = session
}

// Remove removes any Session that is associated with the specified client ID,
// from this registry.
func (registry *SessionRegistry) Remove(id client.ID) {
	registry.mutex.Lock()
	defer registry.mutex.Unlock()

	delete(registry.sessions, id)
}

// Get looks up a Session instance by the given client ID. May return nil.
func (registry *SessionRegistry) Get(id client.ID) *Session {
	registry.mutex.RLock()
	defer registry.mutex.RUnlock()

	return registry.sessions[id]
}
