DIM physicsWorld AS STRING
DIM groundBody AS STRING
DIM boxBody AS STRING
DIM circleBody AS STRING
DIM ballBody AS STRING

physicsWorld = CREATEPHYSICSWORLD2D 0, -9.8

groundBody = CREATEPHYSICSBODY2D physicsWorld, BODY_STATIC, SHAPE_BOX, 50, 2, 0
SETPHYSICSPOSITION2D groundBody, 0, -10
SETPHYSICSDENSITY2D groundBody, 1.0
SETPHYSICSFRICTION2D groundBody, 0.5

boxBody = CREATEPHYSICSBODY2D physicsWorld, BODY_DYNAMIC, SHAPE_BOX, 2, 2, 1.0
SETPHYSICSPOSITION2D boxBody, 0, 5
SETPHYSICSDENSITY2D boxBody, 1.0
SETPHYSICSFRICTION2D boxBody, 0.3
SETPHYSICSRESTITUTION2D boxBody, 0.5

circleBody = CREATEPHYSICSBODY2D physicsWorld, BODY_DYNAMIC, SHAPE_CIRCLE, 1, 1, 1.0
SETPHYSICSPOSITION2D circleBody, 3, 8
SETPHYSICSDENSITY2D circleBody, 0.5
SETPHYSICSFRICTION2D circleBody, 0.2
SETPHYSICSRESTITUTION2D circleBody, 0.8

ballBody = CREATEPHYSICSBODY2D physicsWorld, BODY_DYNAMIC, SHAPE_CIRCLE, 0.5, 0.5, 2.0
SETPHYSICSPOSITION2D ballBody, -3, 10
SETPHYSICSDENSITY2D ballBody, 2.0
SETPHYSICSFRICTION2D ballBody, 0.1
SETPHYSICSRESTITUTION2D ballBody, 0.9

DIM frameCount AS INTEGER
frameCount = 0

WHILE frameCount < 10
    frameCount = frameCount + 1
    
    STEPPHYSICS2D physicsWorld, 1.0/60.0, 8, 3
    
    DIM boxX AS FLOAT
    DIM boxY AS FLOAT
    GETPHYSICSPOSITION2D boxBody
    boxY = POP()
    boxX = POP()
    
    DIM circleX AS FLOAT
    DIM circleY AS FLOAT
    GETPHYSICSPOSITION2D circleBody
    circleY = POP()
    circleX = POP()
    
    DIM ballX AS FLOAT
    DIM ballY AS FLOAT
    GETPHYSICSPOSITION2D ballBody
    ballY = POP()
    ballX = POP()
    
    PRINT "Frame " + STR(frameCount)
    PRINT "Box: " + STR(boxX) + ", " + STR(boxY)
    PRINT "Circle: " + STR(circleX) + ", " + STR(circleY)
    PRINT "Ball: " + STR(ballX) + ", " + STR(ballY)
    
    IF frameCount = 5 THEN
        APPLYPHYSICSFORCE2D boxBody, 10, 0, 0, 0
    ENDIF
    
    IF frameCount = 7 THEN
        APPLYPHYSICSIMPULSE2D circleBody, 5, 5, 0, 0
    ENDIF
    
    DIM hit AS BOOLEAN
    DIM hitX AS FLOAT
    DIM hitY AS FLOAT
    RAYCAST2D physicsWorld, -5, 0, 5, 0
    hitY = POP()
    hitX = POP()
    hit = POP()
    
    IF hit THEN
        PRINT "Raycast hit at: " + STR(hitX) + ", " + STR(hitY)
    ENDIF
    
    DIM colliding AS BOOLEAN
    CHECKCOLLISION2D boxBody, circleBody
    colliding = POP()
    
    IF colliding THEN
        PRINT "Box and Circle are colliding!"
    ENDIF
WEND

DESTROYPHYSICSBODY2D physicsWorld, groundBody
DESTROYPHYSICSBODY2D physicsWorld, boxBody
DESTROYPHYSICSBODY2D physicsWorld, circleBody
DESTROYPHYSICSBODY2D physicsWorld, ballBody
DESTROYPHYSICSWORLD2D physicsWorld

PRINT "Box2D Physics Demo completed!"
