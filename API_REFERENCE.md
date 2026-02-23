# CyberBasic API Reference – All Bindings

All functions callable from BASIC. Names are **case-insensitive**. You can call with or without namespace (e.g. `InitWindow(...)` or `RL.InitWindow(...)` for raylib; use `BOX2D.*` and `BULLET.*` for physics).

---

## 1. Raylib (core) – `raylib_core.go`

InitWindow, SetTargetFPS, WindowShouldClose, CloseWindow, SetWindowPosition, ClearBackground, GetFrameTime, **DeltaTime** (same as GetFrameTime; preferred for frame delta), GetFPS, GetScreenWidth, GetScreenHeight, SetWindowSize, SetWindowTitle, MaximizeWindow, MinimizeWindow, IsWindowReady, IsWindowFullscreen, GetTime, GetRandomValue, SetRandomSeed, SetWindowState, ClearWindowState, GetMonitorCount, GetCurrentMonitor, GetClipboardText, SetClipboardText, TakeScreenshot, OpenURL, IsWindowHidden, IsWindowMinimized, IsWindowMaximized, IsWindowFocused, IsWindowResized, ToggleFullscreen, RestoreWindow, GetRenderWidth, GetRenderHeight, GetMonitorName, GetMonitorWidth, GetMonitorHeight, GetMonitorRefreshRate, WaitTime, EnableEventWaiting, DisableEventWaiting, IsCursorHidden, EnableCursor, DisableCursor, IsCursorOnScreen, IsWindowState, ToggleBorderlessWindowed, SetWindowMonitor, SetWindowMinSize, SetWindowMaxSize, SetWindowOpacity, GetWindowPosition, GetWindowScaleDPI, GetMonitorPosition, GetMonitorPhysicalWidth, GetMonitorPhysicalHeight, SetConfigFlags, SwapScreenBuffer, PollInputEvents, SetCamera2D, BeginMode2D, EndMode2D, GetWorldToScreen2D, GetScreenToWorld2D, BeginBlendMode, EndBlendMode, BeginScissorMode, EndScissorMode, BeginShaderMode, EndShaderMode, LoadShader, LoadShaderFromMemory, UnloadShader, IsShaderValid, **FileExists**

---

## 2. Raylib (input) – `raylib_input.go`

IsMouseButtonPressed, GetMouseX, GetMouseY, IsKeyPressed, IsKeyDown, IsKeyReleased, IsKeyUp, GetKeyPressed, SetExitKey, IsMouseButtonDown, IsMouseButtonReleased, GetMouseWheelMove, SetMousePosition, SetMouseOffset, SetMouseScale, HideCursor, ShowCursor, GetMousePosition, GetVector2X, GetVector2Y, GetVector3Z, IsMouseButtonUp, IsKeyPressedRepeat, GetCharPressed, SetMouseCursor, IsGamepadAvailable, GetGamepadName, IsGamepadButtonPressed, IsGamepadButtonDown, IsGamepadButtonReleased, GetGamepadAxisMovement, IsGamepadButtonUp, GetGamepadButtonPressed, GetGamepadAxisCount, SetGamepadMappings, SetGamepadVibration, GetTouchPointCount, GetTouchX, GetTouchY, GetTouchPosition, GetTouchPointId, GetMouseWheelMoveV  
**Key constants (0-arg):** KEY_NULL, KEY_APOSTROPHE, KEY_COMMA, KEY_MINUS, KEY_PERIOD, KEY_SLASH, KEY_ZERO … KEY_NINE, KEY_SEMICOLON, KEY_EQUAL, KEY_A … KEY_Z, KEY_LEFT_BRACKET, KEY_BACKSLASH, KEY_RIGHT_BRACKET, KEY_GRAVE, KEY_SPACE, KEY_ESCAPE, KEY_ENTER, KEY_TAB, KEY_BACKSPACE, KEY_INSERT, KEY_DELETE, KEY_RIGHT, KEY_LEFT, KEY_DOWN, KEY_UP, KEY_PAGE_UP, KEY_PAGE_DOWN, KEY_HOME, KEY_END, KEY_F1 … KEY_F12  
**Movement:** For simple movement use **GetAxisX()** / **GetAxisY()** (return -1, 0, or 1 for A/D and W/S), e.g. `x = x + speed * GetAxisX()`. For full 2D/3D use **GAME.MoveWASD**, **MoveHorizontal2D**, **Jump2D** (see §13 Raylib (game)).

---

## 3. Raylib (shapes) – `raylib_shapes.go`

SetShapesTexture, GetShapesTextureRectangle, DrawRectangle, DrawCircle, DrawLine, DrawLineV, DrawCircleLines, DrawRectangleLines, DrawTriangle, DrawTriangleLines, DrawPixel, DrawPoly, DrawEllipse, DrawRing, DrawRectangleRounded, DrawGrid, DrawFPS, DrawLineEx, DrawPixelV, DrawCircleSector, DrawCircleGradient, DrawCircleV, DrawEllipseLines, DrawRingLines, DrawRectangleV, DrawRectangleRec, DrawRectanglePro, DrawRectangleLinesEx, DrawRectangleRoundedLines, DrawPolyLines

---

## 4. Raylib (text) – `raylib_text.go`

DrawText, MeasureText, DrawTextEx, DrawTextPro, SetTextLineSpacing, TextCopy, TextIsEqual, TextLength, TextFormat, TextSubtext, TextReplace, TextInsert, TextJoin, TextSplit, GetTextSplitItem, TextAppend, TextFindIndex, TextToUpper, TextToLower, TextToPascal, TextToSnake, TextToCamel, TextToInteger, TextToFloat, GetCodepointCount, GetCodepoint, GetCodepointNext, GetCodepointPrevious, CodepointToUTF8, LoadCodepoints, UnloadCodepoints, GetLoadedCodepoint, LoadUTF8, UnloadUTF8

---

## 5. Raylib (textures) – `raylib_textures.go`

LoadTexture, UnloadTexture, LoadRenderTexture, UnloadRenderTexture, BeginTextureMode, EndTextureMode, DrawTexture, DrawTextureEx, DrawTextureRec, DrawTexturePro, LoadTextureFromImage, LoadTextureCubemap, IsTextureValid, IsRenderTextureValid, UpdateTexture, UpdateTextureRec, GenTextureMipmaps, SetTextureFilter, SetTextureWrap, DrawTextureV, DrawTextureNPatch

---

## 6. Raylib (images) – `raylib_images.go`

LoadImage, LoadImageRaw, LoadImageAnim, GetLoadImageAnimFrames, LoadImageAnimFromMemory, LoadImageFromMemory, LoadImageFromTexture, LoadImageFromScreen, IsImageValid, UnloadImage, ExportImage, ExportImageToMemory, **ExportImageAsCode**, GenImageColor, GenImageGradientLinear, GenImageGradientRadial, GenImageGradientSquare, GenImageChecked, GenImageWhiteNoise, GenImagePerlinNoise, GenImageCellular, GenImageText, ImageCopy, ImageFromImage, ImageFromChannel, ImageText, ImageTextEx, ImageFormat, ImageToPOT, ImageCrop, ImageAlphaCrop, ImageAlphaClear, ImageAlphaMask, ImageAlphaPremultiply, ImageBlurGaussian, ImageKernelConvolution, ImageResize, ImageResizeNN, ImageResizeCanvas, ImageMipmaps, ImageDither, ImageFlipVertical, ImageFlipHorizontal, ImageRotate, ImageRotateCW, ImageRotateCCW, ImageColorTint, ImageColorInvert, ImageColorGrayscale, ImageColorContrast, ImageColorBrightness, ImageColorReplace, LoadImageColors, UnloadImageColors, GetLoadedImageColor, GetImageColor, ImageClearBackground, ImageDrawPixel, ImageDrawPixelV, ImageDrawLine, ImageDrawLineV, ImageDrawLineEx, ImageDrawCircle, ImageDrawCircleV, ImageDrawCircleLines, ImageDrawCircleLinesV, ImageDrawRectangle, ImageDrawRectangleV, ImageDrawRectangleRec, ImageDrawRectangleLines, ImageDrawTriangle, ImageDrawTriangleEx, ImageDrawTriangleLines, ImageDrawTriangleFan, ImageDrawTriangleStrip, ImageDraw

---

## 7. Raylib (3D) – `raylib_3d.go`

SetCamera3D, BeginMode3D, EndMode3D, LoadModel, LoadModelFromMesh, UnloadModel, DrawModel, DrawCube, DrawCubeWires, DrawSphere, DrawSphereWires, DrawPlane, DrawLine3D, DrawPoint3D, DrawCircle3D, DrawCubeV, DrawCylinder, DrawCylinderWires, DrawRay, DrawTriangle3D, DrawTriangleStrip3D, DrawCubeWiresV, DrawSphereEx, DrawCylinderEx, DrawCylinderWiresEx, DrawCapsule, DrawCapsuleWires, DrawModelEx, DrawModelWires, DrawBoundingBox, IsModelValid, GetModelBoundingBox, DrawModelWiresEx, DrawModelPoints, DrawModelPointsEx, DrawBillboard, DrawBillboardRec, SetModelMeshMaterial, DrawBillboardPro, LoadModelAnimations, GetModelAnimationId, UpdateModelAnimation, UpdateModelAnimationBones, UnloadModelAnimation, UnloadModelAnimations, IsModelAnimationValid

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

GetMouseDelta, NewColor, CheckCollisionRecs, CheckCollisionCircles, CheckCollisionCircleRec, CheckCollisionPointRec, CheckCollisionPointCircle, GetCollisionRec, CheckCollisionSpheres, CheckCollisionBoxes, CheckCollisionBoxSphere, GetRayCollisionSphere, GetRayCollisionBox, GetRayCollisionTriangle, GetRayCollisionQuad, GetRayCollisionPointX, GetRayCollisionPointY, GetRayCollisionPointZ, GetRayCollisionNormalX, GetRayCollisionNormalY, GetRayCollisionNormalZ, GetRayCollisionDistance, Fade, ColorAlpha, ColorToInt, GetColor  
**Color constants (0-arg):** White, Black, LightGray, Gray, DarkGray, Yellow, Gold, Orange, Pink, Red, Maroon, Green, Lime, DarkGreen, SkyBlue, Blue, DarkBlue, Purple, Violet, DarkPurple, Beige, Brown, DarkBrown, Magenta, RayWhite, Blank  
ColorIsEqual, ColorNormalize, ColorFromNormalized, ColorToHSV, ColorFromHSV, ColorTint, ColorBrightness, ColorContrast, ColorAlphaBlend, ColorLerp, GetPixelDataSize

---

## 12. Raylib (math) – `raylib_math.go`

Clamp, Lerp, Normalize, Remap, Wrap, FloatEquals  
Vector2Zero, Vector2One, Vector2Add, Vector2AddValue, Vector2Subtract, Vector2SubtractValue, Vector2Length, Vector2LengthSqr, Vector2DotProduct, Vector2Distance, Vector2DistanceSqr, Vector2Angle, Vector2Scale, Vector2Multiply, Vector2Negate, Vector2Divide, Vector2Normalize, Vector2Transform, Vector2Lerp, Vector2Reflect, Vector2Rotate, Vector2MoveTowards, Vector2Invert, Vector2Clamp, Vector2ClampValue, Vector2Equals  
Vector3Zero, Vector3One, Vector3Add, Vector3AddValue, Vector3Subtract, Vector3SubtractValue, Vector3Scale, Vector3Multiply, Vector3CrossProduct, Vector3Perpendicular, Vector3Length, Vector3LengthSqr, Vector3DotProduct, Vector3Distance, Vector3DistanceSqr, Vector3Angle, Vector3Negate, Vector3Divide, Vector3Normalize, Vector3OrthoNormalize, Vector3Transform, Vector3RotateByQuaternion, Vector3RotateByAxisAngle, Vector3Lerp, Vector3Reflect, Vector3Min, Vector3Max, Vector3Barycenter, Vector3Unproject, Vector3Invert, Vector3Clamp, Vector3ClampValue, Vector3Equals, Vector3Refract, Vector3ToFloatV  
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

## 16. Std (file, JSON, HTTP, HELP) – `std.go`

**File:** ReadFile(path) → string or nil on error; WriteFile(path, contents) → boolean; DeleteFile(path) → boolean. **IsNull(value)** → boolean (true when value is null). Use the **Nil** or **Null** literal for missing values.

**JSON:** LoadJSON(path) → handle (string); LoadJSONFromString(str) → handle; GetJSONKey(handle, key) → value (string/number/boolean); SaveJSON(path, handle) → boolean.

**HTTP:** HttpGet(url) → body string or nil; HttpPost(url, body) → response string; DownloadFile(url, path) → boolean.

**Help:** HELP(), ?() – print quick reference and path to this API; return nil.

---

## 17. UI – `raylib_ui.go`

BeginUI(), EndUI(), Label(text), Button(text) → boolean. **Minimal immediate-mode UI:** BeginUI resets layout; Label(text) draws text and advances layout; Button(text) draws a button and returns true when clicked. Call inside your game loop (e.g. inside Main() after ClearBackground). Use for menus, pause screens, and HUD.

---

## 18. Language and built-ins

**User code:** Function/Sub with parameters and Return; call by name. **Modules:** Module Name … End Module; body is Function/Sub only; call as ModuleName.FunctionName(...).

**Events:** On KeyDown("KEY") … End On, On KeyPressed("KEY") … End On. Register a handler for that key; handlers run when PollInputEvents() is called (e.g. in the game loop). Key names: "ESCAPE", "W", "SPACE", etc., or KEY_* constant values.

**Coroutines:** StartCoroutine SubName() – start a fiber at that sub; Yield – switch to next fiber; WaitSeconds(seconds) – block current fiber for N seconds (blocks entire VM). Fibers share the same chunk; each has its own IP, stack, and call stack.

---

## Notes

- **Resource IDs:** Load* functions (LoadImage, LoadTexture, LoadSound, LoadMusicStream, LoadWave, LoadFont, LoadModel, LoadMesh, LoadShader, LoadRenderTexture, LoadAudioStream, etc.) return string IDs (e.g. `img_1`, `sound_1`). Pass these IDs to the matching Unload* and other APIs.
- **Vectors/Matrix/Quaternion:** Pass as flat numbers: Vector2 (x,y), Vector3 (x,y,z), Matrix (16 floats row-major), Quaternion (x,y,z,w). Functions that return vectors/matrices return a list (e.g. [x,y] or 16 values).
- **Colors:** Pass as (r,g,b,a) or use constants (White, Red, etc. – return packed int). NewColor(r,g,b,a) returns packed int.
- **Export*AsCode:** ExportImageAsCode(imageId, fileName), ExportFontAsCode(fontId, fileName), ExportWaveAsCode(waveId, fileName) write C header files; return true on success.
- **Audio callbacks:** SetAudioStreamCallback, AttachAudioStreamProcessor, AttachAudioMixedProcessor (and Detach*) are registered but return an error from BASIC; use UpdateAudioStream to push samples instead.
