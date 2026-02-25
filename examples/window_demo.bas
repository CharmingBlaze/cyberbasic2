// Window demo - current API: InitWindow, game loop, drawing
// Run: cyberbasic examples/window_demo.bas

InitWindow(800, 450, "CyberBasic Window")
SetTargetFPS(60)

WHILE NOT WindowShouldClose()
  ClearBackground(60, 70, 90, 255)
  DrawRectangle(350, 200, 100, 80, 100, 150, 220, 255)
  DrawText("Window Demo", 320, 170, 24, 255, 255, 255, 255)
  DrawText("Close window to exit", 320, 300, 16, 200, 200, 200, 255)
WEND

CloseWindow()
