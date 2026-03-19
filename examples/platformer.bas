// Simple 2D platformer: WASD move, Space jump, gravity.
// Run: cyberbasic examples/platformer.bas

InitWindow(800, 600, "Platformer")
SetTargetFPS(60)

VAR playerX = 400
VAR playerY = 300
VAR vx = 0.0
VAR vy = 0.0
VAR groundY = 500
VAR moveSpeed = 200
VAR jumpVel = -320
VAR gravity = 600

mainloop
  VAR dt = DeltaTime()
  LET vx = moveSpeed * GetAxisX()
  IF IsKeyPressed(KEY_SPACE()) AND playerY >= groundY - 40 THEN
    LET vy = jumpVel
  END IF
  LET vy = vy + gravity * dt
  LET playerX = playerX + vx * dt
  LET playerY = playerY + vy * dt
  IF playerY > groundY - 40 THEN
    LET playerY = groundY - 40
    LET vy = 0
  END IF
  IF playerX < 20 THEN
    LET playerX = 20
  END IF
  IF playerX > 780 THEN
    LET playerX = 780
  END IF

  ClearBackground(30, 30, 50, 255)
  DrawRectangle(playerX - 20, playerY - 40, 40, 40, 100, 200, 255, 255)
  DrawRectangle(0, groundY, 800, 100, 60, 80, 60, 255)
  DrawText("WASD move | Space jump", 10, 10, 18, 255, 255, 255, 255)
  SYNC
endmain

CloseWindow()
