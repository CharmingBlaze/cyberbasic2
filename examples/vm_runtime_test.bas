DIM id AS STRING
DIM img AS STRING
id = "player"
img = "player.png"
LOADIMAGE img
CREATESPRITE id, img, 100, 200
SETSPRITEPOSITION id, 150, 250
DRAWSPRITE id
PRINT "Game opcodes completed."
