REM AGK2-Inspired CyberBasic Demo
REM Demonstrates all the core features working

REM === Variable Declarations === 
DIM score AS INTEGER
DIM lives AS INTEGER
DIM playerName AS STRING
DIM gameActive AS BOOLEAN

REM === Variable Assignments ===
score = 100
lives = 3
playerName = "CyberBasic"
gameActive = TRUE

REM === Arithmetic Operations ===
DIM result AS INTEGER
result = score + 50
result = result * 2
result = result - 25

REM === String Operations ===
DIM message AS STRING
message = "Player: " + playerName

REM === Boolean Logic ===
DIM canContinue AS BOOLEAN
canContinue = gameActive AND (lives > 0)

REM === Control Flow - IF Statement ===
IF canContinue THEN
    PRINT "Game can continue!"
ENDIF

REM === Control Flow - FOR Loop ===
DIM i AS INTEGER
FOR i = 1 TO 5
    PRINT "Loop " + STR(i)
NEXT

REM === Control Flow - WHILE Loop ===
DIM counter AS INTEGER
counter = 0
WHILE counter < 3
    PRINT "While count: " + STR(counter)
    counter = counter + 1
WEND

REM === Game Commands (2D Graphics) ===
LOADIMAGE "player.png"
CREATESPRITE "player", "player.png", 100, 100
SETSPRITEPOSITION "player", 400, 300
DRAWSPRITE "player"

REM === Final Output ===
PRINT "=== AGK2 Demo Results ==="
PRINT "Score: " + STR(score)
PRINT "Lives: " + STR(lives)
PRINT "Player: " + playerName
PRINT "Result: " + STR(result)
PRINT "Message: " + message
PRINT "Can Continue: " + STR(canContinue)
PRINT "Demo completed successfully!"

REM === Program End ===
