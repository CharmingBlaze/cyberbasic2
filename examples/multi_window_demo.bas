// Multi-window demo: main window + one child window, talking via NET (same .bas)
// Run: cyberbasic examples/multi_window_demo.bas

IF IsWindowProcess() THEN
  VAR cid = ConnectToParent()
  IF IsNull(cid) THEN QUIT
  ENDIF
  InitWindow(GetWindowWidth(), GetWindowHeight(), GetWindowTitle())
  SetTargetFPS(60)
  VAR tickCount = 0
  WHILE NOT WindowShouldClose()
    VAR msg = Receive(cid)
    IF NOT IsNull(msg) THEN tickCount = tickCount + 1
    ENDIF
    Send(cid, "ack")
    ClearBackground(40, 40, 50, 255)
    DrawText("Child - ticks: " + Str(tickCount), 20, 30, 20, 255, 255, 255, 255)
  WEND
  CloseWindow()
  QUIT
ENDIF

VAR sid = Host(9999)
IF IsNull(sid) THEN QUIT
ENDIF
InitWindow(800, 600, "Main window")
SetTargetFPS(60)
SpawnWindow(9999, "Child", 400, 300)
VAR cid = AcceptTimeout(sid, 5000)
IF IsNull(cid) THEN
  CloseServer(sid)
  QUIT
ENDIF
VAR frameCount = 0
WHILE NOT WindowShouldClose()
  Send(cid, "tick")
  VAR msg = Receive(cid)
  frameCount = frameCount + 1
  ClearBackground(25, 25, 35, 255)
  DrawText("Main - frame " + Str(frameCount), 20, 30, 20, 200, 220, 255, 255)
  IF NOT IsNull(msg) THEN DrawText("Reply: " + msg, 20, 55, 18, 180, 200, 255, 255)
  ENDIF
WEND
CloseServer(sid)
CloseWindow()
