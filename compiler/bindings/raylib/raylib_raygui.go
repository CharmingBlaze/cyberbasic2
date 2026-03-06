// Package raylib: full raygui bindings (gen2brain/raylib-go/raygui). Requires CGO.
// BASIC functions: GuiLabel, GuiButton, GuiCheckBox, GuiSlider, GuiProgressBar,
// GuiTextBox, GuiDropdownBox, GuiWindowBox, GuiGroupBox, GuiLine, GuiPanel.
package raylib

import (
	"fmt"
	"strings"
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

// applyUIStylePreset applies a named style preset via raygui.SetStyle.
// Presets: "default", "dark", "light", "cyber"
func applyUIStylePreset(name string) {
	// Property IDs: BORDER_COLOR_NORMAL=0, BASE_COLOR_NORMAL=1, TEXT_COLOR_NORMAL=2,
	// BORDER_COLOR_FOCUSED=3, BASE_COLOR_FOCUSED=4, TEXT_COLOR_FOCUSED=5,
	// BORDER_COLOR_PRESSED=6, BASE_COLOR_PRESSED=7, TEXT_COLOR_PRESSED=8,
	// BORDER_COLOR_DISABLED=9, BASE_COLOR_DISABLED=10, TEXT_COLOR_DISABLED=11,
	// BORDER_WIDTH=12, TEXT_PADDING=13, TEXT_ALIGNMENT=14
	// BACKGROUND_COLOR=19 (16+3)
	const (
		borderColorNormal = 0
		baseColorNormal   = 1
		textColorNormal   = 2
		borderColorFocused = 3
		baseColorFocused  = 4
		textColorFocused  = 5
		borderColorPressed = 6
		baseColorPressed  = 7
		textColorPressed  = 8
		borderColorDisabled = 9
		baseColorDisabled = 10
		textColorDisabled = 11
		borderWidth      = 12
		textPadding      = 13
		backgroundColor  = 19
	)
	set := func(ctrl, prop int32, val raygui.PropertyValue) {
		raygui.SetStyle(raygui.ControlID(ctrl), raygui.PropertyID(prop), val)
	}
	c := func(r, g, b, a uint8) raygui.PropertyValue {
		return raygui.NewColorPropertyValue(rl.Color{R: r, G: g, B: b, A: a})
	}
	switch name {
	case "default":
		raygui.LoadStyleDefault()
	case "dark":
		// Dark theme: dark grays, light text
		set(0, borderColorNormal, c(80, 80, 90, 255))
		set(0, baseColorNormal, c(60, 60, 70, 255))
		set(0, textColorNormal, c(220, 220, 230, 255))
		set(0, borderColorFocused, c(100, 100, 120, 255))
		set(0, baseColorFocused, c(80, 80, 100, 255))
		set(0, textColorFocused, c(240, 240, 255, 255))
		set(0, borderColorPressed, c(120, 120, 140, 255))
		set(0, baseColorPressed, c(100, 100, 120, 255))
		set(0, textColorPressed, c(255, 255, 255, 255))
		set(0, borderColorDisabled, c(50, 50, 55, 255))
		set(0, baseColorDisabled, c(40, 40, 45, 255))
		set(0, textColorDisabled, c(120, 120, 130, 255))
		set(0, borderWidth, 1)
		set(0, textPadding, 4)
		set(0, backgroundColor, c(35, 35, 40, 255))
	case "light":
		// Light theme: light grays, dark text
		set(0, borderColorNormal, c(180, 180, 190, 255))
		set(0, baseColorNormal, c(240, 240, 245, 255))
		set(0, textColorNormal, c(40, 40, 50, 255))
		set(0, borderColorFocused, c(100, 150, 200, 255))
		set(0, baseColorFocused, c(230, 238, 248, 255))
		set(0, textColorFocused, c(30, 80, 150, 255))
		set(0, borderColorPressed, c(80, 130, 180, 255))
		set(0, baseColorPressed, c(210, 225, 242, 255))
		set(0, textColorPressed, c(20, 70, 140, 255))
		set(0, borderColorDisabled, c(200, 200, 205, 255))
		set(0, baseColorDisabled, c(230, 230, 235, 255))
		set(0, textColorDisabled, c(150, 150, 160, 255))
		set(0, borderWidth, 1)
		set(0, textPadding, 4)
		set(0, backgroundColor, c(250, 250, 252, 255))
	case "cyber":
		// Cyber/neon: dark base, cyan/magenta accents
		set(0, borderColorNormal, c(0, 200, 220, 255))
		set(0, baseColorNormal, c(15, 25, 45, 255))
		set(0, textColorNormal, c(0, 255, 255, 255))
		set(0, borderColorFocused, c(255, 0, 200, 255))
		set(0, baseColorFocused, c(25, 35, 65, 255))
		set(0, textColorFocused, c(255, 100, 255, 255))
		set(0, borderColorPressed, c(255, 50, 220, 255))
		set(0, baseColorPressed, c(40, 50, 90, 255))
		set(0, textColorPressed, c(255, 150, 255, 255))
		set(0, borderColorDisabled, c(50, 80, 90, 255))
		set(0, baseColorDisabled, c(10, 15, 25, 255))
		set(0, textColorDisabled, c(80, 120, 130, 255))
		set(0, borderWidth, 2)
		set(0, textPadding, 6)
		set(0, backgroundColor, c(5, 10, 20, 255))
	default:
		raygui.LoadStyleDefault()
	}
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
	// GuiTextBoxId(id, x, y, w, h, text) -> currentText (id used as cache key; name avoids collision with GuiTextbox)
	v.RegisterForeign("GuiTextBoxId", func(args []interface{}) (interface{}, error) {
		if len(args) < 6 {
			return nil, fmt.Errorf("GuiTextBoxId(id, x, y, w, h, text) requires 6 arguments")
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

	// Theme and style (raygui)
	v.RegisterForeign("GuiLoadStyle", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("GuiLoadStyle(filePath) requires 1 argument")
		}
		raygui.LoadStyle(toString(args[0]))
		return nil, nil
	})
	v.RegisterForeign("GuiLoadStyleDefault", func(args []interface{}) (interface{}, error) {
		raygui.LoadStyleDefault()
		return nil, nil
	})
	v.RegisterForeign("LoadUIStyle", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("LoadUIStyle(name) requires 1 argument")
		}
		name := strings.ToLower(strings.TrimSpace(toString(args[0])))
		applyUIStylePreset(name)
		return nil, nil
	})
	v.RegisterForeign("GuiSetStyle", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("GuiSetStyle(controlId, propertyId, value) requires 3 arguments")
		}
		control := raygui.ControlID(toInt32(args[0]))
		property := raygui.PropertyID(toInt32(args[1]))
		value := raygui.PropertyValue(toInt32(args[2]))
		raygui.SetStyle(control, property, value)
		return nil, nil
	})
	v.RegisterForeign("GuiGetStyle", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("GuiGetStyle(controlId, propertyId) requires 2 arguments")
		}
		control := raygui.ControlID(toInt32(args[0]))
		property := raygui.PropertyID(toInt32(args[1]))
		v := raygui.GetStyle(control, property)
		return int(v), nil
	})
}
