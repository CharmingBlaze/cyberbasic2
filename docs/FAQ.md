# Frequently asked questions

## Game loop and hybrid

**What is the difference between the hybrid loop and the manual loop?**

- **Manual loop:** You write the full game loop yourself: get delta time, step physics (if any), update game state, clear background, draw 2D/3D/GUI, and (if needed) call BeginDrawing/EndDrawing or BeginMode2D/EndMode3D explicitly. The compiler does not inject any code.
- **Hybrid loop:** You define **update(dt)** and/or **draw()** (Sub or Function) and use a game loop with an **empty body** (`WHILE NOT WindowShouldClose() WEND`). The compiler injects: GetFrameTime, StepAllPhysics2D(dt), StepAllPhysics3D(dt), update(dt), ClearRenderQueues, draw(), FlushRenderQueues. All Draw*/Gui* calls inside draw() are queued and executed in order. Use the hybrid loop when you want a clear update/draw split and automatic physics stepping. See [Program Structure](PROGRAM_STRUCTURE.md#hybrid-updatedraw-loop).

**When should I use the hybrid loop?**

Prefer it for new games when you want automatic physics step and a single place (draw()) for all rendering. Use the manual loop when you need full control over the order of operations or when maintaining legacy-style code.

## Multi-window

**When should I use in-process multi-window vs multi-process (SpawnWindow)?**

- **In-process (WindowCreate, etc.):** Multiple logical viewports in **one process**. They share memory and state; you use Channel/State and WindowSendMessage for communication. Best for editor UIs, panels, and tools inside one app. See [In-process multi-window](MULTI_WINDOW_INPROCESS.md).
- **Multi-process (SpawnWindow, ConnectToParent):** Separate **processes**, each with its own raylib window. They talk over TCP (Send/Receive). Best when you need true separate windows or isolation. See [Multiple windows from one .bas](MULTI_WINDOW.md).

## API and syntax

**Why do some examples use RL.InitWindow or RL.ClearBackground?**

The docs prefer **flat names** (InitWindow, ClearBackground) so examples work the same whether or not a namespace is used. Both forms are valid; use flat names for consistency. Physics and game helpers use namespaces: **BOX2D.***, **BULLET.***, **GAME.***.

**What include syntax should I use?**

Use **`#include "path/to/file.bas"`**. The compiler also accepts **IMPORT "path"** as an alias. Path is relative to the file that contains the directive. See [Libraries and includes](LIBRARIES.md).

**How do I see all commands?**

Run **`cyberbasic --list-commands`** for a short grouped list. For the full reference, see [Command Reference](COMMAND_REFERENCE.md) and [API Reference](../API_REFERENCE.md). In your program you can call **HELP()** or **?()** to print a reminder and paths.

**Why does my draw() not show anything?**

Make sure you are using the **game loop** (`WHILE NOT WindowShouldClose() WEND` or `REPEAT ... UNTIL WindowShouldClose()`) and that you have defined both **update(dt)** and **draw()** (Sub or Function). The hybrid pipeline only runs when the compiler detects this pattern. See [Program Structure – hybrid loop](PROGRAM_STRUCTURE.md#hybrid-updatedraw-loop).

**When to use BOX2D.* vs flat CreateWorld2D?**

Both work. **CreateWorld2D**, **CreateBox2D**, **Step2D**, **GetPositionX2D**, etc. are flat names. **BOX2D.CreateWorld**, **BOX2D.Step**, **BOX2D.GetPositionX** are the namespaced form. Prefer flat names for consistency with the rest of the docs. The hybrid loop calls **StepAllPhysics2D(dt)** automatically for all registered 2D worlds.

## See also

- [Documentation Index](DOCUMENTATION_INDEX.md)
- [Getting Started](GETTING_STARTED.md)
- [Troubleshooting](TROUBLESHOOTING.md)
