CREATEPHYSICSWORLD3D "world", 0, -9.8, 0
CREATEPHYSICSBODY3D "world", "ground", 0, 20, 1, 20, 0
SETPHYSICSPOSITION3D "ground", 0, -10, 0
CREATEPHYSICSBODY3D "world", "box1", 1, 2, 2, 2, 1.0
SETPHYSICSPOSITION3D "box1", 0, 5, 0
CREATEPHYSICSBODY3D "world", "box2", 1, 1, 1, 1, 0.5
SETPHYSICSPOSITION3D "box2", 3, 8, 0
CREATEPHYSICSBODY3D "world", "sphere1", 1, 1, 1, 1, 0.5
SETPHYSICSPOSITION3D "sphere1", -3, 10, 0
CREATEPHYSICSBODY3D "world", "sphere2", 1, 0.5, 0.5, 0.5, 2.0
SETPHYSICSPOSITION3D "sphere2", 5, 12, 0

DIM frameCount AS INTEGER
frameCount = 0

WHILE frameCount < 10
    frameCount = frameCount + 1
    STEPPHYSICS3D "world", 1.0/60.0
    
    GETPHYSICSPOSITION3D "box1"
    DIM box1Z AS FLOAT
    DIM box1Y AS FLOAT
    DIM box1X AS FLOAT
    box1Z = POP()
    box1Y = POP()
    box1X = POP()
    
    GETPHYSICSPOSITION3D "box2"
    DIM box2Z AS FLOAT
    DIM box2Y AS FLOAT
    DIM box2X AS FLOAT
    box2Z = POP()
    box2Y = POP()
    box2X = POP()
    
    GETPHYSICSPOSITION3D "sphere1"
    DIM sphere1Z AS FLOAT
    DIM sphere1Y AS FLOAT
    DIM sphere1X AS FLOAT
    sphere1Z = POP()
    sphere1Y = POP()
    sphere1X = POP()
    
    GETPHYSICSPOSITION3D "sphere2"
    DIM sphere2Z AS FLOAT
    DIM sphere2Y AS FLOAT
    DIM sphere2X AS FLOAT
    sphere2Z = POP()
    sphere2Y = POP()
    sphere2X = POP()
    
    PRINT "Frame " + STR(frameCount)
    PRINT "Box1: " + STR(box1X) + ", " + STR(box1Y) + ", " + STR(box1Z)
    PRINT "Box2: " + STR(box2X) + ", " + STR(box2Y) + ", " + STR(box2Z)
    PRINT "Sphere1: " + STR(sphere1X) + ", " + STR(sphere1Y) + ", " + STR(sphere1Z)
    PRINT "Sphere2: " + STR(sphere2X) + ", " + STR(sphere2Y) + ", " + STR(sphere2Z)
    
    IF frameCount = 3 THEN
        APPLYPHYSICSFORCE3D "box1", 10, 0, 0, 0, 0, 0
    ENDIF
    
    IF frameCount = 5 THEN
        APPLYPHYSICSIMPULSE3D "sphere1", 5, 5, 0, 0, 0, 0
    ENDIF
    
    IF frameCount = 7 THEN
        SETPHYSICSVELOCITY3D "sphere2", -2, 0, 0
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
    
    CHECKCOLLISION3D "box1", "sphere1"
    DIM colliding AS BOOLEAN
    colliding = POP()
    
    IF colliding THEN
        PRINT "Box1 and Sphere1 are colliding!"
    ENDIF
    
    GETPHYSICSVELOCITY3D "box1"
    DIM velZ AS FLOAT
    DIM velY AS FLOAT
    DIM velX AS FLOAT
    velZ = POP()
    velY = POP()
    velX = POP()
    PRINT "Box1 velocity: " + STR(velX) + ", " + STR(velY) + ", " + STR(velZ)
    
    GETPHYSICSROTATION3D "box1"
    DIM rotZ AS FLOAT
    DIM rotY AS FLOAT
    DIM rotX AS FLOAT
    rotZ = POP()
    rotY = POP()
    rotX = POP()
    PRINT "Box1 rotation: " + STR(rotX) + ", " + STR(rotY) + ", " + STR(rotZ)
    
    SETPHYSICSMASS3D "sphere2", 5.0
WEND

DESTROYPHYSICSBODY3D "world", "ground"
DESTROYPHYSICSBODY3D "world", "box1"
DESTROYPHYSICSBODY3D "world", "box2"
DESTROYPHYSICSBODY3D "world", "sphere1"
DESTROYPHYSICSBODY3D "world", "sphere2"
DESTROYPHYSICSWORLD3D "world"

PRINT "Complete Bullet Physics demo completed!"
