package game

import (
	"time"

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
func NewGame(assets *AssetBundle, world *entity.World) *Game {
	return &Game{
		world:         world,
		entityFactory: NewEntityFactory(assets),
		eventBus:      event.NewSerialBus(),
	}
}

// pulse is called every game pulse to process the game.
func (game *Game) pulse(deltaTime time.Duration) error {
	return game.world.Update(deltaTime)
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

// RemovePlayer removes the given Player-like entity.
func (game *Game) RemovePlayer(entity *entity.Entity) {
	game.RemoveEntity(entity)
}

// RemoveNpc removes the given Npc-like entity.
func (game *Game) RemoveNpc(entity *entity.Entity) {
	game.RemoveEntity(entity)
}

// RemoveMonster removes the given Monster-like entity.
func (game *Game) RemoveMonster(entity *entity.Entity) {
	game.RemoveEntity(entity)
}

// RemoveEntity removes the given Entity from the game world.
func (game *Game) RemoveEntity(entity *entity.Entity) {
	game.world.DestroyEntity(entity)
}
