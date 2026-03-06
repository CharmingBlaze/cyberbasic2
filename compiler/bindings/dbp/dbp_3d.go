// Package dbp: 3D Game API - window, camera, objects, math, replication.
//
// Modular 3D commands for FPS, RPGs, sandbox, survival, and multiplayer.
// See docs/3D_GAME_API.md for full reference.
package dbp

import (
	"fmt"
	"math"
	"sync"

	"cyberbasic/compiler/vm"
	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
	// syncObjects tracks which object IDs are marked for replication.
	syncObjects   = make(map[int]bool)
	syncObjectsMu sync.Mutex
)

// register3D adds 3D-specific DBP commands: window aliases, camera queries,
// object creation (MakeCylinder, MakeGrid), object queries, parenting, tags,
// replication, and 3D math helpers.
func register3D(v *vm.VM) {
	register3DWindow(v)
	register3DCamera(v)
	register3DObjects(v)
	register3DMesh(v)
	register3DAnimation(v)
	register3DTerrain(v)
	register3DMath(v)
	register3DReplication(v)
}

// --- Window & Rendering ---
func register3DWindow(v *vm.VM) {
	// Window(width, height, title$): DBP alias for InitWindow
	v.RegisterForeign("Window", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("Window(width, height, title) requires 3 arguments")
		}
		return v.CallForeign("InitWindow", args)
	})
	// CloseWindow: DBP alias
	v.RegisterForeign("CloseWindow", func(args []interface{}) (interface{}, error) {
		rl.CloseWindow()
		return nil, nil
	})
	// SetTargetFPS: already in raylib; add DBP alias if needed
	v.RegisterForeign("SetTargetFPS", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("SetTargetFPS(value) requires 1 argument")
		}
		rl.SetTargetFPS(int32(toInt(args[0])))
		return nil, nil
	})
	// SetFramerate(cap): DBP alias for SetTargetFPS. 0 = uncapped.
	v.RegisterForeign("SetFramerate", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("SetFramerate(cap) requires 1 argument")
		}
		rl.SetTargetFPS(int32(toInt(args[0])))
		return nil, nil
	})
	// Clear, StartDraw, EndDraw, Start3D, End3D: already in dbp.go
}

// --- Camera (standard + queries) ---
func register3DCamera(v *vm.VM) {
	// PointCameraAt(x, y, z): alias for PointCamera / SetCameraTarget
	v.RegisterForeign("PointCameraAt", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("PointCameraAt(x, y, z) requires 3 arguments")
		}
		return v.CallForeign("SetCameraTarget", args)
	})
	// SetCameraFOV(value): set camera field of view
	v.RegisterForeign("SetCameraFOV", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("SetCameraFOV(value) requires 1 argument")
		}
		return v.CallForeign("SetCameraFOV", args)
	})
	// SetCameraRange(near, far): projection near/far - raylib uses fixed values; no-op for compatibility
	v.RegisterForeign("SetCameraRange", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("SetCameraRange(near, far) requires 2 arguments")
		}
		return nil, nil
	})
	// SetCameraUp(x, y, z): set default camera up vector via SetCamera3D
	v.RegisterForeign("SetCameraUp", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("SetCameraUp(x, y, z) requires 3 arguments")
		}
		px, _ := v.CallForeign("GetCameraPositionX", nil)
		py, _ := v.CallForeign("GetCameraPositionY", nil)
		pz, _ := v.CallForeign("GetCameraPositionZ", nil)
		tx, _ := v.CallForeign("GetCameraTargetX", nil)
		ty, _ := v.CallForeign("GetCameraTargetY", nil)
		tz, _ := v.CallForeign("GetCameraTargetZ", nil)
		ux, uy, uz := toFloat32(args[0]), toFloat32(args[1]), toFloat32(args[2])
		_, err := v.CallForeign("SetCamera3D", []interface{}{
			toFloat32(px), toFloat32(py), toFloat32(pz),
			toFloat32(tx), toFloat32(ty), toFloat32(tz),
			ux, uy, uz,
		})
		return nil, err
	})

	// GetCameraX, GetCameraY, GetCameraZ: read camera position (delegate to raylib)
	v.RegisterForeign("GetCameraX", func(args []interface{}) (interface{}, error) {
		return v.CallForeign("GetCameraPositionX", nil)
	})
	v.RegisterForeign("GetCameraY", func(args []interface{}) (interface{}, error) {
		return v.CallForeign("GetCameraPositionY", nil)
	})
	v.RegisterForeign("GetCameraZ", func(args []interface{}) (interface{}, error) {
		return v.CallForeign("GetCameraPositionZ", nil)
	})
	// GetCameraPitch, GetCameraYaw: derived from position and target
	v.RegisterForeign("GetCameraPitch", func(args []interface{}) (interface{}, error) {
		px, _ := v.CallForeign("GetCameraPositionX", nil)
		py, _ := v.CallForeign("GetCameraPositionY", nil)
		pz, _ := v.CallForeign("GetCameraPositionZ", nil)
		tx, _ := v.CallForeign("GetCameraTargetX", nil)
		ty, _ := v.CallForeign("GetCameraTargetY", nil)
		tz, _ := v.CallForeign("GetCameraTargetZ", nil)
		dx := toFloat64(tx) - toFloat64(px)
		dy := toFloat64(ty) - toFloat64(py)
		dz := toFloat64(tz) - toFloat64(pz)
		dist := math.Sqrt(dx*dx + dy*dy + dz*dz)
		if dist < 1e-10 {
			return 0.0, nil
		}
		pitch := math.Asin(dy / dist)
		return pitch * 180 / math.Pi, nil
	})
	v.RegisterForeign("GetCameraYaw", func(args []interface{}) (interface{}, error) {
		px, _ := v.CallForeign("GetCameraPositionX", nil)
		pz, _ := v.CallForeign("GetCameraPositionZ", nil)
		tx, _ := v.CallForeign("GetCameraTargetX", nil)
		tz, _ := v.CallForeign("GetCameraTargetZ", nil)
		dx := toFloat64(tx) - toFloat64(px)
		dz := toFloat64(tz) - toFloat64(pz)
		yaw := math.Atan2(dz, dx)
		return yaw * 180 / math.Pi, nil
	})
}

func toFloat64(v interface{}) float64 {
	switch x := v.(type) {
	case int:
		return float64(x)
	case float64:
		return x
	case float32:
		return float64(x)
	default:
		return 0
	}
}

// --- 3D Object creation (MakeCylinder, MakeGrid) + queries + parenting + tags ---
func register3DObjects(v *vm.VM) {
	// MakeCylinder(id, radius, height): procedural cylinder
	v.RegisterForeign("MakeCylinder", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("MakeCylinder(id, radius, height) requires 3 arguments")
		}
		id := toInt(args[0])
		radius := toFloat32(args[1])
		height := toFloat32(args[2])
		mesh := rl.GenMeshCylinder(radius, height, 16)
		model := rl.LoadModelFromMesh(mesh)
		objectsMu.Lock()
		objects[id] = newDbpObject(model)
		objectsMu.Unlock()
		return nil, nil
	})
	// MakeGrid(id, size, spacing): procedural grid plane
	v.RegisterForeign("MakeGrid", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("MakeGrid(id, size, spacing) requires 3 arguments")
		}
		id := toInt(args[0])
		size := toFloat32(args[1])
		spacing := toFloat32(args[2])
		slices := int(size / spacing)
		if slices < 1 {
			slices = 1
		}
		mesh := rl.GenMeshPlane(size, size, slices, slices)
		model := rl.LoadModelFromMesh(mesh)
		objectsMu.Lock()
		objects[id] = newDbpObject(model)
		objectsMu.Unlock()
		return nil, nil
	})
	// LoadObject(id, file$): DBP arg order - id first
	v.RegisterForeign("LoadObject", func(args []interface{}) (interface{}, error) {
		return v.CallForeign("LoadObjectId", args)
	})

	// --- Object query commands ---
	v.RegisterForeign("GetObjectX", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return 0.0, nil
		}
		id := toInt(args[0])
		objectsMu.Lock()
		obj, ok := objects[id]
		objectsMu.Unlock()
		if !ok {
			return 0.0, nil
		}
		return float64(obj.x), nil
	})
	v.RegisterForeign("GetObjectY", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return 0.0, nil
		}
		id := toInt(args[0])
		objectsMu.Lock()
		obj, ok := objects[id]
		objectsMu.Unlock()
		if !ok {
			return 0.0, nil
		}
		return float64(obj.y), nil
	})
	v.RegisterForeign("GetObjectZ", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return 0.0, nil
		}
		id := toInt(args[0])
		objectsMu.Lock()
		obj, ok := objects[id]
		objectsMu.Unlock()
		if !ok {
			return 0.0, nil
		}
		return float64(obj.z), nil
	})
	v.RegisterForeign("GetObjectPitch", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return 0.0, nil
		}
		id := toInt(args[0])
		objectsMu.Lock()
		obj, ok := objects[id]
		objectsMu.Unlock()
		if !ok {
			return 0.0, nil
		}
		return float64(obj.pitch), nil
	})
	v.RegisterForeign("GetObjectYaw", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return 0.0, nil
		}
		id := toInt(args[0])
		objectsMu.Lock()
		obj, ok := objects[id]
		objectsMu.Unlock()
		if !ok {
			return 0.0, nil
		}
		return float64(obj.yaw), nil
	})
	v.RegisterForeign("GetObjectRoll", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return 0.0, nil
		}
		id := toInt(args[0])
		objectsMu.Lock()
		obj, ok := objects[id]
		objectsMu.Unlock()
		if !ok {
			return 0.0, nil
		}
		return float64(obj.roll), nil
	})
	v.RegisterForeign("GetObjectScaleX", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return 1.0, nil
		}
		id := toInt(args[0])
		objectsMu.Lock()
		obj, ok := objects[id]
		objectsMu.Unlock()
		if !ok {
			return 1.0, nil
		}
		return float64(obj.scaleX), nil
	})
	v.RegisterForeign("GetObjectScaleY", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return 1.0, nil
		}
		id := toInt(args[0])
		objectsMu.Lock()
		obj, ok := objects[id]
		objectsMu.Unlock()
		if !ok {
			return 1.0, nil
		}
		return float64(obj.scaleY), nil
	})
	v.RegisterForeign("GetObjectScaleZ", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return 1.0, nil
		}
		id := toInt(args[0])
		objectsMu.Lock()
		obj, ok := objects[id]
		objectsMu.Unlock()
		if !ok {
			return 1.0, nil
		}
		return float64(obj.scaleZ), nil
	})

	// --- Object parenting ---
	v.RegisterForeign("ParentObject", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("ParentObject(childID, parentID) requires 2 arguments")
		}
		childID := toInt(args[0])
		parentID := toInt(args[1])
		objectsMu.Lock()
		if obj, ok := objects[childID]; ok {
			obj.parentID = parentID
		}
		objectsMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("AttachObject", func(args []interface{}) (interface{}, error) {
		return v.CallForeign("ParentObject", []interface{}{args[0], args[1]})
	})
	v.RegisterForeign("UnparentObject", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("UnparentObject(id) requires 1 argument")
		}
		id := toInt(args[0])
		objectsMu.Lock()
		if obj, ok := objects[id]; ok {
			obj.parentID = -1
		}
		objectsMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("DetachObject", func(args []interface{}) (interface{}, error) {
		return v.CallForeign("UnparentObject", args)
	})

	// --- Object tags ---
	v.RegisterForeign("SetObjectTag", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("SetObjectTag(id, tag) requires 2 arguments")
		}
		id := toInt(args[0])
		tag := toString(args[1])
		objectsMu.Lock()
		if obj, ok := objects[id]; ok {
			obj.tag = tag
		}
		objectsMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("GetObjectTag", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return "", nil
		}
		id := toInt(args[0])
		objectsMu.Lock()
		obj, ok := objects[id]
		objectsMu.Unlock()
		if !ok {
			return "", nil
		}
		return obj.tag, nil
	})
}

// --- Terrain DBP wrappers ---
func register3DTerrain(v *vm.VM) {
	v.RegisterForeign("SetTerrainHeight", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("SetTerrainHeight(terrainId, x, z, height) requires 4 arguments")
		}
		return v.CallForeign("TerrainFlatten", []interface{}{
			toString(args[0]), toFloat64(args[1]), toFloat64(args[2]), 0.5, toFloat64(args[3]),
		})
	})
	v.RegisterForeign("PaintTerrainTexture", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("PaintTerrainTexture(terrainId, x, z, layer) requires 4 arguments")
		}
		return v.CallForeign("TerrainPaint", []interface{}{
			toString(args[0]), toFloat64(args[1]), toFloat64(args[2]), 2.0, toFloat64(args[3]), 0.5,
		})
	})
}

// --- 3D Math helpers ---
func register3DMath(v *vm.VM) {
	v.RegisterForeign("Distance3D", func(args []interface{}) (interface{}, error) {
		if len(args) < 6 {
			return nil, fmt.Errorf("Distance3D(x1,y1,z1, x2,y2,z2) requires 6 arguments")
		}
		x1, y1, z1 := toFloat64(args[0]), toFloat64(args[1]), toFloat64(args[2])
		x2, y2, z2 := toFloat64(args[3]), toFloat64(args[4]), toFloat64(args[5])
		dx, dy, dz := x2-x1, y2-y1, z2-z1
		return math.Sqrt(dx*dx + dy*dy + dz*dz), nil
	})
	v.RegisterForeign("AngleBetween3D", func(args []interface{}) (interface{}, error) {
		if len(args) < 6 {
			return nil, fmt.Errorf("AngleBetween3D(x1,y1,z1, x2,y2,z2) requires 6 arguments")
		}
		x1, y1, z1 := toFloat64(args[0]), toFloat64(args[1]), toFloat64(args[2])
		x2, y2, z2 := toFloat64(args[3]), toFloat64(args[4]), toFloat64(args[5])
		dx, dy, dz := x2-x1, y2-y1, z2-z1
		return math.Atan2(math.Sqrt(dx*dx+dz*dz), dy), nil
	})
	v.RegisterForeign("Normalize3D", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("Normalize3D(x, y, z) requires 3 arguments")
		}
		x, y, z := toFloat64(args[0]), toFloat64(args[1]), toFloat64(args[2])
		len := math.Sqrt(x*x + y*y + z*z)
		if len < 1e-10 {
			return []interface{}{0.0, 0.0, 0.0}, nil
		}
		return []interface{}{x / len, y / len, z / len}, nil
	})
	v.RegisterForeign("Dot3D", func(args []interface{}) (interface{}, error) {
		if len(args) < 6 {
			return nil, fmt.Errorf("Dot3D(x1,y1,z1, x2,y2,z2) requires 6 arguments")
		}
		x1, y1, z1 := toFloat64(args[0]), toFloat64(args[1]), toFloat64(args[2])
		x2, y2, z2 := toFloat64(args[3]), toFloat64(args[4]), toFloat64(args[5])
		return x1*x2 + y1*y2 + z1*z2, nil
	})
	v.RegisterForeign("Cross3D", func(args []interface{}) (interface{}, error) {
		if len(args) < 6 {
			return nil, fmt.Errorf("Cross3D(x1,y1,z1, x2,y2,z2) requires 6 arguments")
		}
		x1, y1, z1 := toFloat64(args[0]), toFloat64(args[1]), toFloat64(args[2])
		x2, y2, z2 := toFloat64(args[3]), toFloat64(args[4]), toFloat64(args[5])
		cx := y1*z2 - z1*y2
		cy := z1*x2 - x1*z2
		cz := x1*y2 - y1*x2
		return []interface{}{cx, cy, cz}, nil
	})

	// Matrix/Quaternion stubs
	v.RegisterForeign("MakeQuaternion", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("MakeQuaternion(id, pitch, yaw, roll) requires 4 arguments")
		}
		return nil, nil
	})
	v.RegisterForeign("RotateObjectQuat", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("RotateObjectQuat(id, quatID) requires 2 arguments")
		}
		return nil, nil
	})
	v.RegisterForeign("GetObjectMatrix", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("GetObjectMatrix(id) requires 1 argument")
		}
		// Return identity 4x4 matrix as flat array
		return []interface{}{1.0, 0.0, 0.0, 0.0, 0.0, 1.0, 0.0, 0.0, 0.0, 0.0, 1.0, 0.0, 0.0, 0.0, 0.0, 1.0}, nil
	})
}

// --- 3D Replication (SyncObject, UnsyncObject, SetObjectOwner, GetObjectOwner) ---
func register3DReplication(v *vm.VM) {
	v.RegisterForeign("SyncObject", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("SyncObject(id) requires 1 argument")
		}
		id := toInt(args[0])
		syncObjectsMu.Lock()
		syncObjects[id] = true
		syncObjectsMu.Unlock()
		objectsMu.Lock()
		if obj, ok := objects[id]; ok {
			obj.syncMe = true
		}
		objectsMu.Unlock()
		// Register with game replication
		v.CallForeign("ReplicatePosition", []interface{}{fmt.Sprintf("obj_%d", id)})
		return nil, nil
	})
	v.RegisterForeign("UnsyncObject", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("UnsyncObject(id) requires 1 argument")
		}
		id := toInt(args[0])
		syncObjectsMu.Lock()
		delete(syncObjects, id)
		syncObjectsMu.Unlock()
		objectsMu.Lock()
		if obj, ok := objects[id]; ok {
			obj.syncMe = false
		}
		objectsMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("SetObjectOwner", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("SetObjectOwner(id, playerID) requires 2 arguments")
		}
		id := toInt(args[0])
		ownerID := toInt(args[1])
		objectsMu.Lock()
		if obj, ok := objects[id]; ok {
			obj.ownerID = ownerID
		}
		objectsMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("GetObjectOwner", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return -1, nil
		}
		id := toInt(args[0])
		objectsMu.Lock()
		obj, ok := objects[id]
		objectsMu.Unlock()
		if !ok {
			return -1, nil
		}
		return obj.ownerID, nil
	})
}
