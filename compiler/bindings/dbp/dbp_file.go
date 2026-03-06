// Package dbp - File I/O: SaveString, LoadString, SaveValue, LoadValue.
package dbp

import (
	"fmt"
	"os"
	"strconv"

	"cyberbasic/compiler/vm"
)

func registerFile(v *vm.VM) {
	v.RegisterForeign("SaveString", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("SaveString(path, text) requires 2 arguments")
		}
		path := toString(args[0])
		text := toString(args[1])
		err := os.WriteFile(path, []byte(text), 0644)
		return err == nil, err
	})
	v.RegisterForeign("LoadString", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("LoadString(path) requires 1 argument")
		}
		data, err := os.ReadFile(toString(args[0]))
		if err != nil {
			return "", err
		}
		return string(data), nil
	})
	v.RegisterForeign("SaveValue", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("SaveValue(path, value) requires 2 arguments")
		}
		path := toString(args[0])
		val := args[1]
		var text string
		switch x := val.(type) {
		case int:
			text = strconv.Itoa(x)
		case float64:
			text = strconv.FormatFloat(x, 'g', -1, 64)
		case string:
			text = x
		default:
			text = fmt.Sprint(val)
		}
		err := os.WriteFile(path, []byte(text), 0644)
		return err == nil, err
	})
	v.RegisterForeign("LoadValue", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("LoadValue(path) requires 1 argument")
		}
		data, err := os.ReadFile(toString(args[0]))
		if err != nil {
			return nil, err
		}
		text := string(data)
		if f, err := strconv.ParseFloat(text, 64); err == nil {
			return f, nil
		}
		return text, nil
	})
}
