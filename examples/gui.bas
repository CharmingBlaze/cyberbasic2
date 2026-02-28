// Simple GUI demo: labels, button, slider
// Run: cyberbasic examples/gui.bas

InitWindow(480, 320, "CyberBasic GUI")
SetTargetFPS(60)

VAR clicks = 0
VAR volume = 50.0

WHILE NOT WindowShouldClose()
  BeginDrawing()
  ClearBackground(50, 55, 65, 255)

  GuiLabel(24, 24, 200, 24, "CyberBasic GUI Demo")
  GuiLabel(24, 56, 200, 22, "Click the button or move the slider:")

  VAR btn = GuiButton(24, 90, 140, 32, "Click Me")
  IF btn THEN
    LET clicks = clicks + 1
  ENDIF
  GuiLabel(24, 130, 220, 22, "Clicks: " + STR(clicks))

  LET volume = GuiSlider(24, 170, 200, 0, 100, volume)
  GuiLabel(240, 166, 80, 22, STR(INT(volume)) + "%")

  GuiLabel(24, 260, 400, 22, "Close the window to exit.")
  EndDrawing()
WEND

CloseWindow()
