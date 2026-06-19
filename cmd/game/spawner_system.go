package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	ecs "github.com/mlange-42/ark/ecs"
)

type SpawnerSystem struct {
	spawnerMapper *ecs.Map3[Position3, Renderable, Spawner]
	spawnerFilter *ecs.Filter2[Position3, Spawner]
	enemyMapper   *ecs.Map6[Position3, Renderable, Enemy, MovementGoal, HasMovement, Health]
}

func (system *SpawnerSystem) Initialize(game *Game) {
	system.spawnerMapper = ecs.NewMap3[Position3, Renderable, Spawner](game.world)
	system.spawnerFilter = ecs.NewFilter2[Position3, Spawner](game.world)
	system.enemyMapper = ecs.NewMap6[Position3, Renderable, Enemy, MovementGoal, HasMovement, Health](game.world)
	spawnerModel := game.assets.Model("spawner")
	for _, position := range spawnerGridPositions() {
		system.spawnerMapper.NewEntity(
			&Position3{
				X: float32(position.X) + gridCellCenter,
				Y: spawnerY,
				Z: float32(position.Z) + gridCellCenter,
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
	defer query.Close()

	for query.Next() {
		position, _ := query.Get()
		spawnPositions = append(spawnPositions, *position)
	}

	for _, spawnPosition := range spawnPositions {
		gridX := int(spawnPosition.X)
		gridZ := int(spawnPosition.Z)

		system.enemyMapper.NewEntity(
			&Position3{
				X: spawnPosition.X,
				Y: spawnPosition.Y,
				Z: spawnPosition.Z,
			},
			&Renderable{
				model:             game.assets.Model("miniMob"),
				scale:             1.0,
				tint:              rl.White,
				shaderTintEnabled: false,
			},
			&Enemy{},
			&MovementGoal{
				nextGridX: gridX,
				nextGridZ: gridZ,
			},
			&HasMovement{
				speed: float32(enemySpeed),
			},
			&Health{
				current: enemyMaxHealth,
				max:     enemyMaxHealth,
			},
		)
	}
}
