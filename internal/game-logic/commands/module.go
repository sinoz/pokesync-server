package commands

import (
	"fmt"

	"gitlab.com/pokesync/game-service/internal/game-service/game"
	"gitlab.com/pokesync/game-service/internal/game-service/game/entity"
)

func showPosition(dk *game.DependencyKit, entity *entity.Entity, arguments []string) error {
	component := entity.GetComponent(game.TransformTag).(*game.TransformComponent)
	position := component.Position

	fmt.Println(position.MapX, position.MapZ, position.LocalX, position.LocalZ)
	return nil
}

// Module is an externally defined module to subscribe chat commands with.
func Module(dk *game.DependencyKit) {
	dk.OnCommand("pos", showPosition)
}
