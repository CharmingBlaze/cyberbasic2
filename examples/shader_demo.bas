// Demo: embedded PBR-style preset shader via shader.pbr(); use handle.id with BeginShaderMode.
// cyberbasic --lint examples/shader_demo.bas
InitWindow(640, 480, "Shader demo")
SetTargetFPS(60)
SetCamera3D(0, 8, 12, 0, 0, 0, 0, 1, 0)
VAR sh = shader.pbr()
VAR sid = sh.id

mainloop
  ClearBackground(30, 30, 40, 255)
  BeginMode3D()
  BeginShaderMode(sid)
  DrawCube(0, 0, 0, 2, 2, 2, 200, 100, 255, 255)
  DrawGrid(10, 1.0)
  EndShaderMode()
  EndMode3D()
  DrawText("shader_demo: shader.pbr + BeginShaderMode", 10, 10, 18, 255, 255, 255, 255)
  SYNC
endmain

CloseWindow()
