package game

// The type of genders.
const (
	Man   Gender = 0
	Woman Gender = 1
)

// Gender is a type of gender of an entity.
type Gender int

// The types of monster colorations
const (
	RegularColour MonsterColoration = 0
	ShinyColour   MonsterColoration = 1
)

// MonsterColoration defines the skin color of the monster ingame,
// which may raise its value if there is an invariance.
type MonsterColoration int

// StatusCondition describes a condition a monster is in.
type StatusCondition int

// These are the supported types of StatusCondition's.
const (
	Healthy   StatusCondition = 0
	Paralyzed StatusCondition = 1
	Poisoned  StatusCondition = 2
	Asleep    StatusCondition = 3
	Frozen    StatusCondition = 4
	Burnt     StatusCondition = 5
)

// MovementType is a type of movement an entity can perform.
type MovementType int

// These are the supported types of movement.
const (
	Walk     MovementType = 0
	Run      MovementType = 1
	Cycle    MovementType = 2
	Teleport MovementType = 3
	Jump     MovementType = 4
	Surf     MovementType = 5
	Dive     MovementType = 6
	Glide    MovementType = 7
)

// Movement is a movement between two points in the game world.
type Movement struct {
	Source      Position
	Destination Position
	Type        MovementType
}
