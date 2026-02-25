// Multiplayer server: Host, AcceptTimeout, JoinRoom, SendToRoom
// Run first: cyberbasic examples/multiplayer_server.bas
// Then run client: cyberbasic examples/multiplayer_client.bas

VAR sid = Host(9999)
IF IsNull(sid) THEN
  PRINT "Failed to host on port 9999"
  QUIT
ENDIF
CreateRoom("lobby")

InitWindow(500, 300, "Multiplayer Server")
SetTargetFPS(60)
VAR lastMsg = "Waiting for clients..."

WHILE NOT WindowShouldClose()
  VAR cid = AcceptTimeout(sid, 50)
  IF NOT IsNull(cid) THEN
    JoinRoom("lobby", cid)
    lastMsg = "New client joined"
  ENDIF
  VAR n = GetConnectionCount()
  SendToRoom("lobby", "clients:" + STR(n))

  ClearBackground(30, 40, 50, 255)
  DrawText("Server - Port 9999", 20, 20, 22, 255, 255, 255, 255)
  DrawText("Connections: " + STR(n), 20, 60, 18, 200, 220, 255, 255)
  DrawText(lastMsg, 20, 100, 16, 180, 180, 180, 255)
  DrawText("Close window to stop server", 20, 260, 14, 150, 150, 150, 255)
WEND

CloseServer(sid)
CloseWindow()
