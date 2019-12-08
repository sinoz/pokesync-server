package commands

import (
	"gitlab.com/pokesync/game-service/internal/game-service/game"
)

// Module is an externally defined module to subscribe chat commands with.
func Module(dk *game.DependencyKit) {
	dk.OnCommand("pos", showPosition)
}
