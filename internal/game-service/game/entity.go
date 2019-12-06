package game

import (
	"gitlab.com/pokesync/game-service/internal/game-service/character"
	"gitlab.com/pokesync/game-service/internal/game-service/game/entity"
)

const (
	PlayerKind  EntityKind = 0
	NpcKind     EntityKind = 1
	MonsterKind EntityKind = 2
	ObjectKind  EntityKind = 3

	Man   Gender = 0
	Woman Gender = 1
)

// ModelID is the ID of an Entity's model.
type ModelID int

// EntityKind represents the type of an Entity.
type EntityKind int

// Gender is a type of gender of an Entity.
type Gender int

// EntityFactory is in charge of producing different types of entities.
type EntityFactory struct {
	assets *AssetBundle
}

// NewEntityFactory constructs a new EntityFactory to produce entities with.
func NewEntityFactory(assets *AssetBundle) *EntityFactory {
	return &EntityFactory{assets: assets}
}

// CreatePlayer creates the set of Component's to create a Player-like Entity from.
func (factory *EntityFactory) CreatePlayer(position Position, direction Direction, gender Gender, displayName character.DisplayName, userGroup character.UserGroup) []entity.Component {
	return []entity.Component{
		&TransformComponent{Position: position},
		&UsernameComponent{DisplayName: displayName},
		&RankComponent{UserGroup: userGroup},
		&KindComponent{Kind: PlayerKind},
	}
}

// CreateNpc creates the set of Component's to create a Npc-like Entity from.
func (factory *EntityFactory) CreateNpc(position Position, direction Direction, modelID ModelID) []entity.Component {
	return []entity.Component{
		&ModelIDComponent{ModelID: modelID},
		&TransformComponent{Position: position},
		&KindComponent{Kind: NpcKind},
		&BlockingComponent{},
		&TrackingComponent{},
	}
}

// CreateMonster creates the set of Component's to create a Monster-like Entity from.
func (factory *EntityFactory) CreateMonster(position Position, direction Direction, modelID ModelID) []entity.Component {
	return []entity.Component{
		&ModelIDComponent{ModelID: modelID},
		&TransformComponent{Position: position},
		&HealthComponent{Max: 1, Current: 1}, // TODO
		&KindComponent{Kind: MonsterKind},
	}
}

// CreateObject creates the set of Component's to create a Object-like Entity from.
func (factory *EntityFactory) CreateObject(position Position, direction Direction) []entity.Component {
	return []entity.Component{}
}
