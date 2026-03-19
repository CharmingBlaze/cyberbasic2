# CyberBASIC2 Core Command Reference

Complete reference for the full required command set. All commands are DBP-style, modern, and compatible with raylib-go. For domain-specific commands with examples, see [2D Game API](2D_GAME_API.md) and [3D Game API](3D_GAME_API.md).

---

## 1. Core Engine & Rendering

| Command | Syntax | Description | Status |
|---------|--------|-------------|--------|
| SYNC | `SYNC` | End frame: update, render 3D→2D→GUI, swap buffers | Implemented |
| SET CLEAR COLOR | `SetClearColor r, g, b` | Set background clear color | Implemented |
| SET VSYNC | `SetVsync onOff` | Enable (1) or disable (0) vertical sync. Call before InitWindow | Implemented |
| SET FRAMERATE | `SetFramerate cap` | Target framerate (0 = uncapped). Alias for SetTargetFPS | Implemented |

**When to use:** SYNC ends each frame when using UseUnifiedRenderer. SetClearColor, SetVsync, SetFramerate before InitWindow or at startup.

---

## 2. Camera System

| Command | Syntax | Description | Status |
|---------|--------|-------------|--------|
| MAKE CAMERA | `MakeCamera id` | Create camera with default parameters | Implemented |
| DELETE CAMERA | `DeleteCamera id` | Remove camera | Implemented |
| POSITION CAMERA | `PositionCamera id, x, y, z` | Set camera world position | Implemented |
| ROTATE CAMERA | `RotateCamera id, pitch, yaw, roll` | Set camera rotation (degrees) | Implemented |
| POINT CAMERA | `POINT CAMERA id, x, y, z` | Point camera at world position | Implemented |
| SET CAMERA FOV | `SetCameraFOV value` | Field of view (degrees) | Implemented |
| SET CAMERA RANGE | `SetCameraRange near, far` | Near/far clip planes | Implemented |
| SET CAMERA ACTIVE | `SetCameraActive id` | Use camera for 3D rendering | Implemented |
| ATTACH CAMERA TO OBJECT | `AttachCameraToObject camID, objID` | Parent camera to object | Implemented |

**When to use:** Create cameras in OnStart; set position/target each frame for follow or orbit cameras. Use SetCameraActive before drawing 3D.

---

## 3. Object System

| Command | Syntax | Description | Status |
|---------|--------|-------------|--------|
| LOAD OBJECT | `LoadObject path, id` or `LoadObjectId id, path` | Load 3D model (GLTF/FBX/OBJ) | Implemented |
| DELETE OBJECT | `DeleteObject id` | Remove object and resources | Implemented |
| MAKE OBJECT CUBE | `MakeCube id, size` | Create cube primitive | Implemented |
| MAKE OBJECT SPHERE | `MakeSphere id, radius` | Create UV sphere | Implemented |
| MAKE OBJECT CYLINDER | `MakeCylinder id, radius, height` | Create cylinder | Implemented |
| MAKE OBJECT PLANE | `MakePlane id, width, depth` | Create plane | Implemented |
| POSITION OBJECT | `PositionObject id, x, y, z` | Set world position | Implemented |
| ROTATE OBJECT | `RotateObject id, pitch, yaw, roll` | Set rotation | Implemented |
| SCALE OBJECT | `ScaleObject id, sx, sy, sz` | Set scale | Implemented |
| MOVE OBJECT | `MoveObject id, distance` | Move along local Z | Implemented |
| TURN OBJECT | `TurnObject id, pitch, yaw, roll` | Add to rotation | Implemented |
| SHOW OBJECT | `ShowObject id` | Make visible | Implemented |
| HIDE OBJECT | `HideObject id` | Make invisible | Implemented |
| CLONE OBJECT | `CloneObject newID, sourceID` | Copy object (shares meshes) | Implemented |
| ATTACH OBJECT | `AttachObject childID, parentID` | Parent object | Implemented |
| DETACH OBJECT | `DetachObject id` | Remove parent | Implemented |

**When to use:** Load objects once; PositionObject, RotateObject, DrawObject each frame. Use AttachObject for hierarchies (e.g. weapon on character).

---

## 4. Material System

| Command | Syntax | Description | Status |
|---------|--------|-------------|--------|
| SET OBJECT COLOR | `SetObjectColor id, r, g, b` | Base color tint | Implemented |
| SET OBJECT ALPHA | `SetObjectAlpha id, a` | Alpha (0-255) | Implemented |
| SET OBJECT TEXTURE | `SetObjectTexture id, textureID` or `id, "file.png"` | Diffuse texture (ID or path) | Implemented |
| SET OBJECT NORMALMAP | `SetObjectNormalmap id, "file.png"` | Normal map | Implemented |
| SET OBJECT ROUGHNESS | `SetObjectRoughness id, value` | Roughness 0-1 | Implemented |
| SET OBJECT METALLIC | `SetObjectMetallic id, value` | Metallic 0-1 | Implemented |
| SET OBJECT EMISSIVE | `SetObjectEmissive id, r, g, b` | Emissive color | Implemented |
| SET OBJECT SHADER | `SetObjectShader id, shaderID` | Custom shader | Implemented |
| SET SHADER UNIFORM | `SetShaderUniform shaderID, name$, value` | Shader uniform | Implemented |

**When to use:** Set material properties after loading; override per-object with SetObjectRoughness, SetObjectMetallic for PBR.

---

## 5. Lighting & Shadows

| Command | Syntax | Description | Status |
|---------|--------|-------------|--------|
| MAKE LIGHT | `MakeLight id, type` | type: 0=point, 1=directional, 2=spot | Implemented |
| DELETE LIGHT | `DeleteLight id` | Remove light | Implemented |
| POSITION LIGHT | `PositionLight id, x, y, z` | Set position | Implemented |
| ROTATE LIGHT | `RotateLight id, pitch, yaw, roll` | Set direction | Implemented |
| SET LIGHT COLOR | `SetLightColor id, r, g, b` | Light color | Implemented |
| SET LIGHT INTENSITY | `SetLightIntensity id, value` | Brightness | Implemented |
| SET LIGHT RANGE | `SetLightRange id, value` | Range for point/spot | Implemented |
| SET LIGHT ANGLE | `SetLightAngle id, degrees` | Cone angle for spot | Implemented |
| ENABLE SHADOWS | `EnableShadows id` | Enable shadow casting (directional, spot, or point) | Implemented |
| DISABLE SHADOWS | `DisableShadows id` | Disable shadow casting | Implemented |
| SET SHADOW CASCADES | `SetShadowCascades count` | Override cascade count (1, 3, or 4) for directional shadows | Implemented |
| SHADOW CASCADE COUNT | `ShadowCascadeCount()` | Return current cascade count | Implemented |

**When to use:** Create lights in OnStart; directional (type 1) for sun/moon, spot (type 2) or point (type 0) for local lights. EnableShadows(id) enables shadows for any light type. SetShadowQuality("low"|"medium"|"high") for performance; SetShadowCascades(1|3|4) for cascade override.

---

## 6. 3D Import & Levels

| Command | Syntax | Description | Status |
|---------|--------|-------------|--------|
| LOAD OBJECT | See Object System | GLTF/FBX/OBJ | Implemented |
| LOAD LEVEL | `LoadLevel id, path` | Load scene | Implemented |
| UNLOAD LEVEL | `UnloadLevel id` | Unload level | Implemented |

---

## 7. Animation System

| Command | Syntax | Description | Status |
|---------|--------|-------------|--------|
| PLAY ANIMATION | `PlayAnimation objectID, animID, clipIndex, speed` or `objectID, animID, speed` | Play skeletal animation | Implemented |
| STOP ANIMATION | `StopAnimation objectID` | Stop animation | Implemented |
| SET ANIMATION FRAME | `SetAnimationFrame objectID, frame` | Jump to frame | Implemented |
| SET ANIMATION SPEED | `SetAnimationSpeed objectID, speed` | Playback speed | Implemented |
| SET ANIMATION LOOP | `SetAnimationLoop objectID, onOff` | Loop on/off | Implemented |
| RESET BONES | `ResetBones objectID` | Reset to bind pose | Implemented |
| LOAD MESH ANIMATION | `LoadMeshAnimation id, folder$, frameCount` | Load mesh sequence | Implemented |
| PLAY MESH ANIMATION | `PlayMeshAnimation id, speed` | Play mesh frames | Implemented |
| SET MESH ANIMATION FRAME | `SetMeshAnimationFrame id, frame` | Set frame index | Implemented |
| SET BONE ROTATION | `SetBoneRotation id, bone$, pitch, yaw, roll` | Manual bone control | Stub |
| SET BONE POSITION | `SetBonePosition id, bone$, x, y, z` | Manual bone offset | Stub |
| IK SOLVE TWOBONE | `IKSolveTwoBone id, upper$, lower$, tx, ty, tz` | Two-bone IK | Experimental |

---

## 8. World Systems

### Water
| Command | Syntax | Status |
|---------|--------|--------|
| MAKE WATER | `MakeWater id, width, depth` | Implemented |
| SET WATER TEXTURE | `SetWaterTexture id, "file.png"` | Implemented |
| SET WATER NORMALMAP | `SetWaterNormalmap id, "file.png"` | Implemented |
| SET WATER SCROLL | `SetWaterScroll id, uSpeed, vSpeed` | Implemented |
| SET WATER WAVE | `SetWaterWave id, strength` | Implemented |

### Terrain
| Command | Syntax | Status |
|---------|--------|--------|
| LOAD HEIGHTMAP | `LoadHeightmap id, file$, width, depth, heightScale` | Implemented |
| SET TERRAIN LAYER | `SetTerrainLayer id, layerIndex, "file.png"` | Implemented |
| SET TERRAIN SPLATMAP | `SetTerrainSplatmap id, "splat.png"` | Implemented |

### Clouds
| Command | Syntax | Status |
|---------|--------|--------|
| SET CLOUDS ON/OFF | `SetCloudsOn` / `SetCloudsOff` | Implemented |
| SET CLOUD TEXTURE | `SetCloudTexture "clouds.png"` | Implemented |
| SET CLOUD SPEED | `SetCloudSpeed value` | Implemented |

### Skybox
| Command | Syntax | Status |
|---------|--------|--------|
| SET SKYBOX | `SetSkybox "sky.png"` | Implemented |
| SET SKYBOX CUBEMAP | `SetSkyboxCubemap right$, left$, top$, bottom$, front$, back$` | Implemented |

---

## 9. Physics (Bullet)

| Command | Syntax | Status |
|---------|--------|--------|
| MAKE RIGIDBODY | `MakeRigidBodyId id, x, y, z, mass` | Implemented |
| DELETE RIGIDBODY | `DestroyBody3D worldId, bodyId` or `DeleteBody3D bodyId` (default world) | Implemented |
| MAKE COLLIDER BOX | `MakeBoxCollider id, sx, sy, sz` | Implemented |
| MAKE COLLIDER SPHERE | `MakeSphereCollider id, radius` | Implemented |
| MAKE COLLIDER CAPSULE | `MakeCapsuleCollider id, radius, height` | Implemented |
| APPLY FORCE | `ApplyForce3D worldId, bodyId, fx, fy, fz` | Implemented |
| APPLY IMPULSE | `ApplyImpulse3D worldId, bodyId, fx, fy, fz` | Implemented |
| SET VELOCITY | `SetRigidBodyVelocity id, x, y, z` | Implemented |
| SET ANGULAR VELOCITY | `SetAngularVelocity id, x, y, z` | Implemented |
| RAYCAST | `Raycast ox, oy, oz, dx, dy, dz` + RayHitX/Y/Z/Body | Implemented |

---

## 10. Sprites & 2D

| Command | Syntax | Status |
|---------|--------|--------|
| LOAD SPRITE | `LoadSprite path, id` | Implemented |
| DELETE SPRITE | `DeleteSprite id` | Implemented |
| DRAW SPRITE | `DrawSprite id, x, y` / `Sprite id, x, y` | Implemented |
| DRAW SPRITE ROTATED | `DrawSpriteRotated id, x, y, angle` | Implemented |
| DRAW SPRITE SCALED | `DrawSpriteScaled id, x, y, sx, sy` | Implemented |
| SET SPRITE COLOR | `SetSpriteColor id, r, g, b, a` | Implemented |

### Spritesheets (Aseprite + Grid)

| Command | Syntax | Status |
|---------|--------|--------|
| LOAD SPRITE SHEET | `LoadSpritesheet id, pngPath, jsonPath` or `id, path, frameW, frameH` | Implemented |
| PLAY SPRITE ANIMATION | `PlaySpriteAnimation id, tagName, speed` | Implemented |
| STOP SPRITE ANIMATION | `StopSpriteAnimation id` | Implemented |
| SET SPRITE FRAME | `SetSpriteFrame id, frame` | Implemented |
| GET SPRITE FRAME | `GetSpriteFrame id` | Implemented |
| DRAW SPRITE FRAME | `DrawSpriteFrame id, frame, x, y` | Implemented |
| GET SLICE RECT | `GetSliceRect id, sliceName` | Implemented |
| GET ANIMATION LENGTH | `GetAnimationLength id, tagName` | Implemented |

---

## 11. GUI System

| Command | Syntax | Status |
|---------|--------|--------|
| BEGIN UI | `BeginUI` | Implemented |
| END UI | `EndUI` | Implemented |
| LABEL | `Label text` (auto layout) | Implemented |
| LABEL AT | `LabelAt x, y, text` | Implemented |
| BUTTON | `Button text` (auto layout) | Implemented |
| BUTTON AT | `ButtonAt x, y, w, h, text` | Implemented |
| CHECKBOX | `Checkbox text, checked` | Implemented |
| CHECKBOX AT | `CheckboxAt x, y, label, checked` | Implemented |
| SLIDER | `Slider text, value, min, max` | Implemented |
| SLIDER AT | `SliderAt x, y, width, min, max, value` | Implemented |

---

## 12. Audio System

| Command | Syntax | Status |
|---------|--------|--------|
| LOAD SOUND | `LoadSound id, path` | Implemented |
| DELETE SOUND | `DeleteSound id` | Implemented |
| PLAY SOUND | `PlaySound id` | Implemented |
| STOP SOUND | `StopSound id` | Implemented |
| PAUSE SOUND | `PauseSound id` | Implemented |
| SET SOUND VOLUME | `SetSoundVolume id, value` | Implemented |
| SET SOUND LOOP | `SetSoundLoop id, onOff` | Implemented |

---

## 13. Input System

| Command | Syntax | Status |
|---------|--------|--------|
| KEYDOWN | `KeyDown(key)` | Implemented |
| KEYUP | `KeyUp(key)` | Implemented |
| MOUSECLICK | `MouseClick(button)` | Implemented |
| MOUSEX / MOUSEY | `MouseX` / `MouseY` | Implemented |
| MOUSEWHEEL | `GetMouseWheelMove` | Implemented |
| GAMEPAD AXIS | `GamepadAxis pad, axis` | Implemented |
| GAMEPAD BUTTON | `GamepadButton pad, button` | Implemented |

---

## 14. Resource Management

| Command | Syntax | Status |
|---------|--------|--------|
| UNLOAD OBJECT | `DeleteObject id` | Implemented |
| UNLOAD TEXTURE | `DeleteTexture id` | Implemented |
| UNLOAD SOUND | `DeleteSound id` | Implemented |
| UNLOAD LEVEL | `UnloadLevel id` | Implemented |

---

## 15. File I/O

| Command | Syntax | Status |
|---------|--------|--------|
| OPEN FILE | `OpenFile(path, mode)` mode: 0=read, 1=write, 2=append | Implemented |
| CLOSE FILE | `CloseFile handle` | Implemented |
| READ LINE | `ReadLine handle` | Implemented |
| WRITE LINE | `WriteLine handle, text` | Implemented |
| READ BYTE | `ReadByte handle` | Implemented |
| WRITE BYTE | `WriteByte handle, value` | Implemented |

---

## 16. Math Library

| Command | Syntax | Status |
|---------|--------|--------|
| VEC3 | `Vec3(x, y, z)` | Implemented |
| DOT | `Dot x1,y1,z1, x2,y2,z2` | Implemented |
| CROSS | `Cross3D` / `Vector3CrossProduct` | Implemented |
| NORMALIZE | `Normalize3D` / `Vector3Normalize` | Implemented |
| DISTANCE | `Distance x1,y1,z1, x2,y2,z2` | Implemented |
| LENGTH | `Vector3Length` | Implemented |
| LERP | `Lerp a, b, t` | Implemented |
| CLAMP | `Clamp value, min, max` | Implemented |

---

## See also

- [2D Game API](2D_GAME_API.md) – Full 2D API with examples
- [3D Game API](3D_GAME_API.md) – Full 3D API with examples
- [DBP Extended](DBP_EXTENDED.md) – Module-by-module implementation details
- [Blender Workflow](BLENDER_WORKFLOW.md) – Export 3D models for CyberBASIC2
- [Aseprite Workflow](ASEPRITE_WORKFLOW.md) – Sprite sheet export with tags and slices
- [World, Water, Terrain](WORLD_WATER_TERRAIN.md) – Water, terrain, skybox, clouds
- [Level Loading](LEVEL_LOADING.md) – Unified 3D loading pipeline
