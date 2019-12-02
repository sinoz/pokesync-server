package game

import (
	"time"

	"gitlab.com/pokesync/game-service/internal/game-service/game/ecs"
)

// MapViewProcessor processes map view changes.
type MapViewProcessor struct {
	// TODO
}

// NewMapViewSystem constructs a new instance of an ecs.System with
// a MapViewProcessor as its internal processor.
func NewMapViewSystem() *ecs.System {
	return ecs.NewSystem(ecs.NewDefaultSystemPolicy(), NewMapViewProcessor())
}

// NewMapViewProcessor constructs a new instance of a MapViewProcessor.
func NewMapViewProcessor() *MapViewProcessor {
	return &MapViewProcessor{}
}

// AddedToWorld is called when the System of this Processor is added
// to the game World.
func (processor *MapViewProcessor) AddedToWorld(world *ecs.World) error {
	return nil
}

// RemovedFromWorld is called when the System of this Processor is removed
// from the game World.
func (processor *MapViewProcessor) RemovedFromWorld(world *ecs.World) error {
	return nil
}

// Update is called every game pulse to check if entities need their map view
// refreshed and if so, refreshes them.
func (processor *MapViewProcessor) Update(world *ecs.World, deltaTime time.Duration) error {
	return nil
}

// Components returns a pack of ComponentTag's the MapViewProcessor has
// interest in.
func (processor *MapViewProcessor) Components() ecs.ComponentTag {
	return MapViewTag
}
