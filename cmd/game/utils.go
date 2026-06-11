package main

import (
	"image/color"

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

type gridCoord struct {
	x int
	z int
}

func spawnerGridPositions() []gridCoord {
	return []gridCoord{
		{x: gridCenterX, z: gridTopRow},
		{x: gridCenterX, z: gridBottomRow},
		{x: gridLeftCol, z: gridCenterZ},
		{x: gridRightCol, z: gridCenterZ},
	}
}
