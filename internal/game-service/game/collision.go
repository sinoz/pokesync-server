package game

import "errors"

// These are the supported types of CollisionBitFlag's.
const (
	Blocked CollisionBitFlag = 1
	Water   CollisionBitFlag = 2
	Grass   CollisionBitFlag = 4
	Door    CollisionBitFlag = 8
)

// CollisionBitFlag is a collision flag with a bit-mask value.
type CollisionBitFlag int

// CollisionMatrix is the grid of CollisionBitFlag's.
type CollisionMatrix struct {
	Flags [][]CollisionBitFlag
}

// NewCollisionMatrix constructs a CollisionMatrix of the specified
// width and length.
func NewCollisionMatrix(width, length int) *CollisionMatrix {
	flags := make([][]CollisionBitFlag, width)
	for i := range flags {
		flags[i] = make([]CollisionBitFlag, length)
	}

	return &CollisionMatrix{Flags: flags}
}

// Add adds the given CollisionBitFlag to the specified tile. May return an
// error if the given coordinates were out of bounds of the matrix.
func (m *CollisionMatrix) Add(x, z int, flag CollisionBitFlag) error {
	if err := m.outOfBounds(x, z); err != nil {
		return err
	}

	m.Flags[x][z] |= flag
	return nil
}

// Remove removes the given CollisionBitFlag from the specified tile. May
// return an error if the given coordinates were out of bounds of the matrix.
func (m *CollisionMatrix) Remove(x, z int, flag CollisionBitFlag) error {
	if err := m.outOfBounds(x, z); err != nil {
		return err
	}

	m.Flags[x][z] &= ^flag
	return nil
}

// Contains returns whether the given CollisionBitFlag is set at the specified
// coordinates. May return an error if the given coordinates were out of bounds
// of the matrix.
func (m CollisionMatrix) Contains(x, z int, flag CollisionBitFlag) (bool, error) {
	if err := m.outOfBounds(x, z); err != nil {
		return false, err
	}

	return m.Flags[x][z]&flag == flag, nil
}

// Get looks up the CollisionBitFlag's that are set for the specified tile.
// May return an error if the given coordinates were out of bounds of the
// matrix.
func (m CollisionMatrix) Get(x, z int) (CollisionBitFlag, error) {
	if err := m.outOfBounds(x, z); err != nil {
		return 0, err
	}

	return m.Flags[x][z], nil
}

// Clear clears the specified tile of all of its collision flags. May return
// an error if the given coordinates were out of bounds of the matrix.
func (m *CollisionMatrix) Clear(x, z int) error {
	if err := m.outOfBounds(x, z); err != nil {
		return err
	}

	m.Flags[x][z] = 0
	return nil
}

// Width returns the width of the matrix, in tiles.
func (m CollisionMatrix) Width() int {
	return len(m.Flags)
}

// Length returns the length of the matrix, in tiles.
func (m CollisionMatrix) Length() int {
	return len(m.Flags[0])
}

// outOfBounds returns whether the given coordinates fall out of bounds
// of this matrix. Returns nil if the given coordinates are ok.
func (m CollisionMatrix) outOfBounds(x, z int) error {
	if x < 0 || z < 0 || x >= m.Width() || z >= m.Length() {
		return errors.New("given coordinates are out of bounds")
	}

	return nil
}
