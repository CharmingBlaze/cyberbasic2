// Smoke: dot-style window + draw namespaces (same foreigns as flat raylib).
// Run: cyberbasic examples/smoke_dot_api.bas
// Lint: cyberbasic --lint examples/smoke_dot_api.bas
WINDOW.INITWINDOW(400, 240, "dot smoke")
WINDOW.SETTARGETFPS(60)
WHILE NOT WINDOW.WINDOWSHOULDCLOSE()
  DRAW.BEGIN()
  DRAW.CLEAR(30, 35, 45, 255)
  DRAW.TEXT("smoke_dot_api ok", 12, 12, 18, 220, 240, 255, 255)
  DRAW.END()
WEND
WINDOW.CLOSE()
