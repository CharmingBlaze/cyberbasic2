#!/usr/bin/env python3
"""
Strip the RL. namespace from raylib bindings so BASIC can call InitWindow()
instead of RL.InitWindow(). Updates all RegisterForeign("RL.XXX", ...) to
RegisterForeign("XXX", ...) in compiler/bindings/raylib/*.go.

Run from repo root: python scripts/strip_raylib_namespace.py
"""
import re
from pathlib import Path

REPO_ROOT = Path(__file__).resolve().parent.parent
RAYLIB_DIR = REPO_ROOT / "compiler" / "bindings" / "raylib"


def strip_namespace(content: str) -> tuple[str, int]:
    """Replace RegisterForeign(\"RL.XXX\" with RegisterForeign(\"XXX\". Returns (new_content, count)."""
    pattern = re.compile(r'RegisterForeign\("RL\.([^"]+)"')
    count = 0
    def repl(m: re.Match) -> str:
        nonlocal count
        count += 1
        return f'RegisterForeign("{m.group(1)}"'
    new_content = pattern.sub(repl, content)
    return new_content, count


def main() -> None:
    if not RAYLIB_DIR.exists():
        print(f"Raylib dir not found: {RAYLIB_DIR}")
        return
    total = 0
    for path in sorted(RAYLIB_DIR.glob("*.go")):
        text = path.read_text(encoding="utf-8")
        new_text, n = strip_namespace(text)
        if n:
            path.write_text(new_text, encoding="utf-8")
            print(f"{path.relative_to(REPO_ROOT)}: stripped RL. from {n} RegisterForeign calls")
            total += n
    if total:
        print(f"Done. Total: {total} bindings updated.")
    else:
        print("No RL. namespace found (already stripped or no matches).")


if __name__ == "__main__":
    main()
