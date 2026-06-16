package main

import (
	"image/color"
	"math"

	gamegrid "go-towerdefense/internal/gamegrid"

	rl "github.com/gen2brain/raylib-go/raylib"
)

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

func clampFloat32(value, min, max float32) float32 {
	return float32(math.Max(float64(min), math.Min(float64(max), float64(value))))
}

func lerpFloat32(a, b, t float32) float32 {
	return a + (b-a)*t
}

func lerpColorRGBA(a, b color.RGBA, t float32) color.RGBA {
	return color.RGBA{
		R: uint8(lerpFloat32(float32(a.R), float32(b.R), t)),
		G: uint8(lerpFloat32(float32(a.G), float32(b.G), t)),
		B: uint8(lerpFloat32(float32(a.B), float32(b.B), t)),
		A: uint8(lerpFloat32(float32(a.A), float32(b.A), t)),
	}
}

func intersectRayGroundPlane(ray rl.Ray) (rl.Vector3, bool) {
	if ray.Direction.Y > -rayParallelEpsilon && ray.Direction.Y < rayParallelEpsilon {
		return rl.Vector3{}, false
	}

	t := (groundPlaneY - ray.Position.Y) / ray.Direction.Y
	if t < 0 {
		return rl.Vector3{}, false
	}

	return rl.NewVector3(
		ray.Position.X+ray.Direction.X*t,
		groundPlaneY,
		ray.Position.Z+ray.Direction.Z*t,
	), true
}

func spawnerGridPositions() []gamegrid.GridCoord {
	return []gamegrid.GridCoord{
		{X: gridCenterX, Z: gridTopRow},
		{X: gridCenterX, Z: gridBottomRow},
		{X: gridLeftCol, Z: gridCenterZ},
		{X: gridRightCol, Z: gridCenterZ},
	}
}
