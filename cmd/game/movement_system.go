package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	ecs "github.com/mlange-42/ark/ecs"
)

type MovementSystem struct {
	filter *ecs.Filter4[Position3, MovementGoal, Movement, Enemy]
}

func (system *MovementSystem) Initialize(game *Game) {
	system.filter = ecs.NewFilter4[Position3, MovementGoal, Movement, Enemy](game.world)
}

func (system *MovementSystem) Update(game *Game) {
	deltaTime := rl.GetFrameTime()
	query := system.filter.Query()
	defer query.Close()

	for query.Next() {
		position, movementGoal, movement, _ := query.Get()

		goalPosition := rl.NewVector3(
			float32(movementGoal.nextGridX)+gridCellCenter,
			position.Y,
			float32(movementGoal.nextGridZ)+gridCellCenter,
		)

		toGoal := rl.Vector3Subtract(goalPosition, rl.Vector3(*position))
		distance := rl.Vector3Length(toGoal)
		if distance <= 0 {
			continue
		}

		maxStep := movement.speed * deltaTime
		if maxStep >= distance {
			*position = Position3(goalPosition)
			continue
		}

		direction := rl.Vector3Scale(rl.Vector3Normalize(toGoal), maxStep)
		*position = Position3(rl.Vector3Add(rl.Vector3(*position), direction))
	}
}
