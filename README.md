# Go Tower Defense

Tiny 3D tower defense game built in Go and raylib.

The player defends a central target while enemies spawn at the map edges and move inward. Towers can be placed on buildable grid cells.

## The Game

Run the game with `go run ./cmd/game` or with `make game`

For repository structure and gameplay rules, see [docs/architecture.md](docs/architecture.md).
For LLM-friendly implementation guidance, see [docs/llm-context.md](docs/llm-context.md).

## TODO:
* Create models for spires, spawners, and mobs
* Implement a health system for enemies, walls, and spires
* When enemies reach the spire, they explode and the spire loses health
* Healthbar overlay for towers, spires, and enemies
* Add towers
  * Missile tower (should be easiest to implement - missile can track target and curve to attack)
  * Mortar tower (similar to missile tower, except projectile explodes when it hits the ground and does AOE damage)
  * Gunner tower (shoots in a straight line and when it hits an enemy, causes damage only to that enemy)
    * Ideally would need quad trees for this, but not strictly required since there won't be that many projectiles
