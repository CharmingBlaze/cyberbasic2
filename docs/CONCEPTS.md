# Concepts: 2D vs 3D vs Physics

One-page conceptual overview for new users. Use this to understand which subsystems to use and when to combine them.

---

## Table of Contents

1. [2D graphics](#2d-graphics)
2. [3D graphics](#3d-graphics)
3. [2D physics (Box2D)](#2d-physics-box2d)
4. [3D physics (Bullet-style)](#3d-physics-bullet-style)
5. [When to combine](#when-to-combine)
6. [See also](#see-also)

---

## 2D graphics

2D rendering for platformers, top-down games, and UI overlays.

- **Sprites and shapes:** Draw rectangles, circles, lines, triangles, and textures.
- **Tilemaps:** Draw grids of tiles for levels.
- **2D camera:** Pan and zoom the view.
- **Text:** Draw labels, scores, and HUD text.

Use 2D when your game is viewed from the side (platformer) or top-down (puzzle, strategy). See [2D Graphics Guide](2D_GRAPHICS_GUIDE.md).

---

## 3D graphics

3D rendering for FPS, third-person, and open worlds.

- **Models and meshes:** Load GLTF/OBJ models, draw cubes, spheres, and meshes.
- **3D camera:** Position, target, and orbit cameras.
- **Lighting:** Directional, point, and spot lights; PBR materials.
- **Level loading:** Load full scenes with hierarchy, materials, and textures.
- **World features:** Water, terrain (heightmap, procedural), sky, clouds.

Use 3D when your game has depth, perspective, and a 3D camera. See [3D Graphics Guide](3D_GRAPHICS_GUIDE.md), [Level Loading](LEVEL_LOADING.md), [World, Water, Terrain](WORLD_WATER_TERRAIN.md).

---

## 2D physics (Box2D)

Box2D for rigid bodies, joints, and collision in 2D.

- **Bodies:** Dynamic, kinematic, and static bodies with mass and velocity.
- **Joints:** Distance, revolute, prismatic, weld, and more.
- **Collision:** Contact callbacks, collision groups, sensors.
- **Forces:** Apply impulses and linear velocity.

Use 2D physics when you need realistic movement, gravity, bouncing, or constraints in a 2D game. See [2D Physics Guide](2D_PHYSICS_GUIDE.md).

---

## 3D physics (Bullet-style)

Bullet-style 3D physics for characters, projectiles, and simple collision.

- **Bodies:** Rigid bodies with mass, velocity, and collision shapes.
- **Raycast:** Cast rays for line-of-sight, shooting, or ground detection.
- **Forces:** Apply impulses and linear velocity.
- **Joints:** Optional; limited support compared to 2D.

Use 3D physics when you need gravity, collision, or raycasting in a 3D world. See [3D Physics Guide](3D_PHYSICS_GUIDE.md).

---

## When to combine

| Combination | Use case |
|-------------|----------|
| **2D + Box2D** | Platformer, top-down shooter, puzzle with physics |
| **3D + Bullet** | FPS, third-person, character movement and projectiles |
| **3D world + water + terrain** | Open world, outdoor scene: load level, add water, add terrain, sky |

You can mix 2D and 3D in the same game (e.g. 3D world with 2D HUD). The hybrid loop (update/draw) and unified renderer handle the draw order. See [Game Development Guide](GAME_DEVELOPMENT_GUIDE.md).

---

## See also

- [Getting Started](GETTING_STARTED.md)
- [2D Graphics Guide](2D_GRAPHICS_GUIDE.md)
- [3D Graphics Guide](3D_GRAPHICS_GUIDE.md)
- [2D Physics Guide](2D_PHYSICS_GUIDE.md)
- [3D Physics Guide](3D_PHYSICS_GUIDE.md)
- [Documentation Index](DOCUMENTATION_INDEX.md)
