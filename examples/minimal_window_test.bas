PRINT "Starting..."
InitWindow(800, 600, "Test Window")
PRINT "InitWindow done."
SetTargetFPS(60)
REPEAT
  BeginDrawing()
  ClearBackground(100, 100, 150, 255)
  DrawText("If you see this, window works!", 50, 50, 20, 255, 255, 255, 255)
  EndDrawing()
UNTIL WindowShouldClose()
CloseWindow()
PRINT "Done."
