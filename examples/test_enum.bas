REM Test ENUM: members are constants (auto 0,1,2 or explicit values)
ENUM Color : Red, Green, Blue
ENUM State : Idle = 0, Running = 10, Done = 20

PRINT "Red = ", Red
PRINT "Green = ", Green
PRINT "Blue = ", Blue
PRINT "Idle = ", Idle
PRINT "Running = ", Running
PRINT "Done = ", Done

DIM c
LET c = Green
IF c = Green THEN
  PRINT "Color is Green"
ENDIF
