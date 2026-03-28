// Smoke: unified renderer post-FX path (offscreen 3D+2D blit when camera.fx has entries).
// Lint: cyberbasic --lint examples/smoke_postfx_unified.bas
InitWindow(640, 480, "Post-FX unified")
SetTargetFPS(60)
SetCamera3D(0, 8, 12, 0, 0, 0, 0, 1, 0)
camera.fx.add("vignette")

mainloop
  ClearBackground(30, 30, 40, 255)
  BeginMode3D()
  DrawCube(0, 0, 0, 2, 2, 2, 200, 100, 255, 255)
  DrawGrid(10, 1.0)
  EndMode3D()
  DrawText("smoke_postfx_unified: camera.fx vignette", 10, 10, 18, 255, 255, 255, 255)
  SYNC
endmain

CloseWindow()
