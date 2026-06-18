package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	ecs "github.com/mlange-42/ark/ecs"
)

type ParticleRenderSystem struct {
	filter *ecs.Filter2[Position3, Particle]
}

func (system *ParticleRenderSystem) Initialize(game *Game) {
	system.filter = ecs.NewFilter2[Position3, Particle](game.world)
}

func (system *ParticleRenderSystem) Update(game *Game) {
	rl.BeginMode3D(game.camera)
	system.renderParticles(game)
	rl.EndMode3D()
}

func (system *ParticleRenderSystem) renderParticles(game *Game) {
	texture := game.assets.Texture("white")
	source := rl.NewRectangle(0, 0, 1, 1)
	up := rl.NewVector3(0, 1, 0)
	origin := rl.NewVector2(0, 0)

	query := system.filter.Query()
	defer query.Close()

	for query.Next() {
		position, particle := query.Get()
		if particle.currentSize <= 0 {
			continue
		}

		size := rl.NewVector2(particle.currentSize, particle.currentSize)
		rl.DrawBillboardPro(
			game.camera,
			texture,
			source,
			rl.Vector3(*position),
			up,
			size,
			origin,
			0,
			particle.currentColor,
		)
	}
}
