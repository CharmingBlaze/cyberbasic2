# Windows, Scaling, and Splitscreen

Guide to window commands, DPI/scaling, and view-based splitscreen (view-to-view) in CyberBasic.

---

## Window commands

Use these to control the main window and display. Call **SetConfigFlags** *before* **InitWindow** if you use flags.

| Command | Description |
|--------|-------------|
| **InitWindow**(width, height, title) | Create the window |
| **CloseWindow**() | Close and unload |
| **SetWindowPosition**(x, y) | Move window on screen |
| **GetWindowPosition**() | Returns [x, y] |
| **SetWindowSize**(width, height) | Resize window |
| **SetWindowTitle**(title) | Set title bar text |
| **SetWindowMinSize**(width, height) | Minimum size when resizable |
| **SetWindowMaxSize**(width, height) | Maximum size when resizable |
| **SetWindowOpacity**(opacity) | 0.0–1.0 |
| **SetWindowState**(flags) | Set state (see FLAG_* below) |
| **ClearWindowState**(flags) | Clear state flags |
| **SetWindowMonitor**(monitor) | Move to monitor index |
| **GetWindowScaleDPI**() | Returns [scaleX, scaleY] for DPI |
| **GetScaleDPI**() | Returns single scale factor (average of X/Y) for UI |
| **MaximizeWindow**(), **MinimizeWindow**(), **RestoreWindow**() | Window state |
| **IsWindowReady**(), **IsWindowFullscreen**(), **IsWindowHidden**(), **IsWindowMinimized**(), **IsWindowMaximized**(), **IsWindowFocused**(), **IsWindowResized**(), **IsWindowState**(flag) | Queries |
| **ToggleFullscreen**(), **ToggleBorderlessWindowed**() | Toggle modes |

**Screen and monitor:**

| Command | Description |
|--------|-------------|
| **GetScreenWidth**(), **GetScreenHeight**() | Logical window size (for layout) |
| **GetRenderWidth**(), **GetRenderHeight**() | Framebuffer size (for render textures, HiDPI) |
| **GetMonitorCount**(), **GetCurrentMonitor**() | Monitor index |
| **GetMonitorName**(monitor), **GetMonitorWidth**(monitor), **GetMonitorHeight**(monitor), **GetMonitorRefreshRate**(monitor), **GetMonitorPosition**(monitor) | Monitor info |

**Config flags (call SetConfigFlags before InitWindow):** Use **FLAG_VSYNC_HINT**(), **FLAG_WINDOW_RESIZABLE**(), **FLAG_WINDOW_UNDECORATED**(), **FLAG_WINDOW_HIDDEN**(), **FLAG_WINDOW_MINIMIZED**(), **FLAG_WINDOW_MAXIMIZED**(), **FLAG_WINDOW_UNFOCUSED**(), **FLAG_WINDOW_TOPMOST**(), **FLAG_WINDOW_ALWAYS_RUN**(), **FLAG_MSAA_4X_HINT**(), **FLAG_INTERLACED_HINT**(), **FLAG_WINDOW_HIGHDPI**(), **FLAG_BORDERLESS_WINDOWED_MODE**(). Combine with OR, e.g. `SetConfigFlags(FLAG_VSYNC_HINT() OR FLAG_WINDOW_RESIZABLE())`.

---

## Scaling and resolution

- **GetScreenWidth** / **GetScreenHeight** — Logical size of the window (use for layout and viewport splits).
- **GetRenderWidth** / **GetRenderHeight** — Actual framebuffer size; use when creating render textures so they match output resolution (important on HiDPI).
- **GetWindowScaleDPI**() — Returns [scaleX, scaleY] for DPI scaling. Use for UI: e.g. scale font size or layout by this factor so text stays readable on high-DPI displays.
- **GetScaleDPI**() — Single scale factor (average of X and Y) for simple UI scaling, e.g. `fontSize = 20 * GetScaleDPI()`.

**Recipe for DPI-scaled UI:** Call `scale = GetScaleDPI()` (or use the two components from GetWindowScaleDPI), then multiply your UI coordinates or font size by `scale` when drawing.

---

## View-to-view (splitscreen)

Raylib has **one** window. "Window-to-window" style is done with **views**: rectangles on screen, each showing a **render texture**. Communication between "windows" is via your variables and which render texture you assign to each view.

**Commands:**

- **CreateView**(viewId, x, y, width, height) — Define a view rectangle.
- **SetViewTarget**(viewId, renderTextureId) — This view will show the given render texture when you call DrawView.
- **DrawView**(viewId) — Draw the view’s render texture into its rectangle on screen.
- **GetViewX**(viewId), **GetViewY**(viewId), **GetViewWidth**(viewId), **GetViewHeight**(viewId) — Read back view rect (0 if missing).
- **SetViewPosition**(viewId, x, y), **SetViewSize**(viewId, width, height), **SetViewRect**(viewId, x, y, width, height) — Resize or move a view at runtime.

**Splitscreen helpers (use after InitWindow):**

- **CreateSplitscreenLeftRight**(viewIdLeft, viewIdRight) — Two views: left half and right half.
- **CreateSplitscreenTopBottom**(viewIdTop, viewIdBottom) — Two views: top half and bottom half.
- **CreateSplitscreenFour**(viewIdTL, viewIdTR, viewIdBL, viewIdBR) — Four quadrants.

You still create render textures (**LoadRenderTexture**) and assign them with **SetViewTarget**; the helpers only create the view rectangles.

---

## Splitscreen recipe (two players)

1. **InitWindow**, then **CreateSplitscreenLeftRight**("left", "right") (or TopBottom / manual CreateView).
2. Create two render textures, e.g. `rt1 = LoadRenderTexture(GetRenderWidth()/2, GetRenderHeight())`, `rt2 = LoadRenderTexture(GetRenderWidth()/2, GetRenderHeight())` for left/right.
3. Each frame:
   - **BeginTextureMode**(rt1), draw player 1’s scene (and own 2D camera if needed), **EndTextureMode**.
   - **BeginTextureMode**(rt2), draw player 2’s scene, **EndTextureMode**.
   - **SetViewTarget**("left", rt1), **SetViewTarget**("right", rt2).
   - **DrawView**("left"), **DrawView**("right").

Shared state (scores, etc.) lives in your variables; each view just shows the render texture you assign.

---

## See also

- [API_REFERENCE.md](../API_REFERENCE.md) — Full list of window, view, and texture functions.
- [2D_GRAPHICS_GUIDE.md](2D_GRAPHICS_GUIDE.md) — Drawing and cameras.
