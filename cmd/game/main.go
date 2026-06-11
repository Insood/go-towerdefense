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
		rl.BeginMode3D(game.camera)

		game.UpdateSystems()
		drawCoordinateSystem()

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

func drawCoordinateSystem() {
	origin := rl.Vector3Zero()
	rl.DrawLine3D(origin, rl.NewVector3(axisLength, 0, 0), rl.Red)
	rl.DrawLine3D(origin, rl.NewVector3(0, axisLength, 0), rl.Green)
	rl.DrawLine3D(origin, rl.NewVector3(0, 0, axisLength), rl.Blue)
}
