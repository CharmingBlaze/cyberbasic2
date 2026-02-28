# CyberBasic

[![Go 1.22+](https://img.shields.io/badge/Go-1.22+-00ADD8?logo=go)](https://go.dev/)

A modern BASIC-inspired language and game engine with full 2D/3D graphics, physics, networking, and tooling. The runtime is written in **Go** and ships as a single binary: no C++ build step, no external engine DLL. Script your game in a readable, familiar dialect; the engine handles rendering, physics, audio, and multiplayer behind a consistent API.

**Repository:** [github.com/CharmingBlaze/cyberbasic2](https://github.com/CharmingBlaze/cyberbasic2) — report issues and contribute there.

---

## Why Go: From C++ to a Single-Binary Runtime

CyberBasic’s core was reimplemented in Go to improve the entire development and distribution story.

- **Maintainability:** Go’s straightforward syntax and standard tooling (gofmt, go test, go build) make the compiler, VM, and bindings easy to navigate and refactor. New contributors can follow data flow without fighting template or macro expansion.
- **Readability:** A single language from CLI to VM to bindings reduces context switching. Package boundaries (lexer, parser, codegen, vm, bindings) keep responsibilities clear.
- **Build speed:** `go build` produces one executable. No C++ compile cycles, no linking against engine libs for the default flow. CI and local iteration stay fast.
- **Contributor experience:** Clone, `go build -o cyberbasic .`, run. Dependencies are Go modules; no Raylib or Bullet build required for the default path (raylib-go and Go wrappers for Box2D/Bullet are used). An optional C layer in `engine/` exists for custom builds; the recommended path is pure Go.

The result is a stable, feature-rich engine that is pleasant to work on and easy to ship.

---

## Technical Identity

CyberBasic is a **bytecode-compiled**, **single-process** engine: source is parsed and compiled to a custom bytecode, then executed by a small VM. All graphics, input, audio, and physics are exposed as foreign functions (Go) called from bytecode. The design favors clarity and predictability: one main thread, explicit update/draw phases, and a hybrid loop that can run physics and render queues for you when you define `update(dt)` and `draw()`.

---

## Core Systems (Integrated and Available)

The engine ships with full integrations for the following. Scripts call into them via a uniform, case-insensitive BASIC API.

| Domain | Technology | Capabilities |
|--------|------------|--------------|
| **Graphics** | Raylib (raylib-go) | 2D/3D primitives, textures, text, fonts, images, shaders, render textures, skybox, blend modes, scissor, multiple 2D/3D cameras |
| **2D engine** | Raylib + custom | Layers, parallax, z-order, backgrounds (static/scrolling/tiled), tilemaps (create/load/save/fill), sprites (transform, animation, batching), 2D particle emitters, texture atlas, 2D culling |
| **3D** | Raylib | Models, meshes (procedural and from file), materials, model animation, DrawMesh with matrix, full camera (position/target/up/FOV, move/rotate/roll), raycasting, collision helpers |
| **2D physics** | Box2D | Worlds, bodies, shapes, joints, raycast, collision callbacks, layer-based collision matrix, gravity |
| **3D physics** | Bullet | Worlds, rigid bodies, collision shapes, raycast, terrain heightfield, water buoyancy hooks, step integration |
| **Multiplayer** | TCP (net) | Host, Accept, Connect, Send, Receive, Disconnect; event callbacks; optional multi-window (spawn processes, connect to parent) |
| **Events** | VM + bindings | On KeyDown/KeyPressed, OnWindowUpdate/Draw/Resize/Close/Message, collision handlers, configurable event loop |
| **GUI** | Raylib + raygui | BeginUI/EndUI, Label, Button, Slider, Checkbox, TextBox, ProgressBar, WindowBox, GroupBox; raygui widgets for immediate-mode UI |
| **Game helpers** | game | Tilemaps, particle systems, 2D camera center/follow, high-level game loop helpers |
| **Scene** | scene | CreateScene, LoadScene, SaveScene, AddToScene/RemoveFromScene; 2D scene save/load (layers, backgrounds, sprites, tilemaps, particles, camera) |
| **Terrain** | terrain | Heightmap load/generate, terrain mesh, create/update/draw, sculpt/paint, TerrainGetHeight/GetNormal/Raycast, collision, friction/bounce |
| **Water** | water | WaterCreate, DrawWater, wave params, WaterGetHeight, buoyancy/density/drag hooks |
| **Vegetation** | vegetation | Tree types and placement, grass patches, wind, LOD/instancing, collision radius |
| **Objects** | objects | Object placement, scatter, paint/erase, raycast, get-at |
| **World** | world | WorldSave/WorldLoad, export/import JSON, chunk save/load, streaming (load/unload by radius) |
| **Navigation** | navigation | NavGrid (A*), NavMesh from terrain, obstacles, NavAgent (speed, radius, destination, waypoints) |
| **Indoor** | indoor | Rooms, portals, doors, levers, switches, buttons, triggers (enter/exit/stay), interactables, pickups, light zones, nav by room, save/load indoor state |
| **Procedural** | procedural | Noise (Perlin, fractal, Simplex), scatter (trees, grass, objects), biomes |
| **Standard library** | std | File I/O, JSON, HTTP, strings, math, HELP; valueutil (e.g. truthiness) |
| **Data** | sql | SQLite: OpenDatabase, Exec, Query, parameterized statements, transactions |
| **ECS** | ecs | Entity-component system: world, entities, components, queries |

The compiler supports modules, user subs/functions, event handlers, and coroutines (StartCoroutine, Yield, WaitSeconds). A **hybrid loop** is available: define `update(dt)` and `draw()`, use a game loop with an empty body, and the compiler injects physics step, render queue clear, draw call, and flush so you never call BeginDrawing/EndDrawing or BeginMode2D/EndMode3D yourself. See [Rendering and the game loop](docs/RENDERING_AND_GAME_LOOP.md).

---

## Architecture

```
CyberBasic/
├── compiler/          # Go compiler and VM
│   ├── lexer/         # Tokenizer
│   ├── parser/        # AST
│   ├── vm/            # Bytecode VM, render queues, fibers
│   ├── runtime/       # Game loop, StepFrame
│   └── bindings/      # Foreign API (all subsystems above)
│       ├── raylib/    # Graphics, window, input, audio, fonts, math, 2D layers/camera/backgrounds, 3D, hybrid flush
│       ├── box2d/     # 2D physics
│       ├── bullet/    # 3D physics
│       ├── game/      # Tilemaps, particles, game helpers
│       ├── scene/     # Scenes, 2D save/load
│       ├── net/       # TCP multiplayer
│       ├── ecs/       # Entity-component system
│       ├── terrain/   # Heightmap, terrain mesh, edit, query
│       ├── water/     # Water mesh, shader, buoyancy
│       ├── vegetation/# Trees, grass
│       ├── objects/   # Object placement
│       ├── world/     # Save/load, streaming
│       ├── navigation/# NavGrid, NavMesh, agents
│       ├── indoor/    # Rooms, doors, triggers, interactables
│       ├── procedural/# Noise, scatter, biomes
│       ├── std/       # File, JSON, HTTP, HELP
│       └── sql/       # SQLite
├── engine/            # Optional C wrapper (custom builds)
├── examples/         # BASIC examples
└── main.go            # CLI: compile + run .bas
```

No circular dependencies; the compiler does not depend on gogen. The VM exposes a small set of primitives (stack, globals, foreign calls, fibers, render queues); each binding package registers its commands and participates in the same runtime.

---

## Hello World and First Game

Minimal runnable program:

```bas
PRINT "Hello, CyberBasic!"
```

Window and input (build first; see [Building](#building)):

```bas
InitWindow(800, 600, "My Game")
SetTargetFPS(60)
VAR x = 400
VAR y = 300
WHILE NOT WindowShouldClose()
  IF IsKeyDown(KEY_W) THEN LET y = y - 4
  IF IsKeyDown(KEY_S) THEN LET y = y + 4
  IF IsKeyDown(KEY_A) THEN LET x = x - 4
  IF IsKeyDown(KEY_D) THEN LET x = x + 4
  ClearBackground(20, 20, 30, 255)
  DrawCircle(x, y, 30, 255, 100, 100, 255)
  DrawText("WASD to move", 10, 10, 20, 255, 255, 255, 255)
WEND
CloseWindow()
```

Run the first-game demo in one command:

```bash
./cyberbasic examples/first_game.bas
```

On Windows: `.\cyberbasic.exe examples\first_game.bas`. Helper scripts: `./run_demo.sh` (Unix) or `.\run_demo.ps1` (PowerShell) build and run the demo.

---

## Documentation

| Resource | Description |
|----------|-------------|
| [Documentation Index](docs/DOCUMENTATION_INDEX.md) | Master index of all guides |
| [Rendering and the game loop](docs/RENDERING_AND_GAME_LOOP.md) | Pipeline, hybrid vs manual loop, rule for `draw()` |
| [Roadmap](ROADMAP.md) | Priorities and planned work |
| [Getting Started](docs/GETTING_STARTED.md) | Build, run, first program |
| [Quick Reference](docs/QUICK_REFERENCE.md) | One-page syntax |
| [Language Spec](LANGUAGE_SPEC.md) | Full language reference |
| [2D Graphics](docs/2D_GRAPHICS_GUIDE.md) | 2D rendering, layers, camera |
| [3D Graphics](docs/3D_GRAPHICS_GUIDE.md) | 3D rendering, camera, models |
| [Game Development](docs/GAME_DEVELOPMENT_GUIDE.md) | Game loop, input, physics, ECS |
| [GUI Guide](docs/GUI_GUIDE.md) | Immediate-mode UI |
| [Multiplayer](docs/MULTIPLAYER.md) | TCP client/server |
| [In-process multi-window](docs/MULTI_WINDOW_INPROCESS.md) | Multiple viewports in one process |
| [API Reference](API_REFERENCE.md) | All bindings |
| [Changelog](CHANGELOG.md) | Version history and release notes |

---

## Building

Default build (no C dependencies for normal run):

```bash
go build -o cyberbasic .
./cyberbasic examples/first_game.bas
```

Optional C engine (Raylib + Bullet):

```bash
cd engine && make
```

---

## Limitations and Roadmap

- **Physics:** Some Box2D/Bullet joint and body-property APIs are stubbed; see API_REFERENCE for the full list. Core worlds, bodies, shapes, raycast, and collision are implemented.
- **UI:** Full widget set (Label, Button, Slider, Checkbox, TextBox, etc.) is available; some advanced layout or theme options may be extended.
- **Audio:** Stream callbacks that require C function pointers are not exposed from BASIC; use UpdateAudioStream and similar push APIs.
- **Roadmap:** Working debugger, more complete physics joints, REPL, VSCode extension, and CI are on the roadmap. See [ROADMAP.md](ROADMAP.md).

---

## Contributing

The project is open to contributions. The codebase is structured so that bindings stay in `compiler/bindings/<package>`, the VM and compiler are in `compiler/`, and documentation lives in `docs/` and the root. If you add a command, register it in the appropriate binding package and add an entry to [docs/COMMAND_REFERENCE.md](docs/COMMAND_REFERENCE.md) and [API_REFERENCE.md](API_REFERENCE.md). Run `go build ./...` and the compiler tests before submitting a PR.

CyberBasic is an ambitious, evolving engine: stable in its core systems and actively improving in tooling, physics completeness, and developer experience. We aim to keep the architecture clean, the docs accurate, and the project welcoming to new contributors.
