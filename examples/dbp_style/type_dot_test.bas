REM DBP-style: TYPE with dot notation test
REM Verifies TYPE...END TYPE, DIM AS, p.x, p.name work

TYPE Player
  x AS FLOAT
  y AS FLOAT
  name AS STRING
END TYPE

DIM p AS Player
p.x = 100
p.y = 200
p.name = "Hero"
PRINT p.name
PRINT "x="
PRINT p.x
PRINT "y="
PRINT p.y
WaitKey()
