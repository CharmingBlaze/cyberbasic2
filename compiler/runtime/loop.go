// Package runtime: frame stepping helpers for hybrid and implicit update/draw loops.
package runtime

import (
	"cyberbasic/compiler/bindings/inputmap"
	"cyberbasic/compiler/bindings/raylib"
	"cyberbasic/compiler/bindings/tween"
	"cyberbasic/compiler/runtime/renderer"
	"cyberbasic/compiler/runtime/time"
	"cyberbasic/compiler/vm"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const maxFixedCatchupSteps = 8

func beginRuntimeFrame(v *vm.VM) (float64, error) {
	rl.PollInputEvents()
	inputmap.TickInputMap(v)
	raylib.CaptureOrbitWheel()
	dt, err := v.CallForeign("GetFrameTime", nil)
	if err != nil {
		return 0, err
	}
	dtVal := float32(0)
	if f, ok := dt.(float64); ok {
		dtVal = float32(f)
	}
	time.Update(dtVal)
	fixedStep := time.GetFixedDeltaTime()
	if fixedStep <= 0 {
		fixedStep = 1.0 / 60.0
	}
	fixedStepArg := float64(fixedStep)
	steps := 0
	for time.GetAccumulator() >= fixedStep && steps < maxFixedCatchupSteps {
		if _, err = v.CallForeign("StepAllPhysics2D", []interface{}{fixedStepArg}); err != nil {
			return 0, err
		}
		if _, err = v.CallForeign("StepAllPhysics3D", []interface{}{fixedStepArg}); err != nil {
			return 0, err
		}
		if label := FixedUpdateLabel(); label != "" {
			if err = v.InvokeSub(label, []interface{}{fixedStepArg}); err != nil {
				return 0, err
			}
		}
		time.ConsumeAccumulator(fixedStep)
		steps++
	}
	if steps == maxFixedCatchupSteps {
		time.ClampAccumulator(fixedStep)
	}
	tween.Tick(float64(dtVal))
	return float64(dtVal), nil
}

// StepFrame runs one hybrid frame: get dt, step fixed physics/callbacks, update(dt), clear queues, draw(), flush queues.
// The VM must have raylib and hybrid bindings registered (GetFrameTime, StepAllPhysics2D/3D, ClearRenderQueues, FlushRenderQueues).
// update(dt) and draw() are invoked if the loaded chunk defines them. Use for the compiler-emitted hybrid loop and tests.
func StepFrame(v *vm.VM) error {
	if v.Chunk() == nil {
		return nil
	}
	chunk := v.Chunk()
	_, hasUpdate := chunk.GetFunction("update")
	_, hasDraw := chunk.GetFunction("draw")

	dt, err := beginRuntimeFrame(v)
	if err != nil {
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

// StepImplicitFrame runs one DBP-style implicit frame using OnUpdate/OnDraw naming.
// When UseUnifiedRenderer is enabled: queue path—OnDraw runs with SetInsideDraw so Draw* calls queue;
// user writes SYNC at end of OnDraw; SYNC runs Frame() which draws from queues and presents.
// When UseUnifiedRenderer is off: BeginDrawing, OnDraw, EndDrawing (legacy path).
func StepImplicitFrame(v *vm.VM) error {
	if v.Chunk() == nil {
		return nil
	}
	chunk := v.Chunk()
	_, hasUpdate := chunk.GetFunction("onupdate")
	_, hasDraw := chunk.GetFunction("ondraw")

	dt, err := beginRuntimeFrame(v)
	if err != nil {
		return err
	}
	if hasUpdate {
		if err = v.InvokeSub("OnUpdate", []interface{}{dt}); err != nil {
			return err
		}
	}

	if renderer.IsUseUnified() {
		// Queue path: OnDraw queues; SYNC in OnDraw runs Frame() and presents.
		if _, err = v.CallForeign("ClearRenderQueues", nil); err != nil {
			return err
		}
		if hasDraw {
			v.SetInsideDraw(true)
			if err = v.InvokeSub("OnDraw", nil); err != nil {
				v.SetInsideDraw(false)
				return err
			}
			v.SetInsideDraw(false)
		}
		return nil
	}

	// Legacy path: direct BeginDrawing/EndDrawing (beginRuntimeFrame already polled; do not poll again or IsKeyPressed gets cleared)
	rl.BeginDrawing()
	rl.ClearBackground(rl.NewColor(0, 0, 0, 255))
	if hasDraw {
		if err = v.InvokeSub("OnDraw", nil); err != nil {
			rl.EndDrawing()
			return err
		}
	}
	rl.EndDrawing()
	return nil
}
