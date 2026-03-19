// UI demo: BeginUI, Label, Button.
// Run: cyberbasic examples/ui_demo.bas

InitWindow(600, 400, "UI Demo")
SetTargetFPS(60)

VAR clicks = 0

mainloop
  ClearBackground(40, 40, 50, 255)
  BeginUI(10, 10, 300)
  Label("CyberBASIC2 UI")
  IF Button("Click me") THEN
    LET clicks = clicks + 1
  END IF
  Label("Clicks: " + Str(clicks))
  EndUI()
  DrawText("Immediate-mode UI: Label, Button", 10, 350, 16, 200, 200, 200, 255)
  SYNC
endmain

CloseWindow()
