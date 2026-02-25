// GUI demo: BeginUI, EndUI, Label, Button, Slider, Checkbox (pure-Go layout UI)
// Run: cyberbasic examples/gui_demo.bas

InitWindow(400, 320, "GUI Demo")
SetTargetFPS(60)

VAR volume = 0.5
VAR enabled = 1
VAR clickCount = 0

WHILE NOT WindowShouldClose()
  ClearBackground(40, 40, 50, 255)
  BeginUI()
  Label("GUI Demo - BeginUI / EndUI")
  VAR clicked = Button("Click me")
  IF clicked THEN
    LET clickCount = clickCount + 1
  ENDIF
  Label("Clicks: " + STR(clickCount))
  volume = Slider("Volume", volume, 0, 1)
  enabled = Checkbox("Enable sound", enabled)
  EndUI()
  DrawText("Close window to exit", 10, 290, 14, 180, 180, 180, 255)
WEND

CloseWindow()
