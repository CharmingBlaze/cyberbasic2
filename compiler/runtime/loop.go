// Package runtime: StepFrame provides a single "one frame" entry point for the hybrid update/draw pipeline.
package runtime

import (
	"cyberbasic/compiler/vm"
)

// StepFrame runs one hybrid frame: get dt, step physics 2D/3D, update(dt), clear queues, draw(), flush queues.
// The VM must have raylib and hybrid bindings registered (GetFrameTime, StepAllPhysics2D/3D, ClearRenderQueues, FlushRenderQueues).
// update(dt) and draw() are invoked if the loaded chunk defines them. Use for a single entry point per frame (e.g. headless testing).
func StepFrame(v *vm.VM) error {
	if v.Chunk() == nil {
		return nil
	}
	chunk := v.Chunk()
	_, hasUpdate := chunk.GetFunction("update")
	_, hasDraw := chunk.GetFunction("draw")

	dt, err := v.CallForeign("GetFrameTime", nil)
	if err != nil {
		return err
	}
	if _, err = v.CallForeign("StepAllPhysics2D", []interface{}{dt}); err != nil {
		return err
	}
	if _, err = v.CallForeign("StepAllPhysics3D", []interface{}{dt}); err != nil {
		return err
	}
	if hasUpdate {
		if err = v.InvokeSub("update", []interface{}{dt}); err != nil {
			return err
		}
	}
	if _, err = v.CallForeign("ClearRenderQueues", nil); err != nil {
		return err
	}
	if hasDraw {
		if err = v.InvokeSub("draw", nil); err != nil {
			return err
		}
	}
	_, err = v.CallForeign("FlushRenderQueues", nil)
	return err
}
