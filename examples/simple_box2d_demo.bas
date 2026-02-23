// Simple Box2D demo - uses BOX2D.* API only (no POP())
let worldId = "w"
BOX2D.CreateWorld(worldId, 0, -9.8)

let groundBody = BOX2D.CreateBody(worldId, "ground", 0, 0, 0, -10, 0, 25, 1)

let boxBody = BOX2D.CreateBody(worldId, "box", 2, 0, 0, 5, 1.0, 1, 1)

let frameCount = 0

WHILE frameCount < 5
    let frameCount = frameCount + 1

    BOX2D.Step(worldId, 1.0/60.0, 8, 3)

    let boxX = BOX2D.GetPositionX(worldId, boxBody)
    let boxY = BOX2D.GetPositionY(worldId, boxBody)

    PRINT "Frame " + STR(frameCount)
    PRINT "Box position: " + STR(boxX) + ", " + STR(boxY)

    IF frameCount = 3 THEN
        BOX2D.ApplyForce(worldId, boxBody, 10, 0)
    ENDIF
WEND

BOX2D.DestroyBody(worldId, groundBody)
BOX2D.DestroyBody(worldId, boxBody)
BOX2D.DestroyWorld(worldId)

PRINT "Simple Box2D demo completed!"
