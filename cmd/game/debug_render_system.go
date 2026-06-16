package main

import (
	"strconv"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type DebugRenderSystem struct{}

func (system *DebugRenderSystem) Initialize(game *Game) {}

func (system *DebugRenderSystem) Update(game *Game) {
	if !debugShowGridDistances {
		return
	}

	rl.DrawFPS(10, 10)
	system.drawGridDistances(game)
	system.drawCoordinateSystem()
}

func (system *DebugRenderSystem) drawCoordinateSystem() {
	origin := rl.Vector3Zero()
	rl.DrawLine3D(origin, rl.NewVector3(axisLength, 0, 0), rl.Red)
	rl.DrawLine3D(origin, rl.NewVector3(0, axisLength, 0), rl.Green)
	rl.DrawLine3D(origin, rl.NewVector3(0, 0, axisLength), rl.Blue)
}

func (system *DebugRenderSystem) drawGridDistances(game *Game) {
	for z := 0; z < game.grid.Length; z++ {
		for x := 0; x < game.grid.Width; x++ {
			cell, ok := game.grid.Cell(x, z)
			if !ok {
				continue
			}

			worldPosition := rl.NewVector3(
				float32(x)+gridCellCenter,
				gridDistanceLabelY,
				float32(z)+gridCellCenter,
			)
			screenPosition := rl.GetWorldToScreen(worldPosition, game.camera)
			label := strconv.Itoa(cell.Distance())
			labelWidth := rl.MeasureText(label, gridDistanceLabelSize)

			drawX := int32(screenPosition.X) - labelWidth/2
			drawY := int32(screenPosition.Y) - gridDistanceLabelOffset

			rl.DrawText(label, drawX+1, drawY+1, gridDistanceLabelSize, rl.Black)
			rl.DrawText(label, drawX, drawY, gridDistanceLabelSize, rl.White)
		}
	}
}
