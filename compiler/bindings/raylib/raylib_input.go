// Package raylib: keyboard and mouse input.
package raylib

import (
	"cyberbasic/compiler/vm"
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func registerInput(v *vm.VM) {
	v.RegisterForeign("IsMouseButtonPressed", func(args []interface{}) (interface{}, error) {
		button := rl.MouseButtonLeft
		if len(args) >= 1 {
			switch toInt32(args[0]) {
			case 1:
				button = rl.MouseButtonRight
			case 2:
				button = rl.MouseButtonMiddle
			default:
				button = rl.MouseButtonLeft
			}
		}
		return rl.IsMouseButtonPressed(button), nil
	})
	v.RegisterForeign("GetMouseX", func(args []interface{}) (interface{}, error) {
		return int(rl.GetMouseX()), nil
	})
	v.RegisterForeign("GetMouseY", func(args []interface{}) (interface{}, error) {
		return int(rl.GetMouseY()), nil
	})
	v.RegisterForeign("IsKeyPressed", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("IsKeyPressed requires (key)")
		}
		return rl.IsKeyPressed(int32(toInt32(args[0]))), nil
	})
	v.RegisterForeign("IsKeyDown", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("IsKeyDown requires (key)")
		}
		return rl.IsKeyDown(int32(toInt32(args[0]))), nil
	})
	v.RegisterForeign("KeyDown", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("KeyDown requires (key)")
		}
		return rl.IsKeyDown(int32(toInt32(args[0]))), nil
	})
	v.RegisterForeign("KeyPressed", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("KeyPressed requires (key)")
		}
		return rl.IsKeyPressed(int32(toInt32(args[0]))), nil
	})
	v.RegisterForeign("IsKeyReleased", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("IsKeyReleased requires (key)")
		}
		return rl.IsKeyReleased(int32(toInt32(args[0]))), nil
	})
	v.RegisterForeign("IsKeyUp", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("IsKeyUp requires (key)")
		}
		return rl.IsKeyUp(int32(toInt32(args[0]))), nil
	})
	v.RegisterForeign("GetKeyPressed", func(args []interface{}) (interface{}, error) {
		return int(rl.GetKeyPressed()), nil
	})
	v.RegisterForeign("SetExitKey", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("SetExitKey requires (key)")
		}
		rl.SetExitKey(int32(toInt32(args[0])))
		return nil, nil
	})
	v.RegisterForeign("IsMouseButtonDown", func(args []interface{}) (interface{}, error) {
		button := rl.MouseButtonLeft
		if len(args) >= 1 {
			switch toInt32(args[0]) {
			case 1:
				button = rl.MouseButtonRight
			case 2:
				button = rl.MouseButtonMiddle
			default:
				button = rl.MouseButtonLeft
			}
		}
		return rl.IsMouseButtonDown(button), nil
	})
	v.RegisterForeign("IsMouseButtonReleased", func(args []interface{}) (interface{}, error) {
		button := rl.MouseButtonLeft
		if len(args) >= 1 {
			switch toInt32(args[0]) {
			case 1:
				button = rl.MouseButtonRight
			case 2:
				button = rl.MouseButtonMiddle
			default:
				button = rl.MouseButtonLeft
			}
		}
		return rl.IsMouseButtonReleased(button), nil
	})
	v.RegisterForeign("GetMouseWheelMove", func(args []interface{}) (interface{}, error) {
		return float64(rl.GetMouseWheelMove()), nil
	})
	v.RegisterForeign("SetMousePosition", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("SetMousePosition requires (x, y)")
		}
		rl.SetMousePosition(int(toInt32(args[0])), int(toInt32(args[1])))
		return nil, nil
	})
	v.RegisterForeign("SetMouseOffset", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("SetMouseOffset requires (offsetX, offsetY)")
		}
		rl.SetMouseOffset(int(toInt32(args[0])), int(toInt32(args[1])))
		return nil, nil
	})
	v.RegisterForeign("SetMouseScale", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("SetMouseScale requires (scaleX, scaleY)")
		}
		rl.SetMouseScale(toFloat32(args[0]), toFloat32(args[1]))
		return nil, nil
	})
	v.RegisterForeign("HideCursor", func(args []interface{}) (interface{}, error) {
		rl.HideCursor()
		return nil, nil
	})
	v.RegisterForeign("ShowCursor", func(args []interface{}) (interface{}, error) {
		rl.ShowCursor()
		return nil, nil
	})
	v.RegisterForeign("GetMousePosition", func(args []interface{}) (interface{}, error) {
		pos := rl.GetMousePosition()
		return []interface{}{float64(pos.X), float64(pos.Y)}, nil
	})
	v.RegisterForeign("GetMouseDelta", func(args []interface{}) (interface{}, error) {
		delta := rl.GetMouseDelta()
		return []interface{}{float64(delta.X), float64(delta.Y)}, nil
	})
	v.RegisterForeign("MouseDown", func(args []interface{}) (interface{}, error) {
		button := rl.MouseButtonLeft
		if len(args) >= 1 {
			switch toInt32(args[0]) {
			case 1:
				button = rl.MouseButtonRight
			case 2:
				button = rl.MouseButtonMiddle
			default:
				button = rl.MouseButtonLeft
			}
		}
		return rl.IsMouseButtonDown(button), nil
	})
	v.RegisterForeign("MousePressed", func(args []interface{}) (interface{}, error) {
		button := rl.MouseButtonLeft
		if len(args) >= 1 {
			switch toInt32(args[0]) {
			case 1:
				button = rl.MouseButtonRight
			case 2:
				button = rl.MouseButtonMiddle
			default:
				button = rl.MouseButtonLeft
			}
		}
		return rl.IsMouseButtonPressed(button), nil
	})
	v.RegisterForeign("MouseReleased", func(args []interface{}) (interface{}, error) {
		button := rl.MouseButtonLeft
		if len(args) >= 1 {
			switch toInt32(args[0]) {
			case 1:
				button = rl.MouseButtonRight
			case 2:
				button = rl.MouseButtonMiddle
			default:
				button = rl.MouseButtonLeft
			}
		}
		return rl.IsMouseButtonReleased(button), nil
	})
	v.RegisterForeign("GamepadConnected", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("GamepadConnected requires (id)")
		}
		return rl.IsGamepadAvailable(int32(toInt32(args[0]))), nil
	})
	v.RegisterForeign("GetGamepadAxis", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("GetGamepadAxis requires (id, axis)")
		}
		return float64(rl.GetGamepadAxisMovement(int32(toInt32(args[0])), int32(toInt32(args[1])))), nil
	})
	v.RegisterForeign("GamepadButtonDown", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("GamepadButtonDown requires (id, button)")
		}
		return rl.IsGamepadButtonDown(int32(toInt32(args[0])), int32(toInt32(args[1]))), nil
	})
	v.RegisterForeign("GetVector2X", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return float64(0), nil
		}
		if s, ok := args[0].([]interface{}); ok && len(s) >= 1 {
			if f, ok := toFloat64Safe(s[0]); ok {
				return f, nil
			}
		}
		return float64(0), nil
	})
	v.RegisterForeign("GetVector2Y", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return float64(0), nil
		}
		if s, ok := args[0].([]interface{}); ok && len(s) >= 2 {
			if f, ok := toFloat64Safe(s[1]); ok {
				return f, nil
			}
		}
		return float64(0), nil
	})
	v.RegisterForeign("GetVector3Z", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return float64(0), nil
		}
		if s, ok := args[0].([]interface{}); ok && len(s) >= 3 {
			if f, ok := toFloat64Safe(s[2]); ok {
				return f, nil
			}
		}
		return float64(0), nil
	})
	v.RegisterForeign("IsMouseButtonUp", func(args []interface{}) (interface{}, error) {
		button := rl.MouseButtonLeft
		if len(args) >= 1 {
			switch toInt32(args[0]) {
			case 1:
				button = rl.MouseButtonRight
			case 2:
				button = rl.MouseButtonMiddle
			default:
				button = rl.MouseButtonLeft
			}
		}
		return rl.IsMouseButtonUp(button), nil
	})
	v.RegisterForeign("IsKeyPressedRepeat", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("IsKeyPressedRepeat requires (key)")
		}
		return rl.IsKeyPressedRepeat(int32(toInt32(args[0]))), nil
	})
	v.RegisterForeign("GetCharPressed", func(args []interface{}) (interface{}, error) {
		return int(rl.GetCharPressed()), nil
	})
	v.RegisterForeign("SetMouseCursor", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("SetMouseCursor requires (cursor)")
		}
		rl.SetMouseCursor(toInt32(args[0]))
		return nil, nil
	})
	v.RegisterForeign("IsGamepadAvailable", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("IsGamepadAvailable requires (gamepad)")
		}
		return rl.IsGamepadAvailable(int32(toInt32(args[0]))), nil
	})
	v.RegisterForeign("GetGamepadName", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("GetGamepadName requires (gamepad)")
		}
		return rl.GetGamepadName(int32(toInt32(args[0]))), nil
	})
	v.RegisterForeign("IsGamepadButtonPressed", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("IsGamepadButtonPressed requires (gamepad, button)")
		}
		return rl.IsGamepadButtonPressed(int32(toInt32(args[0])), int32(toInt32(args[1]))), nil
	})
	v.RegisterForeign("IsGamepadButtonDown", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("IsGamepadButtonDown requires (gamepad, button)")
		}
		return rl.IsGamepadButtonDown(int32(toInt32(args[0])), int32(toInt32(args[1]))), nil
	})
	v.RegisterForeign("IsGamepadButtonReleased", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("IsGamepadButtonReleased requires (gamepad, button)")
		}
		return rl.IsGamepadButtonReleased(int32(toInt32(args[0])), int32(toInt32(args[1]))), nil
	})
	v.RegisterForeign("GetGamepadAxisMovement", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("GetGamepadAxisMovement requires (gamepad, axis)")
		}
		return float64(rl.GetGamepadAxisMovement(int32(toInt32(args[0])), int32(toInt32(args[1])))), nil
	})
	v.RegisterForeign("IsGamepadButtonUp", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("IsGamepadButtonUp requires (gamepad, button)")
		}
		return rl.IsGamepadButtonUp(int32(toInt32(args[0])), int32(toInt32(args[1]))), nil
	})
	v.RegisterForeign("GetGamepadButtonPressed", func(args []interface{}) (interface{}, error) {
		return int(rl.GetGamepadButtonPressed()), nil
	})
	v.RegisterForeign("GetGamepadAxisCount", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("GetGamepadAxisCount requires (gamepad)")
		}
		return int(rl.GetGamepadAxisCount(int32(toInt32(args[0])))), nil
	})
	v.RegisterForeign("SetGamepadMappings", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("SetGamepadMappings requires (mappings)")
		}
		return int(rl.SetGamepadMappings(toString(args[0]))), nil
	})
	v.RegisterForeign("SetGamepadVibration", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("SetGamepadVibration requires (gamepad, leftMotor, rightMotor, duration)")
		}
		rl.SetGamepadVibration(int32(toInt32(args[0])), toFloat32(args[1]), toFloat32(args[2]), toFloat32(args[3]))
		return nil, nil
	})
	v.RegisterForeign("GetTouchPointCount", func(args []interface{}) (interface{}, error) {
		return int(rl.GetTouchPointCount()), nil
	})
	v.RegisterForeign("GetTouchX", func(args []interface{}) (interface{}, error) {
		return int(rl.GetTouchX()), nil
	})
	v.RegisterForeign("GetTouchY", func(args []interface{}) (interface{}, error) {
		return int(rl.GetTouchY()), nil
	})
	v.RegisterForeign("GetTouchPosition", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("GetTouchPosition requires (index)")
		}
		pos := rl.GetTouchPosition(toInt32(args[0]))
		return []interface{}{float64(pos.X), float64(pos.Y)}, nil
	})
	v.RegisterForeign("GetTouchPointId", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("GetTouchPointId requires (index)")
		}
		return int(rl.GetTouchPointId(toInt32(args[0]))), nil
	})
	v.RegisterForeign("GetMouseWheelMoveV", func(args []interface{}) (interface{}, error) {
		v := rl.GetMouseWheelMoveV()
		return []interface{}{float64(v.X), float64(v.Y)}, nil
	})

	// Key constants (0-arg, return key code for IsKeyDown/IsKeyPressed etc. — CharmingBlaze style, unprefixed)
	v.RegisterForeign("KEY_NULL", func(args []interface{}) (interface{}, error) { return int(rl.KeyNull), nil })
	v.RegisterForeign("KEY_APOSTROPHE", func(args []interface{}) (interface{}, error) { return int(rl.KeyApostrophe), nil })
	v.RegisterForeign("KEY_COMMA", func(args []interface{}) (interface{}, error) { return int(rl.KeyComma), nil })
	v.RegisterForeign("KEY_MINUS", func(args []interface{}) (interface{}, error) { return int(rl.KeyMinus), nil })
	v.RegisterForeign("KEY_PERIOD", func(args []interface{}) (interface{}, error) { return int(rl.KeyPeriod), nil })
	v.RegisterForeign("KEY_SLASH", func(args []interface{}) (interface{}, error) { return int(rl.KeySlash), nil })
	v.RegisterForeign("KEY_ZERO", func(args []interface{}) (interface{}, error) { return int(rl.KeyZero), nil })
	v.RegisterForeign("KEY_ONE", func(args []interface{}) (interface{}, error) { return int(rl.KeyOne), nil })
	v.RegisterForeign("KEY_TWO", func(args []interface{}) (interface{}, error) { return int(rl.KeyTwo), nil })
	v.RegisterForeign("KEY_THREE", func(args []interface{}) (interface{}, error) { return int(rl.KeyThree), nil })
	v.RegisterForeign("KEY_FOUR", func(args []interface{}) (interface{}, error) { return int(rl.KeyFour), nil })
	v.RegisterForeign("KEY_FIVE", func(args []interface{}) (interface{}, error) { return int(rl.KeyFive), nil })
	v.RegisterForeign("KEY_SIX", func(args []interface{}) (interface{}, error) { return int(rl.KeySix), nil })
	v.RegisterForeign("KEY_SEVEN", func(args []interface{}) (interface{}, error) { return int(rl.KeySeven), nil })
	v.RegisterForeign("KEY_EIGHT", func(args []interface{}) (interface{}, error) { return int(rl.KeyEight), nil })
	v.RegisterForeign("KEY_NINE", func(args []interface{}) (interface{}, error) { return int(rl.KeyNine), nil })
	v.RegisterForeign("KEY_SEMICOLON", func(args []interface{}) (interface{}, error) { return int(rl.KeySemicolon), nil })
	v.RegisterForeign("KEY_EQUAL", func(args []interface{}) (interface{}, error) { return int(rl.KeyEqual), nil })
	v.RegisterForeign("KEY_A", func(args []interface{}) (interface{}, error) { return int(rl.KeyA), nil })
	v.RegisterForeign("KEY_B", func(args []interface{}) (interface{}, error) { return int(rl.KeyB), nil })
	v.RegisterForeign("KEY_C", func(args []interface{}) (interface{}, error) { return int(rl.KeyC), nil })
	v.RegisterForeign("KEY_D", func(args []interface{}) (interface{}, error) { return int(rl.KeyD), nil })
	v.RegisterForeign("KEY_E", func(args []interface{}) (interface{}, error) { return int(rl.KeyE), nil })
	v.RegisterForeign("KEY_F", func(args []interface{}) (interface{}, error) { return int(rl.KeyF), nil })
	v.RegisterForeign("KEY_G", func(args []interface{}) (interface{}, error) { return int(rl.KeyG), nil })
	v.RegisterForeign("KEY_H", func(args []interface{}) (interface{}, error) { return int(rl.KeyH), nil })
	v.RegisterForeign("KEY_I", func(args []interface{}) (interface{}, error) { return int(rl.KeyI), nil })
	v.RegisterForeign("KEY_J", func(args []interface{}) (interface{}, error) { return int(rl.KeyJ), nil })
	v.RegisterForeign("KEY_K", func(args []interface{}) (interface{}, error) { return int(rl.KeyK), nil })
	v.RegisterForeign("KEY_L", func(args []interface{}) (interface{}, error) { return int(rl.KeyL), nil })
	v.RegisterForeign("KEY_M", func(args []interface{}) (interface{}, error) { return int(rl.KeyM), nil })
	v.RegisterForeign("KEY_N", func(args []interface{}) (interface{}, error) { return int(rl.KeyN), nil })
	v.RegisterForeign("KEY_O", func(args []interface{}) (interface{}, error) { return int(rl.KeyO), nil })
	v.RegisterForeign("KEY_P", func(args []interface{}) (interface{}, error) { return int(rl.KeyP), nil })
	v.RegisterForeign("KEY_Q", func(args []interface{}) (interface{}, error) { return int(rl.KeyQ), nil })
	v.RegisterForeign("KEY_R", func(args []interface{}) (interface{}, error) { return int(rl.KeyR), nil })
	v.RegisterForeign("KEY_S", func(args []interface{}) (interface{}, error) { return int(rl.KeyS), nil })
	v.RegisterForeign("KEY_T", func(args []interface{}) (interface{}, error) { return int(rl.KeyT), nil })
	v.RegisterForeign("KEY_U", func(args []interface{}) (interface{}, error) { return int(rl.KeyU), nil })
	v.RegisterForeign("KEY_V", func(args []interface{}) (interface{}, error) { return int(rl.KeyV), nil })
	v.RegisterForeign("KEY_W", func(args []interface{}) (interface{}, error) { return int(rl.KeyW), nil })
	v.RegisterForeign("KEY_X", func(args []interface{}) (interface{}, error) { return int(rl.KeyX), nil })
	v.RegisterForeign("KEY_Y", func(args []interface{}) (interface{}, error) { return int(rl.KeyY), nil })
	v.RegisterForeign("KEY_Z", func(args []interface{}) (interface{}, error) { return int(rl.KeyZ), nil })
	v.RegisterForeign("KEY_LEFT_BRACKET", func(args []interface{}) (interface{}, error) { return int(rl.KeyLeftBracket), nil })
	v.RegisterForeign("KEY_BACKSLASH", func(args []interface{}) (interface{}, error) { return int(rl.KeyBackSlash), nil })
	v.RegisterForeign("KEY_RIGHT_BRACKET", func(args []interface{}) (interface{}, error) { return int(rl.KeyRightBracket), nil })
	v.RegisterForeign("KEY_GRAVE", func(args []interface{}) (interface{}, error) { return int(rl.KeyGrave), nil })
	v.RegisterForeign("KEY_SPACE", func(args []interface{}) (interface{}, error) { return int(rl.KeySpace), nil })
	v.RegisterForeign("KEY_ESCAPE", func(args []interface{}) (interface{}, error) { return int(rl.KeyEscape), nil })
	v.RegisterForeign("KEY_ENTER", func(args []interface{}) (interface{}, error) { return int(rl.KeyEnter), nil })
	v.RegisterForeign("KEY_TAB", func(args []interface{}) (interface{}, error) { return int(rl.KeyTab), nil })
	v.RegisterForeign("KEY_BACKSPACE", func(args []interface{}) (interface{}, error) { return int(rl.KeyBackspace), nil })
	v.RegisterForeign("KEY_INSERT", func(args []interface{}) (interface{}, error) { return int(rl.KeyInsert), nil })
	v.RegisterForeign("KEY_DELETE", func(args []interface{}) (interface{}, error) { return int(rl.KeyDelete), nil })
	v.RegisterForeign("KEY_RIGHT", func(args []interface{}) (interface{}, error) { return int(rl.KeyRight), nil })
	v.RegisterForeign("KEY_LEFT", func(args []interface{}) (interface{}, error) { return int(rl.KeyLeft), nil })
	v.RegisterForeign("KEY_DOWN", func(args []interface{}) (interface{}, error) { return int(rl.KeyDown), nil })
	v.RegisterForeign("KEY_UP", func(args []interface{}) (interface{}, error) { return int(rl.KeyUp), nil })
	v.RegisterForeign("KEY_PAGE_UP", func(args []interface{}) (interface{}, error) { return int(rl.KeyPageUp), nil })
	v.RegisterForeign("KEY_PAGE_DOWN", func(args []interface{}) (interface{}, error) { return int(rl.KeyPageDown), nil })
	v.RegisterForeign("KEY_HOME", func(args []interface{}) (interface{}, error) { return int(rl.KeyHome), nil })
	v.RegisterForeign("KEY_END", func(args []interface{}) (interface{}, error) { return int(rl.KeyEnd), nil })
	v.RegisterForeign("KEY_F1", func(args []interface{}) (interface{}, error) { return int(rl.KeyF1), nil })
	v.RegisterForeign("KEY_F2", func(args []interface{}) (interface{}, error) { return int(rl.KeyF2), nil })
	v.RegisterForeign("KEY_F3", func(args []interface{}) (interface{}, error) { return int(rl.KeyF3), nil })
	v.RegisterForeign("KEY_F4", func(args []interface{}) (interface{}, error) { return int(rl.KeyF4), nil })
	v.RegisterForeign("KEY_F5", func(args []interface{}) (interface{}, error) { return int(rl.KeyF5), nil })
	v.RegisterForeign("KEY_F6", func(args []interface{}) (interface{}, error) { return int(rl.KeyF6), nil })
	v.RegisterForeign("KEY_F7", func(args []interface{}) (interface{}, error) { return int(rl.KeyF7), nil })
	v.RegisterForeign("KEY_F8", func(args []interface{}) (interface{}, error) { return int(rl.KeyF8), nil })
	v.RegisterForeign("KEY_F9", func(args []interface{}) (interface{}, error) { return int(rl.KeyF9), nil })
	v.RegisterForeign("KEY_F10", func(args []interface{}) (interface{}, error) { return int(rl.KeyF10), nil })
	v.RegisterForeign("KEY_F11", func(args []interface{}) (interface{}, error) { return int(rl.KeyF11), nil })
	v.RegisterForeign("KEY_F12", func(args []interface{}) (interface{}, error) { return int(rl.KeyF12), nil })
}
