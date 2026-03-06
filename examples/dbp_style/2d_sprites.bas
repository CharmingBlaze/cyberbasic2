REM DBP-style 2D sprites with implicit loop (OnStart, OnUpdate, OnDraw)
REM Uses LoadImage + Sprite. If player.png missing, draws a colored rectangle instead.

DIM x AS FLOAT
DIM y AS FLOAT

SUB OnStart()
  x = 400
  y = 300
END SUB

SUB OnUpdate(dt AS FLOAT)
  IF IsKeyDown(KEY_W) THEN y = y - 200 * dt
  IF IsKeyDown(KEY_S) THEN y = y + 200 * dt
  IF IsKeyDown(KEY_A) THEN x = x - 200 * dt
  IF IsKeyDown(KEY_D) THEN x = x + 200 * dt
END SUB

SUB OnDraw()
  ClearBackground(20, 20, 30, 255)
  DrawCircle(x, y, 30, 255, 100, 100, 255)
  DrawText("WASD to move - DBP-style implicit loop", 10, 10, 20, 255, 255, 255, 255)
END SUB
