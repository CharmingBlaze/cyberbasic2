// Minimal nav grid + agent using ai.* (delegates to navigation). Agent moves toward grid corner.
// cyberbasic --lint examples/ai_patrol.bas
InitWindow(640, 480, "AI patrol demo")
SetTargetFPS(60)
SetCamera3D(10, 14, 14, 0, 0, 0, 0, 1, 0)

VAR g = ai.navgridcreate(8, 8)
VAR a = ai.navagentcreate("", g)
ai.navagentsetposition(a, 0, 0, 0)
ai.navagentsetspeed(a, 4.0)
ai.navagentsetdestination(a, 7, 7, 0)

mainloop
  ClearBackground(20, 25, 35, 255)
  ai.navagentupdate(a, DeltaTime())
  VAR x = ai.navagentgetpositionx(a)
  VAR y = ai.navagentgetpositiony(a)
  VAR z = ai.navagentgetpositionz(a)
  BeginMode3D()
  DrawGrid(16, 1.0)
  DrawCube(x, y, z, 0.5, 0.5, 0.5, 80, 220, 140, 255)
  EndMode3D()
  DrawText("ai_patrol: ai.navgridcreate + navagent", 10, 10, 18, 255, 255, 255, 255)
  SYNC
endmain

CloseWindow()
