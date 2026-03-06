// Package dbp - Collision module: Raycast, Spherecast, ObjectCollides, PointInObject.
package dbp

import (
	"fmt"
	"math"
	"strconv"

	"cyberbasic/compiler/vm"
)

const (
	defaultCollisionWorld3D = "default"
	defaultCollisionWorld2D = "default"
	defaultRayMaxDist       = 10000.0
)

func toFloat64Collision(v interface{}) float64 {
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

// registerCollision registers DBP-style collision commands.
func registerCollision(v *vm.VM) {
	// --- 3D Raycast (Bullet) ---
	v.RegisterForeign("Raycast", func(args []interface{}) (interface{}, error) {
		if len(args) < 6 {
			return nil, fmt.Errorf("Raycast(ox, oy, oz, dx, dy, dz) requires 6 arguments")
		}
		ox, oy, oz := toFloat64Collision(args[0]), toFloat64Collision(args[1]), toFloat64Collision(args[2])
		dx, dy, dz := toFloat64Collision(args[3]), toFloat64Collision(args[4]), toFloat64Collision(args[5])
		maxDist := defaultRayMaxDist
		if len(args) >= 7 {
			maxDist = toFloat64Collision(args[6])
		}
		return v.CallForeign("RayCastFromDir3D", []interface{}{
			defaultCollisionWorld3D, ox, oy, oz, dx, dy, dz, maxDist,
		})
	})
	v.RegisterForeign("RayHitX", func(args []interface{}) (interface{}, error) {
		return v.CallForeign("RayHitX3D", nil)
	})
	v.RegisterForeign("RayHitY", func(args []interface{}) (interface{}, error) {
		return v.CallForeign("RayHitY3D", nil)
	})
	v.RegisterForeign("RayHitZ", func(args []interface{}) (interface{}, error) {
		return v.CallForeign("RayHitZ3D", nil)
	})
	v.RegisterForeign("RayHitBody", func(args []interface{}) (interface{}, error) {
		return v.CallForeign("RayHitBody3D", nil)
	})
	v.RegisterForeign("RayHitNormalX", func(args []interface{}) (interface{}, error) {
		return v.CallForeign("RayHitNormalX3D", nil)
	})
	v.RegisterForeign("RayHitNormalY", func(args []interface{}) (interface{}, error) {
		return v.CallForeign("RayHitNormalY3D", nil)
	})
	v.RegisterForeign("RayHitNormalZ", func(args []interface{}) (interface{}, error) {
		return v.CallForeign("RayHitNormalZ3D", nil)
	})

	// --- 2D Raycast (Box2D) ---
	v.RegisterForeign("Raycast2D", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("Raycast2D(ox, oy, dx, dy) requires 4 arguments")
		}
		ox, oy := toFloat64Collision(args[0]), toFloat64Collision(args[1])
		dx, dy := toFloat64Collision(args[2]), toFloat64Collision(args[3])
		dist := defaultRayMaxDist
		if dx*dx+dy*dy > 1e-10 {
			norm := math.Sqrt(dx*dx + dy*dy)
			dx, dy = dx/norm, dy/norm
		}
		tx, ty := ox+dx*dist, oy+dy*dist
		return v.CallForeign("RayCast2D", []interface{}{defaultCollisionWorld2D, ox, oy, tx, ty})
	})
	v.RegisterForeign("RayHitX2D", func(args []interface{}) (interface{}, error) {
		return v.CallForeign("RayHitX2D", nil)
	})
	v.RegisterForeign("RayHitY2D", func(args []interface{}) (interface{}, error) {
		return v.CallForeign("RayHitY2D", nil)
	})
	v.RegisterForeign("RayHitBody2D", func(args []interface{}) (interface{}, error) {
		return v.CallForeign("RayHitBody2D", nil)
	})

	// --- Spherecast (3D): thick ray - uses raycast with radius check; stub returns 0 if no native support ---
	v.RegisterForeign("Spherecast", func(args []interface{}) (interface{}, error) {
		if len(args) < 7 {
			return nil, fmt.Errorf("Spherecast(ox, oy, oz, dx, dy, dz, radius) requires 7 arguments")
		}
		// Delegate to Raycast; spherecast would need radius-aware hit - use raycast as approximation
		ox, oy, oz := toFloat64Collision(args[0]), toFloat64Collision(args[1]), toFloat64Collision(args[2])
		dx, dy, dz := toFloat64Collision(args[3]), toFloat64Collision(args[4]), toFloat64Collision(args[5])
		_ = toFloat64Collision(args[6]) // radius - not used in simple impl
		return v.CallForeign("RayCastFromDir3D", []interface{}{
			defaultCollisionWorld3D, ox, oy, oz, dx, dy, dz, defaultRayMaxDist,
		})
	})

	// --- Object collision (DBP objects by id) ---
	v.RegisterForeign("ObjectCollides", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("ObjectCollides(idA, idB) requires 2 arguments")
		}
		idA, idB := toInt(args[0]), toInt(args[1])
		objectsMu.Lock()
		a, okA := objects[idA]
		b, okB := objects[idB]
		objectsMu.Unlock()
		if !okA || !okB || !a.visible || !b.visible || !a.collision || !b.collision {
			return false, nil
		}
		return aabbOverlap(
			float64(a.x), float64(a.y), float64(a.z),
			float64(a.scaleX), float64(a.scaleY), float64(a.scaleZ),
			float64(b.x), float64(b.y), float64(b.z),
			float64(b.scaleX), float64(b.scaleY), float64(b.scaleZ),
		), nil
	})

	// --- Point in object (DBP object AABB) ---
	v.RegisterForeign("PointInObject", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("PointInObject(objectId, x, y, z) requires 4 arguments")
		}
		id := toInt(args[0])
		px, py, pz := toFloat64Collision(args[1]), toFloat64Collision(args[2]), toFloat64Collision(args[3])
		objectsMu.Lock()
		obj, ok := objects[id]
		objectsMu.Unlock()
		if !ok {
			return false, nil
		}
		hx, hy, hz := float64(obj.scaleX)/2, float64(obj.scaleY)/2, float64(obj.scaleZ)/2
		minX, maxX := float64(obj.x)-hx, float64(obj.x)+hx
		minY, maxY := float64(obj.y)-hy, float64(obj.y)+hy
		minZ, maxZ := float64(obj.z)-hz, float64(obj.z)+hz
		return px >= minX && px <= maxX && py >= minY && py <= maxY && pz >= minZ && pz <= maxZ, nil
	})

	// --- Physics body collision (wrap bullet) ---
	v.RegisterForeign("BodyCollides", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("BodyCollides(bodyIdA$, bodyIdB$) requires 2 arguments")
		}
		return v.CallForeign("CheckCollision3D", args[:2])
	})
}

func aabbOverlap(ax, ay, az, asx, asy, asz, bx, by, bz, bsx, bsy, bsz float64) bool {
	aMinX, aMaxX := ax-asx/2, ax+asx/2
	aMinY, aMaxY := ay-asy/2, ay+asy/2
	aMinZ, aMaxZ := az-asz/2, az+asz/2
	bMinX, bMaxX := bx-bsx/2, bx+bsx/2
	bMinY, bMaxY := by-bsy/2, by+bsy/2
	bMinZ, bMaxZ := bz-bsz/2, bz+bsz/2
	return aMinX <= bMaxX && aMaxX >= bMinX &&
		aMinY <= bMaxY && aMaxY >= bMinY &&
		aMinZ <= bMaxZ && aMaxZ >= bMinZ
}
