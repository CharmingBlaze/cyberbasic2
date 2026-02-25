// SQL demo: OpenDatabase, Exec, Query, GetCell, display in window
// Run: cyberbasic examples/sql_demo.bas

VAR db = OpenDatabase("sql_demo.db")
IF IsNull(db) THEN
  PRINT "Open failed: " + LastError()
  QUIT
ENDIF

Exec(db, "CREATE TABLE IF NOT EXISTS players (id INTEGER PRIMARY KEY, name TEXT, score INT)")
Exec(db, "DELETE FROM players")
Exec(db, "INSERT INTO players (name, score) VALUES ('Alice', 100)")
Exec(db, "INSERT INTO players (name, score) VALUES ('Bob', 200)")
Exec(db, "INSERT INTO players (name, score) VALUES ('Carol', 150)")

Query(db, "SELECT name, score FROM players ORDER BY score DESC")
VAR rows = GetRowCount()
VAR cols = GetColumnCount()

InitWindow(640, 400, "SQL Demo")
SetTargetFPS(60)

WHILE NOT WindowShouldClose()
  ClearBackground(40, 44, 52, 255)
  DrawText("SQLite demo - players table", 20, 20, 20, 255, 255, 255, 255)
  VAR y = 60
  VAR r = 0
  VAR naLabel = "n/a"
  FOR r = 0 TO rows - 1
    VAR rowText = ""
    VAR c = 0
    FOR c = 0 TO cols - 1
      VAR v = GetCell(r, c)
      IF IsNull(v) THEN
        rowText = rowText + naLabel
      ENDIF
      IF NOT IsNull(v) THEN
        rowText = rowText + STR(v)
      ENDIF
      IF c < cols - 1 THEN
        rowText = rowText + "  |  "
      ENDIF
    NEXT c
    DrawText(rowText, 40, y, 18, 200, 220, 255, 255)
    LET y = y + 28
  NEXT r
  DrawText("Close window to exit", 20, 360, 16, 180, 180, 180, 255)
WEND

CloseDatabase(db)
CloseWindow()
