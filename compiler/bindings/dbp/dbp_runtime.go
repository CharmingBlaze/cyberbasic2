// Package dbp - VM/Runtime: StopTask, PauseTask, ResumeTask, FixedUpdate, OnFixedUpdate.
//
// WaitFrames(n) is a language construct - use it in coroutines; it compiles to WaitSeconds(n/60).
package dbp

import (
	"cyberbasic/compiler/runtime"
	"cyberbasic/compiler/vm"
)

func toFloat64Runtime(v interface{}) float64 {
	switch x := v.(type) {
	case int:
		return float64(x)
	case float64:
		return x
	case float32:
		return float64(x)
	default:
		return 0
	}
}

// registerRuntime registers DBP-style runtime/task commands.
func registerRuntime(v *vm.VM) {
	// StopTask, PauseTask, ResumeTask require fiber name tracking in VM - stubs for now
	v.RegisterForeign("StopTask", func(args []interface{}) (interface{}, error) {
		_ = args // task name - would need VM support to stop fiber by name
		return nil, nil
	})
	v.RegisterForeign("PauseTask", func(args []interface{}) (interface{}, error) {
		_ = args
		return nil, nil
	})
	v.RegisterForeign("ResumeTask", func(args []interface{}) (interface{}, error) {
		_ = args
		return nil, nil
	})
	v.RegisterForeign("FixedUpdate", func(args []interface{}) (interface{}, error) {
		if len(args) >= 1 {
			runtime.SetFixedUpdateRate(toFloat64Runtime(args[0]))
		}
		return nil, nil
	})
	v.RegisterForeign("OnFixedUpdate", func(args []interface{}) (interface{}, error) {
		if len(args) >= 1 {
			runtime.SetFixedUpdateLabel(toString(args[0]))
		}
		return nil, nil
	})
	v.RegisterForeign("StepFrame", func(args []interface{}) (interface{}, error) {
		return nil, runtime.StepFrame(v)
	})
}

// FixedUpdateLabel returns the label set by OnFixedUpdate (for game loop integration).
func FixedUpdateLabel() string {
	return runtime.FixedUpdateLabel()
}

// FixedUpdateRate returns the rate set by FixedUpdate.
func FixedUpdateRate() float64 {
	return runtime.FixedUpdateRate()
}
