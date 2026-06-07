package main

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/mlange-42/ark/ecs"
)

type Game struct {
	models     map[string]*rl.Model
	camera     rl.Camera3D
	systems    []System
	world      *ecs.World
	cubeMapper *ecs.Map2[Position3, Renderable]
	cubeSpots  map[gridCell]struct{}
}

type gridCell struct {
	X int
	Z int
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
		models:    make(map[string]*rl.Model),
		camera:    camera,
		world:     ecs.NewWorld(),
		cubeSpots: make(map[gridCell]struct{}),
	}
	game.cubeMapper = ecs.NewMap2[Position3, Renderable](game.world)
	game.loadModels()
	game.AddSystem(&RenderSystem{})
	game.InitializeSystems()
	game.placeGroundPlane()

	rl.SetTargetFPS(targetFPS)
	return game
}

func (game *Game) loadModels() {
	checkered_image := rl.GenImageChecked(2, 2, 1, 1, rl.Red, rl.Green)
	texture := rl.LoadTextureFromImage(checkered_image)
	rl.UnloadImage(checkered_image)

	plane := rl.LoadModelFromMesh(rl.GenMeshPlane(1, 1, 1, 1))
	plane.GetMaterials()[0].GetMap(rl.MapDiffuse).Texture = texture
	game.models["plane"] = &plane

	cube := rl.LoadModelFromMesh(rl.GenMeshCube(1, 1, 1))
	cube.GetMaterials()[0].GetMap(rl.MapDiffuse).Texture = texture
	game.models["cube"] = &cube
}

func (game *Game) placeGroundPlane() {
	mapper := ecs.NewMap2[Position3, Renderable](game.world)

	plane := game.models["plane"]
	for x := -float32(gridLength / 2); x < float32(gridLength/2); x++ {
		for z := -float32(gridWidth / 2); z < float32(gridWidth/2); z++ {
			mapper.NewEntity(
				&Position3{X: x, Y: 0, Z: z},
				&Renderable{model: plane, scale: 1.0, tint: rl.White},
			)
		}
	}
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
}

func (game *Game) TryPlaceCube(x, z int) bool {
	cell := gridCell{X: x, Z: z}
	if _, occupied := game.cubeSpots[cell]; occupied {
		fmt.Printf("cube placement blocked: occupied cell (%d, %d)\n", x, z)
		return false
	}

	if x < -gridLength/2 || x >= gridLength/2 || z < -gridWidth/2 || z >= gridWidth/2 {
		fmt.Printf("cube placement blocked: out of bounds (%d, %d)\n", x, z)
		return false
	}

	game.cubeMapper.NewEntity(
		&Position3{
			X: float32(x),
			Y: groundPlaneY + 0.5,
			Z: float32(z),
		},
		&Renderable{
			model: game.models["cube"],
			scale: 1.0,
			tint:  baseCubeColor,
		},
	)
	game.cubeSpots[cell] = struct{}{}

	fmt.Printf("cube placed at grid (%d, %d)\n", x, z)
	return true
}
