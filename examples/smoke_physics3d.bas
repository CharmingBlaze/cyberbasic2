// Smoke: Bullet-style 3D world create + step (empty world).
// cyberbasic --lint examples/smoke_physics3d.bas
InitWindow(400, 200, "Smoke Physics3D")
SetTargetFPS(60)
CreateWorld3D("smoke3d", 0, -9.8, 0)

mainloop
  Step3D("smoke3d", DeltaTime())
  ClearBackground(35, 40, 55, 255)
  DrawText("Step3D ok", 10, 10, 18, 255, 255, 255, 255)
  SYNC
endmain

DestroyWorld3D("smoke3d")
CloseWindow()
