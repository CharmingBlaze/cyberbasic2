// Package dbp: Instancing - MakeInstance, PositionInstance, DrawInstances.
package dbp

import (
	"fmt"
	"math"
	"sync"

	"cyberbasic/compiler/vm"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type dbpInstance struct {
	x, y, z     float32
	pitch, yaw  float32
	roll        float32
	scaleX      float32
	scaleY      float32
	scaleZ      float32
}

type instanceInfo struct {
	baseID int
	idx    int
}

var (
	instances       = make(map[int][]*dbpInstance) // baseID -> instances
	instanceToBase   = make(map[int]instanceInfo)   // instanceID -> (baseID, idx)
	instancesMu     sync.Mutex
)

// registerInstancing adds MakeInstance, PositionInstance, DrawInstances.
func registerInstancing(v *vm.VM) {
	v.RegisterForeign("MakeInstance", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("MakeInstance(baseID, instanceID) requires 2 arguments")
		}
		baseID := toInt(args[0])
		instanceID := toInt(args[1])
		instancesMu.Lock()
		if _, ok := instances[baseID]; !ok {
			instances[baseID] = make([]*dbpInstance, 0)
		}
		idx := len(instances[baseID])
		instances[baseID] = append(instances[baseID], &dbpInstance{
			scaleX: 1, scaleY: 1, scaleZ: 1,
		})
		instanceToBase[instanceID] = instanceInfo{baseID: baseID, idx: idx}
		instancesMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("PositionInstance", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("PositionInstance(instanceID, x, y, z) requires 4 arguments")
		}
		instanceID := toInt(args[0])
		x, y, z := toFloat32(args[1]), toFloat32(args[2]), toFloat32(args[3])
		instancesMu.Lock()
		info, ok := instanceToBase[instanceID]
		if !ok {
			instancesMu.Unlock()
			return nil, nil
		}
		baseID, idx := info.baseID, info.idx
		list := instances[baseID]
		if idx >= 0 && idx < len(list) && list[idx] != nil {
			list[idx].x, list[idx].y, list[idx].z = x, y, z
		}
		instancesMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("DeleteInstance", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("DeleteInstance(instanceID) requires 1 argument")
		}
		instanceID := toInt(args[0])
		instancesMu.Lock()
		info, ok := instanceToBase[instanceID]
		if ok {
			delete(instanceToBase, instanceID)
			if list, exists := instances[info.baseID]; exists && info.idx >= 0 && info.idx < len(list) {
				list[info.idx] = nil
			}
		}
		instancesMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("DeleteAllInstances", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("DeleteAllInstances(baseID) requires 1 argument")
		}
		baseID := toInt(args[0])
		instancesMu.Lock()
		var toDelete []int
		for instanceID, info := range instanceToBase {
			if info.baseID == baseID {
				toDelete = append(toDelete, instanceID)
			}
		}
		for _, id := range toDelete {
			delete(instanceToBase, id)
		}
		instances[baseID] = make([]*dbpInstance, 0)
		instancesMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("InstanceExists", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("InstanceExists(instanceID) requires 1 argument")
		}
		instanceID := toInt(args[0])
		instancesMu.Lock()
		_, ok := instanceToBase[instanceID]
		instancesMu.Unlock()
		if ok {
			return 1, nil
		}
		return 0, nil
	})
	v.RegisterForeign("DrawInstances", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("DrawInstances(baseID) requires 1 argument")
		}
		baseID := toInt(args[0])
		instancesMu.Lock()
		list := instances[baseID]
		instancesMu.Unlock()
		if list == nil || len(list) == 0 {
			return nil, nil
		}
		objectsMu.Lock()
		obj, ok := objects[baseID]
		objectsMu.Unlock()
		if !ok || !obj.visible {
			return nil, nil
		}
		tint := rl.NewColor(obj.colorR, obj.colorG, obj.colorB, obj.colorA)
		instancesMu.Lock()
		for _, inst := range list {
			if inst == nil {
				continue
			}
			pos := rl.Vector3{X: inst.x, Y: inst.y, Z: inst.z}
			rotAxis := rl.Vector3{X: 0, Y: 1, Z: 0}
			rotAngle := inst.yaw * math.Pi / 180
			scale := rl.Vector3{X: inst.scaleX, Y: inst.scaleY, Z: inst.scaleZ}
			if obj.wireframe {
				rl.DrawModelWiresEx(obj.model, pos, rotAxis, rotAngle, scale, tint)
			} else {
				rl.DrawModelEx(obj.model, pos, rotAxis, rotAngle, scale, tint)
			}
		}
		instancesMu.Unlock()
		return nil, nil
	})
}
