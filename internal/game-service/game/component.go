package game

import (
	"gitlab.com/pokesync/game-service/internal/game-service/character"
	"gitlab.com/pokesync/game-service/internal/game-service/game/entity"
)

const (
	ModelIDTag   entity.ComponentTag = 0
	RankTag      entity.ComponentTag = 1
	UsernameTag  entity.ComponentTag = 2
	HealthTag    entity.ComponentTag = 3
	CanRunTag    entity.ComponentTag = 4
	TransformTag entity.ComponentTag = 5
	KindTag      entity.ComponentTag = 6
	TrackingTag  entity.ComponentTag = 7
	SessionTag   entity.ComponentTag = 8
	MapViewTag   entity.ComponentTag = 9
	BlockingTag  entity.ComponentTag = 10
	PartyBeltTag entity.ComponentTag = 11
	CoinBagTag   entity.ComponentTag = 12
)

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
	Nearby []*entity.Entity
}

// SessionComponent holds a Session instance, which indicates
// that the entity was created out of a request from a client user.
type SessionComponent struct {
	Session *Session
}

// MapViewComponent keeps track of an entity's map view.
type MapViewComponent struct {
}

// PartyBeltComponent is an entity Component that holds the PartyBelt.
type PartyBeltComponent struct {
	PartyBelt *PartyBelt
}

// CoinBagComponent is an entity Component that holds the CoinBag.
type CoinBagComponent struct {
	CoinBag *CoinBag
}

// Tag returns the tag of a Component instance for identification
// and storage purposes.
func (component *ModelIDComponent) Tag() entity.ComponentTag {
	return ModelIDTag
}

// Tag returns the tag of a Component instance for identification
// and storage purposes.
func (component *RankComponent) Tag() entity.ComponentTag {
	return RankTag
}

// Tag returns the tag of a Component instance for identification
// and storage purposes.
func (component *HealthComponent) Tag() entity.ComponentTag {
	return HealthTag
}

// Tag returns the tag of a Component instance for identification
// and storage purposes.
func (component *UsernameComponent) Tag() entity.ComponentTag {
	return UsernameTag
}

// Tag returns the tag of a Component instance for identification
// and storage purposes.
func (component *CanRunComponent) Tag() entity.ComponentTag {
	return CanRunTag
}

// Tag returns the tag of a Component instance for identification
// and storage purposes.
func (component *TransformComponent) Tag() entity.ComponentTag {
	return TransformTag
}

// Tag returns the tag of a Component instance for identification
// and storage purposes.
func (component *KindComponent) Tag() entity.ComponentTag {
	return KindTag
}

// Tag returns the tag of a Component instance for identification
// and storage purposes.
func (component *BlockingComponent) Tag() entity.ComponentTag {
	return BlockingTag
}

// Tag returns the tag of a Component instance for identification
// and storage purposes.
func (component *TrackingComponent) Tag() entity.ComponentTag {
	return TrackingTag
}

// Tag returns the tag of a Component instance for identification
// and storage purposes.
func (component *SessionComponent) Tag() entity.ComponentTag {
	return SessionTag
}

// Tag returns the tag of a Component instance for identification
// and storage purposes.
func (component *MapViewComponent) Tag() entity.ComponentTag {
	return MapViewTag
}

// Tag returns the tag of a Component instance for identification
// and storage purposes.
func (component *PartyBeltComponent) Tag() entity.ComponentTag {
	return PartyBeltTag
}

// Tag returns the tag of a Component instance for identification
// and storage purposes.
func (component *CoinBagComponent) Tag() entity.ComponentTag {
	return CoinBagTag
}
