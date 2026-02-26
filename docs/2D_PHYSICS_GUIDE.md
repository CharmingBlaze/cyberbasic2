# 2D Physics Guide (Box2D)

Complete guide to 2D physics in CyberBasic using Box2D: worlds, bodies, shapes, joints, raycast, collision, and integration with the hybrid loop and GAME.* helpers.

## Table of Contents

1. [Quick start](#quick-start)
2. [BOX2D.* vs legacy names](#box2d-vs-legacy-names)
3. [Worlds and bodies](#worlds-and-bodies)
4. [Body shapes](#body-shapes)
5. [Joints](#joints)
6. [Raycast](#raycast)
7. [Collision](#collision)
8. [Hybrid loop (StepAllPhysics2D)](#hybrid-loop-stepallphysics2d)
9. [GAME.* helpers](#game-helpers)
10. [Full command reference](#full-command-reference)
11. [Example](#example)
12. [See also](#see-also)

---

## Quick start

Create a world, add a ground and a player box, step each frame, and draw at the body position:

```basic
InitWindow(800, 600, "2D Physics")
SetTargetFPS(60)

// Legacy flat API: world name, gravity x, gravity y
CreateWorld2D("w", 0, -10)
CreateBox2D("w", "ground", 400, 550, 800, 20, 0, 0)   // static (mass 0, isDynamic 0)
CreateBox2D("w", "player", 400, 300, 32, 32, 1, 1)   // dynamic

WHILE NOT WindowShouldClose()
    VAR dt = GetFrameTime()
    IF dt > 0.05 THEN LET dt = 0.016
    Step2D("w", dt)

    VAR px = GetPositionX2D("w", "player")
    VAR py = GetPositionY2D("w", "player")

    ClearBackground(30, 30, 50, 255)
    DrawRectangle(px - 16, py - 16, 32, 32, 255, 100, 100, 255)
WEND

DestroyWorld2D("w")
CloseWindow()
```

You can use **BOX2D.*** prefixed names instead: **BOX2D.CreateWorld**(worldId, gravityX, gravityY), **BOX2D.Step**(worldId, dt, velocityIters, positionIters), **BOX2D.GetPositionX**(worldId, bodyId), etc. See [API Reference](../API_REFERENCE.md) section 14 for full signatures.

---

## BOX2D.* vs legacy names

- **Prefixed (BOX2D.\*):** **BOX2D.CreateWorld**(worldId, gravityX, gravityY), **BOX2D.Step**(worldId, dt [, velIters, posIters]), **BOX2D.CreateBody**, **BOX2D.GetPositionX/Y**, **BOX2D.SetLinearVelocity**, **BOX2D.ApplyForce**, etc. Use when you want an explicit namespace.
- **Legacy (flat):** **CreateWorld2D**(worldId, gravityX, gravityY), **Step2D**(worldId, dt), **CreateBox2D**, **CreateCircle2D**, **GetPositionX2D** / **GetPositionY2D**, **SetVelocity2D**, **ApplyForce2D**, **ApplyImpulse2D**, etc. Same behavior; no prefix.

All commands are **case-insensitive**. For the complete list see [API Reference](../API_REFERENCE.md) section 14.

---

## Worlds and bodies

- **Create a world:** **CreateWorld2D**(worldId, gravityX, gravityY) or **BOX2D.CreateWorld**(worldId, gravityX, gravityY). Gravity is in m/s² (e.g. 0, -10 for downward).
- **Step the simulation:** **Step2D**(worldId, dt) or **BOX2D.Step**(worldId, dt, velocityIters, positionIters). Call once per frame with delta time (e.g. GetFrameTime()). Clamp dt (e.g. max 0.05) for stability.
- **Destroy:** **DestroyWorld2D**(worldId) / **BOX2D.DestroyWorld**(worldId). **DestroyBody** / **BOX2D.DestroyBody**(worldId, bodyId) to remove a body.
- **Body IDs:** You pass a string bodyId when creating shapes (e.g. "player", "ground"); use the same id for GetPositionX2D, SetVelocity2D, etc.

---

## Body shapes

- **CreateBox2D**(worldId, bodyId, x, y, width, height, mass, isDynamic) — axis-aligned box. mass 0 = static; isDynamic 1 = dynamic.
- **CreateCircle2D**(worldId, bodyId, x, y, radius, mass, isDynamic) — circle.
- **CreatePolygon2D**(worldId, bodyId, x, y, mass, isDynamic, v1x, v1y, v2x, v2y, v3x, v3y, …) — polygon (vertices relative to body center).
- **CreateEdge2D**(worldId, bodyId, x1, y1, x2, y2) — line segment (static).
- **CreateChain2D** — chain of edges (stubbed in current build; see API Reference).

Use **SetSensor2D**(worldId, bodyId, 1) to make a body a sensor (no physical collision response). Other options: **SetFriction2D**, **SetRestitution2D**, **SetDamping2D**, **SetFixedRotation2D**, **SetGravityScale2D**, **SetMass2D**, **SetBullet2D** (continuous collision).

---

## Joints

- **CreateDistanceJoint2D**(worldId, bodyAId, bodyBId, length) — distance constraint between two bodies (implemented).
- Other joint types (Revolute, Prismatic, Pulley, Gear, Weld, Rope, Wheel) and **SetJointLimits2D** / **SetJointMotor2D** are stubbed; see [API Reference](../API_REFERENCE.md).

---

## Raycast

- **RayCast2D**(worldId, x1, y1, x2, y2) — cast a ray from (x1,y1) to (x2,y2). Returns true if something was hit.
- After a hit: **RayHitX2D**(), **RayHitY2D**() — hit point; **RayHitBody2D**() — body id; **RayHitNormalX2D**(), **RayHitNormalY2D**() — normal.

---

## Collision

After **Step2D** (or **BOX2D.Step**), you can query contacts:

- **GetCollisionCount2D**(worldId) — number of contact pairs this step.
- **GetCollisionOther2D**(index) — the "other" body id in the pair (0-based index).
- **GetCollisionNormalX2D**(index), **GetCollisionNormalY2D**(index) — contact normal.

**Collision callbacks:** Register a Sub to be called when a specific body collides:

- **GAME.SetCollisionHandler**(bodyId, subName) — when bodyId collides, the engine will call the Sub named subName with (otherBodyId).
- After **Step2D**, call **GAME.ProcessCollisions2D**(worldId) to invoke all registered handlers for this frame.

See [Game Development Guide](GAME_DEVELOPMENT_GUIDE.md#collision-callbacks-2d).

---

## Hybrid loop (StepAllPhysics2D)

When you define **update(dt)** and **draw()** and use a game loop (WHILE NOT WindowShouldClose()), the compiler injects an automatic pipeline. Part of that pipeline is **StepAllPhysics2D(dt)** — all registered Box2D worlds are stepped with the same dt. You do **not** need to call **Step2D** or **BOX2D.Step** per world yourself.

See [Program Structure](PROGRAM_STRUCTURE.md#hybrid-updatedraw-loop).

---

## GAME.* helpers

- **GAME.SetCamera2DFollow**(worldId, bodyId, xOffset, yOffset) — set the 2D camera to follow a body. Then each frame call **GAME.UpdateCamera2D**() (e.g. in your loop after physics step).
- **GAME.SyncSpriteToBody2D**(worldId, bodyId, spriteId) — update a sprite’s position to match the body (e.g. call from draw).
- **GAME.SetCollisionHandler**(bodyId, subName), **GAME.ProcessCollisions2D**(worldId) — see [Collision](#collision).

See [Game Development Guide](GAME_DEVELOPMENT_GUIDE.md#game-helpers).

---

## Full command reference

| Command | Arguments | Returns | Description |
|---------|-----------|---------|-------------|
| **BOX2D.CreateWorld** | (worldId, gravityX, gravityY) | — | Create world |
| **BOX2D.Step** | (worldId, dt [, velIters, posIters]) | — | Step simulation |
| **BOX2D.DestroyWorld** | (worldId) | — | Destroy world |
| **BOX2D.CreateBody** | (worldId, bodyId, type, shape, x, y, …) | bodyId | Create body (see binding) |
| **BOX2D.GetPositionX/Y** | (worldId, bodyId) | float | Position |
| **BOX2D.GetAngle** | (worldId, bodyId) | float | Angle (radians) |
| **BOX2D.SetLinearVelocity** | (worldId, bodyId, vx, vy) | — | Set velocity |
| **BOX2D.ApplyForce** | (worldId, bodyId, fx, fy, x, y) | — | Apply force at point |
| **CreateWorld2D** | (worldId, gravityX, gravityY) | — | Legacy: create world |
| **Step2D** | (worldId, dt) | — | Legacy: step |
| **CreateBox2D** | (worldId, bodyId, x, y, w, h, mass, isDynamic) | — | Box shape |
| **CreateCircle2D** | (worldId, bodyId, x, y, radius, mass, isDynamic) | — | Circle shape |
| **GetPositionX2D/Y2D** | (worldId, bodyId) | float | Position |
| **SetVelocity2D** | (worldId, bodyId, vx, vy) | — | Set velocity |
| **ApplyForce2D** | (worldId, bodyId, fx, fy) | — | Apply force |
| **ApplyImpulse2D** | (worldId, bodyId, ix, iy) | — | Apply impulse |
| **CreateDistanceJoint2D** | (worldId, bodyAId, bodyBId, length) | — | Distance joint |
| **RayCast2D** | (worldId, x1, y1, x2, y2) | bool | Ray cast |
| **RayHitX2D/Y2D** | () | float | Hit point |
| **RayHitBody2D** | () | bodyId | Hit body |
| **GetCollisionCount2D** | (worldId) | int | Contact count |
| **GetCollisionOther2D** | (index) | bodyId | Other body in contact |
| **StepAllPhysics2D** | (dt) | — | Step all worlds (hybrid loop) |

For the full list including **SetSensor2D**, **SetFriction2D**, **GetCollisionNormalX2D/Y2D**, and legacy **DestroyWorld2D**, **DestroyBody**, see [API Reference](../API_REFERENCE.md) section 14.

---

## Example

Minimal runnable example using the legacy flat API:

```basic
InitWindow(800, 600, "Box2D Demo")
SetTargetFPS(60)
CreateWorld2D("w", 0, -10)
CreateBox2D("w", "ground", 400, 550, 800, 20, 0, 0)
CreateBox2D("w", "player", 400, 200, 40, 40, 1, 1)

WHILE NOT WindowShouldClose()
    VAR dt = GetFrameTime()
    IF dt > 0.05 THEN LET dt = 0.016
    Step2D("w", dt)
    VAR px = GetPositionX2D("w", "player")
    VAR py = GetPositionY2D("w", "player")
    ClearBackground(30, 30, 50, 255)
    DrawRectangle(px - 20, py - 20, 40, 40, 255, 100, 100, 255)
WEND

DestroyWorld2D("w")
CloseWindow()
```

More examples: [examples/box2d_demo.bas](../examples/box2d_demo.bas), [examples/box2d_physics_demo.bas](../examples/box2d_physics_demo.bas), [examples/README.md](../examples/README.md).

---

## See also

- [API Reference](../API_REFERENCE.md) (section 14 – Box2D)
- [Game Development Guide](GAME_DEVELOPMENT_GUIDE.md) – 2D physics, GAME.* helpers, collision callbacks
- [2D Graphics Guide](2D_GRAPHICS_GUIDE.md) – Drawing primitives and textures
- [Command Reference](COMMAND_REFERENCE.md) – Commands by feature
