REM DBP-style 3D: procedural cube, PositionObject, YRotateObject, SYNC
REM Uses implicit window from runtime - we need InitWindow for manual loop

InitWindow(800, 600, "DBP-Style 3D Cube")
SetTargetFPS(60)

REM Create cube (id=1, size=2) - no .obj file needed
LoadCube(1, 2)
PositionObject(1, 0, 0, 5)

REM Set camera
SetCameraPosition(0, 2, 10)
SetCameraTarget(0, 0, 0)

WHILE NOT WindowShouldClose()
  YRotateObject(1, 1)
  BeginDrawing()
  ClearBackground(30, 30, 40, 255)
  BeginMode3D()
  DrawObject(1)
  EndMode3D()
  DrawText("DBP-style 3D cube - ESC to quit", 10, 10, 20, 255, 255, 255, 255)
  EndDrawing()
WEND

CloseWindow()
