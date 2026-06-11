package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	rl.InitWindow(windowWidth, windowHeight, windowTitle)
	defer rl.CloseWindow()

	game := InitializeGame()
	defer game.UnloadShaders()

	for !rl.WindowShouldClose() {
		game.cameraSystem.Update(game)

		rl.BeginDrawing()
		rl.ClearBackground(bgColor)

		game.UpdateSystems()

		rl.EndDrawing()
	}
}
