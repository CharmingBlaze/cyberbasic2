// Package raylib: textures (rtextures) load, draw, unload.
package raylib

import (
	"cyberbasic/compiler/vm"
	"fmt"
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func registerTextures(v *vm.VM) {
	v.RegisterForeign("LoadTexture", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("LoadTexture requires (fileName)")
		}
		path := toString(args[0])
		tex := rl.LoadTexture(path)
		texMu.Lock()
		texCounter++
		id := fmt.Sprintf("tex_%d", texCounter)
		textures[id] = tex
		texMu.Unlock()
		return id, nil
	})
	v.RegisterForeign("UnloadTexture", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("UnloadTexture requires (id)")
		}
		id := toString(args[0])
		texMu.Lock()
		tex, ok := textures[id]
		delete(textures, id)
		texMu.Unlock()
		if ok {
			rl.UnloadTexture(tex)
		}
		return nil, nil
	})
	v.RegisterForeign("LoadSprite", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("LoadSprite requires (path)")
		}
		path := toString(args[0])
		tex := rl.LoadTexture(path)
		texMu.Lock()
		texCounter++
		id := fmt.Sprintf("tex_%d", texCounter)
		textures[id] = tex
		texMu.Unlock()
		return id, nil
	})
	v.RegisterForeign("DrawSprite", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("DrawSprite requires (spriteId, x, y)")
		}
		id := toString(args[0])
		texMu.Lock()
		tex, ok := textures[id]
		texMu.Unlock()
		if !ok {
			return nil, fmt.Errorf("unknown texture id: %s", id)
		}
		x, y := toInt32(args[1]), toInt32(args[2])
		rl.DrawTexture(tex, x, y, rl.White)
		return nil, nil
	})
	v.RegisterForeign("LoadRenderTexture", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("LoadRenderTexture requires (width, height)")
		}
		rt := rl.LoadRenderTexture(toInt32(args[0]), toInt32(args[1]))
		renderTexMu.Lock()
		renderTexCounter++
		id := fmt.Sprintf("rt_%d", renderTexCounter)
		renderTextures[id] = rt
		renderTexMu.Unlock()
		return id, nil
	})
	v.RegisterForeign("UnloadRenderTexture", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("UnloadRenderTexture requires (id)")
		}
		id := toString(args[0])
		renderTexMu.Lock()
		rt, ok := renderTextures[id]
		delete(renderTextures, id)
		renderTexMu.Unlock()
		if ok {
			rl.UnloadRenderTexture(rt)
		}
		return nil, nil
	})
	v.RegisterForeign("BeginTextureMode", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("BeginTextureMode requires (renderTextureId)")
		}
		id := toString(args[0])
		renderTexMu.Lock()
		rt, ok := renderTextures[id]
		renderTexMu.Unlock()
		if !ok {
			return nil, fmt.Errorf("unknown render texture id: %s", id)
		}
		rl.BeginTextureMode(rt)
		return nil, nil
	})
	v.RegisterForeign("EndTextureMode", func(args []interface{}) (interface{}, error) {
		rl.EndTextureMode()
		return nil, nil
	})
	v.RegisterForeign("DrawTexture", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("DrawTexture requires (id, posX, posY) and optional tint r,g,b,a")
		}
		id := toString(args[0])
		texMu.Lock()
		tex, ok := textures[id]
		texMu.Unlock()
		if !ok {
			return nil, fmt.Errorf("unknown texture id: %s", id)
		}
		posX, posY := toInt32(args[1]), toInt32(args[2])
		c := rl.White
		if len(args) >= 7 {
			c = argsToColor(args, 3)
		}
		rl.DrawTexture(tex, posX, posY, c)
		return nil, nil
	})
	v.RegisterForeign("DrawTextureEx", func(args []interface{}) (interface{}, error) {
		if len(args) < 5 {
			return nil, fmt.Errorf("DrawTextureEx requires (id, posX, posY, rotation, scale) and optional tint")
		}
		id := toString(args[0])
		texMu.Lock()
		tex, ok := textures[id]
		texMu.Unlock()
		if !ok {
			return nil, fmt.Errorf("unknown texture id: %s", id)
		}
		pos := rl.Vector2{X: toFloat32(args[1]), Y: toFloat32(args[2])}
		rotation := toFloat32(args[3])
		scale := toFloat32(args[4])
		c := rl.White
		if len(args) >= 9 {
			c = argsToColor(args, 5)
		}
		rl.DrawTextureEx(tex, pos, rotation, scale, c)
		return nil, nil
	})
	v.RegisterForeign("DrawTextureRec", func(args []interface{}) (interface{}, error) {
		if len(args) < 8 {
			return nil, fmt.Errorf("DrawTextureRec requires (id, srcX, srcY, srcW, srcH, posX, posY, tint)")
		}
		id := toString(args[0])
		texMu.Lock()
		tex, ok := textures[id]
		texMu.Unlock()
		if !ok {
			return nil, fmt.Errorf("unknown texture id: %s", id)
		}
		source := rl.Rectangle{X: toFloat32(args[1]), Y: toFloat32(args[2]), Width: toFloat32(args[3]), Height: toFloat32(args[4])}
		pos := rl.Vector2{X: toFloat32(args[5]), Y: toFloat32(args[6])}
		c := rl.White
		if len(args) >= 12 {
			c = argsToColor(args, 7)
		}
		rl.DrawTextureRec(tex, source, pos, c)
		return nil, nil
	})
	v.RegisterForeign("DrawTexturePro", func(args []interface{}) (interface{}, error) {
		if len(args) < 11 {
			return nil, fmt.Errorf("DrawTexturePro requires (id, srcX,srcY,srcW,srcH, destX,destY,destW,destH, originX,originY, rotation, tint)")
		}
		id := toString(args[0])
		texMu.Lock()
		tex, ok := textures[id]
		texMu.Unlock()
		if !ok {
			return nil, fmt.Errorf("unknown texture id: %s", id)
		}
		sourceRec := rl.Rectangle{X: toFloat32(args[1]), Y: toFloat32(args[2]), Width: toFloat32(args[3]), Height: toFloat32(args[4])}
		destRec := rl.Rectangle{X: toFloat32(args[5]), Y: toFloat32(args[6]), Width: toFloat32(args[7]), Height: toFloat32(args[8])}
		origin := rl.Vector2{X: toFloat32(args[9]), Y: toFloat32(args[10])}
		rotation := toFloat32(args[11])
		c := rl.White
		if len(args) >= 16 {
			c = argsToColor(args, 12)
		}
		rl.DrawTexturePro(tex, sourceRec, destRec, origin, rotation, c)
		return nil, nil
	})
	v.RegisterForeign("DrawEntity", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("DrawEntity requires (entityName$)")
		}
		entityName := toString(args[0])
		g := v.Globals()[strings.ToLower(entityName)]
		if g == nil {
			return nil, fmt.Errorf("entity not found: %s", entityName)
		}
		m, ok := g.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("entity %s is not a map", entityName)
		}
		spriteVal, _ := m["sprite"]
		if spriteVal == nil {
			spriteVal, _ = m["texture"]
		}
		spriteId := toString(spriteVal)
		if spriteId == "" {
			return nil, fmt.Errorf("entity %s has no sprite or texture", entityName)
		}
		texMu.Lock()
		tex, ok := textures[spriteId]
		texMu.Unlock()
		if !ok {
			return nil, fmt.Errorf("entity %s sprite/texture id not loaded: %s", entityName, spriteId)
		}
		x := toFloat32(m["x"])
		y := toFloat32(m["y"])
		if scaleVal, has := m["scale"]; has && scaleVal != nil {
			scale := toFloat32(scaleVal)
			angle := toFloat32(m["angle"])
			rl.DrawTextureEx(tex, rl.Vector2{X: x, Y: y}, angle, scale, rl.White)
		} else {
			rl.DrawTexture(tex, int32(x), int32(y), rl.White)
		}
		return nil, nil
	})

	v.RegisterForeign("LoadTextureFromImage", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("LoadTextureFromImage requires (imageId)")
		}
		imageMu.Lock()
		img, ok := images[toString(args[0])]
		imageMu.Unlock()
		if !ok || img == nil {
			return nil, fmt.Errorf("unknown image id: %s", toString(args[0]))
		}
		tex := rl.LoadTextureFromImage(img)
		texMu.Lock()
		texCounter++
		id := fmt.Sprintf("tex_%d", texCounter)
		textures[id] = tex
		texMu.Unlock()
		return id, nil
	})
	v.RegisterForeign("LoadTextureCubemap", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("LoadTextureCubemap requires (imageId, layout)")
		}
		imageMu.Lock()
		img, ok := images[toString(args[0])]
		imageMu.Unlock()
		if !ok || img == nil {
			return nil, fmt.Errorf("unknown image id: %s", toString(args[0]))
		}
		tex := rl.LoadTextureCubemap(img, toInt32(args[1]))
		texMu.Lock()
		texCounter++
		id := fmt.Sprintf("tex_%d", texCounter)
		textures[id] = tex
		texMu.Unlock()
		return id, nil
	})
	v.RegisterForeign("IsTextureValid", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return false, nil
		}
		texMu.Lock()
		tex, ok := textures[toString(args[0])]
		texMu.Unlock()
		return ok && rl.IsTextureValid(tex), nil
	})
	v.RegisterForeign("IsRenderTextureValid", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return false, nil
		}
		renderTexMu.Lock()
		rt, ok := renderTextures[toString(args[0])]
		renderTexMu.Unlock()
		return ok && rl.IsRenderTextureValid(rt), nil
	})
	v.RegisterForeign("UpdateTexture", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("UpdateTexture requires (textureId, pixels)")
		}
		texMu.Lock()
		tex, ok := textures[toString(args[0])]
		texMu.Unlock()
		if !ok {
			return nil, fmt.Errorf("unknown texture id: %s", toString(args[0]))
		}
		var data []byte
		switch v := args[1].(type) {
		case string:
			data = []byte(v)
		case []byte:
			data = v
		default:
			return nil, fmt.Errorf("pixels must be string or []byte (RGBA)")
		}
		n := len(data) / 4
		pixels := make([]rl.Color, n)
		for i := 0; i < n; i++ {
			pixels[i] = rl.NewColor(data[i*4], data[i*4+1], data[i*4+2], data[i*4+3])
		}
		rl.UpdateTexture(tex, pixels)
		return nil, nil
	})
	v.RegisterForeign("UpdateTextureRec", func(args []interface{}) (interface{}, error) {
		if len(args) < 6 {
			return nil, fmt.Errorf("UpdateTextureRec requires (textureId, x, y, w, h, pixels)")
		}
		texMu.Lock()
		tex, ok := textures[toString(args[0])]
		texMu.Unlock()
		if !ok {
			return nil, fmt.Errorf("unknown texture id: %s", toString(args[0]))
		}
		rec := rl.Rectangle{X: toFloat32(args[1]), Y: toFloat32(args[2]), Width: toFloat32(args[3]), Height: toFloat32(args[4])}
		var data []byte
		switch v := args[5].(type) {
		case string:
			data = []byte(v)
		case []byte:
			data = v
		default:
			return nil, fmt.Errorf("pixels must be string or []byte (RGBA)")
		}
		n := len(data) / 4
		pixels := make([]rl.Color, n)
		for i := 0; i < n; i++ {
			pixels[i] = rl.NewColor(data[i*4], data[i*4+1], data[i*4+2], data[i*4+3])
		}
		rl.UpdateTextureRec(tex, rec, pixels)
		return nil, nil
	})
	v.RegisterForeign("GenTextureMipmaps", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("GenTextureMipmaps requires (textureId)")
		}
		texMu.Lock()
		tex, ok := textures[toString(args[0])]
		if !ok {
			texMu.Unlock()
			return nil, fmt.Errorf("unknown texture id: %s", toString(args[0]))
		}
		rl.GenTextureMipmaps(&tex)
		textures[toString(args[0])] = tex
		texMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("SetTextureFilter", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("SetTextureFilter requires (textureId, filter)")
		}
		texMu.Lock()
		tex, ok := textures[toString(args[0])]
		texMu.Unlock()
		if !ok {
			return nil, fmt.Errorf("unknown texture id: %s", toString(args[0]))
		}
		rl.SetTextureFilter(tex, rl.TextureFilterMode(toInt32(args[1])))
		return nil, nil
	})
	v.RegisterForeign("SetTextureWrap", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("SetTextureWrap requires (textureId, wrap)")
		}
		texMu.Lock()
		tex, ok := textures[toString(args[0])]
		texMu.Unlock()
		if !ok {
			return nil, fmt.Errorf("unknown texture id: %s", toString(args[0]))
		}
		rl.SetTextureWrap(tex, rl.TextureWrapMode(toInt32(args[1])))
		return nil, nil
	})
	v.RegisterForeign("DrawTextureV", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("DrawTextureV requires (textureId, posX, posY) and optional tint")
		}
		texMu.Lock()
		tex, ok := textures[toString(args[0])]
		texMu.Unlock()
		if !ok {
			return nil, fmt.Errorf("unknown texture id: %s", toString(args[0]))
		}
		pos := rl.Vector2{X: toFloat32(args[1]), Y: toFloat32(args[2])}
		c := rl.White
		if len(args) >= 7 {
			c = argsToColor(args, 3)
		}
		rl.DrawTextureV(tex, pos, c)
		return nil, nil
	})
	v.RegisterForeign("DrawTextureNPatch", func(args []interface{}) (interface{}, error) {
		if len(args) < 18 {
			return nil, fmt.Errorf("DrawTextureNPatch requires (texId, srcX,srcY,srcW,srcH, left,top,right,bottom, layout, destX,destY,destW,destH, originX,originY, rotation, tint)")
		}
		texMu.Lock()
		tex, ok := textures[toString(args[0])]
		texMu.Unlock()
		if !ok {
			return nil, fmt.Errorf("unknown texture id: %s", toString(args[0]))
		}
		nPatch := rl.NPatchInfo{
			Source: rl.Rectangle{X: toFloat32(args[1]), Y: toFloat32(args[2]), Width: toFloat32(args[3]), Height: toFloat32(args[4])},
			Left:   toInt32(args[5]), Top: toInt32(args[6]), Right: toInt32(args[7]), Bottom: toInt32(args[8]),
			Layout: rl.NPatchLayout(toInt32(args[9])),
		}
		dest := rl.Rectangle{X: toFloat32(args[10]), Y: toFloat32(args[11]), Width: toFloat32(args[12]), Height: toFloat32(args[13])}
		origin := rl.Vector2{X: toFloat32(args[14]), Y: toFloat32(args[15])}
		rotation := toFloat32(args[16])
		c := rl.White
		if len(args) >= 21 {
			c = argsToColor(args, 17)
		}
		rl.DrawTextureNPatch(tex, nPatch, dest, origin, rotation, c)
		return nil, nil
	})
}
