# Asset Pipeline

Central asset loading, caching, and lifecycle for models and textures. Supports reference counting and shared parsed-model reuse where wired.

---

## Purpose

- **Avoid duplicate loads:** Same path loads once; subsequent loads increment refcount.
- **Unload when done:** `UnloadAsset(path)` decrements refcount; GPU resources freed when refs reach 0.
- **Warm cache:** `PreloadAsset` or `LoadAsset` before `LoadObject` to reuse parsed models.

---

## Architecture

```
User code                    Asset pipeline
─────────────────────────────────────────────────────────────
LoadAsset(path)     →    Parsed-model cache (by path)
PreloadAsset(path)  →    Same as LoadAsset
LoadObject(id,path) →    Uses cache when available
LoadLevel(id,path)  →    Imports directly (bypasses cache)
LoadPrefab(id,path) →    Imports directly (bypasses cache)
UnloadAsset(path)   →    Decrement refcount; unload at 0
```

**Packages:**
- `compiler/runtime/resources` — Model and texture cache with refcounting
- `compiler/runtime/assets` — Parsed-model cache for GLTF/OBJ
- `compiler/bindings/dbp` — LoadObject, LoadLevel, LoadPrefab

---

## API Surface

| Command | Args | Returns | Description |
|---------|------|---------|-------------|
| `LoadAsset` | (path) | path (string) | Load model or texture; increment refcount. Returns path key. |
| `UnloadAsset` | (path) | — | Decrement refcount; unload when refs reach 0. |
| `AssetExists` | (path) | 1 or 0 | True if asset is currently cached. |
| `PreloadAsset` | (path) | — | Synchronous preload; same behavior as LoadAsset. |

---

## Supported Formats

| Type | Extensions |
|------|------------|
| Models | `.gltf`, `.glb`, `.obj` |
| Textures | `.png`, `.jpg`, `.jpeg`, `.bmp`, `.tga`, `.gif` |

---

## Defaults

- **Cache key:** Path string (relative or absolute).
- **Refcount:** Starts at 1 on first load; increments on subsequent LoadAsset/LoadObject for same path.
- **Unload:** Call `UnloadAsset(path)` only after all users (objects, levels) are gone.

---

## Current Integration Status

| Entry Point | Cache-Backed | Notes |
|-------------|--------------|-------|
| `LoadObject(id, path)` | Yes (static) | Uses parsed-model cache for static builds. |
| Animated GLTF | No | Uses raylib native loader for skeletal animation. |
| `LoadLevel(id, path)` | No | Parses through importer directly. |
| `LoadLevelWithHierarchy(id, path)` | No | Same as LoadLevel; preserves hierarchy. |
| `LoadPrefab(id, path)` | No | Not cache-backed. |
| Texture upload (BuildModel) | Partial | Model textures uploaded per build; shared texture cache exists but not fully wired. |

---

## Edge Cases

- **Missing file:** Load fails; no cache entry. Check return/error.
- **Missing textures in GLTF:** Engine uses 1×1 white placeholder.
- **Unload before DeleteObject:** Safe; refcount prevents premature unload. Unload after DeleteObject.
- **Same path, different cases:** Path is case-sensitive; `"Crate.glb"` and `"crate.glb"` are different keys.

---

## Level vs Object Loading

| Command | Use Case | Cache |
|---------|----------|-------|
| `LoadObject(id, path)` | Single object, best cache support | Yes |
| `LoadLevel(id, path)` | Full scene, meshes + materials + textures | No |
| `LoadLevelWithHierarchy(id, path)` | Scene with node hierarchy | No |
| `LoadPrefab(id, path)` | Reusable template | No |

---

## Performance Considerations

- **Preload during loading screen:** Call `PreloadAsset` for critical assets before gameplay.
- **Unload unused:** Call `UnloadAsset` when switching levels or removing objects to free GPU memory.
- **Level loads:** LoadLevel/LoadLevelWithHierarchy parse and build each time; no cache. Use for one-time level load.

---

## Multiplayer / Determinism

Asset loading is local and does not affect network simulation. Load the same assets on server and clients for consistent rendering.

---

## Contributor Notes

- **Parsed-model cache:** `compiler/runtime/assets/assets.go` — `LoadModelForBuild` for cache integration.
- **Resource manager:** `compiler/runtime/resources/manager.go` — `LoadModel`, `UnloadModel`, `GetModel`, `ModelExists`, `TextureExists`.
- **Wiring LoadLevel:** `compiler/bindings/dbp/dbp.go` — `LoadLevel` calls model importer; to add cache, route through `assets.LoadModelForBuild`.
- **Texture reuse:** `BuildModel` in dbp_model.go uploads textures per build; integrate with `resources` for shared texture reuse.

---

## Example

```basic
' Warm the cache for a static object
PreloadAsset "props/crate.glb"

' Later: uses the parsed-model cache
LoadObject 1, "props/crate.glb"

' Level loads still build/import directly
LoadLevelWithHierarchy 1, "levels/test_room.glb"

' When finished
DeleteObject 1
UnloadLevel 1
UnloadAsset "props/crate.glb"
```

---

## See Also

- [Level Loading](LEVEL_LOADING.md)
- [3D Loading Spec](3D_LOADING_SPEC.md)
- [Blender Workflow](BLENDER_WORKFLOW.md)
- [Documentation Index](DOCUMENTATION_INDEX.md)
