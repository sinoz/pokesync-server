package commands

import (
	"fmt"
	"gitlab.com/pokesync/game-service/internal/game-service/game"
)

func showPosition(dk *game.DependencyKit, plr *game.Player, arguments []string) error {
	position := plr.Position()
	fmt.Println(position.MapX, position.MapZ, position.LocalX, position.LocalZ)
	return nil
}
