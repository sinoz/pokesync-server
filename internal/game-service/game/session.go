package game

import (
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

	Email  account.Email
	Player *Player

	commands chan client.Message
	events   chan client.Message
}

// SessionRegistry keeps track of Session's.
type SessionRegistry struct {
	sessions map[client.ID]*Session
}

// NewSession constructs a new instance of a Session.
func NewSession(cl *client.Client, config SessionConfig, email account.Email, player *Player) *Session {
	session := &Session{
		client: cl,

		config: config,

		Email:  email,
		Player: player,

		commands: make(chan client.Message, config.CommandLimit),
		events:   make(chan client.Message, config.EventLimit),
	}

	return session
}

// NewInstalledSession constructs a new Session that installs listeners
// into the given Player's components.
func NewInstalledSession(cl *client.Client, config SessionConfig, email account.Email, plr *Player) *Session {
	session := NewSession(cl, config, email, plr)

	plr.
		GetComponent(MapViewTag).(*MapViewComponent).MapView.
		AddListener(&MapViewSessionListener{session: session})

	plr.
		GetComponent(CoinBagTag).(*CoinBagComponent).CoinBag.
		AddListener(&CoinBagSessionListener{session: session})

	plr.
		GetComponent(PartyBeltTag).(*PartyBeltComponent).PartyBelt.
		AddListener(&PartyBeltSessionListener{session: session})

	return session
}

// NewSessionRegistry constructs a new instance of a SessionRegistry.
func NewSessionRegistry() *SessionRegistry {
	return &SessionRegistry{sessions: make(map[client.ID]*Session)}
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
func (session *Session) QueueCommand(command client.Message) bool {
	select {
	case session.commands <- command:
		return true // successfully stored into the queue

	default:
		return false // message wasn't stored the buffer is full.
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
func (session *Session) QueueEvent(event client.Message) bool {
	select {
	case session.events <- event:
		return true // successfully stored into the queue

	default:
		return false // message wasn't stored the buffer is full.
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
