# BASIC Programming Guide

Step-by-step tutorial for CyberBasic: variables, types, I/O, and handling errors. For game-specific topics see [Game Development Guide](GAME_DEVELOPMENT_GUIDE.md) and the [Documentation Index](DOCUMENTATION_INDEX.md).

## 1. Variables and assignment

Use **VAR** to declare and assign in one line, or **DIM** to declare (optionally with a type), then **LET** to assign:

```basic
VAR x = 10
VAR name = "CyberBasic"
LET x = 20

DIM y AS Float
LET y = 3.14

DIM a[10]        // array of 10 elements
DIM grid[5, 5]   // 2D array
```

Names are **case-insensitive**. Prefer **VAR** for local or one-off variables and **DIM** when you want an explicit type hint (e.g. `AS Float`, `AS Integer`, `AS String`).

## 2. Constants and types

**CONST** gives a name to a fixed value:

```basic
CONST Pi = 3.14159
CONST ScreenW = 800
CONST MaxLives = 3
```

**TYPE ... END TYPE** defines a structured type; use dot notation to access fields:

```basic
TYPE Player
    x AS Float
    y AS Float
    health AS Integer
END TYPE

VAR p = Player()
p.x = 100
p.y = 200
p.health = 3
```

**ENUM** defines named constants:

```basic
ENUM State : Idle = 0, Walk = 1, Jump = 2
VAR s = State.Walk
```

## 3. Control flow

Use **IF ... THEN ... ELSE ... ENDIF**, **WHILE ... WEND**, **FOR ... NEXT**, **REPEAT ... UNTIL**, and **SELECT CASE**:

```basic
IF score > 100 THEN
    PRINT "High score!"
ELSE
    PRINT "Keep trying"
ENDIF

FOR i = 1 TO 10
    PRINT i
NEXT i

SELECT CASE state
    CASE 0 : PRINT "Idle"
    CASE 1 : PRINT "Walking"
    CASE ELSE : PRINT "Other"
END SELECT
```

## 4. Functions and subs

**FUNCTION** returns a value; **SUB** does not:

```basic
FUNCTION Add(a, b)
    RETURN a + b
END FUNCTION

SUB Greet(name)
    PRINT "Hello, " + name
END SUB

VAR sum = Add(2, 3)   // 5
Greet("World")
```

Use **MODULE ... END MODULE** to group functions and subs under a name (e.g. `Math3D.Dot(...)`). See [Quick Reference](QUICK_REFERENCE.md) and [Language Spec](../LANGUAGE_SPEC.md).

## 5. Input and output

- **PRINT** writes to the console (one or more values, separated by commas or semicolons):

```basic
PRINT "Hello"
PRINT 42, " items"
PRINT "Score: "; score
```

- **File I/O:** Use **ReadFile(path)** to read a whole file as a string; it returns **Nil** if the file does not exist or cannot be read. Use **WriteFile(path, contents)** to write a string to a file (returns true/false). **FileExists(path)** checks existence.

```basic
VAR text = ReadFile("config.txt")
IF text <> Nil THEN
    PRINT text
ENDIF
WriteFile("out.txt", "Hello from CyberBasic")
```

- **JSON:** **LoadJSON(path)** or **LoadJSONFromString(str)** to load; **GetJSONKey(handle, key)** to read values; **SaveJSON(path, handle)** to save. See [API_REFERENCE.md](../API_REFERENCE.md).

## 6. Errors and null handling

Many functions signal “no value” or “error” by returning **Nil** (or **Null**). Always check before using the result:

```basic
VAR data = ReadFile("missing.txt")
IF IsNull(data) THEN
    PRINT "File not found"
ELSE
    PRINT data
ENDIF
```

- **Nil** / **Null** – literal for “no value”; use **IsNull(value)** to test.
- Compare with `= Nil` or `<> Nil` when you need an explicit check.

When you call a Sub or Function that might fail (e.g. file or network), check the return value (or use IsNull) and handle the failure in your code; CyberBasic does not have built-in exceptions.

## 7. Next steps

| Goal | Where to go |
|------|-------------|
| **Syntax at a glance** | [Quick Reference](QUICK_REFERENCE.md) |
| **Full language rules** | [Language Spec](../LANGUAGE_SPEC.md) |
| **First 2D/3D game** | [Game Development Guide](GAME_DEVELOPMENT_GUIDE.md) → [2D](2D_GRAPHICS_GUIDE.md) / [3D](3D_GRAPHICS_GUIDE.md) |
| **All docs** | [Documentation Index](DOCUMENTATION_INDEX.md) |
