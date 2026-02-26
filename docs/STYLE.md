# Documentation style guide

Conventions for CyberBasic documentation so all docs feel consistent.

## Heading levels

- **One H1** per document: the main title (e.g. "Getting Started with CyberBasic").
- **##** for major sections (e.g. "Prerequisites", "Build the compiler").
- **###** for subsections when needed.
- Do not skip levels (e.g. do not go from ## to ####).

## Table of Contents

- Add a **Table of Contents** for documents with **4 or more** major sections (##).
- Use anchor links: `[Section name](#section-name)` (lowercase, spaces as hyphens).
- Place the TOC after the intro paragraph, before the first ##.

## Code blocks

- Use **` ```basic `** for all CyberBasic source code.
- Use **` ```bash `** (or no tag) for shell commands when trivial.
- Keep snippets short and runnable where possible.

## Cross-links

- End major docs with a **"See also"** or **"Related"** section.
- Link to [Documentation Index](DOCUMENTATION_INDEX.md) and to the next-step doc (e.g. Quick Reference, Game Development Guide).
- Use relative links within the repo: `[Link text](FILENAME.md)` or `[Link](../PATH.md)` from docs/.

## Terminology

- **Game loop:** Use "game loop" when referring to `WHILE NOT WindowShouldClose() ... WEND` or `REPEAT ... UNTIL WindowShouldClose()`.
- **Delta time:** Use "delta time" for frame elapsed time; refer to **DeltaTime()** or **GetFrameTime()**.
- **Flat API:** Prefer flat names in examples: **InitWindow**, **ClearBackground**, **DrawCircle** (not `RL.InitWindow`). Use namespaced forms only when documenting that API (e.g. **BOX2D.***, **GAME.***).

## Include syntax

- Document the directive as **`#include "path"`**. Optionally mention **IMPORT "path"** as an alias if supported.
- Do not use "INCLUDE" without the hash (the compiler expects `#include` or `IMPORT`).

## See also

- [Documentation Index](DOCUMENTATION_INDEX.md)
- [Getting Started](GETTING_STARTED.md)
