// Package dbp: Level loading - LoadLevel, DrawLevel, UnloadLevel.
package dbp

import (
	"fmt"
	"path/filepath"
	"sync"

	"cyberbasic/compiler/bindings/model"
	"cyberbasic/compiler/vm"
)

const levelObjectIDBase = 100000

type levelRuntime struct {
	model       *model.Model
	objectIDs   []int
	textureIDs  []int
	materialIDs []int
	lightIDs    []int
	colliderIDs []string // physics body IDs for level collision
}

var (
	levels   = make(map[int]*levelRuntime)
	levelsMu sync.Mutex
)

func registerLevel(v *vm.VM) {
	v.RegisterForeign("LoadLevel", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("LoadLevel(id, path) requires 2 arguments")
		}
		id := toInt(args[0])
		path := toString(args[1])
		m, err := model.Load(path)
		if err != nil {
			return nil, fmt.Errorf("LoadLevel: %w", err)
		}
		basePath := filepath.Dir(path)
		objectIDBase := id * levelObjectIDBase
		res, err := BuildModel(m, objectIDBase, basePath)
		if err != nil {
			return nil, fmt.Errorf("LoadLevel build: %w", err)
		}
		for _, oid := range res.ObjectIDs {
			RegisterObjectModel(oid, m)
		}
		levelsMu.Lock()
		levels[id] = &levelRuntime{
			model:       m,
			objectIDs:   res.ObjectIDs,
			textureIDs:  res.TextureIDs,
			materialIDs: res.MaterialIDs,
			lightIDs:    res.LightIDs,
			colliderIDs: nil, // populated by LoadLevelCollision
		}
		levelsMu.Unlock()
		return nil, nil
	})

	v.RegisterForeign("LoadLevelWithHierarchy", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("LoadLevelWithHierarchy(id, path) requires 2 arguments")
		}
		id := toInt(args[0])
		path := toString(args[1])
		m, err := model.Load(path)
		if err != nil {
			return nil, fmt.Errorf("LoadLevelWithHierarchy: %w", err)
		}
		basePath := filepath.Dir(path)
		objectIDBase := id * levelObjectIDBase
		res, err := BuildModelWithHierarchy(m, objectIDBase, basePath)
		if err != nil {
			return nil, fmt.Errorf("LoadLevelWithHierarchy build: %w", err)
		}
		for _, oid := range res.ObjectIDs {
			RegisterObjectModel(oid, m)
		}
		levelsMu.Lock()
		levels[id] = &levelRuntime{
			model:       m,
			objectIDs:   res.ObjectIDs,
			textureIDs:  res.TextureIDs,
			materialIDs: res.MaterialIDs,
			lightIDs:    res.LightIDs,
			colliderIDs: nil,
		}
		levelsMu.Unlock()
		return nil, nil
	})

	v.RegisterForeign("DrawLevel", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("DrawLevel(id) requires 1 argument")
		}
		id := toInt(args[0])
		levelsMu.Lock()
		lr, ok := levels[id]
		levelsMu.Unlock()
		if !ok {
			return nil, fmt.Errorf("LoadLevel: unknown level id %d", id)
		}
		for _, objID := range lr.objectIDs {
			_, _ = v.CallForeign("DrawObject", []interface{}{objID})
		}
		return nil, nil
	})

	v.RegisterForeign("UnloadLevel", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("UnloadLevel(id) requires 1 argument")
		}
		id := toInt(args[0])
		levelsMu.Lock()
		lr, ok := levels[id]
		delete(levels, id)
		levelsMu.Unlock()
		if !ok {
			return nil, nil
		}
		for _, objID := range lr.objectIDs {
			UnregisterObjectModel(objID)
			_, _ = v.CallForeign("DeleteObject", []interface{}{objID})
		}
		for _, texID := range lr.textureIDs {
			_, _ = v.CallForeign("DeleteTexture", []interface{}{texID})
		}
		for _, lightID := range lr.lightIDs {
			_, _ = v.CallForeign("DeleteLight", []interface{}{lightID})
		}
		for _, bid := range lr.colliderIDs {
			_, _ = v.CallForeign("DestroyBody3D", []interface{}{defaultPhysicsWorld3D, bid})
		}
		return nil, nil
	})

	v.RegisterForeign("LoadLevelCollision", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("LoadLevelCollision(id) requires 1 argument")
		}
		id := toInt(args[0])
		levelsMu.Lock()
		lr, ok := levels[id]
		levelsMu.Unlock()
		if !ok {
			return 0, nil
		}
		if lr.model == nil || len(lr.model.Colliders) == 0 {
			return 0, nil
		}
		// Clear any existing colliders (reload)
		for _, bid := range lr.colliderIDs {
			_, _ = v.CallForeign("DestroyBody3D", []interface{}{defaultPhysicsWorld3D, bid})
		}
		lr.colliderIDs = nil
		_, _ = v.CallForeign("PhysicsEnable", nil)
		for i, col := range lr.model.Colliders {
			bodyId := fmt.Sprintf("level_%d_col_%d", id, i)
			x, y, z := float64(col.Transform.X), float64(col.Transform.Y), float64(col.Transform.Z)
			mass := 0.0
			switch col.Type {
			case model.ColliderBox:
				sx, sy, sz := float64(col.SizeX), float64(col.SizeY), float64(col.SizeZ)
				if sx < 0.01 {
					sx = 0.5
				}
				if sy < 0.01 {
					sy = 0.5
				}
				if sz < 0.01 {
					sz = 0.5
				}
				_, err := v.CallForeign("CreateBox3D", []interface{}{
					defaultPhysicsWorld3D, bodyId, x, y, z, sx, sy, sz, mass,
				})
				if err == nil {
					lr.colliderIDs = append(lr.colliderIDs, bodyId)
				}
			case model.ColliderSphere:
				r := float64(col.Radius)
				if r < 0.01 {
					r = 0.5
				}
				_, err := v.CallForeign("CreateSphere3D", []interface{}{
					defaultPhysicsWorld3D, bodyId, x, y, z, r, mass,
				})
				if err == nil {
					lr.colliderIDs = append(lr.colliderIDs, bodyId)
				}
			case model.ColliderCapsule:
				r, h := float64(col.Radius), float64(col.Height)
				if r < 0.01 {
					r = 0.5
				}
				if h < 0.01 {
					h = 1
				}
				_, err := v.CallForeign("CreateCapsule3D", []interface{}{
					defaultPhysicsWorld3D, bodyId, x, y, z, r, h, mass,
				})
				if err == nil {
					lr.colliderIDs = append(lr.colliderIDs, bodyId)
				}
			default:
				// ColliderMesh or unknown: use box from size
				sx, sy, sz := float64(col.SizeX), float64(col.SizeY), float64(col.SizeZ)
				if sx < 0.01 {
					sx = 0.5
				}
				if sy < 0.01 {
					sy = 0.5
				}
				if sz < 0.01 {
					sz = 0.5
				}
				_, err := v.CallForeign("CreateBox3D", []interface{}{
					defaultPhysicsWorld3D, bodyId, x, y, z, sx, sy, sz, mass,
				})
				if err == nil {
					lr.colliderIDs = append(lr.colliderIDs, bodyId)
				}
			}
		}
		return len(lr.colliderIDs), nil
	})

	v.RegisterForeign("GetLevelColliderCount", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("GetLevelColliderCount(id) requires 1 argument")
		}
		id := toInt(args[0])
		levelsMu.Lock()
		lr, ok := levels[id]
		levelsMu.Unlock()
		if !ok {
			return 0, nil
		}
		return len(lr.colliderIDs), nil
	})

	v.RegisterForeign("GetLevelCollider", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("GetLevelCollider(id, index) requires 2 arguments")
		}
		id := toInt(args[0])
		index := toInt(args[1])
		levelsMu.Lock()
		lr, ok := levels[id]
		levelsMu.Unlock()
		if !ok || index < 0 || index >= len(lr.colliderIDs) {
			return "", nil
		}
		return lr.colliderIDs[index], nil
	})

	v.RegisterForeign("GetLevelObjectCount", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("GetLevelObjectCount(id) requires 1 argument")
		}
		id := toInt(args[0])
		levelsMu.Lock()
		lr, ok := levels[id]
		levelsMu.Unlock()
		if !ok {
			return 0, nil
		}
		return len(lr.objectIDs), nil
	})

	v.RegisterForeign("GetLevelObject", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("GetLevelObject(id, index, outVar) requires 3 arguments")
		}
		id := toInt(args[0])
		index := toInt(args[1])
		levelsMu.Lock()
		lr, ok := levels[id]
		levelsMu.Unlock()
		if !ok || index < 0 || index >= len(lr.objectIDs) {
			return nil, nil
		}
		return lr.objectIDs[index], nil
	})
}
