package renderer

import (
	"cyberbasic/compiler/bindings/raylib"
)

// draw2D draws the 2D layer (queued sprites, tilemaps, particles).
func (r *Renderer) draw2D() {
	if v := VM(); v != nil {
		raylib.DrawQueues2D(v)
	}
}
