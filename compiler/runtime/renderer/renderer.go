// Package renderer provides a unified render pipeline: 3D → 2D → GUI in a single Frame().
package renderer

import (
	"sync"

	"cyberbasic/compiler/bindings/raylib"
	"cyberbasic/compiler/vm"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
	useUnifiedMu     sync.RWMutex
	useUnified       bool
	defaultRenderer  *Renderer
	defaultRenderMu  sync.Mutex
	rendererVM       *vm.VM
	rendererVMu      sync.RWMutex
)

// SetVM sets the VM for draw callbacks that need to call foreign functions (e.g. DrawWater).
func SetVM(v *vm.VM) {
	rendererVMu.Lock()
	defer rendererVMu.Unlock()
	rendererVM = v
}

// VM returns the current VM, if set.
func VM() *vm.VM {
	rendererVMu.RLock()
	defer rendererVMu.RUnlock()
	return rendererVM
}

// SetUseUnified enables or disables the unified renderer. When enabled, SYNC runs the full frame.
func SetUseUnified(enabled bool) {
	useUnifiedMu.Lock()
	defer useUnifiedMu.Unlock()
	useUnified = enabled
}

// IsUseUnified returns whether the unified renderer is active.
func IsUseUnified() bool {
	useUnifiedMu.RLock()
	defer useUnifiedMu.RUnlock()
	return useUnified
}

// Default returns the default renderer instance, creating it if needed.
func Default() *Renderer {
	defaultRenderMu.Lock()
	defer defaultRenderMu.Unlock()
	if defaultRenderer == nil {
		defaultRenderer = New()
	}
	return defaultRenderer
}

// Draw3DFunc is called during the 3D pass. Set via SetDraw3D to inject dbp.DrawAllDBPObjects etc.
type Draw3DFunc func()

var draw3DFn Draw3DFunc

// SetDraw3D sets the 3D draw callback. Called from runtime init with dbp.DrawAllDBPObjects.
func SetDraw3D(fn Draw3DFunc) {
	draw3DFn = fn
}

// PreDraw2DFunc is called before the 2D pass (e.g. UpdateSpriteAnimations).
type PreDraw2DFunc func()

var preDraw2DFn PreDraw2DFunc

// SetPreDraw2D sets the pre-2D callback. Called before draw2D (e.g. dbp.UpdateSpriteAnimations).
func SetPreDraw2D(fn PreDraw2DFunc) {
	preDraw2DFn = fn
}

// Renderer is the unified render pipeline. Frame() handles the full frame.
type Renderer struct {
	clearColor rl.Color
}

// New creates a new Renderer with default clear color.
func New() *Renderer {
	return &Renderer{
		clearColor: rl.NewColor(25, 25, 35, 255),
	}
}

// SetClearColor sets the background clear color.
func (r *Renderer) SetClearColor(c rl.Color) {
	r.clearColor = c
}

// Frame runs one full frame: BeginDrawing → 3D pass → 2D pass → GUI pass → EndDrawing.
// Timing is driven only by beginRuntimeFrame or BeginDrawing; do not call time.Update here or physics desyncs.
func (r *Renderer) Frame() {
	// PollInputEvents is called once at frame start (beginRuntimeFrame); do not poll here or IsKeyPressed/IsMouseButtonPressed get cleared
	rl.BeginDrawing()
	rl.ClearBackground(r.clearColor)

	// Shadow pass (when enabled): render from light POV to depth texture
	RenderShadowPass()

	// 3D pass: scene (sky, terrain, water, clouds, objects) + queued 3D items
	cam := raylib.GetCamera3D()
	rl.BeginMode3D(cam)
	r.draw3D()
	if v := VM(); v != nil {
		raylib.DrawQueues3D(v)
	}
	rl.EndMode3D()

	// Pre-2D: sprite animation tick, etc.
	if preDraw2DFn != nil {
		preDraw2DFn()
	}

	// 2D pass: queued 2D items
	r.draw2D()

	// GUI pass: queued GUI items
	r.drawUI()

	rl.EndDrawing()
}

// FrameIfUnified runs Frame() if unified mode is enabled. Returns true if frame was run.
func FrameIfUnified() bool {
	if !IsUseUnified() {
		return false
	}
	Default().Frame()
	return true
}
