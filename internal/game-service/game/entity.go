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

	Man        Gender = 0
	Woman      Gender = 1
	Genderless Gender = 2
)

// ModelID is the ID of an Entity's model.
type ModelID int

// EntityKind represents the type of an Entity.
type EntityKind int

// Gender is a type of gender of an Entity.
type Gender int

// EntityFactory is in charge of producing different types of entities.
type EntityFactory struct {
	world  *entity.World
	assets *AssetBundle
}

// NewEntityFactory constructs a new EntityFactory to produce entities with.
func NewEntityFactory(world *entity.World, assets *AssetBundle) *EntityFactory {
	return &EntityFactory{
		world:  world,
		assets: assets,
	}
}

// CreatePlayer creates the set of Component's to create a Player-like Entity from.
func (factory *EntityFactory) CreatePlayer(position Position, gender Gender, displayName character.DisplayName, userGroup character.UserGroup) *entity.Entity {
	return factory.world.
		CreateEntity().
		With(&TransformComponent{MovementQueue: NewMovementQueue(position)}).
		With(&UsernameComponent{DisplayName: displayName}).
		With(&RankComponent{UserGroup: userGroup}).
		With(&TrackingComponent{}).
		With(&MapViewComponent{MapView: NewMapView()}).
		With(&BicycleComponent{BicycleType: NoBike}).
		With(&CanRunComponent{}).
		With(&KindComponent{Kind: PlayerKind}).
		With(&CoinBagComponent{CoinBag: NewCoinBag()}).
		With(&PartyBeltComponent{PartyBelt: NewPartyBelt()}).
		With(&WaryOfTimeComponent{}).
		Build()
}

// CreateNpc creates the set of Component's to create a Npc-like Entity from.
func (factory *EntityFactory) CreateNpc(position Position, modelID ModelID) *entity.Entity {
	return factory.world.
		CreateEntity().
		With(&ModelIDComponent{ModelID: modelID}).
		With(&TransformComponent{MovementQueue: NewMovementQueue(position)}).
		With(&KindComponent{Kind: NpcKind}).
		With(&BlockingComponent{}).
		With(&TrackingComponent{}).
		Build()
}

// CreateMonster creates the set of Component's to create a Monster-like Entity from.
func (factory *EntityFactory) CreateMonster(position Position, modelID ModelID) *entity.Entity {
	return factory.world.
		CreateEntity().
		With(&ModelIDComponent{ModelID: modelID}).
		With(&TransformComponent{MovementQueue: NewMovementQueue(position)}).
		With(&TrackingComponent{}).
		With(&HealthComponent{Max: 1, Current: 1}). // TODO
		With(&KindComponent{Kind: MonsterKind}).
		Build()
}

// CreateObject creates the set of Component's to create a Object-like Entity from.
func (factory *EntityFactory) CreateObject(position Position) *entity.Entity {
	return factory.world.
		CreateEntity().
		With(&TransformComponent{MovementQueue: NewMovementQueue(position)}).
		Build()
}
