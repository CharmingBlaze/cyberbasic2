// Smoke: engine.<subsystem> reaches the same DotObject as the global (e.g. tween.count).
// cyberbasic --lint examples/smoke_engine.bas
InitWindow(320, 160, "Smoke engine")
SetTargetFPS(60)
VAR n = engine.tween.count()

mainloop
  ClearBackground(25, 30, 40, 255)
  DrawText("engine.tween ok", 10, 10, 14, 255, 255, 255, 255)
  SYNC
endmain

CloseWindow()
