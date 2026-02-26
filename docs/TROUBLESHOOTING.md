# Troubleshooting

Common errors and how to fix them.

## Build and run

**Compiler not found / "cyberbasic" is not recognized**

- The compiler is the binary you built (e.g. `cyberbasic` or `cyberbasic.exe`). Either run it from the directory where it lives, or add that directory to your **PATH**.
- On Windows: add the project root (or the folder containing `cyberbasic.exe`) to the system or user PATH, or use the full path: `.\cyberbasic examples\first_game.bas`.

**"open cb.exe: The process cannot access the file"**

- The executable is locked because another process is using it (e.g. a previous run still open, or another terminal running the same binary). Close the other process, or build to a different output name: `go build -o cyberbasic_new .`

**Build fails with CGO or raylib errors**

- Default build uses raylib-go (Go). If you see C compiler or raylib C errors, check that Go can find the right packages. For pure-Go (no CGO), use `CGO_ENABLED=0 go build -o cyberbasic .` (some features like full raygui require CGO).

## Parsing and compile errors

**"expected identifier" or parse error at a variable name**

- Some names can be parsed as keywords or in a way that breaks the parser. Try renaming the variable (e.g. `msg` → `received`, or avoid single-letter or keyword-like names in that scope).
- Ensure you are not using a reserved word as a variable name.

**"undefined" or "unknown function" at runtime**

- The name is not registered in the VM. Use the **flat API** names from the docs: InitWindow, ClearBackground, WindowShouldClose, etc. Legacy names like INITGRAPHICS or RL.InitWindow may not be registered; see [API Reference](../API_REFERENCE.md) and [Examples README](../examples/README.md) for current API.

## Runtime and examples

**Window opens and closes immediately / nothing draws**

- Make sure you have a game loop: `WHILE NOT WindowShouldClose() ... WEND` (or REPEAT UNTIL), and that you call drawing functions (e.g. ClearBackground, DrawCircle) inside the loop. If you use the hybrid loop, define **draw()** and leave the loop body empty.

**sql_demo.bas fails to open database**

- Run from the project root and avoid redirecting output or running from a restricted directory. The demo creates a SQLite file in the current working directory.

**multiplayer_server.bas / multiplayer_client.bas don’t connect**

- Start the **server** first, then run the **client** in another terminal. Use the same port and (for client) the correct host (e.g. localhost or 127.0.0.1).

## See also

- [Getting Started](GETTING_STARTED.md) – build and first run
- [FAQ](FAQ.md) – hybrid vs manual loop, multi-window, API
- [Documentation Index](DOCUMENTATION_INDEX.md)
