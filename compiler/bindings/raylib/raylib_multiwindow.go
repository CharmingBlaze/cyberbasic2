// Package raylib: multi-window system (logical windows as viewports + render textures).
// WindowCreate, WindowClose, message/channel/state, events, docking, RPC.
package raylib

import (
	"cyberbasic/compiler/vm"
	"fmt"
	"strings"
	"sync"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Window type constants (stored as string: "normal", "popup", "modal", "tool", "child")
const (
	winTypeNormal = "normal"
	winTypePopup  = "popup"
	winTypeModal  = "modal"
	winTypeTool   = "tool"
	winTypeChild  = "child"
)

type windowMessage struct {
	Message string
	Data    interface{}
}

type windowEventHandlers struct {
	Update  string
	Draw    string
	Resize  string
	Close   string
	Message string
}

type windowState struct {
	ID               int
	Title            string
	X, Y, W, H       int32
	Visible          bool
	Focused          bool
	Type             string
	ParentID         int
	RenderTextureID  string
	MessageQueue    []windowMessage
	Handlers        windowEventHandlers
	RegisteredFuncs map[string]string // name -> Sub name
	CameraID        string            // for 3D
}

var (
	multiWindowMu     sync.RWMutex
	windows           = make(map[int]*windowState)
	multiWindowNextID = 1
	currentDrawWindow int = -1

	// Channels: name -> queue of values
	channelMu   sync.Mutex
	channels   = make(map[string][]interface{})

	// State: key -> value
	stateMu sync.Mutex
	state   = make(map[string]interface{})

	// Docking (optional)
	dockMu    sync.Mutex
	dockNodes = make(map[int]*dockNode)
	dockNextID = 1
)

type dockNode struct {
	ID         int
	X, Y, W, H int32
	Direction  string // "horizontal", "vertical"
	Size       float32
	ParentID   int
	ChildA     int
	ChildB     int
	WindowID   int
}

func windowByID(id int) *windowState {
	multiWindowMu.RLock()
	defer multiWindowMu.RUnlock()
	return windows[id]
}

func createWindow(w, h int32, title, winType string, parentID int) (int, error) {
	rt := rl.LoadRenderTexture(w, h)
	renderTexMu.Lock()
	renderTexCounter++
	rtID := fmt.Sprintf("rt_mw_%d", renderTexCounter)
	renderTextures[rtID] = rt
	renderTexMu.Unlock()

	multiWindowMu.Lock()
	id := multiWindowNextID
	multiWindowNextID++
	windows[id] = &windowState{
		ID:              id,
		Title:           title,
		X:               0,
		Y:               0,
		W:               w,
		H:               h,
		Visible:         true,
		Focused:         false,
		Type:            winType,
		ParentID:        parentID,
		RenderTextureID: rtID,
		MessageQueue:    nil,
		Handlers:        windowEventHandlers{},
		RegisteredFuncs: make(map[string]string),
	}
	multiWindowMu.Unlock()
	return id, nil
}

func registerMultiWindow(v *vm.VM) {
	// --- Window creation ---
	v.RegisterForeign("WindowCreate", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("WindowCreate requires (width, height, title)")
		}
		w, h := toInt32(args[0]), toInt32(args[1])
		title := toString(args[2])
		id, err := createWindow(w, h, title, winTypeNormal, 0)
		if err != nil {
			return nil, err
		}
		return id, nil
	})
	v.RegisterForeign("WindowClose", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("WindowClose requires (id)")
		}
		id := toInt32(args[0])
		multiWindowMu.Lock()
		defer multiWindowMu.Unlock()
		win, ok := windows[int(id)]
		if !ok {
			return nil, nil
		}
		if win.RenderTextureID != "" {
			renderTexMu.Lock()
			rt, okRT := renderTextures[win.RenderTextureID]
			delete(renderTextures, win.RenderTextureID)
			renderTexMu.Unlock()
			if okRT {
				rl.UnloadRenderTexture(rt)
			}
		}
		delete(windows, int(id))
		return nil, nil
	})
	v.RegisterForeign("WindowSetTitle", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("WindowSetTitle requires (id, title)")
		}
		id := int(toInt32(args[0]))
		win := windowByID(id)
		if win == nil {
			return nil, nil
		}
		multiWindowMu.Lock()
		win.Title = toString(args[1])
		multiWindowMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("WindowSetSize", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("WindowSetSize requires (id, width, height)")
		}
		id := int(toInt32(args[0]))
		w, h := toInt32(args[1]), toInt32(args[2])
		multiWindowMu.Lock()
		defer multiWindowMu.Unlock()
		win, ok := windows[id]
		if !ok {
			return nil, nil
		}
		if win.RenderTextureID != "" {
			renderTexMu.Lock()
			rt, okRT := renderTextures[win.RenderTextureID]
			delete(renderTextures, win.RenderTextureID)
			renderTexMu.Unlock()
			if okRT {
				rl.UnloadRenderTexture(rt)
			}
		}
		rt := rl.LoadRenderTexture(w, h)
		renderTexMu.Lock()
		renderTexCounter++
		rtID := fmt.Sprintf("rt_mw_%d", renderTexCounter)
		renderTextures[rtID] = rt
		renderTexMu.Unlock()
		win.W, win.H = w, h
		win.RenderTextureID = rtID
		return nil, nil
	})
	v.RegisterForeign("WindowSetPosition", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("WindowSetPosition requires (id, x, y)")
		}
		id := int(toInt32(args[0]))
		win := windowByID(id)
		if win == nil {
			return nil, nil
		}
		multiWindowMu.Lock()
		win.X, win.Y = toInt32(args[1]), toInt32(args[2])
		multiWindowMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("WindowFocus", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("WindowFocus requires (id)")
		}
		id := int(toInt32(args[0]))
		multiWindowMu.Lock()
		for _, w := range windows {
			w.Focused = w.ID == id
		}
		multiWindowMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("WindowIsOpen", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("WindowIsOpen requires (id)")
		}
		id := int(toInt32(args[0]))
		multiWindowMu.RLock()
		_, ok := windows[id]
		multiWindowMu.RUnlock()
		return ok, nil
	})
	v.RegisterForeign("WindowGetWidth", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("WindowGetWidth requires (id)")
		}
		win := windowByID(int(toInt32(args[0])))
		if win == nil {
			return 0, nil
		}
		return int(win.W), nil
	})
	v.RegisterForeign("WindowGetHeight", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("WindowGetHeight requires (id)")
		}
		win := windowByID(int(toInt32(args[0])))
		if win == nil {
			return 0, nil
		}
		return int(win.H), nil
	})
	v.RegisterForeign("WindowGetPositionX", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("WindowGetPositionX requires (id)")
		}
		win := windowByID(int(toInt32(args[0])))
		if win == nil {
			return 0, nil
		}
		return int(win.X), nil
	})
	v.RegisterForeign("WindowGetPositionY", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("WindowGetPositionY requires (id)")
		}
		win := windowByID(int(toInt32(args[0])))
		if win == nil {
			return 0, nil
		}
		return int(win.Y), nil
	})
	v.RegisterForeign("WindowGetPosition", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("WindowGetPosition requires (id)")
		}
		win := windowByID(int(toInt32(args[0])))
		if win == nil {
			return "0,0", nil
		}
		return fmt.Sprintf("%d,%d", win.X, win.Y), nil
	})
	v.RegisterForeign("WindowIsFocused", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("WindowIsFocused requires (id)")
		}
		win := windowByID(int(toInt32(args[0])))
		if win == nil {
			return false, nil
		}
		return win.Focused, nil
	})
	v.RegisterForeign("WindowIsVisible", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("WindowIsVisible requires (id)")
		}
		win := windowByID(int(toInt32(args[0])))
		if win == nil {
			return false, nil
		}
		return win.Visible, nil
	})
	v.RegisterForeign("WindowShow", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("WindowShow requires (id)")
		}
		win := windowByID(int(toInt32(args[0])))
		if win == nil {
			return nil, nil
		}
		multiWindowMu.Lock()
		win.Visible = true
		multiWindowMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("WindowHide", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("WindowHide requires (id)")
		}
		win := windowByID(int(toInt32(args[0])))
		if win == nil {
			return nil, nil
		}
		multiWindowMu.Lock()
		win.Visible = false
		multiWindowMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("WindowCreatePopup", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("WindowCreatePopup requires (width, height, title)")
		}
		w, h := toInt32(args[0]), toInt32(args[1])
		title := toString(args[2])
		return createWindow(w, h, title, winTypePopup, 0)
	})
	v.RegisterForeign("WindowCreateModal", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("WindowCreateModal requires (width, height, title)")
		}
		w, h := toInt32(args[0]), toInt32(args[1])
		title := toString(args[2])
		return createWindow(w, h, title, winTypeModal, 0)
	})
	v.RegisterForeign("WindowCreateToolWindow", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("WindowCreateToolWindow requires (width, height, title)")
		}
		w, h := toInt32(args[0]), toInt32(args[1])
		title := toString(args[2])
		return createWindow(w, h, title, winTypeTool, 0)
	})
	v.RegisterForeign("WindowCreateChild", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("WindowCreateChild requires (parentID, width, height, title)")
		}
		parentID := int(toInt32(args[0]))
		w, h := toInt32(args[1]), toInt32(args[2])
		title := toString(args[3])
		id, err := createWindow(w, h, title, winTypeChild, parentID)
		if err != nil {
			return nil, err
		}
		return id, nil
	})

	// --- Multi-window rendering ---
	v.RegisterForeign("WindowBeginDrawing", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("WindowBeginDrawing requires (id)")
		}
		id := int(toInt32(args[0]))
		if id == 0 {
			currentDrawWindow = 0
			return nil, nil
		}
		win := windowByID(id)
		if win == nil || win.RenderTextureID == "" {
			return nil, fmt.Errorf("unknown window or no render texture: %d", id)
		}
		renderTexMu.Lock()
		rt, ok := renderTextures[win.RenderTextureID]
		renderTexMu.Unlock()
		if !ok {
			return nil, fmt.Errorf("render texture not found for window %d", id)
		}
		rl.BeginTextureMode(rt)
		currentDrawWindow = id
		return nil, nil
	})
	v.RegisterForeign("WindowEndDrawing", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("WindowEndDrawing requires (id)")
		}
		id := int(toInt32(args[0]))
		if id != 0 && currentDrawWindow == id {
			rl.EndTextureMode()
		}
		currentDrawWindow = -1
		return nil, nil
	})
	v.RegisterForeign("WindowClearBackground", func(args []interface{}) (interface{}, error) {
		if len(args) < 5 {
			return nil, fmt.Errorf("WindowClearBackground requires (id, r, g, b, a)")
		}
		id := int(toInt32(args[0]))
		r, g, b, a := toInt32(args[1]), toInt32(args[2]), toInt32(args[3]), toInt32(args[4])
		c := rl.NewColor(uint8(r), uint8(g), uint8(b), uint8(a))
		if id == 0 {
			rl.ClearBackground(c)
			return nil, nil
		}
		if currentDrawWindow == id {
			rl.ClearBackground(c)
		}
		return nil, nil
	})
	v.RegisterForeign("WindowDrawAllToScreen", func(args []interface{}) (interface{}, error) {
		multiWindowMu.RLock()
		winList := make([]*windowState, 0, len(windows))
		for _, w := range windows {
			if w.Visible && w.RenderTextureID != "" {
				winList = append(winList, w)
			}
		}
		multiWindowMu.RUnlock()
		for _, w := range winList {
			renderTexMu.Lock()
			rt, ok := renderTextures[w.RenderTextureID]
			renderTexMu.Unlock()
			if !ok {
				continue
			}
			tex := rt.Texture
			src := rl.Rectangle{X: 0, Y: 0, Width: float32(tex.Width), Height: float32(tex.Height)}
			dest := rl.Rectangle{X: float32(w.X), Y: float32(w.Y), Width: float32(w.W), Height: float32(w.H)}
			rl.DrawTexturePro(tex, src, dest, rl.Vector2{}, 0, rl.White)
		}
		return nil, nil
	})

	// --- Message passing ---
	v.RegisterForeign("WindowSendMessage", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("WindowSendMessage requires (targetID, message, data)")
		}
		targetID := int(toInt32(args[0]))
		msg := toString(args[1])
		data := args[2]
		multiWindowMu.Lock()
		defer multiWindowMu.Unlock()
		win, ok := windows[targetID]
		if !ok {
			return nil, nil
		}
		win.MessageQueue = append(win.MessageQueue, windowMessage{Message: msg, Data: data})
		return nil, nil
	})
	v.RegisterForeign("WindowBroadcast", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("WindowBroadcast requires (message, data)")
		}
		msg := toString(args[0])
		data := args[1]
		multiWindowMu.Lock()
		defer multiWindowMu.Unlock()
		for _, win := range windows {
			win.MessageQueue = append(win.MessageQueue, windowMessage{Message: msg, Data: data})
		}
		return nil, nil
	})
	v.RegisterForeign("WindowReceiveMessage", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("WindowReceiveMessage requires (id)")
		}
		id := int(toInt32(args[0]))
		multiWindowMu.Lock()
		defer multiWindowMu.Unlock()
		win, ok := windows[id]
		if !ok || len(win.MessageQueue) == 0 {
			return nil, nil
		}
		item := win.MessageQueue[0]
		win.MessageQueue = win.MessageQueue[1:]
		return fmt.Sprintf("%s|%v", item.Message, item.Data), nil
	})
	v.RegisterForeign("WindowHasMessage", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("WindowHasMessage requires (id)")
		}
		id := int(toInt32(args[0]))
		multiWindowMu.RLock()
		win, ok := windows[id]
		has := ok && len(win.MessageQueue) > 0
		multiWindowMu.RUnlock()
		return has, nil
	})

	// --- Channels ---
	v.RegisterForeign("ChannelCreate", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("ChannelCreate requires (name)")
		}
		name := toString(args[0])
		channelMu.Lock()
		if _, ok := channels[name]; !ok {
			channels[name] = nil
		}
		channelMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("ChannelSend", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("ChannelSend requires (name, data)")
		}
		name := toString(args[0])
		data := args[1]
		channelMu.Lock()
		channels[name] = append(channels[name], data)
		channelMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("ChannelReceive", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("ChannelReceive requires (name)")
		}
		name := toString(args[0])
		channelMu.Lock()
		defer channelMu.Unlock()
		q, ok := channels[name]
		if !ok || len(q) == 0 {
			return nil, nil
		}
		val := q[0]
		channels[name] = q[1:]
		return val, nil
	})
	v.RegisterForeign("ChannelHasData", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("ChannelHasData requires (name)")
		}
		name := toString(args[0])
		channelMu.Lock()
		q := channels[name]
		has := len(q) > 0
		channelMu.Unlock()
		return has, nil
	})

	// --- Shared state ---
	v.RegisterForeign("StateSet", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("StateSet requires (key, value)")
		}
		key := toString(args[0])
		stateMu.Lock()
		state[key] = args[1]
		stateMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("StateGet", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("StateGet requires (key)")
		}
		key := toString(args[0])
		stateMu.Lock()
		val := state[key]
		stateMu.Unlock()
		return val, nil
	})
	v.RegisterForeign("StateHas", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("StateHas requires (key)")
		}
		key := toString(args[0])
		stateMu.Lock()
		_, ok := state[key]
		stateMu.Unlock()
		return ok, nil
	})
	v.RegisterForeign("StateRemove", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("StateRemove requires (key)")
		}
		key := toString(args[0])
		stateMu.Lock()
		delete(state, key)
		stateMu.Unlock()
		return nil, nil
	})

	// --- Window events ---
	v.RegisterForeign("OnWindowUpdate", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("OnWindowUpdate requires (id, function)")
		}
		id := int(toInt32(args[0]))
		subName := toString(args[1])
		multiWindowMu.Lock()
		defer multiWindowMu.Unlock()
		win, ok := windows[id]
		if !ok {
			return nil, nil
		}
		win.Handlers.Update = subName
		return nil, nil
	})
	v.RegisterForeign("OnWindowDraw", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("OnWindowDraw requires (id, function)")
		}
		id := int(toInt32(args[0]))
		subName := toString(args[1])
		multiWindowMu.Lock()
		defer multiWindowMu.Unlock()
		win, ok := windows[id]
		if !ok {
			return nil, nil
		}
		win.Handlers.Draw = subName
		return nil, nil
	})
	v.RegisterForeign("OnWindowResize", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("OnWindowResize requires (id, function)")
		}
		id := int(toInt32(args[0]))
		subName := toString(args[1])
		multiWindowMu.Lock()
		defer multiWindowMu.Unlock()
		win, ok := windows[id]
		if !ok {
			return nil, nil
		}
		win.Handlers.Resize = subName
		return nil, nil
	})
	v.RegisterForeign("OnWindowClose", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("OnWindowClose requires (id, function)")
		}
		id := int(toInt32(args[0]))
		subName := toString(args[1])
		multiWindowMu.Lock()
		defer multiWindowMu.Unlock()
		win, ok := windows[id]
		if !ok {
			return nil, nil
		}
		win.Handlers.Close = subName
		return nil, nil
	})
	v.RegisterForeign("OnWindowMessage", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("OnWindowMessage requires (id, function)")
		}
		id := int(toInt32(args[0]))
		subName := toString(args[1])
		multiWindowMu.Lock()
		defer multiWindowMu.Unlock()
		win, ok := windows[id]
		if !ok {
			return nil, nil
		}
		win.Handlers.Message = subName
		return nil, nil
	})
	v.RegisterForeign("WindowProcessEvents", func(args []interface{}) (interface{}, error) {
		multiWindowMu.RLock()
		winList := make([]*windowState, 0, len(windows))
		for _, w := range windows {
			winList = append(winList, w)
		}
		multiWindowMu.RUnlock()
		for _, w := range winList {
			if w.Handlers.Update != "" {
				_ = v.InvokeSub(w.Handlers.Update, []interface{}{w.ID})
			}
		}
		return nil, nil
	})
	v.RegisterForeign("WindowDraw", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("WindowDraw requires (id)")
		}
		id := int(toInt32(args[0]))
		win := windowByID(id)
		if win == nil {
			return nil, nil
		}
		if win.Handlers.Draw == "" {
			return nil, nil
		}
		if id != 0 && win.RenderTextureID != "" {
			renderTexMu.Lock()
			rt, ok := renderTextures[win.RenderTextureID]
			renderTexMu.Unlock()
			if ok {
				rl.BeginTextureMode(rt)
				currentDrawWindow = id
				_ = v.InvokeSub(win.Handlers.Draw, []interface{}{id})
				rl.EndTextureMode()
				currentDrawWindow = -1
			}
		} else if id == 0 {
			currentDrawWindow = 0
			_ = v.InvokeSub(win.Handlers.Draw, []interface{}{0})
			currentDrawWindow = -1
		}
		return nil, nil
	})

	// --- 3D: WindowSetCamera, WindowDrawModel, WindowDrawScene ---
	v.RegisterForeign("WindowSetCamera", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("WindowSetCamera requires (id, cameraId)")
		}
		id := int(toInt32(args[0]))
		camID := toString(args[1])
		multiWindowMu.Lock()
		defer multiWindowMu.Unlock()
		win, ok := windows[id]
		if !ok {
			return nil, nil
		}
		win.CameraID = camID
		return nil, nil
	})
	v.RegisterForeign("WindowDrawModel", func(args []interface{}) (interface{}, error) {
		if len(args) < 5 {
			return nil, fmt.Errorf("WindowDrawModel requires (id, modelId, x, y, z)")
		}
		id := int(toInt32(args[0]))
		if currentDrawWindow != id {
			return nil, nil
		}
		modelID := toString(args[1])
		x, y, z := toFloat32(args[2]), toFloat32(args[3]), toFloat32(args[4])
		modelMu.Lock()
		model, ok := models[modelID]
		modelMu.Unlock()
		if !ok {
			return nil, fmt.Errorf("unknown model: %s", modelID)
		}
		rl.DrawModel(model, rl.Vector3{X: x, Y: y, Z: z}, 1, rl.White)
		return nil, nil
	})
	v.RegisterForeign("WindowDrawScene", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("WindowDrawScene requires (id, sceneID)")
		}
		id := int(toInt32(args[0]))
		if currentDrawWindow != id {
			return nil, nil
		}
		sceneID := toString(args[1])
		// Scene drawing would use scene binding; for now no-op if scene package not called from here
		_ = sceneID
		return nil, nil
	})

	// --- RPC ---
	v.RegisterForeign("WindowRegisterFunction", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("WindowRegisterFunction requires (id, name, function)")
		}
		id := int(toInt32(args[0]))
		name := toString(args[1])
		subName := toString(args[2])
		multiWindowMu.Lock()
		defer multiWindowMu.Unlock()
		win, ok := windows[id]
		if !ok {
			return nil, nil
		}
		win.RegisteredFuncs[name] = subName
		return nil, nil
	})
	v.RegisterForeign("WindowCall", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("WindowCall requires (targetID, functionName, ...args)")
		}
		targetID := int(toInt32(args[0]))
		funcName := toString(args[1])
		multiWindowMu.RLock()
		win, ok := windows[targetID]
		var subName string
		if ok {
			subName = win.RegisteredFuncs[funcName]
		}
		multiWindowMu.RUnlock()
		if subName == "" {
			return nil, nil
		}
		callArgs := make([]interface{}, 0, len(args)-2)
		for i := 2; i < len(args); i++ {
			callArgs = append(callArgs, args[i])
		}
		_ = v.InvokeSub(subName, callArgs)
		return nil, nil
	})

	// --- Docking ---
	v.RegisterForeign("DockCreateArea", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("DockCreateArea requires (id)")
		}
		id := int(toInt32(args[0]))
		dockMu.Lock()
		defer dockMu.Unlock()
		dockNodes[id] = &dockNode{ID: id}
		return nil, nil
	})
	v.RegisterForeign("DockSplit", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("DockSplit requires (id, direction)")
		}
		id := int(toInt32(args[0]))
		dir := strings.ToLower(toString(args[1]))
		if dir != "horizontal" && dir != "vertical" {
			return nil, fmt.Errorf("DockSplit direction must be horizontal or vertical")
		}
		dockMu.Lock()
		defer dockMu.Unlock()
		node, ok := dockNodes[id]
		if !ok {
			return nil, fmt.Errorf("unknown dock id: %d", id)
		}
		node.Direction = dir
		node.Size = 0.5
		dockNextID++
		node.ChildA = dockNextID
		dockNodes[dockNextID] = &dockNode{ID: dockNextID, ParentID: id}
		dockNextID++
		node.ChildB = dockNextID
		dockNodes[dockNextID] = &dockNode{ID: dockNextID, ParentID: id}
		return nil, nil
	})
	v.RegisterForeign("DockAttachWindow", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("DockAttachWindow requires (dockID, windowID)")
		}
		dockID := int(toInt32(args[0]))
		windowID := int(toInt32(args[1]))
		dockMu.Lock()
		node, ok := dockNodes[dockID]
		dockMu.Unlock()
		if !ok {
			return nil, fmt.Errorf("unknown dock id: %d", dockID)
		}
		node.WindowID = windowID
		multiWindowMu.Lock()
		win, ok := windows[windowID]
		if ok {
			win.X, win.Y = node.X, node.Y
			win.W, win.H = node.W, node.H
		}
		multiWindowMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("DockSetSize", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("DockSetSize requires (dockID, size)")
		}
		id := int(toInt32(args[0]))
		size := toFloat32(args[1])
		if size <= 0 || size >= 1 {
			return nil, fmt.Errorf("DockSetSize size must be between 0 and 1")
		}
		dockMu.Lock()
		defer dockMu.Unlock()
		node, ok := dockNodes[id]
		if !ok {
			return nil, fmt.Errorf("unknown dock id: %d", id)
		}
		node.Size = size
		return nil, nil
	})
}
