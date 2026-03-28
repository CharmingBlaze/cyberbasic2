// Package physics2d is a high-level 2D physics API (DotObject bodies) over box2d bindings.
package physics2d

import (
	"cyberbasic/compiler/errors"
	"cyberbasic/compiler/vm"
	"fmt"
	"strings"
	"sync"
)

const defaultWorld = "default"

var bodySeqMu sync.Mutex
var bodySeq int

// RegisterPhysics2DHigh registers PhysicsHighWorld, PhysicsHighDynamicBox, PhysicsHighStaticBox, PhysicsHighRaycast2D.
func RegisterPhysics2DHigh(v *vm.VM) {
	v.RegisterForeign("PhysicsHighWorld", func(args []interface{}) (interface{}, error) {
		gx, gy := 0.0, 9.8
		if len(args) >= 2 {
			gx, gy = toF(args[0]), toF(args[1])
		}
		if _, err := v.CallForeign("CreateWorld2D", []interface{}{defaultWorld, gx, gy}); err != nil {
			return nil, err
		}
		WorldEnsured = true
		return nil, nil
	})

	v.RegisterForeign("PhysicsHighDynamicBox", func(args []interface{}) (interface{}, error) {
		if err := ensureWorld(v); err != nil {
			return nil, err
		}
		if len(args) < 4 {
			return nil, fmt.Errorf("PhysicsHighDynamicBox requires (x, y, w, h)")
		}
		x, y, w, h := toF(args[0]), toF(args[1]), toF(args[2]), toF(args[3])
		bodySeqMu.Lock()
		bodySeq++
		id := fmt.Sprintf("ph_%d", bodySeq)
		bodySeqMu.Unlock()
		// CreateBody2D(worldId, bodyId, bodyType dynamic=2, shapeKind box=0, x, y, density, hx, hy)
		hx, hy := w/2, h/2
		if _, err := v.CallForeign("CreateBody2D", []interface{}{defaultWorld, id, 2, 0, x, y, 1.0, hx, hy}); err != nil {
			return nil, err
		}
		return &BodyDot{v: v, world: defaultWorld, id: id}, nil
	})

	v.RegisterForeign("PhysicsHighStaticBox", func(args []interface{}) (interface{}, error) {
		if err := ensureWorld(v); err != nil {
			return nil, err
		}
		if len(args) < 4 {
			return nil, fmt.Errorf("PhysicsHighStaticBox requires (x, y, w, h)")
		}
		x, y, w, h := toF(args[0]), toF(args[1]), toF(args[2]), toF(args[3])
		bodySeqMu.Lock()
		bodySeq++
		id := fmt.Sprintf("ph_%d", bodySeq)
		bodySeqMu.Unlock()
		hx, hy := w/2, h/2
		if _, err := v.CallForeign("CreateBody2D", []interface{}{defaultWorld, id, 0, 0, x, y, 0.0, hx, hy}); err != nil {
			return nil, err
		}
		return &BodyDot{v: v, world: defaultWorld, id: id}, nil
	})

	v.RegisterForeign("PhysicsHighRaycast2D", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("PhysicsHighRaycast2D requires (x1, y1, x2, y2)")
		}
		x1, y1, x2, y2 := toF(args[0]), toF(args[1]), toF(args[2]), toF(args[3])
		dx, dy := x2-x1, y2-y1
		if _, err := v.CallForeign("Physics2DRaycast", []interface{}{x1, y1, dx, dy}); err != nil {
			return nil, err
		}
		return nil, nil
	})

	v.SetGlobal("physics", &physicsModuleDot{v: v})
}

// physicsModuleDot is the v2 namespace PHYSICS.* (global key "physics").
type physicsModuleDot struct {
	v *vm.VM
}

func (p *physicsModuleDot) GetProp(path []string) (vm.Value, error) {
	if len(path) == 1 && strings.ToLower(path[0]) == "dynamic" {
		return &physicsDynamicDot{v: p.v}, nil
	}
	return nil, nil
}

func (p *physicsModuleDot) SetProp([]string, vm.Value) error {
	return fmt.Errorf("physics: namespace is not assignable")
}

func (p *physicsModuleDot) CallMethod(name string, args []vm.Value) (vm.Value, error) {
	ia := valuesToIface(args)
	switch strings.ToLower(name) {
	case "world":
		return p.v.CallForeign("PhysicsHighWorld", ia)
	case "dynamicbox":
		return p.v.CallForeign("PhysicsHighDynamicBox", ia)
	case "staticbox":
		return p.v.CallForeign("PhysicsHighStaticBox", ia)
	case "raycast2d":
		return p.v.CallForeign("PhysicsHighRaycast2D", ia)
	default:
		return nil, &errors.CyberError{
			Code:       errors.ErrDotAccess,
			Message:    fmt.Sprintf("unknown physics method %q", name),
			Suggestion: "Use world, dynamicbox, staticbox, raycast2d (legacy: PhysicsHighWorld, …)",
		}
	}
}

func valuesToIface(a []vm.Value) []interface{} {
	out := make([]interface{}, len(a))
	for i := range a {
		out[i] = a[i]
	}
	return out
}

func ensureWorld(v *vm.VM) error {
	if WorldEnsured {
		return nil
	}
	if RequireExplicitWorld {
		return &errors.CyberError{
			Code:       errors.ErrPhysicsBodyBeforeWorld,
			Message:    "Physics body created before a physics world exists.",
			Suggestion: "Call PhysicsHighWorld() before creating bodies, or use implicit mode (ON UPDATE without InitWindow).",
		}
	}
	if _, err := v.CallForeign("CreateWorld2D", []interface{}{defaultWorld, 0.0, 9.8}); err != nil {
		return err
	}
	WorldEnsured = true
	return nil
}

func toF(a interface{}) float64 {
	switch x := a.(type) {
	case float64:
		return x
	case int:
		return float64(x)
	case int32:
		return float64(x)
	case int64:
		return float64(x)
	default:
		return 0
	}
}

// BodyDot is a vm.DotObject for a 2D body (minimal property set).
type BodyDot struct {
	v     *vm.VM
	world string
	id    string
}

func anyToF(x interface{}) float64 {
	switch v := x.(type) {
	case float64:
		return v
	case float32:
		return float64(v)
	case int:
		return float64(v)
	default:
		return 0
	}
}

func (b *BodyDot) GetProp(path []string) (vm.Value, error) {
	if len(path) == 0 {
		return nil, fmt.Errorf("empty path")
	}
	switch strings.ToLower(path[0]) {
	case "x":
		rx, err := b.v.CallForeign("GetPositionX2D", []interface{}{b.world, b.id})
		if err != nil {
			return nil, err
		}
		return anyToF(rx), nil
	case "y":
		ry, err := b.v.CallForeign("GetPositionY2D", []interface{}{b.world, b.id})
		if err != nil {
			return nil, err
		}
		return anyToF(ry), nil
	case "id":
		return b.id, nil
	case "vx":
		rx, err := b.v.CallForeign("GetVelocityX2D", []interface{}{b.world, b.id})
		if err != nil {
			return nil, err
		}
		return anyToF(rx), nil
	case "vy":
		ry, err := b.v.CallForeign("GetVelocityY2D", []interface{}{b.world, b.id})
		if err != nil {
			return nil, err
		}
		return anyToF(ry), nil
	default:
		return nil, &errors.CyberError{Code: errors.ErrDotAccess, Message: fmt.Sprintf("unknown body property %q", path[0]), Suggestion: "Use x, y, vx, vy, id"}
	}
}

func (b *BodyDot) SetProp(path []string, val vm.Value) error {
	if len(path) != 1 {
		return fmt.Errorf("nested body property not supported")
	}
	switch strings.ToLower(path[0]) {
	case "x", "y":
		xr, _ := b.v.CallForeign("GetPositionX2D", []interface{}{b.world, b.id})
		yr, _ := b.v.CallForeign("GetPositionY2D", []interface{}{b.world, b.id})
		x, y := anyToF(xr), anyToF(yr)
		if strings.ToLower(path[0]) == "x" {
			x = toF(val)
		} else {
			y = toF(val)
		}
		_, err := b.v.CallForeign("SetPosition2D", []interface{}{b.world, b.id, x, y})
		return err
	case "vx", "vy":
		vxr, _ := b.v.CallForeign("GetVelocityX2D", []interface{}{b.world, b.id})
		vyr, _ := b.v.CallForeign("GetVelocityY2D", []interface{}{b.world, b.id})
		vx, vy := anyToF(vxr), anyToF(vyr)
		if strings.ToLower(path[0]) == "vx" {
			vx = toF(val)
		} else {
			vy = toF(val)
		}
		_, err := b.v.CallForeign("SetVelocity2D", []interface{}{b.world, b.id, vx, vy})
		return err
	case "friction":
		_, err := b.v.CallForeign("SetFriction2D", []interface{}{b.world, b.id, toF(val)})
		return err
	default:
		return &errors.CyberError{Code: errors.ErrDotAccess, Message: fmt.Sprintf("unknown or read-only property %q", path[0])}
	}
}

func (b *BodyDot) CallMethod(name string, args []vm.Value) (vm.Value, error) {
	switch strings.ToLower(name) {
	case "delete":
		_, err := b.v.CallForeign("DestroyBody2D", []interface{}{b.world, b.id})
		return nil, err
	default:
		return nil, &errors.CyberError{Code: errors.ErrDotAccess, Message: fmt.Sprintf("unknown body method %q", name), Suggestion: "Use delete or SetProp x/y for position"}
	}
}

// physicsDynamicDot is physics.dynamic.* (nested namespace).
type physicsDynamicDot struct {
	v *vm.VM
}

func (d *physicsDynamicDot) GetProp([]string) (vm.Value, error) { return nil, nil }
func (d *physicsDynamicDot) SetProp([]string, vm.Value) error {
	return fmt.Errorf("physics.dynamic: not assignable")
}

func (d *physicsDynamicDot) CallMethod(name string, args []vm.Value) (vm.Value, error) {
	ia := valuesToIface(args)
	switch strings.ToLower(name) {
	case "box":
		return d.v.CallForeign("PhysicsHighDynamicBox", ia)
	default:
		return nil, fmt.Errorf("physics.dynamic: unknown method %q (box)", name)
	}
}
