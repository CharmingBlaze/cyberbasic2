# Deprecated / moved files

Files in this folder were moved from the repository root during cleanup. They are not part of the main build or test suite.

- **test_*.bas** – Ad-hoc or one-off test scripts; use `examples/` for runnable demos.
- **test_raylib_window.go** – Diagnostic: run with `go run test_raylib_window.go` to verify raylib/OpenGL on the machine; excluded from main build via `// +build ignore`.

Do not rely on these for development. Prefer `examples/first_game.bas`, `examples/hybrid_update_draw_demo.bas`, and the rest of `examples/` for current API usage.
