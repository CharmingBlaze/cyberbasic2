# 3D Graphics Guide

Complete guide to 3D graphics in CyberBasic using the raylib API.

## Table of Contents

1. [Getting started](#getting-started)
2. [3D camera](#3d-camera)
3. [3D drawing frame](#3d-drawing-frame)
4. [3D primitives](#3d-primitives)
5. [Models and meshes](#models-and-meshes)
6. [3D game checklist](#3d-game-checklist)
7. [Complete 3D game example](#complete-3d-game-example)

---

## Getting started

### Basic 3D setup

Every 3D program needs a window, a 3D camera, and a loop that draws between `BeginMode3D()` and `EndMode3D()`:

```basic
InitWindow(800, 600, "My 3D Game")
SetTargetFPS(60)

// Set camera: position (x,y,z), target (x,y,z), up vector (x,y,z)
SetCamera3D(0, 10, 10,  0, 0, 0,  0, 1, 0)

WHILE NOT WindowShouldClose()
    BeginDrawing()
    ClearBackground(20, 20, 30, 255)

    BeginMode3D()
        // All 3D drawing here
        DrawCube(0, 0, 0, 2, 2, 2, 255, 100, 100, 255)
        DrawSphere(5, 1, 0, 1, 100, 200, 255, 255)
    EndMode3D()

    // 2D overlay (UI)
    DrawText("3D Scene", 10, 10, 20, 255, 255, 255, 255)

    EndDrawing()
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

After that, all drawing inside `BeginMode3D()` / `EndMode3D()` uses this camera until you call `SetCamera3D` again.

### Orbit camera (GAME.CameraOrbit)

For a third-person or orbit camera driven by mouse/keyboard, use **GAME.CameraOrbit** each frame. It sets the raylib 3D camera from a target position and orbit angles:

```basic
// Each frame, after updating your target position (e.g. player):
GAME.CameraOrbit(targetX, targetY, targetZ, angleRad, pitchRad, distance)
```

Then call `BeginMode3D()` and draw as usual. See [Game Development Guide](GAME_DEVELOPMENT_GUIDE.md) and [templates/3d_game.bas](../templates/3d_game.bas) for a full loop with Bullet physics and `GAME.MoveWASD`.

---

## 3D drawing frame

All 3D drawing must happen between `BeginMode3D()` and `EndMode3D()`, and that block must be inside `BeginDrawing()` … `EndDrawing()`:

```basic
BeginDrawing()
    ClearBackground(50, 50, 60, 255)

    BeginMode3D()
        // 3D primitives and models here
    EndMode3D()

    // 2D UI here
    DrawText("FPS: " + STR(GetFPS()), 10, 10, 20, 255, 255, 255, 255)

EndDrawing()
```

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

---

## 3D game checklist

Use this checklist to confirm your program is a valid 3D game:

- [ ] **Window:** `InitWindow(width, height, title)` and `SetTargetFPS(60)` (or desired FPS)
- [ ] **Camera:** Either `SetCamera3D(posX, posY, posZ, targetX, targetY, targetZ, upX, upY, upZ)` or, for orbit/follow, `GAME.CameraOrbit(...)` each frame (and optionally `GAME.MoveWASD` for movement)
- [ ] **Loop:** `Main() ... EndMain` or `WHILE NOT WindowShouldClose()` (or `REPEAT` … `UNTIL WindowShouldClose()`); Main() and the WHILE form auto-wrap with BeginDrawing/EndDrawing
- [ ] **Clear:** `ClearBackground(r, g, b, a)` at the start of each frame
- [ ] **3D block:** `BeginMode3D()` … `EndMode3D()` with primitives and/or `DrawModel` inside
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

  RL.BeginDrawing()
  RL.ClearBackground(RL.SkyBlue)
  RL.BeginMode3D()
  RL.DrawCube(0, -0.5, 0, 25, 1, 25, RL.DarkGreen)
  RL.DrawSphere(px, py, pz, 0.5, RL.Red)
  RL.EndMode3D()
  RL.DrawText("WASD move, Mouse look, Space jump", 10, 10, 20, RL.White)
  RL.EndDrawing()
UNTIL RL.WindowShouldClose()

RL.EnableCursor()
RL.CloseWindow()
```

Run it: `cyberbasic templates/3d_game.bas`

---

For more 3D examples see [examples/README.md](../examples/README.md) (e.g. run_3d_physics_demo.bas, mario64.bas). For the full list of 3D and camera functions see [API_REFERENCE.md](../API_REFERENCE.md).
