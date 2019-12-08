package game

import (
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"strconv"
	"strings"

	"gitlab.com/pokesync/game-service/internal/game-service/game/collision"
	"gitlab.com/pokesync/game-service/pkg/bytes"
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

// LoadGridFromConfig loads a Grid of TileMap's using the given WorldConfig
// as a directive. May return an error.
func LoadGridFromConfig(config WorldConfig) (*Grid, error) {
	grid := NewGrid(config.Width, config.Length)

	for regionID, regionIndex := range config.Regions {
		for _, mapIndex := range regionIndex.Maps {
			tileMap, err := LoadTileMap("assets/tile/", regionID, mapIndex)
			if err != nil {
				return nil, err
			}

			grid.TileMaps[mapIndex.MapX][mapIndex.MapZ] = tileMap
		}
	}

	return grid, nil
}

// NewGrid constructs a new Grid of TileMap's of the specified dimensions.
func NewGrid(width, length int) *Grid {
	grid := make([][]*TileMap, width)
	for x := 0; x < width; x++ {
		grid[x] = make([]*TileMap, length)
	}

	return &Grid{TileMaps: grid}
}

// LoadTileMap loads a TileMap from a file. May return an error.
func LoadTileMap(directory string, regionID int, index MapIndex) (*TileMap, error) {
	mapDataFilePath := fmt.Sprint(directory, regionID, "_", index.MapX, "_", index.MapZ, ".dat")
	fileBytes, err := ioutil.ReadFile(mapDataFilePath)
	if err != nil {
		return nil, err
	}

	bytes := bytes.StringWrap(fileBytes)
	iterator := bytes.Iterator()
	if !iterator.CanRead(4) {
		return nil, errors.New("not enough bytes for file stamp")
	}

	dataFileStamp, _ := iterator.ReadString(4)
	if dataFileStamp != "PSTM" {
		return nil, fmt.Errorf("Unexpected data file stamp of %v, expected PSTM", dataFileStamp)
	}

	headerLength, err := iterator.ReadByte()
	if err != nil {
		return nil, err
	}

	if !iterator.CanRead(int(headerLength)) {
		return nil, errors.New("not enough bytes to read header")
	}

	iterator.ReadByte() // export version

	payloadLength, _ := iterator.ReadUInt32()
	mapCount, _ := iterator.ReadUInt16()

	if payloadLength < 0 || !iterator.CanRead(int(payloadLength)) {
		return nil, errors.New("not enough bytes in payload")
	}

	if mapCount > 1 {
		return nil, fmt.Errorf("Unexpected world count of %v, expected just one", mapCount)
	}

	mapWidth, _ := iterator.ReadUInt16()
	mapLength, _ := iterator.ReadUInt16()

	graphicBlockStamp, _ := iterator.ReadString(4)
	if graphicBlockStamp != "GFX0" {
		return nil, fmt.Errorf("Unexpected data file stamp of %v, expected GFX0", graphicBlockStamp)
	}

	layerCount, _ := iterator.ReadByte()
	for i := 0; i < int(layerCount); i++ {
		for z := 0; z < int(mapLength); z++ {
			for x := 0; x < int(mapWidth); x++ {
				iterator.ReadUInt16() // graphic tile id
			}
		}
	}

	collisionBlockStamp, _ := iterator.ReadString(4)
	if collisionBlockStamp != "COLL" {
		return nil, fmt.Errorf("Unexpected data file stamp of %v, expected COLL", collisionBlockStamp)
	}

	collisionMatrix := collision.NewMatrix(int(mapWidth), int(mapLength))
	for z := 0; z < int(mapLength); z++ {
		for x := 0; x < int(mapWidth); x++ {
			flagID, _ := iterator.ReadByte()
			if flagID == 0 {
				// the tile is free
				continue
			}

			var flag collision.BitFlag
			switch flagID {
			case 1:
				flag = collision.Blocked
				break
			case 2:
				flag = collision.Water
				break
			case 3:
				flag = collision.Grass
				break
			case 4:
				flag = collision.Door
				break
			default:
				return nil, fmt.Errorf("Unexpected collision flag id of %v", flagID)
			}

			collisionMatrix.Add(x, z, flag)
		}
	}

	return &TileMap{Index: index, CollisionMatrix: collisionMatrix}, nil
}

// Width returns the width of the TileMap, in tiles.
func (tileMap *TileMap) Width() int {
	return tileMap.CollisionMatrix.Width()
}

// Length returns the length of the TileMap, in tiles.
func (tileMap *TileMap) Length() int {
	return tileMap.CollisionMatrix.Length()
}

// Render builds up a render matrix of the TileMap with all of the
// collision flags, which can then be printed right into the console.
func (tileMap *TileMap) Render() string {
	var bldr strings.Builder
	for z := 0; z < tileMap.Length(); z++ {
		for x := 0; x < tileMap.Width(); x++ {
			flag, _ := tileMap.CollisionMatrix.Get(x, z)
			bldr.WriteString(strconv.Itoa(int(flag)))
		}

		bldr.WriteRune('\n')
	}

	return bldr.String()
}

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

		newMap, err := grid.GetMap(mapX, mapZ)
		if err != nil {
			return position, err
		}

		localX = newMap.Width() - 1
	} else if localX >= mapWidth {
		mapX++
		localX = 0

		_, err := grid.GetMap(mapX, mapZ)
		if err != nil {
			return position, err
		}
	} else if localZ < 0 {
		mapZ--

		newMap, err := grid.GetMap(mapX, mapZ)
		if err != nil {
			return position, err
		}

		localZ = newMap.Length() - 1
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
