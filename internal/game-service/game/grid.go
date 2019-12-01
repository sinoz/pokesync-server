package game

import (
	"fmt"
	"math"

	"gitlab.com/pokesync/game-service/internal/game-service/game/collision"
)

// These are the supported types of Direction's.
const (
	South Direction = 0
	North Direction = 1
	West  Direction = 2
	East  Direction = 3
)

// Position describes an exact tile position on the game map.
type Position struct {
	MapX int
	MapZ int

	LocalX int
	LocalZ int

	Altitude int
}

// TileMap is a map of tiles in the game world for entities to traverse.
type TileMap struct {
	Index           MapIndex
	CollisionMatrix *collision.Matrix
}

// Grid represents the physical grid of TileMap's that together
// make up the game world.
type Grid struct {
	TileMaps [][]*TileMap
}

// Direction is a compass direction.
type Direction int

// GetMap looks up a TileMap at the specified coordinates. May return
// nil if the given coordinates fall out of bounds of the grid.
func (grid *Grid) GetMap(x, z int) (*TileMap, error) {
	if err := grid.checkBoundaries(x, z); err != nil {
		return nil, err
	}

	return grid.TileMaps[x][z], nil
}

// Width returns the width of the grid, in tile maps.
func (grid *Grid) Width() int {
	return len(grid.TileMaps)
}

// Length returns the length of the grid, in tile maps.
func (grid *Grid) Length() int {
	return len(grid.TileMaps[0])
}

// checkBoundaries returns whether the given x/z coordinates fall out
// of bounds of the Grid.
func (grid *Grid) checkBoundaries(x, z int) error {
	if x < 0 || z < 0 || x >= grid.Width() || z >= grid.Length() {
		return fmt.Errorf("coordinates %v %v fall out of bounds", x, z)
	}

	return nil
}

// DistanceBetween calculates the distance between the two given Position
// points on the given world Grid.
func DistanceBetween(p1, p2 Position, grid *Grid) (int, error) {
	mapA, err := grid.GetMap(p1.MapX, p1.MapZ)
	if err != nil {
		return 0, err
	}

	mapB, err := grid.GetMap(p2.MapX, p2.MapZ)
	if err != nil {
		return 0, err
	}

	renderXOfPosA := mapA.Index.RenderX + p1.LocalX
	renderZOfPosA := mapA.Index.RenderZ + p1.LocalZ

	renderXOfPosB := mapB.Index.RenderX + p2.LocalX
	renderZOfPosB := mapB.Index.RenderZ + p2.LocalZ

	deltaX := renderXOfPosB - renderXOfPosA
	deltaZ := renderZOfPosB - renderZOfPosA

	return int(math.Sqrt(float64(deltaX*2) + float64(deltaZ*2))), nil
}

// AddStep adds a single step to the given Position in the specified Direction
// on the given map Grid. If the step falls out of bounds of the Position's
// current map, the location of the current map is updated. May return an error
// if the map for the specified Position does not exist or if the step falls
// into a map that doesn't exist on the Grid.
func AddStep(position Position, direction Direction, grid *Grid) (Position, error) {
	tileMap, err := grid.GetMap(position.MapX, position.MapZ)
	if err != nil {
		return position, err
	}

	mapWidth := tileMap.CollisionMatrix.Width()
	mapLength := tileMap.CollisionMatrix.Length()

	mapX := position.MapX
	mapZ := position.MapZ

	localX := position.LocalX
	localZ := position.LocalZ

	switch direction {
	case North:
		localZ++
		break
	case South:
		localZ--
		break
	case East:
		localX++
		break
	case West:
		localX--
		break
	}

	if localX < 0 {
		mapX--
		localX = mapWidth - 1

		_, err := grid.GetMap(mapX, mapZ)
		if err != nil {
			return position, err
		}
	} else if localX >= mapWidth {
		mapX++
		localX = 0

		_, err := grid.GetMap(mapX, mapZ)
		if err != nil {
			return position, err
		}
	} else if localZ < 0 {
		mapZ--
		localZ = mapLength - 1

		_, err := grid.GetMap(mapX, mapZ)
		if err != nil {
			return position, err
		}
	} else if localZ >= mapLength {
		mapZ++
		localZ = 0

		_, err := grid.GetMap(mapX, mapZ)
		if err != nil {
			return position, err
		}
	}

	return Position{
		MapX:     mapX,
		MapZ:     mapZ,
		LocalX:   localX,
		LocalZ:   localZ,
		Altitude: position.Altitude,
	}, nil
}
