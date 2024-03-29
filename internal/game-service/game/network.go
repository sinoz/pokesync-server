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

type attachFollowerHandler func(plr *Player, slot int) error

type clearFollowerHandler func(plr *Player) error

type switchPartySlotsHandler func(plr *Player, slotFrom, slotTo int) error

type selectPlayerOptionHandler func(plr *Player, entityID entity.ID, slot int) error

type continueDialogueHandler func(plr *Player) error

type submitChatCommandHandler func(plr *Player, trigger string, arguments []string) error

type clickTeleportHandler func(plr *Player, mapX, mapZ, localX, localZ int) error

type faceDirectionHandler func(plr *Player, direction Direction) error

type changeMovementTypeHandler func(plr *Player, movementType MovementType) error

type moveAvatarHandler func(plr *Player, direction Direction) error

type interactWithEntityHandler func(plr *Player, entityID entity.ID) error

type commandHandlerOption func(processor *InboundNetworkProcessor)

// InboundNetworkProcessor processes received messages for entities that
// have a Session associated with them.
type InboundNetworkProcessor struct {
	Logger *zap.SugaredLogger

	handleAttachFollower     attachFollowerHandler
	handleClearFollower      clearFollowerHandler
	handleSwitchPartySlot    switchPartySlotsHandler
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

func withSwitchPartySlotHandler(handler switchPartySlotsHandler) commandHandlerOption {
	return func(processor *InboundNetworkProcessor) {
		processor.handleSwitchPartySlot = handler
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

			var err error
			switch cmd := command.(type) {
			case *transport.AttachFollower:
				err = processor.handleAttachFollower(session.Player, int(cmd.PartySlot))
			case *transport.ClearFollower:
				err = processor.handleClearFollower(session.Player)
			case *transport.SubmitChatCommand:
				err = processor.handleChatCommandSubmit(session.Player, cmd.Trigger, cmd.Arguments)
			case *transport.SwitchPartySlots:
				err = processor.handleSwitchPartySlot(session.Player, int(cmd.SlotFrom), int(cmd.SlotTo))
			case *transport.FaceDirection:
				err = processor.handleFaceDirection(session.Player, Direction(cmd.Direction))
			case *transport.ChangeMovementType:
				err = processor.handleMovementTypeChange(session.Player, MovementType(cmd.Type))
			case *transport.MoveAvatar:
				err = processor.handleAvatarMove(session.Player, Direction(cmd.Direction))
			case *transport.ClickTeleport:
				err = processor.handleClickTeleport(session.Player, int(cmd.MapX), int(cmd.MapZ), int(cmd.LocalX), int(cmd.LocalZ))
			case *transport.SelectPlayerOption:
				err = processor.handlePlayerOptionSelect(session.Player, entity.ID(cmd.PID), int(cmd.Option))
			case *transport.ContinueDialogue:
				err = processor.handleContinueDialogue(session.Player)
			case *transport.InteractWithEntity:
				err = processor.handleEntityInteraction(session.Player, entity.ID(cmd.PID))
			default:
				processor.Logger.Errorf("Unexpected session command of type %v", reflect.TypeOf(cmd))
			}

			if err != nil {
				processor.Logger.Errorf("Error whilst processing command of type %v: %v", reflect.TypeOf(command), err)
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
