// 3D physics demo: raylib + flat 3D physics API (CreateWorld3D, Step3D, etc.)
// Run: cyberbasic examples/run_3d_physics_demo.bas

InitWindow(800, 600, "3D Physics Demo")
SetWindowPosition(120, 80)
RestoreWindow()
SetTargetFPS(60)

// Camera: position (10,10,10), look at (0,0,0), up (0,1,0)
SetCamera3D(10, 10, 10, 0, 0, 0, 0, 1, 0)

// Physics world with gravity down
CreateWorld3D("world", 0, -9.81, 0)
// Ground: flat box at y=-1, half extents (10, 0.5, 10), mass 0 = static
CreateBox3D("world", "ground", 0, -1, 0, 10, 0.5, 10, 0)
// Falling box: at (0, 5, 0), half extents (1,1,1), mass 1
CreateBox3D("world", "box1", 0, 5, 0, 1, 1, 1, 1)
// Falling sphere: at (2, 8, 0), radius 0.5, mass 1
CreateSphere3D("world", "ball1", 2, 8, 0, 0.5, 1)

REPEAT
  Step3D("world", 0.016)

  BeginDrawing()
  ClearBackground(32, 32, 48, 255)

  BeginMode3D()
  // Ground plane (wide flat cube)
  DrawCube(0, -1, 0, 20, 1, 20, 100, 100, 110, 255)
  DrawCubeWires(0, -1, 0, 20, 1, 20, 60, 60, 70, 255)

  let bx = GetPositionX3D("world", "box1")
  let by = GetPositionY3D("world", "box1")
  let bz = GetPositionZ3D("world", "box1")
  DrawCube(bx, by, bz, 2, 2, 2, 255, 180, 80, 255)
  DrawCubeWires(bx, by, bz, 2, 2, 2, 200, 120, 40, 255)

  let sx = GetPositionX3D("world", "ball1")
  let sy = GetPositionY3D("world", "ball1")
  let sz = GetPositionZ3D("world", "ball1")
  DrawSphere(sx, sy, sz, 0.5, 100, 200, 255, 255)
  DrawSphereWires(sx, sy, sz, 0.5, 80, 160, 200, 255)

  DrawGrid(20, 1.0)
  EndMode3D()

  DrawText("3D Physics - flat API + Raylib", 10, 10, 20, 255, 255, 255, 255)
  DrawText("Close window to exit", 10, 35, 16, 200, 200, 200, 255)
  EndDrawing()
UNTIL WindowShouldClose()

DestroyBody3D("world", "ground")
DestroyBody3D("world", "box1")
DestroyBody3D("world", "ball1")
CloseWindow()
PRINT "3D Physics demo finished."
