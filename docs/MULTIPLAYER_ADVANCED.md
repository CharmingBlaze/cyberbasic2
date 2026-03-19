# Multiplayer Advanced: Lockstep, Rollback, Prediction

Deep dive on deterministic multiplayer patterns in CyberBASIC2.

---

## Transport: KCP vs TCP

CyberBASIC2 uses **KCP (reliable UDP)** via kcp-go for Host/Connect. Benefits over TCP:

- **Lower latency** — KCP can achieve lower round-trip latency on lossy or high-latency links.
- **Packet-loss resilience** — Fast retransmission tuned for real-time games.
- **Stream mode** — Same line-based protocol; messages are newline-terminated.

For legacy TCP (e.g. firewalls that block UDP), consider a fallback; the current API uses KCP only.

## Nakama + Direct KCP

You can use **Nakama** for matchmaking (accounts, find opponents) and then **direct KCP** for gameplay. Flow: NakamaAddMatchmaker → OnNakamaMatchmakerMatched → get host/port from your game logic → Connect(host, port). This hybrid requires your match handler to exchange connection details (e.g. via Nakama RPC or match state).

---

## Overview

CyberBASIC2 supports:

- **Lockstep** — Server and clients advance simulation only when all inputs for tick N are received.
- **Rollback** — Save/restore game state; server can tell clients to rollback and resimulate.
- **Prediction** — Client runs local simulation; reconciles with server state when correction arrives.

Use these for RTS, fighting games, or any game requiring deterministic simulation.

---

## Lockstep

### Flow

1. **Server:** `LockstepEnable(60)`. Each tick, wait for `OnLockstepTickReady(tickId)`. Call `LockstepGetInputs(tickId)` to get `{connectionId: inputData}`. Run simulation. Broadcast state (or `T\t<tickId>` is sent automatically when all inputs received).
2. **Client:** `LockstepEnable(60)`. Each tick, call `LockstepSendInput(tickId, inputData)`. When `OnLockstepTickReady(tickId)` fires (server sent `T\t<tickId>`), advance simulation.

### Protocol

- `L\t<tickId>\t<data>` — client sends input for tick
- `T\t<tickId>` — server broadcasts when tick has all inputs

### Example (server)

```basic
LockstepEnable(60)
VAR tick = 0
VAR sid = Host(9999)
VAR cid = Accept(sid)

SUB OnLockstepTickReady(tickId)
  VAR inputs = LockstepGetInputs(tickId)
  // inputs is dict: connectionId -> input string
  // Run simulation with inputs
  tick = tick + 1
END SUB

WHILE NOT WindowShouldClose()
  ProcessNetworkEvents()
  // ...
WEND
```

### Example (client)

```basic
LockstepEnable(60)
VAR tick = 0
VAR cid = Connect("127.0.0.1", 9999)

SUB OnLockstepTickReady(tickId)
  // Advance simulation
  tick = tick + 1
END SUB

WHILE NOT WindowShouldClose()
  ProcessNetworkEvents()
  VAR inputData = Str(KeyDown(KEY_W)) + "," + Str(KeyDown(KEY_S))
  LockstepSendInput(Str(tick), inputData)
  tick = tick + 1
  // ...
WEND
```

---

## Rollback

### Snapshot Handlers

You must register handlers that serialize and restore game state:

```basic
RegisterSnapshotHandler("MySnapshot")
RegisterRestoreHandler("MyRestore")

SUB MySnapshot(tickId)
  // Serialize state to a string (e.g. JSON: "{\"x\":100,\"y\":200}")
  VAR data = "{\"x\":" + Str(playerX) + ",\"y\":" + Str(playerY) + "}"
  SnapshotStoreResult(tickId, data)
END SUB

SUB MyRestore(tickId, data)
  // Parse data and restore playerX, playerY, etc.
  // Use Instr, Mid, Val or a JSON parser to extract values
END SUB
```

### Creating and Restoring Snapshots

- `SnapshotCreate(tickId)` — invokes your snapshot sub; stores result.
- `SnapshotRestore(tickId)` — invokes your restore sub with stored data.

### Server-Initiated Rollback

When the server detects a desync, it calls `RollbackBroadcast(tickId, correctTickId)`. Clients receive `OnRollbackRequired(tickId, correctTickId)` and should call `SnapshotRestore(correctTickId)` then resimulate from that tick.

---

## Prediction

### Flow

1. **Client:** `PredictionEnable()`. Each tick, `PredictionStoreInput(tickId, input)` when sending input. Run local simulation.
2. **Server:** Sends authoritative state periodically.
3. **Client:** When server state arrives, call `PredictionReconcile(tickId, stateJson)`. This invokes your restore handler and fires `OnPredictionCorrected(tickId)`.

### Determinism

Prediction works best when client and server simulation is deterministic. Use `FixedUpdate` and `OnFixedUpdate`; avoid `GetFrameTime()` in simulation logic.

---

## See Also

- [Multiplayer Design](MULTIPLAYER_DESIGN.md)
- [Multiplayer](MULTIPLAYER.md)
- [Tutorial Multiplayer](TUTORIAL_MULTIPLAYER.md)
