# Multiplayer Game Development Tutorial - Complete Guide

Welcome to the complete multiplayer game development tutorial! This guide will teach you how to create networked games that multiple players can enjoy together.

## What You'll Build

By the end of this tutorial, you'll have created:
- Client-server architecture
- Real-time multiplayer games
- Networked physics synchronization
- Chat systems
- Lobby and matchmaking systems
- Network security considerations

---

## Prerequisites

Before starting, make sure you've completed:
- **Module 1**: BASIC Programming Fundamentals (from LEARNING_PATH.md)
- **Module 2**: 2D Game Development (recommended for context)
- Basic understanding of networking concepts

---

## Lesson 1: Understanding Network Architecture

### Client-Server Model

In multiplayer games, we typically use a client-server architecture:
- **Server**: Authoritative game state, processes game logic
- **Client**: Sends input, receives game state, renders graphics

```basic
// Basic Server Example
// Save as: multiplayer_server.bas

// Server setup
VAR server = Host(9999)  // Host on port 9999
IF IsNull(server) THEN
    PRINT "Failed to host server on port 9999"
    QUIT
ENDIF

PRINT "Server started successfully!"
PRINT "Waiting for players to connect..."

// Game state
VAR players[4]  // Support up to 4 players
VAR playerCount = 0
VAR gameStarted = 0

// Player data structure (using arrays)
VAR playerX[4] = [0, 0, 0, 0]
VAR playerY[4] = [0, 0, 0, 0]
VAR playerNames[4] = ["", "", "", ""]
VAR playerConnected[4] = [0, 0, 0, 0]

FUNCTION AddPlayer(client)
    IF playerCount < 4 THEN
        players[playerCount] = client
        playerConnected[playerCount] = 1
        playerNames[playerCount] = "Player" + STR(playerCount + 1)
        playerX[playerCount] = GetRandomValue(100, 700)
        playerY[playerCount] = GetRandomValue(100, 500)
        playerCount = playerCount + 1
        
        PRINT "Player " + STR(playerCount) + " connected!"
        RETURN playerCount - 1
    ENDIF
    RETURN -1  // Server full
END FUNCTION

FUNCTION RemovePlayer(clientIndex)
    IF clientIndex >= 0 AND clientIndex < 4 THEN
        playerConnected[clientIndex] = 0
        playerNames[clientIndex] = ""
        playerCount = playerCount - 1
        PRINT "Player " + STR(clientIndex + 1) + " disconnected"
    ENDIF
END FUNCTION

FUNCTION Broadcast(message)
    // Send message to all connected players
    FOR i = 0 TO 3
        IF playerConnected[i] = 1 THEN
            Send(players[i], message)
        ENDIF
    NEXT i
END FUNCTION

FUNCTION SendGameState()
    // Send current game state to all players
    FOR i = 0 TO 3
        IF playerConnected[i] = 1 THEN
            VAR stateMessage = "STATE:"
            FOR j = 0 TO 3
                IF playerConnected[j] = 1 THEN
                    stateMessage = stateMessage + playerNames[j] + "," + STR(playerX[j]) + "," + STR(playerY[j]) + ";"
                ENDIF
            NEXT j
            Send(players[i], stateMessage)
        ENDIF
    NEXT i
END FUNCTION

// Main server loop
VAR lastStateUpdate = 0
WHILE NOT IsWindowShouldClose()  // Use a condition to keep server running
    VAR currentTime = GetTime()
    
    // Accept new connections
    VAR newClient = Accept(server, 100)  // Wait 100ms for new connection
    IF NOT IsNull(newClient) THEN
        VAR playerIndex = AddPlayer(newClient)
        IF playerIndex >= 0 THEN
            // Send welcome message to new player
            Send(newClient, "WELCOME:" + STR(playerIndex))
            Send(newClient, "NAME:" + playerNames[playerIndex])
            
            // Notify other players
            Broadcast("PLAYER_JOINED:" + playerNames[playerIndex])
            
            // Send current game state
            SendGameState()
        ELSE
            Send(newClient, "SERVER_FULL")
            Disconnect(newClient)
        ENDIF
    ENDIF
    
    // Receive messages from clients
    FOR i = 0 TO 3
        IF playerConnected[i] = 1 THEN
            VAR message = Receive(players[i], 10)  // Non-blocking receive
            IF NOT IsNull(message) THEN
                // Parse client messages
                IF LEFT(message, 9) = "POSITION:" THEN
                    VAR posData = MID(message, 10)
                    // Parse position data (simplified)
                    VAR commaPos = INSTR(posData, ",")
                    IF commaPos > 0 THEN
                        VAR xStr = LEFT(posData, commaPos - 1)
                        VAR yStr = MID(posData, commaPos + 1)
                        playerX[i] = VAL(xStr)
                        playerY[i] = VAL(yStr)
                    ENDIF
                ELSEIF message = "DISCONNECT" THEN
                    Broadcast("PLAYER_LEFT:" + playerNames[i])
                    RemovePlayer(i)
                ELSEIF LEFT(message, 6) = "CHAT:" THEN
                    VAR chatMessage = MID(message, 7)
                    Broadcast("CHAT:" + playerNames[i] + ": " + chatMessage)
                ENDIF
            ENDIF
        ENDIF
    NEXT i
    
    // Send game state updates periodically
    IF currentTime - lastStateUpdate > 0.05 THEN  // 20 FPS updates
        SendGameState()
        lastStateUpdate = currentTime
    ENDIF
    
    // Small delay to prevent 100% CPU usage
    WaitTime(0.01)
WEND

// Cleanup
PRINT "Server shutting down..."
CloseServer(server)
```

```basic
// Basic Client Example
// Save as: multiplayer_client.bas

// Client setup
VAR serverIP = "127.0.0.1"  // localhost
VAR serverPort = 9999

PRINT "Connecting to server " + serverIP + ":" + STR(serverPort) + "..."

VAR connection = Connect(serverIP, serverPort)
IF IsNull(connection) THEN
    PRINT "Failed to connect to server"
    QUIT
ENDIF

PRINT "Connected successfully!"

// Game state
VAR myPlayerIndex = -1
VAR myPlayerName = ""
VAR playerX[4] = [0, 0, 0, 0]
VAR playerY[4] = [0, 0, 0, 0]
VAR playerNames[4] = ["", "", "", ""]
VAR playerConnected[4] = [0, 0, 0, 0]

// Graphics
InitWindow(800, 600, "Multiplayer Client")
SetTargetFPS(60)

// Input
VAR myX = 400
VAR myY = 300
VAR chatInput = ""
VAR chatMode = 0

FUNCTION ParseStateMessage(message)
    // Parse STATE:Name1,x1,y1;Name2,x2,y2;...
    VAR data = MID(message, 7)  // Remove "STATE:"
    
    // Clear current state
    FOR i = 0 TO 3
        playerConnected[i] = 0
        playerNames[i] = ""
    NEXT i
    
    // Parse player data
    VAR playerIndex = 0
    VAR remaining = data
    
    WHILE LEN(remaining) > 0 AND playerIndex < 4
        VAR semicolonPos = INSTR(remaining, ";")
        IF semicolonPos = 0 THEN semicolonPos = LEN(remaining) + 1
        
        VAR playerData = LEFT(remaining, semicolonPos - 1)
        remaining = MID(remaining, semicolonPos + 1)
        
        IF LEN(playerData) > 0 THEN
            VAR comma1 = INSTR(playerData, ",")
            VAR comma2 = INSTR(comma1 + 1, playerData, ",")
            
            IF comma1 > 0 AND comma2 > 0 THEN
                playerNames[playerIndex] = LEFT(playerData, comma1 - 1)
                playerX[playerIndex] = VAL(MID(playerData, comma1 + 1, comma2 - comma1 - 1))
                playerY[playerIndex] = VAL(MID(playerData, comma2 + 1))
                playerConnected[playerIndex] = 1
                playerIndex = playerIndex + 1
            ENDIF
        ENDIF
    WEND
END FUNCTION

WHILE NOT WindowShouldClose()
    // Handle server messages
    VAR message = Receive(connection, 10)
    IF NOT IsNull(message) THEN
        IF LEFT(message, 8) = "WELCOME:" THEN
            myPlayerIndex = VAL(MID(message, 9))
            PRINT "Assigned player index: " + STR(myPlayerIndex)
        ELSEIF LEFT(message, 5) = "NAME:" THEN
            myPlayerName = MID(message, 6)
            PRINT "My name: " + myPlayerName
        ELSEIF LEFT(message, 6) = "STATE:" THEN
            ParseStateMessage(message)
        ELSEIF LEFT(message, 13) = "PLAYER_JOINED:" THEN
            VAR playerName = MID(message, 14)
            PRINT playerName + " joined the game"
        ELSEIF LEFT(message, 11) = "PLAYER_LEFT:" THEN
            VAR playerName = MID(message, 12)
            PRINT playerName + " left the game"
        ELSEIF LEFT(message, 5) = "CHAT:" THEN
            VAR chatMessage = MID(message, 6)
            PRINT chatMessage
        ELSEIF message = "SERVER_FULL" THEN
            PRINT "Server is full"
            EXIT WHILE
        ENDIF
    ENDIF
    
    // Handle input
    IF NOT chatMode THEN
        // Movement input
        VAR oldX = myX
        VAR oldY = myY
        
        IF IsKeyDown(KEY_LEFT) THEN myX = myX - 5
        IF IsKeyDown(KEY_RIGHT) THEN myX = myX + 5
        IF IsKeyDown(KEY_UP) THEN myY = myY - 5
        IF IsKeyDown(KEY_DOWN) THEN myY = myY + 5
        
        // Keep player on screen
        myX = Clamp(myX, 20, 780)
        myY = Clamp(myY, 20, 580)
        
        // Send position update if moved
        IF oldX <> myX OR oldY <> myY THEN
            Send(connection, "POSITION:" + STR(myX) + "," + STR(myY))
        ENDIF
        
        // Enter chat mode
        IF IsKeyPressed(KEY_T) THEN
            chatMode = 1
            chatInput = ""
        ENDIF
    ELSE
        // Chat input mode
        IF IsKeyPressed(KEY_ESCAPE) THEN
            chatMode = 0
            chatInput = ""
        ELSEIF IsKeyPressed(KEY_ENTER) THEN
            IF LEN(chatInput) > 0 THEN
                Send(connection, "CHAT:" + chatInput)
                chatInput = ""
            ENDIF
            chatMode = 0
        ELSEIF IsKeyPressed(KEY_BACKSPACE) THEN
            IF LEN(chatInput) > 0 THEN
                chatInput = LEFT(chatInput, LEN(chatInput) - 1)
            ENDIF
        ENDIF
    ENDIF
    
    // Drawing
    ClearBackground(40, 40, 60, 255)
    
    // Draw all players
    FOR i = 0 TO 3
        IF playerConnected[i] = 1 THEN
            VAR drawX = playerX[i]
            VAR drawY = playerY[i]
            
            // Different colors for different players
            SELECT CASE i
                CASE 0: DrawCircle(drawX, drawY, 20, 255, 100, 100, 255)
                CASE 1: DrawCircle(drawX, drawY, 20, 100, 255, 100, 255)
                CASE 2: DrawCircle(drawX, drawY, 20, 100, 100, 255, 255)
                CASE 3: DrawCircle(drawX, drawY, 20, 255, 255, 100, 255)
            END SELECT
            
            // Draw player name
            DrawText(playerNames[i], drawX - 30, drawY - 35, 12, 255, 255, 255, 255)
            
            // Highlight self
            IF i = myPlayerIndex THEN
                DrawCircleLines(drawX, drawY, 25, 255, 255, 255, 255)
            ENDIF
        ENDIF
    NEXT i
    
    // UI
    DrawText("Multiplayer Client", 10, 10, 20, 255, 255, 255, 255)
    DrawText("Arrow keys to move, T to chat", 10, 35, 16, 200, 200, 200, 255)
    
    IF chatMode THEN
        DrawText("Chat: " + chatInput + "_", 10, 60, 16, 255, 255, 100, 255)
        DrawText("ESC to cancel, ENTER to send", 10, 80, 14, 200, 200, 200, 255)
    ENDIF
    
    DrawText("Connected to: " + serverIP + ":" + STR(serverPort), 10, 570, 14, 200, 200, 200, 255)
WEND

// Cleanup
Send(connection, "DISCONNECT")
Disconnect(connection)
CloseWindow()
PRINT "Disconnected from server"
```

---

## Lesson 2: Real-time Action Game

### Networked 2D Shooter

```basic
// Multiplayer Shooter Server
// Save as: shooter_server.bas

VAR server = Host(9999)
IF IsNull(server) THEN
    PRINT "Failed to start shooter server"
    QUIT
ENDIF

PRINT "Shooter server started on port 9999"

// Game constants
CONST MAX_PLAYERS = 8
CONST ARENA_WIDTH = 800
CONST ARENA_HEIGHT = 600
CONST PLAYER_SIZE = 20
CONST BULLET_SPEED = 10
CONST FIRE_RATE = 0.2  // Seconds between shots

// Player data
VAR players[MAX_PLAYERS]
VAR playerConnected[MAX_PLAYERS] = [0, 0, 0, 0, 0, 0, 0, 0]
VAR playerX[MAX_PLAYERS] = [0, 0, 0, 0, 0, 0, 0, 0]
VAR playerY[MAX_PLAYERS] = [0, 0, 0, 0, 0, 0, 0, 0]
VAR playerAngle[MAX_PLAYERS] = [0, 0, 0, 0, 0, 0, 0, 0]
VAR playerHealth[MAX_PLAYERS] = [100, 100, 100, 100, 100, 100, 100, 100]
VAR playerScore[MAX_PLAYERS] = [0, 0, 0, 0, 0, 0, 0, 0]
VAR playerLastShot[MAX_PLAYERS] = [0, 0, 0, 0, 0, 0, 0, 0]
VAR playerNames[MAX_PLAYERS]
VAR playerCount = 0

// Bullet data
VAR bulletX[100]
VAR bulletY[100]
VAR bulletVX[100]
VAR bulletVY[100]
VAR bulletOwner[100]
VAR bulletActive[100] = [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0]
VAR bulletCount = 0

FUNCTION SpawnBullet(x, y, angle, ownerIndex)
    IF bulletCount < 100 THEN
        bulletX[bulletCount] = x
        bulletY[bulletCount] = y
        bulletVX[bulletCount] = Cos(angle) * BULLET_SPEED
        bulletVY[bulletCount] = Sin(angle) * BULLET_SPEED
        bulletOwner[bulletCount] = ownerIndex
        bulletActive[bulletCount] = 1
        bulletCount = bulletCount + 1
    ENDIF
END FUNCTION

FUNCTION UpdateBullets()
    FOR i = 0 TO 99
        IF bulletActive[i] = 1 THEN
            // Move bullet
            bulletX[i] = bulletX[i] + bulletVX[i]
            bulletY[i] = bulletY[i] + bulletVY[i]
            
            // Check if bullet is out of bounds
            IF bulletX[i] < 0 OR bulletX[i] > ARENA_WIDTH OR bulletY[i] < 0 OR bulletY[i] > ARENA_HEIGHT THEN
                bulletActive[i] = 0
            ENDIF
            
            // Check collision with players
            FOR j = 0 TO MAX_PLAYERS - 1
                IF playerConnected[j] = 1 AND j <> bulletOwner[i] THEN
                    VAR dx = bulletX[i] - playerX[j]
                    VAR dy = bulletY[i] - playerY[j]
                    IF Sqrt(dx*dx + dy*dy) < PLAYER_SIZE THEN
                        // Hit!
                        playerHealth[j] = playerHealth[j] - 10
                        playerScore[bulletOwner[i]] = playerScore[bulletOwner[i]] + 10
                        bulletActive[i] = 0
                        
                        // Send hit message
                        Send(players[j], "HIT:" + STR(bulletOwner[i]))
                        Send(players[bulletOwner[i]], "HIT_PLAYER:" + STR(j))
                        
                        // Check if player is dead
                        IF playerHealth[j] <= 0 THEN
                            playerHealth[j] = 100
                            playerX[j] = GetRandomValue(50, ARENA_WIDTH - 50)
                            playerY[j] = GetRandomValue(50, ARENA_HEIGHT - 50)
                            playerScore[j] = playerScore[j] - 5
                            
                            Broadcast("PLAYER_DIED:" + STR(j))
                        ENDIF
                    ENDIF
                ENDIF
            NEXT j
        ENDIF
    NEXT i
END FUNCTION

FUNCTION SendGameState()
    FOR i = 0 TO MAX_PLAYERS - 1
        IF playerConnected[i] = 1 THEN
            VAR stateMsg = "GAME_STATE:"
            
            // Add player data
            FOR j = 0 TO MAX_PLAYERS - 1
                IF playerConnected[j] = 1 THEN
                    stateMsg = stateMsg + STR(j) + "," + STR(playerX[j]) + "," + STR(playerY[j]) + "," + STR(playerAngle[j]) + "," + STR(playerHealth[j]) + "," + STR(playerScore[j]) + ";"
                ENDIF
            NEXT j
            
            // Add bullet data
            stateMsg = stateMsg + "BULLETS:"
            FOR j = 0 TO 99
                IF bulletActive[j] = 1 THEN
                    stateMsg = stateMsg + STR(bulletX[j]) + "," + STR(bulletY[j]) + ";"
                ENDIF
            NEXT j
            
            Send(players[i], stateMsg)
        ENDIF
    NEXT i
END FUNCTION

// Main server loop
VAR lastUpdate = 0
WHILE 1  // Infinite loop for server
    VAR currentTime = GetTime()
    
    // Accept new connections
    VAR newClient = Accept(server, 50)
    IF NOT IsNull(newClient) THEN
        IF playerCount < MAX_PLAYERS THEN
            // Find empty slot
            FOR i = 0 TO MAX_PLAYERS - 1
                IF playerConnected[i] = 0 THEN
                    players[i] = newClient
                    playerConnected[i] = 1
                    playerNames[i] = "Player" + STR(i + 1)
                    playerX[i] = GetRandomValue(50, ARENA_WIDTH - 50)
                    playerY[i] = GetRandomValue(50, ARENA_HEIGHT - 50)
                    playerHealth[i] = 100
                    playerScore[i] = 0
                    playerCount = playerCount + 1
                    
                    Send(newClient, "CONNECTED:" + STR(i))
                    Send(newClient, "ARENA:" + STR(ARENA_WIDTH) + "," + STR(ARENA_HEIGHT))
                    Broadcast("PLAYER_JOINED:" + STR(i) + "," + playerNames[i])
                    
                    PRINT playerNames[i] + " connected (ID: " + STR(i) + ")"
                    EXIT FOR
                ENDIF
            NEXT i
        ELSE
            Send(newClient, "SERVER_FULL")
            Disconnect(newClient)
        ENDIF
    ENDIF
    
    // Process client messages
    FOR i = 0 TO MAX_PLAYERS - 1
        IF playerConnected[i] = 1 THEN
            VAR msg = Receive(players[i], 10)
            IF NOT IsNull(msg) THEN
                IF LEFT(msg, 9) = "POSITION:" THEN
                    // Parse position and angle
                    VAR data = MID(msg, 10)
                    VAR comma1 = INSTR(data, ",")
                    VAR comma2 = INSTR(comma1 + 1, data, ",")
                    
                    IF comma1 > 0 AND comma2 > 0 THEN
                        playerX[i] = VAL(LEFT(data, comma1 - 1))
                        playerY[i] = VAL(MID(data, comma1 + 1, comma2 - comma1 - 1))
                        playerAngle[i] = VAL(MID(data, comma2 + 1))
                    ENDIF
                ELSEIF msg = "FIRE" THEN
                    // Check fire rate
                    IF currentTime - playerLastShot[i] >= FIRE_RATE THEN
                        SpawnBullet(playerX[i], playerY[i], playerAngle[i], i)
                        playerLastShot[i] = currentTime
                        Broadcast("BULLET_FIRED:" + STR(i))
                    ENDIF
                ELSEIF msg = "DISCONNECT" THEN
                    Broadcast("PLAYER_LEFT:" + STR(i))
                    playerConnected[i] = 0
                    playerCount = playerCount - 1
                    PRINT playerNames[i] + " disconnected"
                ENDIF
            ENDIF
        ENDIF
    NEXT i
    
    // Update game state
    UpdateBullets()
    
    // Send state updates
    IF currentTime - lastUpdate >= 0.033 THEN  // 30 FPS
        SendGameState()
        lastUpdate = currentTime
    ENDIF
    
    WaitTime(0.01)
WEND
```

```basic
// Multiplayer Shooter Client
// Save as: shooter_client.bas

VAR connection = Connect("127.0.0.1", 9999)
IF IsNull(connection) THEN
    PRINT "Failed to connect to shooter server"
    QUIT
ENDIF

// Game state
VAR myPlayerID = -1
VAR arenaWidth = 800
VAR arenaHeight = 600

// Player data
VAR playerConnected[8] = [0, 0, 0, 0, 0, 0, 0, 0]
VAR playerX[8] = [0, 0, 0, 0, 0, 0, 0, 0]
VAR playerY[8] = [0, 0, 0, 0, 0, 0, 0, 0]
VAR playerAngle[8] = [0, 0, 0, 0, 0, 0, 0, 0]
VAR playerHealth[8] = [100, 100, 100, 100, 100, 100, 100, 100]
VAR playerScore[8] = [0, 0, 0, 0, 0, 0, 0, 0]

// Bullet data
VAR bulletX[50]
VAR bulletY[50]
VAR bulletActive[50] = [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0]
VAR bulletCount = 0

// Input
VAR mouseX = 400
VAR mouseY = 300

InitWindow(arenaWidth, arenaHeight, "Multiplayer Shooter")
SetTargetFPS(60)
DisableCursor()

FUNCTION ParseGameState(message)
    // Parse GAME_STATE:player_data;BULLETS:bullet_data
    VAR playerSection = LEFT(message, INSTR(message, "BULLETS:") - 1)
    VAR bulletSection = MID(message, INSTR(message, "BULLETS:") + 8)
    
    // Parse player data
    VAR playerData = MID(playerSection, 12)  // Remove "GAME_STATE:"
    bulletCount = 0
    
    WHILE LEN(playerData) > 0
        VAR semicolon = INSTR(playerData, ";")
        IF semicolon = 0 THEN EXIT WHILE
        
        VAR playerInfo = LEFT(playerData, semicolon - 1)
        playerData = MID(playerData, semicolon + 1)
        
        IF LEN(playerInfo) > 0 THEN
            VAR commas[4]
            VAR start = 1
            VAR count = 0
            
            // Find commas
            FOR i = 1 TO LEN(playerInfo)
                IF MID(playerInfo, i, 1) = "," THEN
                    commas[count] = i
                    count = count + 1
                    IF count >= 4 THEN EXIT FOR
                ENDIF
            NEXT i
            
            IF count >= 4 THEN
                VAR id = VAL(LEFT(playerInfo, commas[0] - 1))
                playerX[id] = VAL(MID(playerInfo, commas[0] + 1, commas[1] - commas[0] - 1))
                playerY[id] = VAL(MID(playerInfo, commas[1] + 1, commas[2] - commas[1] - 1))
                playerAngle[id] = VAL(MID(playerInfo, commas[2] + 1, commas[3] - commas[2] - 1))
                playerHealth[id] = VAL(MID(playerInfo, commas[3] + 1, commas[4] - commas[3] - 1))
                playerScore[id] = VAL(MID(playerInfo, commas[4] + 1))
                playerConnected[id] = 1
            ENDIF
        ENDIF
    WEND
    
    // Parse bullet data
    WHILE LEN(bulletSection) > 0 AND bulletCount < 50
        VAR semicolon = INSTR(bulletSection, ";")
        IF semicolon = 0 THEN EXIT WHILE
        
        VAR bulletInfo = LEFT(bulletSection, semicolon - 1)
        bulletSection = MID(bulletSection, semicolon + 1)
        
        IF LEN(bulletInfo) > 0 THEN
            VAR comma = INSTR(bulletInfo, ",")
            IF comma > 0 THEN
                bulletX[bulletCount] = VAL(LEFT(bulletInfo, comma - 1))
                bulletY[bulletCount] = VAL(MID(bulletInfo, comma + 1))
                bulletActive[bulletCount] = 1
                bulletCount = bulletCount + 1
            ENDIF
        ENDIF
    WEND
END FUNCTION

WHILE NOT WindowShouldClose()
    // Get input
    mouseX = GetMouseX()
    mouseY = GetMouseY()
    
    // Calculate angle to mouse
    VAR angle = ATan2(mouseY - playerY[myPlayerID], mouseX - playerX[myPlayerID])
    
    // Movement
    VAR moveX = 0
    VAR moveY = 0
    VAR moveSpeed = 5
    
    IF IsKeyDown(KEY_A) THEN moveX = -moveSpeed
    IF IsKeyDown(KEY_D) THEN moveX = moveSpeed
    IF IsKeyDown(KEY_W) THEN moveY = -moveSpeed
    IF IsKeyDown(KEY_S) THEN moveY = moveSpeed
    
    // Update position
    playerX[myPlayerID] = playerX[myPlayerID] + moveX
    playerY[myPlayerID] = playerY[myPlayerID] + moveY
    
    // Keep in bounds
    playerX[myPlayerID] = Clamp(playerX[myPlayerID], 20, arenaWidth - 20)
    playerY[myPlayerID] = Clamp(playerY[myPlayerID], 20, arenaHeight - 20)
    
    // Send position update
    Send(connection, "POSITION:" + STR(playerX[myPlayerID]) + "," + STR(playerY[myPlayerID]) + "," + STR(angle))
    
    // Shooting
    IF IsMouseButtonDown(0) THEN
        Send(connection, "FIRE")
    ENDIF
    
    // Receive messages
    VAR msg = Receive(connection, 10)
    IF NOT IsNull(msg) THEN
        IF LEFT(msg, 10) = "CONNECTED:" THEN
            myPlayerID = VAL(MID(msg, 11))
            PRINT "Connected as player " + STR(myPlayerID)
        ELSEIF LEFT(msg, 6) = "ARENA:" THEN
            VAR sizeData = MID(msg, 7)
            VAR comma = INSTR(sizeData, ",")
            arenaWidth = VAL(LEFT(sizeData, comma - 1))
            arenaHeight = VAL(MID(sizeData, comma + 1))
        ELSEIF LEFT(msg, 11) = "GAME_STATE:" THEN
            ParseGameState(msg)
        ELSEIF LEFT(msg, 13) = "PLAYER_DIED:" THEN
            VAR id = VAL(MID(msg, 14))
            PRINT "Player " + STR(id) + " died"
        ENDIF
    ENDIF
    
    // Drawing
    ClearBackground(20, 20, 30, 255)
    
    // Draw arena boundary
    DrawRectangleLines(0, 0, arenaWidth, arenaHeight, 100, 100, 100, 255)
    
    // Draw bullets
    FOR i = 0 TO 49
        IF bulletActive[i] = 1 THEN
            DrawCircle(bulletX[i], bulletY[i], 3, 255, 255, 0, 255)
        ENDIF
    NEXT i
    
    // Draw players
    FOR i = 0 TO 7
        IF playerConnected[i] = 1 THEN
            // Draw player
            DrawCircle(playerX[i], playerY[i], 20, 100, 200, 255, 255)
            
            // Draw health bar
            DrawRectangle(playerX[i] - 20, playerY[i] - 30, 40, 4, 255, 0, 0, 255)
            DrawRectangle(playerX[i] - 20, playerY[i] - 30, Int(40 * playerHealth[i] / 100), 4, 0, 255, 0, 255)
            
            // Draw direction indicator
            VAR dirX = playerX[i] + Cos(playerAngle[i]) * 25
            VAR dirY = playerY[i] + Sin(playerAngle[i]) * 25
            DrawLine(playerX[i], playerY[i], dirX, dirY, 255, 255, 255, 255)
            
            // Draw score
            DrawText("Score: " + STR(playerScore[i]), playerX[i] - 25, playerY[i] + 25, 12, 255, 255, 255, 255)
            
            // Highlight self
            IF i = myPlayerID THEN
                DrawCircleLines(playerX[i], playerY[i], 25, 255, 255, 0, 255)
            ENDIF
        ENDIF
    NEXT i
    
    // Draw crosshair
    DrawLine(mouseX - 10, mouseY, mouseX + 10, mouseY, 255, 255, 255, 255)
    DrawLine(mouseX, mouseY - 10, mouseX, mouseY + 10, 255, 255, 255, 255)
    
    // UI
    DrawText("Health: " + STR(playerHealth[myPlayerID]), 10, 10, 16, 255, 255, 255, 255)
    DrawText("Score: " + STR(playerScore[myPlayerID]), 10, 30, 16, 255, 255, 255, 255)
    DrawText("WASD: Move, Mouse: Aim, Click: Shoot", 10, arenaHeight - 20, 14, 200, 200, 200, 255)
WEND

// Cleanup
Send(connection, "DISCONNECT")
Disconnect(connection)
CloseWindow()
```

---

## Lesson 3: Advanced Networking Concepts

### Lag Compensation and Prediction

```basic
// Client-Side Prediction Example
// Add this to your client code

// Prediction data
VAR predictedX = 400
VAR predictedY = 300
VAR lastServerUpdate = 0
VAR inputSequence = 0
VAR pendingInputs[10]
VAR pendingCount = 0

TYPE InputState
    sequence AS Integer
    moveX AS Float
    moveY AS Float
    timestamp AS Double
END TYPE

FUNCTION AddPendingInput(moveX, moveY)
    IF pendingCount < 10 THEN
        pendingInputs[pendingCount].sequence = inputSequence
        pendingInputs[pendingCount].moveX = moveX
        pendingInputs[pendingCount].moveY = moveY
        pendingInputs[pendingCount].timestamp = GetTime()
        pendingCount = pendingCount + 1
        inputSequence = inputSequence + 1
    ENDIF
END FUNCTION

FUNCTION PredictPosition()
    // Start from last known server position
    predictedX = playerX[myPlayerID]
    predictedY = playerY[myPlayerID]
    
    // Apply all pending inputs
    FOR i = 0 TO pendingCount - 1
        predictedX = predictedX + pendingInputs[i].moveX
        predictedY = predictedY + pendingInputs[i].moveY
    NEXT i
END FUNCTION

FUNCTION Reconcile(serverX, serverY, serverSequence)
    // Remove old inputs that server has processed
    VAR i = 0
    WHILE i < pendingCount
        IF pendingInputs[i].sequence <= serverSequence THEN
            // Remove this input
            FOR j = i TO pendingCount - 2
                pendingInputs[j] = pendingInputs[j + 1]
            NEXT j
            pendingCount = pendingCount - 1
        ELSE
            i = i + 1
        ENDIF
    WEND
    
    // Repredict from new server position
    playerX[myPlayerID] = serverX
    playerY[myPlayerID] = serverY
    PredictPosition()
END FUNCTION
```

---

## Lesson 4: Security and Anti-Cheat

### Basic Server Validation

```basic
// Server-side validation functions
FUNCTION ValidatePosition(x, y, oldX, oldY, deltaTime)
    VAR maxSpeed = 300  // pixels per second
    VAR maxDistance = maxSpeed * deltaTime
    
    VAR distance = Sqrt((x - oldX) * (x - oldX) + (y - oldY) * (y - oldY))
    
    IF distance > maxDistance THEN
        RETURN 0  // Invalid - moving too fast
    ENDIF
    
    // Check boundaries
    IF x < 0 OR x > ARENA_WIDTH OR y < 0 OR y > ARENA_HEIGHT THEN
        RETURN 0  // Invalid - out of bounds
    ENDIF
    
    RETURN 1  // Valid
END FUNCTION

FUNCTION ValidateFire(lastFireTime, currentTime)
    VAR minFireRate = 0.1  // Minimum time between shots
    
    IF currentTime - lastFireTime < minFireRate THEN
        RETURN 0  // Firing too fast
    ENDIF
    
    RETURN 1  // Valid
END FUNCTION

// Rate limiting
VAR messageCount[8] = [0, 0, 0, 0, 0, 0, 0, 0]
VAR lastRateCheck[8] = [0, 0, 0, 0, 0, 0, 0, 0]
CONST MAX_MESSAGES_PER_SECOND = 30

FUNCTION CheckRateLimit(playerIndex)
    VAR currentTime = GetTime()
    
    IF currentTime - lastRateCheck[playerIndex] >= 1.0 THEN
        messageCount[playerIndex] = 0
        lastRateCheck[playerIndex] = currentTime
    ENDIF
    
    messageCount[playerIndex] = messageCount[playerIndex] + 1
    
    IF messageCount[playerIndex] > MAX_MESSAGES_PER_SECOND THEN
        RETURN 0  // Rate limited
    ENDIF
    
    RETURN 1  // OK
END FUNCTION
```

---

## Conclusion

Congratulations! You've now learned:

- **Network architecture** - client-server model
- **Real-time synchronization** - position updates, game state
- **Action game networking** - shooting, collision detection
- **Advanced techniques** - prediction, lag compensation
- **Security considerations** - validation, rate limiting

### Next Steps

1. **Implement matchmaking**: Lobby systems and server browsing
2. **Add voice chat**: Real-time audio communication
3. **Create dedicated servers**: Separate game and server logic
4. **Implement NAT traversal**: Connect through routers
5. **Add persistence**: Save player stats and progress

### Common Multiplayer Patterns

- **First-Person Shooters**: Fast-paced action, low latency
- **MMORPGs**: Massive worlds, complex state synchronization
- **Strategy Games**: Turn-based or real-time, resource management
- **Racing Games**: Synchronized physics, position prediction
- **Co-op Games**: Shared objectives, team coordination

Multiplayer development adds complexity but creates engaging social experiences. Start with simple games and gradually add more advanced features as you learn!

**Happy coding!**
