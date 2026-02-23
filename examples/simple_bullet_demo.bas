CREATEPHYSICSWORLD3D "world", 0, -9.8, 0
CREATEPHYSICSBODY3D "world", "ground", 0, 10, 1, 10, 0
SETPHYSICSPOSITION3D "ground", 0, -5, 0
CREATEPHYSICSBODY3D "world", "box", 1, 2, 2, 2, 1.0
SETPHYSICSPOSITION3D "box", 0, 5, 0
CREATEPHYSICSBODY3D "world", "sphere", 1, 1, 1, 1, 0.5
SETPHYSICSPOSITION3D "sphere", 3, 8, 0

DIM frameCount AS INTEGER
frameCount = 0

WHILE frameCount < 5
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
    
    PRINT frameCount
    PRINT boxX
    PRINT boxY
    PRINT boxZ
    PRINT sphereX
    PRINT sphereY
    PRINT sphereZ
    
    IF frameCount = 3 THEN
        APPLYPHYSICSFORCE3D "box", 10, 0, 0, 0, 0, 0
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
        PRINT "Raycast hit"
    ENDIF
    
    GETPHYSICSVELOCITY3D "box"
    DIM velZ AS FLOAT
    DIM velY AS FLOAT
    DIM velX AS FLOAT
    velZ = POP()
    velY = POP()
    velX = POP()
    PRINT velX
    PRINT velY
    PRINT velZ
WEND

DESTROYPHYSICSBODY3D "world", "ground"
DESTROYPHYSICSBODY3D "world", "box"
DESTROYPHYSICSBODY3D "world", "sphere"
DESTROYPHYSICSWORLD3D "world"

PRINT "Simple Bullet demo completed!"
