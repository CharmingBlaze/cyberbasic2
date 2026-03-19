// Package assets provides ASSETS.LOAD / AssetsGet-style keyed resources.
package assets

import (
	"cyberbasic/compiler/errors"
	"cyberbasic/compiler/vm"
	"fmt"
	"strings"
	"sync"
)

var (
	mu    sync.RWMutex
	store = make(map[string]vm.Value)
)

// RegisterAssets registers AssetsSet, AssetsGet, AssetsUnload, AssetsUnloadAll, AssetsLoaded.
func RegisterAssets(v *vm.VM) {
	v.RegisterForeign("AssetsSet", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("AssetsSet requires (key$, value)")
		}
		key := strings.ToLower(strings.TrimSpace(fmt.Sprint(args[0])))
		mu.Lock()
		store[key] = args[1]
		mu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("AssetsGet", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("AssetsGet requires (key$)")
		}
		key := strings.ToLower(strings.TrimSpace(fmt.Sprint(args[0])))
		mu.RLock()
		val, ok := store[key]
		keys := make([]string, 0, len(store))
		for k := range store {
			keys = append(keys, k)
		}
		mu.RUnlock()
		if !ok {
			sug := errors.Nearest(key, keys, 2)
			msg := fmt.Sprintf("Asset %q not found", key)
			suggestion := "Check ASSETS.LOAD ran first."
			if sug != "" {
				suggestion = fmt.Sprintf("Did you mean %q? %s", sug, suggestion)
			}
			return nil, &errors.CyberError{Code: errors.ErrAssetNotFound, Message: msg, Suggestion: suggestion}
		}
		return val, nil
	})
	v.RegisterForeign("AssetsUnload", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, nil
		}
		key := strings.ToLower(strings.TrimSpace(fmt.Sprint(args[0])))
		mu.Lock()
		delete(store, key)
		mu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("AssetsUnloadAll", func(args []interface{}) (interface{}, error) {
		mu.Lock()
		store = make(map[string]vm.Value)
		mu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("AssetsLoaded", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return false, nil
		}
		key := strings.ToLower(strings.TrimSpace(fmt.Sprint(args[0])))
		mu.RLock()
		_, ok := store[key]
		mu.RUnlock()
		return ok, nil
	})

	v.SetGlobal("assets", &assetsModuleDot{v: v})
}

type assetsModuleDot struct {
	v *vm.VM
}

func (a *assetsModuleDot) GetProp([]string) (vm.Value, error) { return nil, nil }
func (a *assetsModuleDot) SetProp([]string, vm.Value) error {
	return fmt.Errorf("assets: namespace is not assignable")
}

func (a *assetsModuleDot) CallMethod(name string, args []vm.Value) (vm.Value, error) {
	ia := make([]interface{}, len(args))
	for i := range args {
		ia[i] = args[i]
	}
	switch strings.ToLower(name) {
	case "set":
		return a.v.CallForeign("AssetsSet", ia)
	case "get":
		return a.v.CallForeign("AssetsGet", ia)
	case "unload":
		return a.v.CallForeign("AssetsUnload", ia)
	case "unloadall":
		return a.v.CallForeign("AssetsUnloadAll", ia)
	case "loaded":
		return a.v.CallForeign("AssetsLoaded", ia)
	default:
		return nil, fmt.Errorf("unknown assets method %q", name)
	}
}
