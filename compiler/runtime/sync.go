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
// When UseUnifiedRenderer is enabled, runs the full unified frame.
// Otherwise, ends the frame (rl.EndDrawing) for manual/hybrid mode.
// PollInputEvents is called only in BeginDrawing (frame start). Calling it here would
// clear IsKeyPressed/IsMouseButtonPressed before the next frame reads them.
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
	rl.EndDrawing()
}
