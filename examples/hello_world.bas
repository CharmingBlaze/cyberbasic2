REM Hello World in CyberBasic
REM This is a simple example to demonstrate the language

let message = "Hello, CyberBasic World!"

PRINT message

REM Simple graphics example
LOADIMAGE "player.png"
CREATESPRITE "player", "player.png", 400, 300
SETSPRITEPOSITION "player", 400, 300

REM Main game loop
WHILE NOT WINDOWSHOULDCLOSE()
    DRAWSPRITE "player"
    
    REM Move player with arrow keys
    IF ISKEYDOWN(KEY_RIGHT) THEN
        SETSPRITEPOSITION "player", 410, 300
    ENDIF
    
    IF ISKEYDOWN(KEY_LEFT) THEN
        SETSPRITEPOSITION "player", 390, 300
    ENDIF
    
WEND
