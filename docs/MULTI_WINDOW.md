# Multiple windows from one .bas

Run multiple windows from a single .bas file. Each window is a separate process with its own raylib window; they communicate over loopback TCP using the existing **NET** API (Send / Receive).

---

## How it works

- **Main process**: Your script runs from the top. It calls `Host(port)`, `InitWindow(...)`, then **SpawnWindow**(port, title, width, height) to start a child process. The child connects back; the main script gets the connection with **AcceptTimeout**(serverId, timeout) and then uses **Send**(connectionId, text) / **Receive**(connectionId) as usual.

- **Window process**: The same .bas is run again with flags `--window --parent=127.0.0.1:port --title=... --width=... --height=...`. The script detects this with **IsWindowProcess()** and runs a short block: **ConnectToParent()**, **InitWindow**(GetWindowWidth(), GetWindowHeight(), GetWindowTitle()), then its own loop with Receive / Send, then CloseWindow and END.

No new wire protocol: use **Send** and **Receive** (and **SendNumbers** / **ReceiveNumbers** etc.) for “easy talk” between main and child windows.

---

## API

| Function | Where | Purpose |
|----------|--------|---------|
| **GetEnv**(key) | std | Read env var (e.g. CYBERBASIC_*) |
| **IsWindowProcess()** | std | True if this process is a spawned window |
| **GetWindowTitle()** | std | Window title for child (from --title=) |
| **GetWindowWidth()** | std | Width for child (from --width=); default 400 |
| **GetWindowHeight()** | std | Height for child (from --height=); default 300 |
| **SpawnWindow**(port, title, width, height) | std | Start same .bas as child window; returns 1 on success, 0 on failure |
| **ConnectToParent()** | net | Connect to parent; returns connection id or null |

---

## Script pattern (one .bas, two windows)

Put the window-process block at the top so the child exits after its window loop.

```basic
IF IsWindowProcess() THEN
  VAR cid = ConnectToParent()
  IF IsNull(cid) THEN QUIT
  ENDIF
  InitWindow(GetWindowWidth(), GetWindowHeight(), GetWindowTitle())
  WHILE NOT WindowShouldClose()
    VAR msg = Receive(cid)
    Send(cid, "ok")
  WEND
  CloseWindow()
  QUIT
ENDIF

VAR sid = Host(9999)
IF IsNull(sid) THEN QUIT
ENDIF
InitWindow(800, 600, "Main")
SpawnWindow(9999, "Child", 400, 300)
VAR cid = AcceptTimeout(sid, 5000)
IF IsNull(cid) THEN
  CloseServer(sid)
  QUIT
ENDIF
WHILE NOT WindowShouldClose()
  Send(cid, "update")
  VAR msg = Receive(cid)
WEND
CloseServer(sid)
CloseWindow()
```

- Main uses **SpawnWindow**(port, title, width, height) then **AcceptTimeout**(serverId, timeout) to get the connection.
- More windows: call **SpawnWindow** again and **AcceptTimeout** again for each; store connection ids in variables or an array.

---

## See also

- [Windows, scaling, and splitscreen](WINDOWS_AND_VIEWS.md) – Single-window views and splitscreen
- [Multiplayer (TCP)](MULTIPLAYER.md) – Connect, Host, Send, Receive reference
