// Package bullet exposes a 3D physics API to the CyberBasic VM as BULLET.*.
// Implemented in pure Go (no CGO). BASIC can call BULLET.CreateWorld, BULLET.CreateBox, etc.
// Same API can be wired to real Bullet Physics via CGO later.
package bullet

import (
	"bufio"
	"cyberbasic/compiler/vm"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
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

type triangle struct{ v0, v1, v2 vec3 }

type collisionHit struct {
	otherId string
	normal  vec3
}

type body struct {
	id              string
	position        vec3
	velocity        vec3
	rotation        vec3 // Euler angles
	angularVelocity vec3
	pendingTorque   vec3 // accumulated per frame; applied in integrateBody
	scale           vec3    // scale factors (default 1,1,1)
	halfExt         vec3    // half extents for box/cylinder
	radius          float64 // for sphere (0 = box/cylinder)
	mass            float64
	active          bool
	collisions      []collisionHit // filled each Step, cleared at start
	// body properties (used in Step and resolveCollisions)
	friction       float64
	restitution    float64
	linearDamping  float64
	angularDamping float64
	kinematic      bool
	gravityScale   float64
	linearFactor   vec3
	angularFactor  vec3
	ccd            bool
	// mesh collider (static only)
	meshTriangles []triangle
	meshAABBMin   vec3 // local space
	meshAABBMax   vec3 // local space
	// compound collider: multiple axis-aligned boxes in parent space (center offset + half extents)
	compound []struct{ ox, oy, oz, hx, hy, hz float64 }
}

func bodyUsesCapsuleBounds(b *body) bool {
	return b != nil && b.radius > 0 && (b.halfExt.x > b.radius || b.halfExt.y > b.radius || b.halfExt.z > b.radius)
}

func bodyIsMesh(b *body) bool {
	return b != nil && len(b.meshTriangles) > 0
}

func bodyIsCompound(b *body) bool {
	return b != nil && len(b.compound) > 0
}

// bodyBoxParams returns (center, halfExt) for box or mesh bodies. For sphere, halfExt is radius (caller must use sphere overlap).
func bodyBoxParams(b *body) (cx, cy, cz, hx, hy, hz float64) {
	if bodyIsCompound(b) {
		min, max := bodyAABB(b)
		cx = (min.x + max.x) / 2
		cy = (min.y + max.y) / 2
		cz = (min.z + max.z) / 2
		hx = (max.x - min.x) / 2
		hy = (max.y - min.y) / 2
		hz = (max.z - min.z) / 2
		return
	}
	if bodyIsMesh(b) {
		cx = b.position.x + (b.meshAABBMin.x+b.meshAABBMax.x)/2
		cy = b.position.y + (b.meshAABBMin.y+b.meshAABBMax.y)/2
		cz = b.position.z + (b.meshAABBMin.z+b.meshAABBMax.z)/2
		hx = (b.meshAABBMax.x - b.meshAABBMin.x) / 2
		hy = (b.meshAABBMax.y - b.meshAABBMin.y) / 2
		hz = (b.meshAABBMax.z - b.meshAABBMin.z) / 2
		return
	}
	cx, cy, cz = b.position.x, b.position.y, b.position.z
	hx, hy, hz = b.halfExt.x, b.halfExt.y, b.halfExt.z
	return
}

func bodyAABB(b *body) (min, max vec3) {
	if b == nil {
		return vec3{}, vec3{}
	}
	if bodyIsCompound(b) {
		first := true
		var cmin, cmax vec3
		for _, p := range b.compound {
			wx := b.position.x + p.ox
			wy := b.position.y + p.oy
			wz := b.position.z + p.oz
			pmin := vec3{wx - p.hx, wy - p.hy, wz - p.hz}
			pmax := vec3{wx + p.hx, wy + p.hy, wz + p.hz}
			if first {
				cmin, cmax = pmin, pmax
				first = false
				continue
			}
			if pmin.x < cmin.x {
				cmin.x = pmin.x
			}
			if pmin.y < cmin.y {
				cmin.y = pmin.y
			}
			if pmin.z < cmin.z {
				cmin.z = pmin.z
			}
			if pmax.x > cmax.x {
				cmax.x = pmax.x
			}
			if pmax.y > cmax.y {
				cmax.y = pmax.y
			}
			if pmax.z > cmax.z {
				cmax.z = pmax.z
			}
		}
		return cmin, cmax
	}
	if bodyIsMesh(b) {
		return vec3{b.position.x + b.meshAABBMin.x, b.position.y + b.meshAABBMin.y, b.position.z + b.meshAABBMin.z},
			vec3{b.position.x + b.meshAABBMax.x, b.position.y + b.meshAABBMax.y, b.position.z + b.meshAABBMax.z}
	}
	if bodyUsesCapsuleBounds(b) || b.radius == 0 {
		return vec3{b.position.x - b.halfExt.x, b.position.y - b.halfExt.y, b.position.z - b.halfExt.z},
			vec3{b.position.x + b.halfExt.x, b.position.y + b.halfExt.y, b.position.z + b.halfExt.z}
	}
	return vec3{b.position.x - b.radius, b.position.y - b.radius, b.position.z - b.radius},
		vec3{b.position.x + b.radius, b.position.y + b.radius, b.position.z + b.radius}
}

// loadOBJForCollision parses a minimal OBJ (v and f lines) and returns triangles plus AABB in local space.
func loadOBJForCollision(path string) (tris []triangle, minV, maxV vec3, err error) {
	path = filepath.Clean(path)
	f, err := os.Open(path)
	if err != nil {
		return nil, vec3{}, vec3{}, err
	}
	defer f.Close()
	var verts []vec3
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) < 4 {
			continue
		}
		if parts[0] == "v" {
			x, _ := strconv.ParseFloat(parts[1], 64)
			y, _ := strconv.ParseFloat(parts[2], 64)
			z, _ := strconv.ParseFloat(parts[3], 64)
			verts = append(verts, vec3{x, y, z})
			continue
		}
		if parts[0] == "f" {
			indices := make([]int, 0, len(parts)-1)
			for i := 1; i < len(parts); i++ {
				s := parts[i]
				if idx := strings.Index(s, "/"); idx >= 0 {
					s = s[:idx]
				}
				v, e := strconv.ParseInt(s, 10, 64)
				if e != nil {
					continue
				}
				idx := int(v) - 1
				if idx < 0 {
					idx += len(verts) + 1
				}
				indices = append(indices, idx)
			}
			for i := 2; i < len(indices); i++ {
				i0, i1, i2 := indices[0], indices[i-1], indices[i]
				if i0 >= 0 && i0 < len(verts) && i1 >= 0 && i1 < len(verts) && i2 >= 0 && i2 < len(verts) {
					tris = append(tris, triangle{verts[i0], verts[i1], verts[i2]})
				}
			}
		}
	}
	if err = sc.Err(); err != nil {
		return nil, vec3{}, vec3{}, err
	}
	if len(verts) == 0 {
		return nil, vec3{}, vec3{}, fmt.Errorf("obj has no vertices")
	}
	minV = verts[0]
	maxV = verts[0]
	for _, v := range verts {
		if v.x < minV.x {
			minV.x = v.x
		}
		if v.y < minV.y {
			minV.y = v.y
		}
		if v.z < minV.z {
			minV.z = v.z
		}
		if v.x > maxV.x {
			maxV.x = v.x
		}
		if v.y > maxV.y {
			maxV.y = v.y
		}
		if v.z > maxV.z {
			maxV.z = v.z
		}
	}
	return tris, minV, maxV, nil
}

type joint struct {
	kind         string  // "point_to_point" | "fixed" | "hinge" | "slider" | "cone_twist"
	bodyA        string
	bodyB        string
	anchorA      vec3    // in body A local space
	anchorB      vec3    // in body B local space
	axisA        vec3    // axis in body A local space (hinge/slider/conetwist)
	axisB        vec3    // axis in body B local space (hinge/slider/conetwist)
	limitMin     float64 // angle (rad) or position for hinge/slider; cone angle for conetwist
	limitMax     float64
	motorTarget  float64 // target velocity (rad/s or m/s)
	motorMaxForce float64
}

type world struct {
	gravity vec3
	bodies  map[string]*body
	joints  map[string]*joint
	mu      sync.RWMutex
}

const defaultPhysicsWorld = "default"

var (
	bulletNativeAvailable = false // set to true by bullet_native.go when built with -tags bullet
	worlds                = make(map[string]*world)
	worldMu        sync.RWMutex
	physicsBodySeq int
	lastRay        struct {
		hit    bool
		p      vec3
		bodyId string
		normal vec3
	}
	lastRayMu sync.Mutex
)

func bulletFeatureAvailable(name string) bool {
	switch strings.ToLower(strings.TrimSpace(name)) {
	case "sphere", "box", "capsule", "cylinder", "cone", "raycast", "body_properties", "ccd", "kinematic", "torque", "torque_impulse":
		return true
	case "point_to_point_joint", "fixed_joint", "hinge_joint", "slider_joint", "cone_twist_joint",
		"joint_limits", "joint_motor":
		return true
	case "mesh_collider":
		return true
	case "compound", "compound_shapes":
		return true
	case "native", "native_backend", "joints", "heightmap", "exact_mesh_collision":
		return false
	default:
		return false
	}
}

// bodyInertia returns approximate moment of inertia for torque. Sphere: (2/5)*m*r^2; box: m*(hx^2+hy^2+hz^2)/12. Default 1 if unknown.
func bodyInertia(b *body) float64 {
	if b == nil || b.mass <= 0 {
		return 1
	}
	if b.radius > 0 && (b.halfExt.x <= b.radius && b.halfExt.y <= b.radius && b.halfExt.z <= b.radius) {
		return b.mass * b.radius * b.radius * 0.4 // sphere: (2/5)*m*r^2
	}
	hx, hy, hz := b.halfExt.x*2, b.halfExt.y*2, b.halfExt.z*2
	return b.mass * (hx*hx + hy*hy + hz*hz) / 12
}

// eulerRotate applies Euler XYZ rotation to a local-space vector.
func eulerRotate(v vec3, euler vec3) vec3 {
	cx, sx := math.Cos(euler.x), math.Sin(euler.x)
	cy, sy := math.Cos(euler.y), math.Sin(euler.y)
	cz, sz := math.Cos(euler.z), math.Sin(euler.z)
	// R = Rz * Ry * Rx
	x := v.x*cy*cz + v.y*(cx*sz+sx*sy*cz) + v.z*(-sx*sz+cx*sy*cz)
	y := v.x*(-cy*sz) + v.y*(cx*cz-sx*sy*sz) + v.z*(sx*cz+cx*sy*sz)
	z := v.x*sy + v.y*(-sx*cy) + v.z*(cx*cy)
	return vec3{x, y, z}
}

func unsupportedBulletFeatureError(feature string) error {
	return fmt.Errorf("%s is not supported by the shipped Bullet fallback backend; check BulletFeatureAvailable() or BulletNativeAvailable()", feature)
}

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
		if w.joints == nil {
			w.joints = make(map[string]*joint)
		}
		return w
	}
	w := &world{
		gravity: vec3{gx, gy, gz},
		bodies:  make(map[string]*body),
		joints:  make(map[string]*joint),
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

// RegisterBullet registers 3D physics with the VM as flat names only (CreateWorld3D, Step3D, etc.).
// The compiler rewrites BULLET.* calls to these flat names for backward compatibility.
func RegisterBullet(v *vm.VM) {
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
	v.RegisterForeign("BulletBackendName", func(args []interface{}) (interface{}, error) {
		return "purego-fallback", nil
	})
	v.RegisterForeign("BulletBackendMode", func(args []interface{}) (interface{}, error) {
		return "fallback", nil
	})
	v.RegisterForeign("BulletNativeAvailable", func(args []interface{}) (interface{}, error) {
		if bulletNativeAvailable {
			return 1, nil
		}
		return 0, nil
	})
	v.RegisterForeign("BulletFeatureAvailable", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return 0, fmt.Errorf("BulletFeatureAvailable requires (featureName$)")
		}
		if bulletFeatureAvailable(toString(args[0])) {
			return 1, nil
		}
		return 0, nil
	})

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
	v.RegisterForeign("SetWorldGravity3D", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("SetWorldGravity3D requires (worldId, gx, gy, gz)")
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
	v.RegisterForeign("DestroyBody3D", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("DestroyBody3D requires (worldId, bodyId)")
		}
		w := getWorld(toString(args[0]))
		if w == nil {
			return nil, nil
		}
		bid := toString(args[1])
		w.mu.Lock()
		delete(w.bodies, bid)
		w.mu.Unlock()
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
			integrateBody(b, w.gravity, dt)
		}
		solveJoints(w, dt)
		for pass := 0; pass < 3; pass++ {
			resolveCollisions(w)
		}
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
				integrateBody(b, w.gravity, dt)
			}
			solveJoints(w, dt)
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
		height := toFloat64(args[6])
		if height < radius*2 {
			height = radius * 2
		}
		w.mu.Lock()
		w.bodies[bid] = &body{
			id:       bid,
			position: vec3{toFloat64(args[2]), toFloat64(args[3]), toFloat64(args[4])},
			halfExt:  vec3{radius, height / 2, radius},
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
		meshName := toString(args[2])
		w := getWorld(wid)
		if w == nil {
			w = getOrCreateWorld(wid, 0, -9.81, 0)
		}
		b := &body{
			id:       bid,
			position: vec3{0, 0, 0},
			halfExt:  vec3{1, 1, 1},
			mass:     0,
			active:   true,
			scale:    vec3{1, 1, 1},
		}
		if meshName != "" {
			tris, minV, maxV, err := loadOBJForCollision(meshName)
			if err == nil && len(tris) > 0 {
				b.meshTriangles = tris
				b.meshAABBMin = minV
				b.meshAABBMax = maxV
				b.halfExt = vec3{(maxV.x - minV.x) / 2, (maxV.y - minV.y) / 2, (maxV.z - minV.z) / 2}
			}
		}
		w.mu.Lock()
		w.bodies[bid] = b
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
		return nil, fmt.Errorf("CreateHeightmap3D: use CreateStaticMesh3D with a heightfield mesh or a large static CreateBox3D in the pure-Go fallback; heightfield sampling is not implemented here")
	})
	v.RegisterForeign("CreateCompound3D", func(args []interface{}) (interface{}, error) {
		if len(args) < 6 {
			return nil, fmt.Errorf("CreateCompound3D requires (world$, body$, x, y, z, mass)")
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
			halfExt:  vec3{0, 0, 0},
			mass:     toFloat64(args[5]),
			active:   true,
			scale:    vec3{1, 1, 1},
			compound: nil,
		}
		w.mu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("AddShapeToCompound3D", func(args []interface{}) (interface{}, error) {
		if len(args) < 8 {
			return nil, fmt.Errorf("AddShapeToCompound3D requires (world$, body$, ox, oy, oz, sizeX, sizeY, sizeZ)")
		}
		w := getWorld(toString(args[0]))
		if w == nil {
			return nil, fmt.Errorf("world not found")
		}
		bid := toString(args[1])
		sx, sy, sz := toFloat64(args[5]), toFloat64(args[6]), toFloat64(args[7])
		w.mu.Lock()
		b := w.bodies[bid]
		if b == nil {
			w.mu.Unlock()
			return nil, fmt.Errorf("body not found")
		}
		if len(b.meshTriangles) > 0 {
			w.mu.Unlock()
			return nil, fmt.Errorf("AddShapeToCompound3D: cannot add compound shape to mesh body")
		}
		b.compound = append(b.compound, struct{ ox, oy, oz, hx, hy, hz float64 }{
			ox: toFloat64(args[2]), oy: toFloat64(args[3]), oz: toFloat64(args[4]),
			hx: sx / 2, hy: sy / 2, hz: sz / 2,
		})
		w.mu.Unlock()
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
		if len(args) < 5 {
			return nil, fmt.Errorf("ApplyTorque3D requires (worldId, bodyId, tx, ty, tz)")
		}
		b := getBody(getWorld(toString(args[0])), toString(args[1]))
		if b == nil {
			return nil, nil
		}
		b.pendingTorque.x += toFloat64(args[2])
		b.pendingTorque.y += toFloat64(args[3])
		b.pendingTorque.z += toFloat64(args[4])
		return nil, nil
	})
	v.RegisterForeign("ApplyTorqueImpulse3D", func(args []interface{}) (interface{}, error) {
		if len(args) < 5 {
			return nil, fmt.Errorf("ApplyTorqueImpulse3D requires (worldId, bodyId, ix, iy, iz)")
		}
		b := getBody(getWorld(toString(args[0])), toString(args[1]))
		if b == nil {
			return nil, nil
		}
		I := bodyInertia(b)
		b.angularVelocity.x += toFloat64(args[2]) / I
		b.angularVelocity.y += toFloat64(args[3]) / I
		b.angularVelocity.z += toFloat64(args[4]) / I
		return nil, nil
	})

	// Body properties (implemented; used in Step and resolveCollisions)
	v.RegisterForeign("SetFriction3D", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("SetFriction3D requires (worldId, bodyId, friction)")
		}
		b := getBody(getWorld(toString(args[0])), toString(args[1]))
		if b == nil {
			return nil, nil
		}
		b.friction = toFloat64(args[2])
		return nil, nil
	})
	v.RegisterForeign("SetRestitution3D", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("SetRestitution3D requires (worldId, bodyId, restitution)")
		}
		b := getBody(getWorld(toString(args[0])), toString(args[1]))
		if b == nil {
			return nil, nil
		}
		b.restitution = toFloat64(args[2])
		return nil, nil
	})
	v.RegisterForeign("SetDamping3D", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("SetDamping3D requires (worldId, bodyId, linearDamp, angularDamp)")
		}
		b := getBody(getWorld(toString(args[0])), toString(args[1]))
		if b == nil {
			return nil, nil
		}
		b.linearDamping = toFloat64(args[2])
		b.angularDamping = toFloat64(args[3])
		return nil, nil
	})
	v.RegisterForeign("SetKinematic3D", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("SetKinematic3D requires (worldId, bodyId, kinematic)")
		}
		b := getBody(getWorld(toString(args[0])), toString(args[1]))
		if b == nil {
			return nil, nil
		}
		b.kinematic = toFloat64(args[2]) != 0
		return nil, nil
	})
	v.RegisterForeign("SetGravity3D", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("SetGravity3D requires (worldId, bodyId, gravityScale)")
		}
		b := getBody(getWorld(toString(args[0])), toString(args[1]))
		if b == nil {
			return nil, nil
		}
		b.gravityScale = toFloat64(args[2])
		return nil, nil
	})
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
	v.RegisterForeign("GetMass3D", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return 0.0, nil
		}
		b := getBody(getWorld(toString(args[0])), toString(args[1]))
		if b == nil {
			return 0.0, nil
		}
		return b.mass, nil
	})
	v.RegisterForeign("SetLinearFactor3D", func(args []interface{}) (interface{}, error) {
		if len(args) < 5 {
			return nil, fmt.Errorf("SetLinearFactor3D requires (worldId, bodyId, fx, fy, fz)")
		}
		b := getBody(getWorld(toString(args[0])), toString(args[1]))
		if b == nil {
			return nil, nil
		}
		b.linearFactor = vec3{toFloat64(args[2]), toFloat64(args[3]), toFloat64(args[4])}
		return nil, nil
	})
	v.RegisterForeign("SetAngularFactor3D", func(args []interface{}) (interface{}, error) {
		if len(args) < 5 {
			return nil, fmt.Errorf("SetAngularFactor3D requires (worldId, bodyId, ax, ay, az)")
		}
		b := getBody(getWorld(toString(args[0])), toString(args[1]))
		if b == nil {
			return nil, nil
		}
		b.angularFactor = vec3{toFloat64(args[2]), toFloat64(args[3]), toFloat64(args[4])}
		return nil, nil
	})
	v.RegisterForeign("SetCCD3D", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("SetCCD3D requires (worldId, bodyId, enable)")
		}
		b := getBody(getWorld(toString(args[0])), toString(args[1]))
		if b == nil {
			return nil, nil
		}
		b.ccd = toFloat64(args[2]) != 0
		return nil, nil
	})

	// BulletJointsAvailable: 0 = pure-Go (no joints), 1 = joints supported (PointToPoint, Fixed)
	v.RegisterForeign("BulletJointsAvailable", func(args []interface{}) (interface{}, error) {
		return 1, nil
	})
	// 3D joints: Hinge, Slider, ConeTwist implemented in pure-Go fallback.
	v.RegisterForeign("CreateHingeJoint3D", func(args []interface{}) (interface{}, error) {
		if len(args) < 13 {
			return nil, fmt.Errorf("CreateHingeJoint3D requires (worldId, jointId, bodyA, bodyB, ax, ay, az, bx, by, bz, axisAx, axisAy, axisAz)")
		}
		w := getWorld(toString(args[0]))
		if w == nil {
			return nil, fmt.Errorf("world not found")
		}
		jid := toString(args[1])
		ba, bb := toString(args[2]), toString(args[3])
		if getBody(w, ba) == nil || getBody(w, bb) == nil {
			return nil, fmt.Errorf("body not found")
		}
		w.mu.Lock()
		if w.joints == nil {
			w.joints = make(map[string]*joint)
		}
		w.joints[jid] = &joint{
			kind:    "hinge",
			bodyA:   ba,
			bodyB:   bb,
			anchorA: vec3{toFloat64(args[4]), toFloat64(args[5]), toFloat64(args[6])},
			anchorB: vec3{toFloat64(args[7]), toFloat64(args[8]), toFloat64(args[9])},
			axisA:   vec3{toFloat64(args[10]), toFloat64(args[11]), toFloat64(args[12])},
			axisB:   vec3{toFloat64(args[10]), toFloat64(args[11]), toFloat64(args[12])},
			limitMin: -math.Pi, limitMax: math.Pi,
		}
		w.mu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("CreateSliderJoint3D", func(args []interface{}) (interface{}, error) {
		if len(args) < 13 {
			return nil, fmt.Errorf("CreateSliderJoint3D requires (worldId, jointId, bodyA, bodyB, ax, ay, az, bx, by, bz, axisAx, axisAy, axisAz)")
		}
		w := getWorld(toString(args[0]))
		if w == nil {
			return nil, fmt.Errorf("world not found")
		}
		jid := toString(args[1])
		ba, bb := toString(args[2]), toString(args[3])
		if getBody(w, ba) == nil || getBody(w, bb) == nil {
			return nil, fmt.Errorf("body not found")
		}
		w.mu.Lock()
		if w.joints == nil {
			w.joints = make(map[string]*joint)
		}
		w.joints[jid] = &joint{
			kind:     "slider",
			bodyA:    ba,
			bodyB:    bb,
			anchorA:  vec3{toFloat64(args[4]), toFloat64(args[5]), toFloat64(args[6])},
			anchorB:  vec3{toFloat64(args[7]), toFloat64(args[8]), toFloat64(args[9])},
			axisA:    vec3{toFloat64(args[10]), toFloat64(args[11]), toFloat64(args[12])},
			axisB:    vec3{toFloat64(args[10]), toFloat64(args[11]), toFloat64(args[12])},
			limitMin: -1e30, limitMax: 1e30,
		}
		w.mu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("CreateConeTwistJoint3D", func(args []interface{}) (interface{}, error) {
		if len(args) < 13 {
			return nil, fmt.Errorf("CreateConeTwistJoint3D requires (worldId, jointId, bodyA, bodyB, ax, ay, az, bx, by, bz, axisAx, axisAy, axisAz)")
		}
		w := getWorld(toString(args[0]))
		if w == nil {
			return nil, fmt.Errorf("world not found")
		}
		jid := toString(args[1])
		ba, bb := toString(args[2]), toString(args[3])
		if getBody(w, ba) == nil || getBody(w, bb) == nil {
			return nil, fmt.Errorf("body not found")
		}
		w.mu.Lock()
		if w.joints == nil {
			w.joints = make(map[string]*joint)
		}
		w.joints[jid] = &joint{
			kind:     "cone_twist",
			bodyA:    ba,
			bodyB:    bb,
			anchorA:  vec3{toFloat64(args[4]), toFloat64(args[5]), toFloat64(args[6])},
			anchorB:  vec3{toFloat64(args[7]), toFloat64(args[8]), toFloat64(args[9])},
			axisA:    vec3{toFloat64(args[10]), toFloat64(args[11]), toFloat64(args[12])},
			axisB:    vec3{toFloat64(args[10]), toFloat64(args[11]), toFloat64(args[12])},
			limitMin: 0, limitMax: math.Pi / 2,
		}
		w.mu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("CreatePointToPointJoint3D", func(args []interface{}) (interface{}, error) {
		if len(args) < 10 {
			return nil, fmt.Errorf("CreatePointToPointJoint3D requires (worldId, jointId, bodyA, bodyB, ax, ay, az, bx, by, bz)")
		}
		w := getWorld(toString(args[0]))
		if w == nil {
			return nil, fmt.Errorf("world not found")
		}
		jid := toString(args[1])
		ba, bb := toString(args[2]), toString(args[3])
		if getBody(w, ba) == nil || getBody(w, bb) == nil {
			return nil, fmt.Errorf("body not found")
		}
		w.mu.Lock()
		if w.joints == nil {
			w.joints = make(map[string]*joint)
		}
		w.joints[jid] = &joint{
			kind:    "point_to_point",
			bodyA:   ba,
			bodyB:   bb,
			anchorA: vec3{toFloat64(args[4]), toFloat64(args[5]), toFloat64(args[6])},
			anchorB: vec3{toFloat64(args[7]), toFloat64(args[8]), toFloat64(args[9])},
		}
		w.mu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("CreateFixedJoint3D", func(args []interface{}) (interface{}, error) {
		if len(args) < 5 {
			return nil, fmt.Errorf("CreateFixedJoint3D requires (worldId, jointId, bodyA, bodyB)")
		}
		w := getWorld(toString(args[0]))
		if w == nil {
			return nil, fmt.Errorf("world not found")
		}
		jid := toString(args[1])
		ba, bb := toString(args[2]), toString(args[3])
		if getBody(w, ba) == nil || getBody(w, bb) == nil {
			return nil, fmt.Errorf("body not found")
		}
		w.mu.Lock()
		if w.joints == nil {
			w.joints = make(map[string]*joint)
		}
		w.joints[jid] = &joint{
			kind:    "fixed",
			bodyA:   ba,
			bodyB:   bb,
			anchorA: vec3{},
			anchorB: vec3{},
		}
		w.mu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("SetJointLimits3D", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("SetJointLimits3D requires (worldId, jointId, low, high)")
		}
		w := getWorld(toString(args[0]))
		if w == nil || w.joints == nil {
			return nil, nil
		}
		j := w.joints[toString(args[1])]
		if j == nil {
			return nil, nil
		}
		w.mu.Lock()
		j.limitMin = toFloat64(args[2])
		j.limitMax = toFloat64(args[3])
		w.mu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("SetJointMotor3D", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("SetJointMotor3D requires (worldId, jointId, targetVel, maxForce)")
		}
		w := getWorld(toString(args[0]))
		if w == nil || w.joints == nil {
			return nil, nil
		}
		j := w.joints[toString(args[1])]
		if j == nil {
			return nil, nil
		}
		w.mu.Lock()
		j.motorTarget = toFloat64(args[2])
		j.motorMaxForce = toFloat64(args[3])
		w.mu.Unlock()
		return nil, nil
	})

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
			min, max := bodyAABB(b)
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
	// RayCastFromDir3D(worldId, startX, startY, startZ, dirX, dirY, dirZ, maxDist [, excludeBody$]):
	// optional 9th arg skips that body id (e.g. "player") so foot rays don't hit the character.
	v.RegisterForeign("RayCastFromDir3D", func(args []interface{}) (interface{}, error) {
		if len(args) < 8 {
			return nil, fmt.Errorf("RayCastFromDir3D requires (worldId, startX, startY, startZ, dirX, dirY, dirZ, maxDist [, excludeBody$])")
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
		exclude := ""
		if len(args) >= 9 {
			exclude = strings.ToLower(strings.TrimSpace(toString(args[8])))
		}
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
			if exclude != "" && strings.ToLower(id) == exclude {
				continue
			}
			min, max := bodyAABB(b)
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
	v.RegisterForeign("SetBodyPosition", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("SetBodyPosition requires (bodyId, x, y, z)")
		}
		w := getWorld(defaultPhysicsWorld)
		if w == nil {
			return nil, nil
		}
		w.mu.Lock()
		b := w.bodies[toString(args[0])]
		if b != nil {
			b.position = vec3{toFloat64(args[1]), toFloat64(args[2]), toFloat64(args[3])}
		}
		w.mu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("GetBodyPosition", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return []interface{}{0.0, 0.0, 0.0}, nil
		}
		w := getWorld(defaultPhysicsWorld)
		if w == nil {
			return []interface{}{0.0, 0.0, 0.0}, nil
		}
		w.mu.RLock()
		b := w.bodies[toString(args[0])]
		x, y, z := 0.0, 0.0, 0.0
		if b != nil {
			x, y, z = b.position.x, b.position.y, b.position.z
		}
		w.mu.RUnlock()
		return []interface{}{x, y, z}, nil
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

// integrateBody applies gravity (if not kinematic), integrates position, and applies damping. Call with world mutex held.
func integrateBody(b *body, gravity vec3, dt float64) {
	if b.kinematic {
		return
	}
	gs := b.gravityScale
	if gs == 0 {
		gs = 1
	}
	b.velocity.x += gravity.x * gs * dt
	b.velocity.y += gravity.y * gs * dt
	b.velocity.z += gravity.z * gs * dt
	lf := b.linearFactor
	if lf.x == 0 && lf.y == 0 && lf.z == 0 {
		lf = vec3{1, 1, 1}
	}
	b.position.x += b.velocity.x * lf.x * dt
	b.position.y += b.velocity.y * lf.y * dt
	b.position.z += b.velocity.z * lf.z * dt
	if b.linearDamping > 0 {
		damp := 1.0 - b.linearDamping*dt
		if damp < 0 {
			damp = 0
		}
		b.velocity.x *= damp
		b.velocity.y *= damp
		b.velocity.z *= damp
	}
	// Apply pending torque and integrate angular velocity into rotation
	I := bodyInertia(b)
	af := b.angularFactor
	if af.x == 0 && af.y == 0 && af.z == 0 {
		af = vec3{1, 1, 1}
	}
	b.angularVelocity.x += b.pendingTorque.x / I * dt
	b.angularVelocity.y += b.pendingTorque.y / I * dt
	b.angularVelocity.z += b.pendingTorque.z / I * dt
	b.pendingTorque = vec3{}
	b.rotation.x += b.angularVelocity.x * af.x * dt
	b.rotation.y += b.angularVelocity.y * af.y * dt
	b.rotation.z += b.angularVelocity.z * af.z * dt
	if b.angularDamping > 0 {
		damp := 1.0 - b.angularDamping*dt
		if damp < 0 {
			damp = 0
		}
		b.angularVelocity.x *= damp
		b.angularVelocity.y *= damp
		b.angularVelocity.z *= damp
	}
}

// vec3Norm normalizes v and returns length. If length < 1e-9, returns 0.
func vec3Norm(v *vec3) float64 {
	len := math.Sqrt(v.x*v.x + v.y*v.y + v.z*v.z)
	if len < 1e-9 {
		return 0
	}
	v.x /= len
	v.y /= len
	v.z /= len
	return len
}

// vec3Dot returns dot product.
func vec3Dot(a, b vec3) float64 {
	return a.x*b.x + a.y*b.y + a.z*b.z
}

// vec3Cross returns cross product a x b.
func vec3Cross(a, b vec3) vec3 {
	return vec3{
		a.y*b.z - a.z*b.y,
		a.z*b.x - a.x*b.z,
		a.x*b.y - a.y*b.x,
	}
}

// solveJoints runs position-based constraint correction for PointToPoint, Fixed, Hinge, Slider, ConeTwist. Call with w.mu held.
func solveJoints(w *world, dt float64) {
	if w == nil || w.joints == nil {
		return
	}
	const iterations = 6
	const correctionFactor = 0.4
	for iter := 0; iter < iterations; iter++ {
		for _, j := range w.joints {
			a, b := w.bodies[j.bodyA], w.bodies[j.bodyB]
			if a == nil || b == nil || !a.active || !b.active {
				continue
			}
			worldA := vec3{
				a.position.x + eulerRotate(j.anchorA, a.rotation).x,
				a.position.y + eulerRotate(j.anchorA, a.rotation).y,
				a.position.z + eulerRotate(j.anchorA, a.rotation).z,
			}
			worldB := vec3{
				b.position.x + eulerRotate(j.anchorB, b.rotation).x,
				b.position.y + eulerRotate(j.anchorB, b.rotation).y,
				b.position.z + eulerRotate(j.anchorB, b.rotation).z,
			}
			ma, mb := a.mass, b.mass
			if ma <= 0 {
				ma = 1
			}
			if mb <= 0 {
				mb = 1
			}
			total := ma + mb
			wa, wb := mb/total, ma/total
			if a.kinematic {
				wa, wb = 0, 1
			}
			if b.kinematic {
				wa, wb = 1, 0
			}
			if a.kinematic && b.kinematic {
				continue
			}
			switch j.kind {
			case "point_to_point", "hinge", "cone_twist":
				errX := worldA.x - worldB.x
				errY := worldA.y - worldB.y
				errZ := worldA.z - worldB.z
				dist := math.Sqrt(errX*errX + errY*errY + errZ*errZ)
				if dist >= 1e-6 {
					invDist := 1.0 / dist
					nx, ny, nz := errX*invDist, errY*invDist, errZ*invDist
					corr := dist * correctionFactor
					if !a.kinematic {
						a.position.x -= nx * corr * wa
						a.position.y -= ny * corr * wa
						a.position.z -= nz * corr * wa
					}
					if !b.kinematic {
						b.position.x += nx * corr * wb
						b.position.y += ny * corr * wb
						b.position.z += nz * corr * wb
					}
				}
			case "slider":
				axisA := eulerRotate(j.axisA, a.rotation)
				if vec3Norm(&axisA) < 1e-9 {
					axisA = vec3{1, 0, 0}
				}
				diff := vec3{worldB.x - worldA.x, worldB.y - worldA.y, worldB.z - worldA.z}
				slideDist := vec3Dot(diff, axisA)
				clamped := slideDist
				if clamped < j.limitMin {
					clamped = j.limitMin
				}
				if clamped > j.limitMax {
					clamped = j.limitMax
				}
				perp := vec3{
					diff.x - axisA.x*slideDist,
					diff.y - axisA.y*slideDist,
					diff.z - axisA.z*slideDist,
				}
				perpLen := math.Sqrt(perp.x*perp.x + perp.y*perp.y + perp.z*perp.z)
				if perpLen >= 1e-6 {
					invLen := 1.0 / perpLen
					nx, ny, nz := perp.x*invLen, perp.y*invLen, perp.z*invLen
					corrPerp := perpLen * correctionFactor
					if !a.kinematic {
						a.position.x += nx * corrPerp * wa
						a.position.y += ny * corrPerp * wa
						a.position.z += nz * corrPerp * wa
					}
					if !b.kinematic {
						b.position.x -= nx * corrPerp * wb
						b.position.y -= ny * corrPerp * wb
						b.position.z -= nz * corrPerp * wb
					}
				}
				corr := (clamped - slideDist) * correctionFactor
				if !a.kinematic {
					a.position.x += axisA.x * corr * wa
					a.position.y += axisA.y * corr * wa
					a.position.z += axisA.z * corr * wa
				}
				if !b.kinematic {
					b.position.x -= axisA.x * corr * wb
					b.position.y -= axisA.y * corr * wb
					b.position.z -= axisA.z * corr * wb
				}
				if j.motorMaxForce > 0 {
					velErr := j.motorTarget - (b.velocity.x*axisA.x + b.velocity.y*axisA.y + b.velocity.z*axisA.z)
					impulse := velErr * 0.1
					if impulse > j.motorMaxForce*dt {
						impulse = j.motorMaxForce * dt
					}
					if impulse < -j.motorMaxForce*dt {
						impulse = -j.motorMaxForce * dt
					}
					if !b.kinematic {
						b.velocity.x += axisA.x * impulse / mb
						b.velocity.y += axisA.y * impulse / mb
						b.velocity.z += axisA.z * impulse / mb
					}
				}
			case "fixed":
				errX := worldA.x - worldB.x
				errY := worldA.y - worldB.y
				errZ := worldA.z - worldB.z
				dist := math.Sqrt(errX*errX + errY*errY + errZ*errZ)
				if dist >= 1e-6 {
					invDist := 1.0 / dist
					nx, ny, nz := errX*invDist, errY*invDist, errZ*invDist
					corr := dist * correctionFactor
					if !a.kinematic {
						a.position.x -= nx * corr * wa
						a.position.y -= ny * corr * wa
						a.position.z -= nz * corr * wa
					}
					if !b.kinematic {
						b.position.x += nx * corr * wb
						b.position.y += ny * corr * wb
						b.position.z += nz * corr * wb
					}
				}
				b.rotation = a.rotation
				b.angularVelocity = a.angularVelocity
			}
			if j.kind == "hinge" || j.kind == "cone_twist" {
				axisA := eulerRotate(j.axisA, a.rotation)
				axisB := eulerRotate(j.axisB, b.rotation)
				if vec3Norm(&axisA) < 1e-9 {
					axisA = vec3{1, 0, 0}
				}
				if vec3Norm(&axisB) < 1e-9 {
					axisB = axisA
				}
				cross := vec3Cross(axisA, axisB)
				sinAngle := math.Sqrt(cross.x*cross.x + cross.y*cross.y + cross.z*cross.z)
				if sinAngle >= 1e-6 {
					dot := vec3Dot(axisA, axisB)
					angle := math.Atan2(sinAngle, math.Max(-1+1e-9, math.Min(1-1e-9, dot)))
					targetAngle := angle
					if j.kind == "cone_twist" {
						if angle > j.limitMax {
							targetAngle = j.limitMax
						}
					} else if j.kind == "hinge" {
						if angle < j.limitMin {
							targetAngle = j.limitMin
						} else if angle > j.limitMax {
							targetAngle = j.limitMax
						}
					}
					corrAngle := (targetAngle - angle) * correctionFactor * 0.5
					invLen := 1.0 / (math.Sqrt(cross.x*cross.x+cross.y*cross.y+cross.z*cross.z) + 1e-9)
					cross.x *= invLen
					cross.y *= invLen
					cross.z *= invLen
					if !b.kinematic {
						ia := bodyInertia(b)
						b.angularVelocity.x += cross.x * corrAngle / (dt*ia + 1e-9)
						b.angularVelocity.y += cross.y * corrAngle / (dt*ia + 1e-9)
						b.angularVelocity.z += cross.z * corrAngle / (dt*ia + 1e-9)
					}
				}
				if j.motorMaxForce > 0 && j.kind == "hinge" {
					axisAMotor := eulerRotate(j.axisA, a.rotation)
					if vec3Norm(&axisAMotor) < 1e-9 {
						axisAMotor = vec3{1, 0, 0}
					}
					angVelAlongAxis := b.angularVelocity.x*axisAMotor.x + b.angularVelocity.y*axisAMotor.y + b.angularVelocity.z*axisAMotor.z
					velErr := j.motorTarget - angVelAlongAxis
					impulse := velErr * 0.1
					if impulse > j.motorMaxForce*dt {
						impulse = j.motorMaxForce * dt
					}
					if impulse < -j.motorMaxForce*dt {
						impulse = -j.motorMaxForce * dt
					}
					ia := bodyInertia(b)
					if !b.kinematic && ia > 0 {
						b.angularVelocity.x += axisAMotor.x * impulse / ia
						b.angularVelocity.y += axisAMotor.y * impulse / ia
						b.angularVelocity.z += axisAMotor.z * impulse / ia
					}
				}
			}
		}
	}
}

func bodiesAreJoined(w *world, aId, bId string) bool {
	if w == nil || w.joints == nil {
		return false
	}
	for _, j := range w.joints {
		if (j.bodyA == aId && j.bodyB == bId) || (j.bodyA == bId && j.bodyB == aId) {
			return true
		}
	}
	return false
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
			if bodiesAreJoined(w, a.id, b.id) {
				continue
			}
			nx, ny, nz, depth := overlapBodies(a, b)
			if depth <= 0 {
				continue
			}
			a.collisions = append(a.collisions, collisionHit{b.id, vec3{nx, ny, nz}})
			b.collisions = append(b.collisions, collisionHit{a.id, vec3{-nx, -ny, -nz}})
			// Only push dynamic bodies; do not move kinematic bodies
			if !a.kinematic {
				a.position.x += nx * depth
				a.position.y += ny * depth
				a.position.z += nz * depth
				vn := a.velocity.x*nx + a.velocity.y*ny + a.velocity.z*nz
				// restitution: reflect normal velocity (use average of both bodies)
				rest := a.restitution
				if b.restitution > rest {
					rest = b.restitution
				}
				if vn < 0 {
					a.velocity.x -= (1 + rest) * vn * nx
					a.velocity.y -= (1 + rest) * vn * ny
					a.velocity.z -= (1 + rest) * vn * nz
				}
				// friction: reduce tangential velocity
				fric := a.friction
				if b.friction > fric {
					fric = b.friction
				}
				if fric > 0 {
					vx, vy, vz := a.velocity.x, a.velocity.y, a.velocity.z
					vnVal := vx*nx + vy*ny + vz*nz
					tx := vx - vnVal*nx
					ty := vy - vnVal*ny
					tz := vz - vnVal*nz
					scale := 1.0 - fric
					if scale < 0 {
						scale = 0
					}
					a.velocity.x = vnVal*nx + tx*scale
					a.velocity.y = vnVal*ny + ty*scale
					a.velocity.z = vnVal*nz + tz*scale
				}
			}
			if !b.kinematic {
				b.position.x -= nx * depth
				b.position.y -= ny * depth
				b.position.z -= nz * depth
				vn := b.velocity.x*(-nx) + b.velocity.y*(-ny) + b.velocity.z*(-nz)
				rest := b.restitution
				if a.restitution > rest {
					rest = a.restitution
				}
				if vn < 0 {
					b.velocity.x -= (1 + rest) * vn * (-nx)
					b.velocity.y -= (1 + rest) * vn * (-ny)
					b.velocity.z -= (1 + rest) * vn * (-nz)
				}
				fric := b.friction
				if a.friction > fric {
					fric = a.friction
				}
				if fric > 0 {
					vx, vy, vz := b.velocity.x, b.velocity.y, b.velocity.z
					nnx, nny, nnz := -nx, -ny, -nz
					vnVal := vx*nnx + vy*nny + vz*nnz
					tx := vx - vnVal*nnx
					ty := vy - vnVal*nny
					tz := vz - vnVal*nnz
					scale := 1.0 - fric
					if scale < 0 {
						scale = 0
					}
					b.velocity.x = vnVal*nnx + tx*scale
					b.velocity.y = vnVal*nny + ty*scale
					b.velocity.z = vnVal*nnz + tz*scale
				}
			}
		}
	}
}

// overlapBoxWithBody tests an axis-aligned box vs body b (same rules as overlapBodies for b).
func overlapBoxWithBody(ax, ay, az, ahx, ahy, ahz float64, b *body) (nx, ny, nz, depth float64) {
	if b == nil {
		return 0, 0, 0, 0
	}
	if bodyUsesCapsuleBounds(b) {
		br := b.radius
		b.radius = 0
		nx, ny, nz, depth = overlapBoxWithBody(ax, ay, az, ahx, ahy, ahz, b)
		b.radius = br
		return nx, ny, nz, depth
	}
	if bodyIsCompound(b) {
		best := 0.0
		for _, p := range b.compound {
			bx := b.position.x + p.ox
			by := b.position.y + p.oy
			bz := b.position.z + p.oz
			nnx, nny, nnz, dd := overlapBoxBox(ax, ay, az, ahx, ahy, ahz, bx, by, bz, p.hx, p.hy, p.hz)
			if dd > best {
				best = dd
				nx, ny, nz = nnx, nny, nnz
			}
		}
		return nx, ny, nz, best
	}
	bx, by, bz := b.position.x, b.position.y, b.position.z
	bhx, bhy, bhz := b.halfExt.x, b.halfExt.y, b.halfExt.z
	if bodyIsMesh(b) {
		bx, by, bz, bhx, bhy, bhz = bodyBoxParams(b)
	}
	bIsSphere := !bodyIsMesh(b) && b.radius > 0
	if bIsSphere {
		return overlapSphereBox(b.position.x, b.position.y, b.position.z, b.radius, ax, ay, az, ahx, ahy, ahz)
	}
	return overlapBoxBox(ax, ay, az, ahx, ahy, ahz, bx, by, bz, bhx, bhy, bhz)
}

func overlapCompoundWithBody(a, b *body) (nx, ny, nz, depth float64) {
	best := 0.0
	for _, p := range a.compound {
		ax := a.position.x + p.ox
		ay := a.position.y + p.oy
		az := a.position.z + p.oz
		nnx, nny, nnz, dd := overlapBoxWithBody(ax, ay, az, p.hx, p.hy, p.hz, b)
		if dd > best {
			best = dd
			nx, ny, nz = nnx, nny, nnz
		}
	}
	return nx, ny, nz, best
}

// overlapBodies returns contact normal (from a toward b) and penetration depth. Depth > 0 means overlap. Caller holds world lock.
func overlapBodies(a, b *body) (nx, ny, nz, depth float64) {
	if bodyUsesCapsuleBounds(a) {
		aSphere := a.radius
		a.radius = 0
		nx, ny, nz, depth = overlapBodies(a, b)
		a.radius = aSphere
		return nx, ny, nz, depth
	}
	if bodyUsesCapsuleBounds(b) {
		bSphere := b.radius
		b.radius = 0
		nx, ny, nz, depth = overlapBodies(a, b)
		b.radius = bSphere
		return nx, ny, nz, depth
	}
	if bodyIsCompound(a) && bodyIsCompound(b) {
		best := 0.0
		for _, pa := range a.compound {
			ax := a.position.x + pa.ox
			ay := a.position.y + pa.oy
			az := a.position.z + pa.oz
			for _, pb := range b.compound {
				bx := b.position.x + pb.ox
				by := b.position.y + pb.oy
				bz := b.position.z + pb.oz
				nnx, nny, nnz, dd := overlapBoxBox(ax, ay, az, pa.hx, pa.hy, pa.hz, bx, by, bz, pb.hx, pb.hy, pb.hz)
				if dd > best {
					best = dd
					nx, ny, nz = nnx, nny, nnz
				}
			}
		}
		return nx, ny, nz, best
	}
	if bodyIsCompound(a) {
		return overlapCompoundWithBody(a, b)
	}
	if bodyIsCompound(b) {
		nx, ny, nz, depth = overlapCompoundWithBody(b, a)
		return -nx, -ny, -nz, depth
	}
	// Resolve mesh to box params for overlap
	ax, ay, az := a.position.x, a.position.y, a.position.z
	ahx, ahy, ahz := a.halfExt.x, a.halfExt.y, a.halfExt.z
	if bodyIsMesh(a) {
		ax, ay, az, ahx, ahy, ahz = bodyBoxParams(a)
	}
	bx, by, bz := b.position.x, b.position.y, b.position.z
	bhx, bhy, bhz := b.halfExt.x, b.halfExt.y, b.halfExt.z
	if bodyIsMesh(b) {
		bx, by, bz, bhx, bhy, bhz = bodyBoxParams(b)
	}
	// Sphere: radius>0 and not mesh (mesh uses box). Capsule already handled.
	aIsSphere := !bodyIsMesh(a) && a.radius > 0
	bIsSphere := !bodyIsMesh(b) && b.radius > 0
	if aIsSphere && bIsSphere {
		return overlapSphereSphere(a.position.x, a.position.y, a.position.z, a.radius, b.position.x, b.position.y, b.position.z, b.radius)
	}
	if aIsSphere && !bIsSphere {
		return overlapSphereBox(a.position.x, a.position.y, a.position.z, a.radius, bx, by, bz, bhx, bhy, bhz)
	}
	if !aIsSphere && bIsSphere {
		nx, ny, nz, depth := overlapSphereBox(b.position.x, b.position.y, b.position.z, b.radius, ax, ay, az, ahx, ahy, ahz)
		return -nx, -ny, -nz, depth
	}
	return overlapBoxBox(ax, ay, az, ahx, ahy, ahz, bx, by, bz, bhx, bhy, bhz)
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
	if distSq >= sr*sr {
		return 0, 0, 0, 0
	}
	if distSq < 1e-18 {
		// Sphere center inside box: push out along the axis of minimum penetration (depth = distance to push sphere center so surface touches face)
		penR := (bx + hx) - sx - sr
		penL := sx - (bx - hx) - sr
		penU := (by + hy) - sy - sr
		penD := sy - (by - hy) - sr
		penF := (bz + hz) - sz - sr
		penB := sz - (bz - hz) - sr
		best := penR
		nx, ny, nz = 1, 0, 0
		if penL < best {
			best, nx, ny, nz = penL, -1, 0, 0
		}
		if penU < best {
			best, nx, ny, nz = penU, 0, 1, 0
		}
		if penD < best {
			best, nx, ny, nz = penD, 0, -1, 0
		}
		if penF < best {
			best, nx, ny, nz = penF, 0, 0, 1
		}
		if penB < best {
			best, nx, ny, nz = penB, 0, 0, -1
		}
		return nx, ny, nz, best
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

// rayAABB returns the smallest t >= 0 where ray (origin + t*dir) hits the AABB, or -1 if no hit.
// Handles axis-aligned rays (zero dir components) and rays that start inside the box (t=0).
func rayAABB(ox, oy, oz, dx, dy, dz, minX, minY, minZ, maxX, maxY, maxZ float64) float64 {
	const eps = 1e-9
	t0, t1 := 0.0, math.Inf(1)

	slab := func(o, d, mn, mx float64) bool {
		if math.Abs(d) < eps {
			if o < mn-eps || o > mx+eps {
				return false
			}
			return true
		}
		tNear := (mn - o) / d
		tFar := (mx - o) / d
		if tNear > tFar {
			tNear, tFar = tFar, tNear
		}
		if tNear > t0 {
			t0 = tNear
		}
		if tFar < t1 {
			t1 = tFar
		}
		return t0 <= t1
	}

	if !slab(ox, dx, minX, maxX) {
		return -1
	}
	if !slab(oy, dy, minY, maxY) {
		return -1
	}
	if !slab(oz, dz, minZ, maxZ) {
		return -1
	}
	if t0 < 0 {
		t0 = 0
	}
	if t0 > t1 {
		return -1
	}
	return t0
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
		integrateBody(b, w.gravity, timeStep)
	}
	solveJoints(w, timeStep)
	resolveCollisions(w)
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

// ApplyImpulse applies an impulse to the body.
func ApplyImpulse(worldId, bodyId string, ix, iy, iz float64) {
	b := getBody(getWorld(worldId), bodyId)
	if b == nil || b.mass <= 0 {
		return
	}
	m := b.mass
	b.velocity.x += ix / m
	b.velocity.y += iy / m
	b.velocity.z += iz / m
}

// GetRotationX/Y/Z return the body's Euler rotation (radians).
func GetRotationX(worldId, bodyId string) float64 {
	if b := getBody(getWorld(worldId), bodyId); b != nil {
		return b.rotation.x
	}
	return 0
}
func GetRotationY(worldId, bodyId string) float64 {
	if b := getBody(getWorld(worldId), bodyId); b != nil {
		return b.rotation.y
	}
	return 0
}
func GetRotationZ(worldId, bodyId string) float64 {
	if b := getBody(getWorld(worldId), bodyId); b != nil {
		return b.rotation.z
	}
	return 0
}

// SetRotation sets the body's Euler rotation (radians).
func SetRotation(worldId, bodyId string, rx, ry, rz float64) {
	if b := getBody(getWorld(worldId), bodyId); b != nil {
		b.rotation = vec3{rx, ry, rz}
	}
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
	var hitBodyId string
	var hitNorm vec3
	for id, b := range w.bodies {
		if !b.active {
			continue
		}
		min, max := bodyAABB(b)
		t := rayAABB(startX, startY, startZ, dx, dy, dz, min.x, min.y, min.z, max.x, max.y, max.z)
		if t >= 0 && t < maxDist && t < bestT {
			bestT = t
			hit = true
			hitP = vec3{startX + dx*t, startY + dy*t, startZ + dz*t}
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

// GetRayCastHitBody returns the body ID of the last raycast hit.
func GetRayCastHitBody() string {
	lastRayMu.Lock()
	defer lastRayMu.Unlock()
	return lastRay.bodyId
}

// GetRayCastHitNormalX/Y/Z return the normal of the last raycast hit.
func GetRayCastHitNormalX() float64 {
	lastRayMu.Lock()
	defer lastRayMu.Unlock()
	return lastRay.normal.x
}
func GetRayCastHitNormalY() float64 {
	lastRayMu.Lock()
	defer lastRayMu.Unlock()
	return lastRay.normal.y
}
func GetRayCastHitNormalZ() float64 {
	lastRayMu.Lock()
	defer lastRayMu.Unlock()
	return lastRay.normal.z
}

// Flat-name exports for --gen-go (no namespace in generated code).
func CreateWorld3D(worldId string, gx, gy, gz float64)     { CreateWorld(worldId, gx, gy, gz) }
func SetWorldGravity3D(worldId string, gx, gy, gz float64) { SetGravity(worldId, gx, gy, gz) }
func Step3D(worldId string, dt float64)                    { Step(worldId, dt) }
func CreateBox3D(worldId, bodyId string, x, y, z, hx, hy, hz, mass float64) {
	CreateBox(worldId, bodyId, x, y, z, hx, hy, hz, mass)
}
func CreateSphere3D(worldId, bodyId string, x, y, z, radius, mass float64) {
	CreateSphere(worldId, bodyId, x, y, z, radius, mass)
}
func DestroyBody3D(worldId, bodyId string)                  { DestroyBody(worldId, bodyId) }
func SetPosition3D(worldId, bodyId string, x, y, z float64) { SetPosition(worldId, bodyId, x, y, z) }
func GetPositionX3D(worldId, bodyId string) float64         { return GetPositionX(worldId, bodyId) }
func GetPositionY3D(worldId, bodyId string) float64         { return GetPositionY(worldId, bodyId) }
func GetPositionZ3D(worldId, bodyId string) float64         { return GetPositionZ(worldId, bodyId) }
func SetVelocity3D(worldId, bodyId string, vx, vy, vz float64) {
	SetVelocity(worldId, bodyId, vx, vy, vz)
}
func GetVelocityX3D(worldId, bodyId string) float64 { return GetVelocityX(worldId, bodyId) }
func GetVelocityY3D(worldId, bodyId string) float64 { return GetVelocityY(worldId, bodyId) }
func GetVelocityZ3D(worldId, bodyId string) float64 { return GetVelocityZ(worldId, bodyId) }
func GetYaw3D(worldId, bodyId string) float64       { return GetRotationY(worldId, bodyId) }
func GetPitch3D(worldId, bodyId string) float64     { return GetRotationX(worldId, bodyId) }
func GetRoll3D(worldId, bodyId string) float64      { return GetRotationZ(worldId, bodyId) }
func SetRotation3D(worldId, bodyId string, rx, ry, rz float64) {
	SetRotation(worldId, bodyId, rx, ry, rz)
}
func ApplyForce3D(worldId, bodyId string, fx, fy, fz float64) {
	ApplyForce(worldId, bodyId, fx, fy, fz)
}
func ApplyImpulse3D(worldId, bodyId string, ix, iy, iz float64) {
	ApplyImpulse(worldId, bodyId, ix, iy, iz)
}
func RayCastFromDir3D(worldId string, sx, sy, sz, dx, dy, dz, maxDist float64) int {
	return RayCast(worldId, sx, sy, sz, dx, dy, dz, maxDist)
}
func SphereCastFromDir3D(worldId string, sx, sy, sz, dx, dy, dz, radius, maxDist float64) int {
	w := getWorld(worldId)
	if w == nil {
		lastRayMu.Lock()
		lastRay.hit = false
		lastRayMu.Unlock()
		return 0
	}
	norm := math.Sqrt(dx*dx + dy*dy + dz*dz)
	if norm < 1e-9 {
		lastRayMu.Lock()
		lastRay.hit = false
		lastRayMu.Unlock()
		return 0
	}
	dx, dy, dz = dx/norm, dy/norm, dz/norm
	if radius < 0 {
		radius = 0
	}
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
		min, max := bodyAABB(b)
		min.x -= radius
		min.y -= radius
		min.z -= radius
		max.x += radius
		max.y += radius
		max.z += radius
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
		return 1
	}
	return 0
}
func RayHitX3D() float64       { return GetRayCastHitX() }
func RayHitY3D() float64       { return GetRayCastHitY() }
func RayHitZ3D() float64       { return GetRayCastHitZ() }
func RayHitBody3D() string     { return GetRayCastHitBody() }
func RayHitNormalX3D() float64 { return GetRayCastHitNormalX() }
func RayHitNormalY3D() float64 { return GetRayCastHitNormalY() }
func RayHitNormalZ3D() float64 { return GetRayCastHitNormalZ() }
