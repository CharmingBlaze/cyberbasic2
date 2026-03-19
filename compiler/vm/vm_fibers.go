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

// StopTaskByName removes all fibers with the given name from the run queue and sleeping list.
func (vm *VM) StopTaskByName(name string) {
	name = strings.ToLower(name)
	// Remove from fiberQueue
	newQueue := vm.fiberQueue[:0]
	for _, i := range vm.fiberQueue {
		if vm.fiberNames != nil && vm.fiberNames[i] != name {
			newQueue = append(newQueue, i)
		}
	}
	vm.fiberQueue = newQueue
	// Remove from sleeping
	stillSleeping := vm.sleeping[:0]
	for _, e := range vm.sleeping {
		if vm.fiberNames == nil || vm.fiberNames[e.fiberIndex] != name {
			stillSleeping = append(stillSleeping, e)
		}
	}
	vm.sleeping = stillSleeping
	// If current fiber was stopped, switch to next (handled by Run loop when queue empty)
}

// PauseTaskByName moves all fibers with the given name from the queue to sleeping with isPaused=true.
func (vm *VM) PauseTaskByName(name string) {
	name = strings.ToLower(name)
	// Remove matching fibers from queue, add to sleeping as paused
	newQueue := vm.fiberQueue[:0]
	for _, i := range vm.fiberQueue {
		if vm.fiberNames != nil && vm.fiberNames[i] == name {
			vm.sleeping = append(vm.sleeping, sleepEntry{fiberIndex: i, resumeAt: time.Time{}, isPaused: true})
		} else {
			newQueue = append(newQueue, i)
		}
	}
	vm.fiberQueue = newQueue
}

// ResumeTaskByName moves all paused fibers with the given name from sleeping back to the run queue.
func (vm *VM) ResumeTaskByName(name string) {
	name = strings.ToLower(name)
	stillSleeping := vm.sleeping[:0]
	for _, e := range vm.sleeping {
		if e.isPaused && vm.fiberNames != nil && vm.fiberNames[e.fiberIndex] == name {
			vm.fiberQueue = append(vm.fiberQueue, e.fiberIndex)
		} else {
			stillSleeping = append(stillSleeping, e)
		}
	}
	vm.sleeping = stillSleeping
}

// wakeSleeping moves any sleeping fibers whose resumeAt <= now back onto the run queue.
// Paused fibers (isPaused) are never auto-woken; only ResumeTask moves them.
func (vm *VM) wakeSleeping() {
	now := time.Now()
	stillSleeping := vm.sleeping[:0]
	for _, e := range vm.sleeping {
		if e.isPaused {
			stillSleeping = append(stillSleeping, e)
			continue
		}
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
