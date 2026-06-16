package main

import (
	"fmt"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type InputSystem struct{}

func (system *InputSystem) Initialize(game *Game) {}

func (system *InputSystem) Update(game *Game) {
	if rl.IsKeyPressed(rl.KeyF11) {
		debugEnabled = !debugEnabled
	}

	if rl.IsMouseButtonPressed(rl.MouseButtonLeft) {
		system.OnLeftClick(game)
	}

	if rl.IsMouseButtonPressed(rl.MouseButtonRight) {
		system.OnRightClick(game)
	}
}

func (system *InputSystem) OnLeftClick(game *Game) {
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
		return
	}

	fmt.Println("click missed the ground plane")
}

func (system *InputSystem) OnRightClick(game *Game) {
	ray := rl.GetScreenToWorldRay(rl.GetMousePosition(), game.camera)
	if point, ok := intersectRayGroundPlane(ray); ok {
		gridX := int(math.Floor(float64(point.X)))
		gridZ := int(math.Floor(float64(point.Z)))
		position := rl.NewVector3(float32(gridX)+1, 0.25, float32(gridZ)+0.5)
		fmt.Println(position)
		game.SpawnExplosion(position, 100, rl.Red)
	}
}
