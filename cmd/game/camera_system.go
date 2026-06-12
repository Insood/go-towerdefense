package main

import rl "github.com/gen2brain/raylib-go/raylib"

type CameraSystem struct{}

func (system *CameraSystem) Initialize(game *Game) {}

func (system *CameraSystem) Update(game *Game) {
	camera := &game.camera

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
