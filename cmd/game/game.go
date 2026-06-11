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
	grid         GameGrid
	systems      []System
	world        *ecs.World
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
		grid:         NewGameGrid(gridWidth, gridLength),
		world:        ecs.NewWorld(),
	}
	game.cameraSystem.Initialize(game)
	game.grid.Initialize(game.world)
	game.loadShaders()
	game.loadModels()
	game.placeSpire()
	game.AddSystem(&InputSystem{})
	game.AddSystem(&SpawnerSystem{})
	game.AddSystem(&RenderSystem3D{})
	game.AddSystem(&GridDistanceDebugRenderSystem{})
	game.InitializeSystems()
	game.placeModels()

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

	checkeredImage := rl.GenImageChecked(2, 2, 1, 1, rl.Red, rl.Green)
	texture := rl.LoadTextureFromImage(checkeredImage)
	rl.UnloadImage(checkeredImage)

	cube := rl.LoadModelFromMesh(rl.GenMeshCube(1, 1, 1))
	cube.GetMaterials()[0].GetMap(rl.MapDiffuse).Texture = texture
	game.models["cube"] = &cube

	spire := rl.LoadModelFromMesh(rl.GenMeshCube(1, 2, 1))
	spire.GetMaterials()[0].GetMap(rl.MapDiffuse).Texture = texture
	game.models["spire"] = &spire

	spawner := rl.LoadModelFromMesh(rl.GenMeshCube(1, 0.5, 1))
	spawner.GetMaterials()[0].GetMap(rl.MapDiffuse).Texture = texture
	game.models["spawner"] = &spawner
}

func (game *Game) placeModels() {
	modelMapper := ecs.NewMap2[Position3, Renderable](game.world)

	plane := game.models["plane"]
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
	spire := game.models["spire"]

	if !game.grid.ForcePlaceEntity(gridCenterX, gridCenterZ, spireY, spire, rl.White) {
		panic("failed to place spire at the center of the grid")
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
		stem := strings.TrimSuffix(entry.Name(), filepath.Ext(entry.Name()))

		switch ext {
		case ".vs", ".vert":
			current := paths[stem]
			current.vertex = filepath.Join(shaderDir, entry.Name())
			paths[stem] = current
		case ".fs", ".frag":
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
