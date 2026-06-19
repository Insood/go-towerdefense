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

type WaypointPath struct {
	waypoints []rl.Vector3
	index     int
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
