// Package dbp: Pathfinding DBP wrappers - NavMeshLoad, NavMeshFindPath, NavMeshDraw.
package dbp

import (
	"fmt"
	"sync"

	"cyberbasic/compiler/vm"
)

var (
	navMeshMap   = make(map[int]string)
	navMeshSeq   int
	navMeshMapMu sync.Mutex
)

// registerNav adds NavMeshLoad, NavMeshFindPath, NavMeshDraw.
func registerNav(v *vm.VM) {
	v.RegisterForeign("NavMeshLoad", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("NavMeshLoad(id, path) requires 2 arguments")
		}
		id := toInt(args[0])
		path := toString(args[1])
		res, err := v.CallForeign("NavMeshLoadFromFile", []interface{}{path})
		if err != nil {
			return nil, err
		}
		meshId, _ := res.(string)
		if meshId == "" {
			return nil, nil
		}
		navMeshMapMu.Lock()
		navMeshMap[id] = meshId
		navMeshMapMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("NavMeshFindPath", func(args []interface{}) (interface{}, error) {
		if len(args) < 7 {
			return nil, fmt.Errorf("NavMeshFindPath(id, startX, startY, startZ, endX, endY, endZ) requires 7 arguments")
		}
		id := toInt(args[0])
		navMeshMapMu.Lock()
		meshId, ok := navMeshMap[id]
		navMeshMapMu.Unlock()
		if !ok {
			return []interface{}{}, nil
		}
		return v.CallForeign("NavMeshFindPathRaw", []interface{}{
			meshId,
			toFloat64(args[1]), toFloat64(args[2]), toFloat64(args[3]),
			toFloat64(args[4]), toFloat64(args[5]), toFloat64(args[6]),
		})
	})
	v.RegisterForeign("NavMeshDraw", func(args []interface{}) (interface{}, error) {
		// Stub: no debug draw
		return nil, nil
	})
}
