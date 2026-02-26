# CyberBasic Examples

Run any example with:
```bash
cyberbasic examples/first_game.bas
```

**Compile sweep:** All examples in this folder (and under `include_demo/`) and both templates (`2d_game.bas`, `3d_game.bas`) compile successfully with `cyberbasic --compile-only <file>.bas`.

## Verification (how to confirm everything works)

From the project root:

1. **Run Go unit tests:** `go test ./compiler/...` — parser, compiler, and VM tests should pass.
2. **Build the CLI:** `go build -o cb.exe .` then `.\cb.exe --help` to confirm the binary runs.
3. **Compile sweep:** Run `.\cb.exe --compile-only <path>` on each working example and template. Working set: `first_game.bas`, `2d_shapes_demo.bas`, `dot_and_colors_demo.bas`, `window_demo.bas`, `minimal_window_test.bas`, `box2d_demo.bas`, `simple_box2d_demo.bas`, `simplest_box2d_demo.bas`, `ecs_demo.bas`, `sql_demo.bas`, `gui_demo.bas`, `gui_raygui_demo.bas`, `hybrid_update_draw_demo.bas`, `multi_window_gui_demo.bas`, `multi_window_demo.bas`, `multiplayer_server.bas`, `multiplayer_client.bas`, `examples/include_demo/main.bas`, `templates/2d_game.bas`, `templates/3d_game.bas` (and optionally `run_3d_physics_demo.bas`, `minimal_3d_demo.bas`). All should compile without errors.
4. **Runtime check (optional):** Run `.\cb.exe examples\simple_box2d_demo.bas` — it should print and exit. Run windowed examples (e.g. `first_game.bas`, `gui_demo.bas`) manually; they should open a window and draw. **sql_demo.bas** may fail to open the database if run with redirected output or from a restricted working directory; run it normally from the project root. **multiplayer_server.bas** and **multiplayer_client.bas** require two terminals (server first, then client).

## API: current vs legacy

- **Current API** (run correctly): Use `InitWindow`, `ClearBackground`, `DrawText`, `DrawCircle`, `WindowShouldClose`, `CloseWindow`, `SetTargetFPS`, `IsKeyDown`, etc. (Raylib-style flat names). The examples listed in **Working examples by feature** below all use the current API and run correctly.
- **Legacy API** (compile but fail at runtime): Names like `INITGRAPHICS`, `CLEARSCREEN`, `DRAWRECTANGLE`, `RL.InitWindow`, etc. are not registered; use the flat current names instead. **platformer.bas**, **agk2_demo.bas**, **working_demo.bas**, and some others still use legacy names—they compile but fail at runtime.

## Working examples by feature

| Feature | Examples | How to run |
|--------|----------|------------|
| **2D graphics** | first_game.bas, 2d_shapes_demo.bas, dot_and_colors_demo.bas, window_demo.bas | `cyberbasic examples/first_game.bas` (or the other .bas) |
| **3D** | run_3d_physics_demo.bas, mario64.bas, minimal_3d_demo.bas | `cyberbasic examples/run_3d_physics_demo.bas` |
| **Box2D (2D physics)** | box2d_demo.bas, simple_box2d_demo.bas, simplest_box2d_demo.bas | `cyberbasic examples/box2d_demo.bas` |
| **ECS** | ecs_demo.bas | `cyberbasic examples/ecs_demo.bas` |
| **SQL (SQLite)** | sql_demo.bas | `cyberbasic examples/sql_demo.bas` (single process; creates sql_demo.db) |
| **Multiplayer** | multiplayer_server.bas, multiplayer_client.bas | Run server first: `cyberbasic examples/multiplayer_server.bas`; then in another terminal: `cyberbasic examples/multiplayer_client.bas` |
| **Hybrid loop** (update/draw) | hybrid_update_draw_demo.bas | `cyberbasic examples/hybrid_update_draw_demo.bas` — update(dt), draw(), empty loop |
| **Multi-window (in-process)** | multi_window_gui_demo.bas | `cyberbasic examples/multi_window_gui_demo.bas` — WindowCreate, StateSet/Get, OnWindowDraw |
| **Multi-window (same .bas, processes)** | multi_window_demo.bas | `cyberbasic examples/multi_window_demo.bas` — main + child window, talking via NET |
| **GUI** | gui_demo.bas | `cyberbasic examples/gui_demo.bas` (BeginUI / EndUI widgets) |
| **GUI (raygui)** | gui_raygui_demo.bas | `cyberbasic examples/gui_raygui_demo.bas` (Gui* functions; requires CGO) |
| **Includes (#include)** | include_demo/main.bas | From project root: `cyberbasic examples/include_demo/main.bas` |

## Quick start

1. Run **first_game.bas** – minimal 2D game (WASD move a circle).
2. Try **templates/2d_game.bas** or **templates/3d_game.bas** for a copy-paste starting point (see [templates/README.md](../templates/README.md)).
3. For 3D: **mario64.bas** or **run_3d_physics_demo.bas**. For 2D physics: **box2d_demo.bas**.

## Start here

- **first_game.bas** – Minimal game loop: window, input (WASD), DrawCircle, WindowShouldClose
- **minimal_window_test.bas** – Bare window open/close
- **hello_world.bas** – Simple Print and window

## 2D / Shapes

- **dot_and_colors_demo.bas** – Mouse position, colors, rectangles (current API)
- **2d_shapes_demo.bas** – Rectangles, circles, lines, animated circle, FPS (current API)
- **window_demo.bas** – Window and drawing (current API)

## 3D

- **run_3d_physics_demo.bas** – 3D physics (Bullet) + raylib 3D
- **mario64.bas** – Mario64-style camera and movement
- **minimal_3d_demo.bas** – Basic 3D scene

## Physics 2D (Box2D)

- **box2d_demo.bas** – Box2D world, bodies, click to spawn (current API)
- **simple_box2d_demo.bas** – Minimal Box2D (console PRINT; no window)
- **simplest_box2d_demo.bas** – One box, ground, window (current API)

## Physics 3D (Bullet)

- **bullet_demo.bas** – Bullet 3D physics
- **run_3d_physics_demo.bas** – 3D physics demo

## Hybrid loop, multi-window, GUI, includes

- **hybrid_update_draw_demo.bas** – update(dt) and draw(); empty game loop; automatic physics step and render queues. See [docs/PROGRAM_STRUCTURE.md](../docs/PROGRAM_STRUCTURE.md#hybrid-updatedraw-loop).
- **multi_window_gui_demo.bas** – In-process multi-window: WindowCreate, StateSet/StateGet, OnWindowDraw, WindowProcessEvents, WindowDrawAllToScreen. See [docs/MULTI_WINDOW_INPROCESS.md](../docs/MULTI_WINDOW_INPROCESS.md).
- **multi_window_demo.bas** – Multi-process: main spawns child with SpawnWindow; they talk via Send/Receive (ConnectToParent, AcceptTimeout). See [docs/MULTI_WINDOW.md](../docs/MULTI_WINDOW.md).
- **gui_demo.bas** – BeginUI, Label, Button, Slider, Checkbox.
- **gui_raygui_demo.bas** – Full raygui (GuiWindowBox, GuiButton, GuiSlider, etc.); requires CGO.
- **include_demo/main.bas** – Uses `#include "lib/helper.bas"` and calls a Sub from the included file.

## SQL, Multiplayer

- **sql_demo.bas** – OpenDatabase, Exec, Query, GetCell; shows results in a window.
- **multiplayer_server.bas** / **multiplayer_client.bas** – Host, Connect, Send, Receive, rooms.

## Input

- **first_game.bas** – IsKeyDown(KEY_W) etc.; for axis-style movement see **templates/2d_game.bas** (GetAxisX, GetAxisY).

## Other

- **test_enum.bas**, **test_const.bas**, **case_test.bas** – Language features
- **math_matrix_test.bas** – Matrix math
- **features_test.bas** – Assorted features

**Note:** **platformer.bas** and several other demos still use legacy names (INITGRAPHICS, CREATEPHYSICSBODY2D, etc.) that are not in the current VM; they compile but will fail at runtime. Use **box2d_demo.bas**, **first_game.bas**, and the current raylib + BOX2D.* API for 2D physics and graphics.
