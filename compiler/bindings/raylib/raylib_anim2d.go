// Package raylib: 2D sprite (texture) animation: sprite-sheet frame animation with time-based playback.
package raylib

import (
	"fmt"
	"sync"

	"cyberbasic/compiler/vm"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type spriteAnimState struct {
	TextureId    string
	FrameWidth   int32
	FrameHeight  int32
	FramesPerRow int32
	TotalFrames  int32
	FPS          float64
	CurrentTime  float64
	Loop         bool
	CurrentFrame int32
}

var (
	spriteAnimStates   = make(map[string]*spriteAnimState)
	spriteAnimCounter  int
	spriteAnimMu       sync.Mutex
)

func registerAnim2D(v *vm.VM) {
	v.RegisterForeign("CreateSpriteAnimation", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("CreateSpriteAnimation requires (textureId, frameWidth, frameHeight, framesPerRow [, totalFrames])")
		}
		textureId := toString(args[0])
		frameWidth := toInt32(args[1])
		frameHeight := toInt32(args[2])
		framesPerRow := toInt32(args[3])
		if frameWidth <= 0 || frameHeight <= 0 || framesPerRow <= 0 {
			return nil, fmt.Errorf("CreateSpriteAnimation: frameWidth, frameHeight, framesPerRow must be positive")
		}
		texMu.Lock()
		tex, ok := textures[textureId]
		texMu.Unlock()
		if !ok {
			return nil, fmt.Errorf("unknown texture id: %s", textureId)
		}
		var totalFrames int32
		if len(args) >= 5 {
			totalFrames = toInt32(args[4])
			if totalFrames <= 0 {
				totalFrames = (tex.Width / frameWidth) * (tex.Height / frameHeight)
				if totalFrames <= 0 {
					totalFrames = 1
				}
			}
		} else {
			totalFrames = (tex.Width / frameWidth) * (tex.Height / frameHeight)
			if totalFrames <= 0 {
				totalFrames = 1
			}
		}
		spriteAnimMu.Lock()
		spriteAnimCounter++
		id := fmt.Sprintf("anim2d_%d", spriteAnimCounter)
		spriteAnimStates[id] = &spriteAnimState{
			TextureId:    textureId,
			FrameWidth:   frameWidth,
			FrameHeight:  frameHeight,
			FramesPerRow: framesPerRow,
			TotalFrames:  totalFrames,
			FPS:          8,
			Loop:         true,
		}
		spriteAnimMu.Unlock()
		return id, nil
	})

	v.RegisterForeign("SetSpriteAnimationFPS", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("SetSpriteAnimationFPS requires (animId, fps)")
		}
		id := toString(args[0])
		fps := toFloat64(args[1])
		if fps < 0 {
			fps = 0
		}
		spriteAnimMu.Lock()
		defer spriteAnimMu.Unlock()
		s, ok := spriteAnimStates[id]
		if !ok {
			return nil, fmt.Errorf("unknown sprite animation id: %s", id)
		}
		s.FPS = fps
		return nil, nil
	})

	v.RegisterForeign("SetSpriteAnimationLoop", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("SetSpriteAnimationLoop requires (animId, loop)")
		}
		id := toString(args[0])
		loop := false
		if v := args[1]; v != nil {
			switch x := v.(type) {
			case bool:
				loop = x
			case int:
				loop = x != 0
			case float64:
				loop = x != 0
			}
		}
		spriteAnimMu.Lock()
		defer spriteAnimMu.Unlock()
		s, ok := spriteAnimStates[id]
		if !ok {
			return nil, fmt.Errorf("unknown sprite animation id: %s", id)
		}
		s.Loop = loop
		return nil, nil
	})

	v.RegisterForeign("SetSpriteAnimationFrame", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("SetSpriteAnimationFrame requires (animId, frameIndex)")
		}
		id := toString(args[0])
		idx := toInt32(args[1])
		spriteAnimMu.Lock()
		defer spriteAnimMu.Unlock()
		s, ok := spriteAnimStates[id]
		if !ok {
			return nil, fmt.Errorf("unknown sprite animation id: %s", id)
		}
		if idx < 0 {
			idx = 0
		}
		if idx >= s.TotalFrames {
			idx = s.TotalFrames - 1
		}
		s.CurrentFrame = idx
		s.CurrentTime = float64(idx) / s.FPS
		return nil, nil
	})

	v.RegisterForeign("UpdateSpriteAnimation", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("UpdateSpriteAnimation requires (animId, deltaTime)")
		}
		id := toString(args[0])
		dt := toFloat64(args[1])
		spriteAnimMu.Lock()
		defer spriteAnimMu.Unlock()
		s, ok := spriteAnimStates[id]
		if !ok {
			return nil, fmt.Errorf("unknown sprite animation id: %s", id)
		}
		s.CurrentTime += dt
		frameIdx := int32(s.CurrentTime * s.FPS)
		if s.Loop && s.TotalFrames > 0 {
			for frameIdx >= s.TotalFrames {
				frameIdx -= s.TotalFrames
			}
			for frameIdx < 0 {
				frameIdx += s.TotalFrames
			}
		} else {
			if frameIdx >= s.TotalFrames {
				frameIdx = s.TotalFrames - 1
			}
			if frameIdx < 0 {
				frameIdx = 0
			}
		}
		s.CurrentFrame = frameIdx
		return nil, nil
	})

	v.RegisterForeign("GetSpriteAnimationFrame", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return 0, nil
		}
		id := toString(args[0])
		spriteAnimMu.Lock()
		defer spriteAnimMu.Unlock()
		s, ok := spriteAnimStates[id]
		if !ok {
			return 0, nil
		}
		return int(s.CurrentFrame), nil
	})

	v.RegisterForeign("DrawSpriteAnimation", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("DrawSpriteAnimation requires (animId, posX, posY [, scaleX, scaleY, rotation, tint r,g,b,a])")
		}
		id := toString(args[0])
		posX := toFloat32(args[1])
		posY := toFloat32(args[2])
		scaleX := float32(1)
		scaleY := float32(1)
		if len(args) >= 5 {
			scaleX = toFloat32(args[3])
			scaleY = toFloat32(args[4])
		}
		rotation := float32(0)
		tintOffset := 5
		if len(args) >= 6 {
			rotation = toFloat32(args[5])
			tintOffset = 6
		}
		spriteAnimMu.Lock()
		s, ok := spriteAnimStates[id]
		spriteAnimMu.Unlock()
		if !ok {
			return nil, fmt.Errorf("unknown sprite animation id: %s", id)
		}
		texMu.Lock()
		tex, ok := textures[s.TextureId]
		texMu.Unlock()
		if !ok {
			return nil, fmt.Errorf("texture %s not found for sprite animation %s", s.TextureId, id)
		}
		col := int32(s.CurrentFrame) % s.FramesPerRow
		row := int32(s.CurrentFrame) / s.FramesPerRow
		srcX := float32(col * s.FrameWidth)
		srcY := float32(row * s.FrameHeight)
		sourceRec := rl.Rectangle{X: srcX, Y: srcY, Width: float32(s.FrameWidth), Height: float32(s.FrameHeight)}
		destW := float32(s.FrameWidth) * scaleX
		destH := float32(s.FrameHeight) * scaleY
		destRec := rl.Rectangle{X: posX, Y: posY, Width: destW, Height: destH}
		origin := rl.Vector2{X: destW / 2, Y: destH / 2}
		c := rl.White
		if len(args) >= tintOffset+4 {
			c = argsToColor(args, tintOffset)
		}
		rl.DrawTexturePro(tex, sourceRec, destRec, origin, rotation, c)
		return nil, nil
	})

	v.RegisterForeign("DestroySpriteAnimation", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, nil
		}
		id := toString(args[0])
		spriteAnimMu.Lock()
		delete(spriteAnimStates, id)
		spriteAnimMu.Unlock()
		return nil, nil
	})
}
