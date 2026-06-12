package main

import rl "github.com/gen2brain/raylib-go/raylib"

const (
	windowWidth  = 1280
	windowHeight = 720
	windowTitle  = "Go Tower Defense"
	targetFPS    = 60

	groundPlaneY       = float32(0)
	rayParallelEpsilon = float32(0.0001)

	gridWidth               = 13
	gridLength              = 11
	gridBorderWidth         = 2
	gridCenterX             = gridWidth / 2
	gridCenterZ             = gridLength / 2
	gridTopRow              = 0
	gridBottomRow           = gridLength - 1
	gridLeftCol             = 0
	gridRightCol            = gridWidth - 1
	gridCellCenter          = float32(0.5)
	spireY                  = groundPlaneY + 1.0
	spawnerY                = groundPlaneY + 0.25
	gridDistanceLabelY      = float32(0.0)
	gridDistanceLabelSize   = int32(16)
	gridDistanceLabelOffset = int32(0)
	enemyGoalDelta          = float32(0.05)

	baseCubeSize    = float32(1)
	axisLength      = float32(4)
	cameraPanSpeed  = float32(12)
	cameraZoomSpeed = float32(2)
	cameraMinZoom   = float32(3)
	cameraMaxZoom   = float32(20)

	enemySpeed = 0.5
)

var (
	debugShowGridDistances = false

	cameraPosition    = rl.NewVector3(gridWidth/2.0, 12, gridLength*2/3+8)
	cameraTarget      = rl.NewVector3(gridWidth/2.0, 0, gridLength*2/3)
	cameraUp          = rl.NewVector3(0, 1, 0)
	cameraFOVY        = float32(45)
	baseCubePosition  = rl.Vector3Zero()
	baseCubeColor     = rl.NewColor(198, 120, 76, 255)
	buildableGridTint = rl.NewColor(230, 230, 230, 255)
	noBuildGridTint   = rl.NewColor(140, 80, 80, 255)
	hoverPreviewTint  = rl.NewColor(70, 140, 255, 120)

	bgColor = rl.LightGray
)
