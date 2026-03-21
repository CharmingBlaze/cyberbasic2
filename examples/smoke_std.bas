// Smoke: std library + string length (no window).
// cyberbasic --lint examples/smoke_std.bas
LET n = Len("ok")
IF n <> 2 THEN
  PRINT "unexpected Len"
END IF
PRINT "smoke_std ok"
