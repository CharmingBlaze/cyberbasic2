# CyberBASIC2 Documentation Index

Complete guide to all CyberBASIC2 documentation.

## Philosophy

- **[Documentation Philosophy](DOCUMENTATION_PHILOSOPHY.md)** – DBP-style simplicity, flat API, sensible defaults, deterministic multiplayer, contributor-friendly design

## Entry Points

- **[Concepts (2D vs 3D vs physics)](CONCEPTS.md)** – One-page overview for new users
- **[Quick Start Guide](QUICK_START_GUIDE.md)** – Get running in 5 minutes (DBP-style or manual loop)
- **[Getting Started](GETTING_STARTED.md)** – Installation, building, running your first program
- **[Learning Path](LEARNING_PATH.md)** – Complete structured curriculum from beginner to advanced
- **[DBP Parity](DBP_PARITY.md)** – Zero-boilerplate (OnStart/OnUpdate/OnDraw) verification checklist
- **[Quick Reference](QUICK_REFERENCE.md)** – One-page syntax reference for daily use
- **[FAQ](FAQ.md)** – Hybrid vs manual loop, in-process vs multi-process windows, API and include syntax
- **[Troubleshooting](TROUBLESHOOTING.md)** – Common errors (compiler not found, parse errors, runtime, multiplayer)
- **[Roadmap (ROADMAP.md)](../ROADMAP.md)** – Planned features and priorities
- **[Roadmap Implementation Status](ROADMAP_IMPLEMENTATION.md)** – What the remediation pass actually completed and what still remains partial
- **[Changelog (CHANGELOG.md)](../CHANGELOG.md)** – Version history and release notes

## Language Reference

- **[Language Spec (LANGUAGE_SPEC.md)](../LANGUAGE_SPEC.md)** – Full language reference
  - Variables, constants, types (DIM, VAR, LET, TYPE...END TYPE, ENUM)
  - Null literal (Nil/Null), IsNull(value)
  - Control flow (IF/THEN, FOR/NEXT, WHILE/WEND, REPEAT/UNTIL, SELECT CASE)
  - Functions, subs, modules, dot notation
  - Includes (#include "file.bas"), events, coroutines, compound assignment

- **[Libraries and includes](LIBRARIES.md)** – Multi-file projects, reusing .bas files as libraries

- **[Block Structure Guide](BLOCK_STRUCTURE_GUIDE.md)** – END keywords (ENDIF / END IF, ENDFUNCTION / END FUNCTION, ENDSUB / END SUB, END MODULE, ELSEIF), single-word vs two-word forms, examples and best practices

- **[Program Structure](PROGRAM_STRUCTURE.md)** – Comments (`//`), feature list, block quick reference, example skeleton

- **[Rendering and the game loop](RENDERING_AND_GAME_LOOP.md)** – Pipeline, manual vs hybrid loop, rule for draw(), diagram

- **[Command Reference](COMMAND_REFERENCE.md)** – Broader raylib + physics + utilities. **For DBP-style commands, see [Core Command Reference](CORE_COMMAND_REFERENCE.md).**

## Command Reference (DBP-Style)

- **[Core Command Reference](CORE_COMMAND_REFERENCE.md)** – DBP-style commands (SYNC, LoadObject, MakeCamera, etc.) with when-to-use notes
- **[2D Game API](2D_GAME_API.md)** – Domain-specific 2D commands with examples (sprites, tilemaps, physics)
- **[3D Game API](3D_GAME_API.md)** – Domain-specific 3D commands with examples (objects, camera, lighting)
- **[Command Reference](COMMAND_REFERENCE.md)** – Broader raylib + physics + utilities
- **[API Reference](../API_REFERENCE.md)** – Exhaustive flat API list (for lookup)
- **[DBP Extended](DBP_EXTENDED.md)** – Module-by-module implementation details (for contributors)

## Game Development

- **[Learning Path](LEARNING_PATH.md)** – Complete curriculum from basics to advanced games
- **[2D Games Tutorial](TUTORIAL_2D_GAMES.md)** – Complete 2D game development guide
- **[3D Games Tutorial](TUTORIAL_3D_GAMES.md)** – Complete 3D game development guide
- **[GUI Development Tutorial](TUTORIAL_GUI_DEVELOPMENT.md)** – User interfaces and menus
- **[Multiplayer Tutorial](TUTORIAL_MULTIPLAYER.md)** – Network programming and multiplayer games
- **[Game Development Guide](GAME_DEVELOPMENT_GUIDE.md)** – Making games with CyberBASIC2
  - Game loop, input handling
  - GAME.* helpers (camera, movement, collision)
  - 2D physics (Box2D) and current 3D Bullet-shaped fallback physics
  - ECS (entity-component system)
  - Best practices
- **[2D Graphics Guide](2D_GRAPHICS_GUIDE.md)** – Full 2D rendering reference
  - Window and frame (InitWindow, ClearBackground; no auto-wrap, compiles as written)
  - Primitives, textures, text, colors
  - 2D camera (SetCamera2D)
  - 2D game checklist
- **[2D Physics Guide](2D_PHYSICS_GUIDE.md)** – Box2D: worlds, bodies, shapes, joints, raycast, collision, StepAllPhysics2D, GAME.* helpers
- **[3D Graphics Guide](3D_GRAPHICS_GUIDE.md)** – Full 3D rendering reference
  - 3D camera (SetCamera3D, GAME.CameraOrbit)
  - Primitives, models, meshes
  - 3D game checklist
  - **3D editor and level builder** (GetMouseRay, PickGroundPlane, level objects, SaveLevel/LoadLevel)
- **[3D Physics Guide](3D_PHYSICS_GUIDE.md)** – Bullet-shaped 3D physics API: worlds, bodies, position/rotation, forces, raycast, StepAllPhysics3D, GAME.* helpers
- **[Asset Pipeline](ASSET_PIPELINE.md)** – LoadAsset, PreloadAsset, caching, level vs object loading

- **[Windows, scaling, and splitscreen](WINDOWS_AND_VIEWS.md)** – Window commands, DPI/scaling, views and split-screen
  - Window and config flags (FLAG_*), blend modes (BLEND_*)
  - GetScreenWidth/Height vs GetRenderWidth/Height, GetWindowScaleDPI, GetScaleDPI
  - CreateView, SetViewTarget, DrawView; GetViewX/Y/Width/Height, SetViewPosition/SetViewSize/SetViewRect
  - CreateSplitscreenLeftRight, CreateSplitscreenTopBottom, CreateSplitscreenFour; splitscreen recipe

- **[In-process multi-window](MULTI_WINDOW_INPROCESS.md)** – Multiple logical windows (viewports) in one process
  - WindowCreate, WindowClose, WindowSetTitle/Size/Position, WindowGetWidth/Height/Position
  - WindowBeginDrawing, WindowEndDrawing, WindowClearBackground, WindowDrawAllToScreen
  - Messages (WindowSendMessage, WindowReceiveMessage), Channels (ChannelCreate/Send/Receive), State (StateSet/Get/Has/Remove)
  - Events (OnWindowUpdate/Draw/Resize/Close/Message), 3D (WindowSetCamera, WindowDrawModel/Scene), RPC (WindowRegisterFunction, WindowCall), Docking
  - Example: multi_window_gui_demo.bas

- **[Multiple windows from one .bas (multi-process)](MULTI_WINDOW.md)** – Spawn extra windows (child processes) and talk via NET
  - IsWindowProcess(), GetWindowTitle/Width/Height, SpawnWindow(port, title, width, height), ConnectToParent()
  - Script pattern: branch on IsWindowProcess(); main Host + SpawnWindow + AcceptTimeout; child ConnectToParent + InitWindow + loop

- **[ECS Guide](ECS_GUIDE.md)** – Entity-Component System (via library)
  - Create world, entities, components
  - Queries and iteration
  - Example (ecs_demo.bas)

- **[Multiplayer (TCP)](MULTIPLAYER.md)** – Full multiplayer guide: TCP client/server
  - Connect, Send, Receive, Disconnect (client)
  - Host, Accept, CloseServer (server)
  - Event callbacks (OnClientConnect, OnMessage), SendTable/ReceiveTable, RPC, entity sync
- **[Multiplayer Design](MULTIPLAYER_DESIGN.md)** – Current architecture, fixed-step simulation guidance, and future lockstep/rollback gaps

- **[GUI Guide](GUI_GUIDE.md)** – Full GUI guide: immediate-mode UI
  - BeginUI, EndUI, Label, Button, Slider, Checkbox, TextBox, Dropdown, ProgressBar
  - WindowBox, GroupBox, layout and examples

- **[SQL (SQLite)](SQL.md)** – Full SQL guide: OpenDatabase, Exec, Query, parameterized statements, transactions, common patterns

- **[World, Water, Terrain, Clouds](WORLD_WATER_TERRAIN.md)** – Water, terrain, skybox, clouds, sun, time
- **[Level Loading](LEVEL_LOADING.md)** – Unified 3D loading (LOAD LEVEL loads meshes, materials, textures, hierarchy, and collision hooks)
- **[3D Loading Spec](3D_LOADING_SPEC.md)** – Design goals and safe loading behavior for 3D assets

## Workflows

- **[Blender to CyberBASIC2](BLENDER_WORKFLOW.md)** – Export 3D models (GLTF, FBX, OBJ), PBR materials, animation
- **[Aseprite to CyberBASIC2](ASEPRITE_WORKFLOW.md)** – Export sprite sheets with JSON, tags, slices

## API and Reference

- **[API Reference (API_REFERENCE.md)](../API_REFERENCE.md)** – All bindings (raylib, Box2D, Bullet, GAME, ECS, std)
- **[Cheatsheet (CHEATSHEET.md)](../CHEATSHEET.md)** – First 10 lines for 2D and 3D games
- **[Modules (MODULES.md)](../MODULES.md)** – Codebase layout and adding bindings

## Compatibility

- **[DarkBASIC Pro compatibility (DBP_COMPAT.md)](DBP_COMPAT.md)** – Map DBP commands to CyberBASIC2 (existing or new)
- **[DBP gap list (DBP_GAP.md)](DBP_GAP.md)** – Commands added for DBP parity (Left, Right, Mid, Len, Chr, Asc, Str, Val, Rnd, Int, CopyFile, ListDir, ExecuteFile)

## Examples

- **[Examples README](../examples/README.md)** – Index of example programs
- **Examples:** [hello_world.bas](../examples/hello_world.bas), [first_game.bas](../examples/first_game.bas)
- **Templates:** [2D game](../templates/2d_game.bas), [3D game](../templates/3d_game.bas)

---

## Quick Navigation

| I want to… | Start here |
|------------|------------|
| **Get started immediately** | [Quick Start Guide](QUICK_START_GUIDE.md) |
| **Learn the language step-by-step** | [Learning Path](LEARNING_PATH.md) → [Quick Reference](QUICK_REFERENCE.md) → [Language Spec](../LANGUAGE_SPEC.md) |
| **Make a 2D game** | [2D Games Tutorial](TUTORIAL_2D_GAMES.md) → [Game Development Guide](GAME_DEVELOPMENT_GUIDE.md) → [2D Graphics Guide](2D_GRAPHICS_GUIDE.md) |
| **Make a 3D game** | [3D Games Tutorial](TUTORIAL_3D_GAMES.md) → [Game Development Guide](GAME_DEVELOPMENT_GUIDE.md) → [3D Graphics Guide](3D_GRAPHICS_GUIDE.md) |
| **Build a complete game (2D+3D+GUI+multiplayer)** | [Game Development Guide](GAME_DEVELOPMENT_GUIDE.md#combining-2d-3d-gui-and-multiplayer) |
| **Create user interfaces** | [GUI Development Tutorial](TUTORIAL_GUI_DEVELOPMENT.md) → [GUI Guide](GUI_GUIDE.md) |
| **Add multiplayer** | [Multiplayer Tutorial](TUTORIAL_MULTIPLAYER.md) → [Multiplayer Guide](MULTIPLAYER.md) |
| **Use 2D physics (Box2D)** | [2D Physics Guide](2D_PHYSICS_GUIDE.md) |
| **Use 3D physics (Bullet-shaped fallback)** | [3D Physics Guide](3D_PHYSICS_GUIDE.md) |
| **Use the hybrid loop** | [Program Structure](PROGRAM_STRUCTURE.md#hybrid-updatedraw-loop) (define update(dt) and draw()) |
| **Use in-process multi-window** | [In-process multi-window](MULTI_WINDOW_INPROCESS.md) |
| **Use ECS** | [ECS Guide](ECS_GUIDE.md) → [API Reference](../API_REFERENCE.md) |
| **Look up a function** | [API Reference](../API_REFERENCE.md) |
| **Look up DBP-style commands** | [Core Command Reference](CORE_COMMAND_REFERENCE.md) → [2D Game API](2D_GAME_API.md) / [3D Game API](3D_GAME_API.md) |
| **Use zero-boilerplate (OnStart/OnUpdate/OnDraw)** | [DBP Parity](DBP_PARITY.md) → [examples/first_game.bas](../examples/first_game.bas) |
| **Add water to my 3D scene** | [World, Water, Terrain](WORLD_WATER_TERRAIN.md#1-water-commands-simple-to-advanced) |
| **Add terrain (heightmap, procedural)** | [World, Water, Terrain](WORLD_WATER_TERRAIN.md#2-terrain-commands-simple-to-advanced) |
| **Load a 3D level (GLTF/OBJ)** | [Level Loading](LEVEL_LOADING.md) |
| **Build a full 3D world (level + water + terrain + sky)** | [3D Graphics Guide](3D_GRAPHICS_GUIDE.md) → [World, Water, Terrain](WORLD_WATER_TERRAIN.md) → [Level Loading](LEVEL_LOADING.md) |
| **Optimize 3D performance** | [3D Graphics Guide](3D_GRAPHICS_GUIDE.md#optimization-culling-and-pbr) |

**Current shipped feature set:** CyberBASIC2 supports **full 2D** and **full 3D** graphics, **authoritative 2D physics** (Box2D), a **Bullet-shaped 3D physics API backed by the shipped pure-Go fallback**, **full ECS** (entity-component system), **GUI** (BeginUI, Label, Button, Slider, Checkbox, etc.), and **multiplayer** (TCP Connect/Send/Receive, Host/Accept). See the guides above for the exact current scope of each area.

All documentation lives in the `docs/` directory and the project root. Start with [Getting Started](GETTING_STARTED.md). For doc conventions (headings, code blocks, links), see [Documentation style](STYLE.md).
