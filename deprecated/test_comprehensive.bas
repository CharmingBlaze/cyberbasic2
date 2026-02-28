REM Comprehensive test of CyberBasic AGK2-inspired language

REM Test variable declarations
DIM x AS INTEGER
DIM y AS STRING
DIM z AS FLOAT
DIM flag AS BOOLEAN

REM Test assignments
x = 42
y = "Hello CyberBasic!"
z = 3.14159
flag = TRUE

REM Test arithmetic
DIM result AS INTEGER
result = x + 10
result = result * 2
result = result - 5

REM Test string concatenation
DIM message AS STRING
message = "Value: " + STR(x)

REM Test boolean logic
DIM condition AS BOOLEAN
condition = flag AND (x > 40)

REM Test IF statement
IF condition THEN
    PRINT "IF statement works!"
ENDIF

REM Test FOR loop
DIM i AS INTEGER
FOR i = 1 TO 5
    PRINT "Loop iteration: " + STR(i)
NEXT

REM Test WHILE loop
DIM counter AS INTEGER
counter = 0
WHILE counter < 3
    PRINT "While loop: " + STR(counter)
    counter = counter + 1
WEND

REM Test game commands
LOADIMAGE "test.png"
CREATESPRITE "test", "test.png", 100, 100
SETSPRITEPOSITION "test", 200, 200
DRAWSPRITE "test"

REM Test final output
PRINT "=== Test Results ==="
PRINT "x = " + STR(x)
PRINT "y = " + y
PRINT "z = " + STR(z)
PRINT "flag = " + STR(flag)
PRINT "result = " + STR(result)
PRINT "message = " + message
PRINT "condition = " + STR(condition)
PRINT "All tests completed successfully!"
