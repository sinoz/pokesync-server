package ecs

import (
	"time"
)

// World is the world of entities.
type World struct {
	entityManager *entityManager
	systemManager *systemManager
}

// entityListener listens for entity instances and component instances
// being added to and removed from the World.
type entityListener struct {
	world *World
}

// systemComponentListener listens for system instances being added to
// and removed from the world.
type systemComponentListener struct {
	world *World
}

// NewWorld constructs a new instance of a World.
func NewWorld(capacity int) *World {
	world := &World{
		entityManager: newEntityManager(capacity),
		systemManager: newSystemManager(),
	}

	world.entityManager.install(&entityListener{world: world})
	world.systemManager.install(&systemComponentListener{world: world})

	return world
}

// CreateEntity schedules the given Entity to be added to the world.
func (world *World) CreateEntity() *EntityBuilder {
	return newEntityBuilder(world)
}

// DestroyEntity schedules the given Entity to be removed from the world.
// Note that this does not occur until World#Update(dt) is called!
func (world *World) DestroyEntity(entity *Entity) {
	world.entityManager.remove(entity)
}

// AddSystem schedules the given System to be added to the world. Note that
//// this does not occur until World#Update(dt) is called!
func (world *World) AddSystem(system *System) {
	world.systemManager.add(system)
}

// RemoveSystem schedules the given System to be removed from the world. Note
// that this does not occur until World#Update(dt) is called!
func (world *World) RemoveSystem(system *System) {
	world.systemManager.remove(system)
}

// Update processes all of the entities and their components through the systems
// that they are indirectly subscribed to. Entities that are to be removed, are
// removed and entities that are to be added, are added.
func (world *World) Update(deltaTime time.Duration) error {
	if err := world.entityManager.update(deltaTime); err != nil {
		return err
	}

	if err := world.systemManager.update(world, deltaTime); err != nil {
		return err
	}

	return nil
}

// subscribeEntityToSystems subscribes the given Entity to each System within
// the World.
func (world *World) subscribeEntityToSystems(entity *Entity) {
	for _, system := range world.systemManager.systems {
		// make sure to only subscribe an entity to a system if the entity
		// isn't subscribed yet, to avoid duplicates.
		if entity.shouldBeSubscribedTo(system) && !entity.isSubscribedToSystem(system) {
			system.entities = append(system.entities, entity)
		}
	}
}

// unsubscribeEntityFromSystems unsubscribes the given Entity from each System
// within the World.
func (world *World) unsubscribeEntityFromSystems(entity *Entity) {
SystemLoop:
	for _, system := range world.systemManager.systems {
		if entity.shouldBeSubscribedTo(system) {
			for slot, ent := range system.entities {
				if ent.id == entity.id {
					system.entities = append(system.entities[:slot], system.entities[slot+1:]...)
					continue SystemLoop
				}
			}
		}
	}
}

// GetEntitiesFor looks up the list of entities that are indirectly subscribed
// to the System the given Processor is matched to.
func (world *World) GetEntitiesFor(processor Processor) []*Entity {
	for _, system := range world.systemManager.systems {
		if system.Processor == processor {
			return system.entities
		}
	}

	return emptyEntityList
}

func (listener *systemComponentListener) systemAdded(system *System) {
	for _, entity := range listener.world.entityManager.entities {
		if entity != nil && entity.shouldBeSubscribedTo(system) {
			system.entities = append(system.entities, entity)
		}
	}
}

func (listener *systemComponentListener) systemRemoved(system *System) {
	// no need to clear out the System's list of entities. leave it to the gc
}

func (listener *entityListener) entityAdded(entity *Entity) {
	listener.world.subscribeEntityToSystems(entity)
}

func (listener *entityListener) entityRemoved(entity *Entity) {
	listener.world.unsubscribeEntityFromSystems(entity)
}

func (listener *entityListener) componentAdded(entity *Entity, component Component) {
	listener.world.subscribeEntityToSystems(entity)
}

func (listener *entityListener) componentRemoved(entity *Entity, component Component) {
	listener.world.unsubscribeEntityFromSystems(entity)
}
