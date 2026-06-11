package main

import (
	"fmt"
	"math"
	"strconv"

	rl "github.com/gen2brain/raylib-go/raylib"
	ecs "github.com/mlange-42/ark/ecs"
)

type System interface {
	Initialize(*Game)
	Update(*Game)
}

type InputSystem struct{}

func (system *InputSystem) Initialize(game *Game) {}

func (system *InputSystem) Update(game *Game) {
	if !rl.IsMouseButtonPressed(rl.MouseButtonLeft) {
		return
	}

	ray := rl.GetScreenToWorldRay(rl.GetMousePosition(), game.camera)
	if point, ok := intersectRayGroundPlane(ray); ok {
		gridX := int(math.Floor(float64(point.X)))
		gridZ := int(math.Floor(float64(point.Z)))

		fmt.Printf(
			"click world=(%.2f, %.2f, %.2f) grid=(%d, %d)\n",
			point.X,
			point.Y,
			point.Z,
			gridX,
			gridZ,
		)
		game.grid.PlaceEntity(gridX, gridZ, game.models["cube"], baseCubeColor)
		return
	}

	fmt.Println("click missed the ground plane")
}

type CameraSystem struct{}

func (system *CameraSystem) Initialize(game *Game) {}

func (system *CameraSystem) Update(game *Game) {
	camera := &game.camera
	// cameraVector := rl.Vector3Subtract(camera.Target, camera.Position)
	// fmt.Printf(
	// 	"camera pos=(%.2f, %.2f, %.2f) target=(%.2f, %.2f, %.2f) vector=(%.2f, %.2f, %.2f)\n",
	// 	camera.Position.X,
	// 	camera.Position.Y,
	// 	camera.Position.Z,
	// 	camera.Target.X,
	// 	camera.Target.Y,
	// 	camera.Target.Z,
	// 	cameraVector.X,
	// 	cameraVector.Y,
	// 	cameraVector.Z,
	// )

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

type RenderSystem3D struct {
	filter *ecs.Filter2[Position3, Renderable]
}

func (system *RenderSystem3D) Initialize(game *Game) {
	system.filter = ecs.NewFilter2[Position3, Renderable](game.world)
}

func (system *RenderSystem3D) Update(game *Game) {
	rl.BeginMode3D(game.camera)
	system.drawCoordinateSystem()
	system.renderModels()
	rl.EndMode3D()
}

func (system *RenderSystem3D) drawCoordinateSystem() {
	origin := rl.Vector3Zero()
	rl.DrawLine3D(origin, rl.NewVector3(axisLength, 0, 0), rl.Red)
	rl.DrawLine3D(origin, rl.NewVector3(0, axisLength, 0), rl.Green)
	rl.DrawLine3D(origin, rl.NewVector3(0, 0, axisLength), rl.Blue)
}

func (system *RenderSystem3D) renderModels() {
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

type SpawnerSystem struct {
	spawnerMapper      *ecs.Map3[Position3, Renderable, Spawner]
	spawnerFilter      *ecs.Filter2[Position3, Spawner]
	enemyMapper        *ecs.Map3[Position3, Renderable, Enemy]
	movementGoalMapper *ecs.Map1[MovementGoal]
}

func (system *SpawnerSystem) Initialize(game *Game) {
	system.spawnerMapper = ecs.NewMap3[Position3, Renderable, Spawner](game.world)
	system.spawnerFilter = ecs.NewFilter2[Position3, Spawner](game.world)
	system.enemyMapper = ecs.NewMap3[Position3, Renderable, Enemy](game.world)
	system.movementGoalMapper = ecs.NewMap1[MovementGoal](game.world)
	spawnerModel := game.models["spawner"]
	for _, position := range spawnerGridPositions() {
		system.spawnerMapper.NewEntity(
			&Position3{
				X: float32(position.x) + gridCellCenter,
				Y: spawnerY,
				Z: float32(position.z) + gridCellCenter,
			},
			&Renderable{
				model:             spawnerModel,
				scale:             1.0,
				tint:              rl.White,
				shaderTintEnabled: false,
			},
			&Spawner{},
		)
	}
}

func (system *SpawnerSystem) Update(game *Game) {
	if (game.tick+1)%100 != 0 {
		return
	}

	spawnPositions := make([]Position3, 0, 4)
	query := system.spawnerFilter.Query()
	for query.Next() {
		position, _ := query.Get()
		spawnPositions = append(spawnPositions, *position)
	}
	query.Close()

	for _, spawnPosition := range spawnPositions {
		gridX := int(spawnPosition.X)
		gridZ := int(spawnPosition.Z)

		enemyEntity := system.enemyMapper.NewEntity(
			&Position3{
				X: spawnPosition.X,
				Y: spawnPosition.Y,
				Z: spawnPosition.Z,
			},
			&Renderable{
				model:             game.models["miniMob"],
				scale:             1.0,
				tint:              rl.White,
				shaderTintEnabled: false,
			},
			&Enemy{},
		)

		if goalX, goalZ, ok := game.grid.NextLowerDistanceCell(gridX, gridZ); ok {
			system.movementGoalMapper.Add(enemyEntity, &MovementGoal{
				nextGridX: goalX,
				nextGridY: goalZ,
			})
		}
	}
}

func (system *SpawnerSystem) spawnEntity() {

}

type GridDistanceDebugRenderSystem struct{}

func (system *GridDistanceDebugRenderSystem) Initialize(game *Game) {}

func (system *GridDistanceDebugRenderSystem) Update(game *Game) {
	if !debugShowGridDistances {
		return
	}

	for z := 0; z < game.grid.Length; z++ {
		for x := 0; x < game.grid.Width; x++ {
			cell, ok := game.grid.Cell(x, z)
			if !ok {
				continue
			}

			worldPosition := rl.NewVector3(
				float32(x)+gridCellCenter,
				gridDistanceLabelY,
				float32(z)+gridCellCenter,
			)
			screenPosition := rl.GetWorldToScreen(worldPosition, game.camera)
			label := strconv.Itoa(cell.distance)
			labelWidth := rl.MeasureText(label, gridDistanceLabelSize)

			drawX := int32(screenPosition.X) - labelWidth/2
			drawY := int32(screenPosition.Y) - gridDistanceLabelOffset

			rl.DrawText(label, drawX+1, drawY+1, gridDistanceLabelSize, rl.Black)
			rl.DrawText(label, drawX, drawY, gridDistanceLabelSize, rl.White)
		}
	}
}
