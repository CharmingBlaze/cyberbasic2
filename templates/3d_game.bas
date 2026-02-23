// Minimal 3D game template – Bullet physics, orbit camera, WASD + jump
// Run: cyberbasic templates/3d_game.bas

RL.InitWindow(1024, 600, "3D Game")
RL.SetTargetFPS(60)
RL.DisableCursor()

REM Physics: ground at y=0, player sphere
BULLET.CreateWorld("w", 0, -18, 0)
BULLET.CreateSphere("w", "player", 0, 0.5, 0, 0.5, 1)
BULLET.CreateBox("w", "ground", 0, -0.5, 0, 12.5, 0.5, 12.5, 0)

VAR camAngle = 0
VAR camDist = 10
VAR dt = 0.016

REPEAT
  LET dt = RL.GetFrameTime()
  IF dt > 0.05 THEN LET dt = 0.016
  LET delta = RL.GetMouseDelta()
  LET camAngle = camAngle - delta.x * 0.002

  BULLET.Step("w", dt)
  GAME.MoveWASD("w", "player", camAngle, 120, 9, dt)

  LET px = BULLET.GetPositionX("w", "player")
  LET py = BULLET.GetPositionY("w", "player")
  LET pz = BULLET.GetPositionZ("w", "player")
  GAME.CameraOrbit(px, py + 1.5, pz, camAngle, 0.2, camDist)

  RL.ClearBackground(RL.SkyBlue)
  RL.BeginMode3D()
  RL.DrawCube(0, -0.5, 0, 25, 1, 25, RL.DarkGreen)
  RL.DrawSphere(px, py, pz, 0.5, RL.Red)
  RL.EndMode3D()
  RL.DrawText("WASD move, Mouse look, Space jump", 10, 10, 20, RL.White)
UNTIL RL.WindowShouldClose()

RL.EnableCursor()
RL.CloseWindow()
