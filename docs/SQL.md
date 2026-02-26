# SQL (SQLite) in CyberBasic

CyberBasic includes SQLite bindings so you can open a database file, run SQL, and read results. Uses Go's `database/sql` with the pure-Go **modernc.org/sqlite** driver (no CGO). One database file per path; one "current" result set from the last **Query** or **QueryParams**.

## Quick start

```basic
VAR db = OpenDatabase("game.db")
IF IsNull(db) THEN
  PRINT "Open failed: " + LastError()
  END
END IF

Exec(db, "CREATE TABLE IF NOT EXISTS players (id INTEGER PRIMARY KEY, name TEXT, score INT)")
Exec(db, "INSERT INTO players (name, score) VALUES ('Alice', 100)")

Query(db, "SELECT name, score FROM players")
VAR rows = GetRowCount()
VAR cols = GetColumnCount()
FOR r = 0 TO rows - 1
  FOR c = 0 TO cols - 1
    VAR v = GetCell(r, c)
    IF IsNull(v) THEN PRINT "NULL" ELSE PRINT v
  NEXT c
  PRINT ""
NEXT r

CloseDatabase(db)
```

## Commands

| Command | Description |
|--------|-------------|
| **OpenDatabase**(path) | Open SQLite DB at path (e.g. `"game.db"`). Returns dbId or null on failure. |
| **CloseDatabase**(dbId) | Close the database. |
| **Exec**(dbId, sql) | Run INSERT/UPDATE/DELETE/DDL. Returns rows affected, or -1 on error. |
| **Query**(dbId, sql) | Run SELECT; store result. Returns row count, or -1 on error. |
| **GetRowCount**() | Rows in last query result (0 if no query or error). |
| **GetColumnCount**() | Columns in last query result. |
| **GetColumnName**(colIndex) | Column name (0-based). Empty string if out of range. |
| **GetCell**(row, col) | Value at (row, col) as number or string; **null** for SQL NULL — use `IsNull(GetCell(r,c))`. |
| **LastError**() | Last error message (string). Check after Open/Exec/Query failures. |
| **Begin**(dbId) | Start transaction. Returns 1 on success, 0 on error. |
| **Commit**(dbId) | Commit transaction. Returns 1 on success, 0 on error. |
| **Rollback**(dbId) | Rollback transaction. Returns 1 on success, 0 on error. |
| **ExecParams**(dbId, sql, arg1, arg2, …) | Like Exec; use `?` in sql; args (number or string) replace placeholders. Returns rows affected or -1. |
| **QueryParams**(dbId, sql, arg1, arg2, …) | Like Query; use `?` in sql for parameters. Returns row count or -1. |

Row and column indices are **0-based**. Only one result set is active at a time; the next **Query** or **QueryParams** replaces it.

## Parameterized statements

Use **ExecParams** and **QueryParams** with `?` placeholders to avoid SQL injection and to pass values safely:

```basic
// Insert with parameters
ExecParams(db, "INSERT INTO players (name, score) VALUES (?, ?)", "Bob", 200)

// Select with parameter
QueryParams(db, "SELECT name, score FROM players WHERE score > ?", 50)
VAR n = GetRowCount()
FOR r = 0 TO n - 1
  PRINT GetCell(r, 0) + " " + STR(GetCell(r, 1))
NEXT r
```

## Transactions

Wrap multiple writes in a transaction for atomicity:

```basic
Begin(db)
VAR a = Exec(db, "INSERT INTO players (name, score) VALUES ('A', 10)")
VAR b = Exec(db, "INSERT INTO players (name, score) VALUES ('B', 20)")
IF a >= 0 AND b >= 0 THEN Commit(db) ELSE Rollback(db)
```

## Errors

After **OpenDatabase**, **Exec**, or **Query** (or their Param variants), check the return value. If it indicates failure (null for Open, -1 for Exec/Query), call **LastError**() to get the error message:

```basic
VAR db = OpenDatabase("game.db")
IF IsNull(db) THEN
  PRINT "Error: " + LastError()
  END
END IF

VAR affected = Exec(db, "UPDATE players SET score = 0 WHERE id = 99")
IF affected < 0 THEN PRINT "Exec failed: " + LastError()
```

## Common patterns

### Leaderboard (SELECT with ORDER BY and LIMIT)

```basic
Query(db, "SELECT name, score FROM players ORDER BY score DESC LIMIT 10")
VAR n = GetRowCount()
FOR r = 0 TO n - 1
  PRINT GetCell(r, 0) + ": " + STR(GetCell(r, 1))
NEXT r
```

### Save/load game state (single row or JSON in a column)

Store a single row (e.g. one "save" slot) or a JSON string in a TEXT column:

```basic
ExecParams(db, "INSERT OR REPLACE INTO save (id, data) VALUES (1, ?)", jsonString)
QueryParams(db, "SELECT data FROM save WHERE id = ?", 1)
IF GetRowCount() > 0 THEN VAR loaded = GetCell(0, 0)
```

Use **GetJSONKey** or your JSON helpers to read fields from `loaded` if you store a JSON object.

---

## See also

- [API Reference](../API_REFERENCE.md) (section 19) — full SQL command list
- [Command Reference](COMMAND_REFERENCE.md) — commands by feature
- [Getting Started](GETTING_STARTED.md) — setup and first run
