# 2D Graphics Guide

Complete guide to 2D graphics in CyberBasic using the raylib API.

## Table of Contents

1. [Getting started](#getting-started)
2. [Drawing primitives](#drawing-primitives)
3. [Textures and images](#textures-and-images)
4. [Text rendering](#text-rendering)
5. [Colors](#colors)
6. [2D camera](#2d-camera)
7. [2D game checklist](#2d-game-checklist)
8. [Complete 2D game example](#complete-2d-game-example)
9. [Full 2D command reference](#full-2d-command-reference)
10. [See also](#see-also)

---

## Getting started

### Basic setup

Every 2D graphics program needs a window and a game loop:

Use **WHILE NOT WindowShouldClose() ... WEND** (or **REPEAT ... UNTIL WindowShouldClose()**) and **DeltaTime()**:

```basic
InitWindow(800, 600, "My 2D Game")
SetTargetFPS(60)

WHILE NOT WindowShouldClose()
    VAR dt = DeltaTime()
    ClearBackground(20, 20, 30, 255)
    // Your drawing code here; use dt for frame-based movement
WEND

CloseWindow()
```

The compiler does not inject any frame or mode calls; your loop compiles exactly as written (DBPro-style).

### Hybrid update/draw loop

When you define **update(dt)** and **draw()** and use a game loop, 2D draw calls inside **draw()** are queued and flushed automatically (2D, then 3D, then GUI). You do not need to call BeginDrawing/EndDrawing yourself. See [Program Structure](PROGRAM_STRUCTURE.md#hybrid-updatedraw-loop) and [Rendering and the game loop](RENDERING_AND_GAME_LOOP.md).

---

## Drawing primitives

### Rectangles

```basic
// Draw filled rectangle: (x, y, width, height, color)
// Color: either (r, g, b, a) or a packed color constant
DrawRectangle(100, 100, 200, 150, 255, 100, 100, 255)

// Draw rectangle outline
DrawRectangleLines(100, 100, 200, 150, 255, 255, 255, 255)

// Rounded rectangle
DrawRectangleRounded(100, 100, 200, 150, 0.2, 255, 100, 100, 255)
```

### Circles

```basic
// Draw filled circle: (centerX, centerY, radius, r, g, b, a)
DrawCircle(400, 300, 50, 100, 200, 255, 255)

// Draw circle outline
DrawCircleLines(400, 300, 50, 255, 255, 255, 255)
```

### Lines and triangles

```basic
// Draw line: (x1, y1, x2, y2, r, g, b, a)
DrawLine(100, 100, 500, 400, 255, 255, 255, 255)

// Line with thickness (DrawLineEx)
// Triangle: (x1,y1, x2,y2, x3,y3, color)
DrawTriangle(400, 100, 300, 200, 500, 200, 255, 100, 100, 255)
```

### Pixels

```basic
DrawPixel(100, 100, 255, 255, 255, 255)
```

---

## Textures and images

### Loading and drawing textures

```basic
// Load texture from file (returns texture id string)
VAR tex = LoadTexture("sprite.png")

// Draw texture at position (id, posX, posY) and optional tint (r,g,b,a)
DrawTexture(tex, 100, 100)
DrawTexture(tex, 200, 200, 255, 255, 255, 255)

// Draw with rotation and scale: (id, posX, posY, rotation, scale, tint...)
DrawTextureEx(tex, 100, 100, 45, 2.0, 255, 255, 255, 255)

// Draw part of texture (sprite sheet): (id, srcX, srcY, srcW, srcH, posX, posY, tint...)
DrawTextureRec(tex, 0, 0, 32, 32, 100, 100, 255, 255, 255, 255)

// Unload when done
UnloadTexture(tex)
```

### 2D sprite (texture) animation

Use **sprite-sheet** (texture frame) animation with time-based playback. One texture holds a grid of frames; an animation state advances by FPS and draws the current frame.

```basic
// Load texture and create animation: (textureId, frameWidth, frameHeight, framesPerRow [, totalFrames])
VAR tex = LoadTexture("hero_sheet.png")
VAR anim = CreateSpriteAnimation(tex, 32, 32, 4)
// Optional: SetSpriteAnimationFPS(anim, 8), SetSpriteAnimationLoop(anim, TRUE)

WHILE NOT WindowShouldClose()
    // Advance animation by delta time
    UpdateSpriteAnimation(anim, GetFrameTime())
    ClearBackground(20, 20, 30, 255)
    // Draw current frame at (x, y); optional scaleX, scaleY, rotation, tint
    DrawSpriteAnimation(anim, 100, 100)
    DrawSpriteAnimation(anim, 200, 200, 2.0, 2.0, 0, 255, 255, 255, 255)
WEND

DestroySpriteAnimation(anim)
UnloadTexture(tex)
```

**Commands:**

- **CreateSpriteAnimation**(textureId, frameWidth, frameHeight, framesPerRow [, totalFrames]) → animId. If totalFrames is omitted, it is derived from texture size.
- **SetSpriteAnimationFPS**(animId, fps) — playback speed (frames per second).
- **SetSpriteAnimationLoop**(animId, loop) — TRUE = loop, FALSE = one-shot.
- **SetSpriteAnimationFrame**(animId, frameIndex) — set current frame by index (0-based).
- **UpdateSpriteAnimation**(animId, deltaTime) — advance time; call each frame with GetFrameTime().
- **GetSpriteAnimationFrame**(animId) → current frame index (int).
- **DrawSpriteAnimation**(animId, posX, posY [, scaleX, scaleY, rotation, r, g, b, a]) — draw current frame. Position is top-left; rotation is around center of the frame.
- **DestroySpriteAnimation**(animId) — remove state (texture is not unloaded).

---

## Text rendering

### Basic text

```basic
// DrawText(text, x, y, fontSize) and optional (r, g, b, a)
DrawText("Hello, World!", 10, 10, 20)
DrawText("Colored", 10, 40, 20, 255, 255, 0, 255)

// Measure text width (pixels)
VAR w = MeasureText("Hello", 20)
```

### Custom font

```basic
VAR font = LoadFont("font.ttf")
// DrawTextEx(text, x, y, fontSize, spacing, r, g, b, a)
DrawTextEx("Hello", 10, 10, 20, 2, 255, 255, 255, 255)
// When done: UnloadFont(font)
```

---

## Colors

Pass colors as **r, g, b, a** (0–255). You can use raylib color constants (0-arg functions that return a packed color) where the API accepts a single color value:

- `White`, `Black`, `Red`, `Green`, `Blue`
- `LightGray`, `Gray`, `DarkGray`
- `Yellow`, `Gold`, `Orange`, `Pink`, `Maroon`
- `Lime`, `DarkGreen`, `SkyBlue`, `DarkBlue`, `Purple`, `Violet`, `Magenta`
- `RayWhite`, `Blank`, etc.

Example with constants (where the call accepts one color argument):

```basic
DrawRectangle(100, 100, 50, 50, Red)
DrawCircle(200, 200, 30, Green)
```

For calls that take separate r,g,b,a, use numbers: `255, 0, 0, 255` for red.

**NewColor(r, g, b, a)** returns a packed color integer for APIs that accept one color value.

---

## 2D camera

Use a 2D camera for scrolling or zoomed worlds:

```basic
// SetCamera2D(offsetX, offsetY, targetX, targetY, rotation, zoom)
SetCamera2D(400, 300, 400, 300, 0, 1.0)

ClearBackground(20, 20, 30, 255)
// Draw world content (world coordinates)
DrawRectangle(0, 0, 100, 100, 255, 100, 100, 255)
// UI in screen space
DrawText("HUD", 10, 10, 20, 255, 255, 255, 255)
```

**GetWorldToScreen2D(worldX, worldY)** and **GetScreenToWorld2D(screenX, screenY)** convert between world and screen coordinates when using the 2D camera.

---

## 2D game checklist

Use this checklist to confirm your program is a valid 2D game:

- [ ] **Window:** `InitWindow(width, height, title)` and `SetTargetFPS(60)` (or desired FPS)
- [ ] **Loop:** `WHILE NOT WindowShouldClose() ... WEND` (or `REPEAT ... UNTIL WindowShouldClose()`). No auto-wrap; your code compiles as written.
- [ ] **Clear:** `ClearBackground(r, g, b, a)` at the start of each frame
- [ ] **Input:** e.g. `IsKeyDown(KEY_W)`, `GetAxisX()`, `GetAxisY()`, `GetMouseX()`, `GetMouseY()`
- [ ] **Draw:** Primitives (`DrawRectangle`, `DrawCircle`, `DrawLine`, …) or textures (`LoadTexture`, `DrawTexture`, …) and/or `DrawText`
- [ ] **Close:** `CloseWindow()` after the loop

Optional for 2D games:

- **Physics:** `BOX2D.CreateWorld`, `BOX2D.CreateBody`, `BOX2D.Step`, etc. See [Game Development Guide](GAME_DEVELOPMENT_GUIDE.md) and [API_REFERENCE.md](../API_REFERENCE.md).
- **Camera follow:** `GAME.SetCamera2DFollow`, `GAME.UpdateCamera2D` with a Box2D body.

---

## Complete 2D game example

From [templates/2d_game.bas](../templates/2d_game.bas):

```basic
InitWindow(800, 600, "2D Game")
SetTargetFPS(60)

VAR x = 400
VAR y = 300
VAR speed = 4

WHILE NOT WindowShouldClose()
    LET x = x + speed * GetAxisX()
    LET y = y + speed * GetAxisY()

    ClearBackground(20, 20, 30, 255)
    DrawCircle(x, y, 30, 255, 100, 100, 255)
    DrawText("WASD to move", 10, 10, 20, 255, 255, 255, 255)
WEND

CloseWindow()
```

Run it: `cyberbasic templates/2d_game.bas`

---

## Full 2D command reference

All 2D-relevant commands in one place. See [API_REFERENCE.md](../API_REFERENCE.md) for details.

- **Window:** `InitWindow`, `SetTargetFPS`, `GetScreenWidth`, `GetScreenHeight`, `GetScreenCenterX`, `GetScreenCenterY`, `CloseWindow`, `WindowShouldClose`
- **Frame:** `ClearBackground` and your draw calls in the loop. The compiler does not inject any frame or mode calls.
- **Camera:** `SetCamera2D`, `SetCamera2DCenter`, `GetWorldToScreen2D`, `GetScreenToWorld2D`; `GAME.SetCamera2DFollow`, `GAME.UpdateCamera2D`
- **Primitives:** `DrawRectangle`, `DrawRectangleLines`, `DrawRectangleRounded`, `DrawCircle`, `DrawCircleLines`, `DrawLine`, `DrawLineEx`, `DrawTriangle`, `DrawPixel`
- **Textures:** `LoadTexture`, `DrawTexture`, `DrawTextureEx`, `DrawTextureRec`, `UnloadTexture`
- **Text:** `DrawText`, `DrawTextEx`, `MeasureText`, `LoadFont`, `UnloadFont`
- **Math / distance:** `Distance2D`, `Lerp`, `Vector2Lerp`, `Vector2MoveTowards`, `Vector2Distance`
- **Collision:** `CheckCollisionRecs`, `CheckCollisionPointRec`, `CheckCollisionCircles`, `GetCollisionRec` — see [API Reference](../API_REFERENCE.md) for arguments and returns.
- **Game loop / input:** `DeltaTime`, `WHILE NOT WindowShouldClose() … WEND`, `GetAxisX`, `GetAxisY`; `GAME.MoveHorizontal2D`, `GAME.Jump2D`, `GAME.SyncSpriteToBody2D`
- **Colors:** `NewColor`, and color constants (`Red`, `Green`, `White`, etc.)

---

## See also

- [API Reference](../API_REFERENCE.md) — full list of drawing and window functions
- [Game Development Guide](GAME_DEVELOPMENT_GUIDE.md) — 2D physics, GAME.* helpers
- [2D Physics Guide](2D_PHYSICS_GUIDE.md) — Box2D worlds, bodies, collision
- [Command Reference](COMMAND_REFERENCE.md) — commands by feature
- [examples/README.md](../examples/README.md) — more 2D examples
