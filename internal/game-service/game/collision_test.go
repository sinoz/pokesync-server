package game

import (
	"testing"
)

func TestCollisionMatrix_Add(t *testing.T) {
	matrix := NewCollisionMatrix(16, 16)
	_ = matrix.Add(2, 2, Blocked)
	if matrix.Flags[2][2] != Blocked {
		t.Errorf("expected collision flag at coordinates %v %v to equal %v", 2, 2, Blocked)
	}
}

func TestCollisionMatrix_Add_Stacking(t *testing.T) {
	matrix := NewCollisionMatrix(16, 16)

	_ = matrix.Add(7, 5, Water)
	_ = matrix.Add(7, 5, Door)

	isBlocked, _ := matrix.Contains(7, 5, Blocked)
	if isBlocked {
		t.Errorf("expected collision flag at coordinates %v %v to equal %v", 7, 5, 0)
	}

	hasWater, _ := matrix.Contains(7, 5, Water)
	if !hasWater {
		t.Errorf("expected collision flag at coordinates %v %v to equal %v", 7, 5, Water)
	}

	hasDoor, _ := matrix.Contains(7, 5, Door)
	if !hasDoor {
		t.Errorf("expected collision flag at coordinates %v %v to equal %v", 7, 5, Door)
	}
}

func TestCollisionMatrix_Clear(t *testing.T) {
	matrix := NewCollisionMatrix(16, 16)

	_ = matrix.Add(5, 7, Water)
	_ = matrix.Add(3, 2, Blocked)
	_ = matrix.Clear(3, 2)

	if matrix.Flags[3][2] != 0 {
		t.Errorf("expected collision flag at coordinates %v %v to equal %v", 3, 2, 0)
	}
}

func TestCollisionMatrix_Remove(t *testing.T) {
	matrix := NewCollisionMatrix(16, 16)

	_ = matrix.Add(7, 5, Water)
	_ = matrix.Add(7, 5, Blocked)
	_ = matrix.Remove(7, 5, Blocked)

	isBlocked, _ := matrix.Contains(7, 5, Blocked)
	if isBlocked {
		t.Errorf("expected collision flag at coordinates %v %v to equal %v", 7, 5, 0)
	}
}
