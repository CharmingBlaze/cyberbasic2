# CyberBasic Language Spec

Single language reference for CyberBasic. For a one-page cheat sheet see [docs/QUICK_REFERENCE.md](docs/QUICK_REFERENCE.md). For installation and first run see [docs/GETTING_STARTED.md](docs/GETTING_STARTED.md).

Current implementation: Go compiler + VM + raylib/Box2D/Bullet. All names are **case-insensitive**.

---

## 1. Core language

### 1.1 Variables and types

```basic
DIM x
DIM x AS Float
DIM name$ AS String
DIM a[10]

VAR y = 10
LET z = 20
VAR name = "CyberBasic"

CONST Pi = 3.14159
CONST MaxLives = 3

DIM p AS Player
DIM enemies[100] AS Enemy
```

- **VAR** and **LET** declare and assign (VAR is modern style).
- **DIM** declares (optionally with **AS** type); arrays: `DIM a[n]` or `DIM a[m,n]`.
- **CONST** name = expression (compile-time constant).
- **Null literal:** Use **Nil**, **Null**, or **None** (case-insensitive) to represent a missing value. Assign and compare: `VAR x = Nil`, `IF result = Null THEN ...`, `IF result <> None THEN ...`. **IsNull(value)** returns true when the value is null (e.g. `IF IsNull(ReadFile("x.txt")) THEN ...`).

**User types:**

```basic
TYPE Player
    x AS Float
    y AS Float
    vx AS Float
    vy AS Float
    health AS Integer
END TYPE
```

**Entities (single instance):** `ENTITY` defines one named instance stored as a dictionary in globals. Use **END ENTITY** or **ENDENTITY**. Properties are `name = expression` (initial values). Read and write with `entityName.property`.

```basic
ENTITY Player
  x = 100
  y = 200
  health = 3
END ENTITY

Player.x = 50
VAR a = Player.x
```

- There is a single instance per entity name (e.g. one `Player` in globals).
- `entityName.property` reads and writes that instance’s property (via the same dict/GetJSONKey/SetDictKey mechanism as other dictionaries).

**Physics binding:** If an entity has properties `body` (string body id) and `world` (string world id), then reading `entity.x`, `entity.y` (2D) or `entity.x`, `entity.y`, `entity.z` (3D) and `entity.angle` (2D) or `entity.yaw`/`entity.pitch`/`entity.roll` (3D) returns values from the physics engine. Assign `entity.body = bodyId` and `entity.world = worldId` after creating a physics body to enable this.

**Rendering:** **DrawEntity(entityName)** draws the entity’s `sprite` (or `texture`) at `entity.x`, `entity.y`; if the entity has `scale` and `angle`, uses scaled/rotated drawing. Use **DrawEntityEx(entityName, scale, angle)** for an extended form if added. The entity’s `sprite` should be a texture id from **LoadSprite** or **LoadTexture**.

**Dot notation:** `p.x = 100`, `p.y = 200`, `p.health = 100`.

**Enums:** Single-line (`ENUM Name : a, b = 2, c`) or multi-line with **END ENUM** / **ENDENUM**. Name is optional (unnamed enum). Use **Enum.getValue(enumName, valueName)**, **Enum.getName(enumName, value)**, **Enum.hasValue(enumName, valueName)** at runtime.

```basic
ENUM Color : Red, Green, Blue
ENUM State : Idle = 0, Walk = 1, Jump = 2
ENUM Direction
    North, South, East, West
END ENUM
```

**Dictionary literals:** `{ "key": value }` (JSON-style) or `{ key = value }` (BASIC-style). Keys can be string, number, or identifier. Use **GetJSONKey(dict, key)** to read; **CreateDict()** and **SetDictKey** for building; **Dictionary.has/keys/values/size/remove/clear/merge/get** for operations.

### 1.2 Control flow

```basic
IF condition THEN
    ...
ELSE
    ...
ENDIF

WHILE condition
    ...
WEND

FOR i = 1 TO 10
    ...
NEXT i

FOR i = 10 TO 1 STEP -1
    ...
NEXT i

REPEAT
    ...
UNTIL condition

SELECT CASE value
    CASE 1 : ...
    CASE 2 : ...
    CASE ELSE : ...
END SELECT

EXIT FOR
EXIT WHILE

WHILE NOT WindowShouldClose()
    ClearBackground(0, 0, 0, 255)
    // ... draw calls ...
WEND
```

### 1.3 Compound assignment

Use `+=`, `-=`, `*=`, `/=` for in-place updates. LET is optional when the statement starts with a variable name.

```basic
camAngle -= GetMouseDeltaX() * 0.002
camPitch += GetMouseDeltaY() * 0.002
camDist = Clamp(camDist, 3, 25)
```

### 1.4 Functions and subs

```basic
FUNCTION Add(a, b)
    RETURN a + b
END FUNCTION

SUB DrawPlayer()
    // no return
END SUB
```

### 1.5 Modules (namespaces)

```basic
MODULE Math3D
    FUNCTION Dot(a, b)
        RETURN a.x * b.x + a.y * b.y + a.z * b.z
    END FUNCTION
END MODULE

VAR d = Math3D.Dot(v1, v2)
```

Module body contains only Function/Sub. Call as **ModuleName.FunctionName(...)**.

### 1.6 Events (optional)

```basic
ON KeyDown("ESCAPE")
    // handle
END ON

ON KeyPressed("SPACE")
    // fire once
END ON
```

Handlers run when **PollInputEvents()** is called (e.g. in the game loop).

### 1.7 Coroutines (optional)

```basic
StartCoroutine FadeOut()
Yield
WaitSeconds(1.0)
```

**StartCoroutine SubName()** starts a fiber; **Yield** switches fiber; **WaitSeconds(seconds)** blocks current fiber for N seconds.

### 1.8 Comments

```basic
// full-line comment
PRINT "hi"   // inline comment
```

Comments are **only** `//` (line) and `/* */` (block).

### 1.9 Includes and import (modularity)

At the top of a line (optionally after whitespace), use **#include "path"** or **IMPORT "path"** to insert the contents of another `.bas` file. The path is relative to the file containing the line. One directive per line. Use for shared code, libraries, or packaging. Both directives behave the same (file is inlined; cycles are avoided).

```basic
#include "lib/utils.bas"
IMPORT "math.bas"
IMPORT "graphics.bas"
```

---

## 2. Math, vectors, timers, random

- **Math:** `Sin`, `Cos`, `Tan`, `Sqrt`, `Abs`, `Random`, `GetRandomValue`, `SetRandomSeed`. **Std:** `Radians(degrees)`, `Degrees(radians)`, `AngleWrap(angle)` (wrap to [-π, π]). **Raylib:** `Clamp(value, min, max)`, `Lerp(start, end, amount)`, `Wrap(value, min, max)` – see [API_REFERENCE.md](API_REFERENCE.md).
- **Vectors:** Vector2/Vector3 types and helpers; raylib math (Vector2Add, Vector3Scale, etc.) – see [API_REFERENCE.md](API_REFERENCE.md)
- **Time:** `GetFrameTime`, **DeltaTime** (same as GetFrameTime; preferred for frame delta), `SetTargetFPS`, `GetFPS`, `GetTime`

---

## 3. Window, system, input

- **Window:** `InitWindow(width, height, title)`, `SetTargetFPS(fps)`, `WindowShouldClose()`, `CloseWindow()`, `SetWindowSize`, `SetWindowTitle`, etc.
- **Input:** `IsKeyDown(KEY_W)`, `IsKeyPressed(KEY_SPACE)`, `GetAxisX()`, `GetAxisY()`, `GetMouseX()`, `GetMouseY()`, `GetMouseDelta()`, `IsMouseButtonDown`, etc.

See [API_REFERENCE.md](API_REFERENCE.md) for the full list.

---

## 4. 2D graphics

- **Frame:** The compiler does not inject any frame or mode calls; your code compiles as written (DBPro-style). Use `ClearBackground(r, g, b, a)` and your draw calls in the loop.
- **Primitives:** `DrawRectangle`, `DrawCircle`, `DrawLine`, `DrawTriangle`, `DrawPixel`, etc.
- **Textures:** `LoadTexture(path)`, `DrawTexture(id, x, y)`, `DrawTextureEx`, `DrawTextureRec`, `UnloadTexture`
- **Text:** `DrawText(text, x, y, fontSize, r, g, b, a)`, `MeasureText`, `LoadFont`, `DrawTextEx`

See [docs/2D_GRAPHICS_GUIDE.md](docs/2D_GRAPHICS_GUIDE.md) and [API_REFERENCE.md](API_REFERENCE.md).

---

## 5. 3D graphics

- **Camera:** `SetCamera3D(posX, posY, posZ, targetX, targetY, targetZ, upX, upY, upZ)`; for orbit use **GAME.CameraOrbit(...)** each frame
- **3D drawing:** See raylib bindings for 3D primitives and models.
- **Primitives:** `DrawCube`, `DrawSphere`, `DrawPlane`, `DrawLine3D`, `DrawGrid`, etc.
- **Models:** `LoadModel(path)`, `DrawModel(id, x, y, z, scale)`, `DrawModelEx`, `UnloadModel`; **GenMeshCube**, **GenMeshSphere**, **LoadModelFromMesh**

See [docs/3D_GRAPHICS_GUIDE.md](docs/3D_GRAPHICS_GUIDE.md) and [API_REFERENCE.md](API_REFERENCE.md).

---

## 6. Shaders

`LoadShader`, `BeginShaderMode` / `EndShaderMode`, `UnloadShader` – see raylib bindings in [API_REFERENCE.md](API_REFERENCE.md).

---

## 7. Audio

`InitAudioDevice`, `LoadSound`, `PlaySound`, `StopSound`, `SetSoundVolume`, `LoadMusicStream`, `PlayMusicStream`, `UpdateMusicStream`, `StopMusicStream`, `SetMusicVolume`, `UnloadSound`, `UnloadMusicStream`, etc. See [API_REFERENCE.md](API_REFERENCE.md).

---

## 8. File, JSON, HTTP (std)

- **File:** `ReadFile(path)` → string or nil; `WriteFile(path, contents)` → boolean; `DeleteFile(path)` → boolean; **FileExists(path)** (raylib core)
- **JSON:** `LoadJSON(path)`, `LoadJSONFromString(str)`, `GetJSONKey(handle, key)`, `SaveJSON(path, handle)`
- **HTTP:** `HttpGet(url)`, `HttpPost(url, body)`, `DownloadFile(url, path)`

---

## 9. Physics 2D (Box2D)

Use **BOX2D.*** namespace: `BOX2D.CreateWorld`, `BOX2D.Step`, `BOX2D.CreateBody`, `BOX2D.CreateBox`, `BOX2D.CreateCircle`, `BOX2D.GetPositionX`, `BOX2D.GetPositionY`, `BOX2D.SetLinearVelocity`, `BOX2D.ApplyForce`, etc. Legacy names without prefix (CreateWorld2D, Step2D, CreateBox2D, …) also available. See [docs/GAME_DEVELOPMENT_GUIDE.md](docs/GAME_DEVELOPMENT_GUIDE.md) and [API_REFERENCE.md](API_REFERENCE.md).

---

## 10. Physics 3D (Bullet)

Use **BULLET.*** namespace: `BULLET.CreateWorld`, `BULLET.Step`, `BULLET.CreateSphere`, `BULLET.CreateBox`, `BULLET.GetPositionX`, `BULLET.GetPositionY`, `BULLET.GetPositionZ`, `BULLET.SetVelocity`, `BULLET.ApplyForce`, etc. Legacy names (CreateWorld3D, Step3D, …) also available. See [docs/GAME_DEVELOPMENT_GUIDE.md](docs/GAME_DEVELOPMENT_GUIDE.md) and [API_REFERENCE.md](API_REFERENCE.md).

---

## 11. ECS (Entity-Component System)

Use **ECS.*** from the library binding: `ECS.CreateWorld`, `ECS.CreateEntity`, `ECS.AddComponent`, `ECS.GetComponent`, `ECS.Query`, `ECS.DestroyEntity`, `ECS.DestroyWorld`, etc. See [docs/ECS_GUIDE.md](docs/ECS_GUIDE.md) and [API_REFERENCE.md](API_REFERENCE.md).

---

## 12. UI (minimal)

`BeginUI()`, `Label(text)`, `Button(text)` → boolean, `EndUI()`. See [API_REFERENCE.md](API_REFERENCE.md).

---

## 13. Typical game loop (2D)

Use **WHILE NOT WindowShouldClose() ... WEND** (or **REPEAT ... UNTIL WindowShouldClose()**) for the main loop. No auto-wrap; your code compiles as written. **DeltaTime()** returns frame delta (same as GetFrameTime).

```basic
InitWindow(800, 450, "CyberBasic Game")
SetTargetFPS(60)

WHILE NOT WindowShouldClose()
    VAR dt = DeltaTime()
    IF dt > 0.05 THEN LET dt = 0.016
    // Input and logic...
    // BOX2D.Step("w", dt, 8, 3)
    ClearBackground(20, 20, 30, 255)
    DrawRectangle(playerX, playerY, 32, 32, 255, 255, 255, 255)
WEND

CloseWindow()
```

---

## Implementation notes

- **Case-insensitive:** Keywords, identifiers, and built-in/foreign names. `MyVar` and `myvar` are the same.
- **Dynamic typing:** Variables hold values; `AS Type` is an optional hint. No static type checking.
- **Namespaces:** Call raylib as `InitWindow(...)` or `RL.InitWindow(...)`; physics as `BOX2D.*`, `BULLET.*`; game helpers as `GAME.*`; ECS as `ECS.*`.
- **Bytecode constant limit:** A program may use at most 256 constants (e.g. string/number literals, global names); compiled bytecode uses a single byte for constant indices.

For the full API see [API_REFERENCE.md](API_REFERENCE.md). For doc index see [docs/DOCUMENTATION_INDEX.md](docs/DOCUMENTATION_INDEX.md).
