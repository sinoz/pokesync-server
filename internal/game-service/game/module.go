package game

import (
	"gitlab.com/pokesync/game-service/internal/game-service/game/entity"
)

import "go.uber.org/zap"

// CommandCallback is a subscribable callback to register for a chat command.
type CommandCallback func(dk *DependencyKit, entity *entity.Entity, arguments []string) error

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

// OnCommand subscribes the given callback to the given trigger.
func (dk *DependencyKit) OnCommand(trigger string, cb CommandCallback) {
	dk.game.chatCommands.Put(trigger, func(entity *entity.Entity, arguments []string) error {
		return cb(dk, entity, arguments)
	})
}
