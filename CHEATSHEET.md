# CyberBASIC2 – First 10 lines

Minimal snippets to get a game running.

## v2 module namespaces (optional)

Same engine as flat commands; dotted style for readability. Keys are case-insensitive.

| Module | Example |
|--------|---------|
| **window** | `WINDOW.TARGETFPS = 60` |
| **physics** | `physics.world(0, 9.8)` · `VAR b = physics.dynamicbox(100, 300, 40, 40)` |
| **audio** | `VAR s = audio.load("hit.wav")` · `audio.playsoundid(s)` |
| **input** | `input.map.register("fire", 32)` · `IF input.map.pressed("fire") THEN ...` |
| **assets** | `assets.set("hero", texId)` · `VAR t = assets.get("hero")` |
| **scenes** | `scenes.create("level1")` · `scenes.load("level1")` |

See [API_REFERENCE.md](API_REFERENCE.md#module-api-v2-style-vs-legacy-flat-names).

## DBP-style 2D (zero boilerplate)

No `InitWindow`, no `WHILE` loop. See [docs/DBP_PARITY.md](docs/DBP_PARITY.md).

```basic
VAR x = 400
VAR y = 300

SUB OnStart()
  UseUnifiedRenderer
END SUB

SUB OnUpdate(dt)
  LET x = x + 100 * dt * GetAxisX()
  LET y = y + 100 * dt * GetAxisY()
END SUB

SUB OnDraw()
  ClearBackground(20, 20, 30, 255)
  DrawCircle(x, y, 30, 255, 100, 100, 255)
  SYNC
END SUB
```

## Manual 2D game (move a circle)

```basic
InitWindow(800, 600, "My Game")
SetTargetFPS(60)
VAR x = 400
VAR y = 300
mainloop
  VAR dt = DeltaTime()
  LET x = x + 100 * dt * GetAxisX()
  LET y = y + 100 * dt * GetAxisY()
  ClearBackground(20, 20, 30, 255)
  DrawCircle(x, y, 30, 255, 100, 100, 255)
  SYNC
endmain
CloseWindow()
```

## 3D game (Bullet + orbit camera)

```basic
InitWindow(1024, 600, "3D")
SetTargetFPS(60)
DisableCursor()
CreateWorld3D("w", 0, -18, 0)
CreateSphere3D("w", "player", 0, 0.5, 0, 0.5, 1)
CreateBox3D("w", "ground", 0, -0.5, 0, 12.5, 0.5, 12.5, 0)
VAR camAngle = 0
REPEAT
  VAR dt = GetFrameTime()
  Step3D("w", dt)
  GAME.MoveWASD("w", "player", camAngle, 120, 9, dt)
  VAR px = GetPositionX3D("w", "player")
  VAR py = GetPositionY3D("w", "player")
  VAR pz = GetPositionZ3D("w", "player")
  GAME.CameraOrbit(px, py+1.5, pz, camAngle, 0.2, 10)
  ClearBackground(SkyBlue)
  DrawCube(0, -0.5, 0, 25, 1, 25, DarkGreen)
  DrawSphere(px, py, pz, 0.5, Red)
UNTIL WindowShouldClose()
CloseWindow()
```

**Hybrid loop:** Define **`update(dt)`** and **`draw()`** (Sub or Function) and use an empty game loop body; the compiler injects physics step, update, clear, draw, and flush automatically. See [docs/PROGRAM_STRUCTURE.md](docs/PROGRAM_STRUCTURE.md).

**Movement:** Use **GetAxisX()** / **GetAxisY()** for -1/0/1, or **GAME.MoveWASD** / **MoveHorizontal2D** for full 2D/3D. **Delta time:** **DeltaTime()** or **GetFrameTime()**; clamp with `IF dt > 0.05 THEN LET dt = 0.016` for physics.  
**Center screen:** **GetScreenCenterX()**, **GetScreenCenterY()** – center UI or spawn. **Distance:** **Distance2D(x1, y1, x2, y2)** and **Distance3D(x1, y1, z1, x2, y2, z2)** for simple distance.  
**Multi-window (same .bas):** **IsWindowProcess()**, **SpawnWindow(port, title, w, h)**, **ConnectToParent()**; main uses **AcceptTimeout** to get connection, then **Send**/ **Receive**. See [docs/MULTI_WINDOW.md](docs/MULTI_WINDOW.md).  
**Command reference:** [Core Command Reference](docs/CORE_COMMAND_REFERENCE.md), [2D Game API](docs/2D_GAME_API.md), [3D Game API](docs/3D_GAME_API.md). See [API_REFERENCE.md](API_REFERENCE.md) and [examples/README.md](examples/README.md).
