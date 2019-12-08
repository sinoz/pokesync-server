package game

import (
	"reflect"
	"time"

	"go.uber.org/zap"

	"gitlab.com/pokesync/game-service/internal/game-service/game/entity"
	"gitlab.com/pokesync/game-service/internal/game-service/game/transport"
)

const (
	// CommandLimit is the amount of command messages the
	// InboundNetworkProcessor consumes every pulse.
	CommandLimit = 3
)

type attachFollowerHandler func(entity *entity.Entity, slot int) error

type clearFollowerHandler func(entity *entity.Entity) error

type selectPlayerOptionHandler func(entity *entity.Entity, entityID entity.ID, slot int) error

type continueDialogueHandler func(entity *entity.Entity) error

type submitChatCommandHandler func(entity *entity.Entity, trigger string, arguments []string) error

type clickTeleportHandler func(entity *entity.Entity, mapX, mapZ, localX, localZ int) error

type faceDirectionHandler func(entity *entity.Entity, direction Direction) error

type changeMovementTypeHandler func(entity *entity.Entity, movementType MovementType) error

type moveAvatarHandler func(entity *entity.Entity, direction Direction) error

type interactWithEntityHandler func(entity *entity.Entity, entityID entity.ID) error

type commandHandlerOption func(processor *InboundNetworkProcessor)

// InboundNetworkProcessor processes received messages for entities that
// have a Session associated with them.
type InboundNetworkProcessor struct {
	Logger *zap.SugaredLogger

	handleAttachFollower     attachFollowerHandler
	handleClearFollower      clearFollowerHandler
	handleContinueDialogue   continueDialogueHandler
	handleChatCommandSubmit  submitChatCommandHandler
	handleClickTeleport      clickTeleportHandler
	handlePlayerOptionSelect selectPlayerOptionHandler
	handleFaceDirection      faceDirectionHandler
	handleMovementTypeChange changeMovementTypeHandler
	handleAvatarMove         moveAvatarHandler
	handleEntityInteraction  interactWithEntityHandler
}

// OutboundNetworkProcessor processes queued messages for entities that
// have a Session associated with them.
type OutboundNetworkProcessor struct {
	// TODO
}

func withAttachFollowerHandler(handler attachFollowerHandler) commandHandlerOption {
	return func(processor *InboundNetworkProcessor) {
		processor.handleAttachFollower = handler
	}
}

func withClearFollowerHandler(handler clearFollowerHandler) commandHandlerOption {
	return func(processor *InboundNetworkProcessor) {
		processor.handleClearFollower = handler
	}
}

func withContinueDialogueHandler(handler continueDialogueHandler) commandHandlerOption {
	return func(processor *InboundNetworkProcessor) {
		processor.handleContinueDialogue = handler
	}
}

func withSubmitChatCommandHandler(handler submitChatCommandHandler) commandHandlerOption {
	return func(processor *InboundNetworkProcessor) {
		processor.handleChatCommandSubmit = handler
	}
}

func withClickTeleportHandler(handler clickTeleportHandler) commandHandlerOption {
	return func(processor *InboundNetworkProcessor) {
		processor.handleClickTeleport = handler
	}
}

func withSelectPlayerOptionHandler(handler selectPlayerOptionHandler) commandHandlerOption {
	return func(processor *InboundNetworkProcessor) {
		processor.handlePlayerOptionSelect = handler
	}
}

func withDirectionFacingHandler(handler faceDirectionHandler) commandHandlerOption {
	return func(processor *InboundNetworkProcessor) {
		processor.handleFaceDirection = handler
	}
}

func withMoveAvatarHandler(handler moveAvatarHandler) commandHandlerOption {
	return func(processor *InboundNetworkProcessor) {
		processor.handleAvatarMove = handler
	}
}

func withMovementTypeChangeHandler(handler changeMovementTypeHandler) commandHandlerOption {
	return func(processor *InboundNetworkProcessor) {
		processor.handleMovementTypeChange = handler
	}
}

func withEntityInteraction(handler interactWithEntityHandler) commandHandlerOption {
	return func(processor *InboundNetworkProcessor) {
		processor.handleEntityInteraction = handler
	}
}

// NewInboundNetworkSystem constructs a new instance of an entity.System with
// a InboundNetworkProcessor as its internal processor.
func NewInboundNetworkSystem(logger *zap.SugaredLogger, handlerOptions ...commandHandlerOption) *entity.System {
	return entity.NewSystem(entity.NewIntervalPolicy(100*time.Millisecond), NewInboundNetworkProcessor(logger, handlerOptions...))
}

// NewOutboundNetworkSystem constructs a new instance of an entity.System with
// a OutboundNetworkProcessor as its internal processor.
func NewOutboundNetworkSystem() *entity.System {
	return entity.NewSystem(entity.NewDefaultSystemPolicy(), NewOutboundNetworkProcessor())
}

// NewInboundNetworkProcessor constructs a new instance of a
// InboundNetworkProcessor.
func NewInboundNetworkProcessor(logger *zap.SugaredLogger, handlerOptions ...commandHandlerOption) *InboundNetworkProcessor {
	processor := &InboundNetworkProcessor{
		Logger: logger,
	}
	for _, applyOption := range handlerOptions {
		applyOption(processor)
	}

	return processor
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
	for _, ent := range entities {
		sessionComponent := ent.GetComponent(SessionTag).(*SessionComponent)
		session := sessionComponent.session

		for i := 0; i < CommandLimit; i++ {
			command := session.DequeueCommand()
			if command == nil {
				break
			}

			switch cmd := command.(type) {
			case *transport.AttachFollower:
				processor.handleAttachFollower(ent, int(cmd.PartySlot))
			case *transport.ClearFollower:
				processor.handleClearFollower(ent)
			case *transport.SubmitChatCommand:
				processor.handleChatCommandSubmit(ent, cmd.Trigger, cmd.Arguments)
			case *transport.FaceDirection:
				processor.handleFaceDirection(ent, Direction(cmd.Direction))
			case *transport.ChangeMovementType:
				processor.handleMovementTypeChange(ent, MovementType(cmd.Type))
			case *transport.MoveAvatar:
				processor.handleAvatarMove(ent, Direction(cmd.Direction))
			case *transport.ClickTeleport:
				processor.handleClickTeleport(ent, int(cmd.MapX), int(cmd.MapZ), int(cmd.LocalX), int(cmd.LocalZ))
			case *transport.SelectPlayerOption:
				processor.handlePlayerOptionSelect(ent, entity.ID(cmd.PID), int(cmd.Option))
			case *transport.ContinueDialogue:
				processor.handleContinueDialogue(ent)
			case *transport.InteractWithEntity:
				processor.handleEntityInteraction(ent, entity.ID(cmd.PID))
			default:
				processor.Logger.Errorf("Unexpected session command type of %v", reflect.TypeOf(cmd))
			}
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
		session := sessionComponent.session

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
