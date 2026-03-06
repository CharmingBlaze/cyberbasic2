// Package raylib: immediate-mode UI (BeginUI, EndUI, Label, Button, Slider, Checkbox, TextBox, Dropdown, ProgressBar, WindowBox, GroupBox).
// Layout: vertical cursor. All widgets use raylib draw/input; no CGO.
package raylib

import (
	"fmt"
	"strings"
	"sync"

	"cyberbasic/compiler/vm"
	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	uiPadding    = 8
	uiLineH      = 24
	uiButtonH    = 28
	uiStartX     = 10
	uiStartY     = 10
	uiMinBtnW    = 80
	uiSliderW    = 120
	uiSliderH    = 24
	uiCheckBoxW  = 24
	uiTextBoxW   = 200
	uiTextBoxH   = 28
	uiDropdownH  = 28
	uiProgressH  = 24
	uiWindowPad  = 12
	uiGroupPad   = 8
)

const (
	LayoutVertical   = 0
	LayoutHorizontal = 1
	LayoutGrid       = 2
)

var (
	uiX, uiY           int32
	uiLayoutMode       int32 = LayoutVertical
	uiLayoutSpacing    int32 = 4
	uiLayoutPadL       int32 = 0
	uiLayoutPadT       int32 = 0
	uiLayoutPadR       int32 = 0
	uiLayoutPadB       int32 = 0
	uiGridCols         int32 = 2
	uiGridCol          int32
	uiLayoutStartX     int32
	uiLayoutStartY     int32
	uiLayoutRowH       int32
	uiMu               sync.Mutex
	uiTextBoxBuf       = make(map[string]string)
	uiDropdownOpen     = make(map[string]bool)
	uiFocusedText      string
	uiLayoutStack      []layoutState
)

type layoutState struct {
	x, y, mode, spacing, padL, padT, padR, padB int32
	gridCols, gridCol, rowH                    int32
}

// uiLayoutAdvance returns (x, y) for the current widget and advances layout.
// Call with widget width and height.
func uiLayoutAdvance(w, h int32) (x, y int32) {
	uiMu.Lock()
	defer uiMu.Unlock()
	x, y = uiX, uiY
	switch uiLayoutMode {
	case LayoutVertical:
		uiY += h + uiLayoutSpacing
	case LayoutHorizontal:
		uiX += w + uiLayoutSpacing
		if h > uiLayoutRowH {
			uiLayoutRowH = h
		}
	case LayoutGrid:
		if h > uiLayoutRowH {
			uiLayoutRowH = h
		}
		uiGridCol++
		if uiGridCol >= uiGridCols {
			uiGridCol = 0
			uiX = uiLayoutStartX
			uiY += uiLayoutRowH + uiLayoutSpacing
			uiLayoutRowH = 0
		} else {
			uiX += w + uiLayoutSpacing
		}
	default:
		uiY += h + uiLayoutSpacing
	}
	return x, y
}

func registerUI(v *vm.VM) {
	v.RegisterForeign("LAYOUT_VERTICAL", func(args []interface{}) (interface{}, error) { return int(LayoutVertical), nil })
	v.RegisterForeign("LAYOUT_HORIZONTAL", func(args []interface{}) (interface{}, error) { return int(LayoutHorizontal), nil })
	v.RegisterForeign("LAYOUT_GRID", func(args []interface{}) (interface{}, error) { return int(LayoutGrid), nil })
	v.RegisterForeign("BeginUI", func(args []interface{}) (interface{}, error) {
		uiMu.Lock()
		uiX, uiY = uiStartX, uiStartY
		uiLayoutMode = LayoutVertical
		uiLayoutStack = nil
		uiMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("EndUI", func(args []interface{}) (interface{}, error) {
		return nil, nil
	})
	v.RegisterForeign("BeginLayout", func(args []interface{}) (interface{}, error) {
		mode := int32(LayoutVertical)
		if len(args) >= 1 {
			mode = int32(toInt32(args[0]))
		}
		uiMu.Lock()
		uiLayoutMode = mode
		uiLayoutStartX, uiLayoutStartY = uiX, uiY
		uiGridCol = 0
		uiLayoutRowH = 0
		if mode == LayoutGrid && len(args) >= 2 {
			uiGridCols = int32(toInt32(args[1]))
			if uiGridCols < 1 {
				uiGridCols = 1
			}
		}
		uiMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("LayoutSpacing", func(args []interface{}) (interface{}, error) {
		if len(args) >= 1 {
			uiMu.Lock()
			uiLayoutSpacing = int32(toInt32(args[0]))
			if uiLayoutSpacing < 0 {
				uiLayoutSpacing = 0
			}
			uiMu.Unlock()
		}
		return nil, nil
	})
	v.RegisterForeign("LayoutPadding", func(args []interface{}) (interface{}, error) {
		if len(args) >= 4 {
			uiMu.Lock()
			uiLayoutPadL = int32(toInt32(args[0]))
			uiLayoutPadT = int32(toInt32(args[1]))
			uiLayoutPadR = int32(toInt32(args[2]))
			uiLayoutPadB = int32(toInt32(args[3]))
			uiMu.Unlock()
		}
		return nil, nil
	})
	v.RegisterForeign("BeginLayoutGroup", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("BeginLayoutGroup(id, mode) requires 2 arguments")
		}
		mode := int32(toInt32(args[1]))
		uiMu.Lock()
		uiLayoutStack = append(uiLayoutStack, layoutState{
			x: uiX, y: uiY, mode: uiLayoutMode, spacing: uiLayoutSpacing,
			padL: uiLayoutPadL, padT: uiLayoutPadT, padR: uiLayoutPadR, padB: uiLayoutPadB,
			gridCols: uiGridCols, gridCol: uiGridCol, rowH: uiLayoutRowH,
		})
		uiLayoutMode = mode
		uiX += uiLayoutPadL
		uiY += uiLayoutPadT
		uiLayoutStartX, uiLayoutStartY = uiX, uiY
		uiGridCol = 0
		uiLayoutRowH = 0
		if mode == LayoutGrid && len(args) >= 3 {
			uiGridCols = int32(toInt32(args[2]))
		}
		uiMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("EndLayoutGroup", func(args []interface{}) (interface{}, error) {
		uiMu.Lock()
		if len(uiLayoutStack) > 0 {
			prev := uiLayoutStack[len(uiLayoutStack)-1]
			uiLayoutStack = uiLayoutStack[:len(uiLayoutStack)-1]
			uiX = prev.x
			if prev.mode == LayoutHorizontal {
				uiY = prev.y + uiLayoutRowH + prev.spacing + prev.padB
			} else {
				uiY = uiY + prev.padB
			}
			uiLayoutMode = prev.mode
			uiLayoutSpacing = prev.spacing
			uiLayoutPadL, uiLayoutPadT = prev.padL, prev.padT
			uiLayoutPadR, uiLayoutPadB = prev.padR, prev.padB
			uiGridCols, uiGridCol = prev.gridCols, prev.gridCol
			uiLayoutRowH = prev.rowH
		}
		uiMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("Label", func(args []interface{}) (interface{}, error) {
		text := fmt.Sprint(args[0])
		w := int32(rl.MeasureText(text, 20))
		x, y := uiLayoutAdvance(w, uiLineH)
		rl.DrawText(text, x, y, 20, rl.LightGray)
		return nil, nil
	})
	v.RegisterForeign("LabelAt", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("LabelAt(x, y, text) requires 3 arguments")
		}
		x := int32(toInt32(args[0]))
		y := int32(toInt32(args[1]))
		text := fmt.Sprint(args[2])
		rl.DrawText(text, x, y, 20, rl.LightGray)
		return nil, nil
	})
	v.RegisterForeign("Button", func(args []interface{}) (interface{}, error) {
		text := fmt.Sprint(args[0])
		w := rl.MeasureText(text, 20) + 24
		if w < uiMinBtnW {
			w = uiMinBtnW
		}
		x, y := uiLayoutAdvance(w, uiButtonH)
		h := uiButtonH
		mx := rl.GetMouseX()
		my := rl.GetMouseY()
		hit := mx >= x && mx <= x+w && my >= y && my <= y+int32(h)
		clicked := hit && rl.IsMouseButtonPressed(rl.MouseButtonLeft)
		// Draw
		if hit {
			rl.DrawRectangle(int32(x), int32(y), w, int32(h), rl.DarkGray)
		} else {
			rl.DrawRectangle(int32(x), int32(y), w, int32(h), rl.Gray)
		}
		rl.DrawRectangleLines(int32(x), int32(y), w, int32(h), rl.LightGray)
		tx := x + (w-int32(rl.MeasureText(text, 20)))/2
		ty := y + (int32(h)-20)/2
		rl.DrawText(text, tx, ty, 20, rl.White)
		return clicked, nil
	})
	v.RegisterForeign("ButtonAt", func(args []interface{}) (interface{}, error) {
		if len(args) < 5 {
			return nil, fmt.Errorf("ButtonAt(x, y, w, h, text) requires 5 arguments")
		}
		x := int32(toInt32(args[0]))
		y := int32(toInt32(args[1]))
		w := int32(toInt32(args[2]))
		h := int32(toInt32(args[3]))
		text := fmt.Sprint(args[4])
		if w < uiMinBtnW {
			w = uiMinBtnW
		}
		mx := rl.GetMouseX()
		my := rl.GetMouseY()
		hit := mx >= x && mx <= x+w && my >= y && my <= y+h
		clicked := hit && rl.IsMouseButtonPressed(rl.MouseButtonLeft)
		if hit {
			rl.DrawRectangle(x, y, w, h, rl.DarkGray)
		} else {
			rl.DrawRectangle(x, y, w, h, rl.Gray)
		}
		rl.DrawRectangleLines(x, y, w, h, rl.LightGray)
		tx := x + (w-int32(rl.MeasureText(text, 20)))/2
		ty := y + (h-20)/2
		rl.DrawText(text, tx, ty, 20, rl.White)
		return clicked, nil
	})

	// Slider(text, value, min, max) -> value
	v.RegisterForeign("Slider", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("Slider(text, value, min, max) requires 4 arguments")
		}
		text := fmt.Sprint(args[0])
		val := toFloat64(args[1])
		minV := toFloat64(args[2])
		maxV := toFloat64(args[3])
		if maxV <= minV {
			maxV = minV + 1
		}
		labelW := int32(rl.MeasureText(text, 18)) + 8
		w := labelW + int32(uiSliderW)
		h := int32(uiSliderH)
		x, y := uiLayoutAdvance(w, h)
		mx := rl.GetMouseX()
		my := rl.GetMouseY()
		// Track thumb: left part = bar
		barX := x + labelW
		barW := w - labelW
		norm := (val - minV) / (maxV - minV)
		if norm < 0 {
			norm = 0
		}
		if norm > 1 {
			norm = 1
		}
		thumbX := barX + int32(float64(barW-12)*norm)
		dragging := rl.IsMouseButtonDown(rl.MouseButtonLeft) && mx >= barX && mx <= barX+barW && my >= y && my <= y+int32(h)
		if dragging {
			norm = float64(mx-barX) / float64(barW-12)
			if norm < 0 {
				norm = 0
			}
			if norm > 1 {
				norm = 1
			}
			val = minV + norm*(maxV-minV)
		}
		rl.DrawText(text, x, y, 18, rl.LightGray)
		rl.DrawRectangle(barX, y, barW, int32(h), rl.DarkGray)
		rl.DrawRectangle(thumbX, y, 12, int32(h), rl.Gray)
		rl.DrawRectangleLines(barX, y, barW, int32(h), rl.LightGray)
		return val, nil
	})

	// Checkbox(text, checked) -> checked (1 or 0)
	v.RegisterForeign("Checkbox", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("Checkbox(text, checked) requires 2 arguments")
		}
		text := fmt.Sprint(args[0])
		checked := toInt32(args[1]) != 0
		w := int32(uiCheckBoxW) + int32(rl.MeasureText(text, 18)) + 8
		h := int32(uiLineH)
		x, y := uiLayoutAdvance(w, h)
		mx := rl.GetMouseX()
		my := rl.GetMouseY()
		hit := mx >= x && mx <= x+w && my >= y && my <= y+uiCheckBoxW
		if hit && rl.IsMouseButtonPressed(rl.MouseButtonLeft) {
			checked = !checked
		}
		rl.DrawRectangle(x, y, uiCheckBoxW, uiCheckBoxW, rl.Gray)
		if checked {
			rl.DrawRectangle(x+4, y+4, uiCheckBoxW-8, uiCheckBoxW-8, rl.White)
		}
		rl.DrawRectangleLines(x, y, uiCheckBoxW, uiCheckBoxW, rl.LightGray)
		rl.DrawText(text, x+uiCheckBoxW+8, y+2, 18, rl.LightGray)
		if checked {
			return 1, nil
		}
		return 0, nil
	})
	v.RegisterForeign("CheckboxAt", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("CheckboxAt(x, y, label, checked) requires 4 arguments")
		}
		x := int32(toInt32(args[0]))
		y := int32(toInt32(args[1]))
		text := fmt.Sprint(args[2])
		checked := toInt32(args[3]) != 0
		mx := rl.GetMouseX()
		my := rl.GetMouseY()
		hit := mx >= x && mx <= x+uiCheckBoxW+int32(rl.MeasureText(text, 18))+8 && my >= y && my <= y+uiCheckBoxW
		if hit && rl.IsMouseButtonPressed(rl.MouseButtonLeft) {
			checked = !checked
		}
		rl.DrawRectangle(x, y, uiCheckBoxW, uiCheckBoxW, rl.Gray)
		if checked {
			rl.DrawRectangle(x+4, y+4, uiCheckBoxW-8, uiCheckBoxW-8, rl.White)
		}
		rl.DrawRectangleLines(x, y, uiCheckBoxW, uiCheckBoxW, rl.LightGray)
		rl.DrawText(text, x+uiCheckBoxW+8, y+2, 18, rl.LightGray)
		if checked {
			return 1, nil
		}
		return 0, nil
	})
	v.RegisterForeign("SliderAt", func(args []interface{}) (interface{}, error) {
		if len(args) < 6 {
			return nil, fmt.Errorf("SliderAt(x, y, width, min, max, value) requires 6 arguments")
		}
		x := int32(toInt32(args[0]))
		y := int32(toInt32(args[1]))
		barW := int32(toInt32(args[2]))
		minV := toFloat64(args[3])
		maxV := toFloat64(args[4])
		val := toFloat64(args[5])
		if maxV <= minV {
			maxV = minV + 1
		}
		norm := (val - minV) / (maxV - minV)
		if norm < 0 {
			norm = 0
		}
		if norm > 1 {
			norm = 1
		}
		mx := rl.GetMouseX()
		my := rl.GetMouseY()
		dragging := rl.IsMouseButtonDown(rl.MouseButtonLeft) && mx >= x && mx <= x+barW && my >= y && my <= y+int32(uiSliderH)
		if dragging {
			norm = float64(mx-x) / float64(barW-12)
			if norm < 0 {
				norm = 0
			}
			if norm > 1 {
				norm = 1
			}
			val = minV + norm*(maxV-minV)
		}
		rl.DrawRectangle(x, y, barW, int32(uiSliderH), rl.DarkGray)
		thumbX := x + int32(float64(barW-12)*norm)
		rl.DrawRectangle(thumbX, y, 12, int32(uiSliderH), rl.Gray)
		rl.DrawRectangleLines(x, y, barW, int32(uiSliderH), rl.LightGray)
		return val, nil
	})

	// TextBox(id, text) -> text (edited)
	v.RegisterForeign("TextBox", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("TextBox(id, text) requires 2 arguments")
		}
		id := fmt.Sprint(args[0])
		initial := fmt.Sprint(args[1])
		uiMu.Lock()
		if uiFocusedText != id {
			uiTextBoxBuf[id] = initial
		}
		buf := uiTextBoxBuf[id]
		uiMu.Unlock()
		w := int32(uiTextBoxW)
		h := int32(uiTextBoxH)
		x, y := uiLayoutAdvance(w, h)
		mx := rl.GetMouseX()
		my := rl.GetMouseY()
		hit := mx >= x && mx <= x+int32(w) && my >= y && my <= y+int32(h)
		if hit && rl.IsMouseButtonPressed(rl.MouseButtonLeft) {
			uiFocusedText = id
		}
		if uiFocusedText == id && !hit && rl.IsMouseButtonPressed(rl.MouseButtonLeft) {
			uiFocusedText = ""
		}
		if uiFocusedText == id {
			c := rl.GetCharPressed()
			for c > 0 {
				buf += string(rune(c))
				c = rl.GetCharPressed()
			}
			if rl.IsKeyPressed(rl.KeyBackspace) && len(buf) > 0 {
				buf = buf[:len(buf)-1]
			}
			uiTextBoxBuf[id] = buf
		}
		if uiFocusedText == id {
			rl.DrawRectangle(x, y, int32(w), int32(h), rl.DarkGray)
		} else {
			rl.DrawRectangle(x, y, int32(w), int32(h), rl.Gray)
		}
		rl.DrawRectangleLines(x, y, int32(w), int32(h), rl.LightGray)
		rl.DrawText(buf, x+4, y+4, 18, rl.White)
		return buf, nil
	})

	// Dropdown(id, itemsText, activeIndex). itemsText = "Item1;Item2;Item3". Returns new activeIndex (0-based).
	v.RegisterForeign("Dropdown", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("Dropdown(id, itemsText, activeIndex) requires 3 arguments")
		}
		id := fmt.Sprint(args[0])
		itemsText := fmt.Sprint(args[1])
		active := int(toInt32(args[2]))
		items := strings.Split(itemsText, ";")
		if len(items) == 0 {
			items = []string{""}
		}
		if active < 0 || active >= len(items) {
			active = 0
		}
		uiMu.Lock()
		open := uiDropdownOpen[id]
		uiMu.Unlock()
		w := int32(uiTextBoxW)
		barH := int32(uiDropdownH)
		h := barH
		if open {
			h += int32(len(items) * 24)
		}
		x, y := uiLayoutAdvance(w, h)
		mx := rl.GetMouseX()
		my := rl.GetMouseY()
		barHit := mx >= x && mx <= x+int32(w) && my >= y && my <= y+int32(barH)
		if barHit && rl.IsMouseButtonPressed(rl.MouseButtonLeft) {
			open = !open
		}
		if open {
			itemH := 24
			for i, item := range items {
				iy := y + int32(barH) + int32(i*itemH)
				itemHit := mx >= x && mx <= x+int32(w) && my >= iy && my <= iy+int32(itemH)
				if itemHit && rl.IsMouseButtonPressed(rl.MouseButtonLeft) {
					active = i
					open = false
				}
				if i == active {
					rl.DrawRectangle(x, iy, int32(w), int32(itemH), rl.DarkGray)
				} else {
					rl.DrawRectangle(x, iy, int32(w), int32(itemH), rl.Gray)
				}
				rl.DrawText(item, x+4, iy+2, 18, rl.White)
			}
			rl.DrawRectangleLines(x, y, w, barH+int32(len(items)*24), rl.LightGray)
		} else {
			rl.DrawRectangle(x, y, int32(w), int32(barH), rl.Gray)
			rl.DrawRectangleLines(x, y, int32(w), int32(barH), rl.LightGray)
			rl.DrawText(items[active], x+4, y+4, 18, rl.White)
		}
		uiDropdownOpen[id] = open
		return active, nil
	})

	// ProgressBar(text, value, min, max) -> value (draw only)
	v.RegisterForeign("ProgressBar", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("ProgressBar(text, value, min, max) requires 4 arguments")
		}
		text := fmt.Sprint(args[0])
		val := toFloat64(args[1])
		minV := toFloat64(args[2])
		maxV := toFloat64(args[3])
		if maxV <= minV {
			maxV = minV + 1
		}
		norm := (val - minV) / (maxV - minV)
		if norm < 0 {
			norm = 0
		}
		if norm > 1 {
			norm = 1
		}
		labelW := int32(rl.MeasureText(text, 18)) + 8
		w := labelW + int32(uiSliderW)
		h := int32(uiProgressH)
		x, y := uiLayoutAdvance(w, h)
		barX := x + labelW
		barW := w - labelW
		rl.DrawText(text, x, y, 18, rl.LightGray)
		rl.DrawRectangle(barX, y, barW, uiProgressH, rl.DarkGray)
		rl.DrawRectangle(barX, y, int32(float64(barW)*norm), uiProgressH, rl.Gray)
		rl.DrawRectangleLines(barX, y, barW, uiProgressH, rl.LightGray)
		return val, nil
	})

	// WindowBox(title) - container; advances layout and draws a titled box
	v.RegisterForeign("WindowBox", func(args []interface{}) (interface{}, error) {
		title := ""
		if len(args) >= 1 {
			title = fmt.Sprint(args[0])
		}
		w := int32(300)
		h := int32(uiLineH + uiWindowPad)
		x, y := uiLayoutAdvance(w, h)
		uiMu.Lock()
		uiX += uiWindowPad
		uiMu.Unlock()
		rl.DrawRectangle(x, y, 300, 2, rl.LightGray)
		rl.DrawText(title, x, y-2, 18, rl.White)
		return nil, nil
	})
	v.RegisterForeign("EndWindowBox", func(args []interface{}) (interface{}, error) {
		uiMu.Lock()
		uiX -= uiWindowPad
		uiY += uiWindowPad
		uiMu.Unlock()
		return nil, nil
	})

	// GroupBox(text) - container
	v.RegisterForeign("GroupBox", func(args []interface{}) (interface{}, error) {
		text := ""
		if len(args) >= 1 {
			text = fmt.Sprint(args[0])
		}
		w := int32(280)
		h := int32(uiLineH + uiGroupPad)
		x, y := uiLayoutAdvance(w, h)
		uiMu.Lock()
		uiX += uiGroupPad
		uiMu.Unlock()
		rl.DrawRectangleLines(x, y, 280, 2, rl.Gray)
		rl.DrawText(text, x, y-2, 16, rl.LightGray)
		return nil, nil
	})
	v.RegisterForeign("EndGroupBox", func(args []interface{}) (interface{}, error) {
		uiMu.Lock()
		uiX -= uiGroupPad
		uiY += uiGroupPad
		uiMu.Unlock()
		return nil, nil
	})
}
