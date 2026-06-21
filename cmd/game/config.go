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
	gridCenterX     = gridWidth / 2
	gridCenterZ     = gridLength / 2
	gridTopRow      = 0
	gridBottomRow   = gridLength - 1
	gridLeftCol     = 0
	gridRightCol    = gridWidth - 1
	gridCellCenter  = float32(0.5)

	spireY   = groundPlaneY
	spawnerY = groundPlaneY

	gridDistanceLabelY      = float32(0.0)
	gridDistanceLabelSize   = int32(16)
	gridDistanceLabelOffset = int32(0)

	healthBarWidth        = float32(42)
	healthBarHeight       = float32(6)
	healthBarBorderWidth  = float32(1)
	healthBarWorldYOffset = float32(1.4)

	spireMaxHealth = float32(20)

	enemyMaxHealth                   = float32(10)
	enemyMoveSpeed                   = float32(0.5)
	enemyWaypointDelta               = float32(0.05)
	enemyReachedExplositionParticles = 100

	axisLength = float32(4)

	cameraPanSpeed  = float32(12)
	cameraZoomSpeed = float32(2)
	cameraMinZoom   = float32(3)
	cameraMaxZoom   = float32(20)

	gravityAcceleration = float32(1.0)

	explosionParticleCountMin = 8
	explosionParticleCountMax = 14
	explosionSpeedMin         = float32(1.5)
	explosionSpeedMax         = float32(3.5)
	explosionHeightBoostMin   = float32(0.5)
	explosionHeightBoostMax   = float32(2.0)
	explosionSizeMin          = float32(0.08)
	explosionSizeMax          = float32(0.20)
	explosionLifespanMin      = float32(0.35)
	explosionLifespanMax      = float32(0.85)
)

var (
	debugEnabled = false

	cameraPosition             = rl.NewVector3(gridWidth/2.0, 12, gridLength*2/3+8)
	cameraTarget               = rl.NewVector3(gridWidth/2.0, 0, gridLength*2/3)
	cameraUp                   = rl.NewVector3(0, 1, 0)
	cameraFOVY                 = float32(45)
	baseTurretTint             = rl.RayWhite
	buildableGridTint          = rl.NewColor(230, 230, 230, 255)
	noBuildGridTint            = rl.NewColor(140, 80, 80, 255)
	hoverPreviewTintAllowed    = rl.NewColor(70, 140, 255, 120)
	hoverPreviewTintNotAllowed = rl.NewColor(255, 150, 150, 110)
	healthBarBackTint          = rl.NewColor(20, 20, 20, 210)
	healthBarBorderTint        = rl.NewColor(10, 10, 10, 255)

	bgColor = rl.LightGray
)
