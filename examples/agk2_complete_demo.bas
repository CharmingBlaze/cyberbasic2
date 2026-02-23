DIM score AS INTEGER
DIM lives AS INTEGER
DIM playerName AS STRING
DIM gameActive AS BOOLEAN

score = 100
lives = 3
playerName = "CyberBasic"
gameActive = TRUE

DIM result AS INTEGER
result = score + 50
result = result * 2
result = result - 25

DIM canContinue AS BOOLEAN
canContinue = gameActive AND (lives > 0)

IF canContinue THEN
    PRINT "Game can continue!"
ENDIF

DIM i AS INTEGER
FOR i = 1 TO 5
    PRINT "Loop iteration"
NEXT

DIM counter AS INTEGER
counter = 0
WHILE counter < 3
    PRINT "While loop iteration"
    counter = counter + 1
WEND

LOADIMAGE "player.png"
CREATESPRITE "player", "player.png", 100, 100
SETSPRITEPOSITION "player", 400, 300
DRAWSPRITE "player"

PRINT "=== AGK2 Demo Results ==="
PRINT score
PRINT lives
PRINT playerName
PRINT result
PRINT canContinue
PRINT "Demo completed successfully!"
