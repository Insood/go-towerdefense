package main

import (
	"fmt"
	"image/color"

	rl "github.com/gen2brain/raylib-go/raylib"
	ecs "github.com/mlange-42/ark/ecs"
)

type System interface {
	Initialize(*Game)
	Update(*Game)
}

type CameraSystem struct{}

func (system *CameraSystem) Initialize(game *Game) {}

func (system *CameraSystem) Update(game *Game) {
	camera := &game.camera
	cameraVector := rl.Vector3Subtract(camera.Target, camera.Position)
	fmt.Printf(
		"camera pos=(%.2f, %.2f, %.2f) target=(%.2f, %.2f, %.2f) vector=(%.2f, %.2f, %.2f)\n",
		camera.Position.X,
		camera.Position.Y,
		camera.Position.Z,
		camera.Target.X,
		camera.Target.Y,
		camera.Target.Z,
		cameraVector.X,
		cameraVector.Y,
		cameraVector.Z,
	)

	frameStep := cameraPanSpeed * rl.GetFrameTime()

	var moveX float32
	var moveZ float32
	if rl.IsKeyDown(rl.KeyW) {
		moveZ -= 1
	}
	if rl.IsKeyDown(rl.KeyS) {
		moveZ += 1
	}
	if rl.IsKeyDown(rl.KeyD) {
		moveX += 1
	}
	if rl.IsKeyDown(rl.KeyA) {
		moveX -= 1
	}

	pan := cameraMoveOnGround(moveX, moveZ, frameStep)
	camera.Position = rl.Vector3Add(camera.Position, pan)
	camera.Target = rl.Vector3Add(camera.Target, pan)

	wheel := rl.GetMouseWheelMove()
	if wheel != 0 {
		forward := rl.Vector3Subtract(camera.Target, camera.Position)
		if forward.X != 0 || forward.Y != 0 || forward.Z != 0 {
			forward = rl.Vector3Normalize(forward)
			distance := rl.Vector3Distance(camera.Position, camera.Target)
			distance -= wheel * cameraZoomSpeed
			if distance < cameraMinZoom {
				distance = cameraMinZoom
			}
			if distance > cameraMaxZoom {
				distance = cameraMaxZoom
			}

			camera.Position = rl.Vector3Subtract(camera.Target, rl.Vector3Scale(forward, distance))
		}
	}
}

type RenderSystem struct {
	filter *ecs.Filter2[Position3, Renderable]
}

func (system *RenderSystem) Initialize(game *Game) {
	system.filter = ecs.NewFilter2[Position3, Renderable](game.world)
}

func (system *RenderSystem) Update(game *Game) {
	query := system.filter.Query()

	for query.Next() {
		position, renderable := query.Get()

		drawTint := renderable.tint
		if renderable.shaderTintEnabled {
			materials := renderable.model.GetMaterials()
			if len(materials) > 0 {
				shader := materials[0].Shader
				location := rl.GetShaderLocation(shader, "tintColor")
				if location >= 0 {
					tint := colorToVec4(renderable.shaderTint)
					rl.SetShaderValue(
						shader,
						location,
						tint[:],
						rl.ShaderUniformVec4,
					)
				}
			}
			drawTint = rl.White
		}

		rl.DrawModel(*renderable.model, *position, renderable.scale, drawTint)
	}
}

func cameraMoveOnGround(x, z, distance float32) rl.Vector3 {
	return rl.NewVector3(x*distance, 0, z*distance)
}

func colorToVec4(c color.RGBA) [4]float32 {
	return [4]float32{
		float32(c.R) / 255,
		float32(c.G) / 255,
		float32(c.B) / 255,
		float32(c.A) / 255,
	}
}
