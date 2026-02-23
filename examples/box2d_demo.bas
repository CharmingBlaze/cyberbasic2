// Box2D demo: click to spawn boxes that fall to the ground
RL.InitWindow(800, 450, "Box2D Demo - Click to spawn")
RL.SetTargetFPS(60)

BOX2D.CreateWorld("w", 0, -10)
BOX2D.CreateBody("w", "ground", 0, 0, 0, 0, 1, 50, 0.5)
// One initial box
BOX2D.CreateBody("w", "box1", 2, 0, 0, 5, 1, 0.5, 0.5)

let scale = 50
let ox = 400
let oy = 350
REPEAT
  // Click: spawn new box at mouse position (CreateBodyAtScreen does screen->world and auto ID)
  IF RL.IsMouseButtonPressed(0) THEN
    let mx = RL.GetMouseX()
    let my = RL.GetMouseY()
    BOX2D.CreateBodyAtScreen("w", mx, my, scale, ox, oy)
  ENDIF

  BOX2D.Step("w", 0.016, 8, 3)

  RL.BeginDrawing()
  RL.ClearBackground(40, 44, 52, 255)
  RL.DrawText("Click to spawn boxes - close window to exit", 10, 10, 20, 255, 255, 255, 255)
  RL.DrawRectangle(50, 325, 700, 50, 128, 128, 128, 255)

  // Draw all bodies except ground
  let n = BOX2D.GetBodyCount("w")
  FOR i = 0 TO n - 1
    let id = BOX2D.GetBodyId("w", i)
    IF id <> "ground" THEN
      let x = BOX2D.GetPositionX("w", id)
      let y = BOX2D.GetPositionY("w", id)
      let sx = ox + x * scale
      let sy = oy - y * scale
      RL.DrawRectangle(Int(sx - 25), Int(sy - 25), 50, 50, 255, 200, 100, 255)
    ENDIF
  NEXT i

  RL.EndDrawing()
UNTIL RL.WindowShouldClose()

BOX2D.DestroyWorld("w")
RL.CloseWindow()
