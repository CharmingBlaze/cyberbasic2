package renderer

import "cyberbasic/compiler/runtime/camera"

// draw3D draws the 3D scene: sky, terrain, water, clouds, objects, particles.
func (r *Renderer) draw3D() {
	camera.UpdateAttachments()
	if draw3DFn != nil {
		draw3DFn()
	}
}
