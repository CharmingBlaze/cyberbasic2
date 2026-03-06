# Aseprite to CyberBASIC2 Workflow

This guide explains how to export sprite sheets from Aseprite for use in CyberBASIC2.

## Aseprite Export

1. **File > Export > Sprite Sheet**
2. **Output**:
   - Sprite sheet: PNG
   - Texture: JSON Data
3. **Layout**:
   - Packed (or by rows/columns)
   - Trim: optional
4. **JSON Data**:
   - Hash (recommended) or Array
   - Include: Frame Tags, Slices

## Loading in CyberBASIC2

```basic
LOAD SPRITE SHEET 1, "sprite.png", "sprite.json"
```

### Grid Mode (no JSON)

For uniform grids:

```basic
LOAD SPRITE SHEET 1, "sprite.png", 32, 32
```

## Animation by Tag

```basic
LOAD SPRITE SHEET 1, "character.png", "character.json"
PLAY SPRITE ANIMATION 1, "walk", 1.0
```

Tags are defined in Aseprite (e.g. "idle", "walk", "attack", "hurt", "death").

## Animation Commands

| Command | Args | Description |
|---------|------|-------------|
| `PLAY SPRITE ANIMATION` | id, tagName, speed | Play animation by tag |
| `STOP SPRITE ANIMATION` | id | Stop animation |
| `SET SPRITE FRAME` | id, frame | Set frame index |
| `GET SPRITE FRAME` | id | Current frame |
| `GET ANIMATION LENGTH` | id, tagName | Frame count for tag |

## Slices (Hitboxes, UI Regions)

Aseprite slices export to `meta.slices` in the JSON. Use them for hitboxes, UI regions, attach points:

```basic
GET SLICE RECT 1, "Hitbox"
```

Returns a string `"x,y,w,h"` in sprite space. Parse with `SPLIT` or similar.

## Drawing

```basic
DRAW SPRITE FRAME 1, GET SPRITE FRAME 1, 100, 100
```

Or use the current frame for animated sprites:

```basic
PLAY SPRITE ANIMATION 1, "walk", 1.0
SYNC
DRAW SPRITE FRAME 1, GET SPRITE FRAME 1, 100, 100
```

## Animation Tick

Sprite animations advance automatically when `SYNC` runs (unified renderer). No manual tick needed.

## Frame Tags

Common tag names:

- `idle` - Standing still
- `walk` - Walking
- `run` - Running
- `attack` - Attack
- `hurt` - Hit reaction
- `death` - Death

## Direction

- `forward` - Loop from start to end
- `reverse` - Loop from end to start
- `pingpong` - Forward then reverse

## See also

- [2D Game API](2D_GAME_API.md) – Full 2D command reference (spritesheets, animation)
- [Core Command Reference](CORE_COMMAND_REFERENCE.md) – Spritesheet commands
