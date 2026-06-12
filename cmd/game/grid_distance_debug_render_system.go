package main

import (
	"strconv"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type GridDistanceDebugRenderSystem struct{}

func (system *GridDistanceDebugRenderSystem) Initialize(game *Game) {}

func (system *GridDistanceDebugRenderSystem) Update(game *Game) {
	if !debugShowGridDistances {
		return
	}

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
			label := strconv.Itoa(cell.distance)
			labelWidth := rl.MeasureText(label, gridDistanceLabelSize)

			drawX := int32(screenPosition.X) - labelWidth/2
			drawY := int32(screenPosition.Y) - gridDistanceLabelOffset

			rl.DrawText(label, drawX+1, drawY+1, gridDistanceLabelSize, rl.Black)
			rl.DrawText(label, drawX, drawY, gridDistanceLabelSize, rl.White)
		}
	}
}
