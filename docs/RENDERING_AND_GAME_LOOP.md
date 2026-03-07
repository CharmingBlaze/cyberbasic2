# Rendering and the Game Loop

How drawing and the game loop work in CyberBASIC2: three modes, frame pipeline, SYNC behavior, and rules for draw code.

---

## Purpose

- **Unify frame logic:** One frame = one update, one draw, one present.
- **Support multiple styles:** DBP-style (OnStart/OnUpdate/OnDraw), manual loop, hybrid loop.
- **SYNC as frame boundary:** End frame, poll input, present. Use SYNC at the end of your loop or OnDraw. SYNC polls input and presents the frame.

---

## Three Modes

| Mode | When to Use | Frame handling |
|------|-------------|----------------|
| **DBP-style** | Zero boilerplate, OnStart/OnUpdate/OnDraw | Runtime provides window + loop; SYNC at end of OnDraw |
| **Manual loop** | Full control, explicit Begin/End | You write WHILE NOT WindowShouldClose(); compiler injects BeginDrawing/EndDrawing when loop body contains draw calls |
| **Hybrid loop** | update(dt) + draw() with empty body | Compiler replaces body with StepFrame; automatic physics + queue + flush |

---

## DBP-Style (OnStart/OnUpdate/OnDraw)

**Zero boilerplate:** Define `OnStart()`, `OnUpdate(dt)`, and `OnDraw()`. No `InitWindow`, no `WHILE` loop.

- Call `UseUnifiedRenderer` in `OnStart()` for the full 3D→2D→GUI pipeline.
- Call `SYNC` at the end of `OnDraw()` to end the frame.

See [DBP Parity](DBP_PARITY.md).

---

## Manual Loop

You write the full game loop. The compiler injects `BeginDrawing`/`EndDrawing` (and optionally `BeginMode2D`/`EndMode2D`) when the loop body contains draw calls and does not contain `BeginDrawing`/`EndDrawing` (or `Draw`/`update` subs).

- **With SYNC:** If the loop body contains `SYNC`, the compiler omits injected `EndDrawing`; SYNC does it.
- **Without SYNC:** Compiler injects both BeginDrawing and EndDrawing.

Use for full control or when porting existing code.

---

## Hybrid Loop

Define `update(dt)` and/or `draw()` and use an empty body:

```basic
WHILE NOT WindowShouldClose()
WEND
```

The compiler replaces the body with `StepFrame`. You do not call BeginDrawing/EndDrawing or mode Begin/End; the engine does it. Draw calls inside `draw()` are queued and flushed in order (2D then 3D then GUI).

---

## Pipeline (One Frame)

### Hybrid loop (StepFrame)

1. **GetFrameTime** → dt
2. **Accumulate** scaled frame time into fixed-step clock
3. **Fixed steps** (zero or more): StepAllPhysics2D, StepAllPhysics3D, OnFixedUpdate
4. **update(dt)** (if defined)
5. **ClearRenderQueues**
6. **draw()** (if defined) — Draw*/Gui* calls are queued
7. **FlushRenderQueues** — BeginDrawing, ClearBackground, 2D by layer, BeginMode3D…EndMode3D, GUI, EndDrawing

### DBP-style (StepImplicitFrame)

- Same fixed-step path as hybrid
- Uses OnUpdate/OnDraw naming
- When UseUnifiedRenderer: OnDraw queues; SYNC runs full frame (3D→2D→GUI)
- When not UseUnifiedRenderer: StepImplicitFrame does BeginDrawing, OnDraw, EndDrawing

### SYNC Behavior

| Context | SYNC does |
|---------|-----------|
| UseUnifiedRenderer on | Full frame (renderer.Frame: 3D→2D→GUI, swap buffers) |
| UseUnifiedRenderer off | EndDrawing only (BeginDrawing already polled input at frame start) |

**Input polling:** PollInputEvents is called exactly once per frame at the start (BeginDrawing or beginRuntimeFrame). Do not poll again; a second poll clears IsKeyPressed/IsMouseButtonPressed before user code can read them.

---

## Rule for draw()

Inside **`draw()`** (or **`OnDraw()`**), only call:

- **Draw*** (DrawRectangle, DrawCircle, DrawTexture, DrawModel, DrawCube, ClearBackground, etc.)
- **Gui*** (GuiButton, GuiLabel, etc.)
- **SpriteDraw**, **DrawTilemap**, **DrawBackground**, **DrawParticleEmitter**, **SpriteBatchBegin**/**SpriteBatchEnd**, etc.

Do **not** call **BeginMode3D**, **EndMode3D**, **BeginMode2D**, or **EndMode2D**. The engine wraps 2D and 3D blocks during flush. If you do call them, they are ignored.

---

## Defaults

| Setting | Default |
|---------|---------|
| Fixed step rate | 1/60 (60 Hz) |
| Frame delta | Clamped for catch-up |
| Clear color | Set by SetClearColor or renderer default |
| Draw order | 3D → 2D → GUI |

---

## Edge Cases

- **Loop with no draw calls:** No frame wrap injected; window may not render. Add at least one draw call or use hybrid/DBP-style.
- **Loop with BeginDrawing/EndDrawing in body:** Compiler does not inject; you control the frame.
- **Loop with Draw() in body:** Treated as hybrid; StepFrame used.
- **SYNC without UseUnifiedRenderer:** SYNC does EndDrawing only. Ensure BeginDrawing was called earlier in the frame.
- **Multiple SYNC per frame:** Each SYNC ends the frame. Avoid; one SYNC per frame.

---

## Performance Considerations

- **Fixed-step cap:** Catch-up is capped to avoid spiral-of-death frame. See `maxFixedCatchupSteps` in `compiler/runtime/loop.go`.
- **Queue flush:** FlushRenderQueues runs once per frame. Heavy draw calls in draw() are queued; flush does the actual GPU work.

---

## Multiplayer / Determinism

Run simulation from `OnFixedUpdate`; read input from the main frame. Fixed-step ensures deterministic physics. See [Multiplayer Design](MULTIPLAYER_DESIGN.md).

---

## Contributor Notes

- **Frame wrap:** `compiler/codegen_statements.go` — `isGameLoopCondition`, `bodyContainsFrameBoundaries`, `bodyContainsSync`, `emitFrameWrap`
- **StepFrame:** `compiler/runtime/loop.go` — `StepFrame`, `StepImplicitFrame`, `beginRuntimeFrame`
- **SyncFrame:** `compiler/runtime/sync.go` — `SyncFrame`; `FlushRenderQueues` override in `renderer/`
- **Renderer:** `compiler/runtime/renderer/renderer.go` — `Frame()`, `FrameIfUnified()`

---

## Example (Hybrid Loop)

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

See [examples/first_game.bas](../examples/first_game.bas) and [templates/2d_game.bas](../templates/2d_game.bas).

---

## See Also

- [DBP Parity](DBP_PARITY.md)
- [Program Structure](PROGRAM_STRUCTURE.md#hybrid-updatedraw-loop)
- [FAQ](FAQ.md)
- [Command Reference](COMMAND_REFERENCE.md)
- [2D Graphics Guide](2D_GRAPHICS_GUIDE.md), [3D Graphics Guide](3D_GRAPHICS_GUIDE.md)
- [Documentation Index](DOCUMENTATION_INDEX.md)
