BULLET.CreateWorld("world1", 0, -9.81, 0)
BULLET.CreateBox("world1", "ground", 0, 0, 0, 10, 0.5, 10, 0)
BULLET.CreateBox("world1", "box1", 0, 5, 0, 0.5, 0.5, 0.5, 1)

DIM hit AS INTEGER
DIM x AS FLOAT
DIM y AS FLOAT
DIM z AS FLOAT
hit = BULLET.RayCast("world1", 0, 10, 0, 0, -1, 0, 20)
PRINT hit
IF hit THEN
    x = BULLET.GetRayCastHitX()
    y = BULLET.GetRayCastHitY()
    z = BULLET.GetRayCastHitZ()
    PRINT x
    PRINT y
    PRINT z
ENDIF

BULLET.Step("world1", 0.1)
x = BULLET.GetPositionX("world1", "box1")
y = BULLET.GetPositionY("world1", "box1")
z = BULLET.GetPositionZ("world1", "box1")
PRINT x
PRINT y
PRINT z
