package ecs

// ComponentTag is a unique bit mask value that is assigned to
// each type of Component for identification purposes.
type ComponentTag int

// Component represents a data typePack of a specific domain that
// is to describe a part of an Entity.
type Component interface {
	Tag() ComponentTag
}

const (
	// TypeLimit is the amount of different types of components an
	// entity can have.
	TypeLimit = 64
)

// componentStorage stores every Component for every Entity.
type componentStorage struct {
	components map[EntityId]*[TypeLimit]Component
}

// newComponentStorage constructs a new storage of components
// for all of our entities in our little ECS world.
func newComponentStorage() *componentStorage {
	return &componentStorage{components: make(map[EntityId]*[TypeLimit]Component)}
}

func (storage *componentStorage) add(id EntityId, component Component) {
	if storage.components[id] == nil {
		storage.components[id] = &[TypeLimit]Component{}
	}

	row := *storage.components[id]
	row[int(component.Tag())] = component

	storage.components[id] = &row
}

func (storage *componentStorage) remove(id EntityId, tag ComponentTag) {
	if storage.components[id] == nil {
		return
	}

	row := *storage.components[id]
	if row[int(tag)] == nil {
		return
	}

	row[int(tag)] = nil
	storage.components[id] = &row
}

func (storage *componentStorage) getComponentList(id EntityId) *[TypeLimit]Component {
	return storage.components[id]
}

func (storage *componentStorage) getComponent(id EntityId, tag ComponentTag) Component {
	if storage.components[id] == nil {
		return nil
	}

	row := *storage.components[id]
	return row[int(tag)]
}

func (storage *componentStorage) delete(id EntityId) {
	delete(storage.components, id)
}
