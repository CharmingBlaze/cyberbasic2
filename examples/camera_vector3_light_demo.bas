// Demo: CAMERA3D(), VECTOR3(), camera setters, LoadModel/animations, lighting API, NONE
// Uses procedural API (SetCameraPosition etc.) since dot assignment is not supported.

InitWindow(800, 600, "Camera + VECTOR3 + Light demo")
SetTargetFPS(60)

// First-person style camera
VAR fpsCamera = CAMERA3D()
SetCameraPosition(fpsCamera, 0, 2, 0)
SetCameraTarget(fpsCamera, 0, 2, -1)
SetCameraUp(fpsCamera, 0, 1, 0)
SetCameraFovy(fpsCamera, 60)
SetCameraProjection(fpsCamera, CAMERA_PERSPECTIVE())

// Third-person style camera
VAR tpsCamera = CAMERA3D()
SetCameraPosition(tpsCamera, 0, 5, 10)
SetCameraTarget(tpsCamera, 0, 0, 0)
SetCameraUp(tpsCamera, 0, 1, 0)
SetCameraFovy(tpsCamera, 45)
SetCameraProjection(tpsCamera, CAMERA_PERSPECTIVE())

SetCurrentCamera(fpsCamera)

// Lighting (stubs)
ENABLELIGHTING()
VAR light = LIGHT()
SetLightType(light, LIGHT_DIRECTIONAL())
SetLightPosition(light, 0, 10, 0)
SetLightTarget(light, 0, 0, 0)
SetLightColor(light, 255, 255, 255, 255)
SetLightIntensity(light, 1.0)
SETAMBIENTLIGHT(0.2, 0.2, 0.2)

VAR x = NONE
IF IsNull(x) THEN
  x = 1
ENDIF

WHILE NOT WindowShouldClose()
  BeginMode3D()
  DrawCube(0, 0, 0, 2, 2, 2, 200, 100, 50, 255)
  EndMode3D()
  DrawText("CAMERA3D + VECTOR3 + Lighting API", 20, 40, 18, 200, 220, 255, 255)
WEND
CloseWindow()
