# DBP Parity Checklist

CyberBASIC 2 aims for a DarkBASIC Pro–style experience: minimal boilerplate, familiar commands, and "it just works" feel.

## Checklist

- [x] BASIC file with only `PRINT` + `WAITKEY` runs
- [x] BASIC file with no explicit window still opens a window for graphics (when using OnUpdate/OnDraw)
- [x] On Start, On Update, On Draw work without manual loop wiring
- [x] `LoadImage` + `Sprite` works with IDs (DBP-style integer IDs)
- [x] `LoadObject` / `LoadCube` + `PositionObject` + `RotateObject` / `YRotateObject` works
- [x] A simple 3D demo runs in under 30 lines of BASIC
- [ ] Errors point to the right line and don't feel cryptic (ongoing improvement)

## How to Verify

### 1. PRINT + WAITKEY

```basic
PRINT "Hello from CyberBASIC 2"
WaitKey()
```

Run: `cyberbasic examples/hello_world.bas`

### 2. Implicit Window + Loop

```basic
SUB OnStart()
  ' init
  UseUnifiedRenderer()
END SUB

SUB OnUpdate(dt AS FLOAT)
  ' per-frame logic
END SUB

SUB OnDraw()
  ClearBackground(0, 0, 0, 255)
  DrawCircle(400, 300, 50, 255, 100, 100, 255)
  SYNC
END SUB
```

No `InitWindow`, no `WHILE` loop—the runtime provides them. With `UseUnifiedRenderer`, call `SYNC` at the end of `OnDraw` to end the frame.

### 3. DBP-Style 2D

```basic
LoadImage("hero.png", 1)
Sprite(1, x, y)
Cls()
Ink(255, 0, 0)
```

### 4. DBP-Style 3D

```basic
LoadCube(1, 2)
PositionObject(1, 0, 0, 5)
YRotateObject(1, 1)
DrawObject(1)
```

### 5. TYPE + Dot Notation

```basic
TYPE Player
  x AS FLOAT
  y AS FLOAT
  name AS STRING
END TYPE

DIM p AS Player
p.x = 100
p.y = 200
p.name = "Hero"
PRINT p.name
```

## Examples

See [examples/README.md](../examples/README.md):

- `examples/hello_world.bas` – PRINT + WAITKEY
- `examples/first_game.bas` – implicit loop, OnStart/OnUpdate/OnDraw
- `templates/2d_game.bas` – 2D game with WASD movement
- `templates/3d_game.bas` – 3D game with LoadCube, PositionObject, YRotateObject
