# LLM Context

This repo is a small 3D tower defense game built in Go and raylib. If you are an LLM looking for the right place to implement a feature, start here.

## Start Here

Read these files first:

- [README.md](../README.md) for the project summary and current TODOs
- [docs/architecture.md](architecture.md) for runtime flow and ownership rules
- [internal/gamegrid/grid.go](../internal/gamegrid/grid.go) for grid, buildability, and BFS pathing
- [internal/gamegrid/grid_test.go](../internal/gamegrid/grid_test.go) for the intended grid behavior
- [cmd/game/game.go](../cmd/game/game.go) for game boot, system order, and asset lifecycle

## What This Repo Is

- The game uses an ECS architecture.
- Gameplay code should stay small and easy to follow.
- Grid and pathing rules live in `internal/gamegrid`.
- ECS components live in `cmd/game/components.go`.
- One system should generally stay in one file under `cmd/game`.
- Rendering changes belong in the render systems unless a separate pass is clearly needed.

## When You Generate Code

- Update or add tests for the behavior you changed.
- Update the repo docs if the change affects gameplay rules, movement, pathing, rendering order, or system ownership.
- Keep the docs aligned with the implementation, especially when the code is intentionally simple or uses a heuristic.
- Prefer leaving a short explanation in the docs when the behavior might not be obvious from the code alone.

## Core Invariants

- Grid width maps to `X`.
- Grid length maps to `Z`.
- The scene uses a right-handed coordinate system.
- The spire occupies the center cell.
- Spawners occupy the edge-center cells.
- Border tiles are no-build tiles.
- BFS pathing must respect occupied cells.
- A placement must not block every path origin.

## Current Enemy Movement Spec

Enemy movement is intentionally simple and mostly driven by precomputed waypoints.

### How waypoint paths are built

- When an enemy spawns, we compute a BFS path from its spawn cell to the center.
- When a tower is successfully placed, we repath every enemy using the updated grid.
- Waypoints are built in world space from the enemy's current tile center followed by the BFS tile centers.
- The last waypoint is not the exact center cell. It lands on the edge of the center approach so enemies stop short of the spire instead of entering the center cell.

### Why the current shape exists

- The current shape keeps movement readable and avoids visually obvious diagonal shortcuts.
- The code favors a small, understandable implementation over a more general route system.
- Repathing is done from the enemy's current tile so tower placement can update movement without needing persistent graph state on the entity.

### How movement consumes waypoints

- `WaypointSystem` reads the current waypoint index and moves the entity toward that target.
- When the entity gets close enough to the current target, the index advances.
- The `MoveSpeed` component stores per-enemy speed so movement does not depend on a global constant.
- `WaypointPath.distanceToGoal` stores the remaining distance along the current waypoint chain and is reduced as the enemy moves.
- If the path is exhausted, the enemy transitions to `ReachedGoal`.

### Current implementation note

- The waypoint builder includes the center of the tile the entity currently occupies as the first waypoint.
- It may skip that first waypoint if the entity is already close enough to the next waypoint, which keeps the motion from doubling back unnecessarily.
- This is a visual rule, not a full navigation guarantee. If the motion still looks acceptable, prefer keeping the implementation simple.
- `WaypointPath.index` is still worth keeping because it tells movement which waypoint is currently active; `distanceToGoal` is a separate remaining-distance metric, not a replacement for the active target pointer.

## Where To Make Changes

Use this as the default change map when adding a feature:

| Feature type | Usually change these files |
| --- | --- |
| New gameplay rules | `internal/gamegrid/*`, `cmd/game/config.go` |
| New ECS data | `cmd/game/components.go` |
| New gameplay system | `cmd/game/*_system.go` |
| New rendering behavior | `cmd/game/render_system.go` or a dedicated render system file |
| New debug overlay | `cmd/game/debug_render_overlay_system.go` |
| New 3D debug drawing | `cmd/game/debug_render_3d_system.go` |
| New asset loading or cleanup | `cmd/game/game.go`, especially `InitializeGame()` and `UnloadAssets()` |
| Shared math or helper code | `cmd/game/utils.go` |
| Grid/pathing changes | `internal/gamegrid/grid.go` and `internal/gamegrid/grid_test.go` |
| Enemy waypoint changes | `cmd/game/waypoint_system.go`, `cmd/game/game.go`, `cmd/game/utils.go` |

## Implementation Pattern

When adding a feature, prefer this order:

1. Update or add tests for the intended behavior.
2. Put shared rules in the smallest file that owns them.
3. Add or update the relevant ECS component or system.
4. Wire the system into `cmd/game/game.go` in the correct update order.
5. Update the render path only if the feature has a visible effect.

## System Order Matters

The update order in `cmd/game/game.go` is intentional. If you change gameplay timing, make sure the new system runs before or after the systems it depends on.

Common examples:

- Input should run before gameplay reacts to it.
- Waypoint selection and pathing should happen before inertia moves entities.
- Movement should happen before rendering.
- Debug drawing should happen after the world is drawn.

## Things Not To Do

- Do not move core gameplay logic into `main.go`.
- Do not add new packages unless the repo clearly needs them.
- Do not change grid or pathing semantics without updating tests and docs.
- Do not add new asset ownership without updating `UnloadAssets()`.
- Do not make the docs vague if the code can express the rule more directly.

## Good Handoff Rule

If you are about to implement a new feature and are unsure where it belongs, check:

1. `internal/gamegrid` for rules about occupancy, pathing, and buildability.
2. `cmd/game` for ECS data, systems, rendering, and assets.
3. `docs/architecture.md` for update order and repo conventions.
