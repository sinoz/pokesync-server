package game

import (
	"gitlab.com/pokesync/game-service/internal/game-service/character"
	"gitlab.com/pokesync/game-service/internal/game-service/game/entity"
)

// Player is a type of Entity.
type Player struct {
	*entity.Entity
}

// PlayerBy wraps the given Entity as a Player.
func PlayerBy(entity *entity.Entity) *Player {
	return &Player{entity}
}

// Face updates the Player's character sprite to face the speciifed direction.
func (plr *Player) Face(direction Direction) {
	transform := plr.GetComponent(TransformTag).(*TransformComponent)
	transform.MovementQueue.AddDirectionToFace(direction)
}

// Move tells the Player to move towards the specified Direction.
func (plr *Player) Move(direction Direction) {
	transform := plr.GetComponent(TransformTag).(*TransformComponent)
	transform.MovementQueue.AddStep(direction)
}

// MoveTo tells the Player to move to the specified coordinates.
func (plr *Player) MoveTo(mapX, mapZ, localX, localZ int) {
	// TODO
}

// Walk tells the Player to walk from now on.
func (plr *Player) Walk() {
	// TODO
}

// Run tells the Player to run from now on.
func (plr *Player) Run() {
	// TODO
}

// HasBicycle returns whether the Player owns a bicycle to ride on.
func (plr *Player) HasBicycle() bool {
	return false
}

// Cycle tells the Player to start cycling from now on. Returns false
// if the Player does not own a bicycle.
func (plr *Player) Cycle() bool {
	if !plr.HasBicycle() {
		return false
	}

	// TODO
	return true
}

// DisplayName returns the player's display name.
func (plr *Player) DisplayName() character.DisplayName {
	return plr.GetComponent(UsernameTag).(*UsernameComponent).DisplayName
}

// Rank returns the player's rank or UserGroup, which the user is
// associated with.
func (plr *Player) Rank() character.UserGroup {
	return plr.GetComponent(RankTag).(*RankComponent).UserGroup
}

// Position returns the player's current Position on the game map.
func (plr *Player) Position() Position {
	return plr.GetComponent(TransformTag).(*TransformComponent).MovementQueue.Position
}

// CoinBag returns the player's bag of coins.
func (plr *Player) CoinBag() *CoinBag {
	return plr.GetComponent(CoinBagTag).(*CoinBagComponent).CoinBag
}

// PartyBelt returns the player's belt of party monsters.
func (plr *Player) PartyBelt() *PartyBelt {
	return plr.GetComponent(PartyBeltTag).(*PartyBeltComponent).PartyBelt
}
