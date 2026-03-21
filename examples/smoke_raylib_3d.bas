// Smoke: pure raylib 3D path (BeginMode3D, DrawCube, DrawGrid).
// cyberbasic --lint examples/smoke_raylib_3d.bas
InitWindow(640, 480, "Smoke Raylib 3D")
SetTargetFPS(60)
SetCamera3D(0, 8, 12, 0, 0, 0, 0, 1, 0)

mainloop
  ClearBackground(30, 30, 40, 255)
  BeginMode3D()
  DrawCube(0, 0, 0, 2, 2, 2, 200, 100, 255, 255)
  DrawGrid(10, 1.0)
  EndMode3D()
  DrawText("smoke_raylib_3d", 10, 10, 18, 255, 255, 255, 255)
  SYNC
endmain

CloseWindow()
