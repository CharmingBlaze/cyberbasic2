// Package camera provides DBP-style camera management with integer IDs.
package camera

import (
	"math"
	"sync"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
	mu         sync.RWMutex
	cameras    = make(map[int]rl.Camera3D)
	activeID   int = -1
	activeCam  rl.Camera3D
	override   *rl.Camera3D
	overrideMu sync.RWMutex

	attachToObject       = make(map[int]int) // camID -> objID
	attachMu             sync.Mutex
	objectPositionGetter func(int) (float32, float32, float32)
)

// Camera is a 3D camera (alias for raylib Camera3D for clarity).
type Camera = rl.Camera3D

// Make creates a camera with the given integer ID.
func Make(id int) {
	mu.Lock()
	defer mu.Unlock()
	cameras[id] = rl.Camera3D{
		Position:   rl.Vector3{X: 0, Y: 0, Z: 0},
		Target:     rl.Vector3{X: 0, Y: 0, Z: -1},
		Up:         rl.Vector3{X: 0, Y: 1, Z: 0},
		Fovy:       60,
		Projection: rl.CameraPerspective,
	}
}

// SetPosition sets the camera position.
func SetPosition(id int, x, y, z float32) {
	mu.Lock()
	if c, ok := cameras[id]; ok {
		c.Position = rl.Vector3{X: x, Y: y, Z: z}
		cameras[id] = c
		if id == activeID {
			activeCam.Position = c.Position
		}
	}
	mu.Unlock()
}

// SetTarget sets the camera target (look-at point).
func SetTarget(id int, x, y, z float32) {
	mu.Lock()
	if c, ok := cameras[id]; ok {
		c.Target = rl.Vector3{X: x, Y: y, Z: z}
		cameras[id] = c
		if id == activeID {
			activeCam.Target = c.Target
		}
	}
	mu.Unlock()
}

// SetActive sets the active camera for rendering.
func SetActive(id int) {
	mu.Lock()
	if c, ok := cameras[id]; ok {
		activeID = id
		activeCam = c
		overrideMu.Lock()
		override = &activeCam
		overrideMu.Unlock()
	}
	mu.Unlock()
}

// GetActive returns the active camera for the 3D pass.
func GetActive() rl.Camera3D {
	overrideMu.RLock()
	o := override
	overrideMu.RUnlock()
	if o != nil {
		return *o
	}
	mu.RLock()
	c, ok := cameras[activeID]
	mu.RUnlock()
	if ok {
		return c
	}
	return rl.Camera3D{
		Position:   rl.Vector3{X: 0, Y: 2, Z: 10},
		Target:     rl.Vector3{X: 0, Y: 0, Z: 0},
		Up:         rl.Vector3{X: 0, Y: 1, Z: 0},
		Fovy:       60,
		Projection: rl.CameraPerspective,
	}
}

// SetOverride sets an external camera override. Pass nil to clear (use raylib camera3D).
func SetOverride(c *rl.Camera3D) {
	overrideMu.Lock()
	defer overrideMu.Unlock()
	override = c
}

// ClearOverride clears the override so GetCamera3D returns raylib's camera3D.
func ClearOverride() {
	overrideMu.Lock()
	defer overrideMu.Unlock()
	override = nil
}

// Exists returns true if the camera ID exists.
func Exists(id int) bool {
	mu.RLock()
	defer mu.RUnlock()
	_, ok := cameras[id]
	return ok
}

// Delete removes a camera. Clears active if it was the active camera.
func Delete(id int) {
	mu.Lock()
	delete(cameras, id)
	if activeID == id {
		activeID = -1
		overrideMu.Lock()
		override = nil
		overrideMu.Unlock()
	}
	mu.Unlock()
	attachMu.Lock()
	delete(attachToObject, id)
	attachMu.Unlock()
}

// Rotate sets camera rotation (pitch, yaw, roll in degrees) by computing new target from position.
func Rotate(id int, pitch, yaw, roll float32) {
	mu.Lock()
	c, ok := cameras[id]
	if !ok {
		mu.Unlock()
		return
	}
	front := rl.Vector3{
		X: c.Target.X - c.Position.X,
		Y: c.Target.Y - c.Position.Y,
		Z: c.Target.Z - c.Position.Z,
	}
	dist := float32(math.Sqrt(float64(front.X*front.X + front.Y*front.Y + front.Z*front.Z)))
	if dist < 0.001 {
		dist = 1
		front = rl.Vector3{X: 0, Y: 0, Z: -1}
	} else {
		front.X /= dist
		front.Y /= dist
		front.Z /= dist
	}
	// Apply yaw (Y axis), pitch (X axis), roll (Z axis) - simplified
	py, yy, ry := float32(math.Pi)*pitch/180, float32(math.Pi)*yaw/180, float32(math.Pi)*roll/180
	cy, sy := float32(math.Cos(float64(py))), float32(math.Sin(float64(py)))
	cx, sx := float32(math.Cos(float64(yy))), float32(math.Sin(float64(yy)))
	// Rotate front by yaw then pitch
	nx := front.X*cx - front.Z*sx
	nz := front.X*sx + front.Z*cx
	ny := front.Y
	front.X = nx
	front.Y = ny*cy - nz*sy
	front.Z = ny*sy + nz*cy
	_ = ry // roll affects up vector; skip for simplicity
	c.Target = rl.Vector3{
		X: c.Position.X + front.X*dist,
		Y: c.Position.Y + front.Y*dist,
		Z: c.Position.Z + front.Z*dist,
	}
	cameras[id] = c
	if id == activeID {
		activeCam = c
	}
	mu.Unlock()
}

// SetObjectPositionGetter sets the callback to get object world position (for AttachCameraToObject).
func SetObjectPositionGetter(fn func(int) (float32, float32, float32)) {
	objectPositionGetter = fn
}

// SetAttachToObject parents a camera to an object. Camera position follows object each frame.
func SetAttachToObject(camID, objID int) {
	attachMu.Lock()
	if objID < 0 {
		delete(attachToObject, camID)
	} else {
		attachToObject[camID] = objID
	}
	attachMu.Unlock()
}

// UpdateAttachments updates camera positions for cameras attached to objects. Call each frame.
func UpdateAttachments() {
	attachMu.Lock()
	att := make(map[int]int)
	for k, v := range attachToObject {
		att[k] = v
	}
	attachMu.Unlock()
	if objectPositionGetter == nil {
		return
	}
	for camID, objID := range att {
		x, y, z := objectPositionGetter(objID)
		// Offset camera behind and above object
		SetPosition(camID, x, y+2, z+5)
		SetTarget(camID, x, y, z)
	}
}

// HasOverride returns true if an active camera is set via SetActive.
func HasOverride() bool {
	overrideMu.RLock()
	defer overrideMu.RUnlock()
	return override != nil
}
