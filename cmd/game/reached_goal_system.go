package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	ecs "github.com/mlange-42/ark/ecs"
)

type ReachedGoalSystem struct {
	filter *ecs.Filter2[Position3, ReachedGoal]
}

func (system *ReachedGoalSystem) Initialize(game *Game) {
	system.filter = ecs.NewFilter2[Position3, ReachedGoal](game.world)
}

func (system *ReachedGoalSystem) Update(game *Game) {
	type reachedGoalParticleBurst struct {
		entity   ecs.Entity
		position Position3
	}

	bursts := make([]reachedGoalParticleBurst, 0)

	query := system.filter.Query()
	defer query.Close()

	for query.Next() {
		position, _ := query.Get()
		bursts = append(bursts, reachedGoalParticleBurst{
			entity:   query.Entity(),
			position: Position3{X: position.X, Y: position.Y, Z: position.Z},
		})
	}

	for _, burst := range bursts {
		game.world.RemoveEntity(burst.entity)
		game.SpawnExplosion(burst.position, enemyReachedExplositionParticles, rl.Orange)
		rl.PlaySound(game.sounds["pop"])
	}
}
