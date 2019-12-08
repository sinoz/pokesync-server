package game

import (
	"time"

	"gitlab.com/pokesync/game-service/internal/game-service/game/entity"
)

// MovementType is a type of movement an entity can perform.
type MovementType int

// These are the supported types of movement.
const (
	Walk     MovementType = 0
	Run      MovementType = 1
	Cycle    MovementType = 2
	Teleport MovementType = 3
	Jump     MovementType = 4
	Surf     MovementType = 5
	Dive     MovementType = 6
	Glide    MovementType = 7
)

// Movement is a movement between two points in the game world.
type Movement struct {
	Source      Position
	Destination Position
	Type        MovementType
}

const (
	walkingVelocity = 250 * time.Millisecond
	runningVelocity = (1 * time.Second) / 6
	cyclingVelocity = 100 * time.Millisecond
)

// WalkingProcessor processes walking steps.
type WalkingProcessor struct{}

// RunningProcessor processes running steps.
type RunningProcessor struct{}

// CyclingProcessor processes cycling steps.
type CyclingProcessor struct{}

// NewWalkingProcessor TODO
func NewWalkingProcessor() *WalkingProcessor {
	return &WalkingProcessor{}
}

// NewRunningProcessor TODO
func NewRunningProcessor() *RunningProcessor {
	return &RunningProcessor{}
}

// NewCyclingProcessor TODO
func NewCyclingProcessor() *CyclingProcessor {
	return &CyclingProcessor{}
}

// NewWalkingSystem constructs a System that processes walking
// steps for entities.
func NewWalkingSystem() *entity.System {
	return entity.NewSystem(entity.NewIntervalPolicy(walkingVelocity), NewWalkingProcessor())
}

// NewRunningSystem constructs a System that processes running
// steps for entities.
func NewRunningSystem() *entity.System {
	return entity.NewSystem(entity.NewIntervalPolicy(runningVelocity), NewRunningProcessor())
}

// NewCyclingSystem constructs a System that processes cycling
// steps for entities.
func NewCyclingSystem() *entity.System {
	return entity.NewSystem(entity.NewIntervalPolicy(cyclingVelocity), NewCyclingProcessor())
}

// AddedToWorld is called when the System of this Processor is added
// to the game World.
func (processor *WalkingProcessor) AddedToWorld(world *entity.World) error {
	return nil
}

// RemovedFromWorld is called when the System of this Processor is removed
// from the game World.
func (processor *WalkingProcessor) RemovedFromWorld(world *entity.World) error {
	return nil
}

// Update is called every game pulse to check if entities need to take any
// walking steps and if so, applies them.
func (processor *WalkingProcessor) Update(world *entity.World, deltaTime time.Duration) error {
	return nil
}

// AddedToWorld is called when the System of this Processor is added
// to the game World.
func (processor *RunningProcessor) AddedToWorld(world *entity.World) error {
	return nil
}

// RemovedFromWorld is called when the System of this Processor is removed
// from the game World.
func (processor *RunningProcessor) RemovedFromWorld(world *entity.World) error {
	return nil
}

// Update is called every game pulse to check if entities need to take any
// running steps and if so, applies them.
func (processor *RunningProcessor) Update(world *entity.World, deltaTime time.Duration) error {
	return nil
}

// AddedToWorld is called when the System of this Processor is added
// to the game World.
func (processor *CyclingProcessor) AddedToWorld(world *entity.World) error {
	return nil
}

// RemovedFromWorld is called when the System of this Processor is removed
// from the game World.
func (processor *CyclingProcessor) RemovedFromWorld(world *entity.World) error {
	return nil
}

// Update is called every game pulse to check if entities need to take any
// cycling steps and if so, applies them.
func (processor *CyclingProcessor) Update(world *entity.World, deltaTime time.Duration) error {
	return nil
}

// Components returns a pack of ComponentTag's the WalkingProcessor has
// interest in.
func (processor *WalkingProcessor) Components() entity.ComponentTag {
	return TransformTag
}

// Components returns a pack of ComponentTag's the RunningProcessor has
// interest in.
func (processor *RunningProcessor) Components() entity.ComponentTag {
	return TransformTag | CanRunTag
}

// Components returns a pack of ComponentTag's the CyclingProcessor has
// interest in.
func (processor *CyclingProcessor) Components() entity.ComponentTag {
	return TransformTag
}

// faceDirection tells the given Entity to change its currently
// facing direction to the specified one.
func faceDirection() faceDirectionHandler {
	return func(entity *entity.Entity, direction Direction) error {
		return nil
	}
}

func changeMovementType() changeMovementTypeHandler {
	return func(entity *entity.Entity, movementType MovementType) error {
		return nil
	}
}

func moveAvatar() moveAvatarHandler {
	return func(entity *entity.Entity, direction Direction) error {
		return nil
	}
}

func clickTeleport() clickTeleportHandler {
	return func(entity *entity.Entity, mapX, mapZ, localX, localZ int) error {
		return nil
	}
}
