// Package raylib: minimal immediate-mode UI (BeginUI, Label, Button, EndUI).
// Draws with raylib; Button returns true when clicked. Layout: vertical cursor.
package raylib

import (
	"fmt"
	"sync"

	"cyberbasic/compiler/vm"
	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	uiPadding   = 8
	uiLineH     = 24
	uiButtonH   = 28
	uiStartX    = 10
	uiStartY    = 10
	uiMinBtnW   = 80
)

var (
	uiX, uiY    int32
	uiMu        sync.Mutex
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
}
