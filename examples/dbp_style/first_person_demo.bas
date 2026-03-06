REM DBP-style first-person camera demo
REM WASD move, mouse look, procedural cube in scene

InitWindow(800, 600, "DBP First-Person Demo")
SetTargetFPS(60)
DisableCursor()

DIM camX AS FLOAT
DIM camY AS FLOAT
DIM camZ AS FLOAT
DIM yaw AS FLOAT
DIM pitch AS FLOAT

camX = 0
camY = 2
camZ = 10
yaw = 0
pitch = 0

LoadCube(1, 2)
PositionObject(1, 0, 0, 0)

WHILE NOT WindowShouldClose()
  VAR dt = GetFrameTime()
  VAR moveSpeed = 5.0 * dt
  VAR lookSpeed = 0.002

  yaw = yaw - GetMouseDeltaX() * lookSpeed
  pitch = pitch - GetMouseDeltaY() * lookSpeed
  IF pitch > 1.5 THEN pitch = 1.5
  IF pitch < -1.5 THEN pitch = -1.5

  VAR cosYaw = Cos(yaw)
  VAR sinYaw = Sin(yaw)
  IF IsKeyDown(KEY_W) THEN
    camX = camX - sinYaw * moveSpeed
    camZ = camZ - cosYaw * moveSpeed
  END IF
  IF IsKeyDown(KEY_S) THEN
    camX = camX + sinYaw * moveSpeed
    camZ = camZ + cosYaw * moveSpeed
  END IF
  IF IsKeyDown(KEY_A) THEN
    camX = camX - cosYaw * moveSpeed
    camZ = camZ + sinYaw * moveSpeed
  END IF
  IF IsKeyDown(KEY_D) THEN
    camX = camX + cosYaw * moveSpeed
    camZ = camZ - sinYaw * moveSpeed
  END IF

  VAR targetX = camX - Sin(yaw) * Cos(pitch)
  VAR targetY = camY + Sin(pitch)
  VAR targetZ = camZ - Cos(yaw) * Cos(pitch)

  SetCameraPosition(camX, camY, camZ)
  SetCameraTarget(targetX, targetY, targetZ)

  BeginDrawing()
  ClearBackground(30, 30, 50, 255)
  BeginMode3D()
  DrawObject(1)
  DrawGrid(20, 1)
  EndMode3D()
  DrawText("WASD move, mouse look - ESC quit", 10, 10, 20, 255, 255, 255, 255)
  EndDrawing()
WEND

EnableCursor()
CloseWindow()
