// Package dbp: Asset pipeline - LoadAsset, UnloadAsset, AssetExists, PreloadAsset.
package dbp

import (
	"fmt"

	"cyberbasic/compiler/runtime/assets"
	"cyberbasic/compiler/vm"
)

func registerAssets(v *vm.VM) {
	v.RegisterForeign("LoadAsset", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("LoadAsset(path) requires 1 argument")
		}
		path := toString(args[0])
		_, err := assets.LoadAsset(path)
		if err != nil {
			return nil, err
		}
		return path, nil
	})
	v.RegisterForeign("UnloadAsset", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("UnloadAsset(path) requires 1 argument")
		}
		path := toString(args[0])
		assets.UnloadAsset(path)
		return nil, nil
	})
	v.RegisterForeign("AssetExists", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return 0, nil
		}
		path := toString(args[0])
		if assets.AssetExists(path) {
			return 1, nil
		}
		return 0, nil
	})
	v.RegisterForeign("PreloadAsset", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("PreloadAsset(path) requires 1 argument")
		}
		path := toString(args[0])
		_, err := assets.PreloadAsset(path)
		if err != nil {
			return nil, err
		}
		return path, nil
	})
}
