// Multi-window GUI demo: logical windows, messages, channels, state, OnWindowDraw
// Run: cyberbasic examples/multi_window_gui_demo.bas

InitWindow(800, 600, "Multi-Window Demo")
SetTargetFPS(60)

VAR win1 = WindowCreate(300, 200, "Panel A")
VAR win2 = WindowCreate(280, 180, "Panel B")
WindowSetPosition(win1, 50, 80)
WindowSetPosition(win2, 420, 100)

ChannelCreate("events")
StateSet("clicks", 0)
OnWindowDraw(win1, "DrawPanelA")
OnWindowDraw(win2, "DrawPanelB")

SUB DrawPanelA(id)
  WindowClearBackground(id, 60, 80, 120, 255)
  DrawText("Panel A - Window " + Str(id), 10, 10, 20, 255, 255, 255, 255)
  DrawText("Send message to Panel B", 10, 45, 16, 200, 220, 255, 255)
  VAR cnt = StateGet("clicks")
  IF IsNull(cnt) THEN LET cnt = 0
  ENDIF
  DrawText("Clicks: " + Str(cnt), 10, 75, 18, 255, 255, 200, 255)
END SUB

SUB DrawPanelB(id)
  WindowClearBackground(id, 80, 60, 100, 255)
  DrawText("Panel B - Window " + Str(id), 10, 10, 20, 255, 255, 255, 255)
  IF WindowHasMessage(id) THEN
    VAR received = WindowReceiveMessage(id)
    DrawText("Got: " + Str(received), 10, 45, 16, 200, 255, 200, 255)
  ELSE
    DrawText("(no message)", 10, 45, 16, 180, 180, 180, 255)
  ENDIF
  DrawText("Uses Channel + State", 10, 80, 14, 220, 220, 255, 255)
END SUB

WHILE NOT WindowShouldClose()
  WindowProcessEvents()
  WindowBeginDrawing(win1)
  DrawPanelA(win1)
  WindowEndDrawing(win1)
  WindowBeginDrawing(win2)
  DrawPanelB(win2)
  WindowEndDrawing(win2)

  BeginDrawing()
  ClearBackground(25, 30, 40, 255)
  DrawText("Multi-Window Demo - Two panels below", 20, 20, 22, 255, 255, 255, 255)
  WindowDrawAllToScreen()
  EndDrawing()
WEND

WindowClose(win1)
WindowClose(win2)
CloseWindow()
