# DarkBASIC Pro parity – gap list

Only DBP commands that have **no** current equivalent in CyberBASIC2 are listed here. These are the ones we implement; all others are skipped (see [DBP_COMPAT.md](DBP_COMPAT.md) for DBP → CyberBASIC2 mapping).

## Already implemented (stdlib)

- **LEFT$(s, n)** → Left(s, n) in std
- **RIGHT$(s, n)** → Right(s, n) in std
- **MID$(s, p, n)** → Mid(s, start1Based, count) in std (1-based)
- **LEN(s)** → Len(s) in std
- **CHR$(code)** → Chr(code) in std
- **ASC(s)** → Asc(s) in std
- **STR$(x)** → Str(x) in std
- **VAL(s)** → Val(s) in std (float)
- **RND / RND(n)** → Rnd() and Rnd(n) in std
- **INT(x)** → Int(x) in std (truncate)
- **COPY FILE** → CopyFile(src, dst) in std
- **DIR** → ListDir(path) / Dir(path) in std
- **EXECUTE FILE** → ExecuteFile(path) in std

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

See [compiler/bindings/std/std.go](../compiler/bindings/std/std.go) and [DBP_COMPAT.md](DBP_COMPAT.md).
