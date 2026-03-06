# World, Water, Terrain, and Clouds

This document describes the world environment systems: water, terrain, clouds, sky, sun, and time. All systems support simple-to-advanced usage and **always load safely**—missing textures, heightmaps, or shaders never crash the engine.

## Design Goals

- **Always safe to load**: If a map has no sky, no water, no clouds, no atmosphere, the engine runs with sane defaults.
- **Deterministic**: No hidden randomness; any noise is seeded and controllable.
- **Multiplayer-safe**: World visuals (sky, clouds, water) are cosmetic; gameplay state is explicit.
- **Blender-friendly**: Levels can bake sky/water meshes or rely on engine-side systems.

---

## 1. Water Commands (Simple to Advanced)

### Simple Water (Textured Plane)

```basic
MakeWater 1, 500, 500
SetWaterTexture 1, "water.png"
PositionWater 1, 0, 0, 0
DrawWater 1
```

| Command | Args | Description |
|--------|------|-------------|
| `MakeWater` | (id, width, depth) | Create water plane |
| `SetWaterTexture` | (id, path) | Load texture; **fallback: solid blue** if missing |
| `PositionWater` | (id, x, y, z) | Set position |
| `SetWaterLevel` | (id, height) | Set Y position only |
| `SetWaterColor` | (id, r, g, b) | Base color |
| `DrawWater` | (id) | Draw at stored position |

### Animated Water (Scrolling UVs)

```basic
SetWaterScroll 1, 0.01, 0.02
```

| Command | Args | Description |
|--------|------|-------------|
| `SetWaterScroll` | (id, uSpeed, vSpeed) | UV scroll speed (deterministic, time-based) |

### Layered Water (Foam, Normals, Depth Tint)

```basic
SetWaterNormalmap 1, "water_normal.png"
SetWaterFoamTexture 1, "foam.png"
SetWaterDepthColor 1, 0.1, 0.2, 0.4
SetWaterShallowColor 1, 0.2, 0.5, 0.8
```

| Command | Args | Description |
|--------|------|-------------|
| `SetWaterNormalmap` | (id, path) | Normal map; **skip if missing** |
| `SetWaterFoamTexture` | (id, path) | Foam mask; **skip if missing** |
| `SetWaterDepthColor` | (id, r, g, b) | Deep water tint |
| `SetWaterShallowColor` | (id, r, g, b) | Shallow water tint |

### Advanced Water (Waves, Reflection, Refraction)

```basic
SetWaterWaveStrength 1, 0.5
SetWaterWaveSpeed 1, 1.2
SetWaterReflection 1, 1
SetWaterRefraction 1, 1
```

| Command | Args | Description |
|--------|------|-------------|
| `SetWaterWaveStrength` | (id, value) | Wave height |
| `SetWaterWaveSpeed` | (id, value) | Wave animation speed |
| `SetWaterReflection` | (id, onOff) | Reflection; **silently disabled if unsupported** |
| `SetWaterRefraction` | (id, onOff) | Refraction; **silently disabled if unsupported** |

---

## 2. Terrain Commands (Simple to Advanced)

### Simple Terrain (Single Texture)

```basic
MakeTerrain 1, 1000, 1000
SetTerrainTexture 1, "grass.png"
PositionTerrain 1, 0, 0, 0
DrawTerrain 1
```

| Command | Args | Description |
|--------|------|-------------|
| `MakeTerrain` | (id, width, depth) | Create flat terrain |
| `SetTerrainTexture` | (id, path) | Single texture; **fallback: green** if missing |
| `PositionTerrain` | (id, x, y, z) | Set position |
| `DrawTerrain` | (id) | Draw at stored position |

### Heightmap Terrain

```basic
LoadHeightmap 2, "terrain.png", 1000, 1000, 50
SetTerrainTexture 2, "grass.png"
PositionTerrain 2, 0, 0, 0
DrawTerrain 2
```

| Command | Args | Description |
|--------|------|-------------|
| `LoadHeightmap` | (id, path, width, depth, heightScale) | Load grayscale image; **fallback: flat** if missing |

### Multi-Layer Terrain (Splat Maps)

```basic
SetTerrainLayer 2, 0, "grass.png"
SetTerrainLayer 2, 1, "rock.png"
SetTerrainSplatmap 2, "splat.png"
```

| Command | Args | Description |
|--------|------|-------------|
| `SetTerrainLayer` | (id, layerIndex, path) | Layer 0–3 texture |
| `SetTerrainSplatmap` | (id, path) | RGBA splatmap; **fallback: layer 0 only** if missing |

### Procedural Terrain

```basic
MakeTerrain 3, 512, 512
GenerateTerrainNoise 3, 12345, 4, 0.01
SetTerrainTexture 3, "grass.png"
PositionTerrain 3, 0, 0, 0
DrawTerrain 3
```

| Command | Args | Description |
|--------|------|-------------|
| `GenerateTerrainNoise` | (id, seed, octaves, scale) | Procedural heightmap; deterministic via seed |

---

## 3. Clouds and Sky

### Simple Clouds (Texture Dome)

```basic
SetCloudsOn
SetCloudTexture "clouds.png"
SetCloudSpeed 0.001
```

| Command | Args | Description |
|--------|------|-------------|
| `SetCloudsOn` | () | Enable clouds |
| `SetCloudsOff` | () | Disable clouds |
| `SetCloudTexture` | (path) | Sky dome texture; **fallback: disable clouds** if missing |
| `SetCloudSpeed` | (value) | UV scroll speed |
| `SetCloudDensity` | (value) | 0–1 |
| `SetCloudHeight` | (value) | Dome height |
| `SetCloudColor` | (r, g, b) | Tint |

---

## 4. World / Sky / Atmosphere

| Command | Args | Description |
|--------|------|-------------|
| `SetSkybox` | (path) | Load skybox; **fallback: solid color** if missing |
| `SetAmbientLight` | (r, g, b) | Ambient color |
| `SetFog` | (onOff) | Enable fog |
| `SetFogOff` | () | Disable fog |
| `SetFogColor` | (r, g, b) | Fog color |
| `SetFogRange` | (near, far) | Fog distance |
| `SetSunDirection` | (x, y, z) | Sun direction |
| `SetSunColor` | (r, g, b) | Sun color |
| `SetSunIntensity` | (value) | Sun intensity |
| `SetWorldTime` | (hours) | 0–24 |
| `GetWorldTime` | () | Current hours |
| `SetWorldTimeScale` | (value) | Time multiplier |
| `SetWeatherPreset` | (name) | "Clear", "Rain", "Storm", etc. (scripts interpret) |

---

## 5. Safety Rules (Critical)

All load paths are safe:

| Missing | Fallback |
|---------|----------|
| Water texture | Solid blue |
| Terrain texture | Solid green |
| Cloud texture | Disable clouds |
| Skybox | Solid color sky |
| Heightmap | Flat terrain |
| Normal map | Skip normal mapping |
| Foam texture | Skip foam layer |
| Splatmap | Use layer 0 only |
| Shader | Basic textured shader |
| Reflection/refraction | Silently disable |

Invalid parameters (negative, NaN) are clamped to safe defaults. The engine never crashes.

---

## 6. Example: Full Scene

```basic
REM Water
MakeWater 1, 500, 500
SetWaterTexture 1, "water.png"
SetWaterScroll 1, 0.01, 0.02
SetWaterWaveStrength 1, 0.5
PositionWater 1, 0, 0, 0

REM Terrain
LoadHeightmap 2, "terrain.png", 1000, 1000, 50
SetTerrainLayer 2, 0, "grass.png"
SetTerrainLayer 2, 1, "rock.png"
SetTerrainSplatmap 2, "splat.png"
PositionTerrain 2, 0, 0, 0

REM Sky
SetCloudsOn
SetCloudTexture "clouds.png"
SetCloudSpeed 0.001

REM Draw loop
While Not WindowShouldClose()
  StartDraw
  Clear 30, 30, 50
  Start3D
  DrawTerrain 2
  DrawWater 1
  End3D
  EndDraw
Wend
```
