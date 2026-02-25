// Package raylib: view (render-target viewport) helpers for split-screen or picture-in-picture.
// CreateView(viewId, x, y, w, h), SetViewTarget(viewId, renderTextureId), DrawView(viewId).
package raylib

import (
	"cyberbasic/compiler/vm"
	"fmt"
	"sync"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type viewState struct {
	X, Y, W, H   int32
	RenderTexID  string
}

var (
	viewsMu sync.RWMutex
	views   = make(map[string]*viewState)
)

func registerViews(v *vm.VM) {
	v.RegisterForeign("CreateView", func(args []interface{}) (interface{}, error) {
		if len(args) < 5 {
			return nil, fmt.Errorf("CreateView requires (viewId, x, y, width, height)")
		}
		id := toString(args[0])
		viewsMu.Lock()
		views[id] = &viewState{
			X: toInt32(args[1]), Y: toInt32(args[2]),
			W: toInt32(args[3]), H: toInt32(args[4]),
		}
		viewsMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("SetViewTarget", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("SetViewTarget requires (viewId, renderTextureId)")
		}
		id := toString(args[0])
		rtId := toString(args[1])
		viewsMu.Lock()
		defer viewsMu.Unlock()
		vw, ok := views[id]
		if !ok {
			return nil, fmt.Errorf("unknown view: %s", id)
		}
		vw.RenderTexID = rtId
		return nil, nil
	})
	v.RegisterForeign("DrawView", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("DrawView requires (viewId)")
		}
		id := toString(args[0])
		viewsMu.RLock()
		vw, ok := views[id]
		viewsMu.RUnlock()
		if !ok {
			return nil, fmt.Errorf("unknown view: %s", id)
		}
		if vw.RenderTexID == "" {
			return nil, nil
		}
		renderTexMu.Lock()
		rt, ok := renderTextures[vw.RenderTexID]
		renderTexMu.Unlock()
		if !ok {
			return nil, fmt.Errorf("unknown render texture: %s", vw.RenderTexID)
		}
		tex := rt.Texture
		src := rl.Rectangle{X: 0, Y: 0, Width: float32(tex.Width), Height: float32(tex.Height)}
		dest := rl.Rectangle{X: float32(vw.X), Y: float32(vw.Y), Width: float32(vw.W), Height: float32(vw.H)}
		origin := rl.Vector2{X: 0, Y: 0}
		rl.DrawTexturePro(tex, src, dest, origin, 0, rl.White)
		return nil, nil
	})

	// View getters (return 0 if view not found)
	v.RegisterForeign("GetViewX", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("GetViewX requires (viewId)")
		}
		viewsMu.RLock()
		vw, ok := views[toString(args[0])]
		viewsMu.RUnlock()
		if !ok {
			return 0, nil
		}
		return int(vw.X), nil
	})
	v.RegisterForeign("GetViewY", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("GetViewY requires (viewId)")
		}
		viewsMu.RLock()
		vw, ok := views[toString(args[0])]
		viewsMu.RUnlock()
		if !ok {
			return 0, nil
		}
		return int(vw.Y), nil
	})
	v.RegisterForeign("GetViewWidth", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("GetViewWidth requires (viewId)")
		}
		viewsMu.RLock()
		vw, ok := views[toString(args[0])]
		viewsMu.RUnlock()
		if !ok {
			return 0, nil
		}
		return int(vw.W), nil
	})
	v.RegisterForeign("GetViewHeight", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("GetViewHeight requires (viewId)")
		}
		viewsMu.RLock()
		vw, ok := views[toString(args[0])]
		viewsMu.RUnlock()
		if !ok {
			return 0, nil
		}
		return int(vw.H), nil
	})

	// View resize/position (for dynamic splitscreen)
	v.RegisterForeign("SetViewPosition", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("SetViewPosition requires (viewId, x, y)")
		}
		id := toString(args[0])
		viewsMu.Lock()
		defer viewsMu.Unlock()
		vw, ok := views[id]
		if !ok {
			return nil, fmt.Errorf("unknown view: %s", id)
		}
		vw.X, vw.Y = toInt32(args[1]), toInt32(args[2])
		return nil, nil
	})
	v.RegisterForeign("SetViewSize", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("SetViewSize requires (viewId, width, height)")
		}
		id := toString(args[0])
		viewsMu.Lock()
		defer viewsMu.Unlock()
		vw, ok := views[id]
		if !ok {
			return nil, fmt.Errorf("unknown view: %s", id)
		}
		vw.W, vw.H = toInt32(args[1]), toInt32(args[2])
		return nil, nil
	})
	v.RegisterForeign("SetViewRect", func(args []interface{}) (interface{}, error) {
		if len(args) < 5 {
			return nil, fmt.Errorf("SetViewRect requires (viewId, x, y, width, height)")
		}
		id := toString(args[0])
		viewsMu.Lock()
		defer viewsMu.Unlock()
		vw, ok := views[id]
		if !ok {
			return nil, fmt.Errorf("unknown view: %s", id)
		}
		vw.X, vw.Y = toInt32(args[1]), toInt32(args[2])
		vw.W, vw.H = toInt32(args[3]), toInt32(args[4])
		return nil, nil
	})

	// Splitscreen convenience: create two or four views using current screen size (call after InitWindow).
	v.RegisterForeign("CreateSplitscreenLeftRight", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("CreateSplitscreenLeftRight requires (viewIdLeft, viewIdRight)")
		}
		w, h := rl.GetScreenWidth(), rl.GetScreenHeight()
		halfW := int32(w / 2)
		viewsMu.Lock()
		views[toString(args[0])] = &viewState{X: 0, Y: 0, W: halfW, H: int32(h)}
		views[toString(args[1])] = &viewState{X: halfW, Y: 0, W: halfW, H: int32(h)}
		viewsMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("CreateSplitscreenTopBottom", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("CreateSplitscreenTopBottom requires (viewIdTop, viewIdBottom)")
		}
		w, h := rl.GetScreenWidth(), rl.GetScreenHeight()
		halfH := int32(h / 2)
		viewsMu.Lock()
		views[toString(args[0])] = &viewState{X: 0, Y: 0, W: int32(w), H: halfH}
		views[toString(args[1])] = &viewState{X: 0, Y: halfH, W: int32(w), H: halfH}
		viewsMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("CreateSplitscreenFour", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("CreateSplitscreenFour requires (viewIdTL, viewIdTR, viewIdBL, viewIdBR)")
		}
		w, h := rl.GetScreenWidth(), rl.GetScreenHeight()
		halfW, halfH := int32(w/2), int32(h/2)
		viewsMu.Lock()
		views[toString(args[0])] = &viewState{X: 0, Y: 0, W: halfW, H: halfH}
		views[toString(args[1])] = &viewState{X: halfW, Y: 0, W: halfW, H: halfH}
		views[toString(args[2])] = &viewState{X: 0, Y: halfH, W: halfW, H: halfH}
		views[toString(args[3])] = &viewState{X: halfW, Y: halfH, W: halfW, H: halfH}
		viewsMu.Unlock()
		return nil, nil
	})
}
