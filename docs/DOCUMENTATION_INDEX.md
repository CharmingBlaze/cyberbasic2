# CyberBasic Documentation Index

Complete guide to all CyberBasic documentation.

## Getting Started

- **[Getting Started Guide](GETTING_STARTED.md)** – Installation, building, running your first program
- **[Quick Reference](QUICK_REFERENCE.md)** – One-page syntax reference for daily use
- **[BASIC Programming Guide](BASIC_PROGRAMMING_GUIDE.md)** – Step-by-step tutorial: variables, types, I/O, errors
- **[FAQ](FAQ.md)** – Hybrid vs manual loop, in-process vs multi-process windows, API and include syntax
- **[Troubleshooting](TROUBLESHOOTING.md)** – Common errors (compiler not found, parse errors, runtime, multiplayer)

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

- **[Command Reference](COMMAND_REFERENCE.md)** – Structured command set: window, input, math, camera, 3D, 2D, audio, file, game loop, utility

## Game Development

- **[Game Development Guide](GAME_DEVELOPMENT_GUIDE.md)** – Making games with CyberBasic
  - Game loop, input handling
  - GAME.* helpers (camera, movement, collision)
  - 2D/3D physics (Box2D, Bullet)
  - ECS (entity-component system)
  - Best practices

- **[2D Graphics Guide](2D_GRAPHICS_GUIDE.md)** – 2D rendering reference
  - Window and frame (InitWindow, ClearBackground; no auto-wrap, compiles as written)
  - Primitives, textures, text, colors
  - 2D camera (SetCamera2D)
  - 2D game checklist

- **[3D Graphics Guide](3D_GRAPHICS_GUIDE.md)** – 3D rendering reference
  - 3D camera (SetCamera3D, GAME.CameraOrbit)
  - Primitives, models, meshes
  - 3D game checklist
  - **3D editor and level builder** (GetMouseRay, PickGroundPlane, level objects, SaveLevel/LoadLevel)

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

- **[Multiplayer (TCP)](MULTIPLAYER.md)** – Simple TCP client/server
  - Connect, Send, Receive, Disconnect (client)
  - Host, Accept, CloseServer (server)
  - Line-based messages, game-loop usage

- **[GUI Guide](GUI_GUIDE.md)** – Immediate-mode UI
  - BeginUI, EndUI, Label, Button, Slider, Checkbox, TextBox, Dropdown, ProgressBar
  - WindowBox, GroupBox, layout and examples

## API and Reference

- **[API Reference (API_REFERENCE.md)](../API_REFERENCE.md)** – All bindings (raylib, Box2D, Bullet, GAME, ECS, std)
- **[Cheatsheet (CHEATSHEET.md)](../CHEATSHEET.md)** – First 10 lines for 2D and 3D games
- **[Modules (MODULES.md)](../MODULES.md)** – Codebase layout and adding bindings

## Compatibility

- **[DarkBASIC Pro compatibility (DBP_COMPAT.md)](DBP_COMPAT.md)** – Map DBP commands to CyberBasic (existing or new)
- **[DBP gap list (DBP_GAP.md)](DBP_GAP.md)** – Commands added for DBP parity (Left, Right, Mid, Len, Chr, Asc, Str, Val, Rnd, Int, CopyFile, ListDir, ExecuteFile)

## Examples

- **[Examples README](../examples/README.md)** – Index of example programs
- **Templates:** [2D game](../templates/2d_game.bas), [3D game](../templates/3d_game.bas)

---

## Quick Navigation

| I want to… | Start here |
|------------|------------|
| **Learn the language** | [Getting Started](GETTING_STARTED.md) → [Quick Reference](QUICK_REFERENCE.md) → [Language Spec](../LANGUAGE_SPEC.md) |
| **Make a 2D game** | [Game Development Guide](GAME_DEVELOPMENT_GUIDE.md) → [2D Graphics Guide](2D_GRAPHICS_GUIDE.md) → [Cheatsheet](../CHEATSHEET.md) |
| **Make a 3D game** | [Game Development Guide](GAME_DEVELOPMENT_GUIDE.md) → [3D Graphics Guide](3D_GRAPHICS_GUIDE.md) → [Cheatsheet](../CHEATSHEET.md) |
| **Use the hybrid loop** | [Program Structure](PROGRAM_STRUCTURE.md#hybrid-updatedraw-loop) (define update(dt) and draw()) |
| **Use in-process multi-window** | [In-process multi-window](MULTI_WINDOW_INPROCESS.md) |
| **Use ECS** | [ECS Guide](ECS_GUIDE.md) → [API Reference](../API_REFERENCE.md) |
| **Look up a function** | [API Reference](../API_REFERENCE.md) |

**Full feature set:** CyberBasic supports **full 2D** and **full 3D** graphics, **full 2D physics** (Box2D), **full 3D physics** (Bullet), **full ECS** (entity-component system), **GUI** (BeginUI, Label, Button, Slider, Checkbox, etc.), and **multiplayer** (TCP Connect/Send/Receive, Host/Accept). See the guides above for each area.

All documentation lives in the `docs/` directory and the project root. Start with [Getting Started](GETTING_STARTED.md). For doc conventions (headings, code blocks, links), see [Documentation style](STYLE.md).
