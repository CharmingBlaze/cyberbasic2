// Package httpdot exposes http.get (sync), http.get_async + http.await (goroutine-based; blocks on await).
package httpdot

import (
	"cyberbasic/compiler/vm"
	"fmt"
	"strings"
)

// RegisterHTTPDot registers global "http".
func RegisterHTTPDot(v *vm.VM) {
	v.SetGlobal("http", &httpModuleDot{v: v})
}

type httpModuleDot struct {
	v *vm.VM
}

func (h *httpModuleDot) GetProp([]string) (vm.Value, error) { return nil, nil }
func (h *httpModuleDot) SetProp([]string, vm.Value) error {
	return fmt.Errorf("http: namespace is not assignable")
}

func (h *httpModuleDot) CallMethod(name string, args []vm.Value) (vm.Value, error) {
	ia := make([]interface{}, len(args))
	for i := range args {
		ia[i] = args[i]
	}
	switch strings.ToLower(name) {
	case "get":
		return h.v.CallForeign("HttpGet", ia)
	case "get_async":
		return h.v.CallForeign("HttpGetAsync", ia)
	case "await":
		return h.v.CallForeign("HttpAwait", ia)
	default:
		return nil, fmt.Errorf("http: unknown method %q (get, get_async, await)", name)
	}
}
