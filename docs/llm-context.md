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

## Core Invariants

- Grid width maps to `X`.
- Grid length maps to `Z`.
- The scene uses a right-handed coordinate system.
- The spire occupies the center cell.
- Spawners occupy the edge-center cells.
- Border tiles are no-build tiles.
- BFS pathing must respect occupied cells.
- A placement must not block every path origin.

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
- Goal selection and pathing should happen before movement.
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

