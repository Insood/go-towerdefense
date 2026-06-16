# Go Tower Defense

Tiny 3D tower defense game built in Go and raylib.

The player defends a central target while enemies spawn at the map edges and move inward. Towers can be placed on buildable grid cells.

## The Game

Run the game with `go run ./cmd/game` or with `make game`

For repository structure and gameplay rules, see [docs/architecture.md](docs/architecture.md).

## TODO:
* Create models for spires, spawners, and mobs
* Implement a health system for enemies and for the tower
* When enemies reach the tower, they explode and the tower loses health
* Healthbar overlay for towers, spires, and enemies