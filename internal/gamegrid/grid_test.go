package gamegrid

import (
	"testing"

	"github.com/mlange-42/ark/ecs"
)

const (
	testGridWidth       = 13
	testGridLength      = 11
	testGridBorderWidth = 2
	testGridCenterX     = testGridWidth / 2
	testGridCenterZ     = testGridLength / 2
)

func TestGameGridBorderBuildability(t *testing.T) {
	grid := newTestGrid()

	for z := 0; z < grid.Length; z++ {
		for x := 0; x < grid.Width; x++ {
			cell, ok := grid.Cell(x, z)
			if !ok {
				t.Fatalf("expected cell (%d, %d) to exist", x, z)
			}

			wantBuildable := x >= testGridBorderWidth && x < grid.Width-testGridBorderWidth &&
				z >= testGridBorderWidth && z < grid.Length-testGridBorderWidth
			if cell.Buildable() != wantBuildable {
				t.Fatalf("cell (%d, %d) buildable = %v, want %v", x, z, cell.Buildable(), wantBuildable)
			}
		}
	}
}

func TestGameGridBFSRespectsOccupiedCells(t *testing.T) {
	grid := newTestGrid()

	if got := grid.Distance(testGridCenterX+2, testGridCenterZ); got != 2 {
		t.Fatalf("distance before occupying cell = %d, want 2", got)
	}
	occupyCell(t, &grid, testGridCenterX+1, testGridCenterZ)

	if got := grid.Distance(testGridCenterX+1, testGridCenterZ); got != -1 {
		t.Fatalf("occupied cell distance = %d, want -1", got)
	}
	if got := grid.Distance(testGridCenterX+2, testGridCenterZ); got != 4 {
		t.Fatalf("distance around occupied cell = %d, want 4", got)
	}
}

func TestGameGridNextLowerDistanceCell(t *testing.T) {
	grid := newTestGrid()

	nextX, nextZ, ok := grid.NextLowerDistanceCell(testGridCenterX+2, testGridCenterZ)
	if !ok {
		t.Fatal("expected a next lower distance cell")
	}
	if nextX != testGridCenterX+1 || nextZ != testGridCenterZ {
		t.Fatalf("next lower cell = (%d, %d), want (%d, %d)", nextX, nextZ, testGridCenterX+1, testGridCenterZ)
	}

	occupyCell(t, &grid, testGridCenterX+1, testGridCenterZ)
	if _, _, ok := grid.NextLowerDistanceCell(testGridCenterX+1, testGridCenterZ); ok {
		t.Fatal("expected no next lower cell for an occupied/unreachable cell")
	}
}

func TestGameGridRejectsPlacementThatBlocksPathOrigins(t *testing.T) {
	grid := newTestGrid()

	occupyCell(t, &grid, testGridCenterX, testGridCenterZ-1)
	occupyCell(t, &grid, testGridCenterX, testGridCenterZ+1)
	occupyCell(t, &grid, testGridCenterX-1, testGridCenterZ)

	if grid.SetCellEntity(testGridCenterX+1, testGridCenterZ, ecs.Entity{}) {
		t.Fatal("expected placement that blocks all path origins to be rejected")
	}

	cell, ok := grid.Cell(testGridCenterX+1, testGridCenterZ)
	if !ok {
		t.Fatal("expected rejected placement cell to exist")
	}
	if cell.HasEntity() {
		t.Fatal("rejected placement left the cell occupied")
	}
	if got := grid.Distance(testGridCenterX+1, testGridCenterZ); got != 1 {
		t.Fatalf("distance after rejected placement = %d, want 1", got)
	}
	for _, position := range testPathOrigins() {
		if got := grid.Distance(position.X, position.Z); got < 0 {
			t.Fatalf("path origin (%d, %d) distance = %d, want reachable", position.X, position.Z, got)
		}
	}
}

func newTestGrid() GameGrid {
	return NewGameGrid(
		testGridWidth,
		testGridLength,
		testGridBorderWidth,
		GridCoord{X: testGridCenterX, Z: testGridCenterZ},
		testPathOrigins(),
	)
}

func testPathOrigins() []GridCoord {
	return []GridCoord{
		{X: testGridCenterX, Z: 0},
		{X: testGridCenterX, Z: testGridLength - 1},
		{X: 0, Z: testGridCenterZ},
		{X: testGridWidth - 1, Z: testGridCenterZ},
	}
}

func occupyCell(t *testing.T, grid *GameGrid, x, z int) {
	t.Helper()

	if !grid.SetCellEntity(x, z, ecs.Entity{}) {
		t.Fatalf("expected to occupy cell (%d, %d)", x, z)
	}
}
