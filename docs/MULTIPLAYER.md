# Multiplayer (TCP)

CyberBasic includes simple TCP client/server bindings so you can send and receive text messages between programs (e.g. for chat, lobby, or simple game state). No extra libraries; uses Go's standard `net` package. For encrypted channels over the internet, use **ConnectTLS** and **HostTLS**.

## Quick pick: choose your multiplayer type

Pick a use case and copy the snippet. Same commands work for all; only the flow changes.

- **[Two players (1 server, 1 client)](#two-players-1-server-1-client)** – Simplest: one app hosts, one app connects. Good for testing or head-to-head.
- **[Lobby + rooms (many players)](#lobby--rooms-many-players)** – Server creates a room, accepts many clients, broadcasts to the room. Good for chat, co-op, or game lobbies.
- **[Listen server (host plays too)](#listen-server-host-plays-too)** – One app both hosts and plays; use **AcceptTimeout** in the game loop to add clients without blocking.
- **[LAN vs Internet](#lan-vs-internet)** – Same code; only the address and (for internet) TLS and port forwarding change.

---

## Two players (1 server, 1 client)

**Server (run first):**

```basic
VAR sid = Host(9999)
IF IsNull(sid) THEN PRINT "Failed to host" : END
VAR cid = Accept(sid)
IF IsNull(cid) THEN CloseServer(sid) : END
Send(cid, "hello")
VAR msg = Receive(cid)
PRINT "Got: " + msg
Disconnect(cid)
CloseServer(sid)
```

**Client:**

```basic
VAR cid = Connect("127.0.0.1", 9999)
IF IsNull(cid) THEN PRINT "Failed to connect" : END
VAR msg = Receive(cid)
PRINT "Got: " + msg
Send(cid, "world")
Disconnect(cid)
```

Use **127.0.0.1** for same machine; use the server’s LAN IP (e.g. from **GetLocalIP()**) for another computer on the network.

---

## Lobby + rooms (many players)

**Server:** Create a room, then each frame try to accept new clients and add them to the room; broadcast state with **SendToRoom**.

```basic
VAR sid = Host(9999)
IF IsNull(sid) THEN END
CreateRoom("lobby")
VAR state = "welcome"
WHILE NOT WindowShouldClose()
  VAR cid = AcceptTimeout(sid, 50)
  IF NOT IsNull(cid) THEN JoinRoom("lobby", cid)
  SendToRoom("lobby", state)
  // Optional: loop clients and Receive(cid) to handle incoming messages
WEND
CloseServer(sid)
```

**Client:** Connect once, then each frame **Receive** and **Send** as needed.

```basic
VAR cid = Connect("127.0.0.1", 9999)
IF IsNull(cid) THEN END
WHILE NOT WindowShouldClose()
  VAR msg = Receive(cid)
  IF NOT IsNull(msg) THEN PRINT msg
  Send(cid, "my_input")
WEND
Disconnect(cid)
```

---

## Listen server (host plays too)

The same program hosts and plays. Use **AcceptTimeout**(serverId, 10) in the game loop so new clients join without blocking the game.

```basic
VAR sid = Host(9999)
IF IsNull(sid) THEN END
CreateRoom("lobby")
WHILE NOT WindowShouldClose()
  VAR cid = AcceptTimeout(sid, 10)
  IF NOT IsNull(cid) THEN JoinRoom("lobby", cid)
  SendToRoom("lobby", "state " + STR(GetConnectionCount()))
  // ... your game draw and input ...
WEND
CloseServer(sid)
```

---

## LAN vs Internet

- **Same machine:** Use **127.0.0.1** as the host address.
- **LAN (same Wi‑Fi / network):** Use the host’s LAN IP. Call **GetLocalIP()** on the server and show "Connect to: " + GetLocalIP() + " :9999" so players can type it in.
- **Internet:** Use the server’s public IP and set up port forwarding on the router. Prefer **HostTLS** and **ConnectTLS** so traffic is encrypted; see [Security](#security).

Same **Connect** / **Host** / **Send** / **Receive** / rooms API everywhere; only the address (and TLS for internet) changes.

---

## Protocol

Messages are sent as **lines of text** (one message per line). When you call **Send**(connectionId, text), a newline is appended. When you call **Receive**(connectionId), you get one line (without the newline), or null if no data is available (non-blocking) or the connection closed.

**Limits (optimized and secure):** Each message is limited to **256 KB** and must not contain newline or carriage-return characters. **Send** and **SendToRoom** reject oversized or invalid messages (return false or 0). This keeps the protocol predictable and prevents abuse.

Use a simple protocol in your game, for example:

- Send positions: `Send(cid, "pos " + STR(x) + " " + STR(y))`
- Parse on receive: split the string and interpret the first word as command.

## Sending in different forms (JSON and text)

- **Plain text:** Use **Send**(connectionId, text) and **Receive**(connectionId). Good for simple commands like `"pos 100 200"` or chat.
- **JSON:** Use **SendJSON**(connectionId, jsonText) to send a JSON string (it is validated before send; returns 0 if invalid). On the other side use **ReceiveJSON**(connectionId) to read the next line and get it only if it is valid JSON (otherwise null). Parse the returned string with **LoadJSONFromString** and **GetJSONKey** (see standard library). Example:

```basic
// Sender: build JSON string and send
VAR jsonStr = "{\"x\":" + STR(x) + ",\"y\":" + STR(y) + "}"
SendJSON(cid, jsonStr)

// Receiver: receive and parse
VAR jsonStr = ReceiveJSON(cid)
IF NOT IsNull(jsonStr) THEN
  VAR h = LoadJSONFromString(jsonStr)
  VAR x = GetJSONKey(h, "x")
  VAR y = GetJSONKey(h, "y")
END IF
```

- **Broadcast JSON to a room:** **SendToRoomJSON**(roomId, jsonText) — validates JSON and sends to every connection in the room; returns the number of connections the message was sent to (0 if JSON is invalid or too long).
- **Other formats:** Send any single-line text with **Send**; receive with **Receive**. You can use a prefix in the line (e.g. `"TEXT|"` or `"JSON|"`) and split on the receiver to decide how to handle it.

### Sending numbers (integers, floats, multiple)

Use typed send/receive when you only need numbers (no extra parsing):

- **SendInt**(connectionId, value) — send one integer. Returns 1 if sent, 0 on failure.
- **SendFloat**(connectionId, value) — send one float. Returns 1 if sent, 0 on failure.
- **SendNumbers**(connectionId, n1, n2, …) — send up to 16 numbers in one message (e.g. position x,y or x,y,z). Returns 1 if sent, 0 on failure.
- **SendText**(connectionId, text) — same as **Send**; use for clarity when sending plain text.

On the receiver:

- **ReceiveNumbers**(connectionId) — read the next line and parse it as numbers (handles "i 42", "f 3.14", or "n 1 2 3.5"). Returns the **count** of numbers received (0 if no data or parse error). Non-blocking.
- **GetReceivedNumber**(index) — get the number at 0-based index from the last **ReceiveNumbers** call. Returns 0.0 if index is out of range.

Example:

```basic
// Sender
SendNumbers(cid, x, y, health)

// Receiver
VAR n = ReceiveNumbers(cid)
IF n >= 3 THEN
  VAR x = GetReceivedNumber(0)
  VAR y = GetReceivedNumber(1)
  VAR h = GetReceivedNumber(2)
END IF
```

Room variants: **SendToRoomInt**(roomId, value), **SendToRoomFloat**(roomId, value), **SendToRoomNumbers**(roomId, n1, n2, …) — each returns the number of connections the message was sent to (0 if failed).

## Server

1. **Host**(port) — start listening on the given port. Returns a serverId (e.g. "server_1") or null on failure.
2. **Accept**(serverId) — wait for one client to connect (blocking). Returns a connectionId (e.g. "conn_1") or null on error.
3. Use **Send**(connectionId, text) and **Receive**(connectionId) to talk to that client.
4. **CloseServer**(serverId) — stop listening.

For multiple clients, call **Accept** in a loop (or use **AcceptTimeout** in your game loop) and store each connectionId; then each frame call **Receive** on each connection.

## Rooms

Rooms let you group connections (e.g. a "lobby" or "game_1") and broadcast to everyone in the room.

- **CreateRoom**(roomId) — ensure the room exists (idempotent; safe to call multiple times).
- **JoinRoom**(roomId, connectionId) — add a connection to a room. If the room doesn't exist, it is created. Fails if connectionId is invalid.
- **LeaveRoom**(connectionId) — remove the connection from **all** rooms.
- **LeaveRoom**(connectionId, roomId) — remove the connection from that room only. One connection can be in multiple rooms.
- **SendToRoom**(roomId, text) — send a line of text to every connection in the room (same as calling **Send** for each). Returns the number of connections the message was sent to.
- **GetRoomConnectionCount**(roomId) — number of connections in the room (0 if the room doesn't exist).
- **GetRoomConnectionId**(roomId, index) — the connectionId at 0-based index in the room (empty string if index is out of range). Use with **GetRoomConnectionCount** to iterate: `FOR i = 0 TO GetRoomConnectionCount(roomId)-1 ... GetRoomConnectionId(roomId, i) ... NEXT`.

## Convenience

- **IsConnected**(connectionId) — returns 1 if the connection is still in the connection map, 0 otherwise. Use before **Send** to avoid sending to a closed connection.
- **GetConnectionCount**() — returns the total number of connections (on server or client). Useful for "how many players".
- **AcceptTimeout**(serverId, timeoutMs) — like **Accept** but waits at most timeoutMs milliseconds. Returns a new connectionId when a client connects, or null on timeout. Lets the game loop poll for new clients without blocking (e.g. call each frame with a short timeout).
- **GetLocalIP**() — returns this machine’s local IP (e.g. 192.168.1.x) so you can show "Connect to: GetLocalIP() : port" for LAN play.

## Security

- **When to use plain TCP (Connect / Host):** LAN, localhost, or trusted development. Traffic is not encrypted.
- **When to use TLS (ConnectTLS / HostTLS):** Internet or any untrusted network. Traffic is encrypted. You need a certificate and key on the server; clients verify the server by default. For testing you can use a self-signed certificate (e.g. `go run crypto/tls/generate_cert.go` or openssl); for production use a proper certificate (e.g. Let’s Encrypt).
- **Best practices:** Validate and sanitize all received text; never trust the client for game authority (server should decide outcomes); optional token or password in the first message; rate limiting (e.g. limit messages per second per connection) in your BASIC logic. Message size is capped at 256 KB and newlines in payloads are rejected to keep the protocol safe and predictable.

For other transports (e.g. WebSocket), a future binding could use a Go library such as gorilla/websocket; the same room and Send/Receive concepts would apply.

## Client

1. **Connect**(host, port) — connect to a server. Returns connectionId or null on failure. Use **ConnectTLS**(host, port) for encrypted connection.
2. **Send**(connectionId, text) — send a line of text.
3. **Receive**(connectionId) — get the next line (or null if none yet). Non-blocking.
4. **Disconnect**(connectionId) — close the connection.

## Game loop usage

In a game, typically:

- **Server:** Each frame, call **Receive**(cid) for each client; if not null, handle the message. Use **Send**(cid, ...) or **SendToRoom**(roomId, ...) to broadcast state or replies.
- **Client:** Each frame, call **Receive**(cid); if not null, update game state. Use **Send**(cid, ...) to send input or actions.

Receive is non-blocking (returns null if no data), so your game stays responsive.

## API summary

| Function | Description |
|----------|-------------|
| **Connect**(host, port) | Connect to a server. Returns connectionId or null. |
| **ConnectTLS**(host, port) | Connect with TLS encryption. Returns connectionId or null. |
| **Send**(connectionId, text) | Send a line of text (max 256 KB, no newlines). Returns true/false. |
| **SendJSON**(connectionId, jsonText) | Send valid JSON string; returns 1 if sent, 0 if invalid or failed. |
| **Receive**(connectionId) | Read next line (or null). Non-blocking. |
| **ReceiveJSON**(connectionId) | Read next line; return it only if valid JSON, else null. Non-blocking. |
| **Disconnect**(connectionId) | Close the connection (and remove from all rooms). |
| **Host**(port) | Start a server. Returns serverId or null. |
| **HostTLS**(port, certFile, keyFile) | Start a TLS server. Returns serverId or null. |
| **Accept**(serverId) | Wait for a client (blocking). Returns connectionId or null. |
| **AcceptTimeout**(serverId, timeoutMs) | Wait for a client with timeout. Returns connectionId or null. |
| **CloseServer**(serverId) | Stop the server. |
| **CreateRoom**(roomId) | Ensure room exists. Idempotent. |
| **JoinRoom**(roomId, connectionId) | Add connection to room. |
| **LeaveRoom**(connectionId) | Remove connection from all rooms. |
| **LeaveRoom**(connectionId, roomId) | Remove connection from one room. |
| **SendToRoom**(roomId, text) | Send text to every connection in the room (max 256 KB, no newlines). Returns count sent. |
| **SendToRoomJSON**(roomId, jsonText) | Send valid JSON to every connection in the room. Returns count sent (0 if invalid). |
| **SendInt**(connectionId, value) | Send one integer. Returns 1 if sent, 0 on failure. |
| **SendFloat**(connectionId, value) | Send one float. Returns 1 if sent, 0 on failure. |
| **SendNumbers**(connectionId, n1, n2, …) | Send up to 16 numbers in one message. Returns 1 if sent, 0 on failure. |
| **SendText**(connectionId, text) | Same as Send; plain text. Returns true/false. |
| **ReceiveNumbers**(connectionId) | Read next line as numbers; returns count (0 if no data or parse error). Use GetReceivedNumber(index). |
| **GetReceivedNumber**(index) | Get number at 0-based index from last ReceiveNumbers. Returns 0.0 if out of range. |
| **SendToRoomInt**(roomId, value) | Broadcast one integer to room. Returns count sent. |
| **SendToRoomFloat**(roomId, value) | Broadcast one float to room. Returns count sent. |
| **SendToRoomNumbers**(roomId, n1, n2, …) | Broadcast up to 16 numbers to room. Returns count sent. |
| **GetRoomConnectionCount**(roomId) | Number of connections in room. |
| **GetRoomConnectionId**(roomId, index) | ConnectionId at 0-based index in room. |
| **IsConnected**(connectionId) | 1 if connected, 0 otherwise. |
| **GetConnectionCount**() | Total number of connections. |
| **GetLocalIP**() | This machine’s local IP for LAN (e.g. 192.168.1.x). |

See [API_REFERENCE.md](../API_REFERENCE.md) for full details.
