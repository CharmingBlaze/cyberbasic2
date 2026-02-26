# CyberBasic Program Structure

This document summarizes program structure, comments, and the main language features.

## Table of Contents

1. [Comments](#comments)
2. [Feature list (implemented)](#feature-list-implemented)
3. [Block structure (quick reference)](#block-structure-quick-reference)
4. [Example skeleton](#example-skeleton)
5. [Hybrid update/draw loop](#hybrid-updatedraw-loop)

---

## Comments

Use **`//`** for line comments. Everything from `//` to the end of the line is ignored.

```basic
// This is a comment
VAR x = 10   // inline comment
PRINT x
```

---

## Feature list (implemented)

- **Variables:** `VAR`, `DIM`, `LET`; arrays `VAR a[10]`, `DIM b[5,5]`
- **Constants:** `CONST name = value`
- **Types:** `TYPE ... END TYPE`, `EXTENDS`
- **Enums:** `ENUM Name ... END ENUM` (named/unnamed, custom values); `Enum.getValue`, `Enum.getName`, `Enum.hasValue`
- **Control flow:** `IF/THEN/ELSE/ELSEIF/ENDIF`, `FOR/NEXT`, `WHILE/WEND`, `REPEAT/UNTIL`, `SELECT CASE`
- **Loop control:** `EXIT FOR`, `EXIT WHILE`, `BREAK FOR`, `BREAK WHILE`, `CONTINUE FOR`, `CONTINUE WHILE`
- **Procedures:** `SUB`, `FUNCTION`; `END SUB` / `ENDSUB`, `END FUNCTION` / `ENDFUNCTION`
- **Modules:** `MODULE name` / `END MODULE` / `ENDMODULE`
- **Operators:** `+ - * / % \` (integer div), `^` (power), `= <> < <= > >=`, `AND`, `OR`, `XOR`, `NOT`
- **Compound assign:** `+=`, `-=`, `*=`, `/=`
- **String/std:** `Left`, `Right`, `Mid`, `Substr`, `Instr`, `Upper`, `Lower`, `Len`, `Chr`, `Asc`, `Str`, `Val`, `Rnd`, `Rnd(n)`, `Random(n)`, `Int`
- **Assert:** `ASSERT condition [, message]`
- **Null:** `Nil`, `Null`, `None`; `IsNull(value)`
- **JSON/dict:** `LoadJSON`, `ParseJSON`, `GetJSONKey`, dict literal `{"key": value}` or `{key = value}`, `CreateDict`, `SetDictKey`, `Dictionary.has/keys/values/size/remove/clear/merge/get`
- **File I/O:** `ReadFile`, `WriteFile`, `DeleteFile`, `CopyFile`, `ListDir`
- **Includes:** `#include "file.bas"` (or `IMPORT "file.bas"`); path relative to current file
- **Events/coroutines:** `ON ... GOSUB`, `StartCoroutine`, `Yield`, `WaitSeconds`
- **Graphics:** raylib (2D/3D), Box2D, Bullet; automatic frame/mode wrapping in game loops
- **Multi-window:** `SpawnWindow`, `ConnectToParent`, `NET.*`
- **ECS, GUI, multiplayer:** See [ECS_GUIDE.md](ECS_GUIDE.md), [GUI_GUIDE.md](GUI_GUIDE.md), [MULTIPLAYER.md](MULTIPLAYER.md)

---

## Block structure (quick reference)

| Block        | Start        | End              |
|-------------|--------------|------------------|
| IF          | IF ... THEN  | ENDIF or END IF  |
| FOR         | FOR x = a TO b [STEP s] | NEXT   |
| WHILE       | WHILE cond   | WEND             |
| REPEAT      | REPEAT       | UNTIL cond       |
| SELECT CASE | SELECT CASE expr | ENDSELECT    |
| FUNCTION    | FUNCTION name(params) | ENDFUNCTION or END FUNCTION |
| SUB         | SUB name(params) | ENDSUB or END SUB   |
| MODULE      | MODULE name  | ENDMODULE or END MODULE |
| TYPE        | TYPE name    | ENDTYPE          |
| ENUM        | ENUM [name]  | ENDENUM or END ENUM |

---

## Example skeleton

```basic
// My game
#include "constants.bas"

ENUM GameState
    Menu, Playing, Paused
END ENUM

VAR state = 0
VAR config = {"width": 1024, "height": 768}

FUNCTION main()
    InitWindow(config["width"], config["height"], "Game")
    SetTargetFPS(60)
    WHILE NOT WindowShouldClose()
        // Update and draw
        IF state = 0 THEN
            // menu
        ENDIF
    WEND
    CloseWindow()
ENDFUNCTION

main()
```

---

## Hybrid update/draw loop

**When to use:** Prefer the hybrid loop for new games when you want automatic physics stepping and a clear split between update and draw. Use the manual loop when you need full control over the order of operations or legacy code.

If you define **`update(dt)`** and/or **`draw()`** (as Sub or Function) and use a game loop (`WHILE NOT WindowShouldClose()` or `REPEAT ... UNTIL WindowShouldClose()`), the compiler replaces the loop body with an automatic pipeline:

1. **GetFrameTime** → `dt`
2. **StepAllPhysics2D(dt)** and **StepAllPhysics3D(dt)** (all registered worlds)
3. **update(dt)** (if defined)
4. **ClearRenderQueues**
5. **draw()** (if defined) — all Draw*/Gui* calls inside `draw()` are **queued** (2D, 3D, GUI)
6. **FlushRenderQueues** — BeginDrawing, ClearBackground, then render queue2D, queue3D, queueGUI in order, EndDrawing

You do not call `BeginDrawing`/`EndDrawing` or `BeginMode2D`/`BeginMode3D` yourself; the engine does it. Example:

```basic
SUB update(dt)
  REM move player with dt
END SUB
SUB draw()
  ClearBackground(30, 30, 45, 255)
  DrawRectangle(x, y, 40, 40, 255, 100, 100, 255)
  DrawText("Hello", 20, 20, 20, 255, 255, 255, 255)
END SUB
WHILE NOT WindowShouldClose()
WEND
```

See **examples/hybrid_update_draw_demo.bas**. Scripts that do not define `update`/`draw` keep the previous behaviour (manual or compiler-wrapped Begin/End).

---

## See also

- [Documentation Index](DOCUMENTATION_INDEX.md)
- [Getting Started](GETTING_STARTED.md)
- [Game Development Guide](GAME_DEVELOPMENT_GUIDE.md)
- [Libraries and includes](LIBRARIES.md)
