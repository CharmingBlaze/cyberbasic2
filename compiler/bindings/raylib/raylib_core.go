// Package raylib: core window, frame, and system (rcore).
package raylib

import (
	"cyberbasic/compiler/vm"
	"fmt"
	"os"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func registerCore(v *vm.VM) {
	v.RegisterForeign("InitWindow", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("InitWindow requires (width, height, title)")
		}
		rl.InitWindow(toInt32(args[0]), toInt32(args[1]), toString(args[2]))
		return nil, nil
	})
	v.RegisterForeign("SetTargetFPS", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("SetTargetFPS requires (fps)")
		}
		rl.SetTargetFPS(toInt32(args[0]))
		return nil, nil
	})
	v.RegisterForeign("WindowShouldClose", func(args []interface{}) (interface{}, error) {
		return rl.WindowShouldClose(), nil
	})
	v.RegisterForeign("CloseWindow", func(args []interface{}) (interface{}, error) {
		rl.CloseWindow()
		return nil, nil
	})
	v.RegisterForeign("SetWindowPosition", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("SetWindowPosition requires (x, y)")
		}
		rl.SetWindowPosition(int(toInt32(args[0])), int(toInt32(args[1])))
		return nil, nil
	})
	v.RegisterForeign("BeginDrawing", func(args []interface{}) (interface{}, error) {
		rl.BeginDrawing()
		return nil, nil
	})
	v.RegisterForeign("EndDrawing", func(args []interface{}) (interface{}, error) {
		rl.EndDrawing()
		return nil, nil
	})
	// BeginFrame(): alias for BeginDrawing (start frame)
	v.RegisterForeign("BeginFrame", func(args []interface{}) (interface{}, error) {
		rl.BeginDrawing()
		return nil, nil
	})
	// EndFrame(): alias for EndDrawing (end frame)
	v.RegisterForeign("EndFrame", func(args []interface{}) (interface{}, error) {
		rl.EndDrawing()
		return nil, nil
	})
	// SetUpdateFunction(func), SetDrawFunction(func): no-op (use WHILE NOT WindowShouldClose() ... WEND and call your update/draw logic manually).
	v.RegisterForeign("SetUpdateFunction", func(args []interface{}) (interface{}, error) { return nil, nil })
	v.RegisterForeign("SetDrawFunction", func(args []interface{}) (interface{}, error) { return nil, nil })
	// Run(): no-op; run your game loop with WHILE NOT WindowShouldClose() ... WEND.
	v.RegisterForeign("Run", func(args []interface{}) (interface{}, error) { return nil, nil })
	// Background(r, g, b): clear with RGB, alpha 255
	v.RegisterForeign("Background", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			rl.ClearBackground(rl.Black)
			return nil, nil
		}
		r, g, b := toInt32(args[0]), toInt32(args[1]), toInt32(args[2])
		rl.ClearBackground(rl.NewColor(uint8(r), uint8(g), uint8(b), 255))
		return nil, nil
	})
	v.RegisterForeign("ClearBackground", func(args []interface{}) (interface{}, error) {
		if len(args) == 0 {
			rl.ClearBackground(rl.Black)
			return nil, nil
		}
		if len(args) == 1 {
			// Single packed color (e.g. RL.Black, RL.DarkGray)
			switch v := args[0].(type) {
			case int:
				c := rl.NewColor(uint8(v>>16&0xff), uint8(v>>8&0xff), uint8(v&0xff), 255)
				rl.ClearBackground(c)
			case float64:
				u := uint32(v)
				c := rl.NewColor(uint8(u>>16&0xff), uint8(u>>8&0xff), uint8(u&0xff), 255)
				rl.ClearBackground(c)
			default:
				rl.ClearBackground(rl.Black)
			}
			return nil, nil
		}
		if len(args) >= 4 {
			r, g, b, a := toInt32(args[0]), toInt32(args[1]), toInt32(args[2]), toInt32(args[3])
			rl.ClearBackground(rl.NewColor(uint8(r), uint8(g), uint8(b), uint8(a)))
		}
		return nil, nil
	})
	v.RegisterForeign("GetFrameTime", func(args []interface{}) (interface{}, error) {
		return float64(rl.GetFrameTime()), nil
	})
	v.RegisterForeign("DeltaTime", func(args []interface{}) (interface{}, error) {
		return float64(rl.GetFrameTime()), nil
	})
	v.RegisterForeign("GetFPS", func(args []interface{}) (interface{}, error) {
		return int(rl.GetFPS()), nil
	})
	v.RegisterForeign("GetScreenWidth", func(args []interface{}) (interface{}, error) {
		return rl.GetScreenWidth(), nil
	})
	v.RegisterForeign("GetScreenHeight", func(args []interface{}) (interface{}, error) {
		return rl.GetScreenHeight(), nil
	})
	v.RegisterForeign("SetWindowSize", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("SetWindowSize requires (width, height)")
		}
		rl.SetWindowSize(int(toInt32(args[0])), int(toInt32(args[1])))
		return nil, nil
	})
	v.RegisterForeign("SetWindowTitle", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("SetWindowTitle requires (title)")
		}
		rl.SetWindowTitle(toString(args[0]))
		return nil, nil
	})
	v.RegisterForeign("MaximizeWindow", func(args []interface{}) (interface{}, error) {
		rl.MaximizeWindow()
		return nil, nil
	})
	v.RegisterForeign("MinimizeWindow", func(args []interface{}) (interface{}, error) {
		rl.MinimizeWindow()
		return nil, nil
	})
	v.RegisterForeign("IsWindowReady", func(args []interface{}) (interface{}, error) {
		return rl.IsWindowReady(), nil
	})
	v.RegisterForeign("IsWindowFullscreen", func(args []interface{}) (interface{}, error) {
		return rl.IsWindowFullscreen(), nil
	})
	v.RegisterForeign("GetTime", func(args []interface{}) (interface{}, error) {
		return rl.GetTime(), nil
	})
	v.RegisterForeign("GetRandomValue", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("GetRandomValue requires (min, max)")
		}
		min, max := toInt32(args[0]), toInt32(args[1])
		if max < min {
			min, max = max, min
		}
		n := max - min + 1
		if n <= 0 {
			return int(min), nil
		}
		return int(getRand().Int31n(n)) + int(min), nil
	})
	v.RegisterForeign("SetRandomSeed", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("SetRandomSeed requires (seed)")
		}
		seed := toInt32(args[0])
		setRandSeed(int64(seed))
		return nil, nil
	})
	v.RegisterForeign("SeedRandom", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("SeedRandom requires (seed)")
		}
		seed := toInt32(args[0])
		setRandSeed(int64(seed))
		return nil, nil
	})
	v.RegisterForeign("SetWindowState", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("SetWindowState requires (flags)")
		}
		rl.SetWindowState(uint32(toInt32(args[0])))
		return nil, nil
	})
	v.RegisterForeign("ClearWindowState", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("ClearWindowState requires (flags)")
		}
		rl.ClearWindowState(uint32(toInt32(args[0])))
		return nil, nil
	})
	v.RegisterForeign("GetMonitorCount", func(args []interface{}) (interface{}, error) {
		return rl.GetMonitorCount(), nil
	})
	v.RegisterForeign("GetCurrentMonitor", func(args []interface{}) (interface{}, error) {
		return rl.GetCurrentMonitor(), nil
	})
	v.RegisterForeign("GetClipboardText", func(args []interface{}) (interface{}, error) {
		return rl.GetClipboardText(), nil
	})
	v.RegisterForeign("SetClipboardText", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("SetClipboardText requires (text)")
		}
		rl.SetClipboardText(toString(args[0]))
		return nil, nil
	})
	v.RegisterForeign("TakeScreenshot", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("TakeScreenshot requires (fileName)")
		}
		rl.TakeScreenshot(toString(args[0]))
		return nil, nil
	})
	v.RegisterForeign("Screenshot", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("Screenshot requires (path)")
		}
		rl.TakeScreenshot(toString(args[0]))
		return nil, nil
	})
	v.RegisterForeign("IsFullscreen", func(args []interface{}) (interface{}, error) {
		return rl.IsWindowFullscreen(), nil
	})
	v.RegisterForeign("OpenURL", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("OpenURL requires (url)")
		}
		rl.OpenURL(toString(args[0]))
		return nil, nil
	})
	// Extra core from cheatsheet
	v.RegisterForeign("IsWindowHidden", func(args []interface{}) (interface{}, error) {
		return rl.IsWindowHidden(), nil
	})
	v.RegisterForeign("IsWindowMinimized", func(args []interface{}) (interface{}, error) {
		return rl.IsWindowMinimized(), nil
	})
	v.RegisterForeign("IsWindowMaximized", func(args []interface{}) (interface{}, error) {
		return rl.IsWindowMaximized(), nil
	})
	v.RegisterForeign("IsWindowFocused", func(args []interface{}) (interface{}, error) {
		return rl.IsWindowFocused(), nil
	})
	v.RegisterForeign("IsWindowResized", func(args []interface{}) (interface{}, error) {
		return rl.IsWindowResized(), nil
	})
	v.RegisterForeign("ToggleFullscreen", func(args []interface{}) (interface{}, error) {
		rl.ToggleFullscreen()
		return nil, nil
	})
	v.RegisterForeign("RestoreWindow", func(args []interface{}) (interface{}, error) {
		rl.RestoreWindow()
		return nil, nil
	})
	v.RegisterForeign("GetRenderWidth", func(args []interface{}) (interface{}, error) {
		return rl.GetFramebufferWidth(), nil
	})
	v.RegisterForeign("GetRenderHeight", func(args []interface{}) (interface{}, error) {
		return rl.GetFramebufferHeight(), nil
	})
	v.RegisterForeign("GetMonitorName", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("GetMonitorName requires (monitor)")
		}
		return rl.GetMonitorName(int(toInt32(args[0]))), nil
	})
	v.RegisterForeign("GetMonitorWidth", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("GetMonitorWidth requires (monitor)")
		}
		return rl.GetMonitorWidth(int(toInt32(args[0]))), nil
	})
	v.RegisterForeign("GetMonitorHeight", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("GetMonitorHeight requires (monitor)")
		}
		return rl.GetMonitorHeight(int(toInt32(args[0]))), nil
	})
	v.RegisterForeign("GetMonitorRefreshRate", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("GetMonitorRefreshRate requires (monitor)")
		}
		return rl.GetMonitorRefreshRate(int(toInt32(args[0]))), nil
	})
	v.RegisterForeign("WaitTime", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("WaitTime requires (seconds)")
		}
		rl.WaitTime(toFloat64(args[0]))
		return nil, nil
	})
	v.RegisterForeign("EnableEventWaiting", func(args []interface{}) (interface{}, error) {
		rl.EnableEventWaiting()
		return nil, nil
	})
	v.RegisterForeign("DisableEventWaiting", func(args []interface{}) (interface{}, error) {
		rl.DisableEventWaiting()
		return nil, nil
	})
	v.RegisterForeign("IsCursorHidden", func(args []interface{}) (interface{}, error) {
		return rl.IsCursorHidden(), nil
	})
	v.RegisterForeign("EnableCursor", func(args []interface{}) (interface{}, error) {
		rl.EnableCursor()
		return nil, nil
	})
	v.RegisterForeign("DisableCursor", func(args []interface{}) (interface{}, error) {
		rl.DisableCursor()
		return nil, nil
	})
	v.RegisterForeign("IsCursorOnScreen", func(args []interface{}) (interface{}, error) {
		return rl.IsCursorOnScreen(), nil
	})

	// Window state and options
	v.RegisterForeign("IsWindowState", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("IsWindowState requires (flag)")
		}
		return rl.IsWindowState(uint32(toInt32(args[0]))), nil
	})
	v.RegisterForeign("ToggleBorderlessWindowed", func(args []interface{}) (interface{}, error) {
		rl.ToggleBorderlessWindowed()
		return nil, nil
	})
	v.RegisterForeign("SetWindowMonitor", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("SetWindowMonitor requires (monitor)")
		}
		rl.SetWindowMonitor(int(toInt32(args[0])))
		return nil, nil
	})
	v.RegisterForeign("SetWindowMinSize", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("SetWindowMinSize requires (width, height)")
		}
		rl.SetWindowMinSize(int(toInt32(args[0])), int(toInt32(args[1])))
		return nil, nil
	})
	v.RegisterForeign("SetWindowMaxSize", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("SetWindowMaxSize requires (width, height)")
		}
		rl.SetWindowMaxSize(int(toInt32(args[0])), int(toInt32(args[1])))
		return nil, nil
	})
	v.RegisterForeign("SetWindowOpacity", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("SetWindowOpacity requires (opacity)")
		}
		rl.SetWindowOpacity(toFloat32(args[0]))
		return nil, nil
	})
	v.RegisterForeign("GetWindowPosition", func(args []interface{}) (interface{}, error) {
		pos := rl.GetWindowPosition()
		return []interface{}{float64(pos.X), float64(pos.Y)}, nil
	})
	v.RegisterForeign("GetWindowScaleDPI", func(args []interface{}) (interface{}, error) {
		scale := rl.GetWindowScaleDPI()
		return []interface{}{float64(scale.X), float64(scale.Y)}, nil
	})
	v.RegisterForeign("GetScaleDPI", func(args []interface{}) (interface{}, error) {
		scale := rl.GetWindowScaleDPI()
		avg := (float64(scale.X) + float64(scale.Y)) / 2
		return avg, nil
	})
	v.RegisterForeign("GetMonitorPosition", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("GetMonitorPosition requires (monitor)")
		}
		pos := rl.GetMonitorPosition(int(toInt32(args[0])))
		return []interface{}{float64(pos.X), float64(pos.Y)}, nil
	})
	v.RegisterForeign("GetMonitorPhysicalWidth", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("GetMonitorPhysicalWidth requires (monitor)")
		}
		return rl.GetMonitorPhysicalWidth(int(toInt32(args[0]))), nil
	})
	v.RegisterForeign("GetMonitorPhysicalHeight", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("GetMonitorPhysicalHeight requires (monitor)")
		}
		return rl.GetMonitorPhysicalHeight(int(toInt32(args[0]))), nil
	})

	// Frame control
	v.RegisterForeign("SetConfigFlags", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("SetConfigFlags requires (flags)")
		}
		rl.SetConfigFlags(uint32(toInt32(args[0])))
		return nil, nil
	})
	v.RegisterForeign("SwapScreenBuffer", func(args []interface{}) (interface{}, error) {
		rl.SwapScreenBuffer()
		return nil, nil
	})
	v.RegisterForeign("PollInputEvents", func(args []interface{}) (interface{}, error) {
		rl.PollInputEvents()
		if err := v.ProcessEvents(); err != nil {
			return nil, err
		}
		return nil, nil
	})

	// 2D mode (Camera2D: offsetX, offsetY, targetX, targetY, rotation, zoom)
	v.RegisterForeign("SetCamera2D", func(args []interface{}) (interface{}, error) {
		if len(args) < 6 {
			return nil, fmt.Errorf("SetCamera2D requires (offsetX, offsetY, targetX, targetY, rotation, zoom)")
		}
		camera2D.Offset = rl.Vector2{X: toFloat32(args[0]), Y: toFloat32(args[1])}
		camera2D.Target = rl.Vector2{X: toFloat32(args[2]), Y: toFloat32(args[3])}
		camera2D.Rotation = toFloat32(args[4])
		camera2D.Zoom = toFloat32(args[5])
		return nil, nil
	})
	// 2D camera by ID (Phase 1)
	v.RegisterForeign("Camera2DCreate", func(args []interface{}) (interface{}, error) {
		camera2DMu.Lock()
		cameras2DCounter++
		id := fmt.Sprintf("cam2d_%d", cameras2DCounter)
		w, h := float32(rl.GetScreenWidth()), float32(rl.GetScreenHeight())
		cameras2D[id] = rl.Camera2D{
			Offset:   rl.Vector2{X: w / 2, Y: h / 2},
			Target:   rl.Vector2{},
			Rotation: 0,
			Zoom:     1,
		}
		camera2DMu.Unlock()
		return id, nil
	})
	v.RegisterForeign("Camera2DSetPosition", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("Camera2DSetPosition requires (cameraID, x, y)")
		}
		id := toString(args[0])
		x, y := toFloat32(args[1]), toFloat32(args[2])
		camera2DMu.Lock()
		defer camera2DMu.Unlock()
		cam, ok := cameras2D[id]
		if !ok {
			return nil, fmt.Errorf("unknown 2D camera: %s", id)
		}
		cam.Target = rl.Vector2{X: x, Y: y}
		cameras2D[id] = cam
		return nil, nil
	})
	v.RegisterForeign("Camera2DSetZoom", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("Camera2DSetZoom requires (cameraID, zoom)")
		}
		id := toString(args[0])
		zoom := toFloat32(args[1])
		camera2DMu.Lock()
		defer camera2DMu.Unlock()
		cam, ok := cameras2D[id]
		if !ok {
			return nil, fmt.Errorf("unknown 2D camera: %s", id)
		}
		cam.Zoom = zoom
		cameras2D[id] = cam
		return nil, nil
	})
	v.RegisterForeign("Camera2DSetRotation", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("Camera2DSetRotation requires (cameraID, angle)")
		}
		id := toString(args[0])
		angle := toFloat32(args[1])
		camera2DMu.Lock()
		defer camera2DMu.Unlock()
		cam, ok := cameras2D[id]
		if !ok {
			return nil, fmt.Errorf("unknown 2D camera: %s", id)
		}
		cam.Rotation = angle
		cameras2D[id] = cam
		return nil, nil
	})
	v.RegisterForeign("Camera2DMove", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("Camera2DMove requires (cameraID, dx, dy)")
		}
		id := toString(args[0])
		dx, dy := toFloat32(args[1]), toFloat32(args[2])
		camera2DMu.Lock()
		defer camera2DMu.Unlock()
		cam, ok := cameras2D[id]
		if !ok {
			return nil, fmt.Errorf("unknown 2D camera: %s", id)
		}
		cam.Target.X += dx
		cam.Target.Y += dy
		cameras2D[id] = cam
		return nil, nil
	})
	v.RegisterForeign("Camera2DSmoothFollow", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("Camera2DSmoothFollow requires (cameraID, targetX, targetY, speed)")
		}
		id := toString(args[0])
		camera2DMu.Lock()
		camera2DFollowID = id
		camera2DFollowTgtX = toFloat32(args[1])
		camera2DFollowTgtY = toFloat32(args[2])
		camera2DFollowSpd = toFloat32(args[3])
		camera2DMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("BeginCamera2D", func(args []interface{}) (interface{}, error) {
		id := ""
		if len(args) >= 1 {
			id = toString(args[0])
		}
		camera2DMu.Lock()
		currentCamera2DID = id
		camera2DMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("EndCamera2D", func(args []interface{}) (interface{}, error) {
		camera2DMu.Lock()
		currentCamera2DID = ""
		camera2DMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("BeginMode2D", func(args []interface{}) (interface{}, error) {
		if len(args) >= 6 {
			camera2D.Offset = rl.Vector2{X: toFloat32(args[0]), Y: toFloat32(args[1])}
			camera2D.Target = rl.Vector2{X: toFloat32(args[2]), Y: toFloat32(args[3])}
			camera2D.Rotation = toFloat32(args[4])
			camera2D.Zoom = toFloat32(args[5])
		} else {
			// Default: 1:1 screen coords so automatic game-loop 2D works (Zoom=0 would render blank)
			camera2D.Offset = rl.Vector2{}
			camera2D.Target = rl.Vector2{}
			camera2D.Rotation = 0
			camera2D.Zoom = 1
		}
		rl.BeginMode2D(camera2D)
		return nil, nil
	})
	v.RegisterForeign("EndMode2D", func(args []interface{}) (interface{}, error) {
		rl.EndMode2D()
		return nil, nil
	})
	v.RegisterForeign("GetWorldToScreen2D", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("GetWorldToScreen2D requires (worldX, worldY)")
		}
		pos := rl.GetWorldToScreen2D(rl.Vector2{X: toFloat32(args[0]), Y: toFloat32(args[1])}, camera2D)
		return []interface{}{float64(pos.X), float64(pos.Y)}, nil
	})
	v.RegisterForeign("GetScreenToWorld2D", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("GetScreenToWorld2D requires (screenX, screenY)")
		}
		pos := rl.GetScreenToWorld2D(rl.Vector2{X: toFloat32(args[0]), Y: toFloat32(args[1])}, camera2D)
		return []interface{}{float64(pos.X), float64(pos.Y)}, nil
	})

	// Blend and scissor
	v.RegisterForeign("BeginBlendMode", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("BeginBlendMode requires (mode)")
		}
		rl.BeginBlendMode(rl.BlendMode(toInt32(args[0])))
		return nil, nil
	})
	v.RegisterForeign("EndBlendMode", func(args []interface{}) (interface{}, error) {
		rl.EndBlendMode()
		return nil, nil
	})
	v.RegisterForeign("BeginScissorMode", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("BeginScissorMode requires (x, y, width, height)")
		}
		rl.BeginScissorMode(toInt32(args[0]), toInt32(args[1]), toInt32(args[2]), toInt32(args[3]))
		return nil, nil
	})
	v.RegisterForeign("EndScissorMode", func(args []interface{}) (interface{}, error) {
		rl.EndScissorMode()
		return nil, nil
	})

	// Shader mode (use id from LoadShader)
	v.RegisterForeign("BeginShaderMode", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("BeginShaderMode requires (shaderId)")
		}
		id := toString(args[0])
		shaderMu.Lock()
		sh, ok := shaders[id]
		shaderMu.Unlock()
		if !ok {
			return nil, fmt.Errorf("unknown shader id: %s", id)
		}
		currentShaderMu.Lock()
		currentShaderId = id
		currentShaderMu.Unlock()
		rl.BeginShaderMode(sh)
		return nil, nil
	})
	v.RegisterForeign("EndShaderMode", func(args []interface{}) (interface{}, error) {
		currentShaderMu.Lock()
		currentShaderId = ""
		currentShaderMu.Unlock()
		rl.EndShaderMode()
		return nil, nil
	})
	v.RegisterForeign("ApplyShader", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("ApplyShader requires (shaderId)")
		}
		id := toString(args[0])
		shaderMu.Lock()
		sh, ok := shaders[id]
		shaderMu.Unlock()
		if !ok {
			return nil, fmt.Errorf("unknown shader id: %s", id)
		}
		currentShaderMu.Lock()
		currentShaderId = id
		currentShaderMu.Unlock()
		rl.BeginShaderMode(sh)
		return nil, nil
	})
	v.RegisterForeign("RemoveShader", func(args []interface{}) (interface{}, error) {
		currentShaderMu.Lock()
		currentShaderId = ""
		currentShaderMu.Unlock()
		rl.EndShaderMode()
		return nil, nil
	})
	v.RegisterForeign("SetShaderUniform", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("SetShaderUniform requires (shaderId, name, value)")
		}
		id := toString(args[0])
		name := toString(args[1])
		val := toFloat32(args[2])
		shaderMu.Lock()
		sh, ok := shaders[id]
		shaderMu.Unlock()
		if !ok {
			return nil, fmt.Errorf("unknown shader id: %s", id)
		}
		loc := rl.GetShaderLocation(sh, name)
		if loc >= 0 {
			rl.SetShaderValue(sh, loc, []float32{val}, rl.ShaderUniformFloat)
		}
		return nil, nil
	})
	// SetShaderValueMatrix(shaderId, uniformName, m0..m15): set shader uniform to 4x4 matrix (row-major: M0,M4,M8,M12, M1,M5,M9,M13, ...).
	v.RegisterForeign("SetShaderValueMatrix", func(args []interface{}) (interface{}, error) {
		if len(args) < 19 {
			return nil, fmt.Errorf("SetShaderValueMatrix requires (shaderId, uniformName, 16 matrix floats)")
		}
		id := toString(args[0])
		name := toString(args[1])
		shaderMu.Lock()
		sh, ok := shaders[id]
		shaderMu.Unlock()
		if !ok {
			return nil, fmt.Errorf("unknown shader id: %s", id)
		}
		loc := rl.GetShaderLocation(sh, name)
		if loc >= 0 {
			mat := argsToMatrix(args, 2)
			rl.SetShaderValueMatrix(sh, loc, mat)
		}
		return nil, nil
	})
	// SetShaderValueTexture(shaderId, uniformName, textureId): set shader uniform to texture (sampler2d).
	v.RegisterForeign("SetShaderValueTexture", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("SetShaderValueTexture requires (shaderId, uniformName, textureId)")
		}
		id := toString(args[0])
		name := toString(args[1])
		texId := toString(args[2])
		shaderMu.Lock()
		sh, ok := shaders[id]
		shaderMu.Unlock()
		if !ok {
			return nil, fmt.Errorf("unknown shader id: %s", id)
		}
		texMu.Lock()
		tex, okTex := textures[texId]
		texMu.Unlock()
		if !okTex {
			return nil, fmt.Errorf("unknown texture id: %s", texId)
		}
		loc := rl.GetShaderLocation(sh, name)
		if loc >= 0 {
			rl.SetShaderValueTexture(sh, loc, tex)
		}
		return nil, nil
	})
	v.RegisterForeign("LoadShader", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("LoadShader requires (vsFileName, fsFileName)")
		}
		sh := rl.LoadShader(toString(args[0]), toString(args[1]))
		shaderMu.Lock()
		shaderCounter++
		id := fmt.Sprintf("shader_%d", shaderCounter)
		shaders[id] = sh
		shaderMu.Unlock()
		return id, nil
	})
	v.RegisterForeign("LoadShaderFromMemory", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("LoadShaderFromMemory requires (vsCode, fsCode)")
		}
		sh := rl.LoadShaderFromMemory(toString(args[0]), toString(args[1]))
		shaderMu.Lock()
		shaderCounter++
		id := fmt.Sprintf("shader_%d", shaderCounter)
		shaders[id] = sh
		shaderMu.Unlock()
		return id, nil
	})
	v.RegisterForeign("UnloadShader", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("UnloadShader requires (shaderId)")
		}
		id := toString(args[0])
		shaderMu.Lock()
		sh, ok := shaders[id]
		delete(shaders, id)
		shaderMu.Unlock()
		if ok {
			rl.UnloadShader(sh)
		}
		return nil, nil
	})
	v.RegisterForeign("IsShaderValid", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("IsShaderValid requires (shaderId)")
		}
		id := toString(args[0])
		shaderMu.Lock()
		sh, ok := shaders[id]
		shaderMu.Unlock()
		if !ok {
			return false, nil
		}
		return rl.IsShaderValid(sh), nil
	})

	// File / utils (rcore-style; raylib has in utils)
	v.RegisterForeign("FileExists", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("FileExists requires (fileName)")
		}
		_, err := os.Stat(toString(args[0]))
		return err == nil, nil
	})
}
