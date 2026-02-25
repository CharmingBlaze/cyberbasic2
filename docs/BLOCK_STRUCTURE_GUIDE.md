# Block Structure Guide

CyberBASIC uses explicit END keywords for block structure. This guide documents what's supported in the parser and interpreter.

## How It Works

- **Parser level:** The parser validates END keywords during code parsing. If an END keyword is missing or incorrect, a parse error occurs.
- **Interpreter level:** The interpreter works with the parsed AST (Abstract Syntax Tree) and doesn't check END keywords at runtime. By the time code reaches the interpreter, all END keywords have already been validated by the parser.

## Supported END Keywords

CyberBASIC supports both single-word and two-word END forms where noted.

### Single-Word END Keywords

- **ENDFUNCTION** — ends a `FUNCTION` block
- **ENDSUB** — ends a `SUB` block
- **ENDIF** — ends an `IF` block
- **WEND** — ends a `WHILE` block (no `END WHILE` variant)
- **ENDTYPE** — ends a `TYPE` block
- **ENDMODULE** — ends a `MODULE` block
- **ENDSELECT** — ends a `SELECT CASE` block
- **ENDON** — ends an `ON` event block

### Two-Word END Keywords (supported as alternatives)

- **END IF** — same as `ENDIF`
- **END FUNCTION** — same as `ENDFUNCTION`
- **END SUB** — same as `ENDSUB`
- **END MODULE** — same as `ENDMODULE`
- **END TYPE** — same as `ENDTYPE`
- **END SELECT** — same as `ENDSELECT`

### Other Block Endings

- **NEXT** — ends a `FOR` loop (not an END keyword)
- **UNTIL** — ends a `REPEAT` loop

## Control Flow: IF, ELSEIF, ELSE, END IF

- **IF** … **THEN** … **ENDIF** or **END IF**
- **ELSEIF** — optional, multiple branches: `IF … THEN … ELSEIF … THEN … ELSE … END IF`
- **ELSE** — optional final branch

## Examples

### Function (single-word)

```basic
FUNCTION Add(a, b)
    RETURN a + b
ENDFUNCTION
```

### Function (two-word)

```basic
FUNCTION Add(a, b)
    RETURN a + b
END FUNCTION
```

### Sub (single-word)

```basic
SUB PrintHeader(title)
    PRINT "=========="
    PRINT title
    PRINT "=========="
ENDSUB
```

### Sub (two-word)

```basic
SUB PrintHeader(title)
    PRINT "=========="
    PRINT title
END SUB
```

### IF with END IF (two words)

```basic
IF x > 0 THEN
    PRINT "Positive"
END IF
```

### IF with ELSEIF and END IF

```basic
FUNCTION calculate(x, y, op)
    IF op = "add" THEN
        RETURN x + y
    ELSEIF op = "multiply" THEN
        RETURN x * y
    ELSE
        RETURN 0
    END IF
ENDFUNCTION
```

### Mixed style (END IF + ENDFUNCTION)

```basic
FUNCTION calculate(x, y, op)
    IF op = "add" THEN
        RETURN x + y
    ELSEIF op = "multiply" THEN
        RETURN x * y
    ELSE
        RETURN 0
    END IF
ENDFUNCTION
```

### Nested blocks

```basic
FUNCTION processData(data)
    IF data <> "" THEN
        VAR result = PARSEJSON(data)
        IF result <> NIL THEN
            RETURN result
        END IF
    END IF
    RETURN NIL
END FUNCTION
```

### Module (END MODULE or ENDMODULE)

```basic
MODULE Helpers
    FUNCTION double(n)
        RETURN n * 2
    END FUNCTION
    SUB greet(name)
        PRINT "Hello, "; name
    END SUB
END MODULE
```

### Indentation

CyberBASIC requires explicit END keywords; indentation alone is not sufficient. Consistent indentation (2 or 4 spaces) is recommended for readability:

```basic
FUNCTION wellFormatted(a, b)
    IF a > b THEN
        VAR result = a - b
        IF result > 0 THEN
            RETURN result
        END IF
    END IF
    RETURN 0
ENDFUNCTION
```

## Best Practices

1. **Be consistent** — Choose one style (single-word or two-word) and use it consistently in a project.
2. **Indent properly** — Use consistent indentation to show block structure.
3. **Match keywords** — Use the correct END form for each block:
   - `FUNCTION` … `ENDFUNCTION` or `END FUNCTION`
   - `IF` … `ENDIF` or `END IF`
   - `SUB` … `ENDSUB` or `END SUB`
   - `MODULE` … `ENDMODULE` or `END MODULE`
4. **Readability first** — Prefer the style that works best for your team.

## Case Insensitivity

All keywords are case-insensitive:

```basic
function Add(a, b)
    return a + b
endfunction

FUNCTION Subtract(a, b)
    RETURN a - b
ENDFUNCTION

Function Multiply(a, b)
    Return a * b
EndFunction
```

All of the above are valid and equivalent.
