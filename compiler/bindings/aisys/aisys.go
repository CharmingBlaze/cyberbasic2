// Package aisys provides ai.* as a thin facade over navigation (nav grid / mesh / agent).
package aisys

import (
	"cyberbasic/compiler/bindings/dotargs"
	"cyberbasic/compiler/bindings/navigation"
	"cyberbasic/compiler/vm"
	"fmt"
	"strings"
)

// RegisterAisys registers ai.* methods that delegate to navigation foreigns.
func RegisterAisys(v *vm.VM) {
	v.RegisterForeign("AisysVersion", func(args []interface{}) (interface{}, error) {
		return "v1-nav-delegate", nil
	})
	v.SetGlobal("ai", &aiModuleDot{v: v})
}

type aiModuleDot struct {
	v *vm.VM
}

func (a *aiModuleDot) GetProp([]string) (vm.Value, error) { return nil, nil }
func (a *aiModuleDot) SetProp([]string, vm.Value) error {
	return fmt.Errorf("ai: namespace is not assignable")
}

func (a *aiModuleDot) CallMethod(name string, args []vm.Value) (vm.Value, error) {
	low := strings.ToLower(name)
	switch low {
	case "version":
		return a.v.CallForeign("AisysVersion", dotargs.From(args))
	case "agent":
		if len(args) < 1 {
			return nil, fmt.Errorf("ai.agent requires (navAgentId$)")
		}
		return newAIAgentDot(a.v, fmt.Sprint(args[0])), nil
	}
	if fn, ok := navigation.MethodToForeign[low]; ok {
		ia := make([]interface{}, len(args))
		for i := range args {
			ia[i] = args[i]
		}
		return a.v.CallForeign(fn, ia)
	}
	return nil, fmt.Errorf("unknown ai method %q (navigation aliases, agent, version)", name)
}

// aiAgentDot wraps a nav agent id for handle-style calls.
type aiAgentDot struct {
	v  *vm.VM
	id string
}

func newAIAgentDot(v *vm.VM, id string) *aiAgentDot {
	return &aiAgentDot{v: v, id: id}
}

func (d *aiAgentDot) GetProp(path []string) (vm.Value, error) {
	if len(path) != 1 {
		return nil, fmt.Errorf("ai agent: single property only")
	}
	switch strings.ToLower(path[0]) {
	case "id":
		return d.id, nil
	default:
		return nil, nil
	}
}

func (d *aiAgentDot) SetProp([]string, vm.Value) error {
	return fmt.Errorf("ai agent: not assignable")
}

func (d *aiAgentDot) CallMethod(name string, args []vm.Value) (vm.Value, error) {
	ia := append([]interface{}{d.id}, valuesToIface(args)...)
	switch strings.ToLower(name) {
	case "setdestination":
		if len(args) < 3 {
			return nil, fmt.Errorf("setdestination(x, y, z) requires 3 arguments")
		}
		return d.v.CallForeign("NavAgentSetDestination", ia)
	case "update":
		if len(args) < 1 {
			return nil, fmt.Errorf("update(dt) requires dt")
		}
		return d.v.CallForeign("NavAgentUpdate", ia)
	case "setposition":
		if len(args) < 3 {
			return nil, fmt.Errorf("setposition(x, y, z) requires 3 arguments")
		}
		return d.v.CallForeign("NavAgentSetPosition", ia)
	case "setspeed":
		if len(args) < 1 {
			return nil, fmt.Errorf("setspeed(s) requires speed")
		}
		return d.v.CallForeign("NavAgentSetSpeed", ia)
	case "setradius":
		if len(args) < 1 {
			return nil, fmt.Errorf("setradius(r) requires radius")
		}
		return d.v.CallForeign("NavAgentSetRadius", ia)
	case "nextwaypoint":
		return d.v.CallForeign("NavAgentGetNextWaypoint", []interface{}{d.id})
	default:
		return nil, fmt.Errorf("ai agent: use setdestination, update, setposition, setspeed, setradius, nextwaypoint")
	}
}

func valuesToIface(args []vm.Value) []interface{} {
	ia := make([]interface{}, len(args))
	for i := range args {
		ia[i] = args[i]
	}
	return ia
}
