# Game Development Guide

Complete guide to making games with CyberBasic: game loop, input, GAME.* helpers, 2D/3D physics, ECS, and best practices.

## Table of Contents

1. [Game loop pattern](#game-loop-pattern)
2. [Hybrid update/draw loop](#hybrid-updatedraw-loop)
3. [Input](#input)
4. [GAME.* helpers](#game-helpers)
5. [2D physics (Box2D)](#2d-physics-box2d)
6. [3D physics (Bullet)](#3d-physics-bullet)
7. [Collision callbacks (2D)](#collision-callbacks-2d)
8. [ECS (Entity-Component System)](#ecs-entity-component-system)
9. [Quality of life](#quality-of-life)
10. [State machines](#state-machines)
11. [Best practices](#best-practices)
12. [Physics implementation status](#physics-implementation-status)
13. [2D and 3D quick reference](#2d-and-3d-quick-reference)

---

## Game loop pattern

Every game uses the same structure. Use **WHILE NOT WindowShouldClose() ... WEND** (or **REPEAT ... UNTIL WindowShouldClose()**) and **DeltaTime()**:

```basic
InitWindow(800, 600, "My Game")
SetTargetFPS(60)

WHILE NOT WindowShouldClose()
    // 1. Delta time (for movement and physics)
    VAR dt = DeltaTime()
    IF dt > 0.05 THEN LET dt = 0.016   // Clamp for physics stability

    // 2. Input and game logic
    // ...

    // 3. Physics step (if using Box2D or Bullet)
    // Step2D(...) or Step3D(...)

    // 4. Draw
    ClearBackground(20, 20, 30, 255)
    // Draw 2D/3D and UI
WEND

CloseWindow()
```

The compiler does not inject any frame or mode calls; your loop compiles as written (DBPro-style). Exception: the **hybrid loop** (see below).

See [2D Graphics Guide](2D_GRAPHICS_GUIDE.md) and [3D Graphics Guide](3D_GRAPHICS_GUIDE.md) for the 2D and 3D game checklists and full examples.

### Hybrid update/draw loop

If you define **update(dt)** and **draw()** (as Sub or Function) and use a game loop (`WHILE NOT WindowShouldClose()` or `REPEAT ... UNTIL WindowShouldClose()`), the compiler **replaces the loop body** with an automatic pipeline: GetFrameTime → StepAllPhysics2D(dt), StepAllPhysics3D(dt) → update(dt) → ClearRenderQueues → draw() (Draw*/Gui* calls are queued) → FlushRenderQueues. You do not call BeginDrawing/EndDrawing or BeginMode2D/EndMode3D yourself. Prefer this for new games when you want a clear update/draw split and automatic physics stepping. See [Program Structure](PROGRAM_STRUCTURE.md#hybrid-updatedraw-loop) and **examples/hybrid_update_draw_demo.bas**.

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

Set a target Box2D body; the 2D camera follows it each frame. In the loop, after Step2D call **GAME.UpdateCamera2D()**. Example:

```basic
GAME.SetCamera2DFollow(worldId, bodyId, xOffset, yOffset)
// In the loop, after Step2D:
GAME.UpdateCamera2D()
```

### 3D orbit camera

**GAME.CameraOrbit(targetX, targetY, targetZ, angleRad, pitchRad, distance)** sets the 3D camera to orbit around a point. Call it each frame (e.g. after getting the player position from Bullet), then draw.

### 3D WASD + jump (Bullet)

**GAME.MoveWASD(worldId, bodyId, angleRad, speed, jumpVel, dt)** applies horizontal force from WASD (relative to `angleRad`) and jump on Space when near the ground. Use with **Step3D** and **GAME.CameraOrbit** for a full 3D character controller. See [templates/3d_game.bas](../templates/3d_game.bas).

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
CreateWorld2D("w", 0, -10)   // gravity x, y
CreateBox2D("w", "player", 400, 300, 32, 32, 1, 1)   // dynamic box
CreateBox2D("w", "ground", 400, 550, 800, 20, 0, 0)   // static ground

WHILE NOT WindowShouldClose()
    VAR dt = GetFrameTime()
    Step2D("w", dt, 8, 3)   // velocity iters, position iters

    VAR px = GetPositionX2D("w", "player")
    VAR py = GetPositionY2D("w", "player")

    ClearBackground(30, 30, 50, 255)
    DrawRectangle(px - 16, py - 16, 32, 32, 255, 100, 100, 255)
WEND

DestroyWorld2D("w")
```

See [examples/box2d_demo.bas](../examples/box2d_demo.bas) and [API_REFERENCE.md](../API_REFERENCE.md) for all 2D physics (flat) functions.

---

## 3D physics (Bullet)

Create a world, add bodies (sphere, box, etc.), step each frame, read positions for drawing and for **GAME.CameraOrbit** / **GAME.MoveWASD**.

```basic
CreateWorld3D("w", 0, -18, 0)   // gravity x, y, z
CreateSphere3D("w", "player", 0, 0.5, 0, 0.5, 1)   // pos x,y,z, radius, mass
CreateBox3D("w", "ground", 0, -0.5, 0, 12.5, 0.5, 12.5, 0)   // static

WHILE NOT WindowShouldClose()
    VAR dt = GetFrameTime()
    IF dt > 0.05 THEN LET dt = 0.016
    BULLET.Step("w", dt)
    GAME.MoveWASD("w", "player", camAngle, 120, 9, dt)

    VAR px = GetPositionX3D("w", "player")
    VAR py = GetPositionY3D("w", "player")
    VAR pz = GetPositionZ3D("w", "player")
    GAME.CameraOrbit(px, py + 1.5, pz, camAngle, 0.2, 10)

    ClearBackground(50, 50, 60, 255)
    DrawCube(0, -0.5, 0, 25, 1, 25, 0, 128, 0, 255)
    DrawSphere(px, py, pz, 0.5, 255, 0, 0, 255)
WEND

DestroyWorld3D("w")
```

See [templates/3d_game.bas](../templates/3d_game.bas), [examples/run_3d_physics_demo.bas](../examples/run_3d_physics_demo.bas), and [API_REFERENCE.md](../API_REFERENCE.md).

---

## Collision callbacks (2D)

When a Box2D body collides, you can call a BASIC sub by name:

```basic
GAME.SetCollisionHandler(bodyId, "OnHit")
// After Step2D("w", dt, 8, 3):
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

## GUI

Use **BeginUI()** … **EndUI()** each frame and add widgets: **Label**, **Button**, **Slider**, **Checkbox**, **TextBox**, **Dropdown**, **ProgressBar**, **WindowBox** / **EndWindowBox**, **GroupBox** / **EndGroupBox**. Call after ClearBackground in your game loop. See **[GUI Guide](GUI_GUIDE.md)** and [API_REFERENCE.md](../API_REFERENCE.md) (section UI).

---

## Quality of life

- **GAME.AssetPath("hero.png")** → `"assets/hero.png"` so you can keep art in an `assets/` folder.
- **GAME.ClampDelta(0.05)** → use instead of raw `DeltaTime()`/`GetFrameTime()` when stepping physics to avoid large steps after a stall.
- **GAME.ShowDebug()** or **ShowDebug("extra line")** – draw FPS (and optional text) for quick debugging.

---

## State machines

CyberBasic does not have built-in STATE/TRANSITION syntax. You can implement a state machine with a variable and **SELECT CASE**, or by calling different Subs per state:

**Using a state variable and SELECT CASE:**

```basic
DIM gameState AS Integer   // 0=menu, 1=playing, 2=paused, 3=gameover
gameState = 0

WHILE NOT WindowShouldClose()
  SELECT CASE gameState
    CASE 0
      // Draw menu; if Start pressed then gameState = 1
    CASE 1
      // Update and draw game; if Pause then gameState = 2
    CASE 2
      // Draw pause screen; if Resume then gameState = 1
    CASE 3
      // Draw game over; if Restart then gameState = 1
  END SELECT
WEND
```

**Using Subs per state:** keep a variable (e.g. `state`) and call `UpdateMenu()`, `UpdatePlaying()`, `DrawMenu()`, `DrawPlaying()` from the main loop based on `state`; set `state` inside those Subs when transitioning.

---

## Best practices

1. **Use CONST for configuration** – e.g. `CONST ScreenW = 800`, `CONST FPS = 60`, `CONST Gravity = -18`.
2. **Clamp delta time** – `IF dt > 0.05 THEN LET dt = 0.016` (or use `GAME.ClampDelta(0.05)`) before physics step.
3. **Load resources once** – Load textures, fonts, models at startup; unload in cleanup.
4. **Use GetAxisX/GetAxisY** for simple 2D movement.
5. **Organize with functions** – e.g. `UpdatePlayer()`, `DrawPlayer()`, `UpdateEnemies()` called from the main loop.
6. **Draw each frame** – ClearBackground and your draw calls in the loop.

---

## Physics implementation status

- **Box2D (2D):** Worlds, bodies, shapes, raycast, and collision are implemented. All joint types are supported: **CreateDistanceJoint2D**, **CreateRevoluteJoint2D**, **CreatePrismaticJoint2D**, **CreateWeldJoint2D**, **CreateRopeJoint2D**, **CreatePulleyJoint2D**, **CreateGearJoint2D**, **CreateWheelJoint2D**. Use **SetJointLimits2D** and **SetJointMotor2D** to configure joints; **DestroyJoint2D** to remove them.
- **Bullet (3D):** Worlds, bodies, shapes, raycast, step, and collision are implemented. Body properties are supported: **SetFriction3D**, **SetRestitution3D**, **SetDamping3D**, **SetKinematic3D**, **SetGravity3D**, **SetLinearFactor3D**, **SetAngularFactor3D**, **SetCCD3D**. **3D constraint joints** (CreateHingeJoint3D, CreateSliderJoint3D, CreateConeTwistJoint3D, CreatePointToPointJoint3D, CreateFixedJoint3D, SetJointLimits3D, SetJointMotor3D) remain stubbed in the pure-Go engine.

See [API_REFERENCE.md](../API_REFERENCE.md) and [2D Physics Guide](2D_PHYSICS_GUIDE.md) / [3D Physics Guide](3D_PHYSICS_GUIDE.md) for the full list.

---

## 2D and 3D quick reference

### 2D game quick reference

| Task | Commands |
|------|----------|
| Loop & time | `WHILE NOT WindowShouldClose() … WEND`, `DeltaTime()` |
| Window / center | `InitWindow`, `GetScreenWidth`, `GetScreenHeight`, `GetScreenCenterX`, `GetScreenCenterY` |
| Camera | `SetCamera2D`, `SetCamera2DCenter`; with Box2D: `GAME.SetCamera2DFollow`, `GAME.UpdateCamera2D` |
| Draw | `ClearBackground`, draw calls (`DrawRectangle`, `DrawCircle`, `DrawTexture`, `DrawText`). |
| Input | `GetAxisX`, `GetAxisY`, `IsKeyDown` |
| Distance | `Distance2D(x1, y1, x2, y2)` |
| Physics (optional) | `CreateWorld2D`, `Step2D`, `GAME.MoveHorizontal2D`, `GAME.Jump2D` |

Full list: [2D Graphics Guide – Full 2D command reference](2D_GRAPHICS_GUIDE.md#full-2d-command-reference).

### 3D game quick reference

| Task | Commands |
|------|----------|
| Loop & time | `WHILE NOT WindowShouldClose() … WEND`, `DeltaTime()` or `GetFrameTime()` |
| Window / center | `InitWindow`, `GetScreenCenterX`, `GetScreenCenterY` |
| Camera | `SetCamera3D`; orbit: `GAME.CameraOrbit`, `GAME.SetCamera3DOrbit`, `GAME.UpdateCamera3D` |
| Draw | `ClearBackground`, `DrawCube`, `DrawSphere`, `DrawModel`, fog if needed |
| Movement | `GAME.MoveWASD`, `GAME.SnapToGround3D` |
| Distance | `Distance3D(x1, y1, z1, x2, y2, z2)` |
| Physics (optional) | `CreateWorld3D`, `Step3D`, `GetPositionX3D` / `GetPositionY3D` / `GetPositionZ3D` |

Full list: [3D Graphics Guide – Full 3D command reference](3D_GRAPHICS_GUIDE.md#full-3d-command-reference).

---

For more examples see [examples/README.md](../examples/README.md) and the [Documentation Index](DOCUMENTATION_INDEX.md).
