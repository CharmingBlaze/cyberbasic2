CREATEPHYSICSWORLD3D "world", 0, -9.8, 0
CREATEPHYSICSBODY3D "world", "ground", 0, 10, 1, 10, 0
SETPHYSICSPOSITION3D "ground", 0, -5, 0
CREATEPHYSICSBODY3D "world", "box", 1, 2, 2, 2, 1.0
SETPHYSICSPOSITION3D "box", 0, 5, 0
CREATEPHYSICSBODY3D "world", "sphere", 1, 1, 1, 1, 0.5
SETPHYSICSPOSITION3D "sphere", 3, 8, 0
CREATEPHYSICSBODY3D "world", "ball", 1, 0.5, 0.5, 0.5, 2.0
SETPHYSICSPOSITION3D "ball", -3, 10, 0

DIM frameCount AS INTEGER
frameCount = 0

WHILE frameCount < 10
    frameCount = frameCount + 1
    STEPPHYSICS3D "world", 1.0/60.0
    
    GETPHYSICSPOSITION3D "box"
    DIM boxZ AS FLOAT
    DIM boxY AS FLOAT
    DIM boxX AS FLOAT
    boxZ = POP()
    boxY = POP()
    boxX = POP()
    
    GETPHYSICSPOSITION3D "sphere"
    DIM sphereZ AS FLOAT
    DIM sphereY AS FLOAT
    DIM sphereX AS FLOAT
    sphereZ = POP()
    sphereY = POP()
    sphereX = POP()
    
    GETPHYSICSPOSITION3D "ball"
    DIM ballZ AS FLOAT
    DIM ballY AS FLOAT
    DIM ballX AS FLOAT
    ballZ = POP()
    ballY = POP()
    ballX = POP()
    
    PRINT frameCount
    PRINT boxX
    PRINT boxY
    PRINT boxZ
    PRINT sphereX
    PRINT sphereY
    PRINT sphereZ
    PRINT ballX
    PRINT ballY
    PRINT ballZ
    
    IF frameCount = 5 THEN
        APPLYPHYSICSFORCE3D "box", 10, 0, 0, 0, 0, 0
    ENDIF
    
    IF frameCount = 7 THEN
        APPLYPHYSICSIMPULSE3D "sphere", 5, 5, 0, 0, 0, 0
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
        PRINT "Raycast hit at: " + STR(hitX) + ", " + STR(hitY) + ", " + STR(hitZ)
    ENDIF
    
    CHECKCOLLISION3D "box", "sphere"
    DIM colliding AS BOOLEAN
    colliding = POP()
    
    IF colliding THEN
        PRINT "Box and Sphere are colliding!"
    ENDIF
    
    GETPHYSICSVELOCITY3D "box"
    DIM velZ AS FLOAT
    DIM velY AS FLOAT
    DIM velX AS FLOAT
    velZ = POP()
    velY = POP()
    velX = POP()
    PRINT "Box velocity: " + STR(velX) + ", " + STR(velY) + ", " + STR(velZ)
    
    GETPHYSICSROTATION3D "box"
    DIM rotZ AS FLOAT
    DIM rotY AS FLOAT
    DIM rotX AS FLOAT
    rotZ = POP()
    rotY = POP()
    rotX = POP()
    PRINT "Box rotation: " + STR(rotX) + ", " + STR(rotY) + ", " + STR(rotZ)
WEND

DESTROYPHYSICSBODY3D "world", "ground"
DESTROYPHYSICSBODY3D "world", "box"
DESTROYPHYSICSBODY3D "world", "sphere"
DESTROYPHYSICSBODY3D "world", "ball"
DESTROYPHYSICSWORLD3D "world"

PRINT "Bullet Physics demo completed!"
