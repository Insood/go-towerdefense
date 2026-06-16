package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	ecs "github.com/mlange-42/ark/ecs"
)

type DebugRender3DSystem struct {
	enemyFilter *ecs.Filter2[Position3, Enemy]
}

func (system *DebugRender3DSystem) Initialize(game *Game) {
	system.enemyFilter = ecs.NewFilter2[Position3, Enemy](game.world)
}

func (system *DebugRender3DSystem) Update(game *Game) {
	if !debugEnabled {
		return
	}

	query := system.enemyFilter.Query()
	defer query.Close()
	rl.BeginMode3D(game.camera)
	for query.Next() {
		position, _ := query.Get()
		system.drawCoordinateSystemAt(rl.Vector3(*position))
	}
	rl.EndMode3D()
}

func (system *DebugRender3DSystem) drawCoordinateSystemAt(origin rl.Vector3) {
	rl.DrawLine3D(origin, rl.Vector3Add(origin, rl.NewVector3(axisLength, 0, 0)), rl.Red)
	rl.DrawLine3D(origin, rl.Vector3Add(origin, rl.NewVector3(0, axisLength, 0)), rl.Green)
	rl.DrawLine3D(origin, rl.Vector3Add(origin, rl.NewVector3(0, 0, axisLength)), rl.Blue)
}
