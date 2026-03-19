# DarkBASIC Pro → CyberBASIC2 compatibility

This table maps DarkBASIC Pro commands to CyberBASIC2 equivalents. Use it when porting DBP code. **Bold** = same or very similar name in CyberBASIC2.

## Core / language

| DBP | CyberBASIC2 |
|-----|------------|
| PRINT | **Print** (language built-in) |
| IF / ELSE / ENDIF | **IF / THEN / ELSE / END IF** |
| FOR / NEXT | **FOR / TO / STEP / NEXT** |
| REPEAT / UNTIL | **LOOP / UNTIL** or **REPEAT / UNTIL** (if supported) |
| DO / LOOP | **LOOP** (check language docs) |
| FUNCTION / ENDFUNCTION | **Function / End Function** |
| GOSUB / RETURN | Use **Sub** and **Call** |
| DIM | **VAR** or variable declaration |
| DATA / READ / RESTORE | Use arrays or **ReadFile** + parse |
| REM / REMSTART / REMEND | **//** or **/* */** |

## Math / string (stdlib)

| DBP | CyberBASIC2 |
|-----|------------|
| RND | **Rnd()** (0..1 float) or **Rnd(n)** (1..n int) |
| RANDOMIZE | **SetRandomSeed(seed)** (raylib) |
| INT(x) | **Int(x)** |
| ABS | **Clamp** or implement with **If**; or add Abs if needed |
| SIN, COS, TAN, ATAN, etc. | Use raylib math or add to std if needed |
| LEFT$(s, n) | **Left(s, n)** |
| RIGHT$(s, n) | **Right(s, n)** |
| MID$(s, p, n) | **Mid(s, p, n)** (p 1-based) |
| LEN(s) | **Len(s)** or **TextLength(s)** |
| INSTR(s, sub) / INSTR(s, sub, start) | **Instr(s, sub)** or **Instr(s, sub, start)** (1-based; 0 if not found) |
| CHR$(code) | **Chr(code)** |
| ASC(s) | **Asc(s)** |
| STR$(x) | **Str(x)** |
| VAL(s) | **Val(s)** or **TextToInteger(s)** / **TextToFloat(s)** (raylib) |

## File / system

| DBP | CyberBASIC2 |
|-----|------------|
| FILE EXIST | **FileExists(path)** (raylib core) |
| OPEN TO READ / READ / CLOSE FILE | **ReadFile(path)** |
| OPEN TO WRITE / WRITE / CLOSE FILE | **WriteFile(path, contents)** |
| DELETE FILE | **DeleteFile(path)** |
| COPY FILE | **CopyFile(src, dst)** |
| DIR / directory listing | **Dir(path)** or **ListDir(path)** → count; **GetDirItem(i)** for each name |
| EXECUTE FILE | **ExecuteFile(path)** |
| GET DIR$ | Use **ListDir** + **GetDirItem** or current-dir API if added |

## Mouse / keyboard

| DBP | CyberBASIC2 |
|-----|------------|
| MOUSEX, MOUSEY | **GetMouseX()**, **GetMouseY()** |
| MOUSECLICK | **IsMouseButtonPressed(MouseButtonLeft)** etc. |
| HIDE MOUSE / SHOW MOUSE | **HideCursor()** / **ShowCursor()** |
| INKEY$ | **GetCharPressed()** (raylib) or key state |
| ESCAPEKEY | **IsKeyPressed(KEY_ESCAPE)** |
| UPKEY, DOWNKEY, LEFTKEY, RIGHTKEY | **IsKeyDown(KEY_UP)** etc. |
| CONTROLKEY, SHIFTKEY | **IsKeyDown(KEY_LEFT_CONTROL)** etc. |

## 2D (sprites / bitmaps)

| DBP | CyberBASIC2 |
|-----|------------|
| LOAD BITMAP / CREATE BITMAP | **LoadTexture(path)** or **LoadImage** (raylib) |
| SPRITE / DRAW SPRITE | **LoadTexture** + **DrawTexture(x, y, tint)** |
| MOVE SPRITE, SIZE SPRITE | Draw at (x, y) with **DrawTextureEx** for scale/rotation |
| SPRITE HIT | **CheckCollisionPointRec** or manual bounds |
| BACKDROP | **ClearBackground(r,g,b,a)** or full-screen texture |

## 3D (objects / camera / lights)

| DBP | CyberBASIC2 |
|-----|------------|
| LOAD OBJECT | **LoadModel(path)** |
| POSITION OBJECT | **DrawModelEx(modelId, x, y, z, rotX, rotY, rotZ, scale, tint)** (set position via draw) or use ECS/Bullet for body position |
| ROTATE OBJECT | Pass rotation to **DrawModelEx** or physics body |
| MAKE OBJECT BOX/SPHERE/CONE etc. | **GenMeshCube** + **LoadModelFromMesh** or **DrawCube** / **DrawSphere** etc. |
| OBJECT COLLISION | **CheckCollisionBoxes**, **GetRayCollision*** (raylib); or Box2D/Bullet |
| POSITION CAMERA | **SetCamera3D** (position, target, up) |
| POINT CAMERA | Set target in **SetCamera3D** |
| MAKE LIGHT / POSITION LIGHT | raylib has limited built-in lighting; use shaders or **SetShaderValue** for custom |

## Sound / music

| DBP | CyberBASIC2 |
|-----|------------|
| LOAD SOUND / PLAY SOUND | **LoadSound(path)**, **PlaySound(soundId)** |
| LOAD MUSIC / PLAY MUSIC | **LoadMusicStream(path)**, **PlayMusicStream(musicId)**, **UpdateMusicStream(musicId)** |
| SET SOUND VOLUME | **SetSoundVolume(soundId, 0..1)** |

## Network

| DBP | CyberBASIC2 |
|-----|------------|
| CREATE NET GAME / JOIN NET GAME | **Host(port)** → serverId; **Connect(host, port)** → connectionId |
| SEND NET MESSAGE | **Send(connectionId, text)** |
| GET NET MESSAGE | **Receive(connectionId)** (line-based, non-blocking) |
| FREE NET GAME / DISCONNECT | **CloseServer(serverId)** / **Disconnect(connectionId)** |

## GUI

| DBP | CyberBASIC2 |
|-----|------------|
| (dialogs / buttons) | **GuiButton**, **GuiLabel**, **GuiSlider**, etc. (raygui) or **BeginUI** / **Button**, **Label**, etc. (pure-Go UI) |

## Summary

- **Already in CyberBASIC2:** Print, file (ReadFile, WriteFile, DeleteFile, FileExists), mouse/keyboard (GetMouseX/Y, IsKeyPressed, etc.), 3D (LoadModel, SetCamera3D, DrawModel, meshes), 2D (DrawTexture, LoadTexture, shapes), audio (LoadSound, PlayMusicStream, etc.), network (Connect, Send, Receive, Host, Accept), collision (CheckCollision*, GetRayCollision*, Box2D, Bullet-shaped 3D fallback physics), GUI (Gui* raygui, BeginUI/EndUI).
- **Added for DBP parity:** Left, Right, Mid, Len, Instr, Chr, Asc, Str, Val, Rnd, Int (std); CopyFile, Dir (alias for ListDir), ListDir, GetDirItem, ExecuteFile (std). See [DBP_GAP.md](DBP_GAP.md) for the gap list.
