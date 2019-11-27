package game

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
