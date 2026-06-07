# Agent Notes

This repo is for a simple 3D tower defense game built in Go with Raylib.

## Development Rules

- Keep the codebase small and understandable.
- Prefer straightforward implementations over clever abstractions.
- Build gameplay in small, playable slices.
- Keep 3D rendering and gameplay systems separated where practical.
- Prefer clear, testable code.
- Avoid adding dependencies unless they solve a real problem.
- Preserve existing user changes unless explicitly asked to modify them.
- If a change affects gameplay structure, explain the tradeoff briefly.

## Task management and development harness

- Use `go run ./cmd/harness` to work with the local development harness.
- Pull open tasks from the harness when requested
- Keep harness changes minimal unless the user asks for more tooling.

## The Game
- The game code is located in cmd/game
