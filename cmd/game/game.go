package main

import (
	"image/color"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/mlange-42/ark/ecs"
)

type Game struct {
	assets  *AssetManager
	camera  rl.Camera3D
	grid    GameGrid
	systems []System
	world   *ecs.World
	tick    int
}

func InitializeGame() *Game {
	camera := rl.Camera3D{
		Position:   cameraPosition,
		Target:     cameraTarget,
		Up:         cameraUp,
		Fovy:       cameraFOVY,
		Projection: rl.CameraPerspective,
	}

	game := &Game{
		assets: NewAssetManager(),
		camera: camera,
		grid:   NewGameGrid(gridWidth, gridLength),
		world:  ecs.NewWorld(),
	}
	game.grid.Initialize(game.world)
	game.assets.Load()
	game.placeSpire()
	game.AddSystem(&CameraSystem{})
	game.AddSystem(&HoverPreviewSystem{})
	game.AddSystem(&InputSystem{})
	game.AddSystem(&SpawnerSystem{})
	game.AddSystem(&WaypointSystem{})
	game.AddSystem(&GravitySystem{})
	game.AddSystem(&InertiaSystem{})
	game.AddSystem(&ReachedGoalSystem{})
	game.AddSystem(&ParticleSystem{})
	game.AddSystem(&RenderSystem3D{})
	game.AddSystem(&ParticleRenderSystem{})
	game.AddSystem(&HealthBarOverlaySystem{})
	game.AddSystem(&DebugRender3DSystem{})
	game.AddSystem(&DebugRenderOverlaySystem{})
	game.InitializeSystems()
	game.placeModels()

	rl.SetTargetFPS(targetFPS)
	return game
}

func (game *Game) placeModels() {
	modelMapper := ecs.NewMap2[Position3, Renderable](game.world)

	plane := game.assets.Model("plane")
	for x := 0; x < gridWidth; x++ {
		for z := 0; z < gridLength; z++ {
			shaderTint := buildableGridTint
			if cell, ok := game.grid.Cell(x, z); ok && !cell.Buildable() {
				shaderTint = noBuildGridTint
			}

			modelMapper.NewEntity(
				&Position3{X: float32(x) + gridCellCenter, Y: 0, Z: float32(z) + gridCellCenter},
				&Renderable{
					model:             plane,
					scale:             1.0,
					tint:              rl.White,
					shaderTint:        shaderTint,
					shaderTintEnabled: true,
				},
			)
		}
	}
}

func (game *Game) placeSpire() {
	spire := game.assets.Model("spire")

	spireMapper := ecs.NewMap3[Position3, Renderable, Health](game.world)
	spireEntity := spireMapper.NewEntity(
		&Position3{
			X: float32(gridCenterX) + gridCellCenter,
			Y: spireY,
			Z: float32(gridCenterZ) + gridCellCenter,
		},
		&Renderable{
			model:             spire,
			scale:             1.0,
			tint:              rl.White,
			shaderTintEnabled: false,
		},
		&Health{
			current: spireMaxHealth,
			max:     spireMaxHealth,
		},
	)

	if !game.grid.SetCellEntityForce(gridCenterX, gridCenterZ, spireEntity) {
		panic("failed to place spire at the center of the grid")
	}
}

func (game *Game) PlaceTower(x, z int, model *rl.Model, tint color.RGBA) bool {
	if !game.grid.PlaceEntity(x, z, model, tint) {
		return false
	}

	game.RebuildEnemyPaths()
	return true
}

func (game *Game) RebuildEnemyPaths() {
	pathFilter := ecs.NewFilter4[Position3, WaypointPath, Velocity3, Enemy](game.world)
	query := pathFilter.Query()
	defer query.Close()

	for query.Next() {
		position, path, velocity, _ := query.Get()
		gridCoord := gridCoordFromPosition(*position)
		waypoints := buildWaypointPath(game.grid.PathToCenter(gridCoord.X, gridCoord.Z))
		path.waypoints = waypoints
		path.index = 0
		*velocity = Velocity3{}
	}
}

func (game *Game) SpawnExplosion(position Position3, count int, startColor color.RGBA) {
	if count <= 0 {
		return
	}

	particleMapper := ecs.NewMap4[Position3, Velocity3, Particle, HasGravity](game.world)
	for range count {
		theta := float32(rng.Float64() * 2 * math.Pi)
		phi := float32((rng.Float64() * 0.5) * math.Pi)
		speed := explosionSpeedMin + float32(rng.Float64())*(explosionSpeedMax-explosionSpeedMin)
		heightBoost := explosionHeightBoostMin + float32(rng.Float64())*(explosionHeightBoostMax-explosionHeightBoostMin)

		direction := rl.NewVector3(
			float32(math.Cos(float64(theta))*math.Sin(float64(phi))),
			float32(math.Cos(float64(phi))),
			float32(math.Sin(float64(theta))*math.Sin(float64(phi))),
		)
		velocity := rl.Vector3Scale(direction, speed)
		velocity.Y += heightBoost
		particleVelocity := Velocity3(velocity)

		size := explosionSizeMin + float32(rng.Float64())*(explosionSizeMax-explosionSizeMin)
		lifespan := explosionLifespanMin + float32(rng.Float64())*(explosionLifespanMax-explosionLifespanMin)
		endColor := startColor
		endColor.A = 0

		particleMapper.NewEntity(
			&Position3{X: position.X, Y: position.Y, Z: position.Z},
			&particleVelocity,
			&Particle{
				age:          0,
				lifespan:     lifespan,
				startColor:   startColor,
				endColor:     endColor,
				startSize:    size,
				endSize:      0,
				currentColor: startColor,
				currentSize:  size,
			},
			&HasGravity{},
		)
	}
}

func (game *Game) PlaySound(name string) {
	game.assets.PlaySound(name)
}

func (game *Game) AddSystem(system System) {
	game.systems = append(game.systems, system)
}

func (game *Game) InitializeSystems() {
	for _, system := range game.systems {
		system.Initialize(game)
	}
}

func (game *Game) UpdateSystems() {
	for _, system := range game.systems {
		system.Update(game)
	}
	game.tick += 1
}

func (game *Game) UnloadAssets() {
	game.assets.Unload()
}
