# CyberBasic2 – A BASIC-like language with Raylib + Bullet physics

**CyberBasic2** is a modern BASIC-inspired language with full 2D/3D game development capabilities.  
Repository: **[github.com/CharmingBlaze/cyberbasic2](https://github.com/CharmingBlaze/cyberbasic2)** · Report issues and contribute there.

## Architecture

- **Compiler/Interpreter**: Go-based lexer, parser, and bytecode VM
- **Runtime**: Go bindings for **raylib-go** (graphics, window, input, audio), **Box2D** (2D physics), and **Bullet** (3D physics). No C engine required for normal run.
- **Optional**: An `engine/` C layer exists for custom builds; the default flow uses only Go.

## Project Structure

```
CyberBasic/
├── compiler/          # Go compiler components
│   ├── lexer/        # Tokenizer
│   ├── parser/       # AST builder
│   ├── vm/           # Bytecode VM
│   ├── runtime/      # BASIC runtime (window, game loop)
│   └── bindings/     # Foreign API: raylib, box2d, bullet
│       ├── raylib/   # Window, shapes, text, textures, images, 3D, audio, fonts, math
│       ├── box2d/    # 2D physics
│       ├── bullet/   # 3D physics
│       ├── ecs/      # Entity-component system
│       ├── net/      # TCP multiplayer (Connect, Send, Receive, Host, Accept)
│       └── std/      # File, JSON, HTTP, IsNull, HELP
├── engine/           # Optional C wrapper (raylib_wrapper.c, bullet_wrapper.c)
├── examples/         # BASIC example programs
└── main.go           # CLI: compile + run .bas
```

## BASIC Language Features

- Classic BASIC syntax: IF...THEN, FOR...NEXT, WHILE...WEND, REPEAT...UNTIL, SELECT CASE
- Types: DIM, TYPE...END TYPE, ENUM, CONST; dot notation
- **User functions and subs** with parameters and Return; call by name.
- **Modules:** Module Name … End Module (body is Function/Sub only); call as ModuleName.FunctionName(...).
- **Event handlers:** On KeyDown("KEY") … End On, On KeyPressed("KEY") … End On; run when PollInputEvents() is called.
- **Coroutines:** StartCoroutine SubName(), Yield, WaitSeconds(seconds); fibers share the same chunk (WaitSeconds currently blocks the whole VM).
- Raylib-style API (unprefixed or `RL.`): InitWindow, ClearBackground, DrawCircle, LoadImage, LoadSound, PlayMusic, etc.
- **Hybrid loop:** Define **update(dt)** and **draw()** for automatic physics step and render queue; empty game loop body. See [Program Structure](docs/PROGRAM_STRUCTURE.md#hybrid-updatedraw-loop).
- **In-process multi-window:** WindowCreate, WindowClose, Channel/State, OnWindowUpdate/Draw, WindowDrawAllToScreen. See [In-process multi-window](docs/MULTI_WINDOW_INPROCESS.md).
- Export to C code: **ExportImageAsCode**, **ExportFontAsCode**, **ExportWaveAsCode** (write .h with pixel/sample data)
- Physics: BOX2D.* (2D), BULLET.* (3D)

## Limitations

- **UI:** Full immediate-mode GUI: BeginUI, EndUI, Label, Button, Slider, Checkbox, TextBox, Dropdown, ProgressBar, WindowBox, GroupBox (see [docs/GUI_GUIDE.md](docs/GUI_GUIDE.md)).
- **WaitSeconds** is non-blocking: the current fiber yields until the delay elapses; other fibers keep running (scheduler in VM).
- **Audio callbacks** (SetAudioStreamCallback, AttachAudioStreamProcessor, AttachAudioMixedProcessor) are not supported from BASIC because they require a function pointer; use **UpdateAudioStream**(streamId, ...samples) to push PCM instead.
- **Physics stubs:** Box2D joint APIs (CreateRevoluteJoint2D, etc.) and Bullet joint/body-property APIs (CreateHingeJoint3D, SetFriction3D, etc.) are no-op stubs; see API_REFERENCE.md for the full list.

**Full feature set:** Full 2D and 3D graphics, full 2D physics (Box2D), full 3D physics (Bullet), full ECS, full GUI (widgets above), and multiplayer (TCP Connect/Send/Receive, Host/Accept). See **[docs/DOCUMENTATION_INDEX.md](docs/DOCUMENTATION_INDEX.md)** for the full doc index. **First 2D game:** [CHEATSHEET.md](CHEATSHEET.md) and [docs/2D_GRAPHICS_GUIDE.md](docs/2D_GRAPHICS_GUIDE.md). **First 3D game:** [CHEATSHEET.md](CHEATSHEET.md) and [docs/3D_GRAPHICS_GUIDE.md](docs/3D_GRAPHICS_GUIDE.md). **Examples:** [examples/README.md](examples/README.md).

## Documentation

- **[Documentation Index](docs/DOCUMENTATION_INDEX.md)** – Master index of all guides
- **[Getting Started](docs/GETTING_STARTED.md)** – Build, run first program
- **[Quick Reference](docs/QUICK_REFERENCE.md)** – One-page syntax
- **[Language Spec](LANGUAGE_SPEC.md)** – Full language reference
- **First 2D game:** [Cheatsheet](CHEATSHEET.md) + [2D Graphics Guide](docs/2D_GRAPHICS_GUIDE.md)
- **First 3D game:** [Cheatsheet](CHEATSHEET.md) + [3D Graphics Guide](docs/3D_GRAPHICS_GUIDE.md)
- **[Game Development Guide](docs/GAME_DEVELOPMENT_GUIDE.md)** – Game loop, input, physics, ECS
- **[GUI Guide](docs/GUI_GUIDE.md)** – Immediate-mode UI (Button, Slider, Checkbox, TextBox, etc.)
- **[Multiplayer](docs/MULTIPLAYER.md)** – TCP client/server (Connect, Send, Receive, Host, Accept)
- **[In-process multi-window](docs/MULTI_WINDOW_INPROCESS.md)** – Multiple viewports in one process
- **[API Reference](API_REFERENCE.md)** – All bindings

## Building

```bash
# Build and run (uses raylib-go; no C build needed)
go build -o cyberbasic .

# Run a BASIC file
./cyberbasic examples/first_game.bas

# Optional: build the C engine (requires Raylib and Bullet)
cd engine && make
```
