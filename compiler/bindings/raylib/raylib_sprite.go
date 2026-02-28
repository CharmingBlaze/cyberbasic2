// Package raylib: sprite state layer (position, scale, rotation, origin, flip) and draw with DrawTexturePro.
package raylib

import (
	"fmt"
	"sync"

	"cyberbasic/compiler/vm"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type spriteState struct {
	TextureId  string
	X          float32
	Y          float32
	Scale      float32
	ScaleX     float32
	ScaleY     float32
	Rotation   float32
	OriginX    float32
	OriginY    float32
	FlipX      bool
	FlipY      bool
	LayerID    string
	ZIndex     int
	FrameIndex int
	FrameW     int // 0 = use full texture
	FrameH     int
	Playing    bool
	FrameCount int
	AnimSpeed  float32
}

var (
	sprites      = make(map[string]*spriteState)
	spriteCounter int
	spriteMu     sync.Mutex
)

func registerSprite(v *vm.VM) {
	v.RegisterForeign("CreateSprite", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("CreateSprite requires (textureId)")
		}
		texId := toString(args[0])
		texMu.Lock()
		_, ok := textures[texId]
		texMu.Unlock()
		if !ok {
			return nil, fmt.Errorf("unknown texture id: %s", texId)
		}
		spriteMu.Lock()
		spriteCounter++
		id := fmt.Sprintf("sprite_%d", spriteCounter)
		sprites[id] = &spriteState{
			TextureId: texId,
			Scale:     1,
			ScaleX:    1,
			ScaleY:    1,
		}
		spriteMu.Unlock()
		return id, nil
	})
	v.RegisterForeign("SpriteSetPosition", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("SpriteSetPosition requires (spriteId, x, y)")
		}
		id := toString(args[0])
		spriteMu.Lock()
		defer spriteMu.Unlock()
		s, ok := sprites[id]
		if !ok {
			return nil, fmt.Errorf("unknown sprite id: %s", id)
		}
		s.X = toFloat32(args[1])
		s.Y = toFloat32(args[2])
		return nil, nil
	})
	v.RegisterForeign("SpriteSetScale", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("SpriteSetScale requires (spriteId, scale)")
		}
		id := toString(args[0])
		spriteMu.Lock()
		defer spriteMu.Unlock()
		s, ok := sprites[id]
		if !ok {
			return nil, fmt.Errorf("unknown sprite id: %s", id)
		}
		scale := toFloat32(args[1])
		s.Scale = scale
		s.ScaleX = scale
		s.ScaleY = scale
		return nil, nil
	})
	v.RegisterForeign("SpriteSetScaleXY", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("SpriteSetScaleXY requires (spriteId, sx, sy)")
		}
		id := toString(args[0])
		spriteMu.Lock()
		defer spriteMu.Unlock()
		s, ok := sprites[id]
		if !ok {
			return nil, fmt.Errorf("unknown sprite id: %s", id)
		}
		s.ScaleX = toFloat32(args[1])
		s.ScaleY = toFloat32(args[2])
		return nil, nil
	})
	v.RegisterForeign("SpriteSetRotation", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("SpriteSetRotation requires (spriteId, angleRad)")
		}
		id := toString(args[0])
		spriteMu.Lock()
		defer spriteMu.Unlock()
		s, ok := sprites[id]
		if !ok {
			return nil, fmt.Errorf("unknown sprite id: %s", id)
		}
		s.Rotation = toFloat32(args[1])
		return nil, nil
	})
	v.RegisterForeign("SpriteSetOrigin", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("SpriteSetOrigin requires (spriteId, ox, oy)")
		}
		id := toString(args[0])
		spriteMu.Lock()
		defer spriteMu.Unlock()
		s, ok := sprites[id]
		if !ok {
			return nil, fmt.Errorf("unknown sprite id: %s", id)
		}
		s.OriginX = toFloat32(args[1])
		s.OriginY = toFloat32(args[2])
		return nil, nil
	})
	v.RegisterForeign("SpriteSetFlip", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("SpriteSetFlip requires (spriteId, flipX, flipY)")
		}
		id := toString(args[0])
		spriteMu.Lock()
		defer spriteMu.Unlock()
		s, ok := sprites[id]
		if !ok {
			return nil, fmt.Errorf("unknown sprite id: %s", id)
		}
		s.FlipX = toFloat32(args[1]) != 0
		s.FlipY = toFloat32(args[2]) != 0
		return nil, nil
	})
	v.RegisterForeign("SpriteSetLayer", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("SpriteSetLayer requires (spriteId, layerId)")
		}
		id := toString(args[0])
		layerID := toString(args[1])
		spriteMu.Lock()
		defer spriteMu.Unlock()
		s, ok := sprites[id]
		if !ok {
			return nil, fmt.Errorf("unknown sprite id: %s", id)
		}
		s.LayerID = layerID
		return nil, nil
	})
	v.RegisterForeign("SpriteSetZIndex", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("SpriteSetZIndex requires (spriteId, z)")
		}
		id := toString(args[0])
		z := int(toFloat32(args[1]))
		spriteMu.Lock()
		defer spriteMu.Unlock()
		s, ok := sprites[id]
		if !ok {
			return nil, fmt.Errorf("unknown sprite id: %s", id)
		}
		s.ZIndex = z
		return nil, nil
	})
	v.RegisterForeign("SpriteSetFrame", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("SpriteSetFrame requires (spriteId, frameIndex)")
		}
		id := toString(args[0])
		frame := int(toFloat32(args[1]))
		spriteMu.Lock()
		defer spriteMu.Unlock()
		s, ok := sprites[id]
		if !ok {
			return nil, fmt.Errorf("unknown sprite id: %s", id)
		}
		s.FrameIndex = frame
		return nil, nil
	})
	v.RegisterForeign("SpritePlay", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("SpritePlay requires (spriteId)")
		}
		id := toString(args[0])
		spriteMu.Lock()
		defer spriteMu.Unlock()
		s, ok := sprites[id]
		if !ok {
			return nil, fmt.Errorf("unknown sprite id: %s", id)
		}
		s.Playing = true
		return nil, nil
	})
	v.RegisterForeign("SpritePause", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("SpritePause requires (spriteId)")
		}
		id := toString(args[0])
		spriteMu.Lock()
		defer spriteMu.Unlock()
		s, ok := sprites[id]
		if !ok {
			return nil, fmt.Errorf("unknown sprite id: %s", id)
		}
		s.Playing = false
		return nil, nil
	})
	v.RegisterForeign("SpriteStop", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("SpriteStop requires (spriteId)")
		}
		id := toString(args[0])
		spriteMu.Lock()
		defer spriteMu.Unlock()
		s, ok := sprites[id]
		if !ok {
			return nil, fmt.Errorf("unknown sprite id: %s", id)
		}
		s.Playing = false
		s.FrameIndex = 0
		return nil, nil
	})
	v.RegisterForeign("SpriteSetFrameSize", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("SpriteSetFrameSize requires (spriteId, frameWidth, frameHeight)")
		}
		id := toString(args[0])
		w := int(toFloat32(args[1]))
		h := int(toFloat32(args[2]))
		spriteMu.Lock()
		defer spriteMu.Unlock()
		s, ok := sprites[id]
		if !ok {
			return nil, fmt.Errorf("unknown sprite id: %s", id)
		}
		s.FrameW = w
		s.FrameH = h
		return nil, nil
	})
	v.RegisterForeign("SpriteSetFrameCount", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("SpriteSetFrameCount requires (spriteId, count)")
		}
		id := toString(args[0])
		n := int(toFloat32(args[1]))
		spriteMu.Lock()
		defer spriteMu.Unlock()
		s, ok := sprites[id]
		if !ok {
			return nil, fmt.Errorf("unknown sprite id: %s", id)
		}
		s.FrameCount = n
		return nil, nil
	})
	v.RegisterForeign("SpriteSetAnimSpeed", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("SpriteSetAnimSpeed requires (spriteId, framesPerSecond)")
		}
		id := toString(args[0])
		spd := toFloat32(args[1])
		spriteMu.Lock()
		defer spriteMu.Unlock()
		s, ok := sprites[id]
		if !ok {
			return nil, fmt.Errorf("unknown sprite id: %s", id)
		}
		s.AnimSpeed = spd
		return nil, nil
	})
	v.RegisterForeign("SpriteDraw", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("SpriteDraw requires (spriteId) and optional tint")
		}
		id := toString(args[0])
		spriteMu.Lock()
		s, ok := sprites[id]
		if ok && s.Playing && s.FrameCount > 0 && s.AnimSpeed > 0 {
			dt := rl.GetFrameTime()
			s.FrameIndex = int(float32(s.FrameIndex) + s.AnimSpeed*dt) % s.FrameCount
			if s.FrameIndex < 0 {
				s.FrameIndex += s.FrameCount
			}
		}
		spriteMu.Unlock()
		if !ok {
			return nil, fmt.Errorf("unknown sprite id: %s", id)
		}
		texMu.Lock()
		tex, ok := textures[s.TextureId]
		texMu.Unlock()
		if !ok {
			return nil, fmt.Errorf("sprite %s texture not loaded: %s", id, s.TextureId)
		}
		w, h := float32(tex.Width), float32(tex.Height)
		src := rl.Rectangle{X: 0, Y: 0, Width: w, Height: h}
		if s.FrameW > 0 && s.FrameH > 0 {
			idx := s.FrameIndex
			cols := int(w) / s.FrameW
			if cols <= 0 {
				cols = 1
			}
			col := idx % cols
			row := idx / cols
			src = rl.Rectangle{
				X:      float32(col * s.FrameW),
				Y:      float32(row * s.FrameH),
				Width:  float32(s.FrameW),
				Height: float32(s.FrameH),
			}
			w = float32(s.FrameW)
			h = float32(s.FrameH)
		}
		dw := w * s.ScaleX
		dh := h * s.ScaleY
		if s.FlipX {
			dw = -dw
		}
		if s.FlipY {
			dh = -dh
		}
		dest := rl.Rectangle{X: s.X - s.OriginX*s.ScaleX, Y: s.Y - s.OriginY*s.ScaleY, Width: dw, Height: dh}
		origin := rl.Vector2{X: s.OriginX * s.ScaleX, Y: s.OriginY * s.ScaleY}
		c := rl.White
		if len(args) >= 5 {
			c = argsToColor(args, 1)
		}
		rl.DrawTexturePro(tex, src, dest, origin, s.Rotation, c)
		return nil, nil
	})
	v.RegisterForeign("DestroySprite", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, nil
		}
		id := toString(args[0])
		spriteMu.Lock()
		delete(sprites, id)
		spriteMu.Unlock()
		return nil, nil
	})
}

// GetSpriteLayerAndZ returns layerID and zIndex for a sprite (for 2D flush sorting).
func GetSpriteLayerAndZ(spriteID string) (layerID string, zIndex int) {
	spriteMu.Lock()
	defer spriteMu.Unlock()
	s, ok := sprites[spriteID]
	if !ok {
		return "", 0
	}
	return s.LayerID, s.ZIndex
}

func clearSpriteLayer(layerID string) {
	spriteMu.Lock()
	defer spriteMu.Unlock()
	for _, s := range sprites {
		if s.LayerID == layerID {
			s.LayerID = ""
		}
	}
}
