# Agent Notes

This repo is for a simple 3D tower defense game built in Go with Raylib.

## Goals

- Keep the codebase small and understandable.
- Prefer straightforward implementations over clever abstractions.
- Build gameplay in small, playable slices.

## Harness

- Use `go run ./cmd/harness` to work with the local development harness.
- Keep harness changes minimal unless the user asks for more tooling.

## The Game
- The game code is located in cmd/game

## Development Rules

- Prefer clear, testable code.
- Avoid adding dependencies unless they solve a real problem.
- Preserve existing user changes unless explicitly asked to modify them.
- If a change affects gameplay structure, explain the tradeoff briefly.

## Next Game Work

- Start with a minimal playable loop.
- Add basic map, enemies, towers, and waves before expanding features.
- Keep 3D rendering and gameplay systems separated where practical.
