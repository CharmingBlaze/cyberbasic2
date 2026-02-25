// Simplest Box2D demo - one box, ground, window - current API
// Run: cyberbasic examples/simplest_box2d_demo.bas

BOX2D.CreateWorld("world", 0, -9.8)
// Ground: static (0), box (0), x=0 y=0, density 1, half-width 25, half-height 0.5
BOX2D.CreateBody("world", "ground", 0, 0, 0, 0, 1, 25, 0.5)
// Box: dynamic (2), box (0), x=0 y=5, density 1, half 0.5 0.5
BOX2D.CreateBody("world", "box", 2, 0, 0, 5, 1, 0.5, 0.5)

InitWindow(500, 400, "Simplest Box2D")
SetTargetFPS(60)
VAR scale = 40
VAR ox = 250
VAR oy = 300

WHILE NOT WindowShouldClose()
  BOX2D.Step("world", 1.0 / 60.0, 8, 3)

  ClearBackground(44, 48, 56, 255)
  VAR bx = BOX2D.GetPositionX("world", "box")
  VAR by = BOX2D.GetPositionY("world", "box")
  VAR sx = ox + bx * scale
  VAR sy = oy - by * scale
  DrawRectangle(0, 320, 500, 80, 80, 80, 80, 255)
  DrawRectangle(Int(sx - 20), Int(sy - 20), 40, 40, 255, 180, 80, 255)
  DrawText("Simplest Box2D - box falls to ground", 10, 10, 18, 255, 255, 255, 255)
  DrawText("Close window to exit", 10, 375, 14, 180, 180, 180, 255)
WEND

BOX2D.DestroyWorld("world")
CloseWindow()
