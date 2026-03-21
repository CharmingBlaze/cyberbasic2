// Package shadersys provides SHADER.* factories and material-style handles wired to raylib LoadShader / uniforms.
package shadersys

import (
	"cyberbasic/compiler/bindings/dotargs"
	"cyberbasic/compiler/vm"
	"fmt"
	"strings"
	"sync"
)

// RegisterShaderSys registers shader module, ShaderHandleSet, and wires handles to raylib shader ids.
func RegisterShaderSys(v *vm.VM) {
	v.RegisterForeign("ShaderSysVersion", func(args []interface{}) (interface{}, error) {
		return "v2-raylib", nil
	})
	v.RegisterForeign("ShaderHandleSet", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("ShaderHandleSet requires (shaderId, uniformName$, value) or (shaderId, uniformName$, r, g, b, a) for vec4")
		}
		sid := fmt.Sprint(args[0])
		u := fmt.Sprint(args[1])
		if len(args) >= 6 {
			_, err := v.CallForeign("SetShaderUniformVec4", []interface{}{sid, u, args[2], args[3], args[4], args[5]})
			return nil, err
		}
		_, err := v.CallForeign("SetShaderUniform", []interface{}{sid, u, args[2]})
		return nil, err
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
	case "pbr":
		return s.loadPreset(presetPBRFragment, "pbr", args)
	case "toon":
		return s.loadPreset(presetToonFragment, "toon", args)
	case "dissolve":
		h, err := s.loadPreset(presetDissolveFragment, "dissolve", args)
		if err != nil {
			return nil, err
		}
		// Sensible default so geometry is visible before user sets "dissolve"
		_, _ = s.v.CallForeign("SetShaderUniform", []interface{}{h.glID, "dissolve", 1.0})
		return h, nil
	case "load":
		if len(args) < 2 {
			return nil, fmt.Errorf("shader.load requires (vsPath$, fsPath$)")
		}
		vs := fmt.Sprint(args[0])
		fs := fmt.Sprint(args[1])
		r, err := s.v.CallForeign("LoadShader", []interface{}{vs, fs})
		if err != nil {
			return nil, err
		}
		glid := fmt.Sprint(r)
		return newShaderHandle(s.v, glid, "file"), nil
	default:
		return nil, fmt.Errorf("unknown shader method %q (pbr, toon, dissolve, load, version)", name)
	}
}

func (s *shaderModuleDot) loadPreset(fsCode, kind string, ctorArgs []vm.Value) (*shaderHandleDot, error) {
	_ = ctorArgs
	r, err := s.v.CallForeign("LoadShaderFromMemory", []interface{}{presetVertexShader, fsCode})
	if err != nil {
		return nil, err
	}
	glid := fmt.Sprint(r)
	return newShaderHandle(s.v, glid, kind), nil
}

type shaderHandleDot struct {
	v        *vm.VM
	glID     string // raylib foreign id, e.g. shader_1 (use with BeginShaderMode / SetShaderUniform*)
	kind     string
	mu       sync.RWMutex
	unloaded bool
}

func newShaderHandle(v *vm.VM, glID, kind string) *shaderHandleDot {
	return &shaderHandleDot{v: v, glID: glID, kind: kind}
}

func (h *shaderHandleDot) GetProp(path []string) (vm.Value, error) {
	if len(path) == 0 {
		return nil, fmt.Errorf("empty path")
	}
	p := strings.ToLower(path[0])
	h.mu.RLock()
	defer h.mu.RUnlock()
	if h.unloaded {
		return nil, fmt.Errorf("shader handle: unloaded")
	}
	switch p {
	case "id":
		return h.glID, nil
	case "kind":
		return h.kind, nil
	default:
		return nil, nil
	}
}

func (h *shaderHandleDot) SetProp(path []string, val vm.Value) error {
	if len(path) != 1 {
		return fmt.Errorf("shader handle: single property only")
	}
	h.mu.RLock()
	id, unloaded := h.glID, h.unloaded
	h.mu.RUnlock()
	if unloaded || id == "" {
		return fmt.Errorf("shader handle: unloaded")
	}
	p := strings.ToLower(path[0])
	_, err := h.v.CallForeign("ShaderHandleSet", []interface{}{id, p, val})
	return err
}

func (h *shaderHandleDot) CallMethod(name string, args []vm.Value) (vm.Value, error) {
	switch strings.ToLower(name) {
	case "set":
		if len(args) < 2 {
			return nil, fmt.Errorf("set(uniform$, value) or set(uniform$, r, g, b, a) requires 2 or 5 value arguments")
		}
		h.mu.RLock()
		id, unloaded := h.glID, h.unloaded
		h.mu.RUnlock()
		if unloaded || id == "" {
			return nil, fmt.Errorf("shader handle: unloaded")
		}
		u := strings.ToLower(fmt.Sprint(args[0]))
		if len(args) >= 5 {
			_, err := h.v.CallForeign("SetShaderUniformVec4", []interface{}{id, u, args[1], args[2], args[3], args[4]})
			return nil, err
		}
		_, err := h.v.CallForeign("SetShaderUniform", []interface{}{id, u, args[1]})
		return nil, err
	case "unload":
		h.mu.Lock()
		defer h.mu.Unlock()
		if h.unloaded || h.glID == "" {
			return nil, nil
		}
		_, err := h.v.CallForeign("UnloadShader", []interface{}{h.glID})
		h.glID = ""
		h.unloaded = true
		return nil, err
	default:
		return nil, fmt.Errorf("shader handle: use .set, property assignment, or .unload()")
	}
}
