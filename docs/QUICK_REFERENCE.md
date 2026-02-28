# CyberBasic Quick Reference

One-page syntax reference. Names are **case-insensitive**.

## Variables and constants

```basic
VAR x = 10
LET y = 20
DIM name$ AS String
DIM a[10]

CONST Pi = 3.14159
CONST MaxLives = 3

VAR opt = Nil   // or Null; use IsNull(opt) or opt = Nil to check
```

## Types and dot notation

```basic
TYPE Player
    x AS Float
    y AS Float
    health AS Integer
END TYPE

VAR p = Player()
p.x = 100
p.y = 200
```

## Enums

```basic
ENUM Color : Red, Green, Blue
ENUM State : Idle = 0, Walk = 1, Jump = 2
```

## Control flow

```basic
IF condition THEN
    ...
ELSE
    ...
ENDIF

WHILE condition
    ...
WEND

FOR i = 1 TO 10
    ...
NEXT i

FOR i = 10 TO 1 STEP -1
    ...
NEXT i

REPEAT
    ...
UNTIL condition

SELECT CASE value
    CASE 1 : ...
    CASE 2 : ...
    CASE ELSE : ...
END SELECT

EXIT FOR
EXIT WHILE
```

## Functions and subs

```basic
FUNCTION Add(a, b)
    RETURN a + b
END FUNCTION

SUB DrawPlayer()
    // no return
END SUB
```

## Modules (namespaces)

```basic
MODULE Math3D
    FUNCTION Dot(a, b)
        RETURN a.x * b.x + a.y * b.y + a.z * b.z
    END FUNCTION
END MODULE

VAR d = Math3D.Dot(v1, v2)
```

## Compound assignment

```basic
x += 1
y -= 2
a *= 3
b /= 4
```

## Comments

```basic
// full-line comment
PRINT "hi"   // inline comment
```

Block comment `/* ... */` is supported where available.

## Includes

```basic
#include "other.bas"
```

Path is relative to the current file; one per line.

## Window and game loop (minimal)

Use **WHILE NOT WindowShouldClose() ... WEND** (or **REPEAT ... UNTIL WindowShouldClose()**) for the main loop. No auto-wrap; code compiles as written (DBPro-style). **DeltaTime()** for frame delta.

```basic
InitWindow(800, 600, "Title")
SetTargetFPS(60)

WHILE NOT WindowShouldClose()
    ClearBackground(20, 20, 30, 255)
    // draw here; use DeltaTime() for frame-based movement
WEND

CloseWindow()
```

## Hybrid loop (optional)

Define **`update(dt)`** and/or **`draw()`** (Sub or Function) and use a game loop with an **empty body**; the compiler injects GetFrameTime, physics step, update(dt), ClearRenderQueues, draw(), and FlushRenderQueues. You do not call BeginDrawing/EndDrawing yourself. See [Program Structure](PROGRAM_STRUCTURE.md#hybrid-updatedraw-loop).

## Input (common)

```basic
IsKeyDown(KEY_W)
IsKeyPressed(KEY_SPACE)
GetAxisX()   // -1, 0, or 1 for A/D
GetAxisY()   // -1, 0, or 1 for W/S
GetMouseX()
GetMouseY()
```

## Events (optional)

```basic
ON KeyDown("ESCAPE")
    // handle
END ON

ON KeyPressed("SPACE")
    // handle
END ON
```

Call `PollInputEvents()` in your game loop for events to run.

## Coroutines (optional)

```basic
StartCoroutine MySub()
Yield
WaitSeconds(1.0)
```

## Namespaces (no new syntax)

- **Raylib:** `InitWindow`, `DrawCircle`, etc. or `RL.InitWindow`, …
- **2D physics:** `CreateWorld2D`, `Step2D`, `CreateBody2D`, …
- **3D physics:** `CreateWorld3D`, `Step3D`, `CreateBox3D`, …
- **Game helpers:** `GAME.MoveWASD`, `GAME.CameraOrbit`, …

See [API_REFERENCE.md](../API_REFERENCE.md) for the full list.
