// 3D orbit camera + click to drop physics balls
// Right-drag: orbit  Wheel: zoom  Left-click: drop ball (cursor visible for easy aiming)
// Run: cyberbasic examples/orbit_drop_balls_3d.bas

InitWindow(800, 600, "Orbit Camera - Click to Drop Balls")
SetTargetFPS(60)

VAR ballCount = 0
VAR maxBalls = 64

// Physics world and ground (box center y=-0.25 so top is at y=0)
CreateWorld3D("world", 0, -9.81, 0)
CreateBox3D("world", "ground", 0, -0.25, 0, 15, 0.5, 15, 0)

SetCamera3D(0, 5, 14, 0, 0, 0, 0, 1, 0)

WHILE NOT WindowShouldClose()
  VAR dt = GetFrameTime()
  IF dt > 0.05 THEN
    LET dt = 0.016
  ENDIF

  OrbitCamera(0, 0, 0)

  VAR hit = PickGroundPlane()
  VAR dropX = GetRayCollisionPointX()
  VAR dropZ = GetRayCollisionPointZ()
  IF IsMouseButtonPressed(0) AND ballCount < maxBalls AND hit THEN
    LET ballCount = ballCount + 1
    CreateSphere3D("world", "ball" + STR(ballCount), dropX, 10, dropZ, 0.4, 1)
  ENDIF

  VAR substeps = 4
  VAR subdt = dt / substeps
  FOR s = 1 TO substeps
    Step3D("world", subdt)
  NEXT

  BeginDrawing()
  ClearBackground(32, 36, 48, 255)

  BeginMode3D()
  DrawCube(0, -0.25, 0, 30, 0.5, 30, 70, 85, 95, 255)
  DrawCubeWires(0, -0.25, 0, 30, 0.5, 30, 45, 52, 60, 255)

  FOR i = 1 TO ballCount
    VAR bx = GetPositionX3D("world", "ball" + STR(i))
    VAR by = GetPositionY3D("world", "ball" + STR(i))
    VAR bz = GetPositionZ3D("world", "ball" + STR(i))
    VAR k = (i - 1) - INT((i - 1) / 8) * 8
    VAR r = 255
    VAR g = 100
    VAR b = 120
    IF k = 1 THEN
      LET r = 255
      LET g = 180
      LET b = 80
    ENDIF
    IF k = 2 THEN
      LET r = 100
      LET g = 220
      LET b = 120
    ENDIF
    IF k = 3 THEN
      LET r = 80
      LET g = 160
      LET b = 255
    ENDIF
    IF k = 4 THEN
      LET r = 220
      LET g = 100
      LET b = 255
    ENDIF
    IF k = 5 THEN
      LET r = 255
      LET g = 220
      LET b = 80
    ENDIF
    IF k = 6 THEN
      LET r = 100
      LET g = 255
      LET b = 200
    ENDIF
    IF k = 7 THEN
      LET r = 255
      LET g = 120
      LET b = 120
    ENDIF
    DrawSphere(bx, by, bz, 0.4, r, g, b, 255)
    DrawSphereWires(bx, by, bz, 0.4, r * 3 / 4, g * 3 / 4, b * 3 / 4, 255)
  NEXT

  DrawGrid(20, 1.0)
  EndMode3D()

  DrawText("Right-drag: orbit  Wheel: zoom  Left-click: drop ball", 10, 10, 18, 255, 255, 255, 255)
  DrawText("Balls: " + STR(ballCount) + " / " + STR(maxBalls), 10, 32, 16, 200, 200, 200, 255)
  EndDrawing()
WEND

CloseWindow()
