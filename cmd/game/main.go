package main

import (
	"fmt"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	rl.InitWindow(windowWidth, windowHeight, windowTitle)
	defer rl.CloseWindow()

	game := InitializeGame()
	defer game.UnloadShaders()

	for !rl.WindowShouldClose() {
		game.cameraSystem.Update(game)

		if rl.IsMouseButtonPressed(rl.MouseButtonLeft) {
			ray := rl.GetScreenToWorldRay(rl.GetMousePosition(), game.camera)
			if point, ok := intersectRayGroundPlane(ray); ok {
				gridX := int(math.Floor(float64(point.X)))
				gridZ := int(math.Floor(float64(point.Z)))

				fmt.Printf(
					"click world=(%.2f, %.2f, %.2f) grid=(%d, %d)\n",
					point.X,
					point.Y,
					point.Z,
					gridX,
					gridZ,
				)
				game.grid.PlaceEntity(gridX, gridZ, game.models["cube"], baseCubeColor)
			} else {
				fmt.Println("click missed the ground plane")
			}
		}

		rl.BeginDrawing()
		rl.ClearBackground(bgColor)

		game.UpdateSystems()

		rl.EndDrawing()
	}
}
