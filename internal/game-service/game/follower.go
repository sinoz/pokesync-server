package game

func attachFollower() attachFollowerHandler {
	return func(plr *Player, partySlot int) error {
		return nil
	}
}

func clearFollower() clearFollowerHandler {
	return func(plr *Player) error {
		return nil
	}
}
