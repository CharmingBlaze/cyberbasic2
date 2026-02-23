CREATEPHYSICSWORLD2D 0, -9.8
CREATEPHYSICSBODY2D "world", 0, 0, 50, 2, 0
SETPHYSICSPOSITION2D "ground", 0, -10
CREATEPHYSICSBODY2D "world", 1, 0, 2, 2, 1.0
SETPHYSICSPOSITION2D "box", 0, 5
CREATEPHYSICSBODY2D "world", 1, 1, 1, 1, 0.5
SETPHYSICSPOSITION2D "circle", 3, 8

DIM frameCount AS INTEGER
frameCount = 0

WHILE frameCount < 5
    frameCount = frameCount + 1
    STEPPHYSICS2D "world", 1.0/60.0, 8, 3
    
    GETPHYSICSPOSITION2D "box"
    DIM boxY AS FLOAT
    DIM boxX AS FLOAT
    boxY = POP()
    boxX = POP()
    
    GETPHYSICSPOSITION2D "circle"
    DIM circleY AS FLOAT
    DIM circleX AS FLOAT
    circleY = POP()
    circleX = POP()
    
    PRINT "Frame " + STR(frameCount)
    PRINT "Box: " + STR(boxX) + ", " + STR(boxY)
    PRINT "Circle: " + STR(circleX) + ", " + STR(circleY)
    
    IF frameCount = 3 THEN
        APPLYPHYSICSFORCE2D "box", 10, 0, 0, 0
    ENDIF
WEND

DESTROYPHYSICSBODY2D "world", "ground"
DESTROYPHYSICSBODY2D "world", "box"
DESTROYPHYSICSBODY2D "world", "circle"
DESTROYPHYSICSWORLD2D "world"

PRINT "Working Box2D demo completed!"
