# CyberBasic Roadmap

This document is the project’s priority list and planned improvements. Contributions are welcome—see the repo [issues](https://github.com/CharmingBlaze/cyberbasic2/issues) and open a discussion or PR for any item below.

---

## High priority

- **Clear roadmap** – This document; keep it updated as priorities shift.
- **Better onboarding** – Hello World in README, one-command “run the demo,” screenshot or GIF of CyberBasic running.
- **Working UI system** – Expand BeginUI, Label, Button (and related widgets) beyond stubs; minimal immediate-mode UI or raygui integration for in-game editors, debug overlays, menus, ECS inspection.
- **Debugger / inspector** – Breakpoints, step/next, variable watch, call stack view, VM bytecode dump, ECS entity/component inspector (even text-based would be a big improvement).
- **More complete physics** – Box2D and Bullet: joints, body properties, collision callbacks, raycasts (fill remaining stubs).

---

## Language and VM

- **Better error messages** – Line/column reporting, “did you mean?” suggestions, runtime stack traces, type mismatch explanations.
- **String and array improvements** – Dynamic arrays, slicing, built-in string functions (LEFT$, MID$, INSTR, etc.), string interpolation.
- **Non-blocking WaitSeconds and fiber scheduler** – Smooth timers, async-style events, parallel coroutines (scheduler already yields; polish and docs).
- **Optional** – JIT, GC tuning, REPL (line → compile → run, variables, expressions).

---

## Engine and API

- **2D engine and world systems (implemented)** – Layers, parallax, 2D camera by ID, backgrounds, tilemap create/save/load/fill, sprite animation and batching, 2D particle emitters, 2D culling and atlas, scene save/load 2D, physics 2D gravity/raycast/layer collision; terrain/water physics (collision, friction, bounce, density, drag, buoyancy stub); vegetation (tree/grass wind, collision radius); weather, fire/smoke, environment stubs; navigation (NavGrid, NavMesh, NavAgent stubs); indoor (rooms, doors, triggers, interactables stubs); streaming (chunk load/unload stubs); editor (EditorEnable, brush stubs); decals (stubs). See COMMAND_REFERENCE and API_REFERENCE.
- **Physics completeness** – Same as high-priority physics: joints, body props, collision callbacks, raycasts.
- **Asset pipeline** – Sprite sheet packer, tilemap loader, model importer helpers, audio streaming helpers (alongside existing ExportImageAsCode, ExportFontAsCode, ExportWaveAsCode).

---

## Documentation and examples

- **More examples** – Platformer, top-down shooter, simple 3D scene with physics, UI demo (once UI is solid), ECS example with systems.
- **“Why CyberBasic?”** – Short explanation: why BASIC, why Go, why raylib, what makes CyberBasic unique (see README or this section as we add it).
- **README screenshot/GIF** – Add a screenshot or short GIF of CyberBasic running (e.g. first_game.bas or GUI) to improve first impressions.

---

## Testing and CI

- **GitHub Actions** – Run tests on push; VM and parser/lexer regression or golden tests.

---

## Developer experience

- **VSCode extension** – Syntax highlighting, basic indentation, run/compile command.
- **REPL** – Minimal REPL: type a line → compile → run; print variables; evaluate expressions.
- **Package / module distribution** – Simple package manager (e.g. JSON-based index) for community modules, reusable libraries, versioning, dependency resolution (longer-term).

---

## Summary

| Area              | Focus |
|-------------------|--------|
| Onboarding        | Hello World in README, one-command demo, screenshot/GIF |
| UI                | Real immediate-mode UI or raygui for menus, debug, ECS |
| Debugging         | Breakpoints, step, watch, stack, bytecode dump, ECS inspector |
| Physics           | Complete Box2D/Bullet joints, callbacks, raycasts |
| Errors            | Line/column, “did you mean?”, stack traces, type hints |
| Strings/arrays    | Dynamic arrays, slicing, LEFT$/MID$/INSTR, interpolation |
| Scheduler         | Polish WaitSeconds and fiber scheduler |
| Examples          | Platformer, shooter, 3D scene, UI demo, ECS |
| CI                | GitHub Actions, VM/parser tests |
| DX                | VSCode extension, REPL, package distribution |
