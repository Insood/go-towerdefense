package main

import (
	"testing"

	rl "github.com/gen2brain/raylib-go/raylib"
	ecs "github.com/mlange-42/ark/ecs"
)

func TestGunnerTowerSystemUpdateSpawnsProjectileAfterTowerQueryCloses(t *testing.T) {
	game, system := newGunnerTowerSystemTestFixture(t)

	towerMapper := ecs.NewMap2[Position3, GunnerTower](game.world)
	towerMapper.NewEntity(
		&Position3{X: 2.5, Y: 0, Z: 2.5},
		&GunnerTower{
			damage:       gunnerTowerDamage,
			rangeRadius:  gunnerTowerRange,
			speed:        gunnerTowerProjectileSpeed,
			cooldown:     gunnerTowerCooldown,
			fireCooldown: 0,
		},
	)

	enemyMapper := ecs.NewMap3[Position3, WaypointPath, Enemy](game.world)
	enemyMapper.NewEntity(
		&Position3{X: 1.5, Y: 0, Z: 2.5},
		&WaypointPath{distanceToGoal: 2},
		&Enemy{},
	)

	assertNotPanics(t, func() {
		system.Update(game)
	})

	projectileFilter := ecs.NewFilter3[Position3, Velocity3, Projectile](game.world)
	query := projectileFilter.Query()
	defer query.Close()

	count := 0
	for query.Next() {
		count++
	}

	if count != 1 {
		t.Fatalf("projectile count = %d, want 1", count)
	}
}

func TestGunnerTowerSystemUpdateTargetsEnemyClosestToGoal(t *testing.T) {
	game, system := newGunnerTowerSystemTestFixture(t)

	towerMapper := ecs.NewMap2[Position3, GunnerTower](game.world)
	towerMapper.NewEntity(
		&Position3{X: 2.5, Y: 0, Z: 2.5},
		&GunnerTower{
			damage:       gunnerTowerDamage,
			rangeRadius:  gunnerTowerRange,
			speed:        gunnerTowerProjectileSpeed,
			cooldown:     gunnerTowerCooldown,
			fireCooldown: 0,
		},
	)

	enemyMapper := ecs.NewMap3[Position3, WaypointPath, Enemy](game.world)
	enemyMapper.NewEntity(
		&Position3{X: 3.5, Y: 0, Z: 2.5},
		&WaypointPath{distanceToGoal: 8},
		&Enemy{},
	)
	enemyMapper.NewEntity(
		&Position3{X: 1.5, Y: 0, Z: 2.5},
		&WaypointPath{distanceToGoal: 2},
		&Enemy{},
	)

	system.Update(game)

	projectileFilter := ecs.NewFilter3[Position3, Velocity3, Projectile](game.world)
	query := projectileFilter.Query()
	defer query.Close()

	if !query.Next() {
		t.Fatal("expected one projectile to be spawned")
	}

	_, velocity, _ := query.Get()
	got := rl.Vector3(*velocity)
	if got.X >= 0 {
		t.Fatalf("projectile velocity X = %v, want negative X toward the lower distance-to-goal enemy", got.X)
	}
	if got.Z != 0 {
		t.Fatalf("projectile velocity Z = %v, want 0", got.Z)
	}
}

func newGunnerTowerSystemTestFixture(t *testing.T) (*Game, *GunnerTowerSystem) {
	t.Helper()

	game := &Game{
		assets: NewAssetManager(),
		world:  ecs.NewWorld(),
	}
	game.assets.models["miniMob"] = &rl.Model{}

	system := &GunnerTowerSystem{}
	system.Initialize(game)
	return game, system
}

func assertNotPanics(t *testing.T, fn func()) {
	t.Helper()

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("unexpected panic: %v", r)
		}
	}()

	fn()
}
