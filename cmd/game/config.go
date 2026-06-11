package main

import rl "github.com/gen2brain/raylib-go/raylib"

const (
	windowWidth  = 1280
	windowHeight = 720
	windowTitle  = "Go Tower Defense"
	targetFPS    = 60

	groundPlaneY       = float32(0)
	rayParallelEpsilon = float32(0.0001)

	gridWidth       = 13
	gridLength      = 11
	gridBorderWidth = 2

	baseCubeSize    = float32(1)
	axisLength      = float32(4)
	cameraPanSpeed  = float32(12)
	cameraZoomSpeed = float32(2)
	cameraMinZoom   = float32(3)
	cameraMaxZoom   = float32(20)
)

var (
	cameraPosition    = rl.NewVector3(gridWidth/2.0, 8, gridLength*2/3+8)
	cameraTarget      = rl.NewVector3(gridWidth/2.0, 0, gridLength*2/3)
	cameraUp          = rl.NewVector3(0, 1, 0)
	cameraFOVY        = float32(45)
	baseCubePosition  = rl.Vector3Zero()
	baseCubeColor     = rl.NewColor(198, 120, 76, 255)
	buildableGridTint = rl.NewColor(230, 230, 230, 255)
	noBuildGridTint   = rl.NewColor(140, 80, 80, 255)

	bgColor = rl.LightGray
)
