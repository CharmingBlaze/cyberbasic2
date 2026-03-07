# Rendering and the Game Loop

This document explains how drawing and the game loop work in CyberBASIC2: the three modes (DBP-style, manual, hybrid), the pipeline when using `update`/`draw` or `OnStart`/`OnUpdate`/`OnDraw`, and the rule for what to call inside `draw()`.

## Three modes

### DBP-style (OnStart/OnUpdate/OnDraw)

**Zero boilerplate:** Define `OnStart()`, `OnUpdate(dt)`, and `OnDraw()`. No `InitWindow`, no `WHILE` loop. The runtime now uses the same fixed-step accumulator path as the hybrid loop before calling `OnUpdate(dt)` and `OnDraw()`. Call `UseUnifiedRenderer` in `OnStart()` and `SYNC` at the end of `OnDraw()` when you want the unified 3D→2D→GUI pipeline. See [DBP Parity](DBP_PARITY.md).

### Manual loop

You write the full game loop yourself. You get delta time, step physics (if any), update game state, then draw. If you use 2D or 3D mode, you call **BeginDrawing**/**EndDrawing** and optionally **BeginMode2D**/**EndMode2D** or **BeginMode3D**/**EndMode3D** yourself. The compiler does not inject any code.

Use the manual loop when you need full control over the order of operations, or when you are not using `update`/`draw`.

### Hybrid loop

You define **`update(dt)`** and/or **`draw()`** (as Sub or Function) and use a game loop with an **empty body**: `WHILE NOT WindowShouldClose() WEND` (or `REPEAT ... UNTIL WindowShouldClose()`). The compiler **replaces** the loop body with an automatic pipeline. You do not call BeginDrawing/EndDrawing or BeginMode2D/BeginMode3D/EndMode2D/EndMode3D yourself; the engine does it.

Prefer the hybrid loop for new games when you want automatic physics stepping and a clear split between update and draw.

## Pipeline (one frame)

When using the hybrid loop, each frame runs in this order:

1. **GetFrameTime** → `dt`
2. **Accumulate scaled frame time** into the runtime fixed-step clock
3. Run **zero or more fixed steps** at `FixedDeltaTime()`:
   - `StepAllPhysics2D(FixedDeltaTime())`
   - `StepAllPhysics3D(FixedDeltaTime())`
   - `OnFixedUpdate(label$)` callback if registered
4. **update(dt)** (if defined)
5. **ClearRenderQueues**
6. **draw()** (if defined) — all Draw*/Gui* calls inside `draw()` are **queued** (2D, 3D, GUI)
7. **FlushRenderQueues** — the engine then:
   - BeginDrawing
   - ClearBackground
   - 2D by layer (sorted, with parallax/scroll per layer)
   - BeginMode3D … all 3D draws … EndMode3D
   - GUI (2D overlay)
   - EndDrawing

`dt` is still the per-frame delta for frame-rate-dependent update logic. Physics and deterministic-style gameplay should use `FixedUpdate(rate)` plus `OnFixedUpdate(label$)` when you want explicit fixed-step code alongside the automatic physics stepping.



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

## SYNC and UseUnifiedRenderer

When using **UseUnifiedRenderer**, call **SYNC** at the end of each frame. SYNC replaces manual `Start3D`/`End3D`—the engine handles 3D→2D→GUI order. Use `SYNC` in `OnDraw()` for DBP-style programs.

## See also

- [DBP Parity](DBP_PARITY.md) — Zero-boilerplate (OnStart/OnUpdate/OnDraw)
- [Program Structure](PROGRAM_STRUCTURE.md#hybrid-updatedraw-loop)
- [FAQ](FAQ.md) — Hybrid vs manual loop
- [Command Reference](COMMAND_REFERENCE.md) — ClearRenderQueues, FlushRenderQueues
- [2D Graphics Guide](2D_GRAPHICS_GUIDE.md), [3D Graphics Guide](3D_GRAPHICS_GUIDE.md)
