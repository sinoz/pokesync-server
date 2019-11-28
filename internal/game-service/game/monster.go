package game

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

// Monster represents a monster.
type Monster struct {
	ID              MonsterID
	Gender          Gender
	StatusCondition StatusCondition
	Coloration      MonsterColoration
}
