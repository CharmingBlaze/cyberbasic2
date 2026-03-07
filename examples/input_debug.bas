// Input debug: exercise all input commands. Run: cyberbasic examples/input_debug.bas
InitWindow(900, 700, "Input Debug")
SetTargetFPS(60)

WHILE NOT WindowShouldClose()
  ClearBackground(20, 20, 30, 255)
  DrawText("KEYBOARD", 10, 10, 20, 255, 200, 100, 255)
  DrawText("KeyDown(W): " + STR(IsKeyDown(87)), 10, 35, 16, 255, 255, 255, 255)
  DrawText("KeyPressed(Spc): " + STR(IsKeyPressed(32)), 10, 55, 16, 255, 255, 255, 255)
  DrawText("GetKeyPressed: " + STR(GetKeyPressed()), 10, 75, 16, 255, 255, 255, 255)
  DrawText("IsKeyUp(A): " + STR(IsKeyUp(65)), 10, 95, 16, 255, 255, 255, 255)
  DrawText("MOUSE", 10, 130, 20, 255, 200, 100, 255)
  DrawText("X: " + STR(GetMouseX()) + " Y: " + STR(GetMouseY()), 10, 155, 16, 255, 255, 255, 255)
  DrawText("DeltaX: " + STR(GetMouseDeltaX()) + " DeltaY: " + STR(GetMouseDeltaY()), 10, 175, 16, 255, 255, 255, 255)
  DrawText("L: " + STR(IsMouseButtonDown(0)) + " R: " + STR(IsMouseButtonDown(1)) + " M: " + STR(IsMouseButtonDown(2)), 10, 195, 16, 255, 255, 255, 255)
  DrawText("Wheel: " + STR(GetMouseWheelMove()), 10, 215, 16, 255, 255, 255, 255)
  DrawText("GAMEPAD", 10, 250, 20, 255, 200, 100, 255)
  DrawText("Avail: " + STR(IsGamepadAvailable(0)) + " Btn0: " + STR(IsGamepadButtonDown(0, 0)) + " Axis0: " + STR(GetGamepadAxisMovement(0, 0)), 10, 275, 16, 255, 255, 255, 255)
  DrawText("ESC to close", 10, 320, 16, 150, 150, 150, 255)
  SYNC
WEND

CloseWindow()
