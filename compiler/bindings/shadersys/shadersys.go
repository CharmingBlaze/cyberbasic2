// Package shadersys provides SHADER.* factories and material-style handles (v1 stubs; wire to raylib later).
package shadersys

import (
	"cyberbasic/compiler/vm"
	"fmt"
	"strings"
	"sync"

	"cyberbasic/compiler/bindings/dotargs"
)

var shaderSeqMu sync.Mutex
var shaderSeq int

// RegisterShaderSys registers shader module, stub foreigns, and ShaderHandleSet for uniform-style props.
func RegisterShaderSys(v *vm.VM) {
	v.RegisterForeign("ShaderSysVersion", func(args []interface{}) (interface{}, error) {
		return "v1-stub", nil
	})
	v.RegisterForeign("ShaderHandleSet", func(args []interface{}) (interface{}, error) {
		// Reserved for future: (handleId, uniform$, value)
		return nil, nil
	})

	v.SetGlobal("shader", &shaderModuleDot{v: v})
}

type shaderModuleDot struct {
	v *vm.VM
}

func (s *shaderModuleDot) GetProp([]string) (vm.Value, error) { return nil, nil }
func (s *shaderModuleDot) SetProp([]string, vm.Value) error {
	return fmt.Errorf("shader: namespace is not assignable")
}

func (s *shaderModuleDot) CallMethod(name string, args []vm.Value) (vm.Value, error) {
	ia := dotargs.From(args)
	switch strings.ToLower(name) {
	case "version":
		return s.v.CallForeign("ShaderSysVersion", ia)
	case "pbr", "toon", "dissolve":
		shaderSeqMu.Lock()
		shaderSeq++
		id := fmt.Sprintf("sh_%s_%d", strings.ToLower(name), shaderSeq)
		shaderSeqMu.Unlock()
		return newShaderHandle(s.v, id, strings.ToLower(name), args), nil
	case "load":
		if len(args) < 1 {
			return nil, fmt.Errorf("shader.load requires (path$)")
		}
		path := fmt.Sprint(args[0])
		shaderSeqMu.Lock()
		shaderSeq++
		id := fmt.Sprintf("sh_file_%d", shaderSeq)
		shaderSeqMu.Unlock()
		return newShaderHandle(s.v, id, "file:"+path, nil), nil
	default:
		return nil, fmt.Errorf("unknown shader method %q (pbr, toon, dissolve, load, version)", name)
	}
}

type shaderHandleDot struct {
	v      *vm.VM
	id     string
	kind   string
	props  map[string]float64
	colors map[string][4]float64 // r,g,b,a for edgeColor etc.
	strs   map[string]string
	mu     sync.RWMutex
}

func newShaderHandle(v *vm.VM, id, kind string, ctorArgs []vm.Value) *shaderHandleDot {
	h := &shaderHandleDot{
		v: v, id: id, kind: kind,
		props:  make(map[string]float64),
		colors: make(map[string][4]float64),
		strs:   make(map[string]string),
	}
	// Seed ctor dict-like args: first arg may be map from BASIC - stub stores nothing deep
	_ = ctorArgs
	return h
}

func (h *shaderHandleDot) GetProp(path []string) (vm.Value, error) {
	if len(path) == 0 {
		return nil, fmt.Errorf("empty path")
	}
	p := strings.ToLower(path[0])
	h.mu.RLock()
	defer h.mu.RUnlock()
	switch p {
	case "id":
		return h.id, nil
	case "kind":
		return h.kind, nil
	default:
		if v, ok := h.props[p]; ok {
			return v, nil
		}
		return nil, nil
	}
}

func (h *shaderHandleDot) SetProp(path []string, val vm.Value) error {
	if len(path) != 1 {
		return fmt.Errorf("shader handle: single property only")
	}
	p := strings.ToLower(path[0])
	h.mu.Lock()
	defer h.mu.Unlock()
	h.props[p] = toF(val)
	_, _ = h.v.CallForeign("ShaderHandleSet", []interface{}{h.id, p, val})
	return nil
}

func (h *shaderHandleDot) CallMethod(name string, args []vm.Value) (vm.Value, error) {
	switch strings.ToLower(name) {
	case "set":
		if len(args) < 2 {
			return nil, fmt.Errorf("set(uniform$, value) requires 2 arguments")
		}
		u := strings.ToLower(fmt.Sprint(args[0]))
		h.mu.Lock()
		h.props[u] = toF(args[1])
		h.mu.Unlock()
		return h.v.CallForeign("ShaderHandleSet", []interface{}{h.id, u, args[1]})
	default:
		return nil, fmt.Errorf("shader handle: use .set or property assignment")
	}
}

func toF(x vm.Value) float64 {
	switch v := x.(type) {
	case float64:
		return v
	case int:
		return float64(v)
	case int32:
		return float64(v)
	case int64:
		return float64(v)
	default:
		return 0
	}
}
