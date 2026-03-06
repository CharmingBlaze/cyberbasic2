# 3D Loading Design Spec

Design goals and safe loading behavior for the CyberBASIC 3D loading pipeline.

## Design Goals (Non-Negotiables)

- **Never fatal:** Missing lights, materials, textures, collision, animations → use defaults or no-op, never crash
- **Blender-first:** GLTF/FBX/OBJ from Blender should "just work"
- **Deterministic:** Same file → same scene graph → same IDs
- **Multiplayer-safe:** No hidden state, no randomization, no auto-physics, no auto-networking
- **Graceful degradation:** Load and render what we can; skip or default what is missing

## Safe Loading Behavior

### Textures

| Condition | Behavior |
|-----------|----------|
| Texture path empty | Use default 1x1 white |
| Texture load fails | Use default 1x1 white |
| No textures in file | Materials use default texture |

### Materials

| Condition | Behavior |
|-----------|----------|
| No materials in file | Create one default (0.8, 0.8, 0.8) |
| Material index out of range | Use index 0 |
| Invalid BaseColorTextureIndex | Use default texture |

### Meshes

| Condition | Behavior |
|-----------|----------|
| Missing normals | Compute flat normals from triangle faces |
| Missing texcoords | Use zeros |
| Empty mesh (no vertices) | Skip; do not add object |
| Mesh has no material | Use default material |

### Collision

| Condition | Behavior |
|-----------|----------|
| No collision nodes | LoadLevelCollision returns 0 |
| Collider mesh invalid | Use box from bounds or default size |
| LoadLevelCollision on unknown level | Return 0 |

### Animation

| Condition | Behavior |
|-----------|----------|
| LoadAnimation on file with no animations | Succeed; store nothing |
| GetAnimationLength on unknown anim | Return 0 |
| PlayAnimation with no anim | No-op |
| SetAnimationFrame with no anim | No-op |
| GetAnimationFrame with no state | Return 0 |

### Skeleton / IK

| Condition | Behavior |
|-----------|----------|
| No skin in GLTF | Skeleton = nil; treat as static |
| IKSolveTwoBone with no skeleton | No-op |
| IKSolveTwoBone with bone not found | No-op |
| IKEnable on unknown object | Store state (no crash) |

## Blender Pipeline

### Export Settings (GLTF)

- **Format:** glTF 2.0 (.gltf or .glb)
- **Include:** Meshes, Materials, Textures
- **Apply modifiers** before export

### Collision Naming

Name objects in Blender to auto-detect as colliders:

| Pattern | Example |
|---------|---------|
| `col_*` | col_floor, col_wall |
| `collision_*` | collision_trigger |
| `Collision*` | CollisionBox |

Or add custom property `collision: true` in GLTF extras (JSON).

### Collider Type from Name

| Name contains | Type |
|---------------|------|
| `sphere` | Sphere |
| `capsule` | Capsule |
| (default) | Box |

## Programmer Checklist

- [ ] Use LoadLevel before LoadLevelCollision
- [ ] Call PhysicsOn if using LoadLevelCollision
- [ ] UnloadLevel cleans colliders; no manual DestroyBody3D needed
- [ ] SpawnPrefab returns first object ID; use for positioning
- [ ] LoadAnimation succeeds even with empty file; check GetAnimationLength
- [ ] IKEnable(objectID, 1) before IKSolveTwoBone
- [ ] Same asset paths on all clients (multiplayer)

## Object ID Ranges

| Source | Base | Range |
|--------|------|-------|
| Level | levelID × 100000 | 100000, 100001, ... |
| Prefab spawn | 500000 + counter×10000 | 510000, 520000, ... |
| User objects | 1–99999 | Typical |

## Related

- [LEVEL_LOADING.md](LEVEL_LOADING.md) — User-facing API
- [3D_GAME_API.md](3D_GAME_API.md) — Full 3D reference
