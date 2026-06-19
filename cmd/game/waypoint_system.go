package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	ecs "github.com/mlange-42/ark/ecs"
)

type WaypointSystem struct {
	filter              *ecs.Filter5[Position3, WaypointPath, MoveSpeed, Velocity3, Enemy]
	reachedGoalExchange *ecs.Exchange1[ReachedGoal]
}

func (system *WaypointSystem) Initialize(game *Game) {
	system.filter = ecs.NewFilter5[Position3, WaypointPath, MoveSpeed, Velocity3, Enemy](game.world)
	system.reachedGoalExchange = ecs.NewExchange1[ReachedGoal](game.world).Removes(ecs.C[WaypointPath](), ecs.C[Velocity3]())
}

func (system *WaypointSystem) Update(game *Game) {
	deltaTime := rl.GetFrameTime()

	entitiesToTransition := make([]ecs.Entity, 0)
	query := system.filter.Query()
	defer query.Close()

	for query.Next() {
		position, path, moveSpeed, velocity, _ := query.Get()
		if len(path.waypoints) == 0 || path.index >= len(path.waypoints) {
			*velocity = Velocity3{}
			continue
		}

		from := rl.Vector3(*position)
		target := path.waypoints[path.index]
		toTarget := rl.Vector3Subtract(target, from)
		distance := rl.Vector3Length(toTarget)

		// Are we near the next waypoint
		if distance <= enemyWaypointDelta {

			path.index++ // Yes, set the next waypoint

			// Are we at the last waypoint?
			// If so, remove set the ReachedGoal component
			if path.index >= len(path.waypoints) {
				*velocity = Velocity3{}
				entitiesToTransition = append(entitiesToTransition, query.Entity())
				continue
			}

			// Calculate the distance to the next target
			target = path.waypoints[path.index]
			toTarget = rl.Vector3Subtract(target, from)
			distance = rl.Vector3Length(toTarget)
		}

		// Prevent oscilation around the endpoint; if we would overstep the waypoint
		// clamp the velocity so that we hit it exactly
		speed := moveSpeed.value
		if maxStep := speed * deltaTime; maxStep > distance {
			speed = distance / deltaTime
		}

		*velocity = Velocity3(rl.Vector3Scale(rl.Vector3Normalize(toTarget), speed))
	}

	for _, entity := range entitiesToTransition {
		system.reachedGoalExchange.Exchange(entity, &ReachedGoal{})
	}
}
