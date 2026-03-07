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

These are still real limitations after the implementation pass:

- Shadows are now implemented as a single-light directional shadow-map path. Point-light shadows, spot-light shadows, and cascaded shadow maps are still future work.
- GLTF `KHR_lights_punctual` import is still not wired into runtime light creation, so level lighting must still be created explicitly in game code.
- The asset cache is only partially integrated. `LoadLevel`, `LoadLevelWithHierarchy`, and `LoadPrefab` still bypass part of the shared parsed-model cache flow.
- 2D physics now reports its backend as `bytearena-box2d` / `authoritative`, while 3D physics reports `purego-fallback` / `fallback`; there is still no native Bullet backend in this checkout.
- The pure-Go Bullet layer still approximates some non-box/non-sphere behavior, especially capsule overlap and swept queries.
- Multiplayer does not yet provide lockstep, rollback, prediction, matchmaking, or interest management.

## How To Read The Roadmap

- Use `ROADMAP.md` for future priorities.
- Use this file for current implementation truth.
- Use feature guides in `docs/` for API-level behavior and examples.

## Verification Summary

The remediation pass was validated by:

- Formatting the touched Go files with `gofmt`
- Building/testing the affected Go packages with `go test`
- Updating docs so public guidance matches current behavior instead of planned behavior
