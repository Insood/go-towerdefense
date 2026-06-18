package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type AssetManager struct {
	models   map[string]*rl.Model
	textures map[string]rl.Texture2D
	sounds   map[string]rl.Sound
	shaders  map[string]rl.Shader
}

func NewAssetManager() *AssetManager {
	return &AssetManager{
		models:   make(map[string]*rl.Model),
		textures: make(map[string]rl.Texture2D),
		sounds:   make(map[string]rl.Sound),
		shaders:  make(map[string]rl.Shader),
	}
}

func (assets *AssetManager) Load() {
	assets.loadShaders()
	assets.loadTextures()
	assets.loadModels()
	assets.loadSounds()
}

func (assets *AssetManager) Shader(name string) rl.Shader {
	return assets.shaders[name]
}

func (assets *AssetManager) Model(name string) *rl.Model {
	return assets.models[name]
}

func (assets *AssetManager) Texture(name string) rl.Texture2D {
	return assets.textures[name]
}

func (assets *AssetManager) PlaySound(name string) {
	sound, ok := assets.sounds[name]
	if !ok {
		return
	}

	rl.PlaySound(sound)
}

func (assets *AssetManager) Unload() {
	for _, model := range assets.models {
		rl.UnloadModel(*model)
	}
	for _, texture := range assets.textures {
		rl.UnloadTexture(texture)
	}
	for _, sound := range assets.sounds {
		rl.UnloadSound(sound)
	}
	for _, shader := range assets.shaders {
		rl.UnloadShader(shader)
	}
}

func (assets *AssetManager) loadShaders() {
	for name, paths := range shaderAssetPaths() {
		assets.shaders[name] = rl.LoadShader(paths.vertex, paths.fragment)
	}
}

func (assets *AssetManager) loadTextures() {
	whiteImage := rl.GenImageColor(1, 1, rl.White)
	assets.textures["white"] = rl.LoadTextureFromImage(whiteImage)
	rl.UnloadImage(whiteImage)
}

func (assets *AssetManager) loadModels() {
	plane := rl.LoadModelFromMesh(rl.GenMeshPlane(1, 1, 1, 1))
	plane.GetMaterials()[0].Shader = assets.shaders["grid"]
	assets.models["plane"] = &plane

	turret := rl.LoadModel(gameAssetPath("assets", "models", "turret.glb"))
	assets.models["turret"] = &turret

	spire := rl.LoadModel(gameAssetPath("assets", "models", "spire.glb"))
	assets.models["spire"] = &spire

	spawner := rl.LoadModel(gameAssetPath("assets", "models", "spawner.glb"))
	assets.models["spawner"] = &spawner

	mobCheckeredImage := rl.GenImageChecked(2, 2, 1, 1, rl.Orange, rl.Purple)
	assets.textures["miniMob"] = rl.LoadTextureFromImage(mobCheckeredImage)
	rl.UnloadImage(mobCheckeredImage)

	miniMob := rl.LoadModelFromMesh(rl.GenMeshCube(0.25, 0.25, 0.25))
	miniMob.GetMaterials()[0].GetMap(rl.MapDiffuse).Texture = assets.textures["miniMob"]
	assets.models["miniMob"] = &miniMob
}

func (assets *AssetManager) loadSounds() {
	assets.sounds["pop"] = rl.LoadSound(gameAssetPath("assets", "sounds", "pop.wav"))
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
