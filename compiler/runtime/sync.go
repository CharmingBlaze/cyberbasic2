// Package runtime: SyncFrame provides the SYNC command implementation.
package runtime

import (
	"fmt"

	"cyberbasic/compiler/bindings/raylib"
	"cyberbasic/compiler/runtime/renderer"
	"cyberbasic/compiler/vm"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// RegisterFlushOverride registers FlushRenderQueues to use unified renderer when enabled.
func RegisterFlushOverride(v *vm.VM) {
	v.RegisterForeign("FlushRenderQueues", func(args []interface{}) (interface{}, error) {
		if renderer.IsUseUnified() {
			renderer.Default().Frame()
			return nil, nil
		}
		return raylib.FlushRenderQueues(v)
	})
}

// SyncFrame is called by the SYNC/Sync foreign command.
// SYNC = the update: poll input (for next frame), end frame, present.
// When UseUnifiedRenderer is enabled, runs the full unified frame.
// Otherwise: PollInputEvents, CaptureOrbitWheel, rl.EndDrawing.
// The WHILE loop defines the frame; SYNC is the update step at the end.
var syncDebugCount uint64

func SyncFrame() {
	if raylib.DebugRender() {
		syncDebugCount++
		if syncDebugCount%60 == 1 {
			fmt.Println("[DEBUG] SyncFrame (EndDrawing)")
		}
	}
	if renderer.FrameIfUnified() {
		return
	}
	rl.PollInputEvents()
	raylib.CaptureOrbitWheel()
	rl.EndDrawing()
}
