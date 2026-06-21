# Architecture

This repo is a small ECS-based tower defense game built in Go and raylib.

## Directory Layout

- `cmd/game` is the executable game.
- `cmd/game/components.go` defines ECS components.
- `cmd/game/*_system.go` contains one system per file.
- `cmd/game/config.go` stores gameplay constants.
- `cmd/game/utils.go` stores small shared math and rendering helpers.
- `cmd/game/assets` contains runtime assets such as shaders.
- `internal/gamegrid` contains the grid state and BFS pathing rules.

## Runtime Flow

The game loop is intentionally simple:

1. Input system processes mouse and keyboard state.
2. Camera system updates the camera.
3. Hover preview system updates the hovered grid cell.
4. Gameplay systems update spawning, waypoint selection, and movement.
5. Render system draws the 3D scene.
6. Debug systems draw overlays after the world is rendered.

## Game State Rules

- Grid width maps to the X axis.
- Grid length maps to the Z axis.
- The world is right-handed.
- The center cell contains the spire.
- Spawners are placed on the edge-center cells.
- Border cells are no-build cells.
- Occupied cells block BFS pathing.
- Enemies store precomputed waypoint paths that follow BFS distances toward the center.
- Tower placement can trigger enemy path recomputation, but only after placement is accepted.
- Enemy waypoint paths are built to keep visible movement simple and readable:
  - the path starts at the center of the tile the enemy currently occupies
  - the path then follows BFS toward the center
  - the final waypoint lands on the edge of the center approach instead of the exact center
  - when a tower is placed, enemies are repathed from their current tile rather than from a global route state
  - each enemy keeps a waypoint index for the active target and a `distanceToGoal` value for the remaining route length

## Editing Guidance

- If a change affects pathing or buildability, update `internal/gamegrid` first.
- If a change affects render order, update the system list in `cmd/game/game.go`.
- If a change affects an asset lifecycle, update `Game.UnloadAssets()`.
- If a change affects enemy path follow logic, update the waypoint system and enemy path rebuild helpers together.
- If a change affects shared math or color conversion, prefer `cmd/game/utils.go`.
- Prefer one responsibility per system file.
