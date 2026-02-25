// Exact pattern from plan: four IFs, no ENDIF
SUB movePlayer(dx, dy)
END SUB
VAR speed = 2
IF IsKeyDown(KEY_W) THEN movePlayer(0, -speed)
IF IsKeyDown(KEY_S) THEN movePlayer(0, speed)
IF IsKeyDown(KEY_A) THEN movePlayer(-speed, 0)
IF IsKeyDown(KEY_D) THEN movePlayer(speed, 0)
