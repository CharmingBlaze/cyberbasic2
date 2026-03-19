DIM m(3, 3)
m(0, 0) = 1
m(0, 1) = 2
m(0, 2) = 3
m(1, 0) = 4
m(1, 1) = 5
m(1, 2) = 6
m(2, 0) = 7
m(2, 1) = 8
m(2, 2) = 9

PRINT "Bracket access m[1,2]:"
PRINT m[1, 2]

PRINT "Bracket assignment m[2,1] = 99:"
m[2, 1] = 99
PRINT m[2, 1]

PRINT "1D array:"
DIM a(5)
a(0) = 10
a(1) = 20
a(2) = 30
PRINT a[1]
a[2] = 33
PRINT a[2]
