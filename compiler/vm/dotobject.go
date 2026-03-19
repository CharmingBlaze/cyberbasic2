package vm

// DotObject is implemented by VM-visible handles that support property and method dispatch (Phase 3+).
type DotObject interface {
	GetProp(path []string) (Value, error)
	SetProp(path []string, val Value) error
	CallMethod(name string, args []Value) (Value, error)
}
