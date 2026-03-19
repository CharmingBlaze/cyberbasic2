# Nakama Guide

CyberBASIC2 includes optional Nakama client bindings for cloud-hosted multiplayer: accounts, matchmaking, realtime matches, and RPC.

## Prerequisites

- A Nakama server (self-hosted or [Heroic Labs cloud](https://heroiclabs.com))
- Default: `127.0.0.1:7350` with server key `defaultkey`

## Quick Start

```basic
NakamaConnect("127.0.0.1", 7350, "defaultkey", 0)
NakamaAuthenticateDevice("my-device-id", 1, "")
NakamaCreateSocket()
NakamaSocketConnect()
VAR matchId = NakamaCreateMatch("my-match")
IF NOT IsNull(matchId) THEN
  PRINT "Match created: " + matchId
  NakamaProcessEvents()
END IF
```

## API Reference

### Connection

| Command | Args | Description |
|---------|------|-------------|
| **NakamaConnect** | (host, port, serverKey [, useSSL]) | Create client. useSSL=1 for https/wss. |
| **NakamaCreateSocket** | () | Create realtime websocket. Call after authenticate. |
| **NakamaSocketConnect** | () | Open socket connection. |

### Authentication

| Command | Args | Description |
|---------|------|-------------|
| **NakamaAuthenticateDevice** | (deviceId [, create, username]) | Auth with device ID. create=1 creates account if new. |
| **NakamaAuthenticateCustom** | (customId [, create, username]) | Auth with custom ID. |
| **NakamaAuthenticateEmail** | (email, password [, create, username]) | Email/password auth. |

### Matches

| Command | Args | Description |
|---------|------|-------------|
| **NakamaCreateMatch** | ([name]) | Create match. Returns matchId or null. |
| **NakamaJoinMatch** | (matchId [, token]) | Join by matchId or token (from matchmaker). |
| **NakamaLeaveMatch** | (matchId) | Leave match. |
| **NakamaSendMatchState** | (matchId, opCode, data [, reliable]) | Send data to match. opCode is integer; data is string; reliable=1 default. |

### Matchmaking

| Command | Args | Description |
|---------|------|-------------|
| **NakamaAddMatchmaker** | ([minPlayers, maxPlayers, query]) | Join matchmaking pool. Returns ticket. |
| **NakamaRemoveMatchmaker** | (ticket) | Cancel matchmaking. |

### Account & RPC

| Command | Args | Description |
|---------|------|-------------|
| **NakamaGetAccount** | () | Get current account as JSON string. |
| **NakamaRPC** | (id, input) | Call server RPC. Returns response string. |

### Event Processing

| Command | Description |
|---------|-------------|
| **NakamaProcessEvents** | Drain event queue; invoke callbacks. Call once per frame. |

## Callbacks

Define these Subs (case-insensitive) and call **NakamaProcessEvents** each frame:

| Sub | Args | When |
|-----|------|------|
| **OnNakamaMatchData** | (matchId, opCode, data, sender) | Match data received. |
| **OnNakamaMatchJoin** | (matchId, presences) | Players joined match. |
| **OnNakamaMatchLeave** | (matchId, presences) | Players left match. |
| **OnNakamaMatchmakerMatched** | (matchId, token) | Matchmaking found a match. Use token with NakamaJoinMatch. |

## Example: Full Multiplayer Game

```basic
NakamaConnect("127.0.0.1", 7350, "defaultkey", 0)
NakamaAuthenticateDevice("player-" + GetLocalIP(), 1, "")
NakamaCreateSocket()
NakamaSocketConnect()

VAR ticket = NakamaAddMatchmaker(2, 4, "")
VAR matchId = ""
VAR token = ""

SUB OnNakamaMatchmakerMatched(mid, tok)
  matchId = mid
  token = tok
END SUB

mainloop
  NakamaProcessEvents()
  IF matchId <> "" AND token <> "" THEN
    NakamaJoinMatch(matchId, token)
    matchId = ""
    token = ""
  END IF
  // ... game update/draw ...
  SYNC
endmain
```

## Realtime vs Direct KCP

- **Direct KCP (Host/Connect):** LAN, peer-to-peer, no server. Full control.
- **Nakama:** Internet, accounts, matchmaking, leaderboards. Requires Nakama server.

You can use Nakama for matchmaking and then direct KCP for gameplay (hybrid), but that requires additional setup.
