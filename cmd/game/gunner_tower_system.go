package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	ecs "github.com/mlange-42/ark/ecs"
)

type GunnerTowerSystem struct {
	towerFilter      *ecs.Filter2[Position3, GunnerTower]
	enemyFilter      *ecs.Filter3[Position3, WaypointPath, Enemy]
	projectileMapper *ecs.Map4[Position3, Renderable, Velocity3, Projectile]
}

type gunnerTowerTarget struct {
	position       Position3
	distanceToGoal float32
}

type gunnerTowerShot struct {
	position   Position3
	velocity   Velocity3
	projectile Projectile
}

func (system *GunnerTowerSystem) Initialize(game *Game) {
	system.towerFilter = ecs.NewFilter2[Position3, GunnerTower](game.world)
	system.enemyFilter = ecs.NewFilter3[Position3, WaypointPath, Enemy](game.world)
	system.projectileMapper = ecs.NewMap4[Position3, Renderable, Velocity3, Projectile](game.world)
}

func (system *GunnerTowerSystem) Update(game *Game) {
	deltaTime := rl.GetFrameTime()
	shots := make([]gunnerTowerShot, 0)

	query := system.towerFilter.Query()

	for query.Next() {
		position, tower := query.Get()

		if tower.fireCooldown > 0 {
			tower.fireCooldown = clampFloat32(tower.fireCooldown-deltaTime, 0, tower.cooldown)
		}

		if tower.fireCooldown > 0 {
			continue
		}

		target, ok := system.selectTarget(*position, *tower)
		if !ok {
			continue
		}

		shot, ok := system.buildShot(game, *position, *tower, target)
		if ok {
			shots = append(shots, shot)
		}
		tower.fireCooldown = tower.cooldown
	}
	query.Close()

	for _, shot := range shots {
		system.spawnShot(game, shot)
	}
}

func (system *GunnerTowerSystem) selectTarget(towerPosition Position3, tower GunnerTower) (gunnerTowerTarget, bool) {
	query := system.enemyFilter.Query()
	defer query.Close()

	targets := make([]gunnerTowerTarget, 0)
	towerWorldPosition := rl.Vector3(towerPosition)

	for query.Next() {
		position, path, _ := query.Get()
		if path.distanceToGoal <= 0 {
			continue
		}

		candidatePosition := rl.Vector3(*position)
		if horizontalDistance(towerWorldPosition, candidatePosition) > tower.rangeRadius {
			continue
		}

		targets = append(targets, gunnerTowerTarget{
			position:       *position,
			distanceToGoal: path.distanceToGoal,
		})
	}

	return pickGunnerTowerTarget(targets)
}

func pickGunnerTowerTarget(targets []gunnerTowerTarget) (gunnerTowerTarget, bool) {
	if len(targets) == 0 {
		return gunnerTowerTarget{}, false
	}

	best := targets[0]
	for _, target := range targets[1:] {
		if target.distanceToGoal < best.distanceToGoal {
			best = target
		}
	}

	return best, true
}

func (system *GunnerTowerSystem) buildShot(game *Game, towerPosition Position3, tower GunnerTower, target gunnerTowerTarget) (gunnerTowerShot, bool) {
	towerWorldPosition := rl.Vector3(towerPosition)
	targetWorldPosition := rl.Vector3(target.position)
	direction := rl.Vector3Subtract(targetWorldPosition, towerWorldPosition)
	if rl.Vector3Length(direction) <= 0 {
		return gunnerTowerShot{}, false
	}

	velocity := rl.Vector3Scale(rl.Vector3Normalize(direction), tower.speed)
	return gunnerTowerShot{
		position: Position3(towerWorldPosition),
		velocity: Velocity3(velocity),
		projectile: Projectile{
			damage:         tower.damage,
			remainingRange: tower.rangeRadius,
		},
	}, true
}

func (system *GunnerTowerSystem) spawnShot(game *Game, shot gunnerTowerShot) {
	system.projectileMapper.NewEntity(
		&shot.position,
		&Renderable{
			model:             game.assets.Model("miniMob"),
			scale:             gunnerTowerProjectileSize,
			tint:              rl.NewColor(250, 220, 90, 255),
			shaderTintEnabled: false,
		},
		&shot.velocity,
		&shot.projectile,
	)
}
