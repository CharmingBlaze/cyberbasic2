# Multiplayer (TCP)

CyberBasic includes simple TCP client/server bindings so you can send and receive text messages between programs (e.g. for chat, lobby, or simple game state). No extra libraries; uses Go's standard `net` package. For encrypted channels over the internet, use **ConnectTLS** and **HostTLS**.

## Quick pick: choose your multiplayer type

Pick a use case and copy the snippet. Same commands work for all; only the flow changes.

- **[Event-based callbacks (recommended)](#event-based-callbacks)** – Define **OnClientConnect**, **OnClientDisconnect**, **OnMessage** and call **ProcessNetworkEvents()** each frame. Easiest for games and chat.
- **[Two players (1 server, 1 client)](#two-players-1-server-1-client)** – Simplest: one app hosts, one app connects. Good for testing or head-to-head.
- **[Lobby + rooms (many players)](#lobby--rooms-many-players)** – Server creates a room, accepts many clients, broadcasts to the room. Good for chat, co-op, or game lobbies.
- **[Listen server (host plays too)](#listen-server-host-plays-too)** – One app both hosts and plays; use **AcceptTimeout** in the game loop to add clients without blocking.
- **[LAN vs Internet](#lan-vs-internet)** – Same code; only the address and (for internet) TLS and port forwarding change.

---

## Event-based callbacks

You can use **event-based callbacks** instead of polling **Receive**. Define these Subs (names are case-insensitive):

- **OnClientConnect**(id) – called when a new client connects (server side, after **Accept** or **AcceptTimeout** adds the connection).
- **OnClientDisconnect**(id) – called when a connection is closed or lost.
- **OnMessage**(id, msg) – called when a message is received; `id` is the connectionId, `msg` is the message text.

**You must call ProcessNetworkEvents() once per frame** (e.g. at the start of **update**(dt) or at the top of your main loop). It drains the internal event queue and invokes the Subs above if they are defined.

**StartServer**(port) is the preferred way to start a server (same as **Host**(port); returns serverId or null). **Broadcast**(text) sends `text` to every connection (all clients).

**Minimal server example:**

```basic
VAR sid = StartServer(1234)
IF IsNull(sid) THEN PRINT "Failed to start server" : END

SUB OnClientConnect(id)
  PRINT "Client joined: " + id
END SUB

SUB OnClientDisconnect(id)
  PRINT "Client left: " + id
END SUB

SUB OnMessage(id, msg)
  PRINT "From " + id + ": " + msg
  Broadcast(msg)
END SUB

WHILE TRUE
  ProcessNetworkEvents()
  // ... your game update/draw ...
WEND
```

To accept new clients each frame, call **AcceptTimeout**(sid, 50) in the loop; when it returns a connectionId, the next **ProcessNetworkEvents()** will call **OnClientConnect** with that id. You can still use **Receive**(connectionId) to poll for messages if you prefer; messages are queued and either **OnMessage** is called or **Receive** will return them.

**Minimal client example:**

```basic
VAR cid = Connect("127.0.0.1", 1234)
IF IsNull(cid) THEN PRINT "Failed to connect" : END

SUB OnMessage(id, msg)
  PRINT "Got: " + msg
END SUB

WHILE TRUE
  ProcessNetworkEvents()
  Send(cid, "hello")
  // ... your game ...
WEND
Disconnect(cid)
```

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
- **Tables/dictionaries:** **SendTable**(connectionId, data) — `data` is a dictionary (e.g. from **CreateDict** or a dict literal). It is serialized to JSON and sent. Returns 1 if sent, 0 on failure (e.g. message too long). **ReceiveTable**(connectionId) — reads the next message from the queue and, if it is valid JSON, parses it into a dictionary and returns it (so you can use **GetJSONKey** on the result). Returns null if no message or invalid JSON.
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

1. **StartServer**(port) — preferred way to start a server (same as **Host**). Returns a serverId or null on failure.
2. **Host**(port) — start listening on the given port. Returns a serverId (e.g. "server_1") or null on failure.
3. **Accept**(serverId) — wait for one client to connect (blocking). Returns a connectionId (e.g. "conn_1") or null on error.
4. Use **Send**(connectionId, text) and **Receive**(connectionId) to talk to that client.
5. **CloseServer**(serverId) — stop listening.

For event-based handling, define **OnClientConnect**, **OnClientDisconnect**, **OnMessage** and call **ProcessNetworkEvents()** each frame; see [Event-based callbacks](#event-based-callbacks). **Broadcast**(text) sends to every connection.

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

## RPC (remote procedure calls)

You can call a Sub on the other side by name with **SendRPC**(connectionId, name, args...). The receiver must **RegisterRPC**(name, subName) so that when an RPC packet arrives, the engine invokes the registered Sub with the deserialized arguments.

- **RegisterRPC**(name, subName) — when an RPC with this `name` is received, call the Sub `subName` with the arguments. Example: `RegisterRPC("spawnEnemy", "SpawnEnemy")`. The Sub must have the same number of parameters as the args sent.
- **SendRPC**(connectionId, name, arg1, arg2, ...) — send an RPC; the other side’s registered Sub is invoked when **ProcessNetworkEvents()** runs. Returns true if sent.

RPC runs on the same thread as **ProcessNetworkEvents()** (no extra threading). If no handler is registered for an RPC name, the message is still delivered to **OnMessage** (if defined) as raw text.

## Ping and disconnect

- **OnClientDisconnect**(id) is called when a connection is closed or lost (e.g. the reader goroutine gets EOF or error).
- **SendPing**(connectionId) sends a ping line; the other side replies with a pong. Call this periodically (e.g. every few seconds) to measure latency.
- **GetPing**(connectionId) returns the last round-trip time in milliseconds (0 if no pong received yet). Use after **SendPing** and when you receive the pong on the reader side (handled automatically).

## Entity synchronization

You can sync entity position (and optionally other state) so the engine tracks and sends updates.

- **SyncEntity**(connectionId, entityId, x, y) or **SyncEntity**(connectionId, entityId, x, y, z) — send position for `entityId` to one connection.
- **SyncEntityToRoom**(roomId, entityId, x, y) or **SyncEntityToRoom**(roomId, entityId, x, y, z) — send to every connection in the room. Returns the number of connections the message was sent to.
- On the receiver, define **OnEntitySync**(entityId, x, y, z) (4 parameters). It is called when **ProcessNetworkEvents()** processes an entity sync message. You can also read the last synced state with **GetRemoteEntity**(entityId), which returns a dictionary with keys `"x"`, `"y"`, `"z"` (use **GetJSONKey** to read them). No interpolation in this phase; interpolation can be a later enhancement.

## API summary

| Function | Description |
|----------|-------------|
| **ProcessNetworkEvents**() | Drain the network event queue and call OnClientConnect / OnClientDisconnect / OnMessage if defined. Call once per frame. |
| **StartServer**(port) | Start a server (alias for Host). Returns serverId or null. |
| **Broadcast**(text) | Send text to every connection. Same limits as Send. |
| **Connect**(host, port) | Connect to a server. Returns connectionId or null. |
| **ConnectTLS**(host, port) | Connect with TLS encryption. Returns connectionId or null. |
| **Send**(connectionId, text) | Send a line of text (max 256 KB, no newlines). Returns true/false. |
| **SendJSON**(connectionId, jsonText) | Send valid JSON string; returns 1 if sent, 0 if invalid or failed. |
| **SendTable**(connectionId, data) | Serialize dictionary to JSON and send. Returns 1 if sent, 0 on failure. |
| **Receive**(connectionId) | Read next line (or null). Non-blocking. |
| **ReceiveJSON**(connectionId) | Read next line; return it only if valid JSON, else null. Non-blocking. |
| **ReceiveTable**(connectionId) | Read next message; if valid JSON, return as dictionary, else null. Non-blocking. |
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
| **RegisterRPC**(name, subName) | Register a Sub to be called when RPC `name` is received. |
| **SendRPC**(connectionId, name, args...) | Send an RPC; receiver’s registered Sub is invoked with args. |
| **SendPing**(connectionId) | Send a ping; the peer replies with pong. Returns true if sent. |
| **GetPing**(connectionId) | Last RTT in milliseconds (0 if no pong received yet). |
| **SyncEntity**(connectionId, entityId, x, y) / (…, z) | Send entity position to one connection. Returns true if sent. |
| **SyncEntityToRoom**(roomId, entityId, x, y) / (…, z) | Send entity position to every connection in the room. Returns count sent. |
| **GetRemoteEntity**(entityId) | Last synced state (dict with x, y, z). Returns null if none. |

For the complete list of all network commands and signatures see [API Reference](../API_REFERENCE.md) section 18.

---

## See also

- [API Reference](../API_REFERENCE.md) (section 18) — full network command list
- [Command Reference](COMMAND_REFERENCE.md) — commands by feature
- [Getting Started](GETTING_STARTED.md) — setup and first run
