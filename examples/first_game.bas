// First game: 3D spinning cube with mouse orbit camera.
// Run: cyberbasic examples/first_game.bas
// Uses DeltaTime for frame-rate independent motion (clamped for stability).
InitWindow(800, 600, "My Game")
SetTargetFPS(60)

// Create cube (DBP: id, size) and position at origin
LoadCube(1, 2)
PositionObject(1, 0, 0, 0)

// Set initial camera before loop (BeginMode3D uses it on first frame)
SetCamera3D(0, 5, 14, 0, 0, 0, 0, 1, 0)

mainloop
  OrbitCamera(0, 0, 0)
  ClearBackground(25, 25, 35, 255)
  YRotateObject(1, GetFrameTime() * 20.5)
  DrawObject(1)
  DrawGrid(20, 1.0)
  DrawText("Drag: orbit | Middle+drag: zoom | Wheel/PgUp/PgDn | Space: test", 10, 10, 18, 255, 255, 255, 255)
  IF IsKeyPressed(KEY_SPACE()) THEN
    DrawText("SPACE tap", 10, 35, 20, 255, 200, 100, 255)
  END IF
  IF IsKeyDown(KEY_SPACE()) THEN
    DrawText("SPACE held", 10, 58, 20, 200, 255, 100, 255)
  END IF
  SYNC
endmain

DeleteObject(1)
CloseWindow()
