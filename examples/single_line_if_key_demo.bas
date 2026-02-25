// Single-line IF + isKeyDown + movePlayer pattern (no ENDIF per line)
InitWindow(400, 300, "Single-line IF demo")
SetTargetFPS(60)

SUB movePlayer(dx, dy)
  REM could update player position
END SUB

VAR speed = 2
WHILE NOT WindowShouldClose()
  IF IsKeyDown(KEY_W) THEN movePlayer(0, -speed)
  IF IsKeyDown(KEY_S) THEN movePlayer(0, speed)
  IF IsKeyDown(KEY_A) THEN movePlayer(-speed, 0)
  IF IsKeyDown(KEY_D) THEN movePlayer(speed, 0)
  ENDIF
  ClearBackground(30, 30, 40, 255)
  DrawText("WASD: single-line IF", 20, 20, 18, 255, 255, 255, 255)
WEND
CloseWindow()
