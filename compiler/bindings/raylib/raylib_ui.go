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

var (
	uiX, uiY       int32
	uiMu           sync.Mutex
	uiTextBoxBuf   = make(map[string]string)
	uiDropdownOpen = make(map[string]bool)
	uiFocusedText  string
)

func registerUI(v *vm.VM) {
	v.RegisterForeign("BeginUI", func(args []interface{}) (interface{}, error) {
		uiMu.Lock()
		uiX, uiY = uiStartX, uiStartY
		uiMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("EndUI", func(args []interface{}) (interface{}, error) {
		return nil, nil
	})
	v.RegisterForeign("Label", func(args []interface{}) (interface{}, error) {
		text := fmt.Sprint(args[0])
		uiMu.Lock()
		x, y := uiX, uiY
		uiY += uiLineH
		uiMu.Unlock()
		rl.DrawText(text, x, y, 20, rl.LightGray)
		return nil, nil
	})
	v.RegisterForeign("Button", func(args []interface{}) (interface{}, error) {
		text := fmt.Sprint(args[0])
		uiMu.Lock()
		x, y := int32(uiX), int32(uiY)
		uiY += uiButtonH + uiPadding
		uiMu.Unlock()
		w := rl.MeasureText(text, 20) + 24
		if w < uiMinBtnW {
			w = uiMinBtnW
		}
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
		uiMu.Lock()
		x, y := uiX, uiY
		uiY += uiSliderH + uiPadding
		uiMu.Unlock()
		w := int32(uiSliderW)
		h := uiSliderH
		mx := rl.GetMouseX()
		my := rl.GetMouseY()
		// Track thumb: left part = bar
		barX := x + int32(rl.MeasureText(text, 18)) + 8
		barW := w - (barX - x)
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
		uiMu.Lock()
		x, y := uiX, uiY
		uiY += uiLineH + uiPadding
		uiMu.Unlock()
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
		x, y := uiX, uiY
		uiY += uiTextBoxH + uiPadding
		uiMu.Unlock()
		w := uiTextBoxW
		h := uiTextBoxH
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
		x, y := uiX, uiY
		uiY += uiDropdownH + uiPadding
		uiMu.Unlock()
		w := uiTextBoxW
		h := uiDropdownH
		mx := rl.GetMouseX()
		my := rl.GetMouseY()
		barHit := mx >= x && mx <= x+int32(w) && my >= y && my <= y+int32(h)
		if barHit && rl.IsMouseButtonPressed(rl.MouseButtonLeft) {
			open = !open
		}
		if open {
			itemH := 24
			listH := len(items) * itemH
			for i, item := range items {
				iy := y + int32(h) + int32(i*itemH)
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
			uiMu.Lock()
			uiY += int32(listH)
			uiMu.Unlock()
			rl.DrawRectangleLines(x, y, int32(w), int32(h+listH), rl.LightGray)
		} else {
			rl.DrawRectangle(x, y, int32(w), int32(h), rl.Gray)
			rl.DrawRectangleLines(x, y, int32(w), int32(h), rl.LightGray)
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
		uiMu.Lock()
		x, y := uiX, uiY
		uiY += uiProgressH + uiPadding
		uiMu.Unlock()
		w := int32(uiSliderW)
		barX := x + int32(rl.MeasureText(text, 18)) + 8
		barW := w - (barX - x)
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
		uiMu.Lock()
		x, y := uiX, uiY
		uiY += uiLineH + uiWindowPad
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
		uiMu.Lock()
		x, y := uiX, uiY
		uiY += uiLineH + uiGroupPad
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
