DIM s
s = "Hello World"
PRINT "Slice 0:5: " + s[0:5]
PRINT "Slice 6:11: " + s[6:11]
PRINT "Slice 6:: " + s[6:]
PRINT "Char 0: " + s[0]
DIM x AS INTEGER
x = 42
PRINT "Interpolation: Hello {x}"

DIM b(5)
b(1) = 10
b(2) = 20
PRINT b(1)
PRINT b(2)
DIM arr(5)
arr(1) = 10
APPEND arr, 30
PRINT arr(1)
PRINT "Done"
