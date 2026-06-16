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

type Movement struct {
	speed float32
}

type MovementGoal struct {
	nextGridX int
	nextGridZ int
}

type Enemy struct{}

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
