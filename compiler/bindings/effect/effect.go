// Package effect provides EFFECT.* factories and stub camera post-FX registration (v1).
package effect

import (
	"cyberbasic/compiler/vm"
	"fmt"
	"strings"
	"sync"

	"cyberbasic/compiler/bindings/dotargs"
)

var (
	fxMu      sync.Mutex
	fxSlots   []string // effect kind ids for debugging
	cameraFX  []string // entries "id:kind" from CameraFXAddStub / camera.fx.add
	nextFX    int
)

// PostFXEntry is one queued camera post-process for the unified renderer.
type PostFXEntry struct {
	ID   string
	Kind string
}

// SnapshotPostFX returns a copy of the current camera.fx queue (for the renderer pass).
func SnapshotPostFX() []PostFXEntry {
	fxMu.Lock()
	defer fxMu.Unlock()
	out := make([]PostFXEntry, 0, len(cameraFX))
	for _, s := range cameraFX {
		id, kind := s, "unknown"
		if i := strings.Index(s, ":"); i >= 0 {
			id, kind = s[:i], s[i+1:]
		}
		out = append(out, PostFXEntry{ID: id, Kind: kind})
	}
	return out
}

// HasPostFXQueued reports whether camera.fx has any effects to apply this frame.
func HasPostFXQueued() bool {
	fxMu.Lock()
	defer fxMu.Unlock()
	return len(cameraFX) > 0
}

// RegisterEffect registers EFFECT module, stub foreigns for camera FX queue, and globals used by cameradot.
func RegisterEffect(v *vm.VM) {
	v.RegisterForeign("EffectSysVersion", func(args []interface{}) (interface{}, error) {
		return "v1-stub", nil
	})
	v.RegisterForeign("CameraFXAddStub", func(args []interface{}) (interface{}, error) {
		kind := "unknown"
		if len(args) >= 1 {
			kind = fmt.Sprint(args[0])
		}
		fxMu.Lock()
		nextFX++
		id := fmt.Sprintf("fx_%d", nextFX)
		cameraFX = append(cameraFX, id+":"+kind)
		fxSlots = append(fxSlots, kind)
		fxMu.Unlock()
		return id, nil
	})
	v.RegisterForeign("CameraFXClearStub", func(args []interface{}) (interface{}, error) {
		fxMu.Lock()
		cameraFX = cameraFX[:0]
		fxMu.Unlock()
		return nil, nil
	})

	v.SetGlobal("effect", &effectModuleDot{v: v})
}

type effectModuleDot struct {
	v *vm.VM
}

func (e *effectModuleDot) GetProp([]string) (vm.Value, error) { return nil, nil }
func (e *effectModuleDot) SetProp([]string, vm.Value) error {
	return fmt.Errorf("effect: namespace is not assignable")
}

func (e *effectModuleDot) CallMethod(name string, args []vm.Value) (vm.Value, error) {
	ia := dotargs.From(args)
	switch strings.ToLower(name) {
	case "version":
		return e.v.CallForeign("EffectSysVersion", ia)
	case "bloom", "vignette", "dof":
		fxMu.Lock()
		nextFX++
		id := fmt.Sprintf("effect_%s_%d", strings.ToLower(name), nextFX)
		fxSlots = append(fxSlots, name)
		fxMu.Unlock()
		return &effectHandleDot{id: id, kind: strings.ToLower(name), v: e.v}, nil
	default:
		return nil, fmt.Errorf("unknown effect factory %q (stub: bloom, vignette, dof)", name)
	}
}

type effectHandleDot struct {
	v    *vm.VM
	id   string
	kind string
}

func (h *effectHandleDot) GetProp(path []string) (vm.Value, error) {
	if len(path) == 1 && strings.ToLower(path[0]) == "id" {
		return h.id, nil
	}
	return nil, nil
}
func (h *effectHandleDot) SetProp([]string, vm.Value) error { return nil }
func (h *effectHandleDot) CallMethod(string, []vm.Value) (vm.Value, error) {
	return nil, fmt.Errorf("effect handle %q: use camera.fx.add(effect)", h.kind)
}
