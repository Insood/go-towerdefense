package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/mlange-42/ark/ecs"
)

type Game struct {
	models       map[string]*rl.Model
	camera       rl.Camera3D
	cameraSystem *CameraSystem
	shaders      map[string]rl.Shader
	systems      []System
	world        *ecs.World
	cubeMapper   *ecs.Map2[Position3, Renderable]
	cubeSpots    map[gridCell]struct{}
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
		models:       make(map[string]*rl.Model),
		camera:       camera,
		cameraSystem: &CameraSystem{},
		shaders:      make(map[string]rl.Shader),
		world:        ecs.NewWorld(),
		cubeSpots:    make(map[gridCell]struct{}),
	}
	game.cameraSystem.Initialize(game)
	game.loadShaders()
	game.cubeMapper = ecs.NewMap2[Position3, Renderable](game.world)
	game.loadModels()
	game.AddSystem(&RenderSystem{})
	game.InitializeSystems()
	game.placeGroundPlane()

	rl.SetTargetFPS(targetFPS)
	return game
}

func (game *Game) loadShaders() {
	for name, paths := range shaderAssetPaths() {
		game.shaders[name] = rl.LoadShader(paths.vertex, paths.fragment)
	}
}

func (game *Game) loadModels() {
	plane := rl.LoadModelFromMesh(rl.GenMeshPlane(1, 1, 1, 1))
	plane.GetMaterials()[0].Shader = game.shaders["grid"]
	game.models["plane"] = &plane

	checkered_image := rl.GenImageChecked(2, 2, 1, 1, rl.Red, rl.Green)
	texture := rl.LoadTextureFromImage(checkered_image)
	rl.UnloadImage(checkered_image)

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

func (game *Game) UnloadShaders() {
	for _, shader := range game.shaders {
		rl.UnloadShader(shader)
	}
}

type shaderFiles struct {
	vertex   string
	fragment string
}

func shaderAssetPaths() map[string]shaderFiles {
	shaderDir := gameAssetPath("assets", "shaders")
	entries, err := os.ReadDir(shaderDir)
	if err != nil {
		panic(fmt.Errorf("read shader dir %q: %w", shaderDir, err))
	}

	paths := make(map[string]shaderFiles)
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		ext := strings.ToLower(filepath.Ext(entry.Name()))
		switch ext {
		case ".vs", ".vert":
			stem := strings.TrimSuffix(entry.Name(), filepath.Ext(entry.Name()))
			current := paths[stem]
			current.vertex = filepath.Join(shaderDir, entry.Name())
			paths[stem] = current
		case ".fs", ".frag":
			stem := strings.TrimSuffix(entry.Name(), filepath.Ext(entry.Name()))
			current := paths[stem]
			current.fragment = filepath.Join(shaderDir, entry.Name())
			paths[stem] = current
		}
	}

	for name, paths := range paths {
		if paths.vertex == "" || paths.fragment == "" {
			panic(fmt.Errorf("shader %q is missing a vertex or fragment file", name))
		}
	}

	return paths
}

func gameAssetPath(parts ...string) string {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return filepath.Join(parts...)
	}

	base := filepath.Dir(filename)
	segments := make([]string, 0, len(parts)+1)
	segments = append(segments, base)
	segments = append(segments, parts...)
	return filepath.Join(segments...)
}
