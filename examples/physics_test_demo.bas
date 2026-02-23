CREATEPHYSICSWORLD3D "world", 0, -9.8, 0
CREATEPHYSICSBODY3D "world", "ground", 0, 20, 1, 20, 0
SETPHYSICSPOSITION3D "ground", 0, -10, 0
CREATEPHYSICSBODY3D "world", "box1", 1, 2, 2, 2, 1.0
SETPHYSICSPOSITION3D "box1", 0, 5, 0
CREATEPHYSICSBODY3D "world", "sphere1", 1, 1, 1, 1, 0.5
SETPHYSICSPOSITION3D "sphere1", 3, 8, 0

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

PRINT "Initial positions:"
PRINT box1X
PRINT box1Y
PRINT box1Z
PRINT sphere1X
PRINT sphere1Y
PRINT sphere1Z

APPLYPHYSICSFORCE3D "box1", 10, 0, 0, 0, 0, 0
STEPPHYSICS3D "world", 1.0/60.0
GETPHYSICSPOSITION3D "box1"
box1Z = POP()
box1Y = POP()
box1X = POP()

APPLYPHYSICSIMPULSE3D "sphere1", 5, 5, 0, 0, 0, 0
STEPPHYSICS3D "world", 1.0/60.0
GETPHYSICSPOSITION3D "sphere1"
sphere1Z = POP()
sphere1Y = POP()
sphere1X = POP()

PRINT "After forces:"
PRINT box1X
PRINT box1Y
PRINT box1Z
PRINT sphere1X
PRINT sphere1Y
PRINT sphere1Z

RAYCAST3D "world", -5, 0, 0, 5, 0, 0
DIM hit AS BOOLEAN
DIM hitZ AS FLOAT
DIM hitY AS FLOAT
DIM hitX AS FLOAT
hitZ = POP()
hitY = POP()
hitX = POP()
hit = POP()
PRINT "Raycast result:"
PRINT hit

CHECKCOLLISION3D "box1", "sphere1"
DIM colliding AS BOOLEAN
colliding = POP()
PRINT "Collision result:"
PRINT colliding

GETPHYSICSVELOCITY3D "box1"
DIM velZ AS FLOAT
DIM velY AS FLOAT
DIM velX AS FLOAT
velZ = POP()
velY = POP()
velX = POP()
PRINT "Box1 velocity:"
PRINT velX
PRINT velY
PRINT velZ

GETPHYSICSROTATION3D "box1"
DIM rotZ AS FLOAT
DIM rotY AS FLOAT
DIM rotX AS FLOAT
rotZ = POP()
rotY = POP()
rotX = POP()
PRINT "Box1 rotation:"
PRINT rotX
PRINT rotY
PRINT rotZ

SETPHYSICSMASS3D "sphere1", 5.0
STEPPHYSICS3D "world", 1.0/60.0

DESTROYPHYSICSBODY3D "world", "ground"
DESTROYPHYSICSBODY3D "world", "box1"
DESTROYPHYSICSBODY3D "world", "sphere1"
DESTROYPHYSICSWORLD3D "world"

PRINT "Physics test demo completed!"
