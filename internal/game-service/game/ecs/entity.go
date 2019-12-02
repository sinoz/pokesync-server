package ecs

import (
	"time"
)

const (
	// TypeLimit is the amount of different types of components an
	// entity can have.
	TypeLimit = 64
)

// EntityID is the unique id of an entity, which we can use to identify
// Entity's with within our world.
type EntityID int

// listener listens for changes made to an entity's bag of components.
type listener interface {
	entityAdded(entity *Entity)
	entityRemoved(entity *Entity)

	componentAdded(entity *Entity, component Component)
	componentRemoved(entity *Entity, component Component)
}

// componentListener listens for changes made to an Entity's bag
// of components, to publish these changes to the entityManager.
type componentListener struct {
	manager *entityManager
}

// Entity represents an entity in the game world. An entity only holds its own
// id and a bitpack value to reference which components an entity has.
type Entity struct {
	id EntityID

	typePack   int
	components [TypeLimit]Component

	listeners []listener
}

// EntityBuilder builds up an Entity.
type EntityBuilder struct {
	world      *World
	components []Component
}

// entityAddition is the addition of an entity to the world.
type entityAddition struct {
	components []Component
}

// entityRemoval is the removal of an entity from the world.
type entityRemoval struct {
	entity *Entity
}

// componentAddition is the addition of a component to an entity.
type componentAddition struct {
	entity    *Entity
	component Component
}

// componentRemoval is the removal of a component from an entity.
type componentRemoval struct {
	entity    *Entity
	component Component
}

// entityManager keeps track of Entity's that exist within our world.
type entityManager struct {
	entities []*Entity

	entitiesToAdd    []entityAddition
	entitiesToRemove []entityRemoval

	entitiesWithNewComponents []componentAddition
	entitiesWithOldComponents []componentRemoval

	listeners []listener
}

// NewEntity constructs a new Entity without any components.
func NewEntity() *Entity {
	return &Entity{}
}

// newEntityBuilder creates a new instance of an EntityBuilder.
func newEntityBuilder(world *World) *EntityBuilder {
	return &EntityBuilder{world: world}
}

// newEntityManager creates a new manager of entities.
func newEntityManager(capacity int) *entityManager {
	return &entityManager{entities: make([]*Entity, capacity)}
}

// With includes the given series of Component's in the building process of an Entity.
func (bldr *EntityBuilder) With(components ...Component) *EntityBuilder {
	for _, component := range components {
		bldr.components = append(bldr.components, component)
	}

	return bldr
}

// Build schedules the building-and registration process of an Entity to
// be ran by the World.
func (bldr *EntityBuilder) Build() {
	bldr.world.entityManager.add(bldr.components)
}

// Add adds the given Component to this entity's typePack of components
// by adding the component's tag value to the entity's typePack.
func (entity *Entity) Add(component Component) {
	entity.typePack |= int(component.Tag())
	entity.components[component.Tag()] = component

	entity.notifyComponentAdded(component)
}

// Contains checks whether the Entity holds a Component with the
// specified tag.
func (entity *Entity) Contains(tag ComponentTag) bool {
	return (entity.typePack & int(tag)) != 0
}

// GetComponent looks up a Component by its specified tag. May return
// null if there isn't such a Component.
func (entity *Entity) GetComponent(tag ComponentTag) Component {
	return entity.components[tag]
}

// Remove removes the given Component by its tag value.
func (entity *Entity) Remove(component Component) {
	entity.typePack &= ^int(component.Tag())
	entity.components[component.Tag()] = nil
	entity.notifyComponentRemoved(component)
}

// Clear clears this entity of all of its components.
func (entity *Entity) Clear() {
	entity.typePack = 0
}

// notifyComponentAdded notifies all listener's of the given Component
// having been added to this Entity.
func (entity *Entity) notifyComponentAdded(component Component) {
	for _, listener := range entity.listeners {
		listener.componentAdded(entity, component)
	}
}

// notifyComponentRemoved notifies all listener's of the given Component
// having been removed from this Entity.
func (entity *Entity) notifyComponentRemoved(component Component) {
	for _, listener := range entity.listeners {
		listener.componentRemoved(entity, component)
	}
}

// install installs the given listener into this Entity.
func (entity *Entity) install(listener listener) {
	entity.listeners = append(entity.listeners, listener)
}

// uninstall uninstalls the given listener from this Entity.
func (entity *Entity) uninstall(listener listener) {
	for i, l := range entity.listeners {
		if l == listener {
			entity.listeners = append(entity.listeners[:i], entity.listeners[i+1:]...)
			break
		}
	}
}

// isSubscribedToSystem returns whether the given Entity is subscribed to
// the given System.
func (entity *Entity) isSubscribedToSystem(system *System) bool {
	for _, ent := range system.entities {
		if ent == entity {
			return true
		}
	}

	return false
}

// shouldBeSubscribedTo returns whether this Entity has any interest in
// being subscribed to the specified System.
func (entity *Entity) shouldBeSubscribedTo(system *System) bool {
	return (int(system.Processor.Components()) & entity.GetTypePack()) != 0
}

// clearListeners clears this Entity from all of its listener's.
func (entity *Entity) clearListeners() {
	entity.listeners = []listener{}
}

// GetTypePack returns the bitpack value of all of the component
// types the entity has.
func (entity *Entity) GetTypePack() int {
	return entity.typePack
}

// add schedules the given Entity to be added to this entity manager.
func (manager *entityManager) add(components []Component) {
	manager.entitiesToAdd = append(manager.entitiesToAdd, entityAddition{components: components})
}

// remove schedules the given Entity to be removed from this entity manager.
func (manager *entityManager) remove(entity *Entity) {
	manager.entitiesToRemove = append(manager.entitiesToRemove, entityRemoval{entity: entity})
}

// add schedules the given Entity to be added to this entity manager.
func (manager *entityManager) addComponent(entity *Entity, component Component) {
	manager.entitiesWithNewComponents = append(manager.entitiesWithNewComponents, componentAddition{
		entity:    entity,
		component: component,
	})
}

// remove schedules the given Entity to be removed from this entity manager.
func (manager *entityManager) removeComponent(entity *Entity, component Component) {
	manager.entitiesWithOldComponents = append(manager.entitiesWithOldComponents, componentRemoval{
		entity:    entity,
		component: component,
	})
}

// update updates all pending entities for removals or additions to the world.
func (manager *entityManager) update(deltaTime time.Duration) error {
	for _, removal := range manager.entitiesToRemove {
		manager.entities[removal.entity.id] = nil
		removal.entity.clearListeners()

		manager.notifyEntityRemoved(removal.entity)
		manager.entitiesToRemove = manager.entitiesToRemove[1:]
	}

	for _, addition := range manager.entitiesToAdd {
		entityID := manager.getAvailableID()
		if entityID == -1 {
			manager.entitiesToAdd = manager.entitiesToAdd[1:]
			continue
		}

		entity := NewEntity()

		entity.id = entityID
		entity.install(&componentListener{manager: manager})

		for _, component := range addition.components {
			entity.Add(component)
		}

		manager.entities[entityID] = entity
		manager.notifyEntityAdded(entity)

		manager.entitiesToAdd = manager.entitiesToAdd[1:]
	}

	for _, removal := range manager.entitiesWithOldComponents {
		manager.notifyComponentRemoved(removal.entity, removal.component)
		manager.entitiesWithOldComponents = manager.entitiesWithOldComponents[1:]
	}

	for _, addition := range manager.entitiesWithNewComponents {
		manager.notifyComponentAdded(addition.entity, addition.component)
		manager.entitiesWithNewComponents = manager.entitiesWithNewComponents[1:]
	}

	return nil
}

// getAvailableID searches for an available slot within the list of entities.
// Returns -1 if no slot is currently available, which indicates that the
// world has reached its capacity.
func (manager *entityManager) getAvailableID() EntityID {
	for i, entity := range manager.entities {
		if entity == nil {
			return EntityID(i)
		}
	}

	return -1
}

// notifyComponentAdded notifies all listener's of the given Component
// having been added to the given Entity.
func (manager *entityManager) notifyComponentAdded(entity *Entity, component Component) {
	for _, listener := range manager.listeners {
		listener.componentAdded(entity, component)
	}
}

// notifyComponentRemoved notifies all listener's of the given Component
// having been removed from the given Entity.
func (manager *entityManager) notifyComponentRemoved(entity *Entity, component Component) {
	for _, listener := range manager.listeners {
		listener.componentRemoved(entity, component)
	}
}

// notifyEntityAdded notifies all listener's of the given Entity
// having been added to the manager.
func (manager *entityManager) notifyEntityAdded(entity *Entity) {
	for _, listener := range manager.listeners {
		listener.entityAdded(entity)
	}
}

// entityRemoved notifies all listener's of the given Entity
// having been removed from the manager.
func (manager *entityManager) notifyEntityRemoved(entity *Entity) {
	for _, listener := range manager.listeners {
		listener.entityRemoved(entity)
	}
}

// install installs the given listener into this manager.
func (manager *entityManager) install(listener listener) {
	manager.listeners = append(manager.listeners, listener)
}

// uninstall uninstalls the given listener from this manager.
func (manager *entityManager) uninstall(listener listener) {
	for i, l := range manager.listeners {
		if l == listener {
			manager.listeners = append(manager.listeners[:i], manager.listeners[i+1:]...)
			break
		}
	}
}

// clearListeners clears this entityManager from all of its listener's.
func (manager *entityManager) clearListeners() {
	manager.listeners = []listener{}
}

func (listener *componentListener) entityAdded(entity *Entity) {
	// nothing
}

func (listener *componentListener) entityRemoved(entity *Entity) {
	// nothing
}

func (listener *componentListener) componentAdded(entity *Entity, component Component) {
	listener.manager.addComponent(entity, component)
}

func (listener *componentListener) componentRemoved(entity *Entity, component Component) {
	listener.manager.removeComponent(entity, component)
}
