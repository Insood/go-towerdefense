package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	ecs "github.com/mlange-42/ark/ecs"
)

type System interface {
	Initialize(*Game)
	Update(*Game)
}

type RenderSystem struct {
	filter *ecs.Filter2[Position3, Renderable]
}

func (system *RenderSystem) Initialize(game *Game) {
	system.filter = ecs.NewFilter2[Position3, Renderable](game.world)
}

func (system *RenderSystem) Update(game *Game) {
	query := system.filter.Query()

	for query.Next() {
		position, renderable := query.Get()

		rl.DrawModel(*renderable.model, *position, renderable.scale, renderable.tint)
	}
}
