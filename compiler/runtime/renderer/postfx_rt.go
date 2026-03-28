package renderer

import (
	"sync"

	"cyberbasic/compiler/bindings/effect"
	"cyberbasic/compiler/bindings/raylib"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
	postFXRTMu       sync.Mutex
	postFXRT         rl.RenderTexture2D
	postFXRTReady    bool
	postFXW, postFXH int32
)

func unloadPostFXRTLocked() {
	if postFXRTReady {
		rl.UnloadRenderTexture(postFXRT)
		postFXRTReady = false
	}
	postFXW, postFXH = 0, 0
}

func ensurePostFXRenderTexture(w, h int32) {
	postFXRTMu.Lock()
	defer postFXRTMu.Unlock()
	if postFXRTReady && rl.IsRenderTextureValid(postFXRT) && postFXW == w && postFXH == h {
		return
	}
	unloadPostFXRTLocked()
	postFXRT = rl.LoadRenderTexture(w, h)
	postFXRTReady = rl.IsRenderTextureValid(postFXRT)
	postFXW, postFXH = w, h
}

// compositeTintForPostFX returns a simple tint for stub FX kinds (multiplied in one blit).
func compositeTintForPostFX(entries []effect.PostFXEntry) rl.Color {
	tint := rl.White
	for _, e := range entries {
		switch e.Kind {
		case "vignette":
			tint = rl.NewColor(uint8(float32(tint.R)*0.85), uint8(float32(tint.G)*0.85), uint8(float32(tint.B)*0.85), 255)
		case "bloom":
			tint = rl.NewColor(255, 245, 235, 255)
		case "dof":
			tint = rl.NewColor(248, 248, 252, 255)
		}
	}
	return tint
}

func (r *Renderer) drawWorld3DAnd2D() {
	cam := raylib.GetCamera3D()
	rl.BeginMode3D(cam)
	r.draw3D()
	if v := VM(); v != nil {
		raylib.DrawQueues3D(v)
	}
	rl.EndMode3D()

	if preDraw2DFn != nil {
		preDraw2DFn()
	}
	r.draw2D()
}

// compositePostFXScene renders 3D+queued 2D into an offscreen target and blits to the current framebuffer.
func (r *Renderer) compositePostFXScene() {
	entries := effect.SnapshotPostFX()
	w := int32(rl.GetRenderWidth())
	h := int32(rl.GetRenderHeight())
	if w < 1 || h < 1 {
		r.drawWorld3DAnd2D()
		return
	}
	ensurePostFXRenderTexture(w, h)
	postFXRTMu.Lock()
	rt := postFXRT
	ok := postFXRTReady
	postFXRTMu.Unlock()
	if !ok {
		r.drawWorld3DAnd2D()
		return
	}

	rl.BeginTextureMode(rt)
	rl.ClearBackground(r.clearColor)
	r.drawWorld3DAnd2D()
	rl.EndTextureMode()

	src := rl.NewRectangle(0, 0, float32(rt.Texture.Width), -float32(rt.Texture.Height))
	pos := rl.Vector2{X: 0, Y: 0}
	rl.DrawTextureRec(rt.Texture, src, pos, compositeTintForPostFX(entries))
}
