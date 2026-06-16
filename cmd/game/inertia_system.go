package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	ecs "github.com/mlange-42/ark/ecs"
)

type InertiaSystem struct {
	filter *ecs.Filter2[Position3, Velocity3]
}

func (system *InertiaSystem) Initialize(game *Game) {
	system.filter = ecs.NewFilter2[Position3, Velocity3](game.world)
}

func (system *InertiaSystem) Update(game *Game) {
	deltaTime := rl.GetFrameTime()
	query := system.filter.Query()
	defer query.Close()

	for query.Next() {
		position, velocity := query.Get()
		step := rl.Vector3Scale(rl.Vector3(*velocity), deltaTime)
		*position = Position3(rl.Vector3Add(rl.Vector3(*position), step))
	}
}
