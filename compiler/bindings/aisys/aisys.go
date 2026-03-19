// Package aisys is reserved for NAVMESH / AI.AGENT / BTREE (Phase 12).
package aisys

import (
	"cyberbasic/compiler/vm"
	"fmt"
	"strings"
)

// RegisterAisys registers placeholder foreigns (extend in Phase 12).
func RegisterAisys(v *vm.VM) {
	v.RegisterForeign("AisysVersion", func(args []interface{}) (interface{}, error) {
		return "stub", nil
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
	ia := make([]interface{}, len(args))
	for i := range args {
		ia[i] = args[i]
	}
	switch strings.ToLower(name) {
	case "version":
		return a.v.CallForeign("AisysVersion", ia)
	default:
		return nil, fmt.Errorf("unknown ai method %q (stub: version)", name)
	}
}
