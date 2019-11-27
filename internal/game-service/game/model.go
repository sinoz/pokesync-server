package game

const (
	Man   Gender = 0
	Woman Gender = 1

	RegularColour MonsterColoration = 0
	ShinyColour   MonsterColoration = 1

	Healthy   StatusCondition = 0
	Paralyzed StatusCondition = 1
	Poisoned  StatusCondition = 2
	Asleep    StatusCondition = 3
	Frozen    StatusCondition = 4
	Burnt     StatusCondition = 5

	Player  EntityKind = 0
	Npc     EntityKind = 1
	Monster EntityKind = 2
	Object  EntityKind = 3
)

// Gender is a type of gender of an entity.
type Gender int

// MonsterColoration defines the skin color of the monster ingame,
// which may raise its value if there is an invariance.
type MonsterColoration int

// StatusCondition describes a condition a monster is in.
type StatusCondition int

// EntityKind represents the type of an Entity.
type EntityKind int
