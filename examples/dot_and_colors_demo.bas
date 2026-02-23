// Dot notation and color constants demo
// Uses RL.ClearBackground(RL.DarkGray), RL.White, GetMousePosition().x

RL.InitWindow(800, 450, "Dot and Colors Demo")
RL.SetTargetFPS(60)

DIM mx AS Float
DIM pos AS Vector2

REPEAT
  RL.BeginDrawing()
  RL.ClearBackground(RL.DarkGray)
  pos = RL.GetMousePosition()
  mx = pos.x
  RL.DrawText("Mouse X: " + STR(mx), 20, 20, 20, RL.White)
  RL.DrawRectangle(100, 100, 200, 80, RL.Blue)
  RL.DrawRectangle(400, 200, 150, 150, RL.Yellow)
  RL.EndDrawing()
UNTIL RL.WindowShouldClose()

RL.CloseWindow()
PRINT "Done."
