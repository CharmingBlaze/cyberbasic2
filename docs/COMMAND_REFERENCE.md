# CyberBasic – Complete Command Reference

Structured command set for window, input, math, camera, 3D, 2D, audio, file, game loop, and utilities. All names are **case-insensitive**.

---

## Window & system

| Command | Description |
|--------|-------------|
| **InitWindow**(width, height, title) | Open game window |
| **CloseWindow**() | Close window and exit |
| **SetTargetFPS**(fps) | Target frames per second |
| **GetFrameTime**() | Delta time since last frame (seconds) |
| **WindowShouldClose**() | True when user requested close |
| **DisableCursor**() | Hide and confine mouse |
| **EnableCursor**() | Show mouse cursor |

---

## Input

| Command | Description |
|--------|-------------|
| **KeyDown**(key) | True while key held (use KEY_W, KEY_ESCAPE, etc.) |
| **KeyPressed**(key) | True once when key pressed |
| **GetMouseX**() | Mouse X position |
| **GetMouseY**() | Mouse Y position |
| **GetMouseDeltaX**() | Mouse movement X this frame |
| **GetMouseDeltaY**() | Mouse movement Y this frame |
| **GetMouseWheelMove**() | Scroll wheel delta |
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

---

## 3D models

| Command | Description |
|--------|-------------|
| **LoadModel**(path) | Load model from file → model id |
| **LoadModelFromMesh**(meshId) | Create model from mesh → model id |
| **GenMeshCube**(w, h, d) | Create cube mesh → mesh id |
| **LoadCube**(size) | Create cube model (single size) → model id |
| **DrawModel**(id, x, y, z, scale [, tint]) | Draw model at position with scale |
| **DrawModelSimple**(id, x, y, z [, angle]) | Draw at (x,y,z), scale 1; uses SetModelColor / RotateModel state |
| **DrawModelEx**(id, posX, posY, posZ, rotAxisX,Y,Z, rotAngle, scaleX,Y,Z [, tint]) | Full transform and tint |
| **SetModelColor**(modelId, r, g, b, a) | Stored tint for DrawModelSimple |
| **RotateModel**(modelId, speedDegPerSec) | Auto-rotate model each frame |
| **UnloadModel**(id) | Free model |

---

## 2D drawing

| Command | Description |
|--------|-------------|
| **ClearBackground**(r, g, b, a) | Clear screen with RGBA |
| **Background**(r, g, b) | Clear with RGB (alpha 255) |
| **DrawText**(text, x, y, size, r, g, b, a) | Text at (x,y) with font size and color |
| **DrawTextSimple**(text, x, y) | Text at (x,y), size 20, white *(use for on-screen; PRINT = console)* |
| **DrawRectangle**(x, y, w, h, r, g, b, a) | Filled rectangle |
| **DrawCircle**(x, y, radius, r, g, b, a) | Filled circle |

---

## Audio

| Command | Description |
|--------|-------------|
| **LoadSound**(path) | Load sound file → sound id |
| **PlaySound**(soundId) | Start playback |
| **StopSound**(soundId) | Stop playback |
| **SetSoundVolume**(soundId, volume) | Set volume (0.0–1.0) |

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
| **CameraFPS**() | **CameraFree**() | **CameraSetFOV**(fov) | **CameraSetClipping**(near, far) | **CameraShake**(amount, duration) |

### 3D models (extended)
| **LoadModelAnimated**(path) | **PlayModelAnimation**(model, anim) | **SetModelTexture**(model, texture) | **LoadTexture**(path) **UnloadTexture**(id) | **SetObjectPosition**(obj, x,y,z) **SetObjectRotation**(obj, pitch,yaw,roll) **SetObjectScale**(obj, sx,sy,sz) | **ObjectLookAt**(obj, x,y,z) |

### 2D drawing (extended)
| **DrawSprite**(spriteId, x, y) | **LoadSprite**(path) (alias LoadTexture) | **DrawLine**(x1,y1,x2,y2, r,g,b,a) | **DrawTriangle**(…) | **DrawTexture**(id, x, y) | **MeasureText**(text, size) |

### Audio (extended)
| **LoadMusic**(path) | **PlayMusic**(id) **PauseMusic**(id) **ResumeMusic**(id) | **SetMusicVolume**(id, vol) | **IsMusicPlaying**(id) |

### File I/O (extended)
| **LoadJSON**(path) **SaveJSON**(path, object) | **LoadImage**(path) **SaveImage**(imageId, path) | **DirectoryList**(path) (alias ListDir) |

### Gameplay helpers
| **TimerStart**(name) | **TimerElapsed**(name) → seconds | **CollisionBox**(x,y,z, w,h,d) → boxId | **CheckCollision**(boxIdA, boxIdB) | **RayCast**(ox,oy,oz, dx,dy,dz [, boxId]) → distance or −1 |

### Game loop (extended)
| **SetUpdateFunction**(func) **SetDrawFunction**(func) **Run**() | No-op; use `WHILE NOT WindowShouldClose()` … `WEND` and call your update/draw code. |

### Debugging & development
| **ShowFPS**(x, y) | **Log**(value) (stderr) | **Assert**(condition, message) (std) |

---

## Physics commands (default 3D world)

Use `PhysicsEnable()` then `BULLET.Step("default", GetFrameTime())` each frame. Bodies use the default world `"default"`.

| Command | Description |
|--------|-------------|
| **PhysicsEnable**() | Create/enable default 3D physics world |
| **PhysicsDisable**() | Destroy default world |
| **PhysicsSetGravity**(x, y, z) | Set gravity of default world |
| **CreateRigidBody**([model,] mass) | Create box body in default world → bodyId |
| **ApplyForce**(bodyId, fx, fy, fz) | Apply force to body |
| **ApplyImpulse**(bodyId, ix, iy, iz) | Apply impulse |
| **SetBodyVelocity**(bodyId, vx, vy, vz) | Set linear velocity |
| **GetBodyVelocity**(bodyId) | → [vx, vy, vz] |
| **CheckCollision3D**(bodyIdA, bodyIdB) | → true if AABBs overlap |

---

## GUI commands (menus, HUDs)

| Command | Description |
|--------|-------------|
| **GuiButton**(text, x, y, w, h) | Button; returns 1 if clicked else 0 |
| **GuiLabel**(x, y, w, h, text) | Label (or GuiLabel text, x, y for simple) |
| **GuiSlider**(x, y, w, min, max, value) | Slider; returns current value (6-arg). Also (x,y,w,h, textL, textR, value, min, max) |
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
