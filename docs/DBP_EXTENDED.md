# DBP Extended Commands Reference

This document lists all DarkBASIC Pro-style high-level commands in CyberBASIC2, organized by module.

**See [CORE_COMMAND_REFERENCE.md](CORE_COMMAND_REFERENCE.md) for the complete required command set with status.**

## Package Structure

The DBP bindings are split into modular files under `compiler/bindings/dbp/`:

| File | Purpose |
|------|---------|
| `dbp.go` | Core: 2D/3D graphics, objects, scene, input, time, math, FPS camera |
| `dbp_textures.go` | Texture registry (id-based LoadTexture, DeleteTexture, filter, wrap) |
| `dbp_materials.go` | Material registry (MakeMaterial, SetMaterialColor, ApplyMaterial) |
| `dbp_camera.go` | Camera extras (CameraFollow, CameraOrbit, CameraShake, CameraSmooth) |
| `dbp_world.go` | World (SetSkybox, SetAmbientLight, SetFog, clouds, sun, time, weather) |
| `dbp_water.go` | Water (MakeWater, SetWaterTexture, PositionWater, scroll, waves, DrawWater) |
| `dbp_terrain.go` | Terrain (MakeTerrain, LoadHeightmap, SetTerrainLayer, SetTerrainSplatmap, DrawTerrain) |
| `dbp_groups.go` | Object groups (MakeGroup, AddToGroup, PositionGroup, DrawGroup) |
| `dbp_players.go` | Player state (MakePlayer, SetPlayerPosition, MovePlayer) |
| `dbp_audio.go` | Music (LoadMusic, PlayMusic, StopMusic, SetMusicVolume, SetMusicLoop) |
| `dbp_lighting.go` | Light registry (MakeLight, PositionLight, SetLightColor, etc.) |
| `dbp_physics.go` | Physics wrappers (PhysicsOn/Off, MakeRigidBody, SetGravity, etc.) |
| `dbp_collision.go` | Raycast, Spherecast, ObjectCollides, PointInObject |
| `dbp_particles.go` | MakeParticles, SetParticleColor/Size/Speed, EmitParticlesAt |
| `dbp_net.go` | Networking (wrappers over net package) |
| `dbp_file.go` | SaveString, LoadString, SaveValue, LoadValue |
| `dbp_runtime.go` | StopTask, PauseTask, ResumeTask, FixedUpdate, OnFixedUpdate |
| `dbp_replication.go` | Replication (from game package) |
| `dbp_2d.go` | 2D Game API: drawing, sprites, spritesheets (Aseprite + grid), tilemaps, camera, collision, physics, objects, particles |
| `dbp_3d.go` | 3D Game API: window, camera, objects, mesh, animation, terrain, math, replication |
| `dbp_mesh.go` | LoadMesh, GetModelBounds, GetMeshVertexCount, GetMeshTriangleCount |
| `dbp_animation.go` | LoadAnimation, PlayAnimation (multi-clip), StopAnimation, SetAnimationSpeed, SetAnimationLoop, ResetBones, LoadMeshAnimation, PlayMeshAnimation, SetMeshAnimationFrame |
| `dbp_level.go` | LoadLevel, DrawLevel, UnloadLevel, LoadLevelCollision, GetLevelColliderCount, GetLevelCollider, GetLevelObjectCount, GetLevelObject |
| `dbp_prefab.go` | LoadPrefab, SpawnPrefab |
| `dbp_ik.go` | IKEnable, IKSolveTwoBone |
| `dbp_instancing.go` | MakeInstance, PositionInstance, DrawInstances |
| `dbp_nav.go` | NavMeshLoad, NavMeshFindPath, NavMeshDraw |

## Core Commands (dbp.go)

### 2D Graphics
- `LoadImage(path, id)` - Load texture, store at integer id
- `Sprite(id, x, y)` - Draw texture at position
- `Cls()` / `CLS` - Clear screen
- `Ink(r, g, b)` - Set current draw color

### 3D Objects
- `LoadObject(path, id)` / `LoadObjectId(id, path)` - Load model
- `LoadCube(id, size)` / `MakeCube(id, size)` - Procedural cube
- `MakeSphere(id, radius)` - Procedural sphere
- `MakePlane(id, width, height)` - Procedural plane
- `PositionObject(id, x, y, z)` - Set position
- `RotateObject(id, pitch, yaw, roll)` - Set rotation
- `YRotateObject(id, angle)` - Add to yaw
- `ScaleObject(id, sx, sy, sz)` - Set scale
- `MoveObject(id, x, y, z)` - Add to position
- `TurnObject(id, pitch, yaw, roll)` - Add to rotation
- `DrawObject(id)` - Draw object
- `HideObject(id)` / `ShowObject(id)` - Visibility
- `DeleteObject(id)` - Remove and unload

### Object Extras
- `CloneObject(newID, sourceID)` - Clone object
- `CopyObject(newID, sourceID)` - Alias for CloneObject
- `ObjectExists(id)` - Returns 1 if exists, 0 otherwise
- `FixObject(id)` / `UnfixObject(id)` - Fixed/kinematic flag
- `SetObjectColor(id, r, g, b)` - Tint color
- `SetObjectAlpha(id, value)` - Alpha
- `SetObjectTexture(id, textureID)` or `(id, "path.png")` - Apply texture
- `SetObjectNormalmap(id, path)` / `SetObjectRoughness(id, value)` / `SetObjectMetallic(id, value)` / `SetObjectEmissive(id, r, g, b)` - PBR material
- `SetObjectShader(id, shaderID)` - Store shader for object (custom draw)
- `SetObjectWireframe(id, onOff)` - Wireframe mode
- `SetObjectCollision(id, onOff)` - Collision flag

### Scene
- `Start3D` / `End3D` - BeginMode3D / EndMode3D
- `DrawGrid(size, spacing)` - Draw grid
- `Clear(r, g, b)` - ClearBackground
- `BackgroundColor(r, g, b)` - Alias for Clear
- `SetClearColor(r, g, b)` - Set unified renderer clear color
- `SetVsync(onOff)` - Enable/disable VSync (call before InitWindow)
- `SetFramerate(cap)` - Target framerate (0 = uncapped)

### Input
- `KeyDown(key)` / `KeyHit(key)` / `KeyUp(key)`
- `MouseX` / `MouseY` / `MouseMoveX` / `MouseMoveY`
- `MouseButtonDown(btn)` / `MouseButtonHit(btn)` / `MouseButtonUp(btn)`
- `MouseClick(button)` - Alias for click this frame
- `GamepadAxis(pad, axis)` / `GamepadButton(pad, button)` - Gamepad input
- `HideMouse` / `ShowMouse` / `LockMouse` / `UnlockMouse`

### Time & Math
- `DeltaTime` - GetFrameTime
- `FPS` - GetFPS
- `Clamp(value, min, max)` / `Lerp(a, b, t)` / `RandomRange(min, max)`

### Game Loop
- `StartDraw` / `EndDraw` - BeginDrawing / EndDrawing
- `Sync` / `SYNC` - End frame (when UseUnifiedRenderer: full frame)
- `UseUnifiedRenderer` - Enable unified pipeline (3D→2D→GUI, auto water/terrain)
- `SetUseUnifiedRenderer`(enabled) - Toggle unified renderer
- `DeltaTime` - Scaled frame delta (GetFrameTime * TimeScale)
- `FixedDeltaTime` - Fixed physics timestep (1/60)
- `SetTimeScale`(value) / `GetTimeScale` - Time multiplier
- `GetFrameCounter` - Frame count
- `LASTERROR$` - Last error message
- `MakeCamera`(id) / `MAKE CAMERA` - Create camera
- `DeleteCamera`(id) / `RotateCamera`(id, pitch, yaw, roll) / `AttachCameraToObject`(camID, objID)
- `PositionCamera`(id,x,y,z) / `POSITION CAMERA` - Set camera position
- `POINT CAMERA`(id,tx,ty,tz) - Set camera target
- `SetCameraActive`(id) / `SET CAMERA ACTIVE` - Use camera for 3D
- `EscapeKey` - IsKeyDown(KEY_ESCAPE)

### FPS Camera
- `FpsCameraOn` / `FpsCameraOff`
- `FpsCameraPosition(x, y, z)`
- `FpsMoveSpeed(value)` / `FpsLookSpeed(value)`
- `FpsUpdate` - Update camera from mouse/WASD

### Audio (Sounds)
- `LoadSound(id, path)` / `PlaySound(id)` / `StopSound(id)` / `PauseSound(id)` / `SetSoundVolume(id, value)` / `SetSoundLoop(id, onOff)`

## Textures (dbp_textures.go)
- `LoadTexture(id, path)` - Load texture at id
- `DeleteTexture(id)` - Unload and remove
- `TextureExists(id)` - Returns 1 if exists, 0 otherwise
- `SetTextureFilter(id, mode)` / `SetTextureWrap(id, mode)`

## Materials (dbp_materials.go)
- `MakeMaterial(id)` - Create material
- `SetMaterialColor(id, r, g, b)` / `SetMaterialTexture(id, textureID)`
- `ApplyMaterial(id, objectID)` - Apply to object
- `DeleteMaterial(id)` - Unload and remove
- `MaterialExists(id)` - Returns 1 if exists, 0 otherwise

## Camera Extras (dbp_camera.go)
- `DeleteCamera(id)` / `RotateCamera(id, pitch, yaw, roll)` / `AttachCameraToObject(camID, objID)`
- `CameraFollow(objectID, distance)` - Follow object
- `CameraOrbit(x, y, z, angle, pitch, distance)` - Orbit around point
- `CameraShake(amount, duration)` - Screen shake
- `CameraSmooth(value)` - Lerp factor
- `CameraUpdate` - Apply all camera extras (call each frame)

## World (dbp_world.go)
- `SetSkybox(path)` - Load skybox (fallback: solid color if missing)
- `SetSkyboxCubemap(right, left, top, bottom, front, back)` - 6-face cubemap
- `SetAmbientLight(r, g, b)` - Ambient color
- `SetFog(onOff)` / `SetFogOff` / `SetFogColor(r, g, b)` / `SetFogRange(near, far)`
- **Clouds:** `SetCloudsOn` / `SetCloudsOff` / `SetCloudTexture(path)` / `SetCloudSpeed(value)` / `SetCloudDensity(value)` / `SetCloudHeight(value)` / `SetCloudColor(r, g, b)`
- **Sun:** `SetSunDirection(x, y, z)` / `SetSunColor(r, g, b)` / `SetSunIntensity(value)`
- **Time:** `SetWorldTime(hours)` / `GetWorldTime()` / `SetWorldTimeScale(value)` / `SetWeatherPreset(name)`

## Lifecycle Commands (Delete, Hide, Clone, Exists)

All entity types support lifecycle operations. Use `XExists(id)` to check validity before use; `IsNull(value)` for optional returns (e.g. NetConnect).

| Entity | Delete | Hide/Show | Clone | Exists |
|--------|--------|-----------|-------|--------|
| Object | DeleteObject | HideObject, ShowObject | CloneObject, CopyObject | ObjectExists |
| Light | DeleteLight | - | - | LightExists |
| Texture | DeleteTexture | - | - | TextureExists |
| Material | DeleteMaterial | - | - | MaterialExists |
| Water | DeleteWater | HideWater, ShowWater | CloneWater | WaterExists |
| Terrain | DeleteTerrain | HideTerrain, ShowTerrain | CloneTerrain | TerrainExists |
| Group | DeleteGroup | HideGroup, ShowGroup | - | GroupExists |
| Instance | DeleteInstance, DeleteAllInstances | - | - | InstanceExists |
| SpriteObject2D | DeleteSpriteObject | HideSpriteObject, ShowSpriteObject | CloneSpriteObject | SpriteObjectExists |
| Tilemap | DeleteTilemap | HideTilemap, ShowTilemap | - | TilemapExists |
| Spritesheet | DeleteSpritesheet | - | CloneSpritesheet | SpritesheetExists |
| Particles2D | DeleteParticles2D | - | - | Particles2DExists |
| Font | DeleteFont | - | - | FontExists |
| Music | DeleteMusic | - | - | MusicExists |
| Sound | DeleteSound | - | - | SoundExists |
| Mesh | DeleteMesh | - | - | MeshExists |
| Prefab | DeletePrefab | - | - | PrefabExists |

## Water (dbp_water.go)
- `MakeWater(id, width, depth)` - Create water plane
- `SetWaterTexture(id, path)` / `PositionWater(id, x, y, z)` / `SetWaterLevel(id, height)` / `SetWaterColor(id, r, g, b)`
- `SetWaterScroll(id, uSpeed, vSpeed)` - UV scroll
- `SetWaterWave(id, strength)` / `SetWaterWaveStrength(id, value)` / `SetWaterWaveSpeed(id, value)`
- `SetWaterReflection(id, onOff)` / `SetWaterRefraction(id, onOff)`
- `SetWaterNormalmap(id, path)` / `SetWaterFoamTexture(id, path)` / `SetWaterDepthColor(id, r, g, b)` / `SetWaterShallowColor(id, r, g, b)`
- `DrawWater(id)` - Draw at stored position
- `DeleteWater(id)` / `HideWater(id)` / `ShowWater(id)` / `CloneWater(newID, sourceID)` / `WaterExists(id)`

## Terrain (dbp_terrain.go)
- `MakeTerrain(id, width, depth)` - Create flat terrain
- `LoadHeightmap(id, path, width, depth, heightScale)` - From file (fallback: flat if missing)
- `SetTerrainTexture(id, path)` / `PositionTerrain(id, x, y, z)`
- `SetTerrainLayer(id, layerIndex, path)` / `SetTerrainSplatmap(id, path)`
- `GenerateTerrainNoise(id, seed, octaves, scale)` - Procedural (deterministic)
- `DrawTerrain(id)` - Draw at stored position
- `DeleteTerrain(id)` / `HideTerrain(id)` / `ShowTerrain(id)` / `CloneTerrain(newID, sourceID)` / `TerrainExists(id)`

See [docs/WORLD_WATER_TERRAIN.md](WORLD_WATER_TERRAIN.md) for full reference and safety rules.

## Groups (dbp_groups.go)
- `MakeGroup(id)` / `AddToGroup(groupID, objectID)` / `RemoveFromGroup(groupID, objectID)`
- `DeleteGroup(id)` / `HideGroup(id)` / `ShowGroup(id)` / `GroupExists(id)`
- `PositionGroup(groupID, x, y, z)` / `RotateGroup(groupID, pitch, yaw, roll)`
- `DrawGroup(groupID)` / `SyncGroup(groupID)`

## Players (dbp_players.go)
- `MakePlayer(id)` / `SetPlayerPosition(id, x, y, z)` / `SetPlayerAngle(id, pitch, yaw, roll)`
- `MovePlayer(id, x, y, z)` / `TurnPlayer(id, pitch, yaw, roll)` / `SyncPlayer(id)`

## Audio Expanded (dbp_audio.go)
- `LoadMusic(id, path)` / `PlayMusic(id)` / `StopMusic(id)`
- `SetMusicVolume(id, value)` / `SetMusicLoop(id, onOff)`

## Physics (dbp_physics.go)
- **3D (Bullet):** `PhysicsOn` / `PhysicsOff` - Enable/disable default world
- `SetAngularVelocity(id, x, y, z)` - Set rigid body angular velocity
- `SetGravity(x, y, z)` - Set 3D gravity
- `PhysicsStep(dt)` - Step all 3D worlds
- `MakeRigidBody(bodyId$, x, y, z, mass)` - Create sphere rigid body
- `MakeStaticBody(bodyId$, x, y, z, sizeX, sizeY, sizeZ)` - Create static box
- `GetVelocityX/Y/Z(bodyId$)` - Get velocity (use bullet's ApplyForce/ApplyImpulse)
- `GetPositionX/Y/Z(bodyId$)` - Get position
- **2D (Box2D):** `PhysicsOn2D(gx?, gy?)` / `PhysicsOff2D`
- `SetGravity2D(x, y)` / `PhysicsStep2D(dt)`
- `MakeRigidBody2D(bodyId$, x, y, w, h, density)` / `MakeStaticBody2D(bodyId$, x, y, w, h)`
- `GetVelocityX2D/Y2D(bodyId$)` - Use ApplyForce2D/ApplyImpulse2D with world "default"

## Collision (dbp_collision.go)
- `Raycast(ox, oy, oz, dx, dy, dz [, maxDist])` - 3D raycast; use RayHitX/Y/Z, RayHitBody for results
- `Raycast2D(ox, oy, dx, dy)` - 2D raycast; use RayHitX2D, RayHitY2D, RayHitBody2D
- `Spherecast(ox, oy, oz, dx, dy, dz, radius)` - Thick ray (uses raycast internally)
- `ObjectCollides(idA, idB)` - DBP objects AABB overlap (collision flag must be set)
- `PointInObject(objectId, x, y, z)` - Point inside object AABB
- `BodyCollides(bodyIdA$, bodyIdB$)` - Physics body collision (bullet)

## Particles (dbp_particles.go)
- `MakeParticles(id)` - Create particle system (integer id)
- `SetParticleColor(id, r, g, b, a)` / `SetParticleSize(id, size)` / `SetParticleSpeed(id, vx, vy, vz)`
- `SetParticleLifetime(id, seconds)` - Default lifetime
- `EmitParticlesAt(id, x, y, z [, count])` - Emit at position (count default 10)
- `DrawParticles(id)` - Update and draw (call in 3D mode)
- `DeleteParticles(id)` - Remove system

## Networking (net package + dbp_net.go)
- `NetConnect(ip$, port)` - Connect, returns connectionId$
- `NetSend(connectionId$, data$)` / `NetReceive(connectionId$)` - Send/receive text
- `NetDisconnect(connectionId$)` - Close connection
- `NetIsServer()` - 1 if hosting, 0 otherwise
- `NetPlayerID()` - First connection id (or "")
- `NetPing(connectionId$)` / `NetLatency(connectionId$)` - RTT in ms
- `Host(port)` / `Accept(serverId$)` - Server API

## File I/O (VM + dbp_file.go)
- `OpenFile(path, mode)` / `ReadLine(handle)` / `WriteLine(handle, text)` / `CloseFile(handle)` - VM built-ins
- `ReadByte(handle)` / `WriteByte(handle, value)` - Byte I/O
- `SaveString(path, text)` - Write string to file
- `LoadString(path)` - Read file as string
- `SaveValue(path, value)` - Save number or string
- `LoadValue(path)` - Load (returns number if parseable, else string)

## Runtime (dbp_runtime.go + language)
- `WaitFrames(n)` - Language construct: yield for n frames (~n/60 sec) in coroutines
- `StopTask(name)` / `PauseTask(name)` / `ResumeTask(name)` - Stubs (need VM fiber name tracking)

## Replication (game package)
- `ReplicatePosition(entityId$)` / `ReplicateRotation(entityId$)` / `ReplicateScale(entityId$)`
- `ReplicateValue(entityId$, varName$)` - Alias for ReplicateVariable

## Lighting (dbp_lighting.go)
- `MakeLight(id, type)` / `PositionLight(id, x, y, z)` / `RotateLight(id, pitch, yaw, roll)`
- `SetLightColor(id, r, g, b)` / `SetLightIntensity(id, value)` / `SetLightRange(id, value)` / `SetLightAngle(id, degrees)` (spot)
- `EnableShadows(id)` / `DisableShadows(id)` - Per-light shadow flag
- `DeleteLight(id)` / `SyncLight(id)`

Note: Raylib has no built-in dynamic lights. The light registry stores state; visual effect requires a custom shader.

## 2D Game API (dbp_2d.go)

**Full reference:** See [docs/2D_GAME_API.md](2D_GAME_API.md) for the complete 2D API with examples.

### Drawing
- `DrawPixel(x, y, r, g, b)` / `DrawRectOutline(x, y, w, h, r, g, b)` / `DrawCircleOutline(x, y, radius, r, g, b)`
- `DrawTriangle(x1,y1, x2,y2, x3,y3, r,g,b)`

### Sprites
- `DeleteSprite(id)` / `SetSpriteColor(id, r, g, b, a)` / `DrawSpriteRotated(id, x, y, angle)` / `DrawSpriteScaled(id, x, y, sx, sy)` / `DrawSpriteTint(id, x, y, r, g, b)`

### Spritesheets
- `LoadSpritesheet(id, pngPath, jsonPath)` or `(id, path, frameW, frameH)` / `PlaySpriteAnimation(id, tagName, speed)` / `GetSliceRect(id, sliceName)` / `GetAnimationLength(id, tagName)`
- `DrawSpriteFrame(id, frame, x, y)` / `AnimateSprite(id, startFrame, endFrame, speed)`

### Tilemaps
- `LoadTilemap(id, path)` / `DrawTilemap(id)` / `SetTile(id, x, y, tileIndex)` / `GetTile(id, x, y)`

### 2D Camera
- `Camera2DOn` / `Camera2DOff` / `Camera2DPosition(x, y)` / `Camera2DZoom(value)` / `Camera2DRotation(angle)` / `Camera2DFollow(objectId)`

### 2D Collision
- `RectCollides(x1,y1,w1,h1, x2,y2,w2,h2)` / `PointInRect(x,y, rx,ry,rw,rh)` / `CircleCollides(...)` / `PointInCircle(...)`

### 2D Physics
- `Physics2DOn` / `Physics2DOff` / `MakeBody2D(id, mass)` / `MakeStatic2D(id)`
- `SetBody2DPosition(id, x, y)` / `SetBody2DVelocity(id, vx, vy)` / `ApplyForce2D(id, fx, fy)` / `ApplyImpulse2D(id, ix, iy)`
- `GetBody2DX(id)` / `GetBody2DY(id)` / `GetBody2DVX(id)` / `GetBody2DVY(id)`

### 2D Objects (SpriteObject2D)
- `MakeSpriteObject(id, spriteId)` / `PositionObject2D(id, x, y)` / `MoveObject2D(id, dx, dy)` / `RotateObject2D(id, angle)` / `ScaleObject2D(id, sx, sy)`
- `DrawObject2D(id)` / `SyncObject2D(id)`

### UI
- `UITextbox(id, x, y, w, h)` - Editable text box; returns current text

### 2D Math
- `AngleBetween2D(x1,y1, x2,y2)` / `Distance2D(x1,y1, x2,y2)` / `Normalize2D(x, y)` / `Dot2D(x1,y1, x2,y2)`

### 2D Particles
- `MakeParticles2D(id, maxCount)` / `SetParticles2DColor(id, r, g, b)` / `SetParticles2DSize(id, size)` / `SetParticles2DSpeed(id, speed)`
- `EmitParticles2D(id, count [, x, y])` / `DrawParticles2D(id)`

## 3D Game API (dbp_3d.go)

**Full reference:** See [docs/3D_GAME_API.md](3D_GAME_API.md) for the complete 3D API with examples.

### Window
- `Window(width, height, title)` / `CloseWindow` / `SetTargetFPS(value)`

### Camera Queries
- `GetCameraX` / `GetCameraY` / `GetCameraZ` / `GetCameraPitch` / `GetCameraYaw`
- `PointCameraAt(x, y, z)` / `SetCameraFOV(value)` / `SetCameraUp(x, y, z)`

### 3D Object Creation
- `MakeCylinder(id, radius, height)` / `MakeGrid(id, size, spacing)`
- `LoadObject(id, path)` - DBP arg order (id first)

### 3D Object Queries
- `GetObjectX(id)` / `GetObjectY(id)` / `GetObjectZ(id)`
- `GetObjectPitch(id)` / `GetObjectYaw(id)` / `GetObjectRoll(id)`
- `GetObjectScaleX(id)` / `GetObjectScaleY(id)` / `GetObjectScaleZ(id)`

### Parenting & Tags
- `ParentObject(childID, parentID)` / `UnparentObject(id)`
- `SetObjectTag(id, tag)` / `GetObjectTag(id)`

### Replication
- `SyncObject(id)` / `UnsyncObject(id)` / `SetObjectOwner(id, playerID)` / `GetObjectOwner(id)`

### 3D Physics (Collision Shapes, Rigid Body)
- `MakeBoxCollider(id, sx, sy, sz)` / `MakeSphereCollider(id, radius)` / `MakeCapsuleCollider(id, radius, height)`
- `MakeMeshCollider(id, meshID)` - Unsupported in the shipped 3D fallback; returns an explicit error
- `SetBodyPosition(bodyId$, x, y, z)` / `GetBodyPosition(bodyId$)` - Default-world vector helpers for string body IDs
- `SetBodyVelocity(bodyId$, vx, vy, vz)` / `GetBodyVelocity(bodyId$)` - Default-world velocity helpers for string body IDs
- `MakeRigidBodyId(id, x, y, z, mass)` - Create rigid body with int ID
- `SetRigidBodyPosition(id, x, y, z)` / `SetRigidBodyVelocity(id, vx, vy, vz)`
- `GetBodyX/Y/Z(bodyId$)` / `GetBodyVX/VY/VZ(bodyId$)` - Default-world queries for string body IDs
- `GetRigidBodyX/Y/Z(id)` / `GetRigidBodyMass(id)` / `GetRigidBodySpeed(id)` / `GetRigidBodyAngularVelocity(id)`

### Level Loading (dbp_level.go)
- `LoadLevel(id, path)` / `DrawLevel(id)` / `UnloadLevel(id)`
- `LoadLevelCollision(id)` - Create physics colliders from level; returns count
- `GetLevelColliderCount(id)` / `GetLevelCollider(id, index)` - Collider queries
- `GetLevelObjectCount(id)` / `GetLevelObject(id, index)` - Object queries

### Prefab (dbp_prefab.go)
- `LoadPrefab(id, path)` - Load prefab template
- `SpawnPrefab(id, x, y, z)` - Instantiate prefab at position; returns object ID

### IK (dbp_ik.go)
- `IKEnable(objectID, onOff)` - Enable/disable IK for object
- `IKSolveTwoBone(objectID, boneA$, boneB$, targetX, targetY, targetZ)` - Experimental two-bone IK solve request

### Model / Mesh / Animation
- `LoadMesh(id, path)` / `GetModelBounds(objectID)` / `GetMeshVertexCount(id)` / `GetMeshTriangleCount(id)`
- `LoadAnimation(id, path)` / `PlayAnimation(objectID, animID, clipIndex, speed)` / `SetAnimationFrame(objectID, frame)` / `GetAnimationLength(animID, clipIndex)`
- `GetAnimationFrame(objectID)` / `GetAnimationLength(animID)` / `GetAnimationName(animID)` (graceful when no anim)

### Light / Skybox Queries
- `GetLightX/Y/Z(id)` / `GetLightColorR/G/B(id)`
- `GetSkybox()` / `GetAmbientLightR/G/B()`

### Terrain
- `SetTerrainHeight(terrainId, x, z, height)` / `PaintTerrainTexture(terrainId, x, z, layer)`

### 3D Particles
- `MakeParticles3D(id, maxCount)` / `SetParticles3DColor/Size/Speed(id, ...)` / `EmitParticles3D(id, count [, x, y, z])`
- `DrawParticles3D(id)` / `GetParticles3DCount(id)` / `GetParticles3DMax(id)`

### Instancing
- `MakeInstance(baseID, instanceID)` / `PositionInstance(instanceID, x, y, z)` / `DrawInstances(baseID)`

### Pathfinding
- `NavMeshLoad(id, path)` / `NavMeshFindPath(id, startX, startY, startZ, endX, endY, endZ)` / `NavMeshDraw(id)`

### Matrix / Quaternion (stubs)
- `MakeQuaternion(id, pitch, yaw, roll)` / `RotateObjectQuat(id, quatID)` / `GetObjectMatrix(id)`

### Runtime
- `FixedUpdate(rate)` / `OnFixedUpdate(label$)`

### 3D Math
- `Distance3D(x1,y1,z1, x2,y2,z2)` / `AngleBetween3D(...)` / `Normalize3D(x,y,z)` / `Dot3D(...)` / `Cross3D(...)`

## 2D Drawing Aliases (dbp.go)
- `DrawRect(x, y, w, h, r, g, b)` / `DrawCircle(x, y, radius, r, g, b)` / `DrawLine(x1, y1, x2, y2, r, g, b)`
- `DrawSprite(id, x, y)` / `LoadSprite(id, path)` - Same as Sprite/LoadImage

## Text
- `DrawText(text$, x, y, size)` or `DrawText(text$, x, y, size, r, g, b)`
- `LoadFont(id, path)` / `SetFont(id)` - Font registry

## File I/O
- `FileExists(path)` / `AppendFile(path, text)`

## UI
- `UIButton(id, x, y, w, h, text$)` / `UILabel(x, y, text$)` / `UICheckbox(id, x, y, text$)` / `UISlider(id, x, y, w, min, max, value)`
- `LabelAt(x, y, text)` / `ButtonAt(x, y, w, h, text)` / `CheckboxAt(x, y, label, checked)` / `SliderAt(x, y, width, min, max, value)` - Explicit layout

## Debug
- `DebugLog(text$)` / `DebugDrawLine(x1,y1,z1, x2,y2,z2)` / `DebugDrawBox(x, y, z, size)`

## Math
- `Randomize(seed)` / `RandomMinMax(a, b)` / `Distance(x1,y1,z1, x2,y2,z2)` / `AngleBetween(...)` / `Dot(...)`

## Coroutines
- `StartTask SubName()` - Alias for StartCoroutine
- `Wait(seconds)` - Alias for WaitSeconds when followed by parentheses
- `Yield` - Already exists
