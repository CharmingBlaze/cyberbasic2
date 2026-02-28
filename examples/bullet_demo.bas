CreateWorld3D("world1", 0, -9.81, 0)
CreateBox3D("world1", "ground", 0, 0, 0, 10, 0.5, 10, 0)
CreateBox3D("world1", "box1", 0, 5, 0, 0.5, 0.5, 0.5, 1)

DIM hit AS INTEGER
DIM x AS FLOAT
DIM y AS FLOAT
DIM z AS FLOAT
hit = RayCastFromDir3D("world1", 0, 10, 0, 0, -1, 0, 20)
PRINT hit
IF hit THEN
    x = RayHitX3D()
    y = RayHitY3D()
    z = RayHitZ3D()
    PRINT x
    PRINT y
    PRINT z
ENDIF

Step3D("world1", 0.1)
x = GetPositionX3D("world1", "box1")
y = GetPositionY3D("world1", "box1")
z = GetPositionZ3D("world1", "box1")
PRINT x
PRINT y
PRINT z
