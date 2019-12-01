package game

import (
	ecs "gitlab.com/pokesync/ecs/src"
	"gitlab.com/pokesync/game-service/internal/game-service/character"
	"gitlab.com/pokesync/game-service/internal/game-service/game/session"
)

const (
	PIDTag       ecs.ComponentTag = 1 << 0
	ModelIDTag   ecs.ComponentTag = 1 << 1
	RankTag      ecs.ComponentTag = 1 << 2
	UsernameTag  ecs.ComponentTag = 1 << 3
	HealthTag    ecs.ComponentTag = 1 << 4
	CanRunTag    ecs.ComponentTag = 1 << 5
	TransformTag ecs.ComponentTag = 1 << 6
	KindTag      ecs.ComponentTag = 1 << 7
	TrackingTag  ecs.ComponentTag = 1 << 8
	SessionTag   ecs.ComponentTag = 1 << 9
	MapViewTag   ecs.ComponentTag = 1 << 10
	BlockingTag  ecs.ComponentTag = 1 << 11
)

// PIDComponent holds a process id of an entity.
type PIDComponent struct {
	PID PID
}

// ModelIDComponent holds a model id of an entity.
type ModelIDComponent struct {
	ModelID ModelID
}

// RankComponent holds the UserGroup of an entity.
type RankComponent struct {
	UserGroup character.UserGroup
}

// UsernameComponent holds a display name of a player entity.
type UsernameComponent struct {
	DisplayName character.DisplayName
}

// HealthComponent keeps track of how much health an entity has left.
type HealthComponent struct {
	Max     int
	Current int
}

// CanRunComponent marks an entity as being able to run.
type CanRunComponent struct{}

// TransformComponent holds the entity's Position in the game world
// and keeps track of its recent movements.
type TransformComponent struct {
	Position Position
}

// BlockingComponent marks an Entity as blocking all other entities paths.
type BlockingComponent struct{}

// KindComponent holds the EntityKind of the Entity, which is used
// to check what kind of Entity it is (player, npc, monster, obj etc).
type KindComponent struct {
	Kind EntityKind
}

// TrackingComponent keeps track of entities that are nearby the
// Entity this Component is for.
type TrackingComponent struct {
	Nearby []*ecs.Entity
}

// SessionComponent holds a Session instance, which indicates
// that the entity was created out of a request from a client user.
type SessionComponent struct {
	Session *session.Session
}

// MapViewComponent keeps track of an entity's map view.
type MapViewComponent struct {
}

// Tag returns the tag of a Component instance for identification
// and storage purposes.
func (component *PIDComponent) Tag() ecs.ComponentTag {
	return PIDTag
}

// Tag returns the tag of a Component instance for identification
// and storage purposes.
func (component *ModelIDComponent) Tag() ecs.ComponentTag {
	return ModelIDTag
}

// Tag returns the tag of a Component instance for identification
// and storage purposes.
func (component *RankComponent) Tag() ecs.ComponentTag {
	return RankTag
}

// Tag returns the tag of a Component instance for identification
// and storage purposes.
func (component *HealthComponent) Tag() ecs.ComponentTag {
	return HealthTag
}

// Tag returns the tag of a Component instance for identification
// and storage purposes.
func (component *UsernameComponent) Tag() ecs.ComponentTag {
	return UsernameTag
}

// Tag returns the tag of a Component instance for identification
// and storage purposes.
func (component *CanRunComponent) Tag() ecs.ComponentTag {
	return CanRunTag
}

// Tag returns the tag of a Component instance for identification
// and storage purposes.
func (component *TransformComponent) Tag() ecs.ComponentTag {
	return TransformTag
}

// Tag returns the tag of a Component instance for identification
// and storage purposes.
func (component *KindComponent) Tag() ecs.ComponentTag {
	return KindTag
}

// Tag returns the tag of a Component instance for identification
// and storage purposes.
func (component *BlockingComponent) Tag() ecs.ComponentTag {
	return BlockingTag
}

// Tag returns the tag of a Component instance for identification
// and storage purposes.
func (component *TrackingComponent) Tag() ecs.ComponentTag {
	return TrackingTag
}

// Tag returns the tag of a Component instance for identification
// and storage purposes.
func (component *SessionComponent) Tag() ecs.ComponentTag {
	return SessionTag
}

// Tag returns the tag of a Component instance for identification
// and storage purposes.
func (component *MapViewComponent) Tag() ecs.ComponentTag {
	return MapViewTag
}
