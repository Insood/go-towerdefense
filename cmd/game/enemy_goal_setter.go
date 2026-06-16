package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	ecs "github.com/mlange-42/ark/ecs"
)

type EnemyGoalSetter struct {
	filter              *ecs.Filter3[Position3, MovementGoal, Enemy]
	reachedGoalExchange *ecs.Exchange1[ReachedGoal]
}

func (system *EnemyGoalSetter) Initialize(game *Game) {
	system.filter = ecs.NewFilter3[Position3, MovementGoal, Enemy](game.world)
	system.reachedGoalExchange = ecs.NewExchange1[ReachedGoal](game.world).Removes(ecs.C[MovementGoal]())
}

func (system *EnemyGoalSetter) Update(game *Game) {
	entitiesToTransition := make([]ecs.Entity, 0)

	query := system.filter.Query()
	defer query.Close()

	// Enemies are trying to make it to center spire location
	centerPosition := rl.NewVector3(
		float32(gridCenterX)+gridCellCenter,
		0,
		float32(gridCenterZ)+gridCellCenter,
	)

	for query.Next() {
		position, movementGoal, _ := query.Get()

		// Once enemies get close enough to their final target, then we can stop moving them
		if movementGoal.nextGridX == gridCenterX && movementGoal.nextGridZ == gridCenterZ {
			if rl.Vector3Distance(rl.Vector3(*position), centerPosition) <= enemyReachedGoalDelta {
				entitiesToTransition = append(entitiesToTransition, query.Entity())
			}
			continue
		}

		// Most of the time is spent here - waiting until the enemy gets close enough to the
		// next waypoint
		goalPosition := rl.NewVector3(
			float32(movementGoal.nextGridX)+gridCellCenter,
			position.Y,
			float32(movementGoal.nextGridZ)+gridCellCenter,
		)
		if rl.Vector3Distance(rl.Vector3(*position), goalPosition) > enemyGoalDelta {
			continue
		}

		// Once we're here, we're near the waypoint - so can set a new one
		candidates := game.grid.NextLowerDistanceCells(movementGoal.nextGridX, movementGoal.nextGridZ)
		if len(candidates) == 0 {
			// In the future, we can do something like detonate the enemy entity and damage nearby buildings
			// Do not implement just yet
			continue
		}

		// If there are multiple path with equal distances, pick a path at random
		next := candidates[rng.Intn(len(candidates))]
		movementGoal.nextGridX = next.X
		movementGoal.nextGridZ = next.Z
	}

	for _, entity := range entitiesToTransition {
		system.reachedGoalExchange.Exchange(entity, &ReachedGoal{})
	}
}
