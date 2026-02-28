package valueutil

// IsTruthy returns false for nil, false, 0, 0.0, and ""; true otherwise.
// Used by both the VM (bytecode execution) and std bindings (e.g. IIf) so truthiness is consistent.
func IsTruthy(v interface{}) bool {
	if v == nil {
		return false
	}
	switch x := v.(type) {
	case bool:
		return x
	case int:
		return x != 0
	case float64:
		return x != 0
	case string:
		return x != ""
	default:
		return true
	}
}
