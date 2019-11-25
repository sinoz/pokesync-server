package ecs

import "github.com/google/uuid"

// EntityId is the unique id of an entity, which we can use to identify
// Entity's with within our world.
type EntityId uuid.UUID

// Entity represents an entity in the game world. An entity only holds its own
// Id and a bitpack value to reference which components an entity has.
type Entity struct {
	Id       EntityId
	typePack int
}

// entityManager keeps track of Entity's that exist within our engine.
type entityManager struct {
	entities map[EntityId]*Entity
}

// newEntityId generates a new unique EntityId.
func newEntityId() EntityId {
	return EntityId(uuid.New())
}

// NewEntity constructs a new Entity without any components.
func NewEntity() *Entity {
	return &Entity{}
}

// newEntityManager creates a new manager of entities.
func newEntityManager() *entityManager {
	return &entityManager{entities: make(map[EntityId]*Entity)}
}

// Add adds the given Component to this entity's typePack of components
// by adding the component's tag value to the entity's typePack.
func (entity *Entity) Add(component Component) {
	entity.typePack |= int(component.Tag())
}

// Contains checks whether the Entity holds a Component with the
// specified tag.
func (entity *Entity) Contains(tag ComponentTag) bool {
	return (entity.typePack & int(tag)) != 0
}

// Remove removes the given Component by its tag value.
func (entity *Entity) Remove(component Component) {
	entity.RemoveByTag(component.Tag())
}

// Remove removes a Component by the given tag value.
func (entity *Entity) RemoveByTag(tag ComponentTag) {
	entity.typePack &= ^int(tag)
}

// Clear clears this entity of all of its components.
func (entity *Entity) Clear() {
	entity.typePack = 0
}

// GetTypePack returns the bitpack value of all of the component
// types the entity has.
func (entity *Entity) GetTypePack() int {
	return entity.typePack
}
