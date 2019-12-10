package game

// Route is the generated series of directional steps to take.
type Route []Direction

// RouteFinder calculates a path between the two given Position's.
type RouteFinder func(grid *Grid, source, dest Position) (Route, error)

// AStarRouteFinder TODO
func AStarRouteFinder() RouteFinder {
	return func(grid *Grid, source, dest Position) (Route, error) {
		return nil, nil
	}
}
