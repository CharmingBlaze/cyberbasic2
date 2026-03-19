# Roadmap Implementation Status

This document records the current implementation status behind the roadmap audit/remediation pass. It is intentionally about the code that exists now, not the aspirational feature list in `ROADMAP.md`.

## Completed In This Pass

### Runtime and networking

- Fixed-step runtime support is now wired through the main loop.
- `FixedUpdate(rate)` updates the shared fixed timestep.
- `OnFixedUpdate(label$)` runs a BASIC callback on each fixed step.
- Frame delta is clamped before entering the accumulator to avoid runaway catch-up after stalls.
- Fixed-step catch-up is capped to avoid a spiral-of-death frame.
- TCP reader handling no longer treats short idle read timeouts as disconnects.
- Message delivery is single-consumption instead of being duplicated between event processing and `Receive()`.

### Rendering and scene graph

- Hierarchical object transforms now compose full position, rotation, and scale correctly.
- Imported GLTF node hierarchies preserve transform-only ancestors when using level hierarchy loading.
- Imported PBR values are preserved unless the game explicitly overrides them.
- 3D culling now uses bounded-object math instead of a point-only approximation.
- The engine now has a first-pass directional shadow-map path with low/medium/high quality presets for weaker to stronger hardware tiers, and the global shadow path starts enabled by default.

### Gameplay systems

- Animation crossfade now tracks both source and destination clip progress during a blend.
- Skeletal blends use interpolated pose data when bind-pose/bone data exists.
- `CreateCharacterController` now preserves the intended capsule dimensions.
- `GAME.MoveWASD` uses direct velocity control for more stable character motion.
- `Spherecast` now performs a swept-sphere style broad-phase test against inflated body bounds.
- The simple DBP-style 2D physics wrappers now route through internal aliases correctly instead of recursing, and the documented collider helper arities now work as shipped.
- The shipped 3D fallback physics path now exposes explicit backend/feature queries and returns clear unsupported-feature errors for stubbed joints, torque-only APIs, heightmaps, compounds, and DBP mesh colliders.

## Documentation Corrected In This Pass

- `3D_GAME_API.md`
- `3D_GRAPHICS_GUIDE.md`
- `3D_PHYSICS_GUIDE.md`
- `ASSET_PIPELINE.md`
- `COMMAND_REFERENCE.md`
- `MULTIPLAYER.md`
- `TUTORIAL_MULTIPLAYER.md`

## Known Remaining Gaps

These are real limitations. Each has a clear work item for contributors.

### Shadows

| Gap | Current State | Action for Contributors |
|-----|---------------|-------------------------|
| Point-light shadows | **Done** | Single perspective frustum aimed at scene center |
| Spot-light shadows | **Done** | Perspective frustum per spot light |
| Cascaded shadow maps | **Done** (API) | SetShadowCascades, ShadowCascadeCount; SetShadowQuality sets cascade count (low=1, medium=3, high=4) |

**File:** `compiler/runtime/renderer/shadow.go`, `compiler/bindings/dbp/dbp_lighting.go`

### Level Lighting

| Gap | Current State | Action for Contributors |
|-----|---------------|-------------------------|
| GLTF KHR_lights_punctual | **Done** | Parse lights and node extensions in `compiler/bindings/model/gltf.go`; lights appear in LoadLevel via BuildModel |

**File:** `compiler/bindings/model/gltf.go`, `compiler/bindings/dbp/dbp_model.go`

### Asset Cache

| Gap | Current State | Action for Contributors |
|-----|---------------|-------------------------|
| LoadLevel cache | **Done** | Routed through `assets.LoadModelForBuild`; `UnloadLevel` calls `UnloadModelForBuild` |
| LoadLevelWithHierarchy cache | **Done** | Same as LoadLevel |
| LoadPrefab cache | **Done** | `LoadPrefab` uses `assets.LoadModelForBuild`; `DeletePrefab` calls `UnloadModelForBuild` |
| Texture reuse in BuildModel | **Done** | `BuildModel` uses `resources.LoadTexture`/`GetTexture`; levels track `texturePaths` and call `UnloadTexture` on unload |

**Files:** `compiler/bindings/dbp/dbp_level.go`, `compiler/bindings/dbp/dbp_prefab.go`, `compiler/bindings/dbp/dbp_model.go`, `compiler/runtime/assets/assets.go`, `compiler/runtime/resources/manager.go`

### 3D Physics

| Gap | Current State | Action for Contributors |
|-----|---------------|-------------------------|
| Native Bullet backend | Not in checkout | Optional CGO build; wire bullet C lib; see `BulletNativeAvailable()` |
| ApplyTorque3D, ApplyTorqueImpulse3D | **Done** | Implemented in pure-Go fallback |
| PointToPoint, Fixed joints | **Done** | Implemented in pure-Go fallback; `BulletJointsAvailable()` returns 1 |
| CreateStaticMesh3D | **Done** (AABB) | Loads OBJ, uses AABB for broad phase; exact triangle narrow-phase deferred |
| Hinge, Slider, ConeTwist joints | Stubbed; return error | Implement CreateHingeJoint3D, CreateSliderJoint3D in pure-Go or native |
| Capsule/swept queries | Approximated | Improve overlap and spherecast math in `compiler/bindings/bullet` |

**Files:** `compiler/bindings/bullet/bullet.go`, `compiler/bindings/bullet/*.go`

### Multiplayer

| Gap | Current State | Action for Contributors |
|-----|---------------|-------------------------|
| Lockstep | Not implemented | Design: fixed tick rate; clients send input per tick; server broadcasts state |
| Rollback | Not implemented | Requires snapshot save/restore; resimulate from last confirmed state |
| Prediction | Not implemented | Client-side prediction from input; reconcile with server state |
| Matchmaking | Not implemented | External service or custom lobby |
| Interest management | Not implemented | Filter SyncEntity/Replicate by distance or relevance |

**File:** `compiler/bindings/net/net.go`, `docs/MULTIPLAYER_DESIGN.md`

---

## See Also

- [ROADMAP.md](../ROADMAP.md) — Forward-looking priorities
- [Documentation Philosophy](DOCUMENTATION_PHILOSOPHY.md)
- [Documentation Index](DOCUMENTATION_INDEX.md)

## How To Read The Roadmap

- Use `ROADMAP.md` for future priorities.
- Use this file for current implementation truth.
- Use feature guides in `docs/` for API-level behavior and examples.

## Verification Summary

The remediation pass was validated by:

- Formatting the touched Go files with `gofmt`
- Building/testing the affected Go packages with `go test`
- Updating docs so public guidance matches current behavior instead of planned behavior
