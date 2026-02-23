# CyberBasic modular layout

## Overview

The codebase is split into clear modules so you can extend or replace parts without touching the rest. The language and raylib API are a **carbon copy** of [CharmingBlaze/cyberbasic](https://github.com/CharmingBlaze/cyberbasic): same syntax (VAR, LET, CONST, TYPE...END TYPE, dot notation, ENUM) and same **unprefixed** raylib-style commands (InitWindow, SetTargetFPS, IsKeyDown, BeginDrawing, EndDrawing, ClearBackground, DrawCircle, KEY_W, etc.).

## Packages

| Package | Role |
|---------|------|
| **main** | Entry point: CLI, file loading, compiles and runs BASIC. Registers all bindings with the VM. |
| **compiler** | Lexer, parser, AST, compiler (source → bytecode). Single package; internal files by concern. |
| **compiler/vm** | Bytecode VM: execution, stack, opcodes. Physics opcodes deprecated; use foreign calls. |
| **compiler/parser** | Parser and AST (parser.go, ast.go). |
| **compiler/lexer** | Tokenizer (lexer.go, token.go). |
| **compiler/gogen** | Optional Go code generator (--gen-go). |
| **compiler/runtime** | Game runtime: window, sync, physics bridge (used by legacy opcodes if any). |

## Bindings (foreign API)

All bindings register with the VM via `v.RegisterForeign("Name", fn)`. Main loads them in one place:

- **compiler/bindings/raylib** – Graphics, window, input, audio, shapes, text, textures, 3D, fonts, misc.  
  Split into: `raylib.go` (shared helpers + `RegisterRaylib`), `raylib_core.go`, `raylib_shapes.go`, `raylib_textures.go`, `raylib_text.go`, `raylib_fonts.go`, `raylib_input.go`, `raylib_audio.go`, `raylib_3d.go`, `raylib_misc.go`.
- **compiler/bindings/box2d** – 2D physics (BOX2D.*). Single file `box2d.go`.
- **compiler/bindings/bullet** – 3D physics (BULLET.*). Single file `bullet.go`.

In BASIC, call with or without namespace: `InitWindow(800, 450, "Title")` or `RL.InitWindow(...)`, and `BOX2D.CreateWorld("w", 0, -10)`. All raylib functions and constants use the **same names** as the raylib C API (unprefixed): e.g. `InitWindow`, `SetTargetFPS`, `IsKeyDown`, `BeginDrawing`, `EndDrawing`, `ClearBackground`, `DrawCircle`, `DrawText`, `WindowShouldClose`, `CloseWindow`. Key constants: `KEY_W`, `KEY_A`, `KEY_S`, `KEY_D`, `KEY_SPACE`, `KEY_ESCAPE`, `KEY_UP`, `KEY_DOWN`, etc. (0-arg foreigns returning the key code). Color constants: `White`, `Black`, `Red`, `Gray`, etc.

## Adding a new binding module

1. Create `compiler/bindings/<name>/<name>.go` in package `name`.
2. Implement `func Register<Name>(v *vm.VM)` and call `v.RegisterForeign("NAMESPACE.FuncName", func(args []interface{}) (interface{}, error) { ... })`.
3. In `main.go`, import the package and call `name.RegisterName(rt.GetVM())` next to `raylib.RegisterRaylib`, `box2d.RegisterBox2D`, `bullet.RegisterBullet`.

No changes to lexer/parser/compiler are needed; the compiler emits `OpCallForeign` for any unknown name (and strips an optional `rl.` prefix).

## Examples

- **examples/first_game.bas** – First game: InitWindow, SetTargetFPS, WHILE NOT WindowShouldClose(), IsKeyDown(KEY_W), BeginDrawing, ClearBackground, DrawCircle, EndDrawing, CloseWindow.
- **examples/simple_box2d_demo.bas** – BOX2D only, no POP().
- **examples/box2d_demo.bas** – BOX2D + raylib window, click to spawn boxes.
- **examples/run_3d_physics_demo.bas** – BULLET + raylib 3D.

Use `BOX2D.*` and `BULLET.*` for physics; legacy physics keywords and opcodes are deprecated.

## Language reference (modern raylib BASIC)

- **Dynamic typing:** Variables are dynamic. `DIM x`, `DIM y AS Float`, and `LET` store values as `interface{}` in the VM. `AS Type` is an optional hint (initial value and optional types like Vector2/Body only); no static type checking.
- **Variables:** `VAR x = 10` is an alias for `LET x = 10` (assign and create if needed). `DIM x`, `DIM y AS Float`, and `LET` also declare/assign.
- **Constants:** `CONST name = expression` (e.g. `CONST Pi = 3.14159`). Names are in-scope for the program; use as normal identifiers (compiled as `OpLoadConst`). For key codes use built-in `KEY_W`, `KEY_A`, etc.
- **Enums:** `ENUM Name : member1, member2 = 5, member3` — members are integer constants (auto-increment from 0, or explicit `= expr`). Use enum members like any constant (e.g. `LET c = Green`).
- **Control flow:** `IF/THEN/ELSE/ENDIF`, `FOR/NEXT`, `WHILE/WEND`, `REPEAT/UNTIL`, `SELECT CASE/END SELECT`. Blocks can contain assignments, calls (e.g. `Sin`, `RL.InitWindow`), and nested control flow.
- **Dot notation:** `expr.member` compiles to a getter: `.x` / `.y` → `GetVector2X`/`GetVector2Y`, `.z` → `GetVector3Z`; other members → `GetVector2&lt;Member&gt;`. Namespace “constants”: `RL.DarkGray`, `BOX2D.*`, `BULLET.*`, `GAME.KEY_W` etc. compile as 0-arg foreign calls. Qualified calls: `RL.InitWindow(...)`, `BULLET.Step(...)`.
- **Namespaces:** `RL` (raylib), `BOX2D`, `BULLET`, `GAME` (game helpers). Raylib calls use the same names as the C API; call with or without prefix: `InitWindow(...)` or `RL.InitWindow(...)`. Use `BeginDrawing()` and `EndDrawing()` when not inside Main() or auto-wrapped WHILE. Comments: **only** `//` (line) and `/* */` (block). `NOT WindowShouldClose()` works in loop conditions.
- **Compound assignment:** `+=`, `-=`, `*=`, `/=` (e.g. `x += 1`).
- **Loop exit:** `EXIT FOR`, `EXIT WHILE` (jump to after the current loop).
- **Minimal 2D/3D:** Use `RL.*` for window and drawing, `BULLET.*` for 3D physics, `GAME.CameraOrbit` / `GAME.MoveWASD` / `GAME.OnGround` for camera and movement; see examples (e.g. mario64-style).
