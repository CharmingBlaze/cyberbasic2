# 3D Physics Guide (Bullet)

Complete guide to 3D physics in CyberBasic using Bullet: worlds, gravity, rigid bodies, shapes, position/rotation, velocity and forces, raycast, and integration with the hybrid loop and GAME.* helpers.

## Table of Contents

1. [Quick start](#quick-start)
2. [BULLET.* vs legacy names](#bullet-vs-legacy-names)
3. [Worlds and gravity](#worlds-and-gravity)
4. [Body shapes](#body-shapes)
5. [Position and rotation](#position-and-rotation)
6. [Velocity and forces](#velocity-and-forces)
7. [Raycast](#raycast)
8. [Hybrid loop (StepAllPhysics3D)](#hybrid-loop-stepallphysics3d)
9. [GAME.* 3D helpers](#game-3d-helpers)
10. [Full command reference](#full-command-reference)
11. [Example](#example)
12. [See also](#see-also)

---

## Quick start

Create a 3D world, add a ground box and a sphere, step each frame, and draw at the body position:

```basic
InitWindow(800, 600, "3D Physics")
SetTargetFPS(60)
InitAudioDevice()

// Legacy flat API: world id, gravity x, y, z
CreateWorld3D("w", 0, -18, 0)
CreateBox3D("w", "ground", 0, -0.5, 0, 12.5, 0.5, 12.5, 0)   // static (mass 0)
CreateSphere3D("w", "player", 0, 2, 0, 0.5, 1)   // dynamic sphere

VAR camAngle = 0

WHILE NOT WindowShouldClose()
    VAR dt = GetFrameTime()
    IF dt > 0.05 THEN LET dt = 0.016
    Step3D("w", dt)

    VAR px = GetPositionX3D("w", "player")
    VAR py = GetPositionY3D("w", "player")
    VAR pz = GetPositionZ3D("w", "player")

    ClearBackground(50, 50, 60, 255)
    DrawCube(0, -0.5, 0, 25, 1, 25, 0, 128, 0, 255)
    DrawSphere(px, py, pz, 0.5, 255, 0, 0, 255)
WEND

DestroyWorld3D("w")
CloseWindow()
```

You can use **BULLET.*** prefixed names: **BULLET.CreateWorld**(worldId, gravityX, gravityY, gravityZ), **BULLET.SetGravity**(worldId, x, y, z), **BULLET.Step**(worldId, dt), **BULLET.GetPositionX/Y/Z**, etc. See [API Reference](../API_REFERENCE.md) section 15.

---

## BULLET.* vs legacy names

- **Prefixed (BULLET.\*):** **BULLET.CreateWorld**(worldId, gx, gy, gz), **BULLET.SetGravity**(worldId, x, y, z), **BULLET.Step**(worldId, dt), **BULLET.CreateBox**, **BULLET.CreateSphere**, **BULLET.GetPositionX/Y/Z**, **BULLET.SetVelocity**, **BULLET.ApplyForce**, **BULLET.ApplyCentralForce**, **BULLET.ApplyImpulse**, **BULLET.RayCast**, etc.
- **Legacy (flat):** **CreateWorld3D**(worldId, gx, gy, gz), **Step3D**(worldId, dt), **CreateBox3D**, **CreateSphere3D**, **GetPositionX3D/Y3D/Z3D**, **SetPosition3D**, **SetVelocity3D**, **GetYaw3D**, **GetPitch3D**, **GetRoll3D**, **SetRotation3D**, **ApplyForce3D**, **ApplyImpulse3D**, **RayCast3D**, **RayHitX3D/Y3D/Z3D**, **RayHitBody3D**, etc.

All commands are **case-insensitive**. For the complete list see [API Reference](../API_REFERENCE.md) section 15.

---

## Worlds and gravity

- **Create a world:** **CreateWorld3D**(worldId, gravityX, gravityY, gravityZ) or **BULLET.CreateWorld**(worldId, gx, gy, gz). Gravity is in m/s² (e.g. 0, -18, 0 for downward).
- **Set gravity:** **BULLET.SetGravity**(worldId, x, y, z) to change gravity after creation. Legacy **SetGravity3D** is available (see API Reference).
- **Step:** **Step3D**(worldId, dt) or **BULLET.Step**(worldId, dt). Call once per frame; clamp dt (e.g. max 0.05) for stability.
- **Destroy:** **DestroyWorld3D**(worldId) / **BULLET.DestroyWorld**(worldId). **BULLET.DestroyBody**(worldId, bodyId) to remove a body.

---

## Body shapes

- **CreateBox3D**(worldId, bodyId, x, y, z, halfWidth, halfHeight, halfDepth, mass) — box (half-extents). mass 0 = static.
- **CreateSphere3D**(worldId, bodyId, x, y, z, radius, mass) — sphere.
- **CreateCapsule3D**, **CreateCylinder3D**, **CreateCone3D** — capsule, cylinder, cone (legacy names; see API Reference).
- **CreateStaticMesh3D**, **CreateHeightmap3D** — static mesh and heightfield (legacy; see API Reference).
- **CreateCompound3D**, **AddShapeToCompound3D** — compound bodies (multiple shapes). **SetScale3D** for scaling.

**Body properties (implemented):** **SetFriction3D**(worldId, bodyId, friction), **SetRestitution3D**(worldId, bodyId, restitution), **SetDamping3D**(worldId, bodyId, linearDamp, angularDamp), **SetKinematic3D**(worldId, bodyId, kinematic), **SetGravity3D**(worldId, bodyId, gravityScale), **SetLinearFactor3D**(worldId, bodyId, fx, fy, fz), **SetAngularFactor3D**(worldId, bodyId, ax, ay, az), **SetCCD3D**(worldId, bodyId, enable). These are used in the pure-Go engine's Step and collision resolution. **3D joints** (CreateHingeJoint3D, CreateSliderJoint3D, etc.) remain stubbed in the pure-Go engine; see API Reference.

---

## Position and rotation

- **Position:** **GetPositionX3D**(worldId, bodyId), **GetPositionY3D**, **GetPositionZ3D** — or **BULLET.GetPositionX/Y/Z**. **SetPosition3D**(worldId, bodyId, x, y, z) / **BULLET.SetPosition** to teleport.
- **Rotation (Euler):** **GetYaw3D**(worldId, bodyId), **GetPitch3D**, **GetRoll3D** — or **BULLET.GetRotationX/Y/Z**. **SetRotation3D**(worldId, bodyId, yaw, pitch, roll) / **BULLET.SetRotation**(worldId, bodyId, x, y, z) to set rotation.

Use these each frame after **Step3D** to draw your 3D model or to drive **GAME.CameraOrbit** / **GAME.SetCamera3DOrbit**.

---

## Velocity and forces

- **Velocity:** **SetVelocity3D**(worldId, bodyId, vx, vy, vz) / **BULLET.SetVelocity**. **GetVelocityX3D/Y3D/Z3D** or **BULLET.GetVelocityX/Y/Z**.
- **Angular velocity:** **SetAngularVelocity3D**, **GetAngularVelocityX3D/Y3D/Z3D** (legacy; see API Reference).
- **Forces:** **BULLET.ApplyForce**(worldId, bodyId, fx, fy, fz, px, py, pz) — force at world point. **BULLET.ApplyCentralForce**(worldId, bodyId, fx, fy, fz) — at center. **ApplyForce3D**, **ApplyImpulse3D** (legacy). **ApplyTorque3D**, **ApplyTorqueImpulse3D** for rotation.

---

## Raycast

- **RayCast3D**(worldId, originX, originY, originZ, dirX, dirY, dirZ) or **BULLET.RayCast**(worldId, ox, oy, oz, dx, dy, dz) — cast a ray from origin along direction. Returns true if hit.
- After a hit: **RayHitX3D**(), **RayHitY3D**(), **RayHitZ3D**() — hit point; **RayHitBody3D**() — body id. BULLET: **GetRayCastHitX/Y/Z**(), **GetRayCastHitBody**().

---

## Hybrid loop (StepAllPhysics3D)

When you define **update(dt)** and **draw()** and use the automatic game loop, the pipeline calls **StepAllPhysics3D(dt)** — all registered Bullet worlds are stepped with the same dt. You do **not** need to call **Step3D** or **BULLET.Step** per world yourself.

See [Program Structure](PROGRAM_STRUCTURE.md#hybrid-updatedraw-loop).

---

## GAME.* 3D helpers

- **GAME.CameraOrbit**(cx, cy, cz, angle, pitch, distance) — position the 3D camera to orbit around a point (e.g. player position). Call each frame after reading body position.
- **GAME.MoveWASD**(worldId, bodyId, angle, speed, jumpForce, dt) — apply WASD movement and jump to a body; use with **GAME.CameraOrbit** for a third-person controller.
- **GAME.SetCamera3DOrbit**(targetX, targetY, targetZ, …), **GAME.UpdateCamera3D** — alternative camera helpers (see [Game Development Guide](GAME_DEVELOPMENT_GUIDE.md)).

See [Game Development Guide](GAME_DEVELOPMENT_GUIDE.md#3d-physics-bullet).

---

## Full command reference

| Command | Arguments | Returns | Description |
|---------|-----------|---------|-------------|
| **BULLET.CreateWorld** | (worldId, gx, gy, gz) | — | Create world |
| **BULLET.SetGravity** | (worldId, x, y, z) | — | Set gravity |
| **BULLET.Step** | (worldId, dt) | — | Step simulation |
| **BULLET.CreateBox** | (worldId, x, y, z, w, h, d, mass) | bodyId | Box body |
| **BULLET.CreateSphere** | (worldId, x, y, z, radius, mass) | bodyId | Sphere body |
| **BULLET.GetPositionX/Y/Z** | (worldId, bodyId) | float | Position |
| **BULLET.SetPosition** | (worldId, bodyId, x, y, z) | — | Set position |
| **BULLET.GetRotationX/Y/Z** | (worldId, bodyId) | float | Rotation (euler) |
| **BULLET.SetVelocity** | (worldId, bodyId, vx, vy, vz) | — | Set velocity |
| **BULLET.ApplyForce** | (worldId, bodyId, fx, fy, fz, px, py, pz) | — | Apply force at point |
| **BULLET.ApplyCentralForce** | (worldId, bodyId, fx, fy, fz) | — | Apply force at center |
| **BULLET.ApplyImpulse** | (worldId, bodyId, ix, iy, iz, x, y, z) | — | Apply impulse |
| **BULLET.RayCast** | (worldId, ox, oy, oz, dx, dy, dz) | bool | Ray cast |
| **CreateWorld3D** | (worldId, gx, gy, gz) | — | Legacy: create world |
| **Step3D** | (worldId, dt) | — | Legacy: step |
| **CreateBox3D** | (worldId, bodyId, x, y, z, hw, hh, hd, mass) | — | Box |
| **CreateSphere3D** | (worldId, bodyId, x, y, z, radius, mass) | — | Sphere |
| **GetPositionX3D/Y3D/Z3D** | (worldId, bodyId) | float | Position |
| **GetYaw3D / GetPitch3D / GetRoll3D** | (worldId, bodyId) | float | Rotation |
| **SetRotation3D** | (worldId, bodyId, yaw, pitch, roll) | — | Set rotation |
| **SetVelocity3D** | (worldId, bodyId, vx, vy, vz) | — | Set velocity |
| **StepAllPhysics3D** | (dt) | — | Step all worlds (hybrid loop) |

For the full list including **CreateCapsule3D**, **CreateStaticMesh3D**, **GetCollisionCount3D**, **GetCollisionOther3D**, and legacy **DestroyWorld3D** / **DestroyBody**, see [API Reference](../API_REFERENCE.md) section 15.

---

## Example

Minimal runnable example using the legacy flat API:

```basic
InitWindow(800, 600, "3D Physics Demo")
SetTargetFPS(60)
InitAudioDevice()
CreateWorld3D("w", 0, -18, 0)
CreateBox3D("w", "ground", 0, -0.5, 0, 12.5, 0.5, 12.5, 0)
CreateSphere3D("w", "ball", 0, 3, 0, 0.5, 1)

WHILE NOT WindowShouldClose()
    VAR dt = GetFrameTime()
    IF dt > 0.05 THEN LET dt = 0.016
    Step3D("w", dt)
    VAR px = GetPositionX3D("w", "ball")
    VAR py = GetPositionY3D("w", "ball")
    VAR pz = GetPositionZ3D("w", "ball")
    ClearBackground(50, 50, 60, 255)
    DrawCube(0, -0.5, 0, 25, 1, 25, 0, 128, 0, 255)
    DrawSphere(px, py, pz, 0.5, 255, 0, 0, 255)
WEND

DestroyWorld3D("w")
CloseWindow()
```

More examples: [templates/3d_game.bas](../templates/3d_game.bas), [examples/run_3d_physics_demo.bas](../examples/run_3d_physics_demo.bas), [examples/README.md](../examples/README.md).

---

## See also

- [API Reference](../API_REFERENCE.md) (section 15 – Bullet)
- [Game Development Guide](GAME_DEVELOPMENT_GUIDE.md) – 3D physics, GAME.CameraOrbit, GAME.MoveWASD
- [3D Graphics Guide](3D_GRAPHICS_GUIDE.md) – Cameras, models, lighting
- [Command Reference](COMMAND_REFERENCE.md) – Commands by feature
