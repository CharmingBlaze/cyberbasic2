// Spinning cube with high-level commands: MouseOrbitCamera, RotateModel, DrawModelSimple, LoadCube, Background, Print
// Run: cyberbasic examples/spinning_cube_simple.bas

InitWindow(800, 600, "Spinning Cube")
SetTargetFPS(60)
DisableCursor()

VAR cube = LoadCube(2)
SetModelColor(cube, 255, 200, 100, 255)

WHILE NOT WindowShouldClose()
  MouseOrbitCamera()
  RotateModel(cube, 45)
  Background(32, 32, 48)
  DrawModelSimple(cube, 0, 0, 0)
  DrawTextSimple("Mouse: orbit  Wheel: zoom", 10, 10)
WEND

EnableCursor()
UnloadModel(cube)
CloseWindow()
