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

func gridCoordFromPosition(position Position3) gamegrid.GridCoord {
	return gamegrid.GridCoord{
		X: int(math.Floor(float64(position.X))),
		Z: int(math.Floor(float64(position.Z))),
	}
}

func worldPositionForGridCoord(coord gamegrid.GridCoord) rl.Vector3 {
	return rl.NewVector3(float32(coord.X)+gridCellCenter, 0, float32(coord.Z)+gridCellCenter)
}

func buildWaypointPathFromPosition(start rl.Vector3, path []gamegrid.GridCoord) ([]rl.Vector3, int) {
	currentTileCenter := worldPositionForGridCoord(gridCoordFromPosition(Position3(start)))
	waypoints := []rl.Vector3{currentTileCenter}
	if len(path) == 0 {
		return waypoints, 0
	}

	if len(path) > 1 {
		for i := 1; i < len(path)-1; i++ {
			waypoints = append(waypoints, worldPositionForGridCoord(path[i]))
		}

		prev := worldPositionForGridCoord(path[len(path)-2])
		center := worldPositionForGridCoord(path[len(path)-1])
		waypoints = append(waypoints, rl.Vector3Scale(rl.Vector3Add(prev, center), 0.5))
	}

	startIndex := 0
	if len(waypoints) > 1 && manhattanDistance2D(start, waypoints[1]) < 1.0 {
		startIndex = 1
	}

	return waypoints, startIndex
}

func waypointPathDistanceToGoal(start rl.Vector3, waypoints []rl.Vector3, startIndex int) float32 {
	if startIndex < 0 || startIndex >= len(waypoints) {
		return 0
	}

	distance := rl.Vector3Distance(start, waypoints[startIndex])
	for i := startIndex; i < len(waypoints)-1; i++ {
		distance += rl.Vector3Distance(waypoints[i], waypoints[i+1])
	}

	return distance
}

func horizontalDistance(a, b rl.Vector3) float32 {
	dx := a.X - b.X
	dz := a.Z - b.Z
	return float32(math.Sqrt(float64(dx*dx + dz*dz)))
}

func pointInsideHitBox(point, center rl.Vector3, hitBox HitBox) bool {
	halfSize := rl.Vector3Scale(hitBox.size, 0.5)
	return math.Abs(float64(point.X-center.X)) <= float64(halfSize.X) &&
		math.Abs(float64(point.Y-center.Y)) <= float64(halfSize.Y) &&
		math.Abs(float64(point.Z-center.Z)) <= float64(halfSize.Z)
}

func manhattanDistance2D(a, b rl.Vector3) float32 {
	return float32(math.Abs(float64(a.X-b.X)) + math.Abs(float64(a.Z-b.Z)))
}
