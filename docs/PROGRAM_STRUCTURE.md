# CyberBasic Program Structure

This document summarizes program structure, comments, and the main language features.

---

## Comments

Use **`//`** for line comments. Everything from `//` to the end of the line is ignored.

```basic
// This is a comment
VAR x = 10   // inline comment
PRINT x
```

---

## Feature list (implemented)

- **Variables:** `VAR`, `DIM`, `LET`; arrays `VAR a[10]`, `DIM b[5,5]`
- **Constants:** `CONST name = value`
- **Types:** `TYPE ... END TYPE`, `EXTENDS`
- **Enums:** `ENUM Name ... END ENUM` (named/unnamed, custom values); `Enum.getValue`, `Enum.getName`, `Enum.hasValue`
- **Control flow:** `IF/THEN/ELSE/ELSEIF/ENDIF`, `FOR/NEXT`, `WHILE/WEND`, `REPEAT/UNTIL`, `SELECT CASE`
- **Loop control:** `EXIT FOR`, `EXIT WHILE`, `BREAK FOR`, `BREAK WHILE`, `CONTINUE FOR`, `CONTINUE WHILE`
- **Procedures:** `SUB`, `FUNCTION`; `END SUB` / `ENDSUB`, `END FUNCTION` / `ENDFUNCTION`
- **Modules:** `MODULE name` / `END MODULE` / `ENDMODULE`
- **Operators:** `+ - * / % \` (integer div), `^` (power), `= <> < <= > >=`, `AND`, `OR`, `XOR`, `NOT`
- **Compound assign:** `+=`, `-=`, `*=`, `/=`
- **String/std:** `Left`, `Right`, `Mid`, `Substr`, `Instr`, `Upper`, `Lower`, `Len`, `Chr`, `Asc`, `Str`, `Val`, `Rnd`, `Rnd(n)`, `Random(n)`, `Int`
- **Assert:** `ASSERT condition [, message]`
- **Null:** `Nil`, `Null`, `None`; `IsNull(value)`
- **JSON/dict:** `LoadJSON`, `ParseJSON`, `GetJSONKey`, dict literal `{"key": value}` or `{key = value}`, `CreateDict`, `SetDictKey`, `Dictionary.has/keys/values/size/remove/clear/merge/get`
- **File I/O:** `ReadFile`, `WriteFile`, `DeleteFile`, `CopyFile`, `ListDir`
- **Includes:** `INCLUDE "file.bas"`
- **Events/coroutines:** `ON ... GOSUB`, `StartCoroutine`, `Yield`, `WaitSeconds`
- **Graphics:** raylib (2D/3D), Box2D, Bullet; automatic frame/mode wrapping in game loops
- **Multi-window:** `SpawnWindow`, `ConnectToParent`, `NET.*`
- **ECS, GUI, multiplayer:** See [ECS_GUIDE.md](ECS_GUIDE.md), [GUI_GUIDE.md](GUI_GUIDE.md), [MULTIPLAYER.md](MULTIPLAYER.md)

---

## Block structure (quick reference)

| Block        | Start        | End              |
|-------------|--------------|------------------|
| IF          | IF ... THEN  | ENDIF or END IF  |
| FOR         | FOR x = a TO b [STEP s] | NEXT   |
| WHILE       | WHILE cond   | WEND             |
| REPEAT      | REPEAT       | UNTIL cond       |
| SELECT CASE | SELECT CASE expr | ENDSELECT    |
| FUNCTION    | FUNCTION name(params) | ENDFUNCTION or END FUNCTION |
| SUB         | SUB name(params) | ENDSUB or END SUB   |
| MODULE      | MODULE name  | ENDMODULE or END MODULE |
| TYPE        | TYPE name    | ENDTYPE          |
| ENUM        | ENUM [name]  | ENDENUM or END ENUM |

---

## Example skeleton

```basic
// My game
INCLUDE "constants.bas"

ENUM GameState
    Menu, Playing, Paused
END ENUM

VAR state = 0
VAR config = {"width": 1024, "height": 768}

FUNCTION main()
    INITWINDOW(config["width"], config["height"], "Game")
    SETTARGETFPS(60)
    WHILE NOT WindowShouldClose()
        // Update and draw
        IF state = 0 THEN
            // menu
        ENDIF
    WEND
    CLOSEWINDOW()
ENDFUNCTION

main()
```
