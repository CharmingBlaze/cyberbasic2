REM Physics Demo in CyberBasic
REM Demonstrates 2D and 3D physics integration

REM Initialize graphics
INITGRAPHICS 800, 600, "Physics Demo"

REM Initialize 2D physics world
CREATEPHYSICSWORLD2D 0, -9.81

REM Create ground (static body)
CREATEPHYSICSBODY2D "ground", BODY_2D_STATIC, SHAPE_2D_BOX, 800, 20, 0, 300, 1.0

REM Create falling boxes
CREATEPHYSICSBODY2D "box1", BODY_2D_DYNAMIC, SHAPE_2D_BOX, 50, 50, 200, 100, 1.0
CREATEPHYSICSBODY2D "box2", BODY_2D_DYNAMIC, SHAPE_2D_BOX, 50, 50, 300, 150, 1.0
CREATEPHYSICSBODY2D "box3", BODY_2D_DYNAMIC, SHAPE_2D_BOX, 50, 50, 400, 200, 1.0

REM Create bouncing ball
CREATEPHYSICSBODY2D "ball", BODY_2D_DYNAMIC, SHAPE_2D_CIRCLE, 30, 30, 150, 250, 0.8
SETRESTITUTION "ball", 0.8

REM Create 3D physics world
CREATEPHYSICSWORLD3D 0, -9.81, 0

REM Create 3D camera
CREATECAMERA "main", 10, 10, 10
SETCAMERATARGET "main", 0, 0, 0

REM Create 3D objects
LOADMODEL "cube.obj"
CREATEPHYSICSBODY3D "cube3d", BODY_DYNAMIC, SHAPE_BOX, 2, 2, 2, 0, 10, 0, 1.0

REM Main game loop
WHILE NOT WINDOWSHOULDCLOSE()
    REM Update 2D physics
    STEPPHYSICS2D 1.0/60.0, 8, 3
    
    REM Update 3D physics
    STEPPHYSICS3D 1.0/60.0
    
    REM Clear screen
    CLEARSCREEN 135, 206, 235
    
    REM Draw 2D physics objects
    DRAWPHYSICSBODY2D "ground", 100, 100, 100
    DRAWPHYSICSBODY2D "box1", 255, 0, 0
    DRAWPHYSICSBODY2D "box2", 0, 255, 0
    DRAWPHYSICSBODY2D "box3", 0, 0, 255
    DRAWPHYSICSBODY2D "ball", 255, 255, 0
    
    REM Draw 3D scene
    BEGIN3DMODE "main"
    DRAWMODEL3D "cube3d", 0, 5, 0, 1.0
    DRAWWIRECUBE 0, 0, 0, 10, 255, 255, 255
    END3DMODE
    
    REM Handle input
    IF ISKEYPRESSED(KEY_SPACE) THEN
        REM Apply impulse to ball
        APPLYIMPULSE2D "ball", 0, 500, 0, 0
    ENDIF
    
    IF ISKEYPRESSED(KEY_R) THEN
        REM Reset positions
        SETPOSITION2D "box1", 200, 100
        SETPOSITION2D "box2", 300, 150
        SETPOSITION2D "box3", 400, 200
        SETPOSITION2D "ball", 150, 250
    ENDIF
    
    REM Draw debug info
    DRAWTEXT "Press SPACE to apply impulse to ball", 10, 10, 20, 255, 255, 255
    DRAWTEXT "Press R to reset positions", 10, 35, 20, 255, 255, 255
    
    REM Show physics info
    pos = GETPOSITION2D "ball"
    vel = GETVELOCITY2D "ball"
    DRAWTEXT "Ball Position: " + STR(pos.x) + ", " + STR(pos.y), 10, 60, 16, 255, 255, 255
    DRAWTEXT "Ball Velocity: " + STR(vel.x) + ", " + STR(vel.y), 10, 80, 16, 255, 255, 255
    
WEND

REM Cleanup
CLEANUPPHYSICS2D
CLEANUPPHYSICS3D
CLOSEGRAPHICS
