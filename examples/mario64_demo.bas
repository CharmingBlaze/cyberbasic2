// Mario-style 3D playground: Bullet-style 3D physics (fallback engine), orbit camera, star.
// Run: cyberbasic examples/mario64_demo.bas
// Static level uses mass 0 + SetKinematic3D(1) so the ground is not pushed each frame (stable rest).
// Jump uses a short foot ray (below sphere), not repeated ground overlap checks.
// Controls: WASD | Space jump | Mouse orbit | ESC quit

InitWindow(1280, 720, "CyberBasic 64 Demo")
SetTargetFPS(60)

// --- Visual meshes (positions match physics bodies) ---
LoadCube(1, 1.2)
SetObjectColor(1, 220, 55, 45)
PositionObject(1, 0, 1.2, 0)

LoadCube(2, 1)
ScaleObject(2, 36, 0.35, 36)
SetObjectColor(2, 45, 140, 65)
PositionObject(2, 0, -0.175, 0)

LoadCube(3, 1)
ScaleObject(3, 2.2, 2.2, 2.2)
SetObjectColor(3, 255, 210, 60)
PositionObject(3, 7, 1.1, 4)

LoadCube(4, 1)
ScaleObject(4, 1.8, 1.2, 3.5)
SetObjectColor(4, 70, 120, 220)
PositionObject(4, -6, 0.6, 5)

LoadCube(5, 1)
ScaleObject(5, 3, 0.6, 1.5)
SetObjectColor(5, 200, 80, 200)
PositionObject(5, 0, 0.3, -8)

LoadCube(6, 1)
ScaleObject(6, 1.2, 3.5, 1.2)
SetObjectColor(6, 255, 140, 40)
PositionObject(6, -4, 1.75, -3)

LoadCube(7, 1)
ScaleObject(7, 2.5, 0.5, 2.5)
SetObjectColor(7, 180, 180, 40)
PositionObject(7, 5, 0.25, -5)

LoadCube(8, 1)
ScaleObject(8, 1.4, 1.4, 1.4)
SetObjectColor(8, 100, 200, 255)
PositionObject(8, -8, 0.7, -6)

MakeSphere(10, 0.55)
SetObjectColor(10, 255, 230, 60)
PositionObject(10, 0, 1.4, -10)

// --- Physics world (id m64) ---
CreateWorld3D("m64", 0, -20, 0)

// Ground: thin slab, top at y=0; kinematic so it never drifts when resolving contacts
CreateBox3D("m64", "ground", 0, -0.5, 0, 36, 1, 36, 0)
SetKinematic3D("m64", "ground", 1)
SetFriction3D("m64", "ground", 0.9)
SetRestitution3D("m64", "ground", 0)

// Course blocks (static + kinematic)
CreateBox3D("m64", "b3", 7, 1.1, 4, 2.2, 2.2, 2.2, 0)
SetKinematic3D("m64", "b3", 1)
SetFriction3D("m64", "b3", 0.85)
SetRestitution3D("m64", "b3", 0)

CreateBox3D("m64", "b4", -6, 0.6, 5, 1.8, 1.2, 3.5, 0)
SetKinematic3D("m64", "b4", 1)
SetFriction3D("m64", "b4", 0.85)
SetRestitution3D("m64", "b4", 0)

CreateBox3D("m64", "b5", 0, 0.3, -8, 3, 0.6, 1.5, 0)
SetKinematic3D("m64", "b5", 1)
SetFriction3D("m64", "b5", 0.85)
SetRestitution3D("m64", "b5", 0)

CreateBox3D("m64", "b6", -4, 1.75, -3, 1.2, 3.5, 1.2, 0)
SetKinematic3D("m64", "b6", 1)
SetFriction3D("m64", "b6", 0.85)
SetRestitution3D("m64", "b6", 0)

CreateBox3D("m64", "b7", 5, 0.25, -5, 2.5, 0.5, 2.5, 0)
SetKinematic3D("m64", "b7", 1)
SetFriction3D("m64", "b7", 0.85)
SetRestitution3D("m64", "b7", 0)

CreateBox3D("m64", "b8", -8, 0.7, -6, 1.4, 1.4, 1.4, 0)
SetKinematic3D("m64", "b8", 1)
SetFriction3D("m64", "b8", 0.85)
SetRestitution3D("m64", "b8", 0)

// Player sphere r=0.55, mass 1; start above spawn so first Step settles cleanly
CreateSphere3D("m64", "player", 0, 2.2, 0, 0.55, 1)
SetFriction3D("m64", "player", 0.35)
SetRestitution3D("m64", "player", 0)
SetDamping3D("m64", "player", 2.5, 2.5)

VAR moveSpeed = 6.0
VAR gotStar = 0
VAR sx = 0.0
VAR sy = 1.4
VAR sz = -10.0
VAR pr = 0.55
VAR footY = 0.0
VAR hitGround = 0
VAR nhx = 0.0
VAR nhz = 0.0
VAR len2 = 0.0
VAR sl = 0.0

SetCamera3D(0, 6, 14, 0, 0, 0, 0, 1, 0)

mainloop
  VAR dt = GetFrameTime()
  IF dt > 0.05 THEN
    LET dt = 0.05
  END IF

  Step3D("m64", dt)

  VAR px = GetPositionX3D("m64", "player")
  VAR py = GetPositionY3D("m64", "player")
  VAR pz = GetPositionZ3D("m64", "player")
  VAR vx = GetVelocityX3D("m64", "player")
  VAR vy = GetVelocityY3D("m64", "player")
  VAR vz = GetVelocityZ3D("m64", "player")

  IF px > 16 THEN
    SetPosition3D("m64", "player", 16, py, pz)
    LET px = 16
    LET vx = 0
  END IF
  IF px < -16 THEN
    SetPosition3D("m64", "player", -16, py, pz)
    LET px = -16
    LET vx = 0
  END IF
  IF pz > 16 THEN
    SetPosition3D("m64", "player", px, py, 16)
    LET pz = 16
    LET vz = 0
  END IF
  IF pz < -16 THEN
    SetPosition3D("m64", "player", px, py, -16)
    LET pz = -16
    LET vz = 0
  END IF

  LET nhx = 0
  LET nhz = 0
  IF IsKeyDown(KEY_D()) THEN
    LET nhx = nhx + 1
  END IF
  IF IsKeyDown(KEY_A()) THEN
    LET nhx = nhx - 1
  END IF
  IF IsKeyDown(KEY_W()) THEN
    LET nhz = nhz - 1
  END IF
  IF IsKeyDown(KEY_S()) THEN
    LET nhz = nhz + 1
  END IF
  LET len2 = nhx * nhx + nhz * nhz
  IF len2 > 0.0001 THEN
    LET sl = Sqrt(len2)
    LET nhx = nhx / sl * moveSpeed
    LET nhz = nhz / sl * moveSpeed
  END IF

  SetVelocity3D("m64", "player", nhx, vy, nhz)

  // Foot probe: vertical ray (dx=dz=0 fixed in engine). 9th arg skips "player" so we hit level, not self.
  LET footY = py - pr - 0.02
  LET hitGround = RayCastFromDir3D("m64", px, footY, pz, 0, -1, 0, 0.9, "player")
  IF IsKeyPressed(KEY_SPACE()) THEN
    IF hitGround = 1 AND vy < 2.0 THEN
      ApplyImpulse3D("m64", "player", 0, 8.5, 0)
    END IF
  END IF

  PositionObject(1, px, py, pz)

  OrbitCamera(px, py + 0.9, pz)

  ClearBackground(135, 185, 255, 255)
  DrawObjectRange(2, 8)
  IF gotStar = 0 THEN
    DrawObject(10)
  END IF
  DrawObject(1)
  DrawGrid(40, 1.0)

  DrawText("CyberBasic 64  |  Bullet 3D  |  WASD  Space jump  |  Mouse orbit", 12, 12, 18, 30, 30, 40, 255)
  IF gotStar = 1 THEN
    DrawText("STAR GET! You cleared the demo.", 12, 38, 26, 255, 220, 60, 255)
  ELSE
    DrawText("Foot-ray jump (no ground collision polling). Collect the star!", 12, 38, 16, 255, 255, 255, 255)
  END IF

  IF gotStar = 0 THEN
    IF CheckCollisionSpheres(px, py, pz, 0.5, sx, sy, sz, 0.55) THEN
      LET gotStar = 1
    END IF
  END IF

  SYNC
endmain

DestroyWorld3D("m64")
DeleteObjectRange(2, 8)
DeleteObject(1)
DeleteObject(10)
CloseWindow()
