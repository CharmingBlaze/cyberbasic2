// Package modfacade provides a generic vm.DotObject that dispatches CallMethod to RegisterForeign names.
package modfacade

import (
	"cyberbasic/compiler/vm"
	"fmt"
	"strings"
)

// ModuleDot is a singleton namespace: CallMethod(name, args) -> v.CallForeign(foreignName, args).
type ModuleDot struct {
	v *vm.VM
	m map[string]string
}

// New returns a DotObject that dispatches lowercase method names to the given foreign API names.
func New(v *vm.VM, methodToForeign map[string]string) *ModuleDot {
	cp := make(map[string]string, len(methodToForeign))
	for k, val := range methodToForeign {
		cp[strings.ToLower(k)] = val
	}
	return &ModuleDot{v: v, m: cp}
}

func (d *ModuleDot) GetProp([]string) (vm.Value, error) { return nil, nil }

func (d *ModuleDot) SetProp([]string, vm.Value) error {
	return fmt.Errorf("module namespace is not assignable")
}

func (d *ModuleDot) CallMethod(name string, args []vm.Value) (vm.Value, error) {
	fn, ok := d.m[strings.ToLower(name)]
	if !ok {
		return nil, fmt.Errorf("unknown module method %q", name)
	}
	ia := make([]interface{}, len(args))
	for i := range args {
		ia[i] = args[i]
	}
	return d.v.CallForeign(fn, ia)
}
