package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	ecs "github.com/mlange-42/ark/ecs"
)

type RenderSystem3D struct {
	filter        *ecs.Filter2[Position3, Renderable]
	previewFilter *ecs.Filter1[HoverPreview]
}

func (system *RenderSystem3D) Initialize(game *Game) {
	system.filter = ecs.NewFilter2[Position3, Renderable](game.world)
	system.previewFilter = ecs.NewFilter1[HoverPreview](game.world)
}

func (system *RenderSystem3D) Update(game *Game) {
	rl.BeginMode3D(game.camera)
	system.drawCoordinateSystem()
	system.renderModels()
	system.renderHoverPreview(game)
	rl.EndMode3D()
}

func (system *RenderSystem3D) drawCoordinateSystem() {
	origin := rl.Vector3Zero()
	rl.DrawLine3D(origin, rl.NewVector3(axisLength, 0, 0), rl.Red)
	rl.DrawLine3D(origin, rl.NewVector3(0, axisLength, 0), rl.Green)
	rl.DrawLine3D(origin, rl.NewVector3(0, 0, axisLength), rl.Blue)
}

func (system *RenderSystem3D) renderModels() {
	query := system.filter.Query()
	defer query.Close()

	for query.Next() {
		position, renderable := query.Get()

		drawTint := renderable.tint
		if renderable.shaderTintEnabled {
			materials := renderable.model.GetMaterials()
			if len(materials) > 0 {
				shader := materials[0].Shader
				location := rl.GetShaderLocation(shader, "tintColor")
				if location >= 0 {
					tint := colorToVec4(renderable.shaderTint)
					rl.SetShaderValue(
						shader,
						location,
						tint[:],
						rl.ShaderUniformVec4,
					)
				}
			}
			drawTint = rl.White
		}

		rl.DrawModel(*renderable.model, *position, renderable.scale, drawTint)
	}
}

func (system *RenderSystem3D) renderHoverPreview(game *Game) {
	query := system.previewFilter.Query()
	defer query.Close()

	if !query.Next() {
		return
	}

	preview := query.Get()
	if !preview.visible {
		return
	}

	cell, ok := game.grid.Cell(preview.gridX, preview.gridZ)
	if !ok || !cell.Buildable() || cell.HasEntity() {
		return
	}

	rl.DrawModel(
		*game.models["cube"],
		rl.NewVector3(float32(preview.gridX)+gridCellCenter, groundPlaneY+0.5, float32(preview.gridZ)+gridCellCenter),
		1.0,
		hoverPreviewTint,
	)
}
