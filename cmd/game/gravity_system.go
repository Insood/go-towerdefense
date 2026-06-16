package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	ecs "github.com/mlange-42/ark/ecs"
)

type GravitySystem struct {
	filter *ecs.Filter3[Position3, Velocity3, HasGravity]
}

func (system *GravitySystem) Initialize(game *Game) {
	system.filter = ecs.NewFilter3[Position3, Velocity3, HasGravity](game.world)
}

func (system *GravitySystem) Update(game *Game) {
	deltaTime := rl.GetFrameTime()
	query := system.filter.Query()
	defer query.Close()

	for query.Next() {
		_, velocity, _ := query.Get()
		velocity.Y -= gravityAcceleration * deltaTime
	}
}
