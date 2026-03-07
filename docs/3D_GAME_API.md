# 3D Game API Reference

Complete reference for the CyberBASIC2 3D game API: FPS, RPGs, sandbox, survival, and multiplayer worlds. All commands use PascalCase and integer IDs where applicable.

**Multiplayer-safe?** Commands marked with ✓ are safe to use in networked games. Camera and UI are local-only. Physics and object transforms need replication for sync.

---

## 1. Overview

The 3D API is organized across multiple modules:

- **dbp.go:** Core 3D objects, FPS camera, scene, input
- **dbp_3d.go:** Window aliases, camera queries, object creation (MakeCylinder, MakeGrid), object queries, parenting, tags, replication, 3D math
- **dbp_camera.go:** CameraFollow, CameraOrbit, CameraShake, CameraSmooth
- **dbp_groups.go:** Object groups
- **dbp_physics.go:** Physics (Bullet 3D)
- **dbp_collision.go:** Raycast, Spherecast, ObjectCollides, PointInObject
- **dbp_lighting.go:** Lights
- **dbp_materials.go:** Materials
- **dbp_world.go:** Skybox, ambient light, fog

---

## 2. Window & Rendering

| Command | Args | Description |
|---------|------|-------------|
| `Window` | (width, height, title) | DBP alias for InitWindow |
| `CloseWindow` | () | Close the window |
| `SetTargetFPS` | (value) | Target frame rate |
| `Clear` | (r, g, b) | Clear background |
| `StartDraw` | () | BeginDrawing |
| `EndDraw` | () | EndDrawing |
| `Start3D` | () | BeginMode3D |
| `End3D` | () | EndMode3D |

**When using UseUnifiedRenderer:** Call **SYNC** instead of Start3D/End3D. The engine handles 3D→2D→GUI order automatically. Use `SYNC` at the end of `OnDraw()` for DBP-style programs.

**Example:**
```basic
Window 1280, 720, "My 3D Game"
SetTargetFPS 60
While Not WindowShouldClose()
  StartDraw
  Clear 30, 30, 50
  Start3D
  DrawObject 1
  DrawGrid 20, 1
  End3D
  EndDraw
Wend
CloseWindow
```

---

## 3. Camera (Standard)

| Command | Args | Description |
|---------|------|-------------|
| `SetCameraPosition` | (x, y, z) | Set camera position |
| `SetCameraTarget` | (x, y, z) | Set camera target |
| `PointCamera` | (x, y, z) | Alias for SetCameraTarget |
| `PointCameraAt` | (x, y, z) | Set camera to look at point |
| `SetCameraFOV` | (value) | Field of view (degrees) |
| `SetCameraRange` | (near, far) | Projection planes (no-op) |
| `SetCameraUp` | (x, y, z) | Camera up vector |

**Example:**
```basic
SetCameraPosition 0, 5, 10
PointCameraAt 0, 0, 0
SetCameraFOV 60
```

---

## 4. Camera Queries

| Command | Args | Description |
|---------|------|-------------|
| `GetCameraX` | () | Camera position X |
| `GetCameraY` | () | Camera position Y |
| `GetCameraZ` | () | Camera position Z |
| `GetCameraPitch` | () | Camera pitch (degrees) |
| `GetCameraYaw` | () | Camera yaw (degrees) |

**Multiplayer-safe?** ✓ Yes (local only).

**Example:**
```basic
cx = GetCameraX
cy = GetCameraY
cz = GetCameraZ
```

---

## 5. FPS Camera

| Command | Args | Description |
|---------|------|-------------|
| `FpsCameraOn` | () | Enable FPS camera, disable cursor |
| `FpsCameraOff` | () | Disable FPS camera |
| `FpsCameraPosition` | (x, y, z) | Set FPS camera position |
| `FpsMoveSpeed` | (value) | Movement speed |
| `FpsLookSpeed` | (value) | Mouse look sensitivity |
| `FpsUpdate` | () | Update camera from WASD + mouse |

**Multiplayer-safe?** ✓ Yes (camera is local only).

**Example:**
```basic
FpsCameraOn
FpsCameraPosition 0, 2, 10
FpsMoveSpeed 5
FpsLookSpeed 0.002
' In loop:
FpsUpdate
```

---

## 6. 3D Object Creation

| Command | Args | Description |
|---------|------|-------------|
| `MakeCube` | (id, size) | Procedural cube |
| `MakeSphere` | (id, radius) | Procedural sphere |
| `MakeCylinder` | (id, radius, height) | Procedural cylinder |
| `MakePlane` | (id, width, height) | Procedural plane |
| `MakeGrid` | (id, size, spacing) | Procedural grid |
| `LoadObject` | (id, path) | Load model from file (GLTF, OBJ) |
| `LoadPrefab` | (id, path) | Load prefab template |
| `SpawnPrefab` | (id, x, y, z) | Instantiate prefab at position; returns object ID |
| `DeleteObject` | (id) | Unload and remove |
| `CloneObject` | (newID, sourceID) | Duplicate object |
| `CopyObject` | (newID, sourceID) | Alias for CloneObject |
| `ObjectExists` | (id) | Returns 1 if exists, 0 otherwise |

**Example:**
```basic
MakeCube 1, 2
MakeSphere 2, 1.5
LoadObject 4, "character.gltf"
```

---

## 6b. Level Loading

| Command | Args | Description |
|---------|------|-------------|
| `LoadLevel` | (id, path) | Load full level (meshes, materials, textures, and imported light data when present) |
| `LoadLevelWithHierarchy` | (id, path) | Load level with mesh-node hierarchy preserved through the imported node tree |
| `DrawLevel` | (id) | Draw all level objects |
| `UnloadLevel` | (id) | Free level resources |
| `LoadLevelCollision` | (id) | Create physics colliders from level; returns count |
| `GetLevelColliderCount` | (id) | Collider count |
| `GetLevelCollider` | (id, index) | Physics body ID at index |
| `GetLevelObjectCount` | (id) | Object count |
| `GetLevelObject` | (id, index) | Get object ID at index |

**LoadLevel** builds meshes, materials, textures, and DBP objects automatically. Call **LoadLevelCollision** after loading to enable physics colliders. Colliders are detected from GLTF node names (`col_*`, `collision_*`, `Collision*`). Imported GLTF punctual lights are still not a dependable shipping path, so prefer explicit `MakeLight` setup for production scenes.

**Example:**
```basic
LoadLevel 1, "castle.gltf"
LoadLevelCollision 1
PhysicsOn
While Not WindowShouldClose()
  StartDraw
  Clear 30, 30, 50
  Start3D
  DrawLevel 1
  End3D
  EndDraw
Wend
UnloadLevel 1
```

See [LEVEL_LOADING.md](LEVEL_LOADING.md) for full documentation.

---

## 7. 3D Object Transform

| Command | Args | Description |
|---------|------|-------------|
| `PositionObject` | (id, x, y, z) | Set position |
| `RotateObject` | (id, pitch, yaw, roll) | Set rotation |
| `ScaleObject` | (id, sx, sy, sz) | Set scale |
| `MoveObject` | (id, dx, dy, dz) | Add to position |
| `TurnObject` | (id, dpitch, dyaw, droll) | Add to rotation |
| `YRotateObject` | (id, angle) | Add to yaw |

---

## 8. 3D Object Queries

| Command | Args | Description |
|---------|------|-------------|
| `GetObjectX` | (id) | Object X position |
| `GetObjectY` | (id) | Object Y position |
| `GetObjectZ` | (id) | Object Z position |
| `GetObjectPitch` | (id) | Object pitch |
| `GetObjectYaw` | (id) | Object yaw |
| `GetObjectRoll` | (id) | Object roll |
| `GetObjectScaleX` | (id) | Scale X |
| `GetObjectScaleY` | (id) | Scale Y |
| `GetObjectScaleZ` | (id) | Scale Z |

**Multiplayer-safe?** ✓ Yes (read-only).

**Example:**
```basic
x = GetObjectX 1
y = GetObjectY 1
z = GetObjectZ 1
```

---

## 9. Scene Graph and Object Parenting

Objects form a hierarchy: a child's world transform is its local transform composed with the parent's full position, rotation, and scale. `LoadLevelWithHierarchy` preserves transform-only GLTF nodes when rebuilding mesh-node parenting, so nested offsets and rotations now survive import correctly.

| Command | Args | Description |
|---------|------|-------------|
| `ParentObject` | (childID, parentID) | Attach child to parent |
| `AttachObject` | (childID, parentID) | Alias for ParentObject |
| `UnparentObject` | (id) | Detach from parent |
| `DetachObject` | (id) | Alias for UnparentObject |
| `GetObjectLocalX` | (id) | Local X (before parent composition) |
| `GetObjectLocalY` | (id) | Local Y |
| `GetObjectLocalZ` | (id) | Local Z |
| `GetObjectLocalPitch` | (id) | Local pitch |
| `GetObjectLocalYaw` | (id) | Local yaw |
| `GetObjectLocalRoll` | (id) | Local roll |
| `GetObjectLocalScaleX` | (id) | Local scale X |
| `GetObjectLocalScaleY` | (id) | Local scale Y |
| `GetObjectLocalScaleZ` | (id) | Local scale Z |
| `SetObjectLocalTransform` | (id, x,y,z, pitch,yaw,roll, sx,sy,sz) | Set local position, rotation, scale |

**World vs local:** `GetObjectX/Y/Z` return world position (after parent composition). `GetObjectLocalX/Y/Z` return the object's own position. Use local transforms when editing a child relative to its parent.

**Example:**
```basic
ParentObject 2, 1
' Object 2 now follows object 1's transform
x = GetObjectLocalX 2   ' child's offset from parent
SetObjectLocalTransform 2, 0, 1, 0, 0, 0, 0, 1, 1, 1  ' 1 unit above parent
UnparentObject 2
```

---

## 10. Object Tags

| Command | Args | Description |
|---------|------|-------------|
| `SetObjectTag` | (id, tag) | Set string tag |
| `GetObjectTag` | (id) | Get tag |

**Example:**
```basic
SetObjectTag 1, "player"
tag$ = GetObjectTag 1
```

---

## 11. Model / Mesh

| Command | Args | Description |
|---------|------|-------------|
| `LoadMesh` | (id, path) | Load model as mesh |
| `GetModelBounds` | (objectID) | Returns [minX,minY,minZ, maxX,maxY,maxZ] |
| `GetMeshVertexCount` | (id) | Vertex count of first mesh |
| `GetMeshTriangleCount` | (id) | Triangle count of first mesh |

---

## 12. Animation

| Command | Args | Description |
|---------|------|-------------|
| `LoadAnimation` | (id, path) | Load all animation clips from file |
| `PlayAnimation` | (objectID, animID, speed) or (objectID, animID, clipIndex, speed) | Start playing animation |
| `StopAnimation` | (objectID) | Stop animation |
| `SetAnimationFrame` | (objectID, frame) | Set current frame |
| `SetAnimationSpeed` | (objectID, speed) | Playback speed |
| `SetAnimationLoop` | (objectID, onOff) | Loop on/off |
| `GetAnimationFrame` | (objectID) | Current frame |
| `GetAnimationLength` | (animID) or (animID, clipIndex) | Frame count |
| `GetAnimationName` | (animID) or (animID, clipIndex) | Clip name |
| `CrossfadeAnimation` | (objectID, animID, clipIndex, duration) | Blend to new clip over duration seconds |
| `SetAnimationBlend` | (objectID, weight) | Manual blend weight 0–1 |

**CrossfadeAnimation** now blends two skeletal poses over `duration` seconds. `SetAnimationBlend` lets you scrub that same crossfade manually by setting a 0-1 weight.

**Graceful:** LoadAnimation succeeds even when the file has no animations. GetAnimationLength returns 0. PlayAnimation and SetAnimationFrame no-op when anim not found.

**Blender workflow:** See [BLENDER_WORKFLOW.md](BLENDER_WORKFLOW.md).

---

## 13. Drawing Objects

| Command | Args | Description |
|---------|------|-------------|
| `DrawObject` | (id) | Draw object (call between Start3D/End3D) |
| `HideObject` | (id) | Set visible=false |
| `ShowObject` | (id) | Set visible=true |

---

## 13b. Lifecycle Commands (Delete, Hide, Clone, Exists)

| Entity | Delete | Hide/Show | Clone | Exists |
|--------|--------|-----------|-------|--------|
| `Object` | DeleteObject | HideObject, ShowObject | CloneObject, CopyObject | ObjectExists |
| `Light` | DeleteLight | - | - | LightExists |
| `Texture` | DeleteTexture | - | - | TextureExists |
| `Material` | DeleteMaterial | - | - | MaterialExists |
| `Water` | DeleteWater | HideWater, ShowWater | CloneWater | WaterExists |
| `Terrain` | DeleteTerrain | HideTerrain, ShowTerrain | CloneTerrain | TerrainExists |
| `Group` | DeleteGroup | HideGroup, ShowGroup | - | GroupExists |
| `Instance` | DeleteInstance, DeleteAllInstances | - | - | InstanceExists |
| `Mesh` | DeleteMesh | - | - | MeshExists |
| `Prefab` | DeletePrefab | - | - | PrefabExists |
| `Particles` | DeleteParticles | - | - | - |


---

## 14. Object Groups

| Command | Args | Description |
|---------|------|-------------|
| `MakeGroup` | (id) | Create empty group |
| `AddToGroup` | (groupID, objectID) | Add object |
| `RemoveFromGroup` | (groupID, objectID) | Remove object |
| `PositionGroup` | (groupID, x, y, z) | Set all positions |
| `RotateGroup` | (groupID, pitch, yaw, roll) | Set all rotations |
| `DrawGroup` | (groupID) | Draw all objects |
| `SyncGroup` | (groupID) | Sync for multiplayer |

---

## 12b. IK (Inverse Kinematics)

| Command | Args | Description |
|---------|------|-------------|
| `IKEnable` | (objectID, onOff) | Enable/disable IK for object |
| `IKSolveTwoBone` | (objectID, boneA$, boneB$, targetX, targetY, targetZ) | Experimental two-bone IK solve request |

Requires object with skeleton (from GLTF skin). Bones are identified by name. Current status: the solver validates skeleton data and computes an internal two-bone solution, but visible skinned-pose application is not fully integrated yet, so treat this as preparatory/experimental rather than finished runtime IK.

---

## 15. 3D Collision / Hit Testing

| Command | Args | Description |
|---------|------|-------------|
| `Raycast` | (ox, oy, oz, dx, dy, dz [, maxDist]) | Ray vs world from origin plus direction |
| `Spherecast` | (ox, oy, oz, dx, dy, dz, radius) | Swept-sphere query against physics body bounds |
| `ObjectCollides` | (id1, id2) | AABB overlap |
| `PointInObject` | (id, x, y, z) | Point in AABB |

---

## 16. 3D Physics

The current 3D physics layer is a Bullet-shaped API backed by the shipped pure-Go fallback, not native Bullet. Use `BulletBackendName()`, `BulletBackendMode()`, `BulletNativeAvailable()`, `BulletJointsAvailable()`, and `BulletFeatureAvailable(feature$)` when you need to branch on runtime capability.

| Command | Args | Description |
|---------|------|-------------|
| `PhysicsOn` | () | Enable 3D physics |
| `PhysicsOff` | () | Disable |
| `MakeRigidBody` | (bodyId$, x, y, z, mass) | Create rigid body (string ID) |
| `MakeRigidBodyId` | (id, x, y, z, mass) | Create rigid body (int ID) |
| `MakeStaticBody` | (bodyId$, x, y, z, sx, sy, sz) | Create static box |
| `MakeBoxCollider` | (id, sx, sy, sz) | Static box collider (int ID) |
| `MakeSphereCollider` | (id, radius) | Static sphere collider |
| `MakeCapsuleCollider` | (id, radius, height) | Static capsule collider |
| `MakeMeshCollider` | (id, meshID) | Mesh collider helper; unsupported in the shipped fallback |
| `ApplyForce` | (id, fx, fy, fz) | Apply force |
| `ApplyImpulse` | (id, ix, iy, iz) | Apply impulse |
| `SetGravity` | (x, y, z) | Set gravity |
| `SetBodyPosition` | (bodyId$, x, y, z) | Set position for default-world string body |
| `GetBodyPosition` | (bodyId$) | Returns `[x, y, z]` for default-world string body |
| `SetBodyVelocity` | (bodyId$, vx, vy, vz) | Set velocity for default-world string body |
| `GetBodyVelocity` | (bodyId$) | Returns `[vx, vy, vz]` for default-world string body |
| `SetRigidBodyPosition` | (id, x, y, z) | Set position (int ID) |
| `SetRigidBodyVelocity` | (id, vx, vy, vz) | Set velocity |
| `GetBodyX` | (bodyId$) | Position X for default-world string body |
| `GetBodyY` | (bodyId$) | Position Y for default-world string body |
| `GetBodyZ` | (bodyId$) | Position Z for default-world string body |
| `GetBodyVX` | (bodyId$) | Velocity X for default-world string body |
| `GetBodyVY` | (bodyId$) | Velocity Y for default-world string body |
| `GetBodyVZ` | (bodyId$) | Velocity Z for default-world string body |
| `GetRigidBodyX` | (id) | Position X |
| `GetRigidBodyY` | (id) | Position Y |
| `GetRigidBodyZ` | (id) | Position Z |
| `GetRigidBodyMass` | (id) | Body mass |
| `GetRigidBodySpeed` | (id) | Velocity magnitude |
| `GetRigidBodyAngularVelocity` | (id) | Angular velocity magnitude |
| `GetVelocityX` | (id) | Velocity X |
| `GetVelocityY` | (id) | Velocity Y |
| `GetVelocityZ` | (id) | Velocity Z |

Notes:

- `MakeSphereCollider(id, radius)` and `MakeCapsuleCollider(id, radius, height)` now accept the documented argument counts.
- `MakeMeshCollider` now returns a clear unsupported-feature error in the current fallback.
- `GetBodyX/Y/Z` and `GetBodyVX/VY/VZ` provide a more 2D-style default-world query path for string body IDs.
- `Spherecast` is still a swept-bounds helper, not exact Bullet geometry.
- Unsupported low-level 3D features such as joints, torque-only APIs, heightmaps, and compounds now return explicit errors instead of silently succeeding.

---

## 17. 3D Math Helpers

| Command | Args | Description |
|---------|------|-------------|
| `Distance3D` | (x1,y1,z1, x2,y2,z2) | Distance |
| `AngleBetween3D` | (x1,y1,z1, x2,y2,z2) | Angle between vectors |
| `Normalize3D` | (x, y, z) | Returns [nx, ny, nz] |
| `Dot3D` | (x1,y1,z1, x2,y2,z2) | Dot product |
| `Cross3D` | (x1,y1,z1, x2,y2,z2) | Cross product |

**Example:**
```basic
d = Distance3D x1, y1, z1, x2, y2, z2
vec = Normalize3D dx, dy, dz
dot = Dot3D x1, y1, z1, x2, y2, z2
```

---

## 18. Replication (Multiplayer)

| Command | Args | Description |
|---------|------|-------------|
| `SyncObject` | (id) | Mark for position sync |
| `UnsyncObject` | (id) | Unmark |
| `SetObjectOwner` | (id, playerID) | Set owner |
| `GetObjectOwner` | (id) | Get owner |
| `ReplicatePosition` | (entityId) | Game package |
| `ReplicateRotation` | (entityId) | Game package |
| `ReplicateScale` | (entityId) | Game package |

---

## 19. Lighting

| Command | Args | Description |
|---------|------|-------------|
| `MakeLight` | (id, type) | Create light |
| `PositionLight` | (id, x, y, z) | Set position |
| `RotateLight` | (id, pitch, yaw, roll) | Set rotation |
| `SetLightColor` | (id, r, g, b) | Set color |
| `SetLightIntensity` | (id, value) | Set intensity |
| `SetLightRange` | (id, value) | Set range |
| `EnableLightShadows` | (id) | Mark light as preferred shadow caster |
| `DisableLightShadows` | (id) | Disable shadow casting preference for this light |
| `EnableShadows` | () or (id) | Enable global shadows, or mark a light as the active caster |
| `DisableShadows` | () or (id) | Disable global shadows, or clear a light's shadow-caster flag |
| `SetShadowQuality` | (name$) | Apply `low`, `medium`/`mid`, or `high` shadow preset |
| `SetShadowMapSize` | (width, height) | Override shadow map resolution |
| `SetShadowBias` | (bias) | Set depth bias for the shadow compare |
| `DeleteLight` | (id) | Remove |
| `LightExists` | (id) | Returns 1 if exists, 0 otherwise |
| `SyncLight` | (id) | Sync for multiplayer |
| `GetLightX` | (id) | Light position X |
| `GetLightY` | (id) | Light position Y |
| `GetLightZ` | (id) | Light position Z |
| `GetLightColorR` | (id) | Light color R (0-255) |
| `GetLightColorG` | (id) | Light color G |
| `GetLightColorB` | (id) | Light color B |

By default, the global shadow system starts enabled. For the easiest setup, create a directional light with `MakeLight(id, 1)` and orient it with `RotateLight`; the first directional light can then act as the default caster automatically. If you want explicit control over which directional light casts shadows, call `EnableShadows(lightId)` or `EnableLightShadows(lightId)`. Call `DisableShadows()` when you want the lowest-cost non-shadow path.

Current scope: the renderer supports one active directional shadow caster in the main pass. `SetShadowQuality("low"|"medium"|"high")` is the easiest way to scale cost for lower-end, mid-range, or higher-end PCs.

---

## 20. Materials & Textures

Imported GLTF PBR values are preserved by default. Use **SetObjectRoughness**, **SetObjectMetallic**, **SetObjectNormalmap**, and **SetObjectEmissive** only when you want per-object overrides on top of the imported material.

| Command | Args | Description |
|---------|------|-------------|
| `MakeMaterial` | (id) | Create material |
| `CreatePBRMaterial` | (id) | Create material with PBR defaults (metallic=0, roughness=0.5) |
| `SetMaterialPBR` | (id, metallic, roughness, normalTextureID) | Set metallic, roughness, normal map |
| `SetMaterialMetalness` | (id, value) | Set material metalness scalar (0-1) |
| `SetMaterialRoughness` | (id, value) | Set material roughness scalar (0-1) |
| `SetMaterialColor` | (id, r, g, b) | Set color |
| `SetMaterialTexture` | (id, textureID) | Set texture |
| `ApplyMaterial` | (id, objectID) | Apply to object |
| `SetObjectRoughness` | (id, value) | Per-object roughness (0-1) |
| `SetObjectMetallic` | (id, value) | Per-object metallic (0-1) |
| `SetObjectNormalmap` | (id, path) or (id, textureID) | Set normal map |
| `SetObjectEmissive` | (id, r, g, b) | Emissive color (for custom shaders) |
| `LoadTexture` | (id, path) | Load texture |
| `DeleteTexture` | (id) | Unload |
| `TextureExists` | (id) | Returns 1 if exists, 0 otherwise |
| `SetTextureFilter` | (id, mode) | Filter mode |
| `SetTextureWrap` | (id, mode) | Wrap mode |

---

## 21. Skybox / Environment / Clouds / Sun / Time

| Command | Args | Description |
|---------|------|-------------|
| `SetSkybox` | (path) | Load skybox (fallback: solid color if missing) |
| `SetAmbientLight` | (r, g, b) | Ambient color |
| `SetFog` | (onOff) | Enable fog |
| `SetFogOff` | () | Disable fog |
| `SetFogColor` | (r, g, b) | Fog color |
| `SetFogRange` | (near, far) | Fog range |
| `GetSkybox` | () | 1 if loaded, 0 otherwise |
| `GetAmbientLightR` | () | Ambient R (0-255) |
| `GetAmbientLightG` | () | Ambient G |
| `GetAmbientLightB` | () | Ambient B |
| `SetCloudsOn` | () | Enable clouds |
| `SetCloudsOff` | () | Disable clouds |
| `SetCloudTexture` | (path) | Cloud texture (fallback: disable if missing) |
| `SetCloudSpeed` | (value) | UV scroll |
| `SetCloudDensity` | (value) | 0-1 |
| `SetCloudHeight` | (value) | Dome height |
| `SetCloudColor` | (r, g, b) | Tint |
| `SetSunDirection` | (x, y, z) | Sun direction |
| `SetSunColor` | (r, g, b) | Sun color |
| `SetSunIntensity` | (value) | Sun intensity |
| `SetWorldTime` | (hours) | 0-24 |
| `GetWorldTime` | () | Current hours |
| `SetWorldTimeScale` | (value) | Time multiplier |
| `SetWeatherPreset` | (name) | "Clear", "Rain", etc. |

---

## 22. Water

| Command | Args | Description |
|---------|------|-------------|
| `MakeWater` | (id, width, depth) | Create water plane |
| `SetWaterTexture` | (id, path) | Texture (fallback: blue if missing) |
| `PositionWater` | (id, x, y, z) | Position |
| `SetWaterLevel` | (id, height) | Y position |
| `SetWaterColor` | (id, r, g, b) | Base color |
| `SetWaterScroll` | (id, uSpeed, vSpeed) | UV scroll |
| `SetWaterWaveStrength` | (id, value) | Wave height |
| `SetWaterWaveSpeed` | (id, value) | Wave speed |
| `SetWaterReflection` | (id, onOff) | Reflection |
| `SetWaterRefraction` | (id, onOff) | Refraction |
| `SetWaterNormalmap` | (id, path) | Normal map |
| `SetWaterFoamTexture` | (id, path) | Foam texture |
| `SetWaterDepthColor` | (id, r, g, b) | Deep water tint |
| `SetWaterShallowColor` | (id, r, g, b) | Shallow water tint |
| `DrawWater` | (id) | Draw at stored position |
| `DeleteWater` | (id) | Remove and unload |
| `HideWater` | (id) | Set visible=false |
| `ShowWater` | (id) | Set visible=true |
| `CloneWater` | (newID, sourceID) | Duplicate water |
| `WaterExists` | (id) | Returns 1 if exists, 0 otherwise |

---

## 23. Terrain

| Command | Args | Description |
|---------|------|-------------|
| `MakeTerrain` | (id, width, depth) | Create flat terrain |
| `LoadHeightmap` | (id, path, width, depth, heightScale) | From file (fallback: flat if missing) |
| `SetTerrainTexture` | (id, path) | Texture (fallback: green if missing) |
| `PositionTerrain` | (id, x, y, z) | Position |
| `SetTerrainLayer` | (id, layerIndex, path) | Layer 0-3 |
| `SetTerrainSplatmap` | (id, path) | Splatmap (fallback: layer 0 if missing) |
| `GenerateTerrainNoise` | (id, seed, octaves, scale) | Procedural (deterministic) |
| `DrawTerrain` | (id) | Draw at stored position |
| `DeleteTerrain` | (id) | Remove and unload |
| `HideTerrain` | (id) | Set visible=false |
| `ShowTerrain` | (id) | Set visible=true |
| `CloneTerrain` | (newID, sourceID) | Duplicate terrain |
| `TerrainExists` | (id) | Returns 1 if exists, 0 otherwise |
| `TerrainGetHeight` | (terrainId, x, z) | Get height at point |
| `TerrainRaise` | (terrainId, x, z, radius, amount) | Raise terrain |
| `TerrainLower` | (terrainId, x, z, radius, amount) | Lower terrain |
| `SetTerrainHeight` | (terrainId, x, z, height) | Flatten to height |
| `PaintTerrainTexture` | (terrainId, x, z, layer) | Paint texture layer |

See [docs/WORLD_WATER_TERRAIN.md](WORLD_WATER_TERRAIN.md) for full reference and safety rules.

---

## 24. Particles 3D

| Command | Args | Description |
|---------|------|-------------|
| `MakeParticles3D` | (id, maxCount) | Create particle system |
| `SetParticles3DColor` | (id, r, g, b) | Default color |
| `SetParticles3DSize` | (id, size) | Particle size |
| `SetParticles3DSpeed` | (id, vx, vy, vz) | Emission velocity |
| `EmitParticles3D` | (id, count [, x, y, z]) | Emit particles |
| `DrawParticles3D` | (id) | Draw particles |
| `GetParticles3DCount` | (id) | Current particle count |
| `GetParticles3DMax` | (id) | Max capacity |

---

## 25. Instancing

| Command | Args | Description |
|---------|------|-------------|
| `MakeInstance` | (baseID, instanceID) | Create instance of object |
| `PositionInstance` | (instanceID, x, y, z) | Set instance position |
| `DrawInstances` | (baseID) | Draw all instances |

---

## 26. Pathfinding

| Command | Args | Description |
|---------|------|-------------|
| `NavMeshLoad` | (id, path) | Load navmesh (stub) |
| `NavMeshFindPath` | (id, startX, startY, startZ, endX, endY, endZ) | Find path |
| `NavMeshDraw` | (id) | Debug draw (stub) |

---

## 27. Matrix / Quaternion

| Command | Args | Description |
|---------|------|-------------|
| `MakeQuaternion` | (id, pitch, yaw, roll) | Create quaternion (stub) |
| `RotateObjectQuat` | (id, quatID) | Apply quat rotation (stub) |
| `GetObjectMatrix` | (id) | Get 4x4 matrix as array (stub) |

---

## 28. Tasks / Coroutines

| Command | Args | Description |
|---------|------|-------------|
| `StartTask` | (name) | Start coroutine |
| `StopTask` | (name) | Stop coroutine (stub) |
| `PauseTask` | (name) | Pause coroutine (stub) |
| `ResumeTask` | (name) | Resume coroutine (stub) |
| `WaitSeconds` | (value) | Yield for seconds |
| `WaitFrames` | (value) | Yield for frames |
| `Yield` | () | Yield one frame |
| `FixedUpdate` | (rate) | Set fixed update rate (e.g. 60) |
| `OnFixedUpdate` | (label$) | Set label for fixed-step callback |

---

## 29. Quick Reference

| Category | Key Commands |
|----------|--------------|
| Window | Window, CloseWindow, SetTargetFPS, Clear, StartDraw, EndDraw, Start3D, End3D |
| Camera | SetCameraPosition, SetCameraTarget, PointCameraAt, GetCameraX/Y/Z, GetCameraPitch/Yaw |
| FPS Camera | FpsCameraOn, FpsCameraOff, FpsCameraPosition, FpsMoveSpeed, FpsLookSpeed, FpsUpdate |
| Objects | MakeCube, MakeSphere, MakeCylinder, MakeGrid, LoadObject, DeleteObject |
| Transform | PositionObject, RotateObject, ScaleObject, MoveObject, TurnObject |
| Queries | GetObjectX/Y/Z, GetObjectPitch/Yaw/Roll, GetObjectScaleX/Y/Z |
| Parenting | ParentObject, UnparentObject |
| Tags | SetObjectTag, GetObjectTag |
| Groups | MakeGroup, AddToGroup, PositionGroup, RotateGroup, DrawGroup, SyncGroup |
| Collision | Raycast, Spherecast, ObjectCollides, PointInObject |
| Water | MakeWater, SetWaterTexture, PositionWater, SetWaterScroll, SetWaterWaveStrength, DrawWater |
| Terrain | MakeTerrain, LoadHeightmap, SetTerrainTexture, PositionTerrain, SetTerrainLayer, DrawTerrain |
| Physics | PhysicsOn, MakeRigidBodyId, MakeBoxCollider, SetRigidBodyPosition, GetRigidBodyX/Y/Z |
| Math | Distance3D, AngleBetween3D, Normalize3D, Dot3D, Cross3D |
| Replication | SyncObject, UnsyncObject, SetObjectOwner, GetObjectOwner |

---

## 30. Minimal 3D Multiplayer Example

```basic
Window 800, 600, "3D Multiplayer"
SetTargetFPS 60

MakeCube 1, 2
PositionObject 1, 0, 0, 0
SyncObject 1

FpsCameraOn
FpsCameraPosition 0, 2, 10
FpsMoveSpeed 5

While Not WindowShouldClose()
  FpsUpdate
  StartDraw
  Clear 30, 30, 50
  Start3D
  DrawObject 1
  DrawGrid 20, 1
  End3D
  EndDraw
Wend

FpsCameraOff
CloseWindow
```

---

## See also

- [Core Command Reference](CORE_COMMAND_REFERENCE.md) – DBP-style command list
- [Blender Workflow](BLENDER_WORKFLOW.md) – Export 3D models (GLTF, FBX, OBJ)
