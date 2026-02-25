// Multiplayer client: Connect, Send, Receive
// Run server first: cyberbasic examples/multiplayer_server.bas
// Then run: cyberbasic examples/multiplayer_client.bas

VAR cid = Connect("127.0.0.1", 9999)
IF IsNull(cid) THEN
  PRINT "Failed to connect to 127.0.0.1:9999 - is the server running?"
  QUIT
ENDIF

InitWindow(500, 280, "Multiplayer Client")
SetTargetFPS(60)
VAR lastMsg = "(no message yet)"

WHILE NOT WindowShouldClose()
  VAR msg = Receive(cid)
  IF NOT IsNull(msg) THEN
    lastMsg = msg
  ENDIF
  IF IsKeyPressed(KEY_SPACE) THEN
    Send(cid, "hello from client")
  ENDIF

  ClearBackground(40, 35, 50, 255)
  DrawText("Client - Connected", 20, 20, 22, 255, 255, 255, 255)
  DrawText("Last message: " + lastMsg, 20, 60, 16, 220, 200, 255, 255)
  DrawText("Press SPACE to send 'hello from client'", 20, 100, 14, 180, 180, 180, 255)
  DrawText("Close window to disconnect", 20, 240, 14, 150, 150, 150, 255)
WEND

Disconnect(cid)
CloseWindow()
