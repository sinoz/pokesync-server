package entity

import (
	"errors"
)

// List holds a list of entities in a fixed-sized collection.
type List struct {
	entities     []*Entity
	availableIDs []ID

	lowest  int
	highest int

	size int
}

// NewList constructs a list of entities that is bounded
// to the specified capacity.
func NewList(capacity int) *List {
	if capacity < 0 {
		panic("negative capacity")
	}

	pids := make([]ID, capacity)
	for i := 0; i < capacity; i++ {
		pids[i] = ID(i)
	}

	return &List{
		entities: make([]*Entity, capacity+1), // + 1 due to the 'highest' variable

		lowest:  -1,
		highest: -1,

		availableIDs: pids,
	}
}

// Insert inserts the given Entity at the specified index.
func (list *List) Insert(id ID, entity *Entity) error {
	if err := list.checkBoundaries(id); err != nil {
		return err
	}

	if int(id) >= list.highest {
		list.highest = int(id)
	}

	if list.size == 0 {
		list.lowest = int(id)
	} else if int(id) < list.lowest {
		list.lowest = int(id)
	}

	list.entities[id] = entity
	list.size++

	return nil
}

// Clear clears out an Entity that is set at the specified index. Returns the
// entity that was removed and potentially an error.
func (list *List) Clear(id ID) (*Entity, error) {
	if err := list.checkBoundaries(id); err != nil {
		return nil, err
	}

	entity := list.entities[id]
	if entity != nil {
		list.size--
		list.entities[id] = nil

		for i := list.highest; i >= list.lowest; i-- {
			if list.entities[i] != nil {
				list.highest = i
				break
			}
		}

		for i := 0; i <= list.highest; i++ {
			if list.entities[i] != nil {
				list.lowest = i
				break
			}
		}

		if list.size == 0 {
			list.highest = 0
			list.lowest = 0
		}
	}

	return entity, nil
}

// Get looks up an Entity by its ID. May return an error.
func (list *List) Get(id ID) (*Entity, error) {
	if err := list.checkBoundaries(id); err != nil {
		return nil, err
	}

	return list.entities[id], nil
}

// IsEmpty checks if an Entity is occupying the given ID. May return an error.
func (list *List) IsEmpty(id ID) (bool, error) {
	if err := list.checkBoundaries(id); err != nil {
		return false, err
	}

	return list.entities[id] != nil, nil
}

// Find searches for an Entity that satisfies the given predicate.
func (list *List) Find(p func(*Entity) bool) *Entity {
	if list.lowest == -1 || list.highest >= list.size {
		return nil
	}

	for i := list.lowest; i <= list.highest; i++ {
		entity := list.entities[i]
		if entity != nil && p(entity) {
			return entity
		}
	}

	return nil
}

// GetAvailableID obtains an available ID for an Entity.
func (list *List) GetAvailableID() (ID, bool) {
	if len(list.availableIDs) == 0 {
		return ID(0), false
	}

	id := list.availableIDs[0]
	list.availableIDs = list.availableIDs[1:]

	return id, true
}

// ReleaseID releases the given ID back into the stack of ID's. The ID is
// prepended to the top of the stack to avoid a large gap in between entity
// ID's over time.
func (list *List) ReleaseID(id ID) {
	list.availableIDs = append([]ID{id}, list.availableIDs...)
}

// LowestID returns the lowest EntityID that is currently in use.
func (list *List) LowestID() (ID, bool) {
	if list.lowest == -1 {
		return ID(0), false
	}

	return ID(list.lowest), true
}

// HighestID returns the highest EntityID that is currently in use.
func (list *List) HighestID() (ID, bool) {
	if list.highest == -1 {
		return ID(0), false
	}

	return ID(list.highest), true
}

// Size returns the current amount of entities that are registered.
func (list *List) Size() int {
	return list.size
}

// checkBoundaries checks if the given ID is out of bounds or not.
func (list *List) checkBoundaries(id ID) error {
	if id < 0 || int(id) >= len(list.entities) {
		return errors.New("ID out of bounds")
	}

	return nil
}
