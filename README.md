<p align="center">
  <img src="assets/cyberbasic_logo.png" alt="CyberBASIC 2 Logo" width="400">
</p>

# CyberBASIC 2

[![Go 1.22+](https://img.shields.io/badge/Go-1.22+-00ADD8?logo=go)](https://go.dev/)

**A powerful BASIC-inspired language and game engine** — compile your game in seconds, ship a single binary. CyberBASIC 2 brings the simplicity of classic BASIC to modern game development: readable syntax, minimal boilerplate, and full access to industry-standard graphics and physics.

**Inspired by DarkBASIC PRO.** Built for developers who want to focus on making games, not fighting build systems.

---

## Why CyberBASIC 2?

### The Power of BASIC

BASIC has always been about **clarity and accessibility**. CyberBASIC 2 keeps that spirit: no semicolons, no curly braces, no complex type declarations. Write game logic in a language that reads like plain English. Variables, loops, and functions work the way you expect. A few lines of code can open a window, draw a 3D cube, and handle input.

### A Real Compiler

CyberBASIC 2 is **not a script interpreter**. Your `.bas` source is compiled to optimized bytecode and executed by a fast virtual machine. The compiler supports modules, user-defined functions, event handlers, coroutines, and a hybrid game loop that injects physics and rendering automatically. One binary, zero dependencies.

### Written Entirely in Go

The entire stack — lexer, parser, code generator, VM, and all bindings — is written in **Go**. No C++ build step. No external engine DLLs. `go build -o cyberbasic .` produces a single executable. The result: fast iteration, easy contribution, and a runtime that ships anywhere Go runs.

**Architecture (for contributors):** [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md) — compiler pipeline, **`RegisterAll`** binding order, v2 **module DotObject** API vs legacy flat commands, how to add a binding. **Command coverage matrix + generated inventory:** [docs/COMMAND_COVERAGE.md](docs/COMMAND_COVERAGE.md); after changing bindings run `make foreign-audit` and commit `docs/generated/*`.

---

## Full-Stack Game Development

CyberBASIC 2 integrates **full** implementations of industry-standard libraries. No watered-down subsets — you get the complete APIs.

<table>
<tr>
<td width="25%" align="center">
  <img src="https://box2d.org/images/logo.svg" alt="Box2D" width="120"><br>
  <strong>Box2D</strong><br>
  <small>2D Physics</small>
</td>
<td width="25%" align="center">
  <img src="assets/bullet_logo.png" alt="Bullet Physics" width="120"><br>
  <strong>Bullet Physics</strong><br>
  <small>3D Physics</small>
</td>
<td width="25%" align="center">
  <img src="https://raw.githubusercontent.com/raysan5/raylib/master/logo/raylib_256x256.png" alt="raylib" width="120"><br>
  <strong>raylib</strong><br>
  <small>Graphics & Audio</small>
</td>
</tr>
</table>

- **raylib** — Full 2D/3D graphics, textures, shaders, fonts, audio, input. Window management, cameras, render textures, blend modes.
- **Box2D** — Complete 2D physics: worlds, bodies, shapes, joints, raycast, collision callbacks, layer filtering.
- **Bullet-style 3D physics** — Rigid bodies, raycast, and forces are fully supported in the shipped build. Joints and mesh/terrain require the optional native Bullet backend; see [3D Physics Guide](docs/3D_PHYSICS_GUIDE.md).

### What You Can Build

| Domain | Capabilities |
|--------|--------------|
| **2D Games** | Sprites, tilemaps, layers, parallax, particle systems, 2D physics |
| **3D Games** | Models, materials, animations, 3D physics, terrain, water, vegetation |
| **GUI** | Buttons, sliders, textboxes, progress bars, window boxes — immediate-mode UI |
| **Multiplayer** | TCP client/server, Host/Accept, Connect/Send/Receive, event callbacks |
| **World Building** | Terrain sculpting, object placement, navigation (A*, NavMesh), indoor systems |

---

## Quick Start

**Implicit window** (no `InitWindow`): use `ON UPDATE` / `ON DRAW` and optional `window.*` properties. See `examples/implicit_loop.bas`.

**Classic explicit loop** (unchanged):

```basic
InitWindow(800, 600, "My Game")
SetTargetFPS(60)

mainloop
  ClearBackground(25, 25, 35, 255)
  DrawCircle(400, 300, 50, 255, 100, 100, 255)
  DrawText("Hello, CyberBASIC 2!", 280, 260, 20, 255, 255, 255, 255)
  SYNC
endmain

CloseWindow()
```

Run it:

```bash
go build -o cyberbasic .
./cyberbasic examples/first_game.bas
```

With no arguments, the CLI starts the **REPL**. High-level bindings include `PhysicsHighWorld` / `PhysicsHighDynamicBox`, `AssetsGet`, `InputMapRegister`, and `AudioLoadSound` (see `API_REFERENCE.md` / source under `compiler/bindings/`).

For more control (manual loop, custom update/draw), see [Rendering and the game loop](docs/RENDERING_AND_GAME_LOOP.md).

### Examples

| Example | Description | Command |
|---------|-------------|---------|
| [hello_world.bas](examples/hello_world.bas) | Minimal: prints to console | `./cyberbasic examples/hello_world.bas` |
| [first_game.bas](examples/first_game.bas) | 3D cube + orbit camera | `./cyberbasic examples/first_game.bas` |
| [platformer.bas](examples/platformer.bas) | 2D platformer: WASD, jump | `./cyberbasic examples/platformer.bas` |
| [ui_demo.bas](examples/ui_demo.bas) | Immediate-mode UI: Label, Button | `./cyberbasic examples/ui_demo.bas` |
| [input_debug.bas](examples/input_debug.bas) | Keyboard, mouse, gamepad | `./cyberbasic examples/input_debug.bas` |

**Repository:** [github.com/CharmingBlaze/cyberbasic2](https://github.com/CharmingBlaze/cyberbasic2) — **Downloads:** [GitHub Releases](https://github.com/CharmingBlaze/cyberbasic2/releases) (Windows, macOS, Linux).

---

## Why Go: Single-Binary Runtime

CyberBASIC's core was built in Go for maintainability, readability, and fast builds.

- **Maintainability** — Go's straightforward syntax and tooling (gofmt, go test) make the compiler and VM easy to navigate. New contributors can follow data flow without fighting templates or macros.
- **Build speed** — `go build` produces one executable. No C++ compile cycles, no linking against engine libs. CI and local iteration stay fast.
- **Contributor experience** — Clone, `go build -o cyberbasic .`, run. Dependencies are Go modules; raylib-go and Go wrappers for Box2D/Bullet are used. No separate Raylib or Bullet build required.

---

## Technical Identity

CyberBASIC is **bytecode-compiled** and **single-process**: source is parsed and compiled to custom bytecode, then executed by a small VM. Graphics, input, audio, and physics are exposed as foreign functions (Go) called from bytecode. One main thread, explicit update/draw phases, and a hybrid loop that can run physics and render queues automatically when you define `update(dt)` and `draw()`.

---

## Core Systems

| Domain | Technology | Capabilities |
|--------|------------|--------------|
| **Graphics** | Raylib | 2D/3D primitives, textures, text, fonts, shaders, render textures, skybox, blend modes, multiple cameras |
| **2D engine** | Raylib + custom | Layers, parallax, tilemaps, sprites (animation, batching), particle emitters, texture atlas |
| **3D** | Raylib | Models, meshes, materials, animation, raycasting, full camera control |
| **2D physics** | Box2D | Worlds, bodies, shapes, joints, raycast, collision callbacks, gravity |
| **3D physics** | Bullet-style fallback | Worlds, rigid bodies, raycast, forces (full); joints/mesh/terrain optional — see [3D Physics Guide](docs/3D_PHYSICS_GUIDE.md) |
| **Multiplayer** | TCP (net) | Host, Accept, Connect, Send, Receive; event callbacks |
| **GUI** | raygui | Label, Button, Slider, Checkbox, TextBox, ProgressBar, WindowBox |
| **ECS** | ecs | Entity-component system: world, entities, components, queries |
| **Standard library** | std | File I/O, JSON, HTTP, strings, math, HELP |

---

## Architecture

```
CyberBasic/
├── compiler/          # Go compiler and VM
│   ├── lexer/         # Tokenizer
│   ├── parser/        # AST
│   ├── vm/            # Bytecode VM, render queues, fibers
│   ├── runtime/       # Game loop, StepFrame
│   └── bindings/      # Foreign API
│       ├── raylib/    # Graphics, input, audio
│       ├── dbp/       # DBP-style commands (LoadObject, PositionObject, etc.)
│       ├── box2d/     # 2D physics
│       ├── bullet/    # 3D physics
│       ├── game/      # Tilemaps, particles, helpers
│       ├── net/       # TCP multiplayer
│       └── ...
├── examples/          # BASIC examples
└── main.go            # CLI: compile + run .bas
```

---

## Documentation

| Resource | Description |
|----------|-------------|
| [Documentation Index](docs/DOCUMENTATION_INDEX.md) | Master index of all guides |
| [Getting Started](docs/GETTING_STARTED.md) | Build, run, first program |
| [Rendering and Game Loop](docs/RENDERING_AND_GAME_LOOP.md) | Pipeline, mainloop, SYNC |
| [Game Development Guide](docs/GAME_DEVELOPMENT_GUIDE.md) | Game loop, input, physics |
| [2D Graphics](docs/2D_GRAPHICS_GUIDE.md) | 2D rendering, layers, camera |
| [3D Graphics](docs/3D_GRAPHICS_GUIDE.md) | 3D rendering, models |
| [GUI Guide](docs/GUI_GUIDE.md) | Immediate-mode UI |
| [Multiplayer](docs/MULTIPLAYER.md) | TCP client/server |
| [DBP Parity](docs/DBP_PARITY.md) | DarkBASIC Pro–style features |
| [API Reference](API_REFERENCE.md) | All bindings |

---

## Building

```bash
go build -o cyberbasic .
./cyberbasic examples/first_game.bas
```

On Windows: `.\cyberbasic.exe examples\first_game.bas`

---

## Contributing

The project is open to contributions. Bindings live in `compiler/bindings/<package>`, the VM and compiler in `compiler/`, and documentation in `docs/`. Add commands in the appropriate binding package and document them in [COMMAND_REFERENCE.md](docs/COMMAND_REFERENCE.md) and [API_REFERENCE.md](API_REFERENCE.md). Run `go build ./...` and tests before submitting a PR.

CyberBASIC 2 is an ambitious, evolving engine: stable in its core systems and actively improving. We aim to keep the architecture clean, the docs accurate, and the project welcoming to new contributors.
