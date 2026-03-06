// Package dbp: Object groups - batch operations on multiple objects.
//
// Groups are useful for buildings, squads, or modular levels.
// All objects in a group share the same transform operations.
//
// Commands:
//   - MakeGroup(id): Create empty group
//   - AddToGroup(groupID, objectID): Add object to group
//   - RemoveFromGroup(groupID, objectID): Remove object from group
//   - PositionGroup(groupID, x, y, z): Set position of all objects in group
//   - RotateGroup(groupID, pitch, yaw, roll): Set rotation of all objects
//   - DrawGroup(groupID): Draw all objects in group
//   - SyncGroup(groupID): Sync group state for multiplayer (placeholder)
package dbp

import (
	"fmt"
	"math"
	"sync"

	"cyberbasic/compiler/vm"
	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
	groups   = make(map[int][]int) // groupID -> objectIDs
	groupsMu sync.Mutex
)

// registerGroups adds MakeGroup, AddToGroup, RemoveFromGroup, PositionGroup, RotateGroup, DrawGroup, SyncGroup.
func registerGroups(v *vm.VM) {
	v.RegisterForeign("MakeGroup", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("MakeGroup(id) requires 1 argument")
		}
		id := toInt(args[0])
		groupsMu.Lock()
		groups[id] = []int{}
		groupsMu.Unlock()
		return nil, nil
	})

	v.RegisterForeign("AddToGroup", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("AddToGroup(groupID, objectID) requires 2 arguments")
		}
		gid := toInt(args[0])
		oid := toInt(args[1])
		groupsMu.Lock()
		groups[gid] = append(groups[gid], oid)
		groupsMu.Unlock()
		return nil, nil
	})

	v.RegisterForeign("RemoveFromGroup", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("RemoveFromGroup(groupID, objectID) requires 2 arguments")
		}
		gid := toInt(args[0])
		oid := toInt(args[1])
		groupsMu.Lock()
		list := groups[gid]
		newList := make([]int, 0, len(list))
		for _, id := range list {
			if id != oid {
				newList = append(newList, id)
			}
		}
		groups[gid] = newList
		groupsMu.Unlock()
		return nil, nil
	})

	v.RegisterForeign("PositionGroup", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("PositionGroup(groupID, x, y, z) requires 4 arguments")
		}
		gid := toInt(args[0])
		x, y, z := toFloat32(args[1]), toFloat32(args[2]), toFloat32(args[3])
		groupsMu.Lock()
		list := groups[gid]
		groupsMu.Unlock()
		objectsMu.Lock()
		for _, oid := range list {
			if obj, ok := objects[oid]; ok {
				obj.x, obj.y, obj.z = x, y, z
			}
		}
		objectsMu.Unlock()
		return nil, nil
	})

	v.RegisterForeign("RotateGroup", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("RotateGroup(groupID, pitch, yaw, roll) requires 4 arguments")
		}
		gid := toInt(args[0])
		p, y, r := toFloat32(args[1]), toFloat32(args[2]), toFloat32(args[3])
		groupsMu.Lock()
		list := groups[gid]
		groupsMu.Unlock()
		objectsMu.Lock()
		for _, oid := range list {
			if obj, ok := objects[oid]; ok {
				obj.pitch, obj.yaw, obj.roll = p, y, r
			}
		}
		objectsMu.Unlock()
		return nil, nil
	})

	v.RegisterForeign("DrawGroup", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("DrawGroup(groupID) requires 1 argument")
		}
		gid := toInt(args[0])
		groupsMu.Lock()
		list := groups[gid]
		groupsMu.Unlock()
		objectsMu.Lock()
		for _, oid := range list {
			if obj, ok := objects[oid]; ok && obj.visible {
				pos := rl.Vector3{X: obj.x, Y: obj.y, Z: obj.z}
				rotAxis := rl.Vector3{X: 0, Y: 1, Z: 0}
				rotAngle := obj.yaw * float32(math.Pi) / 180
				scale := rl.Vector3{X: obj.scaleX, Y: obj.scaleY, Z: obj.scaleZ}
				tint := rl.NewColor(obj.colorR, obj.colorG, obj.colorB, obj.colorA)
				if obj.wireframe {
					rl.DrawModelWiresEx(obj.model, pos, rotAxis, rotAngle, scale, tint)
				} else {
					rl.DrawModelEx(obj.model, pos, rotAxis, rotAngle, scale, tint)
				}
			}
		}
		objectsMu.Unlock()
		return nil, nil
	})

	v.RegisterForeign("SyncGroup", func(args []interface{}) (interface{}, error) {
		// Placeholder for multiplayer sync - would push group state to net layer
		return nil, nil
	})
}
