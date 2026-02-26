# CyberBasic – First 10 lines

Minimal snippets to get a game running.

## 2D game (move a circle)

Use **WHILE NOT WindowShouldClose() ... WEND** (or **REPEAT ... UNTIL WindowShouldClose()**) for the main loop. No auto-wrap; code compiles as written (DBPro-style). **DeltaTime()** for frame delta.

```basic
InitWindow(800, 600, "My Game")
SetTargetFPS(60)
VAR x = 400
VAR y = 300
WHILE NOT WindowShouldClose()
  VAR dt = DeltaTime()
  LET x = x + 100 * dt * GetAxisX()
  LET y = y + 100 * dt * GetAxisY()
  ClearBackground(20, 20, 30, 255)
  DrawCircle(x, y, 30, 255, 100, 100, 255)
WEND
CloseWindow()
```

## 3D game (Bullet + orbit camera)

```basic
InitWindow(1024, 600, "3D")
SetTargetFPS(60)
DisableCursor()
BULLET.CreateWorld("w", 0, -18, 0)
BULLET.CreateSphere("w", "player", 0, 0.5, 0, 0.5, 1)
BULLET.CreateBox("w", "ground", 0, -0.5, 0, 12.5, 0.5, 12.5, 0)
VAR camAngle = 0
REPEAT
  VAR dt = GetFrameTime()
  BULLET.Step("w", dt)
  GAME.MoveWASD("w", "player", camAngle, 120, 9, dt)
  VAR px = BULLET.GetPositionX("w", "player")
  VAR py = BULLET.GetPositionY("w", "player")
  VAR pz = BULLET.GetPositionZ("w", "player")
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
Full 2D/3D command lists: [2D Graphics Guide](docs/2D_GRAPHICS_GUIDE.md#full-2d-command-reference), [3D Graphics Guide](docs/3D_GRAPHICS_GUIDE.md#full-3d-command-reference). See [API_REFERENCE.md](API_REFERENCE.md) and [examples/README.md](examples/README.md).
