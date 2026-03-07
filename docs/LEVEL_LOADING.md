# Level Loading API

CyberBASIC2 provides a unified 3D loading pipeline. **LOAD LEVEL** loads everything automatically: meshes, materials, textures, hierarchy, and optionally collision. No extra commands required for basic loading.

## Quick Start

```basic
LOAD LEVEL 1, "castle.gltf"
LOAD LEVEL COLLISION 1
PHYSICS ON

WHILE NOT WindowShouldClose()
  START DRAW
    CLEAR 30, 30, 50
    START 3D
      DRAW LEVEL 1
    END 3D
  END DRAW
WEND

UNLOAD LEVEL 1
```

Textures, materials, and hierarchy load automatically. Call `LoadLevelCollision` to enable physics colliders from the level. GLTF punctual lights are not imported into runtime DBP lights yet, so add gameplay lights explicitly after loading if you need them.

## Core Commands

| Command | Args | Description |
|---------|------|-------------|
| `LoadLevel` | (id, path) | Load a full level from file. Parses, uploads to GPU, and creates objects. |
| `DrawLevel` | (id) | Draw all objects in the level. Call between Start3D/End3D. |
| `UnloadLevel` | (id) | Free all level resources (objects, textures, lights, colliders). |

## Level Collision

| Command | Args | Description |
|---------|------|-------------|
| `LoadLevelCollision` | (id) | Create physics bodies for level colliders. Returns count. Call after LoadLevel. |
| `GetLevelColliderCount` | (id) | Return number of colliders (0 if LoadLevelCollision not called). |
| `GetLevelCollider` | (id, index) | Get physics body ID string at index (for Raycast, BodyCollides, etc.). |

**Collision detection in GLTF:** Nodes or meshes named `col_*`, `collision_*`, or `Collision*` are imported as colliders. You can also set `extras.collision: true` on a node. Collider type is inferred from name: `sphere`, `capsule` → sphere/capsule; otherwise box. Size comes from mesh bounds or node scale.

**Example:**
```basic
LOAD LEVEL 1, "arena.gltf"
n = LoadLevelCollision 1
PhysicsOn
' Use GetLevelCollider 1, i for physics queries
```

## Supported Formats

| Format | Extension | Library |
|--------|-----------|---------|
| GLTF | .gltf, .glb | qmuntal/gltf |
| OBJ | .obj | flywave/go-obj |
| FBX | .fbx | Not yet supported |

## What Loads Automatically

- **Meshes** — All geometry (empty meshes skipped)
- **Materials** — PBR (base color, metallic, roughness); default if none
- **Textures** — Base color, normal, metallic-roughness maps (default 1x1 white on failure)
- **Lights** — Not yet imported from GLTF `KHR_lights_punctual`; create runtime lights explicitly
- **Hierarchy** — Node transforms (position, rotation, scale)
- **Skeleton** — From GLTF skin (for animation and experimental IK requests)
- **Colliders** — Detected from naming; use LoadLevelCollision to create physics bodies

## Graceful Degradation

Loading **never fails** on missing data:

- **Missing textures** — Use 1x1 white placeholder
- **Missing materials** — Use default (0.8, 0.8, 0.8)
- **Missing normals** — Compute flat normals from faces
- **Missing texcoords** — Use zeros
- **Empty meshes** — Skipped
- **No collision data** — LoadLevelCollision returns 0 (no error)

## Object ID Allocation

Level objects use IDs: `levelID * 100000 + index`. For level 1 with 50 objects, IDs are 100000–100049. This avoids collision with user-created objects (typically 1–99999).

## Optional Commands

| Command | Args | Description |
|---------|------|-------------|
| `GetLevelObjectCount` | (id) | Return number of objects in level |
| `GetLevelObject` | (id, index) | Get object ID at index (returns value for assignment) |

## LoadPrefab and SpawnPrefab

| Command | Args | Description |
|---------|------|-------------|
| `LoadPrefab` | (id, path) | Load prefab template from file. |
| `SpawnPrefab` | (id, x, y, z) | Instantiate prefab at position. Returns new object ID. |

**Example:**
```basic
LoadPrefab 1, "enemy.gltf"
objId = SpawnPrefab 1, 0, 0, 5
PositionObject objId, 10, 0, 0
```

## LoadObject (Single Model)

`LoadObject(id, path)` also uses the new pipeline. It loads a single model and creates one DBP object at the given id. Supports the same formats (GLTF, OBJ).

## Blender Export Tips

- **GLTF:** File → Export → glTF 2.0 (.gltf or .glb). Include meshes, materials, textures.
- **Collision:** Name objects `col_floor`, `collision_wall`, or `CollisionBox` to auto-detect.
- **Custom collision:** Add `extras: { "collision": true }` to a node (JSON).
- **Skeleton:** Export with skin; bones are available for IK.

## Safe Loading Checklist

- [ ] Same file path on all clients (multiplayer)
- [ ] No random or time-based behavior in load
- [ ] Call LoadLevelCollision after LoadLevel if using physics
- [ ] UnloadLevel cleans up colliders automatically

## Multiplayer Safety

Level loading is deterministic:

- No random numbers
- No time-based behavior
- Same file on all clients yields identical scene
- No auto-physics, auto-animation, or auto-sync unless explicitly requested

## Related

- [3D_GAME_API.md](3D_GAME_API.md) — Full 3D API reference
- [DBP_EXTENDED.md](DBP_EXTENDED.md) — DBP-style command reference
- [3D_LOADING_SPEC.md](3D_LOADING_SPEC.md) — Design spec and safe loading rules
