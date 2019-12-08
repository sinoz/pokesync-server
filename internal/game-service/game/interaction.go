package game

import "gitlab.com/pokesync/game-service/internal/game-service/game/entity"

func selectPlayerOption() selectPlayerOptionHandler {
	return func(entity *entity.Entity, entityID entity.ID, slot int) error {
		return nil
	}
}

func interactWithEntity() interactWithEntityHandler {
	return func(entity *entity.Entity, entityID entity.ID) error {
		return nil
	}
}
