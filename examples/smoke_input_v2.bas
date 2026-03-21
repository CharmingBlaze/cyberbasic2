// Smoke: v2 input action map + key query.
// cyberbasic --lint examples/smoke_input_v2.bas
InitWindow(400, 200, "Smoke Input v2")
SetTargetFPS(60)
INPUT.MAP.REGISTER("quit", KEY_ESCAPE())

mainloop
  ClearBackground(25, 30, 40, 255)
  IF INPUT.MAP.PRESSED("quit") THEN
    DrawText("quit pressed", 10, 10, 16, 255, 100, 100, 255)
  ELSE
    DrawText("ESC = quit action", 10, 10, 16, 255, 255, 255, 255)
  END IF
  SYNC
endmain

CloseWindow()
