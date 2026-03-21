# CyberBasic 2 — codebase architecture

How the compiler, VM, runtime, and bindings fit together. For a broader historical analysis see [MODULAR_COMPILER_ARCHITECTURE.md](MODULAR_COMPILER_ARCHITECTURE.md).

## Pipeline (source → run)

1. **Lexer** (`compiler/lexer`) — token stream.
2. **Parser** (`compiler/parser`) — AST.
3. **Semantic pass** (`compiler/semantic`) — types, entities, user functions.
4. **Codegen** (`compiler/codegen`) — bytecode (`vm.Chunk`): opcodes, constants, foreign call names.
5. **VM** (`compiler/vm`) — stack machine, `OpCallForeign`, `OpGetProp` / `OpSetProp` / `OpCallMethod` for **DotObject** handles.
6. **Runtime** (`compiler/runtime`) — optional implicit window + `ON UPDATE` / `ON DRAW` loop, render/physics orchestration.

## Compiler engine (staged API)

The package [`compiler`](../compiler/compiler.go) exposes the pipeline as **explicit stages** for tools, tests, and future passes (LSP, diagnostics, IR) without duplicating lexer/parser wiring.

| Stage | API | Output |
|--------|-----|--------|
| Lex | `(*Compiler).Tokenize(source)` | `[]lexer.Token` |
| Parse | `ParseTokens(tokens)` or `Parse(source)` | `*parser.Program` (reuse tokens from tooling) |
| Semantic | `(*Compiler).Analyze(program)` | `*semantic.Result` |
| Codegen | `codegen.Emit(program, sem)` | `*vm.Chunk` |

**Single full build path:** `Compile` and `CompileWithOptions` both run the same internal **`fullPipeline`** (lexer → parser → semantic → codegen). Staged methods reuse the same lexer/parser/semantic packages but are for partial runs; they cannot diverge on “what full compile means.”

**How to extend:** Add a new pass in a small package that takes the AST or `semantic.Result` and returns an augmented structure or errors only. Wire it from your driver by calling stages explicitly (e.g. `Parse` → `YourPass` → `Analyze` → `Emit`). Keep **lexer/parser free of `vm` imports**; keep **semantic** free of codegen where possible.

**Roadmap (not required for core builds):** split oversized `codegen` files into subpackages; central table for foreign-name lowering; richer source spans on errors in `compiler/errors`.

**Tests:** [`compiler/compiler_test.go`](../compiler/compiler_test.go) and staged tests in [`compiler/compiler_stages_test.go`](../compiler/compiler_stages_test.go); per-package tests under `lexer/`, `parser/`, `semantic/`, `codegen/`.

## Binding registration (`RegisterAll`)

All foreign functions and global DotObject namespaces are registered in **one place**:

- Go API: **`cyberbasic/compiler/bindings.RegisterAll(v *vm.VM, opts RegisterOptions) error`**
- File: [`compiler/bindings/register.go`](../compiler/bindings/register.go)

`RegisterOptions`:

- **`Source`** — full program text (or REPL session). Used so `physics2d` can set `RequireExplicitWorld` from `runtime.DetectWindowMode` (explicit `InitWindow` + main loop vs implicit DBP-style loop).
- **`SkipRaylib`** — if true, skips raylib registration, flush override, and renderer global hooks (for headless tests).

**Order** (do not reorder lightly; comments in `register.go` mirror this):

1. Raylib + flush override (when not skipped).
2. DBP runtime + renderer hooks (`SetDraw3D`, `SetPreDraw2D`, `SetVM`).
3. Bullet, Box2D, high-level `physics2d`.
4. ECS, net, Nakama, scene, game.
5. DBP 2D overlay, SQL, terrain + DBP terrain overlay, objects + DBP DrawObject overlay, procedural, water + DBP water, vegetation, world, navigation, indoor.
6. Std (global `std`), audiosys, inputmap, assets, **shadersys** (raylib `LoadShader` + embedded presets), **effect** (stub FX + `effect` global), **cameradot** (`camera.fx`), **tween** (`TweenRegister` + frame **Tick** in `runtime`), **aisys** (`ai.*` → `navigation.*` facade), **windowdot** (`window`), **engine** (composition: `engine.<subsystem>` → same global as lowercase name).
7. Reset `physics2d` world flags from `opts.Source`.

Application code should call **`RegisterAll`** after `LoadChunk`, **`std.RegisterEnums(chunk.Enums)`**, and **`vm.SetRuntime(rt)`** — see [`internal/app`](../internal/app).

## Package boundaries (Go)

| Area | Responsibility |
|------|------------------|
| `compiler/vm` | Opcodes, stack, `DotObject`, no game rules |
| `compiler/runtime` | Frame loop, window mode detection, implicit handlers |
| `compiler/bindings/*` | One concern per package; `RegisterX(*vm.VM)` |
| `compiler/errors` | `CyberError`, compile-time `PrettyPrint` |
| `internal/app` | CLI: flags, compile, `RegisterAll`, run, REPL |

Avoid import cycles; keep `errors` free of `vm` imports.

## Module API (v2 style) vs legacy flat commands

**Legacy / flat:** `PhysicsHighWorld`, `CreateWorld2D`, `InputMapRegister`, etc. — still supported.

**v2 modules:** Globals that implement `vm.DotObject` are registered with lowercase keys (VM stores globals case-insensitively). Use dotted calls and property syntax:

| Global (BASIC) | Role |
|----------------|------|
| `window` | `WINDOW.TITLE = "..."`, properties for size, FPS, etc. |
| `physics` | `physics.world(gx, gy)`, `physics.dynamicbox(...)`, `physics.staticbox(...)`, `physics.raycast2d(...)` |
| `audio` | `audio.load(path)`, `audio.playsoundid(id)` |
| `input` | `input.map.register(action, key)`, `input.map.pressed(action)`, … |
| `assets` | `assets.set`, `assets.get`, `assets.unload`, … |
| `shader` | PBR/toon/dissolve/load factories; shader handles with `set` / properties where implemented |
| `ai` | `ai.version()`, navigation aliases, `ai.agent(id$)` |
| `scenes` | `scenes.create`, `scenes.load`, `scenes.save2d`, … |
| `ecs`, `net`, `sql`, `nakama` | Thin `modfacade` namespaces → existing foreigns |
| `terrain`, `water`, `vegetation`, `world`, `navigation`, `indoor`, `procedural`, `objects` | Same pattern |
| `effect` | Stub post-process factories (`bloom`, `vignette`, `dof`) |
| `camera` | `camera.fx.add` / `camera.fx.clear` (stubs) |
| `tween` | `tween.register` → `TweenRegister`; **Tick** runs from `runtime.beginRuntimeFrame` |
| `std` | Small file/env/help facade; `Print` and math remain flat |
| `engine` | `engine.ecs`, `engine.shader`, … — returns the already-registered module object |

Codegen routes dotted calls that are not `rl` / `box2d` / `bullet` / `game` (and not the Box2D/Bullet namespace rewrite map) through **`OpCallMethod`** on the receiver object.

**Handles:** APIs may return values that implement `DotObject` (e.g. high-level physics bodies); use `handle.x = value` via `OpSetProp` where implemented.

## How to add a new binding

1. Implement **`RegisterYourFeature(*vm.VM)`** in `compiler/bindings/yourfeature/` (or extend an existing package).
2. Add **`yourfeature.RegisterYourFeature(v)`** to **`RegisterAll`** in the correct phase (core before overlays; document why in a comment if non-obvious).
3. If you expose a **namespace object**, call **`v.SetGlobal("name", dotObj)`** with a **lowercase** key; add the root name to **`dotObjectRoots`** in `compiler/codegen/dot_emit.go` if the root must always use DotObject semantics (not vector `.x` hacks).
4. Document the command in **`API_REFERENCE.md`** (and **`COMMAND_REFERENCE.md`** if user-facing).
5. Run **`go test ./...`** and **`make examples`** (or `go build` + `--lint` on `examples/*.bas`).

## Game loop modes

- **Explicit:** `InitWindow` + your loop / `mainloop` … `endmain` (see language docs).
- **Implicit:** `ON UPDATE` / `ON DRAW` without opening a window in the classic way; runtime may call `RunImplicitLoop` after top-level `Run` completes.

`DetectWindowMode(source)` classifies the program to tune physics “explicit world required” behavior.
