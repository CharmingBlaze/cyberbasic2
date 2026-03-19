# 3D Physics Guide (Bullet)

Complete guide to 3D physics in CyberBASIC2 using the Bullet-shaped 3D physics API: worlds, gravity, rigid bodies, shapes, position/rotation, velocity and forces, raycast, and integration with the hybrid loop and GAME.* helpers.

**Purpose:** Bullet-style 3D physics for characters, projectiles, and simple collision.

**When to use 3D physics:** Use 3D physics when you need gravity, collision, or raycasting in a 3D world. See [3D Graphics Guide](3D_GRAPHICS_GUIDE.md).

Today, the shipped 3D backend is a pure-Go fallback rather than native Bullet. You can query it at runtime with `BulletBackendName()`, `BulletBackendMode()`, `BulletNativeAvailable()`, `BulletJointsAvailable()`, and `BulletFeatureAvailable(feature$)`, which currently report `purego-fallback`, `fallback`, `0`, `1` (PointToPoint/Fixed joints), and per-feature support flags.

**What's supported today (shipped build, no CGO)**

- **Supported:** CreateWorld3D, SetWorldGravity3D, Step3D, CreateBox3D, CreateSphere3D, CreateCapsule3D, CreateCylinder3D, CreateCone3D; GetPositionX/Y/Z3D, SetPosition3D, SetVelocity3D, ApplyForce3D, ApplyImpulse3D, **ApplyTorque3D**, **ApplyTorqueImpulse3D**; body properties (SetFriction3D, SetRestitution3D, SetDamping3D, SetKinematic3D, SetCCD3D, etc.); **CreatePointToPointJoint3D**, **CreateFixedJoint3D**, **CreateHingeJoint3D**, **CreateSliderJoint3D**, **CreateConeTwistJoint3D**, **SetJointLimits3D**, **SetJointMotor3D**; **CreateStaticMesh3D** (loads OBJ, uses AABB collision); RayCastFromDir3D / RayCast3D and RayHit*; DestroyBody3D, DestroyWorld3D. Good for characters, projectiles, simple collisions, and all constraint joints.
- **Not in fallback:** CreateHeightmap3D, CreateCompound3D, AddShapeToCompound3D, and exact triangle-mesh narrow-phase collision. Call **BulletFeatureAvailable()** or **BulletNativeAvailable()** before using these; see [ROADMAP_IMPLEMENTATION.md](ROADMAP_IMPLEMENTATION.md) for the full gap list.

## Table of Contents

1. [Quick start](#quick-start)
2. [API style (flat names)](#api-style-flat-names)
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

The API uses **flat names** only (no namespace). Legacy **BULLET.*** in source is rewritten at compile time. **3D constraint joints** (CreateHingeJoint3D, CreateSliderJoint3D, CreateConeTwistJoint3D, SetJointLimits3D, SetJointMotor3D) are implemented in the shipped pure-Go fallback. See [API Reference](../API_REFERENCE.md) section 15.

---

## API style (flat names)

- **Flat names:** **CreateWorld3D**, **SetWorldGravity3D**, **Step3D**, **CreateBox3D**, **CreateSphere3D**, **GetPositionX3D** / **GetPositionY3D** / **GetPositionZ3D**, **SetVelocity3D**, **ApplyForce3D**, **ApplyImpulse3D**, **RayCastFromDir3D** or **RayCast3D**, **RayHitX3D** / **RayHitY3D** / **RayHitZ3D**, **RayHitBody3D**, etc. Use these in all new code.

All commands are **case-insensitive**. For the complete list see [API Reference](../API_REFERENCE.md) section 15.

Treat the current 3D backend as:

- good for simple rigid bodies, gravity, velocity-based character helpers, basic raycasts, and broad-phase collision queries
- includes default-world helper commands such as `SetBodyPosition`, `GetBodyPosition`, `SetBodyVelocity`, and `GetBodyVelocity` for quick DBPro-style scripts
- not a full native Bullet replacement yet for constraints, exact mesh collision, or high-fidelity narrow-phase behavior

---

## Worlds and gravity

- **Create a world:** **CreateWorld3D**(worldId, gravityX, gravityY, gravityZ). Gravity is in m/s² (e.g. 0, -18, 0 for downward).
- **Set gravity:** **SetWorldGravity3D**(worldId, x, y, z) to change gravity after creation.
- **Step:** **Step3D**(worldId, dt). Call once per frame; clamp dt (e.g. max 0.05) for stability.
- **Destroy:** **DestroyWorld3D**(worldId). **DestroyBody3D**(worldId, bodyId) to remove a body.

---

## Body shapes

- **CreateBox3D**(worldId, bodyId, x, y, z, halfWidth, halfHeight, halfDepth, mass) — box (half-extents). mass 0 = static.
- **CreateSphere3D**(worldId, bodyId, x, y, z, radius, mass) — sphere.
- **CreateCapsule3D**, **CreateCylinder3D**, **CreateCone3D** — capsule, cylinder, cone (legacy names; see API Reference). In the pure-Go runtime, capsules are still approximated for collision/raycast math, but the requested capsule height is now preserved in the body's bounds instead of being ignored.
- **CreateStaticMesh3D**(worldId, bodyId, meshPath$) loads an OBJ file from meshPath and creates a static body with AABB collision. If the mesh fails to load, a 1×1×1 box placeholder is used.
- **CreateHeightmap3D**, **CreateCompound3D**, **AddShapeToCompound3D** are unsupported in the current fallback and now return explicit errors instead of silently succeeding.
- **SetScale3D** scales fallback bounds; it is not a substitute for native compound or mesh collision.

**Body properties (implemented):** **SetFriction3D**, **SetRestitution3D**, **SetDamping3D**, **SetKinematic3D**, **SetGravity3D**, **SetLinearFactor3D**, **SetAngularFactor3D**, **SetCCD3D**.

**3D joints:** **BulletJointsAvailable**() returns 1 in the shipped backend. All joints are implemented: **CreatePointToPointJoint3D**, **CreateFixedJoint3D**, **CreateHingeJoint3D**, **CreateSliderJoint3D**, **CreateConeTwistJoint3D**, **SetJointLimits3D**, **SetJointMotor3D**. See [3D joints](#3d-joints) below.

---

## Position and rotation

- **Position:** **GetPositionX3D**(worldId, bodyId), **GetPositionY3D**, **GetPositionZ3D**. **SetPosition3D**(worldId, bodyId, x, y, z) to teleport.
- **Rotation (Euler):** **GetYaw3D**(worldId, bodyId), **GetPitch3D**, **GetRoll3D**. **SetRotation3D**(worldId, bodyId, rx, ry, rz) to set rotation.

Use these each frame after **Step3D** to draw your 3D model or to drive **GAME.CameraOrbit** / **GAME.SetCamera3DOrbit**.

---

## Velocity and forces


- **Velocity:** **SetVelocity3D**(worldId, bodyId, vx, vy, vz). **GetVelocityX3D** / **GetVelocityY3D** / **GetVelocityZ3D**.
- **Angular velocity:** **SetAngularVelocity3D**, **GetAngularVelocityX3D/Y3D/Z3D** (see API Reference).
- **Forces:** **ApplyForce3D**(worldId, bodyId, fx, fy, fz). **ApplyImpulse3D**(worldId, bodyId, ix, iy, iz). **ApplyTorque3D** and **ApplyTorqueImpulse3D** for angular forces.

---

## 3D joints

The pure-Go fallback implements all constraint joint types. Use **BulletJointsAvailable**() to confirm (returns 1).

| Joint | Command | Description |
|-------|---------|-------------|
| **Point-to-point** | CreatePointToPointJoint3D(worldId, jointId, bodyA, bodyB, ax, ay, az, bx, by, bz) | Ball joint: anchors align |
| **Fixed** | CreateFixedJoint3D(worldId, jointId, bodyA, bodyB) | Bodies welded together |
| **Hinge** | CreateHingeJoint3D(worldId, jointId, bodyA, bodyB, ax, ay, az, bx, by, bz, axisAx, axisAy, axisAz) | Rotation around one axis only |
| **Slider** | CreateSliderJoint3D(worldId, jointId, bodyA, bodyB, ax, ay, az, bx, by, bz, axisAx, axisAy, axisAz) | Translation along one axis only |
| **Cone twist** | CreateConeTwistJoint3D(worldId, jointId, bodyA, bodyB, ax, ay, az, bx, by, bz, axisAx, axisAy, axisAz) | Ball + cone angle limit |

**Limits and motors:**
- **SetJointLimits3D**(worldId, jointId, low, high) — hinge: angle limits (radians); slider: position limits (meters); conetwist: cone angle (radians)
- **SetJointMotor3D**(worldId, jointId, targetVel, maxForce) — hinge: angular velocity (rad/s); slider: linear velocity (m/s)

Example: hinge joint for a door:

```basic
CreateWorld3D("w", 0, -18, 0)
CreateBox3D("w", "frame", 0, 1, 0, 0.5, 1, 0.1, 0)
CreateBox3D("w", "door", 0.6, 1, 0, 0.5, 1, 0.05, 1)
CreateHingeJoint3D("w", "hinge1", "frame", "door", 0.5, 0, 0, -0.5, 0, 0, 0, 1, 0)
SetJointLimits3D("w", "hinge1", -1.5, 0)
```

---

## Raycast

- **RayCast3D**(worldId, fromX, fromY, fromZ, toX, toY, toZ) — cast a ray from point to point.
- **RayCastFromDir3D**(worldId, sx, sy, sz, dx, dy, dz, maxDist) — cast from origin plus direction. Returns 1 if hit, 0 otherwise.
- After a hit: **RayHitX3D**(), **RayHitY3D**(), **RayHitZ3D**() — hit point; **RayHitBody3D**() — body id; **RayHitNormalX3D**() etc. — hit normal.
- DBP wrappers: **Raycast**(ox, oy, oz, dx, dy, dz [, maxDist]) and **Spherecast**(ox, oy, oz, dx, dy, dz, radius). The spherecast implementation currently sweeps against inflated body AABBs, so treat it as a broad-phase helper rather than exact Bullet geometry.

---

## Hybrid loop (StepAllPhysics3D)

When you define **update(dt)** and **draw()** and use the automatic game loop, the runtime now accumulates time and steps **StepAllPhysics3D** on the fixed timestep from `FixedDeltaTime()` (default 1/60). Use `FixedUpdate(rate)` plus `OnFixedUpdate(label$)` when you want an explicit fixed-step callback alongside physics.

See [Program Structure](PROGRAM_STRUCTURE.md#hybrid-updatedraw-loop).

---

## Character controller

- **CreateCharacterController**(worldId, bodyId, radius, height) — creates a capsule-style body centered at `height/2` with the requested total height preserved in the physics bounds.
- **SetCharacterControllerSpeed**(bodyId, scale) — multiplies the speed in **GAME.MoveWASD** for that body.
- **GAME.OnGround**(worldId, bodyId, planeY, tolerance) — returns 1 if body is near planeY.
- **GAME.SnapToGround**(worldId, bodyId, planeY, offset) — snaps body to ground plane.

---

## GAME.* 3D helpers

- **GAME.CameraOrbit**(cx, cy, cz, angle, pitch, distance) — position the 3D camera to orbit around a point (e.g. player position). Call each frame after reading body position.
- **GAME.MoveWASD**(worldId, bodyId, angle, speed, jumpForce, dt) — applies horizontal movement by setting character velocity and uses `jumpForce` for the Y velocity when grounded; use with **GAME.CameraOrbit** for a third-person controller.
- **GAME.SetCamera3DOrbit**(targetX, targetY, targetZ, …), **GAME.UpdateCamera3D** — alternative camera helpers (see [Game Development Guide](GAME_DEVELOPMENT_GUIDE.md)).

See [Game Development Guide](GAME_DEVELOPMENT_GUIDE.md#3d-physics-bullet).

---

## Full command reference

| Command | Arguments | Returns | Description |
|---------|-----------|---------|-------------|
| **BulletBackendName** | () | string | Backend name (`purego-fallback`) |
| **BulletBackendMode** | () | string | Backend mode (`fallback`) |
| **BulletNativeAvailable** | () | 0/1 | Whether a native Bullet backend is available |
| **BulletJointsAvailable** | () | 0/1 | Whether 3D constraint joints are implemented |
| **BulletFeatureAvailable** | (featureName$) | 0/1 | Query per-feature support in the shipped backend |
| **CreateWorld3D** | (worldId, gx, gy, gz) | — | Create world |
| **SetWorldGravity3D** | (worldId, x, y, z) | — | Set gravity |
| **Step3D** | (worldId, dt) | — | Step simulation |
| **CreateBox3D** | (worldId, bodyId, x, y, z, hx, hy, hz, mass) | — | Box body |
| **CreateSphere3D** | (worldId, bodyId, x, y, z, radius, mass) | — | Sphere body |
| **GetPositionX3D** / **GetPositionY3D** / **GetPositionZ3D** | (worldId, bodyId) | float | Position |
| **SetPosition3D** | (worldId, bodyId, x, y, z) | — | Set position |
| **GetYaw3D** / **GetPitch3D** / **GetRoll3D** | (worldId, bodyId) | float | Rotation (euler) |
| **SetVelocity3D** | (worldId, bodyId, vx, vy, vz) | — | Set velocity |
| **ApplyForce3D** | (worldId, bodyId, fx, fy, fz) | — | Apply force |
| **ApplyImpulse3D** | (worldId, bodyId, ix, iy, iz) | — | Apply impulse |
| **RayCastFromDir3D** | (worldId, sx, sy, sz, dx, dy, dz, maxDist) | 1=hit 0=miss | Ray cast (dir + maxDist) |
| **RayCast3D** | (worldId, fromX, fromY, fromZ, toX, toY, toZ) | 1=hit 0=miss | Ray cast (from–to) |
| **CreateWorld3D** | (worldId, gx, gy, gz) | — | Legacy: create world |
| **Step3D** | (worldId, dt) | — | Legacy: step |
| **CreateBox3D** | (worldId, bodyId, x, y, z, hw, hh, hd, mass) | — | Box |
| **CreateSphere3D** | (worldId, bodyId, x, y, z, radius, mass) | — | Sphere |
| **GetPositionX3D/Y3D/Z3D** | (worldId, bodyId) | float | Position |
| **GetYaw3D / GetPitch3D / GetRoll3D** | (worldId, bodyId) | float | Rotation |
| **SetRotation3D** | (worldId, bodyId, yaw, pitch, roll) | — | Set rotation |
| **SetVelocity3D** | (worldId, bodyId, vx, vy, vz) | — | Set velocity |
| **CreatePointToPointJoint3D** | (worldId, jointId, bodyA, bodyB, ax, ay, az, bx, by, bz) | — | Ball joint |
| **CreateFixedJoint3D** | (worldId, jointId, bodyA, bodyB) | — | Weld bodies |
| **CreateHingeJoint3D** | (worldId, jointId, bodyA, bodyB, ax, ay, az, bx, by, bz, axisX, axisY, axisZ) | — | Hinge joint |
| **CreateSliderJoint3D** | (worldId, jointId, bodyA, bodyB, ax, ay, az, bx, by, bz, axisX, axisY, axisZ) | — | Slider joint |
| **CreateConeTwistJoint3D** | (worldId, jointId, bodyA, bodyB, ax, ay, az, bx, by, bz, axisX, axisY, axisZ) | — | Cone twist joint |
| **SetJointLimits3D** | (worldId, jointId, low, high) | — | Set joint limits |
| **SetJointMotor3D** | (worldId, jointId, targetVel, maxForce) | — | Set joint motor |
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

More examples: [templates/3d_game.bas](../templates/3d_game.bas), [examples/README.md](../examples/README.md).

---

## Native Bullet backend (optional)

The shipped build uses a pure-Go fallback. An optional native Bullet backend can be built with `go build -tags bullet` when Bullet C libraries are installed. When available, **BulletNativeAvailable**() returns 1 and the runtime uses the native Bullet engine for higher-fidelity physics (exact mesh collision, compound shapes, heightmaps). If the native backend fails to initialize, the runtime falls back to the pure-Go implementation. See [ROADMAP_IMPLEMENTATION.md](ROADMAP_IMPLEMENTATION.md) for build instructions.

---

## See also

- [3D Graphics Guide](3D_GRAPHICS_GUIDE.md)
- [Level Loading](LEVEL_LOADING.md)
- [Game Development Guide](GAME_DEVELOPMENT_GUIDE.md)
- [Documentation Index](DOCUMENTATION_INDEX.md)
