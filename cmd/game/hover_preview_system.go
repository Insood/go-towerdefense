package main

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
	ecs "github.com/mlange-42/ark/ecs"
)

type HoverPreviewSystem struct {
	previewMap    *ecs.Map2[HoverPreview, Position3]
	renderableMap *ecs.Map1[Renderable]
	enemyFilter   *ecs.Filter2[Position3, Enemy]
	previewEntity ecs.Entity
}

func (system *HoverPreviewSystem) Initialize(game *Game) {
	system.previewMap = ecs.NewMap2[HoverPreview, Position3](game.world)
	system.renderableMap = ecs.NewMap1[Renderable](game.world)
	system.enemyFilter = ecs.NewFilter2[Position3, Enemy](game.world)
	system.previewEntity = system.previewMap.NewEntity(
		&HoverPreview{},
		&Position3{},
	)
}

func (system *HoverPreviewSystem) Update(game *Game) {
	ray := rl.GetScreenToWorldRay(rl.GetMousePosition(), game.camera)
	preview, position := system.previewMap.Get(system.previewEntity)

	if point, ok := intersectRayGroundPlane(ray); ok {
		preview.gridX = int(math.Floor(float64(point.X)))
		preview.gridZ = int(math.Floor(float64(point.Z)))
		position.X = float32(preview.gridX) + gridCellCenter
		position.Y = groundPlaneY
		position.Z = float32(preview.gridZ) + gridCellCenter

		cell, ok := game.grid.Cell(preview.gridX, preview.gridZ)
		if ok && cell.Buildable() && !cell.HasEntity() {
			tint := hoverPreviewTintAllowed
			if system.gridContainsEnemy(preview.gridX, preview.gridZ) {
				tint = hoverPreviewTintNotAllowed
			}

			if system.renderableMap.HasAll(system.previewEntity) {
				renderable := system.renderableMap.Get(system.previewEntity)
				renderable.tint = tint
			} else {
				system.renderableMap.Add(
					system.previewEntity,
					&Renderable{
						model:             game.assets.Model("turret"),
						scale:             1.0,
						tint:              tint,
						shaderTintEnabled: false,
					},
				)
			}
			return
		}

		if system.renderableMap.HasAll(system.previewEntity) {
			system.renderableMap.Remove(system.previewEntity)
		}
		return
	}

	if system.renderableMap.HasAll(system.previewEntity) {
		system.renderableMap.Remove(system.previewEntity)
	}
}

func (system *HoverPreviewSystem) gridContainsEnemy(gridX, gridZ int) bool {
	query := system.enemyFilter.Query()
	defer query.Close()

	for query.Next() {
		position, _ := query.Get()
		cellX := int(math.Floor(float64(position.X)))
		cellZ := int(math.Floor(float64(position.Z)))
		if cellX == gridX && cellZ == gridZ {
			return true
		}
	}

	return false
}
