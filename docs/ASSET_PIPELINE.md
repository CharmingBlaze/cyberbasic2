# Asset Pipeline

CyberBASIC2 has a real asset cache, but it is currently only wired into part of the loading stack.

## Commands

| Command | Args | Description |
|---------|------|-------------|
| `LoadAsset` | (path) | Load model or texture into the shared cache; returns the path key. |
| `UnloadAsset` | (path) | Decrement the cached refcount for that path. |
| `AssetExists` | (path) | Returns 1 if the asset is currently cached. |
| `PreloadAsset` | (path) | Synchronous preload; same behavior as `LoadAsset`. |

## Supported Formats

- Models: `.gltf`, `.glb`, `.obj`
- Textures: `.png`, `.jpg`, `.jpeg`, `.bmp`, `.tga`, `.gif`

## Current Integration Status

- `LoadObject(id, path)` uses the shared parsed-model cache for static model builds.
- Animated GLTF objects still use raylib's native loader so skeletal animation continues to work.
- `LoadLevel`, `LoadLevelWithHierarchy`, and `LoadPrefab` currently parse through the model importer directly instead of reusing the shared parsed-model cache.
- Texture caching exists in `runtime/resources`, but `BuildModel` still uploads model textures per build instead of reusing the shared texture cache for imported model materials.

## Practical Guidance

1. Use `PreloadAsset` or `LoadAsset` when you want to warm the parsed-model cache before a later `LoadObject`.
2. Keep using `LoadLevel` and `LoadLevelWithHierarchy` normally, but do not assume they are currently cache-backed.
3. Call `UnloadAsset` only after the corresponding object or texture users are gone.

## Level vs Object Loading

- `LoadObject(id, path)`: best current entry point for cache-backed model loading.
- `LoadLevel(id, path)`: imports a whole level and builds DBP objects/materials/textures, but currently bypasses the parsed-model cache.
- `LoadLevelWithHierarchy(id, path)`: same as `LoadLevel`, with node hierarchy preserved between mesh nodes.
- `LoadPrefab(id, path)`: loads a reusable template; currently not cache-backed.

## What Is Cached

- Parsed models: cached by source path.
- Standalone textures loaded through the resource manager: cached by path.
- GPU model builds: not shared; each build still creates its own runtime meshes/material state.

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
