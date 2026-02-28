# 3D Game Development Tutorial - Complete Guide

Welcome to the complete 3D game development tutorial! This guide will take you from basic 3D concepts to creating fully-featured 3D games with physics, cameras, and advanced rendering.

## What You'll Build

By the end of this tutorial, you'll have created:
- A complete 3D exploration game
- First-person camera controls
- 3D physics simulation
- Model loading and texturing
- Lighting and shading effects
- 3D particle systems
- Advanced camera techniques

---

## Prerequisites

Before starting, make sure you've completed:
- **Module 1**: BASIC Programming Fundamentals (from LEARNING_PATH.md)
- **Module 2**: 2D Game Development (recommended for foundation)
- Basic understanding of coordinates and vectors

---

## Lesson 1: Understanding 3D Space

### 3D Coordinate System

In 3D graphics, we work with three axes:
- **X-axis**: Left to right (horizontal)
- **Y-axis**: Bottom to top (vertical) 
- **Z-axis**: Near to far (depth)

The origin (0,0,0) is at the center of your 3D world.

```basic
// 3D Coordinate System Demonstration
InitWindow(1024, 768, "3D Coordinate System")
SetTargetFPS(60)

// Set up camera to view the origin from an angle
SetCamera3D(10, 10, 10, 0, 0, 0, 0, 1, 0)

WHILE NOT WindowShouldClose()
    BeginMode3D()
    ClearBackground(100, 149, 237, 255)
    
    // Draw coordinate axes
    // X-axis (red) - points right
    DrawLine3D(0, 0, 0, 3, 0, 0, 255, 0, 0, 255)
    // Y-axis (green) - points up  
    DrawLine3D(0, 0, 0, 0, 3, 0, 0, 255, 0, 255)
    // Z-axis (blue) - points forward
    DrawLine3D(0, 0, 0, 0, 0, 3, 0, 0, 255, 255)
    
    // Draw objects at different positions
    DrawCube(2, 0, 0, 0.5, 0.5, 0.5, 255, 0, 0, 255)     // Red cube on X-axis
    DrawCube(0, 2, 0, 0.5, 0.5, 0.5, 0, 255, 0, 255)     // Green cube on Y-axis
    DrawCube(0, 0, 2, 0.5, 0.5, 0.5, 0, 0, 255, 255)     // Blue cube on Z-axis
    
    // Draw origin sphere
    DrawSphere(0, 0, 0, 0.3, 255, 255, 255, 255)
    
    // Draw grid for reference
    DrawGrid(10, 1.0)
    
    EndMode3D()
    
    // 2D text overlay
    DrawText("3D Coordinate System", 10, 10, 20, 255, 255, 255, 255)
    DrawText("Red: X-axis, Green: Y-axis, Blue: Z-axis", 10, 35, 16, 200, 200, 200, 255)
    DrawText("White sphere: Origin (0,0,0)", 10, 55, 16, 200, 200, 200, 255)
WEND

CloseWindow()
```

### Essential 3D Primitives

```basic
// 3D Primitives Showcase
InitWindow(1024, 768, "3D Primitives")
SetTargetFPS(60)

SetCamera3D(8, 8, 8, 0, 0, 0, 0, 1, 0)

WHILE NOT WindowShouldClose()
    BeginMode3D()
    ClearBackground(135, 206, 235, 255)
    
    // Basic 3D shapes
    DrawCube(-3, 0, 0, 1.5, 1.5, 1.5, 255, 100, 100, 255)      // Red cube
    DrawSphere(0, 0, 0, 1, 100, 255, 100, 255)                    // Green sphere
    DrawCylinder(3, 0, 0, 1, 2, 100, 100, 255, 255)              // Blue cylinder
    DrawTorus(-3, 0, 3, 0.8, 0.3, 255, 255, 100, 255)            // Yellow torus
    DrawCone(0, 0, 3, 1, 2, 255, 100, 255, 255)                  // Magenta cone
    DrawKnot(3, 0, 3, 1, 100, 255, 255, 255)                     // Cyan knot
    
    // Plane as ground
    DrawPlane(0, -1, 0, 10, 10, 150, 150, 150, 255)
    
    // Wireframe versions
    DrawCubeWires(-3, 0, -3, 1.5, 1.5, 1.5, 255, 255, 255, 255)
    DrawSphereWires(0, 0, -3, 1, 255, 255, 255, 255)
    
    DrawGrid(10, 1.0)
    EndMode3D()
    
    // Labels
    DrawText("3D Primitives Showcase", 10, 10, 20, 255, 255, 255, 255)
    DrawText("Solid shapes in first row, wireframes in second", 10, 35, 16, 200, 200, 200, 255)
WEND

CloseWindow()
```

---

## Lesson 2: Camera Systems

### Static Camera

```basic
// Static Camera Setup
InitWindow(1024, 768, "Static Camera")
SetTargetFPS(60)

// Camera parameters
VAR camX = 10.0
VAR camY = 8.0
VAR camZ = 10.0
VAR targetX = 0.0
VAR targetY = 0.0
VAR targetZ = 0.0
VAR upX = 0.0
VAR upY = 1.0
VAR upZ = 0.0

WHILE NOT WindowShouldClose()
    // Update camera (static in this example)
    SetCamera3D(camX, camY, camZ, targetX, targetY, targetZ, upX, upY, upZ)
    
    BeginMode3D()
    ClearBackground(100, 149, 237, 255)
    
    // Draw scene
    DrawCube(0, 1, 0, 2, 2, 2, 255, 100, 100, 255)
    DrawSphere(3, 1, 0, 1, 100, 255, 100, 255)
    DrawCube(-3, 1, 0, 1.5, 1.5, 1.5, 100, 100, 255, 255)
    DrawPlane(0, 0, 0, 20, 20, 120, 120, 120, 255)
    DrawGrid(20, 1.0)
    
    EndMode3D()
    
    DrawText("Static Camera View", 10, 10, 20, 255, 255, 255, 255)
    DrawText("Camera position: (" + STR(camX) + ", " + STR(camY) + ", " + STR(camZ) + ")", 10, 35, 16, 200, 200, 200, 255)
WEND

CloseWindow()
```

### First-Person Camera Controls

```basic
// First-Person Camera with Mouse Look
InitWindow(1024, 768, "First-Person Camera")
SetTargetFPS(60)
DisableCursor()  // Hide mouse for free look

// Camera state
VAR camX = 0.0
VAR camY = 2.0
VAR camZ = 5.0
VAR camYaw = 0.0    // Horizontal rotation
VAR camPitch = 0.0  // Vertical rotation

WHILE NOT WindowShouldClose()
    VAR dt = GetFrameTime()
    
    // Mouse look
    camYaw = camYaw + GetMouseDeltaX() * 0.003
    camPitch = camPitch + GetMouseDeltaY() * 0.003
    camPitch = Clamp(camPitch, -1.5, 1.5)  // Limit vertical look
    
    // Movement speed
    VAR moveSpeed = 5.0
    
    // Calculate forward and right vectors from yaw
    VAR forwardX = Sin(camYaw)
    VAR forwardZ = Cos(camYaw)
    VAR rightX = Cos(camYaw)
    VAR rightZ = -Sin(camYaw)
    
    // Movement input
    IF IsKeyDown(KEY_W) THEN
        camX = camX + forwardX * moveSpeed * dt
        camZ = camZ + forwardZ * moveSpeed * dt
    ENDIF
    IF IsKeyDown(KEY_S) THEN
        camX = camX - forwardX * moveSpeed * dt
        camZ = camZ - forwardZ * moveSpeed * dt
    ENDIF
    IF IsKeyDown(KEY_A) THEN
        camX = camX - rightX * moveSpeed * dt
        camZ = camZ - rightZ * moveSpeed * dt
    ENDIF
    IF IsKeyDown(KEY_D) THEN
        camX = camX + rightX * moveSpeed * dt
        camZ = camZ + rightZ * moveSpeed * dt
    ENDIF
    
    // Up/down movement
    IF IsKeyDown(KEY_SPACE) THEN
        camY = camY + moveSpeed * dt
    ENDIF
    IF IsKeyDown(KEY_LEFT_SHIFT) THEN
        camY = camY - moveSpeed * dt
    ENDIF
    
    // Calculate look-at position
    VAR lookX = camX + Cos(camPitch) * Sin(camYaw)
    VAR lookY = camY + Sin(camPitch)
    VAR lookZ = camZ + Cos(camPitch) * Cos(camYaw)
    
    // Update camera
    SetCamera3D(camX, camY, camZ, lookX, lookY, lookZ, 0, 1, 0)
    
    // 3D Rendering
    BeginMode3D()
    ClearBackground(135, 206, 235, 255)
    
    // Draw scene objects
    DrawCube(0, 1, 0, 2, 2, 2, 255, 100, 100, 255)
    DrawCube(5, 1, 5, 1.5, 3, 1.5, 100, 255, 100, 255)
    DrawCube(-5, 1, -5, 1, 1, 1, 100, 100, 255, 255)
    DrawSphere(3, 1, -3, 1.5, 255, 255, 100, 255)
    
    // Ground plane
    DrawPlane(0, 0, 0, 50, 50, 100, 100, 100, 255)
    DrawGrid(50, 1.0)
    
    EndMode3D()
    
    // 2D Overlay
    DrawText("First-Person Camera", 10, 10, 20, 255, 255, 255, 255)
    DrawText("WASD: Move horizontally", 10, 35, 16, 200, 200, 200, 255)
    DrawText("Mouse: Look around", 10, 55, 16, 200, 200, 200, 255)
    DrawText("Space/Shift: Move up/down", 10, 75, 16, 200, 200, 200, 255)
    DrawText("Position: (" + STR(Int(camX)) + ", " + STR(Int(camY)) + ", " + STR(Int(camZ)) + ")", 10, 95, 16, 200, 200, 200, 255)
WEND

CloseWindow()
```

### Orbit Camera

```basic
// Orbit Camera for Object Viewing
InitWindow(1024, 768, "Orbit Camera")
SetTargetFPS(60)

// Orbit parameters
VAR targetX = 0.0
VAR targetY = 1.0
VAR targetZ = 0.0
VAR camDistance = 10.0
VAR camAngle = 0.0
VAR camHeight = 5.0

WHILE NOT WindowShouldClose()
    VAR dt = GetFrameTime()
    
    // Auto-rotate or manual control
    IF IsKeyDown(KEY_LEFT) THEN
        camAngle = camAngle - dt * 1.0
    ENDIF
    IF IsKeyDown(KEY_RIGHT) THEN
        camAngle = camAngle + dt * 1.0
    ENDIF
    
    // Zoom control
    IF IsKeyDown(KEY_UP) THEN
        camDistance = camDistance - dt * 5.0
    ENDIF
    IF IsKeyDown(KEY_DOWN) THEN
        camDistance = camDistance + dt * 5.0
    ENDIF
    
    camDistance = Clamp(camDistance, 3.0, 20.0)
    
    // Height control
    IF IsKeyDown(KEY_W) THEN
        camHeight = camHeight + dt * 3.0
    ENDIF
    IF IsKeyDown(KEY_S) THEN
        camHeight = camHeight - dt * 3.0
    ENDIF
    
    // Calculate camera position
    VAR camX = targetX + Cos(camAngle) * camDistance
    VAR camZ = targetZ + Sin(camAngle) * camDistance
    
    // Update camera
    SetCamera3D(camX, camHeight, camZ, targetX, targetY, targetZ, 0, 1, 0)
    
    // 3D Rendering
    BeginMode3D()
    ClearBackground(100, 149, 237, 255)
    
    // Draw central object
    DrawCube(targetX, targetY, targetZ, 2, 2, 2, 255, 100, 100, 255)
    DrawCubeWires(targetX, targetY, targetZ, 2, 2, 2, 255, 255, 255, 255)
    
    // Draw surrounding objects
    DrawSphere(3, 1, 0, 1, 100, 255, 100, 255)
    DrawSphere(-3, 1, 0, 1, 100, 100, 255, 255)
    DrawCube(0, 1, 3, 1, 1, 1, 255, 255, 100, 255)
    DrawCube(0, 1, -3, 1, 1, 1, 255, 100, 255, 255)
    
    // Ground
    DrawPlane(0, 0, 0, 20, 20, 120, 120, 120, 255)
    DrawGrid(20, 1.0)
    
    EndMode3D()
    
    // UI
    DrawText("Orbit Camera", 10, 10, 20, 255, 255, 255, 255)
    DrawText("Arrow Keys: Rotate orbit", 10, 35, 16, 200, 200, 200, 255)
    DrawText("Up/Down: Zoom in/out", 10, 55, 16, 200, 200, 200, 255)
    DrawText("W/S: Camera height", 10, 75, 16, 200, 200, 200, 255)
    DrawText("Distance: " + STR(camDistance), 10, 95, 16, 200, 200, 200, 255)
WEND

CloseWindow()
```

---

## Lesson 3: 3D Physics Integration

### Basic 3D Physics with Bullet

```basic
// 3D Physics with Bullet Engine
InitWindow(1024, 768, "3D Physics Demo")
SetTargetFPS(60)

// Create physics world
BULLET.CreateWorld("world", 0, -9.81, 0)  // Earth gravity

// Create static ground
BULLET.CreateBox("world", "ground", 0, -1, 0, 10, 0.5, 10, 0)

// Create dynamic objects
BULLET.CreateBox("world", "box1", 0, 5, 0, 1, 1, 1, 1)
BULLET.CreateSphere("world", "sphere1", 2, 8, 0, 0.5, 1)
BULLET.CreateBox("world", "box2", -2, 6, 1, 0.5, 2, 0.5, 1)

// Camera
SetCamera3D(10, 8, 10, 0, 0, 0, 0, 1, 0)

WHILE NOT WindowShouldClose()
    VAR dt = GetFrameTime()
    
    // Step physics simulation
    BULLET.Step("world", dt)
    
    BeginMode3D()
    ClearBackground(135, 206, 235, 255)
    
    // Draw ground
    DrawCube(0, -1, 0, 20, 1, 20, 100, 100, 100, 255)
    
    // Draw physics objects at their current positions
    // Box1
    VAR bx1 = BULLET.GetPositionX("world", "box1")
    VAR by1 = BULLET.GetPositionY("world", "box1")
    VAR bz1 = BULLET.GetPositionZ("world", "box1")
    DrawCube(bx1, by1, bz1, 2, 2, 2, 255, 100, 100, 255)
    
    // Sphere1
    VAR sx1 = BULLET.GetPositionX("world", "sphere1")
    VAR sy1 = BULLET.GetPositionY("world", "sphere1")
    VAR sz1 = BULLET.GetPositionZ("world", "sphere1")
    DrawSphere(sx1, sy1, sz1, 0.5, 100, 100, 255, 255)
    
    // Box2
    VAR bx2 = BULLET.GetPositionX("world", "box2")
    VAR by2 = BULLET.GetPositionY("world", "box2")
    VAR bz2 = BULLET.GetPositionZ("world", "box2")
    DrawCube(bx2, by2, bz2, 1, 4, 1, 100, 255, 100, 255)
    
    DrawGrid(20, 1.0)
    EndMode3D()
    
    // UI
    DrawText("3D Physics Demo", 10, 10, 20, 255, 255, 255, 255)
    DrawText("Objects fall and collide realistically", 10, 35, 16, 200, 200, 200, 255)
    DrawText("Click to spawn more objects", 10, 55, 16, 200, 200, 200, 255)
    
    // Spawn new objects on click
    IF IsMouseButtonPressed(0) THEN
        VAR mx = GetMouseX()
        VAR my = GetMouseY()
        // Convert screen coordinates to 3D world coordinates
        VAR ray = GetMouseRay(mx, my)
        VAR spawnX = ray.direction.x * 5.0
        VAR spawnY = 5.0
        VAR spawnZ = ray.direction.z * 5.0
        
        // Random shape
        VAR shapeType = GetRandomValue(0, 2)
        VAR objName = "dynamic" + STR(GetRandomValue(1000, 9999))
        
        SELECT CASE shapeType
            CASE 0: BULLET.CreateBox("world", objName, spawnX, spawnY, spawnZ, 0.5, 0.5, 0.5, 1)
            CASE 1: BULLET.CreateSphere("world", objName, spawnX, spawnY, spawnZ, 0.5, 1)
            CASE 2: BULLET.CreateBox("world", objName, spawnX, spawnY, spawnZ, 0.3, 1.0, 0.3, 1)
        END SELECT
    ENDIF
WEND

// Cleanup
BULLET.DestroyWorld("world")
CloseWindow()
```

### Interactive Physics Player

```basic
// 3D Physics Player Controller
InitWindow(1024, 768, "3D Physics Player")
SetTargetFPS(60)
DisableCursor()

// Physics world
BULLET.CreateWorld("world", 0, -9.81, 0)

// Ground
BULLET.CreateBox("world", "ground", 0, -1, 0, 20, 0.5, 20, 0)

// Platforms
BULLET.CreateBox("world", "platform1", 5, 2, 0, 3, 0.5, 3, 0)
BULLET.CreateBox("world", "platform2", -5, 3, -3, 2, 0.5, 2, 0)
BULLET.CreateBox("world", "platform3", 0, 4, 5, 2.5, 0.5, 2.5, 0)

// Player physics body
BULLET.CreateBox("world", "player", 0, 2, 0, 0.5, 1, 0.5, 1)

// Camera state
VAR camYaw = 0.0
VAR camPitch = 0.0

WHILE NOT WindowShouldClose()
    VAR dt = GetFrameTime()
    
    // Mouse look
    camYaw = camYaw + GetMouseDeltaX() * 0.003
    camPitch = camPitch + GetMouseDeltaY() * 0.003
    camPitch = Clamp(camPitch, -1.5, 1.5)
    
    // Player movement
    VAR moveSpeed = 8.0
    VAR jumpPower = 10.0
    
    // Calculate movement direction
    VAR forwardX = Sin(camYaw)
    VAR forwardZ = Cos(camYaw)
    VAR rightX = Cos(camYaw)
    VAR rightZ = -Sin(camYaw)
    
    VAR moveX = 0
    VAR moveZ = 0
    
    IF IsKeyDown(KEY_W) THEN moveX = moveX + forwardX
    IF IsKeyDown(KEY_S) THEN moveX = moveX - forwardX
    IF IsKeyDown(KEY_A) THEN moveZ = moveZ - rightX
    IF IsKeyDown(KEY_D) THEN moveZ = moveZ + rightX
    
    // Apply movement to player
    IF moveX <> 0 OR moveZ <> 0 THEN
        BULLET.SetVelocity("world", "player", moveX * moveSpeed, 0, moveZ * moveSpeed)
    ENDIF
    
    // Jump
    IF IsKeyPressed(KEY_SPACE) THEN
        VAR currentY = BULLET.GetPositionY("world", "player")
        // Simple ground check
        IF currentY < 3.0 THEN
            BULLET.SetVelocity("world", "player", 0, jumpPower, 0)
        ENDIF
    ENDIF
    
    // Get player position for camera
    VAR playerX = BULLET.GetPositionX("world", "player")
    VAR playerY = BULLET.GetPositionY("world", "player")
    VAR playerZ = BULLET.GetPositionZ("world", "player")
    
    // Calculate camera look position
    VAR lookX = playerX + Cos(camPitch) * Sin(camYaw)
    VAR lookY = playerY + Sin(camPitch)
    VAR lookZ = playerZ + Cos(camPitch) * Cos(camYaw)
    
    // Update camera (first person)
    SetCamera3D(playerX, playerY + 0.8, playerZ, lookX, lookY, lookZ, 0, 1, 0)
    
    // Physics step
    BULLET.Step("world", dt)
    
    // Rendering
    BeginMode3D()
    ClearBackground(135, 206, 235, 255)
    
    // Draw ground
    DrawCube(0, -1, 0, 40, 1, 40, 100, 100, 100, 255)
    
    // Draw platforms
    DrawCube(5, 2, 0, 6, 1, 6, 150, 75, 0, 255)
    DrawCube(-5, 3, -3, 4, 1, 4, 150, 75, 0, 255)
    DrawCube(0, 4, 5, 5, 1, 5, 150, 75, 0, 255)
    
    // Draw player (third person view for debugging)
    DrawCube(playerX, playerY, playerZ, 1, 2, 1, 100, 200, 255, 255)
    
    DrawGrid(20, 1.0)
    EndMode3D()
    
    // UI
    DrawText("3D Physics Player", 10, 10, 20, 255, 255, 255, 255)
    DrawText("WASD: Move, Mouse: Look, Space: Jump", 10, 35, 16, 200, 200, 200, 255)
    DrawText("Position: (" + STR(Int(playerX)) + ", " + STR(Int(playerY)) + ", " + STR(Int(playerZ)) + ")", 10, 55, 16, 200, 200, 200, 255)
WEND

BULLET.DestroyWorld("world")
CloseWindow()
```

---

## Lesson 4: Models and Textures

### Loading 3D Models

```basic
// 3D Model Loading and Display
InitWindow(1024, 768, "3D Models")
SetTargetFPS(60)

// Load models (you'll need actual model files)
VAR cubeModel = GenMeshCube(1, 1, 1)
VAR sphereModel = GenMeshSphere(0.5, 16, 16)
VAR cylinderModel = GenMeshCylinder(0.5, 2, 32)

// Convert meshes to models
VAR cube = LoadModelFromMesh(cubeModel)
VAR sphere = LoadModelFromMesh(sphereModel)
VAR cylinder = LoadModelFromMesh(cylinderModel)

// Load textures (you'll need actual texture files)
VAR texture1 = LoadTexture("textures/brick.png")  // Replace with actual path
VAR texture2 = LoadTexture("textures/metal.png")  // Replace with actual path

// Apply textures to models
SetModelTexture(cube, texture1)
SetModelTexture(sphere, texture2)

// Camera orbit
VAR camAngle = 0.0

WHILE NOT WindowShouldClose()
    VAR dt = GetFrameTime()
    camAngle = camAngle + dt * 0.5
    
    // Calculate camera position
    VAR camX = Cos(camAngle) * 10
    VAR camZ = Sin(camAngle) * 10
    SetCamera3D(camX, 8, camZ, 0, 2, 0, 0, 1, 0)
    
    BeginMode3D()
    ClearBackground(100, 149, 237, 255)
    
    // Draw models
    DrawModel(cube, -3, 2, 0, 2.0)
    DrawModel(sphere, 0, 2, 0, 2.0)
    DrawModel(cylinder, 3, 2, 0, 2.0)
    
    // Draw ground
    DrawPlane(0, 0, 0, 20, 20, 120, 120, 120, 255)
    DrawGrid(20, 1.0)
    
    EndMode3D()
    
    // UI
    DrawText("3D Models with Textures", 10, 10, 20, 255, 255, 255, 255)
    DrawText("Cube with brick texture, Sphere with metal texture", 10, 35, 16, 200, 200, 200, 255)
    DrawText("Cylinder with default color", 10, 55, 16, 200, 200, 200, 255)
WEND

// Cleanup
UnloadModel(cube)
UnloadModel(sphere)
UnloadModel(cylinder)
UnloadTexture(texture1)
UnloadTexture(texture2)
CloseWindow()
```

---

## Lesson 5: Complete 3D Game

### 3D Exploration Game

```basic
// Complete 3D Exploration Game
InitWindow(1024, 768, "3D Exploration Game")
SetTargetFPS(60)
DisableCursor()

// Game state
VAR gameState = "playing"  // "playing", "win"
VAR score = 0
VAR collectiblesFound = 0
VAR totalCollectibles = 5

// Physics world
BULLET.CreateWorld("world", 0, -9.81, 0)

// Ground
BULLET.CreateBox("world", "ground", 0, -1, 0, 50, 0.5, 50, 0)

// Player physics body
BULLET.CreateBox("world", "player", 0, 2, 0, 0.5, 1.5, 0.5, 1)

// Collectibles (orbs)
VAR collectibleX[5] = [5, -8, 3, -5, 0]
VAR collectibleY[5] = [3, 4, 6, 2, 8]
VAR collectibleZ[5] = [3, -2, -5, 6, -3]
VAR collectibleCollected[5] = [0, 0, 0, 0, 0]

// Create physics bodies for collectibles
FOR i = 0 TO 4
    VAR orbName = "orb" + STR(i)
    BULLET.CreateSphere("world", orbName, collectibleX[i], collectibleY[i], collectibleZ[i], 0.5, 0.1)
NEXT i

// Obstacles
BULLET.CreateBox("world", "wall1", 10, 2, 0, 1, 4, 10, 0)
BULLET.CreateBox("world", "wall2", -10, 2, 0, 1, 4, 10, 0)
BULLET.CreateBox("world", "wall3", 0, 2, 10, 10, 4, 1, 0)
BULLET.CreateBox("world", "wall4", 0, 2, -10, 10, 4, 1, 0)

// Camera
VAR camYaw = 0.0
VAR camPitch = 0.0

WHILE NOT WindowShouldClose() AND gameState = "playing"
    VAR dt = GetFrameTime()
    
    // Mouse look
    camYaw = camYaw + GetMouseDeltaX() * 0.003
    camPitch = camPitch + GetMouseDeltaY() * 0.003
    camPitch = Clamp(camPitch, -1.5, 1.5)
    
    // Player movement
    VAR moveSpeed = 10.0
    VAR jumpPower = 12.0
    
    VAR forwardX = Sin(camYaw)
    VAR forwardZ = Cos(camYaw)
    VAR rightX = Cos(camYaw)
    VAR rightZ = -Sin(camYaw)
    
    VAR moveX = 0
    VAR moveZ = 0
    
    IF IsKeyDown(KEY_W) THEN moveX = moveX + forwardX
    IF IsKeyDown(KEY_S) THEN moveX = moveX - forwardX
    IF IsKeyDown(KEY_A) THEN moveZ = moveZ - rightX
    IF IsKeyDown(KEY_D) THEN moveZ = moveZ + rightX
    
    // Apply movement
    IF moveX <> 0 OR moveZ <> 0 THEN
        VAR velX = moveX * moveSpeed
        VAR velZ = moveZ * moveSpeed
        BULLET.SetVelocity("world", "player", velX, 0, velZ)
    ENDIF
    
    // Jump
    IF IsKeyPressed(KEY_SPACE) THEN
        VAR currentY = BULLET.GetPositionY("world", "player")
        IF currentY < 5.0 THEN
            BULLET.SetVelocity("world", "player", 0, jumpPower, 0)
        ENDIF
    ENDIF
    
    // Get player position
    VAR playerX = BULLET.GetPositionX("world", "player")
    VAR playerY = BULLET.GetPositionY("world", "player")
    VAR playerZ = BULLET.GetPositionZ("world", "player")
    
    // Check collectible collection
    FOR i = 0 TO 4
        IF collectibleCollected[i] = 0 THEN
            VAR orbName = "orb" + STR(i)
            VAR orbX = BULLET.GetPositionX("world", orbName)
            VAR orbY = BULLET.GetPositionY("world", orbName)
            VAR orbZ = BULLET.GetPositionZ("world", orbName)
            
            VAR dx = playerX - orbX
            VAR dy = playerY - orbY
            VAR dz = playerZ - orbZ
            VAR distance = Sqrt(dx*dx + dy*dy + dz*dz)
            
            IF distance < 1.5 THEN
                collectibleCollected[i] = 1
                collectiblesFound = collectiblesFound + 1
                score = score + 100
                // Remove the orb from physics world
                BULLET.DestroyBody("world", orbName)
            ENDIF
        ENDIF
    NEXT i
    
    // Check win condition
    IF collectiblesFound >= totalCollectibles THEN
        gameState = "win"
    ENDIF
    
    // Camera setup
    VAR lookX = playerX + Cos(camPitch) * Sin(camYaw)
    VAR lookY = playerY + Sin(camPitch)
    VAR lookZ = playerZ + Cos(camPitch) * Cos(camYaw)
    SetCamera3D(playerX, playerY + 1.2, playerZ, lookX, lookY, lookZ, 0, 1, 0)
    
    // Physics step
    BULLET.Step("world", dt)
    
    // Rendering
    BeginMode3D()
    ClearBackground(100, 149, 237, 255)
    
    // Draw ground
    DrawPlane(0, 0, 0, 100, 100, 80, 120, 80, 255)
    DrawGrid(50, 2.0)
    
    // Draw obstacles
    DrawCube(10, 2, 0, 2, 8, 20, 150, 150, 150, 255)
    DrawCube(-10, 2, 0, 2, 8, 20, 150, 150, 150, 255)
    DrawCube(0, 2, 10, 20, 8, 2, 150, 150, 150, 255)
    DrawCube(0, 2, -10, 20, 8, 2, 150, 150, 150, 255)
    
    // Draw collectibles
    FOR i = 0 TO 4
        IF collectibleCollected[i] = 0 THEN
            VAR orbName = "orb" + STR(i)
            VAR orbX = BULLET.GetPositionX("world", orbName)
            VAR orbY = BULLET.GetPositionY("world", orbName)
            VAR orbZ = BULLET.GetPositionZ("world", orbName)
            
            // Animated floating effect
            VAR floatOffset = Sin(GetTime() * 2.0 + i) * 0.2
            DrawSphere(orbX, orbY + floatOffset, orbZ, 0.5, 255, 215, 0, 255)
            DrawSphereWires(orbX, orbY + floatOffset, orbZ, 0.7, 255, 255, 100, 255)
        ENDIF
    NEXT i
    
    // Draw player (third person for debugging - remove in final game)
    DrawCube(playerX, playerY - 0.75, playerZ, 1, 1.5, 1, 100, 200, 255, 255)
    
    EndMode3D()
    
    // UI
    DrawText("3D Exploration Game", 10, 10, 20, 255, 255, 255, 255)
    DrawText("Collect all orbs to win!", 10, 35, 16, 200, 200, 200, 255)
    DrawText("Orbs: " + STR(collectiblesFound) + "/" + STR(totalCollectibles), 10, 55, 16, 255, 215, 0, 255)
    DrawText("Score: " + STR(score), 10, 75, 16, 255, 255, 255, 255)
    DrawText("WASD: Move, Mouse: Look, Space: Jump", 10, 730, 16, 200, 200, 200, 255)
WEND

// Win screen
IF gameState = "win" THEN
    WHILE NOT WindowShouldClose()
        ClearBackground(50, 150, 50, 255)
        DrawText("YOU WIN!", 400, 300, 50, 255, 255, 255, 255)
        DrawText("Final Score: " + STR(score), 420, 380, 25, 255, 255, 255, 255)
        DrawText("Press ESC to exit", 450, 450, 18, 200, 200, 200, 255)
        
        IF IsKeyPressed(KEY_ESCAPE) THEN
            EXIT WHILE
        ENDIF
    WEND
ENDIF

// Cleanup
BULLET.DestroyWorld("world")
CloseWindow()
```

---

## Lesson 6: Advanced 3D Effects

### 3D Particle System

```basic
// 3D Particle System
InitWindow(1024, 768, "3D Particles")
SetTargetFPS(60)

// Particle system
VAR particleX[200]
VAR particleY[200]
VAR particleZ[200]
VAR velocityX[200]
VAR velocityY[200]
VAR velocityZ[200]
VAR life[200]
VAR maxLife[200]
VAR colorR[200]
VAR colorG[200]
VAR colorB[200]
VAR active[200]

VAR particleCount = 0

FUNCTION CreateParticle(x, y, z, vx, vy, vz, lifetime, r, g, b)
    IF particleCount < 200 THEN
        particleX[particleCount] = x
        particleY[particleCount] = y
        particleZ[particleCount] = z
        velocityX[particleCount] = vx
        velocityY[particleCount] = vy
        velocityZ[particleCount] = vz
        life[particleCount] = lifetime
        maxLife[particleCount] = lifetime
        colorR[particleCount] = r
        colorG[particleCount] = g
        colorB[particleCount] = b
        active[particleCount] = 1
        particleCount = particleCount + 1
    ENDIF
END FUNCTION

SUB UpdateParticles(dt)
    FOR i = 0 TO 199
        IF active[i] = 1 THEN
            // Update position
            particleX[i] = particleX[i] + velocityX[i] * dt
            particleY[i] = particleY[i] + velocityY[i] * dt
            particleZ[i] = particleZ[i] + velocityZ[i] * dt
            
            // Update life
            life[i] = life[i] - dt
            
            // Apply gravity
            velocityY[i] = velocityY[i] + 300 * dt
            
            // Remove dead particles
            IF life[i] <= 0 THEN
                active[i] = 0
            ENDIF
        ENDIF
    NEXT i
END SUB

SUB DrawParticles()
    FOR i = 0 TO 199
        IF active[i] = 1 THEN
            VAR alpha = (life[i] / maxLife[i]) * 255
            VAR size = 0.2 * (life[i] / maxLife[i])
            DrawSphere(particleX[i], particleY[i], particleZ[i], size, colorR[i], colorG[i], colorB[i], Int(alpha))
        ENDIF
    NEXT i
END SUB

// Camera
SetCamera3D(10, 8, 10, 0, 0, 0, 0, 1, 0)

WHILE NOT WindowShouldClose()
    VAR dt = GetFrameTime()
    
    // Create explosion on mouse click
    IF IsMouseButtonPressed(0) THEN
        VAR mx = GetMouseX()
        VAR my = GetMouseY()
        VAR ray = GetMouseRay(mx, my)
        
        // Create burst of particles at click point
        FOR i = 0 TO 29
            VAR angle1 = GetRandomValue(0, 360) * 0.0174533  // Convert to radians
            VAR angle2 = GetRandomValue(0, 360) * 0.0174533
            VAR speed = 200 + GetRandomValue(0, 300)
            
            VAR vx = Sin(angle1) * Cos(angle2) * speed
            VAR vy = Sin(angle2) * speed
            VAR vz = Cos(angle1) * Cos(angle2) * speed
            
            VAR r = GetRandomValue(200, 255)
            VAR g = GetRandomValue(50, 200)
            VAR b = GetRandomValue(0, 100)
            
            CreateParticle(ray.direction.x * 5, 5, ray.direction.z * 5, vx, vy, vz, 3.0, r, g, b)
        NEXT i
    ENDIF
    
    // Update particles
    UpdateParticles(dt)
    
    // Rendering
    BeginMode3D()
    ClearBackground(20, 20, 40, 255)
    
    // Draw ground
    DrawPlane(0, 0, 0, 20, 20, 60, 60, 80, 255)
    DrawGrid(20, 1.0)
    
    // Draw particles
    DrawParticles()
    
    EndMode3D()
    
    // UI
    DrawText("3D Particle System", 10, 10, 20, 255, 255, 255, 255)
    DrawText("Click to create particle explosion", 10, 35, 16, 200, 200, 200, 255)
    DrawText("Active particles: " + STR(particleCount), 10, 55, 16, 200, 200, 200, 255)
WEND

CloseWindow()
```

---

## Conclusion

Congratulations! You've now learned:

- **3D coordinate system** and spatial awareness
- **Camera systems** - first-person, orbit, and static cameras
- **3D physics integration** with Bullet physics engine
- **Model loading** and texturing
- **Complete 3D game** with player controller and objectives
- **Advanced 3D effects** like particle systems

### Next Steps

1. **Expand the 3D game**: Add enemies, weapons, and combat
2. **Try different genres**: First-person shooters, racing games, puzzle games
3. **Learn lighting**: Add dynamic lighting and shadows
4. **Create multiplayer**: Network your 3D games
5. **Optimize performance**: Learn about LOD and culling techniques

### Common 3D Game Patterns

- **First-Person Shooters**: Ray casting, projectile physics, weapon systems
- **Third-Person Adventures**: Character controllers, camera follow, animation
- **Racing Games**: Vehicle physics, track design, AI opponents
- **Puzzle Games**: Object manipulation, physics puzzles, spatial reasoning
- **Simulation Games**: Complex systems, resource management, UI integration

Keep practicing and experimenting with these 3D techniques. The transition from 2D to 3D opens up incredible possibilities for game design and player experiences!

**Happy coding!**
