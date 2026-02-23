// First game: window, WASD input, draw loop
// InitWindow, SetTargetFPS, WHILE NOT WindowShouldClose(), IsKeyDown(KEY_W), BeginDrawing, ClearBackground, DrawCircle, EndDrawing, CloseWindow

InitWindow(800, 600, "My Game")
SetTargetFPS(60)

VAR x = 400
VAR y = 300

WHILE NOT WindowShouldClose()
  IF IsKeyDown(KEY_W) THEN
    LET y = y - 4
  ENDIF
  IF IsKeyDown(KEY_S) THEN
    LET y = y + 4
  ENDIF
  IF IsKeyDown(KEY_A) THEN
    LET x = x - 4
  ENDIF
  IF IsKeyDown(KEY_D) THEN
    LET x = x + 4
  ENDIF

  ClearBackground(20, 20, 30, 255)
  DrawCircle(x, y, 30, 255, 100, 100, 255)
  DrawText("WASD to move", 10, 10, 20, 255, 255, 255, 255)
WEND

CloseWindow()
