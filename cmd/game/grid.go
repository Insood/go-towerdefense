package main

import (
	"fmt"
	"image/color"

	gamegrid "go-towerdefense/internal/gamegrid"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/mlange-42/ark/ecs"
)

type GameGrid struct {
	gamegrid.GameGrid
	world        *ecs.World
	entityMapper *ecs.Map2[Position3, Renderable]
}

func NewGameGrid(width, length int) GameGrid {
	return GameGrid{
		GameGrid: gamegrid.NewGameGrid(
			width,
			length,
			gridBorderWidth,
			gamegrid.GridCoord{X: gridCenterX, Z: gridCenterZ},
			spawnerGridPositions(),
		),
	}
}

func (grid *GameGrid) Initialize(world *ecs.World) {
	grid.world = world
	grid.entityMapper = ecs.NewMap2[Position3, Renderable](world)
}

func (grid *GameGrid) PlaceEntity(x, z int, model *rl.Model, tint color.RGBA) bool {
	return grid.placeEntity(x, z, groundPlaneY+0.5, model, tint, false)
}

func (grid *GameGrid) ForcePlaceEntity(x, z int, y float32, model *rl.Model, tint color.RGBA) bool {
	return grid.placeEntity(x, z, y, model, tint, true)
}

func (grid *GameGrid) placeEntity(x, z int, y float32, model *rl.Model, tint color.RGBA, ignoreBuildable bool) bool {
	if grid.entityMapper == nil {
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
	if !ignoreBuildable && !grid.CanSetCellEntity(x, z) {
		fmt.Printf("entity placement blocked: path blocked (%d, %d)\n", x, z)
		return false
	}

	entity := grid.entityMapper.NewEntity(
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

	if ignoreBuildable {
		if !grid.SetCellEntityForce(x, z, entity) {
			return false
		}
	} else if !grid.SetCellEntity(x, z, entity) {
		return false
	}

	fmt.Printf("entity placed at grid (%d, %d)\n", x, z)
	return true
}
