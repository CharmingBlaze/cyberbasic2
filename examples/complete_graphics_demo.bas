DIM score AS INTEGER
score = 100

INITGRAPHICS3D 1024, 768, "Complete Graphics Demo"
CREATECAMERA "main", 10, 10, 10
BEGIN3DMODE "main"
DRAWMODEL3D "cube", 0, 0, 0, 1.0
DRAWGRID3D 10, 1.0
DRAWAXES3D
END3DMODE

LOADIMAGE "player.png"
CREATESPRITE "player", "player.png", 100, 100
SETSPRITEPOSITION "player", 400, 300
DRAWSPRITE "player"

PRINT "=== Graphics Demo Results ==="
PRINT score
PRINT "3D and 2D graphics working!"
PRINT "Demo completed successfully!"
