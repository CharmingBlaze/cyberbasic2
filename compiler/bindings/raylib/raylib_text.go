// Package raylib: text drawing (rtext) with default font.
package raylib

import (
	"cyberbasic/compiler/vm"
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func registerText(v *vm.VM) {
	v.RegisterForeign("DrawText", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("DrawText requires (text, x, y, fontSize)")
		}
		text := toString(args[0])
		x, y, fontSize := toInt32(args[1]), toInt32(args[2]), toInt32(args[3])
		c := rl.White
		if len(args) >= 8 {
			c = rl.NewColor(uint8(toInt32(args[4])), uint8(toInt32(args[5])), uint8(toInt32(args[6])), uint8(toInt32(args[7])))
		}
		rl.DrawText(text, x, y, fontSize, c)
		return nil, nil
	})
	// DrawTextSimple(text, x, y): draw text at (x,y), font size 20, white (use this for on-screen text; PRINT prints to console)
	v.RegisterForeign("DrawTextSimple", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("DrawTextSimple requires (text, x, y)")
		}
		text := toString(args[0])
		x, y := toInt32(args[1]), toInt32(args[2])
		rl.DrawText(text, x, y, 20, rl.White)
		return nil, nil
	})
	v.RegisterForeign("MeasureText", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("MeasureText requires (text, fontSize)")
		}
		return int(rl.MeasureText(toString(args[0]), toInt32(args[1]))), nil
	})
	v.RegisterForeign("DrawTextEx", func(args []interface{}) (interface{}, error) {
		if len(args) < 5 {
			return nil, fmt.Errorf("DrawTextEx requires (text, x, y, fontSize, spacing) and optional r,g,b,a")
		}
		font := rl.GetFontDefault()
		text := toString(args[0])
		pos := rl.Vector2{X: toFloat32(args[1]), Y: toFloat32(args[2])}
		fontSize := toFloat32(args[3])
		spacing := toFloat32(args[4])
		c := rl.White
		if len(args) >= 9 {
			c = argsToColor(args, 5)
		}
		rl.DrawTextEx(font, text, pos, fontSize, spacing, c)
		return nil, nil
	})
	v.RegisterForeign("DrawTextPro", func(args []interface{}) (interface{}, error) {
		if len(args) < 8 {
			return nil, fmt.Errorf("DrawTextPro requires (fontId, text, x, y, originX, originY, rotation, fontSize, spacing, tint)")
		}
		id := toString(args[0])
		fontMu.Lock()
		font, ok := fonts[id]
		fontMu.Unlock()
		if !ok {
			return nil, fmt.Errorf("unknown font id: %s", id)
		}
		text := toString(args[1])
		position := rl.Vector2{X: toFloat32(args[2]), Y: toFloat32(args[3])}
		origin := rl.Vector2{X: toFloat32(args[4]), Y: toFloat32(args[5])}
		rotation := toFloat32(args[6])
		fontSize := toFloat32(args[7])
		spacing := toFloat32(args[8])
		c := rl.White
		if len(args) >= 13 {
			c = argsToColor(args, 9)
		}
		rl.DrawTextPro(font, text, position, origin, rotation, fontSize, spacing, c)
		return nil, nil
	})
	v.RegisterForeign("SetTextLineSpacing", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("SetTextLineSpacing requires (spacing)")
		}
		rl.SetTextLineSpacing(int(toInt32(args[0])))
		return nil, nil
	})

	// Text string helpers (match C API semantics where possible)
	v.RegisterForeign("TextCopy", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return 0, nil
		}
		s := toString(args[0])
		return len(s), nil
	})
	v.RegisterForeign("TextIsEqual", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return false, nil
		}
		return toString(args[0]) == toString(args[1]), nil
	})
	v.RegisterForeign("TextLength", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return 0, nil
		}
		return len(toString(args[0])), nil
	})
	v.RegisterForeign("TextFormat", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return "", nil
		}
		format := toString(args[0])
		if len(args) == 1 {
			return format, nil
		}
		vals := make([]interface{}, len(args)-1)
		for i := 1; i < len(args); i++ {
			vals[i-1] = args[i]
		}
		return fmt.Sprintf(format, vals...), nil
	})
	v.RegisterForeign("TextSubtext", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return "", nil
		}
		text := toString(args[0])
		pos := int(toInt32(args[1]))
		length := int(toInt32(args[2]))
		if pos < 0 {
			pos = 0
		}
		if pos >= len(text) || length <= 0 {
			return "", nil
		}
		if pos+length > len(text) {
			length = len(text) - pos
		}
		return text[pos : pos+length], nil
	})
	v.RegisterForeign("TextReplace", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return "", nil
		}
		return strings.ReplaceAll(toString(args[0]), toString(args[1]), toString(args[2])), nil
	})
	v.RegisterForeign("TextInsert", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return "", nil
		}
		text := toString(args[0])
		insert := toString(args[1])
		pos := int(toInt32(args[2]))
		if pos < 0 {
			pos = 0
		}
		if pos > len(text) {
			pos = len(text)
		}
		return text[:pos] + insert + text[pos:], nil
	})
	v.RegisterForeign("TextJoin", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return "", nil
		}
		count := int(toInt32(args[0]))
		delimiter := toString(args[1])
		if count <= 0 || len(args) < 2+count {
			return "", nil
		}
		parts := make([]string, count)
		for i := 0; i < count; i++ {
			parts[i] = toString(args[2+i])
		}
		return strings.Join(parts, delimiter), nil
	})
	v.RegisterForeign("TextSplit", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return 0, nil
		}
		text := toString(args[0])
		delimRune := rune(toInt32(args[1]))
		if delimRune == 0 {
			return 0, nil
		}
		parts := strings.Split(text, string(delimRune))
		lastTextSplitMu.Lock()
		lastTextSplit = parts
		lastTextSplitMu.Unlock()
		return len(parts), nil
	})
	v.RegisterForeign("GetTextSplitItem", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return "", nil
		}
		idx := int(toInt32(args[0]))
		lastTextSplitMu.Lock()
		defer lastTextSplitMu.Unlock()
		if idx < 0 || idx >= len(lastTextSplit) {
			return "", nil
		}
		return lastTextSplit[idx], nil
	})
	v.RegisterForeign("TextAppend", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return 0, nil
		}
		return len(toString(args[0])) + len(toString(args[1])), nil
	})
	v.RegisterForeign("TextFindIndex", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return -1, nil
		}
		return strings.Index(toString(args[0]), toString(args[1])), nil
	})
	v.RegisterForeign("TextToUpper", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return "", nil
		}
		return strings.ToUpper(toString(args[0])), nil
	})
	v.RegisterForeign("TextToLower", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return "", nil
		}
		return strings.ToLower(toString(args[0])), nil
	})
	v.RegisterForeign("TextToPascal", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return "", nil
		}
		s := toString(args[0])
		words := strings.Fields(strings.ReplaceAll(strings.ReplaceAll(s, "_", " "), "-", " "))
		for i, w := range words {
			if len(w) > 0 {
				words[i] = strings.ToUpper(w[:1]) + strings.ToLower(w[1:])
			}
		}
		return strings.Join(words, ""), nil
	})
	v.RegisterForeign("TextToSnake", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return "", nil
		}
		var b strings.Builder
		s := toString(args[0])
		for i, r := range s {
			if r >= 'A' && r <= 'Z' {
				if i > 0 {
					b.WriteByte('_')
				}
				b.WriteRune(r + 32)
			} else {
				b.WriteRune(r)
			}
		}
		return b.String(), nil
	})
	v.RegisterForeign("TextToCamel", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return "", nil
		}
		s := toString(args[0])
		words := strings.Fields(strings.ReplaceAll(strings.ReplaceAll(s, "_", " "), "-", " "))
		for i, w := range words {
			if len(w) > 0 {
				if i == 0 {
					words[i] = strings.ToLower(w)
				} else {
					words[i] = strings.ToUpper(w[:1]) + strings.ToLower(w[1:])
				}
			}
		}
		return strings.Join(words, ""), nil
	})
	v.RegisterForeign("TextToInteger", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return 0, nil
		}
		n, _ := strconv.Atoi(toString(args[0]))
		return n, nil
	})
	v.RegisterForeign("TextToFloat", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return 0.0, nil
		}
		f, _ := strconv.ParseFloat(toString(args[0]), 64)
		return f, nil
	})

	// Codepoint / UTF-8 helpers
	v.RegisterForeign("GetCodepointCount", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return 0, nil
		}
		return utf8.RuneCountInString(toString(args[0])), nil
	})
	v.RegisterForeign("GetCodepoint", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return []interface{}{0, 0}, nil
		}
		text := toString(args[0])
		bytePos := int(toInt32(args[1]))
		if bytePos < 0 || bytePos >= len(text) {
			return []interface{}{int('?'), 0}, nil
		}
		r, size := utf8.DecodeRuneInString(text[bytePos:])
		if r == utf8.RuneError {
			return []interface{}{int('?'), 0}, nil
		}
		return []interface{}{int(r), size}, nil
	})
	v.RegisterForeign("GetCodepointNext", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return []interface{}{0, 0}, nil
		}
		text := toString(args[0])
		bytePos := int(toInt32(args[1]))
		if bytePos < 0 || bytePos >= len(text) {
			return []interface{}{int('?'), 0}, nil
		}
		r, size := utf8.DecodeRuneInString(text[bytePos:])
		if r == utf8.RuneError {
			return []interface{}{int('?'), 0}, nil
		}
		return []interface{}{int(r), size}, nil
	})
	v.RegisterForeign("GetCodepointPrevious", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return []interface{}{0, 0}, nil
		}
		text := toString(args[0])
		bytePos := int(toInt32(args[1]))
		if bytePos <= 0 || bytePos > len(text) {
			return []interface{}{int('?'), 0}, nil
		}
		r, size := utf8.DecodeLastRuneInString(text[:bytePos])
		if r == utf8.RuneError {
			return []interface{}{int('?'), 0}, nil
		}
		return []interface{}{int(r), size}, nil
	})
	v.RegisterForeign("CodepointToUTF8", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return "", nil
		}
		return string(rune(toInt32(args[0]))), nil
	})
	v.RegisterForeign("LoadCodepoints", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return 0, nil
		}
		text := toString(args[0])
		runes := []rune(text)
		lastCodepointsMu.Lock()
		lastCodepoints = runes
		lastCodepointsMu.Unlock()
		return len(runes), nil
	})
	v.RegisterForeign("UnloadCodepoints", func(args []interface{}) (interface{}, error) {
		lastCodepointsMu.Lock()
		lastCodepoints = nil
		lastCodepointsMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("GetLoadedCodepoint", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return 0, nil
		}
		idx := int(toInt32(args[0]))
		lastCodepointsMu.Lock()
		defer lastCodepointsMu.Unlock()
		if idx < 0 || idx >= len(lastCodepoints) {
			return 0, nil
		}
		return int(lastCodepoints[idx]), nil
	})
	v.RegisterForeign("LoadUTF8", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return "", nil
		}
		count := int(toInt32(args[0]))
		if count <= 0 || len(args) < 1+count {
			return "", nil
		}
		runes := make([]rune, count)
		for i := 0; i < count; i++ {
			runes[i] = rune(toInt32(args[1+i]))
		}
		return string(runes), nil
	})
	v.RegisterForeign("UnloadUTF8", func(args []interface{}) (interface{}, error) {
		return nil, nil
	})
}
