// Smoke: tween.register on window.targetfps (DotObject + float prop).
// cyberbasic --lint examples/smoke_tween.bas
InitWindow(360, 200, "Smoke Tween")
SetTargetFPS(60)
TweenRegister(window, "targetfps", 60, 72, 0.5)

mainloop
  ClearBackground(30, 35, 45, 255)
  DrawText("tweening targetfps", 10, 10, 16, 255, 255, 255, 255)
  SYNC
endmain

CloseWindow()
