package game

import (
	"gitlab.com/pokesync/game-service/internal/game-service/character"
	"go.uber.org/zap"
)

// CommandCallback is a subscribable callback to register for a chat command.
type CommandCallback func(dk *DependencyKit, plr *Player, arguments []string) error

// DependencyKit holds a bundle of dependencies a Module may require
// throughout installation.
type DependencyKit struct {
	assets *AssetBundle
	game   *Game
	logger *zap.SugaredLogger
}

// Module is a module that can be defined externally and installed
// into the game service core.
type Module func(dk *DependencyKit)

// CreatePlayer creates a new Player-like Entity with the specified details.
func (dk *DependencyKit) CreatePlayer(position Position, gender Gender, displayName character.DisplayName, userGroup character.UserGroup) *Player {
	return dk.game.CreatePlayer(position, gender, displayName, userGroup)
}

// CreateMonster creates a new Monster-like Entity with the specified details.
func (dk *DependencyKit) CreateMonster(position Position, data MonsterData) *Monster {
	return dk.game.CreateMonster(position, data)
}

// CreateNpc creates a new Monster-like Entity with the specified details.
func (dk *DependencyKit) CreateNpc(id ModelID, position Position) *Npc {
	return dk.game.CreateNpc(id, position)
}

// AddPlayer attempts to add the given Player into the game world.
func (dk *DependencyKit) AddPlayer(plr *Player) bool {
	return dk.game.AddPlayer(plr)
}

// AddNpc attempts to add the given Npc into the game world.
func (dk *DependencyKit) AddNpc(npc *Npc) bool {
	return dk.game.AddNpc(npc)
}

// AddMonster attempts to add the given Monster into the game world.
func (dk *DependencyKit) AddMonster(mon *Monster) bool {
	return dk.game.AddMonster(mon)
}

// RemovePlayer removes the given Player from the game world.
func (dk *DependencyKit) RemovePlayer(plr *Player) {
	dk.game.RemovePlayer(plr)
}

// RemoveNpc removes the given Npc from the game world.
func (dk *DependencyKit) RemoveNpc(npc *Npc) {
	dk.game.RemoveNpc(npc)
}

// RemoveMonster removes the given Monster from the game world.
func (dk *DependencyKit) RemoveMonster(mon *Monster) {
	dk.game.RemoveMonster(mon)
}

// OnCommand subscribes the given callback to the given trigger.
func (dk *DependencyKit) OnCommand(trigger string, cb CommandCallback) {
	dk.game.chatCommands.Put(trigger, func(plr *Player, arguments []string) error {
		return cb(dk, plr, arguments)
	})
}
