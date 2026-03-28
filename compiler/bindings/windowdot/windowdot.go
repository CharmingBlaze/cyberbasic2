// Package windowdot provides the WINDOW DotObject for implicit/explicit window state.
package windowdot

import (
	"cyberbasic/compiler/errors"
	"cyberbasic/compiler/vm"
	"fmt"
	"strings"
	"sync"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// RegisterWindowDot installs the global WINDOW handle on the VM (lowercase key "window").
func RegisterWindowDot(v *vm.VM) {
	v.SetGlobal("window", NewWindowDot(v))
}

// WindowDot implements vm.DotObject for WINDOW.* properties.
type WindowDot struct {
	v  *vm.VM
	mu sync.RWMutex
	pendingTitle      string
	pendingW, pendingH int32
	pendingFullscreen bool
	pendingVSync      bool
	pendingIconPath   string
	pendingTargetFPS  int32
	hasPending        bool
}

// NewWindowDot creates a WINDOW handle with defaults matching implicit loop.
func NewWindowDot(v *vm.VM) *WindowDot {
	return &WindowDot{
		v:                v,
		pendingTitle:     "CyberBasic 2",
		pendingW:         1280,
		pendingH:         720,
		pendingTargetFPS: 60,
	}
}

// ImplicitSize returns width, height, title for RunImplicitLoop before InitWindow.
func (w *WindowDot) ImplicitSize() (int32, int32, string) {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.pendingW, w.pendingH, w.pendingTitle
}

// TargetFPS returns pending target FPS for SetTargetFPS after InitWindow.
func (w *WindowDot) TargetFPS() int32 {
	w.mu.RLock()
	defer w.mu.RUnlock()
	if w.pendingTargetFPS <= 0 {
		return 60
	}
	return w.pendingTargetFPS
}

// GetProp implements vm.DotObject.
func (w *WindowDot) GetProp(path []string) (vm.Value, error) {
	if len(path) == 0 {
		return nil, fmt.Errorf("empty property path")
	}
	p := strings.ToLower(path[0])
	switch p {
	case "title":
		w.mu.RLock()
		t := w.pendingTitle
		w.mu.RUnlock()
		return t, nil
	case "width":
		if rl.IsWindowReady() {
			return float64(rl.GetScreenWidth()), nil
		}
		w.mu.RLock()
		v := float64(w.pendingW)
		w.mu.RUnlock()
		return v, nil
	case "height":
		if rl.IsWindowReady() {
			return float64(rl.GetScreenHeight()), nil
		}
		w.mu.RLock()
		v := float64(w.pendingH)
		w.mu.RUnlock()
		return v, nil
	case "fps":
		if rl.IsWindowReady() {
			return float64(rl.GetFPS()), nil
		}
		return float64(0), nil
	case "deltatime", "dt":
		return float64(rl.GetFrameTime()), nil
	case "mousex":
		if rl.IsWindowReady() {
			return float64(rl.GetMouseX()), nil
		}
		return float64(0), nil
	case "mousey":
		if rl.IsWindowReady() {
			return float64(rl.GetMouseY()), nil
		}
		return float64(0), nil
	case "screenwidth":
		return float64(rl.GetRenderWidth()), nil
	case "screenheight":
		return float64(rl.GetRenderHeight()), nil
	case "fullscreen":
		if rl.IsWindowReady() {
			return rl.IsWindowFullscreen(), nil
		}
		w.mu.RLock()
		f := w.pendingFullscreen
		w.mu.RUnlock()
		return f, nil
	case "vsync":
		w.mu.RLock()
		v := w.pendingVSync
		w.mu.RUnlock()
		return v, nil
	case "targetfps":
		w.mu.RLock()
		v := float64(w.pendingTargetFPS)
		w.mu.RUnlock()
		return v, nil
	case "icon":
		w.mu.RLock()
		s := w.pendingIconPath
		w.mu.RUnlock()
		return s, nil
	default:
		return nil, &errors.CyberError{
			Code:       errors.ErrDotAccess,
			Message:    fmt.Sprintf("Unknown WINDOW property %q", path[0]),
			Suggestion: "Use title, width, height, fps, deltatime, mousex, mousey, fullscreen, vsync, targetfps, icon",
		}
	}
}

// SetProp implements vm.DotObject.
func (w *WindowDot) SetProp(path []string, val vm.Value) error {
	if len(path) != 1 {
		return fmt.Errorf("WINDOW: nested property set not supported yet")
	}
	p := strings.ToLower(path[0])
	w.mu.Lock()
	defer w.mu.Unlock()
	w.hasPending = true
	switch p {
	case "title":
		w.pendingTitle = fmt.Sprint(val)
		if rl.IsWindowReady() {
			rl.SetWindowTitle(w.pendingTitle)
		}
	case "width":
		w.pendingW = int32(toFloat(val))
		if rl.IsWindowReady() {
			rl.SetWindowSize(int(w.pendingW), int(rl.GetScreenHeight()))
		}
	case "height":
		w.pendingH = int32(toFloat(val))
		if rl.IsWindowReady() {
			rl.SetWindowSize(int(rl.GetScreenWidth()), int(w.pendingH))
		}
	case "fullscreen":
		w.pendingFullscreen = toBool(val)
		if rl.IsWindowReady() && w.pendingFullscreen != rl.IsWindowFullscreen() {
			rl.ToggleFullscreen()
		}
	case "vsync":
		w.pendingVSync = toBool(val)
	case "icon":
		w.pendingIconPath = fmt.Sprint(val)
	case "targetfps":
		w.pendingTargetFPS = int32(toFloat(val))
		if rl.IsWindowReady() {
			rl.SetTargetFPS(w.pendingTargetFPS)
		}
	default:
		return &errors.CyberError{
			Code:       errors.ErrDotAccess,
			Message:    fmt.Sprintf("Unknown or read-only WINDOW property %q", path[0]),
			Suggestion: "Writable: title, width, height, fullscreen, vsync, icon, targetfps",
		}
	}
	return nil
}

// CallMethod implements vm.DotObject (flat-command aliases).
func (w *WindowDot) CallMethod(name string, args []vm.Value) (vm.Value, error) {
	if w.v == nil {
		return nil, fmt.Errorf("window: VM not wired")
	}
	ia := make([]interface{}, len(args))
	for i := range args {
		ia[i] = args[i]
	}
	switch strings.ToLower(name) {
	case "initwindow":
		return w.v.CallForeign("InitWindow", ia)
	case "close", "closewindow":
		return w.v.CallForeign("CloseWindow", ia)
	case "resize", "setwindowsize":
		return w.v.CallForeign("SetWindowSize", ia)
	case "setwindowminsize":
		return w.v.CallForeign("SetWindowMinSize", ia)
	case "togglefullscreen":
		return w.v.CallForeign("ToggleFullscreen", ia)
	case "settargetfps":
		return w.v.CallForeign("SetTargetFPS", ia)
	case "windowshouldclose":
		return w.v.CallForeign("WindowShouldClose", ia)
	default:
		return nil, &errors.CyberError{
			Code:       errors.ErrDotAccess,
			Message:    fmt.Sprintf("Unknown WINDOW method %q", name),
			Suggestion: "Use initwindow, close, resize, setwindowminsize, togglefullscreen, settargetfps, windowshouldclose, or properties (title, width, …)",
		}
	}
}

func toFloat(v vm.Value) float64 {
	switch x := v.(type) {
	case float64:
		return x
	case int:
		return float64(x)
	case int32:
		return float64(x)
	case int64:
		return float64(x)
	default:
		return 0
	}
}

func toBool(v vm.Value) bool {
	switch x := v.(type) {
	case bool:
		return x
	case int:
		return x != 0
	case float64:
		return x != 0
	default:
		return false
	}
}
