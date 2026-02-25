// Package raylib: 3D editor and level builder (mouse ray, plane pick, grid snap, level objects).
package raylib

import (
	"cyberbasic/compiler/vm"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"sync"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
	lastMouseRay   rl.Ray
	lastMouseRayMu sync.Mutex

	levelObjects   = make(map[string]*levelObject)
	levelObjectsMu sync.RWMutex
	levelObjectSeq int
)

type levelObject struct {
	ModelId string
	X, Y, Z float64
	RotX, RotY, RotZ float64
	ScaleX, ScaleY, ScaleZ float64
}

func toFloat64Editor(v interface{}) float64 {
	switch x := v.(type) {
	case int:
		return float64(x)
	case int32:
		return float64(x)
	case float64:
		return x
	case float32:
		return float64(x)
	default:
		return 0
	}
}

func registerEditor(v *vm.VM) {
	// ---- Mouse ray (3D picking) ----
	v.RegisterForeign("GetMouseRay", func(args []interface{}) (interface{}, error) {
		mousePos := rl.GetMousePosition()
		ray := rl.GetMouseRay(mousePos, camera3D)
		lastMouseRayMu.Lock()
		lastMouseRay = ray
		lastMouseRayMu.Unlock()
		return 1, nil
	})
	v.RegisterForeign("GetMouseRayOriginX", func(args []interface{}) (interface{}, error) {
		lastMouseRayMu.Lock()
		defer lastMouseRayMu.Unlock()
		return float64(lastMouseRay.Position.X), nil
	})
	v.RegisterForeign("GetMouseRayOriginY", func(args []interface{}) (interface{}, error) {
		lastMouseRayMu.Lock()
		defer lastMouseRayMu.Unlock()
		return float64(lastMouseRay.Position.Y), nil
	})
	v.RegisterForeign("GetMouseRayOriginZ", func(args []interface{}) (interface{}, error) {
		lastMouseRayMu.Lock()
		defer lastMouseRayMu.Unlock()
		return float64(lastMouseRay.Position.Z), nil
	})
	v.RegisterForeign("GetMouseRayDirectionX", func(args []interface{}) (interface{}, error) {
		lastMouseRayMu.Lock()
		defer lastMouseRayMu.Unlock()
		return float64(lastMouseRay.Direction.X), nil
	})
	v.RegisterForeign("GetMouseRayDirectionY", func(args []interface{}) (interface{}, error) {
		lastMouseRayMu.Lock()
		defer lastMouseRayMu.Unlock()
		return float64(lastMouseRay.Direction.Y), nil
	})
	v.RegisterForeign("GetMouseRayDirectionZ", func(args []interface{}) (interface{}, error) {
		lastMouseRayMu.Lock()
		defer lastMouseRayMu.Unlock()
		return float64(lastMouseRay.Direction.Z), nil
	})

	// ---- Ray-plane collision ----
	v.RegisterForeign("GetRayCollisionPlane", func(args []interface{}) (interface{}, error) {
		if len(args) < 12 {
			return nil, fmt.Errorf("GetRayCollisionPlane requires (rayPosX,Y,Z, rayDirX,Y,Z, planeX,Y,Z, planeNormX,Y,Z)")
		}
		rayPos := rl.Vector3{X: toFloat32(args[0]), Y: toFloat32(args[1]), Z: toFloat32(args[2])}
		rayDir := rl.Vector3{X: toFloat32(args[3]), Y: toFloat32(args[4]), Z: toFloat32(args[5])}
		planePos := rl.Vector3{X: toFloat32(args[6]), Y: toFloat32(args[7]), Z: toFloat32(args[8])}
		planeNorm := rl.Vector3{X: toFloat32(args[9]), Y: toFloat32(args[10]), Z: toFloat32(args[11])}
		denom := rayDir.X*planeNorm.X + rayDir.Y*planeNorm.Y + rayDir.Z*planeNorm.Z
		var hit bool
		var point rl.Vector3
		var dist float32
		if math.Abs(float64(denom)) > 1e-6 {
			t := ((planePos.X-rayPos.X)*planeNorm.X + (planePos.Y-rayPos.Y)*planeNorm.Y + (planePos.Z-rayPos.Z)*planeNorm.Z) / denom
			if t >= 0 {
				point = rl.Vector3{X: rayPos.X + t*rayDir.X, Y: rayPos.Y + t*rayDir.Y, Z: rayPos.Z + t*rayDir.Z}
				dist = t
				hit = true
			}
		}
		lastRayCollisionMu.Lock()
		lastRayCollision = rl.RayCollision{Hit: hit, Point: point, Normal: planeNorm, Distance: dist}
		lastRayCollisionMu.Unlock()
		if hit {
			return 1, nil
		}
		return 0, nil
	})
	v.RegisterForeign("PickGroundPlane", func(args []interface{}) (interface{}, error) {
		mousePos := rl.GetMousePosition()
		ray := rl.GetMouseRay(mousePos, camera3D)
		lastMouseRayMu.Lock()
		lastMouseRay = ray
		lastMouseRayMu.Unlock()
		planePos := rl.Vector3{X: 0, Y: 0, Z: 0}
		planeNorm := rl.Vector3{X: 0, Y: 1, Z: 0}
		denom := ray.Direction.X*planeNorm.X + ray.Direction.Y*planeNorm.Y + ray.Direction.Z*planeNorm.Z
		var hit bool
		var point rl.Vector3
		var dist float32
		if math.Abs(float64(denom)) > 1e-6 {
			t := ((planePos.X-ray.Position.X)*planeNorm.X + (planePos.Y-ray.Position.Y)*planeNorm.Y + (planePos.Z-ray.Position.Z)*planeNorm.Z) / denom
			if t >= 0 {
				point = rl.Vector3{X: ray.Position.X + t*ray.Direction.X, Y: ray.Position.Y + t*ray.Direction.Y, Z: ray.Position.Z + t*ray.Direction.Z}
				dist = t
				hit = true
			}
		}
		lastRayCollisionMu.Lock()
		lastRayCollision = rl.RayCollision{Hit: hit, Point: point, Normal: planeNorm, Distance: dist}
		lastRayCollisionMu.Unlock()
		if hit {
			return 1, nil
		}
		return 0, nil
	})

	// ---- Grid snap ----
	v.RegisterForeign("SnapToGridX", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("SnapToGridX requires (x, gridSize)")
		}
		x := toFloat64Editor(args[0])
		grid := toFloat64Editor(args[1])
		if grid <= 0 {
			return x, nil
		}
		return math.Round(x/grid) * grid, nil
	})
	v.RegisterForeign("SnapToGridY", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("SnapToGridY requires (y, gridSize)")
		}
		y := toFloat64Editor(args[0])
		grid := toFloat64Editor(args[1])
		if grid <= 0 {
			return y, nil
		}
		return math.Round(y/grid) * grid, nil
	})
	v.RegisterForeign("SnapToGridZ", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("SnapToGridZ requires (z, gridSize)")
		}
		z := toFloat64Editor(args[0])
		grid := toFloat64Editor(args[1])
		if grid <= 0 {
			return z, nil
		}
		return math.Round(z/grid) * grid, nil
	})

	// ---- Level objects ----
	v.RegisterForeign("CreateLevelObject", func(args []interface{}) (interface{}, error) {
		if len(args) < 11 {
			return nil, fmt.Errorf("CreateLevelObject requires (id, modelId, x, y, z, rotX, rotY, rotZ, scaleX, scaleY, scaleZ)")
		}
		id := toString(args[0])
		levelObjectsMu.Lock()
		levelObjects[id] = &levelObject{
			ModelId: toString(args[1]),
			X: toFloat64Editor(args[2]), Y: toFloat64Editor(args[3]), Z: toFloat64Editor(args[4]),
			RotX: toFloat64Editor(args[5]), RotY: toFloat64Editor(args[6]), RotZ: toFloat64Editor(args[7]),
			ScaleX: toFloat64Editor(args[8]), ScaleY: toFloat64Editor(args[9]), ScaleZ: toFloat64Editor(args[10]),
		}
		levelObjectsMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("SetLevelObjectTransform", func(args []interface{}) (interface{}, error) {
		if len(args) < 10 {
			return nil, fmt.Errorf("SetLevelObjectTransform requires (id, x, y, z, rotX, rotY, rotZ, scaleX, scaleY, scaleZ)")
		}
		id := toString(args[0])
		levelObjectsMu.Lock()
		if o, ok := levelObjects[id]; ok {
			o.X, o.Y, o.Z = toFloat64Editor(args[1]), toFloat64Editor(args[2]), toFloat64Editor(args[3])
			o.RotX, o.RotY, o.RotZ = toFloat64Editor(args[4]), toFloat64Editor(args[5]), toFloat64Editor(args[6])
			o.ScaleX, o.ScaleY, o.ScaleZ = toFloat64Editor(args[7]), toFloat64Editor(args[8]), toFloat64Editor(args[9])
		}
		levelObjectsMu.Unlock()
		return nil, nil
	})
	// SetObjectPosition(id, x, y, z): set level object position
	v.RegisterForeign("SetObjectPosition", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("SetObjectPosition requires (id, x, y, z)")
		}
		id := toString(args[0])
		levelObjectsMu.Lock()
		if o, ok := levelObjects[id]; ok {
			o.X = toFloat64Editor(args[1])
			o.Y = toFloat64Editor(args[2])
			o.Z = toFloat64Editor(args[3])
		}
		levelObjectsMu.Unlock()
		return nil, nil
	})
	// RotateObject(id, pitch, yaw, roll): set level object rotation (radians; RotX=pitch, RotY=yaw, RotZ=roll)
	v.RegisterForeign("RotateObject", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("RotateObject requires (id, pitch, yaw, roll)")
		}
		id := toString(args[0])
		levelObjectsMu.Lock()
		if o, ok := levelObjects[id]; ok {
			o.RotX = toFloat64Editor(args[1])
			o.RotY = toFloat64Editor(args[2])
			o.RotZ = toFloat64Editor(args[3])
		}
		levelObjectsMu.Unlock()
		return nil, nil
	})
	// ScaleObject(id, sx, sy, sz): set level object scale
	v.RegisterForeign("ScaleObject", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("ScaleObject requires (id, scaleX, scaleY, scaleZ)")
		}
		id := toString(args[0])
		levelObjectsMu.Lock()
		if o, ok := levelObjects[id]; ok {
			o.ScaleX = toFloat64Editor(args[1])
			o.ScaleY = toFloat64Editor(args[2])
			o.ScaleZ = toFloat64Editor(args[3])
		}
		levelObjectsMu.Unlock()
		return nil, nil
	})
	// SetObjectRotation(obj, pitch, yaw, roll): alias for RotateObject
	v.RegisterForeign("SetObjectRotation", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("SetObjectRotation requires (obj, pitch, yaw, roll)")
		}
		id := toString(args[0])
		levelObjectsMu.Lock()
		if o, ok := levelObjects[id]; ok {
			o.RotX = toFloat64Editor(args[1])
			o.RotY = toFloat64Editor(args[2])
			o.RotZ = toFloat64Editor(args[3])
		}
		levelObjectsMu.Unlock()
		return nil, nil
	})
	// SetObjectScale(obj, sx, sy, sz): alias for ScaleObject
	v.RegisterForeign("SetObjectScale", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("SetObjectScale requires (obj, sx, sy, sz)")
		}
		id := toString(args[0])
		levelObjectsMu.Lock()
		if o, ok := levelObjects[id]; ok {
			o.ScaleX = toFloat64Editor(args[1])
			o.ScaleY = toFloat64Editor(args[2])
			o.ScaleZ = toFloat64Editor(args[3])
		}
		levelObjectsMu.Unlock()
		return nil, nil
	})
	// ObjectLookAt(obj, x, y, z): point object's forward (-Z) at world position (sets rotation from obj position to target)
	v.RegisterForeign("ObjectLookAt", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("ObjectLookAt requires (obj, x, y, z)")
		}
		id := toString(args[0])
		tx := toFloat64Editor(args[1])
		ty := toFloat64Editor(args[2])
		tz := toFloat64Editor(args[3])
		levelObjectsMu.RLock()
		o, ok := levelObjects[id]
		levelObjectsMu.RUnlock()
		if !ok || o == nil {
			return nil, nil
		}
		dx := tx - o.X
		dy := ty - o.Y
		dz := tz - o.Z
		distXZ := math.Sqrt(dx*dx + dz*dz)
		yaw := math.Atan2(-dx, dz)
		pitch := math.Atan2(dy, distXZ)
		levelObjectsMu.Lock()
		o.RotX = pitch
		o.RotY = yaw
		o.RotZ = 0
		levelObjectsMu.Unlock()
		return nil, nil
	})
	// DrawObject(id): draw level object by id (alias for DrawLevelObject)
	v.RegisterForeign("DrawObject", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("DrawObject requires (id)")
		}
		id := toString(args[0])
		levelObjectsMu.RLock()
		o, ok := levelObjects[id]
		levelObjectsMu.RUnlock()
		if !ok || o == nil {
			return nil, nil
		}
		modelMu.Lock()
		model, ok := models[o.ModelId]
		modelMu.Unlock()
		if !ok {
			return nil, nil
		}
		pos := rl.Vector3{X: float32(o.X), Y: float32(o.Y), Z: float32(o.Z)}
		rotAxis := rl.Vector3{X: 0, Y: 1, Z: 0}
		rotAngle := float32(o.RotY)
		scale := rl.Vector3{X: float32(o.ScaleX), Y: float32(o.ScaleY), Z: float32(o.ScaleZ)}
		rl.DrawModelEx(model, pos, rotAxis, rotAngle, scale, rl.White)
		return nil, nil
	})
	getLevelObject := func(id string) *levelObject {
		levelObjectsMu.RLock()
		defer levelObjectsMu.RUnlock()
		return levelObjects[id]
	}
	v.RegisterForeign("GetLevelObjectX", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("GetLevelObjectX requires (id)")
		}
		o := getLevelObject(toString(args[0]))
		if o == nil {
			return 0.0, nil
		}
		return o.X, nil
	})
	v.RegisterForeign("GetLevelObjectY", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("GetLevelObjectY requires (id)")
		}
		o := getLevelObject(toString(args[0]))
		if o == nil {
			return 0.0, nil
		}
		return o.Y, nil
	})
	v.RegisterForeign("GetLevelObjectZ", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("GetLevelObjectZ requires (id)")
		}
		o := getLevelObject(toString(args[0]))
		if o == nil {
			return 0.0, nil
		}
		return o.Z, nil
	})
	v.RegisterForeign("GetLevelObjectRotX", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("GetLevelObjectRotX requires (id)")
		}
		o := getLevelObject(toString(args[0]))
		if o == nil {
			return 0.0, nil
		}
		return o.RotX, nil
	})
	v.RegisterForeign("GetLevelObjectRotY", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("GetLevelObjectRotY requires (id)")
		}
		o := getLevelObject(toString(args[0]))
		if o == nil {
			return 0.0, nil
		}
		return o.RotY, nil
	})
	v.RegisterForeign("GetLevelObjectRotZ", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("GetLevelObjectRotZ requires (id)")
		}
		o := getLevelObject(toString(args[0]))
		if o == nil {
			return 0.0, nil
		}
		return o.RotZ, nil
	})
	v.RegisterForeign("GetLevelObjectScaleX", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("GetLevelObjectScaleX requires (id)")
		}
		o := getLevelObject(toString(args[0]))
		if o == nil {
			return 1.0, nil
		}
		return o.ScaleX, nil
	})
	v.RegisterForeign("GetLevelObjectScaleY", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("GetLevelObjectScaleY requires (id)")
		}
		o := getLevelObject(toString(args[0]))
		if o == nil {
			return 1.0, nil
		}
		return o.ScaleY, nil
	})
	v.RegisterForeign("GetLevelObjectScaleZ", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("GetLevelObjectScaleZ requires (id)")
		}
		o := getLevelObject(toString(args[0]))
		if o == nil {
			return 1.0, nil
		}
		return o.ScaleZ, nil
	})
	v.RegisterForeign("GetLevelObjectModelId", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("GetLevelObjectModelId requires (id)")
		}
		o := getLevelObject(toString(args[0]))
		if o == nil {
			return "", nil
		}
		return o.ModelId, nil
	})
	v.RegisterForeign("DeleteLevelObject", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("DeleteLevelObject requires (id)")
		}
		id := toString(args[0])
		levelObjectsMu.Lock()
		delete(levelObjects, id)
		levelObjectsMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("GetLevelObjectCount", func(args []interface{}) (interface{}, error) {
		levelObjectsMu.RLock()
		n := len(levelObjects)
		levelObjectsMu.RUnlock()
		return n, nil
	})
	v.RegisterForeign("GetLevelObjectId", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("GetLevelObjectId requires (index)")
		}
		idx := int(toFloat64Editor(args[0]))
		levelObjectsMu.RLock()
		var i int
		var out string
		for id := range levelObjects {
			if i == idx {
				out = id
				break
			}
			i++
		}
		levelObjectsMu.RUnlock()
		return out, nil
	})
	v.RegisterForeign("DrawLevelObject", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("DrawLevelObject requires (id)")
		}
		id := toString(args[0])
		levelObjectsMu.RLock()
		o, ok := levelObjects[id]
		levelObjectsMu.RUnlock()
		if !ok || o == nil {
			return nil, nil
		}
		modelMu.Lock()
		model, ok := models[o.ModelId]
		modelMu.Unlock()
		if !ok {
			return nil, nil
		}
		pos := rl.Vector3{X: float32(o.X), Y: float32(o.Y), Z: float32(o.Z)}
		rotAxis := rl.Vector3{X: 0, Y: 1, Z: 0}
		rotAngle := float32(o.RotY)
		scale := rl.Vector3{X: float32(o.ScaleX), Y: float32(o.ScaleY), Z: float32(o.ScaleZ)}
		rl.DrawModelEx(model, pos, rotAxis, rotAngle, scale, rl.White)
		return nil, nil
	})
	v.RegisterForeign("DuplicateLevelObject", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("DuplicateLevelObject requires (id)")
		}
		id := toString(args[0])
		levelObjectsMu.Lock()
		o, ok := levelObjects[id]
		if !ok || o == nil {
			levelObjectsMu.Unlock()
			return "", nil
		}
		levelObjectSeq++
		newId := fmt.Sprintf("obj_%d", levelObjectSeq)
		levelObjects[newId] = &levelObject{
			ModelId: o.ModelId,
			X: o.X, Y: o.Y, Z: o.Z,
			RotX: o.RotX, RotY: o.RotY, RotZ: o.RotZ,
			ScaleX: o.ScaleX, ScaleY: o.ScaleY, ScaleZ: o.ScaleZ,
		}
		levelObjectsMu.Unlock()
		return newId, nil
	})

	// SaveLevel / LoadLevel
	type levelObjectJSON struct {
		ID      string  `json:"id"`
		ModelId string  `json:"modelId"`
		X       float64 `json:"x"`
		Y       float64 `json:"y"`
		Z       float64 `json:"z"`
		RotX    float64 `json:"rotX"`
		RotY    float64 `json:"rotY"`
		RotZ    float64 `json:"rotZ"`
		ScaleX  float64 `json:"scaleX"`
		ScaleY  float64 `json:"scaleY"`
		ScaleZ  float64 `json:"scaleZ"`
	}
	v.RegisterForeign("SaveLevel", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("SaveLevel requires (path)")
		}
		path := toString(args[0])
		levelObjectsMu.RLock()
		var list []levelObjectJSON
		for id, o := range levelObjects {
			list = append(list, levelObjectJSON{
				ID: id, ModelId: o.ModelId,
				X: o.X, Y: o.Y, Z: o.Z,
				RotX: o.RotX, RotY: o.RotY, RotZ: o.RotZ,
				ScaleX: o.ScaleX, ScaleY: o.ScaleY, ScaleZ: o.ScaleZ,
			})
		}
		levelObjectsMu.RUnlock()
		raw, err := json.MarshalIndent(list, "", "  ")
		if err != nil {
			return nil, err
		}
		if err := os.WriteFile(path, raw, 0644); err != nil {
			return nil, err
		}
		return nil, nil
	})
	v.RegisterForeign("LoadLevel", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("LoadLevel requires (path)")
		}
		path := toString(args[0])
		raw, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}
		var list []levelObjectJSON
		if err := json.Unmarshal(raw, &list); err != nil {
			return nil, err
		}
		levelObjectsMu.Lock()
		levelObjects = make(map[string]*levelObject)
		for _, item := range list {
			levelObjects[item.ID] = &levelObject{
				ModelId: item.ModelId,
				X: item.X, Y: item.Y, Z: item.Z,
				RotX: item.RotX, RotY: item.RotY, RotZ: item.RotZ,
				ScaleX: item.ScaleX, ScaleY: item.ScaleY, ScaleZ: item.ScaleZ,
			}
		}
		levelObjectsMu.Unlock()
		return nil, nil
	})

	// ---- Camera readback (for editor UI / save) ----
	v.RegisterForeign("GetCameraPositionX", func(args []interface{}) (interface{}, error) {
		return float64(camera3D.Position.X), nil
	})
	v.RegisterForeign("GetCameraPositionY", func(args []interface{}) (interface{}, error) {
		return float64(camera3D.Position.Y), nil
	})
	v.RegisterForeign("GetCameraPositionZ", func(args []interface{}) (interface{}, error) {
		return float64(camera3D.Position.Z), nil
	})
	v.RegisterForeign("GetCameraTargetX", func(args []interface{}) (interface{}, error) {
		return float64(camera3D.Target.X), nil
	})
	v.RegisterForeign("GetCameraTargetY", func(args []interface{}) (interface{}, error) {
		return float64(camera3D.Target.Y), nil
	})
	v.RegisterForeign("GetCameraTargetZ", func(args []interface{}) (interface{}, error) {
		return float64(camera3D.Target.Z), nil
	})
}
