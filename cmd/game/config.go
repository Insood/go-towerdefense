package main

import rl "github.com/gen2brain/raylib-go/raylib"

const (
	windowWidth  = 1280
	windowHeight = 720
	windowTitle  = "Go Tower Defense"
	targetFPS    = 60

	groundPlaneY       = float32(0)
	rayParallelEpsilon = float32(0.0001)

	gridSize    = 20
	gridSpacing = float32(1)
	gridWidth   = 10
	gridLength  = 20

	baseCubeSize    = float32(1)
	cameraPanSpeed  = float32(12)
	cameraZoomSpeed = float32(2)
	cameraMinZoom   = float32(3)
	cameraMaxZoom   = float32(20)
)

var (
	cameraPosition = rl.NewVector3(8, 6, 0)
	cameraTarget   = rl.Vector3Zero()
	cameraUp       = rl.NewVector3(0, 1, 0)
	cameraFOVY     = float32(45)

	baseCubePosition = rl.Vector3Zero()
	baseCubeColor    = rl.NewColor(198, 120, 76, 255)

	bgColor = rl.LightGray
)
