# DarkBASIC Pro parity – gap list

Only DBP commands that have **no** current equivalent in CyberBasic are listed here. These are the ones we implement; all others are skipped (see [DBP_COMPAT.md](DBP_COMPAT.md) for DBP → CyberBasic mapping).

## String / math (stdlib)

| DBP command | CyberBasic equivalent (before) | Action |
|-------------|--------------------------------|--------|
| LEFT$(s, n) | — | **Add** Left(s, n) in std |
| RIGHT$(s, n) | — | **Add** Right(s, n) in std |
| MID$(s, p, n) | TextSubtext(s, p-1, n) exists but 0-based | **Add** Mid(s, start1Based, count) in std (1-based for DBP familiarity) |
| LEN(s) | TextLength(s) | **Add** alias Len(s) in std |
| CHR$(code) | — | **Add** Chr(code) in std |
| ASC(s) | — | **Add** Asc(s) in std |
| STR$(x) | TextFormat("%v", x) possible | **Add** Str(x) in std |
| VAL(s) | TextToInteger/TextToFloat | **Add** alias Val(s) in std (float) |
| RND / RND(n) | GetRandomValue(0,999)/GetRandomValue(1,n) in raylib | **Add** Rnd() and Rnd(n) in std (no raylib dep) |
| INT(x) | — | **Add** Int(x) in std (truncate) |

## File / system (stdlib)

| DBP command | CyberBasic equivalent (before) | Action |
|-------------|--------------------------------|--------|
| COPY FILE | — | **Add** CopyFile(src, dst) in std |
| DIR / directory list | — | **Add** Dir(path) or ListDir(path) in std (return list of names) |
| EXECUTE FILE | — | **Add** ExecuteFile(path) in std |

## Already covered (skip)

- **PRINT** → language built-in Print.
- **FILE EXIST** → FileExists (raylib_core).
- **MOUSEX, MOUSEY** → GetMouseX, GetMouseY.
- **LOAD OBJECT, 3D, camera, etc.** → LoadModel, SetCamera3D, DrawModel, etc.
- **Sound/Music** → LoadSound, PlayMusicStream, etc.
- **Network** → Connect, Send, Receive, Host, Accept (net.go).
- **Control flow** → IF, FOR, LOOP, etc. in language.
- **Collision** → CheckCollision*, GetRayCollision*, Box2D/Bullet.
- **GUI** → Gui* (raygui), BeginUI/EndUI.

Implementation: add the **Add** items above in [compiler/bindings/std/std.go](../compiler/bindings/std/std.go) and document in [DBP_COMPAT.md](DBP_COMPAT.md).
