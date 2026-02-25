// 2D Shapes Demo in CyberBasic - current API
// Demonstrates basic 2D graphics and shapes

InitWindow(800, 600, "2D Shapes Demo")
SetTargetFPS(60)

VAR time = 0.0

WHILE NOT WindowShouldClose()
  LET time = time + GetFrameTime()
  ClearBackground(135, 206, 235, 255)

  // Red rectangle
  DrawRectangle(100, 100, 200, 150, 255, 0, 0, 255)
  // Green circle
  DrawCircle(400, 200, 50, 0, 255, 0, 255)
  // Blue square (rectangle)
  DrawRectangle(550, 150, 100, 100, 0, 0, 255, 255)
  // Yellow line (thin rectangle)
  DrawRectangle(100, 300, 200, 5, 255, 255, 0, 255)
  // Purple text
  DrawText("2D Shapes Demo!", 300, 50, 24, 255, 0, 255, 255)
  // Animated circle
  VAR x = 400 + SIN(time * 2) * 100
  DrawCircle(x, 400, 30, 255, 165, 0, 255)
  // Grid pattern
  VAR i = 0
  FOR i = 0 TO 8
    DrawRectangle(i * 100, 500, 1, 100, 128, 128, 128, 255)
  NEXT i
  FOR i = 0 TO 10
    DrawRectangle(0, 500 + i * 10, 800, 1, 128, 128, 128, 255)
  NEXT i
  // FPS
  VAR fps = GetFPS()
  DrawText("FPS: " + STR(fps), 10, 10, 16, 255, 255, 255, 255)
WEND

CloseWindow()
