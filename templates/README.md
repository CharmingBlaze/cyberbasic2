# CyberBasic game templates

Minimal runnable games to copy or run as a starting point.

- **2d_game.bas** – 2D: window, WASD movement via GetAxisX/GetAxisY, draw loop.
- **3d_game.bas** – 3D: Bullet world, player sphere, GAME.CameraOrbit + GAME.MoveWASD, draw loop.

Run:
```bash
cyberbasic templates/2d_game.bas
cyberbasic templates/3d_game.bas
```

Copy a template to your project and extend (add sprites, physics bodies, UI, etc.).

For the **hybrid loop** (automatic physics step and render queue), define `update(dt)` and `draw()` and use a game loop with an empty body; see [Program Structure – hybrid loop](../docs/PROGRAM_STRUCTURE.md#hybrid-updatedraw-loop) and examples that use the hybrid style.
