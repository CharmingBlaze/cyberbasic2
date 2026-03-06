// Package dbp: Prefab loading - LoadPrefab, SpawnPrefab.
package dbp

import (
	"fmt"
	"path/filepath"
	"sync"

	"cyberbasic/compiler/bindings/model"
	"cyberbasic/compiler/vm"
)

const prefabObjectIDBase = 500000

type prefabRuntime struct {
	model    *model.Model
	basePath string
}

var (
	prefabs           = make(map[int]*prefabRuntime)
	prefabsMu         sync.Mutex
	prefabSpawnCounter int
)

func registerPrefab(v *vm.VM) {
	v.RegisterForeign("LoadPrefab", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("LoadPrefab(id, path) requires 2 arguments")
		}
		id := toInt(args[0])
		path := toString(args[1])
		m, err := model.Load(path)
		if err != nil {
			return nil, fmt.Errorf("LoadPrefab: %w", err)
		}
		basePath := filepath.Dir(path)
		prefabsMu.Lock()
		prefabs[id] = &prefabRuntime{model: m, basePath: basePath}
		prefabsMu.Unlock()
		return nil, nil
	})

	v.RegisterForeign("SpawnPrefab", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("SpawnPrefab(id, x, y, z) requires 4 arguments")
		}
		id := toInt(args[0])
		x := toFloat32(args[1])
		y := toFloat32(args[2])
		z := toFloat32(args[3])
		prefabsMu.Lock()
		pr, ok := prefabs[id]
		prefabsMu.Unlock()
		if !ok {
			return 0, nil
		}
		prefabsMu.Lock()
		prefabSpawnCounter++
		objectIDBase := prefabObjectIDBase + prefabSpawnCounter*10000
		prefabsMu.Unlock()
		res, err := BuildModel(pr.model, objectIDBase, pr.basePath)
		if err != nil {
			return 0, fmt.Errorf("SpawnPrefab build: %w", err)
		}
		if len(res.ObjectIDs) == 0 {
			return 0, nil
		}
		// Offset all spawned objects by spawn position and register for IK
		for _, objID := range res.ObjectIDs {
			RegisterObjectModel(objID, pr.model)
			objectsMu.Lock()
			obj, ok := objects[objID]
			objectsMu.Unlock()
			if ok {
				obj.x += x
				obj.y += y
				obj.z += z
			}
		}
		return res.ObjectIDs[0], nil
	})
}
