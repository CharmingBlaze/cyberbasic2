# 3D Graphics Guide

Complete guide to 3D graphics in CyberBasic using the raylib API.

## Table of Contents

1. [Getting started](#getting-started)
2. [3D camera](#3d-camera)
3. [3D primitives](#3d-primitives)
4. [Models and meshes](#models-and-meshes)
5. [3D game checklist](#3d-game-checklist)
6. [Complete 3D game example](#complete-3d-game-example)
7. [3D editor and level builder](#3d-editor-and-level-builder)
8. [Full 3D command reference](#full-3d-command-reference)
9. [See also](#see-also)

---

## Getting started

### Basic 3D setup

Every 3D program needs a window, a 3D camera, and a loop with your draw calls. The compiler does not inject any frame or mode calls (DBPro-style). When you use the **hybrid loop** (define **update(dt)** and **draw()**), 3D draw calls in **draw()** are queued and flushed in order (2D then 3D then GUI). See [Program Structure](PROGRAM_STRUCTURE.md#hybrid-updatedraw-loop).

```basic
InitWindow(800, 600, "My 3D Game")
SetTargetFPS(60)

// Set camera: position (x,y,z), target (x,y,z), up vector (x,y,z)
SetCamera3D(0, 10, 10,  0, 0, 0,  0, 1, 0)

WHILE NOT WindowShouldClose()
    ClearBackground(20, 20, 30, 255)
    DrawCube(0, 0, 0, 2, 2, 2, 255, 100, 100, 255)
    DrawSphere(5, 1, 0, 1, 100, 200, 255, 255)
    DrawText("3D Scene", 10, 10, 20, 255, 255, 255, 255)
WEND

CloseWindow()
```

---

## 3D camera

### SetCamera3D

Set the 3D camera with nine numbers:

- **Position:** posX, posY, posZ (camera location)
- **Target:** targetX, targetY, targetZ (point the camera looks at)
- **Up:** upX, upY, upZ (usually 0, 1, 0 for Y-up)

```basic
SetCamera3D(0, 10, 10,  0, 0, 0,  0, 1, 0)
```

3D drawing uses this camera until you call `SetCamera3D` again.

### Orbit camera (GAME.CameraOrbit)

For a third-person or orbit camera driven by mouse/keyboard, use **GAME.CameraOrbit** each frame. It sets the raylib 3D camera from a target position and orbit angles:

```basic
// Each frame, after updating your target position (e.g. player):
GAME.CameraOrbit(targetX, targetY, targetZ, angleRad, pitchRad, distance)
```

Then draw as usual. See [Game Development Guide](GAME_DEVELOPMENT_GUIDE.md) and [templates/3d_game.bas](../templates/3d_game.bas) for a full loop with Bullet physics and `GAME.MoveWASD`.

---

---

## 3D primitives

### Cubes

```basic
// DrawCube(posX, posY, posZ, width, height, length, r, g, b, a)
DrawCube(0, 0, 0, 2, 2, 2, 255, 100, 100, 255)

// Wireframe
DrawCubeWires(0, 0, 0, 2, 2, 2, 255, 255, 255, 255)
```

### Spheres

```basic
// DrawSphere(posX, posY, posZ, radius, r, g, b, a)
DrawSphere(0, 0, 0, 1, 100, 200, 255, 255)

// Wireframe (optional rings, slices)
DrawSphereWires(0, 0, 0, 1, 16, 16, 255, 255, 255, 255)
```

### Planes and grid

```basic
// DrawPlane(centerPos, sizeX, sizeZ, color)
DrawPlane(0, 0, 0, 10, 10, 128, 128, 128, 255)

// Grid (slices, spacing)
DrawGrid(10, 1.0)
```

### Lines and points

```basic
// 3D line: (x1,y1,z1, x2,y2,z2, color)
DrawLine3D(0, 0, 0, 5, 5, 5, 255, 0, 0, 255)

DrawPoint3D(0, 0, 0, 255, 255, 255, 255)
```

---

## Distance fog

Enable distance-based fog so 3D geometry fades with distance. Set fog parameters, then wrap your 3D drawing with **BeginFog()** and **EndFog()**:

```basic
SetFog(1, 0.03, 200, 220, 255)   // enable, density, r, g, b (0–255)
// Or: SetFogDensity(0.03), SetFogColor(200, 220, 255), EnableFog()

BeginFog()
DrawCube(0, 0, 0, 2, 2, 2, 255, 100, 100, 255)
DrawSphere(5, 1, 0, 1, 100, 200, 255, 255)
EndFog()
```

- **SetFog(enable, density, r, g, b):** One call to enable and set density and color. `enable` 1 = on, 0 = off. Density typically 0.02–0.1.
- **SetFogDensity(density), SetFogColor(r, g, b):** Change fog without toggling.
- **EnableFog(), DisableFog(), IsFogEnabled():** Toggle and query.
- **BeginFog() / EndFog():** Applies fog to the following draw calls until EndFog().

---

## Models and meshes

### Loading and drawing models

```basic
VAR model = LoadModel("character.obj")

// DrawModel(modelId, posX, posY, posZ, scale) and optional tint
DrawModel(model, 0, 0, 0, 1.0)
DrawModel(model, 5, 0, 0, 2.0, 255, 255, 255, 255)

UnloadModel(model)
```

Supported formats include `.obj`, `.gltf`, `.glb`, `.iqm`.

### DrawModelEx (position, rotation axis, angle, scale)

```basic
// DrawModelEx(id, posX, posY, posZ, rotX, rotY, rotZ, angleDeg, scaleX, scaleY, scaleZ, tint...)
DrawModelEx(model, 0, 0, 0, 0, 1, 0, 45, 1, 1, 1, 255, 255, 255, 255)
```

### Generated meshes and models

```basic
// Generate mesh (returns mesh id)
VAR cubeMesh = GenMeshCube(2, 2, 2)
VAR sphereMesh = GenMeshSphere(1, 16, 16)
VAR planeMesh = GenMeshPlane(10, 10, 1, 1)

// Create model from mesh
VAR cubeModel = LoadModelFromMesh(cubeMesh)

DrawModel(cubeModel, 0, 0, 0, 1.0)

UnloadModel(cubeModel)
// Meshes can be unloaded with UnloadMesh(meshId) when no longer needed
```

### 3D model animation

Load skeletal animations from a file (e.g. `.iqm`, `.gltf`), get frame count, and either drive frames manually or use **time-based animation state** for playback.

**Load animations and get frame count:**

```basic
VAR model = LoadModel("character.glb")
VAR animCount = LoadModelAnimations("character.glb")
VAR animId = GetModelAnimationId(0)
VAR frameCount = GetModelAnimationFrameCount(animId)

// Manual: advance frame each frame and update model pose
VAR frame = 0
WHILE NOT WindowShouldClose()
    UpdateModelAnimation(model, animId, frame)
    frame = (frame + 1) % frameCount
    // ... draw ...
WEND
```

**Time-based playback (recommended):**

```basic
VAR model = LoadModel("character.glb")
VAR animCount = LoadModelAnimations("character.glb")
VAR animId = GetModelAnimationId(0)
VAR state = CreateModelAnimState(model, animId, 24, TRUE)

WHILE NOT WindowShouldClose()
    UpdateModelAnimState(state, GetFrameTime())
    ClearBackground(20, 20, 30, 255)
    BeginMode3D(...)
    DrawModel(model, 0, 0, 0, 1.0)
    EndMode3D()
WEND

DestroyModelAnimState(state)
UnloadModelAnimations(animId)
UnloadModel(model)
```

**Commands:**

- **GetModelAnimationFrameCount**(animId) → number of frames (int). Use for manual looping: `frame = (frame + 1) % GetModelAnimationFrameCount(animId)`.
- **CreateModelAnimState**(modelId, animId, fps [, loop]) → stateId. Default loop = TRUE.
- **UpdateModelAnimState**(stateId, deltaTime) — advance time and update model pose; call each frame with GetFrameTime().
- **SetModelAnimStateFrame**(stateId, frameIndex) — set current frame by index and update pose.
- **GetModelAnimStateFrame**(stateId) → current frame index.
- **DestroyModelAnimState**(stateId) — remove state.

Rendering is unchanged: call **DrawModel**(modelId, ...) after updating the animation state.

---

## 3D game checklist

Use this checklist to confirm your program is a valid 3D game:

- [ ] **Window:** `InitWindow(width, height, title)` and `SetTargetFPS(60)` (or desired FPS)
- [ ] **Camera:** Either `SetCamera3D(posX, posY, posZ, targetX, targetY, targetZ, upX, upY, upZ)` or, for orbit/follow, `GAME.CameraOrbit(...)` each frame (and optionally `GAME.MoveWASD` for movement)
- [ ] **Loop:** `WHILE NOT WindowShouldClose()` … `WEND` (or `REPEAT` … `UNTIL WindowShouldClose()`). No auto-wrap; code compiles as written.
- [ ] **Clear:** `ClearBackground(r, g, b, a)` at the start of each frame
- [ ] **Draw:** Your 3D primitives and/or `DrawModel` in the loop
- [ ] **Input:** e.g. `IsKeyDown(KEY_W)`, mouse delta for camera, `GetFrameTime()` for delta time
- [ ] **Close:** `CloseWindow()` after the loop

Optional for 3D games:

- **Physics:** `BULLET.CreateWorld`, `BULLET.CreateSphere`/`CreateBox`, `BULLET.Step`, `BULLET.GetPositionX/Y/Z`, and **GAME.MoveWASD**, **GAME.CameraOrbit** for character and camera. See [Game Development Guide](GAME_DEVELOPMENT_GUIDE.md) and [templates/3d_game.bas](../templates/3d_game.bas).

---

## Complete 3D game example

From [templates/3d_game.bas](../templates/3d_game.bas) (Bullet physics + orbit camera):

```basic
RL.InitWindow(1024, 600, "3D Game")
RL.SetTargetFPS(60)
RL.DisableCursor()

BULLET.CreateWorld("w", 0, -18, 0)
BULLET.CreateSphere("w", "player", 0, 0.5, 0, 0.5, 1)
BULLET.CreateBox("w", "ground", 0, -0.5, 0, 12.5, 0.5, 12.5, 0)

VAR camAngle = 0
VAR camDist = 10
VAR dt = 0.016

REPEAT
  LET dt = RL.GetFrameTime()
  IF dt > 0.05 THEN LET dt = 0.016
  LET delta = RL.GetMouseDelta()
  LET camAngle = camAngle - delta.x * 0.002

  BULLET.Step("w", dt)
  GAME.MoveWASD("w", "player", camAngle, 120, 9, dt)

  LET px = BULLET.GetPositionX("w", "player")
  LET py = BULLET.GetPositionY("w", "player")
  LET pz = BULLET.GetPositionZ("w", "player")
  GAME.CameraOrbit(px, py + 1.5, pz, camAngle, 0.2, camDist)

  RL.ClearBackground(RL.SkyBlue)
  RL.DrawCube(0, -0.5, 0, 25, 1, 25, RL.DarkGreen)
  RL.DrawSphere(px, py, pz, 0.5, RL.Red)
  RL.DrawText("WASD move, Mouse look, Space jump", 10, 10, 20, RL.White)
UNTIL RL.WindowShouldClose()

RL.EnableCursor()
RL.CloseWindow()
```

Run it: `cyberbasic templates/3d_game.bas`

---

## 3D editor and level builder

You can build a 3D editor or level builder using mouse picking, ground/plane picking, grid snap, and level objects.

**Typical loop:**

1. **Camera** – Use `SetCamera3D` or `GAME.CameraOrbit` (e.g. orbit around a fixed target for editor view).
2. **Picking** – Call `GetMouseRay()` each frame, then use the ray origin/direction (GetMouseRayOriginX/Y/Z, GetMouseRayDirectionX/Y/Z) with `GetRayCollisionSphere`, `GetRayCollisionBox`, or `GetRayCollisionMesh` to select an object. For placement on the ground, use `PickGroundPlane()` then `GetRayCollisionPointX/Y/Z()` for the hit point.
3. **Snap** – Use `SnapToGridX(x, gridSize)`, `SnapToGridY`, `SnapToGridZ` to snap placement to a grid.
4. **Level objects** – `CreateLevelObject(id, modelId, x, y, z, rotX, rotY, rotZ, scaleX, scaleY, scaleZ)` to add; `SetLevelObjectTransform` to move/rotate/scale; `DrawLevelObject(id)` for each object; `SaveLevel(path)` and `LoadLevel(path)` to persist the level. Use `DuplicateLevelObject(id)` to clone an object.
5. **Visuals** – `DrawGrid(slices, spacing)` for a grid; `DrawModelWires` or `DrawBoundingBox` for selection highlight. **Camera readback:** `GetCameraPositionX/Y/Z`, `GetCameraTargetX/Y/Z` for saving the view.

See [API_REFERENCE.md](../API_REFERENCE.md) (3D editor and level builder) for all commands.

---

## Full 3D command reference

All 3D-relevant commands in one place. See [API_REFERENCE.md](../API_REFERENCE.md) for details.

- **Window / screen:** `InitWindow`, `SetTargetFPS`, `GetScreenWidth`, `GetScreenHeight`, `GetScreenCenterX`, `GetScreenCenterY`, `CloseWindow`, `WindowShouldClose`
- **Frame:** `ClearBackground` and your draw calls in the loop. No auto frame or mode injection.
- **Camera:** `SetCamera3D`; `GAME.CameraOrbit`, `GAME.SetCamera3DOrbit`, `GAME.UpdateCamera3D`
- **Primitives:** `DrawCube`, `DrawCubeWires`, `DrawSphere`, `DrawSphereWires`, `DrawPlane`, `DrawGrid`, `DrawLine3D`, `DrawPoint3D`
- **Models:** `LoadModel`, `DrawModel`, `DrawModelEx`, `UnloadModel`; `LoadModelFromMesh`, `GenMeshCube`, `GenMeshSphere`, `GenMeshPlane`, `UnloadMesh`
- **Fog:** `SetFog`, `SetFogDensity`, `SetFogColor`, `EnableFog`, `DisableFog`, `IsFogEnabled`, `BeginFog`, `EndFog`
- **Math / distance:** `Distance3D`, `Vector3Distance`, `Vector3Lerp`
- **Game loop / movement:** `DeltaTime`, `GetFrameTime`, `WHILE NOT WindowShouldClose() … WEND`; `GAME.MoveWASD`, `GAME.SnapToGround3D`
- **2D overlay:** `DrawText`, `DrawTexture`, etc. (same as 2D)
- **Editor / level builder:** `GetMouseRay`, GetMouseRayOriginX/Y/Z, GetMouseRayDirectionX/Y/Z, `GetRayCollisionPlane`, `PickGroundPlane`, `SnapToGridX/Y/Z`, CreateLevelObject, SetLevelObjectTransform, GetLevelObject*, DeleteLevelObject, GetLevelObjectCount, GetLevelObjectId, DrawLevelObject, SaveLevel, LoadLevel, DuplicateLevelObject, GetCameraPositionX/Y/Z, GetCameraTargetX/Y/Z
- **Scene:** `CreateScene`, `LoadScene`, `SetCurrentScene` — manage scenes; see [API Reference](../API_REFERENCE.md).
- **Billboards:** `DrawBillboard`, `DrawBillboardRec`, `DrawBillboardPro` — texture-always-facing-camera; see [API Reference](../API_REFERENCE.md).

---

## See also

- [API Reference](../API_REFERENCE.md) — full list of 3D and camera functions
- [Game Development Guide](GAME_DEVELOPMENT_GUIDE.md) — 3D physics, GAME.CameraOrbit, GAME.MoveWASD
- [3D Physics Guide](3D_PHYSICS_GUIDE.md) — Bullet worlds, bodies, forces
- [Command Reference](COMMAND_REFERENCE.md) — commands by feature
- [examples/README.md](../examples/README.md) — 3D examples (e.g. run_3d_physics_demo.bas)
