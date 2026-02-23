INITGRAPHICS3D 1024, 768, "3D Physics Demo"
CREATECAMERA "main", 10, 10, 10

CREATEPHYSICSWORLD3D "world", 0, -9.8, 0
CREATEPHYSICSBODY3D "world", "ground", 0, 20, 1, 20, 0
SETPHYSICSPOSITION3D "ground", 0, -10, 0
CREATEPHYSICSBODY3D "world", "box1", 1, 2, 2, 2, 1.0
SETPHYSICSPOSITION3D "box1", 0, 5, 0
CREATEPHYSICSBODY3D "world", "sphere1", 1, 1, 1, 1, 0.5
SETPHYSICSPOSITION3D "sphere1", 3, 8, 0

frameCount = 0

WHILE frameCount < 100
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
    box1Z = POP()
    box1Y = POP()
    box1X = POP()
    
    GETPHYSICSPOSITION3D "sphere1"
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
    
    IF frameCount = 75 THEN
        APPLYPHYSICSIMPULSE3D "sphere1", 10, 10, 0, 0, 0, 0
    ENDIF
    
    SYNC
WEND

DESTROYPHYSICSBODY3D "world", "ground"
DESTROYPHYSICSBODY3D "world", "box1"
DESTROYPHYSICSBODY3D "world", "sphere1"
DESTROYPHYSICSWORLD3D "world"
CLOSEGRAPHICS3D

PRINT "Working visual demo completed!"
