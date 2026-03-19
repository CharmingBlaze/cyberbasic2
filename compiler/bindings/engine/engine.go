// Package engine provides engine.subsystem composition: GetProp returns other registered module DotObjects.
package engine

import (
	"fmt"
	"strings"

	"cyberbasic/compiler/vm"
)

type engineDot struct {
	v *vm.VM
}

// RegisterEngine must run after all other SetGlobal module registrations.
func RegisterEngine(v *vm.VM) {
	v.SetGlobal("engine", &engineDot{v: v})
}

func (e *engineDot) GetProp(path []string) (vm.Value, error) {
	if len(path) != 1 {
		return nil, fmt.Errorf("engine: use engine.subsystem (e.g. engine.ecs)")
	}
	key := strings.ToLower(path[0])
	g := e.v.Globals()
	val, ok := g[key]
	if !ok {
		return nil, fmt.Errorf("engine: unknown subsystem %q", path[0])
	}
	return val, nil
}

func (e *engineDot) SetProp([]string, vm.Value) error {
	return fmt.Errorf("engine: namespace is not assignable")
}

func (e *engineDot) CallMethod(string, []vm.Value) (vm.Value, error) {
	return nil, fmt.Errorf("engine: use property access (engine.ecs, engine.net, …)")
}
