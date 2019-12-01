package game

import (
	ecs "gitlab.com/pokesync/ecs/src"
	"gitlab.com/pokesync/game-service/internal/game-service/character"
	"gitlab.com/pokesync/game-service/pkg/event"
)

// Game represents the game, mkay.
type Game struct {
	world         *ecs.World
	eventBus      event.Bus
	entityList    *EntityList
	entityFactory *EntityFactory
}

// NewGame constructs a new Game.
func NewGame(assets *AssetBundle, entityCapacity int) *Game {
	return &Game{
		world:         ecs.NewWorld(entityCapacity),
		eventBus:      event.NewSerialBus(),
		entityList:    NewEntityList(entityCapacity),
		entityFactory: NewEntityFactory(assets),
	}
}

// AddPlayer adds a player Entity with the specified details.
func (game *Game) AddPlayer(pid PID, position Position, gender Gender, displayName character.DisplayName, userGroup character.UserGroup) {
	components := game.entityFactory.CreatePlayer(pid, position, South, gender, displayName, userGroup)

	game.
		world.
		CreateEntity().
		With(components...).
		Build()
}

// AddNpc adds a npc-like Entity with the specified details.
func (game *Game) AddNpc(pid PID, modelID ModelID, position Position) {
	components := game.entityFactory.CreateNpc(pid, position, South, modelID)

	game.
		world.
		CreateEntity().
		With(components...).
		Build()
}

// AddMonster adds a monster-like Entity with the specified details.
func (game *Game) AddMonster(pid PID, modelID ModelID, position Position) {
	components := game.entityFactory.CreateMonster(pid, position, South, modelID)

	game.
		world.
		CreateEntity().
		With(components...).
		Build()
}
