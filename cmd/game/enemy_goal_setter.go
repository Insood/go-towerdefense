package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	ecs "github.com/mlange-42/ark/ecs"
)

type EnemyGoalSetter struct {
	filter             *ecs.Filter3[Position3, MovementGoal, Enemy]
	movementGoalMapper *ecs.Map1[MovementGoal]
}

func (system *EnemyGoalSetter) Initialize(game *Game) {
	system.filter = ecs.NewFilter3[Position3, MovementGoal, Enemy](game.world)
	system.movementGoalMapper = ecs.NewMap1[MovementGoal](game.world)
}

func (system *EnemyGoalSetter) Update(game *Game) {
	entitiesToClear := make([]ecs.Entity, 0)

	query := system.filter.Query()
	for query.Next() {
		position, movementGoal, _ := query.Get()

		goalPosition := rl.NewVector3(
			float32(movementGoal.nextGridX)+gridCellCenter,
			position.Y,
			float32(movementGoal.nextGridY)+gridCellCenter,
		)
		if rl.Vector3Distance(*position, goalPosition) > enemyGoalDelta {
			continue
		}

		nextGridX, nextGridY, ok := game.grid.NextLowerDistanceCell(movementGoal.nextGridX, movementGoal.nextGridY)
		if !ok || (nextGridX == gridCenterX && nextGridY == gridCenterZ) {
			entitiesToClear = append(entitiesToClear, query.Entity())
			continue
		}

		movementGoal.nextGridX = nextGridX
		movementGoal.nextGridY = nextGridY
	}
	query.Close()

	for _, entity := range entitiesToClear {
		system.movementGoalMapper.Remove(entity)
	}
}
