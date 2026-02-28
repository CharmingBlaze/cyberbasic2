# Changelog

All notable changes to CyberBasic are documented here. The project follows a single main branch; version tags may be added for releases.

---

## [Unreleased] – release preparation

### Physics, UI, and audio (full implementation)

- **Box2D:** All joint types implemented (Revolute, Prismatic, Weld, Rope, Pulley, Gear, Wheel); joint ID storage; **SetJointLimits2D**, **SetJointMotor2D**, **DestroyJoint2D**. Distance joint returns jointId.
- **Bullet:** Body properties implemented and used in Step and collision: friction, restitution, linear/angular damping, kinematic, gravity scale, linear/angular factor, CCD. Setters: SetFriction3D, SetRestitution3D, SetDamping3D, SetKinematic3D, SetGravity3D, SetLinearFactor3D, SetAngularFactor3D, SetCCD3D. 3D constraint joints remain stubs.
- **UI (raygui):** **GuiLoadStyle**(filePath), **GuiLoadStyleDefault**(), **GuiSetStyle**(controlId, propertyId, value), **GuiGetStyle**(controlId, propertyId) for theme and layout.
- **Audio:** Documented that stream callbacks requiring C function pointers are not exposed from BASIC; use **UpdateAudioStream** to push samples.
- Documentation: API_REFERENCE, COMMAND_REFERENCE, 2D_PHYSICS_GUIDE, 3D_PHYSICS_GUIDE, GUI_GUIDE, README, and GAME_DEVELOPMENT_GUIDE updated.

### Cleanup and documentation

- Moved root-level ad-hoc test scripts (`test_*.bas`) and the raylib diagnostic (`test_raylib_window.go`) into `deprecated/` with a README. These are not part of the main build or test suite.
- Updated `.gitignore` to exclude local artifact files (`out.txt`, `e2.txt`, `err.txt`, `o1.txt`, `o2.txt`, `e1.txt`, `*.log`).
- README rewritten to present the project as a modern Go-based engine: technical identity, C++ to Go rationale (maintainability, build speed, contributor experience), and a full table of integrated systems (Raylib, Box2D, Bullet, net, GUI, events, terrain, water, vegetation, world, navigation, indoor, ECS, std, sql, procedural).
- GETTING_STARTED.md updated to point to the Go-based architecture and the main README for rationale.
- This changelog added for release visibility.

### Architecture (current)

- **Compiler:** Go lexer, parser, codegen (statements, expressions, calls, util). Modular layout; no C++.
- **VM:** Bytecode VM with stack, globals, foreign calls, fibers, render queues (2D/3D/GUI). Packages: vm, vm_ops, vm_run, vm_foreign, vm_fibers, bytecode, runtime_iface.
- **Bindings:** raylib (graphics, input, audio, 2D layers/camera/backgrounds, 3D, hybrid flush), box2d, bullet, game, scene, net, ecs, terrain, water, vegetation, objects, world, navigation, indoor, procedural, std, sql. All registered from `main.go`.
- **Default build:** `go build -o cyberbasic .` produces one binary; no C compiler required. Optional C engine in `engine/` for custom builds.

---

## Older history

For earlier work (language features, 2D/3D engine systems, physics, multiplayer, GUI, terrain/water/vegetation, navigation, indoor, streaming, editor stubs, documentation), see the git history and the [Roadmap](ROADMAP.md).
