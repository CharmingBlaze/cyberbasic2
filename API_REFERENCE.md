# CyberBasic API Reference – All Bindings

All functions callable from BASIC. Names are **case-insensitive**. You can call with or without namespace (e.g. `InitWindow(...)` or `RL.InitWindow(...)` for raylib; use `BOX2D.*` and `BULLET.*` for physics).

[COMMAND_REFERENCE.md](docs/COMMAND_REFERENCE.md) groups commands by feature (Window, Input, Math, Camera, 2D, 3D, etc.) for task-based lookup. This document lists **all bindings by source file**. Each section below uses a table for quick lookup.

**Maintaining this document:** When adding a command, (1) register it in the appropriate file under [compiler/bindings/](compiler/bindings/); (2) add one row to the corresponding section table (Command, Arguments, Returns, Description).

---

## Table of Contents

- [1. Raylib (core)](#1-raylib-core--raylib_corego)
- [2. Raylib (input)](#2-raylib-input--raylib_inputgo)
- [3. Raylib (shapes)](#3-raylib-shapes--raylib_shapesgo)
- [4. Raylib (text)](#4-raylib-text--raylib_textgo)
- [5. Raylib (textures)](#5-raylib-textures--raylib_texturesgo)
- [6. Raylib (images)](#6-raylib-images--raylib_imagesgo)
- [7. Raylib (3D)](#7-raylib-3d--raylib_3dgo)
- [8. Raylib (mesh)](#8-raylib-mesh--raylib_meshgo)
- [9. Raylib (audio)](#9-raylib-audio--raylib_audiogo)
- [10. Raylib (fonts)](#10-raylib-fonts--raylib_fontsgo)
- [11. Raylib (misc)](#11-raylib-misc--raylib_miscgo)
- [12. Raylib (math)](#12-raylib-math--raylib_mathgo)
- [13. Raylib (game)](#13-raylib-game--raylib_gamego)
- [14. Box2D](#14-box2d--box2dgo)
- [15. Bullet (3D physics)](#15-bullet-3d-physics--bulletgo)
- [16. ECS](#16-ecs--ecsgo)
- [17. Std](#17-std-file-string-math-json-enum-dictionary-http-help-multi-window--stdgo)
- [18. Multiplayer (TCP)](#18-multiplayer-tcp--netgo)
- [19. SQL](#19-sql--sqlgo)
- [20. UI](#20-ui--raylib_uigo-and-full-raygui--raylib_rayguigo)
- [21. Language and built-ins](#21-language-and-built-ins)
- [22. Multi-window (in-process)](#22-multi-window-in-process--raylib_multiwindowgo)
- [Notes](#notes)

---

## 1. Raylib (core) – `raylib_core.go`

The compiler does not inject frame or mode calls; your code compiles as written. Exception: when you define **update(dt)** and **draw()** and use a game loop, the compiler injects the [hybrid loop](docs/PROGRAM_STRUCTURE.md#hybrid-updatedraw-loop).

### Window and app

| Command | Arguments | Returns | Description |
|---------|-----------|---------|-------------|
| **InitWindow** | (width, height, title) | — | Open game window |
| **CloseWindow** | () | — | Close window |
| **WindowShouldClose** | () | bool | True when user requested close |
| **SetTargetFPS** | (fps) | — | Target frames per second |
| **SetWindowPosition** | (x, y) | — | Set window position |
| **SetWindowSize** | (w, h) | — | Set window size |
| **SetWindowTitle** | (title) | — | Set window title |
| **SetWindowMinSize** | (w, h) | — | Set minimum window size |
| **SetWindowMaxSize** | (w, h) | — | Set maximum window size |
| **SetWindowOpacity** | (opacity) | — | Set window opacity (0–1) |
| **GetScreenWidth** | () | int | Screen width in pixels |
| **GetScreenHeight** | () | int | Screen height in pixels |
| **GetRenderWidth** | () | int | Render width |
| **GetRenderHeight** | () | int | Render height |
| **GetWindowPosition** | () | x, y | Window position |
| **GetWindowScaleDPI** | () | float | Window DPI scale |
| **GetScaleDPI** | () | float | Single DPI scale for UI |
| **MaximizeWindow** | () | — | Maximize window |
| **MinimizeWindow** | () | — | Minimize window |
| **RestoreWindow** | () | — | Restore window |
| **ToggleFullscreen** | () | — | Toggle fullscreen |
| **ToggleBorderlessWindowed** | () | — | Toggle borderless windowed |
| **IsWindowReady** | () | bool | True if window is ready |
| **IsWindowFullscreen** | () | bool | True if fullscreen |
| **IsWindowHidden** | () | bool | True if hidden |
| **IsWindowMinimized** | () | bool | True if minimized |
| **IsWindowMaximized** | () | bool | True if maximized |
| **IsWindowFocused** | () | bool | True if focused |
| **IsWindowResized** | () | bool | True if resized this frame |
| **IsWindowState** | (flag) | bool | True if window has flag |
| **SetWindowState** | (flag) | — | Set window flag |
| **ClearWindowState** | (flag) | — | Clear window flag |
| **SetWindowMonitor** | (monitor) | — | Set window to monitor |
| **SetConfigFlags** | (flags) | — | Set config flags (before InitWindow) |

### Time, FPS, random

| Command | Arguments | Returns | Description |
|---------|-----------|---------|-------------|
| **GetFrameTime** | () | float | Delta time since last frame (seconds) |
| **DeltaTime** | () | float | Same as GetFrameTime; preferred for frame delta |
| **GetFPS** | () | int | Current FPS |
| **GetTime** | () | double | Time in seconds (raylib) |
| **GetRandomValue** | (min, max) | int | Random int in [min, max] |
| **SetRandomSeed** | (seed) | — | Set random seed |
| **WaitTime** | (seconds) | — | Block for N seconds |
| **EnableEventWaiting** | () | — | Enable event waiting (no busy loop) |
| **DisableEventWaiting** | () | — | Disable event waiting |

### Frame, clear, 2D mode

| Command | Arguments | Returns | Description |
|---------|-----------|---------|-------------|
| **ClearBackground** | (r, g, b, a) | — | Clear screen with RGBA |
| **Background** | (r, g, b) | — | Clear with RGB (alpha 255) |
| **BeginFrame** | () | — | Start frame (alias BeginDrawing) |
| **EndFrame** | () | — | End frame (alias EndDrawing) |
| **SetCamera2D** | (camera2D) | — | Set 2D camera |
| **BeginMode2D** | (camera2D) | — | Begin 2D mode. In the hybrid loop (when using update/draw), the engine wraps 2D automatically; calling these in draw() has no effect. |
| **EndMode2D** | () | — | End 2D mode. In the hybrid loop, calling in draw() has no effect. |
| **Camera2DCreate** | () | cameraID | Create 2D camera by ID |
| **Camera2DSetPosition** | (cameraID, x, y) | — | Set camera target |
| **Camera2DSetZoom** | (cameraID, zoom) | — | Set zoom |
| **Camera2DSetRotation** | (cameraID, angle) | — | Set rotation (rad) |
| **Camera2DMove** | (cameraID, dx, dy) | — | Move target by offset |
| **Camera2DSmoothFollow** | (cameraID, targetX, targetY, speed) | — | Smooth follow (call each frame) |
| **BeginCamera2D** | (cameraID) | — | Set active 2D camera for flush |
| **EndCamera2D** | () | — | Clear active 2D camera |
| **GetWorldToScreen2D** | (pos, camera) | x, y | World to screen 2D |
| **GetScreenToWorld2D** | (pos, camera) | x, y | Screen to world 2D |
| **BeginBlendMode** | (mode) | — | Begin blend mode |
| **EndBlendMode** | () | — | End blend mode |
| **BeginScissorMode** | (x, y, w, h) | — | Begin scissor test |
| **EndScissorMode** | () | — | End scissor test |
| **BeginShaderMode** | (shaderId) | — | Begin shader |
| **EndShaderMode** | () | — | End shader |
| **LoadShader** | (vsPath, fsPath) | id | Load shader |
| **LoadShaderFromMemory** | (vsCode, fsCode) | id | Load shader from strings |
| **UnloadShader** | (id) | — | Unload shader |
| **IsShaderValid** | (id) | bool | True if shader valid |
| **SwapScreenBuffer** | () | — | Swap buffers |
| **PollInputEvents** | () | — | Poll input (call each frame) |

### Monitors, clipboard, cursor, misc

| Command | Arguments | Returns | Description |
|---------|-----------|---------|-------------|
| **GetMonitorCount** | () | int | Number of monitors |
| **GetCurrentMonitor** | () | int | Current monitor index |
| **GetMonitorName** | (index) | string | Monitor name |
| **GetMonitorWidth** | (index) | int | Monitor width |
| **GetMonitorHeight** | (index) | int | Monitor height |
| **GetMonitorRefreshRate** | (index) | int | Refresh rate (Hz) |
| **GetMonitorPosition** | (index) | x, y | Monitor position |
| **GetMonitorPhysicalWidth** | (index) | int | Physical width (mm) |
| **GetMonitorPhysicalHeight** | (index) | int | Physical height (mm) |
| **GetClipboardText** | () | string | Clipboard text |
| **SetClipboardText** | (text) | — | Set clipboard text |
| **TakeScreenshot** | (path) | — | Save screenshot to file |
| **OpenURL** | (url) | — | Open URL in browser |
| **IsCursorHidden** | () | bool | True if cursor hidden |
| **EnableCursor** | () | — | Show cursor |
| **DisableCursor** | () | — | Hide cursor |
| **IsCursorOnScreen** | () | bool | True if cursor on screen |
| **FileExists** | (path) | bool | True if file exists |

**Config/blend flags (0-arg):** FLAG_VSYNC_HINT, FLAG_FULLSCREEN_MODE, FLAG_WINDOW_RESIZABLE, FLAG_WINDOW_UNDECORATED, FLAG_WINDOW_HIDDEN, FLAG_WINDOW_MINIMIZED, FLAG_WINDOW_MAXIMIZED, FLAG_WINDOW_UNFOCUSED, FLAG_WINDOW_TOPMOST, FLAG_WINDOW_ALWAYS_RUN, FLAG_MSAA_4X_HINT, FLAG_INTERLACED_HINT, FLAG_WINDOW_HIGHDPI, FLAG_BORDERLESS_WINDOWED_MODE; BLEND_ALPHA, BLEND_ADDITIVE, BLEND_MULTIPLIED, BLEND_ADD_COLORS, BLEND_SUBTRACT_COLORS, BLEND_CUSTOM. See [Windows, scaling, and splitscreen](docs/WINDOWS_AND_VIEWS.md).

### Layer system (raylib_layers.go)

| Command | Arguments | Returns | Description |
|---------|-----------|---------|-------------|
| **LayerCreate** | (name, order) | layerID | Create layer (order = draw priority) |
| **LayerSetOrder** | (layerID, order) | — | Change draw order |
| **LayerSetVisible** | (layerID, flag) | — | Hide (0) or show layer |
| **LayerSetParallax** | (layerID, parallaxX, parallaxY) | — | Parallax factors |
| **LayerSetScroll** | (layerID, scrollX, scrollY) | — | Scroll offset |
| **LayerClear** | (layerID) | — | Clear all drawables from layer |
| **LayerSortSprites** | (layerID) | — | No-op (flush sorts by z) |
| **SpriteSetLayer** | (spriteID, layerID) | — | Assign sprite to layer (raylib_sprite.go) |
| **SpriteSetZIndex** | (spriteID, z) | — | Z-order within layer |
| **TilemapSetLayer** | (tilemapID, layerID) | — | Assign tilemap to layer (game) |
| **TilemapSetParallax** | (tilemapID, px, py) | — | Parallax for tilemap |
| **ParticleSetLayer** | (particleID, layerID) | — | Assign particle system to layer (game) |

### Background system (raylib_background.go)

| Command | Arguments | Returns | Description |
|---------|-----------|---------|-------------|
| **BackgroundCreate** | (textureId) | backgroundId | Create background |
| **BackgroundSetColor** | (backgroundId, r, g, b, a) | — | Tint |
| **BackgroundSetTexture** | (backgroundId, textureId) | — | Set texture |
| **BackgroundSetScroll** | (backgroundId, speedX, speedY) | — | Scroll speed |
| **BackgroundSetOffset** | (backgroundId, offsetX, offsetY) | — | Offset |
| **BackgroundSetParallax** | (backgroundId, px, py) | — | Parallax |
| **BackgroundSetTiled** | (backgroundId, flag) | — | Enable tiling |
| **BackgroundSetTileSize** | (backgroundId, width, height) | — | Tile size |
| **BackgroundAddLayer** | (backgroundId, textureId, px, py) | — | Add layer |
| **BackgroundRemoveLayer** | (backgroundId, layerIndex) | — | Remove layer |
| **DrawBackground** | (backgroundId) | — | Draw (queued in 2D) |

### Hybrid loop (raylib_hybrid.go)

| Command | Arguments | Returns | Description |
|---------|-----------|---------|-------------|
| **ClearRenderQueues** | () | — | Clear 2D/3D/GUI render queues |
| **FlushRenderQueues** | () | — | Execute queued draw commands and present |
| **StepAllPhysics2D** | (dt) | — | Step all Box2D worlds |
| **StepAllPhysics3D** | (dt) | — | Step all Bullet worlds |
| **rect** | (…) | — | Alias of DrawRectangle |
| **circle** | (…) | — | Alias of DrawCircle |
| **cube** | (…) | — | Alias of DrawCube |
| **button** | (…) | — | Alias of GuiButton |
| **sprite** | (…) | — | Alias of DrawTexture |

When **update(dt)** and/or **draw()** are defined and the main loop is a game loop, the compiler invokes them automatically (GetFrameTime → physics step → update(dt) → ClearRenderQueues → draw() → FlushRenderQueues). See [Program Structure](docs/PROGRAM_STRUCTURE.md#hybrid-updatedraw-loop).

---

## 2. Raylib (input) – `raylib_input.go`

| Command | Arguments | Returns | Description |
|---------|-----------|---------|-------------|
| **IsKeyPressed** | (key) | bool | True once when key pressed |
| **IsKeyDown** | (key) | bool | True while key held |
| **KeyPressed** | (key) | bool | Alias of IsKeyPressed |
| **KeyDown** | (key) | bool | Alias of IsKeyDown |
| **IsKeyReleased** | (key) | bool | True once when key released |
| **IsKeyUp** | (key) | bool | True when key not pressed |
| **IsKeyPressedRepeat** | (key) | bool | True when key repeat (hold) |
| **GetKeyPressed** | () | int | Last key pressed (code) |
| **GetCharPressed** | () | int | Last char pressed (unicode) |
| **SetExitKey** | (key) | — | Key that triggers WindowShouldClose |
| **GetMouseX** | () | int | Mouse X position |
| **GetMouseY** | () | int | Mouse Y position |
| **GetMousePosition** | () | x, y | Mouse position as vector |
| **GetMouseDeltaX** | () | float | Mouse movement X this frame |
| **GetMouseDeltaY** | () | float | Mouse movement Y this frame |
| **GetMouseWheelMove** | () | float | Scroll wheel delta |
| **GetMouseWheelMoveV** | () | x, y | Scroll wheel (x, y) |
| **IsMouseButtonPressed** | (button) | bool | True once when button pressed |
| **IsMouseButtonDown** | (button) | bool | True while button held |
| **IsMouseButtonReleased** | (button) | bool | True once when released |
| **IsMouseButtonUp** | (button) | bool | True when button not pressed |
| **SetMousePosition** | (x, y) | — | Set mouse position |
| **SetMouseOffset** | (x, y) | — | Set mouse offset |
| **SetMouseScale** | (scaleX, scaleY) | — | Set mouse scale |
| **HideCursor** | () | — | Hide cursor |
| **ShowCursor** | () | — | Show cursor |
| **SetMouseCursor** | (cursor) | — | Set cursor shape |
| **GetVector2X** | (v) | float | X component of Vector2 |
| **GetVector2Y** | (v) | float | Y component of Vector2 |
| **GetVector3Z** | (v) | float | Z component of Vector3 |
| **IsGamepadAvailable** | (gamepad) | bool | True if gamepad connected |
| **GetGamepadName** | (gamepad) | string | Gamepad name |
| **IsGamepadButtonPressed** | (gamepad, button) | bool | True once when pressed |
| **IsGamepadButtonDown** | (gamepad, button) | bool | True while held |
| **IsGamepadButtonReleased** | (gamepad, button) | bool | True once when released |
| **IsGamepadButtonUp** | (gamepad, button) | bool | True when not pressed |
| **GetGamepadButtonPressed** | (gamepad) | int | Last button pressed |
| **GetGamepadAxisMovement** | (gamepad, axis) | float | Axis value (-1 to 1) |
| **GetGamepadAxisCount** | (gamepad) | int | Number of axes |
| **SetGamepadMappings** | (mappings) | int | Set gamepad mapping string |
| **SetGamepadVibration** | (gamepad, left, right) | — | Set vibration (0–1) |
| **GetTouchPointCount** | () | int | Touch point count |
| **GetTouchX** | (index) | int | Touch X |
| **GetTouchY** | (index) | int | Touch Y |
| **GetTouchPosition** | (index) | x, y | Touch position |
| **GetTouchPointId** | (index) | int | Touch point ID |

**Key constants (0-arg):** KEY_NULL, KEY_APOSTROPHE, KEY_COMMA, KEY_MINUS, KEY_PERIOD, KEY_SLASH, KEY_ZERO … KEY_NINE, KEY_SEMICOLON, KEY_EQUAL, KEY_A … KEY_Z, KEY_LEFT_BRACKET, KEY_BACKSLASH, KEY_RIGHT_BRACKET, KEY_GRAVE, KEY_SPACE, KEY_ESCAPE, KEY_ENTER, KEY_TAB, KEY_BACKSPACE, KEY_INSERT, KEY_DELETE, KEY_RIGHT, KEY_LEFT, KEY_DOWN, KEY_UP, KEY_PAGE_UP, KEY_PAGE_DOWN, KEY_HOME, KEY_END, KEY_F1 … KEY_F12.

**Movement:** **GetAxisX()** / **GetAxisY()** return -1, 0, or 1 for A/D and W/S. For full 2D/3D use **GAME.MoveWASD**, **MoveHorizontal2D**, **Jump2D** (see §13 Raylib (game)).

---

## 3. Raylib (shapes) – `raylib_shapes.go`

| Command | Arguments | Returns | Description |
|---------|-----------|---------|-------------|
| **SetShapesTexture** | (textureId, source) | — | Set texture for shapes |
| **GetShapesTextureRectangle** | () | rect | Get shapes texture rect |
| **DrawRectangle** | (x, y, w, h, r, g, b, a) | — | Filled rectangle |
| **DrawRectangleV** | (pos, size, color) | — | Filled rectangle (vector) |
| **DrawRectangleRec** | (rec, color) | — | Filled rectangle (rec) |
| **DrawRectanglePro** | (rec, origin, rotation, color) | — | Filled rectangle (rotated) |
| **DrawRectangleLines** | (x, y, w, h, color) | — | Rectangle outline |
| **DrawRectangleLinesEx** | (rec, thick, color) | — | Rectangle outline (thickness) |
| **DrawRectangleRounded** | (rec, roundness, segments, color) | — | Rounded filled rectangle |
| **DrawRectangleRoundedLines** | (rec, roundness, segments, thick, color) | — | Rounded rectangle outline |
| **DrawCircle** | (x, y, radius, r, g, b, a) | — | Filled circle |
| **DrawCircleV** | (center, radius, color) | — | Filled circle (vector) |
| **DrawCircleLines** | (x, y, radius, color) | — | Circle outline |
| **DrawCircleSector** | (center, radius, startAngle, endAngle, segments, color) | — | Filled circle sector |
| **DrawCircleGradient** | (x, y, radius, color1, color2) | — | Gradient circle |
| **DrawCircleLinesV** | (center, radius, color) | — | Circle outline (vector) |
| **DrawEllipse** | (x, y, radiusH, radiusV, color) | — | Filled ellipse |
| **DrawEllipseLines** | (x, y, radiusH, radiusV, color) | — | Ellipse outline |
| **DrawRing** | (center, innerRadius, outerRadius, startAngle, endAngle, segments, color) | — | Filled ring |
| **DrawRingLines** | (center, innerRadius, outerRadius, startAngle, endAngle, segments, color) | — | Ring outline |
| **DrawLine** | (x1, y1, x2, y2, color) | — | Line (2D) |
| **DrawLineV** | (start, end, color) | — | Line (vectors) |
| **DrawLineEx** | (start, end, thick, color) | — | Line with thickness |
| **DrawTriangle** | (v1, v2, v3, color) | — | Filled triangle |
| **DrawTriangleLines** | (v1, v2, v3, color) | — | Triangle outline |
| **DrawPoly** | (center, sides, radius, rotation, color) | — | Filled polygon |
| **DrawPolyLines** | (center, sides, radius, rotation, color) | — | Polygon outline |
| **DrawPixel** | (x, y, color) | — | Single pixel |
| **DrawPixelV** | (pos, color) | — | Single pixel (vector) |
| **DrawGrid** | (slices, spacing) | — | Draw 3D grid |
| **DrawFPS** | (x, y) | — | Draw FPS counter |

---

## 4. Raylib (text) – `raylib_text.go`

| Command | Arguments | Returns | Description |
|---------|-----------|---------|-------------|
| **DrawText** | (text, x, y, size, r, g, b, a) | — | Draw text at (x,y) with size and color |
| **DrawTextSimple** | (text, x, y) | — | Draw at (x,y), size 20, white (on-screen; use PRINT for console) |
| **DrawTextEx** | (fontId, text, pos, size, spacing, tint) | — | Draw text with font |
| **DrawTextPro** | (fontId, text, pos, origin, rotation, size, spacing, tint) | — | Draw text (rotated) |
| **MeasureText** | (text, size) | int | Text width in pixels |
| **MeasureTextEx** | (fontId, text, size, spacing) | width, height | Text size |
| **SetTextLineSpacing** | (spacing) | — | Line spacing |
| **TextCopy** | (dst, src) | int | Copy text |
| **TextIsEqual** | (text1, text2) | bool | True if equal |
| **TextLength** | (text) | int | Length in bytes |
| **TextFormat** | (format, …) | string | Formatted string |
| **TextSubtext** | (text, offset, length) | string | Substring |
| **TextReplace** | (text, old, new) | string | Replace in string |
| **TextInsert** | (text, insert, position) | string | Insert at position |
| **TextJoin** | (textList, count, delimiter) | string | Join strings |
| **TextSplit** | (text, delimiter) | count | Split; use GetTextSplitItem |
| **GetTextSplitItem** | (index) | string | Item from last TextSplit |
| **TextAppend** | (text, append) | — | Append in place |
| **TextFindIndex** | (text, find) | int | Index of find in text |
| **TextToUpper** | (text) | string | Uppercase |
| **TextToLower** | (text) | string | Lowercase |
| **TextToPascal** | (text) | string | PascalCase |
| **TextToSnake** | (text) | string | snake_case |
| **TextToCamel** | (text) | string | camelCase |
| **TextToInteger** | (text) | int | Parse integer |
| **TextToFloat** | (text) | float | Parse float |
| **GetCodepointCount** | (text) | int | Codepoint count |
| **GetCodepoint** | (text, index) | int | Codepoint at index |
| **GetCodepointNext** | (text, index) | int | Next codepoint |
| **GetCodepointPrevious** | (text, index) | int | Previous codepoint |
| **CodepointToUTF8** | (codepoint) | string | Codepoint to UTF-8 |
| **LoadCodepoints** | (text) | count | Load codepoints; use GetLoadedCodepoint |
| **UnloadCodepoints** | (codepoints) | — | Unload |
| **GetLoadedCodepoint** | (index) | int | Codepoint at index |
| **LoadUTF8** | (text) | count | Load UTF-8; use GetLoadedCodepoint |
| **UnloadUTF8** | (codepoints) | — | Unload |

---

## 5. Raylib (textures) – `raylib_textures.go`

| Command | Arguments | Returns | Description |
|---------|-----------|---------|-------------|
| **LoadTexture** | (path) | id | Load texture from file |
| **UnloadTexture** | (id) | — | Unload texture |
| **GetTextureWidth** | (textureId) | int | Texture width in pixels |
| **GetTextureHeight** | (textureId) | int | Texture height in pixels |
| **GetTextureSize** | (textureId) | [w, h] | Texture dimensions |
| **LoadRenderTexture** | (width, height) | id | Create render texture |
| **UnloadRenderTexture** | (id) | — | Unload render texture |
| **BeginTextureMode** | (targetId) | — | Begin drawing to texture |
| **EndTextureMode** | () | — | End drawing to texture |
| **DrawTexture** | (id, x, y) | — | Draw texture at (x,y) |
| **DrawTextureV** | (id, pos) | — | Draw texture (position) |
| **DrawTextureEx** | (id, pos, rotation, scale, tint) | — | Draw texture (rotated, scaled) |
| **DrawTextureRec** | (id, source, pos, tint) | — | Draw texture (source rect) |
| **DrawTexturePro** | (id, source, dest, origin, rotation, tint) | — | Draw texture (full transform) |
| **DrawTextureFlipH** | (textureId, x, y [, tint]) | — | Draw texture flipped horizontally |
| **DrawTextureFlipV** | (textureId, x, y [, tint]) | — | Draw texture flipped vertically |
| **DrawTextureNPatch** | (id, nPatchInfo, dest, origin, rotation, tint) | — | Draw 9-patch texture |
| **LoadTextureFromImage** | (imageId) | id | Create texture from image |
| **LoadTextureCubemap** | (imageId, layout) | id | Load cubemap |
| **IsTextureValid** | (id) | bool | True if valid |
| **IsRenderTextureValid** | (id) | bool | True if valid |
| **UpdateTexture** | (id, pixels) | — | Update texture data |
| **UpdateTextureRec** | (id, rec, pixels) | — | Update texture region |
| **GenTextureMipmaps** | (id) | — | Generate mipmaps |
| **SetTextureFilter** | (id, filter) | — | Set filter mode |
| **SetTextureWrap** | (id, wrap) | — | Set wrap mode |

### 2D sprite animation (raylib_anim2d.go)

| Command | Arguments | Returns | Description |
|---------|-----------|---------|-------------|
| **CreateSpriteAnimation** | (textureId, frameWidth, frameHeight, framesPerRow [, totalFrames]) | animId | Create sprite animation |
| **SetSpriteAnimationFPS** | (animId, fps) | — | Set FPS |
| **SetSpriteAnimationLoop** | (animId, loop) | — | Set loop (bool) |
| **SetSpriteAnimationFrame** | (animId, frameIndex) | — | Set current frame |
| **UpdateSpriteAnimation** | (animId, deltaTime) | — | Advance by delta time |
| **GetSpriteAnimationFrame** | (animId) | int | Current frame index |
| **DrawSpriteAnimation** | (animId, posX, posY [, scaleX, scaleY, rotation, r,g,b,a]) | — | Draw animated sprite |
| **DestroySpriteAnimation** | (animId) | — | Free animation |

### Sprite (raylib_sprite.go)

| Command | Arguments | Returns | Description |
|---------|-----------|---------|-------------|
| **CreateSprite** | (textureId) | spriteId | Create sprite from texture |
| **SpriteSetPosition** | (spriteId, x, y) | — | Set position |
| **SpriteSetScale** | (spriteId, scale) | — | Set uniform scale |
| **SpriteSetScaleXY** | (spriteId, sx, sy) | — | Set X/Y scale |
| **SpriteSetRotation** | (spriteId, angleRad) | — | Set rotation |
| **SpriteSetOrigin** | (spriteId, ox, oy) | — | Set pivot (in texture pixels) |
| **SpriteSetFlip** | (spriteId, flipX, flipY) | — | Set flip (0/1 or bool) |
| **SpriteDraw** | (spriteId [, tint]) | — | Draw sprite with current transform |
| **DestroySprite** | (spriteId) | — | Free sprite |

---

## 6. Raylib (images) – `raylib_images.go`

| Command | Arguments | Returns | Description |
|---------|-----------|---------|-------------|
| **LoadImage** | (path) | id | Load image from file |
| **LoadImageRaw** | (width, height, format, data) | id | Load from raw data |
| **LoadImageAnim** | (path) | id | Load animated image |
| **GetLoadImageAnimFrames** | (path) | int | Frame count of animated image |
| **LoadImageFromMemory** | (data) | id | Load from memory |
| **LoadImageFromTexture** | (textureId) | id | Create image from texture |
| **LoadImageFromScreen** | () | id | Capture screen |
| **IsImageValid** | (id) | bool | True if valid |
| **UnloadImage** | (id) | — | Unload image |
| **ExportImage** | (id, path) | bool | Export to file |
| **ExportImageToMemory** | (id) | string | Export to memory |
| **ExportImageAsCode** | (imageId, fileName) | bool | Export as C header |
| **GenImageColor** | (width, height, color) | id | Solid color image |
| **GenImageGradientLinear** | (w, h, start, end) | id | Linear gradient |
| **GenImageGradientRadial** | (w, h, density, inner, outer) | id | Radial gradient |
| **ImageCopy** | (id) | id | Copy image |
| **ImageCrop** | (id, crop) | — | Crop to rect |
| **ImageResize** | (id, newW, newH) | — | Resize image |
| **ImageResizeNN** | (id, newW, newH) | — | Resize (nearest-neighbor) |
| **ImageFlipVertical** | (id) | — | Flip vertically |
| **ImageFlipHorizontal** | (id) | — | Flip horizontally |
| **ImageRotate** | (id, degrees) | — | Rotate image |
| **ImageColorTint** | (id, color) | — | Tint image |
| **ImageClearBackground** | (id, color) | — | Clear image with color |
| **LoadImageColors** | (id) | count | Get pixel colors; use GetLoadedImageColor |
| **UnloadImageColors** | (colors) | — | Unload color array |
| **GetLoadedImageColor** | (index) | color | Color at index |
| **GetImageColor** | (id, x, y) | color | Pixel at (x,y) |

Other image commands: ImageFromImage, ImageFromChannel, ImageText, ImageTextEx, ImageFormat, ImageToPOT, ImageAlphaCrop, ImageAlphaClear, ImageAlphaMask, ImageAlphaPremultiply, ImageBlurGaussian, ImageKernelConvolution, ImageResizeCanvas, ImageMipmaps, ImageDither, ImageRotateCW, ImageRotateCCW, ImageColorInvert, ImageColorGrayscale, ImageColorContrast, ImageColorBrightness, ImageColorReplace, ImageDrawPixel, ImageDrawPixelV, ImageDrawLine, ImageDrawLineV, ImageDrawLineEx, ImageDrawCircle, ImageDrawCircleV, ImageDrawCircleLines, ImageDrawCircleLinesV, ImageDrawRectangle, ImageDrawRectangleV, ImageDrawRectangleRec, ImageDrawRectangleLines, ImageDrawTriangle, ImageDrawTriangleEx, ImageDrawTriangleLines, ImageDrawTriangleFan, ImageDrawTriangleStrip. GenImageGradientSquare, GenImageChecked, GenImageWhiteNoise, GenImagePerlinNoise, GenImageCellular, GenImageText.

---

## 7. Raylib (3D) – `raylib_3d.go`

### 3D mode and primitives

| Command | Arguments | Returns | Description |
|---------|-----------|---------|-------------|
| **SetCamera3D** | (camera) | — | Set 3D camera |
| **BeginMode3D** | (camera) | — | Begin 3D mode. In the hybrid loop (when using update/draw), the engine wraps 3D automatically; calling these in draw() has no effect. |
| **EndMode3D** | () | — | End 3D mode. In the hybrid loop, calling in draw() has no effect. |
| **DrawCube** | (posX, posY, posZ, width, height, length, color) | — | Filled 3D cube |
| **DrawCubeV** | (position, size, color) | — | Filled cube (vectors) |
| **DrawCubeWires** | (posX, posY, posZ, width, height, length, color) | — | Cube outline |
| **DrawCubeWiresV** | (position, size, color) | — | Cube outline (vectors) |
| **DrawSphere** | (centerX, centerY, centerZ, radius, color) | — | Filled sphere |
| **DrawSphereEx** | (center, radius, rings, slices, color) | — | Filled sphere (subdiv) |
| **DrawSphereWires** | (centerX, centerY, centerZ, radius, rings, slices, color) | — | Sphere outline |
| **DrawPlane** | (center, size, color) | — | Draw plane |
| **DrawLine3D** | (start, end, color) | — | 3D line |
| **DrawPoint3D** | (position, color) | — | 3D point |
| **DrawCircle3D** | (center, radius, rotationAxis, rotationAngle, color) | — | 3D circle |
| **DrawCylinder** | (position, radiusTop, radiusBottom, height, slices, color) | — | Filled cylinder |
| **DrawCylinderEx** | (start, end, radiusStart, radiusEnd, slices, color) | — | Cylinder (caps) |
| **DrawCylinderWires** | (position, radiusTop, radiusBottom, height, slices, color) | — | Cylinder outline |
| **DrawCapsule** | (start, end, radius, slices, rings, color) | — | Filled capsule |
| **DrawCapsuleWires** | (start, end, radius, slices, rings, color) | — | Capsule outline |
| **DrawRay** | (ray, color) | — | Draw ray |
| **DrawTriangle3D** | (v1, v2, v3, color) | — | 3D triangle |
| **DrawTriangleStrip3D** | (points, color) | — | 3D triangle strip |
| **DrawBoundingBox** | (box, color) | — | Draw AABB |

### Models

| Command | Arguments | Returns | Description |
|---------|-----------|---------|-------------|
| **LoadModel** | (path) | id | Load model from file |
| **LoadModelFromMesh** | (meshId) | id | Create model from mesh |
| **UnloadModel** | (id) | — | Unload model |
| **DrawModel** | (id, posX, posY, posZ, scale [, tint]) | — | Draw model at position |
| **DrawModelEx** | (id, position, rotationAxis, rotationAngle, scale [, tint]) | — | Draw model (full transform) |
| **DrawModelWires** | (id, position, scale, tint) | — | Draw model wireframe |
| **DrawModelWiresEx** | (id, position, rotationAxis, rotationAngle, scale, tint) | — | Model wireframe (rotated) |
| **DrawModelPoints** | (id, position, scale, tint) | — | Draw model as points |
| **DrawModelPointsEx** | (id, position, rotationAxis, rotationAngle, scale, tint) | — | Model points (rotated) |
| **IsModelValid** | (id) | bool | True if valid |
| **GetModelBoundingBox** | (id) | box | Model AABB |
| **SetModelMeshMaterial** | (modelId, meshId, materialId) | — | Set mesh material |
| **LoadCube** | (size) | id | Create cube model |
| **SetModelColor** | (modelId, r, g, b, a) | — | Stored tint for DrawModelSimple |
| **SetModelPosition** | (modelId, x, y, z) | — | Store position for DrawModelWithState |
| **SetModelRotation** | (modelId, axisX, axisY, axisZ, angleRad) | — | Store rotation |
| **SetModelScale** | (modelId, sx, sy, sz) | — | Store scale |
| **DrawModelWithState** | (modelId [, tint]) | — | Draw model using stored transform |
| **RotateModel** | (modelId, speedDegPerSec) | — | Add rotation each frame |
| **DrawModelSimple** | (id, x, y, z [, angle]) | — | Draw at (x,y,z), scale 1; uses SetModelColor/RotateModel |

### Model animation

| Command | Arguments | Returns | Description |
|---------|-----------|---------|-------------|
| **LoadModelAnimations** | (path) | count | Load animations; use GetModelAnimationId |
| **GetModelAnimationId** | (path, index) | id | Animation id |
| **UpdateModelAnimation** | (modelId, animId, frame) | — | Set animation frame |
| **UpdateModelAnimationBones** | (modelId, animId, frame) | — | Update bones |
| **UnloadModelAnimation** | (animId) | — | Unload animation |
| **UnloadModelAnimations** | (animIds) | — | Unload animations |
| **IsModelAnimationValid** | (modelId, animId) | bool | True if valid |
| **GetModelAnimationFrameCount** | (animId) | int | Frame count |
| **CreateModelAnimState** | (modelId, animId, fps [, loop]) | stateId | Create anim state |
| **UpdateModelAnimState** | (stateId, deltaTime) | — | Advance by delta |
| **SetModelAnimStateFrame** | (stateId, frameIndex) | — | Set frame |
| **GetModelAnimStateFrame** | (stateId) | float | Current frame |
| **DestroyModelAnimState** | (stateId) | — | Free state |

### Camera (global and objects)

| Command | Arguments | Returns | Description |
|---------|-----------|---------|-------------|
| **CameraOrbit** | (targetX, targetY, targetZ, angleRad, pitchRad, distance) | — | Orbit camera; updates state |
| **CameraZoom** | (amount) | — | Adjust orbit distance (e.g. GetMouseWheelMove()) |
| **CameraRotate** | (deltaX, deltaY) | — | Mouse-delta rotation (2 args) |
| **CameraRotate** | (pitchRad, yawRad, rollRad) | — | Absolute rotation (3 args) |
| **SetCameraTarget** | (x, y, z) | — | Orbit look-at target (3 args) |
| **SetCameraTarget** | (cameraId, x, y, z) | — | Set target for named camera (4 args) |
| **UpdateCamera** | () | — | Apply orbit state to camera |
| **MouseOrbitCamera** | () | — | One call: mouse → rotate, wheel → zoom, update |
| **MouseLook** | () | — | FPS-style camera from mouse |
| **CameraLookAt** | (x, y, z) | — | Look at point |
| **CameraMove** | (dx, dy, dz) | — | Move camera |
| **SetCameraFOV** | (fov) | — | Set field of view |
| **SetCameraPosition** | (x, y, z) | — | Global camera position (3 args) |
| **SetCameraPosition** | (cameraId, x, y, z) | — | Named camera position (4 args) |
| **CAMERA3D** | () | cameraId | Create camera object |
| **SetCameraUp** | (cameraId, x, y, z) | — | Set camera up vector |
| **SetCameraFovy** | (cameraId, fovy) | — | Set FOV (y) |
| **SetCameraProjection** | (cameraId, projection) | — | Set projection type |
| **SetCurrentCamera** | (cameraId) | — | Set active camera |
| **CAMERA_PERSPECTIVE** | () | int | Perspective projection |
| **CAMERA_ORTHOGRAPHIC** | () | int | Orthographic projection |
| **GetCameraPositionX/Y/Z** | () | float | Current camera position |
| **GetCameraTargetX/Y/Z** | () | float | Current camera target |

### Billboards, fog, scene, views

| Command | Arguments | Returns | Description |
|---------|-----------|---------|-------------|
| **DrawBillboard** | (camera, textureId, position, size, tint) | — | Draw billboard |
| **DrawBillboardRec** | (camera, textureId, source, position, size, tint) | — | Billboard (source rect) |
| **DrawBillboardPro** | (camera, textureId, source, position, up, size, origin, rotation, tint) | — | Billboard (full) |
| **SetFog** | (enable, density, r, g, b) | — | Set fog (raylib_fog.go) |
| **SetFogDensity** | (density) | — | Fog density (0.02–0.05) |
| **SetFogColor** | (r, g, b) | — | Fog color |
| **EnableFog** | () | — | Enable fog |
| **DisableFog** | () | — | Disable fog |
| **IsFogEnabled** | () | int | 1 if enabled, 0 otherwise |
| **BeginFog** | () | — | Begin fog (before drawing) |
| **EndFog** | () | — | End fog |
| **CreateScene** | (sceneId) | — | Create scene (scene.go) |
| **LoadScene** | (sceneId) | — | Load scene |
| **UnloadScene** | (sceneId) | — | Unload scene |
| **SetCurrentScene** | (sceneId) | — | Set current scene |
| **GetCurrentScene** | () | sceneId | Current scene id or "" |
| **SetSceneWorld** | (sceneId, worldId) | — | Set scene physics world |
| **SaveScene** | (sceneId, path) | — | Save scene to JSON |
| **LoadSceneFromFile** | (path) | sceneId | Load scene from file |
| **CreateView** | (viewId, x, y, width, height) | — | Create viewport (raylib_views.go) |
| **SetViewTarget** | (viewId, renderTextureId) | — | Set view render target |
| **DrawView** | (viewId) | — | Draw view to screen |
| **GetViewX/Y/Width/Height** | (viewId) | int | View rect |
| **SetViewPosition** | (viewId, x, y) | — | Set view position |
| **SetViewSize** | (viewId, width, height) | — | Set view size |
| **SetViewRect** | (viewId, x, y, width, height) | — | Set view rect |
| **CreateSplitscreenLeftRight** | (viewIdLeft, viewIdRight) | — | Split left/right |
| **CreateSplitscreenTopBottom** | (viewIdTop, viewIdBottom) | — | Split top/bottom |
| **CreateSplitscreenFour** | (viewIdTL, viewIdTR, viewIdBL, viewIdBR) | — | Four-way split |

See [Windows, scaling, and splitscreen](docs/WINDOWS_AND_VIEWS.md).

### 3D editor and level builder (raylib_editor.go)

Picking: **GetMouseRay**() (updates internal ray), **GetMouseRayOriginX/Y/Z**(), **GetMouseRayDirectionX/Y/Z**(). **GetRayCollisionPlane**(…), **GetRayCollisionPointX/Y/Z**(), **PickGroundPlane**() → 1 if hit on y=0. **SnapToGridX/Y/Z**(value, gridSize). Level objects: **CreateLevelObject**(id, modelId, x, y, z, rotX, rotY, rotZ, scaleX, scaleY, scaleZ), **SetObjectPosition**(id, x, y, z), **RotateObject**(id, pitch, yaw, roll), **ScaleObject**(id, sx, sy, sz), **DrawObject**(id) / **DrawLevelObject**(id), **SetLevelObjectTransform**(id, …), **GetLevelObjectX/Y/Z**, **GetLevelObjectRotX/RotY/RotZ**, **GetLevelObjectScaleX/ScaleY/ScaleZ**, **GetLevelObjectModelId**(id), **DeleteLevelObject**(id), **GetLevelObjectCount**(), **GetLevelObjectId**(index), **SaveLevel**(path), **LoadLevel**(path), **DuplicateLevelObject**(id) → newId.

**Lighting (stubs):** ENABLELIGHTING(), LIGHT() → lightId, LIGHT_DIRECTIONAL() → 0, SetLightType, SetLightPosition, SetLightTarget, SetLightColor, SetLightIntensity, SETAMBIENTLIGHT(r, g, b).

---

## 8. Raylib (mesh) – `raylib_mesh.go`

| Command | Arguments | Returns | Description |
|---------|-----------|---------|-------------|
| **GenMeshPoly** | (sides, radius) | id | Polygon mesh |
| **GenMeshPlane** | (width, length, resX, resZ) | id | Plane mesh |
| **GenMeshCube** | (width, height, length) | id | Cube mesh |
| **GenMeshSphere** | (radius, rings, slices) | id | Sphere mesh |
| **GenMeshHemiSphere** | (radius, rings, slices) | id | Hemisphere mesh |
| **GenMeshCylinder** | (radius, height, slices) | id | Cylinder mesh |
| **GenMeshCone** | (radius, height, slices) | id | Cone mesh |
| **GenMeshTorus** | (radius, size, radSeg, sides) | id | Torus mesh |
| **GenMeshKnot** | (radius, size, radSeg, sides) | id | Knot mesh |
| **GenMeshHeightmap** | (imageId, size) | id | Heightmap mesh |
| **GenMeshCubicmap** | (imageId, cubeSize) | id | Cubicmap mesh |
| **UploadMesh** | (meshId) | — | Upload to GPU |
| **UnloadMesh** | (id) | — | Unload mesh |
| **GetMeshBoundingBox** | (id) | box | Mesh AABB |
| **ExportMesh** | (id, path) | bool | Export mesh |
| **DrawMesh** | (id, materialId, posX,Y,Z, scaleX,Y,Z) | — | Draw mesh (position + scale) |
| **DrawMeshMatrix** | (meshId, materialId, m0..m15) | — | Draw mesh with full 4x4 matrix (row-major) |
| **UpdateMeshBuffer** | (id, index, data) | — | Update mesh buffer |
| **DrawMeshInstanced** | (id, materialId, transforms) | — | Draw instanced |
| **LoadMaterialDefault** | () | id | Default material |
| **IsMaterialValid** | (id) | bool | True if valid |
| **UnloadMaterial** | (id) | — | Unload material |
| **SetMaterialTexture** | (materialId, mapType, textureId) | — | Set material texture |
| **LoadMaterials** | (path) | count | Load materials |
| **GetMaterialIdFromLoad** | (path, index) | id | Material id |
| **GetRayCollisionMesh** | (ray, meshId, transform) | hit | Ray vs mesh |

---

## 9. Raylib (audio) – `raylib_audio.go`

| Command | Arguments | Returns | Description |
|---------|-----------|---------|-------------|
| **InitAudioDevice** | () | — | Initialize audio |
| **CloseAudioDevice** | () | — | Close audio |
| **IsAudioDeviceReady** | () | bool | True if ready |
| **LoadSound** | (path) | id | Load sound |
| **UnloadSound** | (id) | — | Unload sound |
| **PlaySound** | (id) | — | Play sound |
| **StopSound** | (id) | — | Stop sound |
| **PauseSound** | (id) | — | Pause sound |
| **ResumeSound** | (id) | — | Resume sound |
| **IsSoundPlaying** | (id) | bool | True if playing |
| **SetSoundVolume** | (id, volume) | — | Volume (0–1) |
| **SetSoundPitch** | (id, pitch) | — | Pitch |
| **SetSoundPan** | (id, pan) | — | Pan (-1 to 1) |
| **LoadMusicStream** | (path) | id | Load music |
| **UnloadMusicStream** | (id) | — | Unload music |
| **PlayMusicStream** | (id) | — | Play music |
| **StopMusicStream** | (id) | — | Stop music |
| **PauseMusicStream** | (id) | — | Pause music |
| **ResumeMusicStream** | (id) | — | Resume music |
| **IsMusicStreamPlaying** | (id) | bool | True if playing |
| **UpdateMusicStream** | (id) | — | Update stream (call each frame) |
| **SetMusicVolume** | (id, volume) | — | Volume (0–1) |
| **SetMusicPitch** | (id, pitch) | — | Pitch |
| **SetMusicPan** | (id, pan) | — | Pan |
| **SeekMusicStream** | (id, position) | — | Seek position |
| **GetMusicTimeLength** | (id) | float | Length in seconds |
| **GetMusicTimePlayed** | (id) | float | Time played |
| **IsMusicValid** | (id) | bool | True if valid |
| **SetMasterVolume** | (volume) | — | Master volume (0–1) |
| **GetMasterVolume** | () | float | Master volume |
| **LoadWave** | (path) | id | Load wave |
| **UnloadWave** | (id) | — | Unload wave |
| **LoadSoundFromWave** | (waveId) | id | Create sound from wave |
| **LoadAudioStream** | (sampleRate, sampleSize, channels) | id | Create audio stream |
| **UnloadAudioStream** | (id) | — | Unload stream |
| **UpdateAudioStream** | (id, data) | — | Push samples |
| **PlayAudioStream** | (id) | — | Play stream |
| **StopAudioStream** | (id) | — | Stop stream |
| **SetAudioStreamVolume** | (id, volume) | — | Volume |
| **SetAudioStreamPitch** | (id, pitch) | — | Pitch |
| **SetAudioStreamPan** | (id, pan) | — | Pan |
| **ExportWaveAsCode** | (waveId, fileName) | bool | Export as C header |

Not supported from BASIC (return error): SetAudioStreamCallback, AttachAudioStreamProcessor, DetachAudioStreamProcessor, AttachAudioMixedProcessor, DetachAudioMixedProcessor.

---

## 10. Raylib (fonts) – `raylib_fonts.go`

| Command | Arguments | Returns | Description |
|---------|-----------|---------|-------------|
| **GetFontDefault** | () | id | Default font |
| **LoadFont** | (path) | id | Load font |
| **LoadFontEx** | (path, size, …) | id | Load font (size) |
| **UnloadFont** | (id) | — | Unload font |
| **LoadFontFromImage** | (imageId, key, fontSize) | id | Load from image |
| **LoadFontFromMemory** | (fileType, data, fontSize, …) | id | Load from memory |
| **IsFontValid** | (id) | bool | True if valid |
| **LoadFontData** | (data, fontSize, …) | count | Load font data |
| **UnloadFontData** | (chars) | — | Unload font data |
| **GenImageFontAtlas** | (chars, fontSize, …) | id | Generate atlas |
| **DrawTextExFont** | (fontId, text, pos, size, spacing, tint) | — | Draw text with font |
| **MeasureTextEx** | (fontId, text, size, spacing) | width, height | Text size |
| **ExportFontAsCode** | (fontId, fileName) | bool | Export as C header |
| **DrawTextCodepoint** | (fontId, codepoint, pos, size, tint) | — | Draw one codepoint |
| **DrawTextCodepoints** | (fontId, codepoints, pos, size, spacing, tint) | — | Draw codepoints |
| **GetGlyphIndex** | (fontId, codepoint) | int | Glyph index |
| **GetGlyphInfo** | (fontId, codepoint) | info | Glyph info |
| **GetGlyphAtlasRec** | (fontId, codepoint) | rect | Glyph rect in atlas |

---

## 11. Raylib (misc) – `raylib_misc.go`

| Command | Arguments | Returns | Description |
|---------|-----------|---------|-------------|
| **GetMouseDelta** | () | x, y | Mouse delta this frame |
| **NewColor** | (r, g, b, a) | color | Create color (0–255) |
| **Color** | (r, g, b, a) | color | Same as NewColor |
| **Fade** | (color, alpha) | color | Color with alpha |
| **ColorAlpha** | (color, alpha) | color | Set alpha |
| **ColorToInt** | (color) | int | Packed color int |
| **GetColor** | (hexValue) | color | Unpack hex to color |
| **ColorIsEqual** | (c1, c2) | bool | True if equal |
| **ColorNormalize** | (color) | r, g, b, a | Normalized (0–1) |
| **ColorFromNormalized** | (r, g, b, a) | color | From normalized |
| **ColorToHSV** | (color) | h, s, v | Color to HSV |
| **ColorFromHSV** | (h, s, v) | color | HSV to color |
| **ColorTint** | (color, tint) | color | Tint color |
| **ColorBrightness** | (color, factor) | color | Brightness |
| **ColorContrast** | (color, contrast) | color | Contrast |
| **ColorAlphaBlend** | (dst, src) | color | Blend colors |
| **ColorLerp** | (c1, c2, t) | color | Lerp colors |
| **GetPixelDataSize** | (width, height, format) | int | Pixel buffer size |
| **CheckCollisionRecs** | (rec1, rec2) | bool | Rectangle overlap |
| **CheckCollisionCircles** | (center1, r1, center2, r2) | bool | Circle overlap |
| **CheckCollisionCircleRec** | (center, radius, rec) | bool | Circle vs rect |
| **CheckCollisionPointRec** | (point, rec) | bool | Point in rect |
| **CheckCollisionPointCircle** | (point, center, radius) | bool | Point in circle |
| **GetCollisionRec** | (rec1, rec2) | rec | Overlap rect |
| **CheckCollisionSpheres** | (c1, r1, c2, r2) | bool | Sphere overlap |
| **CheckCollisionBoxes** | (box1, box2) | bool | AABB overlap |
| **CheckCollisionBoxSphere** | (box, center, radius) | bool | Box vs sphere |
| **GetRayCollisionSphere** | (ray, center, radius) | hit | Ray vs sphere |
| **GetRayCollisionBox** | (ray, box) | hit | Ray vs box |
| **GetRayCollisionTriangle** | (ray, p1, p2, p3) | hit | Ray vs triangle |
| **GetRayCollisionQuad** | (ray, p1, p2, p3, p4) | hit | Ray vs quad |
| **GetRayCollisionPointX/Y/Z** | () | float | Hit point from last ray |
| **GetRayCollisionNormalX/Y/Z** | () | float | Hit normal |
| **GetRayCollisionDistance** | () | float | Hit distance |

**Color constants (0-arg):** White, Black, LightGray, Gray, DarkGray, Yellow, Gold, Orange, Pink, Red, Maroon, Green, Lime, DarkGreen, SkyBlue, Blue, DarkBlue, Purple, Violet, DarkPurple, Beige, Brown, DarkBrown, Magenta, RayWhite, Blank.

---

## 12. Raylib (math) – `raylib_math.go`

| Command | Arguments | Returns | Description |
|---------|-----------|---------|-------------|
| **Clamp** | (value, min, max) | float | Clamp to [min, max] |
| **Lerp** | (a, b, t) | float | Linear interpolate |
| **Normalize** | (value, start, end) | float | Normalize to 0–1 |
| **Remap** | (value, inputStart, inputEnd, outputStart, outputEnd) | float | Remap range |
| **Wrap** | (value, min, max) | float | Wrap value |
| **FloatEquals** | (a, b) | bool | Float equality |
| **VECTOR3** | (x, y, z) | [x,y,z] | 3D vector (use with DrawModel) |
| **Vector2Zero** | () | vec | Zero vector |
| **Vector2One** | () | vec | One vector |
| **Vector2Add** | (v1, v2) | vec | Add vectors |
| **Vector2Subtract** | (v1, v2) | vec | Subtract |
| **Vector2Length** | (v) | float | Length |
| **Vector2Distance** | (v1, v2) | float | Distance |
| **Vector2Scale** | (v, scale) | vec | Scale |
| **Vector2Normalize** | (v) | vec | Normalize |
| **Vector2Lerp** | (v1, v2, t) | vec | Lerp |
| **Vector2Rotate** | (v, angle) | vec | Rotate |
| **Vector3Zero** | () | vec | Zero vector |
| **Vector3One** | () | vec | One vector |
| **Vector3Add** | (v1, v2) | vec | Add |
| **Vector3Subtract** | (v1, v2) | vec | Subtract |
| **Vector3Scale** | (v, scale) | vec | Scale |
| **Vector3Length** | (v) | float | Length |
| **Vector3Distance** | (v1, v2) | float | Distance |
| **Vector3Normalize** | (v) | vec | Normalize |
| **Vector3CrossProduct** | (v1, v2) | vec | Cross product |
| **Vector3DotProduct** | (v1, v2) | float | Dot product |
| **Vector3Lerp** | (v1, v2, t) | vec | Lerp |
| **MatrixIdentity** | () | mat | Identity matrix |
| **MatrixMultiply** | (left, right) | mat | Multiply |
| **MatrixTranslate** | (x, y, z) | mat | Translation |
| **MatrixRotate** | (axis, angle) | mat | Rotation |
| **MatrixRotateX** | (angle) | mat | Rotate X |
| **MatrixRotateY** | (angle) | mat | Rotate Y |
| **MatrixRotateZ** | (angle) | mat | Rotate Z |
| **MatrixScale** | (x, y, z) | mat | Scale |
| **QuaternionIdentity** | () | quat | Identity quaternion |
| **QuaternionMultiply** | (q1, q2) | quat | Multiply |
| **QuaternionLerp** | (q1, q2, t) | quat | Lerp |
| **QuaternionSlerp** | (q1, q2, t) | quat | Spherical lerp |
| **QuaternionFromEuler** | (pitch, yaw, roll) | quat | From Euler |
| **QuaternionToEuler** | (q) | pitch, yaw, roll | To Euler |

Other Vector2/Vector3/Matrix/Quaternion helpers: Vector2AddValue, Vector2SubtractValue, Vector2LengthSqr, Vector2DotProduct, Vector2DistanceSqr, Vector2Angle, Vector2Multiply, Vector2Negate, Vector2Divide, Vector2Transform, Vector2Reflect, Vector2MoveTowards, Vector2Invert, Vector2Clamp, Vector2ClampValue, Vector2Equals; Vector3AddValue, Vector3SubtractValue, Vector3LengthSqr, Vector3Perpendicular, Vector3OrthoNormalize, Vector3Transform, Vector3RotateByQuaternion, Vector3RotateByAxisAngle, Vector3Reflect, Vector3Min, Vector3Max, Vector3Barycenter, Vector3Unproject, Vector3Invert, Vector3Clamp, Vector3ClampValue, Vector3Equals, Vector3Refract, Vector3ToFloatV; MatrixDeterminant, MatrixTrace, MatrixTranspose, MatrixInvert, MatrixAdd, MatrixSubtract, MatrixRotateXYZ, MatrixRotateZYX, MatrixFrustum, MatrixPerspective, MatrixOrtho, MatrixLookAt, MatrixToFloatV; QuaternionAdd, QuaternionAddValue, QuaternionSubtract, QuaternionSubtractValue, QuaternionLength, QuaternionNormalize, QuaternionInvert, QuaternionScale, QuaternionDivide, QuaternionNlerp, QuaternionFromVector3ToVector3, QuaternionFromMatrix, QuaternionToMatrix, QuaternionFromAxisAngle, QuaternionToAxisAngle, QuaternionToEuler, QuaternionTransform, QuaternionEquals.

---

## 13. Raylib (game) – `raylib_game.go`

| Command | Arguments | Returns | Description |
|---------|-----------|---------|-------------|
| **GAME.CameraOrbit** | (…) | — | Orbit camera (GAME namespace) |
| **GAME.MoveWASD** | (…) | — | WASD movement |
| **GAME.OnGround** | (…) | bool | On ground check |
| **GAME.SnapToGround** | (…) | — | Snap to ground |
| **MoveWASD3D** | (…) | — | 3D WASD movement |
| **SnapToGround3D** | (…) | — | 3D snap to ground |
| **IsOnGround3D** | (…) | bool | 3D on ground |
| **CameraOrbit3D** | (…) | — | 3D orbit camera |
| **MoveHorizontal2D** | (…) | — | 2D horizontal move |
| **Jump2D** | (…) | — | 2D jump |
| **IsOnGround2D** | (…) | bool | 2D on ground |
| **ClampVelocity2D** | (…) | — | Clamp 2D velocity |
| **MoveVertical2D** | (…) | — | 2D vertical move |
| **CameraFollow2D** | (…) | — | 2D camera follow |
| **SnapToPlatform2D** | (…) | — | Snap to platform |
| **Jump3D** | (…) | — | 3D jump |
| **ClampVelocity3D** | (…) | — | Clamp 3D velocity |
| **CameraFollow3D** | (…) | — | 3D camera follow |
| **GAME.GetAxisX** | () | int | -1, 0, or 1 for A/D |
| **GAME.GetAxisY** | () | int | -1, 0, or 1 for W/S |
| **GetAxisX** | () | int | Same (no namespace) |
| **GetAxisY** | () | int | Same (no namespace) |
| **GAME.SyncSpriteToBody2D** | (worldId, bodyId, spriteId) | — | Set sprite to Box2D body; call in draw |
| **GAME.SetCamera2DFollow** | (worldId, bodyId, xOffset, yOffset) | — | 2D follow preset |
| **GAME.UpdateCamera2D** | () | — | Update 2D camera each frame |
| **GAME.SetCamera3DOrbit** | (worldId, bodyId, distance, heightOffset) | — | 3D orbit preset |
| **GAME.UpdateCamera3D** | (angleRad, pitchRad) | — | Update 3D camera |
| **GAME.SetCollisionHandler** | (bodyId, subName) | — | When bodyId collides, call Sub subName(otherBodyId) |
| **GAME.ProcessCollisions2D** | (worldId) | — | Invoke handlers; call after BOX2D.Step |
| **GAME.AssetPath** | (filename) | string | "assets/" + filename |
| **GAME.ClampDelta** | (maxDt) | float | min(GetFrameTime(), maxDt) |
| **GAME.ShowDebug** | () | — | Draw FPS |
| **ShowDebug** | (extraText) | — | Draw FPS and extra line |
| **GetScreenCenterX** | () | float | Screen center X |
| **GetScreenCenterY** | () | float | Screen center Y |
| **Distance2D** | (x1, y1, x2, y2) | float | 2D distance |
| **Distance3D** | (x1, y1, z1, x2, y2, z2) | float | 3D distance |
| **SetCamera2DCenter** | (worldX, worldY) | — | 2D camera center at (worldX, worldY) |
| **Camera3DMoveForward** | (amount) | — | Move camera and target along look direction |
| **Camera3DMoveBackward** | (amount) | — | Move backward |
| **Camera3DMoveRight** | (amount) | — | Move along right |
| **Camera3DMoveLeft** | (amount) | — | Move along left |
| **Camera3DMoveUp** | (amount) | — | Move along camera up |
| **Camera3DMoveDown** | (amount) | — | Move along camera down |
| **Camera3DRotateYaw** | (angleRad) | — | Rotate position around target (Y axis) |
| **Camera3DRotatePitch** | (angleRad) | — | Tilt camera up/down |
| **Camera3DRotateRoll** | (angleRad) | — | Rotate camera up vector around forward axis |

**Key constants (0-arg):** GAME.KEY_W, GAME.KEY_A, GAME.KEY_S, GAME.KEY_D, GAME.KEY_SPACE.

---

## 14. Box2D – `box2d.go`

Use **BOX2D.*** prefix or legacy flat names.

| Command | Arguments | Returns | Description |
|---------|-----------|---------|-------------|
| **BOX2D.CreateWorld** | () | worldId | Create 2D world |
| **BOX2D.Step** | (worldId, dt) | — | Step simulation |
| **BOX2D.DestroyWorld** | (worldId) | — | Destroy world |
| **BOX2D.CreateBody** | (worldId, x, y, type) | bodyId | Create body |
| **BOX2D.DestroyBody** | (worldId, bodyId) | — | Destroy body |
| **BOX2D.GetBodyCount** | (worldId) | int | Body count |
| **BOX2D.GetBodyId** | (worldId, index) | bodyId | Body at index |
| **BOX2D.CreateBodyAtScreen** | (worldId, screenX, screenY, …) | bodyId | Create at screen pos |
| **BOX2D.GetPosition** | (worldId, bodyId) | x, y | Body position |
| **BOX2D.GetPositionX** | (worldId, bodyId) | float | Position X |
| **BOX2D.GetPositionY** | (worldId, bodyId) | float | Position Y |
| **BOX2D.GetAngle** | (worldId, bodyId) | float | Body angle (rad) |
| **BOX2D.SetLinearVelocity** | (worldId, bodyId, vx, vy) | — | Set velocity |
| **BOX2D.GetLinearVelocity** | (worldId, bodyId) | vx, vy | Get velocity |
| **BOX2D.SetTransform** | (worldId, bodyId, x, y, angle) | — | Set position and angle |
| **BOX2D.ApplyForce** | (worldId, bodyId, fx, fy, x, y) | — | Apply force |
| **CreateWorld2D** | (worldName$, gravityX, gravityY) | — | Create 2D world |
| **Physics2DCreateWorld** | (gravityX, gravityY) | — | Create world named "default" |
| **DestroyWorld2D** | (worldId) | — | Destroy world |
| **Step2D** | (worldId, dt) | — | Step world |
| **Physics2DStep** | (dt) | — | Step world "default" |
| **CreateBox2D** | (worldId, x, y, w, h, …) | bodyId | Create box body |
| **CreateCircle2D** | (worldId, x, y, radius, …) | bodyId | Create circle body |
| **CreatePolygon2D** | (…) | bodyId | Create polygon |
| **CreateEdge2D** | (…) | bodyId | Create edge |
| **CreateChain2D** | (…) | bodyId | Create chain |
| **GetPositionX2D** | (worldId, bodyId) | float | Position X |
| **GetPositionY2D** | (worldId, bodyId) | float | Position Y |
| **SetPosition2D** | (worldId, bodyId, x, y) | — | Set position |
| **GetAngle2D** | (worldId, bodyId) | float | Angle |
| **SetAngle2D** | (worldId, bodyId, angle) | — | Set angle |
| **GetVelocityX2D** | (worldId, bodyId) | float | Velocity X |
| **GetVelocityY2D** | (worldId, bodyId) | float | Velocity Y |
| **SetVelocity2D** | (worldId, bodyId, vx, vy) | — | Set velocity |
| **ApplyForce2D** | (worldId, bodyId, fx, fy, …) | — | Apply force |
| **ApplyImpulse2D** | (…) | — | Apply impulse |
| **CreateDistanceJoint2D** | (worldId, bodyAId, bodyBId, length) | — | Distance joint (implemented) |
| **RayCast2D** | (worldId, x1, y1, x2, y2) | bool | Ray cast; use RayHit* for result |
| **RayHitX2D** | () | float | Hit X |
| **RayHitY2D** | () | float | Hit Y |
| **RayHitBody2D** | () | bodyId | Hit body |
| **GetCollisionCount2D** | (worldId) | int | Collision count (after Step) |
| **GetCollisionOther2D** | (index) | bodyId | Other body in collision |

Other: SetSensor2D, ApplyTorque2D, SetAngularVelocity2D, GetAngularVelocity2D, SetFriction2D, SetRestitution2D, SetDamping2D, SetFixedRotation2D, SetGravityScale2D, SetMass2D, SetBullet2D, GetCollisionNormalX2D, GetCollisionNormalY2D. Joints (Revolute, Prismatic, etc.) are stubbed.

---

## 15. Bullet (3D physics) – `bullet.go`

Use **BULLET.*** prefix or legacy flat names.

| Command | Arguments | Returns | Description |
|---------|-----------|---------|-------------|
| **BULLET.CreateWorld** | () | worldId | Create 3D world |
| **BULLET.SetGravity** | (worldId, x, y, z) | — | Set gravity |
| **BULLET.Step** | (worldId, dt) | — | Step simulation |
| **BULLET.DestroyWorld** | (worldId) | — | Destroy world |
| **BULLET.CreateBox** | (worldId, x, y, z, w, h, d, mass) | bodyId | Create box body |
| **BULLET.CreateSphere** | (worldId, x, y, z, radius, mass) | bodyId | Create sphere body |
| **BULLET.DestroyBody** | (worldId, bodyId) | — | Destroy body |
| **BULLET.SetPosition** | (worldId, bodyId, x, y, z) | — | Set position |
| **BULLET.GetPositionX/Y/Z** | (worldId, bodyId) | float | Position |
| **BULLET.SetVelocity** | (worldId, bodyId, vx, vy, vz) | — | Set velocity |
| **BULLET.GetVelocityX/Y/Z** | (worldId, bodyId) | float | Velocity |
| **BULLET.GetRotationX/Y/Z** | (worldId, bodyId) | float | Rotation (euler) |
| **BULLET.SetRotation** | (worldId, bodyId, x, y, z) | — | Set rotation |
| **BULLET.ApplyForce** | (worldId, bodyId, fx, fy, fz, x, y, z) | — | Apply force |
| **BULLET.ApplyCentralForce** | (worldId, bodyId, fx, fy, fz) | — | Apply central force |
| **BULLET.ApplyImpulse** | (worldId, bodyId, ix, iy, iz, x, y, z) | — | Apply impulse |
| **BULLET.RayCast** | (worldId, ox, oy, oz, dx, dy, dz) | bool | Ray cast |
| **BULLET.GetRayCastHitX/Y/Z** | () | float | Hit point |
| **BULLET.GetRayCastHitBody** | () | bodyId | Hit body |
| **CreateWorld3D** | () | worldId | Legacy: create world |
| **DestroyWorld3D** | (worldId) | — | Legacy: destroy |
| **Step3D** | (worldId, dt) | — | Legacy: step |
| **CreateSphere3D** | (worldId, x, y, z, radius, mass) | bodyId | Create sphere |
| **CreateBox3D** | (worldId, x, y, z, w, h, d, mass) | bodyId | Create box |
| **GetPositionX3D/Y3D/Z3D** | (worldId, bodyId) | float | Position |
| **SetPosition3D** | (worldId, bodyId, x, y, z) | — | Set position |
| **GetYaw3D** | (worldId, bodyId) | float | Yaw |
| **GetPitch3D** | (worldId, bodyId) | float | Pitch |
| **GetRoll3D** | (worldId, bodyId) | float | Roll |
| **SetRotation3D** | (worldId, bodyId, yaw, pitch, roll) | — | Set rotation |
| **SetVelocity3D** | (worldId, bodyId, vx, vy, vz) | — | Set velocity |
| **ApplyForce3D** | (worldId, bodyId, fx, fy, fz, …) | — | Apply force |
| **ApplyImpulse3D** | (…) | — | Apply impulse |
| **RayCast3D** | (worldId, ox, oy, oz, dx, dy, dz) | bool | Ray cast |
| **RayHitX3D/Y3D/Z3D** | () | float | Hit point |
| **RayHitBody3D** | () | bodyId | Hit body |
| **GetCollisionCount3D** | (worldId) | int | Collision count |
| **GetCollisionOther3D** | (index) | bodyId | Other body |

Other legacy: CreateCapsule3D, CreateStaticMesh3D, CreateCylinder3D, CreateCone3D, CreateHeightmap3D, CreateCompound3D, AddShapeToCompound3D, SetScale3D, GetVelocityX3D/Y3D/Z3D, SetAngularVelocity3D, GetAngularVelocityX3D/Y3D/Z3D, ApplyTorque3D, ApplyTorqueImpulse3D, SetMass3D. Stubs: SetFriction3D, SetRestitution3D, SetDamping3D, SetKinematic3D, joints (CreateHingeJoint3D, etc.).

---

## 16. ECS – `ecs.go`

All commands use **ECS.** prefix. See [ECS_GUIDE.md](docs/ECS_GUIDE.md).

| Command | Arguments | Returns | Description |
|---------|-----------|---------|-------------|
| **ECS.CreateWorld** | () | worldId | Create ECS world |
| **ECS.DestroyWorld** | (worldId) | — | Destroy world |
| **ECS.CreateEntity** | (worldId) | entityId | Create entity |
| **ECS.DestroyEntity** | (worldId, entityId) | — | Destroy entity |
| **ECS.AddComponent** | (worldId, entityId, componentType [, args…]) | — | Add component |
| **ECS.HasComponent** | (worldId, entityId, componentType) | bool | True if has component |
| **ECS.RemoveComponent** | (worldId, entityId, componentType) | — | Remove component |
| **ECS.SetTransform** | (worldId, entityId, x, y, z) | — | Set position |
| **ECS.GetTransformX/Y/Z** | (worldId, entityId) | number | Get position |
| **ECS.PlaceEntity** | (worldId, entityId, x, y, z) | — | Same as SetTransform |
| **ECS.GetWorldPositionX/Y/Z** | (worldId, entityId) | number | World position (with Parent chain) |
| **ECS.GetHealthCurrent** | (worldId, entityId) | number | Current health |
| **ECS.GetHealthMax** | (worldId, entityId) | number | Max health |
| **ECS.QueryCount** | (worldId, componentType1 [, …]) | count | Count entities with component(s) |
| **ECS.QueryEntity** | (worldId, componentType, index) | entityId or "" | Entity at index |

Component types: Transform(x,y,z), Sprite(textureId, visible), Health(current, max), Parent(parentEntityId).

---

## 17. Std (file, string, math, JSON, Enum, Dictionary, HTTP, HELP, multi-window) – `std.go`

### File

| Command | Arguments | Returns | Description |
|---------|-----------|---------|-------------|
| **ReadFile** | (path) | string or nil | Read file; nil on error |
| **WriteFile** | (path, contents) | bool | Write file |
| **LoadText** | (path) | string | Alias for ReadFile |
| **SaveText** | (path, text) | bool | Alias for WriteFile |
| **DeleteFile** | (path) | bool | Delete file |
| **CopyFile** | (src, dst) | bool | Copy file |
| **ListDir** | (path) | count | Directory entries; use GetDirItem(index) |
| **GetDirItem** | (index) | string | Entry at 0-based index |
| **ExecuteFile** | (path) | bool | Start process |
| **IsNull** | (value) | bool | True when value is null |

FileExists(path) is in raylib core. Use **Nil** or **Null** for missing values.

### Enum (requires ENUM in script)

| Command | Arguments | Returns | Description |
|---------|-----------|---------|-------------|
| **Enum.getValue** | (enumName, valueName) | int | Value for name |
| **Enum.getName** | (enumName, value) | string | Name for value, or "" |
| **Enum.hasValue** | (enumName, valueName) | bool | True if value exists |

### Dictionary

| Command | Arguments | Returns | Description |
|---------|-----------|---------|-------------|
| **CreateDict** | () | dict | New empty map |
| **SetDictKey** | (dict, key, value) | dict | Set key (mutates); use GetJSONKey to read |
| **Dictionary.has** | (dict, key) | bool | True if key exists |
| **Dictionary.keys** | (dict) | array | All keys |
| **Dictionary.values** | (dict) | array | All values |
| **Dictionary.size** | (dict) | int | Number of keys |
| **Dictionary.remove** | (dict, key) | dict | Remove key |
| **Dictionary.clear** | (dict) | dict | Clear dict |
| **Dictionary.merge** | (dict1, dict2) | dict | New merged dict |
| **Dictionary.get** | (dict, key [, default]) | value | Value or default |

### String (DBP-style)

| Command | Arguments | Returns | Description |
|---------|-----------|---------|-------------|
| **Left** | (s, n) | string | First n characters |
| **Right** | (s, n) | string | Last n characters |
| **Mid** | (s, start1Based [, count]) | string | Substring (1-based start) |
| **Substr** | (s, start0Based [, count]) | string | Substring (0-based) |
| **Instr** | (s, sub) | int | 1-based index or 0 |
| **Upper** | (s) | string | Uppercase |
| **Lower** | (s) | string | Lowercase |
| **Len** | (s) | int | Character count |
| **Chr** | (code) | string | One character |
| **Asc** | (s) | int | Code of first char |
| **Str** | (x) | string | Number to string |
| **Val** | (s) | float | String to float |

### Math and debug

| Command | Arguments | Returns | Description |
|---------|-----------|---------|-------------|
| **Rnd** | () | float | [0, 1) |
| **Rnd** | (n) | int | 1..n inclusive |
| **Random** | (n) | int | 0..n-1 |
| **Random** | (min, max) | int | [min, max] |
| **Int** | (x) | int | Truncate |
| **Radians** | (degrees) | float | To radians |
| **Degrees** | (radians) | float | To degrees |
| **AngleWrap** / **WrapAngle** | (angle) | float | Wrap to [-π, π] |
| **TimeNow** | () | float | Seconds since epoch |
| **PrintDebug** | (value) | — | Print to stderr |
| **Assert** | (condition [, message]) | — | Abort if falsy |
| **HELP** / **?** | () | — | Print quick reference |

### JSON and HTTP

| Command | Arguments | Returns | Description |
|---------|-----------|---------|-------------|
| **LoadJSON** | (path) | handle | Load JSON; use GetJSONKey |
| **LoadJSONFromString** | (str) | handle | Parse JSON string |
| **GetJSONKey** | (handle, key) | value | Get value |
| **SaveJSON** | (path, handle) | bool | Save to file |
| **HttpGet** | (url) | string or nil | GET request |
| **HttpPost** | (url, body) | string | POST request |
| **DownloadFile** | (url, path) | bool | Download to file |

### Multi-window (multiple processes)

| Command | Arguments | Returns | Description |
|---------|-----------|---------|-------------|
| **GetEnv** | (key) | string | Environment variable |
| **IsWindowProcess** | () | bool | True if run with --window |
| **GetWindowTitle** | () | string | Child window title |
| **GetWindowWidth** | () | int | Child window width |
| **GetWindowHeight** | () | int | Child window height |
| **SpawnWindow** | (port, title, width, height) | 1 or 0 | Start same .bas as child |

See [docs/MULTI_WINDOW.md](docs/MULTI_WINDOW.md).

---

## 18. Multiplayer (TCP) – `net.go`

See [docs/MULTIPLAYER.md](docs/MULTIPLAYER.md).

| Command | Arguments | Returns | Description |
|---------|-----------|---------|-------------|
| **Connect** | (host, port) | connectionId or null | Connect to server |
| **ConnectToParent** | () | connectionId or null | Connect using CYBERBASIC_PARENT (spawned windows) |
| **ConnectTLS** | (host, port) | connectionId or null | Encrypted connect |
| **Send** | (connectionId, text) | bool | Send line (max 256 KB, no newlines) |
| **SendText** | (connectionId, text) | — | Same as Send |
| **SendJSON** | (connectionId, jsonText) | 1 or 0 | Send JSON line |
| **SendInt** | (connectionId, value) | 1 or 0 | Send int |
| **SendFloat** | (connectionId, value) | 1 or 0 | Send float |
| **SendNumbers** | (connectionId, n1, n2, …) | 1 or 0 | Send up to 16 numbers |
| **Receive** | (connectionId) | string or null | Next line |
| **ReceiveJSON** | (connectionId) | string or null | Next line if valid JSON |
| **ReceiveNumbers** | (connectionId) | count | Numbers parsed; use GetReceivedNumber(index) |
| **GetReceivedNumber** | (index) | float | Number at index (0.0 if out of range) |
| **Disconnect** | (connectionId) | — | Close connection |
| **Host** | (port) | serverId or null | Start server |
| **HostTLS** | (port, certFile, keyFile) | serverId or null | Encrypted server |
| **Accept** | (serverId) | connectionId | Blocking accept |
| **AcceptTimeout** | (serverId, timeoutMs) | connectionId or null | Accept with timeout |
| **CloseServer** | (serverId) | — | Close server |
| **CreateRoom** | (roomId) | — | Ensure room exists |
| **JoinRoom** | (roomId, connectionId) | — | Add connection to room |
| **LeaveRoom** | (connectionId) | — | Remove from all rooms |
| **LeaveRoom** | (connectionId, roomId) | — | Remove from one room |
| **SendToRoom** | (roomId, text) | int | Send text to room; returns count sent |
| **SendToRoomJSON** | (roomId, jsonText) | int | Send JSON to room |
| **SendToRoomInt** | (roomId, value) | int | Broadcast int |
| **SendToRoomFloat** | (roomId, value) | int | Broadcast float |
| **SendToRoomNumbers** | (roomId, n1, n2, …) | int | Broadcast numbers |
| **GetRoomConnectionCount** | (roomId) | int | Connections in room |
| **GetRoomConnectionId** | (roomId, index) | connectionId or "" | Connection at index |
| **IsConnected** | (connectionId) | 1 or 0 | True if in conns |
| **GetConnectionCount** | () | int | Total connections |
| **GetLocalIP** | () | string | Local IP for LAN |

---

## 19. SQL – `sql.go`

See [docs/SQL.md](docs/SQL.md).

| Command | Arguments | Returns | Description |
|---------|-----------|---------|-------------|
| **OpenDatabase** | (path) | dbId or null | Open SQLite file (e.g. "game.db") |
| **CloseDatabase** | (dbId) | — | Close database |
| **Exec** | (dbId, sql) | int | Run INSERT/UPDATE/DELETE/DDL; rows affected or -1 |
| **ExecParams** | (dbId, sql, arg1, arg2, …) | int | Same with ? placeholders |
| **Query** | (dbId, sql) | int | Run SELECT; row count or -1; use GetRowCount/GetCell |
| **QueryParams** | (dbId, sql, arg1, arg2, …) | int | Parameterized SELECT |
| **GetRowCount** | () | int | Rows in last query result |
| **GetColumnCount** | () | int | Columns |
| **GetColumnName** | (colIndex) | string | Column name (0-based) |
| **GetCell** | (row, col) | value or null | Value at (row, col); 0-based |
| **Begin** | (dbId) | 1 or 0 | Start transaction |
| **Commit** | (dbId) | 1 or 0 | Commit |
| **Rollback** | (dbId) | 1 or 0 | Rollback |
| **LastError** | () | string | Last error message |

---

## 20. UI – `raylib_ui.go` and full raygui – `raylib_raygui.go`

### Pure-Go layout (raylib_ui.go, no CGO)

BeginUI() resets cursor; widgets advance vertically. EndUI().

| Command | Arguments | Returns | Description |
|---------|-----------|---------|-------------|
| **BeginUI** | () | — | Start UI layout |
| **EndUI** | () | — | End UI layout |
| **Label** | (text) | — | Label widget |
| **Button** | (text) | bool | Button; true if clicked |
| **Slider** | (text, value, min, max) | float | Slider; returns value |
| **Checkbox** | (text, checked) | 1 or 0 | Checkbox |
| **TextBox** | (id, text) | string | Editable; use same id each frame |
| **Dropdown** | (id, itemsText, activeIndex) | int | itemsText = "A;B;C" |
| **ProgressBar** | (text, value, min, max) | float | Progress bar |
| **WindowBox** | (title) | — | Start window box |
| **EndWindowBox** | () | — | End window box |
| **GroupBox** | (text) | — | Start group |
| **EndGroupBox** | () | — | End group |

### Full raygui (raylib_raygui.go; requires CGO)

All coordinates and sizes in pixels. See [docs/GUI_GUIDE.md](docs/GUI_GUIDE.md).

| Command | Arguments | Returns | Description |
|---------|-----------|---------|-------------|
| **GuiLabel** | (x, y, w, h, text) | — | Label |
| **GuiButton** | (x, y, w, h, text) | 1 or 0 | 1 if clicked |
| **GuiCheckBox** | (x, y, w, h, text, checked) | 1 or 0 | Checkbox |
| **GuiSlider** | (x, y, w, h, textLeft, textRight, value, min, max) | float | Slider |
| **GuiProgressBar** | (x, y, w, h, textLeft, textRight, value, min, max) | float | Progress bar |
| **GuiTextBox** | (id, x, y, w, h, text) | string | Editable (id = cache key) |
| **GuiDropdownBox** | (id, x, y, w, h, itemsText, active) | int | itemsText e.g. "One;Two;Three" |
| **GuiWindowBox** | (x, y, w, h, title) | 1 or 0 | 1 if close clicked |
| **GuiGroupBox** | (x, y, w, h, text) | — | Group box |
| **GuiLine** | (x, y, w, h, text) | — | Line |
| **GuiPanel** | (x, y, w, h, text) | — | Panel |

---

## 21. Language and built-ins

| Concept | Description |
|--------|-------------|
| **Function/Sub** | Define with parameters and Return; call by name. |
| **Module** | `Module Name … End Module`; body is Function/Sub only; call as ModuleName.FunctionName(...). |
| **Single-line IF** | Consecutive `IF condition THEN statement` lines do not require ENDIF between them. Use ENDIF when the next line is not another IF. |
| **On KeyDown / On KeyPressed** | `On KeyDown("KEY") … End On`; handlers run when PollInputEvents() is called. Key names: "ESCAPE", "W", "SPACE", or KEY_* constants. |
| **StartCoroutine** | (subName) — start fiber at that sub. |
| **Yield** | — switch to next fiber. |
| **WaitSeconds** | (seconds) — yield current fiber for N seconds (non-blocking; other fibers keep running). |

Fibers share the same chunk; each has its own IP, stack, and call stack.

---

## 22. Multi-window (in-process) – `raylib_multiwindow.go`

Logical windows (viewports) in one process; ID 0 = main screen. See [docs/MULTI_WINDOW_INPROCESS.md](docs/MULTI_WINDOW_INPROCESS.md).

| Command | Arguments | Returns | Description |
|---------|-----------|---------|-------------|
| **WindowCreate** | (width, height, title) | id | Create window |
| **WindowCreatePopup** | (…) | id | Create popup |
| **WindowCreateModal** | (…) | id | Create modal |
| **WindowCreateToolWindow** | (width, height, title) | id | Create tool window |
| **WindowCreateChild** | (parentID, width, height, title) | id | Create child |
| **WindowClose** | (id) | — | Close window |
| **WindowIsOpen** | (id) | bool | True if open |
| **WindowSetTitle** | (id, title) | — | Set title |
| **WindowSetSize** | (id, width, height) | — | Set size |
| **WindowSetPosition** | (id, x, y) | — | Set position |
| **WindowGetWidth** | (id) | int | Width |
| **WindowGetHeight** | (id) | int | Height |
| **WindowGetPositionX/Y** | (id) | int | Position |
| **WindowFocus** | (id) | — | Focus window |
| **WindowIsFocused** | (id) | bool | True if focused |
| **WindowIsVisible** | (id) | bool | True if visible |
| **WindowShow** | (id) | — | Show |
| **WindowHide** | (id) | — | Hide |
| **WindowBeginDrawing** | (id) | — | Begin draw to window |
| **WindowEndDrawing** | (id) | — | End draw |
| **WindowClearBackground** | (id, r, g, b, a) | — | Clear window |
| **WindowDrawAllToScreen** | () | — | Draw all windows to screen |
| **WindowSendMessage** | (targetID, message, data) | — | Send message |
| **WindowBroadcast** | (message, data) | — | Broadcast to all |
| **WindowReceiveMessage** | (id) | "message\|data" or null | Receive message |
| **WindowHasMessage** | (id) | bool | True if message queued |
| **ChannelCreate** | (name) | — | Create channel |
| **ChannelSend** | (name, data) | — | Send to channel |
| **ChannelReceive** | (name) | value or null | Receive from channel |
| **ChannelHasData** | (name) | bool | True if data |
| **StateSet** | (key, value) | — | Set state |
| **StateGet** | (key) | value | Get state |
| **StateHas** | (key) | bool | True if key exists |
| **StateRemove** | (key) | — | Remove key |
| **OnWindowUpdate** | (id, subName) | — | Register update callback |
| **OnWindowDraw** | (id, subName) | — | Register draw callback |
| **OnWindowResize** | (id, subName) | — | Register resize callback |
| **OnWindowClose** | (id, subName) | — | Register close callback |
| **OnWindowMessage** | (id, subName) | — | Register message callback |
| **WindowProcessEvents** | () | — | Process window events |
| **WindowDraw** | () | — | Draw all windows |
| **WindowSetCamera** | (id, cameraId) | — | Set window camera |
| **WindowDrawModel** | (id, modelId, x, y, z, scale [, r,g,b,a]) | — | Draw model in window |
| **WindowDrawScene** | (id, sceneId) | — | Draw scene |
| **WindowRegisterFunction** | (windowId, name, subName) | — | Register RPC function |
| **WindowCall** | (targetWindowId, name, arg1, arg2, …) | — | Call window function |
| **DockCreateArea** | (id, x, y, width, height) | — | Create dock area |
| **DockSplit** | (areaId, direction, size) | — | Split dock |
| **DockAttachWindow** | (areaId, windowId) | — | Attach window to dock |
| **DockSetSize** | (nodeId, size) | — | Set dock node size |

---

## Terrain – `compiler/bindings/terrain`

| Command | Arguments | Returns | Description |
|---------|-----------|---------|-------------|
| **LoadHeightmap** | (imageId) | heightmap id | Create heightmap from image |
| **GenHeightmap** | (width, depth, noiseScale) | heightmap id | Procedural heightmap |
| **GenHeightmapPerlin** | (width, depth, offsetX, offsetY, scale) | heightmap id | Perlin heightmap |
| **GenTerrainMesh** | (heightmapId, sizeX, sizeZ, heightScale [, lod]) | mesh id | Build terrain mesh |
| **TerrainCreate** | (heightmapId, sizeX, sizeZ, heightScale) | terrain id | Create terrain |
| **TerrainUpdate** | (terrainId) | — | Rebuild mesh |
| **DrawTerrain** | (terrainId, posX, posY, posZ) | — | Draw terrain (Render3D) |
| **SetTerrainMaterial** / **SetTerrainTexture** | (terrainId, id) | — | Set material/texture |
| **SetTerrainLOD** | (terrainId, lodLevel) | — | Set LOD |
| **TerrainRaise** / **TerrainLower** / **TerrainSmooth** / **TerrainFlatten** / **TerrainPaint** | (terrainId, x, z, radius, …) | — | Edit heightmap |
| **TerrainGetHeight** | (terrainId, x, z) | float | Height at (x,z) |
| **TerrainGetNormal** | (terrainId, x, z) | [nx, ny, nz] | Normal at (x,z) |
| **TerrainRaycast** | (terrainId, ox, oy, oz, dx, dy, dz) | [hit, dist, hx, hy, hz] | Ray vs terrain |
| **TerrainEnableCollision** / **TerrainSetFriction** / **TerrainSetBounce** | (terrainId, …) | — | Physics (state stored) |

---

## Water – `compiler/bindings/water`

| Command | Arguments | Returns | Description |
|---------|-----------|---------|-------------|
| **WaterCreate** | (width, depth, tileSize) | water id | Create water plane |
| **DrawWater** | (waterId [, posX, posY, posZ]) | — | Draw water (Render3D) |
| **SetWaterPosition** / **SetWaterWaveSpeed** / **SetWaterWaveHeight** / **SetWaterWaveFrequency** / **SetWaterTime** | (waterId, …) | — | State |
| **WaterGetHeight** | (waterId, x, z) | float | Wave height at (x,z) |
| **SetWaterTexture** / **SetWaterReflectionTexture** / **SetWaterRefractionTexture** / **SetWaterNormalMap** / **SetWaterColor** / **SetWaterShininess** | (waterId, …) | — | Rendering params |
| **WaterEnableFoam** / **WaterSetFoamIntensity** / **WaterSetDepthFade** / **WaterSetTransparency** | (waterId, …) | — | Advanced |
| **WaterSetDensity** / **WaterSetDrag** | (waterId, …) | — | Physics (buoyancy) |
| **WaterApplyBuoyancy** | (bodyId, waterId) | — | Apply buoyancy (stub) |

---

## Vegetation – `compiler/bindings/vegetation`

| Command | Arguments | Returns | Description |
|---------|-----------|---------|-------------|
| **TreeTypeCreate** | (modelId, trunkTexId, leafTexId) | type id | Tree type |
| **TreeSystemCreate** | () | system id | Tree system |
| **TreePlace** | (systemId, typeId, x, y, z, scale, rotation) | tree id | Place tree |
| **TreeRemove** / **TreeSetPosition** / **TreeSetScale** / **TreeSetRotation** | (treeId, …) | — | Edit tree |
| **TreeSystemSetLOD** / **TreeSystemEnableInstancing** | (systemId, …) | — | LOD/instancing |
| **TreeGetAt** | (systemId, x, z) | tree id or "" | Nearest tree |
| **DrawTrees** | (systemId) | — | Draw trees (Render3D) |
| **GrassCreate** | (textureId, density, patchSize) | grass id | Grass system |
| **GrassSetWind** / **GrassSetHeight** / **GrassSetColor** / **GrassPaint** / **GrassErase** / **GrassSetDensity** / **GrassSetLOD** / **GrassEnableInstancing** | (grassId, …) | — | Grass state |
| **GrassSetBendAmount** / **GrassSetInteraction** | (grassId, …) | — | Wind bend, interaction |
| **DrawGrass** | (grassId) | — | Draw grass (Render3D) |
| **TreeEnableCollision** / **TreeSetCollisionRadius** / **TreeSetWind** / **TreeApplyWind** / **TreeRaycast** | (…) | — | Tree physics (stubs) |

---

## Objects – `compiler/bindings/objects`

| Command | Arguments | Returns | Description |
|---------|-----------|---------|-------------|
| **ObjectPlace** | (modelId, x, y, z, scale, rotation) | object id | Place object |
| **ObjectRemove** / **ObjectSetTransform** | (objectId, …) | — | Edit/remove |
| **ObjectRandomScatter** | (modelId, areaX, areaZ, count, minScale, maxScale) | [id,…] | Scatter |
| **ObjectPaint** / **ObjectErase** | (…) | — | Paint/erase by area |
| **ObjectGetAt** | (x, z) | object id | Nearest at (x,z) |
| **ObjectRaycast** | (ox, oy, oz, dx, dy, dz) | [hit, objectId, hx, hy, hz] | Ray vs objects |
| **DrawObject** / **DrawAllObjects** | (objectId) / () | — | Draw |

---

## World – `compiler/bindings/world`

| Command | Arguments | Returns | Description |
|---------|-----------|---------|-------------|
| **WorldSave** / **WorldLoad** | (path) | — | Save/load world (e.g. objects) |
| **WorldExportJSON** / **WorldImportJSON** | (path) | — | Export/import JSON |
| **WorldStreamEnable** / **WorldStreamSetRadius** / **WorldStreamSetCenter** | (…) | — | Chunk streaming (stubs) |
| **WorldLoadChunk** / **WorldUnloadChunk** / **WorldIsChunkLoaded** / **WorldGetLoadedChunks** | (…) | — | Chunk API (stubs) |

---

## Navigation – `compiler/bindings/navigation`

| Command | Arguments | Returns | Description |
|---------|-----------|---------|-------------|
| **NavGridCreate** / **NavGridSetWalkable** / **NavGridSetCost** / **NavGridFindPath** | (…) | gridId / path | Grid pathfinding (stubs) |
| **NavMeshCreateFromTerrain** / **NavMeshAddObstacle** / **NavMeshRemoveObstacle** / **NavMeshFindPath** | (…) | meshId / path | NavMesh (stubs) |
| **NavAgentCreate** / **NavAgentSetSpeed** / **NavAgentSetRadius** / **NavAgentSetDestination** / **NavAgentGetNextWaypoint** | (…) | agentId / waypoint | Agents (stubs) |

---

## Indoor – `compiler/bindings/indoor`

| Command | Arguments | Returns | Description |
|---------|-----------|---------|-------------|
| **RoomCreate** / **RoomSetBounds** / **RoomAddPortal** | (…) | roomId | Rooms (stubs) |
| **PortalCreate** / **PortalSetOpen** | (…) | portalId | Portals (stubs) |
| **DoorCreate** / **DoorSetOpen** / **DoorToggle** / **DoorSetLocked** | (…) | doorId | Doors (stubs) |
| **LeverCreate** / **ButtonCreate** / **SwitchCreate** / **TriggerCreate** / **InteractableCreate** / **PickupCreate** / **LightZoneCreate** | (…) | id | Interaction (stubs) |
| **WorldSaveInteractables** / **WorldLoadInteractables** | (path) | — | Save/load (stubs) |

---

## Procedural – `compiler/bindings/procedural`

| Command | Arguments | Returns | Description |
|---------|-----------|---------|-------------|
| **NoisePerlin2D** | (x, y, scale) | float | Perlin-style noise [0,1] |
| **NoiseFractal2D** | (x, y, octaves, persistence, lacunarity) | float | Fractal noise |
| **NoiseSimplex2D** | (x, y, scale) | float | Simplex-style noise |
| **ScatterTrees** | (treeSystemId, treeTypeId, areaX, areaZ, density) | — | Scatter trees |
| **ScatterGrass** | (grassId, centerX, centerZ, radius, density) | — | Scatter grass |
| **ScatterObjects** | (modelId, areaX, areaZ, count [, minScale, maxScale]) | — | Scatter objects |

---

## Optimization (raylib 3D)

| Command | Arguments | Returns | Description |
|---------|-----------|---------|-------------|
| **SetCullingDistance** | (distance) | — | Max draw distance |
| **EnableFrustumCulling** | (flag) | — | Toggle frustum culling |

---

## Notes

- **Conventions:** All names are case-insensitive. Flat names (e.g. InitWindow) and optional namespaces (e.g. RL.InitWindow, BOX2D.*) are supported.
- **Alternate names / Aliases:** Many commands have aliases for familiarity: IsKeyDown/KeyDown, IsKeyPressed/KeyPressed, DrawRect/DrawRectangleLines, DrawRectFill/DrawRectangle, DrawCircleFill/DrawCircle, BeginCamera2D/BeginMode2D, EndCamera2D/EndMode2D, BeginCamera3D/BeginMode3D, EndCamera3D/EndMode3D, UILabel/Label, UIButton/Button, UpdateMusic/UpdateMusicStream, UnloadMusic/UnloadMusicStream, SetImageColor/ImageDrawPixel. See [Command Reference](docs/COMMAND_REFERENCE.md) for the full list.
- **Constant limit:** Bytecode uses 1-byte constant indices; each chunk supports at most 256 constants. Programs that exceed this (e.g. very large literals or many distinct identifiers) will fail at compile time with "too many constants".
- **Resource IDs:** Load* functions (LoadImage, LoadTexture, LoadSound, LoadMusicStream, LoadWave, LoadFont, LoadModel, LoadMesh, LoadShader, LoadRenderTexture, LoadAudioStream, etc.) return string IDs (e.g. `img_1`, `sound_1`). Pass these IDs to the matching Unload* and other APIs.
- **Vectors/Matrix/Quaternion:** Pass as flat numbers: Vector2 (x,y), Vector3 (x,y,z), Matrix (16 floats row-major), Quaternion (x,y,z,w). Functions that return vectors/matrices return a list (e.g. [x,y] or 16 values).
- **Colors:** Pass as (r,g,b,a) or use constants (White, Red, etc. – return packed int). NewColor(r,g,b,a) returns packed int.
- **Export*AsCode:** ExportImageAsCode(imageId, fileName), ExportFontAsCode(fontId, fileName), ExportWaveAsCode(waveId, fileName) write C header files; return true on success.
- **Audio callbacks:** SetAudioStreamCallback, AttachAudioStreamProcessor, AttachAudioMixedProcessor (and Detach*) are registered but return an error from BASIC; use UpdateAudioStream to push samples instead.

---

## See also

- [Documentation Index](docs/DOCUMENTATION_INDEX.md)
- [Command Reference](docs/COMMAND_REFERENCE.md) – commands grouped by feature
- [Getting Started](docs/GETTING_STARTED.md)
