# CyberBasic Examples

Run any example with:
```bash
cyberbasic examples/first_game.bas
```

## Quick start

1. Run **first_game.bas** – minimal 2D game (WASD move a circle).
2. Try **templates/2d_game.bas** or **templates/3d_game.bas** for a copy-paste starting point (see [templates/README.md](../templates/README.md)).
3. For 3D: **mario64.bas** or **run_3d_physics_demo.bas**. For 2D physics: **box2d_demo.bas**.

## Start here

- **first_game.bas** – Minimal game loop: window, input (WASD), DrawCircle, WindowShouldClose
- **minimal_window_test.bas** – Bare window open/close
- **hello_world.bas** – Simple Print and window

## 2D / Shapes

- **dot_and_colors_demo.bas** – Dots and color constants
- **2d_shapes_demo.bas** – Rectangles, circles, lines
- **window_demo.bas** – Window and drawing

## 3D

- **run_3d_physics_demo.bas** – 3D physics (Bullet) + raylib 3D
- **mario64.bas** – Mario64-style camera and movement
- **minimal_3d_demo.bas** – Basic 3D scene

## Physics 2D (Box2D)

- **box2d_demo.bas** – Box2D world, bodies, click to spawn
- **simple_box2d_demo.bas** – Minimal Box2D
- **simplest_box2d_demo.bas** – Minimal Box2D setup

## Physics 3D (Bullet)

- **bullet_demo.bas** – Bullet 3D physics
- **run_3d_physics_demo.bas** – 3D physics demo

## Input

- **first_game.bas** – IsKeyDown(KEY_W) etc.; for axis-style movement see **templates/2d_game.bas** (GetAxisX, GetAxisY).

## Other

- **test_enum.bas**, **test_const.bas**, **case_test.bas** – Language features
- **math_matrix_test.bas** – Matrix math
- **features_test.bas** – Assorted features

**Note:** **platformer.bas** uses high-level runtime opcodes (INITGRAPHICS, CREATEPHYSICSBODY2D, etc.) that may not match the current VM; use **box2d_demo.bas** and raylib + BOX2D.* for 2D physics games.
