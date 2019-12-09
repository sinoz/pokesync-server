package game

import (
	"fmt"
	"time"

	"gitlab.com/pokesync/game-service/internal/game-service/game/entity"
)

const (
	Walk     MovementType = 0
	Run      MovementType = 1
	Cycle    MovementType = 2
	Teleport MovementType = 3
	Jump     MovementType = 4
	Surf     MovementType = 5
	Dive     MovementType = 6
	Glide    MovementType = 7

	walkingVelocity = 250 * time.Millisecond
	runningVelocity = (1 * time.Second) / 6
	cyclingVelocity = 100 * time.Millisecond
)

// MovementType is a type of movement an entity can perform.
type MovementType int

// Movement is a movement between two points in the game world.
type Movement struct {
	Source      Position
	Destination Position
	Type        MovementType
}

// MovementQueue is the queue of movement-related steps.
type MovementQueue struct {
	Position Position

	Facing       Direction
	MovementType MovementType

	targetPoint Position

	directionsToFace []Direction
	stepsToTake      []Direction
}

// WalkingProcessor processes walking steps.
type WalkingProcessor struct {
	Grid *Grid
}

// RunningProcessor processes running steps.
type RunningProcessor struct {
	Grid *Grid
}

// CyclingProcessor processes cycling steps.
type CyclingProcessor struct {
	Grid *Grid
}

// NewMovementQueue constructs a new instance of a MovementQueue.
func NewMovementQueue(position Position) *MovementQueue {
	return &MovementQueue{
		Position:     position,
		Facing:       South,
		MovementType: Walk,
	}
}

// NewWalkingProcessor TODO
func NewWalkingProcessor(grid *Grid) *WalkingProcessor {
	return &WalkingProcessor{Grid: grid}
}

// NewRunningProcessor TODO
func NewRunningProcessor(grid *Grid) *RunningProcessor {
	return &RunningProcessor{Grid: grid}
}

// NewCyclingProcessor TODO
func NewCyclingProcessor(grid *Grid) *CyclingProcessor {
	return &CyclingProcessor{Grid: grid}
}

// NewWalkingSystem constructs a System that processes walking
// steps for entities.
func NewWalkingSystem(grid *Grid) *entity.System {
	return entity.NewSystem(entity.NewIntervalPolicy(walkingVelocity), NewWalkingProcessor(grid))
}

// NewRunningSystem constructs a System that processes running
// steps for entities.
func NewRunningSystem(grid *Grid) *entity.System {
	return entity.NewSystem(entity.NewIntervalPolicy(runningVelocity), NewRunningProcessor(grid))
}

// NewCyclingSystem constructs a System that processes cycling
// steps for entities.
func NewCyclingSystem(grid *Grid) *entity.System {
	return entity.NewSystem(entity.NewIntervalPolicy(cyclingVelocity), NewCyclingProcessor(grid))
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
	entities := world.GetEntitiesFor(processor)
	for _, ent := range entities {
		transform := ent.GetComponent(TransformTag).(*TransformComponent)
		if transform.MovementQueue.MovementType != Walk {
			continue
		}

		if err := takeMovementSimulationStep(ent, transform.MovementQueue, processor.Grid); err != nil {
			return err
		}
	}

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
	entities := world.GetEntitiesFor(processor)
	for _, ent := range entities {
		transform := ent.GetComponent(TransformTag).(*TransformComponent)
		if transform.MovementQueue.MovementType != Run {
			continue
		}

		if err := takeMovementSimulationStep(ent, transform.MovementQueue, processor.Grid); err != nil {
			return err
		}
	}

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
	entities := world.GetEntitiesFor(processor)
	for _, ent := range entities {
		transform := ent.GetComponent(TransformTag).(*TransformComponent)
		if transform.MovementQueue.MovementType != Cycle {
			continue
		}

		if err := takeMovementSimulationStep(ent, transform.MovementQueue, processor.Grid); err != nil {
			return err
		}
	}

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

// takeMovementSimulationStep takes a single movement-based step within
// the simulation of the game, for the given Entity on the given map Grid.
func takeMovementSimulationStep(ent *entity.Entity, movementQueue *MovementQueue, grid *Grid) error {
	direction := movementQueue.PollDirectionToFace()
	if direction != nil {
		fmt.Println(*direction)
	}

	nextStep := movementQueue.PollStep()
	if nextStep != nil {
		oldPos := movementQueue.Position
		newPos, err := AddStep(oldPos, *nextStep, grid)
		if err != nil {
			return err
		}

		if oldPos.MapX != newPos.MapX || oldPos.MapZ != newPos.MapZ {
			if ent.Contains(MapViewTag) {
				mapView := ent.GetComponent(MapViewTag).(*MapViewComponent).MapView
				mapView.Refresh(newPos.MapX, newPos.MapZ)
			}
		}

		movementQueue.Position = newPos
		fmt.Println(newPos)

		// TODO add to tracking
	}

	return nil
}

// MoveTo sets the given point on the map as the target for the Entity
// to walk towards. The route to reach the target destination is progressively
// generated on every movement tick.
func (queue *MovementQueue) MoveTo(mapX, mapZ, localX, localZ int) {
	queue.targetPoint = Position{
		MapX:   mapX,
		MapZ:   mapZ,
		LocalX: localX,
		LocalZ: localZ,
	}
}

// AddStep adds the given Direction as the next step to take.
func (queue *MovementQueue) AddStep(direction Direction) {
	queue.stepsToTake = append(queue.stepsToTake, direction)
}

// AddDirectionToFace adds the given Direction as the next direction to face.
func (queue *MovementQueue) AddDirectionToFace(direction Direction) {
	queue.directionsToFace = append(queue.directionsToFace, direction)
}

// PollDirectionToFace polls the next direction to face from the queue.
// May return nil if the queue is empty.
func (queue *MovementQueue) PollDirectionToFace() *Direction {
	if len(queue.directionsToFace) == 0 {
		return nil
	}

	direction := queue.directionsToFace[0]
	queue.directionsToFace = queue.directionsToFace[1:]
	return &direction
}

// PollStep polls the next step to take from the queue. May return nil
// if the queue is empty.
func (queue *MovementQueue) PollStep() *Direction {
	if len(queue.stepsToTake) == 0 {
		return nil
	}

	step := queue.stepsToTake[0]
	queue.stepsToTake = queue.stepsToTake[1:]
	return &step
}

func faceDirection() faceDirectionHandler {
	return func(plr *Player, direction Direction) error {
		plr.Face(direction)

		return nil
	}
}

func changeMovementType() changeMovementTypeHandler {
	return func(plr *Player, movementType MovementType) error {
		switch movementType {
		case Walk:
			plr.Walk()
		case Run:
			plr.Run()
		case Cycle:
			plr.Cycle()
		default:
			return fmt.Errorf("unexpected movement type of value %v", movementType)
		}

		return nil
	}
}

func moveAvatar() moveAvatarHandler {
	return func(plr *Player, direction Direction) error {
		plr.Move(direction)

		return nil
	}
}

func clickTeleport() clickTeleportHandler {
	return func(plr *Player, mapX, mapZ, localX, localZ int) error {
		return nil
	}
}
