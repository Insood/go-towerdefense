package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	ecs "github.com/mlange-42/ark/ecs"
)

type HealthBarOverlaySystem struct {
	filter *ecs.Filter2[Position3, Health]
}

func (system *HealthBarOverlaySystem) Initialize(game *Game) {
	system.filter = ecs.NewFilter2[Position3, Health](game.world)
}

func (system *HealthBarOverlaySystem) Update(game *Game) {
	query := system.filter.Query()
	defer query.Close()

	for query.Next() {
		position, health := query.Get()
		if health.max <= 0 || health.current >= health.max {
			continue
		}

		system.drawHealthBar(game, *position, *health)
	}
}

func (system *HealthBarOverlaySystem) drawHealthBar(game *Game, position Position3, health Health) {
	ratio := clampFloat32(health.current/health.max, 0, 1)

	worldPosition := rl.NewVector3(position.X, position.Y+healthBarWorldYOffset, position.Z)
	screenPosition := rl.GetWorldToScreen(worldPosition, game.camera)

	barX := int32(screenPosition.X - healthBarWidth/2)
	barY := int32(screenPosition.Y - healthBarHeight/2)
	barWidth := int32(healthBarWidth)
	barHeight := int32(healthBarHeight)

	rl.DrawRectangle(barX, barY, barWidth, barHeight, healthBarBackTint)
	rl.DrawRectangle(barX, barY, barWidth, 1, healthBarBorderTint)
	rl.DrawRectangle(barX, barY+barHeight-1, barWidth, 1, healthBarBorderTint)
	rl.DrawRectangle(barX, barY, 1, barHeight, healthBarBorderTint)
	rl.DrawRectangle(barX+barWidth-1, barY, 1, barHeight, healthBarBorderTint)

	innerX := barX + int32(healthBarBorderWidth)
	innerY := barY + int32(healthBarBorderWidth)
	innerWidth := barWidth - int32(healthBarBorderWidth*2)
	innerHeight := barHeight - int32(healthBarBorderWidth*2)
	fillWidth := int32(float32(innerWidth) * ratio)
	if fillWidth > 0 && innerHeight > 0 {
		rl.DrawRectangle(innerX, innerY, fillWidth, innerHeight, healthBarFillTint(ratio))
	}
}

func healthBarFillTint(ratio float32) rl.Color {
	switch {
	case ratio > 0.66:
		return rl.NewColor(80, 210, 95, 255)
	case ratio > 0.33:
		return rl.NewColor(220, 180, 70, 255)
	default:
		return rl.NewColor(225, 80, 70, 255)
	}
}
