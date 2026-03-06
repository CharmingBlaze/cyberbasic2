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

Run: `cyberbasic examples/dbp_style/hello_world.bas`

### 2. Implicit Window + Loop

```basic
SUB OnStart()
  ' init
END SUB

SUB OnUpdate(dt AS FLOAT)
  ' per-frame logic
END SUB

SUB OnDraw()
  ClearBackground(0, 0, 0, 255)
  DrawCircle(400, 300, 50, 255, 100, 100, 255)
END SUB
```

No `InitWindow`, no `WHILE` loop—the runtime provides them.

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

See `examples/dbp_style/`:

- `hello_world.bas` – PRINT + WAITKEY
- `2d_sprites.bas` – implicit loop, OnStart/OnUpdate/OnDraw
- `3d_cube_spin.bas` – LoadCube, PositionObject, YRotateObject
- `first_person_demo.bas` – 3D camera, WASD + mouse
- `simple_platformer.bas` – 2D platformer with implicit loop
- `type_dot_test.bas` – TYPE with dot notation
