package game

import (
	"math"
	"time"

	"gitlab.com/pokesync/game-service/internal/game-service/game/entity"
	"gitlab.com/pokesync/game-service/internal/game-service/game/transport"
)

const (
	MapViewWidth  = 3
	MapViewLength = 3

	SearchMargin = 1
)

// MapViewChangeListener TODO
type MapViewChangeListener interface {
	Refreshed(mapX, mapZ int)
}

// MapViewSessionListener is a MapViewChangeListener that listens
// for changes made to the MapView to visually apply these
// changes to the Client as well.
type MapViewSessionListener struct {
	session *Session
}

// MapView is the viewable and roamable area of tile maps
// an avatar can traverse without refreshing.
type MapView struct {
	TileMaps         [][]*TileMap
	pendingRefreshes []MapRefresh
	listeners        []MapViewChangeListener
}

// MapRefresh is a refresh of the game map.
type MapRefresh struct {
	MapX int
	MapZ int
}

// MapViewProcessor processes map view changes.
type MapViewProcessor struct {
	grid *Grid
}

// NewMapView constructs a new map view.
func NewMapView() *MapView {
	tileMaps := make([][]*TileMap, MapViewWidth)
	for x := 0; x < MapViewWidth; x++ {
		tileMaps[x] = make([]*TileMap, MapViewLength)
	}

	return &MapView{TileMaps: tileMaps}
}

// NewMapViewSystem constructs a new instance of an entity.System with
// a MapViewProcessor as its internal processor.
func NewMapViewSystem(grid *Grid) *entity.System {
	return entity.NewSystem(entity.NewDefaultSystemPolicy(), NewMapViewProcessor(grid))
}

// NewMapViewProcessor constructs a new instance of a MapViewProcessor.
func NewMapViewProcessor(grid *Grid) *MapViewProcessor {
	return &MapViewProcessor{grid: grid}
}

// AddedToWorld is called when the System of this Processor is added
// to the game World.
func (processor *MapViewProcessor) AddedToWorld(world *entity.World) error {
	return nil
}

// RemovedFromWorld is called when the System of this Processor is removed
// from the game World.
func (processor *MapViewProcessor) RemovedFromWorld(world *entity.World) error {
	return nil
}

// Update is called every game pulse to check if entities need their map view
// refreshed and if so, refreshes them.
func (processor *MapViewProcessor) Update(world *entity.World, deltaTime time.Duration) error {
	entities := world.GetEntitiesFor(processor)
	for _, ent := range entities {
		mapView := ent.GetComponent(MapViewTag).(*MapViewComponent).MapView
		refresh := mapView.PollRefresh()
		if refresh == nil {
			continue
		}

		lowerBoundMapX := int(math.Max(0, float64(refresh.MapX-SearchMargin)))
		lowerBoundMapZ := int(math.Max(0, float64(refresh.MapZ-SearchMargin)))

		upperBoundMapX := int(math.Min(float64(processor.grid.Width()), float64(refresh.MapX+SearchMargin)))
		upperBoundMapZ := int(math.Min(float64(processor.grid.Length()), float64(refresh.MapZ+SearchMargin)))

		for x := lowerBoundMapX; x <= upperBoundMapX; x++ {
			for z := lowerBoundMapZ; z <= upperBoundMapZ; z++ {
				gridX := x - lowerBoundMapX
				gridZ := z - lowerBoundMapZ

				tileMap, _ := processor.grid.GetMap(x, z)
				mapView.TileMaps[gridX][gridZ] = tileMap
			}
		}

		mapView.notifyRefresh(refresh.MapX, refresh.MapZ)
	}

	return nil
}

// Components returns a pack of ComponentTag's the MapViewProcessor has
// interest in.
func (processor *MapViewProcessor) Components() entity.ComponentTag {
	return MapViewTag
}

func (mapView *MapView) PollRefresh() *MapRefresh {
	if len(mapView.pendingRefreshes) == 0 {
		return nil
	}

	refresh := mapView.pendingRefreshes[0]
	mapView.pendingRefreshes = mapView.pendingRefreshes[1:]
	return &refresh
}

func (mapView *MapView) Refresh(mapX, mapZ int) {
	mapView.pendingRefreshes = append(mapView.pendingRefreshes, MapRefresh{
		MapX: mapX,
		MapZ: mapZ,
	})
}

func (mapView *MapView) AddListener(listener MapViewChangeListener) {
	mapView.listeners = append(mapView.listeners, listener)
}

func (mapView *MapView) RemoveListener(listener MapViewChangeListener) {
	for i, l := range mapView.listeners {
		if l == listener {
			mapView.listeners = append(mapView.listeners[:i], mapView.listeners[i+1:]...)
			break
		}
	}
}

func (mapView *MapView) notifyRefresh(mapX, mapZ int) {
	for _, l := range mapView.listeners {
		l.Refreshed(mapX, mapZ)
	}
}

// Refreshed sends a visual update to the Session.
func (listener *MapViewSessionListener) Refreshed(mapX, mapZ int) {
	listener.session.QueueEvent(&transport.RefreshMap{
		MapX: uint16(mapX),
		MapZ: uint16(mapZ),
	})
}
