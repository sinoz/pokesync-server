package character

// DisplayName is a name of a user's account that is exposed
// ingame to other users.
type DisplayName string

// UserGroup is a type of group a user may belong to.
type UserGroup int

var (
	Regular       UserGroup = 0
	Patron        UserGroup = 1
	Moderator     UserGroup = 2
	Administrator UserGroup = 3
	GameDesigner  UserGroup = 4
	WebDeveloper  UserGroup = 5
	GameDeveloper UserGroup = 6
)

// Profile represents a player character's last saved game state.
type Profile struct {
	DisplayName DisplayName
	UserGroup   UserGroup
	Gender      int

	MapX   int
	MapZ   int
	LocalX int
	LocalZ int
}
