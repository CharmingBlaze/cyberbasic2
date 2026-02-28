package vm

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

// ForeignFunc is a function callable from BASIC via the foreign API (e.g. RL.InitWindow).
// Args are passed in order; return value is pushed onto the VM stack (nil = no return).
type ForeignFunc func(args []interface{}) (interface{}, error)

// VM represents the virtual machine that executes bytecode
type VM struct {
	chunk          *Chunk
	ip             int // instruction pointer
	stack          []Value
	callStack      []int   // return addresses for user Sub/Function calls
	globals        map[string]Value
	running        bool
	runtime        GameRuntime // optional: when set, game opcodes call runtime instead of no-op
	foreign        map[string]ForeignFunc
	// Entity getter/setter: when set, entityName.prop read/write can be intercepted (e.g. for physics).
	entityGetters  map[string]func(entityName, prop string) (Value, bool) // key: prop name (lowercase) or "entity.prop"
	entitySetters  map[string]func(entityName, prop string, v Value)
	timerZero      time.Time
	fileHandles    map[int]*os.File
	fileReaders    map[int]*bufio.Reader // for ReadLine
	nextFileHandle int
	eventHandlers       []eventHandler
	collisionHandlers   map[string]string // bodyId -> subName for 2D collision callbacks
	fibers              []fiberState
	fiberQueue          []int
	currentFiber        int
	sleeping            []sleepEntry // fibers waiting for resume time (non-blocking WaitSeconds)

	// Hybrid update/draw: when inside draw(), render commands are queued instead of executed.
	insideDraw         bool
	drawFrameStack     []bool // parallel to callStack: true if that frame is draw()
	renderCommandType  map[string]RenderType
	renderQueue2D    []RenderQueueItem
	renderQueue3D    []RenderQueueItem
	renderQueueGUI   []RenderQueueItem
}

// RenderType classifies a foreign command for the hybrid render queue (2D, 3D, or GUI).
type RenderType int

const (
	RenderNone RenderType = iota
	Render2D
	Render3D
	RenderGUI
)

// RenderQueueItem is one entry in a render queue (name + args for deferred foreign call).
type RenderQueueItem struct {
	Name string
	Args []interface{}
}

type sleepEntry struct {
	fiberIndex int
	resumeAt   time.Time
}

type eventHandler struct {
	eventType string
	key       string
	handlerIP int
}

type fiberState struct {
	ip            int
	stack         []Value
	callStack     []int
	drawFrameStack []bool
}

// NewVM creates a new virtual machine instance
func NewVM() *VM {
	return &VM{
		stack:   make([]Value, 0),
		globals: make(map[string]Value),
		running: false,
		foreign: make(map[string]ForeignFunc),
	}
}

// Chunk returns the currently loaded chunk (nil if none). Used by runtime.StepFrame to detect update/draw.
func (vm *VM) Chunk() *Chunk {
	return vm.chunk
}

// Globals returns the global variables map (for tests and inspection).
func (vm *VM) Globals() map[string]Value {
	return vm.globals
}

// RegisterEntityGetter registers a getter for entity property reads. Key is prop name (lowercase) or "entityname.prop".
// When entityName.prop is loaded, if a getter is registered for that prop (or entity.prop), it is called; otherwise globals[entity] map is used.
func (vm *VM) RegisterEntityGetter(key string, fn func(entityName, prop string) (Value, bool)) {
	if vm.entityGetters == nil {
		vm.entityGetters = make(map[string]func(entityName, prop string) (Value, bool))
	}
	vm.entityGetters[strings.ToLower(key)] = fn
}

// RegisterEntitySetter registers a setter for entity property writes. Key is prop name (lowercase) or "entityname.prop".
func (vm *VM) RegisterEntitySetter(key string, fn func(entityName, prop string, v Value)) {
	if vm.entitySetters == nil {
		vm.entitySetters = make(map[string]func(entityName, prop string, v Value))
	}
	vm.entitySetters[strings.ToLower(key)] = fn
}

// LoadChunk loads a bytecode chunk into the VM
func (vm *VM) LoadChunk(chunk *Chunk) {
	vm.chunk = chunk
	vm.ip = 0
	vm.stack = make([]Value, 0)
	vm.callStack = vm.callStack[:0]
	vm.eventHandlers = vm.eventHandlers[:0]
	vm.collisionHandlers = make(map[string]string)
	vm.fibers = []fiberState{{ip: 0, stack: []Value{}, callStack: []int{}}}
	vm.fiberQueue = []int{0}
	vm.currentFiber = 0
	vm.sleeping = vm.sleeping[:0]
	vm.timerZero = time.Now()
	vm.fileHandles = make(map[int]*os.File)
	vm.fileReaders = make(map[int]*bufio.Reader)
	vm.nextFileHandle = 1
	vm.insideDraw = false
	vm.drawFrameStack = vm.drawFrameStack[:0]
	vm.renderQueue2D = nil
	vm.renderQueue3D = nil
	vm.renderQueueGUI = nil
	if vm.renderCommandType == nil {
		vm.renderCommandType = make(map[string]RenderType)
	}
}

// SetRuntime sets the game runtime used by game opcodes. If nil, game opcodes no-op (or debug print).
func (vm *VM) SetRuntime(r GameRuntime) {
	vm.runtime = r
}

// GetRuntime returns the game runtime (may be nil). Used by foreign APIs such as GAME.SyncSpriteToBody2D.
func (vm *VM) GetRuntime() GameRuntime {
	return vm.runtime
}

// RegisterCollisionHandler registers a Sub to call when bodyId has a collision (2D). Used by GAME.SetCollisionHandler.
func (vm *VM) RegisterCollisionHandler(bodyId, subName string) {
	if vm.collisionHandlers == nil {
		vm.collisionHandlers = make(map[string]string)
	}
	vm.collisionHandlers[strings.ToLower(bodyId)] = strings.ToLower(subName)
}

// GetCollisionHandlers returns a copy of bodyId -> subName for collision callbacks.
func (vm *VM) GetCollisionHandlers() map[string]string {
	out := make(map[string]string)
	for k, v := range vm.collisionHandlers {
		out[k] = v
	}
	return out
}

// InvokeSub calls a BASIC Sub by name with the given arguments. Sub sees them as first, second, ... param (stack[0]=first). Returns when the Sub returns.
func (vm *VM) InvokeSub(name string, args []interface{}) error {
	if vm.chunk == nil {
		return nil
	}
	subIP, ok := vm.chunk.GetFunction(strings.ToLower(name))
	if !ok {
		return nil
	}
	savedIP := vm.ip
	// Match OpCallUser: stack becomes [args[0], args[1], ...] for LoadVar 0, 1, ...
	argVals := make([]Value, len(args))
	for i, a := range args {
		argVals[i] = a
	}
	vm.stack = append(vm.stack[:0], argVals...)
	returnAddr := len(vm.chunk.Code)
	vm.callStack = append(vm.callStack, vm.ip)
	isDraw := strings.ToLower(name) == "draw"
	if isDraw {
		vm.drawFrameStack = append(vm.drawFrameStack, true)
		vm.insideDraw = true
	}
	vm.ip = subIP
	for vm.ip < len(vm.chunk.Code) {
		if err := vm.Step(); err != nil {
			if isDraw {
				vm.drawFrameStack = vm.drawFrameStack[:len(vm.drawFrameStack)-1]
				vm.insideDraw = false
				for _, b := range vm.drawFrameStack {
					if b {
						vm.insideDraw = true
						break
					}
				}
			}
			return err
		}
		if vm.ip == returnAddr {
			break
		}
	}
	if isDraw {
		vm.drawFrameStack = vm.drawFrameStack[:len(vm.drawFrameStack)-1]
		vm.insideDraw = false
		for _, b := range vm.drawFrameStack {
			if b {
				vm.insideDraw = true
				break
			}
		}
	}
	vm.ip = savedIP
	return nil
}

// SetForeignRegistry sets the map of foreign API functions (e.g. "RL.InitWindow" -> wrapper).
// Names are case-insensitive (canonical form: lowercase).
func (vm *VM) SetForeignRegistry(registry map[string]ForeignFunc) {
	vm.foreign = make(map[string]ForeignFunc)
	for k, v := range registry {
		vm.foreign[strings.ToLower(k)] = v
	}
}

// RegisterForeign adds one foreign function (e.g. "RL.InitWindow", wrapper).
func (vm *VM) RegisterForeign(name string, fn ForeignFunc) {
	if vm.foreign == nil {
		vm.foreign = make(map[string]ForeignFunc)
	}
	vm.foreign[strings.ToLower(name)] = fn
}

// Run executes the loaded bytecode
func (vm *VM) Run() error {
	if vm.chunk == nil {
		return fmt.Errorf("no chunk loaded")
	}

	vm.running = true

	for vm.running {
		vm.wakeSleeping()
		if len(vm.fiberQueue) == 0 {
			if len(vm.sleeping) == 0 {
				break
			}
			nextWake := vm.sleeping[0].resumeAt
			for _, e := range vm.sleeping[1:] {
				if e.resumeAt.Before(nextWake) {
					nextWake = e.resumeAt
				}
			}
			if d := time.Until(nextWake); d > 0 {
				time.Sleep(d)
			}
			continue
		}
		if vm.ip >= len(vm.chunk.Code) {
			break
		}
		if err := vm.Step(); err != nil {
			return vm.errWithStack(err)
		}
	}

	return nil
}

// Step executes one instruction and returns (for event handler invocation)
func (vm *VM) Step() error {
	if vm.chunk == nil || vm.ip >= len(vm.chunk.Code) {
		return nil
	}
	instruction := vm.chunk.Code[vm.ip]
	vm.ip++
	err := vm.executeInstruction(instruction)
	if err != nil && vm.chunk != nil {
		line := vm.chunk.LineAt(vm.ip - 1)
		if line > 0 {
			err = fmt.Errorf("line %d: %w", line, err)
		}
	}
	return err
}

// StackFrame is one frame in a stack trace (IP and source line).
type StackFrame struct {
	IP   int
	Line int
}

// StackTrace returns the current call stack for debugging (current IP first, then return addresses).
func (vm *VM) StackTrace() []StackFrame {
	if vm.chunk == nil {
		return nil
	}
	var frames []StackFrame
	if vm.ip >= 0 && vm.ip <= len(vm.chunk.Code) {
		frames = append(frames, StackFrame{IP: vm.ip, Line: vm.chunk.LineAt(vm.ip)})
	}
	for i := len(vm.callStack) - 1; i >= 0; i-- {
		ip := vm.callStack[i]
		frames = append(frames, StackFrame{IP: ip, Line: vm.chunk.LineAt(ip)})
	}
	return frames
}

