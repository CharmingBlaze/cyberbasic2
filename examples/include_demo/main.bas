// Include demo: #include "lib/helper.bas" and call Sub from included file
// Run from project root: cyberbasic examples/include_demo/main.bas

#include "lib/helper.bas"

InitWindow(500, 280, "Include Demo")
SetTargetFPS(60)

WHILE NOT WindowShouldClose()
  ClearBackground(35, 40, 55, 255)
  DrawText("Main program - message from included lib:", 20, 20, 18, 255, 255, 255, 255)
  DrawHelperText(20, 60, 20, 200, 220, 255, 255)
  DrawText("Close window to exit", 20, 250, 14, 180, 180, 180, 255)
WEND

CloseWindow()
