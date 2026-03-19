package vm

import "fmt"

// ErrBreakpoint is returned when execution hits a breakpoint.
type ErrBreakpoint struct {
	Line int
}

func (e *ErrBreakpoint) Error() string {
	return fmt.Sprintf("breakpoint at line %d", e.Line)
}

// SetBreakpoints sets the line numbers at which execution should stop.
func (vm *VM) SetBreakpoints(lines map[int]bool) {
	vm.breakpoints = lines
}

// AddBreakpoint adds a breakpoint at the given line.
func (vm *VM) AddBreakpoint(line int) {
	if vm.breakpoints == nil {
		vm.breakpoints = make(map[int]bool)
	}
	vm.breakpoints[line] = true
}

// RemoveBreakpoint removes the breakpoint at the given line.
func (vm *VM) RemoveBreakpoint(line int) {
	if vm.breakpoints != nil {
		delete(vm.breakpoints, line)
	}
}

// ClearBreakpoints removes all breakpoints.
func (vm *VM) ClearBreakpoints() {
	vm.breakpoints = nil
}

// SetDebugMode enables or disables breakpoint checking during Run.
func (vm *VM) SetDebugMode(enabled bool) {
	vm.debugMode = enabled
}

// CurrentLine returns the source line at the current IP (0 if unknown).
func (vm *VM) CurrentLine() int {
	if vm.chunk == nil || vm.ip < 0 || vm.ip >= len(vm.chunk.Code) {
		return 0
	}
	return vm.chunk.LineAt(vm.ip)
}

// CallDepth returns the current call stack depth (0 = top level).
func (vm *VM) CallDepth() int {
	return len(vm.callStack)
}

// StepOver executes until the next line in the current function (does not step into calls).
func (vm *VM) StepOver() error {
	startLine := vm.CurrentLine()
	startDepth := vm.CallDepth()
	for {
		if err := vm.Step(); err != nil {
			return err
		}
		if vm.chunk == nil || vm.ip >= len(vm.chunk.Code) {
			return nil
		}
		line := vm.CurrentLine()
		depth := vm.CallDepth()
		if depth < startDepth {
			return nil
		}
		if depth == startDepth && line != startLine && line > 0 {
			return nil
		}
	}
}

// StepOut executes until the current function returns.
func (vm *VM) StepOut() error {
	startDepth := vm.CallDepth()
	for {
		if err := vm.Step(); err != nil {
			return err
		}
		if vm.chunk == nil || vm.ip >= len(vm.chunk.Code) {
			return nil
		}
		if vm.CallDepth() < startDepth {
			return nil
		}
	}
}

// WatchValue returns the value of a variable by name (from stack slots / globals).
func (vm *VM) WatchValue(name string) (interface{}, bool) {
	if vm.chunk == nil {
		return nil, false
	}
	idx, ok := vm.chunk.GetVariable(name)
	if !ok {
		return nil, false
	}
	if idx >= len(vm.stack) {
		return nil, false
	}
	return vm.stack[idx], true
}
