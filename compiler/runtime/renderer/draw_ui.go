package renderer

import (
	"cyberbasic/compiler/bindings/raylib"
)

// drawUI draws the GUI layer (queued raygui items).
func (r *Renderer) drawUI() {
	if v := VM(); v != nil {
		raylib.DrawQueuesGUI(v)
	}
}
