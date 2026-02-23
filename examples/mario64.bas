REM 3D Mario 64-style: orbit camera (mouse), WASD move, CONST + GAME helpers

CONST KEY_W = 87
CONST KEY_A = 65
CONST KEY_S = 83
CONST KEY_D = 68
CONST KEY_SPACE = 32

RL.InitWindow(1024, 600, "Mario 64 Style")
RL.SetTargetFPS(60)
RL.DisableCursor()

REM Physics: ground top at y=0, player sphere radius 0.5 (center 0.5 on ground)
BULLET.CreateWorld("w", 0, -18, 0)
BULLET.CreateSphere("w", "player", 0, 0.5, 0, 0.5, 1)
BULLET.CreateBox("w", "ground", 0, -0.5, 0, 12.5, 0.5, 12.5, 0)
BULLET.CreateBox("w", "plat1", 4, 0.5, 0, 2, 1, 2, 0)
BULLET.CreateBox("w", "plat2", -3, 0.5, 2, 1.5, 1, 1.5, 0)
BULLET.CreateBox("w", "plat3", 0, 0.5, -4, 3, 1, 2, 0)

DIM camAngle AS Float
LET camAngle = 0
DIM camDist AS Float
LET camDist = 10
DIM mouseSens AS Float
LET mouseSens = 0.002
DIM moveForce AS Float
LET moveForce = 120
DIM jumpVel AS Float
LET jumpVel = 9
DIM playerRad AS Float
LET playerRad = 0.5

REPEAT
  LET dt = RL.GetFrameTime()
  IF dt > 0.05 THEN
    LET dt = 0.016
  ENDIF

  REM Mouse: orbit camera
  LET delta = RL.GetMouseDelta()
  LET dx = delta.x
  LET camAngle = camAngle - dx * mouseSens

  BULLET.Step("w", dt)

  LET px = BULLET.GetPositionX("w", "player")
  LET py = BULLET.GetPositionY("w", "player")
  LET pz = BULLET.GetPositionZ("w", "player")

  REM Snap to platform tops when over them
  IF px > 3 AND px < 5 AND pz > -1 AND pz < 1 THEN
    GAME.SnapToGround("w", "player", 1.5, playerRad)
  ENDIF
  IF px > -3.75 AND px < -2.25 AND pz > 1.25 AND pz < 2.75 THEN
    GAME.SnapToGround("w", "player", 1.5, playerRad)
  ENDIF
  IF px > -1.5 AND px < 1.5 AND pz > -5 AND pz < -3 THEN
    GAME.SnapToGround("w", "player", 1.5, playerRad)
  ENDIF

  REM WASD + jump (ground at y=0)
  GAME.MoveWASD("w", "player", camAngle, moveForce, jumpVel, dt)

  REM Orbit camera around player (re-read position after snaps)
  LET px = BULLET.GetPositionX("w", "player")
  LET py = BULLET.GetPositionY("w", "player")
  LET pz = BULLET.GetPositionZ("w", "player")
  GAME.CameraOrbit(px, py + 1.5, pz, camAngle, 0.2, camDist)

  RL.BeginDrawing()
  RL.ClearBackground(RL.SkyBlue)
  RL.BeginMode3D()
  RL.DrawCube(0, -0.5, 0, 25, 1, 25, RL.DarkGreen)
  RL.DrawCube(4, 0.5, 0, 2, 1, 2, RL.Gray)
  RL.DrawCube(-3, 0.5, 2, 1.5, 1, 1.5, RL.Brown)
  RL.DrawCube(0, 0.5, -4, 3, 1, 2, RL.Gold)
  LET px2 = BULLET.GetPositionX("w", "player")
  LET py2 = BULLET.GetPositionY("w", "player")
  LET pz2 = BULLET.GetPositionZ("w", "player")
  RL.DrawSphere(px2, py2, pz2, 0.5, RL.Red)
  RL.EndMode3D()
  RL.DrawText("WASD move, Mouse look, Space jump", 10, 10, 20, RL.White)
  RL.EndDrawing()
UNTIL RL.WindowShouldClose()

RL.EnableCursor()
RL.CloseWindow()
PRINT "Bye!"
