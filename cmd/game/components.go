package main

import (
	"image/color"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Position3 rl.Vector3
type Velocity3 rl.Vector3

type Renderable struct {
	model             *rl.Model
	scale             float32
	tint              color.RGBA
	shaderTint        color.RGBA
	shaderTintEnabled bool
}

type Spawner struct{}

type HoverPreview struct {
	gridX int
	gridZ int
}

type Enemy struct{}

type Health struct {
	current float32
	max     float32
}

type HitBox struct {
	size rl.Vector3
}

type WaypointPath struct {
	waypoints      []rl.Vector3
	index          int
	distanceToGoal float32
}

type GunnerTower struct {
	damage       float32
	rangeRadius  float32
	speed        float32
	cooldown     float32
	fireCooldown float32
}

type Projectile struct {
	damage         float32
	remainingRange float32
}

type MoveSpeed struct {
	value float32
}

type HasGravity struct{}

type ReachedGoal struct{}

type Particle struct {
	age          float32
	lifespan     float32
	startColor   color.RGBA
	endColor     color.RGBA
	startSize    float32
	endSize      float32
	currentColor color.RGBA
	currentSize  float32
}
