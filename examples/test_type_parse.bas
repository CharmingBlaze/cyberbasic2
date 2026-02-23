REM Test TYPE/ENDTYPE parsing (Phase 1)
TYPE Key
    W = 87
    A = 65
    S = 83
    D = 68
    Space = 32
ENDTYPE

TYPE Color
    White
    Red
    Gray
ENDTYPE

PRINT "Parse and compile OK"
PRINT "Key.W = ", Key.W
PRINT "Key.Space = ", Key.Space
PRINT "Color.White = ", Color.White
PRINT "Color.Red = ", Color.Red
