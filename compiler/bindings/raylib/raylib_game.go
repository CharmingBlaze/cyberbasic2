// Package raylib: GAME.* helpers for 3D games (camera orbit, WASD movement, ground check).
package raylib

import (
	"fmt"
	"math"
	"sync"

	"cyberbasic/compiler/bindings/box2d"
	"cyberbasic/compiler/bindings/bullet"
	"cyberbasic/compiler/vm"
	rl "github.com/gen2brain/raylib-go/raylib"
)

// Key constants for use with IsKeyDown (ASCII/raylib values)
const (
	KeyW     = 87
	KeyA     = 65
	KeyS     = 83
	KeyD     = 68
	KeySpace = 32
)

var (
	cameraOrbitAngle3D  float64
	cameraOrbitPitch3D  float64
	cameraOrbitMu       sync.Mutex
	camera2DFollowWorld string
	camera2DFollowBody  string
	camera2DFollowXOff  float32
	camera2DFollowYOff  float32
	camera3DOrbitWorld  string
	camera3DOrbitBody   string
	camera3DOrbitDist   float32
	camera3DOrbitHeight float32
)

func registerGame(v *vm.VM) {
	// Camera: orbit around target. angleRad=yaw, pitchRad=pitch, distance=radius.
	v.RegisterForeign("GAME.CameraOrbit", func(args []interface{}) (interface{}, error) {
		if len(args) < 6 {
			return nil, fmt.Errorf("CameraOrbit requires (targetX, targetY, targetZ, angleRad, pitchRad, distance)")
		}
		tx := toFloat32(args[0])
		ty := toFloat32(args[1])
		tz := toFloat32(args[2])
		angle := toFloat32(args[3])
		pitch := toFloat32(args[4])
		dist := toFloat32(args[5])
		// Spherical to Cartesian: x = dist*cos(pitch)*sin(angle), y = dist*sin(pitch), z = dist*cos(pitch)*cos(angle)
		cp := float32(math.Cos(float64(pitch)))
		sp := float32(math.Sin(float64(pitch)))
		ca := float32(math.Cos(float64(angle)))
		sa := float32(math.Sin(float64(angle)))
		ex := tx + dist*cp*sa
		ey := ty + dist*sp
		ez := tz + dist*cp*ca
		camera3D.Position = rl.Vector3{X: ex, Y: ey, Z: ez}
		camera3D.Target = rl.Vector3{X: tx, Y: ty, Z: tz}
		camera3D.Up = rl.Vector3{X: 0, Y: 1, Z: 0}
		camera3D.Fovy = 60.0
		camera3D.Projection = rl.CameraPerspective
		return nil, nil
	})

	// MoveWASD: apply horizontal force from WASD relative to angleRad, jump if Space and on ground. Uses bullet.* and RL.IsKeyDown.
	v.RegisterForeign("GAME.MoveWASD", func(args []interface{}) (interface{}, error) {
		if len(args) < 6 {
			return nil, fmt.Errorf("MoveWASD requires (worldId, bodyId, angleRad, speed, jumpVel, dt)")
		}
		worldId := fmt.Sprint(args[0])
		bodyId := fmt.Sprint(args[1])
		angle := toFloat32(args[2])
		speed := toFloat32(args[3])
		jumpVel := toFloat32(args[4])
		cx := float32(math.Cos(float64(angle)))
		cz := float32(math.Sin(float64(angle)))
		moveX, moveZ := float32(0), float32(0)
		if rl.IsKeyDown(KeyW) {
			moveX += cx
			moveZ += cz
		}
		if rl.IsKeyDown(KeyS) {
			moveX -= cx
			moveZ -= cz
		}
		if rl.IsKeyDown(KeyD) {
			moveX += cz
			moveZ -= cx
		}
		if rl.IsKeyDown(KeyA) {
			moveX -= cz
			moveZ += cx
		}
		len2 := moveX*moveX + moveZ*moveZ
		if len2 > 0.0001 {
			len := float32(math.Sqrt(float64(len2)))
			moveX /= len
			moveZ /= len
			fx := float64(moveX * speed)
			fz := float64(moveZ * speed)
			bullet.ApplyForce(worldId, bodyId, fx, 0, fz)
		}
		py := bullet.GetPositionY(worldId, bodyId)
		if py <= 0.6 && rl.IsKeyDown(KeySpace) {
			vx := bullet.GetVelocityX(worldId, bodyId)
			vz := bullet.GetVelocityZ(worldId, bodyId)
			bullet.SetVelocity(worldId, bodyId, vx, float64(jumpVel), vz)
		}
		return nil, nil
	})

	v.RegisterForeign("GAME.OnGround", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("OnGround requires (worldId, bodyId, planeY, tolerance)")
		}
		worldId := fmt.Sprint(args[0])
		bodyId := fmt.Sprint(args[1])
		planeY := toFloat64(args[2])
		tol := toFloat64(args[3])
		py := bullet.GetPositionY(worldId, bodyId)
		// Assume body radius/halfExt ~0.5 for sphere; for box use 0.5 as default bottom
		if py <= planeY+tol && py >= planeY-tol {
			return 1, nil
		}
		return 0, nil
	})

	v.RegisterForeign("GAME.SnapToGround", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("SnapToGround requires (worldId, bodyId, planeY, bodyBottomOffset)")
		}
		worldId := fmt.Sprint(args[0])
		bodyId := fmt.Sprint(args[1])
		planeY := toFloat64(args[2])
		offset := toFloat64(args[3]) // e.g. 0.5 for sphere radius
		px := bullet.GetPositionX(worldId, bodyId)
		py := bullet.GetPositionY(worldId, bodyId)
		pz := bullet.GetPositionZ(worldId, bodyId)
		vx := bullet.GetVelocityX(worldId, bodyId)
		vz := bullet.GetVelocityZ(worldId, bodyId)
		if py < planeY+offset {
			bullet.SetPosition(worldId, bodyId, px, planeY+offset, pz)
			bullet.SetVelocity(worldId, bodyId, vx, 0, vz)
		}
		return nil, nil
	})

	// --- Flat 3D game helpers (no namespace) ---
	v.RegisterForeign("MoveWASD3D", func(args []interface{}) (interface{}, error) {
		if len(args) < 6 {
			return nil, fmt.Errorf("MoveWASD3D requires (world$, body$, camAngle, moveForce, jumpVel, dt)")
		}
		worldId := fmt.Sprint(args[0])
		bodyId := fmt.Sprint(args[1])
		angle := toFloat64(args[2])
		speed := toFloat64(args[3])
		jumpVel := toFloat64(args[4])
		cx := math.Cos(angle)
		cz := math.Sin(angle)
		moveX, moveZ := 0.0, 0.0
		if rl.IsKeyDown(KeyW) {
			moveX += cx
			moveZ += cz
		}
		if rl.IsKeyDown(KeyS) {
			moveX -= cx
			moveZ -= cz
		}
		if rl.IsKeyDown(KeyD) {
			moveX += cz
			moveZ -= cx
		}
		if rl.IsKeyDown(KeyA) {
			moveX -= cz
			moveZ += cx
		}
		len2 := moveX*moveX + moveZ*moveZ
		if len2 > 1e-6 {
			len := math.Sqrt(len2)
			moveX /= len
			moveZ /= len
			bullet.ApplyForce(worldId, bodyId, moveX*speed, 0, moveZ*speed)
		}
		py := bullet.GetPositionY(worldId, bodyId)
		if py <= 0.6 && rl.IsKeyDown(KeySpace) {
			vx := bullet.GetVelocityX(worldId, bodyId)
			vz := bullet.GetVelocityZ(worldId, bodyId)
			bullet.SetVelocity(worldId, bodyId, vx, jumpVel, vz)
		}
		return nil, nil
	})
	v.RegisterForeign("SnapToGround3D", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("SnapToGround3D requires (world$, body$, groundY, radius)")
		}
		worldId := fmt.Sprint(args[0])
		bodyId := fmt.Sprint(args[1])
		groundY := toFloat64(args[2])
		radius := toFloat64(args[3])
		px := bullet.GetPositionX(worldId, bodyId)
		py := bullet.GetPositionY(worldId, bodyId)
		pz := bullet.GetPositionZ(worldId, bodyId)
		vx := bullet.GetVelocityX(worldId, bodyId)
		vz := bullet.GetVelocityZ(worldId, bodyId)
		if py < groundY+radius {
			bullet.SetPosition(worldId, bodyId, px, groundY+radius, pz)
			bullet.SetVelocity(worldId, bodyId, vx, 0, vz)
		}
		return nil, nil
	})
	v.RegisterForeign("IsOnGround3D", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("IsOnGround3D requires (world$, body$)")
		}
		worldId := fmt.Sprint(args[0])
		bodyId := fmt.Sprint(args[1])
		px := bullet.GetPositionX(worldId, bodyId)
		py := bullet.GetPositionY(worldId, bodyId)
		pz := bullet.GetPositionZ(worldId, bodyId)
		// Ray from just below feet downward so we don't hit our own body
		hit := bullet.RayCast(worldId, px, py-0.6, pz, 0, -1, 0, 0.5)
		if hit != 0 {
			return 1, nil
		}
		return 0, nil
	})
	v.RegisterForeign("CameraOrbit3D", func(args []interface{}) (interface{}, error) {
		if len(args) < 6 {
			return nil, fmt.Errorf("CameraOrbit3D requires (world$, body$, heightOffset, distance, mouseSens, dt)")
		}
		worldId := fmt.Sprint(args[0])
		bodyId := fmt.Sprint(args[1])
		heightOffset := toFloat64(args[2])
		dist := toFloat32(args[3])
		mouseSens := toFloat64(args[4])
		_ = toFloat64(args[5])
		px := bullet.GetPositionX(worldId, bodyId)
		py := bullet.GetPositionY(worldId, bodyId)
		pz := bullet.GetPositionZ(worldId, bodyId)
		tx := float32(px)
		ty := float32(py + heightOffset)
		tz := float32(pz)
		delta := rl.GetMouseDelta()
		cameraOrbitMu.Lock()
		cameraOrbitAngle3D += float64(delta.X) * mouseSens * 0.002
		cameraOrbitPitch3D -= float64(delta.Y) * mouseSens * 0.002
		if cameraOrbitPitch3D > 1.4 {
			cameraOrbitPitch3D = 1.4
		}
		if cameraOrbitPitch3D < -1.4 {
			cameraOrbitPitch3D = -1.4
		}
		angle := cameraOrbitAngle3D
		pitch := cameraOrbitPitch3D
		cameraOrbitMu.Unlock()
		cp := float32(math.Cos(pitch))
		sp := float32(math.Sin(pitch))
		ca := float32(math.Cos(angle))
		sa := float32(math.Sin(angle))
		ex := tx + dist*cp*sa
		ey := ty + dist*sp
		ez := tz + dist*cp*ca
		camera3D.Position = rl.Vector3{X: ex, Y: ey, Z: ez}
		camera3D.Target = rl.Vector3{X: tx, Y: ty, Z: tz}
		camera3D.Up = rl.Vector3{X: 0, Y: 1, Z: 0}
		camera3D.Fovy = 60.0
		camera3D.Projection = rl.CameraPerspective
		return nil, nil
	})
	// SetCamera3DOrbit(worldId, bodyId, distance, heightOffset): store target; call UpdateCamera3D(angleRad, pitchRad) each frame (e.g. from mouse).
	v.RegisterForeign("GAME.SetCamera3DOrbit", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("SetCamera3DOrbit requires (worldId, bodyId, distance, heightOffset)")
		}
		camera3DOrbitWorld = fmt.Sprint(args[0])
		camera3DOrbitBody = fmt.Sprint(args[1])
		camera3DOrbitDist = toFloat32(args[2])
		camera3DOrbitHeight = toFloat32(args[3])
		return nil, nil
	})
	v.RegisterForeign("GAME.UpdateCamera3D", func(args []interface{}) (interface{}, error) {
		if camera3DOrbitWorld == "" || camera3DOrbitBody == "" || len(args) < 2 {
			return nil, nil
		}
		px := bullet.GetPositionX(camera3DOrbitWorld, camera3DOrbitBody)
		py := bullet.GetPositionY(camera3DOrbitWorld, camera3DOrbitBody)
		pz := bullet.GetPositionZ(camera3DOrbitWorld, camera3DOrbitBody)
		angle := toFloat32(args[0])
		pitch := toFloat32(args[1])
		tx := float32(px)
		ty := float32(py) + camera3DOrbitHeight
		tz := float32(pz)
		cp := float32(math.Cos(float64(pitch)))
		sp := float32(math.Sin(float64(pitch)))
		ca := float32(math.Cos(float64(angle)))
		sa := float32(math.Sin(float64(angle)))
		ex := tx + camera3DOrbitDist*cp*sa
		ey := ty + camera3DOrbitDist*sp
		ez := tz + camera3DOrbitDist*cp*ca
		camera3D.Position = rl.Vector3{X: ex, Y: ey, Z: ez}
		camera3D.Target = rl.Vector3{X: tx, Y: ty, Z: tz}
		camera3D.Up = rl.Vector3{X: 0, Y: 1, Z: 0}
		camera3D.Fovy = 60.0
		camera3D.Projection = rl.CameraPerspective
		return nil, nil
	})

	// --- Flat 2D game helpers (no namespace) ---
	v.RegisterForeign("MoveHorizontal2D", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("MoveHorizontal2D requires (world$, body$, direction, speed)")
		}
		worldId := fmt.Sprint(args[0])
		bodyId := fmt.Sprint(args[1])
		direction := toFloat64(args[2])
		speed := toFloat64(args[3])
		_, vy := box2d.GetLinearVelocity(worldId, bodyId)
		box2d.SetLinearVelocity(worldId, bodyId, direction*speed, vy)
		return nil, nil
	})
	v.RegisterForeign("Jump2D", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("Jump2D requires (world$, body$, impulse)")
		}
		worldId := fmt.Sprint(args[0])
		bodyId := fmt.Sprint(args[1])
		impulse := toFloat64(args[2])
		vx, vy := box2d.GetLinearVelocity(worldId, bodyId)
		box2d.SetLinearVelocity(worldId, bodyId, vx, vy+impulse)
		return nil, nil
	})
	v.RegisterForeign("IsOnGround2D", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("IsOnGround2D requires (world$, body$)")
		}
		worldId := fmt.Sprint(args[0])
		bodyId := fmt.Sprint(args[1])
		px, py := box2d.GetPosition(worldId, bodyId)
		// Ray from slightly below center downward to avoid hitting self
		hit, _, _, _, _, _ := box2d.RayCastQuery(worldId, px, py-0.05, px, py-0.5)
		if hit {
			return 1, nil
		}
		return 0, nil
	})
	v.RegisterForeign("ClampVelocity2D", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("ClampVelocity2D requires (world$, body$, maxX, maxY)")
		}
		worldId := fmt.Sprint(args[0])
		bodyId := fmt.Sprint(args[1])
		maxX := toFloat64(args[2])
		maxY := toFloat64(args[3])
		vx, vy := box2d.GetLinearVelocity(worldId, bodyId)
		if vx > maxX {
			vx = maxX
		}
		if vx < -maxX {
			vx = -maxX
		}
		if vy > maxY {
			vy = maxY
		}
		if vy < -maxY {
			vy = -maxY
		}
		box2d.SetLinearVelocity(worldId, bodyId, vx, vy)
		return nil, nil
	})
	v.RegisterForeign("MoveVertical2D", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("MoveVertical2D requires (world$, body$, direction, speed)")
		}
		worldId := fmt.Sprint(args[0])
		bodyId := fmt.Sprint(args[1])
		direction := toFloat64(args[2])
		speed := toFloat64(args[3])
		vx, _ := box2d.GetLinearVelocity(worldId, bodyId)
		box2d.SetLinearVelocity(worldId, bodyId, vx, direction*speed)
		return nil, nil
	})
	v.RegisterForeign("CameraFollow2D", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("CameraFollow2D requires (world$, body$, xOffset, yOffset)")
		}
		worldId := fmt.Sprint(args[0])
		bodyId := fmt.Sprint(args[1])
		ox := toFloat32(args[2])
		oy := toFloat32(args[3])
		px, py := box2d.GetPosition(worldId, bodyId)
		camera2D.Target = rl.Vector2{X: float32(px) + ox, Y: float32(py) + oy}
		_ = bodyId
		return nil, nil
	})
	// SetCamera2DFollow(worldId, bodyId, xOffset, yOffset): store target; call UpdateCamera2D() each frame to apply.
	v.RegisterForeign("GAME.SetCamera2DFollow", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("SetCamera2DFollow requires (worldId, bodyId, xOffset, yOffset)")
		}
		camera2DFollowWorld = fmt.Sprint(args[0])
		camera2DFollowBody = fmt.Sprint(args[1])
		camera2DFollowXOff = toFloat32(args[2])
		camera2DFollowYOff = toFloat32(args[3])
		return nil, nil
	})
	v.RegisterForeign("GAME.UpdateCamera2D", func(args []interface{}) (interface{}, error) {
		if camera2DFollowWorld == "" || camera2DFollowBody == "" {
			return nil, nil
		}
		px, py := box2d.GetPosition(camera2DFollowWorld, camera2DFollowBody)
		camera2D.Target = rl.Vector2{X: float32(px) + camera2DFollowXOff, Y: float32(py) + camera2DFollowYOff}
		return nil, nil
	})
	v.RegisterForeign("SetCamera2DFollow", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("SetCamera2DFollow requires (worldId, bodyId, xOffset, yOffset)")
		}
		camera2DFollowWorld = fmt.Sprint(args[0])
		camera2DFollowBody = fmt.Sprint(args[1])
		camera2DFollowXOff = toFloat32(args[2])
		camera2DFollowYOff = toFloat32(args[3])
		return nil, nil
	})
	v.RegisterForeign("UpdateCamera2D", func(args []interface{}) (interface{}, error) {
		if camera2DFollowWorld == "" || camera2DFollowBody == "" {
			return nil, nil
		}
		px, py := box2d.GetPosition(camera2DFollowWorld, camera2DFollowBody)
		camera2D.Target = rl.Vector2{X: float32(px) + camera2DFollowXOff, Y: float32(py) + camera2DFollowYOff}
		return nil, nil
	})
	// SyncSpriteToBody2D: set sprite position to Box2D body position (world coords → screen via camera). Use in draw loop.
	v.RegisterForeign("GAME.SyncSpriteToBody2D", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("SyncSpriteToBody2D requires (worldId, bodyId, spriteId)")
		}
		worldId := fmt.Sprint(args[0])
		bodyId := fmt.Sprint(args[1])
		spriteId := fmt.Sprint(args[2])
		px, py := box2d.GetPosition(worldId, bodyId)
		screen := rl.GetWorldToScreen2D(rl.Vector2{X: float32(px), Y: float32(py)}, camera2D)
		if rt := v.GetRuntime(); rt != nil {
			_ = rt.SetSpritePosition(spriteId, float64(screen.X), float64(screen.Y))
		}
		return nil, nil
	})
	v.RegisterForeign("SyncSpriteToBody2D", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("SyncSpriteToBody2D requires (worldId, bodyId, spriteId)")
		}
		worldId := fmt.Sprint(args[0])
		bodyId := fmt.Sprint(args[1])
		spriteId := fmt.Sprint(args[2])
		px, py := box2d.GetPosition(worldId, bodyId)
		screen := rl.GetWorldToScreen2D(rl.Vector2{X: float32(px), Y: float32(py)}, camera2D)
		if rt := v.GetRuntime(); rt != nil {
			_ = rt.SetSpritePosition(spriteId, float64(screen.X), float64(screen.Y))
		}
		return nil, nil
	})
	// SetCollisionHandler(bodyId, subName): when bodyId collides, call Sub subName(otherBodyId). Call ProcessCollisions2D after BOX2D.Step.
	v.RegisterForeign("GAME.SetCollisionHandler", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("SetCollisionHandler requires (bodyId, subName)")
		}
		v.RegisterCollisionHandler(fmt.Sprint(args[0]), fmt.Sprint(args[1]))
		return nil, nil
	})
	v.RegisterForeign("SetCollisionHandler", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("SetCollisionHandler requires (bodyId, subName)")
		}
		v.RegisterCollisionHandler(fmt.Sprint(args[0]), fmt.Sprint(args[1]))
		return nil, nil
	})
	// ProcessCollisions2D(worldId): call registered collision Subs for each collision this frame. Call after BOX2D.Step.
	v.RegisterForeign("GAME.ProcessCollisions2D", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("ProcessCollisions2D requires (worldId)")
		}
		worldId := fmt.Sprint(args[0])
		for bodyId, subName := range v.GetCollisionHandlers() {
			n := box2d.GetCollisionCountForBody(worldId, bodyId)
			for i := 0; i < n; i++ {
				other := box2d.GetCollisionOtherForBody(worldId, bodyId, i)
				_ = v.InvokeSub(subName, []interface{}{other})
			}
		}
		return nil, nil
	})
	v.RegisterForeign("ProcessCollisions2D", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("ProcessCollisions2D requires (worldId)")
		}
		worldId := fmt.Sprint(args[0])
		for bodyId, subName := range v.GetCollisionHandlers() {
			n := box2d.GetCollisionCountForBody(worldId, bodyId)
			for i := 0; i < n; i++ {
				other := box2d.GetCollisionOtherForBody(worldId, bodyId, i)
				_ = v.InvokeSub(subName, []interface{}{other})
			}
		}
		return nil, nil
	})
	v.RegisterForeign("SnapToPlatform2D", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("SnapToPlatform2D requires (world$, body$, platformY, tolerance)")
		}
		worldId := fmt.Sprint(args[0])
		bodyId := fmt.Sprint(args[1])
		platformY := toFloat64(args[2])
		tol := toFloat64(args[3])
		px, py := box2d.GetPosition(worldId, bodyId)
		if py < platformY+tol && py > platformY-tol {
			vx, _ := box2d.GetLinearVelocity(worldId, bodyId)
			angle := box2d.GetAngle(worldId, bodyId)
			box2d.SetTransform(worldId, bodyId, px, platformY, angle)
			box2d.SetLinearVelocity(worldId, bodyId, vx, 0)
		}
		return nil, nil
	})

	// Additional 3D helpers from super-set
	v.RegisterForeign("Jump3D", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("Jump3D requires (world$, body$, impulse)")
		}
		worldId := fmt.Sprint(args[0])
		bodyId := fmt.Sprint(args[1])
		impulse := toFloat64(args[2])
		vx := bullet.GetVelocityX(worldId, bodyId)
		vy := bullet.GetVelocityY(worldId, bodyId)
		vz := bullet.GetVelocityZ(worldId, bodyId)
		bullet.SetVelocity(worldId, bodyId, vx, vy+impulse, vz)
		return nil, nil
	})
	v.RegisterForeign("ClampVelocity3D", func(args []interface{}) (interface{}, error) {
		if len(args) < 5 {
			return nil, fmt.Errorf("ClampVelocity3D requires (world$, body$, maxX, maxY, maxZ)")
		}
		worldId := fmt.Sprint(args[0])
		bodyId := fmt.Sprint(args[1])
		maxX := toFloat64(args[2])
		maxY := toFloat64(args[3])
		maxZ := toFloat64(args[4])
		vx := bullet.GetVelocityX(worldId, bodyId)
		vy := bullet.GetVelocityY(worldId, bodyId)
		vz := bullet.GetVelocityZ(worldId, bodyId)
		if vx > maxX {
			vx = maxX
		}
		if vx < -maxX {
			vx = -maxX
		}
		if vy > maxY {
			vy = maxY
		}
		if vy < -maxY {
			vy = -maxY
		}
		if vz > maxZ {
			vz = maxZ
		}
		if vz < -maxZ {
			vz = -maxZ
		}
		bullet.SetVelocity(worldId, bodyId, vx, vy, vz)
		return nil, nil
	})
	v.RegisterForeign("CameraFollow3D", func(args []interface{}) (interface{}, error) {
		if len(args) < 5 {
			return nil, fmt.Errorf("CameraFollow3D requires (world$, body$, heightOffset, distance, smooth)")
		}
		worldId := fmt.Sprint(args[0])
		bodyId := fmt.Sprint(args[1])
		heightOffset := toFloat64(args[2])
		dist := toFloat32(args[3])
		smooth := toFloat32(args[4])
		px := bullet.GetPositionX(worldId, bodyId)
		py := bullet.GetPositionY(worldId, bodyId)
		pz := bullet.GetPositionZ(worldId, bodyId)
		tx := float32(px)
		ty := float32(py + heightOffset)
		tz := float32(pz)
		// Smooth follow: blend current target toward body
		curr := camera3D.Target
		if smooth <= 0 {
			smooth = 0.1
		}
		blend := 1.0 - smooth
		camera3D.Target = rl.Vector3{
			X: curr.X*blend + tx*(1-blend),
			Y: curr.Y*blend + ty*(1-blend),
			Z: curr.Z*blend + tz*(1-blend),
		}
		pos := camera3D.Position
		dx := camera3D.Target.X - pos.X
		dy := camera3D.Target.Y - pos.Y
		dz := camera3D.Target.Z - pos.Z
		len := float32(math.Sqrt(float64(dx*dx + dy*dy + dz*dz)))
		if len > 0.001 {
			dx /= len
			dy /= len
			dz /= len
			camera3D.Position = rl.Vector3{
				X: camera3D.Target.X - dx*dist,
				Y: camera3D.Target.Y - dy*dist,
				Z: camera3D.Target.Z - dz*dist,
			}
		}
		camera3D.Up = rl.Vector3{X: 0, Y: 1, Z: 0}
		camera3D.Fovy = 60.0
		camera3D.Projection = rl.CameraPerspective
		_ = bodyId
		return nil, nil
	})

	// Asset path: return "assets/" + filename for LoadTexture(AssetPath("hero.png")) convention
	v.RegisterForeign("GAME.AssetPath", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return "assets/", nil
		}
		return "assets/" + fmt.Sprint(args[0]), nil
	})
	v.RegisterForeign("AssetPath", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return "assets/", nil
		}
		return "assets/" + fmt.Sprint(args[0]), nil
	})
	// ClampDelta(maxDt): return min(GetFrameTime(), maxDt) for stable physics step
	v.RegisterForeign("GAME.ClampDelta", func(args []interface{}) (interface{}, error) {
		dt := rl.GetFrameTime()
		if len(args) >= 1 {
			maxDt := toFloat32(args[0])
			if maxDt > 0 && float32(dt) > maxDt {
				return float64(maxDt), nil
			}
		}
		return float64(dt), nil
	})
	v.RegisterForeign("ClampDelta", func(args []interface{}) (interface{}, error) {
		dt := rl.GetFrameTime()
		if len(args) >= 1 {
			maxDt := toFloat32(args[0])
			if maxDt > 0 && float32(dt) > maxDt {
				return float64(maxDt), nil
			}
		}
		return float64(dt), nil
	})
	// ShowDebug: draw FPS at (10,10); optional second line ShowDebug(line2) at (10, 34)
	v.RegisterForeign("GAME.ShowDebug", func(args []interface{}) (interface{}, error) {
		rl.DrawFPS(10, 10)
		if len(args) >= 1 {
			rl.DrawText(fmt.Sprint(args[0]), 10, 34, 20, rl.White)
		}
		return nil, nil
	})
	v.RegisterForeign("ShowDebug", func(args []interface{}) (interface{}, error) {
		rl.DrawFPS(10, 10)
		if len(args) >= 1 {
			rl.DrawText(fmt.Sprint(args[0]), 10, 34, 20, rl.White)
		}
		return nil, nil
	})

	// Input axes: -1, 0, or 1 for movement (W/S -> Y, A/D -> X). Use: x = x + speed*GetAxisX(), y = y + speed*GetAxisY()
	v.RegisterForeign("GAME.GetAxisX", func(args []interface{}) (interface{}, error) {
		n := 0
		if rl.IsKeyDown(KeyD) {
			n++
		}
		if rl.IsKeyDown(KeyA) {
			n--
		}
		return n, nil
	})
	v.RegisterForeign("GAME.GetAxisY", func(args []interface{}) (interface{}, error) {
		n := 0
		if rl.IsKeyDown(KeyS) {
			n++
		}
		if rl.IsKeyDown(KeyW) {
			n--
		}
		return n, nil
	})
	v.RegisterForeign("GetAxisX", func(args []interface{}) (interface{}, error) {
		n := 0
		if rl.IsKeyDown(KeyD) {
			n++
		}
		if rl.IsKeyDown(KeyA) {
			n--
		}
		return n, nil
	})
	v.RegisterForeign("GetAxisY", func(args []interface{}) (interface{}, error) {
		n := 0
		if rl.IsKeyDown(KeyS) {
			n++
		}
		if rl.IsKeyDown(KeyW) {
			n--
		}
		return n, nil
	})

	// Key constants (0-arg, return key code for use with RL.IsKeyDown)
	v.RegisterForeign("GAME.KEY_W", func(args []interface{}) (interface{}, error) { return KeyW, nil })
	v.RegisterForeign("GAME.KEY_A", func(args []interface{}) (interface{}, error) { return KeyA, nil })
	v.RegisterForeign("GAME.KEY_S", func(args []interface{}) (interface{}, error) { return KeyS, nil })
	v.RegisterForeign("GAME.KEY_D", func(args []interface{}) (interface{}, error) { return KeyD, nil })
	v.RegisterForeign("GAME.KEY_SPACE", func(args []interface{}) (interface{}, error) { return KeySpace, nil })
}
