# CyberBasic – First 10 lines

Minimal snippets to get a game running.

## 2D game (move a circle)

The loop `WHILE NOT WindowShouldClose() ... WEND` and **Main() ... EndMain** are automatically wrapped with a draw frame; you don't call BeginDrawing/EndDrawing. Prefer **Main() ... EndMain** and **DeltaTime()** for the main game loop.

```basic
InitWindow(800, 600, "My Game")
SetTargetFPS(60)
VAR x = 400
VAR y = 300
Main()
  VAR dt = DeltaTime()
  LET x = x + 100 * dt * GetAxisX()
  LET y = y + 100 * dt * GetAxisY()
  ClearBackground(20, 20, 30, 255)
  DrawCircle(x, y, 30, 255, 100, 100, 255)
EndMain
CloseWindow()
```

## 3D game (Bullet + orbit camera)

```basic
RL.InitWindow(1024, 600, "3D")
RL.SetTargetFPS(60)
RL.DisableCursor()
BULLET.CreateWorld("w", 0, -18, 0)
BULLET.CreateSphere("w", "player", 0, 0.5, 0, 0.5, 1)
BULLET.CreateBox("w", "ground", 0, -0.5, 0, 12.5, 0.5, 12.5, 0)
VAR camAngle = 0
REPEAT
  LET dt = RL.GetFrameTime()
  BULLET.Step("w", dt)
  GAME.MoveWASD("w", "player", camAngle, 120, 9, dt)
  LET px = BULLET.GetPositionX("w", "player")
  LET py = BULLET.GetPositionY("w", "player")
  LET pz = BULLET.GetPositionZ("w", "player")
  GAME.CameraOrbit(px, py+1.5, pz, camAngle, 0.2, 10)
  RL.BeginDrawing()
  RL.ClearBackground(RL.SkyBlue)
  RL.BeginMode3D()
  RL.DrawCube(0, -0.5, 0, 25, 1, 25, RL.DarkGreen)
  RL.DrawSphere(px, py, pz, 0.5, RL.Red)
  RL.EndMode3D()
  RL.EndDrawing()
UNTIL RL.WindowShouldClose()
RL.CloseWindow()
```

**Movement:** Use **GetAxisX()** / **GetAxisY()** for -1/0/1, or **GAME.MoveWASD** / **MoveHorizontal2D** for full 2D/3D. **Delta time:** **DeltaTime()** or **GetFrameTime()**; clamp with `IF dt > 0.05 THEN LET dt = 0.016` for physics. See [API_REFERENCE.md](API_REFERENCE.md) and [examples/README.md](examples/README.md).
