package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	rl.InitWindow(1280, 720, "Go Tower Defense")
	defer rl.CloseWindow()

	rl.SetTargetFPS(60)

	camera := rl.Camera3D{
		Position:   rl.NewVector3(8, 8, 8),
		Target:     rl.NewVector3(0, 1, 0),
		Up:         rl.NewVector3(0, 1, 0),
		Fovy:       45,
		Projection: rl.CameraPerspective,
	}

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()
		rl.ClearBackground(rl.NewColor(20, 24, 32, 255))

		rl.BeginMode3D(camera)
		rl.DrawGrid(20, 1)
		rl.DrawCube(rl.NewVector3(0, 1, 0), 1, 2, 1, rl.NewColor(198, 120, 76, 255))
		rl.DrawCubeWires(rl.NewVector3(0, 1, 0), 1, 2, 1, rl.Black)
		rl.DrawSphere(rl.NewVector3(2, 1.2, 2), 0.35, rl.NewColor(90, 160, 255, 255))
		rl.EndMode3D()

		rl.EndDrawing()
	}
}
