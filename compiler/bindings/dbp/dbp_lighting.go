// Package dbp: Lighting - light registry (raylib has no built-in dynamic lights).
//
// Lights are stored with position, color, intensity, range. Visual effect
// requires a custom shader; for now this provides the registry API.
//
// Commands:
//   - MakeLight(id, type): Create light (type: 0=point, 1=directional, 2=spot)
//   - PositionLight(id, x, y, z): Set position
//   - RotateLight(id, pitch, yaw, roll): Set direction (for directional/spot)
//   - SetLightColor(id, r, g, b): Light color
//   - SetLightIntensity(id, value): Intensity multiplier
//   - SetLightRange(id, value): Range/distance
//   - DeleteLight(id): Remove light
//   - SyncLight(id): Sync for multiplayer (placeholder)
package dbp

import (
	"fmt"
	"sync"

	"cyberbasic/compiler/vm"
)

type dbpLight struct {
	lightType  int
	x, y, z    float32
	pitch, yaw float32
	roll       float32
	r, g, b    float32
	intensity  float32
	range_     float32
	angle      float32 // cone angle for spot lights (degrees)
	shadows    bool    // shadow casting enabled
}

var (
	lights   = make(map[int]*dbpLight)
	lightsMu sync.Mutex
)

// registerLighting adds MakeLight, PositionLight, RotateLight, SetLightColor, SetLightIntensity, SetLightRange, DeleteLight, SyncLight.
func registerLighting(v *vm.VM) {
	v.RegisterForeign("MakeLight", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("MakeLight(id, type) requires 2 arguments")
		}
		id := toInt(args[0])
		typ := toInt(args[1])
		lightsMu.Lock()
		lights[id] = &dbpLight{lightType: typ, r: 1, g: 1, b: 1, intensity: 1, range_: 10}
		lightsMu.Unlock()
		return nil, nil
	})

	v.RegisterForeign("PositionLight", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("PositionLight(id, x, y, z) requires 4 arguments")
		}
		id := toInt(args[0])
		x, y, z := toFloat32(args[1]), toFloat32(args[2]), toFloat32(args[3])
		lightsMu.Lock()
		if l, ok := lights[id]; ok {
			l.x, l.y, l.z = x, y, z
		}
		lightsMu.Unlock()
		return nil, nil
	})

	v.RegisterForeign("RotateLight", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("RotateLight(id, pitch, yaw, roll) requires 4 arguments")
		}
		id := toInt(args[0])
		p, y, r := toFloat32(args[1]), toFloat32(args[2]), toFloat32(args[3])
		lightsMu.Lock()
		if l, ok := lights[id]; ok {
			l.pitch, l.yaw, l.roll = p, y, r
		}
		lightsMu.Unlock()
		return nil, nil
	})

	v.RegisterForeign("SetLightColor", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("SetLightColor(id, r, g, b) requires 4 arguments")
		}
		id := toInt(args[0])
		r, g, b := toFloat32(args[1]), toFloat32(args[2]), toFloat32(args[3])
		if r > 1 {
			r /= 255
		}
		if g > 1 {
			g /= 255
		}
		if b > 1 {
			b /= 255
		}
		lightsMu.Lock()
		if l, ok := lights[id]; ok {
			l.r, l.g, l.b = r, g, b
		}
		lightsMu.Unlock()
		return nil, nil
	})

	v.RegisterForeign("SetLightIntensity", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("SetLightIntensity(id, value) requires 2 arguments")
		}
		id := toInt(args[0])
		val := toFloat32(args[1])
		lightsMu.Lock()
		if l, ok := lights[id]; ok {
			l.intensity = val
		}
		lightsMu.Unlock()
		return nil, nil
	})

	v.RegisterForeign("SetLightRange", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("SetLightRange(id, value) requires 2 arguments")
		}
		id := toInt(args[0])
		val := toFloat32(args[1])
		lightsMu.Lock()
		if l, ok := lights[id]; ok {
			l.range_ = val
		}
		lightsMu.Unlock()
		return nil, nil
	})

	v.RegisterForeign("SetLightAngle", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("SetLightAngle(id, degrees) requires 2 arguments")
		}
		id := toInt(args[0])
		deg := toFloat32(args[1])
		lightsMu.Lock()
		if l, ok := lights[id]; ok {
			l.angle = deg
		}
		lightsMu.Unlock()
		return nil, nil
	})

	// EnableLightShadows(id): Enables shadow casting for a light. DBP-style; use EnableShadows() for global.
	v.RegisterForeign("EnableLightShadows", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("EnableLightShadows(id) requires 1 argument")
		}
		id := toInt(args[0])
		lightsMu.Lock()
		if l, ok := lights[id]; ok {
			l.shadows = true
		}
		lightsMu.Unlock()
		return nil, nil
	})
	// DisableLightShadows(id): Disables shadow casting for a light.
	v.RegisterForeign("DisableLightShadows", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("DisableLightShadows(id) requires 1 argument")
		}
		id := toInt(args[0])
		lightsMu.Lock()
		if l, ok := lights[id]; ok {
			l.shadows = false
		}
		lightsMu.Unlock()
		return nil, nil
	})
	// EnableShadows(id) / DisableShadows(id): DBP aliases for per-light shadow control.
	v.RegisterForeign("EnableShadows", func(args []interface{}) (interface{}, error) {
		if len(args) >= 1 {
			return v.CallForeign("EnableLightShadows", args)
		}
		return nil, fmt.Errorf("EnableShadows(id) requires 1 argument")
	})
	v.RegisterForeign("DisableShadows", func(args []interface{}) (interface{}, error) {
		if len(args) >= 1 {
			return v.CallForeign("DisableLightShadows", args)
		}
		return nil, fmt.Errorf("DisableShadows(id) requires 1 argument")
	})

	v.RegisterForeign("DeleteLight", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("DeleteLight(id) requires 1 argument")
		}
		id := toInt(args[0])
		lightsMu.Lock()
		delete(lights, id)
		lightsMu.Unlock()
		return nil, nil
	})

	v.RegisterForeign("SyncLight", func(args []interface{}) (interface{}, error) {
		return nil, nil
	})
	v.RegisterForeign("LightExists", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("LightExists(id) requires 1 argument")
		}
		id := toInt(args[0])
		lightsMu.Lock()
		_, ok := lights[id]
		lightsMu.Unlock()
		if ok {
			return 1, nil
		}
		return 0, nil
	})

	// Light queries
	v.RegisterForeign("GetLightX", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return 0.0, nil
		}
		id := toInt(args[0])
		lightsMu.Lock()
		l, ok := lights[id]
		lightsMu.Unlock()
		if !ok {
			return 0.0, nil
		}
		return float64(l.x), nil
	})
	v.RegisterForeign("GetLightY", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return 0.0, nil
		}
		id := toInt(args[0])
		lightsMu.Lock()
		l, ok := lights[id]
		lightsMu.Unlock()
		if !ok {
			return 0.0, nil
		}
		return float64(l.y), nil
	})
	v.RegisterForeign("GetLightZ", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return 0.0, nil
		}
		id := toInt(args[0])
		lightsMu.Lock()
		l, ok := lights[id]
		lightsMu.Unlock()
		if !ok {
			return 0.0, nil
		}
		return float64(l.z), nil
	})
	v.RegisterForeign("GetLightColorR", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return 0.0, nil
		}
		id := toInt(args[0])
		lightsMu.Lock()
		l, ok := lights[id]
		lightsMu.Unlock()
		if !ok {
			return 0.0, nil
		}
		return float64(l.r * 255), nil
	})
	v.RegisterForeign("GetLightColorG", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return 0.0, nil
		}
		id := toInt(args[0])
		lightsMu.Lock()
		l, ok := lights[id]
		lightsMu.Unlock()
		if !ok {
			return 0.0, nil
		}
		return float64(l.g * 255), nil
	})
	v.RegisterForeign("GetLightColorB", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return 0.0, nil
		}
		id := toInt(args[0])
		lightsMu.Lock()
		l, ok := lights[id]
		lightsMu.Unlock()
		if !ok {
			return 0.0, nil
		}
		return float64(l.b * 255), nil
	})
}
