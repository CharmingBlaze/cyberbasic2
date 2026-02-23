CREATEPHYSICSWORLD2D 0, -9.8

CREATEPHYSICSBODY2D "world", 0, 0, 50, 2, 0
SETPHYSICSPOSITION2D "ground", 0, -10
SETPHYSICSDENSITY2D "ground", 1.0
SETPHYSICSFRICTION2D "ground", 0.5

CREATEPHYSICSBODY2D "world", 1, 0, 2, 2, 1.0
SETPHYSICSPOSITION2D "box", 0, 5
SETPHYSICSDENSITY2D "box", 1.0
SETPHYSICSFRICTION2D "box", 0.3
SETPHYSICSRESTITUTION2D "box", 0.5

CREATEPHYSICSBODY2D "world", 1, 1, 1, 1, 0.5
SETPHYSICSPOSITION2D "circle", 3, 8
SETPHYSICSDENSITY2D "circle", 0.5
SETPHYSICSFRICTION2D "circle", 0.2
SETPHYSICSRESTITUTION2D "circle", 0.8

CREATEPHYSICSBODY2D "world", 1, 1, 0.5, 0.5, 2.0
SETPHYSICSPOSITION2D "ball", -3, 10
SETPHYSICSDENSITY2D "ball", 2.0
SETPHYSICSFRICTION2D "ball", 0.1
SETPHYSICSRESTITUTION2D "ball", 0.9

DIM frameCount AS INTEGER
frameCount = 0

WHILE frameCount < 10
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
    
    GETPHYSICSPOSITION2D "ball"
    DIM ballY AS FLOAT
    DIM ballX AS FLOAT
    ballY = POP()
    ballX = POP()
    
    PRINT "Frame " + STR(frameCount)
    PRINT "Box: " + STR(boxX) + ", " + STR(boxY)
    PRINT "Circle: " + STR(circleX) + ", " + STR(circleY)
    PRINT "Ball: " + STR(ballX) + ", " + STR(ballY)
    
    IF frameCount = 5 THEN
        APPLYPHYSICSFORCE2D "box", 10, 0, 0, 0
    ENDIF
    
    IF frameCount = 7 THEN
        APPLYPHYSICSIMPULSE2D "circle", 5, 5, 0, 0
    ENDIF
    
    RAYCAST2D "world", -5, 0, 5, 0
    DIM hit AS BOOLEAN
    DIM hitY AS FLOAT
    DIM hitX AS FLOAT
    hitY = POP()
    hitX = POP()
    hit = POP()
    
    IF hit THEN
        PRINT "Raycast hit at: " + STR(hitX) + ", " + STR(hitY)
    ENDIF
    
    CHECKCOLLISION2D "box", "circle"
    DIM colliding AS BOOLEAN
    colliding = POP()
    
    IF colliding THEN
        PRINT "Box and Circle are colliding!"
    ENDIF
    
    GETPHYSICSVELOCITY2D "box"
    DIM velY AS FLOAT
    DIM velX AS FLOAT
    velY = POP()
    velX = POP()
    PRINT "Box velocity: " + STR(velX) + ", " + STR(velY)
    
    GETPHYSICSANGLE2D "box"
    DIM angle AS FLOAT
    angle = POP()
    PRINT "Box angle: " + STR(angle)
WEND

DESTROYPHYSICSBODY2D "world", "ground"
DESTROYPHYSICSBODY2D "world", "box"
DESTROYPHYSICSBODY2D "world", "circle"
DESTROYPHYSICSBODY2D "world", "ball"
DESTROYPHYSICSWORLD2D "world"

PRINT "Complete Box2D Physics Demo completed!"
