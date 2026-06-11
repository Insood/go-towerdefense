package main

import (
	"image/color"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Position3 = rl.Vector3

type Renderable struct {
	model             *rl.Model
	scale             float32
	tint              color.RGBA
	shaderTint        color.RGBA
	shaderTintEnabled bool
}

type Spawner struct{}

type MovementGoal struct {
	nextGridX int
	nextGridY int
}

type Enemy struct{}
