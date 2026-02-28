# Changelog

All notable changes to CyberBasic are documented here. The project follows a single main branch; version tags may be added for releases.

---

## [Unreleased] – release preparation

### Physics API: flat names only (namespace removal)

- **Physics API simplified:** 2D and 3D physics now use **flat names** only: **CreateWorld2D**, **Step2D**, **CreateBody2D**, **GetPositionX2D** / **GetPositionY2D**, etc., and **CreateWorld3D**, **Step3D**, **CreateBox3D**, **GetPositionX3D**, etc. The **BOX2D.*** and **BULLET.*** namespaces are no longer registered in the VM.
- **Backward compatibility:** The compiler rewrites legacy `BOX2D.*` and `BULLET.*` calls to the corresponding flat names at compile time, so existing scripts continue to work.
- **Box2D flat API:** Added **CreateBody2D**, **DestroyBody2D**, **GetBodyCount2D**, **GetBodyId2D**, **CreateBodyAtScreen2D** to the flat API; removed all BOX2D.* VM registrations.
- **Bullet flat API:** Added **SetWorldGravity3D**, **DestroyBody3D**, **RayCastFromDir3D**; removed all BULLET.* VM registrations. **3D constraint joints** (CreateHingeJoint3D, etc.) remain stubbed in the pure-Go engine and are documented as not implemented.
- **Gogen:** Generated Go code now emits flat-style physics calls (e.g. `box2d.CreateWorld2D`, `bullet.Step3D`) instead of namespaced names.
- **Documentation:** API_REFERENCE, COMMAND_REFERENCE, 2D_PHYSICS_GUIDE, 3D_PHYSICS_GUIDE, FAQ, GAME_DEVELOPMENT_GUIDE, and related docs updated to use flat names only; one-line deprecation note for BOX2D/BULLET namespaces.
- **Examples and templates:** All BOX2D.* and BULLET.* usage in examples and templates replaced with flat names.

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
