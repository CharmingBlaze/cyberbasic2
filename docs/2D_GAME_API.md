# 2D Game API Reference

Complete reference for the CyberBASIC2 2D game API: platformers, RPGs, shooters, UI, and multiplayer. All commands use PascalCase and integer IDs where applicable.

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
| `SpriteExists` | (id) | Returns 1 if exists, 0 otherwise |
| `DrawSpriteRotated` | (id, x, y, angle) | Draw with rotation (degrees) |
| `DrawSpriteScaled` | (id, x, y, sx, sy) | Draw with scale |
| `DrawSpriteTint` | (id, x, y, r, g, b) | Draw with tint color |
| `SetSpriteColor` | (id, r, g, b, a) | Set persistent tint for sprite (used by Sprite, DrawSpriteRotated, DrawSpriteScaled) |

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
| `LoadSpritesheet` | (id, pngPath, jsonPath) or (id, path, frameW, frameH) | Load Aseprite (PNG+JSON) or grid spritesheet |
| `PlaySpriteAnimation` | (id, tagName, speed) | Play animation by tag (Aseprite) |
| `StopSpriteAnimation` | (id) | Stop sprite animation |
| `SetSpriteFrame` | (id, frame) | Set current frame index |
| `GetSpriteFrame` | (id) | Current frame index |
| `NextSpriteFrame` | (id) | Advance to next frame (wraps) |
| `DrawSpriteFrame` | (id, frame, x, y) | Draw specific frame |
| `AnimateSprite` | (id, startFrame, endFrame, speed) | Configure animation range (grid mode) |
| `GetSliceRect` | (id, sliceName) | Returns "x,y,w,h" for slice at current frame |
| `GetAnimationLength` | (id, tagName) | Frame count for tag |
| `DeleteSpritesheet` | (id) | Unload spritesheet |
| `CloneSpritesheet` | (newID, sourceID) | Duplicate spritesheet (shares texture) |
| `SpritesheetExists` | (id) | Returns 1 if exists, 0 otherwise |

**Aseprite workflow:** See [ASEPRITE_WORKFLOW.md](ASEPRITE_WORKFLOW.md).

**Example (Aseprite):**
```basic
LoadSpritesheet 1, "character.png", "character.json"
PlaySpriteAnimation 1, "walk", 1.0
' In draw loop (SYNC advances animation):
DrawSpriteFrame 1, GetSpriteFrame 1, 100, 100
```

**Example (grid):**
```basic
LoadSpritesheet 1, "hero.png", 32, 32
AnimateSprite 1, 0, 8, 10
NextSpriteFrame 1
DrawSpriteFrame 1, GetSpriteFrame 1, 100, 100
```

---

## 5. Tilemaps

| Command | Args | Description |
|---------|------|-------------|
| `LoadTilemap` | (id, path) | Load tilemap JSON from file; maps DBP id to internal map |
| `DrawTilemap` | (id) or (id, x, y) | Draw tilemap at its origin or with a draw offset |
| `SetTile` | (id, x, y, tileIndex) | Set tile at grid position |
| `GetTile` | (id, x, y) | Get tile value (assign to variable) |
| `TilemapSetTileset` | (mapId, texturePath) | Assign or replace the atlas texture used when drawing non-zero tiles |
| `DeleteTilemap` | (id) | Remove tilemap |
| `HideTilemap` | (id) | Set visible=false |
| `ShowTilemap` | (id) | Set visible=true |
| `TilemapExists` | (id) | Returns 1 if exists, 0 otherwise |

Shipping workflow:

- `LoadTilemap(id, path)` loads a JSON map, not a PNG image.
- The JSON supports `width`, `height`, `tiles`, `solid`, and either `tileSize` or `tileWidth` plus `tileHeight`.
- Optional `tileset` in the JSON, or `TilemapSetTileset(mapId, texturePath)` at runtime, enables atlas-based textured drawing.
- Tile value `0` is empty. Non-zero tiles draw from the atlas using `tileIndex - 1`.
- If no tileset is assigned, the renderer falls back to gray debug rectangles so maps still remain visible while prototyping.

**Multiplayer-safe?** `SetTile` modifies shared state. Use `SyncEntity` or your own messages for real networking; `Replicate*` markers are not an automatic replication layer today.

**Example:**
```basic
LoadTilemap 1, "levels/level1.json"
TilemapSetTileset "tm_1", "tiles/dungeon.png"
SetTile 1, 5, 5, 2
t = GetTile 1, 5, 5
DrawTilemap 1, 64, 32
```

Example tilemap JSON:

```json
{
  "tileWidth": 32,
  "tileHeight": 32,
  "width": 8,
  "height": 6,
  "tileset": "tiles/dungeon.png",
  "solid": [1, 2],
  "tiles": [
    [1,1,1,1,1,1,1,1],
    [1,0,0,0,0,0,0,1],
    [1,0,2,2,0,0,0,1],
    [1,0,0,0,0,3,0,1],
    [1,0,0,0,0,0,0,1],
    [1,1,1,1,1,1,1,1]
  ]
}
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

Uses the default Box2D world. These commands are simple DBPro-style wrappers over the authoritative 2D backend; use `Box2DBackendName()` and `Box2DBackendMode()` if you want to inspect the active backend at runtime. Requires `Physics2DOn` before use.

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
| `DeleteSpriteObject` | (id) | Remove sprite object |
| `HideSpriteObject` | (id) | Set visible=false |
| `ShowSpriteObject` | (id) | Set visible=true |
| `CloneSpriteObject` | (newID, sourceID) | Duplicate sprite object |
| `SpriteObjectExists` | (id) | Returns 1 if exists, 0 otherwise |

**Multiplayer-safe?** `SyncObject2D` is only a local marker. Use `SyncEntity`, RPC, or custom network messages for actual transport today.

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
| `DeleteParticles2D` | (id) | Remove particle system |
| `Particles2DExists` | (id) | Returns 1 if exists, 0 otherwise |

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

For networked games, use the explicit transport helpers:

- `SyncEntity(connectionId, entityId, x, y)` – Send a 2D position update directly
- `SyncEntityToRoom(roomId, entityId, x, y)` – Broadcast a 2D position update
- `SendRPC(connectionId, name, ...)` – Send gameplay events or commands
- `ReplicatePosition(entityId)` / `ReplicateRotation(entityId)` / `ReplicateScale(entityId)` – Marker-only today, not automatic state replication

See `docs/MULTIPLAYER.md` and `docs/MULTIPLAYER_DESIGN.md` for the current shipping model.

---

## Lifecycle Commands (Delete, Hide, Clone, Exists)

| Entity | Delete | Hide/Show | Clone | Exists |
|--------|--------|-----------|-------|--------|
| Sprite | DeleteSprite | - | - | SpriteExists |
| Spritesheet | DeleteSpritesheet | - | CloneSpritesheet | SpritesheetExists |
| Tilemap | DeleteTilemap | HideTilemap, ShowTilemap | - | TilemapExists |
| SpriteObject2D | DeleteSpriteObject | HideSpriteObject, ShowSpriteObject | CloneSpriteObject | SpriteObjectExists |
| Particles2D | DeleteParticles2D | - | - | Particles2DExists |

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

---

## See also

- [Core Command Reference](CORE_COMMAND_REFERENCE.md) – DBP-style command list
- [Aseprite Workflow](ASEPRITE_WORKFLOW.md) – Sprite sheet export with tags and slices
