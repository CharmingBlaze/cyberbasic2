// Smoke: DrawRectangleGradientH / V (2D fills).
// cyberbasic --lint examples/smoke_rectangle_gradient.bas
InitWindow(480, 160, "Smoke gradients")
SetTargetFPS(60)

mainloop
  ClearBackground(15, 15, 20, 255)
  DrawRectangleGradientH(0, 0, 240, 160, 80, 0, 120, 255, 200, 60, 40, 255)
  DrawRectangleGradientV(240, 0, 240, 160, 30, 30, 80, 255, 180, 220, 255, 255)
  DrawText("H | V gradients", 12, 12, 18, 255, 255, 255, 255)
  SYNC
endmain

CloseWindow()
