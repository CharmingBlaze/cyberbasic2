REM Simple Platformer Game in CyberBasic
REM Demonstrates 2D game development with physics

REM Game constants
CONST GRAVITY = -9.81
CONST JUMP_FORCE = 500
CONST MOVE_SPEED = 200
CONST MAX_FALL_SPEED = 400

REM Initialize game
INITGRAPHICS 1024, 768, "CyberBasic Platformer"
CREATEPHYSICSWORLD2D 0, GRAVITY

REM Create player
CREATEPHYSICSBODY2D "player", BODY_2D_DYNAMIC, SHAPE_2D_BOX, 32, 48, 100, 300, 1.0
SETRESTITUTION "player", 0.0
SETFRICTION "player", 0.8

REM Create level platforms
CREATEPHYSICSBODY2D "ground", BODY_2D_STATIC, SHAPE_2D_BOX, 1024, 40, 512, 20, 1.0
CREATEPHYSICSBODY2D "platform1", BODY_2D_STATIC, SHAPE_2D_BOX, 200, 20, 300, 200, 1.0
CREATEPHYSICSBODY2D "platform2", BODY_2D_STATIC, SHAPE_2D_BOX, 200, 20, 600, 300, 1.0
CREATEPHYSICSBODY2D "platform3", BODY_2D_STATIC, SHAPE_2D_BOX, 150, 20, 450, 400, 1.0
CREATEPHYSICSBODY2D "platform4", BODY_2D_STATIC, SHAPE_2D_BOX, 180, 20, 750, 450, 1.0

REM Create collectibles
CREATEPHYSICSBODY2D "coin1", BODY_2D_DYNAMIC, SHAPE_2D_CIRCLE, 16, 16, 300, 250, 0.5
CREATEPHYSICSBODY2D "coin2", BODY_2D_DYNAMIC, SHAPE_2D_CIRCLE, 16, 16, 600, 350, 0.5
CREATEPHYSICSBODY2D "coin3", BODY_2D_DYNAMIC, SHAPE_2D_CIRCLE, 16, 16, 450, 450, 0.5
CREATEPHYSICSBODY2D "coin4", BODY_2D_DYNAMIC, SHAPE_2D_CIRCLE, 16, 16, 750, 500, 0.5

REM Set coin properties
SETRESTITUTION "coin1", 0.5
SETRESTITUTION "coin2", 0.5
SETRESTITUTION "coin3", 0.5
SETRESTITUTION "coin4", 0.5

REM Game variables
DIM score AS INTEGER
score = 0
DIM onGround AS BOOLEAN
onGround = FALSE
DIM canJump AS BOOLEAN
canJump = TRUE

REM Load assets
LOADIMAGE "player.png"
LOADIMAGE "coin.png"
LOADIMAGE "platform.png"

REM Create sprites
CREATESPRITE "playerSprite", "player.png", 100, 300
CREATESPRITE "coinSprite1", "coin.png", 300, 250
CREATESPRITE "coinSprite2", "coin.png", 600, 350
CREATESPRITE "coinSprite3", "coin.png", 450, 450
CREATESPRITE "coinSprite4", "coin.png", 750, 500

REM Main game loop
WHILE NOT WINDOWSHOULDCLOSE()
    REM Update physics
    STEPPHYSICS2D 1.0/60.0, 8, 3
    
    REM Get player position and velocity
    playerPos = GETPOSITION2D "player"
    playerVel = GETVELOCITY2D "player"
    
    REM Check if player is on ground
    onGround = CHECKGROUNDCONTACT "player"
    
    REM Handle input
    IF ISKEYDOWN(KEY_A) OR ISKEYDOWN(KEY_LEFT) THEN
        SETVELOCITY2D "player", -MOVE_SPEED, playerVel.y
    ELSEIF ISKEYDOWN(KEY_D) OR ISKEYDOWN(KEY_RIGHT) THEN
        SETVELOCITY2D "player", MOVE_SPEED, playerVel.y
    ELSE
        SETVELOCITY2D "player", 0, playerVel.y
    ENDIF
    
    REM Jump
    IF (ISKEYPRESSED(KEY_W) OR ISKEYPRESSED(KEY_UP) OR ISKEYPRESSED(KEY_SPACE)) AND onGround AND canJump THEN
        APPLYIMPULSE2D "player", 0, JUMP_FORCE, 0, 0
        canJump = FALSE
    ENDIF
    
    REM Reset jump when key is released
    IF ISKEYUP(KEY_W) AND ISKEYUP(KEY_UP) AND ISKEYUP(KEY_SPACE) THEN
        canJump = TRUE
    ENDIF
    
    REM Limit fall speed
    IF playerVel.y < -MAX_FALL_SPEED THEN
        SETVELOCITY2D "player", playerVel.x, -MAX_FALL_SPEED
    ENDIF
    
    REM Check coin collection
    IF CHECKCOLLISION("player", "coin1") THEN
        score = score + 10
        SETACTIVE "coin1", FALSE
        SETVISIBLE "coinSprite1", FALSE
    ENDIF
    
    IF CHECKCOLLISION("player", "coin2") THEN
        score = score + 10
        SETACTIVE "coin2", FALSE
        SETVISIBLE "coinSprite2", FALSE
    ENDIF
    
    IF CHECKCOLLISION("player", "coin3") THEN
        score = score + 10
        SETACTIVE "coin3", FALSE
        SETVISIBLE "coinSprite3", FALSE
    ENDIF
    
    IF CHECKCOLLISION("player", "coin4") THEN
        score = score + 10
        SETACTIVE "coin4", FALSE
        SETVISIBLE "coinSprite4", FALSE
    ENDIF
    
    REM Reset game if all coins collected
    IF score >= 40 THEN
        RESETGAME()
    ENDIF
    
    REM Render
    CLEARSCREEN 135, 206, 235
    
    REM Draw platforms
    DRAWRECTANGLE 0, 380, 1024, 40, 139, 69, 19, 255
    DRAWRECTANGLE 200, 180, 200, 20, 139, 69, 19, 255
    DRAWRECTANGLE 500, 280, 200, 20, 139, 69, 19, 255
    DRAWRECTANGLE 375, 380, 150, 20, 139, 69, 19, 255
    DRAWRECTANGLE 660, 430, 180, 20, 139, 69, 19, 255
    
    REM Draw coins
    IF GETACTIVE "coin1" THEN DRAWCIRCLE 300, 250, 8, 255, 215, 0, 255
    IF GETACTIVE "coin2" THEN DRAWCIRCLE 600, 350, 8, 255, 215, 0, 255
    IF GETACTIVE "coin3" THEN DRAWCIRCLE 450, 450, 8, 255, 215, 0, 255
    IF GETACTIVE "coin4" THEN DRAWCIRCLE 750, 500, 8, 255, 215, 0, 255
    
    REM Draw player
    DRAWRECTANGLE playerPos.x - 16, playerPos.y - 24, 32, 48, 255, 0, 0, 255
    
    REM Draw UI
    DRAWTEXT "Score: " + STR(score), 10, 10, 24, 255, 255, 255
    DRAWTEXT "Use Arrow Keys/WASD to move, Space/W/Up to jump", 10, 40, 16, 255, 255, 255
    DRAWTEXT "Collect all coins to win!", 10, 60, 16, 255, 255, 255
    
    REM Draw debug info
    DRAWTEXT "Position: " + STR(INT(playerPos.x)) + ", " + STR(INT(playerPos.y)), 10, 90, 12, 200, 200, 200
    DRAWTEXT "Velocity: " + STR(INT(playerVel.x)) + ", " + STR(INT(playerVel.y)), 10, 105, 12, 200, 200, 200
    DRAWTEXT "On Ground: " + IIF(onGround, "YES", "NO"), 10, 120, 12, 200, 200, 200
    
WEND

REM Game reset function
FUNCTION RESETGAME()
    score = 0
    SETPOSITION2D "player", 100, 300
    SETVELOCITY2D "player", 0, 0
    
    SETACTIVE "coin1", TRUE
    SETACTIVE "coin2", TRUE
    SETACTIVE "coin3", TRUE
    SETACTIVE "coin4", TRUE
    
    SETVISIBLE "coinSprite1", TRUE
    SETVISIBLE "coinSprite2", TRUE
    SETVISIBLE "coinSprite3", TRUE
    SETVISIBLE "coinSprite4", TRUE
END FUNCTION

REM Cleanup
CLEANUPPHYSICS2D
CLOSEGRAPHICS
