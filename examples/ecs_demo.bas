// ECS demo: create world, entities with Transform/Health (and Sprite), query and draw
// See docs/ECS_GUIDE.md for the ECS API.

InitWindow(800, 600, "ECS Demo")
SetTargetFPS(60)

VAR wid = ECS.CreateWorld()
VAR e1 = ECS.CreateEntity(wid)
ECS.AddComponent(wid, e1, "Transform", 200, 300, 0)
ECS.AddComponent(wid, e1, "Health", 100, 100)

VAR e2 = ECS.CreateEntity(wid)
ECS.AddComponent(wid, e2, "Transform", 500, 250, 0)
ECS.AddComponent(wid, e2, "Health", 80, 80)
ECS.AddComponent(wid, e2, "Sprite", "hero", 1)

WHILE NOT WindowShouldClose()
  ClearBackground(30, 30, 40, 255)
  VAR n = ECS.QueryCount(wid, "Transform")
  VAR i = 0
  FOR i = 0 TO n - 1
    VAR eid = ECS.QueryEntity(wid, "Transform", i)
    VAR x = ECS.GetTransformX(wid, eid)
    VAR y = ECS.GetTransformY(wid, eid)
    DrawCircle(x, y, 25, 100, 200, 255, 255)
    VAR h = ECS.GetHealthCurrent(wid, eid)
    VAR mx = ECS.GetHealthMax(wid, eid)
    DrawText("HP " + STR(h) + "/" + STR(mx), x - 20, y - 45, 14, 255, 255, 255, 255)
  NEXT i
  DrawText("ECS Demo: " + STR(n) + " entities with Transform", 10, 10, 18, 255, 255, 255, 255)
WEND

ECS.DestroyWorld(wid)
CloseWindow()
