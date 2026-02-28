# CyberBasic – Complete Command Reference

Structured command set for window, input, math, camera, 3D, 2D, audio, file, game loop, and utilities. All names are **case-insensitive**. For the complete list of all bindings and source files, see [API Reference](../API_REFERENCE.md).

---

## Window & system

| Command | Description |
|--------|-------------|
| **InitWindow**(width, height, title) | Open game window |
| **CloseWindow**() | Close window and exit |
| **SetTargetFPS**(fps) | Target frames per second |
| **GetFrameTime**() | Delta time since last frame (seconds) |
| **GetTime**() | Seconds since window init (float) |
| **WaitSeconds**(seconds) | Yield current fiber for N seconds (non-blocking; other fibers run) |
| **WindowShouldClose**() | True when user requested close |
| **DisableCursor**() | Hide and confine mouse |
| **EnableCursor**() | Show mouse cursor |

---

## Game loop (hybrid)

When you define **update(dt)** and **draw()** (Sub or Function) and use a game loop (`WHILE NOT WindowShouldClose()` or `REPEAT ... UNTIL WindowShouldClose()`), the compiler injects an automatic pipeline. You do not call BeginDrawing/EndDrawing yourself.

| Command | Description |
|--------|-------------|
| **ClearRenderQueues**() | Clear 2D, 3D, and GUI render queues (called automatically before draw()) |
| **FlushRenderQueues**() | Execute queued draw commands and present frame (called automatically after draw()) |
| **StepAllPhysics2D**(dt) | Step all registered Box2D worlds (called automatically with frame delta) |
| **StepAllPhysics3D**(dt) | Step all registered Bullet worlds (called automatically with frame delta) |

See [Program Structure](PROGRAM_STRUCTURE.md#hybrid-updatedraw-loop).

---

## Multi-window (in-process)

Logical windows (viewports) in one process. Window ID **0** = main screen. See [In-process multi-window](MULTI_WINDOW_INPROCESS.md) for the full guide.

| Command | Description |
|--------|-------------|
| **WindowCreate**(width, height, title) | Create window → id |
| **WindowClose**(id) | Close window |
| **WindowIsOpen**(id) | True if window exists |
| **WindowSetTitle**(id, title) | **WindowSetSize**(id, w, h) | **WindowSetPosition**(id, x, y) |
| **WindowGetWidth**(id) **WindowGetHeight**(id) | **WindowGetPositionX**(id) **WindowGetPositionY**(id) |
| **WindowBeginDrawing**(id) **WindowEndDrawing**(id) | Draw into window (0 = main) |
| **WindowClearBackground**(id, r, g, b, a) | **WindowDrawAllToScreen**() |
| **WindowSendMessage**(targetID, message, data) | **WindowReceiveMessage**(id) | **WindowHasMessage**(id) |
| **ChannelCreate**(name) **ChannelSend**(name, data) **ChannelReceive**(name) **ChannelHasData**(name) |
| **StateSet**(key, value) **StateGet**(key) **StateHas**(key) **StateRemove**(key) |

---

## Input

| Command | Description |
|--------|-------------|
| **KeyDown**(key) / **IsKeyDown**(key) | True while key held (use KEY_W, KEY_ESCAPE, etc.) |
| **KeyPressed**(key) / **IsKeyPressed**(key) | True once when key pressed |
| **IsKeyReleased**(key) | True once when key released |
| **GetKeyPressed**() | Last key pressed (code) |
| **GetMouseX**() **GetMouseY**() | Mouse position |
| **GetMousePosition**() | Mouse position as [x, y] |
| **IsMouseButtonDown**(button) | True while button held |
| **IsMouseButtonPressed**(button) | True once when button pressed |
| **GetMouseDeltaX**() **GetMouseDeltaY**() | Mouse movement this frame |
| **GetMouseWheelMove**() | Scroll wheel delta |
| **IsGamepadAvailable**(id) | True if gamepad connected |
| **IsGamepadButtonPressed**(id, button) | True once when gamepad button pressed |
| **GetGamepadAxisMovement**(id, axis) | Gamepad axis value |
| **MouseOrbitCamera**() | One call: orbit + zoom from mouse, then update camera |

---

## Math

| Command | Description |
|--------|-------------|
| **Clamp**(value, min, max) | Clamp value to [min, max] |
| **Lerp**(a, b, t) | Linear interpolate: a + (b−a)*t |
| **WrapAngle**(angle) | Wrap angle (radians) to [−π, π] |
| **Vec2**(x, y) | 2D vector [x, y] |
| **Vec3**(x, y, z) | 3D vector [x, y, z] (alias VECTOR3) |
| **Color**(r, g, b, a) | RGBA color (0–255) |

---

## Camera

| Command | Description |
|--------|-------------|
| **CameraOrbit**(cx, cy, cz, angle, pitch, distance) | Orbit camera around target; updates internal state |
| **CameraZoom**(amount) | Adjust orbit distance (e.g. GetMouseWheelMove()) |
| **CameraRotate**(dx, dy) | Rotate from mouse delta (2 args) |
| **SetCameraPosition**(x, y, z) | Set global camera position (3 args) |
| **SetCameraTarget**(x, y, z) | Set orbit/look-at target (3 args) |
| **UpdateCamera**() | Apply orbit state to camera |
| **Camera3DSetPosition**(x, y, z) | Alias; or (cameraId, x, y, z) for named camera |
| **Camera3DSetTarget**(x, y, z) | Alias; or (cameraId, x, y, z) |
| **Camera3DSetUp**(cameraId, x, y, z) | Camera up vector |
| **Camera3DSetFOV**(fov) | Global FOV (degrees) |
| **Camera3DSetProjection**(cameraId, type) | CAMERA_PERSPECTIVE / CAMERA_ORTHOGRAPHIC |
| **Camera3DMoveForward**(amount) | Move camera and target along look direction |
| **Camera3DMoveBackward**(amount) | Move backward (opposite of forward) |
| **Camera3DMoveRight**(amount) **Camera3DMoveLeft**(amount) | Move along right / left |
| **Camera3DMoveUp**(amount) **Camera3DMoveDown**(amount) | Move along camera up / down |
| **Camera3DRotateYaw**(angleRad) | Rotate position around target on Y axis |
| **Camera3DRotatePitch**(angleRad) | Tilt camera up/down (clamped) |
| **Camera3DRotateRoll**(angleRad) | Rotate camera up vector around forward axis |
| **BeginCamera2D**(cameraID) **EndCamera2D**() | Set active 2D camera by ID; use default if no ID |
| **Camera2DCreate**() | Create 2D camera → cameraID |
| **Camera2DSetPosition**(cameraID, x, y) | Set camera target (world position) |
| **Camera2DSetZoom**(cameraID, zoom) | Set zoom level |
| **Camera2DSetRotation**(cameraID, angle) | Set rotation (radians) |
| **Camera2DMove**(cameraID, dx, dy) | Move camera target by offset |
| **Camera2DSmoothFollow**(cameraID, targetX, targetY, speed) | Smooth follow; update each frame |
| **BeginCamera3D** **EndCamera3D** | Aliases of BeginMode3D, EndMode3D (3D camera) |

---

## Layer system (2D)

| Command | Description |
|--------|-------------|
| **LayerCreate**(name, order) | Create layer → layerID (order = draw priority, lower first) |
| **LayerSetOrder**(layerID, order) | Change draw order |
| **LayerSetVisible**(layerID, flag) | Hide (0) or show (non‑zero) layer |
| **LayerSetParallax**(layerID, parallaxX, parallaxY) | Parallax factors (e.g. 0.5 = half speed) |
| **LayerSetScroll**(layerID, scrollX, scrollY) | Scroll offset for layer |
| **LayerClear**(layerID) | Remove all sprites/tilemaps/particles from this layer |
| **LayerSortSprites**(layerID) | No‑op (flush sorts by z automatically) |
| **SpriteSetLayer**(spriteID, layerID) | Assign sprite to layer |
| **SpriteSetZIndex**(spriteID, z) | Draw order within layer (higher = on top) |
| **TilemapSetLayer**(tilemapID, layerID) | Assign tilemap to layer |
| **TilemapSetParallax**(tilemapID, px, py) | Parallax for tilemap |
| **ParticleSetLayer**(particleID, layerID) | Assign particle system to layer |

---

## Background system (2D)

| Command | Description |
|--------|-------------|
| **BackgroundCreate**(textureId) | Create background from texture → backgroundId |
| **BackgroundSetColor**(backgroundId, r, g, b, a) | Tint (0–1 or 0–255) |
| **BackgroundSetTexture**(backgroundId, textureId) | Change texture |
| **BackgroundSetScroll**(backgroundId, speedX, speedY) | Scroll speed |
| **BackgroundSetOffset**(backgroundId, offsetX, offsetY) | Base offset |
| **BackgroundSetParallax**(backgroundId, px, py) | Parallax factors |
| **BackgroundSetTiled**(backgroundId, flag) | Enable tiling |
| **BackgroundSetTileSize**(backgroundId, width, height) | Tile size for tiling |
| **BackgroundAddLayer**(backgroundId, textureId, parallaxX, parallaxY) | Add extra layer |
| **BackgroundRemoveLayer**(backgroundId, layerIndex) | Remove layer by index |
| **DrawBackground**(backgroundId) | Draw background (call in draw(); uses current 2D camera) |

---

## Tilemap (2D) – extended

| Command | Description |
|--------|-------------|
| **TilemapCreate**(tileWidth, tileHeight, mapWidth, mapHeight) | Create empty tilemap → mapId |
| **LoadTilemap**(path) | Load or create tilemap from path → mapId |
| **TilemapLoad**(path) | Alias of LoadTilemap |
| **TilemapSave**(tilemapId, path) | Save tilemap to JSON file |
| **TilemapFill**(tilemapId, tileId) | Fill all cells with tileId |
| **SetTile**(mapId, x, y, tileId) **TilemapSetTile**(…) | Set cell (x,y) to tileId |
| **GetTile**(mapId, x, y) **TilemapGetTile**(…) | Get tile at (x,y) |
| **DrawTilemap**(mapId) | Draw tile grid (respects layer/parallax) |
| **TilemapSetLayer**(tilemapId, layerId) **TilemapSetParallax**(tilemapId, px, py) | Layer and parallax |

---

## Sprite animation and batching

| Command | Description |
|--------|-------------|
| **SpriteSetFrame**(spriteId, frameIndex) | Set current frame (for atlas) |
| **SpriteSetFrameSize**(spriteId, frameWidth, frameHeight) | Set frame size (grid) |
| **SpriteSetFrameCount**(spriteId, count) | Number of frames (for wrap) |
| **SpriteSetAnimSpeed**(spriteId, framesPerSecond) | Animation speed |
| **SpritePlay**(spriteId) | Start playing animation |
| **SpritePause**(spriteId) | Pause animation |
| **SpriteStop**(spriteId) | Stop and reset to frame 0 |
| **SpriteBatchBegin**() | Start collecting SpriteDraw/DrawTexture |
| **SpriteBatchEnd**() | Flush batch (draw grouped by texture) |

---

## 2D particle emitter

| Command | Description |
|--------|-------------|
| **ParticleEmitterCreate**(textureId) | Create 2D emitter → emitterId |
| **ParticleEmitterSetRate**(emitterId, rate) | Particles per second |
| **ParticleEmitterSetLifetime**(emitterId, min, max) | Lifetime range (seconds) |
| **ParticleEmitterSetVelocity**(emitterId, vx, vy) | Base velocity |
| **ParticleEmitterSetSpread**(emitterId, angleRad) | Spread angle |
| **ParticleEmitterSetColor**(emitterId, r, g, b, a) | Particle color (0–1) |
| **ParticleEmitterSetLayer**(emitterId, layerId) | Assign to layer |
| **DrawParticleEmitter**(emitterId) | Update and draw (call in draw()) |

---

## 2D culling and atlas

| Command | Description |
|--------|-------------|
| **Enable2DCulling**(flag) | Skip sprites outside camera + margin |
| **SetCullingMargin**(pixels) | Margin around camera view |
| **AtlasLoad**(path) | Load atlas JSON + texture → atlasId |
| **AtlasGetRegion**(atlasId, name) | Get [x, y, w, h] for region name |
| **AtlasGetTextureId**(atlasId) | Get texture id for DrawTextureRec |

---

## 2D scene save/load

| Command | Description |
|--------|-------------|
| **SceneSave2D**(path) | Save 2D scene to JSON (version + placeholder data) |
| **SceneLoad2D**(path) | Load 2D scene from JSON (stub: validates format) |

---

## Physics 2D helpers

| Command | Description |
|--------|-------------|
| **Physics2DSetGravity**(worldId, x, y) | Set world gravity vector |
| **Physics2DRaycast**(originX, originY, dirX, dirY) | Raycast in world "default"; use RayHit* for result |
| **Physics2DSetLayerCollision**(layerA, layerB, flag) | Set whether two layers collide (for contact filter) |

---

## Terrain and water physics

| Command | Description |
|--------|-------------|
| **TerrainEnableCollision**(terrainId, flag) | Enable/disable terrain as physics collider (state stored) |
| **TerrainSetFriction**(terrainId, value) | Set friction for terrain collider |
| **TerrainSetBounce**(terrainId, value) | Set bounce for terrain collider |
| **WaterSetDensity**(waterId, value) | Set water density (affects buoyancy) |
| **WaterSetDrag**(waterId, linear, angular) | Set drag when submerged |
| **WaterApplyBuoyancy**(bodyId, waterId) | Apply buoyancy to body (stub; call each frame) |

---

## Vegetation physics and weather

| Command | Description |
|--------|-------------|
| **TreeEnableCollision**(treeId, flag) | Enable collision (capsule) on tree |
| **TreeSetCollisionRadius**(treeId, radius) | Set capsule radius |
| **TreeSetWind**(typeId, strength, speed) | Wind for tree type (shader) |
| **TreeApplyWind**(treeId, vx, vy, vz) | Apply wind vector (stub) |
| **TreeRaycast**(systemId, ox, oy, oz, dx, dy, dz) | Ray vs trees (stub) |
| **GrassSetBendAmount**(grassId, value) | Wind bend amount |
| **GrassSetInteraction**(grassId, flag) | Player displacement |
| **WeatherSetType**(type) **WeatherSetIntensity**(…) **WeatherSetWindDirection**(…) **WeatherSetWindSpeed**(…) **WeatherSetFogDensity**(…) **WeatherSetLightningFrequency**(…) | Weather stubs |
| **FireCreate** **FireSetSpreadRate** **FireSetSmokeEmitter** **FireSetLight** **SmokeSetDissolveRate** **SmokeSetRiseSpeed** | Fire/smoke stubs |
| **EnvironmentSetGlobalWind** **EnvironmentSetTemperature** **EnvironmentSetHumidity** **EnvironmentAffectParticles/Water/Vegetation** | Environment stubs |

---

## Navigation

| Command | Description |
|--------|-------------|
| **NavGridCreate**(width, height) **NavGridSetWalkable**(gridId, x, y, flag) **NavGridSetCost**(…) **NavGridFindPath**(gridId, startX, startY, endX, endY) | Grid pathfinding (stubs) |
| **NavMeshCreateFromTerrain**(terrainId) **NavMeshAddObstacle** **NavMeshRemoveObstacle** **NavMeshFindPath**(…) | NavMesh (stubs) |
| **NavAgentCreate** **NavAgentSetSpeed** **NavAgentSetRadius** **NavAgentSetDestination** **NavAgentGetNextWaypoint** | Agents (stubs) |

---

## Sky, time, clouds

| Command | Description |
|--------|-------------|
| **TimeSet**(hour) **TimeGet**() **TimeSetSpeed**(multiplier) | Time of day (stubs) |
| **SkyboxCreate** **SkyboxSetTexture** **SkyboxSetRotation** **SkyboxSetTint** **DrawSkybox** | Skybox (stubs) |
| **CloudLayerCreate** **CloudLayerSetTexture** **CloudLayerSetHeight** **DrawCloudLayer** | Clouds (stubs) |

---

## Indoor and interaction

| Command | Description |
|--------|-------------|
| **RoomCreate** **RoomSetBounds** **RoomAddPortal** **PortalCreate** **PortalSetOpen** | Rooms/portals (stubs) |
| **DoorCreate** **DoorSetOpen** **DoorToggle** **DoorSetLocked** **LeverCreate** **ButtonCreate** **SwitchCreate** | Doors/levers (stubs) |
| **TriggerCreate** **InteractableCreate** **PickupCreate** **LightZoneCreate** **WorldSaveInteractables** **WorldLoadInteractables** | Triggers/interact (stubs) |

---

## Streaming and editor

| Command | Description |
|--------|-------------|
| **WorldStreamEnable**(flag) **WorldStreamSetRadius**(…) **WorldStreamSetCenter**(x,y,z) **WorldLoadChunk**(chunkX, chunkZ) **WorldUnloadChunk**(…) **WorldIsChunkLoaded**(…) **WorldGetLoadedChunks** | Chunk streaming (stubs) |
| **EditorEnable**(flag) **EditorSetMode** **EditorSetBrushSize** **EditorSetBrushStrength** **EditorSetBrushFalloff** **EditorSetBrushShape** | Editor tools (stubs) |

---

## Decals

| Command | Description |
|--------|-------------|
| **DecalCreate**(textureId, x, y, z, size) **DecalSetLifetime** **DecalRemove** | Decals (stubs) |

---

## 3D models

| Command | Description |
|--------|-------------|
| **LoadModel**(path) | Load model from file → model id |
| **LoadModelFromMesh**(meshId) | Create model from mesh → model id |
| **DrawMesh**(meshId, materialId, posX,Y,Z, scaleX,Y,Z) | Draw mesh with position and scale |
| **DrawMeshMatrix**(meshId, materialId, m0..m15) | Draw mesh with full 4x4 matrix (row-major) |
| **DrawMeshInstanced**(meshId, materialId, instanceCount, …matrix floats) | Draw multiple instances |
| **GenMeshCube**(w, h, d) | Create cube mesh → mesh id |
| **LoadCube**(size) | Create cube model (single size) → model id |
| **DrawModel**(id, x, y, z, scale [, tint]) | Draw model at position with scale |
| **DrawModelSimple**(id, x, y, z [, angle]) | Draw at (x,y,z), scale 1; uses SetModelColor / RotateModel state |
| **DrawModelEx**(id, posX, posY, posZ, rotAxisX,Y,Z, rotAngle, scaleX,Y,Z [, tint]) | Full transform and tint |
| **DrawCube**(posX, posY, posZ, width, height, length, color) | Filled 3D cube (primitive) |
| **cube**(…) | Alias of DrawCube |
| **SetModelColor**(modelId, r, g, b, a) | Stored tint for DrawModelSimple |
| **SetModelPosition**(modelId, x, y, z) | **SetModelRotation**(modelId, axisX, axisY, axisZ, angleRad) | **SetModelScale**(modelId, sx, sy, sz) | Store transform for DrawModelWithState |
| **DrawModelWithState**(modelId [, tint]) | Draw model using stored position/rotation/scale |
| **RotateModel**(modelId, speedDegPerSec) | Auto-rotate model each frame |
| **UnloadModel**(id) | Free model |
| **SetModelShader**(modelId, shaderId) | Use custom shader on model's first material |
| **SetModelTexture**(modelId, textureId) | Set diffuse texture on first material |
| **SetMaterialFloat**(modelId, paramName, value) | Set float uniform on material shader |
| **SetMaterialVector**(modelId, paramName, x, y, z) | Set vec3 uniform on material shader |
| **LoadModelAnimations**(path) | Load animations from file; use GetModelAnimationId(n) for ids |
| **UpdateModelAnimation**(modelId, animId, frame) | Set animation frame |
| **IsModelAnimationValid**(modelId, animId) | True if anim applies to model |
| **UnloadModelAnimations**(animId, …) | Unload one or more animation ids |
| **DrawText3D**(fontId, text, x, y, z, fontSize, spacing [, r,g,b,a]) | Draw text at 3D position (projected to screen) |

---

## 2D drawing

| Command | Description |
|--------|-------------|
| **ClearBackground**(r, g, b, a) | Clear screen with RGBA |
| **Background**(r, g, b) | Clear with RGB (alpha 255) |
| **DrawText**(text, x, y, size, r, g, b, a) | Text at (x,y) with font size and color |
| **DrawTextSimple**(text, x, y) | Text at (x,y), size 20, white *(use for on-screen; PRINT = console)* |
| **MeasureText**(text, size) | Width of text in pixels |
| **DrawRectangle**(x, y, w, h, r, g, b, a) | Filled rectangle |
| **DrawRect**(…) **DrawRectFill**(…) | Alias: outline = DrawRectangleLines, filled = DrawRectangle |
| **rect**(…) | Alias of DrawRectangle |
| **DrawCircle**(x, y, radius, r, g, b, a) | Filled circle |
| **DrawCircleFill**(…) | Alias of DrawCircle |
| **circle**(…) | Alias of DrawCircle |
| **DrawPixel**(x, y, r, g, b, a) | Single pixel |
| **DrawLine**(x1, y1, x2, y2, r, g, b, a) | Line segment |
| **LoadTexture**(path) **UnloadTexture**(id) | Load / free texture |
| **DrawTexture**(id, x, y [, tint]) | Draw texture at position |
| **DrawTextureEx**(id, x, y, rotation, scale [, tint]) | With rotation and scale |
| **DrawTextureRec**(id, srcX, srcY, srcW, srcH, x, y [, tint]) | Draw part of texture (sprite frame) |
| **DrawTexturePro**(id, src…, dest…, origin, rotation, tint) | Full control (source/dest rects) |
| **DrawTextureFlipH**(textureId, x, y [, tint]) | Draw texture flipped horizontally |
| **DrawTextureFlipV**(textureId, x, y [, tint]) | Draw texture flipped vertically |
| **DrawTextureNPatch**(…) | Nine-patch / 9-slice drawing |
| **GetTextureWidth**(textureId) **GetTextureHeight**(textureId) | Texture dimensions |
| **GetTextureSize**(textureId) | Returns [width, height] |
| **LoadTilemap**(path) **DrawTilemap**(mapId [, x, y]) | Tilemap (game package) |
| **LoadFont**(path) | Load font → font id (use with DrawTextExFont, MeasureTextEx) |

### Sprite (high-level 2D transform)

| Command | Description |
|--------|-------------|
| **CreateSprite**(textureId) | Create sprite from texture → spriteId |
| **SpriteSetPosition**(spriteId, x, y) | **SpriteSetScale**(spriteId, scale) | **SpriteSetScaleXY**(spriteId, sx, sy) |
| **SpriteSetRotation**(spriteId, angleRad) | **SpriteSetOrigin**(spriteId, ox, oy) | **SpriteSetFlip**(spriteId, flipX, flipY) |
| **SpriteDraw**(spriteId [, tint]) | Draw sprite with current transform |
| **DestroySprite**(spriteId) | Free sprite |

---

## Audio

| Command | Description |
|--------|-------------|
| **LoadSound**(path) **UnloadSound**(id) | Load / free sound |
| **PlaySound**(soundId) **StopSound**(soundId) | Start / stop playback |
| **SetSoundVolume**(soundId, volume) | Set volume (0.0–1.0) |
| **LoadMusic**(path) **PlayMusicStream**(id) **StopMusicStream**(id) **UnloadMusicStream**(id) | Music stream |
| **UpdateMusic**(id) / **UpdateMusicStream**(id) | Call each frame while music plays (buffer streaming) |
| **UnloadMusic**(id) | Alias of UnloadMusicStream |
| **PauseMusicStream** **ResumeMusicStream** | Pause / resume |

---

## File

| Command | Description |
|--------|-------------|
| **FileExists**(path) | True if file exists (raylib core) |
| **LoadText**(path) | Read entire file as string (std) |
| **SaveText**(path, text) | Write string to file (std) |
| **ReadFile**(path) | Same as LoadText; returns nil on error |
| **WriteFile**(path, contents) | Same as SaveText |

---

## Game loop

| Command | Description |
|--------|-------------|
| **BeginFrame**() | Start frame (alias BeginDrawing) |
| **EndFrame**() | End frame (alias EndDrawing) |
| **RunGameLoop**(UpdateFunction) | *Not implemented* — use `WHILE NOT WindowShouldClose()` … `WEND` and call your update/draw code inside the loop. |

---

## Utility

| Command | Description |
|--------|-------------|
| **PrintDebug**(value) | Print value to stderr for debugging |
| **TimeNow**() | Seconds since epoch (float) |
| **Random**(min, max) | Integer in [min, max] inclusive |
| **Random**(n) | Integer in 0..n−1 |
| **Str**(x) | Number to string |
| **Val**(s) | String to number |

---

## Extended commands

### Window & system (extended)
| **SetWindowTitle**(title) | **SetWindowSize**(w, h) | **ToggleFullscreen**() | **IsFullscreen**() | **GetScreenWidth**() | **GetScreenHeight**() | **Screenshot**(path) (alias TakeScreenshot) |

### Input (extended)
| **MouseDown**(button) | **MousePressed**(button) | **MouseReleased**(button) | **GetMousePosition**() → [x,y] | **GetMouseDelta**() → [dx,dy] | **GamepadConnected**(id) | **GetGamepadAxis**(id, axis) | **GamepadButtonDown**(id, button) |

### Math (extended)
| **Sin**(x) **Cos**(x) **Tan**(x) **Sqrt**(x) | **RandomFloat**(min, max) | **RandomInt**(min, max) | **Vec3Add**(x1,y1,z1, x2,y2,z2) **Vec3Sub** **Vec3Scale**(x,y,z, s) **Vec3Normalize**(x,y,z) |

### Camera (extended)
| **CameraFPS**() | **CameraFree**() | **CameraSetFOV**(fov) | **CameraSetClipping**(near, far) | **CameraShake**(amount, duration) | **CAMERA3D**() → cameraId | **SetCurrentCamera**(cameraId) | **BeginMode3D** **EndMode3D** |

### Shaders & lighting
| **LoadShader**(vsPath, fsPath) **UnloadShader**(id) | **BeginShaderMode**(shaderId) **EndShaderMode**() | **SetShaderUniform**(id, name, value) | **SetShaderValueMatrix**(id, name, m0…m15) | **SetShaderValueTexture**(id, name, textureId) |

### 3D raycasting & collision
| **GetMouseRay**() | **GetRayCollisionMesh**(rayPos 3, rayDir 3, meshId, pos 3, scale 3) | **GetRayCollisionModel**(ray 6, modelId, pos 3, scale 3) | **GetRayCollisionTriangle**(ray 6, p1 3, p2 3, p3 3) | Use **GetRayCollisionPointX/Y/Z**(), **GetRayCollisionDistance**() for last hit |

### 3D models (extended)
| **LoadModelAnimated**(path) | **PlayModelAnimation**(model, anim) | **SetModelTexture**(model, texture) | **LoadTexture**(path) **UnloadTexture**(id) | **SetObjectPosition**(obj, x,y,z) **SetObjectRotation**(obj, pitch,yaw,roll) **SetObjectScale**(obj, sx,sy,sz) | **ObjectLookAt**(obj, x,y,z) |

### 2D drawing (extended)
| **DrawSprite**(spriteId, x, y) | **LoadSprite**(path) (alias LoadTexture) | **DrawLine**(x1,y1,x2,y2, r,g,b,a) | **DrawTriangle**(…) | **DrawTexture**(id, x, y) **sprite**(…) (alias) | **MeasureText**(text, size) |

### Audio (extended)
| **LoadMusic**(path) | **PlayMusic**(id) **PauseMusic**(id) **ResumeMusic**(id) | **SetMusicVolume**(id, vol) | **IsMusicPlaying**(id) |

### File I/O (extended)
| **LoadJSON**(path) **SaveJSON**(path, object) | **LoadImage**(path) **SaveImage**(imageId, path) | **DirectoryList**(path) (alias ListDir) |

### Image-level commands (CPU-side; edit pixels in RAM)

| Command | Description |
|--------|-------------|
| **LoadImage**(path) **UnloadImage**(imageId) | Load / free image |
| **GetImageColor**(imageId, x, y) | Get pixel (returns r,g,b,a) |
| **SetImageColor**(imageId, x, y, r, g, b, a) | Set pixel (alias ImageDrawPixel) |
| **ImageClearBackground**(imageId, r, g, b, a) | Fill image with color |
| **ImageDrawPixel**(imageId, x, y, r, g, b, a) | Draw one pixel |
| **ImageDrawRectangle**(imageId, x, y, w, h, r, g, b, a) | Filled rectangle on image |
| **ImageDrawCircle**(imageId, cx, cy, radius, r, g, b, a) | Circle on image |
| **ImageResize**(imageId, newW, newH) **ImageResizeNN**(…) | Resize (bilinear / nearest-neighbor) |
| **ImageRotateCW**(imageId) **ImageFlipHorizontal**(imageId) **ImageFlipVertical**(imageId) | Transform |
| **ImageCrop**(imageId, x, y, w, h) **ImageColorTint**(imageId, r, g, b, a) | Crop / tint |
| **LoadTextureFromImage**(imageId) | Upload image to GPU → texture id |
| **ExportImage**(imageId, path) **ExportImageAsCode**(imageId, path) | Save to file / export as code |
| **GenImageColor**(w, h, r, g, b, a) **GenImageGradientLinear**(…) **GenImageChecked**(…) etc. | Generate images |

### Gameplay helpers
| **TimerStart**(name) | **TimerElapsed**(name) → seconds | **CollisionBox**(x,y,z, w,h,d) → boxId | **CheckCollision**(boxIdA, boxIdB) | **RayCast**(ox,oy,oz, dx,dy,dz [, boxId]) → distance or −1 |

### Game loop (extended)
| **SetUpdateFunction**(func) **SetDrawFunction**(func) **Run**() | No-op; use `WHILE NOT WindowShouldClose()` … `WEND` and call your update/draw code. |

### Debugging & development
| **ShowFPS**(x, y) | **Log**(value) (stderr) | **Assert**(condition, message) (std) |

---

## Physics commands

### 2D physics (Box2D)

| Command | Description |
|--------|-------------|
| **CreateWorld2D**(worldName$, gravityX, gravityY) | Create 2D physics world |
| **Physics2DCreateWorld**(gravityX, gravityY) | Same using world name `"default"` |
| **Step2D**(worldName$, dt) **StepAllPhysics2D**(dt) | Step world(s) |
| **Physics2DStep**(dt) | Step world `"default"` |
| **CreateBox2D**(world$, body$, x, y, w, h, mass, isDynamic) | Create box body → bodyId |
| **CreateCircle2D**(world$, body$, x, y, radius, mass, isDynamic) | Create circle body |
| **GetPositionX2D**(world$, body$) **GetPositionY2D**(…) | Body position |
| **SetVelocity2D**(world$, body$, vx, vy) **ApplyForce2D**(…) **ApplyImpulse2D**(…) | Velocity and forces |
| **SetCollisionHandler**(bodyId, subName) | When bodyId collides, call Sub subName(otherBodyId) |
| **ProcessCollisions2D**(worldId) | Dispatch collision callbacks (call after Step2D) |
| **BOX2D.CreateWorld** **BOX2D.Step** **BOX2D.CreateBody** etc. | Same API with BOX2D. prefix |

### 3D physics (Bullet)

| Command | Description |
|--------|-------------|
| **CreateWorld3D**(worldName$, gx, gy, gz) **Step3D**(worldName$, dt) **StepAllPhysics3D**(dt) | Create and step 3D world |
| **PhysicsEnable**() | Create/enable default 3D physics world |
| **PhysicsDisable**() | Destroy default world |
| **PhysicsSetGravity**(x, y, z) | Set gravity of default world |
| **CreateRigidBody**([model,] mass) | Create box body in default world → bodyId |
| **ApplyForce**(bodyId, fx, fy, fz) **ApplyImpulse**(bodyId, ix, iy, iz) | Force / impulse |
| **SetBodyVelocity**(bodyId, vx, vy, vz) **GetBodyVelocity**(bodyId) | Linear velocity |
| **CheckCollision3D**(bodyIdA, bodyIdB) | → true if AABBs overlap |
| **BULLET.CreateWorld** **BULLET.Step** **BULLET.CreateBox** **BULLET.RayCast** etc. | Same API with BULLET. prefix |

---

## GUI commands (menus, HUDs)

### Immediate-mode UI (BeginUI / EndUI)

| Command | Description |
|--------|-------------|
| **BeginUI**() **EndUI**() | Start / end UI layout (vertical cursor) |
| **Label**(text) / **UILabel**(text) | Draw label at current layout position |
| **Button**(text) / **UIButton**(text) | Button; returns 1 if clicked else 0 |
| **Slider**(…) / **UISlider**(…) | Slider |
| **Checkbox**(text, checked) / **UICheckbox**(…) | Checkbox |
| **TextBox**(…) / **UITextBox**(…) | Text input |
| **ProgressBar**(…) / **UIProgressBar**(…) | Progress bar |
| **WindowBox**(title, x, y, w, h) **GroupBox**(text, x, y, w, h) | Panels |

### Raygui (positioned widgets)

| Command | Description |
|--------|-------------|
| **GuiButton**(x, y, w, h, text) | Button; returns 1 if clicked else 0 |
| **button**(…) | Alias of GuiButton |
| **GuiLabel**(x, y, w, h, text) | Label |
| **GuiSlider**(x, y, w, min, max, value) | Slider; returns current value |
| **GuiCheckbox**(text, x, y, checked) | Checkbox; returns 1 if checked else 0 |
| **GuiTextbox**(x, y, w, text) | Text box; returns current text |

---

## Particle system

| Command | Description |
|--------|-------------|
| **CreateParticleSystem**() | → systemId |
| **EmitParticles**(systemId, count) | Spawn particles (use default color/lifetime/velocity) |
| **SetParticleColor**(systemId, r, g, b, a) | Default color for new particles |
| **SetParticleLifetime**(systemId, seconds) | Default lifetime |
| **SetParticleVelocity**(systemId, vx, vy, vz) | Default velocity |
| **DrawParticles**(systemId) | Update and draw (call inside 3D mode) |

---

## Scene management

| Command | Description |
|--------|-------------|
| **LoadScene**(name) | Set current scene by id |
| **SaveScene**(sceneId, path) | Save scene metadata to JSON |
| **AddToScene**(objectId) | Add object to current scene |
| **RemoveFromScene**(objectId) | Remove from current scene |
| **SceneExists**(name) | → true if scene exists |

---

## Entity–Component–System (ECS)

Unprefixed commands use default world `"default"`. Use **CreateEntity**() then **AddComponent**(entityId, type, …).

| Command | Description |
|--------|-------------|
| **CreateEntity**() | Create entity in default world → entityId |
| **AddComponent**(entityId, type [, …]) | Add Transform, Sprite, Health, Parent (type + args) |
| **GetComponent**(entityId, type) | → map (e.g. Transform → {x,y,z}) |
| **RemoveComponent**(entityId, type) | Remove component |
| **RunSystem**(name) | No-op; iterate via ECS.QueryCount / ECS.QueryEntity in script |

---

## AI commands

Store position with **AISetPosition**(entityId, x, y, z). Each frame call **AIUpdate**(entityId) then **GetAIPosition**(entityId) to draw.

| Command | Description |
|--------|-------------|
| **AIMoveTo**(entityId, x, y, z) | Set movement target |
| **AISetSpeed**(entityId, speed) | Movement speed |
| **AISetPosition**(entityId, x, y, z) | Set current position |
| **GetAIPosition**(entityId) | → [x, y, z] |
| **AIUpdate**(entityId) | Move toward target this frame |
| **AIWander**(entityId, radius) | Set wander radius |
| **AIChase**(entityId, targetEntityId) | Set target to another entity’s position |
| **AIFlee**(entityId, targetEntityId) | Set target away from entity |

---

## Networking

Aliases for Host, Connect, Send, Receive, IsConnected.

| Command | Description |
|--------|-------------|
| **NetHost**(port) | Start server → serverId |
| **NetConnect**(ip, port) | Connect → connectionId |
| **NetSend**(connectionId, data) | Send text |
| **NetReceive**(connectionId) | → received text or nil |
| **NetIsConnected**(connectionId) | → 1 if connected else 0 |

---

## Coroutine / async (stubs)

VM does not support yielding; these are no-ops. Use timers and state in script instead.

| Command | Description |
|--------|-------------|
| **CoroutineStart**(func) | No-op |
| **CoroutineYield**() | No-op |
| **CoroutineWait**(seconds) | No-op |
| **CoroutineStop**(id) | No-op |

---

## Animation (lerp over time)

Call each frame; returns or drives current value. **AnimateValue**(key, target, duration) returns current value.

| Command | Description |
|--------|-------------|
| **AnimateValue**(key, target, durationSeconds) | Lerp value; returns current each call |
| **AnimateColor**(key, r, g, b, a, durationSeconds) | Lerp color |
| **AnimatePosition**(key, x, y, z, durationSeconds) | Lerp position |
| **AnimateRotation**(key, pitch, yaw, roll, durationSeconds) | Lerp rotation |

---

## Utility / engine

| Command | Description |
|--------|-------------|
| **SetLogLevel**(level) | 0=off, 1=error, 2=warn, 3=info, 4=debug (std) |
| **GetTime**() | Raylib time (raylib) |
| **DeltaTime**() | Same as GetFrameTime() (raylib) |
| **SeedRandom**(seed) | Set random seed (raylib) |
| **UUID**() | New UUID string (std) |

---

## Lighting

State is stored for use with custom shaders; raylib has no built-in lighting.

| Command | Description |
|--------|-------------|
| **CreateLight**(type, x, y, z) | Create light → lightId |
| **SetLightColor**(lightId, r, g, b) | Set light color (0–255) |
| **SetLightIntensity**(lightId, amount) | Set intensity |
| **SetLightDirection**(lightId, x, y, z) | Direction vector |
| **EnableShadows**() / **DisableShadows**() | Toggle shadow state |

---

## Material & shader

| Command | Description |
|--------|-------------|
| **LoadShader**(vertexPath, fragmentPath) | Load shader → shaderId |
| **SetShaderUniform**(shaderId, name, value) | Set float uniform |
| **ApplyShader**(shaderId) | Same as BeginShaderMode |
| **RemoveShader**() | End current shader mode |
| **SetMaterialTexture**(modelId, textureId) | Set model diffuse texture |
| **SetMaterialColor**(modelId, r, g, b, a) | Set model tint |

---

## Terrain

| Command | Description |
|--------|-------------|
| **GenerateTerrain**(width, depth, scale) | Create empty height grid → terrainId |
| **LoadHeightmap**(imageId) | Build terrain from image (gray = height) → terrainId |
| **SetTerrainTexture**(terrainId, textureId) | Associate texture |
| **GetTerrainHeight**(terrainId, x, z) | Sample height at world x,z |

---

## Skybox

| Command | Description |
|--------|-------------|
| **LoadSkybox**(folderPath) | Load cubemap (stub) |
| **SetSkyColor**(r, g, b) | Clear/sky color (0–255) |
| **EnableSkybox**() / **DisableSkybox**() | Toggle skybox draw |

---

## Post-processing

State only; actual effects require render-to-texture and shaders.

| Command | Description |
|--------|-------------|
| **EnableBloom**() | Set bloom on |
| **SetBloomIntensity**(amount) | Bloom strength |
| **EnableMotionBlur**() | Motion blur state |
| **EnableCRTFilter**() | CRT effect state |
| **EnablePixelate**(size) | Pixelate block size |

---

## Advanced GUI / UI layout

| Command | Description |
|--------|-------------|
| **GuiPanel**(x, y, w, h [, text]) | Panel (4 or 5 args) |
| **GuiWindow**(title, x, y, w, h) | Window box → 1 if close clicked |
| **GuiList**(items, x, y, w, h) | List; items = "A;B;C" → selected index |
| **GuiDropdown**(items, x, y, w) | Dropdown → selected index |
| **GuiProgressBar**(x, y, w, value) | Full 9-arg form; **GuiProgressBarSimple**(x, y, w, value) for 0–1 |

---

## Tilemap (2D)

| Command | Description |
|--------|-------------|
| **LoadTilemap**(path) | Create tilemap → mapId (path stub) |
| **DrawTilemap**(mapId) | Draw tile grid |
| **SetTile**(mapId, x, y, tileID) | Set tile at cell |
| **GetTile**(mapId, x, y) | Get tile ID at cell |
| **TilemapCollision**(mapId, x, y) | True if tile at world (x,y) is solid (mark tiles solid via data) |

---

## Pathfinding

| Command | Description |
|--------|-------------|
| **PathfindGrid**(mapId, startX, startY, endX, endY) | → path (stub: empty) |
| **PathfindNavmesh**(mesh, start, end) | → path (stub) |
| **FollowPath**(entityId, path) | No-op |

---

## Scripting / event (stubs)

VM cannot pass function references; use key/mouse checks in main loop instead.

| Command | Description |
|--------|-------------|
| **OnKeyPress**(key, function) | No-op |
| **OnMouseClick**(button, function) | No-op |
| **OnUpdate**(function) | No-op |
| **OnDraw**(function) | No-op |
| **OnCollision**(entityA, entityB, function) | No-op |

---

## Debugging & development

| Command | Description |
|--------|-------------|
| **DebugDrawGrid**([slices, spacing]) | Draw 3D grid (default 10, 1) |
| **DebugDrawBounds**(objectId) | Draw bounds (stub) |
| **DebugLog**(value) | Print to stderr |
| **DebugWatch**(variable) | No-op |

---

## Procedural generation

| Command | Description |
|--------|-------------|
| **Noise2D**(x, y) | Value noise in [0,1] (deterministic) |
| **Noise3D**(x, y, z) | 3D value noise in [0,1] |
| **GenerateDungeon**(width, height) | Create tilemap with random rooms → mapId (0=wall, 1=floor) |
| **GenerateTree**(seed) | Return tree id string (deterministic from seed) |
| **GenerateCity**(size) | Create city tilemap (streets + buildings) → mapId |

---

## Save / load system

| Command | Description |
|--------|-------------|
| **SaveGame**(path, data) | Save data to path (data = string, dict, or JSON handle) |
| **LoadGame**(path) | Load JSON from path → handle |
| **Autosave**(intervalSeconds [, path]) | Set autosave interval (script calls SaveGame when needed) |
| **SaveExists**(path) | True if file exists |

---

## Localization

| Command | Description |
|--------|-------------|
| **LoadLanguage**(path) | Load JSON key→value from path; code from filename (e.g. en.json → "en") |
| **SetLanguage**(code) | Set current language code |
| **Translate**(key) | Return translated string for key, or key if missing |

---

## Terrain sculpting

Modify height grid of terrains created with **GenerateTerrain** / **LoadHeightmap**. Brush size and strength are global.

| Command | Description |
|--------|-------------|
| **TerrainRaise**(terrainId, x, z, radius, amount) | Raise heights in radius (world x,z) |
| **TerrainLower**(terrainId, x, z, radius, amount) | Lower heights |
| **TerrainSmooth**(terrainId, x, z, radius) | Smooth by averaging neighbors |
| **TerrainFlatten**(terrainId, x, z, radius, height) | Blend toward target height (0–1) |
| **TerrainPaint**(terrainId, x, z, radius, textureID) | Set terrain texture |
| **TerrainSetMaterial**(terrainId, material) | Set material id |
| **TerrainBrushSetSize**(size) | Brush radius in grid cells |
| **TerrainBrushSetStrength**(strength) | Blend strength 0–1 |
| **TerrainUndo**(terrainId) | Restore last saved heights |

---

## Dialogue system

Load JSON: `{ "nodeId": { "text": "...", "next": "otherId", "choices": [{"text":"...", "next":"..."}] } }`. Use **DialogueStart**(id), then **DialogueNext**() or **DialogueChoice**(index). Draw text/choices in your UI from current node.

| Command | Description |
|--------|-------------|
| **DialogueLoad**(path) | Load dialogue nodes from JSON |
| **DialogueStart**(id) | Set current node |
| **DialogueNext**() | Advance to node["next"] |
| **DialogueChoice**(index) | Go to choices[index]["next"] |
| **DialogueShowText**(text) / **DialogueShowChoices**(choices) | No-op; draw in your UI |
| **DialogueSetVar**(name, value) / **DialogueGetVar**(name) | Dialogue variables |

---

## Inventory system

| Command | Description |
|--------|-------------|
| **InventoryCreate**(size) | Create inventory → invId |
| **InventoryAddItem**(invId, itemID, amount) | Add or stack item |
| **InventoryRemoveItem**(invId, itemID, amount) | Remove amount |
| **InventoryHasItem**(invId, itemID) | True if any amount |
| **ItemDefine**(id, name, icon, stackSize) | Define item type |
| **ItemSetProperty**(itemId, key, value) | Extra properties |
| **InventoryDraw**(invId, x, y) | Draw slot grid (5 columns) |

---

## Physics joints & ragdolls

Stubs; use **BULLET.*** for real 3D joints.

| Command | Description |
|--------|-------------|
| **CreateHingeJoint**(bodyA, bodyB, anchor, axis) | Stub → "" |
| **CreateBallJoint**(bodyA, bodyB, anchor) | Stub |
| **CreateSliderJoint**(bodyA, bodyB, axis) | Stub |
| **CreateRagdoll**(model) | Stub |
| **RagdollEnable** / **RagdollDisable** | No-op |

---

## AI behavior trees

Build trees from selector/sequence/action/condition nodes; **AIRun** is a stub.

| Command | Description |
|--------|-------------|
| **AIBehaviorTreeCreate**() | → tree id |
| **AISelector**(child1, child2, …) | Priority node → node id |
| **AISequence**(child1, child2, …) | Sequence node → node id |
| **AIAction**(functionName) / **AICondition**(functionName) | Leaf nodes |
| **AIRun**(treeId, entityId) | No-op |

---

## Multiplayer replication

State flags for what to sync; **NetStartServer** / **NetStartClient** are stubs. Use **Host** / **Connect** and **Send** / **Receive** for real networking.

| Command | Description |
|--------|-------------|
| **NetStartServer**(port) / **NetStartClient**(ip, port) | Stubs |
| **ReplicateVariable**(entityId, varName) | Mark variable for sync |
| **ReplicatePosition**(entityId) / **ReplicateRotation**(entityId) | Mark transform sync |
| **RPC**(functionName, args) | No-op |

---

## Shader graph

Node-based shader building; **ShaderGraphCompile** returns empty string (stub).

| Command | Description |
|--------|-------------|
| **ShaderNodeTexture**(textureId) | Texture node → node id |
| **ShaderNodeColor**(r, g, b, a) | Color node |
| **ShaderNodeAdd**(a, b) / **ShaderNodeMultiply**(a, b) | Math nodes |
| **ShaderNodeTime**() | Time node |
| **ShaderGraphCreate**() | → graph id |
| **ShaderGraphConnect**(graphId, outputNodeId, inputNodeId) | Connect nodes |
| **ShaderGraphCompile**(graphId) | Stub → "" |

---

## Animation state machines

| Command | Description |
|--------|-------------|
| **AnimStateCreate**(name) | → state id |
| **AnimStateSetClip**(stateId, clipId) | Bind clip to state |
| **AnimTransition**(fromStateId, toStateId, condition) | Add transition |
| **AnimSetParameter**(name, value) | Set param for conditions |
| **AnimSetState**(entityId, stateId) | Set entity state |
| **AnimUpdate**(entityId) | No-op |

---

## Terrain (heightmap + mesh)

| Command | Description |
|--------|-------------|
| **LoadHeightmap**(imageId) | Create heightmap from image (grayscale) → heightmap id |
| **GenHeightmap**(width, depth, noiseScale) | Procedural heightmap → heightmap id |
| **GenHeightmapPerlin**(width, depth, offsetX, offsetY, scale) | Perlin noise heightmap → heightmap id |
| **GenTerrainMesh**(heightmapId, sizeX, sizeZ, heightScale [, lod]) | Build mesh from heightmap → mesh id |
| **TerrainCreate**(heightmapId, sizeX, sizeZ, heightScale) | Create terrain → terrain id |
| **TerrainUpdate**(terrainId) | Rebuild mesh from heightmap |
| **DrawTerrain**(terrainId, posX, posY, posZ) | Draw terrain at position (Render3D) |
| **SetTerrainTexture**(terrainId, textureId) | **SetTerrainMaterial**(terrainId, materialId) |
| **SetTerrainLOD**(terrainId, lodLevel) | Set LOD for next TerrainUpdate |
| **TerrainRaise**(terrainId, x, z, radius, amount) | **TerrainLower**(…) | **TerrainSmooth**(…) | **TerrainFlatten**(terrainId, x, z, radius, targetHeight) |
| **TerrainPaint**(terrainId, x, z, radius, paintValue [, blend]) | Blend paint value in disk |
| **TerrainGetHeight**(terrainId, x, z) | World Y at (x,z) |
| **TerrainGetNormal**(terrainId, x, z) | Normal vector at (x,z) → [nx, ny, nz] |
| **TerrainRaycast**(terrainId, ox, oy, oz, dx, dy, dz) | Ray vs terrain → [hit, dist, hx, hy, hz] |

---

## Water

| Command | Description |
|--------|-------------|
| **WaterCreate**(width, depth, tileSize) | Create water plane → water id |
| **DrawWater**(waterId [, posX, posY, posZ]) | Draw water (Render3D) |
| **SetWaterPosition**(waterId, x, y, z) | **SetWaterWaveSpeed**(waterId, speed) | **SetWaterWaveHeight**(waterId, height) | **SetWaterWaveFrequency**(waterId, freq) |
| **SetWaterTime**(waterId, time) | **WaterGetHeight**(waterId, x, z) | Wave formula height |
| **SetWaterTexture** / **SetWaterReflectionTexture** / **SetWaterRefractionTexture** / **SetWaterNormalMap** | Store texture refs |
| **SetWaterColor**(waterId, r, g, b, a) | **SetWaterShininess**(waterId, shininess) |
| **WaterEnableFoam**(waterId, enabled) | **WaterSetFoamIntensity**(…) | **WaterSetDepthFade**(…) | **WaterSetTransparency**(…) |

---

## Vegetation (trees and grass)

| Command | Description |
|--------|-------------|
| **TreeTypeCreate**(modelId, trunkTexId, leafTexId) | → tree type id |
| **TreeSystemCreate**() | → tree system id |
| **TreePlace**(systemId, typeId, x, y, z, scale, rotation) | Place tree → tree id |
| **TreeRemove**(treeId) | **TreeSetPosition**(treeId, x, y, z) | **TreeSetScale**(treeId, scale) | **TreeSetRotation**(treeId, rotation) |
| **TreeSystemSetLOD**(systemId, near, mid, far) | **TreeSystemEnableInstancing**(systemId, on) |
| **TreeGetAt**(systemId, x, z) | Nearest tree id at (x,z) |
| **DrawTrees**(systemId) | Draw all trees in system (Render3D) |
| **GrassCreate**(textureId, density, patchSize) | → grass id |
| **GrassSetWind**(grassId, speed, strength) | **GrassSetHeight**(grassId, height) | **GrassSetColor**(grassId, r, g, b, a) |
| **GrassPaint**(grassId, x, z, radius, density) | **GrassErase**(grassId, x, z, radius) | **GrassSetDensity**(grassId, density) |
| **GrassSetLOD**(grassId, dist) | **GrassEnableInstancing**(grassId, on) |
| **DrawGrass**(grassId) | Draw grass (Render3D) |

---

## Object placement

| Command | Description |
|--------|-------------|
| **ObjectPlace**(modelId, x, y, z, scale, rotation) | Place object → object id |
| **ObjectRemove**(objectId) | **ObjectSetTransform**(objectId, x, y, z, scaleX, scaleY, scaleZ, rotAxisX, rotAxisY, rotAxisZ, rotAngle) |
| **ObjectRandomScatter**(modelId, areaX, areaZ, count, minScale, maxScale) | Scatter objects → [id,…] |
| **ObjectPaint**(modelId, x, z, radius, density) | **ObjectErase**(x, z, radius) |
| **ObjectGetAt**(x, z) | Nearest object id at (x,z) |
| **ObjectRaycast**(ox, oy, oz, dx, dy, dz) | Ray vs objects → [hit, objectId, hx, hy, hz] |
| **DrawObject**(objectId) | **DrawAllObjects**() |

---

## World save/load

| Command | Description |
|--------|-------------|
| **WorldSave**(path) | Save world (objects, etc.) to file |
| **WorldLoad**(path) | Load world from file |
| **WorldExportJSON**(path) | Export as JSON |
| **WorldImportJSON**(path) | Import from JSON |

---

## Procedural generation

| Command | Description |
|--------|-------------|
| **NoisePerlin2D**(x, y, scale) | **NoiseFractal2D**(x, y, octaves, persistence, lacunarity) | **NoiseSimplex2D**(x, y, scale) | Noise value [0,1] |
| **ScatterTrees**(treeSystemId, treeTypeId, areaX, areaZ, density) | **ScatterGrass**(grassId, centerX, centerZ, radius, density) |
| **ScatterObjects**(modelId, areaX, areaZ, count [, minScale, maxScale]) |

---

## Optimization

| Command | Description |
|--------|-------------|
| **SetCullingDistance**(distance) | Max draw distance for culling |
| **EnableFrustumCulling**(flag) | Enable/disable frustum culling |

---

## Mesh (low-level)

| Command | Description |
|--------|-------------|
| **MeshCreate**(vertices, normals, uvs, indices) | Create mesh from arrays → mesh id |
| **MeshUpdate**(meshId) | Re-upload mesh to GPU |
| **MeshSetVertices**(meshId, array) | **MeshSetNormals**(meshId, array) | **MeshSetUVs**(meshId, array) | **MeshSetIndices**(meshId, array) |

---

## Minimal 3D example

```basic
InitWindow(800, 600, "Game")
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
```

See [API_REFERENCE.md](../API_REFERENCE.md) for the full binding list and [examples/spinning_cube_simple.bas](../examples/spinning_cube_simple.bas) for a runnable version.
