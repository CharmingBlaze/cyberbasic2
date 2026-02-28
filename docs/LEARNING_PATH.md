# CyberBasic Learning Path - From Zero to Game Developer

Welcome to your complete journey into game development with CyberBasic! This guide takes you from absolute beginner to advanced game developer, teaching you everything you need to create 2D games, 3D games, GUI applications, and multiplayer experiences.

## What You'll Learn

- **Module 1**: BASIC Programming Fundamentals
- **Module 2**: 2D Game Development 
- **Module 3**: 3D Game Development
- **Module 4**: Physics & Animation
- **Module 5**: GUI & User Interfaces
- **Module 6**: Multiplayer Games
- **Module 7**: Advanced Topics

---

## Module 1: BASIC Programming Fundamentals

### Lesson 1.1: Your First Program
**Goal**: Understand the absolute basics - variables, printing, and simple operations

```basic
// Hello World - your first CyberBasic program
PRINT "Welcome to CyberBasic!"
PRINT "Let's start game development together!"

// Variables - storing information
VAR name = "Player"
VAR score = 0
VAR lives = 3

// Basic math
VAR bonus = 100
score = score + bonus
PRINT "Your score: " + STR(score)
```

**What you learned:**
- `PRINT` displays text and numbers
- `VAR` creates variables (use `LET` to update them)
- `STR()` converts numbers to strings for printing
- Basic math operations (+, -, *, /)

### Lesson 1.2: Making Decisions
**Goal**: Control your program's flow with conditions

```basic
VAR age = 16
VAR hasTicket = 1

IF age >= 18 THEN
    PRINT "You can enter the club!"
ELSEIF age >= 16 AND hasTicket = 1 THEN
    PRINT "You can enter with supervision!"
ELSE
    PRINT "Sorry, you're too young"
ENDIF

// Multiple conditions with SELECT CASE
VAR day = 3
SELECT CASE day
    CASE 1: PRINT "Monday"
    CASE 2: PRINT "Tuesday" 
    CASE 3: PRINT "Wednesday"
    CASE ELSE: PRINT "Other day"
END SELECT
```

**What you learned:**
- `IF...THEN...ELSE...ENDIF` for conditional logic
- Comparison operators: =, <>, <, >, <=, >=
- Logical operators: AND, OR, NOT
- `SELECT CASE` for multiple conditions

### Lesson 1.3: Loops and Repetition
**Goal**: Repeat actions efficiently

```basic
// FOR loops - repeat a specific number of times
PRINT "Counting to 5:"
FOR i = 1 TO 5
    PRINT "Count: " + STR(i)
NEXT i

// WHILE loops - repeat while condition is true
VAR counter = 0
WHILE counter < 3
    PRINT "While loop iteration: " + STR(counter)
    counter = counter + 1
WEND

// REPEAT loops - repeat until condition becomes true
VAR attempts = 0
REPEAT
    attempts = attempts + 1
    PRINT "Attempt: " + STR(attempts)
UNTIL attempts >= 3
```

**What you learned:**
- `FOR...NEXT` loops for counted repetition
- `WHILE...WEND` loops for condition-based repetition  
- `REPEAT...UNTIL` loops that always run at least once

### Lesson 1.4: Functions and Organization
**Goal**: Reuse code with functions and modules

```basic
// Creating your own functions
FUNCTION CalculateDamage(base, level)
    RETURN base * (1 + level * 0.1)
END FUNCTION

// Using functions
VAR playerLevel = 5
VAR baseDamage = 10
VAR totalDamage = CalculateDamage(baseDamage, playerLevel)
PRINT "Total damage: " + STR(totalDamage)

// Subroutines for actions
SUB ShowHealth(current, maximum)
    VAR percent = (current / maximum) * 100
    PRINT "Health: " + STR(current) + "/" + STR(maximum) + " (" + STR(percent) + "%)"
END SUB

ShowHealth(75, 100)
```

**What you learned:**
- `FUNCTION...END FUNCTION` returns values
- `SUB...END SUB` performs actions
- Parameters pass data to functions
- `RETURN` sends values back

---

## Module 2: 2D Game Development

### Lesson 2.1: Your First Graphics Window
**Goal**: Create a window and draw basic shapes

```basic
// Window setup - every graphics program needs this
InitWindow(800, 600, "My First Game")
SetTargetFPS(60)  // Smooth 60 FPS animation

// Main game loop
WHILE NOT WindowShouldClose()
    // Clear screen with dark blue
    ClearBackground(20, 20, 40, 255)
    
    // Draw shapes
    DrawCircle(400, 300, 50, 255, 100, 100, 255)  // Red circle
    DrawRectangle(100, 100, 200, 150, 100, 255, 100, 255)  // Green rectangle
    DrawText("Hello Graphics!", 300, 50, 30, 255, 255, 255, 255)
    
    // Show FPS
    DrawText("FPS: " + STR(GetFPS()), 10, 10, 20, 255, 255, 0, 255)
WEND

// Clean up
CloseWindow()
```

**What you learned:**
- `InitWindow()` creates the game window
- `SetTargetFPS()` controls frame rate
- `WindowShouldClose()` detects when user wants to quit
- `ClearBackground()` fills screen with color
- Drawing commands for shapes and text
- RGB color values (0-255)

### Lesson 2.2: Player Movement and Input
**Goal**: Make things move with keyboard input

```basic
InitWindow(800, 600, "Move the Player!")
SetTargetFPS(60)

// Player position
VAR playerX = 400.0
VAR playerY = 300.0
VAR playerSpeed = 200.0  // pixels per second

WHILE NOT WindowShouldClose()
    // Get time since last frame
    VAR dt = GetFrameTime()
    
    // Movement with smooth delta-time
    IF IsKeyDown(KEY_RIGHT) THEN
        playerX = playerX + playerSpeed * dt
    ENDIF
    IF IsKeyDown(KEY_LEFT) THEN
        playerX = playerX - playerSpeed * dt
    ENDIF
    IF IsKeyDown(KEY_UP) THEN
        playerY = playerY - playerSpeed * dt
    ENDIF
    IF IsKeyDown(KEY_DOWN) THEN
        playerY = playerY + playerSpeed * dt
    ENDIF
    
    // Keep player on screen
    playerX = Clamp(playerX, 20, 780)
    playerY = Clamp(playerY, 20, 580)
    
    // Drawing
    ClearBackground(30, 30, 50, 255)
    DrawCircle(playerX, playerY, 20, 100, 200, 255, 255)
    DrawText("Use Arrow Keys to Move", 10, 10, 20, 255, 255, 255, 255)
WEND

CloseWindow()
```

**What you learned:**
- `GetFrameTime()` for smooth, frame-rate independent movement
- `IsKeyDown()` for continuous key detection
- `Clamp()` to limit values within a range
- Delta-time movement for consistent speed regardless of FPS

### Lesson 2.3: Game Objects and Collision
**Goal**: Create multiple objects and detect when they touch

```basic
InitWindow(800, 600, "Collision Detection")
SetTargetFPS(60)

// Player
VAR playerX = 400
VAR playerY = 300
VAR playerRadius = 20

// Collectible items
VAR coinX[5] = [100, 200, 300, 400, 500]
VAR coinY[5] = [200, 150, 250, 180, 220]
VAR coinCollected[5] = [0, 0, 0, 0, 0]
VAR score = 0

WHILE NOT WindowShouldClose()
    // Player movement
    IF IsKeyDown(KEY_LEFT) THEN playerX = playerX - 5
    IF IsKeyDown(KEY_RIGHT) THEN playerX = playerX + 5
    IF IsKeyDown(KEY_UP) THEN playerY = playerY - 5
    IF IsKeyDown(KEY_DOWN) THEN playerY = playerY + 5
    
    // Check collision with each coin
    FOR i = 0 TO 4
        IF coinCollected[i] = 0 THEN
            // Distance calculation
            VAR dx = playerX - coinX[i]
            VAR dy = playerY - coinY[i]
            VAR distance = Sqrt(dx * dx + dy * dy)
            
            // Collision if distance < sum of radii
            IF distance < playerRadius + 15 THEN
                coinCollected[i] = 1
                score = score + 10
            ENDIF
        ENDIF
    NEXT i
    
    // Drawing
    ClearBackground(40, 40, 60, 255)
    
    // Draw player
    DrawCircle(playerX, playerY, playerRadius, 100, 150, 255, 255)
    
    // Draw coins
    FOR i = 0 TO 4
        IF coinCollected[i] = 0 THEN
            DrawCircle(coinX[i], coinY[i], 15, 255, 215, 0, 255)
        ENDIF
    NEXT i
    
    // UI
    DrawText("Score: " + STR(score), 10, 10, 30, 255, 255, 255, 255)
    DrawText("Arrow keys to move, collect coins!", 10, 50, 20, 200, 200, 200, 255)
WEND

CloseWindow()
```

**What you learned:**
- Arrays for multiple game objects
- Distance-based collision detection
- `Sqrt()` for mathematical calculations
- Game state management (collected items)
- Score tracking and display

### Lesson 2.4: Complete 2D Game
**Goal**: Put it all together in a complete game

```basic
// Complete 2D Game: Space Shooter
InitWindow(800, 600, "Space Shooter")
SetTargetFPS(60)

// Game state
VAR gameOver = 0
VAR score = 0
VAR level = 1

// Player
VAR playerX = 400
VAR playerY = 500
VAR playerSpeed = 300

// Enemies
VAR enemyX[10]
VAR enemyY[10]
VAR enemySpeed[10]
VAR enemyActive[10]

// Initialize enemies
FOR i = 0 TO 9
    enemyX[i] = GetRandomValue(50, 750)
    enemyY[i] = GetRandomValue(-500, -50)
    enemySpeed[i] = 50 + GetRandomValue(0, 100)
    enemyActive[i] = 1
NEXT i

// Bullets
VAR bulletX[20]
VAR bulletY[20]
VAR bulletActive[20]

WHILE NOT WindowShouldClose() AND gameOver = 0
    VAR dt = GetFrameTime()
    
    // Player movement
    IF IsKeyDown(KEY_LEFT) THEN playerX = playerX - playerSpeed * dt
    IF IsKeyDown(KEY_RIGHT) THEN playerX = playerX + playerSpeed * dt
    IF IsKeyDown(KEY_UP) THEN playerY = playerY - playerSpeed * dt
    IF IsKeyDown(KEY_DOWN) THEN playerY = playerY + playerSpeed * dt
    
    // Keep player on screen
    playerX = Clamp(playerX, 20, 780)
    playerY = Clamp(playerY, 20, 580)
    
    // Shooting
    IF IsKeyPressed(KEY_SPACE) THEN
        // Find inactive bullet
        FOR i = 0 TO 19
            IF bulletActive[i] = 0 THEN
                bulletX[i] = playerX
                bulletY[i] = playerY - 20
                bulletActive[i] = 1
                EXIT FOR
            ENDIF
        NEXT i
    ENDIF
    
    // Update bullets
    FOR i = 0 TO 19
        IF bulletActive[i] = 1 THEN
            bulletY[i] = bulletY[i] - 400 * dt
            IF bulletY[i] < 0 THEN
                bulletActive[i] = 0
            ENDIF
        ENDIF
    NEXT i
    
    // Update enemies
    FOR i = 0 TO 9
        IF enemyActive[i] = 1 THEN
            enemyY[i] = enemyY[i] + enemySpeed[i] * dt
            IF enemyY[i] > 650 THEN
                enemyY[i] = GetRandomValue(-500, -50)
                enemyX[i] = GetRandomValue(50, 750)
            ENDIF
        ENDIF
    NEXT i
    
    // Check collisions
    // Bullets vs enemies
    FOR b = 0 TO 19
        IF bulletActive[b] = 1 THEN
            FOR e = 0 TO 9
                IF enemyActive[e] = 1 THEN
                    VAR dx = bulletX[b] - enemyX[e]
                    VAR dy = bulletY[b] - enemyY[e]
                    IF Sqrt(dx * dx + dy * dy) < 25 THEN
                        bulletActive[b] = 0
                        enemyActive[e] = 0
                        score = score + 100
                    ENDIF
                ENDIF
            NEXT e
        ENDIF
    NEXT b
    
    // Player vs enemies
    FOR i = 0 TO 9
        IF enemyActive[i] = 1 THEN
            VAR dx = playerX - enemyX[i]
            VAR dy = playerY - enemyY[i]
            IF Sqrt(dx * dx + dy * dy) < 30 THEN
                gameOver = 1
            ENDIF
        ENDIF
    NEXT i
    
    // Drawing
    ClearBackground(10, 10, 30, 255)
    
    // Draw player (triangle shape)
    DrawTriangle(playerX, playerY - 15, playerX - 10, playerY + 15, playerX + 10, playerY + 15, 100, 200, 255, 255)
    
    // Draw enemies
    FOR i = 0 TO 9
        IF enemyActive[i] = 1 THEN
            DrawRectangle(enemyX[i] - 10, enemyY[i] - 10, 20, 20, 255, 100, 100, 255)
        ENDIF
    NEXT i
    
    // Draw bullets
    FOR i = 0 TO 19
        IF bulletActive[i] = 1 THEN
            DrawCircle(bulletX[i], bulletY[i], 3, 255, 255, 100, 255)
        ENDIF
    NEXT i
    
    // UI
    DrawText("Score: " + STR(score), 10, 10, 25, 255, 255, 255, 255)
    DrawText("Arrow keys: move, Space: shoot", 10, 40, 16, 200, 200, 200, 255)
    
WEND

// Game over screen
IF gameOver = 1 THEN
    ClearBackground(50, 10, 10, 255)
    DrawText("GAME OVER", 300, 250, 50, 255, 100, 100, 255)
    DrawText("Final Score: " + STR(score), 320, 320, 25, 255, 255, 255, 255)
    DrawText("Press any key to exit...", 280, 380, 16, 200, 200, 200, 255)
    
    // Wait for key press before closing
    WHILE NOT WindowShouldClose()
        IF IsKeyPressed(KEY_ANY) THEN
            EXIT WHILE
        ENDIF
    WEND
ENDIF

CloseWindow()
PRINT "Game completed! Final score: " + STR(score)
```

**What you learned:**
- Complete game structure with game states
- Multiple object types (player, enemies, bullets)
- Complex collision detection
- Score tracking and game over conditions
- Organized game loop with clear sections

---

## Module 3: 3D Game Development

### Lesson 3.1: Your First 3D Scene
**Goal**: Enter the world of 3D graphics

```basic
InitWindow(1024, 768, "My First 3D Scene")
SetTargetFPS(60)

// Set up 3D camera
// Position (10, 10, 10), Look at (0, 0, 0), Up direction (0, 1, 0)
SetCamera3D(10, 10, 10, 0, 0, 0, 0, 1, 0)

WHILE NOT WindowShouldClose()
    // Begin 3D drawing mode
    BeginMode3D()
    
    // Clear background
    ClearBackground(100, 149, 237, 255)  // Sky blue
    
    // Draw 3D objects
    DrawCube(0, 0, 0, 2, 2, 2, 255, 0, 0, 255)      // Red cube
    DrawSphere(4, 0, 0, 1, 0, 255, 0, 255)           // Green sphere
    DrawPlane(0, -1, 0, 10, 10, 100, 100, 100, 255) // Gray ground
    
    // Draw grid for reference
    DrawGrid(10, 1.0)
    
    // End 3D mode
    EndMode3D()
    
    // Draw 2D text on top
    DrawText("Welcome to 3D!", 10, 10, 30, 255, 255, 255, 255)
    DrawText("Red cube, Green sphere, Gray ground", 10, 50, 20, 200, 200, 200, 255)
WEND

CloseWindow()
```

**What you learned:**
- `SetCamera3D()` positions the 3D camera
- `BeginMode3D()` and `EndMode3D()` wrap 3D drawing
- 3D primitives: `DrawCube()`, `DrawSphere()`, `DrawPlane()`
- `DrawGrid()` for spatial reference
- 3D coordinates (X, Y, Z)

### Lesson 3.2: 3D Camera Control
**Goal**: Move the camera around your 3D world

```basic
InitWindow(1024, 768, "3D Camera Control")
SetTargetFPS(60)
DisableCursor()  // Hide mouse for FPS controls

// Camera parameters
VAR cameraX = 10.0
VAR cameraY = 10.0
VAR cameraZ = 10.0
VAR cameraYaw = 0.0
VAR cameraPitch = 0.0

WHILE NOT WindowShouldClose()
    VAR dt = GetFrameTime()
    
    // Mouse look
    cameraYaw = cameraYaw + GetMouseDeltaX() * 0.003
    cameraPitch = cameraPitch + GetMouseDeltaY() * 0.003
    cameraPitch = Clamp(cameraPitch, -1.5, 1.5)  // Limit up/down look
    
    // Movement based on camera direction
    VAR moveSpeed = 5.0
    VAR forwardX = Sin(cameraYaw) * Cos(cameraPitch)
    VAR forwardZ = Cos(cameraYaw) * Cos(cameraPitch)
    VAR rightX = Cos(cameraYaw)
    VAR rightZ = -Sin(cameraYaw)
    
    IF IsKeyDown(KEY_W) THEN
        cameraX = cameraX + forwardX * moveSpeed * dt
        cameraZ = cameraZ + forwardZ * moveSpeed * dt
    ENDIF
    IF IsKeyDown(KEY_S) THEN
        cameraX = cameraX - forwardX * moveSpeed * dt
        cameraZ = cameraZ - forwardZ * moveSpeed * dt
    ENDIF
    IF IsKeyDown(KEY_A) THEN
        cameraX = cameraX - rightX * moveSpeed * dt
        cameraZ = cameraZ - rightZ * moveSpeed * dt
    ENDIF
    IF IsKeyDown(KEY_D) THEN
        cameraX = cameraX + rightX * moveSpeed * dt
        cameraZ = cameraZ + rightZ * moveSpeed * dt
    ENDIF
    
    // Calculate look-at position
    VAR lookX = cameraX + forwardX
    VAR lookY = cameraY + Sin(cameraPitch)
    VAR lookZ = cameraZ + forwardZ
    
    // Update camera
    SetCamera3D(cameraX, cameraY, cameraZ, lookX, lookY, lookZ, 0, 1, 0)
    
    // 3D Drawing
    BeginMode3D()
    ClearBackground(135, 206, 235, 255)
    
    // Draw scene objects
    DrawCube(0, 1, 0, 2, 2, 2, 255, 100, 100, 255)
    DrawSphere(5, 1, 0, 1.5, 100, 255, 100, 255)
    DrawCube(-5, 1, -5, 3, 1, 3, 100, 100, 255, 255)
    DrawPlane(0, 0, 0, 20, 20, 150, 150, 150, 255)
    DrawGrid(20, 1.0)
    
    EndMode3D()
    
    // 2D Overlay
    DrawText("WASD: Move, Mouse: Look", 10, 10, 20, 255, 255, 255, 255)
    DrawText("ESC to exit", 10, 35, 16, 200, 200, 200, 255)
WEND

CloseWindow()
```

**What you learned:**
- First-person camera controls
- Mouse delta for smooth looking
- 3D movement math with trigonometry
- `DisableCursor()` for mouse control
- Camera position and look-at calculation

### Lesson 3.3: 3D Models and Textures
**Goal**: Load and display 3D models with textures

```basic
InitWindow(1024, 768, "3D Models and Textures")
SetTargetFPS(60)

// Load 3D model (you'll need a .obj or .gltf file)
VAR modelId = LoadModel("models/cube.obj")  // Replace with actual model path
VAR textureId = LoadTexture("models/texture.png")  // Replace with actual texture

// Set model texture
SetModelTexture(modelId, textureId)

// Camera orbit around center
VAR cameraAngle = 0.0
VAR cameraDistance = 15.0
VAR cameraHeight = 10.0

WHILE NOT WindowShouldClose()
    VAR dt = GetFrameTime()
    
    // Rotate camera around model
    cameraAngle = cameraAngle + dt * 0.5
    
    // Calculate camera position for orbit
    VAR camX = Cos(cameraAngle) * cameraDistance
    VAR camZ = Sin(cameraAngle) * cameraDistance
    
    SetCamera3D(camX, cameraHeight, camZ, 0, 2, 0, 0, 1, 0)
    
    // 3D Drawing
    BeginMode3D()
    ClearBackground(120, 120, 150, 255)
    
    // Draw loaded model
    DrawModel(modelId, 0, 2, 0, 1.0)  // Scale = 1.0
    
    // Draw some primitive shapes too
    DrawCube(5, 1, 0, 2, 2, 2, 255, 200, 100, 255)
    DrawSphere(-5, 1, 0, 1.5, 100, 200, 255, 255)
    DrawPlane(0, 0, 0, 30, 30, 100, 100, 100, 255)
    DrawGrid(30, 1.0)
    
    EndMode3D()
    
    // UI
    DrawText("3D Model Showcase", 10, 10, 25, 255, 255, 255, 255)
    DrawText("Camera auto-orbits around center", 10, 40, 16, 200, 200, 200, 255)
WEND

// Clean up resources
UnloadModel(modelId)
UnloadTexture(textureId)
CloseWindow()
```

**What you learned:**
- Loading 3D models with `LoadModel()`
- Loading textures with `LoadTexture()`
- Applying textures to models
- Resource management (unload when done)
- Camera orbit for model viewing

---

## Module 4: Physics & Animation

### Lesson 4.1: 2D Physics with Box2D
**Goal**: Add realistic physics to your 2D games

```basic
InitWindow(800, 600, "2D Physics Demo")
SetTargetFPS(60)

// Create physics world with gravity
CreateWorld2D("world", 0, -10)  // Gravity pointing down

// Create ground (static body)
CreateBody2D("world", "ground", 0, 0, 0, 0, 1, 50, 0.5)

// Create some dynamic boxes
CreateBody2D("world", "box1", 2, 0, 5, 0, 1, 0.5, 0.5)  // Dynamic box
CreateBody2D("world", "box2", 2, 2, 8, 0, 1, 0.5, 0.5)  // Another box
CreateBody2D("world", "ball1", 2, -2, 7, 0, 1, 0.5)     // Circle

// Physics to screen conversion
VAR scale = 50  // 1 physics unit = 50 pixels
VAR offsetX = 400
VAR offsetY = 500

WHILE NOT WindowShouldClose()
    // Step physics simulation
    Step2D("world", 0.016, 8, 3)  // 16ms timestep, 8 velocity iterations, 3 position iterations
    
    // Get physics positions and draw
    ClearBackground(40, 44, 52, 255)
    
    // Draw ground
    DrawRectangle(50, 475, 700, 50, 128, 128, 128, 255)
    
    // Draw box1
    VAR x1 = GetPositionX2D("world", "box1")
    VAR y1 = GetPositionY2D("world", "box1")
    VAR screenX1 = offsetX + x1 * scale
    VAR screenY1 = offsetY - y1 * scale
    DrawRectangle(screenX1 - 25, screenY1 - 25, 50, 50, 255, 100, 100, 255)
    
    // Draw box2
    VAR x2 = GetPositionX2D("world", "box2")
    VAR y2 = GetPositionY2D("world", "box2")
    VAR screenX2 = offsetX + x2 * scale
    VAR screenY2 = offsetY - y2 * scale
    DrawRectangle(screenX2 - 25, screenY2 - 25, 50, 50, 100, 255, 100, 255)
    
    // Draw ball
    VAR xb = GetPositionX2D("world", "ball1")
    VAR yb = GetPositionY2D("world", "ball1")
    VAR screenXb = offsetX + xb * scale
    VAR screenYb = offsetY - yb * scale
    DrawCircle(screenXb, screenYb, 25, 100, 100, 255, 255)
    
    // UI
    DrawText("2D Physics Demo - Box2D", 10, 10, 20, 255, 255, 255, 255)
    DrawText("Objects fall and collide realistically", 10, 35, 16, 200, 200, 200, 255)
WEND

// Clean up
DestroyWorld2D("world")
CloseWindow()
```

**What you learned:**
- Creating physics worlds with gravity
- Static vs dynamic bodies
- Physics simulation stepping
- Converting physics coordinates to screen coordinates
- Different physics shapes (boxes, circles)

### Lesson 4.2: 3D Physics with Bullet
**Goal**: Add realistic 3D physics

```basic
InitWindow(1024, 768, "3D Physics Demo")
SetTargetFPS(60)

// Create 3D physics world
CreateWorld3D("world", 0, -9.81, 0)  // Earth gravity

// Create ground (static)
CreateBox3D("world", "ground", 0, -1, 0, 10, 0.5, 10, 0)

// Create falling objects
CreateBox3D("world", "box1", 0, 5, 0, 1, 1, 1, 1)
CreateSphere3D("world", "sphere1", 2, 8, 0, 0.5, 1)
CreateBox3D("world", "box2", -2, 6, 1, 0.5, 2, 0.5, 1)

// Camera
SetCamera3D(10, 8, 10, 0, 0, 0, 0, 1, 0)

WHILE NOT WindowShouldClose()
    // Step physics
    Step3D("world", 0.016)
    
    // 3D Drawing
    BeginMode3D()
    ClearBackground(135, 206, 235, 255)
    
    // Draw ground
    DrawCube(0, -1, 0, 20, 1, 20, 100, 100, 100, 255)
    
    // Draw physics objects at their current positions
    // Box1
    VAR bx1 = GetPositionX3D("world", "box1")
    VAR by1 = GetPositionY3D("world", "box1")
    VAR bz1 = GetPositionZ3D("world", "box1")
    DrawCube(bx1, by1, bz1, 2, 2, 2, 255, 100, 100, 255)
    
    // Sphere1
    VAR sx1 = GetPositionX3D("world", "sphere1")
    VAR sy1 = GetPositionY3D("world", "sphere1")
    VAR sz1 = GetPositionZ3D("world", "sphere1")
    DrawSphere(sx1, sy1, sz1, 0.5, 100, 100, 255, 255)
    
    // Box2
    VAR bx2 = GetPositionX3D("world", "box2")
    VAR by2 = GetPositionY3D("world", "box2")
    VAR bz2 = GetPositionZ3D("world", "box2")
    DrawCube(bx2, by2, bz2, 1, 4, 1, 100, 255, 100, 255)
    
    DrawGrid(20, 1.0)
    EndMode3D()
    
    // UI
    DrawText("3D Physics Demo - Bullet", 10, 10, 20, 255, 255, 255, 255)
    DrawText("Objects fall and collide in 3D space", 10, 35, 16, 200, 200, 200, 255)
WEND

// Clean up
DestroyWorld3D("world")
CloseWindow()
```

**What you learned:**
- 3D physics world creation
- Different 3D physics shapes
- Getting 3D positions from physics
- Drawing objects at physics positions
- 3D physics simulation

---

## Module 5: GUI & User Interfaces

### Lesson 5.1: Basic GUI Elements
**Goal**: Create user interfaces with buttons, sliders, and text

```basic
InitWindow(600, 400, "GUI Demo")
SetTargetFPS(60)

// GUI state variables
VAR sliderValue = 0.5
VAR checkboxState = 0
VAR buttonClicked = 0
VAR textInput = "Hello"

WHILE NOT WindowShouldClose()
    ClearBackground(60, 60, 80, 255)
    
    // Begin GUI mode
    BeginUI()
    
    // Labels
    Label("=== GUI Controls Demo ===")
    Label("Slider Value: " + STR(sliderValue))
    
    // Slider
    sliderValue = Slider("Adjust Me", sliderValue, 0, 1)
    
    // Checkbox
    checkboxState = Checkbox("Enable Feature", checkboxState)
    IF checkboxState = 1 THEN
        Label("✓ Feature is enabled!")
    ENDIF
    
    // Button
    buttonClicked = Button("Click Me!")
    IF buttonClicked THEN
        PRINT "Button was clicked!"
    ENDIF
    
    // Text input (basic)
    Label("Current text: " + textInput)
    
    // End GUI mode
    EndUI()
    
    // Additional drawing outside GUI
    DrawText("ESC to exit", 10, 360, 16, 200, 200, 200, 255)
WEND

CloseWindow()
```

**What you learned:**
- `BeginUI()` and `EndUI()` for GUI mode
- Basic GUI controls: Label, Button, Slider, Checkbox
- GUI state management
- Mixing GUI with regular drawing

### Lesson 5.2: Game Menu System
**Goal**: Create a complete game menu

```basic
InitWindow(800, 600, "Game Menu")
SetTargetFPS(60)

// Menu state
VAR menuState = "main"  // "main", "options", "game", "credits"
VAR selectedOption = 0
VAR musicVolume = 0.7
VAR soundVolume = 0.8
VAR fullscreen = 0

WHILE NOT WindowShouldClose()
    ClearBackground(20, 20, 40, 255)
    
    BeginUI()
    
    SELECT CASE menuState
        CASE "main"
            Label("=== CYBERBASIC GAME ===")
            Label("")
            
            IF Button("Start Game") THEN
                menuState = "game"
            ENDIF
            
            IF Button("Options") THEN
                menuState = "options"
            ENDIF
            
            IF Button("Credits") THEN
                menuState = "credits"
            ENDIF
            
            IF Button("Exit") THEN
                EXIT WHILE
            ENDIF
            
        CASE "options"
            Label("=== OPTIONS ===")
            Label("")
            
            Label("Music Volume")
            musicVolume = Slider("Music", musicVolume, 0, 1)
            
            Label("Sound Volume")
            soundVolume = Slider("Sound", soundVolume, 0, 1)
            
            fullscreen = Checkbox("Fullscreen", fullscreen)
            
            Label("")
            IF Button("Back") THEN
                menuState = "main"
            ENDIF
            
        CASE "credits"
            Label("=== CREDITS ===")
            Label("")
            Label("Game developed with CyberBasic")
            Label("Graphics: Raylib")
            Label("Physics: Box2D & Bullet")
            Label("")
            Label("Thank you for playing!")
            Label("")
            
            IF Button("Back") THEN
                menuState = "main"
            ENDIF
            
        CASE "game"
            Label("=== GAME IN PROGRESS ===")
            Label("")
            Label("This is where your game would be!")
            Label("")
            Label("Current Settings:")
            Label("Music: " + STR(Int(musicVolume * 100)) + "%")
            Label("Sound: " + STR(Int(soundVolume * 100)) + "%")
            Label("Fullscreen: " + IIF(fullscreen = 1, "ON", "OFF"))
            Label("")
            
            IF Button("Back to Menu") THEN
                menuState = "main"
            ENDIF
    END SELECT
    
    EndUI()
    
WEND

CloseWindow()
```

**What you learned:**
- Menu state management
- Navigation between different screens
- Settings persistence
- Complex UI layouts

---

## Module 6: Multiplayer Games

### Lesson 6.1: Basic Network Connection
**Goal**: Connect two computers over the network

```basic
// Server code - save as server.bas
VAR server = Host(9999)  // Host on port 9999
IF IsNull(server) THEN
    PRINT "Failed to host server"
    QUIT
ENDIF

PRINT "Server hosting on port 9999..."
PRINT "Waiting for client to connect..."

VAR client = Accept(server, 5000)  // Wait 5 seconds for client
IF IsNull(client) THEN
    PRINT "No client connected"
    CloseServer(server)
    QUIT
ENDIF

PRINT "Client connected!"

// Communication loop
FOR i = 1 TO 5
    Send(client, "Message " + STR(i) + " from server")
    VAR response = Receive(client, 1000)  // Wait 1 second for response
    IF NOT IsNull(response) THEN
        PRINT "Client says: " + response
    ENDIF
NEXT i

Disconnect(client)
CloseServer(server)
PRINT "Server session ended"
```

```basic
// Client code - save as client.bas
VAR connection = Connect("127.0.0.1", 9999)  // Connect to localhost
IF IsNull(connection) THEN
    PRINT "Failed to connect to server"
    QUIT
ENDIF

PRINT "Connected to server!"

// Communication loop
FOR i = 1 TO 5
    VAR message = Receive(connection, 1000)  // Wait 1 second
    IF NOT IsNull(message) THEN
        PRINT "Server says: " + message
        Send(connection, "Response " + STR(i) + " from client")
    ENDIF
NEXT i

Disconnect(connection)
PRINT "Client session ended"
```

**What you learned:**
- Server-client architecture
- `Host()` and `Connect()` functions
- `Send()` and `Receive()` for communication
- Basic network error handling

### Lesson 6.2: Multiplayer Game
**Goal**: Create a simple multiplayer game

```basic
// This would be a complete multiplayer game example
// For brevity, showing the structure:

InitWindow(800, 600, "Multiplayer Game")
SetTargetFPS(60)

// Network setup
VAR isServer = 1  // Set to 0 for client
VAR connection

IF isServer = 1 THEN
    connection = Host(9999)
    PRINT "Server waiting for player..."
    VAR client = Accept(connection, 30000)
    IF IsNull(client) THEN
        PRINT "No one connected"
        QUIT
    ENDIF
    PRINT "Player connected!"
ELSE
    connection = Connect("127.0.0.1", 9999)
    IF IsNull(connection) THEN
        PRINT "Could not connect to server"
        QUIT
    ENDIF
    PRINT "Connected to game!"
ENDIF

// Game state
VAR myX = 100
VAR myY = 300
VAR otherX = 700
VAR otherY = 300

WHILE NOT WindowShouldClose()
    // Handle input
    IF IsKeyDown(KEY_LEFT) THEN myX = myX - 3
    IF IsKeyDown(KEY_RIGHT) THEN myX = myX + 3
    IF IsKeyDown(KEY_UP) THEN myY = myY - 3
    IF IsKeyDown(KEY_DOWN) THEN myY = myY + 3
    
    // Send position
    Send(connection, "POS:" + STR(myX) + "," + STR(myY))
    
    // Receive other player position
    VAR msg = Receive(connection, 100)
    IF NOT IsNull(msg) THEN
        IF LEFT(msg, 4) = "POS:" THEN
            VAR posData = MID(msg, 5)
            // Parse position (simplified)
            // In real code, you'd parse the comma-separated values
            otherX = 400  // Placeholder
            otherY = 300  // Placeholder
        ENDIF
    ENDIF
    
    // Drawing
    ClearBackground(40, 40, 60, 255)
    DrawCircle(myX, myY, 20, 100, 200, 255, 255)    // Me (blue)
    DrawCircle(otherX, otherY, 20, 255, 100, 100, 255) // Other (red)
    
    DrawText("You are blue, other player is red", 10, 10, 16, 255, 255, 255, 255)
WEND

// Cleanup
IF isServer = 1 THEN
    CloseServer(connection)
ELSE
    Disconnect(connection)
ENDIF
CloseWindow()
```

**What you learned:**
- Real-time network communication
- Synchronizing game state
- Server vs client responsibilities
- Network game loop structure

---

## Module 7: Advanced Topics

### Lesson 7.1: Entity-Component System (ECS)
**Goal**: Organize complex games with ECS

```basic
// Create ECS world
VAR world = ECS.CreateWorld()

// Create entities
VAR player = ECS.CreateEntity(world)
VAR enemy = ECS.CreateEntity(world)

// Add components to player
ECS.AddComponent(world, player, "position", "{ \"x\": 100, \"y\": 200 }")
ECS.AddComponent(world, player, "velocity", "{ \"x\": 0, \"y\": 0 }")
ECS.AddComponent(world, player, "health", "{ \"value\": 100 }")

// Add components to enemy
ECS.AddComponent(world, enemy, "position", "{ \"x\": 300, \"y\": 200 }")
ECS.AddComponent(world, enemy, "health", "{ \"value\": 50 }")

// Game loop with ECS
WHILE NOT WindowShouldClose()
    // Query all entities with position and velocity
    VAR movingEntities = ECS.Query(world, "position,velocity")
    
    // Update positions based on velocity
    FOR EACH entity IN movingEntities
        VAR pos = ECS.GetComponent(world, entity, "position")
        VAR vel = ECS.GetComponent(world, entity, "velocity")
        
        // Update position (simplified)
        ECS.SetComponent(world, entity, "position", "{ \"x\": " + STR(pos.x + vel.x) + " }")
    NEXT
    
    // Render all entities with position
    VAR renderableEntities = ECS.Query(world, "position")
    FOR EACH entity IN renderableEntities
        VAR pos = ECS.GetComponent(world, entity, "position")
        DrawCircle(pos.x, pos.y, 20, 255, 255, 255, 255)
    NEXT
WEND

// Cleanup
ECS.DestroyWorld(world)
```

### Lesson 7.2: Data Persistence
**Goal**: Save and load game data

```basic
// Save game
SUB SaveGame(filename, score, level, playerHealth)
    VAR gameData = "{"
    gameData = gameData + "\"score\": " + STR(score) + ","
    gameData = gameData + "\"level\": " + STR(level) + ","
    gameData = gameData + "\"playerHealth\": " + STR(playerHealth)
    gameData = gameData + "}"
    
    VAR success = WriteFile(filename, gameData)
    IF success THEN
        PRINT "Game saved successfully!"
    ELSE
        PRINT "Failed to save game"
    ENDIF
END SUB

// Load game
FUNCTION LoadGame(filename)
    VAR data = ReadFile(filename)
    IF IsNull(data) THEN
        PRINT "No save file found"
        RETURN "{ \"score\": 0, \"level\": 1, \"playerHealth\": 100 }"
    ENDIF
    
    PRINT "Game loaded successfully!"
    RETURN data
END FUNCTION

// Usage
VAR currentScore = 1500
VAR currentLevel = 3
VAR currentHealth = 75

// Save
SaveGame("savegame.json", currentScore, currentLevel, currentHealth)

// Load
VAR saveData = LoadGame("savegame.json")
VAR loadedScore = GetJSONKey(saveData, "score")
PRINT "Loaded score: " + STR(loadedScore)
```

---

## Graduation Projects

Now that you've completed all modules, try these capstone projects:

### Project 1: Complete 2D Platformer
- Player with jump and double-jump
- Multiple enemy types
- Collectible coins and power-ups
- Level progression system
- Sound effects and music
- Save/load system

### Project 2: 3D Exploration Game
- First-person camera controls
- Large 3D world to explore
- Collectible items scattered around
- Simple physics puzzles
- Day/night cycle
- Inventory system

### Project 3: Multiplayer Arena Game
- 2-4 player support
- Different character classes
- Power-ups and weapons
- Score tracking
- Chat system
- Spectator mode

### Project 4: GUI-Heavy Strategy Game
- Complex menu system
- Resource management
- Building placement
- Unit control
- Save/load multiple games
- Settings and options

---

## Additional Resources

### Documentation References
- [API Reference](../API_REFERENCE.md) - Complete function listing
- [Language Spec](../LANGUAGE_SPEC.md) - Detailed language features
- [Examples](../examples/README.md) - Working code examples

### Community and Support
- GitHub Repository: [github.com/CharmingBlaze/cyberbasic2](https://github.com/CharmingBlaze/cyberbasic2)
- Report issues and contribute there

### Best Practices
1. **Start Simple**: Begin with basic concepts before adding complexity
2. **Test Often**: Run your code frequently to catch mistakes early
3. **Use Comments**: Document your code for future reference
4. **Organize Code**: Use functions and modules to keep code clean
5. **Learn from Examples**: Study the provided examples for techniques

---

## Congratulations!

You've completed the CyberBasic Learning Path! You now have:

- **Solid programming foundation** with BASIC syntax
- **2D game development skills** - graphics, input, collision
- **3D game development skills** - cameras, models, rendering
- **Physics integration** - both 2D and 3D
- **GUI development** - menus, interfaces, user interaction
- **Multiplayer capabilities** - network programming
- **Advanced techniques** - ECS, data persistence, optimization

You're ready to create amazing games! Start with small projects and gradually work your way up to more complex games. The game development journey is endless - keep learning, keep creating, and most importantly, have fun!

**Happy coding with CyberBasic!**
