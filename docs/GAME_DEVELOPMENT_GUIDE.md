# Game Development Guide

Complete guide to making games with CyberBasic: game loop, input, GAME.* helpers, 2D/3D physics, ECS, and best practices.

## Table of Contents

1. [Game loop pattern](#game-loop-pattern)
2. [Input](#input)
3. [GAME.* helpers](#game-helpers)
4. [2D physics (Box2D)](#2d-physics-box2d)
5. [3D physics (Bullet)](#3d-physics-bullet)
6. [Collision callbacks (2D)](#collision-callbacks-2d)
7. [ECS (Entity-Component System)](#ecs-entity-component-system)
8. [Quality of life](#quality-of-life)
9. [Best practices](#best-practices)
10. [Physics stubs](#physics-stubs)

---

## Game loop pattern

Every game uses the same structure. Prefer **Main() ... EndMain** and **DeltaTime()**:

```basic
InitWindow(800, 600, "My Game")
SetTargetFPS(60)

Main()
    // 1. Delta time (for movement and physics)
    VAR dt = DeltaTime()
    IF dt > 0.05 THEN LET dt = 0.016   // Clamp for physics stability

    // 2. Input and game logic
    // ...

    // 3. Physics step (if using Box2D or Bullet)
    // BOX2D.Step(...) or BULLET.Step(...)

    // 4. Draw (no BeginDrawing/EndDrawing needed inside Main)
    ClearBackground(20, 20, 30, 255)
    // Draw 2D/3D and UI
EndMain

CloseWindow()
```

You can also use `WHILE NOT WindowShouldClose() ... WEND`; both forms are automatically wrapped with BeginDrawing/EndDrawing.

See [2D Graphics Guide](2D_GRAPHICS_GUIDE.md) and [3D Graphics Guide](3D_GRAPHICS_GUIDE.md) for the 2D and 3D game checklists and full examples.

---

## Input

### Keyboard

```basic
IsKeyDown(KEY_W)      // Held this frame
IsKeyPressed(KEY_SPACE)   // Just pressed (one shot)
IsKeyReleased(KEY_ESCAPE)
```

Key constants: `KEY_W`, `KEY_A`, `KEY_S`, `KEY_D`, `KEY_SPACE`, `KEY_ESCAPE`, `KEY_UP`, `KEY_DOWN`, `KEY_LEFT`, `KEY_RIGHT`, `KEY_ENTER`, `KEY_F1` … `KEY_F12`, etc. See [API_REFERENCE.md](../API_REFERENCE.md).

### Movement axes (simple)

**GetAxisX()** and **GetAxisY()** return -1, 0, or 1 for A/D and W/S:

```basic
VAR speed = 4
LET x = x + speed * GetAxisX()
LET y = y + speed * GetAxisY()
```

### Mouse

```basic
GetMouseX()
GetMouseY()
GetMousePosition()   // Returns [x, y]
GetMouseDelta()      // Returns [dx, dy] for relative movement (e.g. camera)
IsMouseButtonDown(MOUSE_LEFT_BUTTON)
IsMouseButtonPressed(MOUSE_LEFT_BUTTON)
```

### Events (optional)

Register handlers and call **PollInputEvents()** in your loop:

```basic
ON KeyDown("ESCAPE")
    // handle escape
END ON

ON KeyPressed("SPACE")
    // fire once
END ON

WHILE NOT WindowShouldClose()
    PollInputEvents()
    // ...
WEND
```

---

## GAME.* helpers

These functions simplify camera and movement in 2D and 3D games.

### 2D camera follow

Set a target Box2D body; the 2D camera follows it each frame. In the loop, after BOX2D.Step call **GAME.UpdateCamera2D()**, then use BeginMode2D() / EndMode2D() as usual. Example:

```basic
GAME.SetCamera2DFollow(worldId, bodyId, xOffset, yOffset)
// In the loop, after BOX2D.Step:
GAME.UpdateCamera2D()
```

### 3D orbit camera

**GAME.CameraOrbit(targetX, targetY, targetZ, angleRad, pitchRad, distance)** sets the 3D camera to orbit around a point. Call it each frame (e.g. after getting the player position from Bullet), then **BeginDrawing()**, **BeginMode3D()**, draw, **EndMode3D()**, **EndDrawing()**.

### 3D WASD + jump (Bullet)

**GAME.MoveWASD(worldId, bodyId, angleRad, speed, jumpVel, dt)** applies horizontal force from WASD (relative to `angleRad`) and jump on Space when near the ground. Use with **BULLET.Step** and **GAME.CameraOrbit** for a full 3D character controller. See [templates/3d_game.bas](../templates/3d_game.bas).

### Other helpers

- **MoveHorizontal2D**, **Jump2D**, **IsOnGround2D** – 2D platformer movement
- **GAME.SyncSpriteToBody2D(worldId, bodyId, spriteId)** – sync a sprite position to a Box2D body (call in draw loop)
- **GAME.ClampDelta(maxDt)** – returns min(DeltaTime(), maxDt) (or min(GetFrameTime(), maxDt)) for stable physics
- **GAME.AssetPath(filename)** – returns `"assets/" + filename` (e.g. `LoadTexture(GAME.AssetPath("hero.png"))`)
- **GAME.ShowDebug()** or **ShowDebug(extraText)** – draw FPS (and optional second line)

---

## 2D physics (Box2D)

Create a world, add bodies, step each frame, read positions for drawing.

```basic
BOX2D.CreateWorld("w", 0, -10)   // gravity x, y
BOX2D.CreateBody("w", "player", 1)   // dynamic body
BOX2D.CreateBox2D("w", "player", 400, 300, 32, 32)   // box shape
BOX2D.CreateBox("w", "ground", 400, 550, 800, 20)   // static ground

WHILE NOT WindowShouldClose()
    VAR dt = GetFrameTime()
    BOX2D.Step("w", dt, 8, 3)   // velocity iters, position iters

    VAR px = BOX2D.GetPositionX("w", "player")
    VAR py = BOX2D.GetPositionY("w", "player")

    BeginDrawing()
    ClearBackground(30, 30, 50, 255)
    DrawRectangle(px - 16, py - 16, 32, 32, 255, 100, 100, 255)
    EndDrawing()
WEND

BOX2D.DestroyWorld("w")
```

See [examples/box2d_demo.bas](../examples/box2d_demo.bas) and [API_REFERENCE.md](../API_REFERENCE.md) for all BOX2D.* functions.

---

## 3D physics (Bullet)

Create a world, add bodies (sphere, box, etc.), step each frame, read positions for drawing and for **GAME.CameraOrbit** / **GAME.MoveWASD**.

```basic
BULLET.CreateWorld("w", 0, -18, 0)   // gravity x, y, z
BULLET.CreateSphere("w", "player", 0, 0.5, 0, 0.5, 1)   // pos x,y,z, radius, mass
BULLET.CreateBox("w", "ground", 0, -0.5, 0, 12.5, 0.5, 12.5, 0)   // static

WHILE NOT WindowShouldClose()
    VAR dt = GetFrameTime()
    IF dt > 0.05 THEN LET dt = 0.016
    BULLET.Step("w", dt)
    GAME.MoveWASD("w", "player", camAngle, 120, 9, dt)

    VAR px = BULLET.GetPositionX("w", "player")
    VAR py = BULLET.GetPositionY("w", "player")
    VAR pz = BULLET.GetPositionZ("w", "player")
    GAME.CameraOrbit(px, py + 1.5, pz, camAngle, 0.2, 10)

    BeginDrawing()
    ClearBackground(50, 50, 60, 255)
    BeginMode3D()
    DrawCube(0, -0.5, 0, 25, 1, 25, 0, 128, 0, 255)
    DrawSphere(px, py, pz, 0.5, 255, 0, 0, 255)
    EndMode3D()
    EndDrawing()
WEND

BULLET.DestroyWorld("w")
```

See [templates/3d_game.bas](../templates/3d_game.bas), [examples/run_3d_physics_demo.bas](../examples/run_3d_physics_demo.bas), and [API_REFERENCE.md](../API_REFERENCE.md).

---

## Collision callbacks (2D)

When a Box2D body collides, you can call a BASIC sub by name:

```basic
GAME.SetCollisionHandler(bodyId, "OnHit")
// After BOX2D.Step("w", dt, 8, 3):
GAME.ProcessCollisions2D("w")

SUB OnHit(otherBodyId)
    // Handle collision with otherBodyId
END SUB
```

---

## ECS (Entity-Component System)

CyberBasic provides ECS via a library binding. You create worlds, entities, add/get/remove components, and query entities by component type. Use it for organizing many game objects (enemies, projectiles, props) with shared logic.

- **Create world:** `ECS.CreateWorld()` → worldId  
- **Create entity:** `ECS.CreateEntity(worldId)` → entityId  
- **Add/Get/Remove component:** `ECS.AddComponent(worldId, entityId, componentType, ...args)`, `ECS.GetComponent(...)`, `ECS.HasComponent`, `ECS.RemoveComponent`  
- **Query:** `ECS.Query(worldId, componentType1, ...)` → list of entity IDs to loop over  
- **Destroy:** `ECS.DestroyEntity(worldId, entityId)`, `ECS.DestroyWorld(worldId)`

See **[ECS Guide](ECS_GUIDE.md)** for the full API and [examples/ecs_demo.bas](../examples/ecs_demo.bas) for a runnable example.

---

## Quality of life

- **GAME.AssetPath("hero.png")** → `"assets/hero.png"` so you can keep art in an `assets/` folder.
- **GAME.ClampDelta(0.05)** → use instead of raw `DeltaTime()`/`GetFrameTime()` when stepping physics to avoid large steps after a stall.
- **GAME.ShowDebug()** or **ShowDebug("extra line")** – draw FPS (and optional text) for quick debugging.

---

## Best practices

1. **Use CONST for configuration** – e.g. `CONST ScreenW = 800`, `CONST FPS = 60`, `CONST Gravity = -18`.
2. **Clamp delta time** – `IF dt > 0.05 THEN LET dt = 0.016` (or use `GAME.ClampDelta(0.05)`) before physics step.
3. **Load resources once** – Load textures, fonts, models at startup; unload in cleanup.
4. **Use GetAxisX/GetAxisY** for simple 2D movement.
5. **Organize with functions** – e.g. `UpdatePlayer()`, `DrawPlayer()`, `UpdateEnemies()` called from the main loop.
6. **Use BeginDrawing/EndDrawing** – All drawing between `BeginDrawing()` and `EndDrawing()` each frame.

---

## Physics stubs

Some Box2D and Bullet APIs are **stubbed (no-op)** until implemented:

- **Box2D:** Most joints (Revolute, Prismatic, Pulley, Gear, Weld, Rope, Wheel, SetJointLimits2D, SetJointMotor2D) are stubs. **CreateDistanceJoint2D** is implemented.
- **Bullet:** SetFriction3D, SetRestitution3D, SetDamping3D, SetKinematic3D, SetGravity3D, SetLinearFactor3D, SetAngularFactor3D, SetCCD3D, and joint creation (CreateHingeJoint3D, CreateSliderJoint3D, etc.) and SetJointLimits3D / SetJointMotor3D are stubs.

Core features (world, bodies, step, position/velocity, apply force, raycast) work. See [API_REFERENCE.md](../API_REFERENCE.md) for the full list.

---

For more examples see [examples/README.md](../examples/README.md) and the [Documentation Index](DOCUMENTATION_INDEX.md).
