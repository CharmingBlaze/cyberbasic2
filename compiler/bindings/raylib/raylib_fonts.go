// Package raylib: custom fonts load, draw, measure (rtext).
package raylib

import (
	"fmt"
	"os"
	"strings"

	"cyberbasic/compiler/vm"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const fontDefaultID = "font_default"

func registerFonts(v *vm.VM) {
	v.RegisterForeign("GetFontDefault", func(args []interface{}) (interface{}, error) {
		fontMu.Lock()
		if _, ok := fonts[fontDefaultID]; !ok {
			fonts[fontDefaultID] = rl.GetFontDefault()
		}
		fontMu.Unlock()
		return fontDefaultID, nil
	})
	v.RegisterForeign("LoadFont", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("LoadFont requires (fileName)")
		}
		path := toString(args[0])
		font := rl.LoadFont(path)
		fontMu.Lock()
		fontCounter++
		id := fmt.Sprintf("font_%d", fontCounter)
		fonts[id] = font
		fontMu.Unlock()
		return id, nil
	})
	v.RegisterForeign("LoadFontEx", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("LoadFontEx requires (fileName, fontSize, fontChars)")
		}
		path := toString(args[0])
		fontSize := toInt32(args[1])
		var font rl.Font
		font = rl.LoadFontEx(path, fontSize, nil)
		fontMu.Lock()
		fontCounter++
		id := fmt.Sprintf("font_%d", fontCounter)
		fonts[id] = font
		fontMu.Unlock()
		return id, nil
	})
	v.RegisterForeign("DrawTextExFont", func(args []interface{}) (interface{}, error) {
		if len(args) < 6 {
			return nil, fmt.Errorf("DrawTextExFont requires (fontId, text, x, y, fontSize, spacing) and optional tint")
		}
		id := toString(args[0])
		fontMu.Lock()
		font, ok := fonts[id]
		fontMu.Unlock()
		if !ok {
			return nil, fmt.Errorf("unknown font id: %s", id)
		}
		text := toString(args[1])
		pos := rl.Vector2{X: toFloat32(args[2]), Y: toFloat32(args[3])}
		fontSize := toFloat32(args[4])
		spacing := toFloat32(args[5])
		c := rl.White
		if len(args) >= 10 {
			c = argsToColor(args, 6)
		}
		rl.DrawTextEx(font, text, pos, fontSize, spacing, c)
		return nil, nil
	})
	v.RegisterForeign("MeasureTextEx", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("MeasureTextEx requires (fontId, text, fontSize, spacing)")
		}
		id := toString(args[0])
		fontMu.Lock()
		font, ok := fonts[id]
		fontMu.Unlock()
		if !ok {
			return nil, fmt.Errorf("unknown font id: %s", id)
		}
		text := toString(args[1])
		fontSize := toFloat32(args[2])
		spacing := toFloat32(args[3])
		vec := rl.MeasureTextEx(font, text, fontSize, spacing)
		return []interface{}{float64(vec.X), float64(vec.Y)}, nil
	})
	v.RegisterForeign("UnloadFont", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("UnloadFont requires (id)")
		}
		id := toString(args[0])
		if id == fontDefaultID {
			return nil, nil
		}
		fontMu.Lock()
		font, ok := fonts[id]
		delete(fonts, id)
		fontMu.Unlock()
		if ok {
			rl.UnloadFont(font)
		}
		return nil, nil
	})
	v.RegisterForeign("LoadFontFromImage", func(args []interface{}) (interface{}, error) {
		if len(args) < 6 {
			return nil, fmt.Errorf("LoadFontFromImage requires (imageId, keyR, keyG, keyB, keyA, firstChar)")
		}
		imageMu.Lock()
		img, ok := images[toString(args[0])]
		imageMu.Unlock()
		if !ok || img == nil {
			return nil, fmt.Errorf("unknown image id: %s", toString(args[0]))
		}
		key := argsToColor(args, 1)
		firstChar := toInt32(args[5])
		font := rl.LoadFontFromImage(*img, key, firstChar)
		fontMu.Lock()
		fontCounter++
		id := fmt.Sprintf("font_%d", fontCounter)
		fonts[id] = font
		fontMu.Unlock()
		return id, nil
	})
	v.RegisterForeign("LoadFontFromMemory", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("LoadFontFromMemory requires (fileType, data, dataSize, fontSize)")
		}
		var data []byte
		switch d := args[1].(type) {
		case string:
			data = []byte(d)
		case []byte:
			data = d
		default:
			return nil, fmt.Errorf("data must be string or []byte")
		}
		dataSize := int(toInt32(args[2]))
		if dataSize < len(data) {
			data = data[:dataSize]
		}
		fontSize := toInt32(args[3])
		font := rl.LoadFontFromMemory(toString(args[0]), data, fontSize, nil)
		fontMu.Lock()
		fontCounter++
		id := fmt.Sprintf("font_%d", fontCounter)
		fonts[id] = font
		fontMu.Unlock()
		return id, nil
	})
	v.RegisterForeign("IsFontValid", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return false, nil
		}
		fontMu.Lock()
		font, ok := fonts[toString(args[0])]
		fontMu.Unlock()
		return ok && rl.IsFontValid(font), nil
	})
	v.RegisterForeign("LoadFontData", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("LoadFontData requires (data, dataSize, fontSize, type)")
		}
		var data []byte
		switch d := args[0].(type) {
		case string:
			data = []byte(d)
		case []byte:
			data = d
		default:
			return nil, fmt.Errorf("data must be string or []byte")
		}
		dataSize := int(toInt32(args[1]))
		if dataSize < len(data) {
			data = data[:dataSize]
		}
		fontSize := toInt32(args[2])
		typ := toInt32(args[3])
		glyphs := rl.LoadFontData(data, fontSize, nil, 0, typ)
		lastFontDataMu.Lock()
		lastFontData = glyphs
		lastFontDataMu.Unlock()
		return len(glyphs), nil
	})
	v.RegisterForeign("GenImageFontAtlas", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("GenImageFontAtlas requires (fontSize, padding, packMethod)")
		}
		lastFontDataMu.Lock()
		glyphs := lastFontData
		lastFontDataMu.Unlock()
		if len(glyphs) == 0 {
			return "", nil
		}
		glyphRecs := make([]*rl.Rectangle, len(glyphs))
		for i := range glyphRecs {
			glyphRecs[i] = &rl.Rectangle{}
		}
		img := rl.GenImageFontAtlas(glyphs, glyphRecs, toInt32(args[0]), toInt32(args[1]), toInt32(args[2]))
		p := new(rl.Image)
		*p = img
		imageMu.Lock()
		imageCounter++
		id := fmt.Sprintf("img_%d", imageCounter)
		images[id] = p
		imageMu.Unlock()
		return id, nil
	})
	v.RegisterForeign("UnloadFontData", func(args []interface{}) (interface{}, error) {
		lastFontDataMu.Lock()
		if len(lastFontData) > 0 {
			rl.UnloadFontData(lastFontData)
			lastFontData = nil
		}
		lastFontDataMu.Unlock()
		return nil, nil
	})
	// ExportFontAsCode: export font atlas texture as C header plus font metadata (raylib-go has no native API; we generate .h from font texture image).
	v.RegisterForeign("ExportFontAsCode", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("ExportFontAsCode requires (fontId, fileName)")
		}
		id := toString(args[0])
		fontMu.Lock()
		font, ok := fonts[id]
		fontMu.Unlock()
		if !ok {
			return false, nil
		}
		img := rl.LoadImageFromTexture(font.Texture)
		if img == nil {
			return false, nil
		}
		defer rl.UnloadImage(img)
		cols := rl.LoadImageColors(img)
		if cols == nil {
			return false, nil
		}
		defer rl.UnloadImageColors(cols)
		w, h := img.Width, img.Height
		var b strings.Builder
		b.WriteString("// Exported by CyberBasic ExportFontAsCode (font atlas RGBA + metadata)\n")
		b.WriteString("#ifndef FONT_EXPORT_H\n#define FONT_EXPORT_H\n\n")
		fmt.Fprintf(&b, "static const int FONT_BASE_SIZE = %d;\n", font.BaseSize)
		fmt.Fprintf(&b, "static const int FONT_CHARS_COUNT = %d;\n", font.CharsCount)
		fmt.Fprintf(&b, "static const int FONT_CHARS_PADDING = %d;\n", font.CharsPadding)
		fmt.Fprintf(&b, "static const int FONT_ATLAS_WIDTH = %d;\n", w)
		fmt.Fprintf(&b, "static const int FONT_ATLAS_HEIGHT = %d;\n", h)
		b.WriteString("static const unsigned char FONT_ATLAS_DATA[] = {\n")
		for i, c := range cols {
			if i > 0 {
				b.WriteByte(',')
			}
			if i%16 == 0 {
				b.WriteString("\n    ")
			}
			fmt.Fprintf(&b, "%d,%d,%d,%d", c.R, c.G, c.B, c.A)
		}
		b.WriteString("\n};\n\n#endif\n")
		if err := os.WriteFile(toString(args[1]), []byte(b.String()), 0644); err != nil {
			return false, err
		}
		return true, nil
	})
	v.RegisterForeign("DrawTextCodepoint", func(args []interface{}) (interface{}, error) {
		if len(args) < 6 {
			return nil, fmt.Errorf("DrawTextCodepoint requires (fontId, codepoint, posX, posY, fontSize, tint)")
		}
		fontMu.Lock()
		font, ok := fonts[toString(args[0])]
		fontMu.Unlock()
		if !ok {
			return nil, fmt.Errorf("unknown font id: %s", toString(args[0]))
		}
		codepoint := rune(toInt32(args[1]))
		pos := rl.Vector2{X: toFloat32(args[2]), Y: toFloat32(args[3])}
		fontSize := toFloat32(args[4])
		c := rl.White
		if len(args) >= 10 {
			c = argsToColor(args, 5)
		}
		rl.DrawTextCodepoint(font, codepoint, pos, fontSize, c)
		return nil, nil
	})
	v.RegisterForeign("DrawTextCodepoints", func(args []interface{}) (interface{}, error) {
		if len(args) < 7 {
			return nil, fmt.Errorf("DrawTextCodepoints requires (fontId, codepointCount, posX, posY, fontSize, spacing, tint, ...codepoints)")
		}
		fontMu.Lock()
		font, ok := fonts[toString(args[0])]
		fontMu.Unlock()
		if !ok {
			return nil, fmt.Errorf("unknown font id: %s", toString(args[0]))
		}
		count := int(toInt32(args[1]))
		if count <= 0 || len(args) < 7+count {
			return nil, fmt.Errorf("DrawTextCodepoints needs codepointCount and that many codepoint values")
		}
		codepoints := make([]rune, count)
		for i := 0; i < count; i++ {
			codepoints[i] = rune(toInt32(args[7+i]))
		}
		pos := rl.Vector2{X: toFloat32(args[2]), Y: toFloat32(args[3])}
		fontSize := toFloat32(args[4])
		spacing := toFloat32(args[5])
		c := rl.White
		if len(args) >= 7+count+4 {
			c = argsToColor(args, 7+count)
		}
		rl.DrawTextCodepoints(font, codepoints, pos, fontSize, spacing, c)
		return nil, nil
	})
	v.RegisterForeign("GetGlyphIndex", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("GetGlyphIndex requires (fontId, codepoint)")
		}
		fontMu.Lock()
		font, ok := fonts[toString(args[0])]
		fontMu.Unlock()
		if !ok {
			return nil, fmt.Errorf("unknown font id: %s", toString(args[0]))
		}
		return int(rl.GetGlyphIndex(font, toInt32(args[1]))), nil
	})
	v.RegisterForeign("GetGlyphInfo", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("GetGlyphInfo requires (fontId, codepoint)")
		}
		fontMu.Lock()
		font, ok := fonts[toString(args[0])]
		fontMu.Unlock()
		if !ok {
			return nil, fmt.Errorf("unknown font id: %s", toString(args[0]))
		}
		info := rl.GetGlyphInfo(font, toInt32(args[1]))
		return []interface{}{int(info.Value), int(info.OffsetX), int(info.OffsetY), int(info.AdvanceX)}, nil
	})
	v.RegisterForeign("GetGlyphAtlasRec", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("GetGlyphAtlasRec requires (fontId, codepoint)")
		}
		fontMu.Lock()
		font, ok := fonts[toString(args[0])]
		fontMu.Unlock()
		if !ok {
			return nil, fmt.Errorf("unknown font id: %s", toString(args[0]))
		}
		rec := rl.GetGlyphAtlasRec(font, toInt32(args[1]))
		return []interface{}{float64(rec.X), float64(rec.Y), float64(rec.Width), float64(rec.Height)}, nil
	})
}
