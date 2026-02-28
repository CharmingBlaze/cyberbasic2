# Getting Started with CyberBasic

Install, build, and run your first CyberBasic program.

## Table of Contents

1. [Prerequisites](#prerequisites)
2. [Build the compiler](#build-the-compiler)
3. [Run your first program](#run-your-first-program)
4. [Where is the compiler?](#where-is-the-compiler)
5. [Next steps](#next-steps)
6. [Troubleshooting](#troubleshooting)

---

## Prerequisites

- **Go 1.19+** – [golang.org/dl](https://golang.org/dl/)
- No C compiler required for normal use. The runtime is written in Go and uses raylib-go, so you get a single binary. An optional **C engine** in `engine/` exists for custom builds; the recommended path is the pure-Go build. See the main [README](../README.md) for the rationale behind the Go-based architecture.

## Build the compiler

From the project root:

```bash
go build -o cyberbasic .
```

On Windows you get `cyberbasic.exe`; on Unix you get `cyberbasic`. The executable is the compiler and runtime: it compiles `.bas` files to bytecode and runs them.

## Run your first program

**Minimal (no window):**

```bash
./cyberbasic examples/hello_world.bas
```

**First game (window + WASD circle):**

```bash
./cyberbasic examples/first_game.bas
```

This opens a window, and you move a circle with WASD. See [Cheatsheet](../CHEATSHEET.md) for the “first 10 lines” of a 2D and 3D game.

## Exit codes

When running a script, the compiler exits with:

- **0** – Success (compilation and, if not `--compile-only`, program run completed).
- **1** – Compile or file error (file not found, lex/parse error, compilation error).
- **2** – Runtime error (program ran but hit an error during execution).

Use these in scripts or CI to detect failure (e.g. `./cyberbasic game.bas || echo "Build or run failed"`).

## Where is the compiler?

After building, the compiler is the binary you produced:

- **Default:** `./cyberbasic` (or `cyberbasic.exe` on Windows) in the current directory.
- To use from anywhere, add the project root (or a directory containing `cyberbasic`) to your `PATH`.
- Run `./cyberbasic --help` for options; use `./cyberbasic --list-commands` to print built-in command names. Use `./cyberbasic --lint your.bas` (or `--compile-only`) to check your program without running it. Full reference: [Command Reference](COMMAND_REFERENCE.md) and [API Reference](../API_REFERENCE.md).

## Next steps

| Goal | Where to go |
|------|-------------|
| **Syntax at a glance** | [Quick Reference](QUICK_REFERENCE.md) |
| **Full language rules** | [Language Spec](../LANGUAGE_SPEC.md) |
| **All docs** | [Documentation Index](DOCUMENTATION_INDEX.md) |
| **First 2D game** | [Cheatsheet](../CHEATSHEET.md) → [2D Graphics Guide](2D_GRAPHICS_GUIDE.md) |
| **First 3D game** | [Cheatsheet](../CHEATSHEET.md) → [3D Graphics Guide](3D_GRAPHICS_GUIDE.md) |
| **Hybrid loop** (auto physics + render) | Define `update(dt)` and `draw()`; see [Program Structure](PROGRAM_STRUCTURE.md#hybrid-updatedraw-loop) |
| **Game loop + input + physics** | [Game Development Guide](GAME_DEVELOPMENT_GUIDE.md) |
| **Minimal game templates** | [templates/2d_game.bas](../templates/2d_game.bas), [templates/3d_game.bas](../templates/3d_game.bas) – see [templates/README](../templates/README.md) |
| **See all commands** | Run `cyberbasic --list-commands` or open [Command Reference](COMMAND_REFERENCE.md) |

Start with [Quick Reference](QUICK_REFERENCE.md) and [examples/first_game.bas](../examples/first_game.bas), then explore the [Documentation Index](DOCUMENTATION_INDEX.md).

## Troubleshooting

- **Compiler not found:** Add the project root (or the directory containing `cyberbasic` / `cyberbasic.exe`) to your `PATH`, or run from that directory.
- **"open cb.exe: The process cannot access the file":** The executable is in use (e.g. another terminal is running it). Close the other process or build to a different name (e.g. `go build -o cyberbasic_new .`).
- **Parse error / "expected identifier":** Check variable and symbol names; avoid names that look like keywords. Rename if in doubt (e.g. `msg` → `received` if it clashes in a given scope).
- More: see [FAQ](FAQ.md) and [Troubleshooting](TROUBLESHOOTING.md).

## See also

- [Documentation Index](DOCUMENTATION_INDEX.md)
- [Quick Reference](QUICK_REFERENCE.md)
- [Program Structure](PROGRAM_STRUCTURE.md) – hybrid update/draw loop
- [Game Development Guide](GAME_DEVELOPMENT_GUIDE.md)
