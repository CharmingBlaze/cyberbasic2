# Multiplayer Tutorial

This tutorial shows the current multiplayer workflow. **Transport:** KCP (reliable UDP) for Host/Connect. Optional **Nakama** for cloud accounts and matchmaking.

## What Exists Today

- KCP client/server connections with `Host`, `Accept`, `AcceptTimeout`, `Connect`, and `Disconnect`
- Text, JSON, table, RPC, ping, and entity-sync helpers
- Event-driven delivery through `ProcessNetworkEvents()`
- Polling delivery through `Receive*()` calls
- Fixed-step runtime callbacks through `FixedUpdate(rate)` and `OnFixedUpdate(label$)`

## Advanced Features (Implemented)

- **Lockstep** — `LockstepEnable`, `LockstepSendInput`, `LockstepGetInputs`, `OnLockstepTickReady`
- **Rollback** — `RegisterSnapshotHandler`, `RegisterRestoreHandler`, `SnapshotCreate`, `SnapshotRestore`, `RollbackBroadcast`
- **Prediction** — `PredictionEnable`, `PredictionStoreInput`, `PredictionReconcile`
- **Matchmaking** — `MatchmakingHost`, `MatchmakingDiscover`, `MatchmakingJoin` (LAN broadcast)
- **Interest management** — `SetInterestFilter`, `SetEntityInterestZone`
- **Nakama (optional)** — `NakamaConnect`, `NakamaAuthenticateDevice`, `NakamaCreateMatch`, `NakamaJoinMatch`, etc. See [NAKAMA_GUIDE.md](NAKAMA_GUIDE.md).

See [MULTIPLAYER_ADVANCED.md](MULTIPLAYER_ADVANCED.md) for lockstep/rollback/prediction patterns. For full command list, see [Multiplayer Guide](MULTIPLAYER.md).

## Nakama Quick Start

For cloud matchmaking and accounts:

```basic
NakamaConnect("127.0.0.1", 7350, "defaultkey", 0)
NakamaAuthenticateDevice("my-device-id", 1, "")
NakamaCreateSocket()
NakamaSocketConnect()
VAR matchId = NakamaCreateMatch("")
// Or: NakamaAddMatchmaker(2, 4, "") then OnNakamaMatchmakerMatched(matchId, token) → NakamaJoinMatch(matchId, token)
mainloop
  NakamaProcessEvents()
  // ...
endmain
```

## Pick One Delivery Style

Use one of these patterns for message delivery:

1. Callback-driven: define `OnClientConnect`, `OnClientDisconnect`, `OnMessage`, then call `ProcessNetworkEvents()` once per frame.
2. Polling-driven: skip message callbacks and read with `Receive`, `ReceiveJSON`, `ReceiveTable`, or `ReceiveNumbers`.

Do not mix both for the same messages. `ProcessNetworkEvents()` consumes queued messages before `Receive()` can see them.

## Minimal Server

```basic
VAR sid = Host(9999)
IF IsNull(sid) THEN
  PRINT "Failed to host"
  END
ENDIF

SUB OnClientConnect(id)
  PRINT "Connected: " + id
END SUB

SUB OnClientDisconnect(id)
  PRINT "Disconnected: " + id
END SUB

SUB OnMessage(id, msg)
  PRINT id + ": " + msg
  Send(id, "echo " + msg)
END SUB

WHILE NOT WindowShouldClose()
  VAR cid = AcceptTimeout(sid, 10)
  ProcessNetworkEvents()
WEND

CloseServer(sid)
```

## Minimal Client

```basic
VAR cid = Connect("127.0.0.1", 9999)
IF IsNull(cid) THEN
  PRINT "Connect failed"
  END
ENDIF

SUB OnMessage(id, msg)
  PRINT "Server said: " + msg
END SUB

WHILE NOT WindowShouldClose()
  ProcessNetworkEvents()

  IF IsKeyPressed(KEY_SPACE) THEN
    Send(cid, "ping")
  ENDIF
WEND

Disconnect(cid)
```

## Polling Example

```basic
VAR sid = Host(9999)
IF IsNull(sid) THEN END

VAR cid = Accept(sid)
IF IsNull(cid) THEN END

WHILE TRUE
  VAR msg = Receive(cid)
  IF NOT IsNull(msg) THEN
    PRINT "Got: " + msg
    Send(cid, "ack")
  ENDIF
WEND
```

## Fixed-Step Simulation

For gameplay sync, run game logic on the fixed step instead of raw frame time:

```basic
FixedUpdate 60
OnFixedUpdate "NetFixedStep"

SUB NetFixedStep(dt)
  ' dt matches FixedDeltaTime()
  Send(myConnection, "input " + STR(playerInput))
END SUB
```

The runtime now steps fixed physics callbacks on the accumulator-driven timestep, so this is the correct place for deterministic simulation work.

## Rooms

```basic
CreateRoom("lobby")

VAR sid = Host(9999)
WHILE NOT WindowShouldClose()
  VAR cid = AcceptTimeout(sid, 10)
  IF NOT IsNull(cid) THEN
    JoinRoom("lobby", cid)
  ENDIF

  SendToRoom("lobby", "heartbeat")
WEND
```

## RPC

```basic
RegisterRPC("spawn_enemy", "SpawnEnemy")

SUB SpawnEnemy(x, y)
  PRINT "Spawn enemy at " + STR(x) + ", " + STR(y)
END SUB

SendRPC(cid, "spawn_enemy", 100, 200)
```

RPC handlers are invoked when `ProcessNetworkEvents()` runs.

## Entity Sync

```basic
SyncEntity(cid, "player_1", x, y, z)

SUB OnEntitySync(entityId, x, y, z)
  PRINT entityId + " -> " + STR(x) + "," + STR(y) + "," + STR(z)
END SUB
```

## Matchmaking (LAN Discovery)

**Host a room:**

```basic
VAR sid = MatchmakingHost(9999, "My Game", 4)
IF IsNull(sid) THEN END
VAR cid = Accept(sid)
' ... game loop ...
```

**Discover and join:**

```basic
VAR rooms = MatchmakingDiscover(3000)
VAR n = GetJSONKey(rooms, "count")
FOR i = 0 TO n - 1
  VAR r = GetJSONKey(rooms, Str(i))
  PRINT GetJSONKey(r, "roomName") + " - " + Str(GetJSONKey(r, "playerCount")) + " players"
NEXT
VAR r0 = GetJSONKey(rooms, "0")
VAR cid = MatchmakingJoin(GetJSONKey(r0, "host"), GetJSONKey(r0, "port"))
```

## Practical Notes

- `Receive(connectionId)` is non-blocking and takes no timeout argument.
- Idle TCP connections no longer disconnect just because no payload arrived during a short polling window.
- `SendPing` / `GetPing` are the current latency helpers.
- Use `HostTLS` / `ConnectTLS` when you need encrypted transport.

### Commands you learned

- **Server:** Host, Accept, AcceptTimeout, CloseServer, ProcessNetworkEvents
- **Client:** Connect, Disconnect, Send, Receive
- **Messages:** SendJSON, SendTable, ReceiveJSON, ReceiveTable
- **Rooms:** CreateRoom, JoinRoom, SendToRoom
- **RPC:** RegisterRPC, SendRPC
- **Entity sync:** SyncEntity, OnEntitySync

Full reference: [Multiplayer Guide](MULTIPLAYER.md).

## Next Reading

- [MULTIPLAYER.md](MULTIPLAYER.md) for the API guide
- [MULTIPLAYER_DESIGN.md](MULTIPLAYER_DESIGN.md) for architecture
- [MULTIPLAYER_ADVANCED.md](MULTIPLAYER_ADVANCED.md) for lockstep, rollback, prediction
- [3D_GAME_API.md](3D_GAME_API.md) and [2D_GAME_API.md](2D_GAME_API.md) for gameplay-facing command references
