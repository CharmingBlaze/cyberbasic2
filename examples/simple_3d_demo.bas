DIM score AS INTEGER
score = 1000
DIM lives AS INTEGER
lives = 5
DIM playerName AS STRING
playerName = "CyberBasic3D"
DIM gameActive AS BOOLEAN
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

BEGIN3DMODE "main"
DRAWMODEL3D "player", 0, 5, 0, 1.0
DRAWMODEL3D "enemy", 10, 5, 0, 1.0
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

WHILE frameCount < 3
    frameCount = frameCount + 1
    
    CLEARSCREEN 135, 206, 235
    
    BEGIN3DMODE "main"
    DRAWMODEL3D "player", 0, 5, 0, 1.0
    DRAWMODEL3D "enemy", 10, 5, 0, 1.0
    DRAWMODEL3D "ground", 0, -1, 0, 1.0
    DRAWGRID3D 25, 2.0
    DRAWAXES3D
    END3DMODE
    
    DRAWSPRITE "healthBar"
    
    fps = GETFPS()
    DRAWTEXT "FPS: " + STR(fps), 10, 50, 16, 255, 255, 255, 255
    DRAWTEXT "Score: " + STR(score), 10, 70, 16, 255, 255, 255, 255
    DRAWTEXT "Lives: " + STR(lives), 10, 90, 16, 255, 255, 255, 255
WEND

PRINT "=== Final Game Stats ==="
PRINT "Final Score: " + STR(score)
PRINT "Lives Remaining: " + STR(lives)
PRINT "Player: " + playerName
PRINT "Result: " + STR(result)
PRINT "Can Continue: " + STR(canContinue)
PRINT "Frames Processed: " + STR(frameCount)
PRINT "CyberBasic 3D Demo completed!"
