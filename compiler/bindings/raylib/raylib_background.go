// Package raylib: 2D background system (static, scrolling, parallax, tiled, multi-layer).
package raylib

import (
	"fmt"
	"sync"

	"cyberbasic/compiler/vm"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type backgroundLayer struct {
	TextureID  string
	ParallaxX  float32
	ParallaxY  float32
}

type backgroundState struct {
	TextureID string
	ColorR    float32
	ColorG    float32
	ColorB    float32
	ColorA    float32
	ScrollX   float32
	ScrollY   float32
	OffsetX   float32
	OffsetY   float32
	ParallaxX float32
	ParallaxY float32
	Tiled     bool
	TileW     int
	TileH     int
	Layers    []backgroundLayer
}

var (
	backgrounds     = make(map[string]*backgroundState)
	backgroundSeq   int
	backgroundMu    sync.RWMutex
)

func registerBackground(v *vm.VM) {
	v.RegisterForeign("BackgroundCreate", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("BackgroundCreate requires (textureId)")
		}
		texID := toString(args[0])
		texMu.Lock()
		_, ok := textures[texID]
		texMu.Unlock()
		if !ok {
			return nil, fmt.Errorf("unknown texture id: %s", texID)
		}
		backgroundMu.Lock()
		backgroundSeq++
		id := fmt.Sprintf("bg_%d", backgroundSeq)
		backgrounds[id] = &backgroundState{
			TextureID: texID,
			ColorR:    1, ColorG: 1, ColorB: 1, ColorA: 1,
			ParallaxX: 1, ParallaxY: 1,
			TileW:     0, TileH: 0,
		}
		backgroundMu.Unlock()
		return id, nil
	})
	v.RegisterForeign("BackgroundSetColor", func(args []interface{}) (interface{}, error) {
		if len(args) < 5 {
			return nil, fmt.Errorf("BackgroundSetColor requires (backgroundId, r, g, b, a)")
		}
		id := toString(args[0])
		backgroundMu.Lock()
		defer backgroundMu.Unlock()
		bg := backgrounds[id]
		if bg == nil {
			return nil, fmt.Errorf("unknown background: %s", id)
		}
		bg.ColorR = toFloat32(args[1])
		bg.ColorG = toFloat32(args[2])
		bg.ColorB = toFloat32(args[3])
		bg.ColorA = toFloat32(args[4])
		return nil, nil
	})
	v.RegisterForeign("BackgroundSetTexture", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("BackgroundSetTexture requires (backgroundId, textureId)")
		}
		id := toString(args[0])
		texID := toString(args[1])
		texMu.Lock()
		_, ok := textures[texID]
		texMu.Unlock()
		if !ok {
			return nil, fmt.Errorf("unknown texture id: %s", texID)
		}
		backgroundMu.Lock()
		defer backgroundMu.Unlock()
		bg := backgrounds[id]
		if bg == nil {
			return nil, fmt.Errorf("unknown background: %s", id)
		}
		bg.TextureID = texID
		return nil, nil
	})
	v.RegisterForeign("BackgroundSetScroll", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("BackgroundSetScroll requires (backgroundId, speedX, speedY)")
		}
		id := toString(args[0])
		backgroundMu.Lock()
		defer backgroundMu.Unlock()
		bg := backgrounds[id]
		if bg == nil {
			return nil, fmt.Errorf("unknown background: %s", id)
		}
		bg.ScrollX = toFloat32(args[1])
		bg.ScrollY = toFloat32(args[2])
		return nil, nil
	})
	v.RegisterForeign("BackgroundSetOffset", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("BackgroundSetOffset requires (backgroundId, offsetX, offsetY)")
		}
		id := toString(args[0])
		backgroundMu.Lock()
		defer backgroundMu.Unlock()
		bg := backgrounds[id]
		if bg == nil {
			return nil, fmt.Errorf("unknown background: %s", id)
		}
		bg.OffsetX = toFloat32(args[1])
		bg.OffsetY = toFloat32(args[2])
		return nil, nil
	})
	v.RegisterForeign("BackgroundSetParallax", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("BackgroundSetParallax requires (backgroundId, px, py)")
		}
		id := toString(args[0])
		backgroundMu.Lock()
		defer backgroundMu.Unlock()
		bg := backgrounds[id]
		if bg == nil {
			return nil, fmt.Errorf("unknown background: %s", id)
		}
		bg.ParallaxX = toFloat32(args[1])
		bg.ParallaxY = toFloat32(args[2])
		return nil, nil
	})
	v.RegisterForeign("BackgroundSetTiled", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("BackgroundSetTiled requires (backgroundId, flag)")
		}
		id := toString(args[0])
		backgroundMu.Lock()
		defer backgroundMu.Unlock()
		bg := backgrounds[id]
		if bg == nil {
			return nil, fmt.Errorf("unknown background: %s", id)
		}
		bg.Tiled = toFloat32(args[1]) != 0
		return nil, nil
	})
	v.RegisterForeign("BackgroundSetTileSize", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("BackgroundSetTileSize requires (backgroundId, width, height)")
		}
		id := toString(args[0])
		backgroundMu.Lock()
		defer backgroundMu.Unlock()
		bg := backgrounds[id]
		if bg == nil {
			return nil, fmt.Errorf("unknown background: %s", id)
		}
		bg.TileW = int(toFloat32(args[1]))
		bg.TileH = int(toFloat32(args[2]))
		return nil, nil
	})
	v.RegisterForeign("BackgroundAddLayer", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("BackgroundAddLayer requires (backgroundId, textureId, parallaxX, parallaxY)")
		}
		id := toString(args[0])
		texID := toString(args[1])
		texMu.Lock()
		_, ok := textures[texID]
		texMu.Unlock()
		if !ok {
			return nil, fmt.Errorf("unknown texture id: %s", texID)
		}
		backgroundMu.Lock()
		defer backgroundMu.Unlock()
		bg := backgrounds[id]
		if bg == nil {
			return nil, fmt.Errorf("unknown background: %s", id)
		}
		bg.Layers = append(bg.Layers, backgroundLayer{
			TextureID: texID,
			ParallaxX: toFloat32(args[2]),
			ParallaxY: toFloat32(args[3]),
		})
		return nil, nil
	})
	v.RegisterForeign("BackgroundRemoveLayer", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("BackgroundRemoveLayer requires (backgroundId, layerIndex)")
		}
		id := toString(args[0])
		idx := int(toFloat32(args[1]))
		backgroundMu.Lock()
		defer backgroundMu.Unlock()
		bg := backgrounds[id]
		if bg == nil {
			return nil, fmt.Errorf("unknown background: %s", id)
		}
		if idx < 0 || idx >= len(bg.Layers) {
			return nil, fmt.Errorf("layer index out of range")
		}
		bg.Layers = append(bg.Layers[:idx], bg.Layers[idx+1:]...)
		return nil, nil
	})
	v.RegisterForeign("DrawBackground", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("DrawBackground requires (backgroundId)")
		}
		id := toString(args[0])
		DrawBackground(id, getCurrentCamera2D())
		return nil, nil
	})
}

// DrawBackground draws a background (used during 2D flush). Call with current camera for parallax.
func DrawBackground(backgroundID string, cam rl.Camera2D) {
	backgroundMu.RLock()
	bg := backgrounds[backgroundID]
	backgroundMu.RUnlock()
	if bg == nil {
		return
	}
	texMu.Lock()
	tex := textures[bg.TextureID]
	texMu.Unlock()
	if tex.ID == 0 {
		return
	}
	sw := float32(rl.GetScreenWidth())
	sh := float32(rl.GetScreenHeight())
	tint := rl.NewColor(
		uint8(bg.ColorR*255), uint8(bg.ColorG*255), uint8(bg.ColorB*255), uint8(bg.ColorA*255))
	// Apply parallax to offset
	ox := bg.OffsetX + cam.Target.X*bg.ParallaxX + bg.ScrollX
	oy := bg.OffsetY + cam.Target.Y*bg.ParallaxY + bg.ScrollY
	if bg.Tiled && bg.TileW > 0 && bg.TileH > 0 {
		tw := float32(bg.TileW)
		th := float32(bg.TileH)
		for y := -th; y < sh+th; y += th {
			for x := -tw; x < sw+tw; x += tw {
				rl.DrawTextureEx(tex, rl.Vector2{X: x - ox, Y: y - oy}, 0, 1, tint)
			}
		}
	} else {
		rl.DrawTextureEx(tex, rl.Vector2{X: -ox, Y: -oy}, 0, 1, tint)
	}
	for _, layer := range bg.Layers {
		texMu.Lock()
		lt := textures[layer.TextureID]
		texMu.Unlock()
		if lt.ID == 0 {
			continue
		}
		lox := cam.Target.X * layer.ParallaxX
		loy := cam.Target.Y * layer.ParallaxY
		rl.DrawTextureEx(lt, rl.Vector2{X: -lox, Y: -loy}, 0, 1, tint)
	}
}
