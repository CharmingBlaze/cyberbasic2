// Package drawdot exposes draw.* DotObject aliases for common 2D raylib drawing foreigns.
package drawdot

import (
	"cyberbasic/compiler/bindings/dotargs"
	"cyberbasic/compiler/vm"
	"fmt"
	"strings"
)

// RegisterDrawDot registers global "draw" for draw.begin, draw.end, draw.clear, draw.text, draw.circle, draw.rectangle.
func RegisterDrawDot(v *vm.VM) {
	v.SetGlobal("draw", &drawModuleDot{v: v})
}

type drawModuleDot struct {
	v *vm.VM
}

func (d *drawModuleDot) GetProp([]string) (vm.Value, error) { return nil, nil }
func (d *drawModuleDot) SetProp([]string, vm.Value) error {
	return fmt.Errorf("draw: namespace is not assignable")
}

func (d *drawModuleDot) CallMethod(name string, args []vm.Value) (vm.Value, error) {
	ia := dotargs.From(args)
	switch strings.ToLower(name) {
	case "begin":
		return d.v.CallForeign("BeginDrawing", ia)
	case "end":
		return d.v.CallForeign("EndDrawing", ia)
	case "clear", "clearbackground":
		return d.v.CallForeign("ClearBackground", ia)
	case "text", "drawtext":
		return d.v.CallForeign("DrawText", ia)
	case "circle", "drawcircle":
		return d.v.CallForeign("DrawCircle", ia)
	case "rectangle", "rect", "drawrectangle":
		return d.v.CallForeign("DrawRectangle", ia)
	case "line", "drawline":
		return d.v.CallForeign("DrawLine", ia)
	default:
		return nil, fmt.Errorf("draw: unknown method %q (begin, end, clear, text, circle, rectangle, line)", name)
	}
}
