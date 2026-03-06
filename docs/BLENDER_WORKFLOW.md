# Blender to CyberBASIC2 Workflow

This guide explains how to export 3D models from Blender for use in CyberBASIC2.

## Supported Formats

| Format | Support | Use Case |
|--------|---------|----------|
| **GLTF/GLB** | Full | Primary format. Meshes, materials, skeleton, animations. |
| **OBJ** | Mesh only | Prototyping, terrain, simple static meshes. |
| **FBX** | Stub | Use Blender to export GLTF instead. |

## Blender Export Settings (GLTF)

1. **File > Export > glTF 2.0 (.glb/.gltf)**
2. **Format**: GLB (binary, single file) or GLTF (JSON + external files)
3. **Include**:
   - Limit to: Visible Objects (or All)
   - Apply Modifiers
   - UVs
   - Normals
   - Tangents (for normal maps)
4. **Transform**:
   - Forward: -Z Forward
   - Up: Y Up
   - Scale: 1.00
5. **Geometry**:
   - UVs
   - Normals
   - Tangents
   - Vertex Colors (if used)
6. **Animation**:
   - Bake animation
   - Group by NLA Track (optional)
   - Optimize: Remove redundant keys

## PBR Materials

CyberBASIC2 supports PBR materials from GLTF:

- **Base color** (albedo texture)
- **Normal map**
- **Metallic-roughness** (combined texture or scalar values)
- **Emissive** (color and texture)

In Blender, use Principled BSDF. Export will preserve these.

## Loading in CyberBASIC2

```basic
LOAD OBJECT 1, "character.glb"
LOAD ANIMATION 1, "character.glb"
PLAY ANIMATION 1, 1, 0, 1.0
```

- `LOAD OBJECT 1, "character.glb"` - Loads model. If the file has animations, uses raylib for full skeletal support.
- `LOAD ANIMATION 1, "character.glb"` - Loads all animation clips.
- `PLAY ANIMATION objectID, animID, clipIndex, speed` - Plays clip at index (0 = first). Use 3 args for backward compatibility: `(objectID, animID, speed)` uses clip 0.

## Animation Commands

| Command | Args | Description |
|---------|------|-------------|
| `PLAY ANIMATION` | objectID, animID, clipIndex, speed | Play skeletal animation |
| `STOP ANIMATION` | objectID | Stop animation |
| `SET ANIMATION FRAME` | objectID, frame | Set frame index |
| `SET ANIMATION SPEED` | objectID, speed | Playback speed |
| `SET ANIMATION LOOP` | objectID, onOff | Loop on/off |
| `GET ANIMATION LENGTH` | animID, clipIndex | Frame count |
| `GET ANIMATION NAME` | animID, clipIndex | Clip name |

## Material Overrides

```basic
SET OBJECT TEXTURE 1, "custom_diffuse.png"
SET OBJECT NORMALMAP 1, "custom_normal.png"
SET OBJECT ROUGHNESS 1, 0.5
SET OBJECT METALLIC 1, 0.2
SET OBJECT EMISSIVE 1, 1.0, 0.0, 0.0
```

## OBJ for Simple Models

For static meshes, terrain, or simple props:

```basic
LOAD OBJECT 1, "terrain.obj"
```

OBJ supports vertices, normals, UVs, and diffuse texture (MTL).

## FBX

FBX import is not yet implemented. Use Blender to export to GLTF for best results.

## See also

- [3D Game API](3D_GAME_API.md) – Full 3D command reference
- [Core Command Reference](CORE_COMMAND_REFERENCE.md) – Object, animation, material commands
