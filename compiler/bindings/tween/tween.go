// Package tween provides TWEEN v2 module and TweenRegister foreign; ticks from runtime frame loop.
package tween

import (
	"cyberbasic/compiler/bindings/dotargs"
	"cyberbasic/compiler/vm"
	"fmt"
	"strings"
	"sync"
)

type tweenEntry struct {
	target vm.DotObject
	prop   string
	fromV  float64
	toV    float64
	dur    float64
	t      float64
	done   bool
}

var (
	tweenMu sync.Mutex
	tweens  []*tweenEntry
)

func toF(x interface{}) float64 {
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

// RegisterTween registers TweenRegister foreign and global tween module (stats helper).
func RegisterTween(v *vm.VM) {
	v.RegisterForeign("TweenRegister", func(args []interface{}) (interface{}, error) {
		if len(args) < 5 {
			return nil, fmt.Errorf("TweenRegister(target, prop$, from, to, seconds) requires 5 arguments")
		}
		d, ok := args[0].(vm.DotObject)
		if !ok {
			return nil, fmt.Errorf("TweenRegister: target must be a handle (DotObject)")
		}
		prop := strings.ToLower(strings.TrimSpace(fmt.Sprint(args[1])))
		fromV := toF(args[2])
		toV := toF(args[3])
		dur := toF(args[4])
		if dur <= 0 {
			dur = 0.001
		}
		tweenMu.Lock()
		tweens = append(tweens, &tweenEntry{target: d, prop: prop, fromV: fromV, toV: toV, dur: dur})
		tweenMu.Unlock()
		return len(tweens), nil
	})

	v.SetGlobal("tween", &tweenModuleDot{v: v})
}

type tweenModuleDot struct {
	v *vm.VM
}

func (t *tweenModuleDot) GetProp([]string) (vm.Value, error) { return nil, nil }
func (t *tweenModuleDot) SetProp([]string, vm.Value) error {
	return fmt.Errorf("tween: namespace is not assignable")
}

func (t *tweenModuleDot) CallMethod(name string, args []vm.Value) (vm.Value, error) {
	switch strings.ToLower(name) {
	case "count":
		tweenMu.Lock()
		n := len(tweens)
		tweenMu.Unlock()
		return float64(n), nil
	case "register":
		return t.v.CallForeign("TweenRegister", dotargs.From(args))
	default:
		return nil, fmt.Errorf("tween: methods: register, count (or flat TweenRegister)")
	}
}

// Tick advances active tweens and applies SetProp on targets (call from runtime once per frame).
func Tick(dt float64) {
	if dt <= 0 {
		return
	}
	tweenMu.Lock()
	defer tweenMu.Unlock()
	for _, e := range tweens {
		if e == nil || e.done {
			continue
		}
		e.t += dt
		alpha := e.t / e.dur
		if alpha > 1 {
			alpha = 1
		}
		val := e.fromV + (e.toV-e.fromV)*alpha
		_ = e.target.SetProp([]string{e.prop}, val)
		if alpha >= 1 {
			e.done = true
		}
	}
}
