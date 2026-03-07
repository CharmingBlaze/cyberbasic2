# Documentation style guide

Conventions for CyberBASIC2 documentation so all docs feel consistent.

## Heading levels

- **One H1** per document: the main title (e.g. "Getting Started with CyberBASIC2").
- **##** for major sections (e.g. "Prerequisites", "Build the compiler").
- **###** for subsections when needed.
- Do not skip levels (e.g. do not go from ## to ####).

## Table of Contents

- Add a **Table of Contents** for documents with **4 or more** major sections (##).
- Use anchor links: `[Section name](#section-name)` (lowercase, spaces as hyphens).
- Place the TOC after the intro paragraph, before the first ##.

## Code blocks

- Use **` ```basic `** for all CyberBASIC2 source code.
- Use **` ```bash `** (or no tag) for shell commands when trivial.
- Keep snippets short and runnable where possible.

## Cross-links

- End major docs with a **"See also"** or **"Related"** section.
- Link to [Documentation Index](DOCUMENTATION_INDEX.md) and to the next-step doc (e.g. Quick Reference, Game Development Guide).
- Use relative links within the repo: `[Link text](FILENAME.md)` or `[Link](../PATH.md)` from docs/.

## Terminology

- **Project name:** Use **CyberBASIC2** for the project and product. Use **CyberBASIC** only when referring to the language generically in prose.
- **Game loop:** Use "game loop" when referring to `WHILE NOT WindowShouldClose() ... WEND` or `REPEAT ... UNTIL WindowShouldClose()`.
- **Delta time:** Use "delta time" for frame elapsed time; refer to **DeltaTime()** or **GetFrameTime()**.
- **Flat API:** Prefer flat names in examples: **InitWindow**, **ClearBackground**, **DrawCircle** (not `RL.InitWindow`). Use namespaced forms only when documenting that API (e.g. **GAME.***). Physics uses flat names (CreateWorld2D, Step3D, …).

## API reference

- **API_REFERENCE.md:** One table per section; one row per command; columns: Command, Arguments, Returns, Description. When adding a binding, add a row to the matching section.

## Include syntax

- Document the directive as **`#include "path"`**. Optionally mention **IMPORT "path"** as an alias if supported.
- Do not use "INCLUDE" without the hash (the compiler expects `#include` or `IMPORT`).

## Subsystem documentation template

When documenting a subsystem (e.g. Asset Pipeline, Multiplayer, Physics), include these sections for production-ready clarity:

1. **Purpose** — What the subsystem does and why it exists.
2. **Architecture** — High-level flow, packages involved, data flow.
3. **API Surface** — Commands/functions with args, returns, description.
4. **Defaults** — Default values, behaviors, override points.
5. **Edge Cases** — Error conditions, boundary behavior, gotchas.
6. **Performance Considerations** — When to preload, batch, or avoid.
7. **Multiplayer / Determinism** — How the subsystem interacts with networked games.
8. **Contributor Notes** — File paths, extension points, testing.

See [Asset Pipeline](ASSET_PIPELINE.md) and [Rendering and the Game Loop](RENDERING_AND_GAME_LOOP.md) for examples.

## See also

- [Documentation Index](DOCUMENTATION_INDEX.md)
- [Getting Started](GETTING_STARTED.md)
