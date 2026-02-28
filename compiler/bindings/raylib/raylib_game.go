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
		orbitStateMu.Lock()
		orbitTargetX, orbitTargetY, orbitTargetZ = tx, ty, tz
		orbitAngle, orbitPitch, orbitDistance = angle, pitch, dist
		orbitInitialized = true
		orbitStateMu.Unlock()
		return nil, nil
	})
	// CameraOrbit (unprefixed alias); also updates orbit state for CameraZoom/CameraRotate/UpdateCamera
	v.RegisterForeign("CameraOrbit", func(args []interface{}) (interface{}, error) {
		if len(args) < 6 {
			return nil, fmt.Errorf("CameraOrbit requires (targetX, targetY, targetZ, angleRad, pitchRad, distance)")
		}
		tx := toFloat32(args[0])
		ty := toFloat32(args[1])
		tz := toFloat32(args[2])
		angle := toFloat32(args[3])
		pitch := toFloat32(args[4])
		dist := toFloat32(args[5])
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
		orbitStateMu.Lock()
		orbitTargetX, orbitTargetY, orbitTargetZ = tx, ty, tz
		orbitAngle, orbitPitch, orbitDistance = angle, pitch, dist
		orbitInitialized = true
		orbitStateMu.Unlock()
		return nil, nil
	})
	// CameraZoom(amount): adjust orbit distance with clamping (e.g. amount = GetMouseWheelMove(); default min 3, max 25)
	v.RegisterForeign("CameraZoom", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, nil
		}
		amount := toFloat32(args[0])
		orbitStateMu.Lock()
		orbitDistance -= amount * 1.5
		if orbitDistance < 3 {
			orbitDistance = 3
		}
		if orbitDistance > 25 {
			orbitDistance = 25
		}
		orbitStateMu.Unlock()
		return nil, nil
	})
	// CameraRotate(deltaX, deltaY): mouse-delta rotation; or CameraRotateDelta(deltaX, deltaY) alias
	v.RegisterForeign("CameraRotateDelta", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("CameraRotateDelta requires (deltaX, deltaY)")
		}
		dx := toFloat32(args[0])
		dy := toFloat32(args[1])
		const sens = float32(0.002)
		orbitStateMu.Lock()
		orbitAngle -= dx * sens
		orbitPitch += dy * sens
		if orbitPitch > 1.4 {
			orbitPitch = 1.4
		}
		if orbitPitch < -1.4 {
			orbitPitch = -1.4
		}
		orbitStateMu.Unlock()
		return nil, nil
	})
	// UpdateCamera(): apply orbit state to camera (position from target + angle/pitch/distance)
	v.RegisterForeign("UpdateCamera", func(args []interface{}) (interface{}, error) {
		orbitStateMu.Lock()
		tx, ty, tz := orbitTargetX, orbitTargetY, orbitTargetZ
		angle, pitch, dist := orbitAngle, orbitPitch, orbitDistance
		if dist == 0 {
			dist = 8
			orbitDistance = 8
		}
		orbitStateMu.Unlock()
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
		return nil, nil
	})
	// MouseOrbitCamera(): single call that does CameraRotateDelta(GetMouseDeltaX(), GetMouseDeltaY()), CameraZoom(GetMouseWheelMove()), UpdateCamera()
	v.RegisterForeign("MouseOrbitCamera", func(args []interface{}) (interface{}, error) {
		dx := float32(rl.GetMouseDelta().X)
		dy := float32(rl.GetMouseDelta().Y)
		wheel := float32(rl.GetMouseWheelMove())
		const sens = float32(0.002)
		orbitStateMu.Lock()
		if orbitDistance == 0 {
			orbitDistance = 8
		}
		orbitAngle -= dx * sens
		orbitPitch += dy * sens
		if orbitPitch > 1.4 {
			orbitPitch = 1.4
		}
		if orbitPitch < -1.4 {
			orbitPitch = -1.4
		}
		orbitDistance -= wheel * 1.5
		if orbitDistance < 3 {
			orbitDistance = 3
		}
		if orbitDistance > 25 {
			orbitDistance = 25
		}
		tx, ty, tz := orbitTargetX, orbitTargetY, orbitTargetZ
		angle, pitch, dist := orbitAngle, orbitPitch, orbitDistance
		orbitStateMu.Unlock()
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
		return nil, nil
	})
	// OrbitCamera(targetX, targetY, targetZ): orbit when right-mouse is held (drag = rotate), wheel = zoom anytime. Left click stays free for dropping etc.
	v.RegisterForeign("OrbitCamera", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("OrbitCamera requires (targetX, targetY, targetZ)")
		}
		tx := toFloat32(args[0])
		ty := toFloat32(args[1])
		tz := toFloat32(args[2])
		rightDown := rl.IsMouseButtonDown(rl.MouseButtonRight)
		dx := float32(rl.GetMouseDelta().X)
		dy := float32(rl.GetMouseDelta().Y)
		wheel := float32(rl.GetMouseWheelMove())
		const sens = float32(0.003)
		const zoomMul = float32(1.2)
		const pitchMax = float32(1.4)
		const distMin, distMax = float32(4), float32(35)
		orbitStateMu.Lock()
		orbitTargetX, orbitTargetY, orbitTargetZ = tx, ty, tz
		if !orbitInitialized {
			orbitAngle = 0
			orbitPitch = 0.25
			orbitDistance = 14
			orbitInitialized = true
		}
		if rightDown {
			orbitAngle -= dx * sens
			orbitPitch += dy * sens
			if orbitPitch > pitchMax {
				orbitPitch = pitchMax
			}
			if orbitPitch < -pitchMax {
				orbitPitch = -pitchMax
			}
		}
		orbitDistance -= wheel * zoomMul
		if orbitDistance < distMin {
			orbitDistance = distMin
		}
		if orbitDistance > distMax {
			orbitDistance = distMax
		}
		tx, ty, tz = orbitTargetX, orbitTargetY, orbitTargetZ
		angle, pitch, dist := orbitAngle, orbitPitch, orbitDistance
		orbitStateMu.Unlock()
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
	// MouseLook(): FPS-style camera; rotate view from mouse delta (uses camera position as eye, updates target)
	v.RegisterForeign("MouseLook", func(args []interface{}) (interface{}, error) {
		dx := float32(rl.GetMouseDelta().X)
		dy := float32(rl.GetMouseDelta().Y)
		const sens = float32(0.002)
		mouseLookMu.Lock()
		mouseLookYaw -= dx * sens
		mouseLookPitch += dy * sens
		if mouseLookPitch > 1.4 {
			mouseLookPitch = 1.4
		}
		if mouseLookPitch < -1.4 {
			mouseLookPitch = -1.4
		}
		yaw, pitch := mouseLookYaw, mouseLookPitch
		mouseLookMu.Unlock()
		px, py, pz := camera3D.Position.X, camera3D.Position.Y, camera3D.Position.Z
		cp := float32(math.Cos(float64(pitch)))
		sp := float32(math.Sin(float64(pitch)))
		ca := float32(math.Cos(float64(yaw)))
		sa := float32(math.Sin(float64(yaw)))
		const fwdDist float32 = 10
		tx := px + fwdDist*cp*sa
		ty := py + fwdDist*sp
		tz := pz + fwdDist*cp*ca
		camera3D.Target = rl.Vector3{X: tx, Y: ty, Z: tz}
		return nil, nil
	})
	// CameraLookAt(x, y, z): set camera target so camera looks at point
	v.RegisterForeign("CameraLookAt", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("CameraLookAt requires (x, y, z)")
		}
		camera3D.Target = rl.Vector3{X: toFloat32(args[0]), Y: toFloat32(args[1]), Z: toFloat32(args[2])}
		return nil, nil
	})
	// CameraMove(dx, dy, dz): move camera position and target by delta
	v.RegisterForeign("CameraMove", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("CameraMove requires (dx, dy, dz)")
		}
		dx, dy, dz := toFloat32(args[0]), toFloat32(args[1]), toFloat32(args[2])
		camera3D.Position.X += dx
		camera3D.Position.Y += dy
		camera3D.Position.Z += dz
		camera3D.Target.X += dx
		camera3D.Target.Y += dy
		camera3D.Target.Z += dz
		return nil, nil
	})
	// Camera3DMoveForward(amount): move camera position and target along forward (target - position).
	v.RegisterForeign("Camera3DMoveForward", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("Camera3DMoveForward requires (amount)")
		}
		amount := toFloat32(args[0])
		fwd := rl.Vector3Subtract(camera3D.Target, camera3D.Position)
		dist := rl.Vector3Length(fwd)
		if dist < 1e-6 {
			return nil, nil
		}
		fwd = rl.Vector3Scale(fwd, amount/dist)
		camera3D.Position = rl.Vector3Add(camera3D.Position, fwd)
		camera3D.Target = rl.Vector3Add(camera3D.Target, fwd)
		return nil, nil
	})
	// Camera3DMoveRight(amount): move camera position and target along right vector (cross(forward, up)).
	v.RegisterForeign("Camera3DMoveRight", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("Camera3DMoveRight requires (amount)")
		}
		amount := toFloat32(args[0])
		fwd := rl.Vector3Normalize(rl.Vector3Subtract(camera3D.Target, camera3D.Position))
		right := rl.Vector3CrossProduct(fwd, camera3D.Up)
		right = rl.Vector3Normalize(right)
		delta := rl.Vector3Scale(right, amount)
		camera3D.Position = rl.Vector3Add(camera3D.Position, delta)
		camera3D.Target = rl.Vector3Add(camera3D.Target, delta)
		return nil, nil
	})
	// Camera3DMoveBackward(amount): move camera and target backward (opposite of forward).
	v.RegisterForeign("Camera3DMoveBackward", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("Camera3DMoveBackward requires (amount)")
		}
		return v.CallForeign("Camera3DMoveForward", []interface{}{-toFloat64(args[0])})
	})
	// Camera3DMoveLeft(amount): move camera and target left (opposite of right).
	v.RegisterForeign("Camera3DMoveLeft", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("Camera3DMoveLeft requires (amount)")
		}
		return v.CallForeign("Camera3DMoveRight", []interface{}{-toFloat64(args[0])})
	})
	// Camera3DMoveUp(amount): move camera and target along camera up vector.
	v.RegisterForeign("Camera3DMoveUp", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("Camera3DMoveUp requires (amount)")
		}
		amount := toFloat32(args[0])
		up := rl.Vector3Normalize(camera3D.Up)
		delta := rl.Vector3Scale(up, amount)
		camera3D.Position = rl.Vector3Add(camera3D.Position, delta)
		camera3D.Target = rl.Vector3Add(camera3D.Target, delta)
		return nil, nil
	})
	// Camera3DMoveDown(amount): move camera and target opposite to up vector.
	v.RegisterForeign("Camera3DMoveDown", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("Camera3DMoveDown requires (amount)")
		}
		amount := -toFloat32(args[0])
		up := rl.Vector3Normalize(camera3D.Up)
		delta := rl.Vector3Scale(up, amount)
		camera3D.Position = rl.Vector3Add(camera3D.Position, delta)
		camera3D.Target = rl.Vector3Add(camera3D.Target, delta)
		return nil, nil
	})
	// Camera3DRotateYaw(angleRad): rotate camera position around target on world Y axis.
	v.RegisterForeign("Camera3DRotateYaw", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("Camera3DRotateYaw requires (angleRad)")
		}
		angle := toFloat32(args[0])
		rel := rl.Vector3Subtract(camera3D.Position, camera3D.Target)
		dx, dz := rel.X, rel.Z
		c := float32(math.Cos(float64(angle)))
		s := float32(math.Sin(float64(angle)))
		camera3D.Position.X = camera3D.Target.X + dx*c - dz*s
		camera3D.Position.Z = camera3D.Target.Z + dx*s + dz*c
		return nil, nil
	})
	// Camera3DRotatePitch(angleRad): rotate camera position toward/away from target (pitch).
	v.RegisterForeign("Camera3DRotatePitch", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("Camera3DRotatePitch requires (angleRad)")
		}
		angle := toFloat32(args[0])
		rel := rl.Vector3Subtract(camera3D.Position, camera3D.Target)
		dist := rl.Vector3Length(rel)
		if dist < 1e-6 {
			return nil, nil
		}
		// Current pitch from horizontal (Y); increase = look up.
		currPitch := float32(math.Asin(float64(rel.Y / dist)))
		newPitch := currPitch + angle
		const maxPitch = float32(1.4)
		if newPitch > maxPitch {
			newPitch = maxPitch
		}
		if newPitch < -maxPitch {
			newPitch = -maxPitch
		}
		cp := float32(math.Cos(float64(newPitch)))
		sp := float32(math.Sin(float64(newPitch)))
		horizLen := float32(math.Sqrt(float64(rel.X*rel.X + rel.Z*rel.Z)))
		if horizLen < 1e-6 {
			horizLen = 1
		}
		camera3D.Position.X = camera3D.Target.X + dist*cp*rel.X/horizLen
		camera3D.Position.Y = camera3D.Target.Y + dist*sp
		camera3D.Position.Z = camera3D.Target.Z + dist*cp*rel.Z/horizLen
		return nil, nil
	})
	// Camera3DRotateRoll(angleRad): rotate camera's up vector around the forward axis (position-target).
	v.RegisterForeign("Camera3DRotateRoll", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("Camera3DRotateRoll requires (angleRad)")
		}
		angle := toFloat32(args[0])
		fwd := rl.Vector3Subtract(camera3D.Position, camera3D.Target)
		dist := rl.Vector3Length(fwd)
		if dist < 1e-6 {
			return nil, nil
		}
		fwd = rl.Vector3Scale(fwd, 1/dist)
		up := camera3D.Up
		cr := float32(math.Cos(float64(angle)))
		sr := float32(math.Sin(float64(angle)))
		right := rl.Vector3CrossProduct(fwd, up)
		camera3D.Up = rl.Vector3Add(
			rl.Vector3Scale(up, cr),
			rl.Vector3Scale(right, -sr),
		)
		camera3D.Up = rl.Vector3Normalize(camera3D.Up)
		return nil, nil
	})
	// CameraRotate(deltaX, deltaY): mouse-delta rotation; or CameraRotate(pitchRad, yawRad, rollRad): absolute rotation
	v.RegisterForeign("CameraRotate", func(args []interface{}) (interface{}, error) {
		if len(args) == 2 {
			dx := toFloat32(args[0])
			dy := toFloat32(args[1])
			const sens = float32(0.002)
			orbitStateMu.Lock()
			orbitAngle -= dx * sens
			orbitPitch += dy * sens
			if orbitPitch > 1.4 {
				orbitPitch = 1.4
			}
			if orbitPitch < -1.4 {
				orbitPitch = -1.4
			}
			orbitStateMu.Unlock()
			return nil, nil
		}
		if len(args) < 3 {
			return nil, fmt.Errorf("CameraRotate requires (deltaX, deltaY) or (pitchRad, yawRad, rollRad)")
		}
		pitch := toFloat32(args[0])
		yaw := toFloat32(args[1])
		roll := toFloat32(args[2])
		tx, ty, tz := camera3D.Target.X, camera3D.Target.Y, camera3D.Target.Z
		px, py, pz := camera3D.Position.X, camera3D.Position.Y, camera3D.Position.Z
		dx, dy, dz := px-tx, py-ty, pz-tz
		dist := float32(math.Sqrt(float64(dx*dx + dy*dy + dz*dz)))
		if dist < 1e-6 {
			return nil, nil
		}
		// Current spherical angles (yaw around Y, pitch from horizontal)
		currYaw := float32(math.Atan2(float64(dx), float64(dz)))
		currPitch := float32(math.Asin(float64(dy / dist)))
		currYaw += yaw
		currPitch += pitch
		cp := float32(math.Cos(float64(currPitch)))
		sp := float32(math.Sin(float64(currPitch)))
		ca := float32(math.Cos(float64(currYaw)))
		sa := float32(math.Sin(float64(currYaw)))
		camera3D.Position.X = tx + dist*cp*sa
		camera3D.Position.Y = ty + dist*sp
		camera3D.Position.Z = tz + dist*cp*ca
		if math.Abs(float64(roll)) > 1e-6 {
			// Rotate Up vector in view plane (simplified: rotate around forward)
			fwd := rl.Vector3{X: -dx / dist, Y: -dy / dist, Z: -dz / dist}
			up := camera3D.Up
			cr := float32(math.Cos(float64(roll)))
			sr := float32(math.Sin(float64(roll)))
			right := rl.Vector3{
				X: up.Y*fwd.Z - up.Z*fwd.Y,
				Y: up.Z*fwd.X - up.X*fwd.Z,
				Z: up.X*fwd.Y - up.Y*fwd.X,
			}
			camera3D.Up = rl.Vector3{
				X: up.X*cr + right.X*sr,
				Y: up.Y*cr + right.Y*sr,
				Z: up.Z*cr + right.Z*sr,
			}
		}
		return nil, nil
	})
	// SetCameraFOV(fov): set global camera field of view (degrees)
	v.RegisterForeign("SetCameraFOV", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("SetCameraFOV requires (fov)")
		}
		camera3D.Fovy = toFloat32(args[0])
		return nil, nil
	})
	v.RegisterForeign("CameraSetFOV", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("CameraSetFOV requires (fov)")
		}
		camera3D.Fovy = toFloat32(args[0])
		return nil, nil
	})
	v.RegisterForeign("CameraSetClipping", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("CameraSetClipping requires (near, far)")
		}
		cameraNearZ = toFloat32(args[0])
		cameraFarZ = toFloat32(args[1])
		return nil, nil
	})
	v.RegisterForeign("CameraShake", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("CameraShake requires (amount, duration)")
		}
		cameraShakeMu.Lock()
		cameraShakeAmount = toFloat32(args[0])
		cameraShakeDuration = toFloat32(args[1])
		cameraShakeMu.Unlock()
		return nil, nil
	})
	// CameraFPS(): first-person preset; set position and target. Call MouseLook() each frame for mouse, use WASD to move.
	v.RegisterForeign("CameraFPS", func(args []interface{}) (interface{}, error) {
		camera3D.Position = rl.Vector3{X: 0, Y: 2, Z: 10}
		camera3D.Target = rl.Vector3{X: 0, Y: 0, Z: 0}
		camera3D.Up = rl.Vector3{X: 0, Y: 1, Z: 0}
		camera3D.Fovy = 60
		camera3D.Projection = rl.CameraPerspective
		mouseLookMu.Lock()
		mouseLookYaw = 0
		mouseLookPitch = 0
		mouseLookMu.Unlock()
		return nil, nil
	})
	// CameraFree(): free-fly noclip preset; same as CameraFPS for now.
	v.RegisterForeign("CameraFree", func(args []interface{}) (interface{}, error) {
		camera3D.Position = rl.Vector3{X: 0, Y: 2, Z: 10}
		camera3D.Target = rl.Vector3{X: 0, Y: 0, Z: 0}
		camera3D.Up = rl.Vector3{X: 0, Y: 1, Z: 0}
		camera3D.Fovy = 60
		camera3D.Projection = rl.CameraPerspective
		mouseLookMu.Lock()
		mouseLookYaw = 0
		mouseLookPitch = 0
		mouseLookMu.Unlock()
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

	// CollisionBox(x, y, z, w, h, d): create AABB with center (x,y,z) and size (w,h,d). Returns box id for CheckCollision/RayCast.
	v.RegisterForeign("CollisionBox", func(args []interface{}) (interface{}, error) {
		if len(args) < 6 {
			return nil, fmt.Errorf("CollisionBox requires (x, y, z, w, h, d)")
		}
		cx := toFloat64(args[0])
		cy := toFloat64(args[1])
		cz := toFloat64(args[2])
		w := toFloat64(args[3])
		h := toFloat64(args[4])
		d := toFloat64(args[5])
		hw, hh, hd := w/2, h/2, d/2
		collisionBoxMu.Lock()
		collisionBoxSeq++
		id := fmt.Sprintf("box_%d", collisionBoxSeq)
		collisionBoxes[id] = struct{ Cx, Cy, Cz, Hw, Hh, Hd float64 }{cx, cy, cz, hw, hh, hd}
		collisionBoxMu.Unlock()
		return id, nil
	})
	v.RegisterForeign("CheckCollision", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("CheckCollision requires (boxIdA, boxIdB)")
		}
		aId := toString(args[0])
		bId := toString(args[1])
		collisionBoxMu.Lock()
		a, okA := collisionBoxes[aId]
		b, okB := collisionBoxes[bId]
		collisionBoxMu.Unlock()
		if !okA || !okB {
			return false, nil
		}
		axMin, axMax := a.Cx-a.Hw, a.Cx+a.Hw
		ayMin, ayMax := a.Cy-a.Hh, a.Cy+a.Hh
		azMin, azMax := a.Cz-a.Hd, a.Cz+a.Hd
		bxMin, bxMax := b.Cx-b.Hw, b.Cx+b.Hw
		byMin, byMax := b.Cy-b.Hh, b.Cy+b.Hh
		bzMin, bzMax := b.Cz-b.Hd, b.Cz+b.Hd
		overlap := axMin <= bxMax && axMax >= bxMin && ayMin <= byMax && ayMax >= byMin && azMin <= bzMax && azMax >= bzMin
		return overlap, nil
	})
	// RayCast(originX, originY, originZ, dirX, dirY, dirZ [, boxId]): if boxId given, test vs that box and return hit distance (or -1). Without boxId, stores ray for GetRayCollision*.
	v.RegisterForeign("RayCast", func(args []interface{}) (interface{}, error) {
		if len(args) < 6 {
			return nil, fmt.Errorf("RayCast requires (originX,Y,Z, dirX,dirY,dirZ)")
		}
		ray := rl.Ray{
			Position:  rl.Vector3{X: toFloat32(args[0]), Y: toFloat32(args[1]), Z: toFloat32(args[2])},
			Direction: rl.Vector3{X: toFloat32(args[3]), Y: toFloat32(args[4]), Z: toFloat32(args[5])},
		}
		if len(args) >= 7 {
			boxId := toString(args[6])
			collisionBoxMu.Lock()
			b, ok := collisionBoxes[boxId]
			collisionBoxMu.Unlock()
			if !ok {
				return float64(-1), nil
			}
			box := rl.BoundingBox{
				Min: rl.Vector3{X: float32(b.Cx - b.Hw), Y: float32(b.Cy - b.Hh), Z: float32(b.Cz - b.Hd)},
				Max: rl.Vector3{X: float32(b.Cx + b.Hw), Y: float32(b.Cy + b.Hh), Z: float32(b.Cz + b.Hd)},
			}
			coll := rl.GetRayCollisionBox(ray, box)
			if !coll.Hit {
				return float64(-1), nil
			}
			return float64(coll.Distance), nil
		}
		// No boxId: could store ray for later GetRayCollision*; for now return 0
		_ = ray
		return float64(0), nil
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

	// 2D/3D convenience helpers (no namespace; for full 2D/3D games)
	v.RegisterForeign("GetScreenCenterX", func(args []interface{}) (interface{}, error) {
		return int32(rl.GetScreenWidth() / 2), nil
	})
	v.RegisterForeign("GetScreenCenterY", func(args []interface{}) (interface{}, error) {
		return int32(rl.GetScreenHeight() / 2), nil
	})
	v.RegisterForeign("Distance2D", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("Distance2D requires (x1, y1, x2, y2)")
		}
		v1 := rl.Vector2{X: toFloat32(args[0]), Y: toFloat32(args[1])}
		v2 := rl.Vector2{X: toFloat32(args[2]), Y: toFloat32(args[3])}
		return float64(rl.Vector2Distance(v1, v2)), nil
	})
	v.RegisterForeign("Distance3D", func(args []interface{}) (interface{}, error) {
		if len(args) < 6 {
			return nil, fmt.Errorf("Distance3D requires (x1, y1, z1, x2, y2, z2)")
		}
		v1 := rl.Vector3{X: toFloat32(args[0]), Y: toFloat32(args[1]), Z: toFloat32(args[2])}
		v2 := rl.Vector3{X: toFloat32(args[3]), Y: toFloat32(args[4]), Z: toFloat32(args[5])}
		return float64(rl.Vector3Distance(v1, v2)), nil
	})
	v.RegisterForeign("SetCamera2DCenter", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("SetCamera2DCenter requires (worldX, worldY)")
		}
		w := rl.GetScreenWidth()
		h := rl.GetScreenHeight()
		camera2D.Offset = rl.Vector2{X: float32(w) / 2, Y: float32(h) / 2}
		camera2D.Target = rl.Vector2{X: toFloat32(args[0]), Y: toFloat32(args[1])}
		camera2D.Rotation = 0
		camera2D.Zoom = 1
		return nil, nil
	})

	// Key constants (0-arg, return key code for use with RL.IsKeyDown)
	v.RegisterForeign("GAME.KEY_W", func(args []interface{}) (interface{}, error) { return KeyW, nil })
	v.RegisterForeign("GAME.KEY_A", func(args []interface{}) (interface{}, error) { return KeyA, nil })
	v.RegisterForeign("GAME.KEY_S", func(args []interface{}) (interface{}, error) { return KeyS, nil })
	v.RegisterForeign("GAME.KEY_D", func(args []interface{}) (interface{}, error) { return KeyD, nil })
	v.RegisterForeign("GAME.KEY_SPACE", func(args []interface{}) (interface{}, error) { return KeySpace, nil })
}
