package vm

import "strings"

// OpCode represents bytecode instructions
type OpCode int

const (
	// Stack operations
	OpPush OpCode = iota
	OpPop
	OpDup
	OpSwap

	// Variable operations
	OpLoadVar
	OpStoreVar
	OpLoadGlobal
	OpStoreGlobal

	// Literals
	OpLoadConst
	OpLoadString

	// Arithmetic operations
	OpAdd
	OpSub
	OpMul
	OpDiv
	OpMod
	OpPower  // ^ exponentiation
	OpIntDiv // \ integer division
	OpNeg

	// Comparison operations
	OpEqual
	OpNotEqual
	OpLess
	OpLessEqual
	OpGreater
	OpGreaterEqual

	// Logical operations
	OpAnd
	OpOr
	OpXor
	OpNot

	// Control flow
	OpJump
	OpJumpIfFalse
	OpJumpIfTrue
	OpCall
	OpReturn
	OpReturnVal
	OpPrint
	OpStr
	OpInitGraphics3D
	OpBegin3DMode
	OpEnd3DMode
	OpDrawModel3D
	OpDrawGrid3D
	OpDrawAxes3D
	OpCreatePhysicsWorld2D
	OpDestroyPhysicsWorld2D
	OpStepPhysics2D
	OpCreatePhysicsBody2D
	OpDestroyPhysicsBody2D
	OpSetPhysicsPosition2D
	OpGetPhysicsPosition2D
	OpSetPhysicsAngle2D
	OpGetPhysicsAngle2D
	OpSetPhysicsVelocity2D
	OpGetPhysicsVelocity2D
	OpApplyPhysicsForce2D
	OpApplyPhysicsImpulse2D
	OpSetPhysicsDensity2D
	OpSetPhysicsFriction2D
	OpSetPhysicsRestitution2D
	OpRayCast2D
	OpCheckCollision2D
	OpQueryAABB2D
	OpCreatePhysicsWorld3D
	OpDestroyPhysicsWorld3D
	OpStepPhysics3D
	OpCreatePhysicsBody3D
	OpDestroyPhysicsBody3D
	OpSetPhysicsPosition3D
	OpGetPhysicsPosition3D
	OpSetPhysicsRotation3D
	OpGetPhysicsRotation3D
	OpSetPhysicsVelocity3D
	OpGetPhysicsVelocity3D
	OpApplyPhysicsForce3D
	OpApplyPhysicsImpulse3D
	OpSetPhysicsMass3D
	OpCheckCollision3D
	OpQueryAABB3D

	// Game-specific operations
	OpLoadImage
	OpCreateSprite
	OpSetSpritePosition
	OpDrawSprite
	OpLoadModel
	OpCreateCamera
	OpSetCameraPosition
	OpDrawModel
	OpPlayMusic
	OpPlaySound
	OpLoadSound
	OpCreatePhysicsBody
	OpSetVelocity
	OpApplyForce
	OpRayCast3D

	// Window/game loop
	OpSync
	OpShouldClose

	// Runtime helpers: random, sleep/wait, int, timer
	OpRandom     // no args -> push float 0..1
	OpRandomN    // 1 arg (max) -> push float 0..max
	OpSleep      // 1 arg (milliseconds) -> block
	OpInt        // 1 arg -> push integer (truncate)
	OpTimer      // no args -> push float seconds since start/reset
	OpResetTimer // no args -> reset timer zero

	// Math: Sin, Cos, Tan, Sqrt, Abs (1 arg), Lerp (3 args: a, b, t)
	OpSin
	OpCos
	OpTan
	OpSqrt
	OpAbs
	OpLerp
	// Noise: 2D Perlin-style and Simplex-style (x, y) -> float
	OpNoise2D

	// More math: Floor, Ceil, Round, Min, Max, Clamp, Pow, Exp, Log, Log10, Atan2, Sign, Deg2Rad, Rad2Deg
	OpFloor
	OpCeil
	OpRound
	OpMin
	OpMax
	OpClamp
	OpPow
	OpExp
	OpLog
	OpLog10
	OpAtan2
	OpSign
	OpDeg2Rad
	OpRad2Deg

	// Distance and radius (game math)
	OpDistance2D // x1, y1, x2, y2 -> float
	OpDistance3D // x1, y1, z1, x2, y2, z2 -> float
	OpDistSq2D   // x1, y1, x2, y2 -> float (squared distance, no sqrt)
	OpDistSq3D   // x1, y1, z1, x2, y2, z2 -> float
	OpInRadius2D // x1, y1, x2, y2, radius -> boolean
	OpInRadius3D // x1, y1, z1, x2, y2, z2, radius -> boolean
	OpAngle2D    // x1, y1, x2, y2 -> radians from (x1,y1) toward (x2,y2)

	// Matrix: MatMul(resultName, aName, bName) - R = A*B, names as constant indices
	OpMatMul

	// String: Left(s, n), Right(s, n), Mid(s, start, n), Len(s)
	OpLeftStr
	OpRightStr
	OpMidStr
	OpLenStr

	// File: EOF(handle) -> boolean
	OpEOF

	// File I/O: OpenFile(path, mode), ReadLine(handle), WriteLine(handle, text), CloseFile(handle)
	OpOpenFile  // path, mode (0=read, 1=write, 2=append) -> handle
	OpReadLine  // handle -> string
	OpWriteLine // handle, text
	OpCloseFile // handle

	// Arrays (multidimensional)
	OpCreateArray // nDims, constIdx..., varIndex -> allocate and store
	OpLoadArray   // varIndex -> pop indices..., push element
	OpStoreArray  // varIndex -> pop value, pop indices..., store

	// Foreign API: call into Go libraries (raylib, etc.)
	OpCallForeign
	// User-defined Sub/Function call by name (nameConstIndex, argCount)
	OpCallUser
	// Event handler registration: eventTypeConstIndex, keyConstIndex, handlerOffset (2 bytes)
	OpRegisterEvent
	// Coroutines: StartCoroutine (2-byte target offset), Yield, WaitSeconds (1 arg)
	OpStartCoroutine
	OpYield
	OpWaitSeconds

	// Special
	OpQuit // exit program (like QUIT / END)
	OpHalt
)

// Value represents a value in the VM
type Value interface{}

// EnumMembers maps enum value name (lowercase) to integer value.
type EnumMembers map[string]int64

// Chunk represents a compiled bytecode chunk
type Chunk struct {
	Code      []byte
	Lines     []int // source line per byte in Code (same length as Code; 0 = unknown)
	Constants []Value
	Variables map[string]int
	VarDims   map[string][]int // array dimensions per variable (nil = scalar)
	Functions map[string]int   // user Sub/Function name (lowercase) -> code offset
	// Enums: enum name (lowercase) -> member name (lowercase) -> value; used by Enum.getValue/getName/hasValue at runtime
	Enums map[string]EnumMembers
	currentLine int // used by compiler when emitting; Write records this into Lines
}

// NewChunk creates a new bytecode chunk
func NewChunk() *Chunk {
	return &Chunk{
		Code:      make([]byte, 0),
		Lines:     make([]int, 0),
		Constants: make([]Value, 0),
		Variables: make(map[string]int),
		VarDims:   make(map[string][]int),
		Functions: make(map[string]int),
		Enums:     make(map[string]EnumMembers),
	}
}

// SetLine sets the source line for subsequently emitted bytes (used by compiler).
func (c *Chunk) SetLine(line int) {
	c.currentLine = line
}

// LineAt returns the source line for the given instruction pointer (0 if unknown or out of range).
func (c *Chunk) LineAt(ip int) int {
	if ip < 0 || ip >= len(c.Lines) {
		return 0
	}
	return c.Lines[ip]
}

// Write adds a byte to the chunk and records the current source line for this instruction.
func (c *Chunk) Write(b byte) {
	c.Code = append(c.Code, b)
	c.Lines = append(c.Lines, c.currentLine)
}

// WriteJumpOffset writes a 2-byte signed offset (int16 little-endian) for backward-compatible long jumps.
func (c *Chunk) WriteJumpOffset(offset int) {
	off := int16(offset)
	c.Code = append(c.Code, byte(off), byte(off>>8))
	c.Lines = append(c.Lines, c.currentLine, c.currentLine)
}

// PatchJumpOffset patches the last 2 bytes at the given position with the given offset (int16).
func (c *Chunk) PatchJumpOffset(at int, offset int) {
	off := int16(offset)
	if at+2 > len(c.Code) {
		return
	}
	c.Code[at] = byte(off)
	c.Code[at+1] = byte(off >> 8)
}

// WriteConstant adds a constant to the chunk and returns its index
func (c *Chunk) WriteConstant(value Value) int {
	c.Constants = append(c.Constants, value)
	return len(c.Constants) - 1
}

// AddVariable adds a variable and returns its index (name normalized to lowercase)
func (c *Chunk) AddVariable(name string) int {
	key := strings.ToLower(name)
	if idx, exists := c.Variables[key]; exists {
		return idx
	}
	idx := len(c.Variables)
	c.Variables[key] = idx
	return idx
}

// GetVariable returns the index of a variable (name normalized to lowercase)
func (c *Chunk) GetVariable(name string) (int, bool) {
	idx, exists := c.Variables[strings.ToLower(name)]
	return idx, exists
}

// SetVarDims records array dimensions for a variable (used by compiler and VM).
func (c *Chunk) SetVarDims(name string, dims []int) {
	c.VarDims[strings.ToLower(name)] = dims
}

// GetVarDims returns dimensions for a variable (nil = scalar).
func (c *Chunk) GetVarDims(name string) ([]int, bool) {
	dims, ok := c.VarDims[strings.ToLower(name)]
	return dims, ok
}

// GetFunction returns the code offset for a user function/sub (name lowercase). Used by VM for OpCallUser.
func (c *Chunk) GetFunction(name string) (int, bool) {
	off, ok := c.Functions[strings.ToLower(name)]
	return off, ok
}
