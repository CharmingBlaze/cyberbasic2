REM DBP-style simple 2D platformer with implicit loop
REM Arrow keys or WASD to move, gravity, ground collision

DIM x AS FLOAT
DIM y AS FLOAT
DIM vx AS FLOAT
DIM vy AS FLOAT
DIM onGround AS INTEGER

SUB OnStart()
  x = 400
  y = 400
  vx = 0
  vy = 0
  onGround = 0
END SUB

SUB OnUpdate(dt AS FLOAT)
  VAR speed = 250.0
  VAR gravity = 600.0
  VAR jumpVel = -350.0
  VAR groundY = 450.0

  IF IsKeyDown(KEY_LEFT) OR IsKeyDown(KEY_A) THEN vx = -speed
  IF IsKeyDown(KEY_RIGHT) OR IsKeyDown(KEY_D) THEN vx = speed
  IF NOT (IsKeyDown(KEY_LEFT) OR IsKeyDown(KEY_A) OR IsKeyDown(KEY_RIGHT) OR IsKeyDown(KEY_D)) THEN vx = 0

  IF (IsKeyPressed(KEY_SPACE) OR IsKeyPressed(KEY_UP) OR IsKeyPressed(KEY_W)) AND onGround THEN
    vy = jumpVel
    onGround = 0
  END IF

  vy = vy + gravity * dt
  x = x + vx * dt
  y = y + vy * dt

  IF y >= groundY THEN
    y = groundY
    vy = 0
    onGround = 1
  END IF

  IF x < 20 THEN x = 20
  IF x > 780 THEN x = 780
END SUB

SUB OnDraw()
  ClearBackground(40, 44, 52, 255)
  DrawRectangle(0, 460, 800, 140, 60, 60, 70, 255)
  DrawCircle(x, y, 25, 255, 200, 100, 255)
  DrawText("Arrow/WASD move, Space/Up jump", 10, 10, 20, 255, 255, 255, 255)
END SUB
