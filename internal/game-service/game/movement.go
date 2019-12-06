package game

import "gitlab.com/pokesync/game-service/internal/game-service/game/entity"

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

// faceDirection tells the given Entity to change its currently
// facing direction to the specified one.
func faceDirection(entity *entity.Entity, direction Direction) {
	// TODO
}

// changeMovementType tells the given Entity to change its current
// MovementType to that of the specified one.
func changeMovementType(entity *entity.Entity, movementType MovementType) {
	// TODO
}

// moveAvatar tells the given avatar Entity to move to a tile that
// lays in the specified direction from the Entity's current position.
func moveAvatar(entity *entity.Entity, direction Direction) {
	// TODO
}
