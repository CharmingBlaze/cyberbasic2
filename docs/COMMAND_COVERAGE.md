# Command coverage — 2D/3D games and applications

This document is the **human-facing coverage matrix** for CyberBasic 2: what you need for games and desktop-style apps, where it lives, and how we track completeness.

## Machine-readable inventory

These files live in **`docs/generated/`** and are **committed** to the repo (CI fails if `foreignaudit` output drifts). The root `.gitignore` only ignores a top-level **`/generated/`** directory, not this path.

- **[`generated/foreign_commands.json`](generated/foreign_commands.json)** — every unique `RegisterForeign` name under `compiler/bindings/`, with source file paths (and `packages` when the same name is registered from more than one Go package).
- **[`generated/FOREIGN_COMMANDS_INDEX.md`](generated/FOREIGN_COMMANDS_INDEX.md)** — the same data as a browsable Markdown index.
- **[`generated/raylib_parity.json`](generated/raylib_parity.json)** — diff helper: top-level exported **functions** in `github.com/gen2brain/raylib-go/raylib` vs names registered from `compiler/bindings/raylib/`. It excludes methods, `rlgl`-level entry points, and symbols exposed only under different names or via DBP/game wrappers.

Regenerate locally (or in CI) with:

```bash
make foreign-audit
```

## Pillar matrix (games and apps)

Status legend: **Yes** = production path in tree; **Stub** = callable but minimal / placeholder; **Partial** = subset or expert-only.

**Dot parity** (subset): **Yes** = common game path has `namespace.method` / handle dot API over the same foreigns; **Partial** = some dot surface, large flat-only remainder; **—** = flat / modfacade only for now.

| Pillar | What you use (examples) | v2 global (optional) | Dot parity | Primary bindings / docs | Smoke example |
|--------|-------------------------|----------------------|------------|-------------------------|---------------|
| **Runtime — explicit loop** | `InitWindow`, `mainloop` / `SYNC`, `CloseWindow` | `window` (title, size, fps, …) | Yes | [`compiler/bindings/raylib`](../compiler/bindings/raylib), [COMMAND_REFERENCE — Game loop](COMMAND_REFERENCE.md#game-loop) | [`examples/first_game.bas`](../examples/first_game.bas), [`examples/smoke_dot_api.bas`](../examples/smoke_dot_api.bas) |
| **Runtime — implicit DBP** | `ON UPDATE` / `ON DRAW`, `WINDOW.TITLE` | `window` | Partial | [`compiler/runtime`](../compiler/runtime), [PROGRAM_STRUCTURE](PROGRAM_STRUCTURE.md) | [`examples/implicit_loop.bas`](../examples/implicit_loop.bas) |
| **Runtime — hybrid** | `update(dt)` / `draw()`, render queues | — | Partial | [RENDERING_AND_GAME_LOOP](RENDERING_AND_GAME_LOOP.md), [COMMAND_REFERENCE — hybrid](COMMAND_REFERENCE.md#game-loop-hybrid) | [`examples/platformer.bas`](../examples/platformer.bas) (pattern) |
| **2D drawing** | `DrawRectangle`, `DrawCircle`, `DrawText`, textures, layers | `draw`, `texture`, `sprite` | Partial | `raylib_shapes`, `raylib_text`, `raylib_textures`, [COMMAND_REFERENCE — 2D](COMMAND_REFERENCE.md) | [`examples/platformer.bas`](../examples/platformer.bas), [`examples/smoke_rectangle_gradient.bas`](../examples/smoke_rectangle_gradient.bas), [`examples/smoke_dot_api.bas`](../examples/smoke_dot_api.bas) |
| **2D collision (geometric)** | `CheckCollisionRecs`, `CheckCollisionCircleLine`, … | — | — | [`raylib_misc.go`](../compiler/bindings/raylib/raylib_misc.go), [COMMAND_REFERENCE — 2D geometric collision](COMMAND_REFERENCE.md#2d-geometric-collision-raylib-shapes) | [`examples/smoke_2d_collision.bas`](../examples/smoke_2d_collision.bas) |
| **3D — raylib primitives** | `BeginMode3D`, `DrawCube`, `DrawGrid`, `SetCamera3D`, `CameraMoveForward`, … | `camera`, `model`, `shapes3d` | Partial → broader | [`raylib_3d.go`](../compiler/bindings/raylib/raylib_3d.go), [`raylibdot`](../compiler/bindings/raylibdot) | [`examples/smoke_raylib_3d.bas`](../examples/smoke_raylib_3d.bas), [`examples/smoke_rcamera.bas`](../examples/smoke_rcamera.bas) |
| **3D — DBP / scene objects** | `LoadObject`, `DrawObject`, `PositionObject`, … | `object` (full DBP handle surface) | Partial | [`compiler/bindings/dbp`](../compiler/bindings/dbp), [`objectdot`](../compiler/bindings/objectdot/objectdot.go) | [`examples/first_game.bas`](../examples/first_game.bas), [`examples/smoke_object_dot_extra.bas`](../examples/smoke_object_dot_extra.bas) |
| **Input — raw** | `IsKeyDown`, `IsMouseButtonPressed`, gamepad | — | — | `raylib_input` | [`examples/input_debug.bas`](../examples/input_debug.bas) |
| **Input — actions** | `InputMapRegister`, `InputPressed` | `input` | Partial | [`inputmap`](../compiler/bindings/inputmap/inputmap.go) | [`examples/smoke_input_v2.bas`](../examples/smoke_input_v2.bas) |
| **Audio** | `InitAudioDevice`, `LoadSound`, `PlaySound`, music streams | `audio.*` | Yes (sound handle) | `raylib_audio`, [`audiosys`](../compiler/bindings/audiosys) | [`examples/smoke_audio.bas`](../examples/smoke_audio.bas) |
| **UI** | `raygui` controls, custom UI foreigns | — | — | `raylib_raygui`, `raylib_ui` | [`examples/ui_demo.bas`](../examples/ui_demo.bas) |
| **Physics 2D** | `CreateWorld2D`, `Step2D`, bodies/joints | `physics.*`, `box2d.*` | Partial | [`box2d`](../compiler/bindings/box2d), [`box2ddot`](../compiler/bindings/box2ddot/box2ddot.go), [`physics2d`](../compiler/bindings/physics2d) | [`examples/smoke_physics2d.bas`](../examples/smoke_physics2d.bas) |
| **Physics 3D** | `CreateWorld3D`, `Step3D`, … | `bullet.*` | Partial | [`bullet`](../compiler/bindings/bullet), [`bulletdot`](../compiler/bindings/bulletdot/bulletdot.go) | [`examples/smoke_physics3d.bas`](../examples/smoke_physics3d.bas) |
| **Assets (key/value)** | `AssetsSet` / `AssetsGet` | `assets.set` / `get` | Yes | [`assets`](../compiler/bindings/assets/assets.go) | [`examples/smoke_assets_v2.bas`](../examples/smoke_assets_v2.bas) |
| **Scenes** | `CreateScene`, `LoadScene`, … | `scenes.*` | Yes | [`scene`](../compiler/bindings/scene) | [`examples/smoke_scenes.bas`](../examples/smoke_scenes.bas) |
| **World building** | terrain, water, vegetation, world, nav, indoor, procedural | `terrain`, `water`, … | Partial | respective packages under `compiler/bindings/` | _(large-game)_ |
| **Std / apps** | files, JSON, HTTP, HELP | `std.*` (expanded), `file.*`, `http` | Partial | [`std`](../compiler/bindings/std/std_v2map.go), [`filedot`](../compiler/bindings/filedot/filedot.go), [`httpdot`](../compiler/bindings/httpdot/httpdot.go) | [`examples/smoke_std.bas`](../examples/smoke_std.bas) |
| **Net / SQL / Nakama** | multiplayer and persistence | `net`, `sql`, `nakama` | Yes | respective packages | _(app-specific)_ |
| **Shaders / FX** | `shader.pbr` / `toon` / `dissolve` (embedded GLSL), `BeginShaderMode` + uniforms | `shader`, `effect`, `camera.fx` | Yes (shader handle) | [`shadersys`](../compiler/bindings/shadersys), [`effect`](../compiler/bindings/effect) | [`examples/shader_demo.bas`](../examples/shader_demo.bas); **effect / camera.fx** still stub — see below |
| **AI / behaviour** | `ai.*` → `navigation.*`; optional `ai.agent` handle | `ai`, `navigation` | Yes | [`aisys`](../compiler/bindings/aisys), [`navigation`](../compiler/bindings/navigation) | [`examples/ai_patrol.bas`](../examples/ai_patrol.bas); **BTREE** deferred — [BTREE_PHASE12.md](BTREE_PHASE12.md) |
| **Tween** | `TweenRegister`, frame tick | `tween.register`, `tween.count` | Yes | [`tween`](../compiler/bindings/tween), [`runtime/loop`](../compiler/runtime/loop.go) | [`examples/smoke_tween.bas`](../examples/smoke_tween.bas) |
| **Composition** | `engine.ecs`, `engine.net`, … | `engine` | Yes | [`engine`](../compiler/bindings/engine/engine.go) | [`examples/smoke_engine.bas`](../examples/smoke_engine.bas) |

## Stub and partial APIs (honest expectations)

- **`ai`**: **`ai.version()`** plus **navigation aliases** (`navgridcreate`, `navagentcreate`, …) and **`ai.agent(id$)`**; **BTREE** / full behaviour-tree syntax is **not** in-tree — see [BTREE_PHASE12.md](BTREE_PHASE12.md).
- **`shader`**: presets are **minimal lit / toon / dissolve** fragments (not full PBR); use **`shader.load`** for custom files. **`effect` / `camera.fx`**: still **stub** until a render-graph style post chain exists.
- **Raylib parity**: not every `raylib-go` top-level function is wrapped as a foreign; see `raylib_parity.json` (`in_raylib_not_in_bindings_raylib`). New game-relevant symbols are added in **tranches** — recent batches: **2D collision** helpers (`raylib_misc.go`), **rcamera** helpers + **GetCameraForward/Right/Up**, **DrawRectangleGradientH/V** (`raylib_shapes.go`). Remaining unbound entries are mostly rlgl/low-level, duplicates under other names, or niche APIs.

## Related docs

- [API_REFERENCE.md](../API_REFERENCE.md) — full binding tables by source file.
- [COMMAND_REFERENCE.md](COMMAND_REFERENCE.md) — task-oriented command groups.
- [ARCHITECTURE.md](ARCHITECTURE.md) — `RegisterAll` order and module policy.
