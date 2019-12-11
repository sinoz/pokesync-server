package commands

import (
	"gitlab.com/pokesync/game-service/internal/game-service/game"
	"strconv"
)

func addToParty(dk *game.DependencyKit, plr *game.Player, arguments []string) error {
	if len(arguments) < 1 {
		return nil
	}

	modelID, err := strconv.Atoi(arguments[0])
	if err != nil {
		return err
	}

	plr.PartyBelt().Add(dk.CreateMonster(plr.Position(), game.MonsterData{
		ModelID:         game.ModelID(modelID),
		Gender:          game.Man,
		StatusCondition: game.Healthy,
		Coloration:      game.RegularColour,
	}))

	return nil
}

func clearParty(dk *game.DependencyKit, plr *game.Player, arguments []string) error {
	plr.PartyBelt().ClearAll()
	return nil
}
