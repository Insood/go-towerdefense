# Agent Notes

This repo is for a simple 3D tower defense game built in Go with Raylib.

For the higher-level layout and runtime flow, see [docs/architecture.md](docs/architecture.md).

## Development Rules

- Keep the codebase small and understandable.
- Prefer straightforward implementations over clever abstractions.
- Build gameplay in small, playable slices.
- Keep 3D rendering and gameplay systems separated where practical.
- Prefer clear, testable code.
- Avoid adding dependencies unless they solve a real problem.
- Preserve existing user changes unless explicitly asked to modify them.
- If a change affects gameplay structure, explain the tradeoff briefly.
- When you generate code, update or add tests for the behavior and update docs if the change affects gameplay rules, movement, pathing, rendering order, or system ownership.

## Change Placement

- Put gameplay rules and grid behavior in `internal/gamegrid`.
- Put ECS components in `cmd/game/components.go`.
- Put one system per file in `cmd/game/*_system.go`.
- Put rendering changes in the render system unless the feature clearly needs a separate pass.
- Put constants and tuning values in `cmd/game/config.go`.
- Put shared math and utility helpers in `cmd/game/utils.go`.
- Put asset loading and cleanup in `cmd/game/game.go`.
- Keep `main.go` as small as possible.

## Behavior Invariants

- Grid width is `X`.
- Grid length is `Z`.
- The scene uses a right-handed coordinate system.
- The spire occupies the center cell.
- Spawners occupy the edge-center cells.
- Border tiles are no-build tiles.
- BFS pathing must respect occupied cells.

## Do Not

- Do not move core gameplay logic into `main.go`.
- Do not add new dependencies unless they solve a real problem.
- Do not split a system into extra packages unless the repo clearly needs it.
- Do not introduce new asset ownership without updating `UnloadAssets()`.
- Do not change grid or pathing semantics without updating the docs and tests.
