package game

import (
	"time"

	"gitlab.com/pokesync/game-service/internal/game-service/game/entity"
	"gitlab.com/pokesync/game-service/pkg/event"
)

// MapViewProcessor processes map view changes.
type MapViewProcessor struct {
	eventBus event.Bus
}

// NewMapViewSystem constructs a new instance of an entity.System with
// a MapViewProcessor as its internal processor.
func NewMapViewSystem(eventBus event.Bus) *entity.System {
	return entity.NewSystem(entity.NewDefaultSystemPolicy(), NewMapViewProcessor(eventBus))
}

// NewMapViewProcessor constructs a new instance of a MapViewProcessor.
func NewMapViewProcessor(eventBus event.Bus) *MapViewProcessor {
	return &MapViewProcessor{eventBus: eventBus}
}

// AddedToWorld is called when the System of this Processor is added
// to the game World.
func (processor *MapViewProcessor) AddedToWorld(world *entity.World) error {
	return nil
}

// RemovedFromWorld is called when the System of this Processor is removed
// from the game World.
func (processor *MapViewProcessor) RemovedFromWorld(world *entity.World) error {
	return nil
}

// Update is called every game pulse to check if entities need their map view
// refreshed and if so, refreshes them.
func (processor *MapViewProcessor) Update(world *entity.World, deltaTime time.Duration) error {
	return nil
}

// Components returns a pack of ComponentTag's the MapViewProcessor has
// interest in.
func (processor *MapViewProcessor) Components() entity.ComponentTag {
	return MapViewTag
}
