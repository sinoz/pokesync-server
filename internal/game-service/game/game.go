package game

import (
	"gitlab.com/pokesync/game-service/internal/game-service/character"
	"gitlab.com/pokesync/game-service/internal/game-service/game/entity"
	"gitlab.com/pokesync/game-service/pkg/event"
)

// Game represents the game, mkay.
type Game struct {
	world         *entity.World
	entityFactory *EntityFactory
	eventBus      event.Bus
	grid          *Grid
}

// NewGame constructs a new Game.
func NewGame(assets *AssetBundle, entityCapacity int) *Game {
	return &Game{
		world:         entity.NewWorld(entityCapacity),
		entityFactory: NewEntityFactory(assets),
		eventBus:      event.NewSerialBus(),
	}
}

// AddPlayer adds a player Entity with the specified details.
func (game *Game) AddPlayer(position Position, gender Gender, displayName character.DisplayName, userGroup character.UserGroup) (*entity.Entity, bool) {
	components := game.entityFactory.CreatePlayer(position, South, gender, displayName, userGroup)

	return game.
		world.
		CreateEntity().
		With(components...).
		Build()
}

// AddNpc adds a npc-like Entity with the specified details.
func (game *Game) AddNpc(modelID ModelID, position Position) (*entity.Entity, bool) {
	components := game.entityFactory.CreateNpc(position, South, modelID)

	return game.
		world.
		CreateEntity().
		With(components...).
		Build()
}

// AddMonster adds a monster-like Entity with the specified details.
func (game *Game) AddMonster(modelID ModelID, position Position) (*entity.Entity, bool) {
	components := game.entityFactory.CreateMonster(position, South, modelID)

	return game.
		world.
		CreateEntity().
		With(components...).
		Build()
}
