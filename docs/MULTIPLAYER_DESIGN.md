# Multiplayer Design

Current multiplayer architecture, implementation status, deterministic patterns, and roadmap gaps.

---

## Purpose

- **KCP (reliable UDP) transport:** Connect, Host, Send, Receive for client/server games.
- **Deterministic foundation:** Fixed-step simulation so server and clients can stay in sync.
- **Event + polling:** ProcessNetworkEvents for callbacks; Receive* for polling.

---

## Architecture

```
Client                          Server
──────                          ──────
Connect(host, port)    →        Host(port)
Send(msg)              →        Accept / AcceptTimeout
Receive()              ←        Send(msg)
ProcessNetworkEvents   ←        ProcessNetworkEvents
OnMessage, OnEntitySync         OnClientConnect, OnMessage
```

**Packages:** `compiler/bindings/net` — Host, Connect, Send, Receive, RPC, SyncEntity, ProcessNetworkEvents.

---

## Transport: kcp-go

CyberBASIC2 uses **KCP (reliable UDP)** via [kcp-go](https://github.com/xtaci/kcp-go) for lower latency and better packet-loss resilience than TCP. KCP provides ordered, reliable delivery over UDP with configurable trade-offs for latency vs throughput. Benefits:

- **Lower latency** — KCP can achieve lower round-trip latency than TCP on lossy or high-latency links.
- **Packet-loss resilience** — Fast retransmission and congestion control tuned for real-time games.
- **Stream mode** — Same line-based protocol as before; messages are newline-terminated.

## Current Transport API

CyberBASIC2 ships a KCP line-based networking layer:

- `Host`, `Accept`, `AcceptTimeout`, `CloseServer`
- `Connect`, `Disconnect`, `IsConnected`
- `Send`, `Receive`
- `SendJSON`, `ReceiveJSON`
- `SendTable`, `ReceiveTable`
- `SendNumbers`, `ReceiveNumbers`
- `RegisterRPC`, `SendRPC`
- `SyncEntity`
- `ProcessNetworkEvents`

Each outgoing message is sent as one newline-terminated text record.

## Delivery Model

There are two supported receive styles:

1. Event-driven delivery through `ProcessNetworkEvents()` and callbacks like `OnClientConnect`, `OnClientDisconnect`, `OnMessage`, and `OnEntitySync`.
2. Polling delivery through `Receive*()` functions.

Important rule: messages are now consumed exactly once. If `ProcessNetworkEvents()` handles a queued message first, `Receive()` will not see that same message later.

## Runtime Behavior

- Reader goroutines keep TCP connections alive during idle periods instead of treating short read timeouts as disconnects.
- Disconnect cleanup is centralized so a dead socket does not emit duplicate disconnect events.
- `AcceptTimeout` now uses listener deadlines for both plain TCP and TLS-hosted servers.
- Numeric receive buffers are stored per connection, so multi-connection polling no longer shares one global scratch buffer unless you intentionally use the legacy one-argument `GetReceivedNumber(index)` form.
- RPC and entity-sync payloads ride on the same message queue as normal text messages.
- Callback errors from `OnClientConnect`, `OnClientDisconnect`, `OnMessage`, `OnEntitySync`, or registered RPC handlers now propagate out of `ProcessNetworkEvents()` instead of being silently swallowed.

## Simulation Model

The runtime now supports a fixed-step loop:

- `FixedUpdate(rate)` sets the fixed-step frequency.
- `OnFixedUpdate(label$)` registers the callback to run each step.
- `FixedDeltaTime()` returns the current fixed timestep.

Recommended pattern:

- Read network input on the main frame loop.
- Queue or store authoritative input/state.
- Advance gameplay simulation from `OnFixedUpdate`.

This gives multiplayer code a stable timestep without requiring lockstep or rollback infrastructure yet.

## Implemented Advanced Features

The following are now implemented:

- **Lockstep** — `LockstepEnable(tickRate)`, `LockstepSendInput(tickId, data)`, `LockstepGetInputs(tickId)`, `OnLockstepTickReady(tickId)`. Server collects inputs per tick; when all clients have sent, broadcasts tick-ready. Clients wait for `T\t<tickId>` before advancing.
- **Rollback** — `RegisterSnapshotHandler(sub)`, `RegisterRestoreHandler(sub)`, `SnapshotCreate(tickId)`, `SnapshotStoreResult(tickId, data)`, `SnapshotRestore(tickId)`, `RollbackBroadcast(tickId, correctTickId)`, `OnRollbackRequired(tickId, correctTickId)`. Handlers serialize/restore game state; server can broadcast rollback to clients.
- **Prediction** — `PredictionEnable()`, `PredictionStoreInput(tickId, input)`, `PredictionReconcile(tickId, stateJson)`, `OnPredictionCorrected(tickId)`. Client-side prediction with server reconciliation.
- **Matchmaking** — `MatchmakingHost(port, roomName, maxPlayers)`, `MatchmakingDiscover(timeoutMs)`, `MatchmakingJoin(host, port)`. LAN broadcast discovery; returns table of rooms with host, port, roomName, playerCount.
- **Nakama (optional cloud)** — `NakamaConnect`, `NakamaAuthenticateDevice`, `NakamaCreateMatch`, `NakamaJoinMatch`, etc. Cloud backend for accounts, matchmaking, and realtime matches. See [NAKAMA_GUIDE.md](NAKAMA_GUIDE.md).
- **Interest management** — `SetInterestFilter(connectionId, "distance", maxDist, ox, oy, oz)`, `SetInterestFilter(connectionId, "zone", zoneId)`, `SetEntityInterestZone(entityId, zoneId)`. Filters `SyncEntity` and `SyncEntityToRoom` by distance or zone.

## Current Limitations

The following are not implemented yet:

- Snapshot interpolation buffers
- Reliable-ordered vs unreliable channels
- Built-in host migration
- Automatic high-level replication beyond explicit `SyncEntity` helper messages
- Automatic use of `ReplicatePosition` / `ReplicateRotation` / `ReplicateScale` / `ReplicateValue` markers

## Recommended Usage Today

CyberBASIC2 multiplayer is currently best suited for:

- Local network tools and prototypes
- Small co-op or client/server experiments
- Turn-based games
- Lightweight authoritative servers with modest player counts

For fast-action networked games, treat the current API as a foundation layer rather than a finished replication stack.
`SyncEntity` is real today; the `Replicate*` helpers are not a complete shipping replication system yet.

## Suggested Architecture

For projects shipping on the current engine:

1. Use one authority, usually the server.
2. Send compact input or intent messages from clients.
3. Run gameplay state changes inside `OnFixedUpdate`.
4. Broadcast coarse state with `SyncEntity` or custom text/JSON messages.
5. Reconcile visuals client-side if needed.

## Nakama (optional cloud)

When you need cloud-hosted accounts, matchmaking, or realtime matches:

1. **Connect:** `NakamaConnect(host, port, serverKey [, useSSL])`
2. **Authenticate:** `NakamaAuthenticateDevice(deviceId [, create, username])` or `NakamaAuthenticateEmail`, `NakamaAuthenticateCustom`
3. **Socket:** `NakamaCreateSocket()`, `NakamaSocketConnect()`
4. **Match:** `NakamaCreateMatch([name])`, `NakamaJoinMatch(matchId [, token])`, `NakamaSendMatchState(matchId, opCode, data [, reliable])`
5. **Matchmaking:** `NakamaAddMatchmaker([minPlayers, maxPlayers, query])`, `OnNakamaMatchmakerMatched(matchId, token)`
6. **Process:** Call `NakamaProcessEvents()` each frame to invoke `OnNakamaMatchData`, `OnNakamaMatchJoin`, `OnNakamaMatchLeave`, `OnNakamaMatchmakerMatched`

See [NAKAMA_GUIDE.md](NAKAMA_GUIDE.md) for full API and examples.

## Status Summary

Implemented now:

- Stable KCP (reliable UDP) connect/host/send/receive path
- Single-consumption message queue semantics
- Callback and polling APIs
- RPC and entity-sync helpers
- KCP transport with SetDeadline for timed accept
- Per-connection numeric receive buffers
- Fixed-step callback support for simulation

Still roadmap work:

- Higher-level replication and session services
- Automatic Replicate* integration with transport

---

## Lockstep

Fixed tick rate; server advances only when all client inputs for tick N are received.

- **Server:** `LockstepEnable(60)`; each frame call `ProcessNetworkEvents()`. When `OnLockstepTickReady(tickId)` fires, call `LockstepGetInputs(tickId)` to get `{connectionId: inputData}`. Run simulation, broadcast state.
- **Client:** `LockstepEnable(60)`; each tick call `LockstepSendInput(tickId, inputData)` to send input. When `OnLockstepTickReady(tickId)` fires (after server broadcasts `T\t<tickId>`), advance simulation.
- **Protocol:** `L\t<tickId>\t<data>` for input; `T\t<tickId>` for tick-ready broadcast.

## Rollback and Prediction

- **Snapshot:** Register `RegisterSnapshotHandler(subName)` and `RegisterRestoreHandler(subName)`. The snapshot sub receives `(tickId)` and must call `SnapshotStoreResult(tickId, jsonOrString)` before returning. The restore sub receives `(tickId, data)`.
- **Rollback:** `SnapshotCreate(tickId)` invokes the snapshot handler; `SnapshotRestore(tickId)` invokes the restore handler. Server can `RollbackBroadcast(tickId, correctTickId)`; clients receive `OnRollbackRequired(tickId, correctTickId)`.
- **Prediction:** `PredictionEnable()`; `PredictionStoreInput(tickId, input)` when sending input. When server state arrives, call `PredictionReconcile(tickId, stateJson)` to restore and re-simulate; `OnPredictionCorrected(tickId)` fires.

## Matchmaking

- **Host:** `MatchmakingHost(port, roomName, maxPlayers)` — starts KCP server and UDP broadcast every 1s on port 47777. Broadcast format: `CB_ROOM\t<port>\t<roomName>\t<count>/<max>`.
- **Discover:** `MatchmakingDiscover(timeoutMs)` — listens for broadcasts, returns table with `count`, `"0"`, `"1"`, ... where each entry has `host`, `port`, `roomName`, `playerCount`.
- **Join:** `MatchmakingJoin(host, port)` — alias for `Connect(host, port)`.

## Interest Management

- **SetInterestFilter(connectionId, "distance", maxDist, originX, originY, originZ)** — only send to this connection when entity is within `maxDist` of origin.
- **SetInterestFilter(connectionId, "zone", zoneId)** — only send when entity's zone (from `SetEntityInterestZone`) matches.
- **SetInterestFilter(connectionId, "all", "")** — disable filter.
- **SetEntityInterestZone(entityId, zoneId)** — assign entity to zone.
- `SyncEntity` and `SyncEntityToRoom` check each connection's filter before sending.

---

## Determinism Guidance

For deterministic or semi-deterministic multiplayer:

1. **Fixed step:** Use `FixedUpdate(rate)` and `OnFixedUpdate(label$)`. Run physics and gameplay state changes there.
2. **Input over state:** Clients send input (keys, commands); server runs simulation and broadcasts state. Avoid sending raw floats for positions every frame if you need determinism.
3. **Same order:** Ensure server and clients process inputs in the same order. ProcessNetworkEvents processes queued messages in order.
4. **No frame-dependent logic in sim:** Avoid `GetFrameTime()` or frame count in OnFixedUpdate. Use `FixedDeltaTime()` for step size.

---

## Contributor Notes

- **Net package:** `compiler/bindings/net/net.go` — RegisterNet, Host, Connect, Send, Receive, ProcessNetworkEvents
- **RPC:** `RegisterRPC(name, handler)`; handlers invoked when ProcessNetworkEvents runs
- **SyncEntity:** Sends entity position; receiver gets OnEntitySync(entityId, x, y, z)
- **TLS:** HostTLS, ConnectTLS deprecated (aliased to Host, Connect; use KCP)
- **Testing:** Use loopback (127.0.0.1) for local tests; no mock transport in tests yet

---

## See Also

- [Multiplayer](MULTIPLAYER.md) — Full API and examples
- [Tutorial Multiplayer](TUTORIAL_MULTIPLAYER.md)
- [Rendering and the Game Loop](RENDERING_AND_GAME_LOOP.md) — Fixed-step integration
- [Documentation Index](DOCUMENTATION_INDEX.md)
