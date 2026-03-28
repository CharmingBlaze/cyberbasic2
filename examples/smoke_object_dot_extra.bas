// Smoke: DBP object handle dot methods beyond position/draw (compile-check).
// cyberbasic --lint examples/smoke_object_dot_extra.bas
// Uses a non-existent model path; intended for lint / API surface check only.
VAR obj = OBJECT.LOAD("missing_model_for_lint.glb")
obj.SETCOLOR(200, 100, 50)
obj.SETWIREFRAME(1)
obj.SETALPHA(240)
obj.SETCOLLISION(1)
obj.FIX()
obj.SHOW()
obj.MOVE(0, 0.1, 0)
obj.YROTATE(0.5)
obj.DELETE()
PRINT "smoke_object_dot_extra ok"
