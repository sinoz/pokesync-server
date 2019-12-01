package collision

import "errors"

// These are the supported types of BitFlag's.
const (
	Blocked BitFlag = 1
	Water   BitFlag = 2
	Grass   BitFlag = 4
	Door    BitFlag = 8
)

// BitFlag is a collision flag with a bit-mask value.
type BitFlag int

// Matrix is the grid of CollisionBitFlag's.
type Matrix struct {
	Flags [][]BitFlag
}

// NewMatrix constructs a Matrix of the specified width and length.
func NewMatrix(width, length int) *Matrix {
	flags := make([][]BitFlag, width)
	for x := range flags {
		flags[x] = make([]BitFlag, length)
	}

	return &Matrix{Flags: flags}
}

// Add adds the given BitFlag to the specified tile. May return an
// error if the given coordinates were out of bounds of the matrix.
func (m *Matrix) Add(x, z int, flag BitFlag) error {
	if err := m.outOfBounds(x, z); err != nil {
		return err
	}

	m.Flags[x][z] |= flag
	return nil
}

// Remove removes the given BitFlag from the specified tile. May
// return an error if the given coordinates were out of bounds of the matrix.
func (m *Matrix) Remove(x, z int, flag BitFlag) error {
	if err := m.outOfBounds(x, z); err != nil {
		return err
	}

	m.Flags[x][z] &= ^flag
	return nil
}

// Contains returns whether the given BitFlag is set at the specified
// coordinates. May return an error if the given coordinates were out
// of bounds of the matrix.
func (m Matrix) Contains(x, z int, flag BitFlag) (bool, error) {
	if err := m.outOfBounds(x, z); err != nil {
		return false, err
	}

	return m.Flags[x][z]&flag == flag, nil
}

// Get looks up the BitFlag's that are set for the specified tile. May return
// an error if the given coordinates were out of bounds of the matrix.
func (m Matrix) Get(x, z int) (BitFlag, error) {
	if err := m.outOfBounds(x, z); err != nil {
		return 0, err
	}

	return m.Flags[x][z], nil
}

// Clear clears the specified tile of all of its collision flags. May return
// an error if the given coordinates were out of bounds of the matrix.
func (m *Matrix) Clear(x, z int) error {
	if err := m.outOfBounds(x, z); err != nil {
		return err
	}

	m.Flags[x][z] = 0
	return nil
}

// Width returns the width of the matrix, in tiles.
func (m Matrix) Width() int {
	return len(m.Flags)
}

// Length returns the length of the matrix, in tiles.
func (m Matrix) Length() int {
	return len(m.Flags[0])
}

// outOfBounds returns whether the given coordinates fall out of bounds
// of this matrix. Returns nil if the given coordinates are ok.
func (m Matrix) outOfBounds(x, z int) error {
	if x < 0 || z < 0 || x >= m.Width() || z >= m.Length() {
		return errors.New("given coordinates are out of bounds")
	}

	return nil
}
