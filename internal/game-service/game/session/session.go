package session

import (
	"sync"

	"gitlab.com/pokesync/game-service/internal/game-service/game/entity"

	"gitlab.com/pokesync/game-service/internal/game-service/account"
	"gitlab.com/pokesync/game-service/internal/game-service/client"
)

// Config holds configurations specific to Session's.
type Config struct {
	CommandLimit int
	EventLimit   int
}

// Session acts as an interface between the client networking layer
// and the game simulation layer.
type Session struct {
	client *client.Client

	config Config

	Account account.Account
	Entity  *entity.Entity

	commands chan client.Message
	events   chan client.Message
}

// Registry keeps track of Session's.
type Registry struct {
	sessions map[client.ID]*Session
	mutex    *sync.RWMutex
}

// NewSession constructs a new instance of a Session.
func NewSession(cl *client.Client, config Config, account account.Account, entity *entity.Entity) *Session {
	session := &Session{
		client: cl,

		config: config,

		Account: account,
		Entity:  entity,

		commands: make(chan client.Message, config.CommandLimit),
		events:   make(chan client.Message, config.EventLimit),
	}

	return session
}

// NewInstalledSession constructs a new Session that installs listeners
// into the given Entity's components.
func NewInstalledSession(cl *client.Client, config Config, account account.Account, entity *entity.Entity) *Session {
	session := NewSession(cl, config, account, entity)

	return session
}

// NewRegistry constructs a new instance of a Registry.
func NewRegistry() *Registry {
	return &Registry{
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
func (registry *Registry) Put(id client.ID, session *Session) {
	registry.mutex.Lock()
	defer registry.mutex.Unlock()

	registry.sessions[id] = session
}

// Remove removes any Session that is associated with the specified client ID,
// from this registry.
func (registry *Registry) Remove(id client.ID) {
	registry.mutex.Lock()
	defer registry.mutex.Unlock()

	delete(registry.sessions, id)
}

// Get looks up a Session instance by the given client ID. May return nil.
func (registry *Registry) Get(id client.ID) *Session {
	registry.mutex.RLock()
	defer registry.mutex.RUnlock()

	return registry.sessions[id]
}
