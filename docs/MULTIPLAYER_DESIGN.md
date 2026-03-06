# Multiplayer Design

This document describes the current multiplayer architecture, what is implemented today, and what still belongs to future roadmap work.

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
- RPC and entity-sync payloads ride on the same message queue as normal text messages.

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

## Recommended Usage Today

CyberBASIC2 multiplayer is currently best suited for:

- Local network tools and prototypes
- Small co-op or client/server experiments
- Turn-based games
- Lightweight authoritative servers with modest player counts

For fast-action networked games, treat the current API as a foundation layer rather than a finished replication stack.

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
- Fixed-step callback support for simulation

Still roadmap work:

- True deterministic multiplayer model
- Rollback / prediction toolchain
- Higher-level replication and session services
