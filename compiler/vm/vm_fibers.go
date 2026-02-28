package vm

import (
	"fmt"
	"strings"
	"time"
)

// errWithStack appends a short stack trace (line numbers) to the error for debugging.
func (vm *VM) errWithStack(err error) error {
	if err == nil {
		return nil
	}
	trace := vm.StackTrace()
	if len(trace) == 0 {
		return err
	}
	var parts []string
	for _, f := range trace {
		if f.Line > 0 {
			parts = append(parts, fmt.Sprintf("line %d", f.Line))
		}
	}
	if len(parts) > 0 {
		return fmt.Errorf("%w\nstack: %s", err, strings.Join(parts, "; "))
	}
	return err
}

// wakeSleeping moves any sleeping fibers whose resumeAt <= now back onto the run queue.
func (vm *VM) wakeSleeping() {
	now := time.Now()
	stillSleeping := vm.sleeping[:0]
	for _, e := range vm.sleeping {
		if !e.resumeAt.After(now) {
			vm.fiberQueue = append(vm.fiberQueue, e.fiberIndex)
		} else {
			stillSleeping = append(stillSleeping, e)
		}
	}
	vm.sleeping = stillSleeping
}

// ProcessEvents invokes registered On KeyDown/KeyPressed handlers when the runtime reports matching key state.
// Call after PollInputEvents in the game loop. Returns the first error from a handler if any.
func (vm *VM) ProcessEvents() error {
	if vm.runtime == nil || vm.chunk == nil {
		return nil
	}
	depth := len(vm.callStack)
	for _, h := range vm.eventHandlers {
		trigger := false
		switch h.eventType {
		case "keydown":
			trigger = vm.runtime.IsKeyDown(h.key)
		case "keypressed":
			trigger = vm.runtime.IsKeyPressed(h.key)
		}
		if !trigger {
			continue
		}
		vm.callStack = append(vm.callStack, vm.ip)
		vm.ip = h.handlerIP
		for len(vm.callStack) > depth && vm.running && vm.ip < len(vm.chunk.Code) {
			if err := vm.Step(); err != nil {
				return err
			}
		}
	}
	return nil
}
