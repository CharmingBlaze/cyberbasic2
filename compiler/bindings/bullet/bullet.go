// Package bullet exposes a 3D physics API to the CyberBasic VM as BULLET.*.
// Implemented in pure Go (no CGO). BASIC can call BULLET.CreateWorld, BULLET.CreateBox, etc.
// Same API can be wired to real Bullet Physics via CGO later.
package bullet

import (
	"cyberbasic/compiler/vm"
	"fmt"
	"math"
	"strconv"
	"sync"
)


func toFloat64(v interface{}) float64 {
	switch x := v.(type) {
	case int:
		return float64(x)
	case float64:
		return x
	case string:
		f, _ := strconv.ParseFloat(x, 64)
		return f
	default:
		return 0
	}
}

func toString(v interface{}) string {
	if v == nil {
		return ""
	}
	return fmt.Sprint(v)
}

func toInt(v interface{}) int {
	switch x := v.(type) {
	case int:
		return x
	case float64:
		return int(x)
	case string:
		n, _ := strconv.Atoi(x)
		return n
	default:
		return 0
	}
}

// --- Pure-Go 3D physics (Bullet-like API) ---

type vec3 struct{ x, y, z float64 }

type collisionHit struct {
	otherId string
	normal  vec3
}

type body struct {
	id              string
	position        vec3
	velocity        vec3
	rotation        vec3   // Euler angles
	angularVelocity vec3
	scale           vec3   // scale factors (default 1,1,1)
	halfExt         vec3   // half extents for box/cylinder
	radius          float64 // for sphere (0 = box/cylinder)
	mass            float64
	active          bool
	collisions      []collisionHit // filled each Step, cleared at start
}

type world struct {
	gravity vec3
	bodies  map[string]*body
	mu      sync.RWMutex
}

const defaultPhysicsWorld = "default"

var (
	worlds    = make(map[string]*world)
	worldMu   sync.RWMutex
	physicsBodySeq int
	lastRay  struct {
		hit    bool
		p      vec3
		bodyId string
		normal vec3
	}
	lastRayMu sync.Mutex
)

func getWorld(id string) *world {
	worldMu.RLock()
	defer worldMu.RUnlock()
	return worlds[id]
}

func getOrCreateWorld(id string, gx, gy, gz float64) *world {
	worldMu.Lock()
	defer worldMu.Unlock()
	if w, ok := worlds[id]; ok {
		w.gravity = vec3{gx, gy, gz}
		return w
	}
	w := &world{
		gravity: vec3{gx, gy, gz},
		bodies:  make(map[string]*body),
	}
	worlds[id] = w
	return w
}
func getBody(w *world, id string) *body {
	if w == nil {
		return nil
	}
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.bodies[id]
}

// RegisterBullet registers Bullet-style physics functions with the VM (BULLET.*).
func RegisterBullet(v *vm.VM) {
	// World
	v.RegisterForeign("BULLET.CreateWorld", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("CreateWorld requires (worldId, gravityX, gravityY, gravityZ)")
		}
		wid := toString(args[0])
		gx, gy, gz := toFloat64(args[1]), toFloat64(args[2]), toFloat64(args[3])
		getOrCreateWorld(wid, gx, gy, gz)
		return nil, nil
	})
	v.RegisterForeign("BULLET.SetGravity", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("SetGravity requires (worldId, gx, gy, gz)")
		}
		w := getWorld(toString(args[0]))
		if w == nil {
			return nil, fmt.Errorf("world not found")
		}
		w.mu.Lock()
		w.gravity = vec3{toFloat64(args[1]), toFloat64(args[2]), toFloat64(args[3])}
		w.mu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("BULLET.Step", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("Step requires (worldId, timeStep)")
		}
		w := getWorld(toString(args[0]))
		if w == nil {
			return nil, fmt.Errorf("world not found")
		}
		dt := toFloat64(args[1])
		w.mu.Lock()
		for _, b := range w.bodies {
			if !b.active || b.mass <= 0 {
				continue
			}
			b.velocity.x += w.gravity.x * dt
			b.velocity.y += w.gravity.y * dt
			b.velocity.z += w.gravity.z * dt
			b.position.x += b.velocity.x * dt
			b.position.y += b.velocity.y * dt
			b.position.z += b.velocity.z * dt
		}
		resolveCollisions(w)
		w.mu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("BULLET.DestroyWorld", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("DestroyWorld requires (worldId)")
		}
		worldMu.Lock()
		delete(worlds, toString(args[0]))
		worldMu.Unlock()
		return nil, nil
	})

	// Bodies - box
	v.RegisterForeign("BULLET.CreateBox", func(args []interface{}) (interface{}, error) {
		if len(args) < 9 {
			return nil, fmt.Errorf("CreateBox requires (worldId, bodyId, x, y, z, halfExX, halfExY, halfExZ, mass)")
		}
		wid := toString(args[0])
		bid := toString(args[1])
		w := getWorld(wid)
		if w == nil {
			w = getOrCreateWorld(wid, 0, -9.81, 0)
		}
		w.mu.Lock()
		w.bodies[bid] = &body{
			id:       bid,
			position: vec3{toFloat64(args[2]), toFloat64(args[3]), toFloat64(args[4])},
			halfExt:  vec3{toFloat64(args[5]), toFloat64(args[6]), toFloat64(args[7])},
			mass:     toFloat64(args[8]),
			active:   true,
			scale:    vec3{1, 1, 1},
		}
		w.mu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("BULLET.CreateSphere", func(args []interface{}) (interface{}, error) {
		if len(args) < 7 {
			return nil, fmt.Errorf("CreateSphere requires (worldId, bodyId, x, y, z, radius, mass)")
		}
		wid := toString(args[0])
		bid := toString(args[1])
		w := getWorld(wid)
		if w == nil {
			w = getOrCreateWorld(wid, 0, -9.81, 0)
		}
		w.mu.Lock()
		w.bodies[bid] = &body{
			id:       bid,
			position: vec3{toFloat64(args[2]), toFloat64(args[3]), toFloat64(args[4])},
			radius:   toFloat64(args[5]),
			mass:     toFloat64(args[6]),
			active:   true,
			scale:    vec3{1, 1, 1},
		}
		w.mu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("BULLET.DestroyBody", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("DestroyBody requires (worldId, bodyId)")
		}
		w := getWorld(toString(args[0]))
		if w == nil {
			return nil, nil
		}
		w.mu.Lock()
		delete(w.bodies, toString(args[1]))
		w.mu.Unlock()
		return nil, nil
	})

	// Position / velocity
	v.RegisterForeign("BULLET.SetPosition", func(args []interface{}) (interface{}, error) {
		if len(args) < 5 {
			return nil, fmt.Errorf("SetPosition requires (worldId, bodyId, x, y, z)")
		}
		b := getBody(getWorld(toString(args[0])), toString(args[1]))
		if b == nil {
			return nil, fmt.Errorf("body not found")
		}
		b.position = vec3{toFloat64(args[2]), toFloat64(args[3]), toFloat64(args[4])}
		return nil, nil
	})
	v.RegisterForeign("BULLET.GetPositionX", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("GetPositionX requires (worldId, bodyId)")
		}
		b := getBody(getWorld(toString(args[0])), toString(args[1]))
		if b == nil {
			return 0.0, nil
		}
		return b.position.x, nil
	})
	v.RegisterForeign("BULLET.GetPositionY", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("GetPositionY requires (worldId, bodyId)")
		}
		b := getBody(getWorld(toString(args[0])), toString(args[1]))
		if b == nil {
			return 0.0, nil
		}
		return b.position.y, nil
	})
	v.RegisterForeign("BULLET.GetPositionZ", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("GetPositionZ requires (worldId, bodyId)")
		}
		b := getBody(getWorld(toString(args[0])), toString(args[1]))
		if b == nil {
			return 0.0, nil
		}
		return b.position.z, nil
	})
	v.RegisterForeign("BULLET.SetVelocity", func(args []interface{}) (interface{}, error) {
		if len(args) < 5 {
			return nil, fmt.Errorf("SetVelocity requires (worldId, bodyId, vx, vy, vz)")
		}
		b := getBody(getWorld(toString(args[0])), toString(args[1]))
		if b == nil {
			return nil, fmt.Errorf("body not found")
		}
		b.velocity = vec3{toFloat64(args[2]), toFloat64(args[3]), toFloat64(args[4])}
		return nil, nil
	})
	v.RegisterForeign("BULLET.GetVelocityX", func(args []interface{}) (interface{}, error) {
		b := getBody(getWorld(toString(args[0])), toString(args[1]))
		if b == nil {
			return 0.0, nil
		}
		return b.velocity.x, nil
	})
	v.RegisterForeign("BULLET.GetVelocityY", func(args []interface{}) (interface{}, error) {
		b := getBody(getWorld(toString(args[0])), toString(args[1]))
		if b == nil {
			return 0.0, nil
		}
		return b.velocity.y, nil
	})
	v.RegisterForeign("BULLET.GetVelocityZ", func(args []interface{}) (interface{}, error) {
		b := getBody(getWorld(toString(args[0])), toString(args[1]))
		if b == nil {
			return 0.0, nil
		}
		return b.velocity.z, nil
	})
	v.RegisterForeign("BULLET.GetRotationX", func(args []interface{}) (interface{}, error) {
		b := getBody(getWorld(toString(args[0])), toString(args[1]))
		if b == nil {
			return 0.0, nil
		}
		return b.rotation.x, nil
	})
	v.RegisterForeign("BULLET.GetRotationY", func(args []interface{}) (interface{}, error) {
		b := getBody(getWorld(toString(args[0])), toString(args[1]))
		if b == nil {
			return 0.0, nil
		}
		return b.rotation.y, nil
	})
	v.RegisterForeign("BULLET.GetRotationZ", func(args []interface{}) (interface{}, error) {
		b := getBody(getWorld(toString(args[0])), toString(args[1]))
		if b == nil {
			return 0.0, nil
		}
		return b.rotation.z, nil
	})
	v.RegisterForeign("BULLET.SetRotation", func(args []interface{}) (interface{}, error) {
		if len(args) < 5 {
			return nil, fmt.Errorf("SetRotation requires (worldId, bodyId, rotX, rotY, rotZ)")
		}
		b := getBody(getWorld(toString(args[0])), toString(args[1]))
		if b == nil {
			return nil, fmt.Errorf("body not found")
		}
		b.rotation = vec3{toFloat64(args[2]), toFloat64(args[3]), toFloat64(args[4])}
		return nil, nil
	})
	v.RegisterForeign("BULLET.ApplyForce", func(args []interface{}) (interface{}, error) {
		if len(args) < 6 {
			return nil, fmt.Errorf("ApplyForce requires (worldId, bodyId, fx, fy, fz)")
		}
		b := getBody(getWorld(toString(args[0])), toString(args[1]))
		if b == nil || b.mass <= 0 {
			return nil, nil
		}
		dt := 1.0 / 60.0
		b.velocity.x += toFloat64(args[2]) / b.mass * dt
		b.velocity.y += toFloat64(args[3]) / b.mass * dt
		b.velocity.z += toFloat64(args[4]) / b.mass * dt
		return nil, nil
	})
	v.RegisterForeign("BULLET.ApplyCentralForce", func(args []interface{}) (interface{}, error) {
		if len(args) < 6 {
			return nil, fmt.Errorf("ApplyCentralForce requires (worldId, bodyId, fx, fy, fz)")
		}
		b := getBody(getWorld(toString(args[0])), toString(args[1]))
		if b == nil || b.mass <= 0 {
			return nil, nil
		}
		dt := 1.0 / 60.0
		b.velocity.x += toFloat64(args[2]) / b.mass * dt
		b.velocity.y += toFloat64(args[3]) / b.mass * dt
		b.velocity.z += toFloat64(args[4]) / b.mass * dt
		return nil, nil
	})
	v.RegisterForeign("BULLET.ApplyImpulse", func(args []interface{}) (interface{}, error) {
		if len(args) < 6 {
			return nil, fmt.Errorf("ApplyImpulse requires (worldId, bodyId, ix, iy, iz)")
		}
		b := getBody(getWorld(toString(args[0])), toString(args[1]))
		if b == nil || b.mass <= 0 {
			return nil, nil
		}
		m := b.mass
		b.velocity.x += toFloat64(args[2]) / m
		b.velocity.y += toFloat64(args[3]) / m
		b.velocity.z += toFloat64(args[4]) / m
		return nil, nil
	})

	// Raycast: returns 1 if hit, 0 otherwise. Use GetRayCastHitX/Y/Z for hit point.
	v.RegisterForeign("BULLET.RayCast", func(args []interface{}) (interface{}, error) {
		if len(args) < 8 {
			return nil, fmt.Errorf("RayCast requires (worldId, startX, startY, startZ, dirX, dirY, dirZ, maxDist)")
		}
		w := getWorld(toString(args[0]))
		if w == nil {
			lastRayMu.Lock()
			lastRay.hit = false
			lastRayMu.Unlock()
			return 0, nil
		}
		sx, sy, sz := toFloat64(args[1]), toFloat64(args[2]), toFloat64(args[3])
		dx, dy, dz := toFloat64(args[4]), toFloat64(args[5]), toFloat64(args[6])
		maxDist := toFloat64(args[7])
		norm := math.Sqrt(dx*dx + dy*dy + dz*dz)
		if norm < 1e-9 {
			lastRayMu.Lock()
			lastRay.hit = false
			lastRayMu.Unlock()
			return 0, nil
		}
		dx, dy, dz = dx/norm, dy/norm, dz/norm

		w.mu.RLock()
		var bestT float64 = 1e30
		hit := false
		var hitP vec3
		var hitBodyId string
		var hitNorm vec3
		for id, b := range w.bodies {
			if !b.active {
				continue
			}
			var min, max vec3
			if b.radius > 0 {
				min = vec3{b.position.x - b.radius, b.position.y - b.radius, b.position.z - b.radius}
				max = vec3{b.position.x + b.radius, b.position.y + b.radius, b.position.z + b.radius}
			} else {
				min = vec3{b.position.x - b.halfExt.x, b.position.y - b.halfExt.y, b.position.z - b.halfExt.z}
				max = vec3{b.position.x + b.halfExt.x, b.position.y + b.halfExt.y, b.position.z + b.halfExt.z}
			}
			t := rayAABB(sx, sy, sz, dx, dy, dz, min.x, min.y, min.z, max.x, max.y, max.z)
			if t >= 0 && t < maxDist && t < bestT {
				bestT = t
				hit = true
				hitP = vec3{sx + dx*t, sy + dy*t, sz + dz*t}
				hitBodyId = id
				hitNorm = vec3{-dx, -dy, -dz}
			}
		}
		w.mu.RUnlock()

		lastRayMu.Lock()
		lastRay.hit = hit
		lastRay.p = hitP
		lastRay.bodyId = hitBodyId
		lastRay.normal = hitNorm
		lastRayMu.Unlock()
		if hit {
			return 1, nil
		}
		return 0, nil
	})
	v.RegisterForeign("BULLET.GetRayCastHitX", func(args []interface{}) (interface{}, error) {
		lastRayMu.Lock()
		defer lastRayMu.Unlock()
		return lastRay.p.x, nil
	})
	v.RegisterForeign("BULLET.GetRayCastHitY", func(args []interface{}) (interface{}, error) {
		lastRayMu.Lock()
		defer lastRayMu.Unlock()
		return lastRay.p.y, nil
	})
	v.RegisterForeign("BULLET.GetRayCastHitZ", func(args []interface{}) (interface{}, error) {
		lastRayMu.Lock()
		defer lastRayMu.Unlock()
		return lastRay.p.z, nil
	})
	v.RegisterForeign("BULLET.GetRayCastHitBody", func(args []interface{}) (interface{}, error) {
		lastRayMu.Lock()
		defer lastRayMu.Unlock()
		return lastRay.bodyId, nil
	})
	v.RegisterForeign("BULLET.GetRayCastHitNormalX", func(args []interface{}) (interface{}, error) {
		lastRayMu.Lock()
		defer lastRayMu.Unlock()
		return lastRay.normal.x, nil
	})
	v.RegisterForeign("BULLET.GetRayCastHitNormalY", func(args []interface{}) (interface{}, error) {
		lastRayMu.Lock()
		defer lastRayMu.Unlock()
		return lastRay.normal.y, nil
	})
	v.RegisterForeign("BULLET.GetRayCastHitNormalZ", func(args []interface{}) (interface{}, error) {
		lastRayMu.Lock()
		defer lastRayMu.Unlock()
		return lastRay.normal.z, nil
	})

	// --- Flat 3D commands (no namespace, case-insensitive via VM) ---
	registerFlat3D(v)
}

// registerEntityGetters3D registers getters for entity.x, entity.y, entity.z, and rotation when the entity has "body" and "world" (3D physics).
func registerEntityGetters3D(v *vm.VM) {
	getWorldBody := func(entityName string) (worldId, bodyId string, ok bool) {
		g := v.Globals()[entityName]
		if g == nil {
			return "", "", false
		}
		m, ok := g.(map[string]interface{})
		if !ok {
			return "", "", false
		}
		w, _ := m["world"]
		b, _ := m["body"]
		if w == nil || b == nil {
			if w, _ := m["worldid"]; w != nil {
				if b, _ := m["bodyid"]; b != nil {
					return toString(w), toString(b), true
				}
			}
			return "", "", false
		}
		return toString(w), toString(b), true
	}
	v.RegisterEntityGetter("x", func(entityName, prop string) (vm.Value, bool) {
		worldId, bodyId, ok := getWorldBody(entityName)
		if !ok {
			return nil, false
		}
		return GetPositionX(worldId, bodyId), true
	})
	v.RegisterEntityGetter("y", func(entityName, prop string) (vm.Value, bool) {
		worldId, bodyId, ok := getWorldBody(entityName)
		if !ok {
			return nil, false
		}
		return GetPositionY(worldId, bodyId), true
	})
	v.RegisterEntityGetter("z", func(entityName, prop string) (vm.Value, bool) {
		worldId, bodyId, ok := getWorldBody(entityName)
		if !ok {
			return nil, false
		}
		return GetPositionZ(worldId, bodyId), true
	})
	// Rotation: expose as yaw, pitch, roll (common names) or rotation.x/y/z
	v.RegisterEntityGetter("yaw", func(entityName, prop string) (vm.Value, bool) {
		worldId, bodyId, ok := getWorldBody(entityName)
		if !ok {
			return nil, false
		}
		b := getBody(getWorld(worldId), bodyId)
		if b == nil {
			return nil, false
		}
		return b.rotation.y, true
	})
	v.RegisterEntityGetter("pitch", func(entityName, prop string) (vm.Value, bool) {
		worldId, bodyId, ok := getWorldBody(entityName)
		if !ok {
			return nil, false
		}
		b := getBody(getWorld(worldId), bodyId)
		if b == nil {
			return nil, false
		}
		return b.rotation.x, true
	})
	v.RegisterEntityGetter("roll", func(entityName, prop string) (vm.Value, bool) {
		worldId, bodyId, ok := getWorldBody(entityName)
		if !ok {
			return nil, false
		}
		b := getBody(getWorld(worldId), bodyId)
		if b == nil {
			return nil, false
		}
		return b.rotation.z, true
	})
}

// registerFlat3D registers flat CreateWorld3D, Step3D, CreateBox3D, etc. (no BULLET. prefix).
func registerFlat3D(v *vm.VM) {
	// World
	v.RegisterForeign("CreateWorld3D", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("CreateWorld3D requires (worldName$, gravityX, gravityY, gravityZ)")
		}
		getOrCreateWorld(toString(args[0]), toFloat64(args[1]), toFloat64(args[2]), toFloat64(args[3]))
		return nil, nil
	})
	v.RegisterForeign("DestroyWorld3D", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("DestroyWorld3D requires (worldName$)")
		}
		worldMu.Lock()
		delete(worlds, toString(args[0]))
		worldMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("Step3D", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("Step3D requires (worldName$, dt)")
		}
		w := getWorld(toString(args[0]))
		if w == nil {
			return nil, fmt.Errorf("world not found")
		}
		dt := toFloat64(args[1])
		w.mu.Lock()
		for _, b := range w.bodies {
			if !b.active || b.mass <= 0 {
				continue
			}
			b.velocity.x += w.gravity.x * dt
			b.velocity.y += w.gravity.y * dt
			b.velocity.z += w.gravity.z * dt
			b.position.x += b.velocity.x * dt
			b.position.y += b.velocity.y * dt
			b.position.z += b.velocity.z * dt
		}
		resolveCollisions(w)
		w.mu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("StepAllPhysics3D", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("StepAllPhysics3D requires (dt)")
		}
		dt := toFloat64(args[0])
		worldMu.RLock()
		ids := make([]string, 0, len(worlds))
		for id := range worlds {
			ids = append(ids, id)
		}
		worldMu.RUnlock()
		for _, id := range ids {
			w := getWorld(id)
			if w == nil {
				continue
			}
			w.mu.Lock()
			for _, b := range w.bodies {
				if !b.active || b.mass <= 0 {
					continue
				}
				b.velocity.x += w.gravity.x * dt
				b.velocity.y += w.gravity.y * dt
				b.velocity.z += w.gravity.z * dt
				b.position.x += b.velocity.x * dt
				b.position.y += b.velocity.y * dt
				b.position.z += b.velocity.z * dt
			}
			resolveCollisions(w)
			w.mu.Unlock()
		}
		return nil, nil
	})

	// Bodies
	v.RegisterForeign("CreateSphere3D", func(args []interface{}) (interface{}, error) {
		if len(args) < 7 {
			return nil, fmt.Errorf("CreateSphere3D requires (world$, body$, x, y, z, radius, mass)")
		}
		wid := toString(args[0])
		bid := toString(args[1])
		w := getWorld(wid)
		if w == nil {
			w = getOrCreateWorld(wid, 0, -9.81, 0)
		}
		w.mu.Lock()
		w.bodies[bid] = &body{
			id:       bid,
			position: vec3{toFloat64(args[2]), toFloat64(args[3]), toFloat64(args[4])},
			radius:   toFloat64(args[5]),
			mass:     toFloat64(args[6]),
			active:   true,
			scale:    vec3{1, 1, 1},
		}
		w.mu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("CreateBox3D", func(args []interface{}) (interface{}, error) {
		if len(args) < 9 {
			return nil, fmt.Errorf("CreateBox3D requires (world$, body$, x, y, z, sizeX, sizeY, sizeZ, mass)")
		}
		wid := toString(args[0])
		bid := toString(args[1])
		w := getWorld(wid)
		if w == nil {
			w = getOrCreateWorld(wid, 0, -9.81, 0)
		}
		sx, sy, sz := toFloat64(args[5]), toFloat64(args[6]), toFloat64(args[7])
		w.mu.Lock()
		w.bodies[bid] = &body{
			id:       bid,
			position: vec3{toFloat64(args[2]), toFloat64(args[3]), toFloat64(args[4])},
			halfExt:  vec3{sx / 2, sy / 2, sz / 2},
			mass:     toFloat64(args[8]),
			active:   true,
			scale:    vec3{1, 1, 1},
		}
		w.mu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("CreateCapsule3D", func(args []interface{}) (interface{}, error) {
		if len(args) < 8 {
			return nil, fmt.Errorf("CreateCapsule3D requires (world$, body$, x, y, z, radius, height, mass)")
		}
		wid := toString(args[0])
		bid := toString(args[1])
		w := getWorld(wid)
		if w == nil {
			w = getOrCreateWorld(wid, 0, -9.81, 0)
		}
		radius := toFloat64(args[5])
		w.mu.Lock()
		w.bodies[bid] = &body{
			id:       bid,
			position: vec3{toFloat64(args[2]), toFloat64(args[3]), toFloat64(args[4])},
			radius:   radius,
			mass:     toFloat64(args[7]),
			active:   true,
			scale:    vec3{1, 1, 1},
		}
		w.mu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("CreateStaticMesh3D", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("CreateStaticMesh3D requires (world$, body$, meshName$)")
		}
		wid := toString(args[0])
		bid := toString(args[1])
		w := getWorld(wid)
		if w == nil {
			w = getOrCreateWorld(wid, 0, -9.81, 0)
		}
		w.mu.Lock()
		w.bodies[bid] = &body{
			id:       bid,
			position: vec3{0, 0, 0},
			halfExt:  vec3{1, 1, 1},
			mass:     0,
			active:   true,
			scale:    vec3{1, 1, 1},
		}
		w.mu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("CreateCylinder3D", func(args []interface{}) (interface{}, error) {
		if len(args) < 8 {
			return nil, fmt.Errorf("CreateCylinder3D requires (world$, body$, x, y, z, radius, height, mass)")
		}
		wid := toString(args[0])
		bid := toString(args[1])
		w := getWorld(wid)
		if w == nil {
			w = getOrCreateWorld(wid, 0, -9.81, 0)
		}
		r := toFloat64(args[5])
		h := toFloat64(args[6])
		w.mu.Lock()
		w.bodies[bid] = &body{
			id:       bid,
			position: vec3{toFloat64(args[2]), toFloat64(args[3]), toFloat64(args[4])},
			halfExt:  vec3{r, r, h / 2},
			mass:     toFloat64(args[7]),
			active:   true,
			scale:    vec3{1, 1, 1},
		}
		w.mu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("CreateCone3D", func(args []interface{}) (interface{}, error) {
		if len(args) < 8 {
			return nil, fmt.Errorf("CreateCone3D requires (world$, body$, x, y, z, radius, height, mass)")
		}
		wid := toString(args[0])
		bid := toString(args[1])
		w := getWorld(wid)
		if w == nil {
			w = getOrCreateWorld(wid, 0, -9.81, 0)
		}
		r := toFloat64(args[5])
		h := toFloat64(args[6])
		w.mu.Lock()
		w.bodies[bid] = &body{
			id:       bid,
			position: vec3{toFloat64(args[2]), toFloat64(args[3]), toFloat64(args[4])},
			halfExt:  vec3{r, r, h / 2},
			mass:     toFloat64(args[7]),
			active:   true,
			scale:    vec3{1, 1, 1},
		}
		w.mu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("CreateHeightmap3D", func(args []interface{}) (interface{}, error) {
		return nil, nil
	})
	v.RegisterForeign("CreateCompound3D", func(args []interface{}) (interface{}, error) {
		return nil, nil
	})
	v.RegisterForeign("AddShapeToCompound3D", func(args []interface{}) (interface{}, error) {
		return nil, nil
	})

	// Position
	v.RegisterForeign("GetPositionX3D", func(args []interface{}) (interface{}, error) {
		b := getBody(getWorld(toString(args[0])), toString(args[1]))
		if b == nil {
			return 0.0, nil
		}
		return b.position.x, nil
	})
	v.RegisterForeign("GetPositionY3D", func(args []interface{}) (interface{}, error) {
		b := getBody(getWorld(toString(args[0])), toString(args[1]))
		if b == nil {
			return 0.0, nil
		}
		return b.position.y, nil
	})
	v.RegisterForeign("GetPositionZ3D", func(args []interface{}) (interface{}, error) {
		b := getBody(getWorld(toString(args[0])), toString(args[1]))
		if b == nil {
			return 0.0, nil
		}
		return b.position.z, nil
	})
	v.RegisterForeign("SetPosition3D", func(args []interface{}) (interface{}, error) {
		if len(args) < 5 {
			return nil, fmt.Errorf("SetPosition3D requires (world$, body$, x, y, z)")
		}
		b := getBody(getWorld(toString(args[0])), toString(args[1]))
		if b == nil {
			return nil, fmt.Errorf("body not found")
		}
		b.position = vec3{toFloat64(args[2]), toFloat64(args[3]), toFloat64(args[4])}
		return nil, nil
	})

	// Rotation (Euler: yaw, pitch, roll -> store as x=pitch, y=yaw, z=roll)
	v.RegisterForeign("GetYaw3D", func(args []interface{}) (interface{}, error) {
		b := getBody(getWorld(toString(args[0])), toString(args[1]))
		if b == nil {
			return 0.0, nil
		}
		return b.rotation.y, nil
	})
	v.RegisterForeign("GetPitch3D", func(args []interface{}) (interface{}, error) {
		b := getBody(getWorld(toString(args[0])), toString(args[1]))
		if b == nil {
			return 0.0, nil
		}
		return b.rotation.x, nil
	})
	v.RegisterForeign("GetRoll3D", func(args []interface{}) (interface{}, error) {
		b := getBody(getWorld(toString(args[0])), toString(args[1]))
		if b == nil {
			return 0.0, nil
		}
		return b.rotation.z, nil
	})
	v.RegisterForeign("SetRotation3D", func(args []interface{}) (interface{}, error) {
		if len(args) < 5 {
			return nil, fmt.Errorf("SetRotation3D requires (world$, body$, yaw, pitch, roll)")
		}
		b := getBody(getWorld(toString(args[0])), toString(args[1]))
		if b == nil {
			return nil, fmt.Errorf("body not found")
		}
		b.rotation = vec3{toFloat64(args[3]), toFloat64(args[2]), toFloat64(args[4])}
		return nil, nil
	})
	v.RegisterForeign("SetScale3D", func(args []interface{}) (interface{}, error) {
		if len(args) < 5 {
			return nil, fmt.Errorf("SetScale3D requires (world$, body$, sx, sy, sz)")
		}
		b := getBody(getWorld(toString(args[0])), toString(args[1]))
		if b == nil {
			return nil, fmt.Errorf("body not found")
		}
		b.scale = vec3{toFloat64(args[2]), toFloat64(args[3]), toFloat64(args[4])}
		return nil, nil
	})

	// Velocity
	v.RegisterForeign("GetVelocityX3D", func(args []interface{}) (interface{}, error) {
		b := getBody(getWorld(toString(args[0])), toString(args[1]))
		if b == nil {
			return 0.0, nil
		}
		return b.velocity.x, nil
	})
	v.RegisterForeign("GetVelocityY3D", func(args []interface{}) (interface{}, error) {
		b := getBody(getWorld(toString(args[0])), toString(args[1]))
		if b == nil {
			return 0.0, nil
		}
		return b.velocity.y, nil
	})
	v.RegisterForeign("GetVelocityZ3D", func(args []interface{}) (interface{}, error) {
		b := getBody(getWorld(toString(args[0])), toString(args[1]))
		if b == nil {
			return 0.0, nil
		}
		return b.velocity.z, nil
	})
	v.RegisterForeign("SetVelocity3D", func(args []interface{}) (interface{}, error) {
		if len(args) < 5 {
			return nil, fmt.Errorf("SetVelocity3D requires (world$, body$, vx, vy, vz)")
		}
		b := getBody(getWorld(toString(args[0])), toString(args[1]))
		if b == nil {
			return nil, fmt.Errorf("body not found")
		}
		b.velocity = vec3{toFloat64(args[2]), toFloat64(args[3]), toFloat64(args[4])}
		return nil, nil
	})
	v.RegisterForeign("SetAngularVelocity3D", func(args []interface{}) (interface{}, error) {
		if len(args) < 5 {
			return nil, fmt.Errorf("SetAngularVelocity3D requires (world$, body$, avx, avy, avz)")
		}
		b := getBody(getWorld(toString(args[0])), toString(args[1]))
		if b == nil {
			return nil, fmt.Errorf("body not found")
		}
		b.angularVelocity = vec3{toFloat64(args[2]), toFloat64(args[3]), toFloat64(args[4])}
		return nil, nil
	})
	v.RegisterForeign("GetAngularVelocityX3D", func(args []interface{}) (interface{}, error) {
		b := getBody(getWorld(toString(args[0])), toString(args[1]))
		if b == nil {
			return 0.0, nil
		}
		return b.angularVelocity.x, nil
	})
	v.RegisterForeign("GetAngularVelocityY3D", func(args []interface{}) (interface{}, error) {
		b := getBody(getWorld(toString(args[0])), toString(args[1]))
		if b == nil {
			return 0.0, nil
		}
		return b.angularVelocity.y, nil
	})
	v.RegisterForeign("GetAngularVelocityZ3D", func(args []interface{}) (interface{}, error) {
		b := getBody(getWorld(toString(args[0])), toString(args[1]))
		if b == nil {
			return 0.0, nil
		}
		return b.angularVelocity.z, nil
	})

	// Forces
	v.RegisterForeign("ApplyForce3D", func(args []interface{}) (interface{}, error) {
		if len(args) < 6 {
			return nil, fmt.Errorf("ApplyForce3D requires (world$, body$, fx, fy, fz)")
		}
		b := getBody(getWorld(toString(args[0])), toString(args[1]))
		if b == nil || b.mass <= 0 {
			return nil, nil
		}
		dt := 1.0 / 60.0
		b.velocity.x += toFloat64(args[2]) / b.mass * dt
		b.velocity.y += toFloat64(args[3]) / b.mass * dt
		b.velocity.z += toFloat64(args[4]) / b.mass * dt
		return nil, nil
	})
	v.RegisterForeign("ApplyImpulse3D", func(args []interface{}) (interface{}, error) {
		if len(args) < 6 {
			return nil, fmt.Errorf("ApplyImpulse3D requires (world$, body$, ix, iy, iz)")
		}
		b := getBody(getWorld(toString(args[0])), toString(args[1]))
		if b == nil || b.mass <= 0 {
			return nil, nil
		}
		m := b.mass
		b.velocity.x += toFloat64(args[2]) / m
		b.velocity.y += toFloat64(args[3]) / m
		b.velocity.z += toFloat64(args[4]) / m
		return nil, nil
	})
	v.RegisterForeign("ApplyTorque3D", func(args []interface{}) (interface{}, error) {
		return nil, nil
	})
	v.RegisterForeign("ApplyTorqueImpulse3D", func(args []interface{}) (interface{}, error) {
		return nil, nil
	})

	// Body properties (stubbed; engine does not model these yet)
	v.RegisterForeign("SetFriction3D", func(args []interface{}) (interface{}, error) { return nil, nil })
	v.RegisterForeign("SetRestitution3D", func(args []interface{}) (interface{}, error) { return nil, nil })
	v.RegisterForeign("SetDamping3D", func(args []interface{}) (interface{}, error) { return nil, nil })
	v.RegisterForeign("SetKinematic3D", func(args []interface{}) (interface{}, error) { return nil, nil })
	v.RegisterForeign("SetGravity3D", func(args []interface{}) (interface{}, error) { return nil, nil })
	v.RegisterForeign("SetMass3D", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("SetMass3D requires (world$, body$, mass)")
		}
		b := getBody(getWorld(toString(args[0])), toString(args[1]))
		if b == nil {
			return nil, nil
		}
		b.mass = toFloat64(args[2])
		return nil, nil
	})
	v.RegisterForeign("SetLinearFactor3D", func(args []interface{}) (interface{}, error) { return nil, nil })
	v.RegisterForeign("SetAngularFactor3D", func(args []interface{}) (interface{}, error) { return nil, nil })
	v.RegisterForeign("SetCCD3D", func(args []interface{}) (interface{}, error) { return nil, nil })

	// Joints (stubbed)
	v.RegisterForeign("CreateHingeJoint3D", func(args []interface{}) (interface{}, error) { return nil, nil })
	v.RegisterForeign("CreateSliderJoint3D", func(args []interface{}) (interface{}, error) { return nil, nil })
	v.RegisterForeign("CreateConeTwistJoint3D", func(args []interface{}) (interface{}, error) { return nil, nil })
	v.RegisterForeign("CreatePointToPointJoint3D", func(args []interface{}) (interface{}, error) { return nil, nil })
	v.RegisterForeign("CreateFixedJoint3D", func(args []interface{}) (interface{}, error) { return nil, nil })
	v.RegisterForeign("SetJointLimits3D", func(args []interface{}) (interface{}, error) { return nil, nil })
	v.RegisterForeign("SetJointMotor3D", func(args []interface{}) (interface{}, error) { return nil, nil })

	// Raycast (from->to)
	v.RegisterForeign("RayCast3D", func(args []interface{}) (interface{}, error) {
		if len(args) < 7 {
			return nil, fmt.Errorf("RayCast3D requires (world$, fromX, fromY, fromZ, toX, toY, toZ)")
		}
		w := getWorld(toString(args[0]))
		if w == nil {
			lastRayMu.Lock()
			lastRay.hit = false
			lastRayMu.Unlock()
			return 0, nil
		}
		fx, fy, fz := toFloat64(args[1]), toFloat64(args[2]), toFloat64(args[3])
		tx, ty, tz := toFloat64(args[4]), toFloat64(args[5]), toFloat64(args[6])
		dx, dy, dz := tx-fx, ty-fy, tz-fz
		maxDist := math.Sqrt(dx*dx + dy*dy + dz*dz)
		if maxDist < 1e-9 {
			lastRayMu.Lock()
			lastRay.hit = false
			lastRayMu.Unlock()
			return 0, nil
		}
		dx, dy, dz = dx/maxDist, dy/maxDist, dz/maxDist
		w.mu.RLock()
		var bestT float64 = 1e30
		hit := false
		var hitP vec3
		var hitBodyId string
		var hitNorm vec3
		for id, b := range w.bodies {
			if !b.active {
				continue
			}
			var min, max vec3
			if b.radius > 0 {
				min = vec3{b.position.x - b.radius, b.position.y - b.radius, b.position.z - b.radius}
				max = vec3{b.position.x + b.radius, b.position.y + b.radius, b.position.z + b.radius}
			} else {
				min = vec3{b.position.x - b.halfExt.x, b.position.y - b.halfExt.y, b.position.z - b.halfExt.z}
				max = vec3{b.position.x + b.halfExt.x, b.position.y + b.halfExt.y, b.position.z + b.halfExt.z}
			}
			t := rayAABB(fx, fy, fz, dx, dy, dz, min.x, min.y, min.z, max.x, max.y, max.z)
			if t >= 0 && t < maxDist && t < bestT {
				bestT = t
				hit = true
				hitP = vec3{fx + dx*t, fy + dy*t, fz + dz*t}
				hitBodyId = id
				hitNorm = vec3{-dx, -dy, -dz}
			}
		}
		w.mu.RUnlock()
		lastRayMu.Lock()
		lastRay.hit = hit
		lastRay.p = hitP
		lastRay.bodyId = hitBodyId
		lastRay.normal = hitNorm
		lastRayMu.Unlock()
		if hit {
			return 1, nil
		}
		return 0, nil
	})
	v.RegisterForeign("RayHitX3D", func(args []interface{}) (interface{}, error) {
		lastRayMu.Lock()
		defer lastRayMu.Unlock()
		return lastRay.p.x, nil
	})
	v.RegisterForeign("RayHitY3D", func(args []interface{}) (interface{}, error) {
		lastRayMu.Lock()
		defer lastRayMu.Unlock()
		return lastRay.p.y, nil
	})
	v.RegisterForeign("RayHitZ3D", func(args []interface{}) (interface{}, error) {
		lastRayMu.Lock()
		defer lastRayMu.Unlock()
		return lastRay.p.z, nil
	})
	v.RegisterForeign("RayHitBody3D", func(args []interface{}) (interface{}, error) {
		lastRayMu.Lock()
		defer lastRayMu.Unlock()
		return lastRay.bodyId, nil
	})
	v.RegisterForeign("RayHitNormalX3D", func(args []interface{}) (interface{}, error) {
		lastRayMu.Lock()
		defer lastRayMu.Unlock()
		return lastRay.normal.x, nil
	})
	v.RegisterForeign("RayHitNormalY3D", func(args []interface{}) (interface{}, error) {
		lastRayMu.Lock()
		defer lastRayMu.Unlock()
		return lastRay.normal.y, nil
	})
	v.RegisterForeign("RayHitNormalZ3D", func(args []interface{}) (interface{}, error) {
		lastRayMu.Lock()
		defer lastRayMu.Unlock()
		return lastRay.normal.z, nil
	})

	// Collision events
	v.RegisterForeign("GetCollisionCount3D", func(args []interface{}) (interface{}, error) {
		b := getBody(getWorld(toString(args[0])), toString(args[1]))
		if b == nil {
			return 0, nil
		}
		return len(b.collisions), nil
	})
	v.RegisterForeign("GetCollisionOther3D", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return "", nil
		}
		b := getBody(getWorld(toString(args[0])), toString(args[1]))
		if b == nil {
			return "", nil
		}
		idx := toInt(args[2])
		if idx < 0 || idx >= len(b.collisions) {
			return "", nil
		}
		return b.collisions[idx].otherId, nil
	})
	v.RegisterForeign("GetCollisionNormalX3D", func(args []interface{}) (interface{}, error) {
		b := getBody(getWorld(toString(args[0])), toString(args[1]))
		if b == nil {
			return 0.0, nil
		}
		idx := toInt(args[2])
		if idx < 0 || idx >= len(b.collisions) {
			return 0.0, nil
		}
		return b.collisions[idx].normal.x, nil
	})
	v.RegisterForeign("GetCollisionNormalY3D", func(args []interface{}) (interface{}, error) {
		b := getBody(getWorld(toString(args[0])), toString(args[1]))
		if b == nil {
			return 0.0, nil
		}
		idx := toInt(args[2])
		if idx < 0 || idx >= len(b.collisions) {
			return 0.0, nil
		}
		return b.collisions[idx].normal.y, nil
	})
	v.RegisterForeign("GetCollisionNormalZ3D", func(args []interface{}) (interface{}, error) {
		b := getBody(getWorld(toString(args[0])), toString(args[1]))
		if b == nil {
			return 0.0, nil
		}
		idx := toInt(args[2])
		if idx < 0 || idx >= len(b.collisions) {
			return 0.0, nil
		}
		return b.collisions[idx].normal.z, nil
	})

	// Entity property getters: when an entity has "body" and "world" (3D), entity.x/y/z and rotation come from physics.
	registerEntityGetters3D(v)

	// --- High-level physics (default world "default") ---
	v.RegisterForeign("PhysicsEnable", func(args []interface{}) (interface{}, error) {
		getOrCreateWorld(defaultPhysicsWorld, 0, -9.81, 0)
		return nil, nil
	})
	v.RegisterForeign("PhysicsDisable", func(args []interface{}) (interface{}, error) {
		worldMu.Lock()
		delete(worlds, defaultPhysicsWorld)
		worldMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("PhysicsSetGravity", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("PhysicsSetGravity requires (x, y, z)")
		}
		w := getOrCreateWorld(defaultPhysicsWorld, toFloat64(args[0]), toFloat64(args[1]), toFloat64(args[2]))
		w.mu.Lock()
		w.gravity = vec3{toFloat64(args[0]), toFloat64(args[1]), toFloat64(args[2])}
		w.mu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("CreateRigidBody", func(args []interface{}) (interface{}, error) {
		mass := 1.0
		if len(args) >= 2 {
			mass = toFloat64(args[1])
		}
		worldMu.Lock()
		physicsBodySeq++
		bid := fmt.Sprintf("body_%d", physicsBodySeq)
		worldMu.Unlock()
		w := getOrCreateWorld(defaultPhysicsWorld, 0, -9.81, 0)
		w.mu.Lock()
		w.bodies[bid] = &body{
			id:       bid,
			position: vec3{0, 0, 0},
			halfExt:  vec3{0.5, 0.5, 0.5},
			mass:     mass,
			active:   true,
			scale:    vec3{1, 1, 1},
		}
		w.mu.Unlock()
		return bid, nil
	})
	v.RegisterForeign("ApplyForce", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("ApplyForce requires (bodyId, fx, fy, fz)")
		}
		w := getWorld(defaultPhysicsWorld)
		if w == nil {
			return nil, nil
		}
		w.mu.Lock()
		b := w.bodies[toString(args[0])]
		if b != nil && b.mass > 0 {
			dt := 1.0 / 60.0
			b.velocity.x += toFloat64(args[1]) / b.mass * dt
			b.velocity.y += toFloat64(args[2]) / b.mass * dt
			b.velocity.z += toFloat64(args[3]) / b.mass * dt
		}
		w.mu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("ApplyImpulse", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("ApplyImpulse requires (bodyId, ix, iy, iz)")
		}
		w := getWorld(defaultPhysicsWorld)
		if w == nil {
			return nil, nil
		}
		w.mu.Lock()
		b := w.bodies[toString(args[0])]
		if b != nil && b.mass > 0 {
			m := b.mass
			b.velocity.x += toFloat64(args[1]) / m
			b.velocity.y += toFloat64(args[2]) / m
			b.velocity.z += toFloat64(args[3]) / m
		}
		w.mu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("SetBodyVelocity", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("SetBodyVelocity requires (bodyId, vx, vy, vz)")
		}
		w := getWorld(defaultPhysicsWorld)
		if w == nil {
			return nil, nil
		}
		w.mu.Lock()
		b := w.bodies[toString(args[0])]
		if b != nil {
			b.velocity = vec3{toFloat64(args[1]), toFloat64(args[2]), toFloat64(args[3])}
		}
		w.mu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("GetBodyVelocity", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return []interface{}{0.0, 0.0, 0.0}, nil
		}
		w := getWorld(defaultPhysicsWorld)
		if w == nil {
			return []interface{}{0.0, 0.0, 0.0}, nil
		}
		w.mu.RLock()
		b := w.bodies[toString(args[0])]
		vx, vy, vz := 0.0, 0.0, 0.0
		if b != nil {
			vx, vy, vz = b.velocity.x, b.velocity.y, b.velocity.z
		}
		w.mu.RUnlock()
		return []interface{}{vx, vy, vz}, nil
	})
	v.RegisterForeign("CheckCollision3D", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("CheckCollision3D requires (bodyIdA, bodyIdB)")
		}
		w := getWorld(defaultPhysicsWorld)
		if w == nil {
			return false, nil
		}
		w.mu.RLock()
		a := w.bodies[toString(args[0])]
		b := w.bodies[toString(args[1])]
		w.mu.RUnlock()
		if a == nil || b == nil {
			return false, nil
		}
		_, _, _, depth := overlapBodies(a, b)
		return depth > 0, nil
	})
}

// resolveCollisions resolves overlaps between bodies and records collision events. Call with w.mu held.
func resolveCollisions(w *world) {
	if w == nil {
		return
	}
	bodyList := make([]*body, 0, len(w.bodies))
	for _, b := range w.bodies {
		if b.active {
			b.collisions = nil
			bodyList = append(bodyList, b)
		}
	}
	for i, a := range bodyList {
		if a.mass <= 0 {
			continue
		}
		for j, b := range bodyList {
			if i == j {
				continue
			}
			nx, ny, nz, depth := overlapBodies(a, b)
			if depth <= 0 {
				continue
			}
			a.collisions = append(a.collisions, collisionHit{b.id, vec3{nx, ny, nz}})
			b.collisions = append(b.collisions, collisionHit{a.id, vec3{-nx, -ny, -nz}})
			// Push a out along normal; zero a's velocity along normal
			a.position.x += nx * depth
			a.position.y += ny * depth
			a.position.z += nz * depth
			vn := a.velocity.x*nx + a.velocity.y*ny + a.velocity.z*nz
			if vn < 0 {
				a.velocity.x -= vn * nx
				a.velocity.y -= vn * ny
				a.velocity.z -= vn * nz
			}
		}
	}
}

// overlapBodies returns contact normal (from a toward b) and penetration depth. Depth > 0 means overlap. Caller holds world lock.
func overlapBodies(a, b *body) (nx, ny, nz, depth float64) {
	if a.radius > 0 && b.radius > 0 {
		return overlapSphereSphere(a.position.x, a.position.y, a.position.z, a.radius, b.position.x, b.position.y, b.position.z, b.radius)
	}
	if a.radius > 0 && b.radius == 0 {
		return overlapSphereBox(a.position.x, a.position.y, a.position.z, a.radius, b.position.x, b.position.y, b.position.z, b.halfExt.x, b.halfExt.y, b.halfExt.z)
	}
	if a.radius == 0 && b.radius > 0 {
		nx, ny, nz, depth := overlapSphereBox(b.position.x, b.position.y, b.position.z, b.radius, a.position.x, a.position.y, a.position.z, a.halfExt.x, a.halfExt.y, a.halfExt.z)
		return -nx, -ny, -nz, depth
	}
	return overlapBoxBox(a.position.x, a.position.y, a.position.z, a.halfExt.x, a.halfExt.y, a.halfExt.z, b.position.x, b.position.y, b.position.z, b.halfExt.x, b.halfExt.y, b.halfExt.z)
}

func overlapSphereSphere(ax, ay, az, ar, bx, by, bz, br float64) (nx, ny, nz, depth float64) {
	dx := bx - ax
	dy := by - ay
	dz := bz - az
	distSq := dx*dx + dy*dy + dz*dz
	sum := ar + br
	sumSq := sum * sum
	if distSq >= sumSq || distSq < 1e-18 {
		return 0, 0, 0, 0
	}
	dist := math.Sqrt(distSq)
	depth = sum - dist
	nx = dx / dist
	ny = dy / dist
	nz = dz / dist
	return nx, ny, nz, depth
}

func overlapSphereBox(sx, sy, sz, sr, bx, by, bz, hx, hy, hz float64) (nx, ny, nz, depth float64) {
	clx := clamp(sx, bx-hx, bx+hx)
	cly := clamp(sy, by-hy, by+hy)
	clz := clamp(sz, bz-hz, bz+hz)
	dx := sx - clx
	dy := sy - cly
	dz := sz - clz
	distSq := dx*dx + dy*dy + dz*dz
	if distSq >= sr*sr || distSq < 1e-18 {
		return 0, 0, 0, 0
	}
	dist := math.Sqrt(distSq)
	depth = sr - dist
	nx = dx / dist
	ny = dy / dist
	nz = dz / dist
	return nx, ny, nz, depth
}

func clamp(v, lo, hi float64) float64 {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

func overlapBoxBox(ax, ay, az, ahx, ahy, ahz, bx, by, bz, bhx, bhy, bhz float64) (nx, ny, nz, depth float64) {
	aminX, amaxX := ax-ahx, ax+ahx
	aminY, amaxY := ay-ahy, ay+ahy
	aminZ, amaxZ := az-ahz, az+ahz
	bminX, bmaxX := bx-bhx, bx+bhx
	bminY, bmaxY := by-bhy, by+bhy
	bminZ, bmaxZ := bz-bhz, bz+bhz
	penX := math.Min(amaxX, bmaxX) - math.Max(aminX, bminX)
	penY := math.Min(amaxY, bmaxY) - math.Max(aminY, bminY)
	penZ := math.Min(amaxZ, bmaxZ) - math.Max(aminZ, bminZ)
	if penX <= 0 || penY <= 0 || penZ <= 0 {
		return 0, 0, 0, 0
	}
	// Smallest axis = contact normal
	if penX <= penY && penX <= penZ {
		if ax < bx {
			return 1, 0, 0, penX
		}
		return -1, 0, 0, penX
	}
	if penY <= penZ {
		if ay < by {
			return 0, 1, 0, penY
		}
		return 0, -1, 0, penY
	}
	if az < bz {
		return 0, 0, 1, penZ
	}
	return 0, 0, -1, penZ
}

// rayAABB returns ray parameter t for intersection with AABB, or -1 if no hit.
func rayAABB(ox, oy, oz, dx, dy, dz, minX, minY, minZ, maxX, maxY, maxZ float64) float64 {
	tmin := (minX - ox) / dx
	tmax := (maxX - ox) / dx
	if dx < 0 {
		tmin, tmax = tmax, tmin
	}
	tyMin := (minY - oy) / dy
	tyMax := (maxY - oy) / dy
	if dy < 0 {
		tyMin, tyMax = tyMax, tyMin
	}
	if tmin > tyMax || tyMin > tmax {
		return -1
	}
	if tyMin > tmin {
		tmin = tyMin
	}
	if tyMax < tmax {
		tmax = tyMax
	}
	tzMin := (minZ - oz) / dz
	tzMax := (maxZ - oz) / dz
	if dz < 0 {
		tzMin, tzMax = tzMax, tzMin
	}
	if tmin > tzMax || tzMin > tmax {
		return -1
	}
	if tzMin > tmin {
		tmin = tzMin
	}
	if tmin < 0 {
		return -1
	}
	return tmin
}

// Exported API for Go codegen (same behavior as VM foreign calls)

// CreateWorld creates or updates a physics world with the given gravity.
func CreateWorld(worldId string, gx, gy, gz float64) {
	getOrCreateWorld(worldId, gx, gy, gz)
}

// SetGravity sets the gravity of an existing world.
func SetGravity(worldId string, gx, gy, gz float64) {
	if w := getWorld(worldId); w != nil {
		w.mu.Lock()
		w.gravity = vec3{gx, gy, gz}
		w.mu.Unlock()
	}
}

// Step advances the simulation by timeStep seconds.
func Step(worldId string, timeStep float64) {
	w := getWorld(worldId)
	if w == nil {
		return
	}
	w.mu.Lock()
	for _, b := range w.bodies {
		if !b.active || b.mass <= 0 {
			continue
		}
		b.velocity.x += w.gravity.x * timeStep
		b.velocity.y += w.gravity.y * timeStep
		b.velocity.z += w.gravity.z * timeStep
		b.position.x += b.velocity.x * timeStep
		b.position.y += b.velocity.y * timeStep
		b.position.z += b.velocity.z * timeStep
	}
	w.mu.Unlock()
}

// CreateBox adds a box rigid body. halfEx* are half extents.
func CreateBox(worldId, bodyId string, x, y, z, halfExX, halfExY, halfExZ, mass float64) {
	w := getWorld(worldId)
	if w == nil {
		w = getOrCreateWorld(worldId, 0, -9.81, 0)
	}
	w.mu.Lock()
	w.bodies[bodyId] = &body{
		id:       bodyId,
		position: vec3{x, y, z},
		halfExt:  vec3{halfExX, halfExY, halfExZ},
		mass:     mass,
		active:   true,
		scale:    vec3{1, 1, 1},
	}
	w.mu.Unlock()
}

// CreateSphere adds a sphere rigid body.
func CreateSphere(worldId, bodyId string, x, y, z, radius, mass float64) {
	w := getWorld(worldId)
	if w == nil {
		w = getOrCreateWorld(worldId, 0, -9.81, 0)
	}
	w.mu.Lock()
	w.bodies[bodyId] = &body{
		id:       bodyId,
		position: vec3{x, y, z},
		radius:   radius,
		mass:     mass,
		active:   true,
		scale:    vec3{1, 1, 1},
	}
	w.mu.Unlock()
}

// DestroyBody removes a body from the world.
func DestroyBody(worldId, bodyId string) {
	if w := getWorld(worldId); w != nil {
		w.mu.Lock()
		delete(w.bodies, bodyId)
		w.mu.Unlock()
	}
}

// SetPosition sets the position of a body.
func SetPosition(worldId, bodyId string, x, y, z float64) {
	if b := getBody(getWorld(worldId), bodyId); b != nil {
		b.position = vec3{x, y, z}
	}
}

// GetPositionX returns the X position of the body.
func GetPositionX(worldId, bodyId string) float64 {
	if b := getBody(getWorld(worldId), bodyId); b != nil {
		return b.position.x
	}
	return 0
}

// GetPositionY returns the Y position of the body.
func GetPositionY(worldId, bodyId string) float64 {
	if b := getBody(getWorld(worldId), bodyId); b != nil {
		return b.position.y
	}
	return 0
}

// GetPositionZ returns the Z position of the body.
func GetPositionZ(worldId, bodyId string) float64 {
	if b := getBody(getWorld(worldId), bodyId); b != nil {
		return b.position.z
	}
	return 0
}

// SetVelocity sets the linear velocity of a body.
func SetVelocity(worldId, bodyId string, vx, vy, vz float64) {
	if b := getBody(getWorld(worldId), bodyId); b != nil {
		b.velocity = vec3{vx, vy, vz}
	}
}

// GetVelocityX returns the X velocity of a body.
func GetVelocityX(worldId, bodyId string) float64 {
	if b := getBody(getWorld(worldId), bodyId); b != nil {
		return b.velocity.x
	}
	return 0
}

// GetVelocityY returns the Y velocity of a body.
func GetVelocityY(worldId, bodyId string) float64 {
	if b := getBody(getWorld(worldId), bodyId); b != nil {
		return b.velocity.y
	}
	return 0
}

// GetVelocityZ returns the Z velocity of a body.
func GetVelocityZ(worldId, bodyId string) float64 {
	if b := getBody(getWorld(worldId), bodyId); b != nil {
		return b.velocity.z
	}
	return 0
}

// ApplyForce applies a central force (impulse over one frame at 60fps).
func ApplyForce(worldId, bodyId string, fx, fy, fz float64) {
	b := getBody(getWorld(worldId), bodyId)
	if b == nil || b.mass <= 0 {
		return
	}
	dt := 1.0 / 60.0
	b.velocity.x += fx / b.mass * dt
	b.velocity.y += fy / b.mass * dt
	b.velocity.z += fz / b.mass * dt
}

// RayCast performs a raycast; returns 1 if hit, 0 otherwise. Use GetRayCastHitX/Y/Z for hit point.
func RayCast(worldId string, startX, startY, startZ, dirX, dirY, dirZ, maxDist float64) int {
	w := getWorld(worldId)
	if w == nil {
		lastRayMu.Lock()
		lastRay.hit = false
		lastRayMu.Unlock()
		return 0
	}
	norm := math.Sqrt(dirX*dirX + dirY*dirY + dirZ*dirZ)
	if norm < 1e-9 {
		lastRayMu.Lock()
		lastRay.hit = false
		lastRayMu.Unlock()
		return 0
	}
	dx, dy, dz := dirX/norm, dirY/norm, dirZ/norm

	w.mu.RLock()
	var bestT float64 = 1e30
	hit := false
	var hitP vec3
	for _, b := range w.bodies {
		if !b.active {
			continue
		}
		var min, max vec3
		if b.radius > 0 {
			min = vec3{b.position.x - b.radius, b.position.y - b.radius, b.position.z - b.radius}
			max = vec3{b.position.x + b.radius, b.position.y + b.radius, b.position.z + b.radius}
		} else {
			min = vec3{b.position.x - b.halfExt.x, b.position.y - b.halfExt.y, b.position.z - b.halfExt.z}
			max = vec3{b.position.x + b.halfExt.x, b.position.y + b.halfExt.y, b.position.z + b.halfExt.z}
		}
		t := rayAABB(startX, startY, startZ, dx, dy, dz, min.x, min.y, min.z, max.x, max.y, max.z)
		if t >= 0 && t < maxDist && t < bestT {
			bestT = t
			hit = true
			hitP = vec3{startX + dx*t, startY + dy*t, startZ + dz*t}
		}
	}
	w.mu.RUnlock()

	lastRayMu.Lock()
	lastRay.hit = hit
	lastRay.p = hitP
	lastRayMu.Unlock()
	if hit {
		return 1
	}
	return 0
}

// GetRayCastHitX returns the X of the last raycast hit point.
func GetRayCastHitX() float64 {
	lastRayMu.Lock()
	defer lastRayMu.Unlock()
	return lastRay.p.x
}

// GetRayCastHitY returns the Y of the last raycast hit point.
func GetRayCastHitY() float64 {
	lastRayMu.Lock()
	defer lastRayMu.Unlock()
	return lastRay.p.y
}

// GetRayCastHitZ returns the Z of the last raycast hit point.
func GetRayCastHitZ() float64 {
	lastRayMu.Lock()
	defer lastRayMu.Unlock()
	return lastRay.p.z
}
