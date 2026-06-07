package main

import (
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
	frameStep := cameraPanSpeed * rl.GetFrameTime()

	var moveX float32
	var moveZ float32
	if rl.IsKeyDown(rl.KeyW) {
		moveZ += 1
	}
	if rl.IsKeyDown(rl.KeyS) {
		moveZ -= 1
	}
	if rl.IsKeyDown(rl.KeyD) {
		moveX += 1
	}
	if rl.IsKeyDown(rl.KeyA) {
		moveX -= 1
	}

	pan := cameraMoveOnGround(camera, moveX, moveZ, frameStep)
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

		rl.DrawModel(*renderable.model, *position, renderable.scale, renderable.tint)
	}
}

func cameraMoveOnGround(camera *rl.Camera3D, x, z, distance float32) rl.Vector3 {
	forward := rl.Vector3Subtract(camera.Target, camera.Position)
	forward.Y = 0
	if forward.X == 0 && forward.Z == 0 {
		return rl.Vector3Zero()
	}
	forward = rl.Vector3Normalize(forward)

	right := rl.Vector3CrossProduct(forward, camera.Up)
	right.Y = 0
	if right.X == 0 && right.Z == 0 {
		return rl.Vector3Zero()
	}
	right = rl.Vector3Normalize(right)

	move := rl.Vector3Zero()
	if x != 0 {
		move = rl.Vector3Add(move, rl.Vector3Scale(right, x*distance))
	}
	if z != 0 {
		move = rl.Vector3Add(move, rl.Vector3Scale(forward, z*distance))
	}

	return move
}
