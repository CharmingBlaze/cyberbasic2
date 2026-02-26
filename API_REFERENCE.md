# CyberBasic API Reference – All Bindings

All functions callable from BASIC. Names are **case-insensitive**. You can call with or without namespace (e.g. `InitWindow(...)` or `RL.InitWindow(...)` for raylib; use `BOX2D.*` and `BULLET.*` for physics).

For a **structured list of high-level and extended commands** (Window, Input, Math, Camera, 3D, 2D, Audio, File, Game loop, Debug, Gameplay helpers), see [COMMAND_REFERENCE.md](docs/COMMAND_REFERENCE.md).

---

## 1. Raylib (core) – `raylib_core.go`

InitWindow, SetTargetFPS, WindowShouldClose, CloseWindow, SetWindowPosition, ClearBackground, **BeginFrame**() / **EndFrame**() (alias BeginDrawing/EndDrawing), **Background**(r, g, b) (clear with RGB, alpha 255), GetFrameTime, **DeltaTime** (same as GetFrameTime; preferred for frame delta), GetFPS, GetScreenWidth, GetScreenHeight, SetWindowSize, SetWindowTitle, MaximizeWindow, MinimizeWindow, IsWindowReady, IsWindowFullscreen, GetTime, GetRandomValue, SetRandomSeed, SetWindowState, ClearWindowState, GetMonitorCount, GetCurrentMonitor, GetClipboardText, SetClipboardText, TakeScreenshot, OpenURL, IsWindowHidden, IsWindowMinimized, IsWindowMaximized, IsWindowFocused, IsWindowResized, ToggleFullscreen, RestoreWindow, GetRenderWidth, GetRenderHeight, GetMonitorName, GetMonitorWidth, GetMonitorHeight, GetMonitorRefreshRate, WaitTime, EnableEventWaiting, DisableEventWaiting, IsCursorHidden, EnableCursor, DisableCursor, IsCursorOnScreen, IsWindowState, ToggleBorderlessWindowed, SetWindowMonitor, SetWindowMinSize, SetWindowMaxSize, SetWindowOpacity, GetWindowPosition, GetWindowScaleDPI, **GetScaleDPI** (single DPI scale for UI), GetMonitorPosition, GetMonitorPhysicalWidth, GetMonitorPhysicalHeight, SetConfigFlags, SwapScreenBuffer, PollInputEvents, SetCamera2D, BeginMode2D, EndMode2D, GetWorldToScreen2D, GetScreenToWorld2D, BeginBlendMode, EndBlendMode, BeginScissorMode, EndScissorMode, BeginShaderMode, EndShaderMode, LoadShader, LoadShaderFromMemory, UnloadShader, IsShaderValid, **FileExists**

**Config/blend flags (0-arg, use with SetConfigFlags or BeginBlendMode):** FLAG_VSYNC_HINT, FLAG_FULLSCREEN_MODE, FLAG_WINDOW_RESIZABLE, FLAG_WINDOW_UNDECORATED, FLAG_WINDOW_HIDDEN, FLAG_WINDOW_MINIMIZED, FLAG_WINDOW_MAXIMIZED, FLAG_WINDOW_UNFOCUSED, FLAG_WINDOW_TOPMOST, FLAG_WINDOW_ALWAYS_RUN, FLAG_MSAA_4X_HINT, FLAG_INTERLACED_HINT, FLAG_WINDOW_HIGHDPI, FLAG_BORDERLESS_WINDOWED_MODE; BLEND_ALPHA, BLEND_ADDITIVE, BLEND_MULTIPLIED, BLEND_ADD_COLORS, BLEND_SUBTRACT_COLORS, BLEND_CUSTOM. See [Windows, scaling, and splitscreen](docs/WINDOWS_AND_VIEWS.md).

**Note:** The compiler does not inject any frame or mode calls; your code compiles exactly as written (DBPro-style). Exception: when you define **update(dt)** and **draw()** (Sub or Function) and use a game loop, the compiler injects the **hybrid loop** (see below).

**Hybrid loop (raylib_hybrid.go):** **ClearRenderQueues**(), **FlushRenderQueues**() – clear and then execute the 2D/3D/GUI render queues. **StepAllPhysics2D**(dt), **StepAllPhysics3D**(dt) – step all registered Box2D/Bullet worlds. When **update(dt)** and/or **draw()** are defined and the main loop is a game loop, the compiler invokes them automatically (GetFrameTime → physics step → update(dt) → ClearRenderQueues → draw() → FlushRenderQueues). Draw*/Gui* calls inside draw() are queued and flushed in order. See [Program Structure](docs/PROGRAM_STRUCTURE.md#hybrid-updatedraw-loop).

---

## 2. Raylib (input) – `raylib_input.go`

IsMouseButtonPressed, GetMouseX, GetMouseY, IsKeyPressed, IsKeyDown, **KeyPressed(key)** (alias), **KeyDown(key)** (alias), IsKeyReleased, IsKeyUp, GetKeyPressed, SetExitKey, IsMouseButtonDown, IsMouseButtonReleased, **GetMouseWheelMove()**, **GetMouseDeltaX()**, **GetMouseDeltaY()**, SetMousePosition, SetMouseOffset, SetMouseScale, HideCursor, ShowCursor, GetMousePosition, GetVector2X, GetVector2Y, GetVector3Z, IsMouseButtonUp, IsKeyPressedRepeat, GetCharPressed, SetMouseCursor, IsGamepadAvailable, GetGamepadName, IsGamepadButtonPressed, IsGamepadButtonDown, IsGamepadButtonReleased, GetGamepadAxisMovement, IsGamepadButtonUp, GetGamepadButtonPressed, GetGamepadAxisCount, SetGamepadMappings, SetGamepadVibration, GetTouchPointCount, GetTouchX, GetTouchY, GetTouchPosition, GetTouchPointId, GetMouseWheelMoveV  
**Key constants (0-arg):** KEY_NULL, KEY_APOSTROPHE, KEY_COMMA, KEY_MINUS, KEY_PERIOD, KEY_SLASH, KEY_ZERO … KEY_NINE, KEY_SEMICOLON, KEY_EQUAL, KEY_A … KEY_Z, KEY_LEFT_BRACKET, KEY_BACKSLASH, KEY_RIGHT_BRACKET, KEY_GRAVE, KEY_SPACE, KEY_ESCAPE, KEY_ENTER, KEY_TAB, KEY_BACKSPACE, KEY_INSERT, KEY_DELETE, KEY_RIGHT, KEY_LEFT, KEY_DOWN, KEY_UP, KEY_PAGE_UP, KEY_PAGE_DOWN, KEY_HOME, KEY_END, KEY_F1 … KEY_F12  
**Movement:** For simple movement use **GetAxisX()** / **GetAxisY()** (return -1, 0, or 1 for A/D and W/S), e.g. `x = x + speed * GetAxisX()`. For full 2D/3D use **GAME.MoveWASD**, **MoveHorizontal2D**, **Jump2D** (see §13 Raylib (game)).

---

## 3. Raylib (shapes) – `raylib_shapes.go`

SetShapesTexture, GetShapesTextureRectangle, DrawRectangle, DrawCircle, DrawLine, DrawLineV, DrawCircleLines, DrawRectangleLines, DrawTriangle, DrawTriangleLines, DrawPixel, DrawPoly, DrawEllipse, DrawRing, DrawRectangleRounded, DrawGrid, DrawFPS, DrawLineEx, DrawPixelV, DrawCircleSector, DrawCircleGradient, DrawCircleV, DrawEllipseLines, DrawRingLines, DrawRectangleV, DrawRectangleRec, DrawRectanglePro, DrawRectangleLinesEx, DrawRectangleRoundedLines, DrawPolyLines

---

## 4. Raylib (text) – `raylib_text.go`

DrawText, **DrawTextSimple**(text, x, y) — draw at (x,y), font size 20, white (for on-screen text; use PRINT for console). MeasureText, DrawTextEx, DrawTextPro, SetTextLineSpacing, TextCopy, TextIsEqual, TextLength, TextFormat, TextSubtext, TextReplace, TextInsert, TextJoin, TextSplit, GetTextSplitItem, TextAppend, TextFindIndex, TextToUpper, TextToLower, TextToPascal, TextToSnake, TextToCamel, TextToInteger, TextToFloat, GetCodepointCount, GetCodepoint, GetCodepointNext, GetCodepointPrevious, CodepointToUTF8, LoadCodepoints, UnloadCodepoints, GetLoadedCodepoint, LoadUTF8, UnloadUTF8

---

## 5. Raylib (textures) – `raylib_textures.go`

LoadTexture, UnloadTexture, LoadRenderTexture, UnloadRenderTexture, BeginTextureMode, EndTextureMode, DrawTexture, DrawTextureEx, DrawTextureRec, DrawTexturePro, LoadTextureFromImage, LoadTextureCubemap, IsTextureValid, IsRenderTextureValid, UpdateTexture, UpdateTextureRec, GenTextureMipmaps, SetTextureFilter, SetTextureWrap, DrawTextureV, DrawTextureNPatch

**2D sprite (texture) animation (raylib_anim2d.go):** CreateSpriteAnimation(textureId, frameWidth, frameHeight, framesPerRow [, totalFrames]) → animId. SetSpriteAnimationFPS(animId, fps), SetSpriteAnimationLoop(animId, loop), SetSpriteAnimationFrame(animId, frameIndex), UpdateSpriteAnimation(animId, deltaTime), GetSpriteAnimationFrame(animId) → frame index, DrawSpriteAnimation(animId, posX, posY [, scaleX, scaleY, rotation, r,g,b,a]), DestroySpriteAnimation(animId). Use for sprite-sheet animation: one texture, grid of frames, time-based playback.

---

## 6. Raylib (images) – `raylib_images.go`

LoadImage, LoadImageRaw, LoadImageAnim, GetLoadImageAnimFrames, LoadImageAnimFromMemory, LoadImageFromMemory, LoadImageFromTexture, LoadImageFromScreen, IsImageValid, UnloadImage, ExportImage, ExportImageToMemory, **ExportImageAsCode**, GenImageColor, GenImageGradientLinear, GenImageGradientRadial, GenImageGradientSquare, GenImageChecked, GenImageWhiteNoise, GenImagePerlinNoise, GenImageCellular, GenImageText, ImageCopy, ImageFromImage, ImageFromChannel, ImageText, ImageTextEx, ImageFormat, ImageToPOT, ImageCrop, ImageAlphaCrop, ImageAlphaClear, ImageAlphaMask, ImageAlphaPremultiply, ImageBlurGaussian, ImageKernelConvolution, ImageResize, ImageResizeNN, ImageResizeCanvas, ImageMipmaps, ImageDither, ImageFlipVertical, ImageFlipHorizontal, ImageRotate, ImageRotateCW, ImageRotateCCW, ImageColorTint, ImageColorInvert, ImageColorGrayscale, ImageColorContrast, ImageColorBrightness, ImageColorReplace, LoadImageColors, UnloadImageColors, GetLoadedImageColor, GetImageColor, ImageClearBackground, ImageDrawPixel, ImageDrawPixelV, ImageDrawLine, ImageDrawLineV, ImageDrawLineEx, ImageDrawCircle, ImageDrawCircleV, ImageDrawCircleLines, ImageDrawCircleLinesV, ImageDrawRectangle, ImageDrawRectangleV, ImageDrawRectangleRec, ImageDrawRectangleLines, ImageDrawTriangle, ImageDrawTriangleEx, ImageDrawTriangleLines, ImageDrawTriangleFan, ImageDrawTriangleStrip, ImageDraw

---

## 7. Raylib (3D) – `raylib_3d.go`

SetCamera3D, BeginMode3D, EndMode3D, LoadModel, LoadModelFromMesh, UnloadModel, DrawModel, DrawCube, DrawCubeWires, DrawSphere, DrawSphereWires, DrawPlane, DrawLine3D, DrawPoint3D, DrawCircle3D, DrawCubeV, DrawCylinder, DrawCylinderWires, DrawRay, DrawTriangle3D, DrawTriangleStrip3D, DrawCubeWiresV, DrawSphereEx, DrawCylinderEx, DrawCylinderWiresEx, DrawCapsule, DrawCapsuleWires, DrawModelEx, DrawModelWires, DrawBoundingBox, IsModelValid, GetModelBoundingBox, DrawModelWiresEx, DrawModelPoints, DrawModelPointsEx, DrawBillboard, DrawBillboardRec, SetModelMeshMaterial, DrawBillboardPro, LoadModelAnimations, GetModelAnimationId, UpdateModelAnimation, UpdateModelAnimationBones, UnloadModelAnimation, UnloadModelAnimations, IsModelAnimationValid. **GetModelAnimationFrameCount**(animId) → frame count (int). **Model animation state (time-based):** CreateModelAnimState(modelId, animId, fps [, loop]) → stateId, UpdateModelAnimState(stateId, deltaTime), SetModelAnimStateFrame(stateId, frameIndex), GetModelAnimStateFrame(stateId) → frame, DestroyModelAnimState(stateId). Use UpdateModelAnimState each frame with GetFrameTime(), then DrawModel as usual. **Camera (global):** **CameraOrbit**(targetX, targetY, targetZ, angleRad, pitchRad, distance) — also updates orbit state for Zoom/Rotate/UpdateCamera. **CameraZoom**(amount) — adjust orbit distance with clamping (e.g. amount = GetMouseWheelMove()). **CameraRotate**(deltaX, deltaY) — mouse-delta rotation (2 args); or **CameraRotate**(pitchRad, yawRad, rollRad) — absolute (3 args). **SetCameraTarget**(x, y, z) — set orbit/look-at target (3 args) or SetCameraTarget(cameraId, x, y, z) (4 args). **UpdateCamera**() — apply orbit state to camera. **MouseOrbitCamera**() — one call: mouse delta → rotate, wheel → zoom, then update camera. **MouseLook**() — FPS-style camera from mouse delta. **CameraLookAt**(x, y, z), **CameraMove**(dx, dy, dz), **SetCameraFOV**(fov). **SetCameraPosition**(x, y, z) — 3 args set global camera; **SetCameraPosition**(cameraId, x, y, z) — 4 args set named camera. **Camera objects:** CAMERA3D() → cameraId; SetCameraTarget(cameraId, x, y, z), SetCameraUp(cameraId, x, y, z), SetCameraFovy(cameraId, fovy), SetCameraProjection(cameraId, projection), SetCurrentCamera(cameraId); CAMERA_PERSPECTIVE(), CAMERA_ORTHOGRAPHIC() → int. **DrawModel** accepts (id, posX, posY, posZ, scale [, tint]) or (id, VECTOR3(x,y,z), scale [, tint]). **High-level model:** **LoadCube**(size) → model id (cube mesh). **SetModelColor**(modelId, r, g, b, a) — stored tint for DrawModelSimple. **RotateModel**(modelId, speedDegPerSec) — add rotation each frame. **DrawModelSimple**(id, x, y, z [, angle]) — draw at (x,y,z), scale 1, axis Y; uses SetModelColor tint and RotateModel angle when angle omitted. **Lighting (stubs):** ENABLELIGHTING(), LIGHT() → lightId, LIGHT_DIRECTIONAL() → 0, SetLightType/SetLightPosition/SetLightTarget/SetLightColor/SetLightIntensity(lightId, …), SETAMBIENTLIGHT(r, g, b).

**Fog (raylib_fog.go):** SetFog(enable, density, r, g, b), SetFogDensity(density), SetFogColor(r, g, b), EnableFog(), DisableFog(), IsFogEnabled() → 1 or 0, BeginFog(), EndFog(). Use BeginFog() before drawing, EndFog() after. Density around 0.02–0.05; color 0–255.

**Scene (scene.go):** CreateScene(sceneId), LoadScene(sceneId), UnloadScene(sceneId), SetCurrentScene(sceneId), GetCurrentScene() → sceneId or "", SetSceneWorld(sceneId, worldId), SaveScene(sceneId, path), LoadSceneFromFile(path) → sceneId. SaveScene writes scene metadata (id, worldId) to JSON; LoadSceneFromFile reads and creates the scene and sets it current.

**Views (raylib_views.go):** CreateView(viewId, x, y, width, height), SetViewTarget(viewId, renderTextureId), DrawView(viewId), GetViewX/Y/Width/Height(viewId), SetViewPosition(viewId, x, y), SetViewSize(viewId, width, height), SetViewRect(viewId, x, y, width, height), CreateSplitscreenLeftRight(viewIdLeft, viewIdRight), CreateSplitscreenTopBottom(viewIdTop, viewIdBottom), CreateSplitscreenFour(viewIdTL, viewIdTR, viewIdBL, viewIdBR). Use for split-screen or picture-in-picture. See [Windows, scaling, and splitscreen](docs/WINDOWS_AND_VIEWS.md).

**3D editor and level builder (raylib_editor.go):** **Picking:** GetMouseRay() (updates internal ray from mouse + current 3D camera), GetMouseRayOriginX/Y/Z(), GetMouseRayDirectionX/Y/Z(). Use with GetRayCollisionSphere/Box/Mesh for selection. **Plane pick:** GetRayCollisionPlane(rayPosX,Y,Z, rayDirX,Y,Z, planeX,Y,Z, planeNormX,Y,Z) → 1 if hit, 0 otherwise; then GetRayCollisionPointX/Y/Z() for hit point. PickGroundPlane() → 1 if hit on y=0 plane, 0 otherwise (hit point via GetRayCollisionPointX/Y/Z). **Snap:** SnapToGridX(x, gridSize), SnapToGridY(y, gridSize), SnapToGridZ(z, gridSize). **Level objects:** CreateLevelObject(id, modelId, x, y, z, rotX, rotY, rotZ, scaleX, scaleY, scaleZ), **SetObjectPosition**(id, x, y, z), **RotateObject**(id, pitch, yaw, roll), **ScaleObject**(id, sx, sy, sz), **DrawObject**(id) (alias for DrawLevelObject), SetLevelObjectTransform(id, x, y, z, rotX, rotY, rotZ, scaleX, scaleY, scaleZ), GetLevelObjectX/Y/Z(id), GetLevelObjectRotX/RotY/RotZ(id), GetLevelObjectScaleX/ScaleY/ScaleZ(id), GetLevelObjectModelId(id), DeleteLevelObject(id), GetLevelObjectCount() → count, GetLevelObjectId(index) → id, DrawLevelObject(id), SaveLevel(path), LoadLevel(path), DuplicateLevelObject(id) → newId. **Camera readback:** GetCameraPositionX/Y/Z(), GetCameraTargetX/Y/Z(). For selection highlight use DrawModelWires or DrawBoundingBox(GetModelBoundingBox(...)); DrawGrid(slices, spacing) for editor grid.

---

## 8. Raylib (mesh) – `raylib_mesh.go`

GenMeshPoly, GenMeshPlane, GenMeshCube, GenMeshSphere, GenMeshHemiSphere, GenMeshCylinder, GenMeshCone, GenMeshTorus, GenMeshKnot, GenMeshHeightmap, GenMeshCubicmap, UploadMesh, UnloadMesh, GetMeshBoundingBox, ExportMesh, DrawMesh, UpdateMeshBuffer, DrawMeshInstanced, LoadMaterialDefault, IsMaterialValid, UnloadMaterial, SetMaterialTexture, LoadMaterials, GetMaterialIdFromLoad, GetRayCollisionMesh

---

## 9. Raylib (audio) – `raylib_audio.go`

InitAudioDevice, CloseAudioDevice, IsAudioDeviceReady, LoadSound, PlaySound, StopSound, SetSoundVolume, UnloadSound, LoadMusicStream, PlayMusicStream, UpdateMusicStream, StopMusicStream, SetMusicVolume, UnloadMusicStream, SetMasterVolume, GetMasterVolume, PauseSound, ResumeSound, IsSoundPlaying, SetSoundPitch, SetSoundPan, PauseMusicStream, ResumeMusicStream, IsMusicStreamPlaying, SeekMusicStream, SetMusicPitch, SetMusicPan, GetMusicTimeLength, GetMusicTimePlayed, IsMusicValid, LoadMusicStreamFromMemory, LoadWave, LoadWaveFromMemory, IsWaveValid, UnloadWave, ExportWave, WaveCopy, WaveCrop, WaveFormat, LoadWaveSamples, UnloadWaveSamples, **ExportWaveAsCode**, LoadSoundFromWave, LoadSoundAlias, IsSoundValid, UpdateSound, UnloadSoundAlias, LoadAudioStream, IsAudioStreamValid, UnloadAudioStream, UpdateAudioStream, IsAudioStreamProcessed, PlayAudioStream, PauseAudioStream, ResumeAudioStream, IsAudioStreamPlaying, StopAudioStream, SetAudioStreamVolume, SetAudioStreamPitch, SetAudioStreamPan, SetAudioStreamBufferSizeDefault  
**Not supported from BASIC (return error):** SetAudioStreamCallback, AttachAudioStreamProcessor, DetachAudioStreamProcessor, AttachAudioMixedProcessor, DetachAudioMixedProcessor

---

## 10. Raylib (fonts) – `raylib_fonts.go`

GetFontDefault, LoadFont, LoadFontEx, DrawTextExFont, MeasureTextEx, UnloadFont, LoadFontFromImage, LoadFontFromMemory, IsFontValid, LoadFontData, GenImageFontAtlas, UnloadFontData, **ExportFontAsCode**, DrawTextCodepoint, DrawTextCodepoints, GetGlyphIndex, GetGlyphInfo, GetGlyphAtlasRec

---

## 11. Raylib (misc) – `raylib_misc.go`

GetMouseDelta, NewColor, **Color(r, g, b, a)** (same as NewColor), CheckCollisionRecs, CheckCollisionCircles, CheckCollisionCircleRec, CheckCollisionPointRec, CheckCollisionPointCircle, GetCollisionRec, CheckCollisionSpheres, CheckCollisionBoxes, CheckCollisionBoxSphere, GetRayCollisionSphere, GetRayCollisionBox, GetRayCollisionTriangle, GetRayCollisionQuad, GetRayCollisionPointX, GetRayCollisionPointY, GetRayCollisionPointZ, GetRayCollisionNormalX, GetRayCollisionNormalY, GetRayCollisionNormalZ, GetRayCollisionDistance, Fade, ColorAlpha, ColorToInt, GetColor  
**Color constants (0-arg):** White, Black, LightGray, Gray, DarkGray, Yellow, Gold, Orange, Pink, Red, Maroon, Green, Lime, DarkGreen, SkyBlue, Blue, DarkBlue, Purple, Violet, DarkPurple, Beige, Brown, DarkBrown, Magenta, RayWhite, Blank  
ColorIsEqual, ColorNormalize, ColorFromNormalized, ColorToHSV, ColorFromHSV, ColorTint, ColorBrightness, ColorContrast, ColorAlphaBlend, ColorLerp, GetPixelDataSize

---

## 12. Raylib (math) – `raylib_math.go`

Clamp, Lerp, Normalize, Remap, Wrap, FloatEquals  
Vector2Zero, Vector2One, Vector2Add, Vector2AddValue, Vector2Subtract, Vector2SubtractValue, Vector2Length, Vector2LengthSqr, Vector2DotProduct, Vector2Distance, Vector2DistanceSqr, Vector2Angle, Vector2Scale, Vector2Multiply, Vector2Negate, Vector2Divide, Vector2Normalize, Vector2Transform, Vector2Lerp, Vector2Reflect, Vector2Rotate, Vector2MoveTowards, Vector2Invert, Vector2Clamp, Vector2ClampValue, Vector2Equals  
**VECTOR3(x, y, z)** → [x, y, z] (use with DrawModel etc.). Vector3Zero, Vector3One, Vector3Add, Vector3AddValue, Vector3Subtract, Vector3SubtractValue, Vector3Scale, Vector3Multiply, Vector3CrossProduct, Vector3Perpendicular, Vector3Length, Vector3LengthSqr, Vector3DotProduct, Vector3Distance, Vector3DistanceSqr, Vector3Angle, Vector3Negate, Vector3Divide, Vector3Normalize, Vector3OrthoNormalize, Vector3Transform, Vector3RotateByQuaternion, Vector3RotateByAxisAngle, Vector3Lerp, Vector3Reflect, Vector3Min, Vector3Max, Vector3Barycenter, Vector3Unproject, Vector3Invert, Vector3Clamp, Vector3ClampValue, Vector3Equals, Vector3Refract, Vector3ToFloatV  
MatrixDeterminant, MatrixTrace, MatrixTranspose, MatrixInvert, MatrixIdentity, MatrixAdd, MatrixSubtract, MatrixMultiply, MatrixTranslate, MatrixRotate, MatrixRotateX, MatrixRotateY, MatrixRotateZ, MatrixRotateXYZ, MatrixRotateZYX, MatrixScale, MatrixFrustum, MatrixPerspective, MatrixOrtho, MatrixLookAt, MatrixToFloatV  
QuaternionAdd, QuaternionAddValue, QuaternionSubtract, QuaternionSubtractValue, QuaternionIdentity, QuaternionLength, QuaternionNormalize, QuaternionInvert, QuaternionMultiply, QuaternionScale, QuaternionDivide, QuaternionLerp, QuaternionNlerp, QuaternionSlerp, QuaternionFromVector3ToVector3, QuaternionFromMatrix, QuaternionToMatrix, QuaternionFromAxisAngle, QuaternionToAxisAngle, QuaternionFromEuler, QuaternionToEuler, QuaternionTransform, QuaternionEquals

---

## 13. Raylib (game) – `raylib_game.go`

GAME.CameraOrbit, GAME.MoveWASD, GAME.OnGround, GAME.SnapToGround, MoveWASD3D, SnapToGround3D, IsOnGround3D, CameraOrbit3D, MoveHorizontal2D, Jump2D, IsOnGround2D, ClampVelocity2D, MoveVertical2D, CameraFollow2D, SnapToPlatform2D, Jump3D, ClampVelocity3D, CameraFollow3D  
**Input axes (0-arg):** GAME.GetAxisX, GAME.GetAxisY, GetAxisX, GetAxisY – return -1, 0, or 1 for A/D and W/S.  
**2D sprite–physics:** GAME.SyncSpriteToBody2D(worldId, bodyId, spriteId) – set sprite position to Box2D body (world→screen via camera). Call in draw loop.  
**Camera presets:** GAME.SetCamera2DFollow(worldId, bodyId, xOffset, yOffset) then GAME.UpdateCamera2D() each frame. GAME.SetCamera3DOrbit(worldId, bodyId, distance, heightOffset) then GAME.UpdateCamera3D(angleRad, pitchRad) each frame.  
**Collision callbacks:** GAME.SetCollisionHandler(bodyId, subName) – when bodyId collides, call Sub subName(otherBodyId). GAME.ProcessCollisions2D(worldId) – invoke handlers for this frame; call after BOX2D.Step.  
**Quality of life:** GAME.AssetPath(filename) → "assets/" + filename (for LoadTexture(AssetPath("hero.png"))). GAME.ClampDelta(maxDt) → min(GetFrameTime(), maxDt). GAME.ShowDebug() draws FPS; ShowDebug(extraText) draws FPS and a second line.  
**Key constants (0-arg):** GAME.KEY_W, GAME.KEY_A, GAME.KEY_S, GAME.KEY_D, GAME.KEY_SPACE  

**2D/3D helpers (convenience, no namespace):** **GetScreenCenterX()**, **GetScreenCenterY()** → screen center (width/2, height/2). **Distance2D(x1, y1, x2, y2)** → distance (wrapper around Vector2Distance). **Distance3D(x1, y1, z1, x2, y2, z2)** → distance (wrapper around Vector3Distance). **SetCamera2DCenter(worldX, worldY)** – set 2D camera so (worldX, worldY) is at screen center (offset = half width/height, rotation 0, zoom 1); useful when not using Box2D. These are conveniences for full 2D/3D games.

---

## 14. Box2D – `box2d.go`

**BOX2D.* (prefixed):** BOX2D.CreateWorld, BOX2D.Step, BOX2D.DestroyWorld, BOX2D.CreateBody, BOX2D.DestroyBody, BOX2D.GetBodyCount, BOX2D.GetBodyId, BOX2D.CreateBodyAtScreen, BOX2D.GetPosition, BOX2D.GetPositionX, BOX2D.GetPositionY, BOX2D.GetAngle, BOX2D.SetLinearVelocity, BOX2D.GetLinearVelocity, BOX2D.SetTransform, BOX2D.ApplyForce  
**Legacy (no prefix):** CreateWorld2D, DestroyWorld2D, Step2D, CreateBox2D, CreateCircle2D, CreatePolygon2D, CreateEdge2D, CreateChain2D, SetSensor2D, GetPositionX2D, GetPositionY2D, SetPosition2D, GetAngle2D, SetAngle2D, GetVelocityX2D, GetVelocityY2D, SetVelocity2D, ApplyForce2D, ApplyImpulse2D, ApplyTorque2D, SetAngularVelocity2D, GetAngularVelocity2D, SetFriction2D, SetRestitution2D, SetDamping2D, SetFixedRotation2D, SetGravityScale2D, SetMass2D, SetBullet2D  
**Joints:** CreateDistanceJoint2D(worldId, bodyAId, bodyBId, length) – implemented; rest (Revolute, Prismatic, Pulley, Gear, Weld, Rope, Wheel, SetJointLimits2D, SetJointMotor2D) are stubbed.  
RayCast2D, RayHitX2D, RayHitY2D, RayHitBody2D, RayHitNormalX2D, RayHitNormalY2D, GetCollisionCount2D, GetCollisionOther2D, GetCollisionNormalX2D, GetCollisionNormalY2D

---

## 15. Bullet (3D physics) – `bullet.go`

**BULLET.* (prefixed):** BULLET.CreateWorld, BULLET.SetGravity, BULLET.Step, BULLET.DestroyWorld, BULLET.CreateBox, BULLET.CreateSphere, BULLET.DestroyBody, BULLET.SetPosition, BULLET.GetPositionX, BULLET.GetPositionY, BULLET.GetPositionZ, BULLET.SetVelocity, BULLET.GetVelocityX, BULLET.GetVelocityY, BULLET.GetVelocityZ, BULLET.GetRotationX, BULLET.GetRotationY, BULLET.GetRotationZ, BULLET.SetRotation, BULLET.ApplyForce, BULLET.ApplyCentralForce, BULLET.ApplyImpulse, BULLET.RayCast, BULLET.GetRayCastHitX, BULLET.GetRayCastHitY, BULLET.GetRayCastHitZ, BULLET.GetRayCastHitBody, BULLET.GetRayCastHitNormalX, BULLET.GetRayCastHitNormalY, BULLET.GetRayCastHitNormalZ  
**Legacy (no prefix):** CreateWorld3D, DestroyWorld3D, Step3D, CreateSphere3D, CreateBox3D, CreateCapsule3D, CreateStaticMesh3D, CreateCylinder3D, CreateCone3D, CreateHeightmap3D, CreateCompound3D, AddShapeToCompound3D, GetPositionX3D, GetPositionY3D, GetPositionZ3D, SetPosition3D, GetYaw3D, GetPitch3D, GetRoll3D, SetRotation3D, SetScale3D, GetVelocityX3D, GetVelocityY3D, GetVelocityZ3D, SetVelocity3D, SetAngularVelocity3D, GetAngularVelocityX3D, GetAngularVelocityY3D, GetAngularVelocityZ3D, ApplyForce3D, ApplyImpulse3D, ApplyTorque3D, ApplyTorqueImpulse3D, SetMass3D  
**Stub (no-op):** SetFriction3D, SetRestitution3D, SetDamping3D, SetKinematic3D, SetGravity3D, SetLinearFactor3D, SetAngularFactor3D, SetCCD3D, CreateHingeJoint3D, CreateSliderJoint3D, CreateConeTwistJoint3D, CreatePointToPointJoint3D, CreateFixedJoint3D, SetJointLimits3D, SetJointMotor3D  
RayCast3D, RayHitX3D, RayHitY3D, RayHitZ3D, RayHitBody3D, RayHitNormalX3D, RayHitNormalY3D, RayHitNormalZ3D, GetCollisionCount3D, GetCollisionOther3D, GetCollisionNormalX3D, GetCollisionNormalY3D, GetCollisionNormalZ3D

---

## 16. ECS – `ecs.go`

**ECS.*** (prefixed): ECS.CreateWorld() → worldId, ECS.DestroyWorld(worldId), ECS.CreateEntity(worldId) → entityId, ECS.DestroyEntity(worldId, entityId), ECS.AddComponent(worldId, entityId, componentType [, args...]), ECS.HasComponent(worldId, entityId, componentType) → boolean, ECS.RemoveComponent(worldId, entityId, componentType). **Transform:** ECS.SetTransform(worldId, entityId, x, y, z), ECS.GetTransformX/Y/Z(worldId, entityId) → number. **Placement:** ECS.PlaceEntity(worldId, entityId, x, y, z) — same as SetTransform. **Scene graph:** ECS.GetWorldPositionX/Y/Z(worldId, entityId) → number (world position including Parent chain). Add component **Parent** with parent entity ID: ECS.AddComponent(worldId, entityId, "Parent", parentEntityId). **Health:** ECS.GetHealthCurrent/Max(worldId, entityId) → number. **Query:** ECS.QueryCount(worldId, componentType1 [, componentType2...]) → count, ECS.QueryEntity(worldId, componentType, index) → entityId or "". Component types: Transform(x,y,z), Sprite(textureId, visible), Health(current, max), Parent(parentEntityId). See [ECS_GUIDE.md](docs/ECS_GUIDE.md).

---

## 17. Std (file, string, math, JSON, Enum, Dictionary, HTTP, HELP, multi-window) – `std.go`

**File:** ReadFile(path) → string or nil on error; WriteFile(path, contents) → boolean; **LoadText(path)** → string (alias for ReadFile); **SaveText(path, text)** → boolean (alias for WriteFile); DeleteFile(path) → boolean; **CopyFile(src, dst)** → boolean; **ListDir(path)** → count (use GetDirItem(index) to read entries, 0-based); **GetDirItem(index)** → string; **ExecuteFile(path)** → boolean (start process). **IsNull(value)** → boolean (true when value is null). Use the **Nil** or **Null** literal for missing values. **FileExists(path)** is in raylib core.

**Enum (runtime; requires ENUM declared in script):** **Enum.getValue(enumName, valueName)** → int; **Enum.getName(enumName, value)** → string (name for that value, or ""); **Enum.hasValue(enumName, valueName)** → boolean. Enum names and value names are case-insensitive.

**Dictionary:** **CreateDict()** → new empty map; **SetDictKey(dict, key, value)** → dict (mutates and returns). Use **GetJSONKey(dict, key)** to read (dict can be a map from CreateDict or a dict literal). **Dictionary.has(dict, key)** → boolean; **Dictionary.keys(dict)** → array of keys; **Dictionary.values(dict)** → array of values; **Dictionary.size(dict)** → int; **Dictionary.remove(dict, key)** → dict; **Dictionary.clear(dict)** → dict; **Dictionary.merge(dict1, dict2)** → new merged dict; **Dictionary.get(dict, key [, default])** → value or default.

**Multi-window (same .bas, multiple processes):** **GetEnv(key)** → string (env var); **IsWindowProcess()** → true if run with --window; **GetWindowTitle()**, **GetWindowWidth()**, **GetWindowHeight()** → title/width/height for child window (from --title=, --width=, --height=); **SpawnWindow(port, title, width, height)** → 1 on success, 0 on failure (starts same .bas as child; main then AcceptTimeout to get connection). See [docs/MULTI_WINDOW.md](docs/MULTI_WINDOW.md).

**String (DBP-style):** **Left(s, n)** → first n characters; **Right(s, n)** → last n characters; **Mid(s, start1Based, [count])** → substring (start 1-based; if count omitted, rest of string); **Substr(s, start0Based, [count])** → substring (start 0-based); **Instr(s, sub)** → 1-based index of sub in s, or 0 if not found; **Upper(s)**, **Lower(s)** → string in upper/lower case; **Len(s)** → character count; **Chr(code)** → string of one character; **Asc(s)** → code of first character; **Str(x)** → number to string; **Val(s)** → string to float.

**Math (DBP-style):** **Rnd()** → float in [0,1); **Rnd(n)** → integer in 1..n inclusive; **Random(n)** → integer in 0..n-1; **Random(min, max)** → integer in [min, max] inclusive; **Int(x)** → truncate to integer. **Radians(degrees)** → radians; **Degrees(radians)** → degrees; **AngleWrap(angle)** / **WrapAngle(angle)** → angle in radians wrapped to [-π, π]. **TimeNow()** → seconds since epoch (float). **PrintDebug(value)** → print value to stderr for debugging. **Raylib math:** Clamp(value, min, max), Lerp(a, b, t), **Vec2(x, y)**, **Vec3(x, y, z)** (alias for VECTOR3), **Color(r, g, b, a)** (raylib misc).

**Assert:** **Assert(condition, [message])** — if condition is falsy, stops execution with error (message or "assertion failed"). Use statement: `ASSERT x > 0, "x must be positive"`.

**JSON:** LoadJSON(path) → handle (string); LoadJSONFromString(str) → handle; GetJSONKey(handle, key) → value (string/number/boolean); SaveJSON(path, handle) → boolean.

**HTTP:** HttpGet(url) → body string or nil; HttpPost(url, body) → response string; DownloadFile(url, path) → boolean.

**Help:** HELP(), ?() – print quick reference and path to this API; return nil.

---

## 18. Multiplayer (TCP) – `net.go`

**Client:** Connect(host, port) → connectionId or null; **ConnectToParent()** → connectionId or null (connects using CYBERBASIC_PARENT env, for spawned window processes); **ConnectTLS**(host, port) → connectionId or null (encrypted; use for internet). Send(connectionId, text) → boolean (max 256 KB, no newlines); **SendText**(connectionId, text) — same as Send. **SendJSON**(connectionId, jsonText) → 1 if sent, 0 if invalid; **SendInt**(connectionId, value) → 1/0; **SendFloat**(connectionId, value) → 1/0; **SendNumbers**(connectionId, n1, n2, …) → 1/0 (up to 16 numbers). Receive(connectionId) → next line or null; **ReceiveJSON**(connectionId) → next line if valid JSON, else null; **ReceiveNumbers**(connectionId) → count of numbers parsed (0 if no data or parse error); **GetReceivedNumber**(index) → number at index from last ReceiveNumbers (0.0 if out of range). Disconnect(connectionId).

**Server:** Host(port) → serverId or null; **HostTLS**(port, certFile, keyFile) → serverId or null (encrypted). Accept(serverId) → connectionId when a client connects (blocking); CloseServer(serverId). Messages are line-based (one per Send/Receive). Use plain Connect/Host for LAN or dev; use ConnectTLS/HostTLS for encrypted channels (e.g. internet).

**Rooms:** CreateRoom(roomId) — ensure room exists (idempotent). JoinRoom(roomId, connectionId) — add connection to room (creates room if missing). LeaveRoom(connectionId) — remove connection from all rooms; LeaveRoom(connectionId, roomId) — remove from that room only. SendToRoom(roomId, text) — send text to every connection (max 256 KB, no newlines); returns number sent. **SendToRoomJSON**(roomId, jsonText) — send valid JSON to room; returns number sent (0 if invalid). **SendToRoomInt**(roomId, value), **SendToRoomFloat**(roomId, value), **SendToRoomNumbers**(roomId, n1, n2, …) — broadcast typed numbers to room; each returns count sent. GetRoomConnectionCount(roomId) — number of connections in room (0 if missing). GetRoomConnectionId(roomId, index) — connectionId at 0-based index, or empty string if out of range.

**Convenience:** IsConnected(connectionId) — 1 if connection is in conns, 0 otherwise. GetConnectionCount() — total number of connections. AcceptTimeout(serverId, timeoutMs) — like Accept but with a timeout; returns connectionId or null. GetLocalIP() — this machine’s local IP (e.g. 192.168.1.x) for LAN; show "Connect to: GetLocalIP() : port".

See [docs/MULTIPLAYER.md](docs/MULTIPLAYER.md).

---

## 19. SQL – `sql.go`

**Connection:** OpenDatabase(path) → dbId (string) or null (SQLite file at path; e.g. `"game.db"`). CloseDatabase(dbId).

**Execute:** Exec(dbId, sql) — run INSERT/UPDATE/DELETE/DDL; returns rows affected (int), or -1 on error. ExecParams(dbId, sql, arg1, arg2, …) — same with `?` placeholders; args are numbers or strings.

**Query:** Query(dbId, sql) — run SELECT; stores result set internally; returns row count (int), or -1 on error. QueryParams(dbId, sql, arg1, arg2, …) — parameterized SELECT.

**Result access:** GetRowCount() — rows in last query result. GetColumnCount() — columns. GetColumnName(colIndex) — name (0-based). GetCell(row, col) — value at (row, col) as number or string; returns null for SQL NULL (use IsNull(GetCell(r,c))). Row/col are 0-based.

**Transactions:** Begin(dbId), Commit(dbId), Rollback(dbId) — return 1 on success, 0 on error.

**Errors:** LastError() — returns last error message (string); clear after successful calls.

See [docs/SQL.md](docs/SQL.md).

---

## 20. UI – `raylib_ui.go` and full raygui – `raylib_raygui.go`

**Layout (pure-Go, no CGO):** BeginUI() resets cursor; widgets advance a vertical cursor. EndUI().

**Widgets (pure-Go):** Label(text); Button(text) → boolean (clicked); Slider(text, value, min, max) → value; Checkbox(text, checked) → 1 or 0; TextBox(id, text) → text (editable, use same id each frame); Dropdown(id, itemsText, activeIndex) → activeIndex (itemsText = "A;B;C"); ProgressBar(text, value, min, max) → value; WindowBox(title) / EndWindowBox(); GroupBox(text) / EndGroupBox(). Call inside your game loop (e.g. after ClearBackground).

**Full raygui (gen2brain/raylib-go/raygui; requires CGO):** GuiLabel(x, y, w, h, text); GuiButton(x, y, w, h, text) → 1 if clicked else 0; GuiCheckBox(x, y, w, h, text, checked) → 1 or 0; GuiSlider(x, y, w, h, textLeft, textRight, value, min, max) → value; GuiProgressBar(x, y, w, h, textLeft, textRight, value, min, max) → value; GuiTextBox(id, x, y, w, h, text) → currentText (id is cache key); GuiDropdownBox(id, x, y, w, h, itemsText, active) → newActive (itemsText e.g. "One;Two;Three"); GuiWindowBox(x, y, w, h, title) → 1 if close clicked else 0; GuiGroupBox(x, y, w, h, text); GuiLine(x, y, w, h, text); GuiPanel(x, y, w, h, text). All coordinates and sizes in pixels. See [docs/GUI_GUIDE.md](docs/GUI_GUIDE.md).

---

## 21. Language and built-ins

**User code:** Function/Sub with parameters and Return; call by name. **Modules:** Module Name … End Module; body is Function/Sub only; call as ModuleName.FunctionName(...).

**Single-line IF:** Consecutive `IF condition THEN statement` lines (e.g. `IF IsKeyDown(KEY_W) THEN movePlayer(0, -speed)` on each line) do not require ENDIF between them. Use ENDIF when the next line is not another IF.

**Events:** On KeyDown("KEY") … End On, On KeyPressed("KEY") … End On. Register a handler for that key; handlers run when PollInputEvents() is called (e.g. in the game loop). Key names: "ESCAPE", "W", "SPACE", etc., or KEY_* constant values.

**Coroutines:** StartCoroutine SubName() – start a fiber at that sub; Yield – switch to next fiber; WaitSeconds(seconds) – block current fiber for N seconds (blocks entire VM). Fibers share the same chunk; each has its own IP, stack, and call stack.

---

## 22. Multi-window (in-process) – `raylib_multiwindow.go`

Logical windows (viewports) in one process; ID 0 = main screen. See [docs/MULTI_WINDOW_INPROCESS.md](docs/MULTI_WINDOW_INPROCESS.md).

**Creation/lifecycle:** WindowCreate(width, height, title) → id; WindowCreatePopup, WindowCreateModal, WindowCreateToolWindow(width, height, title); WindowCreateChild(parentID, width, height, title); WindowClose(id); WindowIsOpen(id); WindowSetTitle(id, title); WindowSetSize(id, width, height); WindowSetPosition(id, x, y); WindowGetWidth(id), WindowGetHeight(id); WindowGetPositionX(id), WindowGetPositionY(id), WindowGetPosition(id); WindowFocus(id); WindowIsFocused(id); WindowIsVisible(id); WindowShow(id), WindowHide(id).

**Rendering:** WindowBeginDrawing(id), WindowEndDrawing(id), WindowClearBackground(id, r, g, b, a), WindowDrawAllToScreen().

**Messages:** WindowSendMessage(targetID, message, data), WindowBroadcast(message, data), WindowReceiveMessage(id) → "message|data" or null, WindowHasMessage(id).

**Channels:** ChannelCreate(name), ChannelSend(name, data), ChannelReceive(name) → value or null, ChannelHasData(name).

**State:** StateSet(key, value), StateGet(key), StateHas(key), StateRemove(key).

**Events:** OnWindowUpdate(id, subName), OnWindowDraw(id, subName), OnWindowResize(id, subName), OnWindowClose(id, subName), OnWindowMessage(id, subName); WindowProcessEvents(), WindowDraw().

**3D/RPC:** WindowSetCamera(id, cameraId), WindowDrawModel(id, modelId, x, y, z, scale [, r,g,b,a]), WindowDrawScene(id, sceneId); WindowRegisterFunction(windowId, name, subName), WindowCall(targetWindowId, name, arg1, arg2, …).

**Docking:** DockCreateArea(id, x, y, width, height), DockSplit(areaId, direction, size), DockAttachWindow(areaId, windowId), DockSetSize(nodeId, size).

---

## Notes

- **Resource IDs:** Load* functions (LoadImage, LoadTexture, LoadSound, LoadMusicStream, LoadWave, LoadFont, LoadModel, LoadMesh, LoadShader, LoadRenderTexture, LoadAudioStream, etc.) return string IDs (e.g. `img_1`, `sound_1`). Pass these IDs to the matching Unload* and other APIs.
- **Vectors/Matrix/Quaternion:** Pass as flat numbers: Vector2 (x,y), Vector3 (x,y,z), Matrix (16 floats row-major), Quaternion (x,y,z,w). Functions that return vectors/matrices return a list (e.g. [x,y] or 16 values).
- **Colors:** Pass as (r,g,b,a) or use constants (White, Red, etc. – return packed int). NewColor(r,g,b,a) returns packed int.
- **Export*AsCode:** ExportImageAsCode(imageId, fileName), ExportFontAsCode(fontId, fileName), ExportWaveAsCode(waveId, fileName) write C header files; return true on success.
- **Audio callbacks:** SetAudioStreamCallback, AttachAudioStreamProcessor, AttachAudioMixedProcessor (and Detach*) are registered but return an error from BASIC; use UpdateAudioStream to push samples instead.
