package character

import (
	"time"
)

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
	DisplayName  DisplayName `json:"displayName"`
	UserGroup    UserGroup   `json:"userGroup"`
	LastLoggedIn *time.Time  `json:"lastLoggedIn"`

	Gender int `json:"gender"`

	MapX   int `json:"mapX"`
	MapZ   int `json:"mapZ"`
	LocalX int `json:"localX"`
	LocalZ int `json:"localZ"`
}
