package game

import (
	"gitlab.com/pokesync/game-service/internal/game-service/game/entity"
	"gitlab.com/pokesync/game-service/internal/game-service/game"
)

// Npc is a type of Entity.
type Npc struct {
	*entity.Entity
}

// NpcBy wraps the given Entity as a Npc.
func NpcBy(entity *entity.Entity) *Npc {
	return &Npc{entity}
}

// Face updates the Npc's character sprite to face the speciifed direction.
func (npc *Npc) Face(direction Direction) {
	// TODO
}

// MoveTo tells the Npc to move to the specified coordinates.
func (npc *Npc) MoveTo(mapX, mapZ, localX, localZ int) {
	// TODO
}

// ModelID returns the npc's model id.
func (npc *Npc) ModelID() game.ModelID {
	return npc.GetComponent(ModelIDTag).(*ModelIDComponent).ModelID
}

// Position returns the npc's current Position on the game map.
func (npc *Npc) Position() Position {
	return npc.GetComponent(TransformTag).(*TransformComponent).Position
}
