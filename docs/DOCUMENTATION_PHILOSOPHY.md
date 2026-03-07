# CyberBASIC2 Philosophy and Design Principles

This document defines the core philosophy that guides CyberBASIC2's design, documentation, and implementation. All subsystems, APIs, and docs should align with these principles.

---

## 1. DBP-Style Simplicity

**Principle:** Games should be easy to write. Programmers should spend time on game logic, not engine boilerplate.

- **Zero boilerplate when possible:** OnStart/OnUpdate/OnDraw, SYNC, UseUnifiedRenderer—no InitWindow, no manual WHILE loop required.
- **Sensible defaults:** Window opens at 800×600, 60 FPS, VSync on. Physics steps at 1/60. You override only what you need.
- **Flat, prefix-based commands:** `InitWindow`, `DrawCircle`, `LoadObject`—no nested namespaces in common use. Optional `RL.*` and `GAME.*` for clarity.
- **Familiar BASIC feel:** VAR, LET, DIM, IF/THEN, FOR/NEXT, SUB/FUNCTION. Dot notation for types. No surprises.

**Documentation:** Every guide shows the simplest path first. Advanced options come after.

---

## 2. Flat, Prefix-Based API

**Principle:** Commands are flat by default. Namespaces are optional qualifiers, not requirements.

- **Raylib:** `InitWindow`, `ClearBackground`, `DrawCircle`—same names as raylib C API. `RL.InitWindow` is an alias.
- **Physics:** Flat names: `CreateWorld2D`, `Step3D`, `CreateSphere3D`. No `BULLET.` required.
- **Game helpers:** `GAME.MoveWASD`, `GAME.CameraOrbit`—namespace when it clarifies domain.
- **DBP-style:** `LoadObject`, `PositionObject`, `DrawObject`—integer IDs, familiar to DBP users.

**Documentation:** Examples use flat names unless the namespace adds clarity.

---

## 3. Sensible Defaults

**Principle:** Out of the box, things work. Override for power users.

| Subsystem | Default | Override |
|-----------|---------|----------|
| Window | 800×600, VSync on | InitWindow, SetConfigFlags |
| FPS | 60 | SetTargetFPS |
| Physics step | 1/60 | FixedUpdate(rate) |
| Clear color | Black | SetClearColor |
| 2D camera | Identity | SetCamera2D |
| 3D camera | Default FOV, near/far | SetCamera3D, SetCameraFOV |
| Shadows | Single directional, medium | SetShadowQuality, EnableShadows |

**Documentation:** Each subsystem doc lists defaults and override points.

---

## 4. Deterministic Multiplayer Foundation

**Principle:** The engine is built for deterministic simulation. Multiplayer code should be predictable.

- **Fixed-step simulation:** `FixedUpdate(rate)`, `OnFixedUpdate(label$)`, `FixedDeltaTime()`. Physics and gameplay advance on a stable timestep.
- **Frame delta clamped:** Prevents runaway catch-up after stalls.
- **Single-consumption messages:** Network messages consumed once; no duplicate delivery.
- **Authoritative server pattern:** One authority; clients send input; server broadcasts state.

**Current scope:** TCP transport, RPC, SyncEntity. Lockstep, rollback, prediction are roadmap items.

**Documentation:** Multiplayer docs explain fixed-step usage and deterministic patterns.

---

## 5. Contributor-Friendly Architecture

**Principle:** The codebase is modular. Contributors can add bindings, fix bugs, or extend subsystems without touching unrelated code.

- **Bindings:** One package per domain (raylib, box2d, bullet, dbp, net, etc.). Register foreign functions; no parser changes for new commands.
- **Runtime:** Clear separation: VM, runtime, renderer, time, fixed-step. Each has a single responsibility.
- **Tests:** VM, parser, lexer have tests. New bindings should have at least a smoke test.

**Documentation:** MODULES.md, STYLE.md, and subsystem docs include contributor notes.

---

## 6. Modern Engine Internals

**Principle:** Use modern, maintained libraries. Avoid legacy or unmaintained dependencies.

- **Graphics:** Raylib (raylib-go). 2D/3D, shaders, render textures, multiple cameras.
- **2D physics:** Box2D (bytearena/box2d). Authoritative, full-featured.
- **3D physics:** Bullet-shaped API. Shipped pure-Go fallback; native Bullet optional.
- **Assets:** GLTF/GLB, OBJ for 3D. PNG, JPG for textures. Aseprite JSON for sprites. Blender for authoring.
- **Shadows:** Directional shadow maps. Low/medium/high quality presets.

**Documentation:** Asset pipeline, Blender workflow, Aseprite workflow document authoring and loading.

---

## 7. No Placeholders in Production Docs

**Principle:** Documentation describes what exists. Future work goes in ROADMAP.md or "Known limitations" sections.

- **Implemented:** Full description, examples, edge cases.
- **Partial:** "Current status" section with what works and what does not.
- **Not implemented:** Listed in ROADMAP_IMPLEMENTATION.md "Known Remaining Gaps" or ROADMAP.md. No vague "TODO" in user-facing docs.

**Documentation:** Every doc ends with a clear "Status" or "Limitations" section when applicable.

---

## 8. Real Game Developer Workflows

**Principle:** Docs are written for developers building real games.

- **Tutorials:** Step-by-step from empty file to playable game.
- **Examples:** Runnable, copy-pasteable. first_game.bas, 2d_game.bas, 3d_game.bas.
- **Patterns:** Common patterns (WASD movement, camera follow, collision response) shown with code.
- **Troubleshooting:** Common errors, fixes, and where to look.

**Documentation:** Learning Path, Tutorials, and Game Development Guide are workflow-centric.

---

## See Also

- [Documentation Index](DOCUMENTATION_INDEX.md)
- [Getting Started](GETTING_STARTED.md)
- [ROADMAP.md](../ROADMAP.md)
- [Documentation Style](STYLE.md)
