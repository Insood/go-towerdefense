package main

import (
	"image/color"
	"math"
	"testing"

	gamegrid "go-towerdefense/internal/gamegrid"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func TestCameraMoveOnGround(t *testing.T) {
	got := cameraMoveOnGround(2, -3, 4)
	want := rl.NewVector3(8, 0, -12)

	assertVector3Equal(t, got, want)
}

func TestColorToVec4(t *testing.T) {
	got := colorToVec4(color.RGBA{R: 255, G: 128, B: 0, A: 64})

	assertFloat32ApproxEqual(t, got[0], 1.0)
	assertFloat32ApproxEqual(t, got[1], 128.0/255.0)
	assertFloat32ApproxEqual(t, got[2], 0.0)
	assertFloat32ApproxEqual(t, got[3], 64.0/255.0)
}

func TestClampFloat32(t *testing.T) {
	if got := clampFloat32(5, 1, 10); got != 5 {
		t.Fatalf("clampFloat32 = %v, want 5", got)
	}
	if got := clampFloat32(-1, 1, 10); got != 1 {
		t.Fatalf("clampFloat32 = %v, want 1", got)
	}
	if got := clampFloat32(20, 1, 10); got != 10 {
		t.Fatalf("clampFloat32 = %v, want 10", got)
	}
}

func TestLerpFloat32(t *testing.T) {
	if got := lerpFloat32(10, 20, 0.25); got != 12.5 {
		t.Fatalf("lerpFloat32 = %v, want 12.5", got)
	}
}

func TestLerpColorRGBA(t *testing.T) {
	got := lerpColorRGBA(color.RGBA{R: 0, G: 10, B: 20, A: 30}, color.RGBA{R: 100, G: 110, B: 120, A: 130}, 0.5)
	want := color.RGBA{R: 50, G: 60, B: 70, A: 80}

	if got != want {
		t.Fatalf("lerpColorRGBA = %#v, want %#v", got, want)
	}
}

func TestIntersectRayGroundPlaneHitsPlane(t *testing.T) {
	ray := rl.Ray{
		Position:  rl.NewVector3(1, 5, 2),
		Direction: rl.NewVector3(0, -1, 0),
	}

	got, ok := intersectRayGroundPlane(ray)
	if !ok {
		t.Fatal("intersectRayGroundPlane returned no hit, want hit")
	}

	assertVector3Equal(t, got, rl.NewVector3(1, groundPlaneY, 2))
}

func TestIntersectRayGroundPlaneRejectsParallelRay(t *testing.T) {
	ray := rl.Ray{
		Position:  rl.NewVector3(1, 5, 2),
		Direction: rl.NewVector3(1, 0, 0),
	}

	if _, ok := intersectRayGroundPlane(ray); ok {
		t.Fatal("intersectRayGroundPlane returned hit for parallel ray")
	}
}

func TestIntersectRayGroundPlaneRejectsBehindOrigin(t *testing.T) {
	ray := rl.Ray{
		Position:  rl.NewVector3(1, -1, 2),
		Direction: rl.NewVector3(0, -1, 0),
	}

	if _, ok := intersectRayGroundPlane(ray); ok {
		t.Fatal("intersectRayGroundPlane returned hit behind the ray origin")
	}
}

func TestSpawnerGridPositions(t *testing.T) {
	got := spawnerGridPositions()
	want := []gamegrid.GridCoord{
		{X: gridCenterX, Z: gridTopRow},
		{X: gridCenterX, Z: gridBottomRow},
		{X: gridLeftCol, Z: gridCenterZ},
		{X: gridRightCol, Z: gridCenterZ},
	}

	if len(got) != len(want) {
		t.Fatalf("len(spawnerGridPositions) = %d, want %d", len(got), len(want))
	}

	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("spawnerGridPositions[%d] = %#v, want %#v", i, got[i], want[i])
		}
	}
}

func TestGridCoordFromPosition(t *testing.T) {
	got := gridCoordFromPosition(Position3{X: 2.9, Y: 0, Z: 4.1})
	want := gamegrid.GridCoord{X: 2, Z: 4}

	if got != want {
		t.Fatalf("gridCoordFromPosition = %#v, want %#v", got, want)
	}
}

func TestWorldPositionForGridCoord(t *testing.T) {
	got := worldPositionForGridCoord(gamegrid.GridCoord{X: 2, Z: 4})
	want := rl.NewVector3(2+gridCellCenter, 0, 4+gridCellCenter)

	assertVector3Equal(t, got, want)
}

func TestBuildWaypointPathFromPositionReturnsStartWhenPathIsEmpty(t *testing.T) {
	start := rl.NewVector3(2.5, 0, 4.5)

	path, startIndex := buildWaypointPathFromPosition(start, nil)
	if len(path) != 1 {
		t.Fatalf("path length = %d, want 1", len(path))
	}
	if startIndex != 0 {
		t.Fatalf("start index = %d, want 0", startIndex)
	}

	assertVector3Equal(t, path[0], start)
}

func TestBuildWaypointPathFromPositionReturnsStartAndFinalEdgeForSingleStepPath(t *testing.T) {
	start := rl.NewVector3(2.5, 0, 3.5)
	path := []gamegrid.GridCoord{
		{X: 2, Z: 3},
		{X: 2, Z: 4},
	}

	waypoints, startIndex := buildWaypointPathFromPosition(start, path)
	if len(waypoints) != 2 {
		t.Fatalf("path length = %d, want 2", len(waypoints))
	}
	if startIndex != 1 {
		t.Fatalf("start index = %d, want 1", startIndex)
	}

	assertVector3Equal(t, waypoints[0], start)
	assertVector3Equal(t, waypoints[1], rl.NewVector3(2.5, 0, 4.0))
}

func TestBuildWaypointPathFromPositionIncludesIntermediateCellsAndFinalEdge(t *testing.T) {
	start := rl.NewVector3(1.5, 0, 1.5)
	path := []gamegrid.GridCoord{
		{X: 1, Z: 1},
		{X: 1, Z: 2},
		{X: 1, Z: 3},
	}

	waypoints, startIndex := buildWaypointPathFromPosition(start, path)
	if len(waypoints) != 3 {
		t.Fatalf("path length = %d, want 3", len(waypoints))
	}
	if startIndex != 0 {
		t.Fatalf("start index = %d, want 0", startIndex)
	}

	assertVector3Equal(t, waypoints[0], start)
	assertVector3Equal(t, waypoints[1], rl.NewVector3(1.5, 0, 2.5))
	assertVector3Equal(t, waypoints[2], rl.NewVector3(1.5, 0, 3.0))
}

func TestBuildWaypointPathFromPositionSkipsTileCenterWhenAlreadyCloseToSecondWaypoint(t *testing.T) {
	start := rl.NewVector3(1.6, 0, 1.5)
	path := []gamegrid.GridCoord{
		{X: 1, Z: 1},
		{X: 2, Z: 1},
		{X: 3, Z: 1},
	}

	waypoints, startIndex := buildWaypointPathFromPosition(start, path)
	if startIndex != 1 {
		t.Fatalf("start index = %d, want 1", startIndex)
	}

	assertVector3Equal(t, waypoints[0], rl.NewVector3(1.5, 0, 1.5))
	assertVector3Equal(t, waypoints[1], rl.NewVector3(2.5, 0, 1.5))
}

func TestBuildWaypointPathFromPositionRetainsTileCenterWhenFarFromSecondWaypoint(t *testing.T) {
	start := rl.NewVector3(1.2, 0, 1.5)
	path := []gamegrid.GridCoord{
		{X: 1, Z: 1},
		{X: 2, Z: 1},
		{X: 3, Z: 1},
	}

	_, startIndex := buildWaypointPathFromPosition(start, path)
	if startIndex != 0 {
		t.Fatalf("start index = %d, want 0", startIndex)
	}
}

func TestWaypointPathDistanceToGoalFromStartingWaypoint(t *testing.T) {
	start := rl.NewVector3(1.5, 0, 1.5)
	waypoints := []rl.Vector3{
		rl.NewVector3(1.5, 0, 1.5),
		rl.NewVector3(1.5, 0, 2.5),
		rl.NewVector3(1.5, 0, 3.0),
	}

	got := waypointPathDistanceToGoal(start, waypoints, 0)
	assertFloat32ApproxEqual(t, got, 1.5)
}

func TestWaypointPathDistanceToGoalFromSkippedStartingWaypoint(t *testing.T) {
	start := rl.NewVector3(1.6, 0, 1.5)
	waypoints := []rl.Vector3{
		rl.NewVector3(1.5, 0, 1.5),
		rl.NewVector3(2.5, 0, 1.5),
		rl.NewVector3(3.0, 0, 1.5),
	}

	got := waypointPathDistanceToGoal(start, waypoints, 1)
	assertFloat32ApproxEqual(t, got, 1.4)
}

func TestWaypointPathDistanceToGoalReturnsZeroWhenIndexIsOutOfRange(t *testing.T) {
	got := waypointPathDistanceToGoal(rl.NewVector3(1.5, 0, 1.5), nil, 0)
	if got != 0 {
		t.Fatalf("waypointPathDistanceToGoal = %v, want 0", got)
	}
}

func TestManhattanDistance2D(t *testing.T) {
	got := manhattanDistance2D(rl.NewVector3(1, 0, 2), rl.NewVector3(4, 0, -1))
	if got != 6 {
		t.Fatalf("manhattanDistance2D = %v, want 6", got)
	}
}

func assertVector3Equal(t *testing.T, got, want rl.Vector3) {
	t.Helper()

	if got != want {
		t.Fatalf("vector = %#v, want %#v", got, want)
	}
}

func assertFloat32ApproxEqual(t *testing.T, got, want float32) {
	t.Helper()

	if math.Abs(float64(got-want)) > 1e-5 {
		t.Fatalf("float = %v, want %v", got, want)
	}
}
