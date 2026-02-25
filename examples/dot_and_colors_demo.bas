// Dot notation and color constants demo - flat API
// Uses ClearBackground, GetMouseX, DrawText, DrawRectangle

InitWindow(800, 450, "Dot and Colors Demo")
SetTargetFPS(60)

WHILE NOT WindowShouldClose()
  VAR mx = GetMouseX()
  ClearBackground(80, 80, 80, 255)
  DrawText("Mouse X: " + STR(mx), 20, 20, 20, 255, 255, 255, 255)
  DrawRectangle(100, 100, 200, 80, 0, 100, 255, 255)
  DrawRectangle(400, 200, 150, 150, 255, 220, 0, 255)
  DrawText("Close window to exit", 20, 420, 14, 200, 200, 200, 255)
WEND

CloseWindow()
PRINT "Done."
