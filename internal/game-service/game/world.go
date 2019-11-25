package game

// Position describes an exact tile position on the game map.
type Position struct {
	X int
	Z int
	Y int
}

// A map of tiles in the game world for entities to traverse.
type TileMap struct {
	Index           MapIndex
	CollisionMatrix *CollisionMatrix
}

// Grid represents the physical grid of TileMap's that together
// make up the game world.
type Grid struct {
	TileMaps [][]*TileMap
}

// Direction is a compass direction.
type Direction int

// These are the supported types of Direction's.
const (
	South Direction = 0
	North Direction = 1
	West  Direction = 2
	East  Direction = 3
)
