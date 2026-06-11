package main

import (
	"fmt"
	"image/color"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/mlange-42/ark/ecs"
)

type GameGrid struct {
	Width      int
	Length     int
	world      *ecs.World
	cubeMapper *ecs.Map2[Position3, Renderable]
	cells      []GameGridCell
}

type GameGridCell struct {
	entity    ecs.Entity
	distance  int
	occupied  bool
	buildable bool
}

func NewGameGrid(width, length int) GameGrid {
	grid := GameGrid{
		Width:  width,
		Length: length,
		cells:  make([]GameGridCell, width*length),
	}
	grid.initializeBuildableCells()
	grid.initializeDistances()
	return grid
}

func (grid *GameGrid) Initialize(world *ecs.World) {
	grid.world = world
	grid.cubeMapper = ecs.NewMap2[Position3, Renderable](world)
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

func (grid *GameGrid) SetCellEntity(x, z int, entity ecs.Entity) bool {
	cell, ok := grid.Cell(x, z)
	if !ok || !cell.Buildable() || cell.HasEntity() {
		return false
	}

	cell.SetEntity(entity)
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

func (grid *GameGrid) PlaceEntity(x, z int, model *rl.Model, tint color.RGBA) bool {
	return grid.placeEntity(x, z, groundPlaneY+0.5, model, tint, false)
}

func (grid *GameGrid) ForcePlaceEntity(x, z int, y float32, model *rl.Model, tint color.RGBA) bool {
	return grid.placeEntity(x, z, y, model, tint, true)
}

func (grid *GameGrid) placeEntity(x, z int, y float32, model *rl.Model, tint color.RGBA, ignoreBuildable bool) bool {
	if grid.cubeMapper == nil {
		panic("game grid is not initialized")
	}

	cell, ok := grid.Cell(x, z)
	if !ok {
		fmt.Printf("entity placement blocked: out of bounds (%d, %d)\n", x, z)
		return false
	}
	if !ignoreBuildable && !cell.Buildable() {
		fmt.Printf("entity placement blocked: no-build zone (%d, %d)\n", x, z)
		return false
	}
	if cell.HasEntity() {
		fmt.Printf("entity placement blocked: occupied cell (%d, %d)\n", x, z)
		return false
	}

	entity := grid.cubeMapper.NewEntity(
		&Position3{
			X: float32(x) + 0.5,
			Y: y,
			Z: float32(z) + 0.5,
		},
		&Renderable{
			model:             model,
			scale:             1.0,
			tint:              tint,
			shaderTint:        rl.White,
			shaderTintEnabled: false,
		},
	)
	cell.SetEntity(entity)
	if ignoreBuildable {
		cell.buildable = false
	}

	fmt.Printf("entity placed at grid (%d, %d)\n", x, z)
	return true
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

func (grid *GameGrid) index(x, z int) (int, bool) {
	if x < 0 || x >= grid.Width || z < 0 || z >= grid.Length {
		return 0, false
	}

	return z*grid.Width + x, true
}

func (grid *GameGrid) initializeBuildableCells() {
	for z := 0; z < grid.Length; z++ {
		for x := 0; x < grid.Width; x++ {
			cell := &grid.cells[z*grid.Width+x]
			cell.buildable = x >= gridBorderWidth && x < grid.Width-gridBorderWidth &&
				z >= gridBorderWidth && z < grid.Length-gridBorderWidth
		}
	}
}

func (grid *GameGrid) initializeDistances() {
	for z := 0; z < grid.Length; z++ {
		for x := 0; x < grid.Width; x++ {
			cell := &grid.cells[z*grid.Width+x]
			cell.distance = manhattanDistance(x, z, gridCenterX, gridCenterZ)
		}
	}
}
