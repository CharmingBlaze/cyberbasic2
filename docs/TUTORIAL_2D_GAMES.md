# 2D Game Development Tutorial - Complete Guide

Welcome to the complete 2D game development tutorial! This guide will take you from basic shapes to a fully-featured 2D game with physics, animation, and polish.

## What You'll Build

By the end of this tutorial, you'll have created:
- A complete 2D platformer game
- Player movement and physics
- Enemy AI and collision
- Collectibles and scoring
- Particle effects and animations
- Sound and music
- Menu system

---

## Prerequisites

Before starting, make sure you've completed:
- **Module 1**: BASIC Programming Fundamentals (from LEARNING_PATH.md)
- Basic understanding of variables, functions, and loops

---

## Lesson 1: Basic 2D Setup

### Understanding the 2D Coordinate System

In CyberBasic (like most graphics systems), the 2D coordinate system works like this:
- **Origin (0,0)** is at the **top-left corner**
- **X increases** going **right**
- **Y increases** going **down**
- Screen coordinates are in **pixels**

```basic
// Coordinate system demonstration
InitWindow(800, 600, "2D Coordinates")
SetTargetFPS(60)

WHILE NOT WindowShouldClose()
    ClearBackground(20, 20, 40, 255)
    
    // Draw axes to show coordinate system
    DrawLine(0, 0, 800, 0, 255, 0, 0, 255)  // X-axis (red)
    DrawLine(0, 0, 0, 600, 0, 255, 0, 255)  // Y-axis (green)
    
    // Draw objects at specific coordinates
    DrawCircle(100, 100, 20, 255, 255, 255, 255)  // Top-left area
    DrawCircle(700, 100, 20, 255, 255, 0, 255)    // Top-right
    DrawCircle(100, 500, 20, 255, 0, 255, 255)    // Bottom-left
    DrawCircle(700, 500, 20, 0, 255, 255, 255)    // Bottom-right
    
    // Center of screen
    DrawCircle(400, 300, 25, 255, 100, 100, 255)
    
    // Labels
    DrawText("(0,0) Top-Left", 10, 10, 16, 255, 255, 255, 255)
    DrawText("(800,600) Bottom-Right", 600, 570, 16, 255, 255, 255, 255)
    DrawText("Center: (400,300)", 340, 330, 16, 255, 255, 255, 255)
WEND

CloseWindow()
```

### Essential 2D Functions

Let's create a reference program showing all essential 2D drawing functions:

```basic
// 2D Drawing Functions Reference
InitWindow(1024, 768, "2D Drawing Reference")
SetTargetFPS(60)

WHILE NOT WindowShouldClose()
    ClearBackground(30, 30, 50, 255)
    
    // Basic shapes
    DrawRectangle(50, 50, 100, 60, 255, 100, 100, 255)
    DrawCircle(200, 80, 30, 100, 255, 100, 255)
    DrawLine(250, 50, 350, 110, 255, 255, 100, 255)
    DrawTriangle(400, 50, 370, 110, 430, 110, 255, 100, 255, 255)
    DrawPixel(500, 80, 255, 255, 255, 255)
    
    // Text rendering
    DrawText("Hello CyberBasic!", 50, 150, 30, 255, 255, 255, 255)
    DrawText("Small text", 50, 190, 16, 200, 200, 200, 255)
    
    // Outline versions
    DrawRectangleLines(50, 250, 100, 60, 255, 200, 100, 255)
    DrawCircleLines(200, 280, 30, 100, 200, 255, 255)
    DrawTriangleLines(400, 250, 370, 310, 430, 310, 255, 200, 255, 255)
    
    // Labels
    DrawText("Filled Shapes", 50, 20, 20, 255, 255, 255, 255)
    DrawText("Text", 50, 130, 16, 255, 255, 255, 255)
    DrawText("Outline Shapes", 50, 220, 20, 255, 255, 255, 255)
    
    // Color palette reference
    DrawText("Colors:", 600, 50, 20, 255, 255, 255, 255)
    DrawRectangle(600, 80, 30, 30, 255, 0, 0, 255)      // Red
    DrawRectangle(640, 80, 30, 30, 0, 255, 0, 255)      // Green
    DrawRectangle(680, 80, 30, 30, 0, 0, 255, 255)      // Blue
    DrawRectangle(720, 80, 30, 30, 255, 255, 0, 255)   // Yellow
    DrawRectangle(760, 80, 30, 30, 255, 0, 255, 255)   // Magenta
    DrawRectangle(800, 80, 30, 30, 0, 255, 255, 255)   // Cyan
    DrawRectangle(840, 80, 30, 30, 255, 255, 255, 255)  // White
    DrawRectangle(880, 80, 30, 30, 0, 0, 0, 255)        // Black
    
WEND

CloseWindow()
```

---

## Lesson 2: Player Movement

### Smooth Movement with Delta Time

Delta time ensures your game runs at the same speed regardless of frame rate:

```basic
// Smooth Player Movement
InitWindow(800, 600, "Smooth Movement")
SetTargetFPS(60)

// Player properties
VAR playerX = 400.0
VAR playerY = 300.0
VAR playerSpeed = 200.0  // pixels per second
VAR playerRadius = 20.0

WHILE NOT WindowShouldClose()
    // Get delta time (time since last frame)
    VAR dt = GetFrameTime()
    
    // Movement with delta time for smooth, frame-rate independent movement
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
    
    // Keep player on screen (boundary checking)
    playerX = Clamp(playerX, playerRadius, 800 - playerRadius)
    playerY = Clamp(playerY, playerRadius, 600 - playerRadius)
    
    // Drawing
    ClearBackground(40, 40, 60, 255)
    
    // Draw player
    DrawCircle(playerX, playerY, playerRadius, 100, 200, 255, 255)
    
    // Draw player direction indicator
    DrawCircle(playerX, playerY, 5, 255, 255, 255, 255)
    
    // UI
    DrawText("Use Arrow Keys to Move", 10, 10, 20, 255, 255, 255, 255)
    DrawText("Position: (" + STR(Int(playerX)) + ", " + STR(Int(playerY)) + ")", 10, 35, 16, 200, 200, 200, 255)
    DrawText("FPS: " + STR(GetFPS()), 10, 55, 16, 200, 200, 200, 255)
WEND

CloseWindow()
```

### Advanced Movement with Acceleration

Let's add acceleration and friction for more realistic movement:

```basic
// Advanced Movement with Physics
InitWindow(800, 600, "Physics Movement")
SetTargetFPS(60)

// Player physics properties
VAR playerX = 400.0
VAR playerY = 300.0
VAR velocityX = 0.0
VAR velocityY = 0.0
VAR acceleration = 500.0   // pixels per second squared
VAR maxSpeed = 300.0      // pixels per second
VAR friction = 0.9        // friction coefficient (0-1)

WHILE NOT WindowShouldClose()
    VAR dt = GetFrameTime()
    
    // Apply acceleration based on input
    IF IsKeyDown(KEY_RIGHT) THEN
        velocityX = velocityX + acceleration * dt
    ENDIF
    IF IsKeyDown(KEY_LEFT) THEN
        velocityX = velocityX - acceleration * dt
    ENDIF
    IF IsKeyDown(KEY_UP) THEN
        velocityY = velocityY - acceleration * dt
    ENDIF
    IF IsKeyDown(KEY_DOWN) THEN
        velocityY = velocityY + acceleration * dt
    ENDIF
    
    // Apply friction
    velocityX = velocityX * friction
    velocityY = velocityY * friction
    
    // Limit to max speed
    VAR currentSpeed = Sqrt(velocityX * velocityX + velocityY * velocityY)
    IF currentSpeed > maxSpeed THEN
        velocityX = (velocityX / currentSpeed) * maxSpeed
        velocityY = (velocityY / currentSpeed) * maxSpeed
    ENDIF
    
    // Update position
    playerX = playerX + velocityX * dt
    playerY = playerY + velocityY * dt
    
    // Screen boundaries
    playerX = Clamp(playerX, 20, 780)
    playerY = Clamp(playerY, 20, 580)
    
    // Stop at boundaries
    IF playerX <= 20 OR playerX >= 780 THEN velocityX = 0
    IF playerY <= 20 OR playerY >= 580 THEN velocityY = 0
    
    // Drawing
    ClearBackground(30, 30, 50, 255)
    
    // Draw player
    DrawCircle(playerX, playerY, 20, 100, 200, 255, 255)
    
    // Draw velocity vector (for debugging)
    DrawLine(playerX, playerY, playerX + velocityX * 0.2, playerY + velocityY * 0.2, 255, 255, 0, 255)
    
    // UI
    DrawText("Arrow Keys: Apply acceleration", 10, 10, 18, 255, 255, 255, 255)
    DrawText("Speed: " + STR(Int(currentSpeed)) + " / " + STR(Int(maxSpeed)), 10, 35, 16, 200, 200, 200, 255)
    DrawText("Velocity: (" + STR(Int(velocityX)) + ", " + STR(Int(velocityY)) + ")", 10, 55, 16, 200, 200, 200, 255)
WEND

CloseWindow()
```

---

## Lesson 3: Collision Detection

### Rectangle Collision

```basic
// Rectangle Collision Detection
InitWindow(800, 600, "Rectangle Collision")
SetTargetFPS(60)

// Player rectangle
VAR playerX = 100
VAR playerY = 300
VAR playerWidth = 40
VAR playerHeight = 40
VAR playerSpeed = 200

// Obstacles
VAR obstacleX[3] = [300, 500, 650]
VAR obstacleY[3] = [200, 350, 150]
VAR obstacleWidth[3] = [60, 80, 50]
VAR obstacleHeight[3] = [100, 40, 120]

// Collision state
VAR isColliding = 0

WHILE NOT WindowShouldClose()
    VAR dt = GetFrameTime()
    
    // Player movement
    IF IsKeyDown(KEY_RIGHT) THEN playerX = playerX + playerSpeed * dt
    IF IsKeyDown(KEY_LEFT) THEN playerX = playerX - playerSpeed * dt
    IF IsKeyDown(KEY_UP) THEN playerY = playerY - playerSpeed * dt
    IF IsKeyDown(KEY_DOWN) THEN playerY = playerY + playerSpeed * dt
    
    // Keep player on screen
    playerX = Clamp(playerX, 0, 800 - playerWidth)
    playerY = Clamp(playerY, 0, 600 - playerHeight)
    
    // Check collision with each obstacle
    isColliding = 0
    FOR i = 0 TO 2
        // AABB (Axis-Aligned Bounding Box) collision
        IF playerX < obstacleX[i] + obstacleWidth[i] AND
           playerX + playerWidth > obstacleX[i] AND
           playerY < obstacleY[i] + obstacleHeight[i] AND
           playerY + playerHeight > obstacleY[i] THEN
            isColliding = 1
        ENDIF
    NEXT i
    
    // Drawing
    ClearBackground(40, 40, 60, 255)
    
    // Draw obstacles
    FOR i = 0 TO 2
        DrawRectangle(obstacleX[i], obstacleY[i], obstacleWidth[i], obstacleHeight[i], 200, 100, 100, 255)
    NEXT i
    
    // Draw player (color changes when colliding)
    IF isColliding = 1 THEN
        DrawRectangle(playerX, playerY, playerWidth, playerHeight, 255, 100, 100, 255)
        DrawText("COLLISION!", 10, 10, 30, 255, 100, 100, 255)
    ELSE
        DrawRectangle(playerX, playerY, playerWidth, playerHeight, 100, 200, 255, 255)
        DrawText("No Collision", 10, 10, 30, 100, 255, 100, 255)
    ENDIF
    
    DrawText("Arrow keys to move", 10, 50, 16, 200, 200, 200, 255)
WEND

CloseWindow()
```

### Circle Collision

```basic
// Circle Collision Detection
InitWindow(800, 600, "Circle Collision")
SetTargetFPS(60)

// Player circle
VAR playerX = 400
VAR playerY = 300
VAR playerRadius = 25
VAR playerSpeed = 200

// Collectible circles
VAR collectibleX[5] = [150, 250, 400, 550, 650]
VAR collectibleY[5] = [200, 400, 150, 350, 250]
VAR collectibleRadius[5] = [20, 15, 25, 18, 22]
VAR collectibleCollected[5] = [0, 0, 0, 0, 0]

VAR score = 0

WHILE NOT WindowShouldClose()
    VAR dt = GetFrameTime()
    
    // Player movement
    IF IsKeyDown(KEY_RIGHT) THEN playerX = playerX + playerSpeed * dt
    IF IsKeyDown(KEY_LEFT) THEN playerX = playerX - playerSpeed * dt
    IF IsKeyDown(KEY_UP) THEN playerY = playerY - playerSpeed * dt
    IF IsKeyDown(KEY_DOWN) THEN playerY = playerY + playerSpeed * dt
    
    // Keep player on screen
    playerX = Clamp(playerX, playerRadius, 800 - playerRadius)
    playerY = Clamp(playerY, playerRadius, 600 - playerRadius)
    
    // Check collision with collectibles
    FOR i = 0 TO 4
        IF collectibleCollected[i] = 0 THEN
            // Distance between centers
            VAR dx = playerX - collectibleX[i]
            VAR dy = playerY - collectibleY[i]
            VAR distance = Sqrt(dx * dx + dy * dy)
            
            // Collision if distance < sum of radii
            IF distance < playerRadius + collectibleRadius[i] THEN
                collectibleCollected[i] = 1
                score = score + 10
            ENDIF
        ENDIF
    NEXT i
    
    // Drawing
    ClearBackground(30, 30, 50, 255)
    
    // Draw collectibles
    FOR i = 0 TO 4
        IF collectibleCollected[i] = 0 THEN
            DrawCircle(collectibleX[i], collectibleY[i], collectibleRadius[i], 255, 215, 0, 255)
            // Inner circle for visual effect
            DrawCircle(collectibleX[i], collectibleY[i], collectibleRadius[i] * 0.6, 255, 255, 100, 255)
        ENDIF
    NEXT i
    
    // Draw player
    DrawCircle(playerX, playerY, playerRadius, 100, 200, 255, 255)
    DrawCircle(playerX, playerY, playerRadius * 0.4, 255, 255, 255, 255)
    
    // UI
    DrawText("Score: " + STR(score), 10, 10, 25, 255, 255, 255, 255)
    DrawText("Arrow keys to move, collect coins!", 10, 45, 16, 200, 200, 200, 255)
    
    // Check win condition
    IF score = 50 THEN
        DrawText("YOU WIN!", 300, 250, 50, 255, 215, 0, 255)
    ENDIF
WEND

CloseWindow()
```

---

## Lesson 4: Building a Complete Platformer

Now let's put everything together into a complete platformer game:

```basic
// Complete 2D Platformer Game
InitWindow(1024, 768, "CyberBasic Platformer")
SetTargetFPS(60)

// Game state
VAR gameState = "menu"  // "menu", "playing", "gameover", "win"
VAR score = 0
VAR lives = 3
VAR level = 1

// Player properties
VAR playerX = 100
VAR playerY = 500
VAR playerWidth = 30
VAR playerHeight = 40
VAR velocityX = 0
VAR velocityY = 0
VAR speed = 200
VAR jumpPower = 400
VAR onGround = 0
VAR facingRight = 1

// Physics constants
VAR gravity = 800
VAR friction = 0.8

// Platforms
VAR platformX[10] = [0, 200, 400, 600, 800, 150, 350, 550, 750, 450]
VAR platformY[10] = [700, 600, 550, 500, 450, 400, 350, 300, 250, 150]
VAR platformWidth[10] = [1024, 150, 120, 100, 150, 100, 80, 100, 120, 200]
VAR platformHeight[10] = [100, 20, 20, 20, 20, 20, 20, 20, 20, 20]

// Collectibles
VAR coinX[8] = [250, 450, 650, 180, 380, 580, 750, 550]
VAR coinY[8] = [550, 500, 450, 350, 300, 250, 200, 100]
VAR coinCollected[8] = [0, 0, 0, 0, 0, 0, 0, 0]

// Enemies
VAR enemyX[3] = [250, 450, 650]
VAR enemyY[3] = [570, 520, 470]
VAR enemyWidth = 25
VAR enemyHeight = 25
VAR enemyDirection[3] = [1, -1, 1]
VAR enemySpeed[3] = [50, 75, 60]
VAR enemyPatrolDistance[3] = [60, 80, 70]
VAR enemyStartX[3] = [250, 450, 650]

// Goal
VAR goalX = 550
VAR goalY = 100
VAR goalWidth = 40
VAR goalHeight = 50

FUNCTION CheckPlatformCollision(x, y, w, h)
    // Check collision with all platforms
    FOR i = 0 TO 9
        IF x < platformX[i] + platformWidth[i] AND
           x + w > platformX[i] AND
           y < platformY[i] + platformHeight[i] AND
           y + h > platformY[i] THEN
            RETURN i  // Return platform index
        ENDIF
    NEXT i
    RETURN -1  // No collision
END FUNCTION

FUNCTION CheckCoinCollision()
    FOR i = 0 TO 7
        IF coinCollected[i] = 0 THEN
            VAR dx = (playerX + playerWidth/2) - coinX[i]
            VAR dy = (playerY + playerHeight/2) - coinY[i]
            IF Sqrt(dx*dx + dy*dy) < 25 THEN
                coinCollected[i] = 1
                score = score + 10
                RETURN 1
            ENDIF
        ENDIF
    NEXT i
    RETURN 0
END FUNCTION

FUNCTION CheckEnemyCollision()
    FOR i = 0 TO 2
        IF playerX < enemyX[i] + enemyWidth AND
           playerX + playerWidth > enemyX[i] AND
           playerY < enemyY[i] + enemyHeight AND
           playerY + playerHeight > enemyY[i] THEN
            RETURN 1
        ENDIF
    NEXT i
    RETURN 0
END FUNCTION

FUNCTION CheckGoalCollision()
    IF playerX < goalX + goalWidth AND
       playerX + playerWidth > goalX AND
       playerY < goalY + goalHeight AND
       playerY + playerHeight > goalY THEN
        RETURN 1
    ENDIF
    RETURN 0
END FUNCTION

SUB ResetPlayer()
    playerX = 100
    playerY = 500
    velocityX = 0
    velocityY = 0
    onGround = 0
END SUB

SUB ResetLevel()
    ResetPlayer()
    FOR i = 0 TO 7
        coinCollected[i] = 0
    NEXT i
    FOR i = 0 TO 2
        enemyX[i] = enemyStartX[i]
    NEXT i
END SUB

// Main game loop
WHILE NOT WindowShouldClose()
    VAR dt = GetFrameTime()
    
    SELECT CASE gameState
        CASE "menu"
            ClearBackground(20, 20, 40, 255)
            
            // Title
            DrawText("CYBERBASIC PLATFORMER", 250, 150, 40, 255, 255, 255, 255)
            DrawText("A Complete 2D Game", 320, 200, 20, 200, 200, 200, 255)
            
            // Instructions
            DrawText("Arrow Keys: Move", 350, 300, 18, 255, 255, 255, 255)
            DrawText("Space: Jump", 380, 330, 18, 255, 255, 255, 255)
            DrawText("Collect coins and reach the goal!", 280, 380, 16, 200, 200, 200, 255)
            DrawText("Avoid red enemies!", 380, 410, 16, 200, 200, 200, 255)
            
            DrawText("Press SPACE to Start", 340, 500, 25, 255, 215, 0, 255)
            
            IF IsKeyPressed(KEY_SPACE) THEN
                gameState = "playing"
                ResetLevel()
            ENDIF
            
        CASE "playing"
            // Input handling
            IF IsKeyDown(KEY_LEFT) THEN
                velocityX = velocityX - speed * dt
                facingRight = 0
            ENDIF
            IF IsKeyDown(KEY_RIGHT) THEN
                velocityX = velocityX + speed * dt
                facingRight = 1
            ENDIF
            IF IsKeyDown(KEY_SPACE) AND onGround = 1 THEN
                velocityY = -jumpPower
                onGround = 0
            ENDIF
            
            // Apply physics
            velocityX = velocityX * friction
            velocityY = velocityY + gravity * dt
            
            // Update position
            playerX = playerX + velocityX
            playerY = playerY + velocityY
            
            // Platform collision
            VAR platformIndex = CheckPlatformCollision(playerX, playerY, playerWidth, playerHeight)
            IF platformIndex >= 0 THEN
                // Landing on top of platform
                IF velocityY > 0 AND playerY < platformY[platformIndex] THEN
                    playerY = platformY[platformIndex] - playerHeight
                    velocityY = 0
                    onGround = 1
                ENDIF
            ENDIF
            
            // Screen boundaries
            playerX = Clamp(playerX, 0, 1024 - playerWidth)
            
            // Fall detection
            IF playerY > 800 THEN
                lives = lives - 1
                IF lives <= 0 THEN
                    gameState = "gameover"
                ELSE
                    ResetPlayer()
                ENDIF
            ENDIF
            
            // Update enemies
            FOR i = 0 TO 2
                enemyX[i] = enemyX[i] + enemyDirection[i] * enemySpeed[i] * dt
                IF Abs(enemyX[i] - enemyStartX[i]) > enemyPatrolDistance[i] THEN
                    enemyDirection[i] = -enemyDirection[i]
                ENDIF
            NEXT i
            
            // Check collisions
            CheckCoinCollision()
            IF CheckEnemyCollision() THEN
                lives = lives - 1
                IF lives <= 0 THEN
                    gameState = "gameover"
                ELSE
                    ResetPlayer()
                ENDIF
            ENDIF
            IF CheckGoalCollision() THEN
                gameState = "win"
            ENDIF
            
            // Drawing
            ClearBackground(135, 206, 235, 255)  // Sky blue
            
            // Draw platforms
            FOR i = 0 TO 9
                DrawRectangle(platformX[i], platformY[i], platformWidth[i], platformHeight[i], 139, 69, 19, 255)
                // Platform top highlight
                DrawRectangle(platformX[i], platformY[i], platformWidth[i], 5, 160, 82, 45, 255)
            NEXT i
            
            // Draw coins
            FOR i = 0 TO 7
                IF coinCollected[i] = 0 THEN
                    DrawCircle(coinX[i], coinY[i], 12, 255, 215, 0, 255)
                    DrawCircle(coinX[i], coinY[i], 8, 255, 255, 100, 255)
                ENDIF
            NEXT i
            
            // Draw enemies
            FOR i = 0 TO 2
                DrawRectangle(enemyX[i], enemyY[i], enemyWidth, enemyHeight, 255, 100, 100, 255)
                // Enemy eyes
                DrawCircle(enemyX[i] + 8, enemyY[i] + 8, 3, 255, 255, 255, 255)
                DrawCircle(enemyX[i] + 17, enemyY[i] + 8, 3, 255, 255, 255, 255)
            NEXT i
            
            // Draw goal
            DrawRectangle(goalX, goalY, goalWidth, goalHeight, 0, 255, 0, 255)
            DrawText("GOAL", goalX + 2, goalY + 15, 12, 255, 255, 255, 255)
            
            // Draw player
            DrawRectangle(playerX, playerY, playerWidth, playerHeight, 100, 200, 255, 255)
            // Player face
            IF facingRight = 1 THEN
                DrawCircle(playerX + 20, playerY + 12, 3, 255, 255, 255, 255)
                DrawCircle(playerX + 25, playerY + 12, 3, 255, 255, 255, 255)
            ELSE
                DrawCircle(playerX + 5, playerY + 12, 3, 255, 255, 255, 255)
                DrawCircle(playerX + 10, playerY + 12, 3, 255, 255, 255, 255)
            ENDIF
            
            // UI
            DrawText("Score: " + STR(score), 10, 10, 20, 255, 255, 255, 255)
            DrawText("Lives: " + STR(lives), 10, 35, 20, 255, 255, 255, 255)
            DrawText("Level: " + STR(level), 10, 60, 20, 255, 255, 255, 255)
            
        CASE "gameover"
            ClearBackground(50, 10, 10, 255)
            DrawText("GAME OVER", 350, 250, 50, 255, 100, 100, 255)
            DrawText("Final Score: " + STR(score), 380, 330, 25, 255, 255, 255, 255)
            DrawText("Press SPACE to return to menu", 300, 400, 18, 200, 200, 200, 255)
            
            IF IsKeyPressed(KEY_SPACE) THEN
                gameState = "menu"
                score = 0
                lives = 3
                level = 1
            ENDIF
            
        CASE "win"
            ClearBackground(10, 50, 10, 255)
            DrawText("YOU WIN!", 380, 250, 50, 100, 255, 100, 255)
            DrawText("Final Score: " + STR(score), 380, 330, 25, 255, 255, 255, 255)
            DrawText("Press SPACE to return to menu", 300, 400, 18, 200, 200, 200, 255)
            
            IF IsKeyPressed(KEY_SPACE) THEN
                gameState = "menu"
                score = 0
                lives = 3
                level = 1
            ENDIF
    END SELECT
    
WEND

CloseWindow()
```

---

## Lesson 5: Visual Polish

### Particle Effects

```basic
// Particle System for Visual Effects
InitWindow(800, 600, "Particle Effects")
SetTargetFPS(60)

// Particle structure (using arrays for simplicity)
VAR particleX[100]
VAR particleY[100]
VAR velocityX[100]
VAR velocityY[100]
VAR life[100]
VAR maxLife[100]
VAR colorR[100]
VAR colorG[100]
VAR colorB[100]
VAR active[100]

VAR particleCount = 0

FUNCTION CreateParticle(x, y, vx, vy, lifetime, r, g, b)
    IF particleCount < 100 THEN
        particleX[particleCount] = x
        particleY[particleCount] = y
        velocityX[particleCount] = vx
        velocityY[particleCount] = vy
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
    FOR i = 0 TO 99
        IF active[i] = 1 THEN
            // Update position
            particleX[i] = particleX[i] + velocityX[i] * dt
            particleY[i] = particleY[i] + velocityY[i] * dt
            
            // Update life
            life[i] = life[i] - dt
            
            // Apply gravity
            velocityY[i] = velocityY[i] + 200 * dt
            
            // Remove dead particles
            IF life[i] <= 0 THEN
                active[i] = 0
            ENDIF
        ENDIF
    NEXT i
END SUB

SUB DrawParticles()
    FOR i = 0 TO 99
        IF active[i] = 1 THEN
            VAR alpha = (life[i] / maxLife[i]) * 255
            VAR size = 3 * (life[i] / maxLife[i])
            DrawCircle(particleX[i], particleY[i], size, colorR[i], colorG[i], colorB[i], Int(alpha))
        ENDIF
    NEXT i
END SUB

WHILE NOT WindowShouldClose()
    VAR dt = GetFrameTime()
    
    // Create explosion effect on mouse click
    IF IsMouseButtonPressed(0) THEN
        VAR mx = GetMouseX()
        VAR my = GetMouseY()
        // Create burst of particles
        FOR i = 0 TO 19
            VAR angle = (i / 20.0) * 6.28318  // 2 * PI
            VAR speed = 100 + GetRandomValue(0, 200)
            VAR vx = Cos(angle) * speed
            VAR vy = Sin(angle) * speed - 100
            VAR r = GetRandomValue(200, 255)
            VAR g = GetRandomValue(100, 200)
            VAR b = GetRandomValue(0, 100)
            CreateParticle(mx, my, vx, vy, 2.0, r, g, b)
        NEXT i
    ENDIF
    
    // Update and draw particles
    UpdateParticles(dt)
    
    // Drawing
    ClearBackground(20, 20, 40, 255)
    DrawParticles()
    
    // UI
    DrawText("Click to create particle explosion!", 10, 10, 20, 255, 255, 255, 255)
    DrawText("Active particles: " + STR(particleCount), 10, 35, 16, 200, 200, 200, 255)
WEND

CloseWindow()
```

---

## Lesson 6: Sound and Music

```basic
// Sound Effects and Music
InitWindow(800, 600, "Sound Demo")
SetTargetFPS(60)

// Initialize audio
InitAudioDevice()

// Load sounds (you'll need actual sound files)
VAR jumpSound = LoadSound("sounds/jump.wav")
VAR coinSound = LoadSound("sounds/coin.wav")
VAR music = LoadMusicStream("sounds/background_music.mp3")

// Play music
PlayMusicStream(music)

// Game objects with sound
VAR playerX = 400
VAR playerY = 300
VAR playerRadius = 20

VAR coinX = 300
VAR coinY = 200
VAR coinCollected = 0

WHILE NOT WindowShouldClose()
    VAR dt = GetFrameTime()
    
    // Update music stream
    UpdateMusicStream(music)
    
    // Player movement
    IF IsKeyDown(KEY_LEFT) THEN playerX = playerX - 200 * dt
    IF IsKeyDown(KEY_RIGHT) THEN playerX = playerX + 200 * dt
    IF IsKeyDown(KEY_UP) THEN playerY = playerY - 200 * dt
    IF IsKeyDown(KEY_DOWN) THEN playerY = playerY + 200 * dt
    
    // Jump sound effect
    IF IsKeyPressed(KEY_SPACE) THEN
        PlaySound(jumpSound)
        playerY = playerY - 50
    ENDIF
    
    // Coin collection sound
    IF coinCollected = 0 THEN
        VAR dx = playerX - coinX
        VAR dy = playerY - coinY
        IF Sqrt(dx*dx + dy*dy) < playerRadius + 15 THEN
            coinCollected = 1
            PlaySound(coinSound)
        ENDIF
    ENDIF
    
    // Drawing
    ClearBackground(40, 40, 60, 255)
    
    // Draw player
    DrawCircle(playerX, playerY, playerRadius, 100, 200, 255, 255)
    
    // Draw coin
    IF coinCollected = 0 THEN
        DrawCircle(coinX, coinY, 15, 255, 215, 0, 255)
    ENDIF
    
    // UI
    DrawText("Arrow keys: move, Space: jump (sound!)", 10, 10, 18, 255, 255, 255, 255)
    DrawText("Collect the coin for sound effect!", 10, 35, 16, 200, 200, 200, 255)
    IF coinCollected = 1 THEN
        DrawText("Coin collected! +10 points", 10, 60, 16, 255, 215, 0, 255)
    ENDIF
WEND

// Cleanup
UnloadSound(jumpSound)
UnloadSound(coinSound)
UnloadMusicStream(music)
CloseAudioDevice()
CloseWindow()
```

---

## Conclusion

Congratulations! You've now learned:

- **2D coordinate system** and basic drawing
- **Smooth player movement** with delta time
- **Physics-based movement** with acceleration and friction
- **Collision detection** for rectangles and circles
- **Complete platformer game** with multiple game objects
- **Particle effects** for visual polish
- **Sound integration** for audio feedback

### Next Steps

1. **Expand the platformer**: Add more levels, power-ups, and enemy types
2. **Try different genres**: Puzzle games, top-down shooters, RPGs
3. **Learn 3D graphics**: Move on to the 3D tutorial
4. **Add multiplayer**: Network your 2D games
5. **Create original games**: Use these techniques to make something unique

### Common 2D Game Patterns

- **Platformers**: Jump mechanics, gravity, platforms
- **Top-down shooters**: 8-directional movement, projectiles
- **Puzzle games**: Grid-based logic, drag and drop
- **RPGs**: Stats systems, inventory, dialogue
- **Racing games**: Lap tracking, checkpoints, AI opponents

Keep practicing and experimenting with these techniques. The more you create, the better you'll become at game development with CyberBasic!

**Happy coding!**
