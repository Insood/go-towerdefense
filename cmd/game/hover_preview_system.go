package main

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
	ecs "github.com/mlange-42/ark/ecs"
)

type HoverPreviewSystem struct {
	previewMap    *ecs.Map1[HoverPreview]
	previewEntity ecs.Entity
}

func (system *HoverPreviewSystem) Initialize(game *Game) {
	system.previewMap = ecs.NewMap1[HoverPreview](game.world)
	system.previewEntity = system.previewMap.NewEntity(&HoverPreview{})
}

func (system *HoverPreviewSystem) Update(game *Game) {
	ray := rl.GetScreenToWorldRay(rl.GetMousePosition(), game.camera)
	preview := system.previewMap.Get(system.previewEntity)

	if point, ok := intersectRayGroundPlane(ray); ok {
		preview.gridX = int(math.Floor(float64(point.X)))
		preview.gridZ = int(math.Floor(float64(point.Z)))
		preview.visible = true
		return
	}

	preview.visible = false
}
