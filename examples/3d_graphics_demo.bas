REM 3D Graphics Demo in CyberBasic
REM Demonstrates 3D graphics capabilities

REM Initialize 3D graphics
INITGRAPHICS3D 1024, 768, "3D Graphics Demo"

REM Create 3D camera
CREATECAMERA "main", 10, 10, 10
SETCAMERATARGET "main", 0, 0, 0

REM Load 3D models
LOADMODEL "cube.obj"
LOADMODEL "sphere.obj"
LOADMODEL "plane.obj"

REM Create 3D objects
CREATEPHYSICSBODY3D "cube1", BODY_DYNAMIC, SHAPE_BOX, 2, 2, 2, 0, 5, 0, 1.0
CREATEPHYSICSBODY3D "sphere1", BODY_DYNAMIC, SHAPE_SPHERE, 1, 1, 1, 3, 8, 0, 1.0
CREATEPHYSICSBODY3D "ground", BODY_STATIC, SHAPE_BOX, 20, 1, 20, 0, -1, 0, 0.0

REM Main 3D rendering loop
WHILE NOT WINDOWSHOULDCLOSE()
    REM Clear screen with sky blue
    CLEARSCREEN 135, 206, 235
    
    REM Begin 3D mode
    BEGIN3DMODE "main"
    
    REM Draw 3D models
    DRAWMODEL3D "cube1", 0, 5, 0, 1.0
    DRAWMODEL3D "sphere1", 3, 8, 0, 1.0
    DRAWMODEL3D "ground", 0, -1, 0, 1.0
    
    REM Draw 3D grid for reference
    DRAWGRID3D 10, 1.0
    
    REM Draw 3D axes
    DRAWAXES3D
    
    REM End 3D mode
    END3DMODE
    
    REM Handle input
    IF ISKEYDOWN(KEY_W) THEN
        MOVECAMERAFORWARD "main", 0.1
    ENDIF
    
    IF ISKEYDOWN(KEY_S) THEN
        MOVECAMERABACKWARD "main", 0.1
    ENDIF
    
    IF ISKEYDOWN(KEY_A) THEN
        MOVECAMERALEFT "main", 0.1
    ENDIF
    
    IF ISKEYDOWN(KEY_D) THEN
        MOVECAMERARIGHT "main", 0.1
    ENDIF
    
    REM Update 3D physics
    STEPPHYSICS3D 1.0/60.0
    
    REM Show FPS
    fps = GETFPS()
    DRAWTEXT "FPS: " + STR(fps), 10, 10, 20, 255, 255, 255, 255
    
    REM Show camera position
    camX = GETCAMERAX "main"
    camY = GETCAMERAY "main"
    camZ = GETCAMERAZ "main"
    DRAWTEXT "Camera: " + STR(camX) + ", " + STR(camY) + ", " + STR(camZ), 10, 35, 16, 255, 255, 255, 255
    
    REM Exit on ESC
    IF ISKEYDOWN(KEY_ESCAPE) THEN
        EXIT
    ENDIF
    
WEND

REM Cleanup 3D resources
CLEANUPPHYSICS3D
CLOSEGRAPHICS3D
