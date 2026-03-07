# CyberBASIC2 Examples

Runnable examples to get started quickly.

## Quick Start

| Example | Description | Run |
|---------|-------------|-----|
| [hello_world.bas](hello_world.bas) | Minimal: prints to console | `cyberbasic examples/hello_world.bas` |
| [first_game.bas](first_game.bas) | 3D spinning cube + mouse orbit camera (DBPRO-style: just SYNC) | `cyberbasic examples/first_game.bas` |
| [input_debug.bas](input_debug.bas) | Test all input: keyboard, mouse, gamepad (mainloop + SYNC) | `cyberbasic examples/input_debug.bas` |

**DBPRO-style:** Use `SYNC` at the end of your loop; the compiler injects BeginDrawing. Input, physics, and rendering work seamlessly.

**Debug:** `cyberbasic --debug examples/first_game.bas` prints render trace (BeginDrawing, SyncFrame, etc.) to console.

## Templates

For minimal game starters with more structure, see [templates/](../templates/):

- **templates/2d_game.bas** – 2D: mainloop, WASD via GetAxisX/GetAxisY
- **templates/3d_game.bas** – 3D: Bullet physics, GAME.CameraOrbit, GAME.MoveWASD

Run: `cyberbasic templates/2d_game.bas` or `cyberbasic templates/3d_game.bas`

## More Examples

Additional examples (physics, GUI, multiplayer, ECS) are documented in the [Documentation Index](../docs/DOCUMENTATION_INDEX.md). Use the [templates](../templates/) and [docs](../docs/) as starting points for 2D, 3D, and hybrid-loop games.
