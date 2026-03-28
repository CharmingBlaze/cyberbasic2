// Package inputmap provides INPUT.MAP-style action queries (additive).
package inputmap

import (
	"cyberbasic/compiler/vm"
	"fmt"
	"strings"
	"sync"
)

type actionState struct {
	wasDown bool
	isDown  bool
}

var (
	mu       sync.Mutex
	bindings = make(map[string][]int32)
	states   = make(map[string]*actionState)
)

// TickInputMap should be called once per frame from runtime (after PollInput).
func TickInputMap(v *vm.VM) {
	mu.Lock()
	defer mu.Unlock()
	for act := range bindings {
		st := states[act]
		if st == nil {
			st = &actionState{}
			states[act] = st
		}
		st.wasDown = st.isDown
		st.isDown = false
		for _, kc := range bindings[act] {
			r, err := v.CallForeign("IsKeyDown", []interface{}{int(kc)})
			if err == nil {
				if b, ok := r.(bool); ok && b {
					st.isDown = true
					break
				}
			}
		}
	}
}

func toKeyCode(a interface{}) int32 {
	switch x := a.(type) {
	case int:
		return int32(x)
	case int32:
		return x
	case int64:
		return int32(x)
	case float64:
		return int32(x)
	default:
		return 0
	}
}

// RegisterInputmap registers InputMapRegister, InputPressed, InputHeld, InputReleased.
func RegisterInputmap(v *vm.VM) {
	v.RegisterForeign("InputMapRegister", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("InputMapRegister requires (action$, keyCode)")
		}
		act := strings.ToLower(fmt.Sprint(args[0]))
		kc := toKeyCode(args[1])
		mu.Lock()
		bindings[act] = append(bindings[act], kc)
		mu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("InputPressed", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return false, nil
		}
		act := strings.ToLower(fmt.Sprint(args[0]))
		mu.Lock()
		st := states[act]
		mu.Unlock()
		if st == nil {
			return false, nil
		}
		return !st.wasDown && st.isDown, nil
	})
	v.RegisterForeign("InputHeld", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return false, nil
		}
		act := strings.ToLower(fmt.Sprint(args[0]))
		mu.Lock()
		st := states[act]
		mu.Unlock()
		if st == nil {
			return false, nil
		}
		return st.isDown, nil
	})
	v.RegisterForeign("InputReleased", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return false, nil
		}
		act := strings.ToLower(fmt.Sprint(args[0]))
		mu.Lock()
		st := states[act]
		mu.Unlock()
		if st == nil {
			return false, nil
		}
		return st.wasDown && !st.isDown, nil
	})

	v.SetGlobal("input", &inputRootDot{v: v})
}

// inputRootDot exposes INPUT.MAP.* (global key "input").
type inputRootDot struct {
	v *vm.VM
}

func (r *inputRootDot) GetProp(path []string) (vm.Value, error) {
	if len(path) == 1 && strings.ToLower(path[0]) == "map" {
		return &inputMapDot{v: r.v}, nil
	}
	return nil, fmt.Errorf("input: unknown property %v (use map)", path)
}

func (r *inputRootDot) SetProp([]string, vm.Value) error {
	return fmt.Errorf("input: namespace is not assignable")
}

func (r *inputRootDot) CallMethod(name string, args []vm.Value) (vm.Value, error) {
	ia := make([]interface{}, len(args))
	for i := range args {
		ia[i] = args[i]
	}
	switch strings.ToLower(name) {
	case "pressed":
		return r.v.CallForeign("InputPressed", ia)
	case "held":
		return r.v.CallForeign("InputHeld", ia)
	case "released":
		return r.v.CallForeign("InputReleased", ia)
	default:
		return nil, fmt.Errorf("input: use input.map.* or input.pressed / held / released (action$)")
	}
}

type inputMapDot struct {
	v *vm.VM
}

func (m *inputMapDot) GetProp([]string) (vm.Value, error) { return nil, nil }
func (m *inputMapDot) SetProp([]string, vm.Value) error {
	return fmt.Errorf("input.map: not assignable")
}

func (m *inputMapDot) CallMethod(name string, args []vm.Value) (vm.Value, error) {
	ia := make([]interface{}, len(args))
	for i := range args {
		ia[i] = args[i]
	}
	switch strings.ToLower(name) {
	case "register":
		return m.v.CallForeign("InputMapRegister", ia)
	case "pressed":
		return m.v.CallForeign("InputPressed", ia)
	case "held":
		return m.v.CallForeign("InputHeld", ia)
	case "released":
		return m.v.CallForeign("InputReleased", ia)
	default:
		return nil, fmt.Errorf("input.map: unknown method %q", name)
	}
}
