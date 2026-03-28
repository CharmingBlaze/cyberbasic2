// Smoke: blend/scissor + IsFileDropped / GetDroppedFilePaths parity tranche (Phase C).
// cyberbasic --lint examples/smoke_raylib_tranche_c.bas
InitWindow(420, 280, "Raylib tranche C")
SetTargetFPS(60)
VAR bm = BLEND_ALPHA()

mainloop
  ClearBackground(35, 38, 48, 255)
  BeginBlendMode(bm)
  DrawRectangle(30, 30, 120, 90, 180, 90, 90, 100)
  EndBlendMode()
  BeginScissorMode(40, 40, 220, 120)
  DrawText("scissor + blend", 48, 72, 18, 240, 240, 250, 255)
  EndScissorMode()
  IF IsFileDropped() THEN
    DrawText(GetDroppedFilePaths(), 8, 8, 12, 255, 220, 80, 255)
  END IF
  SYNC
endmain

CloseWindow()
