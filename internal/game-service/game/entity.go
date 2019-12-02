package game

import (
	"errors"

	"gitlab.com/pokesync/game-service/internal/game-service/character"
	"gitlab.com/pokesync/game-service/internal/game-service/game/ecs"
)

const (
	PlayerKind  EntityKind = 0
	NpcKind     EntityKind = 1
	MonsterKind EntityKind = 2
	ObjectKind  EntityKind = 3

	Man        Gender = 0
	Woman      Gender = 1
	Genderless Gender = 2
)

// PID is a process identifier that is used to identify entities with.
type PID int

// ModelID is the ID of an Entity's model.
type ModelID int

// EntityKind represents the type of an Entity.
type EntityKind int

// Gender is a type of gender of an Entity.
type Gender int

// EntityList holds a list of entities in a fixed-sized collection.
type EntityList struct {
	entities []*ecs.Entity

	availablePIDs []PID

	lowest  int
	highest int

	size int
}

// EntityFactory is in charge of producing different types of entities.
type EntityFactory struct {
	assets *AssetBundle
}

// NewEntityList constructs a list of entities that is bounded
// to the specified capacity.
func NewEntityList(capacity int) *EntityList {
	if capacity < 0 {
		panic("negative capacity")
	}

	pids := make([]PID, capacity)
	for i := 0; i < capacity; i++ {
		pids[i] = PID(i)
	}

	return &EntityList{
		entities: make([]*ecs.Entity, capacity+1), // + 1 due to the 'highest' variable

		lowest:  -1,
		highest: -1,

		availablePIDs: pids,
	}
}

// NewEntityFactory constructs a new EntityFactory to produce entities with.
func NewEntityFactory(assets *AssetBundle) *EntityFactory {
	return &EntityFactory{assets: assets}
}

// CreatePlayer creates the set of Component's to create a Player-like Entity from.
func (factory *EntityFactory) CreatePlayer(pid PID, position Position, direction Direction, gender Gender, displayName character.DisplayName, userGroup character.UserGroup) []ecs.Component {
	return []ecs.Component{
		&PIDComponent{PID: pid},
		&TransformComponent{Position: position},
		&KindComponent{Kind: PlayerKind},
	}
}

// CreateNpc creates the set of Component's to create a Npc-like Entity from.
func (factory *EntityFactory) CreateNpc(pid PID, position Position, direction Direction, modelID ModelID) []ecs.Component {
	return []ecs.Component{
		&PIDComponent{PID: pid},
		&ModelIDComponent{ModelID: modelID},
		&TransformComponent{Position: position},
		&KindComponent{Kind: NpcKind},
		&BlockingComponent{},
		&TrackingComponent{},
	}
}

// CreateMonster creates the set of Component's to create a Monster-like Entity from.
func (factory *EntityFactory) CreateMonster(pid PID, position Position, direction Direction, modelID ModelID) []ecs.Component {
	return []ecs.Component{
		&PIDComponent{PID: pid},
		&ModelIDComponent{ModelID: modelID},
		&TransformComponent{Position: position},
		&HealthComponent{Max: 1, Current: 1}, // TODO
		&KindComponent{Kind: MonsterKind},
	}
}

// CreateObject creates the set of Component's to create a Object-like Entity from.
func (factory *EntityFactory) CreateObject(pid PID, position Position, direction Direction) []ecs.Component {
	return []ecs.Component{}
}

// Insert inserts the given Entity at the specified index.
func (list *EntityList) Insert(pid PID, entity *ecs.Entity) error {
	if err := list.checkBoundaries(pid); err != nil {
		return err
	}

	if int(pid) >= list.highest {
		list.highest = int(pid)
	}

	if list.size == 0 {
		list.lowest = int(pid)
	} else if int(pid) < list.lowest {
		list.lowest = int(pid)
	}

	list.entities[pid] = entity
	list.size++

	return nil
}

// Clear clears out an Entity that is set at the specified index. Returns the
// entity that was removed and potentially an error.
func (list *EntityList) Clear(pid PID) (*ecs.Entity, error) {
	if err := list.checkBoundaries(pid); err != nil {
		return nil, err
	}

	entity := list.entities[pid]
	if entity != nil {
		list.size--
		list.entities[pid] = nil

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

		return entity, nil
	}

	return nil, nil
}

// Get looks up an Entity by its PID. May return an error.
func (list *EntityList) Get(pid PID) (*ecs.Entity, error) {
	if err := list.checkBoundaries(pid); err != nil {
		return nil, err
	}

	return list.entities[pid], nil
}

// IsEmpty checks if an Entity is occupying the given PID. May return an error.
func (list *EntityList) IsEmpty(pid PID) (bool, error) {
	if err := list.checkBoundaries(pid); err != nil {
		return false, err
	}

	return list.entities[pid] != nil, nil
}

// Find searches for an Entity that satisfies the given predicate.
func (list *EntityList) Find(p func(*ecs.Entity) bool) *ecs.Entity {
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

// GetAvailablePID obtains an available PID for an Entity.
func (list *EntityList) GetAvailablePID() (PID, bool) {
	if len(list.availablePIDs) == 0 {
		return PID(0), false
	}

	pid := list.availablePIDs[0]
	list.availablePIDs = list.availablePIDs[1:]
	return pid, true
}

// ReleasePID releases the given PID back into the queue of PID's.
func (list *EntityList) ReleasePID(pid PID) {
	list.availablePIDs = append(list.availablePIDs, pid)
}

// LowestPID returns the lowest PID that is currently in use.
func (list *EntityList) LowestPID() (PID, bool) {
	if list.lowest == -1 {
		return PID(0), false
	}

	return PID(list.lowest), true
}

// HighestPID returns the highest PID that is currently in use.
func (list *EntityList) HighestPID() (PID, bool) {
	if list.highest == -1 {
		return PID(0), false
	}

	return PID(list.highest), true
}

// Size returns the current amount of entities that are registered.
func (list *EntityList) Size() int {
	return list.size
}

// checkBoundaries checks if the given PID is out of bounds or not.
func (list *EntityList) checkBoundaries(pid PID) error {
	if pid < 0 || int(pid) >= len(list.entities) {
		return errors.New("PID out of bounds")
	}

	return nil
}
