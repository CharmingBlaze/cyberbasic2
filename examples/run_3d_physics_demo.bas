// 3D physics demo: raylib + BULLET foreign API (no VM opcodes)
// Run: cyberbasic examples/run_3d_physics_demo.bas

InitWindow(800, 600, "3D Physics Demo")
SetWindowPosition(120, 80)
RestoreWindow()
SetTargetFPS(60)

// Camera: position (10,10,10), look at (0,0,0), up (0,1,0)
SetCamera3D(10, 10, 10, 0, 0, 0, 0, 1, 0)

// Physics world with gravity down
BULLET.CreateWorld("world", 0, -9.81, 0)
// Ground: flat box at y=-1, half extents (10, 0.5, 10), mass 0 = static
BULLET.CreateBox("world", "ground", 0, -1, 0, 10, 0.5, 10, 0)
// Falling box: at (0, 5, 0), half extents (1,1,1), mass 1
BULLET.CreateBox("world", "box1", 0, 5, 0, 1, 1, 1, 1)
// Falling sphere: at (2, 8, 0), radius 0.5, mass 1
BULLET.CreateSphere("world", "ball1", 2, 8, 0, 0.5, 1)

REPEAT
  BULLET.Step("world", 0.016)

  BeginDrawing()
  ClearBackground(32, 32, 48, 255)

  BeginMode3D()
  // Ground plane (wide flat cube)
  DrawCube(0, -1, 0, 20, 1, 20, 100, 100, 110, 255)
  DrawCubeWires(0, -1, 0, 20, 1, 20, 60, 60, 70, 255)

  let bx = BULLET.GetPositionX("world", "box1")
  let by = BULLET.GetPositionY("world", "box1")
  let bz = BULLET.GetPositionZ("world", "box1")
  DrawCube(bx, by, bz, 2, 2, 2, 255, 180, 80, 255)
  DrawCubeWires(bx, by, bz, 2, 2, 2, 200, 120, 40, 255)

  let sx = BULLET.GetPositionX("world", "ball1")
  let sy = BULLET.GetPositionY("world", "ball1")
  let sz = BULLET.GetPositionZ("world", "ball1")
  DrawSphere(sx, sy, sz, 0.5, 100, 200, 255, 255)
  DrawSphereWires(sx, sy, sz, 0.5, 80, 160, 200, 255)

  DrawGrid(20, 1.0)
  EndMode3D()

  DrawText("3D Physics - BULLET + Raylib", 10, 10, 20, 255, 255, 255, 255)
  DrawText("Close window to exit", 10, 35, 16, 200, 200, 200, 255)
  EndDrawing()
UNTIL WindowShouldClose()

BULLET.DestroyBody("world", "ground")
BULLET.DestroyBody("world", "box1")
BULLET.DestroyBody("world", "ball1")
CloseWindow()
PRINT "3D Physics demo finished."
