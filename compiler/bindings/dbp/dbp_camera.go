// Package dbp: Camera extras - follow, orbit, shake, smooth.
//
// These commands extend the FPS camera with DBP-style camera behaviors:
//   - CameraFollow(objectID, distance): Camera follows an object at given distance
//   - CameraOrbit(x, y, z, angle, pitch, distance): Orbit around a target point
//   - CameraShake(amount, duration): Add screen shake effect
//   - CameraSmooth(value): Lerp factor for smooth camera movement
//
// Call FpsUpdate or use these in your draw loop. Camera state is applied
// when the next frame is rendered.
package dbp

import (
	"fmt"
	"math"
	"sync"
	"time"

	"cyberbasic/compiler/vm"
	rl "github.com/gen2brain/raylib-go/raylib"
)

// Camera follow state: follows an object by ID.
var (
	followTarget   int   = -1
	followDistance float32 = 5.0
	followMu       sync.Mutex
)

// Camera orbit state: orbits around a target point.
var (
	orbitTargetX   float32
	orbitTargetY   float32
	orbitTargetZ   float32
	orbitAngle     float32
	orbitPitch     float32
	orbitDistance  float32 = 10.0
	orbitActive    bool
	orbitMu        sync.Mutex
)

// Camera shake state: random offset for screen shake.
var (
	shakeAmount  float32
	shakeEndTime time.Time
	shakeActive  bool
	shakeMu      sync.Mutex
)

// Camera smooth: lerp factor (0=instant, 1=no movement).
var (
	smoothFactor float32 = 0.1
	smoothMu     sync.Mutex
)

// lastCamPos stores the last camera position for smooth lerping.
var (
	lastCamX, lastCamY, lastCamZ float32
	lastCamMu                    sync.Mutex
)

// registerCameraExtras adds CameraFollow, CameraOrbit, CameraShake, CameraSmooth.
func registerCameraExtras(v *vm.VM) {
	// CameraFollow(objectID, distance): Camera follows the object at given distance.
	v.RegisterForeign("CameraFollow", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("CameraFollow(objectID, distance) requires 2 arguments")
		}
		followMu.Lock()
		followTarget = toInt(args[0])
		followDistance = toFloat32(args[1])
		followMu.Unlock()
		return nil, nil
	})

	// CameraOrbit(x, y, z, angle, pitch, distance): Orbit around target point.
	v.RegisterForeign("CameraOrbit", func(args []interface{}) (interface{}, error) {
		if len(args) < 6 {
			return nil, fmt.Errorf("CameraOrbit(x, y, z, angle, pitch, distance) requires 6 arguments")
		}
		orbitMu.Lock()
		orbitTargetX = toFloat32(args[0])
		orbitTargetY = toFloat32(args[1])
		orbitTargetZ = toFloat32(args[2])
		orbitAngle = toFloat32(args[3])
		orbitPitch = toFloat32(args[4])
		orbitDistance = toFloat32(args[5])
		orbitActive = true
		orbitMu.Unlock()
		return nil, nil
	})

	// CameraShake(amount, duration): Add shake; duration in seconds.
	v.RegisterForeign("CameraShake", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("CameraShake(amount, duration) requires 2 arguments")
		}
		shakeMu.Lock()
		shakeAmount = toFloat32(args[0])
		dur := toFloat32(args[1])
		if dur > 0 {
			shakeEndTime = time.Now().Add(time.Duration(dur * float32(time.Second)))
			shakeActive = true
		}
		shakeMu.Unlock()
		return nil, nil
	})

	// CameraUpdate: Apply follow, orbit, shake, smooth. Call each frame when using those modes.
	v.RegisterForeign("CameraUpdate", func(args []interface{}) (interface{}, error) {
		ApplyCameraExtras(v)
		return nil, nil
	})

	// CameraSmooth(value): Lerp factor 0-1. Lower = snappier.
	v.RegisterForeign("CameraSmooth", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("CameraSmooth(value) requires 1 argument")
		}
		smoothMu.Lock()
		smoothFactor = toFloat32(args[0])
		if smoothFactor <= 0 {
			smoothFactor = 0.01
		}
		if smoothFactor > 1 {
			smoothFactor = 1
		}
		smoothMu.Unlock()
		return nil, nil
	})
}

// ApplyCameraExtras applies follow, orbit, shake, and smooth to the current camera.
// Call this each frame after FpsUpdate or when using these camera modes.
func ApplyCameraExtras(v *vm.VM) {
	// Get current camera position from FPS state
	fpsCameraMu.Lock()
	camX, camY, camZ := fpsCamX, fpsCamY, fpsCamZ
	fpsCameraMu.Unlock()

	// Initialize lastCam on first use
	lastCamMu.Lock()
	if lastCamX == 0 && lastCamY == 0 && lastCamZ == 0 {
		lastCamX, lastCamY, lastCamZ = camX, camY, camZ
	}
	lastCamMu.Unlock()

	// Apply follow: position camera behind object
	followMu.Lock()
	ft := followTarget
	fd := followDistance
	followMu.Unlock()
	if ft >= 0 {
		objectsMu.Lock()
		obj, ok := objects[ft]
		if ok {
			// Position camera behind object (simple offset)
			camX = obj.x
			camY = obj.y + 2
			camZ = obj.z + fd
		}
		objectsMu.Unlock()
	}

	// Apply orbit: spherical coordinates around target
	orbitMu.Lock()
	if orbitActive {
		angle := orbitAngle * float32(math.Pi) / 180
		pitch := orbitPitch * float32(math.Pi) / 180
		camX = orbitTargetX + orbitDistance*float32(math.Cos(float64(pitch)))*float32(math.Sin(float64(angle)))
		camY = orbitTargetY + orbitDistance*float32(math.Sin(float64(pitch)))
		camZ = orbitTargetZ + orbitDistance*float32(math.Cos(float64(pitch)))*float32(math.Cos(float64(angle)))
	}
	orbitMu.Unlock()

	// Apply shake: add random offset (only while duration active)
	shakeMu.Lock()
	if shakeActive && shakeAmount > 0 && time.Now().Before(shakeEndTime) {
		amt := shakeAmount * (float32(rl.GetRandomValue(-100, 100)) / 100)
		camX += amt
		camY += amt
		camZ += amt
	} else if shakeActive && time.Now().After(shakeEndTime) {
		shakeActive = false
	}
	shakeMu.Unlock()

	// Apply smooth: lerp from last position
	smoothMu.Lock()
	sf := smoothFactor
	smoothMu.Unlock()
	lastCamMu.Lock()
	lastCamX = lastCamX + (camX-lastCamX)*sf
	lastCamY = lastCamY + (camY-lastCamY)*sf
	lastCamZ = lastCamZ + (camZ-lastCamZ)*sf
	camX, camY, camZ = lastCamX, lastCamY, lastCamZ
	lastCamMu.Unlock()

	// Set camera position and target
	v.CallForeign("SetCameraPosition", []interface{}{camX, camY, camZ})
	orbitMu.Lock()
	oa := orbitActive
	ox, oy, oz := orbitTargetX, orbitTargetY, orbitTargetZ
	orbitMu.Unlock()
	if oa {
		v.CallForeign("SetCameraTarget", []interface{}{ox, oy, oz})
	} else if ft >= 0 {
		objectsMu.Lock()
		obj, ok := objects[ft]
		objectsMu.Unlock()
		if ok {
			v.CallForeign("SetCameraTarget", []interface{}{obj.x, obj.y, obj.z})
		}
	}
}
