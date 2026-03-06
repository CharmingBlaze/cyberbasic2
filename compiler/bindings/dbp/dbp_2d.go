// Package dbp - 2D Game API: drawing, sprites, spritesheets, tilemaps, camera, collision, physics, objects, particles.
package dbp

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"strings"
	"sync"

	"cyberbasic/compiler/bindings/aseprite"
	"cyberbasic/compiler/vm"
	rl "github.com/gen2brain/raylib-go/raylib"
)

// --- Spritesheet registry ---
type spritesheetEntry struct {
	tex          rl.Texture2D
	frameW       int
	frameH       int
	frameCount   int
	currentFrame int
	animStart    int
	animEnd      int
	animSpeed    float32
	animAccum    float32
	playing      bool
	// Aseprite mode: non-nil when loaded from JSON
	aseprite   *aseprite.Sheet
	currentTag string // tag name when playing
}

var (
	spritesheets       = make(map[int]*spritesheetEntry)
	spritesheetsMu     sync.Mutex
	spritesheetTexRefs = make(map[uint32]int) // texture ID -> ref count
	spritesheetTexMu   sync.Mutex

	// Tilemap id mapping: DBP int id -> game string id (e.g. "tm_1")
	tilemapId2Str    = make(map[int]string)
	tilemapId2StrMu  sync.Mutex
	tilemapVisible   = make(map[int]bool) // default true when not set
	tilemapVisibleMu sync.Mutex

	// Camera2D default state
	camera2DOn       bool
	camera2DTargetX  float32
	camera2DTargetY  float32
	camera2DZoom     float32 = 1.0
	camera2DRotation float32
	camera2DFollowID int = -1
	camera2DMu       sync.Mutex

	// SpriteObject2D registry
	spriteObjects2D   = make(map[int]*spriteObject2D)
	spriteObjects2DMu sync.Mutex

	// Particles2D registry
	particles2D   = make(map[int]*particles2DSystem)
	particles2DMu sync.Mutex
)

type spriteObject2D struct {
	spriteId int
	x        float32
	y        float32
	angle    float32
	sx       float32
	sy       float32
	syncMe   bool
	visible  bool
}

type particle2D struct {
	x, y   float32
	vx, vy float32
	life   float32
	maxLife float32
	r, g, b uint8
	size   float32
}

type particles2DSystem struct {
	particles []particle2D
	maxCount  int
	r, g, b   uint8
	size      float32
	speed     float32
}

func toFloat64_2d(v interface{}) float64 {
	switch x := v.(type) {
	case int:
		return float64(x)
	case float64:
		return x
	case string:
		f, _ := strconv.ParseFloat(x, 64)
		return f
	default:
		return 0
	}
}

// Register2D registers all 2D game API commands. Call after game.RegisterGame so SetTile/GetTile overwrite.
func Register2D(v *vm.VM) {
	register2DDrawing(v)
	register2DSprites(v)
	register2DSpritesheets(v)
	register2DTilemaps(v)
	register2DCamera(v)
	register2DCollision(v)
	register2DPhysics(v)
	register2DObjects(v)
	register2DUI(v)
	register2DMath(v)
	register2DParticles(v)
}

func register2DDrawing(v *vm.VM) {
	v.RegisterForeign("DrawPixel", func(args []interface{}) (interface{}, error) {
		if len(args) < 5 {
			return nil, fmt.Errorf("DrawPixel(x, y, r, g, b) requires 5 arguments")
		}
		x, y := int32(toInt(args[0])), int32(toInt(args[1]))
		r, g, b := toInt(args[2])&0xff, toInt(args[3])&0xff, toInt(args[4])&0xff
		rl.DrawPixel(x, y, rl.NewColor(uint8(r), uint8(g), uint8(b), 255))
		return nil, nil
	})
	v.RegisterForeign("DrawRectOutline", func(args []interface{}) (interface{}, error) {
		if len(args) < 7 {
			return nil, fmt.Errorf("DrawRectOutline(x, y, w, h, r, g, b) requires 7 arguments")
		}
		x, y, w, h := int32(toInt(args[0])), int32(toInt(args[1])), int32(toInt(args[2])), int32(toInt(args[3]))
		r, g, b := toInt(args[4])&0xff, toInt(args[5])&0xff, toInt(args[6])&0xff
		rl.DrawRectangleLines(x, y, w, h, rl.NewColor(uint8(r), uint8(g), uint8(b), 255))
		return nil, nil
	})
	v.RegisterForeign("DrawCircleOutline", func(args []interface{}) (interface{}, error) {
		if len(args) < 7 {
			return nil, fmt.Errorf("DrawCircleOutline(x, y, radius, r, g, b) requires 7 arguments")
		}
		x, y := toFloat32(args[0]), toFloat32(args[1])
		radius := toFloat32(args[2])
		r, g, b := toInt(args[3])&0xff, toInt(args[4])&0xff, toInt(args[5])&0xff
		rl.DrawCircleLines(int32(x), int32(y), radius, rl.NewColor(uint8(r), uint8(g), uint8(b), 255))
		return nil, nil
	})
	v.RegisterForeign("DrawTriangle", func(args []interface{}) (interface{}, error) {
		if len(args) < 10 {
			return nil, fmt.Errorf("DrawTriangle(x1,y1, x2,y2, x3,y3, r,g,b) requires 10 arguments")
		}
		v1 := rl.Vector2{X: toFloat32(args[0]), Y: toFloat32(args[1])}
		v2 := rl.Vector2{X: toFloat32(args[2]), Y: toFloat32(args[3])}
		v3 := rl.Vector2{X: toFloat32(args[4]), Y: toFloat32(args[5])}
		r, g, b := toInt(args[6])&0xff, toInt(args[7])&0xff, toInt(args[8])&0xff
		rl.DrawTriangle(v1, v2, v3, rl.NewColor(uint8(r), uint8(g), uint8(b), 255))
		return nil, nil
	})
}

func register2DSprites(v *vm.VM) {
	v.RegisterForeign("SetSpriteColor", func(args []interface{}) (interface{}, error) {
		if len(args) < 5 {
			return nil, fmt.Errorf("SetSpriteColor(id, r, g, b, a) requires 5 arguments")
		}
		id := toInt(args[0])
		r, g, b, a := toInt(args[1])&0xff, toInt(args[2])&0xff, toInt(args[3])&0xff, toInt(args[4])&0xff
		spriteColorsMu.Lock()
		spriteColors[id] = rl.NewColor(uint8(r), uint8(g), uint8(b), uint8(a))
		spriteColorsMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("DeleteSprite", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("DeleteSprite(id) requires 1 argument")
		}
		id := toInt(args[0])
		imagesMu.Lock()
		tex, ok := images[id]
		if ok {
			rl.UnloadTexture(tex)
			delete(images, id)
		}
		imagesMu.Unlock()
		spriteColorsMu.Lock()
		delete(spriteColors, id)
		spriteColorsMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("DrawSpriteRotated", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("DrawSpriteRotated(id, x, y, angle) requires 4 arguments")
		}
		id := toInt(args[0])
		x, y := toFloat32(args[1]), toFloat32(args[2])
		angle := toFloat32(args[3])
		imagesMu.Lock()
		tex, ok := images[id]
		imagesMu.Unlock()
		if !ok {
			return nil, fmt.Errorf("unknown sprite id %d", id)
		}
		pos := rl.Vector2{X: x, Y: y}
		rl.DrawTextureEx(tex, pos, angle, 1, getSpriteColor(id))
		return nil, nil
	})
	v.RegisterForeign("DrawSpriteScaled", func(args []interface{}) (interface{}, error) {
		if len(args) < 5 {
			return nil, fmt.Errorf("DrawSpriteScaled(id, x, y, sx, sy) requires 5 arguments")
		}
		id := toInt(args[0])
		x, y := toFloat32(args[1]), toFloat32(args[2])
		sx, sy := toFloat32(args[3]), toFloat32(args[4])
		imagesMu.Lock()
		tex, ok := images[id]
		imagesMu.Unlock()
		if !ok {
			return nil, fmt.Errorf("unknown sprite id %d", id)
		}
		src := rl.Rectangle{X: 0, Y: 0, Width: float32(tex.Width), Height: float32(tex.Height)}
		dest := rl.Rectangle{X: x, Y: y, Width: float32(tex.Width) * sx, Height: float32(tex.Height) * sy}
		rl.DrawTexturePro(tex, src, dest, rl.Vector2{}, 0, getSpriteColor(id))
		return nil, nil
	})
	v.RegisterForeign("DrawSpriteTint", func(args []interface{}) (interface{}, error) {
		if len(args) < 6 {
			return nil, fmt.Errorf("DrawSpriteTint(id, x, y, r, g, b) requires 6 arguments")
		}
		id := toInt(args[0])
		x, y := int32(toInt(args[1])), int32(toInt(args[2]))
		r, g, b := toInt(args[3])&0xff, toInt(args[4])&0xff, toInt(args[5])&0xff
		imagesMu.Lock()
		tex, ok := images[id]
		imagesMu.Unlock()
		if !ok {
			return nil, fmt.Errorf("unknown sprite id %d", id)
		}
		rl.DrawTexture(tex, x, y, rl.NewColor(uint8(r), uint8(g), uint8(b), 255))
		return nil, nil
	})
}

func register2DSpritesheets(v *vm.VM) {
	v.RegisterForeign("LoadSpritesheet", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("LoadSpritesheet(id, pngPath, jsonPath) or (id, path, frameW, frameH) requires 3-4 arguments")
		}
		id := toInt(args[0])
		path := toString(args[1])
		tex := rl.LoadTexture(path)
		if tex.ID == 0 {
			return nil, fmt.Errorf("LoadSpritesheet: failed to load texture %s", path)
		}
		// 3 args: (id, pngPath, jsonPath) - Aseprite mode
		if len(args) == 3 {
			jsonPath := toString(args[2])
			if strings.HasSuffix(strings.ToLower(jsonPath), ".json") {
				sheet, err := aseprite.Load(jsonPath)
				if err != nil {
					rl.UnloadTexture(tex)
					return nil, fmt.Errorf("LoadSpritesheet: %w", err)
				}
				frameCount := len(sheet.Frames)
				if frameCount == 0 {
					rl.UnloadTexture(tex)
					return nil, fmt.Errorf("LoadSpritesheet: no frames in %s", jsonPath)
				}
				spritesheetsMu.Lock()
				spritesheets[id] = &spritesheetEntry{
					tex:          tex,
					frameCount:   frameCount,
					aseprite:     sheet,
					frameW:       0, // per-frame in aseprite
					frameH:       0,
				}
				spritesheetsMu.Unlock()
				spritesheetTexMu.Lock()
				spritesheetTexRefs[tex.ID]++
				spritesheetTexMu.Unlock()
				return nil, nil
			}
		}
		// 4 args: (id, path, frameW, frameH) - grid mode
		if len(args) < 4 {
			rl.UnloadTexture(tex)
			return nil, fmt.Errorf("LoadSpritesheet(id, path, frameW, frameH) requires 4 arguments for grid mode")
		}
		fw, fh := toInt(args[2]), toInt(args[3])
		if fw <= 0 {
			fw = 32
		}
		if fh <= 0 {
			fh = 32
		}
		cols := int(tex.Width) / fw
		rows := int(tex.Height) / fh
		count := cols * rows
		spritesheetsMu.Lock()
		spritesheets[id] = &spritesheetEntry{tex: tex, frameW: fw, frameH: fh, frameCount: count}
		spritesheetsMu.Unlock()
		spritesheetTexMu.Lock()
		spritesheetTexRefs[tex.ID]++
		spritesheetTexMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("SetSpriteFrame", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("SetSpriteFrame(id, frame) requires 2 arguments")
		}
		id := toInt(args[0])
		frame := toInt(args[1])
		spritesheetsMu.Lock()
		if s, ok := spritesheets[id]; ok {
			if frame >= 0 && frame < s.frameCount {
				s.currentFrame = frame
			}
		}
		spritesheetsMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("NextSpriteFrame", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("NextSpriteFrame(id) requires 1 argument")
		}
		id := toInt(args[0])
		spritesheetsMu.Lock()
		if s, ok := spritesheets[id]; ok {
			s.currentFrame = (s.currentFrame + 1) % s.frameCount
		}
		spritesheetsMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("DrawSpriteFrame", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("DrawSpriteFrame(id, frame, x, y) requires 4 arguments")
		}
		id := toInt(args[0])
		frame := toInt(args[1])
		x, y := toFloat32(args[2]), toFloat32(args[3])
		spritesheetsMu.Lock()
		s, ok := spritesheets[id]
		spritesheetsMu.Unlock()
		if !ok {
			return nil, fmt.Errorf("unknown spritesheet id %d", id)
		}
		if frame < 0 || frame >= s.frameCount {
			return nil, nil
		}
		var src rl.Rectangle
		if s.aseprite != nil && frame < len(s.aseprite.Frames) {
			fr := s.aseprite.Frames[frame]
			src = rl.Rectangle{
				X: float32(fr.X), Y: float32(fr.Y),
				Width: float32(fr.W), Height: float32(fr.H),
			}
		} else {
			cols := int(s.tex.Width) / s.frameW
			if cols <= 0 {
				cols = 1
			}
			col := frame % cols
			row := frame / cols
			src = rl.Rectangle{
				X:      float32(col * s.frameW),
				Y:      float32(row * s.frameH),
				Width:  float32(s.frameW),
				Height: float32(s.frameH),
			}
		}
		pos := rl.Vector2{X: x, Y: y}
		rl.DrawTextureRec(s.tex, src, pos, rl.White)
		return nil, nil
	})
	v.RegisterForeign("AnimateSprite", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("AnimateSprite(id, startFrame, endFrame, speed) requires 4 arguments")
		}
		id := toInt(args[0])
		start, end := toInt(args[1]), toInt(args[2])
		speed := toFloat32(args[3])
		spritesheetsMu.Lock()
		if s, ok := spritesheets[id]; ok {
			s.animStart = start
			s.animEnd = end
			s.animSpeed = speed
			s.currentFrame = start
		}
		spritesheetsMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("DeleteSpritesheet", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("DeleteSpritesheet(id) requires 1 argument")
		}
		id := toInt(args[0])
		spritesheetsMu.Lock()
		s, ok := spritesheets[id]
		if ok {
			delete(spritesheets, id)
		}
		spritesheetsMu.Unlock()
		if ok && s.tex.ID != 0 {
			spritesheetTexMu.Lock()
			spritesheetTexRefs[s.tex.ID]--
			if spritesheetTexRefs[s.tex.ID] <= 0 {
				delete(spritesheetTexRefs, s.tex.ID)
				spritesheetTexMu.Unlock()
				rl.UnloadTexture(s.tex)
			} else {
				spritesheetTexMu.Unlock()
			}
		}
		return nil, nil
	})
	v.RegisterForeign("CloneSpritesheet", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("CloneSpritesheet(newID, sourceID) requires 2 arguments")
		}
		newID := toInt(args[0])
		srcID := toInt(args[1])
		spritesheetsMu.Lock()
		src, ok := spritesheets[srcID]
		if !ok {
			spritesheetsMu.Unlock()
			return nil, fmt.Errorf("unknown spritesheet id %d", srcID)
		}
		clone := &spritesheetEntry{
			tex:          src.tex,
			frameW:       src.frameW,
			frameH:       src.frameH,
			frameCount:   src.frameCount,
			currentFrame: src.currentFrame,
			animStart:    src.animStart,
			animEnd:      src.animEnd,
			animSpeed:    src.animSpeed,
			animAccum:    src.animAccum,
			aseprite:     src.aseprite,
		}
		spritesheets[newID] = clone
		spritesheetsMu.Unlock()
		spritesheetTexMu.Lock()
		spritesheetTexRefs[src.tex.ID]++
		spritesheetTexMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("SpritesheetExists", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("SpritesheetExists(id) requires 1 argument")
		}
		id := toInt(args[0])
		spritesheetsMu.Lock()
		_, ok := spritesheets[id]
		spritesheetsMu.Unlock()
		if ok {
			return 1, nil
		}
		return 0, nil
	})
	// PlaySpriteAnimation(id, tagName, speed): Play animation by tag (Aseprite).
	v.RegisterForeign("PlaySpriteAnimation", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("PlaySpriteAnimation(id, tagName, speed) requires 3 arguments")
		}
		id := toInt(args[0])
		tagName := toString(args[1])
		speed := toFloat32(args[2])
		spritesheetsMu.Lock()
		ss, ok := spritesheets[id]
		if !ok || ss.aseprite == nil {
			spritesheetsMu.Unlock()
			return nil, nil
		}
		from, to, ok := ss.aseprite.GetTagFrameRange(tagName, ss.frameCount)
		if !ok {
			spritesheetsMu.Unlock()
			return nil, nil
		}
		ss.animStart = from
		ss.animEnd = to
		ss.currentFrame = from
		ss.animSpeed = speed
		ss.playing = true
		ss.animAccum = 0
		ss.currentTag = tagName
		spritesheetsMu.Unlock()
		return nil, nil
	})
	// StopSpriteAnimation(id): Stop sprite animation.
	v.RegisterForeign("StopSpriteAnimation", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("StopSpriteAnimation(id) requires 1 argument")
		}
		id := toInt(args[0])
		spritesheetsMu.Lock()
		if ss, ok := spritesheets[id]; ok {
			ss.playing = false
		}
		spritesheetsMu.Unlock()
		return nil, nil
	})
	// GetSpriteFrame(id): Return current frame index.
	v.RegisterForeign("GetSpriteFrame", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return 0, nil
		}
		id := toInt(args[0])
		spritesheetsMu.Lock()
		ss, ok := spritesheets[id]
		spritesheetsMu.Unlock()
		if !ok {
			return 0, nil
		}
		return ss.currentFrame, nil
	})
	// GetSliceRect(id, sliceName): Return "x,y,w,h" for slice bounds at current frame. (BASIC: no byref; use string)
	v.RegisterForeign("GetSliceRect", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return "0,0,0,0", nil
		}
		id := toInt(args[0])
		sliceName := toString(args[1])
		spritesheetsMu.Lock()
		ss, ok := spritesheets[id]
		spritesheetsMu.Unlock()
		if !ok || ss.aseprite == nil {
			return "0,0,0,0", nil
		}
		x, y, w, h := ss.aseprite.GetSliceBounds(sliceName, ss.currentFrame)
		return fmt.Sprintf("%d,%d,%d,%d", x, y, w, h), nil
	})
	// GetAnimationLength(id, tagName): Return frame count for tag (Aseprite).
	v.RegisterForeign("GetAnimationLength", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return 0, nil
		}
		id := toInt(args[0])
		tagName := toString(args[1])
		spritesheetsMu.Lock()
		ss, ok := spritesheets[id]
		spritesheetsMu.Unlock()
		if !ok || ss.aseprite == nil {
			return 0, nil
		}
		from, to, ok := ss.aseprite.GetTagFrameRange(tagName, ss.frameCount)
		if !ok {
			return 0, nil
		}
		return to - from + 1, nil
	})
}

// UpdateSpriteAnimations advances time for all playing Aseprite sprite animations.
// Called from renderer pre-2D pass each frame.
func UpdateSpriteAnimations() {
	dt := float32(0.016)
	if rl.IsWindowReady() {
		dt = rl.GetFrameTime()
	}
	spritesheetsMu.Lock()
	defer spritesheetsMu.Unlock()
	for _, ss := range spritesheets {
		if !ss.playing || ss.aseprite == nil || len(ss.aseprite.Frames) == 0 {
			continue
		}
		frameIdx := ss.currentFrame
		if frameIdx < 0 || frameIdx >= len(ss.aseprite.Frames) {
			continue
		}
		durMs := ss.aseprite.Frames[frameIdx].DurationMs
		if durMs <= 0 {
			durMs = 100
		}
		ss.animAccum += dt * 1000 * ss.animSpeed
		for ss.animAccum >= float32(durMs) {
			ss.animAccum -= float32(durMs)
			from, to := ss.animStart, ss.animEnd
			dir := "forward"
			if t, ok := ss.aseprite.Tags[ss.currentTag]; ok {
				dir = t.Direction
			}
			if dir == "reverse" {
				ss.currentFrame--
				if ss.currentFrame < from {
					ss.currentFrame = to
				}
			} else {
				ss.currentFrame++
				if ss.currentFrame > to {
					ss.currentFrame = from
				}
			}
		}
	}
}

func register2DTilemaps(v *vm.VM) {
	v.RegisterForeign("LoadTilemap", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("LoadTilemap(id, path) requires 2 arguments")
		}
		id := toInt(args[0])
		path := toString(args[1])
		result, err := v.CallForeign("LoadTilemap", []interface{}{path})
		if err != nil {
			return nil, err
		}
		strId := toString(result)
		tilemapId2StrMu.Lock()
		tilemapId2Str[id] = strId
		tilemapId2StrMu.Unlock()
		tilemapVisibleMu.Lock()
		tilemapVisible[id] = true
		tilemapVisibleMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("DrawTilemap", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("DrawTilemap(id [, x, y]) requires 1+ arguments")
		}
		id := toInt(args[0])
		tilemapVisibleMu.Lock()
		vis := tilemapVisible[id]
		tilemapVisibleMu.Unlock()
		if !vis {
			return nil, nil
		}
		tilemapId2StrMu.Lock()
		strId, ok := tilemapId2Str[id]
		tilemapId2StrMu.Unlock()
		if !ok {
			strId = toString(args[0])
		}
		_, err := v.CallForeign("DrawTilemap", []interface{}{strId})
		return nil, err
	})
	v.RegisterForeign("SetTile", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("SetTile(id, x, y, tileIndex) requires 4 arguments")
		}
		id := toInt(args[0])
		tilemapId2StrMu.Lock()
		strId, ok := tilemapId2Str[id]
		tilemapId2StrMu.Unlock()
		if !ok {
			strId = toString(args[0])
		}
		return v.CallForeign("SetTileByMapId", []interface{}{strId, args[1], args[2], args[3]})
	})
	v.RegisterForeign("GetTile", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("GetTile(id, x, y) requires 3 arguments")
		}
		id := toInt(args[0])
		tilemapId2StrMu.Lock()
		strId, ok := tilemapId2Str[id]
		tilemapId2StrMu.Unlock()
		if !ok {
			strId = toString(args[0])
		}
		return v.CallForeign("GetTileByMapId", []interface{}{strId, args[1], args[2]})
	})
	v.RegisterForeign("DeleteTilemap", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("DeleteTilemap(id) requires 1 argument")
		}
		id := toInt(args[0])
		tilemapId2StrMu.Lock()
		delete(tilemapId2Str, id)
		tilemapId2StrMu.Unlock()
		tilemapVisibleMu.Lock()
		delete(tilemapVisible, id)
		tilemapVisibleMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("HideTilemap", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("HideTilemap(id) requires 1 argument")
		}
		id := toInt(args[0])
		tilemapVisibleMu.Lock()
		tilemapVisible[id] = false
		tilemapVisibleMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("ShowTilemap", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("ShowTilemap(id) requires 1 argument")
		}
		id := toInt(args[0])
		tilemapVisibleMu.Lock()
		tilemapVisible[id] = true
		tilemapVisibleMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("TilemapExists", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("TilemapExists(id) requires 1 argument")
		}
		id := toInt(args[0])
		tilemapId2StrMu.Lock()
		_, ok := tilemapId2Str[id]
		tilemapId2StrMu.Unlock()
		if ok {
			return 1, nil
		}
		return 0, nil
	})
}

func register2DCamera(v *vm.VM) {
	v.RegisterForeign("Camera2DOn", func(args []interface{}) (interface{}, error) {
		camera2DMu.Lock()
		camera2DOn = true
		tx, ty := camera2DTargetX, camera2DTargetY
		zoom := camera2DZoom
		rot := camera2DRotation
		fid := camera2DFollowID
		camera2DMu.Unlock()
		if fid >= 0 {
			spriteObjects2DMu.Lock()
			if o, ok := spriteObjects2D[fid]; ok {
				tx, ty = o.x, o.y
			}
			spriteObjects2DMu.Unlock()
		}
		return v.CallForeign("BeginMode2D", []interface{}{float32(0), float32(0), tx, ty, rot, zoom})
	})
	v.RegisterForeign("Camera2DOff", func(args []interface{}) (interface{}, error) {
		camera2DMu.Lock()
		camera2DOn = false
		camera2DMu.Unlock()
		return v.CallForeign("EndMode2D", nil)
	})
	v.RegisterForeign("Camera2DPosition", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("Camera2DPosition(x, y) requires 2 arguments")
		}
		camera2DMu.Lock()
		camera2DTargetX = toFloat32(args[0])
		camera2DTargetY = toFloat32(args[1])
		camera2DMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("Camera2DZoom", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("Camera2DZoom(value) requires 1 argument")
		}
		camera2DMu.Lock()
		camera2DZoom = toFloat32(args[0])
		camera2DMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("Camera2DRotation", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("Camera2DRotation(angle) requires 1 argument")
		}
		camera2DMu.Lock()
		camera2DRotation = toFloat32(args[0])
		camera2DMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("Camera2DFollow", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("Camera2DFollow(objectId) requires 1 argument")
		}
		camera2DMu.Lock()
		camera2DFollowID = toInt(args[0])
		if camera2DFollowID < 0 {
			camera2DFollowID = -1
		}
		camera2DMu.Unlock()
		return nil, nil
	})
}

func register2DCollision(v *vm.VM) {
	v.RegisterForeign("RectCollides", func(args []interface{}) (interface{}, error) {
		if len(args) < 8 {
			return nil, fmt.Errorf("RectCollides(x1,y1,w1,h1, x2,y2,w2,h2) requires 8 arguments")
		}
		return v.CallForeign("CheckCollisionRecs", args[:8])
	})
	v.RegisterForeign("PointInRect", func(args []interface{}) (interface{}, error) {
		if len(args) < 6 {
			return nil, fmt.Errorf("PointInRect(x,y, rx,ry,rw,rh) requires 6 arguments")
		}
		return v.CallForeign("CheckCollisionPointRec", args[:6])
	})
	v.RegisterForeign("CircleCollides", func(args []interface{}) (interface{}, error) {
		if len(args) < 6 {
			return nil, fmt.Errorf("CircleCollides(x1,y1,r1, x2,y2,r2) requires 6 arguments")
		}
		return v.CallForeign("CheckCollisionCircles", args[:6])
	})
	v.RegisterForeign("PointInCircle", func(args []interface{}) (interface{}, error) {
		if len(args) < 5 {
			return nil, fmt.Errorf("PointInCircle(x,y, cx,cy,r) requires 5 arguments")
		}
		return v.CallForeign("CheckCollisionPointCircle", args[:5])
	})
}

func register2DPhysics(v *vm.VM) {
	v.RegisterForeign("Physics2DOn", func(args []interface{}) (interface{}, error) {
		return v.CallForeign("PhysicsOn2D", args)
	})
	v.RegisterForeign("Physics2DOff", func(args []interface{}) (interface{}, error) {
		return v.CallForeign("PhysicsOff2D", args)
	})
	v.RegisterForeign("MakeBody2D", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("MakeBody2D(id, mass) requires 2 arguments")
		}
		id := toString(args[0])
		mass := toFloat64_2d(args[1])
		if mass <= 0 {
			mass = 1
		}
		return v.CallForeign("MakeRigidBody2D", []interface{}{id, 0, 0, 1, 1, mass})
	})
	v.RegisterForeign("MakeStatic2D", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("MakeStatic2D(id) requires 1 argument")
		}
		id := toString(args[0])
		return v.CallForeign("MakeStaticBody2D", []interface{}{id, 0, 0, 10, 1})
	})
	v.RegisterForeign("SetBody2DPosition", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("SetBody2DPosition(id, x, y) requires 3 arguments")
		}
		sid := toString(args[0])
		x, y := args[1], args[2]
		return v.CallForeign("SetPosition2D", []interface{}{"default", sid, x, y})
	})
	v.RegisterForeign("SetBody2DVelocity", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("SetBody2DVelocity(id, vx, vy) requires 3 arguments")
		}
		sid := toString(args[0])
		vx, vy := args[1], args[2]
		return v.CallForeign("SetVelocity2D", []interface{}{"default", sid, vx, vy})
	})
	v.RegisterForeign("ApplyForce2D", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("ApplyForce2D(id, fx, fy) requires 3 arguments")
		}
		sid := toString(args[0])
		fx, fy := args[1], args[2]
		return v.CallForeign("ApplyForce2D", []interface{}{"default", sid, fx, fy})
	})
	v.RegisterForeign("ApplyImpulse2D", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("ApplyImpulse2D(id, ix, iy) requires 3 arguments")
		}
		sid := toString(args[0])
		ix, iy := args[1], args[2]
		return v.CallForeign("ApplyImpulse2D", []interface{}{"default", sid, ix, iy})
	})
	v.RegisterForeign("GetBody2DX", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return 0.0, nil
		}
		return v.CallForeign("GetPositionX2D", []interface{}{"default", toString(args[0])})
	})
	v.RegisterForeign("GetBody2DY", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return 0.0, nil
		}
		return v.CallForeign("GetPositionY2D", []interface{}{"default", toString(args[0])})
	})
	v.RegisterForeign("GetBody2DVX", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return 0.0, nil
		}
		return v.CallForeign("GetVelocityX2D", []interface{}{"default", toString(args[0])})
	})
	v.RegisterForeign("GetBody2DVY", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return 0.0, nil
		}
		return v.CallForeign("GetVelocityY2D", []interface{}{"default", toString(args[0])})
	})
}

func register2DObjects(v *vm.VM) {
	v.RegisterForeign("MakeSpriteObject", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("MakeSpriteObject(id, spriteId) requires 2 arguments")
		}
		id := toInt(args[0])
		spriteId := toInt(args[1])
		spriteObjects2DMu.Lock()
		spriteObjects2D[id] = &spriteObject2D{spriteId: spriteId, x: 0, y: 0, angle: 0, sx: 1, sy: 1, visible: true}
		spriteObjects2DMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("PositionObject2D", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("PositionObject2D(id, x, y) requires 3 arguments")
		}
		id := toInt(args[0])
		x, y := toFloat32(args[1]), toFloat32(args[2])
		spriteObjects2DMu.Lock()
		if o, ok := spriteObjects2D[id]; ok {
			o.x, o.y = x, y
		}
		spriteObjects2DMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("MoveObject2D", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("MoveObject2D(id, dx, dy) requires 3 arguments")
		}
		id := toInt(args[0])
		dx, dy := toFloat32(args[1]), toFloat32(args[2])
		spriteObjects2DMu.Lock()
		if o, ok := spriteObjects2D[id]; ok {
			o.x += dx
			o.y += dy
		}
		spriteObjects2DMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("RotateObject2D", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("RotateObject2D(id, angle) requires 2 arguments")
		}
		id := toInt(args[0])
		angle := toFloat32(args[1])
		spriteObjects2DMu.Lock()
		if o, ok := spriteObjects2D[id]; ok {
			o.angle = angle
		}
		spriteObjects2DMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("ScaleObject2D", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("ScaleObject2D(id, sx, sy) requires 3 arguments")
		}
		id := toInt(args[0])
		sx, sy := toFloat32(args[1]), toFloat32(args[2])
		spriteObjects2DMu.Lock()
		if o, ok := spriteObjects2D[id]; ok {
			o.sx, o.sy = sx, sy
		}
		spriteObjects2DMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("DrawObject2D", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("DrawObject2D(id) requires 1 argument")
		}
		id := toInt(args[0])
		spriteObjects2DMu.Lock()
		o, ok := spriteObjects2D[id]
		spriteObjects2DMu.Unlock()
		if !ok || !o.visible {
			return nil, nil
		}
		imagesMu.Lock()
		tex, texOk := images[o.spriteId]
		imagesMu.Unlock()
		if !texOk {
			return nil, nil
		}
		src := rl.Rectangle{X: 0, Y: 0, Width: float32(tex.Width), Height: float32(tex.Height)}
		dest := rl.Rectangle{
			X: o.x, Y: o.y,
			Width:  float32(tex.Width) * o.sx,
			Height: float32(tex.Height) * o.sy,
		}
		origin := rl.Vector2{X: dest.Width / 2, Y: dest.Height / 2}
		rl.DrawTexturePro(tex, src, dest, origin, o.angle, rl.White)
		return nil, nil
	})
	v.RegisterForeign("SyncObject2D", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("SyncObject2D(id) requires 1 argument")
		}
		id := toInt(args[0])
		spriteObjects2DMu.Lock()
		if o, ok := spriteObjects2D[id]; ok {
			o.syncMe = true
		}
		spriteObjects2DMu.Unlock()
		return v.CallForeign("ReplicatePosition", []interface{}{fmt.Sprintf("obj2d_%d", id)})
	})
	v.RegisterForeign("DeleteSpriteObject", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("DeleteSpriteObject(id) requires 1 argument")
		}
		id := toInt(args[0])
		spriteObjects2DMu.Lock()
		delete(spriteObjects2D, id)
		spriteObjects2DMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("HideSpriteObject", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("HideSpriteObject(id) requires 1 argument")
		}
		id := toInt(args[0])
		spriteObjects2DMu.Lock()
		if o, ok := spriteObjects2D[id]; ok {
			o.visible = false
		}
		spriteObjects2DMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("ShowSpriteObject", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("ShowSpriteObject(id) requires 1 argument")
		}
		id := toInt(args[0])
		spriteObjects2DMu.Lock()
		if o, ok := spriteObjects2D[id]; ok {
			o.visible = true
		}
		spriteObjects2DMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("CloneSpriteObject", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("CloneSpriteObject(newID, sourceID) requires 2 arguments")
		}
		newID := toInt(args[0])
		srcID := toInt(args[1])
		spriteObjects2DMu.Lock()
		src, ok := spriteObjects2D[srcID]
		if !ok {
			spriteObjects2DMu.Unlock()
			return nil, fmt.Errorf("unknown sprite object id %d", srcID)
		}
		clone := &spriteObject2D{
			spriteId: src.spriteId,
			x:        src.x, y: src.y,
			angle: src.angle,
			sx: src.sx, sy: src.sy,
			syncMe:  false,
			visible: src.visible,
		}
		spriteObjects2D[newID] = clone
		spriteObjects2DMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("SpriteObjectExists", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("SpriteObjectExists(id) requires 1 argument")
		}
		id := toInt(args[0])
		spriteObjects2DMu.Lock()
		_, ok := spriteObjects2D[id]
		spriteObjects2DMu.Unlock()
		if ok {
			return 1, nil
		}
		return 0, nil
	})
}

func register2DUI(v *vm.VM) {
	v.RegisterForeign("UITextbox", func(args []interface{}) (interface{}, error) {
		if len(args) < 5 {
			return nil, fmt.Errorf("UITextbox(id, x, y, w, h) requires 5 arguments")
		}
		id := toString(args[0])
		x, y, w, h := args[1], args[2], args[3], args[4]
		return v.CallForeign("GuiTextBoxId", []interface{}{id, x, y, w, h, ""})
	})
}

func register2DMath(v *vm.VM) {
	v.RegisterForeign("AngleBetween2D", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("AngleBetween2D(x1,y1, x2,y2) requires 4 arguments")
		}
		x1, y1 := toFloat64_2d(args[0]), toFloat64_2d(args[1])
		x2, y2 := toFloat64_2d(args[2]), toFloat64_2d(args[3])
		return math.Atan2(y2-y1, x2-x1), nil
	})
	v.RegisterForeign("Distance2D", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("Distance2D(x1,y1, x2,y2) requires 4 arguments")
		}
		x1, y1 := toFloat64_2d(args[0]), toFloat64_2d(args[1])
		x2, y2 := toFloat64_2d(args[2]), toFloat64_2d(args[3])
		dx, dy := x2-x1, y2-y1
		return math.Sqrt(dx*dx + dy*dy), nil
	})
	v.RegisterForeign("Normalize2D", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("Normalize2D(x, y) requires 2 arguments")
		}
		x, y := toFloat64_2d(args[0]), toFloat64_2d(args[1])
		len := math.Sqrt(x*x + y*y)
		if len < 1e-10 {
			return []interface{}{0.0, 0.0}, nil
		}
		return []interface{}{x / len, y / len}, nil
	})
	v.RegisterForeign("Dot2D", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("Dot2D(x1,y1, x2,y2) requires 4 arguments")
		}
		x1, y1 := toFloat64_2d(args[0]), toFloat64_2d(args[1])
		x2, y2 := toFloat64_2d(args[2]), toFloat64_2d(args[3])
		return x1*x2 + y1*y2, nil
	})
}

func register2DParticles(v *vm.VM) {
	v.RegisterForeign("MakeParticles2D", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("MakeParticles2D(id, maxCount) requires 2 arguments")
		}
		id := toInt(args[0])
		maxCount := toInt(args[1])
		if maxCount <= 0 {
			maxCount = 100
		}
		particles2DMu.Lock()
		particles2D[id] = &particles2DSystem{
			particles: make([]particle2D, 0, maxCount),
			maxCount:  maxCount,
			r: 255, g: 255, b: 255,
			size:  2, speed: 50,
		}
		particles2DMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("SetParticles2DColor", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("SetParticles2DColor(id, r, g, b) requires 4 arguments")
		}
		id := toInt(args[0])
		r, g, b := toInt(args[1])&0xff, toInt(args[2])&0xff, toInt(args[3])&0xff
		particles2DMu.Lock()
		if p, ok := particles2D[id]; ok {
			p.r, p.g, p.b = uint8(r), uint8(g), uint8(b)
		}
		particles2DMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("SetParticles2DSize", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("SetParticles2DSize(id, size) requires 2 arguments")
		}
		id := toInt(args[0])
		size := toFloat32(args[1])
		particles2DMu.Lock()
		if p, ok := particles2D[id]; ok {
			p.size = size
		}
		particles2DMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("SetParticles2DSpeed", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("SetParticles2DSpeed(id, speed) requires 2 arguments")
		}
		id := toInt(args[0])
		speed := toFloat32(args[1])
		particles2DMu.Lock()
		if p, ok := particles2D[id]; ok {
			p.speed = speed
		}
		particles2DMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("EmitParticles2D", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("EmitParticles2D(id, count [, x, y]) requires 2+ arguments")
		}
		id := toInt(args[0])
		count := toInt(args[1])
		x, y := float32(0), float32(0)
		if len(args) >= 4 {
			x, y = toFloat32(args[2]), toFloat32(args[3])
		}
		particles2DMu.Lock()
		p, ok := particles2D[id]
		if !ok {
			particles2DMu.Unlock()
			return nil, nil
		}
		for i := 0; i < count && len(p.particles) < p.maxCount; i++ {
			angle := float32(rand.Float64() * 2 * math.Pi)
			vx := float32(math.Cos(float64(angle))) * p.speed * (0.5 + float32(rand.Float64())*0.5)
			vy := float32(math.Sin(float64(angle))) * p.speed * (0.5 + float32(rand.Float64())*0.5)
			life := 0.5 + float32(rand.Float64())*1.5
			p.particles = append(p.particles, particle2D{
				x: x, y: y, vx: vx, vy: vy,
				life: life, maxLife: life,
				r: p.r, g: p.g, b: p.b,
				size: p.size,
			})
		}
		particles2DMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("DrawParticles2D", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("DrawParticles2D(id) requires 1 argument")
		}
		id := toInt(args[0])
		dt := rl.GetFrameTime()
		particles2DMu.Lock()
		p, ok := particles2D[id]
		if !ok {
			particles2DMu.Unlock()
			return nil, nil
		}
		live := p.particles[:0]
		for _, part := range p.particles {
			part.x += part.vx * dt
			part.y += part.vy * dt
			part.life -= dt
			if part.life > 0 {
				live = append(live, part)
			}
		}
		p.particles = live
		toDraw := make([]particle2D, len(live))
		copy(toDraw, live)
		particles2DMu.Unlock()
		for _, part := range toDraw {
			alpha := uint8(255 * part.life / part.maxLife)
			c := rl.NewColor(part.r, part.g, part.b, alpha)
			rl.DrawCircle(int32(part.x), int32(part.y), part.size, c)
		}
		return nil, nil
	})
	v.RegisterForeign("DeleteParticles2D", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("DeleteParticles2D(id) requires 1 argument")
		}
		id := toInt(args[0])
		particles2DMu.Lock()
		delete(particles2D, id)
		particles2DMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("Particles2DExists", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("Particles2DExists(id) requires 1 argument")
		}
		id := toInt(args[0])
		particles2DMu.Lock()
		_, ok := particles2D[id]
		particles2DMu.Unlock()
		if ok {
			return 1, nil
		}
		return 0, nil
	})
}
