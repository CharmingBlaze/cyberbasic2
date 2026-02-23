# ECS (Entity-Component System) Guide

CyberBasic provides an ECS via the **ECS.*** binding. Use it to organize many game objects (enemies, projectiles, props) with shared logic: create a world, create entities, add components (Transform, Sprite, Health), and query entities by component type.

## Setup

ECS is always available; no extra setup. Call **ECS.CreateWorld** to get a world ID, then create entities and add components.

## API Summary

| Function | Description |
|----------|-------------|
| **ECS.CreateWorld** | () → worldId (string) |
| **ECS.DestroyWorld** | (worldId) |
| **ECS.CreateEntity** | (worldId) → entityId (string) |
| **ECS.DestroyEntity** | (worldId, entityId) |
| **ECS.AddComponent** | (worldId, entityId, componentType [, args...]) |
| **ECS.HasComponent** | (worldId, entityId, componentType) → boolean |
| **ECS.RemoveComponent** | (worldId, entityId, componentType) |
| **ECS.GetTransformX/Y/Z** | (worldId, entityId) → number |
| **ECS.SetTransform** | (worldId, entityId, x, y, z) |
| **ECS.GetHealthCurrent/Max** | (worldId, entityId) → number |
| **ECS.QueryCount** | (worldId, componentType1 [, componentType2...]) → count |
| **ECS.QueryEntity** | (worldId, componentType, index) → entityId or empty string |

**Component types:** **Transform** (x, y, z), **Sprite** (textureId, visible), **Health** (current, max). Add with **ECS.AddComponent(worldId, entityId, "Transform", x, y, z)** (and similarly for Sprite, Health).

## Example

```basic
VAR wid = ECS.CreateWorld()
VAR e1 = ECS.CreateEntity(wid)
ECS.AddComponent(wid, e1, "Transform", 100, 200, 0)
ECS.AddComponent(wid, e1, "Health", 100, 100)

VAR e2 = ECS.CreateEntity(wid)
ECS.AddComponent(wid, e2, "Transform", 300, 150, 0)
ECS.AddComponent(wid, e2, "Sprite", "hero", 1)

// Query all entities with Transform
VAR n = ECS.QueryCount(wid, "Transform")
VAR i = 0
FOR i = 0 TO n - 1
  VAR eid = ECS.QueryEntity(wid, "Transform", i)
  VAR x = ECS.GetTransformX(wid, eid)
  VAR y = ECS.GetTransformY(wid, eid)
  // draw or update entity...
NEXT i

ECS.DestroyWorld(wid)
```

## When to use ECS

- Many similar objects (bullets, enemies, particles) that share behavior but different data.
- You want to query "all entities with Transform and Health" and iterate.
- You prefer composition (add/remove components) over deep inheritance.

For a runnable demo see [examples/ecs_demo.bas](../examples/ecs_demo.bas) (if present). Full API list: [API_REFERENCE.md](../API_REFERENCE.md).
