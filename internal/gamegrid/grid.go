package gamegrid

import "github.com/mlange-42/ark/ecs"

type GridCoord struct {
	X int
	Z int
}

type GameGrid struct {
	Width       int
	Length      int
	center      GridCoord
	borderWidth int
	pathOrigins []GridCoord
	cells       []GameGridCell
}

type GameGridCell struct {
	entity    ecs.Entity
	distance  int
	occupied  bool
	buildable bool
}

func NewGameGrid(width, length, borderWidth int, center GridCoord, pathOrigins []GridCoord) GameGrid {
	grid := GameGrid{
		Width:       width,
		Length:      length,
		center:      center,
		borderWidth: borderWidth,
		pathOrigins: append([]GridCoord(nil), pathOrigins...),
		cells:       make([]GameGridCell, width*length),
	}
	grid.initializeBuildableCells()
	grid.RecalculateDistances()
	return grid
}

func (grid *GameGrid) Cell(x, z int) (*GameGridCell, bool) {
	index, ok := grid.index(x, z)
	if !ok {
		return nil, false
	}

	return &grid.cells[index], true
}

func (grid *GameGrid) CellEntity(x, z int) (ecs.Entity, bool) {
	cell, ok := grid.Cell(x, z)
	if !ok {
		return ecs.Entity{}, false
	}

	return cell.Entity()
}

func (grid *GameGrid) CanSetCellEntity(x, z int) bool {
	cell, ok := grid.Cell(x, z)
	if !ok || !cell.Buildable() || cell.HasEntity() {
		return false
	}

	return grid.canOccupyCellWithoutBlockingPathOrigins(x, z)
}

func (grid *GameGrid) SetCellEntity(x, z int, entity ecs.Entity) bool {
	if !grid.CanSetCellEntity(x, z) {
		return false
	}

	cell, _ := grid.Cell(x, z)
	cell.SetEntity(entity)
	grid.RecalculateDistances()
	return true
}

func (grid *GameGrid) SetCellEntityForce(x, z int, entity ecs.Entity) bool {
	cell, ok := grid.Cell(x, z)
	if !ok || cell.HasEntity() {
		return false
	}

	cell.SetEntity(entity)
	cell.buildable = false
	grid.RecalculateDistances()
	return true
}

func (grid *GameGrid) SetCellBuildable(x, z int, buildable bool) bool {
	cell, ok := grid.Cell(x, z)
	if !ok {
		return false
	}

	cell.buildable = buildable
	return true
}

func (grid *GameGrid) Distance(x, z int) int {
	cell, ok := grid.Cell(x, z)
	if !ok {
		return -1
	}

	return cell.distance
}

func (grid *GameGrid) NextLowerDistanceCells(x, z int) []GridCoord {
	cell, ok := grid.Cell(x, z)
	if !ok || cell.distance < 0 {
		return nil
	}

	bestDistance := cell.distance
	candidates := make([]GridCoord, 0, 4)

	for _, delta := range [][2]int{{0, -1}, {0, 1}, {-1, 0}, {1, 0}} {
		nextX := x + delta[0]
		nextZ := z + delta[1]
		nextCell, ok := grid.Cell(nextX, nextZ)
		if !ok || nextCell.distance < 0 {
			continue
		}
		if nextCell.distance > bestDistance {
			continue
		}

		if nextCell.distance < bestDistance {
			bestDistance = nextCell.distance
			candidates = candidates[:0]
		}

		candidates = append(candidates, GridCoord{X: nextX, Z: nextZ})
	}

	return candidates
}

func (grid *GameGrid) PathToCenter(x, z int) []GridCoord {
	cell, ok := grid.Cell(x, z)
	if !ok || cell.distance < 0 {
		return nil
	}

	path := make([]GridCoord, 0, cell.distance+1)
	path = append(path, GridCoord{X: x, Z: z})
	currentX, currentZ := x, z

	for {
		candidates := grid.NextLowerDistanceCells(currentX, currentZ)
		if len(candidates) == 0 {
			break
		}

		next := candidates[0]
		path = append(path, next)
		if next.X == grid.center.X && next.Z == grid.center.Z {
			break
		}

		currentX = next.X
		currentZ = next.Z
	}

	return path
}

func (cell *GameGridCell) SetEntity(entity ecs.Entity) {
	cell.entity = entity
	cell.occupied = true
}

func (cell *GameGridCell) ClearEntity() {
	cell.entity = ecs.Entity{}
	cell.occupied = false
}

func (cell *GameGridCell) Entity() (ecs.Entity, bool) {
	if !cell.occupied {
		return ecs.Entity{}, false
	}

	return cell.entity, true
}

func (cell *GameGridCell) HasEntity() bool {
	return cell.occupied
}

func (cell *GameGridCell) Buildable() bool {
	return cell.buildable
}

func (cell *GameGridCell) Distance() int {
	return cell.distance
}

func (grid *GameGrid) index(x, z int) (int, bool) {
	if x < 0 || x >= grid.Width || z < 0 || z >= grid.Length {
		return 0, false
	}

	return z*grid.Width + x, true
}

func (grid *GameGrid) canOccupyCellWithoutBlockingPathOrigins(x, z int) bool {
	cell, ok := grid.Cell(x, z)
	if !ok || cell.HasEntity() {
		return false
	}

	cell.SetEntity(ecs.Entity{})
	grid.RecalculateDistances()
	pathsOpen := grid.pathOriginsReachable()
	cell.ClearEntity()
	grid.RecalculateDistances()

	return pathsOpen
}

func (grid *GameGrid) pathOriginsReachable() bool {
	for _, position := range grid.pathOrigins {
		if grid.Distance(position.X, position.Z) < 0 {
			return false
		}
	}

	return true
}

func (grid *GameGrid) initializeBuildableCells() {
	for z := 0; z < grid.Length; z++ {
		for x := 0; x < grid.Width; x++ {
			cell := &grid.cells[z*grid.Width+x]
			cell.buildable = x >= grid.borderWidth && x < grid.Width-grid.borderWidth &&
				z >= grid.borderWidth && z < grid.Length-grid.borderWidth
		}
	}
}

func (grid *GameGrid) RecalculateDistances() {
	for i := range grid.cells {
		grid.cells[i].distance = -1
	}

	startIndex, ok := grid.index(grid.center.X, grid.center.Z)
	if !ok {
		return
	}

	queue := make([]GridCoord, 0, len(grid.cells))
	queue = append(queue, grid.center)
	grid.cells[startIndex].distance = 0

	for head := 0; head < len(queue); head++ {
		current := queue[head]
		currentIndex, _ := grid.index(current.X, current.Z)
		currentDistance := grid.cells[currentIndex].distance

		for _, delta := range [][2]int{{0, -1}, {0, 1}, {-1, 0}, {1, 0}} {
			nextX := current.X + delta[0]
			nextZ := current.Z + delta[1]
			nextIndex, ok := grid.index(nextX, nextZ)
			if !ok {
				continue
			}

			nextCell := &grid.cells[nextIndex]
			if nextCell.distance >= 0 {
				continue
			}
			if nextCell.occupied && !(nextX == grid.center.X && nextZ == grid.center.Z) {
				continue
			}

			nextCell.distance = currentDistance + 1
			if !nextCell.occupied {
				queue = append(queue, GridCoord{X: nextX, Z: nextZ})
			}
		}
	}
}
