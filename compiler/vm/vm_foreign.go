package vm

import (
	"fmt"
	"strings"
)

// RegisterRenderType registers a command name for the hybrid render queue (2D, 3D, or GUI).
func (vm *VM) RegisterRenderType(name string, typ RenderType) {
	if vm.renderCommandType == nil {
		vm.renderCommandType = make(map[string]RenderType)
	}
	vm.renderCommandType[strings.ToLower(name)] = typ
}

// PushRenderCommand appends a command to the appropriate render queue (used when insideDraw and OpCallForeign).
func (vm *VM) PushRenderCommand(name string, args []interface{}, typ RenderType) {
	argsCopy := make([]interface{}, len(args))
	copy(argsCopy, args)
	item := RenderQueueItem{Name: name, Args: argsCopy}
	switch typ {
	case Render2D:
		vm.renderQueue2D = append(vm.renderQueue2D, item)
	case Render3D:
		vm.renderQueue3D = append(vm.renderQueue3D, item)
	case RenderGUI:
		vm.renderQueueGUI = append(vm.renderQueueGUI, item)
	}
}

// ClearRenderQueues clears all render queues (called at start of each frame in hybrid loop).
func (vm *VM) ClearRenderQueues() {
	vm.renderQueue2D = vm.renderQueue2D[:0]
	vm.renderQueue3D = vm.renderQueue3D[:0]
	vm.renderQueueGUI = vm.renderQueueGUI[:0]
}

// GetRenderQueues returns the three queues for FlushRenderQueues (2D, 3D, GUI).
func (vm *VM) GetRenderQueues() (q2D, q3D, qGUI []RenderQueueItem) {
	return vm.renderQueue2D, vm.renderQueue3D, vm.renderQueueGUI
}

// CallForeign invokes a foreign function by name with the given args (used when flushing render queues).
func (vm *VM) CallForeign(name string, args []interface{}) (interface{}, error) {
	fn := vm.foreign[strings.ToLower(name)]
	if fn == nil {
		return nil, fmt.Errorf("unknown foreign function: %s", name)
	}
	return fn(args)
}

// SetInsideDraw sets whether we are inside the user's draw() call (so render commands are queued).
func (vm *VM) SetInsideDraw(b bool) {
	vm.insideDraw = b
}
