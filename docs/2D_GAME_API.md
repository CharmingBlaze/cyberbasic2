# 2D Game API Reference

Complete reference for the CyberBasic 2D game API: platformers, RPGs, shooters, UI, and multiplayer. All commands use PascalCase and integer IDs where applicable.

**Multiplayer-safe?** Commands marked with ✓ are safe to use in networked games (no server state or deterministic). Commands that modify shared state may need replication.

---

## 1. Overview

The 2D API is organized in `compiler/bindings/dbp/dbp_2d.go` and registered via `Register2D()` after `game.RegisterGame()` so tilemap commands can override game's defaults.

**Core modules:**
- Drawing (pixels, lines, shapes)
- Sprites (LoadImage/Sprite, rotated, scaled, tinted)
- Spritesheets/Animation
- Tilemaps
- 2D Camera
- 2D Collision
- 2D Physics (Box2D)
- SpriteObject2D (position + rotation + scale)
- Text, UI, Math, Particles2D
- Tasks/Coroutines

---

## 2. 2D Drawing Basics

| Command | Args | Description |
|---------|------|-------------|
| `DrawPixel` | (x, y, r, g, b) | Draw a single pixel |
| `DrawLine` | (x1, y1, x2, y2, r, g, b) | Draw a line (dbp.go) |
| `DrawRect` | (x, y, w, h, r, g, b) | Filled rectangle (dbp.go) |
| `DrawRectOutline` | (x, y, w, h, r, g, b) | Rectangle outline |
| `DrawCircle` | (x, y, radius, r, g, b) | Filled circle (dbp.go) |
| `DrawCircleOutline` | (x, y, radius, r, g, b) | Circle outline |
| `DrawTriangle` | (x1,y1, x2,y2, x3,y3, r,g,b) | Filled triangle |

**Multiplayer-safe?** ✓ Yes (draw-only).

**Example:**
```basic
Ink 255, 0, 0
DrawRect 100, 100, 50, 50
DrawCircleOutline 200, 200, 30
DrawRectOutline 250, 250, 40, 40
```

---

## 3. Sprites

| Command | Args | Description |
|---------|------|-------------|
| `LoadImage` | (path, id) | Load texture and store at integer id. Use LoadSprite as alias. |
| `LoadSprite` | (path, id) | Alias for LoadImage |
| `Sprite` | (id, x, y) | Draw sprite at position |
| `DeleteSprite` | (id) | Unload sprite from memory |
| `DrawSpriteRotated` | (id, x, y, angle) | Draw with rotation (degrees) |
| `DrawSpriteScaled` | (id, x, y, sx, sy) | Draw with scale |
| `DrawSpriteTint` | (id, x, y, r, g, b) | Draw with tint color |

**Multiplayer-safe?** ✓ Draw commands are safe. LoadImage/DeleteSprite are typically client-local.

**Example:**
```basic
LoadImage "player.png", 1
Sprite 1, 100, 200
DrawSpriteRotated 1, 150, 150, 45
DrawSpriteScaled 1, 200, 200, 2, 2
```

---

## 4. Spritesheets / Animation

| Command | Args | Description |
|---------|------|-------------|
| `LoadSpritesheet` | (id, path, frameW, frameH) | Load spritesheet with frame dimensions |
| `SetSpriteFrame` | (id, frame) | Set current frame index |
| `NextSpriteFrame` | (id) | Advance to next frame (wraps) |
| `DrawSpriteFrame` | (id, frame, x, y) | Draw specific frame |
| `AnimateSprite` | (id, startFrame, endFrame, speed) | Configure animation range and speed |

**Multiplayer-safe?** ✓ Yes (draw-only).

**Example:**
```basic
LoadSpritesheet 1, "hero.png", 32, 32
AnimateSprite 1, 0, 8, 10
' In draw loop:
NextSpriteFrame 1
DrawSpriteFrame 1, 0, 100, 100
' Or set frame explicitly:
SetSpriteFrame 1, 3
DrawSpriteFrame 1, 3, 100, 100
```

---

## 5. Tilemaps

| Command | Args | Description |
|---------|------|-------------|
| `LoadTilemap` | (id, path) | Load tilemap from file; maps id to internal map |
| `DrawTilemap` | (id) | Draw tilemap |
| `SetTile` | (id, x, y, tileIndex) | Set tile at grid position |
| `GetTile` | (id, x, y) | Get tile value (assign to variable) |

**Multiplayer-safe?** SetTile modifies shared state; use replication for sync.

**Example:**
```basic
LoadTilemap 1, "level1.png"
SetTile 1, 5, 5, 2
t = GetTile 1, 5, 5
DrawTilemap 1
```

---

## 6. 2D Camera

| Command | Args | Description |
|---------|------|-------------|
| `Camera2DOn` | () | Begin 2D camera mode |
| `Camera2DOff` | () | End 2D camera mode |
| `Camera2DPosition` | (x, y) | Set camera target position |
| `Camera2DZoom` | (value) | Set zoom level |
| `Camera2DRotation` | (angle) | Set rotation (degrees) |
| `Camera2DFollow` | (objectId) | Follow sprite object ID |

**Multiplayer-safe?** ✓ Yes (client-local view).

**Example:**
```basic
Camera2DOn
Camera2DPosition 400, 300
Camera2DZoom 1.5
Camera2DFollow 1
```

---

## 7. 2D Collision

| Command | Args | Description |
|---------|------|-------------|
| `RectCollides` | (x1,y1,w1,h1, x2,y2,w2,h2) | Returns true if rectangles overlap |
| `PointInRect` | (x,y, rx,ry,rw,rh) | Returns true if point inside rect |
| `CircleCollides` | (x1,y1,r1, x2,y2,r2) | Returns true if circles overlap |
| `PointInCircle` | (x,y, cx,cy,r) | Returns true if point inside circle |

**Multiplayer-safe?** ✓ Yes (pure computation).

**Example:**
```basic
If RectCollides(x1,y1,w1,h1, x2,y2,w2,h2) Then
  ' collision
End If
```

---

## 8. 2D Physics

Uses Box2D. Requires `Physics2DOn` before use.

| Command | Args | Description |
|---------|------|-------------|
| `Physics2DOn` | () | Enable 2D physics |
| `Physics2DOff` | () | Disable 2D physics |
| `MakeBody2D` | (id, mass) | Create rigid body (default size 1x1) |
| `MakeStatic2D` | (id) | Create static body |
| `SetBody2DPosition` | (id, x, y) | Set position |
| `SetBody2DVelocity` | (id, vx, vy) | Set velocity |
| `ApplyForce2D` | (id, fx, fy) | Apply force to center |
| `ApplyImpulse2D` | (id, ix, iy) | Apply impulse |
| `GetBody2DX` | (id) | Get X position |
| `GetBody2DY` | (id) | Get Y position |
| `GetBody2DVX` | (id) | Get velocity X |
| `GetBody2DVY` | (id) | Get velocity Y |

**Multiplayer-safe?** Physics state needs replication for sync.

**Example:**
```basic
Physics2DOn
MakeBody2D "player", 1
SetBody2DPosition "player", 100, 100
ApplyForce2D "player", 10, 0
x = GetBody2DX "player"
```

---

## 9. 2D Objects (SpriteObject2D)

Combine sprite + position/rotation/scale for easy game objects.

| Command | Args | Description |
|---------|------|-------------|
| `MakeSpriteObject` | (id, spriteId) | Create object with sprite |
| `PositionObject2D` | (id, x, y) | Set position |
| `MoveObject2D` | (id, dx, dy) | Add to position |
| `RotateObject2D` | (id, angle) | Set rotation (degrees) |
| `ScaleObject2D` | (id, sx, sy) | Set scale |
| `DrawObject2D` | (id) | Draw sprite at object transform |
| `SyncObject2D` | (id) | Mark for replication |

**Multiplayer-safe?** SyncObject2D registers for replication; use with ReplicatePosition.

**Example:**
```basic
LoadImage "ship.png", 1
MakeSpriteObject 1, 1
PositionObject2D 1, 200, 200
RotateObject2D 1, 90
DrawObject2D 1
SyncObject2D 1
```

---

## 10. Text

| Command | Args | Description |
|---------|------|-------------|
| `DrawText` | (text, x, y, size, r, g, b) | Draw text |
| `LoadFont` | (path, id) | Load font |
| `SetFont` | (id) | Set current font |

**Multiplayer-safe?** ✓ Yes.

---

## 11. UI

| Command | Args | Description |
|---------|------|-------------|
| `UIButton` | (id, x, y, w, h, text) | Button; returns true when clicked |
| `UILabel` | (id, x, y, text) | Label |
| `UICheckbox` | (id, x, y, text, checked) | Checkbox |
| `UISlider` | (id, x, y, w, h, text, value, min, max) | Slider |
| `UITextbox` | (id, x, y, w, h) | Editable text; returns current text |

**Multiplayer-safe?** ✓ Yes (client-local UI).

**Example:**
```basic
name = UITextbox "tb1", 100, 100, 200, 30
If UIButton "btn1", 100, 150, 100, 30, "Start" Then
  ' clicked
End If
```

---

## 12. 2D Math Helpers

| Command | Args | Description |
|---------|------|-------------|
| `AngleBetween2D` | (x1,y1, x2,y2) | Returns angle in radians (atan2) |
| `Distance2D` | (x1,y1, x2,y2) | Returns distance |
| `Normalize2D` | (x, y) | Returns [nx, ny] as array |
| `Dot2D` | (x1,y1, x2,y2) | Returns dot product |

**Multiplayer-safe?** ✓ Yes (pure computation).

**Example:**
```basic
angle = AngleBetween2D x1, y1, x2, y2
d = Distance2D x1, y1, x2, y2
vec = Normalize2D dx, dy
```

---

## 13. 2D Particles

| Command | Args | Description |
|---------|------|-------------|
| `MakeParticles2D` | (id, maxCount) | Create particle system |
| `SetParticles2DColor` | (id, r, g, b) | Set default color |
| `SetParticles2DSize` | (id, size) | Set particle size |
| `SetParticles2DSpeed` | (id, speed) | Set velocity magnitude |
| `EmitParticles2D` | (id, count [, x, y]) | Spawn particles |
| `DrawParticles2D` | (id) | Update and draw particles |

**Multiplayer-safe?** ✓ Yes (client-local effects).

**Example:**
```basic
MakeParticles2D 1, 100
SetParticles2DColor 1, 255, 0, 0
SetParticles2DSize 1, 4
EmitParticles2D 1, 10, 200, 200
DrawParticles2D 1
```

---

## 14. Tasks / Coroutines

| Command | Args | Description |
|---------|------|-------------|
| `StartTask` | (funcName) | Start coroutine |
| `StopTask` | (funcName) | Stop coroutine (stub) |
| `WaitSeconds` | (seconds) | Yield for N seconds |
| `WaitFrames` | (frames) | Yield for N frames |
| `Yield` | () | Yield one frame |

**Multiplayer-safe?** ✓ Yes (client-local scheduling).

**Example:**
```basic
StartTask "SpawnEnemies"
' ...
Sub SpawnEnemies
  For i = 0 To 10
    ' spawn
    WaitSeconds 1
  Next
End Sub
```

---

## 15. Multiplayer Essentials

For networked games, use these with the replication system:

- `SyncObject2D(id)` – Mark object for position sync
- `ReplicatePosition(entityId)` – Register entity for replication
- `ReplicateRotation(entityId)` – Register rotation
- `ReplicateScale(entityId)` – Register scale

See `docs/DBP_EXTENDED.md` and the replication module for full setup.

---

## Quick Reference

| Category | Key Commands |
|----------|--------------|
| Drawing | DrawPixel, DrawRect, DrawRectOutline, DrawCircle, DrawCircleOutline, DrawTriangle |
| Sprites | LoadImage, Sprite, DeleteSprite, DrawSpriteRotated, DrawSpriteScaled, DrawSpriteTint |
| Spritesheets | LoadSpritesheet, SetSpriteFrame, NextSpriteFrame, DrawSpriteFrame, AnimateSprite |
| Tilemaps | LoadTilemap, DrawTilemap, SetTile, GetTile |
| Camera | Camera2DOn, Camera2DOff, Camera2DPosition, Camera2DZoom, Camera2DFollow |
| Collision | RectCollides, PointInRect, CircleCollides, PointInCircle |
| Physics | Physics2DOn, MakeBody2D, SetBody2DPosition, ApplyForce2D, GetBody2DX/Y |
| Objects | MakeSpriteObject, PositionObject2D, DrawObject2D, SyncObject2D |
| UI | UITextbox, UIButton, UILabel, UICheckbox, UISlider |
| Math | AngleBetween2D, Distance2D, Normalize2D, Dot2D |
| Particles | MakeParticles2D, EmitParticles2D, DrawParticles2D |
| Tasks | StartTask, WaitSeconds, WaitFrames, Yield |
