# Libraries and includes

How to structure multi-file CyberBasic projects and reuse `.bas` code as libraries.

## Table of Contents

1. [Include syntax](#include-syntax)
2. [Example layout](#example-layout)
3. [Rules](#rules)
4. [See also](#see-also)

---

## Include syntax

At the top of a `.bas` file (optionally after whitespace), use:

```basic
#include "path/to/file.bas"
```

- **One directive per line.** Multiple files can be included on separate lines.
- **Path** is relative to the file that contains the `#include` line.
- The compiler inserts the contents of the included file at that position (before parsing). Use for shared Subs, Functions, constants, and types.

## Example layout

```
mygame/
├── main.bas          // entry: InitWindow, game loop (WHILE NOT WindowShouldClose), includes
├── lib/
│   ├── utils.bas     // shared helpers (e.g. Clamp, Lerp)
│   └── player.bas     // player state and movement
└── assets/
    └── ...
```

**main.bas:**

```basic
#include "lib/utils.bas"
#include "lib/player.bas"

InitWindow(800, 600, "My Game")
SetTargetFPS(60)
// ... use Sub and Function names from utils.bas and player.bas
WHILE NOT WindowShouldClose()
  UpdatePlayer()
  DrawPlayer()
WEND
CloseWindow()
```

**lib/utils.bas:**

```basic
// Shared helpers (no executable entry; called from main)
Function Clamp(x, minVal, maxVal)
  IF x < minVal THEN RETURN minVal
  IF x > maxVal THEN RETURN maxVal
  RETURN x
End Function
```

**lib/player.bas:**

```basic
DIM playerX AS Float
DIM playerY AS Float
Sub UpdatePlayer()
  // ... move playerX, playerY from input
End Sub
Sub DrawPlayer()
  DrawCircle(playerX, playerY, 20, 255, 255, 255, 255)
End Sub
```

## Rules

- Included files are **not** namespaced: Subs, Functions, and global variables share the same global scope as the file that includes them. Avoid duplicate names across included files.
- No `#include` cycles: if `a.bas` includes `b.bas`, `b.bas` must not include `a.bas` (or anything that eventually includes `a.bas`).
- Paths use forward slashes; the compiler resolves relative to the including file’s directory.

## See also

- **Language Spec:** [LANGUAGE_SPEC.md](../LANGUAGE_SPEC.md) – Includes and libraries (section 5)
- **Quick Reference:** [QUICK_REFERENCE.md](QUICK_REFERENCE.md) – One-line include reminder
- **Documentation Index:** [DOCUMENTATION_INDEX.md](DOCUMENTATION_INDEX.md)
