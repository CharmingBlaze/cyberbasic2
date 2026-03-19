# CyberBasic 2 — Refactor Master Plan

**Document:** living roadmap for v2-era features.  
**How to use:** Complete phases in order. After each phase: `go build ./...`, `go test ./...`, run regression examples, then commit.

### Phase progress (check off as you go)

| Done | Phase | Name |
|------|-------|------|
| [ ] | 1 | Repo hygiene + project rules |
| [ ] | 2 | Error system rewrite |
| [ ] | 3 | Dot notation layer |
| [ ] | 4 | Window API (optional boilerplate) |
| [ ] | 5 | High-level Physics API |
| [ ] | 6 | High-level Audio API |
| [ ] | 7 | Input action mapping |
| [ ] | 8 | Asset manager |
| [ ] | 9 | Scene system (declarative) |
| [ ] | 10 | String interpolation + TRY/CATCH |
| [ ] | 11 | Shader + post-FX API |
| [ ] | 12 | AI / Navmesh API |
| [ ] | 13 | Live reload + REPL |
| [ ] | 14 | Examples + docs overhaul |

### Phase dependencies (high level)

- **1** — no code deps.
- **2** — used by later phases for rich errors (e.g. assets “did you mean”).
- **3** — prerequisite for **4**, **5**, **6**, **8** (dot objects), **11** (ShaderObject).
- **4** — pairs with **3** for `WINDOW.*`; enables implicit loop story.
- **9** — benefits from **4**, **5**, **8** (optional but realistic).
- **10–13** — mostly after **3**; **13** touches main + runtime heavily.
- **14** — last; assumes prior phases landed (or stub APIs where a phase is deferred).

---

## Baseline: codebase today (before Phase 1)

This section tracks **what already exists** so phases don’t assume a blank slate.

### Layout (simplified)

```
cyberbasic2/
├── compiler/
│   ├── errors/          # PrettyPrint + locerr (compile-time snippet); NOT CyberError yet
│   ├── lexer/, parser/, codegen (via compiler package)
│   ├── vm/              # bytecode, render queues, fibers, OpLoadEntityProp/OpStoreEntityProp
│   ├── runtime/         # game loop, StepFrame
│   └── bindings/
│       ├── raylib/      # graphics, input, audio, UI
│       ├── dbp/         # DBP-style commands (large surface)
│       ├── box2d/, bullet/
│       ├── game/, net/, ecs/, scene/, navigation/, terrain/, world/, ...
│       ├── std/, sql/, nakama/, model/, objects/, procedural/, water/, vegetation/, indoor/, aseprite/
│       └── ...
├── engine/              # separate Makefile only (not root)
├── examples/
└── main.go              # CLI, compile+run, binding registration
```

### Registration (no single `registry.go`)

- Each binding package exposes `Register*(v *vm.VM)` (or similar).
- **[main.go](main.go)** imports packages and calls those registrars (duplicated in two compile/run paths). New bindings: add import + both registration blocks.

### Already partially aligned with later phases

| Area | Today | Plan phase |
|------|--------|------------|
| Compile errors | `compiler/errors` + `PrettyPrint` | **2** — add `CyberError`, codes, `Format`, `Nearest`; unify runtime errors |
| Entity `entity.prop` | `OpLoadEntityProp` / `OpStoreEntityProp`, getters/setters | **3** — general `DotObject` + `OpGetProp` / `OpSetProp` / `OpCallMethod` (additive) |
| Scenes (flat API) | `compiler/bindings/scene` — CreateScene, LoadScene, … | **9** — declarative `SCENE` / `END SCENE` + stack; **coexist** with existing commands |
| Navigation | `compiler/bindings/navigation` — grids, NavMesh, agents | **12** — extend or wrap; avoid duplicating A* if sufficient |
| REPL | `cyberbasic --repl` → `runREPL()` in main | **13** — no-args REPL, buffer + state, CLEAR; optional **breaking** change below |
| No filename | Tries default `examples/run_3d_physics_demo.bas`, else usage | **13** — target: no-args → REPL (document migration for anyone relying on default demo) |

### Regression examples (current repo)

These **.bas** files exist and should stay green after every phase:

- `examples/hello_world.bas`
- `examples/first_game.bas`
- `examples/platformer.bas`
- `examples/ui_demo.bas`
- `examples/input_debug.bas`

Also present: `examples/test_*.bas`, `examples/test_lang_features.bas` — include in `make examples` if practical.

### Root Makefile / ARCHITECTURE / scripts

- **Root `Makefile`:** not present yet (only `engine/Makefile`). Phase 1 adds root targets.
- **`docs/ARCHITECTURE.md`:** not present yet. Phase 1 adds it.
- **Run scripts:** `run_demo.ps1`, `run_demo.sh`, `run_3d_demo.bat` may still be at repo root until Phase 1 moves them to `scripts/`.

---

## Phase overview

| Phase | What | Risk |
|-------|------|------|
| 1 | Repo hygiene + project rules | None |
| 2 | Error system rewrite | Low |
| 3 | Dot notation layer | Medium |
| 4 | Window API — optional boilerplate | Medium |
| 5 | High-level Physics API | Medium |
| 6 | High-level Audio API | Low |
| 7 | Input action mapping | Low |
| 8 | Asset manager | Medium |
| 9 | Scene system (declarative + stack) | High |
| 10 | String interpolation + TRY/CATCH | Medium |
| 11 | Shader + post-FX API | Medium |
| 12 | AI / Navmesh API | High |
| 13 | Live reload + REPL | High |
| 14 | Examples + docs overhaul | Low |

---

## Phase 1 — Repo hygiene + project rules

**Goal:** Clean the repo and document design and conventions for future work.

### Tasks

- Delete debug artifacts from root if present: `e1.txt`, `e2.txt`, `err.txt`, `o1.txt`, `o2.txt`, `out.txt`, `debug_tokens.bas`.
- Move `run_3d_demo.bat`, `run_demo.ps1`, `run_demo.sh` into `scripts/` (create folder if needed).
- Update references in `README.md` and any docs that mention those paths.
- Add project conventions file at repo root: **`CONTRIBUTING.md`** or **`PROJECT_RULES.md`** (pick one).
- Create **`docs/ARCHITECTURE.md`** (compiler pipeline, binding registration, game loop, render queue, loop modes after Phase 4, how to add a binding).
- Add **root `Makefile`** with: `build`, `test`, `run`, `clean`, `examples`.

### Project conventions (put in CONTRIBUTING.md or PROJECT_RULES.md)

- **Stack:** Go 1.22+. Raylib (raylib-go), Box2D, Bullet. Bytecode VM — do not replace with tree-walk interpreter.
- **Naming:** Go packages lowercase, no underscores. BASIC: ALLCAPS or dot notation. Go types PascalCase. Errors: plain English + line + suggested fix.
- **Principles:** (1) Never remove existing BASIC commands; only add. (2) New features optional; no breaking `.bas` programs. (3) Dot notation layers on flat commands. (4) Rich errors. (5) Smart defaults. (6) InitWindow, CloseWindow, mainloop, SYNC are sacred.
- **Architecture:** New APIs under `compiler/bindings/<name>/`. Register in **main.go** (both code paths). VM in `compiler/vm/`; loop in `compiler/runtime/`. Keep **main.go** as CLI/orchestration only.
- **Testing:** `_test.go` per new binding; table-driven lexer/parser tests; examples as integration checks.
- **Avoid:** VM rewrite; custom GC; unnecessary CGo; renaming/removing .bas commands; mandatory new syntax.

### Makefile targets

- `build`: `go build -o cyberbasic .` (on Windows, document `cyberbasic.exe` or use conditional).
- `test`: `go test ./...`
- `run`: `FILE=path/to/file.bas ./cyberbasic $(FILE)` — document PowerShell/cmd equivalents.
- `clean`: remove `cyberbasic`, `cyberbasic.exe`, and other build artifacts you introduce.
- `examples`: compile (and optionally run) each `examples/*.bas`; report pass/fail.

**Phase 1 constraint:** Do not change `.go` compiler/VM/runtime logic—only hygiene, docs, Makefile, scripts, conventions file.

---

## Phase 2 — Error system rewrite

**Goal:** Runtime and compile-time errors consistently report: what went wrong + line/context + suggested fix.

### What to build

Extend **`compiler/errors/`** (today: `PrettyPrint` + locerr for compile errors):

- **`CyberError`:** `Code`, `Message`, `Line`, `Column`, `Snippet`, `Suggestion`.
- **`Format(filename)`** (or similar) producing:
  ```
  Error on line 14 in platformer.bas:
    VAR body = PHYSICS.DYNAMIC.BOX(100, 100, 32, 64)

  Physics body created before a physics world exists.
  Fix: add PHYSICS.WORLD() before creating any physics bodies,
       or use PHYSICS.SIMPLE() for automatic world setup.
  ```
- **`Nearest(key, knownKeys []string)`** — e.g. Levenshtein ≤ 2 for “did you mean”.

### Error codes (minimum set)

| Code / situation | Message | Suggestion |
|------------------|---------|------------|
| Physics body before world | Physics body created before a physics world exists | Add PHYSICS.WORLD() before this line |
| Asset missing | Asset '{key}' not found | Did you mean '{closest}'? Ensure ASSETS.LOAD() ran first |
| BeginDraw outside loop | BEGINDRAW called outside a draw context | Use mainloop/endmain or ON DRAW |
| Undefined variable | Variable '{name}' used before it was declared | Add VAR {name} = ... before line {line} |
| Type mismatch | Expected {type}, got {actual} for argument {n} | Use STR(), INT(), or FLOAT() |
| Missing end keyword | IF on line {n} has no matching ENDIF | Add ENDIF after your IF block |

### Implementation tasks

- Add `ErrorCode` iota: `ErrPhysicsBodyBeforeWorld`, `ErrAssetNotFound`, `ErrBeginDrawOutsideLoop`, `ErrUndefinedVariable`, `ErrTypeMismatch`, `ErrMissingEndKeyword`, …
- Integrate with existing **`PrettyPrint`** where it still applies (compile path); avoid two incompatible error stories long-term.
- Replace `fmt.Errorf` / ad-hoc prints in **`compiler/vm/`** and **`compiler/runtime/`** with structured errors + `Format`; **do not** change *when* errors trigger, only *how* they surface.
- Add **`compiler/errors/errors_test.go`** — table-driven tests for `Format()` and `Nearest()`.

---

## Phase 3 — Dot notation layer

**Goal:** `body.velocity.x = 200` alongside flat commands like DBP `POSITION OBJECT …`. Both remain valid.

### Design

- Parser: `a.b.c = val` → `OpSetProp` path `["b","c"]`; read → `OpGetProp`; `a.m(args)` → `OpCallMethod`.
- VM: **`DotObject`** interface — `GetProp`, `SetProp`, `CallMethod` on `[]string` paths / method names.
- Non-dot values: return a rich error (Phase 2) with a fix suggestion.

### Relation to existing VM

- **`OpLoadEntityProp` / `OpStoreEntityProp`** stay for current `entity.prop` behavior; dot notation is an **additional** path for values that implement `DotObject` (and/or a unified dispatch layer if you refactor carefully).

### Later consumers

Physics bodies, 3D handles, lights, cameras, audio sources, **`WINDOW`** (Phase 4).

### Implementation tasks

- Parser + codegen for dot chains and assignment.
- VM opcodes + tests in e.g. **`compiler/vm/dot_test.go`**: `obj.x`, `obj.position.x`, `obj.velocity.x = 200`, `obj.anim.PLAY("run")`, unknown property error.

---

## Phase 4 — Window API (optional boilerplate)

**Goal:** `InitWindow` + `mainloop`/`SYNC`/`endmain` unchanged. Programs **without** them can still run (implicit window or console).

### Modes

- **EXPLICIT:** source contains `InitWindow(` → user controls window + loop.
- **IMPLICIT:** `ON UPDATE` or `ON DRAW` without `InitWindow(` → auto `InitWindow(1280, 720, "CyberBasic 2")` + driven loop.
- **CONSOLE:** neither → no window; PRINT to stdout.

### WINDOW dot object (needs Phase 3)

- R/W: `title`, `width`, `height`, `fullscreen`, `vsync`, `icon`, `targetfps`
- Read-only: `fps`, `deltatime`, `mousex`, `mousey`, `screenwidth`, `screenheight`

### Implementation tasks

- Detect mode at runtime startup from AST or pre-scan.
- Implement WINDOW properties (pre-loop vs in-loop semantics as in original spec).
- **Do not** alter existing InitWindow, CloseWindow, mainloop, SYNC, SetTargetFPS behaviour.
- Tests: **`compiler/runtime/window_test.go`** — explicit / implicit / console detection.

---

## Phase 5 — High-level Physics API

**Goal:** `PHYSICS.DYNAMIC.BOX(…)` style API; **all** existing Box2D flat commands unchanged.

### API surface (summary)

- World: `PHYSICS.WORLD()`, options map, `PHYSICS3D.WORLD()` (Bullet).
- Bodies: STATIC/DYNAMIC/KINEMATIC BOX/CIRCLE/POLYGON factories returning **DotObject**.
- Properties: `position`, `velocity`, `rotation`, `friction`, `bounciness`, `mass`, `active`, `tag`; methods `ApplyForce`, `ApplyImpulse`, `Destroy`.
- `PHYSICS.RAYCAST`, `PHYSICS.OVERLAP_CIRCLE`; optional collision func callbacks.

### Implementation tasks

- New package **`compiler/bindings/physics2d/`** wrapping **`compiler/bindings/box2d`**.
- IMPLICIT: auto-world on first use; EXPLICIT: require `PHYSICS.WORLD()` first → **CyberError** if violated.
- Register in **main.go**; **`physics2d_test.go`**.

---

## Phase 6 — High-level Audio API

**Goal:** `AUDIO.SOUND`, music crossfade, buses, pools — additive over Raylib audio.

### Implementation tasks

- Package **`compiler/bindings/audiosys/`**; DotObject sound handles; linear crossfade; register + **`audiosys_test.go`**.

---

## Phase 7 — Input action mapping

**Goal:** `INPUT.MAP` + `INPUT.PRESSED` / `HELD` / `RELEASED` / `AXIS` / `REMAP`.

### Implementation tasks

- Package **`compiler/bindings/inputmap/`**; per-frame `ActionState` in **StepFrame**; preserve `IsKeyDown` / `IsKeyPressed` / mouse APIs; **`inputmap_test.go`**.

---

## Phase 8 — Asset manager

**Goal:** `ASSETS.LOAD`, `ASSETS["key"]`, async load, unload — with **Nearest** on missing keys.

### Implementation tasks

- Package **`compiler/bindings/assets/`**; runtime-visible manager for Phase 9; **`assets_test.go`**.

---

## Phase 9 — Scene system (declarative)

**Goal:** `SCENE` … `END SCENE` with `ON ENTER` / `ON UPDATE` / `ON DRAW` / `ON EXIT`; stack: `SCENES.PUSH` / `POP` / `SWAP` / `RESET`.

### Coexistence

- Existing **`CreateScene` / `LoadScene`** (package **scene**) remain; declarative scenes are **additional**. Name new foreign functions carefully (e.g. `ScenesPush`) vs parser keywords `SCENES.PUSH`.

### Implementation tasks

- Parser AST (e.g. `SceneDecl`); runtime **SceneManager** + stack; StepFrame uses top Update, stacked Draws; fallback when no `SCENE` blocks; **`scene_test.go`** (or under `compiler/runtime`).

---

## Phase 10 — String interpolation + TRY/CATCH

**Goal:** `$"Hello {name}"` and `TRY` / `CATCH` / `FINALLY` / `END TRY`.

### Implementation tasks

- Lexer token for `$"..."` with `{expr}` segments; parser node; VM concat via existing string conversion.
- Parser `TryStmt`; VM: **recover** path mapping **CyberError** to catch binding (dot object: message, line, code); FINALLY always runs; nesting tests.

---

## Phase 11 — Shader + post-FX API

**Goal:** `SHADER.LOAD` / PBR / TOON / DISSOLVE; `EFFECT.*`; camera FX list; `TWEEN` in StepFrame.

### Implementation tasks

- **`compiler/bindings/shadersys/`**, `builtins.go` for embedded GLSL strings; render-texture post chain; **`shadersys_test.go`**.

---

## Phase 12 — AI / Navmesh API

**Goal:** `NAVMESH.BUILD`, `AI.AGENT`, steering, **BTREE** syntax + runner.

### Implementation tasks

- **`compiler/bindings/aisys/`**; v1 grid/waypoint A* acceptable; **reuse** **`navigation`** where possible instead of second copy of pathfinding; parser additions for BTREE; **`aisys_test.go`**.

---

## Phase 13 — Live reload + REPL

**Goal:** `--dev` + fsnotify + `runtime.HotSwap`; **no-args** runs full REPL with buffer + persistent globals + CLEAR/EXIT.

### Live reload

- Watch `.bas`; re-parse; on success swap update/draw/function bodies under **RWMutex**; on structural change fast-restart; on failure print error, keep last good version.

### REPL — current vs target

- **Today:** `cyberbasic --repl` only; no file may run a default demo or print usage.
- **Target:** `cyberbasic` with **no args** starts REPL (per original spec). **Breaking:** document in CHANGELOG/README; optional `cyberbasic --demo` if you need to preserve old default-demo behaviour.

### Implementation tasks

- Add **fsnotify** dependency; **`runtime.HotSwap`**; extend **`runREPL`** or replace with spec’d behaviour; integration test for input sequence.

---

## Phase 14 — Examples + docs overhaul

**Goal:** One runnable example per major feature; docs extended (not removed).

### New / updated examples (from plan)

`hello_world` (minimal implicit hello when Phase 4+ exist), `implicit_loop`, `explicit_loop`, `physics_2d`, `physics_3d`, `input_actions`, `scenes`, `assets`, `shaders`, `audio`, `string_interp`, `try_catch`, `ai_patrol`, `repl_intro`, `pong`, `platformer_full`.

**Note:** `hello_world.bas` already exists — update in place or add a second file only if you need to keep both “console hello” and “one-line window hello.”

### Docs

- `README.md`, `LANGUAGE_SPEC.md`, `API_REFERENCE.md`, `docs/GETTING_STARTED.md`, `CHEATSHEET.md` — extend with new syntax and APIs.

---

## Testing strategy

After every phase (use Makefile once Phase 1 lands):

```bash
make test          # go test ./...
make build         # go build -o cyberbasic .
make examples      # all examples/*.bas
```

Until Makefile exists: run the same commands manually.

### Regression baseline

Must compile and run:

- `examples/hello_world.bas`
- `examples/first_game.bas`
- `examples/platformer.bas`
- `examples/ui_demo.bas`
- `examples/input_debug.bas`

---

## Commit strategy

One commit per phase:

```
phase/01-repo-hygiene
phase/02-error-system
phase/03-dot-notation
phase/04-window-api
phase/05-physics-api
phase/06-audio-api
phase/07-input-actions
phase/08-asset-manager
phase/09-scene-system
phase/10-string-interp-try-catch
phase/11-shader-api
phase/12-ai-pathfinding
phase/13-live-reload-repl
phase/14-examples-docs
```

- Tag **`v2.0.0-alpha`** after Phase 10.
- Tag **`v2.0.0`** after Phase 14.

---

## Success criteria

After all phases, **both** styles work:

**Implicit-style (after Phases 3–7+):**

```basic
WINDOW.TITLE = "Hello"
ON UPDATE
    IF INPUT.PRESSED("quit") THEN CLOSEWINDOW()
END ON
ON DRAW
    CLEARBACKGROUND(20, 20, 30, 255)
    DRAWTEXT($"FPS: {WINDOW.FPS}", 10, 10, 20, 255, 255, 255, 255)
END ON
```

**Explicit-style (unchanged from v1 behaviour):**

```basic
InitWindow(800, 600, "Hello")
SetTargetFPS(60)
mainloop
    ClearBackground(20, 20, 30, 255)
    DrawText("Hello!", 10, 10, 20, 255, 255, 255, 255)
    SYNC
endmain
CloseWindow()
```

Both are valid CyberBasic 2.
