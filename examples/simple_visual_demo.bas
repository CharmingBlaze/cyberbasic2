INITGRAPHICS3D 1024, 768, "3D Physics Demo"
CREATECAMERA "main", 10, 10, 10

CREATEPHYSICSWORLD3D "world", 0, -9.8, 0
CREATEPHYSICSBODY3D "world", "ground", 0, 20, 1, 20, 0
SETPHYSICSPOSITION3D "ground", 0, -10, 0
CREATEPHYSICSBODY3D "world", "box1", 1, 2, 2, 2, 1.0
SETPHYSICSPOSITION3D "box1", 0, 5, 0
CREATEPHYSICSBODY3D "world", "sphere1", 1, 1, 1, 1, 0.5
SETPHYSICSPOSITION3D "sphere1", 3, 8, 0

DIM frameCount AS INTEGER
frameCount = 0

WHILE frameCount < 300
    frameCount = frameCount + 1
    
    CLEARSCREEN 135, 206, 235
    
    BEGIN3DMODE "main"
    
    DRAWMODEL3D "ground", 0, -10, 0, 1.0
    DRAWMODEL3D "box1", 0, 5, 0, 1.0
    DRAWMODEL3D "sphere1", 3, 8, 0, 1.0
    
    DRAWGRID3D 25, 2.0
    DRAWAXES3D
    
    END3DMODE
    
    STEPPHYSICS3D "world", 1.0/60.0
    
    GETPHYSICSPOSITION3D "box1"
    DIM box1Z AS FLOAT
    DIM box1Y AS FLOAT
    DIM box1X AS FLOAT
    box1Z = POP()
    box1Y = POP()
    box1X = POP()
    
    GETPHYSICSPOSITION3D "sphere1"
    DIM sphere1Z AS FLOAT
    DIM sphere1Y AS FLOAT
    DIM sphere1X AS FLOAT
    sphere1Z = POP()
    sphere1Y = POP()
    sphere1X = POP()
    
    DRAWTEXT frameCount, 10, 10, 16, 255, 255, 255, 255
    DRAWTEXT box1X, 10, 30, 16, 255, 255, 255, 255
    DRAWTEXT box1Y, 10, 50, 16, 255, 255, 255, 255
    DRAWTEXT box1Z, 10, 70, 16, 255, 255, 255, 255
    DRAWTEXT sphere1X, 10, 90, 16, 255, 255, 255, 255
    DRAWTEXT sphere1Y, 10, 110, 16, 255, 255, 255, 255
    DRAWTEXT sphere1Z, 10, 130, 16, 255, 255, 255, 255
    
    IF frameCount = 50 THEN
        APPLYPHYSICSFORCE3D "box1", 15, 0, 0, 0, 0, 0
    ENDIF
    
    IF frameCount = 100 THEN
        APPLYPHYSICSIMPULSE3D "sphere1", 10, 10, 0, 0, 0, 0
    ENDIF
    
    RAYCAST3D "world", -5, 0, 0, 5, 0, 0
    DIM hit AS BOOLEAN
    DIM hitZ AS FLOAT
    DIM hitY AS FLOAT
    DIM hitX AS FLOAT
    hitZ = POP()
    hitY = POP()
    hitX = POP()
    hit = POP()
    
    IF hit THEN
        DRAWTEXT "Raycast Hit!", 10, 150, 16, 255, 255, 0, 255
    ENDIF
    
    CHECKCOLLISION3D "box1", "sphere1"
    DIM colliding AS BOOLEAN
    colliding = POP()
    
    IF colliding THEN
        DRAWTEXT "Collision Detected!", 10, 170, 16, 255, 0, 0, 255
    ENDIF
    
    GETPHYSICSVELOCITY3D "box1"
    DIM velZ AS FLOAT
    DIM velY AS FLOAT
    DIM velX AS FLOAT
    velZ = POP()
    velY = POP()
    velX = POP()
    DRAWTEXT velX, 10, 190, 16, 255, 255, 0, 255
    DRAWTEXT velY, 10, 210, 16, 255, 255, 0, 255
    DRAWTEXT velZ, 10, 230, 16, 255, 255, 0, 255
    
    SYNC
WEND

DESTROYPHYSICSBODY3D "world", "ground"
DESTROYPHYSICSBODY3D "world", "box1"
DESTROYPHYSICSBODY3D "world", "sphere1"
DESTROYPHYSICSWORLD3D "world"
CLOSEGRAPHICS3D

PRINT "Visual demo completed!"
