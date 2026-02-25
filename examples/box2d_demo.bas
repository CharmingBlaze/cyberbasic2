// Box2D demo: click to spawn boxes that fall to the ground - flat API
// Run: cyberbasic examples/box2d_demo.bas

InitWindow(800, 450, "Box2D Demo - Click to spawn")
SetTargetFPS(60)

BOX2D.CreateWorld("w", 0, -10)
BOX2D.CreateBody("w", "ground", 0, 0, 0, 0, 1, 50, 0.5)
// One initial box: type 2=dynamic, shape 0=box, x=2 y=0, density 1, hx=0.5 hy=0.5
BOX2D.CreateBody("w", "box1", 2, 0, 2, 0, 1, 0.5, 0.5)

VAR scale = 50
VAR ox = 400
VAR oy = 350

WHILE NOT WindowShouldClose()
  IF IsMouseButtonPressed(0) THEN
    VAR mx = GetMouseX()
    VAR my = GetMouseY()
    BOX2D.CreateBodyAtScreen("w", mx, my, scale, ox, oy)
  ENDIF

  BOX2D.Step("w", 0.016, 8, 3)

  ClearBackground(40, 44, 52, 255)
  DrawText("Click to spawn boxes - close window to exit", 10, 10, 20, 255, 255, 255, 255)
  DrawRectangle(50, 325, 700, 50, 128, 128, 128, 255)

  VAR n = BOX2D.GetBodyCount("w")
  VAR i = 0
  FOR i = 0 TO n - 1
    VAR id = BOX2D.GetBodyId("w", i)
    IF id <> "ground" THEN
      VAR x = BOX2D.GetPositionX("w", id)
      VAR y = BOX2D.GetPositionY("w", id)
      VAR sx = ox + x * scale
      VAR sy = oy - y * scale
      DrawRectangle(Int(sx - 25), Int(sy - 25), 50, 50, 255, 200, 100, 255)
    ENDIF
  NEXT i
WEND

BOX2D.DestroyWorld("w")
CloseWindow()
