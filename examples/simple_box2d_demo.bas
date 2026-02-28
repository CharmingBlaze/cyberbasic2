// Simple Box2D demo - flat API (CreateWorld2D, Step2D, etc.)
let worldId = "w"
CreateWorld2D(worldId, 0, -9.8)

let groundBody = CreateBody2D(worldId, "ground", 0, 0, 0, -10, 0, 25, 1)

let boxBody = CreateBody2D(worldId, "box", 2, 0, 0, 5, 1.0, 1, 1)

let frameCount = 0

WHILE frameCount < 5
    let frameCount = frameCount + 1

    Step2D(worldId, 1.0/60.0, 8, 3)

    let boxX = GetPositionX2D(worldId, boxBody)
    let boxY = GetPositionY2D(worldId, boxBody)

    PRINT "Frame " + STR(frameCount)
    PRINT "Box position: " + STR(boxX) + ", " + STR(boxY)

    IF frameCount = 3 THEN
        ApplyForce2D(worldId, boxBody, 10, 0)
    ENDIF
WEND

DestroyBody2D(worldId, groundBody)
DestroyBody2D(worldId, boxBody)
DestroyWorld2D(worldId)

PRINT "Simple Box2D demo completed!"
