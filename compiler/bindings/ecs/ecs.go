// Package ecs provides an Entity-Component System binding for CyberBasic.
// BASIC can call ECS.CreateWorld, ECS.CreateEntity, ECS.AddComponent, ECS.Query, etc.
// Implemented as a minimal in-memory ECS; can be replaced by github.com/mlange-42/arche/ecs if desired.
package ecs

import (
	"fmt"
	"strconv"
	"sync"

	"cyberbasic/compiler/vm"
)

func toFloat64(v interface{}) float64 {
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

func toInt(v interface{}) int {
	switch x := v.(type) {
	case int:
		return x
	case float64:
		return int(x)
	case string:
		n, _ := strconv.Atoi(x)
		return n
	default:
		return 0
	}
}

func toString(v interface{}) string {
	if v == nil {
		return ""
	}
	return fmt.Sprint(v)
}

// Built-in component data (stored as interface{} per type name)
type compTransform struct{ X, Y, Z float64 }
type compSprite struct {
	TextureId string
	Visible   bool
}
type compHealth struct{ Current, Max float64 }

var (
	worldsMu sync.RWMutex
	worlds   = make(map[string]*worldState)
	worldSeq int
)

type worldState struct {
	mu        sync.RWMutex
	entities  []string
	nextEid   int
	comps     map[string]map[string]interface{} // entityId -> componentType -> data
	compNames map[string]bool                    // allowed component type names
}

func newWorldState() *worldState {
	ws := &worldState{
		entities:  []string{},
		comps:     make(map[string]map[string]interface{}),
		compNames: make(map[string]bool),
	}
	for _, n := range []string{"Transform", "Sprite", "Health"} {
		ws.compNames[n] = true
	}
	return ws
}

// RegisterECS registers all ECS.* foreign functions with the VM.
func RegisterECS(v *vm.VM) {
	v.RegisterForeign("ECS.CreateWorld", func(args []interface{}) (interface{}, error) {
		worldsMu.Lock()
		worldSeq++
		id := fmt.Sprintf("w%d", worldSeq)
		worlds[id] = newWorldState()
		worldsMu.Unlock()
		return id, nil
	})

	v.RegisterForeign("ECS.DestroyWorld", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("ECS.DestroyWorld requires (worldId)")
		}
		id := toString(args[0])
		worldsMu.Lock()
		delete(worlds, id)
		worldsMu.Unlock()
		return nil, nil
	})

	v.RegisterForeign("ECS.CreateEntity", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("ECS.CreateEntity requires (worldId)")
		}
		wid := toString(args[0])
		worldsMu.RLock()
		ws, ok := worlds[wid]
		worldsMu.RUnlock()
		if !ok {
			return nil, fmt.Errorf("unknown world: %s", wid)
		}
		ws.mu.Lock()
		ws.nextEid++
		eid := fmt.Sprintf("e%d", ws.nextEid)
		ws.entities = append(ws.entities, eid)
		ws.comps[eid] = make(map[string]interface{})
		ws.mu.Unlock()
		return eid, nil
	})

	v.RegisterForeign("ECS.DestroyEntity", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("ECS.DestroyEntity requires (worldId, entityId)")
		}
		wid := toString(args[0])
		eid := toString(args[1])
		worldsMu.RLock()
		ws, ok := worlds[wid]
		worldsMu.RUnlock()
		if !ok {
			return nil, fmt.Errorf("unknown world: %s", wid)
		}
		ws.mu.Lock()
		delete(ws.comps, eid)
		for i, id := range ws.entities {
			if id == eid {
				ws.entities = append(ws.entities[:i], ws.entities[i+1:]...)
				break
			}
		}
		ws.mu.Unlock()
		return nil, nil
	})

	v.RegisterForeign("ECS.AddComponent", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("ECS.AddComponent requires (worldId, entityId, componentType, ...args)")
		}
		wid := toString(args[0])
		eid := toString(args[1])
		ctype := toString(args[2])
		worldsMu.RLock()
		ws, ok := worlds[wid]
		worldsMu.RUnlock()
		if !ok {
			return nil, fmt.Errorf("unknown world: %s", wid)
		}
		if !ws.compNames[ctype] {
			return nil, fmt.Errorf("unknown component type: %s (use Transform, Sprite, Health)", ctype)
		}
		var data interface{}
		switch ctype {
		case "Transform":
			x, y, z := 0.0, 0.0, 0.0
			if len(args) >= 6 {
				x, y, z = toFloat64(args[3]), toFloat64(args[4]), toFloat64(args[5])
			}
			data = &compTransform{X: x, Y: y, Z: z}
		case "Sprite":
			texId := ""
			visible := true
			if len(args) >= 4 {
				texId = toString(args[3])
			}
			if len(args) >= 5 {
				visible = toInt(args[4]) != 0
			}
			data = &compSprite{TextureId: texId, Visible: visible}
		case "Health":
			cur, max := 100.0, 100.0
			if len(args) >= 5 {
				cur, max = toFloat64(args[3]), toFloat64(args[4])
			}
			data = &compHealth{Current: cur, Max: max}
		default:
			return nil, fmt.Errorf("unknown component type: %s", ctype)
		}
		ws.mu.Lock()
		if ws.comps[eid] == nil {
			ws.comps[eid] = make(map[string]interface{})
		}
		ws.comps[eid][ctype] = data
		ws.mu.Unlock()
		return nil, nil
	})

	v.RegisterForeign("ECS.HasComponent", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("ECS.HasComponent requires (worldId, entityId, componentType)")
		}
		wid := toString(args[0])
		eid := toString(args[1])
		ctype := toString(args[2])
		worldsMu.RLock()
		ws, ok := worlds[wid]
		worldsMu.RUnlock()
		if !ok {
			return nil, fmt.Errorf("unknown world: %s", wid)
		}
		ws.mu.RLock()
		_, has := ws.comps[eid][ctype]
		ws.mu.RUnlock()
		return has, nil
	})

	v.RegisterForeign("ECS.RemoveComponent", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("ECS.RemoveComponent requires (worldId, entityId, componentType)")
		}
		wid := toString(args[0])
		eid := toString(args[1])
		ctype := toString(args[2])
		worldsMu.RLock()
		ws, ok := worlds[wid]
		worldsMu.RUnlock()
		if !ok {
			return nil, fmt.Errorf("unknown world: %s", wid)
		}
		ws.mu.Lock()
		delete(ws.comps[eid], ctype)
		ws.mu.Unlock()
		return nil, nil
	})

	// Getters for Transform
	v.RegisterForeign("ECS.GetTransformX", func(args []interface{}) (interface{}, error) {
		return getCompFloat(args, "Transform", func(c *compTransform) float64 { return c.X })
	})
	v.RegisterForeign("ECS.GetTransformY", func(args []interface{}) (interface{}, error) {
		return getCompFloat(args, "Transform", func(c *compTransform) float64 { return c.Y })
	})
	v.RegisterForeign("ECS.GetTransformZ", func(args []interface{}) (interface{}, error) {
		return getCompFloat(args, "Transform", func(c *compTransform) float64 { return c.Z })
	})
	v.RegisterForeign("ECS.SetTransform", func(args []interface{}) (interface{}, error) {
		if len(args) < 6 {
			return nil, fmt.Errorf("ECS.SetTransform requires (worldId, entityId, x, y, z)")
		}
		wid, eid := toString(args[0]), toString(args[1])
		x, y, z := toFloat64(args[2]), toFloat64(args[3]), toFloat64(args[4])
		worldsMu.RLock()
		ws, ok := worlds[wid]
		worldsMu.RUnlock()
		if !ok {
			return nil, fmt.Errorf("unknown world: %s", wid)
		}
		ws.mu.Lock()
		defer ws.mu.Unlock()
		if m, ok := ws.comps[eid]["Transform"].(*compTransform); ok {
			m.X, m.Y, m.Z = x, y, z
		}
		return nil, nil
	})

	// Getters for Health
	v.RegisterForeign("ECS.GetHealthCurrent", func(args []interface{}) (interface{}, error) {
		return getCompFloatHealth(args, func(c *compHealth) float64 { return c.Current })
	})
	v.RegisterForeign("ECS.GetHealthMax", func(args []interface{}) (interface{}, error) {
		return getCompFloatHealth(args, func(c *compHealth) float64 { return c.Max })
	})

	// Query: ECS.Query(worldId, componentType1, componentType2, ...) -> count then entity IDs
	// Returns number of entities that have ALL the given components, then we need a way to iterate.
	// Simpler: ECS.Query(worldId, componentType) returns a string "e1,e2,e3" so BASIC can Split or we add ECS.QueryCount and ECS.QueryEntity(worldId, componentType, index).
	v.RegisterForeign("ECS.QueryCount", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("ECS.QueryCount requires (worldId, componentType1, ...)")
		}
		wid := toString(args[0])
		worldsMu.RLock()
		ws, ok := worlds[wid]
		worldsMu.RUnlock()
		if !ok {
			return nil, fmt.Errorf("unknown world: %s", wid)
		}
		types := make([]string, 0, len(args)-1)
		for i := 1; i < len(args); i++ {
			types = append(types, toString(args[i]))
		}
		ws.mu.RLock()
		defer ws.mu.RUnlock()
		var count int
		for _, eid := range ws.entities {
			comps := ws.comps[eid]
			all := true
			for _, t := range types {
				if _, has := comps[t]; !has {
					all = false
					break
				}
			}
			if all {
				count++
			}
		}
		return count, nil
	})

	v.RegisterForeign("ECS.QueryEntity", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("ECS.QueryEntity requires (worldId, componentType, index)")
		}
		wid := toString(args[0])
		ctype := toString(args[1])
		idx := toInt(args[2])
		worldsMu.RLock()
		ws, ok := worlds[wid]
		worldsMu.RUnlock()
		if !ok {
			return nil, fmt.Errorf("unknown world: %s", wid)
		}
		ws.mu.RLock()
		defer ws.mu.RUnlock()
		var n int
		for _, eid := range ws.entities {
			if _, has := ws.comps[eid][ctype]; has {
				if n == idx {
					return eid, nil
				}
				n++
			}
		}
		return nil, nil
	})
}

func getCompFloat(args []interface{}, compType string, get func(*compTransform) float64) (interface{}, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("ECS getter requires (worldId, entityId)")
	}
	wid := toString(args[0])
	eid := toString(args[1])
	worldsMu.RLock()
	ws, ok := worlds[wid]
	worldsMu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("unknown world: %s", wid)
	}
	ws.mu.RLock()
	defer ws.mu.RUnlock()
	c, ok := ws.comps[eid][compType]
	if !ok {
		return 0.0, nil
	}
	if t, ok := c.(*compTransform); ok {
		return get(t), nil
	}
	return 0.0, nil
}

func getCompFloatHealth(args []interface{}, get func(*compHealth) float64) (interface{}, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("ECS getter requires (worldId, entityId)")
	}
	wid := toString(args[0])
	eid := toString(args[1])
	worldsMu.RLock()
	ws, ok := worlds[wid]
	worldsMu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("unknown world: %s", wid)
	}
	ws.mu.RLock()
	defer ws.mu.RUnlock()
	c, ok := ws.comps[eid]["Health"]
	if !ok {
		return 0.0, nil
	}
	if t, ok := c.(*compHealth); ok {
		return get(t), nil
	}
	return 0.0, nil
}

func init() {
	// Register Health getters (used above; we already have getCompFloat for Transform)
	// getCompFloatHealth is defined but we need to wire it. We already registered ECS.GetHealthCurrent and ECS.GetHealthMax with inline getters - fix: use a helper that works with compHealth.
	_ = getCompFloatHealth
}
