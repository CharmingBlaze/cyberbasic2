// Smoke: audio device lifecycle (no sound file required).
// cyberbasic --lint examples/smoke_audio.bas
InitWindow(320, 200, "Smoke Audio")
SetTargetFPS(60)
InitAudioDevice()

mainloop
  ClearBackground(20, 20, 20, 255)
  DrawText("Audio device ready", 10, 10, 16, 200, 255, 200, 255)
  SYNC
endmain

CloseAudioDevice()
CloseWindow()
