// Smoke: raylib 2D geometric collision (AABB overlap).
// cyberbasic --lint examples/smoke_2d_collision.bas
InitWindow(480, 200, "Smoke 2D collision")
SetTargetFPS(60)

mainloop
  ClearBackground(40, 40, 50, 255)
  VAR hit = CheckCollisionRecs(10, 10, 80, 60, 50, 30, 80, 60)
  IF hit THEN
    DrawRectangle(10, 10, 80, 60, 50, 200, 50, 255)
    DrawText("overlap", 10, 100, 18, 255, 255, 255, 255)
  ELSE
    DrawRectangle(10, 10, 80, 60, 200, 50, 50, 255)
    DrawText("no overlap", 10, 100, 18, 255, 255, 255, 255)
  END IF
  DrawRectangle(50, 30, 80, 60, 80, 80, 200, 180)
  SYNC
endmain

CloseWindow()
