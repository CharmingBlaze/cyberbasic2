// Smoke: rcamera-style helpers on default SetCamera3D camera.
// cyberbasic --lint examples/smoke_rcamera.bas
InitWindow(480, 300, "Smoke rcamera")
SetTargetFPS(60)
SetCamera3D(0, 6, 10, 0, 0, 0, 0, 1, 0)

mainloop
  ClearBackground(20, 22, 30, 255)
  CameraMoveForward(GetFrameTime() * 0.15, 1)
  BeginMode3D()
  DrawCube(0, 0, 0, 1.5, 1.5, 1.5, 100, 180, 255, 255)
  DrawGrid(12, 1.0)
  EndMode3D()
  DrawText("CameraMoveForward (world plane)", 8, 8, 16, 255, 255, 255, 255)
  SYNC
endmain

CloseWindow()
