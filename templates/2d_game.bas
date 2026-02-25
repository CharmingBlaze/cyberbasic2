// Minimal 2D game template – WHILE NOT WindowShouldClose() loop, DeltaTime(), GetAxisX/GetAxisY
// Run: cyberbasic templates/2d_game.bas

InitWindow(800, 600, "2D Game")
SetTargetFPS(60)

VAR x = 400
VAR y = 300
VAR speed = 100

WHILE NOT WindowShouldClose()
  VAR dt = DeltaTime()
  LET x = x + speed * dt * GetAxisX()
  LET y = y + speed * dt * GetAxisY()
  ClearBackground(20, 20, 30, 255)
  DrawCircle(x, y, 30, 255, 100, 100, 255)
  DrawText("WASD to move", 10, 10, 20, 255, 255, 255, 255)
WEND

CloseWindow()
