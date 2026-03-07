// Package dbp - Physics module: DBP-style wrappers over Box2D and Bullet.
package dbp

import (
	"fmt"
	"math"
	"strconv"
	"sync"

	"cyberbasic/compiler/vm"
)

const (
	defaultPhysicsWorld3D = "default"
	defaultPhysicsWorld2D = "default"
)

var (
	physicsBodyMap   = make(map[int]string)
	physicsBodyMapMu sync.Mutex
)

func toFloat64Physics(v interface{}) float64 {
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

// registerPhysics registers DBP-style physics commands that wrap Box2D and Bullet.
func registerPhysics(v *vm.VM) {
	// --- 3D Physics (Bullet) ---
	v.RegisterForeign("PhysicsOn", func(args []interface{}) (interface{}, error) {
		_, err := v.CallForeign("PhysicsEnable", nil)
		return nil, err
	})
	v.RegisterForeign("PhysicsOff", func(args []interface{}) (interface{}, error) {
		_, err := v.CallForeign("PhysicsDisable", nil)
		return nil, err
	})
	v.RegisterForeign("SetGravity", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("SetGravity(x, y, z) requires 3 arguments")
		}
		return v.CallForeign("PhysicsSetGravity", args[:3])
	})
	v.RegisterForeign("PhysicsStep", func(args []interface{}) (interface{}, error) {
		dt := 1.0 / 60.0
		if len(args) >= 1 {
			dt = toFloat64Physics(args[0])
		}
		return v.CallForeign("StepAllPhysics3D", []interface{}{dt})
	})
	v.RegisterForeign("MakeRigidBody", func(args []interface{}) (interface{}, error) {
		if len(args) < 5 {
			return nil, fmt.Errorf("MakeRigidBody(bodyId$, x, y, z, mass) requires 5 arguments")
		}
		bodyId := toString(args[0])
		x, y, z := toFloat64Physics(args[1]), toFloat64Physics(args[2]), toFloat64Physics(args[3])
		mass := toFloat64Physics(args[4])
		radius := 0.5
		if mass <= 0 {
			radius = 0.5
			mass = 1.0
		}
		_, err := v.CallForeign("CreateSphere3D", []interface{}{
			defaultPhysicsWorld3D, bodyId, x, y, z, radius, mass,
		})
		if err != nil {
			return nil, err
		}
		_, _ = v.CallForeign("PhysicsEnable", nil)
		return bodyId, nil
	})
	v.RegisterForeign("MakeStaticBody", func(args []interface{}) (interface{}, error) {
		if len(args) < 7 {
			return nil, fmt.Errorf("MakeStaticBody(bodyId$, x, y, z, sizeX, sizeY, sizeZ) requires 7 arguments")
		}
		bodyId := toString(args[0])
		x, y, z := toFloat64Physics(args[1]), toFloat64Physics(args[2]), toFloat64Physics(args[3])
		sx, sy, sz := toFloat64Physics(args[4]), toFloat64Physics(args[5]), toFloat64Physics(args[6])
		if sx <= 0 {
			sx = 1
		}
		if sy <= 0 {
			sy = 1
		}
		if sz <= 0 {
			sz = 1
		}
		_, err := v.CallForeign("CreateBox3D", []interface{}{
			defaultPhysicsWorld3D, bodyId, x, y, z, sx, sy, sz, 0.0,
		})
		if err != nil {
			return nil, err
		}
		_, _ = v.CallForeign("PhysicsEnable", nil)
		return bodyId, nil
	})
	// ApplyForce/ApplyImpulse: use bullet's ApplyForce, ApplyImpulse (bodyId, fx, fy, fz)
	// DBP aliases that add default world - not needed for bullet which uses "default" internally
	v.RegisterForeign("GetVelocityX", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return 0.0, nil
		}
		return v.CallForeign("GetVelocityX3D", []interface{}{defaultPhysicsWorld3D, toString(args[0])})
	})
	v.RegisterForeign("GetVelocityY", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return 0.0, nil
		}
		return v.CallForeign("GetVelocityY3D", []interface{}{defaultPhysicsWorld3D, toString(args[0])})
	})
	v.RegisterForeign("GetVelocityZ", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return 0.0, nil
		}
		return v.CallForeign("GetVelocityZ3D", []interface{}{defaultPhysicsWorld3D, toString(args[0])})
	})
	v.RegisterForeign("GetPositionX", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return 0.0, nil
		}
		return v.CallForeign("GetPositionX3D", []interface{}{defaultPhysicsWorld3D, toString(args[0])})
	})
	v.RegisterForeign("GetPositionY", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return 0.0, nil
		}
		return v.CallForeign("GetPositionY3D", []interface{}{defaultPhysicsWorld3D, toString(args[0])})
	})
	v.RegisterForeign("GetPositionZ", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return 0.0, nil
		}
		return v.CallForeign("GetPositionZ3D", []interface{}{defaultPhysicsWorld3D, toString(args[0])})
	})
	v.RegisterForeign("GetBodyX", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return 0.0, nil
		}
		return v.CallForeign("GetPositionX3D", []interface{}{defaultPhysicsWorld3D, toString(args[0])})
	})
	v.RegisterForeign("GetBodyY", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return 0.0, nil
		}
		return v.CallForeign("GetPositionY3D", []interface{}{defaultPhysicsWorld3D, toString(args[0])})
	})
	v.RegisterForeign("GetBodyZ", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return 0.0, nil
		}
		return v.CallForeign("GetPositionZ3D", []interface{}{defaultPhysicsWorld3D, toString(args[0])})
	})
	v.RegisterForeign("GetBodyVX", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return 0.0, nil
		}
		return v.CallForeign("GetVelocityX3D", []interface{}{defaultPhysicsWorld3D, toString(args[0])})
	})
	v.RegisterForeign("GetBodyVY", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return 0.0, nil
		}
		return v.CallForeign("GetVelocityY3D", []interface{}{defaultPhysicsWorld3D, toString(args[0])})
	})
	v.RegisterForeign("GetBodyVZ", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return 0.0, nil
		}
		return v.CallForeign("GetVelocityZ3D", []interface{}{defaultPhysicsWorld3D, toString(args[0])})
	})

	// --- Collision shapes (int ID -> Bullet body mapping) ---
	physicsBodyId := func(id int) string {
		physicsBodyMapMu.Lock()
		defer physicsBodyMapMu.Unlock()
		if b, ok := physicsBodyMap[id]; ok {
			return b
		}
		return ""
	}
	physicsBodyAlloc := func(id int) string {
		physicsBodyMapMu.Lock()
		defer physicsBodyMapMu.Unlock()
		bid := fmt.Sprintf("body_%d", id)
		physicsBodyMap[id] = bid
		return bid
	}

	v.RegisterForeign("MakeBoxCollider", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("MakeBoxCollider(id, sx, sy, sz) requires 4 arguments")
		}
		id := toInt(args[0])
		sx, sy, sz := toFloat64Physics(args[1]), toFloat64Physics(args[2]), toFloat64Physics(args[3])
		if sx <= 0 {
			sx = 1
		}
		if sy <= 0 {
			sy = 1
		}
		if sz <= 0 {
			sz = 1
		}
		bid := physicsBodyAlloc(id)
		_, err := v.CallForeign("CreateBox3D", []interface{}{
			defaultPhysicsWorld3D, bid, 0.0, 0.0, 0.0, sx, sy, sz, 0.0,
		})
		if err != nil {
			return nil, err
		}
		_, _ = v.CallForeign("PhysicsEnable", nil)
		return nil, nil
	})
	v.RegisterForeign("MakeSphereCollider", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("MakeSphereCollider(id, radius) requires 2 arguments")
		}
		id := toInt(args[0])
		radius := toFloat64Physics(args[1])
		if radius <= 0 {
			radius = 0.5
		}
		bid := physicsBodyAlloc(id)
		_, err := v.CallForeign("CreateSphere3D", []interface{}{
			defaultPhysicsWorld3D, bid, 0.0, 0.0, 0.0, radius, 0.0,
		})
		if err != nil {
			return nil, err
		}
		_, _ = v.CallForeign("PhysicsEnable", nil)
		return nil, nil
	})
	v.RegisterForeign("MakeCapsuleCollider", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("MakeCapsuleCollider(id, radius, height) requires 3 arguments")
		}
		id := toInt(args[0])
		radius := toFloat64Physics(args[1])
		height := toFloat64Physics(args[2])
		if radius <= 0 {
			radius = 0.5
		}
		if height <= 0 {
			height = 1
		}
		bid := physicsBodyAlloc(id)
		_, err := v.CallForeign("CreateCapsule3D", []interface{}{
			defaultPhysicsWorld3D, bid, 0.0, 0.0, 0.0, radius, height, 0.0,
		})
		if err != nil {
			return nil, err
		}
		_, _ = v.CallForeign("PhysicsEnable", nil)
		return nil, nil
	})
	v.RegisterForeign("MakeMeshCollider", func(args []interface{}) (interface{}, error) {
		return nil, fmt.Errorf("MakeMeshCollider is not supported by the shipped 3D fallback backend; use simpler colliders or gate with BulletFeatureAvailable(\"mesh_collider\")")
	})

	v.RegisterForeign("MakeRigidBodyId", func(args []interface{}) (interface{}, error) {
		if len(args) < 5 {
			return nil, fmt.Errorf("MakeRigidBodyId(id, x, y, z, mass) requires 5 arguments")
		}
		id := toInt(args[0])
		x, y, z := toFloat64Physics(args[1]), toFloat64Physics(args[2]), toFloat64Physics(args[3])
		mass := toFloat64Physics(args[4])
		if mass <= 0 {
			mass = 1.0
		}
		bid := physicsBodyAlloc(id)
		_, err := v.CallForeign("CreateSphere3D", []interface{}{
			defaultPhysicsWorld3D, bid, x, y, z, 0.5, mass,
		})
		if err != nil {
			return nil, err
		}
		_, _ = v.CallForeign("PhysicsEnable", nil)
		return id, nil
	})

	// --- Rigid body queries (int ID) ---
	v.RegisterForeign("GetRigidBodyMass", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return 0.0, nil
		}
		bid := physicsBodyId(toInt(args[0]))
		if bid == "" {
			return 0.0, nil
		}
		return v.CallForeign("GetMass3D", []interface{}{defaultPhysicsWorld3D, bid})
	})
	v.RegisterForeign("GetRigidBodySpeed", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return 0.0, nil
		}
		bid := physicsBodyId(toInt(args[0]))
		if bid == "" {
			return 0.0, nil
		}
		vx, _ := v.CallForeign("GetVelocityX3D", []interface{}{defaultPhysicsWorld3D, bid})
		vy, _ := v.CallForeign("GetVelocityY3D", []interface{}{defaultPhysicsWorld3D, bid})
		vz, _ := v.CallForeign("GetVelocityZ3D", []interface{}{defaultPhysicsWorld3D, bid})
		vxf := toFloat64Physics(vx)
		vyf := toFloat64Physics(vy)
		vzf := toFloat64Physics(vz)
		return math.Sqrt(vxf*vxf + vyf*vyf + vzf*vzf), nil
	})
	v.RegisterForeign("GetRigidBodyAngularVelocity", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return 0.0, nil
		}
		bid := physicsBodyId(toInt(args[0]))
		if bid == "" {
			return 0.0, nil
		}
		ax, _ := v.CallForeign("GetAngularVelocityX3D", []interface{}{defaultPhysicsWorld3D, bid})
		ay, _ := v.CallForeign("GetAngularVelocityY3D", []interface{}{defaultPhysicsWorld3D, bid})
		az, _ := v.CallForeign("GetAngularVelocityZ3D", []interface{}{defaultPhysicsWorld3D, bid})
		axf := toFloat64Physics(ax)
		ayf := toFloat64Physics(ay)
		azf := toFloat64Physics(az)
		return math.Sqrt(axf*axf + ayf*ayf + azf*azf), nil
	})
	v.RegisterForeign("SetRigidBodyPosition", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("SetRigidBodyPosition(id, x, y, z) requires 4 arguments")
		}
		bid := physicsBodyId(toInt(args[0]))
		if bid == "" {
			return nil, nil
		}
		return v.CallForeign("SetPosition3D", []interface{}{
			defaultPhysicsWorld3D, bid, toFloat64Physics(args[1]), toFloat64Physics(args[2]), toFloat64Physics(args[3]),
		})
	})
	v.RegisterForeign("SetRigidBodyVelocity", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("SetRigidBodyVelocity(id, vx, vy, vz) requires 4 arguments")
		}
		bid := physicsBodyId(toInt(args[0]))
		if bid == "" {
			return nil, nil
		}
		return v.CallForeign("SetVelocity3D", []interface{}{
			defaultPhysicsWorld3D, bid, toFloat64Physics(args[1]), toFloat64Physics(args[2]), toFloat64Physics(args[3]),
		})
	})
	v.RegisterForeign("SetAngularVelocity", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("SetAngularVelocity(id, x, y, z) requires 4 arguments")
		}
		bid := physicsBodyId(toInt(args[0]))
		if bid == "" {
			return nil, nil
		}
		return v.CallForeign("SetAngularVelocity3D", []interface{}{
			defaultPhysicsWorld3D, bid, toFloat64Physics(args[1]), toFloat64Physics(args[2]), toFloat64Physics(args[3]),
		})
	})
	v.RegisterForeign("GetRigidBodyX", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return 0.0, nil
		}
		bid := physicsBodyId(toInt(args[0]))
		if bid == "" {
			return 0.0, nil
		}
		return v.CallForeign("GetPositionX3D", []interface{}{defaultPhysicsWorld3D, bid})
	})
	v.RegisterForeign("GetRigidBodyY", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return 0.0, nil
		}
		bid := physicsBodyId(toInt(args[0]))
		if bid == "" {
			return 0.0, nil
		}
		return v.CallForeign("GetPositionY3D", []interface{}{defaultPhysicsWorld3D, bid})
	})
	v.RegisterForeign("GetRigidBodyZ", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return 0.0, nil
		}
		bid := physicsBodyId(toInt(args[0]))
		if bid == "" {
			return 0.0, nil
		}
		return v.CallForeign("GetPositionZ3D", []interface{}{defaultPhysicsWorld3D, bid})
	})

	// --- 2D Physics (Box2D) ---
	v.RegisterForeign("PhysicsOn2D", func(args []interface{}) (interface{}, error) {
		gx, gy := 0.0, -9.8
		if len(args) >= 2 {
			gx, gy = toFloat64Physics(args[0]), toFloat64Physics(args[1])
		}
		return v.CallForeign("CreateWorld2D", []interface{}{defaultPhysicsWorld2D, gx, gy})
	})
	v.RegisterForeign("PhysicsOff2D", func(args []interface{}) (interface{}, error) {
		return v.CallForeign("DestroyWorld2D", []interface{}{defaultPhysicsWorld2D})
	})
	v.RegisterForeign("SetGravity2D", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("SetGravity2D(x, y) requires 2 arguments")
		}
		return v.CallForeign("Physics2DSetGravity", []interface{}{defaultPhysicsWorld2D, args[0], args[1]})
	})
	v.RegisterForeign("PhysicsStep2D", func(args []interface{}) (interface{}, error) {
		dt := 1.0 / 60.0
		if len(args) >= 1 {
			dt = toFloat64Physics(args[0])
		}
		return v.CallForeign("Step2D", []interface{}{defaultPhysicsWorld2D, dt})
	})
	v.RegisterForeign("MakeRigidBody2D", func(args []interface{}) (interface{}, error) {
		if len(args) < 6 {
			return nil, fmt.Errorf("MakeRigidBody2D(bodyId$, x, y, w, h, density) requires 6 arguments")
		}
		bodyId := toString(args[0])
		x, y := toFloat64Physics(args[1]), toFloat64Physics(args[2])
		w, h := toFloat64Physics(args[3]), toFloat64Physics(args[4])
		density := toFloat64Physics(args[5])
		if w <= 0 {
			w = 1
		}
		if h <= 0 {
			h = 1
		}
		if density <= 0 {
			density = 1.0
		}
		_, err := v.CallForeign("CreateBox2D", []interface{}{
			defaultPhysicsWorld2D, bodyId, x, y, w, h, density, 1,
		})
		return bodyId, err
	})
	v.RegisterForeign("MakeStaticBody2D", func(args []interface{}) (interface{}, error) {
		if len(args) < 5 {
			return nil, fmt.Errorf("MakeStaticBody2D(bodyId$, x, y, w, h) requires 5 arguments")
		}
		bodyId := toString(args[0])
		x, y := toFloat64Physics(args[1]), toFloat64Physics(args[2])
		w, h := toFloat64Physics(args[3]), toFloat64Physics(args[4])
		if w <= 0 {
			w = 1
		}
		if h <= 0 {
			h = 1
		}
		_, err := v.CallForeign("CreateBox2D", []interface{}{
			defaultPhysicsWorld2D, bodyId, x, y, w, h, 0, 0,
		})
		return bodyId, err
	})
	// ApplyForce2D/ApplyImpulse2D: use box2d's (world$, body$, fx, fy) - pass defaultPhysicsWorld2D
	v.RegisterForeign("GetVelocityX2D", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return 0.0, nil
		}
		if len(args) >= 2 {
			return v.CallForeign("GetVelocityX2DByBodyId", []interface{}{toString(args[0]), toString(args[1])})
		}
		return v.CallForeign("GetVelocityX2DByBodyId", []interface{}{defaultPhysicsWorld2D, toString(args[0])})
	})
	v.RegisterForeign("GetVelocityY2D", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return 0.0, nil
		}
		if len(args) >= 2 {
			return v.CallForeign("GetVelocityY2DByBodyId", []interface{}{toString(args[0]), toString(args[1])})
		}
		return v.CallForeign("GetVelocityY2DByBodyId", []interface{}{defaultPhysicsWorld2D, toString(args[0])})
	})
}
