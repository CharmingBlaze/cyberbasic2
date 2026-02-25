// Package raylib: full raygui bindings (gen2brain/raylib-go/raygui). Requires CGO.
// BASIC functions: GuiLabel, GuiButton, GuiCheckBox, GuiSlider, GuiProgressBar,
// GuiTextBox, GuiDropdownBox, GuiWindowBox, GuiGroupBox, GuiLine, GuiPanel.
package raylib

import (
	"fmt"
	"sync"

	"cyberbasic/compiler/vm"
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/gen2brain/raylib-go/raygui"
)

var (
	rayguiTextCache   = make(map[string]string)
	rayguiActiveCache = make(map[string]int32)
	rayguiMu          sync.Mutex
)

func rect(x, y, w, h float32) rl.Rectangle {
	return rl.Rectangle{X: x, Y: y, Width: w, Height: h}
}

func registerRaygui(v *vm.VM) {
	// GuiLabel(x, y, w, h, text)
	v.RegisterForeign("GuiLabel", func(args []interface{}) (interface{}, error) {
		if len(args) < 5 {
			return nil, fmt.Errorf("GuiLabel(x, y, w, h, text) requires 5 arguments")
		}
		b := rect(toFloat32(args[0]), toFloat32(args[1]), toFloat32(args[2]), toFloat32(args[3]))
		raygui.Label(b, toString(args[4]))
		return nil, nil
	})

	// GuiButton(x, y, w, h, text) -> 1 if clicked else 0
	v.RegisterForeign("GuiButton", func(args []interface{}) (interface{}, error) {
		if len(args) < 5 {
			return nil, fmt.Errorf("GuiButton(x, y, w, h, text) requires 5 arguments")
		}
		b := rect(toFloat32(args[0]), toFloat32(args[1]), toFloat32(args[2]), toFloat32(args[3]))
		clicked := raygui.Button(b, toString(args[4]))
		if clicked {
			return 1, nil
		}
		return 0, nil
	})

	// GuiCheckBox(x, y, w, h, text, checked) -> 1 if checked else 0
	v.RegisterForeign("GuiCheckBox", func(args []interface{}) (interface{}, error) {
		if len(args) < 6 {
			return nil, fmt.Errorf("GuiCheckBox(x, y, w, h, text, checked) requires 6 arguments")
		}
		b := rect(toFloat32(args[0]), toFloat32(args[1]), toFloat32(args[2]), toFloat32(args[3]))
		checked := toFloat64(args[5]) != 0
		out := raygui.CheckBox(b, toString(args[4]), checked)
		if out {
			return 1, nil
		}
		return 0, nil
	})
	// GuiCheckbox(text, x, y, checked) -> 1 if checked else 0 (simple 4-arg; default size)
	v.RegisterForeign("GuiCheckbox", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("GuiCheckbox(text, x, y, checked) requires 4 arguments")
		}
		b := rect(toFloat32(args[1]), toFloat32(args[2]), 24, 24)
		checked := toFloat64(args[3]) != 0
		out := raygui.CheckBox(b, toString(args[0]), checked)
		if out {
			return 1, nil
		}
		return 0, nil
	})

	// GuiSlider(x, y, w, h, textLeft, textRight, value, minValue, maxValue) or GuiSlider(x, y, w, min, max, value)
	v.RegisterForeign("GuiSlider", func(args []interface{}) (interface{}, error) {
		if len(args) >= 6 {
			minV := toFloat32(args[3])
			maxV := toFloat32(args[4])
			val := toFloat32(args[5])
			if maxV <= minV {
				maxV = minV + 1
			}
			b := rect(toFloat32(args[0]), toFloat32(args[1]), toFloat32(args[2]), 20)
			out := raygui.Slider(b, "", "", val, minV, maxV)
			return float64(out), nil
		}
		if len(args) < 9 {
			return nil, fmt.Errorf("GuiSlider(x, y, w, min, max, value) or (x, y, w, h, textLeft, textRight, value, min, max)")
		}
		b := rect(toFloat32(args[0]), toFloat32(args[1]), toFloat32(args[2]), toFloat32(args[3]))
		val := toFloat32(args[6])
		minV := toFloat32(args[7])
		maxV := toFloat32(args[8])
		if maxV <= minV {
			maxV = minV + 1
		}
		out := raygui.Slider(b, toString(args[4]), toString(args[5]), val, minV, maxV)
		return float64(out), nil
	})

	// GuiProgressBar(x, y, w, h, textLeft, textRight, value, minValue, maxValue) -> value
	v.RegisterForeign("GuiProgressBar", func(args []interface{}) (interface{}, error) {
		if len(args) < 9 {
			return nil, fmt.Errorf("GuiProgressBar(x, y, w, h, textLeft, textRight, value, min, max) requires 9 arguments")
		}
		b := rect(toFloat32(args[0]), toFloat32(args[1]), toFloat32(args[2]), toFloat32(args[3]))
		val := toFloat32(args[6])
		minV := toFloat32(args[7])
		maxV := toFloat32(args[8])
		if maxV <= minV {
			maxV = minV + 1
		}
		out := raygui.ProgressBar(b, toString(args[4]), toString(args[5]), val, minV, maxV)
		return float64(out), nil
	})

	// GuiTextbox(x, y, w, text) -> currentText (id = "tb_x_y" for cache)
	v.RegisterForeign("GuiTextbox", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("GuiTextbox(x, y, w, text) requires 4 arguments")
		}
		x, y, w := toFloat32(args[0]), toFloat32(args[1]), toFloat32(args[2])
		id := fmt.Sprintf("tb_%g_%g", x, y)
		b := rect(x, y, w, 24)
		rayguiMu.Lock()
		s := rayguiTextCache[id]
		if args[3] != nil {
			s = toString(args[3])
		}
		rayguiMu.Unlock()
		const maxTextSize = 256
		if len(s) > maxTextSize {
			s = s[:maxTextSize]
		}
		raygui.TextBox(b, &s, maxTextSize, true)
		rayguiMu.Lock()
		rayguiTextCache[id] = s
		rayguiMu.Unlock()
		return s, nil
	})
	// GuiTextBox(id, x, y, w, h, text) -> currentText (id used as cache key)
	v.RegisterForeign("GuiTextBox", func(args []interface{}) (interface{}, error) {
		if len(args) < 6 {
			return nil, fmt.Errorf("GuiTextBox(id, x, y, w, h, text) requires 6 arguments")
		}
		id := toString(args[0])
		b := rect(toFloat32(args[1]), toFloat32(args[2]), toFloat32(args[3]), toFloat32(args[4]))
		rayguiMu.Lock()
		s := rayguiTextCache[id]
		if args[5] != nil {
			s = toString(args[5])
		}
		rayguiMu.Unlock()
		const maxTextSize = 256
		if len(s) > maxTextSize {
			s = s[:maxTextSize]
		}
		raygui.TextBox(b, &s, maxTextSize, true)
		rayguiMu.Lock()
		rayguiTextCache[id] = s
		rayguiMu.Unlock()
		return s, nil
	})

	// GuiDropdownBox(id, x, y, w, h, itemsText, active) -> newActive (itemsText e.g. "One;Two;Three")
	v.RegisterForeign("GuiDropdownBox", func(args []interface{}) (interface{}, error) {
		if len(args) < 7 {
			return nil, fmt.Errorf("GuiDropdownBox(id, x, y, w, h, itemsText, active) requires 7 arguments")
		}
		id := toString(args[0])
		b := rect(toFloat32(args[1]), toFloat32(args[2]), toFloat32(args[3]), toFloat32(args[4]))
		itemsText := toString(args[5])
		rayguiMu.Lock()
		active := rayguiActiveCache[id]
		if args[6] != nil {
			active = toInt32(args[6])
		}
		rayguiMu.Unlock()
		raygui.DropdownBox(b, itemsText, &active, false)
		rayguiMu.Lock()
		rayguiActiveCache[id] = active
		rayguiMu.Unlock()
		return int(active), nil
	})

	// GuiWindowBox(x, y, w, h, title) -> 1 if close clicked else 0
	v.RegisterForeign("GuiWindowBox", func(args []interface{}) (interface{}, error) {
		if len(args) < 5 {
			return nil, fmt.Errorf("GuiWindowBox(x, y, w, h, title) requires 5 arguments")
		}
		b := rect(toFloat32(args[0]), toFloat32(args[1]), toFloat32(args[2]), toFloat32(args[3]))
		closed := raygui.WindowBox(b, toString(args[4]))
		if closed {
			return 1, nil
		}
		return 0, nil
	})

	// GuiGroupBox(x, y, w, h, text)
	v.RegisterForeign("GuiGroupBox", func(args []interface{}) (interface{}, error) {
		if len(args) < 5 {
			return nil, fmt.Errorf("GuiGroupBox(x, y, w, h, text) requires 5 arguments")
		}
		b := rect(toFloat32(args[0]), toFloat32(args[1]), toFloat32(args[2]), toFloat32(args[3]))
		raygui.GroupBox(b, toString(args[4]))
		return nil, nil
	})

	// GuiLine(x, y, w, h, text)
	v.RegisterForeign("GuiLine", func(args []interface{}) (interface{}, error) {
		if len(args) < 5 {
			return nil, fmt.Errorf("GuiLine(x, y, w, h, text) requires 5 arguments")
		}
		b := rect(toFloat32(args[0]), toFloat32(args[1]), toFloat32(args[2]), toFloat32(args[3]))
		raygui.Line(b, toString(args[4]))
		return nil, nil
	})

	// GuiPanel(x, y, w, h [, text])
	v.RegisterForeign("GuiPanel", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("GuiPanel(x, y, w, h) or (x, y, w, h, text)")
		}
		b := rect(toFloat32(args[0]), toFloat32(args[1]), toFloat32(args[2]), toFloat32(args[3]))
		text := ""
		if len(args) >= 5 {
			text = toString(args[4])
		}
		raygui.Panel(b, text)
		return nil, nil
	})
	// GuiWindow(title, x, y, w, h) -> 1 if close clicked else 0
	v.RegisterForeign("GuiWindow", func(args []interface{}) (interface{}, error) {
		if len(args) < 5 {
			return nil, fmt.Errorf("GuiWindow(title, x, y, w, h) requires 5 arguments")
		}
		title := toString(args[0])
		b := rect(toFloat32(args[1]), toFloat32(args[2]), toFloat32(args[3]), toFloat32(args[4]))
		closed := raygui.WindowBox(b, title)
		if closed {
			return 1, nil
		}
		return 0, nil
	})
	// GuiList(items, x, y, w, h) -> selected index (items = "Item1;Item2;Item3"); uses dropdown-style list
	v.RegisterForeign("GuiList", func(args []interface{}) (interface{}, error) {
		if len(args) < 5 {
			return nil, fmt.Errorf("GuiList(items, x, y, w, h) requires 5 arguments")
		}
		items := toString(args[0])
		b := rect(toFloat32(args[1]), toFloat32(args[2]), toFloat32(args[3]), toFloat32(args[4]))
		id := fmt.Sprintf("list_%g_%g", args[1], args[2])
		rayguiMu.Lock()
		active := rayguiActiveCache[id]
		rayguiMu.Unlock()
		raygui.DropdownBox(b, items, &active, false)
		rayguiMu.Lock()
		rayguiActiveCache[id] = active
		rayguiMu.Unlock()
		return int(active), nil
	})
	// GuiDropdown(items, x, y, w) -> selected index (items = "A;B;C")
	v.RegisterForeign("GuiDropdown", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("GuiDropdown(items, x, y, w) requires 4 arguments")
		}
		items := toString(args[0])
		b := rect(toFloat32(args[1]), toFloat32(args[2]), toFloat32(args[3]), 24)
		id := fmt.Sprintf("dd_%g_%g", args[1], args[2])
		rayguiMu.Lock()
		active := rayguiActiveCache[id]
		rayguiMu.Unlock()
		raygui.DropdownBox(b, items, &active, false)
		rayguiMu.Lock()
		rayguiActiveCache[id] = active
		rayguiMu.Unlock()
		return int(active), nil
	})
	// GuiProgressBar(x, y, w, value) -> value (0-1); simple 4-arg form
	v.RegisterForeign("GuiProgressBarSimple", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("GuiProgressBar(x, y, w, value) requires 4 arguments")
		}
		b := rect(toFloat32(args[0]), toFloat32(args[1]), toFloat32(args[2]), 24)
		val := toFloat32(args[3])
		out := raygui.ProgressBar(b, "", "", val, 0, 1)
		return float64(out), nil
	})
}
