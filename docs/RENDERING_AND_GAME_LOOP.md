# Rendering and the Game Loop

This document explains how drawing and the game loop work in CyberBasic: the two modes (manual vs hybrid), the pipeline when using `update`/`draw`, and the rule for what to call inside `draw()`.

## Two modes

### Manual loop

You write the full game loop yourself. You get delta time, step physics (if any), update game state, then draw. If you use 2D or 3D mode, you call **BeginDrawing**/**EndDrawing** and optionally **BeginMode2D**/**EndMode2D** or **BeginMode3D**/**EndMode3D** yourself. The compiler does not inject any code.

Use the manual loop when you need full control over the order of operations, or when you are not using `update`/`draw`.

### Hybrid loop

You define **`update(dt)`** and/or **`draw()`** (as Sub or Function) and use a game loop with an **empty body**: `WHILE NOT WindowShouldClose() WEND` (or `REPEAT ... UNTIL WindowShouldClose()`). The compiler **replaces** the loop body with an automatic pipeline. You do not call BeginDrawing/EndDrawing or BeginMode2D/BeginMode3D/EndMode2D/EndMode3D yourself; the engine does it.

Prefer the hybrid loop for new games when you want automatic physics stepping and a clear split between update and draw.

## Pipeline (one frame)

When using the hybrid loop, each frame runs in this order:

1. **GetFrameTime** → `dt`
2. **StepAllPhysics2D(dt)** and **StepAllPhysics3D(dt)** (all registered worlds)
3. **update(dt)** (if defined)
4. **ClearRenderQueues**
5. **draw()** (if defined) — all Draw*/Gui* calls inside `draw()` are **queued** (2D, 3D, GUI)
6. **FlushRenderQueues** — the engine then:
   - BeginDrawing
   - ClearBackground
   - 2D by layer (sorted, with parallax/scroll per layer)
   - BeginMode3D … all 3D draws … EndMode3D
   - GUI (2D overlay)
   - EndDrawing

```mermaid
flowchart LR
  subgraph frame [One frame]
    A[GetFrameTime] --> B[StepAllPhysics2D/3D]
    B --> C[update(dt)]
    C --> D[ClearRenderQueues]
    D --> E[draw()]
    E --> F[FlushRenderQueues]
  end
  subgraph flush [Flush]
    F --> G[BeginDrawing]
    G --> H[2D by layer]
    H --> I[BeginMode3D ... 3D draws ... EndMode3D]
    I --> J[GUI]
    J --> K[EndDrawing]
  end
```

## Rule for draw()

Inside **`draw()`**, only call:

- **Draw*** (DrawRectangle, DrawCircle, DrawTexture, DrawModel, DrawCube, ClearBackground, etc.)
- **Gui*** (GuiButton, GuiLabel, etc.)
- **SpriteDraw**, **DrawTilemap**, **DrawBackground**, **DrawParticleEmitter**, **SpriteBatchBegin**/**SpriteBatchEnd**, etc.

Do **not** call **BeginMode3D**, **EndMode3D**, **BeginMode2D**, or **EndMode2D** in `draw()`. The engine wraps 2D and 3D blocks during flush. If you do call them in `draw()`, they are **ignored** (forgiving behavior), so existing or mistaken code does not break.

## Example (hybrid loop)

```basic
SUB update(dt)
  REM move player, update state using dt
END SUB
SUB draw()
  ClearBackground(30, 30, 45, 255)
  DrawRectangle(x, y, 40, 40, 255, 100, 100, 255)
  DrawText("Hello", 20, 20, 20, 255, 255, 255, 255)
END SUB
WHILE NOT WindowShouldClose()
WEND
```

See **examples/hybrid_update_draw_demo.bas**.

## See also

- [Program Structure](PROGRAM_STRUCTURE.md#hybrid-updatedraw-loop)
- [FAQ](FAQ.md) — Hybrid vs manual loop
- [Command Reference](COMMAND_REFERENCE.md) — ClearRenderQueues, FlushRenderQueues
- [2D Graphics Guide](2D_GRAPHICS_GUIDE.md), [3D Graphics Guide](3D_GRAPHICS_GUIDE.md)
