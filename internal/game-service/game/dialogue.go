package game

import "gitlab.com/pokesync/game-service/internal/game-service/game/entity"

func continueDialogue() continueDialogueHandler {
	return func(entity *entity.Entity) error {
		return nil
	}
}
