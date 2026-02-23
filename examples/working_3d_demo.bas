DIM score AS INTEGER
DIM lives AS INTEGER
DIM playerName AS STRING
DIM gameActive AS BOOLEAN

score = 1000
lives = 5
playerName = "CyberBasic3D"
gameActive = TRUE

DIM result AS INTEGER
result = score + 500
result = result * 2
result = result - 100

DIM canContinue AS BOOLEAN
canContinue = gameActive AND (lives > 0)

IF canContinue THEN
    PRINT "Game initialized successfully!"
ENDIF

DIM i AS INTEGER
FOR i = 1 TO 3
    PRINT "Loading game assets..."
NEXT

INITGRAPHICS3D 1280, 720, "CyberBasic 3D Game"
CREATECAMERA "main", 15, 10, 15
SETCAMERATARGET "main", 0, 0, 0

LOADMODEL "player.obj"
LOADMODEL "enemy.obj"

CREATEPHYSICSBODY3D "playerBody", BODY_DYNAMIC, SHAPE_SPHERE, 1, 1, 1, 0, 5, 0, 1.0
CREATEPHYSICSBODY3D "enemyBody", BODY_DYNAMIC, SHAPE_BOX, 2, 2, 2, 10, 5, 0, 1.0
CREATEPHYSICSBODY3D "ground", BODY_STATIC, SHAPE_BOX, 50, 1, 50, 0, -1, 0, 0.0

BEGIN3DMODE "main"
DRAWMODEL3D "playerBody", 0, 5, 0, 1.0
DRAWMODEL3D "enemyBody", 10, 5, 0, 1.0
DRAWMODEL3D "ground", 0, -1, 0, 1.0
DRAWGRID3D 25, 2.0
DRAWAXES3D
END3DMODE

LOADIMAGE "ui.png"
CREATESPRITE "healthBar", "ui.png", 100, 20
SETSPRITEPOSITION "healthBar", 10, 10
DRAWSPRITE "healthBar"

DIM frameCount AS INTEGER
frameCount = 0

WHILE frameCount < 5
    frameCount = frameCount + 1
    
    CLEARSCREEN 135, 206, 235
    
    BEGIN3DMODE "main"
    DRAWMODEL3D "playerBody", 0, 5, 0, 1.0
    DRAWMODEL3D "enemyBody", 10, 5, 0, 1.0
    DRAWMODEL3D "ground", 0, -1, 0, 1.0
    DRAWGRID3D 25, 2.0
    DRAWAXES3D
    END3DMODE
    
    DRAWSPRITE "healthBar"
    
    STEPPHYSICS3D 1.0/60.0
    
    fps = GETFPS()
    DRAWTEXT "FPS: " + STR(fps), 10, 50, 16, 255, 255, 255, 255
    DRAWTEXT "Score: " + STR(score), 10, 70, 16, 255, 255, 255, 255
    DRAWTEXT "Lives: " + STR(lives), 10, 90, 16, 255, 255, 255, 255
WEND

CLEANUPPHYSICS3D
CLOSEGRAPHICS3D

PRINT "=== Final Game Stats ==="
PRINT "Final Score: " + STR(score)
PRINT "Lives Remaining: " + STR(lives)
PRINT "Player: " + playerName
PRINT "Result: " + STR(result)
PRINT "Can Continue: " + STR(canContinue)
PRINT "Frames Processed: " + STR(frameCount)
PRINT "CyberBasic 3D Demo completed!"
