package session

import (
	"testing"
	"time"

	"gitlab.com/pokesync/game-service/internal/game-service/game/entity"

	"gitlab.com/pokesync/game-service/internal/game-service/account"
	"gitlab.com/pokesync/game-service/internal/game-service/client"
	"gitlab.com/pokesync/game-service/internal/game-service/game/transport"
)

func TestSession_QueueCommand(t *testing.T) {
	account := account.Account{}
	config := Config{CommandLimit: 16, EventLimit: 16}

	command := &transport.SubmitChatMessage{}

	session := NewSession(nil, config, account, entity.NewEntity())
	session.QueueCommand(command)

	select {
	case read := <-session.commands:
		if read != command {
			t.Error("invalid command")
		}

	case <-time.After(1 * time.Second):
		t.Error("timeout")
	}
}

func TestSession_QueueEvent(t *testing.T) {
	account := account.Account{}
	config := Config{CommandLimit: 16, EventLimit: 16}

	event := &transport.DisplayChatMessage{}

	session := NewSession(nil, config, account, entity.NewEntity())
	session.QueueEvent(event)

	select {
	case read := <-session.events:
		if read != event {
			t.Error("invalid event")
		}

	case <-time.After(1 * time.Second):
		t.Error("timeout")
	}
}

func TestSession_DequeueCommand(t *testing.T) {
	account := account.Account{}
	config := Config{CommandLimit: 16, EventLimit: 16}

	session := NewSession(nil, config, account, entity.NewEntity())
	commandCh := make(chan client.Message)

	go func() {
		commandCh <- session.DequeueCommand()
	}()

	select {
	case command := <-commandCh:
		if command != nil {
			t.Error("expected nil to come out")
		}

	case <-time.After(1 * time.Second):
		t.Error("timeout")
	}
}

func TestSession_DequeueEvent(t *testing.T) {
	account := account.Account{}
	config := Config{CommandLimit: 16, EventLimit: 16}

	session := NewSession(nil, config, account, entity.NewEntity())
	eventCh := make(chan client.Message)

	go func() {
		eventCh <- session.DequeueEvent()
	}()

	select {
	case evt := <-eventCh:
		if evt != nil {
			t.Error("expected nil to come out")
		}

	case <-time.After(1 * time.Second):
		t.Error("timeout")
	}
}
