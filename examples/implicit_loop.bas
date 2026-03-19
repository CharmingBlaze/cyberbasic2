// implicit_loop.bas - ON UPDATE / ON DRAW implicit window (no InitWindow)
window.targetfps = 60
ON UPDATE
// dt is parameter to OnUpdate
END ON
ON DRAW
ClearBackground(40, 40, 50, 255)
DrawText("Implicit loop", 20, 20, 20, 255, 255, 255, 255)
SYNC
END ON
