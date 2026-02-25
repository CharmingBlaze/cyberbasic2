// Spinning cube: move mouse to orbit, middle wheel to zoom in/out
// Run: cyberbasic examples/spinning_cube.bas

InitWindow(800, 600, "Spinning Cube")
SetTargetFPS(60)
DisableCursor()

// Cube from mesh (no external file)
VAR mesh = GenMeshCube(2, 2, 2)
VAR cube = LoadModelFromMesh(mesh)

VAR camAngle = 0
VAR camPitch = 0
VAR camDist = 8
VAR cubeAngle = 0

WHILE NOT WindowShouldClose()
  camAngle -= GetMouseDeltaX() * 0.002
  camPitch += GetMouseDeltaY() * 0.002
  camPitch = Clamp(camPitch, -1.4, 1.4)
  camDist -= GetMouseWheelMove() * 1.5
  camDist = Clamp(camDist, 3, 25)
  cubeAngle += GetFrameTime() * 45
  CameraOrbit(0, 0, 0, camAngle, camPitch, camDist)
  ClearBackground(32, 32, 48, 255)
  DrawModelEx(cube, 0, 0, 0, 0, 1, 0, cubeAngle, 1, 1, 1, 255, 200, 100, 255)
  DrawText("Mouse: orbit  Wheel: zoom", 10, 10, 20, 255, 255, 255, 255)
WEND

EnableCursor()
UnloadModel(cube)
CloseWindow()
