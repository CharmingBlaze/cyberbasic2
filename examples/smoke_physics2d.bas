// Smoke: Box2D world create + step (no bodies).
// cyberbasic --lint examples/smoke_physics2d.bas
InitWindow(400, 240, "Smoke Physics2D")
SetTargetFPS(60)
CreateWorld2D("smoke", 0, -10)

mainloop
  Step2D("smoke", DeltaTime())
  ClearBackground(40, 45, 55, 255)
  DrawText("Box2D stepped", 10, 10, 18, 255, 255, 255, 255)
  SYNC
endmain

DestroyWorld2D("smoke")
CloseWindow()
