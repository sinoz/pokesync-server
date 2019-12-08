package game

import "gitlab.com/pokesync/game-service/internal/game-service/game/entity"

func attachFollower() attachFollowerHandler {
	return func(entity *entity.Entity, partySlot int) error {
		return nil
	}
}

func clearFollower() clearFollowerHandler {
	return func(entity *entity.Entity) error {
		return nil
	}
}
