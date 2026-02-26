// Hybrid update/draw demo: automatic lifecycle and command queue
// Define update(dt) and draw(); the main loop is replaced by: physics step, update(dt), draw (queued), flush.
// Run: cyberbasic examples/hybrid_update_draw_demo.bas

InitWindow(800, 600, "Hybrid Update/Draw Demo")
SetTargetFPS(60)

VAR x = 400.0
VAR y = 300.0

SUB update(dt)
  IF IsKeyDown(KEY_W) THEN LET y = y - 200.0 * dt
  ENDIF
  IF IsKeyDown(KEY_S) THEN LET y = y + 200.0 * dt
  ENDIF
  IF IsKeyDown(KEY_A) THEN LET x = x - 200.0 * dt
  ENDIF
  IF IsKeyDown(KEY_D) THEN LET x = x + 200.0 * dt
  ENDIF
END SUB

SUB draw()
  ClearBackground(30, 30, 45, 255)
  DrawRectangle(Int(x) - 20, Int(y) - 20, 40, 40, 255, 100, 100, 255)
  DrawText("WASD to move - Hybrid loop", 20, 20, 20, 255, 255, 255, 255)
END SUB

WHILE NOT WindowShouldClose()
WEND

CloseWindow()
