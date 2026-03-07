# Multiplayer Design

Current multiplayer architecture, implementation status, deterministic patterns, and roadmap gaps.

---

## Purpose

- **TCP transport:** Connect, Host, Send, Receive for client/server games.
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

## Current Transport

CyberBASIC2 currently ships a TCP line-based networking layer:

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

## Current Limitations

The following features are not implemented yet:

- Deterministic lockstep scheduling
- Rollback save/restore and resimulation
- Snapshot interpolation buffers
- Interest management / area-of-relevance filtering
- Reliable-ordered vs unreliable channels
- Matchmaking, lobby discovery, NAT traversal
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

## Status Summary

Implemented now:

- Stable TCP connect/host/send/receive path
- Single-consumption message queue semantics
- Callback and polling APIs
- RPC and entity-sync helpers
- TLS-capable timed accept path
- Per-connection numeric receive buffers
- Fixed-step callback support for simulation

Still roadmap work:

- True deterministic multiplayer model
- Rollback / prediction toolchain
- Higher-level replication and session services

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
- **TLS:** HostTLS, ConnectTLS for encrypted transport
- **Testing:** Use loopback (127.0.0.1) for local tests; no mock transport in tests yet

---

## See Also

- [Multiplayer](MULTIPLAYER.md) — Full API and examples
- [Tutorial Multiplayer](TUTORIAL_MULTIPLAYER.md)
- [Rendering and the Game Loop](RENDERING_AND_GAME_LOOP.md) — Fixed-step integration
- [Documentation Index](DOCUMENTATION_INDEX.md)
