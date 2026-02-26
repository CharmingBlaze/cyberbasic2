# In-process multi-window system

Multiple **logical windows** (viewports) in a single process. Each window is backed by a render texture and can have its own drawing, messages, channels, state, and event handlers. **Window ID 0** is the main screen (the actual raylib window).

For **multiple processes** (separate windows as separate processes talking over TCP), see [Multiple windows from one .bas](MULTI_WINDOW.md) instead.

## Table of Contents

1. [Overview](#overview)
2. [Window creation and lifecycle](#window-creation-and-lifecycle)
3. [Rendering](#rendering)
4. [Messages](#messages)
5. [Channels and state](#channels-and-state)
6. [Events](#events)
7. [3D and RPC](#3d-and-rpc)
8. [Docking](#docking)
9. [Example](#example)
10. [See also](#see-also)

---

## Overview

- **One process, one main window (ID 0).** Additional windows are logical viewports (render textures) drawn onto the main window or arranged via docking.
- Use **WindowCreate**(width, height, title) to create a normal window; **WindowCreatePopup**, **WindowCreateModal**, **WindowCreateToolWindow** for typed windows; **WindowCreateChild**(parentID, width, height, title) for child windows.
- **WindowClose**(id) closes and frees a window. **WindowIsOpen**(id) returns whether the window still exists.
- All names are **case-insensitive**.

---

## Window creation and lifecycle

| Function | Description |
|----------|-------------|
| **WindowCreate**(width, height, title) | Create normal window → window id (integer) |
| **WindowCreatePopup**(width, height, title) | Create popup window |
| **WindowCreateModal**(width, height, title) | Create modal window |
| **WindowCreateToolWindow**(width, height, title) | Create tool window |
| **WindowCreateChild**(parentID, width, height, title) | Create child window under parent |
| **WindowClose**(id) | Close and free window (no-op if already closed) |
| **WindowIsOpen**(id) | True if window exists |
| **WindowSetTitle**(id, title) | Set window title |
| **WindowSetSize**(id, width, height) | Resize window (recreates render texture) |
| **WindowSetPosition**(id, x, y) | Set position (for drawing onto main screen) |
| **WindowGetWidth**(id), **WindowGetHeight**(id) | Get size |
| **WindowGetPositionX**(id), **WindowGetPositionY**(id) | Get position |
| **WindowGetPosition**(id) | Returns "x,y" string |
| **WindowFocus**(id) | Set this window as focused |
| **WindowIsFocused**(id) | True if this window has focus |
| **WindowIsVisible**(id) | True if visible |
| **WindowShow**(id), **WindowHide**(id) | Show or hide window |

---

## Rendering

Drawing into a window: **WindowBeginDrawing**(id), draw with normal raylib 2D/3D/GUI calls, then **WindowEndDrawing**(id). For the main screen use id **0**. To composite all windows onto the main screen, call **WindowDrawAllToScreen**() (draws all visible windows’ render textures at their positions).

| Function | Description |
|----------|-------------|
| **WindowBeginDrawing**(id) | Start drawing to window (0 = main screen) |
| **WindowEndDrawing**(id) | End drawing to window |
| **WindowClearBackground**(id, r, g, b, a) | Clear window background |
| **WindowDrawAllToScreen**() | Draw all visible windows onto the main screen |

---

## Messages

Per-window message queue. Send to one window or broadcast to all.

| Function | Description |
|----------|-------------|
| **WindowSendMessage**(targetID, message, data) | Send (message, data) to one window |
| **WindowBroadcast**(message, data) | Send to all windows |
| **WindowReceiveMessage**(id) | Receive next message for window → "message|data" string or null |
| **WindowHasMessage**(id) | True if window has pending messages |

---

## Channels and state

**Channels** are named queues shared across windows (useful for producer/consumer). **State** is a key-value store shared across windows.

| Function | Description |
|----------|-------------|
| **ChannelCreate**(name) | Ensure channel exists |
| **ChannelSend**(name, data) | Append to channel |
| **ChannelReceive**(name) | Remove and return next value (or null) |
| **ChannelHasData**(name) | True if channel has data |
| **StateSet**(key, value) | Set shared state |
| **StateGet**(key) | Get value (or null) |
| **StateHas**(key) | True if key exists |
| **StateRemove**(key) | Remove key |

---

## Events

Register handler **Sub** names for window events. Then call **WindowProcessEvents**() to run update handlers and **WindowDraw**() to run draw handlers (or call your own loop).

| Function | Description |
|----------|-------------|
| **OnWindowUpdate**(id, subName) | Call Sub when window is updated |
| **OnWindowDraw**(id, subName) | Call Sub when window is drawn |
| **OnWindowResize**(id, subName) | Call Sub on resize |
| **OnWindowClose**(id, subName) | Call Sub when window is closed |
| **OnWindowMessage**(id, subName) | Call Sub when message received |
| **WindowProcessEvents**() | Invoke registered update handlers |
| **WindowDraw**() | Invoke registered draw handlers (e.g. BeginDrawing, draw, EndDrawing per window) |

---

## 3D and RPC

For 3D in a logical window: **WindowSetCamera**(id, cameraId) to use a camera for that window; **WindowDrawModel**(id, modelId, x, y, z, scale [, r, g, b, a]), **WindowDrawScene**(id, sceneId) to draw into that window.

**RPC:** Register a Sub to be called from another window: **WindowRegisterFunction**(windowId, name, subName); **WindowCall**(targetWindowId, name, arg1, arg2, …) to invoke it.

| Function | Description |
|----------|-------------|
| **WindowSetCamera**(id, cameraId) | Set 3D camera for window |
| **WindowDrawModel**(id, modelId, x, y, z, scale [, r, g, b, a]) | Draw model into window |
| **WindowDrawScene**(id, sceneId) | Draw scene into window |
| **WindowRegisterFunction**(windowId, name, subName) | Register Sub as callable by name |
| **WindowCall**(targetWindowId, name, arg1, arg2, …) | Call registered function on target window |

---

## Docking

Optional layout: create dock areas, split them, and attach windows. Useful for editor-style UIs.

| Function | Description |
|----------|-------------|
| **DockCreateArea**(id, x, y, width, height) | Create dock area |
| **DockSplit**(areaId, direction, size) | Split area (direction e.g. "horizontal", "vertical") |
| **DockAttachWindow**(areaId, windowId) | Attach window to area |
| **DockSetSize**(nodeId, size) | Set split size |

---

## Example

See **examples/multi_window_gui_demo.bas** for a full example: two panels, **StateSet**/ **StateGet**, **OnWindowDraw**, **WindowProcessEvents**, **WindowBeginDrawing**/ **WindowEndDrawing**, **WindowDrawAllToScreen**.

---

## See also

- [Multiple windows from one .bas (multi-process)](MULTI_WINDOW.md)
- [API Reference](../API_REFERENCE.md) – full binding list
- [Documentation Index](DOCUMENTATION_INDEX.md)
- [Windows, scaling, and splitscreen](WINDOWS_AND_VIEWS.md)
