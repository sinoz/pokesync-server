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
	BicycleTag   entity.ComponentTag = 4
	CanRunTag    entity.ComponentTag = 5
	TransformTag entity.ComponentTag = 6
	KindTag      entity.ComponentTag = 7
	TrackingTag  entity.ComponentTag = 8
	SessionTag   entity.ComponentTag = 9
	MapViewTag   entity.ComponentTag = 10
	BlockingTag  entity.ComponentTag = 11
	PartyBeltTag entity.ComponentTag = 12
	CoinBagTag   entity.ComponentTag = 13
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

// BicycleComponent holds the type of Bicycle an Entity owns.
type BicycleComponent struct {
	BicycleType BicycleType
}

// CanRunComponent marks an entity as being able to run.
type CanRunComponent struct{}

// TransformComponent holds the entity's Position in the game world
// and keeps track of its recent movements.
type TransformComponent struct {
	MovementQueue *MovementQueue
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
	nearby []*entity.Entity
}

// SessionComponent holds a Session instance, which indicates
// that the entity was created out of a request from a client user.
type SessionComponent struct {
	session *Session
}

// MapViewComponent keeps track of an entity's map view.
type MapViewComponent struct {
	MapView *MapView
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
func (component *BicycleComponent) Tag() entity.ComponentTag {
	return BicycleTag
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
