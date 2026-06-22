package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	ecs "github.com/mlange-42/ark/ecs"
)

type ProjectileSystem struct {
	projectileFilter *ecs.Filter3[Position3, Velocity3, Projectile]
	enemyFilter      *ecs.Filter4[Position3, HitBox, Health, Enemy]
}

type projectileHitEvent struct {
	entity   ecs.Entity
	position Position3
}

func (system *ProjectileSystem) Initialize(game *Game) {
	system.projectileFilter = ecs.NewFilter3[Position3, Velocity3, Projectile](game.world)
	system.enemyFilter = ecs.NewFilter4[Position3, HitBox, Health, Enemy](game.world)
}

func (system *ProjectileSystem) Update(game *Game) {
	deltaTime := rl.GetFrameTime()

	projectilesToRemove := make([]ecs.Entity, 0)
	enemyDeaths := make([]projectileHitEvent, 0)
	deadEnemies := make(map[ecs.Entity]struct{})

	projectileQuery := system.projectileFilter.Query()
	defer projectileQuery.Close()

	for projectileQuery.Next() {
		position, velocity, projectile := projectileQuery.Get()

		if projectile.remainingRange <= 0 {
			projectilesToRemove = append(projectilesToRemove, projectileQuery.Entity())
			continue
		}

		projectilePosition := rl.Vector3(*position)
		projectileEntity := projectileQuery.Entity()
		if system.applyProjectileHit(&projectilePosition, projectile.damage, projectileEntity, &projectilesToRemove, &enemyDeaths, deadEnemies) {
			continue
		}

		stepDistance := rl.Vector3Length(rl.Vector3(*velocity)) * deltaTime
		projectile.remainingRange = clampFloat32(projectile.remainingRange-stepDistance, 0, projectile.remainingRange)
		if projectile.remainingRange <= 0 {
			projectilesToRemove = append(projectilesToRemove, projectileEntity)
		}
	}

	for _, entity := range projectilesToRemove {
		game.world.RemoveEntity(entity)
	}

	for _, death := range enemyDeaths {
		game.world.RemoveEntity(death.entity)
		game.SpawnExplosion(death.position, enemyReachedExplositionParticles, rl.Orange)
		game.PlaySound("pop")
	}
}

func (system *ProjectileSystem) applyProjectileHit(
	projectilePosition *rl.Vector3,
	damage float32,
	projectileEntity ecs.Entity,
	projectilesToRemove *[]ecs.Entity,
	enemyDeaths *[]projectileHitEvent,
	deadEnemies map[ecs.Entity]struct{},
) bool {
	query := system.enemyFilter.Query()
	defer query.Close()

	for query.Next() {
		position, hitBox, health, _ := query.Get()
		enemyEntity := query.Entity()
		if _, ok := deadEnemies[enemyEntity]; ok {
			continue
		}

		enemyPosition := rl.Vector3(*position)
		if !pointInsideHitBox(*projectilePosition, enemyPosition, *hitBox) {
			continue
		}

		health.current -= damage
		*projectilesToRemove = append(*projectilesToRemove, projectileEntity)
		if health.current <= 0 {
			deadEnemies[enemyEntity] = struct{}{}
			*enemyDeaths = append(*enemyDeaths, projectileHitEvent{
				entity:   enemyEntity,
				position: *position,
			})
		}

		return true
	}

	return false
}
