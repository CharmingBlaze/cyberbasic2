DIM score AS INTEGER
score = 1000
DIM lives AS INTEGER
lives = 5
DIM gameActive AS BOOLEAN
gameActive = TRUE

IF gameActive THEN
    PRINT "Game initialized!"
ENDIF

INITGRAPHICS3D 1024, 768, "3D Demo"
CREATECAMERA "main", 10, 10, 10

BEGIN3DMODE "main"
DRAWMODEL3D "cube", 0, 0, 0, 1.0
DRAWGRID3D 10, 1.0
DRAWAXES3D
END3DMODE

LOADIMAGE "player.png"
CREATESPRITE "player", "player.png", 100, 100
DRAWSPRITE "player"

PRINT score
PRINT lives
PRINT "3D and 2D graphics working!"
