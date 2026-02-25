// Package raylib: window/config and blend mode flag constants (0-arg, return value for SetConfigFlags/BeginBlendMode).
package raylib

import (
	"cyberbasic/compiler/vm"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func registerFlags(v *vm.VM) {
	// Window/config flags (use with SetConfigFlags; call before InitWindow). Combine with OR.
	v.RegisterForeign("FLAG_VSYNC_HINT", func(args []interface{}) (interface{}, error) { return int(rl.FlagVsyncHint), nil })
	v.RegisterForeign("FLAG_FULLSCREEN_MODE", func(args []interface{}) (interface{}, error) { return int(rl.FlagFullscreenMode), nil })
	v.RegisterForeign("FLAG_WINDOW_RESIZABLE", func(args []interface{}) (interface{}, error) { return int(rl.FlagWindowResizable), nil })
	v.RegisterForeign("FLAG_WINDOW_UNDECORATED", func(args []interface{}) (interface{}, error) { return int(rl.FlagWindowUndecorated), nil })
	v.RegisterForeign("FLAG_WINDOW_HIDDEN", func(args []interface{}) (interface{}, error) { return int(rl.FlagWindowHidden), nil })
	v.RegisterForeign("FLAG_WINDOW_MINIMIZED", func(args []interface{}) (interface{}, error) { return int(rl.FlagWindowMinimized), nil })
	v.RegisterForeign("FLAG_WINDOW_MAXIMIZED", func(args []interface{}) (interface{}, error) { return int(rl.FlagWindowMaximized), nil })
	v.RegisterForeign("FLAG_WINDOW_UNFOCUSED", func(args []interface{}) (interface{}, error) { return int(rl.FlagWindowUnfocused), nil })
	v.RegisterForeign("FLAG_WINDOW_TOPMOST", func(args []interface{}) (interface{}, error) { return int(rl.FlagWindowTopmost), nil })
	v.RegisterForeign("FLAG_WINDOW_ALWAYS_RUN", func(args []interface{}) (interface{}, error) { return int(rl.FlagWindowAlwaysRun), nil })
	v.RegisterForeign("FLAG_MSAA_4X_HINT", func(args []interface{}) (interface{}, error) { return int(rl.FlagMsaa4xHint), nil })
	v.RegisterForeign("FLAG_INTERLACED_HINT", func(args []interface{}) (interface{}, error) { return int(rl.FlagInterlacedHint), nil })
	v.RegisterForeign("FLAG_WINDOW_HIGHDPI", func(args []interface{}) (interface{}, error) { return int(rl.FlagWindowHighdpi), nil })
	v.RegisterForeign("FLAG_BORDERLESS_WINDOWED_MODE", func(args []interface{}) (interface{}, error) { return int(rl.FlagBorderlessWindowedMode), nil })

	// Blend modes (use with BeginBlendMode)
	v.RegisterForeign("BLEND_ALPHA", func(args []interface{}) (interface{}, error) { return int(rl.BlendAlpha), nil })
	v.RegisterForeign("BLEND_ADDITIVE", func(args []interface{}) (interface{}, error) { return int(rl.BlendAdditive), nil })
	v.RegisterForeign("BLEND_MULTIPLIED", func(args []interface{}) (interface{}, error) { return int(rl.BlendMultiplied), nil })
	v.RegisterForeign("BLEND_ADD_COLORS", func(args []interface{}) (interface{}, error) { return int(rl.BlendAddColors), nil })
	v.RegisterForeign("BLEND_SUBTRACT_COLORS", func(args []interface{}) (interface{}, error) { return int(rl.BlendSubtractColors), nil })
	v.RegisterForeign("BLEND_CUSTOM", func(args []interface{}) (interface{}, error) { return int(rl.BlendCustom), nil })
}
