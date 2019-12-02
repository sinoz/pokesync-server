package game

import (
	"time"

	"gitlab.com/pokesync/game-service/internal/game-service/game/entity"
)

const (
	// CommandLimit is the amount of command messages the
	// InboundNetworkProcessor consumes every pulse.
	CommandLimit = 3
)

// InboundNetworkProcessor processes received messages for entities that
// have a Session associated with them.
type InboundNetworkProcessor struct {
	// TODO
}

// OutboundNetworkProcessor processes queued messages for entities that
// have a Session associated with them.
type OutboundNetworkProcessor struct {
	// TODO
}

// NewInboundNetworkSystem constructs a new instance of an entity.System with
// a InboundNetworkProcessor as its internal processor.
func NewInboundNetworkSystem() *entity.System {
	return entity.NewSystem(entity.NewIntervalPolicy(100*time.Millisecond), NewInboundNetworkProcessor())
}

// NewOutboundNetworkSystem constructs a new instance of an entity.System with
// a OutboundNetworkProcessor as its internal processor.
func NewOutboundNetworkSystem() *entity.System {
	return entity.NewSystem(entity.NewDefaultSystemPolicy(), NewOutboundNetworkProcessor())
}

// NewInboundNetworkProcessor constructs a new instance of a
// InboundNetworkProcessor.
func NewInboundNetworkProcessor() *InboundNetworkProcessor {
	return &InboundNetworkProcessor{}
}

// NewOutboundNetworkProcessor constructs a new instance of a
// OutboundNetworkProcessor.
func NewOutboundNetworkProcessor() *OutboundNetworkProcessor {
	return &OutboundNetworkProcessor{}
}

// AddedToWorld is called when the System of this Processor is added
// to the game World.
func (processor *InboundNetworkProcessor) AddedToWorld(world *entity.World) error {
	return nil
}

// RemovedFromWorld is called when the System of this Processor is removed
// from the game World.
func (processor *InboundNetworkProcessor) RemovedFromWorld(world *entity.World) error {
	return nil
}

// AddedToWorld is called when the System of this Processor is added
// to the game World.
func (processor *OutboundNetworkProcessor) AddedToWorld(world *entity.World) error {
	return nil
}

// RemovedFromWorld is called when the System of this Processor is removed
// from the game World.
func (processor *OutboundNetworkProcessor) RemovedFromWorld(world *entity.World) error {
	return nil
}

// Update is called every game pulse to check if entities need any received
// messages processed.
func (processor *InboundNetworkProcessor) Update(world *entity.World, deltaTime time.Duration) error {
	entities := world.GetEntitiesFor(processor)
	for _, entity := range entities {
		sessionComponent := entity.GetComponent(SessionTag).(*SessionComponent)
		session := sessionComponent.Session

		for i := 0; i < CommandLimit; i++ {
			command := session.DequeueCommand()
			if command == nil {
				break
			}

			// TODO switch over command type
		}
	}

	return nil
}

// Update is called every game pulse to check if entities need any queued
// messages processed.
func (processor *OutboundNetworkProcessor) Update(world *entity.World, deltaTime time.Duration) error {
	entities := world.GetEntitiesFor(processor)
	for _, entity := range entities {
		sessionComponent := entity.GetComponent(SessionTag).(*SessionComponent)
		session := sessionComponent.Session

		var eventCount = 0
		for {
			event := session.DequeueEvent()
			if event == nil {
				break
			}

			session.Send(event)
			eventCount++
		}

		if eventCount > 0 {
			session.Flush()
		}
	}

	return nil
}

// Components returns a pack of ComponentTag's the InboundNetworkProcessor has
// interest in.
func (processor *InboundNetworkProcessor) Components() entity.ComponentTag {
	return SessionTag
}

// Components returns a pack of ComponentTag's the OutboundNetworkProcessor has
// interest in.
func (processor *OutboundNetworkProcessor) Components() entity.ComponentTag {
	return SessionTag
}
