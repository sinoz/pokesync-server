package character

// DisplayName is a name of a user's account that is exposed
// ingame to other users.
type DisplayName string

// Profile represents a player character's last saved game state.
type Profile struct {
	DisplayName DisplayName

	MapX   int
	MapZ   int
	LocalX int
	LocalZ int
}
