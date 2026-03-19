// Package dotargs converts []vm.Value to []interface{} for CallForeign.
package dotargs

import "cyberbasic/compiler/vm"

// From converts VM stack values to foreign call arguments.
func From(a []vm.Value) []interface{} {
	out := make([]interface{}, len(a))
	for i := range a {
		out[i] = a[i]
	}
	return out
}
