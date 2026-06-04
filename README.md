# Go Tower Defense

This is a repo for a tiny tower defense game made using Go and raylib.

## The Game

Run the game with `go run ./cmd/game` or with `make game`

## Tiny LLM harness

It gives you:

- a current goal
- tasks grouped under goals
- an active goal that new tasks attach to
- a persistent JSON state file

## Harness Commands

```
make harness
make test

go run ./cmd/harness goal Build the core 3D tower defense prototype
go run ./cmd/harness add Create the map grid
go run ./cmd/harness goal Add basic enemy waves
go run ./cmd/harness switch 0
go run ./cmd/harness add Spawn basic enemy waves
go run ./cmd/harness list
go run ./cmd/harness next
go run ./cmd/harness done 0
go run ./cmd/harness summary
```

`goal` creates a new goal and makes it active.

`switch` moves the active goal to an existing goal by ID.

State is stored in `.tdharness/state.json`.

Run the harness with `go run ./cmd/harness`.

Run the same shortcuts with `make harness` and `make test`.
