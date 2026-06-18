package main

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
	ecs "github.com/mlange-42/ark/ecs"
)

type HoverPreviewSystem struct {
	previewMap    *ecs.Map2[HoverPreview, Position3]
	renderableMap *ecs.Map1[Renderable]
	previewEntity ecs.Entity
}

func (system *HoverPreviewSystem) Initialize(game *Game) {
	system.previewMap = ecs.NewMap2[HoverPreview, Position3](game.world)
	system.renderableMap = ecs.NewMap1[Renderable](game.world)
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
			if system.renderableMap.HasAll(system.previewEntity) {
				renderable := system.renderableMap.Get(system.previewEntity)
				renderable.tint = hoverPreviewTint
			} else {
				system.renderableMap.Add(
					system.previewEntity,
					&Renderable{
						model:             game.assets.Model("turret"),
						scale:             1.0,
						tint:              hoverPreviewTint,
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
