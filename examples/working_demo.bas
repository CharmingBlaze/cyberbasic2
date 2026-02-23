REM AGK2-Inspired CyberBasic Working Demo

REM Variable Declarations
DIM score AS INTEGER
DIM lives AS INTEGER
DIM playerName AS STRING
DIM gameActive AS BOOLEAN

REM Variable Assignments
score = 100
lives = 3
playerName = "CyberBasic"
gameActive = TRUE

REM Arithmetic Operations
DIM result AS INTEGER
result = score + 50
result = result * 2
result = result - 25

REM Boolean Logic
DIM canContinue AS BOOLEAN
canContinue = gameActive AND (lives > 0)

REM Control Flow - IF Statement
IF canContinue THEN
    PRINT "Game can continue!"
ENDIF

REM Control Flow - FOR Loop
DIM i AS INTEGER
FOR i = 1 TO 5
    PRINT "Loop iteration"
NEXT

REM Control Flow - WHILE Loop
DIM counter AS INTEGER
counter = 0
WHILE counter < 3
    PRINT "While loop iteration"
    counter = counter + 1
WEND

REM Game Commands
LOADIMAGE "player.png"
CREATESPRITE "player", "player.png", 100, 100
SETSPRITEPOSITION "player", 400, 300
DRAWSPRITE "player"

REM Final Output
PRINT "=== Demo Results ==="
PRINT score
PRINT lives
PRINT playerName
PRINT result
PRINT canContinue
PRINT "Demo completed successfully!"
