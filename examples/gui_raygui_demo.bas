// Comprehensive raygui demo: Update() / Draw() refactor, single clear per frame to reduce flicker
// Run: cyberbasic examples/gui_raygui_demo.bas

SetConfigFlags(FLAG_VSYNC_HINT())
InitWindow(920, 640, "RayGui Full Demo")
SetTargetFPS(60)

VAR clickCount = 0
VAR checked1 = 0
VAR checked2 = 0
VAR sliderVal = 50.0
VAR progressVal = 0.35
VAR dropdownSel = 0
VAR dropdownBoxSel = 0
VAR listSel = 0
VAR windowClosed = 0
VAR tbContent = "Type here..."
VAR tbContent2 = "Named textbox"

SUB Update()
  // Per-frame logic that doesn't draw (e.g. physics, input-only). GUI state is updated in Draw() when Gui* return values.
END SUB

SUB Draw()
  BeginDrawing()
  ClearBackground(45, 52, 64, 255)

  // --- Main window box (title bar) ---
  VAR closeBtn = GuiWindowBox(10, 10, 440, 300, "Controls Panel #1")
  IF closeBtn THEN
    LET windowClosed = 1
  ENDIF

  GuiLabel(24, 40, 200, 22, "Label: Hello from RayGui")
  GuiLabel(24, 68, 180, 22, "Another label")

  VAR btn = GuiButton(24, 96, 120, 28, "Click Me")
  IF btn THEN
    LET clickCount = clickCount + 1
  ENDIF
  VAR resetBtn = GuiButton(156, 96, 100, 28, "Reset")
  IF resetBtn THEN
    LET clickCount = 0
  ENDIF
  GuiLabel(24, 130, 250, 22, "Clicks: " + STR(clickCount))

  LET checked1 = GuiCheckBox(24, 158, 24, 24, "Option A", checked1)
  LET checked2 = GuiCheckbox("Option B", 24, 188, checked2)

  LET sliderVal = GuiSlider(24, 222, 180, 0, 100, sliderVal)
  GuiLabel(210, 218, 80, 22, STR(INT(sliderVal)))

  LET progressVal = GuiProgressBar(24, 252, 200, 24, "0%", "100%", progressVal, 0, 1)
  GuiProgressBarSimple(240, 252, 120, progressVal)

  LET tbContent = GuiTextbox(24, 284, 180, tbContent)
  LET tbContent2 = GuiTextbox(240, 96, 180, tbContent2)

  // --- Second panel ---
  GuiPanel(460, 10, 440, 300, "")
  GuiGroupBox(476, 30, 200, 120, "Group: Options")
  GuiLine(476, 38, 200, 2, "")
  GuiCheckbox("Group opt 1", 486, 52, 0)
  GuiCheckbox("Group opt 2", 486, 78, 1)
  GuiButton(486, 108, 80, 24, "Apply")

  GuiLine(476, 158, 200, 2, "Section divider")
  LET dropdownSel = GuiDropdown("Red;Green;Blue;Yellow", 486, 168, 140)
  GuiLabel(476, 200, 200, 22, "Dropdown index: " + STR(dropdownSel))

  LET listSel = GuiList("Apple;Banana;Cherry;Date;Elderberry", 700, 30, 180, 120)
  GuiLabel(700, 158, 180, 22, "List selected: " + STR(listSel))

  LET dropdownBoxSel = GuiDropdownBox("dd1", 476, 230, 160, 28, "One;Two;Three;Four", dropdownBoxSel)

  // --- Bottom ---
  GuiWindowBox(10, 320, 440, 200, "Window Box 2")
  GuiLabel(24, 348, 300, 22, "GuiLabel in second window")
  GuiSlider(24, 378, 250, 0, 10, 5.0)
  GuiProgressBarSimple(24, 408, 250, 0.7)
  GuiButton(284, 378, 100, 28, "OK")

  GuiPanel(460, 320, 440, 200, "Bottom Panel")
  GuiGroupBox(476, 338, 200, 80, "Final group")
  GuiLabel(476, 358, 180, 22, "All raygui controls in one demo")
  GuiLine(476, 388, 200, 2, "")
  GuiButton(476, 398, 120, 28, "Quit (close window)")

  DrawText("RayGui Full Demo - close window to exit", 12, 610, 14, 180, 180, 180, 255)
  EndDrawing()
END SUB

WHILE NOT WindowShouldClose()
  Update()
  Draw()
WEND

CloseWindow()
