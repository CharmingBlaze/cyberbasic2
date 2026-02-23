REM 2D Shapes Demo in CyberBasic
REM Demonstrates basic 2D graphics and shapes

REM Initialize graphics
INITGRAPHICS 800, 600, "2D Shapes Demo"

REM Main loop
WHILE NOT WINDOWSHOULDCLOSE()
    REM Clear screen with sky blue
    CLEARSCREEN 135, 206, 235
    
    REM Draw basic shapes
    REM Red rectangle
    DRAWRECTANGLE 100, 100, 200, 150, 255, 0, 0, 255
    
    REM Green circle
    DRAWCIRCLE 400, 200, 50, 0, 255, 0, 255
    
    REM Blue triangle (using rectangles to simulate)
    DRAWRECTANGLE 550, 150, 100, 100, 0, 0, 255, 255
    
    REM Yellow line (using small rectangles)
    DRAWRECTANGLE 100, 300, 200, 5, 255, 255, 0, 255
    
    REM Purple text
    DRAWTEXT "2D Shapes Demo!", 300, 50, 24, 255, 0, 255, 255
    
    REM Draw animated circle
    x = 400 + SIN(TIME() * 2) * 100
    DRAWCIRCLE x, 400, 30, 255, 165, 0, 255
    
    REM Draw grid pattern
    FOR i = 0 TO 8
        DRAWRECTANGLE i * 100, 500, 1, 100, 128, 128, 128, 255
    NEXT
    
    FOR i = 0 TO 10
        DRAWRECTANGLE 0, 500 + i * 10, 800, 1, 128, 128, 128, 255
    NEXT
    
    REM Show FPS
    fps = GETFPS()
    DRAWTEXT "FPS: " + STR(fps), 10, 10, 16, 255, 255, 255, 255
    
    REM Exit on ESC key
    IF ISKEYDOWN(KEY_ESCAPE) THEN
        EXIT
    ENDIF
    
WEND

REM Cleanup
CLOSEGRAPHICS
