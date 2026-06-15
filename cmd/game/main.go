package main

import (
	"math/rand"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

func main() {
	rl.InitWindow(windowWidth, windowHeight, windowTitle)
	defer rl.CloseWindow()

	game := InitializeGame()
	defer game.UnloadAssets()

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()
		rl.ClearBackground(bgColor)

		game.UpdateSystems()

		rl.EndDrawing()
	}
}
