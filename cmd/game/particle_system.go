package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	ecs "github.com/mlange-42/ark/ecs"
)

type ParticleSystem struct {
	filter  *ecs.Filter2[Particle, Renderable]
	expired []ecs.Entity
}

func (system *ParticleSystem) Initialize(game *Game) {
	system.filter = ecs.NewFilter2[Particle, Renderable](game.world)
	system.expired = system.expired[:0]
}

func (system *ParticleSystem) Update(game *Game) {
	deltaTime := rl.GetFrameTime()
	query := system.filter.Query()
	defer query.Close()

	system.expired = system.expired[:0]

	for query.Next() {
		particle, renderable := query.Get()
		particle.age += deltaTime

		if particle.lifespan <= 0 {
			particle.currentColor = particle.endColor
			particle.currentSize = particle.endSize
			renderable.scale = particle.currentSize
			renderable.tint = particle.currentColor
			system.expired = append(system.expired, query.Entity())
			continue
		}

		t := clampFloat32(particle.age/particle.lifespan, 0, 1)
		particle.currentColor = lerpColorRGBA(particle.startColor, particle.endColor, t)
		particle.currentSize = lerpFloat32(particle.startSize, particle.endSize, t)
		renderable.scale = particle.currentSize
		renderable.tint = particle.currentColor

		if particle.age >= particle.lifespan {
			system.expired = append(system.expired, query.Entity())
		}
	}

	for _, entity := range system.expired {
		game.world.RemoveEntity(entity)
	}
}
