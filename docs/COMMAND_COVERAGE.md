# Command coverage — 2D/3D games and applications

This document is the **human-facing coverage matrix** for CyberBasic 2: what you need for games and desktop-style apps, where it lives, and how we track completeness.

## Machine-readable inventory

- **[`generated/foreign_commands.json`](generated/foreign_commands.json)** — every unique `RegisterForeign` name under `compiler/bindings/`, with source file paths (and `packages` when the same name is registered from more than one Go package).
- **[`generated/FOREIGN_COMMANDS_INDEX.md`](generated/FOREIGN_COMMANDS_INDEX.md)** — the same data as a browsable Markdown index.
- **[`generated/raylib_parity.json`](generated/raylib_parity.json)** — diff helper: top-level exported **functions** in `github.com/gen2brain/raylib-go/raylib` vs names registered from `compiler/bindings/raylib/`. It excludes methods, `rlgl`-level entry points, and symbols exposed only under different names or via DBP/game wrappers.

Regenerate locally (or in CI) with:

```bash
make foreign-audit
```

## Pillar matrix (games and apps)

Status legend: **Yes** = production path in tree; **Stub** = callable but minimal / placeholder; **Partial** = subset or expert-only.

| Pillar | What you use (examples) | v2 global (optional) | Primary bindings / docs | Smoke example |
|--------|-------------------------|----------------------|-------------------------|---------------|
| **Runtime — explicit loop** | `InitWindow`, `mainloop` / `SYNC`, `CloseWindow` | `window` (title, size, fps, …) | [`compiler/bindings/raylib`](../compiler/bindings/raylib), [COMMAND_REFERENCE — Game loop](COMMAND_REFERENCE.md#game-loop) | [`examples/first_game.bas`](../examples/first_game.bas) |
| **Runtime — implicit DBP** | `ON UPDATE` / `ON DRAW`, `WINDOW.TITLE` | `window` | [`compiler/runtime`](../compiler/runtime), [PROGRAM_STRUCTURE](PROGRAM_STRUCTURE.md) | [`examples/implicit_loop.bas`](../examples/implicit_loop.bas) |
| **Runtime — hybrid** | `update(dt)` / `draw()`, render queues | — | [RENDERING_AND_GAME_LOOP](RENDERING_AND_GAME_LOOP.md), [COMMAND_REFERENCE — hybrid](COMMAND_REFERENCE.md#game-loop-hybrid) | [`examples/platformer.bas`](../examples/platformer.bas) (pattern) |
| **2D drawing** | `DrawRectangle`, `DrawCircle`, `DrawText`, textures, layers | — | `raylib_shapes`, `raylib_text`, `raylib_textures`, [COMMAND_REFERENCE — 2D](COMMAND_REFERENCE.md) | [`examples/platformer.bas`](../examples/platformer.bas), [`examples/smoke_rectangle_gradient.bas`](../examples/smoke_rectangle_gradient.bas) |
| **2D collision (geometric)** | `CheckCollisionRecs`, `CheckCollisionCircleLine`, … | — | [`raylib_misc.go`](../compiler/bindings/raylib/raylib_misc.go), [COMMAND_REFERENCE — 2D geometric collision](COMMAND_REFERENCE.md#2d-geometric-collision-raylib-shapes) | [`examples/smoke_2d_collision.bas`](../examples/smoke_2d_collision.bas) |
| **3D — raylib primitives** | `BeginMode3D`, `DrawCube`, `DrawGrid`, `SetCamera3D`, `CameraMoveForward`, … | — | [`raylib_3d.go`](../compiler/bindings/raylib/raylib_3d.go) | [`examples/smoke_raylib_3d.bas`](../examples/smoke_raylib_3d.bas), [`examples/smoke_rcamera.bas`](../examples/smoke_rcamera.bas) |
| **3D — DBP / scene objects** | `LoadObject`, `DrawObject`, `PositionObject`, … | — | [`compiler/bindings/dbp`](../compiler/bindings/dbp), [3D_GAME_API](3D_GAME_API.md) | [`examples/first_game.bas`](../examples/first_game.bas) |
| **Input — raw** | `IsKeyDown`, `IsMouseButtonPressed`, gamepad | — | `raylib_input` | [`examples/input_debug.bas`](../examples/input_debug.bas) |
| **Input — actions** | `InputMapRegister`, `InputPressed` | `input.map.register` / `pressed` | [`inputmap`](../compiler/bindings/inputmap/inputmap.go) | [`examples/smoke_input_v2.bas`](../examples/smoke_input_v2.bas) |
| **Audio** | `InitAudioDevice`, `LoadSound`, `PlaySound`, music streams | `audio.*` | `raylib_audio`, [`audiosys`](../compiler/bindings/audiosys) | [`examples/smoke_audio.bas`](../examples/smoke_audio.bas) |
| **UI** | `raygui` controls, custom UI foreigns | — | `raylib_raygui`, `raylib_ui` | [`examples/ui_demo.bas`](../examples/ui_demo.bas) |
| **Physics 2D** | `CreateWorld2D`, `Step2D`, bodies/joints | `physics.*` | [`box2d`](../compiler/bindings/box2d), [`physics2d`](../compiler/bindings/physics2d) | [`examples/smoke_physics2d.bas`](../examples/smoke_physics2d.bas) |
| **Physics 3D** | `CreateWorld3D`, `Step3D`, … | — | [`bullet`](../compiler/bindings/bullet) | [`examples/smoke_physics3d.bas`](../examples/smoke_physics3d.bas) |
| **Assets (key/value)** | `AssetsSet` / `AssetsGet` | `assets.set` / `get` | [`assets`](../compiler/bindings/assets/assets.go) | [`examples/smoke_assets_v2.bas`](../examples/smoke_assets_v2.bas) |
| **Scenes** | `CreateScene`, `LoadScene`, … | `scenes.*` | [`scene`](../compiler/bindings/scene) | [`examples/smoke_scenes.bas`](../examples/smoke_scenes.bas) |
| **World building** | terrain, water, vegetation, world, nav, indoor, procedural | `terrain`, `water`, … | respective packages under `compiler/bindings/` | _(large-game)_ |
| **Std / apps** | files, JSON, HTTP, HELP | `std.*` (subset) | [`std`](../compiler/bindings/std) | [`examples/smoke_std.bas`](../examples/smoke_std.bas) |
| **Net / SQL / Nakama** | multiplayer and persistence | `net`, `sql`, `nakama` | respective packages | _(app-specific)_ |
| **Shaders / FX** | `shader.pbr` / `toon` / `dissolve` (embedded GLSL), `BeginShaderMode` + uniforms | `shader`, `effect`, `camera.fx` | [`shadersys`](../compiler/bindings/shadersys), [`effect`](../compiler/bindings/effect) | [`examples/shader_demo.bas`](../examples/shader_demo.bas); **effect / camera.fx** still stub — see below |
| **AI / behaviour** | `ai.*` → `navigation.*`; optional `ai.agent` handle | `ai`, `navigation` | [`aisys`](../compiler/bindings/aisys), [`navigation`](../compiler/bindings/navigation) | [`examples/ai_patrol.bas`](../examples/ai_patrol.bas); **BTREE** deferred — [BTREE_PHASE12.md](BTREE_PHASE12.md) |
| **Tween** | `TweenRegister`, frame tick | `tween.register`, `tween.count` | [`tween`](../compiler/bindings/tween), [`runtime/loop`](../compiler/runtime/loop.go) | [`examples/smoke_tween.bas`](../examples/smoke_tween.bas) |
| **Composition** | `engine.ecs`, `engine.net`, … | `engine` | [`engine`](../compiler/bindings/engine/engine.go) | [`examples/smoke_engine.bas`](../examples/smoke_engine.bas) |

## Stub and partial APIs (honest expectations)

- **`ai`**: **`ai.version()`** plus **navigation aliases** (`navgridcreate`, `navagentcreate`, …) and **`ai.agent(id$)`**; **BTREE** / full behaviour-tree syntax is **not** in-tree — see [BTREE_PHASE12.md](BTREE_PHASE12.md).
- **`shader`**: presets are **minimal lit / toon / dissolve** fragments (not full PBR); use **`shader.load`** for custom files. **`effect` / `camera.fx`**: still **stub** until a render-graph style post chain exists.
- **Raylib parity**: not every `raylib-go` top-level function is wrapped as a foreign; see `raylib_parity.json` (`in_raylib_not_in_bindings_raylib`). New game-relevant symbols are added in **tranches** — recent batches: **2D collision** helpers (`raylib_misc.go`), **rcamera** helpers + **GetCameraForward/Right/Up**, **DrawRectangleGradientH/V** (`raylib_shapes.go`). Remaining unbound entries are mostly rlgl/low-level, duplicates under other names, or niche APIs.

## Related docs

- [API_REFERENCE.md](../API_REFERENCE.md) — full binding tables by source file.
- [COMMAND_REFERENCE.md](COMMAND_REFERENCE.md) — task-oriented command groups.
- [ARCHITECTURE.md](ARCHITECTURE.md) — `RegisterAll` order and module policy.
