package main

import (
	"fmt"
	"image/color"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Position3 = rl.Vector3
type Renderable struct {
	model *rl.Model
	scale float32
	tint  color.RGBA
}

func main() {
	rl.InitWindow(windowWidth, windowHeight, windowTitle)
	defer rl.CloseWindow()

	game := InitializeGame()

	for !rl.WindowShouldClose() {
		game.cameraSystem.Update(game)

		if rl.IsMouseButtonPressed(rl.MouseButtonLeft) {
			ray := rl.GetScreenToWorldRay(rl.GetMousePosition(), game.camera)
			if point, ok := intersectRayGroundPlane(ray); ok {
				gridX := int(math.Round(float64(point.X)))
				gridZ := int(math.Round(float64(point.Z)))

				fmt.Printf(
					"click world=(%.2f, %.2f, %.2f) grid=(%d, %d)\n",
					point.X,
					point.Y,
					point.Z,
					gridX,
					gridZ,
				)
				game.TryPlaceCube(gridX, gridZ)
			} else {
				fmt.Println("click missed the ground plane")
			}
		}

		rl.BeginDrawing()
		rl.ClearBackground(bgColor)
		rl.BeginMode3D(game.camera)

		game.UpdateSystems()

		rl.DrawGrid(gridSize, gridSpacing)

		rl.EndMode3D()

		rl.EndDrawing()
	}
}

func intersectRayGroundPlane(ray rl.Ray) (rl.Vector3, bool) {
	if ray.Direction.Y > -rayParallelEpsilon && ray.Direction.Y < rayParallelEpsilon {
		return rl.Vector3{}, false
	}

	t := (groundPlaneY - ray.Position.Y) / ray.Direction.Y
	if t < 0 {
		return rl.Vector3{}, false
	}

	return rl.NewVector3(
		ray.Position.X+ray.Direction.X*t,
		groundPlaneY,
		ray.Position.Z+ray.Direction.Z*t,
	), true
}
