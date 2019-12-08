package game

import "gitlab.com/pokesync/game-service/internal/game-service/game/entity"

const (
	RegularColour MonsterColoration = 0
	ShinyColour   MonsterColoration = 1

	Healthy   StatusCondition = 0
	Paralyzed StatusCondition = 1
	Poisoned  StatusCondition = 2
	Asleep    StatusCondition = 3
	Frozen    StatusCondition = 4
	Burnt     StatusCondition = 5
)

// MonsterID is a unique id of a monster.
type MonsterID int

// MonsterColoration defines the skin color of the monster ingame,
// which may raise its value if there is an invariance.
type MonsterColoration int

// StatusCondition describes a condition a monster is in.
type StatusCondition int

// Monster is a type of Entity.
type Monster struct {
	*entity.Entity

	ID              MonsterID
	Gender          Gender
	StatusCondition StatusCondition
	Coloration      MonsterColoration
}

// MonsterBy wraps the given Entity as a Monster.
func MonsterBy(entity *entity.Entity, id MonsterID, gender Gender, condition StatusCondition, coloration MonsterColoration) *Monster {
	return &Monster{entity, id, gender, condition, coloration}
}

// Position returns the monster's current Position on the game map.
func (mon *Monster) Position() Position {
	return mon.GetComponent(TransformTag).(*TransformComponent).Position
}
