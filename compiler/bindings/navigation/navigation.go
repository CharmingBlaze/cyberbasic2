// Package navigation provides NavGrid, NavMesh, and NavAgent stubs for pathfinding.
package navigation

import (
	"fmt"

	"cyberbasic/compiler/vm"
)

func toString(v interface{}) string {
	if v == nil {
		return ""
	}
	return fmt.Sprint(v)
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

// RegisterNavigation registers NavGrid, NavMesh, NavAgent commands (stubs).
func RegisterNavigation(v *vm.VM) {
	v.RegisterForeign("NavGridCreate", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("NavGridCreate requires (width, height)")
		}
		return fmt.Sprintf("navgrid_%d", 0), nil
	})
	v.RegisterForeign("NavGridSetWalkable", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("NavGridSetWalkable requires (gridId, x, y, flag)")
		}
		return nil, nil
	})
	v.RegisterForeign("NavGridSetCost", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("NavGridSetCost requires (gridId, x, y, cost)")
		}
		return nil, nil
	})
	v.RegisterForeign("NavGridFindPath", func(args []interface{}) (interface{}, error) {
		if len(args) < 5 {
			return nil, fmt.Errorf("NavGridFindPath requires (gridId, startX, startY, endX, endY)")
		}
		return []interface{}{}, nil
	})
	v.RegisterForeign("NavMeshCreateFromTerrain", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("NavMeshCreateFromTerrain requires (terrainId)")
		}
		return fmt.Sprintf("navmesh_%d", 0), nil
	})
	v.RegisterForeign("NavMeshAddObstacle", func(args []interface{}) (interface{}, error) { return nil, nil })
	v.RegisterForeign("NavMeshRemoveObstacle", func(args []interface{}) (interface{}, error) { return nil, nil })
	v.RegisterForeign("NavMeshFindPath", func(args []interface{}) (interface{}, error) {
		if len(args) < 7 {
			return nil, fmt.Errorf("NavMeshFindPath requires (meshId, ox, oy, oz, dx, dy, dz)")
		}
		return []interface{}{}, nil
	})
	v.RegisterForeign("NavAgentCreate", func(args []interface{}) (interface{}, error) {
		return fmt.Sprintf("navagent_%d", 0), nil
	})
	v.RegisterForeign("NavAgentSetSpeed", func(args []interface{}) (interface{}, error) { return nil, nil })
	v.RegisterForeign("NavAgentSetRadius", func(args []interface{}) (interface{}, error) { return nil, nil })
	v.RegisterForeign("NavAgentSetDestination", func(args []interface{}) (interface{}, error) { return nil, nil })
	v.RegisterForeign("NavAgentGetNextWaypoint", func(args []interface{}) (interface{}, error) {
		return []interface{}{0.0, 0.0, 0.0}, nil
	})
}
