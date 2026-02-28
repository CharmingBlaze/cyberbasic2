package vm

import (
	"bufio"
	"fmt"
	"math"
	"math/rand"
	"os"
	"strings"
	"time"
)

// executeInstruction executes a single bytecode instruction
func (vm *VM) executeInstruction(instruction byte) error {
	op := OpCode(instruction)

	switch op {
	case OpPush:
		if vm.ip >= len(vm.chunk.Code) {
			return fmt.Errorf("unexpected end of code")
		}
		value := vm.chunk.Code[vm.ip]
		vm.ip++
		vm.push(value)

	case OpPop:
		vm.pop()

	case OpDup:
		if len(vm.stack) == 0 {
			return fmt.Errorf("stack underflow")
		}
		vm.push(vm.peek())

	case OpSwap:
		if len(vm.stack) < 2 {
			return fmt.Errorf("stack underflow")
		}
		a := vm.pop()
		b := vm.pop()
		vm.push(a)
		vm.push(b)

	case OpLoadConst:
		if vm.ip >= len(vm.chunk.Code) {
			return fmt.Errorf("unexpected end of code")
		}
		constIndex := int(vm.chunk.Code[vm.ip])
		vm.ip++

		if constIndex >= len(vm.chunk.Constants) {
			return fmt.Errorf("constant index out of bounds")
		}
		vm.push(vm.chunk.Constants[constIndex])

	case OpLoadString:
		if vm.ip >= len(vm.chunk.Code) {
			return fmt.Errorf("unexpected end of code")
		}
		constIndex := int(vm.chunk.Code[vm.ip])
		vm.ip++

		if constIndex >= len(vm.chunk.Constants) {
			return fmt.Errorf("constant index out of bounds")
		}
		vm.push(vm.chunk.Constants[constIndex])

	case OpLoadVar:
		if vm.ip >= len(vm.chunk.Code) {
			return fmt.Errorf("unexpected end of code")
		}
		varIndex := int(vm.chunk.Code[vm.ip])
		vm.ip++

		// Ensure stack has space for variables
		for varIndex >= len(vm.stack) {
			vm.push(nil) // Initialize with nil
		}

		vm.push(vm.stack[varIndex])

	case OpStoreVar:
		if vm.ip >= len(vm.chunk.Code) {
			return fmt.Errorf("unexpected end of code")
		}
		varIndex := int(vm.chunk.Code[vm.ip])
		vm.ip++

		var value Value
		if len(vm.stack) > 0 {
			value = vm.pop()
		} else {
			// Defensive: avoid crash when compiler emitted StoreVar without preceding value (e.g. first VAR in some paths)
			value = nil
		}

		// Ensure stack is large enough for variable storage (append so we don't overwrite existing slots)
		for len(vm.stack) <= varIndex {
			vm.stack = append(vm.stack, nil)
		}
		vm.stack[varIndex] = value

	case OpLoadGlobal:
		if vm.ip >= len(vm.chunk.Code) {
			return fmt.Errorf("unexpected end of code")
		}
		constIndex := int(vm.chunk.Code[vm.ip])
		vm.ip++

		if constIndex >= len(vm.chunk.Constants) {
			return fmt.Errorf("constant index out of bounds")
		}

		varName, ok := vm.chunk.Constants[constIndex].(string)
		if !ok {
			return fmt.Errorf("global name must be a string")
		}
		// Case-insensitive: canonical form is lowercase
		key := strings.ToLower(varName)

		value, exists := vm.globals[key]
		if !exists {
			// Allow 0-arg foreign functions as constants (e.g. KEY_W, KEY_A)
			if fn := vm.foreign[key]; fn != nil {
				result, err := fn(nil)
				if err != nil {
					return fmt.Errorf("global/constant %s: %w", varName, err)
				}
				vm.push(result)
				break
			}
			return fmt.Errorf("undefined global variable: %s", varName)
		}
		vm.push(value)

	case OpStoreGlobal:
		if vm.ip >= len(vm.chunk.Code) {
			return fmt.Errorf("unexpected end of code")
		}
		constIndex := int(vm.chunk.Code[vm.ip])
		vm.ip++

		if constIndex >= len(vm.chunk.Constants) {
			return fmt.Errorf("constant index out of bounds")
		}

		varName, ok := vm.chunk.Constants[constIndex].(string)
		if !ok {
			return fmt.Errorf("global name must be a string")
		}
		// Case-insensitive: canonical form is lowercase
		key := strings.ToLower(varName)

		if len(vm.stack) == 0 {
			return fmt.Errorf("stack underflow")
		}

		vm.globals[key] = vm.pop()

	case OpLoadEntityProp:
		if vm.ip+2 > len(vm.chunk.Code) {
			return fmt.Errorf("unexpected end of code for OpLoadEntityProp")
		}
		entityIdx := int(vm.chunk.Code[vm.ip])
		vm.ip++
		propIdx := int(vm.chunk.Code[vm.ip])
		vm.ip++
		if entityIdx >= len(vm.chunk.Constants) || propIdx >= len(vm.chunk.Constants) {
			return fmt.Errorf("constant index out of bounds for entity prop")
		}
		entityName, _ := vm.chunk.Constants[entityIdx].(string)
		propName, _ := vm.chunk.Constants[propIdx].(string)
		entityLower := strings.ToLower(entityName)
		propLower := strings.ToLower(propName)
		if vm.entityGetters != nil {
			if fn := vm.entityGetters[propLower]; fn != nil {
				if v, ok := fn(entityLower, propLower); ok {
					vm.push(v)
					break
				}
			}
			if fn := vm.entityGetters[entityLower+"."+propLower]; fn != nil {
				if v, ok := fn(entityLower, propLower); ok {
					vm.push(v)
					break
				}
			}
		}
		// Fallback: load from globals[entityName] map
		g, exists := vm.globals[entityLower]
		if !exists {
			return fmt.Errorf("entity not found: %s", entityName)
		}
		m, ok := g.(map[string]interface{})
		if !ok {
			return fmt.Errorf("entity %s is not a map", entityName)
		}
		if val, has := m[propName]; has {
			vm.push(val)
		} else if val, has := m[propLower]; has {
			vm.push(val)
		} else {
			vm.push(nil)
		}

	case OpStoreEntityProp:
		if vm.ip+2 > len(vm.chunk.Code) {
			return fmt.Errorf("unexpected end of code for OpStoreEntityProp")
		}
		if len(vm.stack) == 0 {
			return fmt.Errorf("stack underflow for OpStoreEntityProp")
		}
		value := vm.pop()
		entityIdx := int(vm.chunk.Code[vm.ip])
		vm.ip++
		propIdx := int(vm.chunk.Code[vm.ip])
		vm.ip++
		if entityIdx >= len(vm.chunk.Constants) || propIdx >= len(vm.chunk.Constants) {
			return fmt.Errorf("constant index out of bounds for entity prop")
		}
		entityName, _ := vm.chunk.Constants[entityIdx].(string)
		propName, _ := vm.chunk.Constants[propIdx].(string)
		entityLower := strings.ToLower(entityName)
		propLower := strings.ToLower(propName)
		if vm.entitySetters != nil {
			if fn := vm.entitySetters[propLower]; fn != nil {
				fn(entityLower, propLower, value)
				break
			}
			if fn := vm.entitySetters[entityLower+"."+propLower]; fn != nil {
				fn(entityLower, propLower, value)
				break
			}
		}
		// Fallback: set on globals[entityName] map
		g, exists := vm.globals[entityLower]
		if !exists {
			return fmt.Errorf("entity not found: %s", entityName)
		}
		m, ok := g.(map[string]interface{})
		if !ok {
			return fmt.Errorf("entity %s is not a map", entityName)
		}
		m[propName] = value

	case OpAdd:
		b := vm.pop()
		a := vm.pop()
		result, err := vm.add(a, b)
		if err != nil {
			return err
		}
		vm.push(result)

	case OpSub:
		b := vm.pop()
		a := vm.pop()
		result, err := vm.subtract(a, b)
		if err != nil {
			return err
		}
		vm.push(result)

	case OpMul:
		b := vm.pop()
		a := vm.pop()
		result, err := vm.multiply(a, b)
		if err != nil {
			return err
		}
		vm.push(result)

	case OpDiv:
		b := vm.pop()
		a := vm.pop()
		result, err := vm.divide(a, b)
		if err != nil {
			return err
		}
		vm.push(result)

	case OpMod:
		b := vm.pop()
		a := vm.pop()
		result, err := vm.modulo(a, b)
		if err != nil {
			return err
		}
		vm.push(result)

	case OpPower:
		b := vm.pop()
		a := vm.pop()
		result, err := vm.power(a, b)
		if err != nil {
			return err
		}
		vm.push(result)

	case OpIntDiv:
		b := vm.pop()
		a := vm.pop()
		result, err := vm.intDiv(a, b)
		if err != nil {
			return err
		}
		vm.push(result)

	case OpNeg:
		a := vm.pop()
		result, err := vm.negate(a)
		if err != nil {
			return err
		}
		vm.push(result)

	case OpEqual:
		b := vm.pop()
		a := vm.pop()
		vm.push(a == b)

	case OpNotEqual:
		b := vm.pop()
		a := vm.pop()
		vm.push(a != b)

	case OpLess:
		b := vm.pop()
		a := vm.pop()
		result, err := vm.less(a, b)
		if err != nil {
			return err
		}
		vm.push(result)

	case OpLessEqual:
		b := vm.pop()
		a := vm.pop()
		result, err := vm.lessEqual(a, b)
		if err != nil {
			return err
		}
		vm.push(result)

	case OpGreater:
		b := vm.pop()
		a := vm.pop()
		result, err := vm.greater(a, b)
		if err != nil {
			return err
		}
		vm.push(result)

	case OpGreaterEqual:
		b := vm.pop()
		a := vm.pop()
		result, err := vm.greaterEqual(a, b)
		if err != nil {
			return err
		}
		vm.push(result)

	case OpAnd:
		b := vm.pop()
		a := vm.pop()
		vm.push(vm.isTruthy(a) && vm.isTruthy(b))

	case OpOr:
		b := vm.pop()
		a := vm.pop()
		vm.push(vm.isTruthy(a) || vm.isTruthy(b))

	case OpXor:
		b := vm.pop()
		a := vm.pop()
		va, vb := vm.isTruthy(a), vm.isTruthy(b)
		vm.push(va != vb)

	case OpNot:
		a := vm.pop()
		vm.push(!vm.isTruthy(a))

	case OpJump:
		if vm.ip+2 > len(vm.chunk.Code) {
			return fmt.Errorf("unexpected end of code for OpJump")
		}
		// 2-byte signed offset (int16 little-endian) so backward jumps work for large loop bodies
		low := uint16(vm.chunk.Code[vm.ip])
		high := uint16(vm.chunk.Code[vm.ip+1])
		offset := int(int16(low | (high << 8)))
		vm.ip += 2
		vm.ip += offset

	case OpJumpIfFalse:
		if vm.ip+2 > len(vm.chunk.Code) {
			return fmt.Errorf("unexpected end of code for OpJumpIfFalse")
		}
		low := uint16(vm.chunk.Code[vm.ip])
		high := uint16(vm.chunk.Code[vm.ip+1])
		offset := int(int16(low | (high << 8)))
		vm.ip += 2

		if !vm.isTruthy(vm.pop()) {
			vm.ip += offset
		}

	case OpJumpIfTrue:
		if vm.ip+2 > len(vm.chunk.Code) {
			return fmt.Errorf("unexpected end of code for OpJumpIfTrue")
		}
		low := uint16(vm.chunk.Code[vm.ip])
		high := uint16(vm.chunk.Code[vm.ip+1])
		offset := int(int16(low | (high << 8)))
		vm.ip += 2

		if vm.isTruthy(vm.pop()) {
			vm.ip += offset
		}

	case OpCall:
		if vm.ip >= len(vm.chunk.Code) {
			return fmt.Errorf("unexpected end of code")
		}
		argCount := int(vm.chunk.Code[vm.ip])
		vm.ip++

		// Special handling for PRINT and STR
		if len(vm.stack) > 0 {
			// Check if this is a PRINT or STR call by looking at the previous instruction
			if vm.ip > 0 {
				prevInstr := vm.chunk.Code[vm.ip-2]
				if prevInstr == byte(OpPrint) || prevInstr == byte(OpStr) {
					// This is a PRINT or STR call, handle specially
					return nil
				}
			}
		}

		// For now, just pop the arguments
		for i := 0; i < argCount; i++ {
			vm.pop()
		}

	case OpCallUser:
		if vm.ip+2 > len(vm.chunk.Code) {
			return fmt.Errorf("unexpected end of code for OpCallUser")
		}
		nameConstIndex := int(vm.chunk.Code[vm.ip])
		vm.ip++
		argCount := int(vm.chunk.Code[vm.ip])
		vm.ip++
		if nameConstIndex < 0 || nameConstIndex >= len(vm.chunk.Constants) {
			return fmt.Errorf("invalid constant index for user call: %d", nameConstIndex)
		}
		nameVal := vm.chunk.Constants[nameConstIndex]
		name, ok := nameVal.(string)
		if !ok {
			return fmt.Errorf("user call name must be string constant, got %T", nameVal)
		}
		name = strings.ToLower(name)
		targetIP, ok := vm.chunk.GetFunction(name)
		if !ok {
			return fmt.Errorf("unknown user function: %s", name)
		}
		if len(vm.stack) < argCount {
			return fmt.Errorf("stack underflow for user call %s: need %d args, have %d", name, argCount, len(vm.stack))
		}
		// Pop args and replace stack with just those args so callee sees stack[0]=first, stack[1]=second, ...
		args := make([]Value, argCount)
		for i := argCount - 1; i >= 0; i-- {
			args[i] = vm.pop()
		}
		vm.callStack = append(vm.callStack, vm.ip)
		isDraw := (name == "draw")
		vm.drawFrameStack = append(vm.drawFrameStack, isDraw)
		if isDraw {
			vm.insideDraw = true
		}
		vm.stack = append(vm.stack[:0], args...)
		vm.ip = targetIP

	case OpReturn:
		if len(vm.callStack) == 0 {
			// Fiber or main ended: remove from queue if fiber, else halt
			if len(vm.fiberQueue) > 1 {
				// Remove current fiber from queue
				newQueue := make([]int, 0, len(vm.fiberQueue)-1)
				for _, i := range vm.fiberQueue {
					if i != vm.currentFiber {
						newQueue = append(newQueue, i)
					}
				}
				vm.fiberQueue = newQueue
				if len(vm.fiberQueue) == 0 {
					vm.running = false
					return nil
				}
				vm.currentFiber = vm.fiberQueue[0]
				next := &vm.fibers[vm.currentFiber]
				vm.ip = next.ip
				vm.stack = append(vm.stack[:0], next.stack...)
				vm.callStack = append(vm.callStack[:0], next.callStack...)
				vm.drawFrameStack = append(vm.drawFrameStack[:0], next.drawFrameStack...)
				vm.insideDraw = false
				for _, b := range vm.drawFrameStack {
					if b {
						vm.insideDraw = true
						break
					}
				}
			} else {
				vm.running = false
			}
			return nil
		}
		vm.ip = vm.callStack[len(vm.callStack)-1]
		vm.callStack = vm.callStack[:len(vm.callStack)-1]
		if len(vm.drawFrameStack) > 0 {
			wasDraw := vm.drawFrameStack[len(vm.drawFrameStack)-1]
			vm.drawFrameStack = vm.drawFrameStack[:len(vm.drawFrameStack)-1]
			if wasDraw {
				vm.insideDraw = false
				for _, b := range vm.drawFrameStack {
					if b {
						vm.insideDraw = true
						break
					}
				}
			}
		}

	case OpReturnVal:
		if len(vm.stack) == 0 {
			return fmt.Errorf("stack underflow for OpReturnVal")
		}
		val := vm.pop()
		if len(vm.callStack) == 0 {
			if len(vm.fiberQueue) > 1 {
				newQueue := make([]int, 0, len(vm.fiberQueue)-1)
				for _, i := range vm.fiberQueue {
					if i != vm.currentFiber {
						newQueue = append(newQueue, i)
					}
				}
				vm.fiberQueue = newQueue
				if len(vm.fiberQueue) == 0 {
					vm.running = false
					vm.stack = vm.stack[:0]
					vm.push(val)
					return nil
				}
				vm.currentFiber = vm.fiberQueue[0]
				next := &vm.fibers[vm.currentFiber]
				vm.ip = next.ip
				vm.stack = append(vm.stack[:0], next.stack...)
				vm.callStack = append(vm.callStack[:0], next.callStack...)
				vm.drawFrameStack = append(vm.drawFrameStack[:0], next.drawFrameStack...)
				vm.insideDraw = false
				for _, b := range vm.drawFrameStack {
					if b {
						vm.insideDraw = true
						break
					}
				}
				vm.push(val)
			} else {
				vm.running = false
				vm.stack = vm.stack[:0]
				vm.push(val)
			}
			return nil
		}
		vm.ip = vm.callStack[len(vm.callStack)-1]
		vm.callStack = vm.callStack[:len(vm.callStack)-1]
		if len(vm.drawFrameStack) > 0 {
			wasDraw := vm.drawFrameStack[len(vm.drawFrameStack)-1]
			vm.drawFrameStack = vm.drawFrameStack[:len(vm.drawFrameStack)-1]
			if wasDraw {
				vm.insideDraw = false
				for _, b := range vm.drawFrameStack {
					if b {
						vm.insideDraw = true
						break
					}
				}
			}
		}
		vm.stack = vm.stack[:0]
		vm.push(val)

	case OpRegisterEvent:
		if vm.ip+4 > len(vm.chunk.Code) {
			return fmt.Errorf("unexpected end of code for OpRegisterEvent")
		}
		eventTypeIdx := int(vm.chunk.Code[vm.ip])
		vm.ip++
		keyIdx := int(vm.chunk.Code[vm.ip])
		vm.ip++
		low := uint16(vm.chunk.Code[vm.ip])
		high := uint16(vm.chunk.Code[vm.ip+1])
		vm.ip += 2
		handlerIP := int(low | (high << 8))
		if eventTypeIdx >= 0 && eventTypeIdx < len(vm.chunk.Constants) {
			if kIdx := keyIdx; kIdx >= 0 && kIdx < len(vm.chunk.Constants) {
				eventType, _ := vm.chunk.Constants[eventTypeIdx].(string)
				key, _ := vm.chunk.Constants[kIdx].(string)
				vm.eventHandlers = append(vm.eventHandlers, eventHandler{eventType: eventType, key: key, handlerIP: handlerIP})
			}
		}

	case OpStartCoroutine:
		if vm.ip+2 > len(vm.chunk.Code) {
			return fmt.Errorf("unexpected end of code for OpStartCoroutine")
		}
		low := uint16(vm.chunk.Code[vm.ip])
		high := uint16(vm.chunk.Code[vm.ip+1])
		vm.ip += 2
		targetIP := int(low | (high << 8))
		vm.fibers = append(vm.fibers, fiberState{ip: targetIP, stack: []Value{}, callStack: []int{}, drawFrameStack: nil})
		vm.fiberQueue = append(vm.fiberQueue, len(vm.fibers)-1)

	case OpYield:
		// Save current state, rotate queue, load next fiber
		vm.fibers[vm.currentFiber] = fiberState{
			ip:            vm.ip,
			stack:         append([]Value(nil), vm.stack...),
			callStack:     append([]int(nil), vm.callStack...),
			drawFrameStack: append([]bool(nil), vm.drawFrameStack...),
		}
		if len(vm.fiberQueue) < 2 {
			break
		}
		vm.fiberQueue = append(vm.fiberQueue[1:], vm.fiberQueue[0])
		vm.currentFiber = vm.fiberQueue[0]
		next := &vm.fibers[vm.currentFiber]
		vm.ip = next.ip
		vm.stack = append(vm.stack[:0], next.stack...)
		vm.callStack = append(vm.callStack[:0], next.callStack...)
		vm.drawFrameStack = append(vm.drawFrameStack[:0], next.drawFrameStack...)
		vm.insideDraw = false
		for _, b := range vm.drawFrameStack {
			if b {
				vm.insideDraw = true
				break
			}
		}

	case OpWaitSeconds:
		if len(vm.stack) == 0 {
			return fmt.Errorf("stack underflow for WaitSeconds")
		}
		secVal := vm.pop()
		var sec float64
		switch v := secVal.(type) {
		case float64:
			sec = v
		case int:
			sec = float64(v)
		default:
			sec = 1
		}
		if sec <= 0 {
			break
		}
		// Non-blocking: save fiber state, remove from queue, add to sleeping, switch to next fiber
		vm.fibers[vm.currentFiber] = fiberState{
			ip:            vm.ip,
			stack:         append([]Value(nil), vm.stack...),
			callStack:     append([]int(nil), vm.callStack...),
			drawFrameStack: append([]bool(nil), vm.drawFrameStack...),
		}
		resumeAt := time.Now().Add(time.Duration(sec * float64(time.Second)))
		vm.sleeping = append(vm.sleeping, sleepEntry{fiberIndex: vm.currentFiber, resumeAt: resumeAt})
		newQueue := make([]int, 0, len(vm.fiberQueue)-1)
		for _, i := range vm.fiberQueue {
			if i != vm.currentFiber {
				newQueue = append(newQueue, i)
			}
		}
		vm.fiberQueue = newQueue
		if len(vm.fiberQueue) == 0 {
			break
		}
		vm.currentFiber = vm.fiberQueue[0]
		next := &vm.fibers[vm.currentFiber]
		vm.ip = next.ip
		vm.stack = append(vm.stack[:0], next.stack...)
		vm.callStack = append(vm.callStack[:0], next.callStack...)
		vm.drawFrameStack = append(vm.drawFrameStack[:0], next.drawFrameStack...)
		vm.insideDraw = false
		for _, b := range vm.drawFrameStack {
			if b {
				vm.insideDraw = true
				break
			}
		}

	case OpCallForeign:
		if vm.ip+2 > len(vm.chunk.Code) {
			return fmt.Errorf("unexpected end of code for OpCallForeign")
		}
		constIndex := int(vm.chunk.Code[vm.ip])
		vm.ip++
		argCount := int(vm.chunk.Code[vm.ip])
		vm.ip++
		if constIndex < 0 || constIndex >= len(vm.chunk.Constants) {
			return fmt.Errorf("invalid constant index for foreign call: %d", constIndex)
		}
		nameVal := vm.chunk.Constants[constIndex]
		name, ok := nameVal.(string)
		if !ok {
			return fmt.Errorf("foreign call name must be string constant, got %T", nameVal)
		}
		if len(vm.stack) < argCount {
			return fmt.Errorf("stack underflow for foreign call %s: need %d args, have %d", name, argCount, len(vm.stack))
		}
		args := make([]interface{}, argCount)
		for i := argCount - 1; i >= 0; i-- {
			args[i] = vm.pop()
		}
		// Hybrid draw: when inside draw(), queue render commands instead of executing.
		if vm.insideDraw && vm.renderCommandType != nil {
			if typ := vm.renderCommandType[strings.ToLower(name)]; typ != RenderNone {
				vm.PushRenderCommand(name, args, typ)
				break
			}
		}
		fn := vm.foreign[strings.ToLower(name)]
		if fn == nil {
			return fmt.Errorf("unknown foreign function: %s", name)
		}
		result, err := fn(args)
		if err != nil {
			return fmt.Errorf("foreign call %s: %w", name, err)
		}
		if result != nil {
			vm.push(result)
		}

	case OpPrint:
		// Pop argument and print it
		if len(vm.stack) == 0 {
			return fmt.Errorf("stack underflow for PRINT")
		}
		value := vm.pop()
		fmt.Printf("%v\n", value)

	case OpStr:
		// Pop argument and convert to string
		if len(vm.stack) == 0 {
			return fmt.Errorf("stack underflow for STR")
		}
		value := vm.pop()
		var strValue string
		switch v := value.(type) {
		case string:
			strValue = v
		case int:
			strValue = fmt.Sprintf("%d", v)
		case float64:
			strValue = fmt.Sprintf("%.6f", v)
		case bool:
			strValue = fmt.Sprintf("%t", v)
		default:
			strValue = fmt.Sprintf("%v", v)
		}
		vm.push(strValue)

	case OpInitGraphics3D, OpBegin3DMode, OpEnd3DMode, OpDrawModel3D, OpDrawGrid3D, OpDrawAxes3D,
		OpCreatePhysicsWorld2D, OpDestroyPhysicsWorld2D, OpStepPhysics2D, OpCreatePhysicsBody2D, OpDestroyPhysicsBody2D,
		OpSetPhysicsPosition2D, OpGetPhysicsPosition2D, OpSetPhysicsAngle2D, OpGetPhysicsAngle2D,
		OpSetPhysicsVelocity2D, OpGetPhysicsVelocity2D, OpApplyPhysicsForce2D, OpApplyPhysicsImpulse2D,
		OpSetPhysicsDensity2D, OpSetPhysicsFriction2D, OpSetPhysicsRestitution2D, OpRayCast2D, OpCheckCollision2D, OpQueryAABB2D,
		OpCreatePhysicsWorld3D, OpDestroyPhysicsWorld3D, OpStepPhysics3D, OpCreatePhysicsBody3D, OpDestroyPhysicsBody3D,
		OpSetPhysicsPosition3D, OpGetPhysicsPosition3D, OpSetPhysicsRotation3D, OpGetPhysicsRotation3D,
		OpSetPhysicsVelocity3D, OpGetPhysicsVelocity3D, OpApplyPhysicsForce3D, OpApplyPhysicsImpulse3D,
		OpSetPhysicsMass3D, OpCheckCollision3D, OpQueryAABB3D:
		return fmt.Errorf("deprecated opcode %d: use BOX2D.*/BULLET.* and raylib instead", op)

	case OpSync:
		if vm.runtime != nil {
			if err := vm.runtime.Sync(); err != nil {
				return err
			}
		}

	case OpShouldClose:
		if vm.runtime != nil {
			vm.push(vm.runtime.ShouldClose())
		} else {
			vm.push(false)
		}

	case OpRandom:
		vm.push(rand.Float64())

	case OpRandomN:
		if len(vm.stack) == 0 {
			return fmt.Errorf("stack underflow for Random(n)")
		}
		max := valueToFloat64(vm.pop())
		vm.push(rand.Float64() * max)

	case OpSleep:
		if len(vm.stack) == 0 {
			return fmt.Errorf("stack underflow for Sleep")
		}
		ms := valueToFloat64(vm.pop())
		time.Sleep(time.Duration(ms) * time.Millisecond)

	case OpInt:
		if len(vm.stack) == 0 {
			return fmt.Errorf("stack underflow for Int")
		}
		v := vm.pop()
		vm.push(int(valueToFloat64(v)))

	case OpTimer:
		vm.push(time.Since(vm.timerZero).Seconds())

	case OpResetTimer:
		vm.timerZero = time.Now()

	case OpSin:
		if len(vm.stack) == 0 {
			return fmt.Errorf("stack underflow for Sin")
		}
		x := valueToFloat64(vm.pop())
		vm.push(math.Sin(x))
	case OpCos:
		if len(vm.stack) == 0 {
			return fmt.Errorf("stack underflow for Cos")
		}
		x := valueToFloat64(vm.pop())
		vm.push(math.Cos(x))
	case OpTan:
		if len(vm.stack) == 0 {
			return fmt.Errorf("stack underflow for Tan")
		}
		x := valueToFloat64(vm.pop())
		vm.push(math.Tan(x))
	case OpSqrt:
		if len(vm.stack) == 0 {
			return fmt.Errorf("stack underflow for Sqrt")
		}
		x := valueToFloat64(vm.pop())
		vm.push(math.Sqrt(x))
	case OpAbs:
		if len(vm.stack) == 0 {
			return fmt.Errorf("stack underflow for Abs")
		}
		v := vm.pop()
		switch x := v.(type) {
		case float64:
			vm.push(math.Abs(x))
		case int:
			if x < 0 {
				vm.push(-x)
			} else {
				vm.push(x)
			}
		default:
			vm.push(math.Abs(valueToFloat64(v)))
		}
	case OpLerp:
		if len(vm.stack) < 3 {
			return fmt.Errorf("stack underflow for Lerp")
		}
		t := valueToFloat64(vm.pop())
		b := valueToFloat64(vm.pop())
		a := valueToFloat64(vm.pop())
		vm.push(a + (b-a)*t)
	case OpNoise2D:
		if len(vm.stack) < 2 {
			return fmt.Errorf("stack underflow for Noise2D")
		}
		y := valueToFloat64(vm.pop())
		x := valueToFloat64(vm.pop())
		vm.push(simplex2D(x, y))

	case OpFloor:
		if len(vm.stack) == 0 {
			return fmt.Errorf("stack underflow for Floor")
		}
		vm.push(math.Floor(valueToFloat64(vm.pop())))
	case OpCeil:
		if len(vm.stack) == 0 {
			return fmt.Errorf("stack underflow for Ceil")
		}
		vm.push(math.Ceil(valueToFloat64(vm.pop())))
	case OpRound:
		if len(vm.stack) == 0 {
			return fmt.Errorf("stack underflow for Round")
		}
		vm.push(math.Round(valueToFloat64(vm.pop())))
	case OpMin:
		if len(vm.stack) < 2 {
			return fmt.Errorf("stack underflow for Min")
		}
		b := valueToFloat64(vm.pop())
		a := valueToFloat64(vm.pop())
		vm.push(math.Min(a, b))
	case OpMax:
		if len(vm.stack) < 2 {
			return fmt.Errorf("stack underflow for Max")
		}
		b := valueToFloat64(vm.pop())
		a := valueToFloat64(vm.pop())
		vm.push(math.Max(a, b))
	case OpClamp:
		if len(vm.stack) < 3 {
			return fmt.Errorf("stack underflow for Clamp")
		}
		hi := valueToFloat64(vm.pop())
		lo := valueToFloat64(vm.pop())
		x := valueToFloat64(vm.pop())
		vm.push(math.Max(lo, math.Min(hi, x))) // clamp x to [lo, hi]
	case OpPow:
		if len(vm.stack) < 2 {
			return fmt.Errorf("stack underflow for Pow")
		}
		exp := valueToFloat64(vm.pop())
		base := valueToFloat64(vm.pop())
		vm.push(math.Pow(base, exp))
	case OpExp:
		if len(vm.stack) == 0 {
			return fmt.Errorf("stack underflow for Exp")
		}
		vm.push(math.Exp(valueToFloat64(vm.pop())))
	case OpLog:
		if len(vm.stack) == 0 {
			return fmt.Errorf("stack underflow for Log")
		}
		vm.push(math.Log(valueToFloat64(vm.pop())))
	case OpLog10:
		if len(vm.stack) == 0 {
			return fmt.Errorf("stack underflow for Log10")
		}
		vm.push(math.Log10(valueToFloat64(vm.pop())))
	case OpAtan2:
		if len(vm.stack) < 2 {
			return fmt.Errorf("stack underflow for Atan2")
		}
		x := valueToFloat64(vm.pop())
		y := valueToFloat64(vm.pop())
		vm.push(math.Atan2(y, x))
	case OpSign:
		if len(vm.stack) == 0 {
			return fmt.Errorf("stack underflow for Sign")
		}
		x := valueToFloat64(vm.pop())
		if x < 0 {
			vm.push(-1.0)
		} else if x > 0 {
			vm.push(1.0)
		} else {
			vm.push(0.0)
		}
	case OpDeg2Rad:
		if len(vm.stack) == 0 {
			return fmt.Errorf("stack underflow for Deg2Rad")
		}
		vm.push(valueToFloat64(vm.pop()) * math.Pi / 180)
	case OpRad2Deg:
		if len(vm.stack) == 0 {
			return fmt.Errorf("stack underflow for Rad2Deg")
		}
		vm.push(valueToFloat64(vm.pop()) * 180 / math.Pi)

	case OpDistance2D:
		if len(vm.stack) < 4 {
			return fmt.Errorf("stack underflow for Distance2D")
		}
		y2 := valueToFloat64(vm.pop())
		x2 := valueToFloat64(vm.pop())
		y1 := valueToFloat64(vm.pop())
		x1 := valueToFloat64(vm.pop())
		dx, dy := x2-x1, y2-y1
		vm.push(math.Sqrt(dx*dx + dy*dy))
	case OpDistance3D:
		if len(vm.stack) < 6 {
			return fmt.Errorf("stack underflow for Distance3D")
		}
		z2 := valueToFloat64(vm.pop())
		y2 := valueToFloat64(vm.pop())
		x2 := valueToFloat64(vm.pop())
		z1 := valueToFloat64(vm.pop())
		y1 := valueToFloat64(vm.pop())
		x1 := valueToFloat64(vm.pop())
		dx, dy, dz := x2-x1, y2-y1, z2-z1
		vm.push(math.Sqrt(dx*dx + dy*dy + dz*dz))
	case OpDistSq2D:
		if len(vm.stack) < 4 {
			return fmt.Errorf("stack underflow for DistSq2D")
		}
		y2 := valueToFloat64(vm.pop())
		x2 := valueToFloat64(vm.pop())
		y1 := valueToFloat64(vm.pop())
		x1 := valueToFloat64(vm.pop())
		dx, dy := x2-x1, y2-y1
		vm.push(dx*dx + dy*dy)
	case OpDistSq3D:
		if len(vm.stack) < 6 {
			return fmt.Errorf("stack underflow for DistSq3D")
		}
		z2 := valueToFloat64(vm.pop())
		y2 := valueToFloat64(vm.pop())
		x2 := valueToFloat64(vm.pop())
		z1 := valueToFloat64(vm.pop())
		y1 := valueToFloat64(vm.pop())
		x1 := valueToFloat64(vm.pop())
		dx, dy, dz := x2-x1, y2-y1, z2-z1
		vm.push(dx*dx + dy*dy + dz*dz)
	case OpInRadius2D:
		if len(vm.stack) < 5 {
			return fmt.Errorf("stack underflow for InRadius2D")
		}
		radius := valueToFloat64(vm.pop())
		y2 := valueToFloat64(vm.pop())
		x2 := valueToFloat64(vm.pop())
		y1 := valueToFloat64(vm.pop())
		x1 := valueToFloat64(vm.pop())
		dx, dy := x2-x1, y2-y1
		vm.push(dx*dx+dy*dy <= radius*radius)
	case OpInRadius3D:
		if len(vm.stack) < 7 {
			return fmt.Errorf("stack underflow for InRadius3D")
		}
		radius := valueToFloat64(vm.pop())
		z2 := valueToFloat64(vm.pop())
		y2 := valueToFloat64(vm.pop())
		x2 := valueToFloat64(vm.pop())
		z1 := valueToFloat64(vm.pop())
		y1 := valueToFloat64(vm.pop())
		x1 := valueToFloat64(vm.pop())
		dx, dy, dz := x2-x1, y2-y1, z2-z1
		vm.push(dx*dx+dy*dy+dz*dz <= radius*radius)
	case OpAngle2D:
		if len(vm.stack) < 4 {
			return fmt.Errorf("stack underflow for Angle2D")
		}
		y2 := valueToFloat64(vm.pop())
		x2 := valueToFloat64(vm.pop())
		y1 := valueToFloat64(vm.pop())
		x1 := valueToFloat64(vm.pop())
		vm.push(math.Atan2(y2-y1, x2-x1))

	case OpMatMul:
		if vm.ip+3 > len(vm.chunk.Code) {
			return fmt.Errorf("unexpected end of code for OpMatMul")
		}
		rci := int(vm.chunk.Code[vm.ip])
		aci := int(vm.chunk.Code[vm.ip+1])
		bci := int(vm.chunk.Code[vm.ip+2])
		vm.ip += 3
		if rci >= len(vm.chunk.Constants) || aci >= len(vm.chunk.Constants) || bci >= len(vm.chunk.Constants) {
			return fmt.Errorf("OpMatMul: invalid constant index")
		}
		rName, _ := vm.chunk.Constants[rci].(string)
		aName, _ := vm.chunk.Constants[aci].(string)
		bName, _ := vm.chunk.Constants[bci].(string)
		if rName == "" || aName == "" || bName == "" {
			return fmt.Errorf("OpMatMul: names must be strings")
		}
		rName = strings.ToLower(rName)
		aName = strings.ToLower(aName)
		bName = strings.ToLower(bName)
		raIdx, ok := vm.chunk.Variables[rName]
		if !ok {
			return fmt.Errorf("OpMatMul: result variable %s not found", rName)
		}
		aaIdx, ok := vm.chunk.Variables[aName]
		if !ok {
			return fmt.Errorf("OpMatMul: matrix A %s not found", aName)
		}
		baIdx, ok := vm.chunk.Variables[bName]
		if !ok {
			return fmt.Errorf("OpMatMul: matrix B %s not found", bName)
		}
		rd := vm.chunk.VarDims[rName]
		ad := vm.chunk.VarDims[aName]
		bd := vm.chunk.VarDims[bName]
		if len(rd) != 2 || len(ad) != 2 || len(bd) != 2 {
			return fmt.Errorf("OpMatMul: all three must be 2D arrays")
		}
		n, m, p := ad[0], ad[1], bd[1]
		if bd[0] != m {
			return fmt.Errorf("OpMatMul: A columns (%d) != B rows (%d)", m, bd[0])
		}
		if rd[0] != n || rd[1] != p {
			return fmt.Errorf("OpMatMul: result must be %d×%d", n, p)
		}
		for raIdx >= len(vm.stack) {
			vm.stack = append(vm.stack, nil)
		}
		for aaIdx >= len(vm.stack) {
			vm.stack = append(vm.stack, nil)
		}
		for baIdx >= len(vm.stack) {
			vm.stack = append(vm.stack, nil)
		}
		aArr, ok := vm.stack[aaIdx].([]Value)
		if !ok {
			return fmt.Errorf("OpMatMul: %s is not an array", aName)
		}
		bArr, ok := vm.stack[baIdx].([]Value)
		if !ok {
			return fmt.Errorf("OpMatMul: %s is not an array", bName)
		}
		rArr, ok := vm.stack[raIdx].([]Value)
		if !ok {
			return fmt.Errorf("OpMatMul: %s is not an array", rName)
		}
		for i := 0; i < n; i++ {
			for j := 0; j < p; j++ {
				var sum float64
				for k := 0; k < m; k++ {
					sum += valueToFloat64(aArr[i*m+k]) * valueToFloat64(bArr[k*p+j])
				}
				rArr[i*p+j] = sum
			}
		}

	case OpLeftStr:
		if len(vm.stack) < 2 {
			return fmt.Errorf("stack underflow for Left")
		}
		n := valueToInt(vm.pop())
		s := valueToString(vm.pop())
		if n < 0 {
			n = 0
		}
		if n > len(s) {
			n = len(s)
		}
		vm.push(s[:n])
	case OpRightStr:
		if len(vm.stack) < 2 {
			return fmt.Errorf("stack underflow for Right")
		}
		n := valueToInt(vm.pop())
		s := valueToString(vm.pop())
		if n < 0 {
			n = 0
		}
		if n > len(s) {
			n = len(s)
		}
		vm.push(s[len(s)-n:])
	case OpMidStr:
		if len(vm.stack) < 3 {
			return fmt.Errorf("stack underflow for Mid")
		}
		count := valueToInt(vm.pop())
		start := valueToInt(vm.pop())
		s := valueToString(vm.pop())
		if start < 1 {
			start = 1
		}
		if count < 0 {
			count = 0
		}
		start0 := start - 1
		if start0 >= len(s) {
			vm.push("")
		} else {
			end := start0 + count
			if end > len(s) {
				end = len(s)
			}
			vm.push(s[start0:end])
		}
	case OpLenStr:
		if len(vm.stack) == 0 {
			return fmt.Errorf("stack underflow for Len")
		}
		vm.push(len(valueToString(vm.pop())))
	case OpEOF:
		if len(vm.stack) == 0 {
			return fmt.Errorf("stack underflow for EOF")
		}
		h := valueToInt(vm.pop())
		rd := vm.fileReaders[h]
		if rd == nil {
			vm.push(true)
		} else {
			_, err := rd.Peek(1)
			vm.push(err != nil)
		}

	case OpOpenFile:
		if len(vm.stack) < 2 {
			return fmt.Errorf("stack underflow for OpenFile")
		}
		mode := valueToInt(vm.pop())
		path := valueToString(vm.pop())
		flags := os.O_RDONLY
		if mode == 1 {
			flags = os.O_WRONLY | os.O_CREATE | os.O_TRUNC
		} else if mode == 2 {
			flags = os.O_WRONLY | os.O_CREATE | os.O_APPEND
		}
		f, err := os.OpenFile(path, flags, 0644)
		if err != nil {
			return fmt.Errorf("OpenFile %s: %w", path, err)
		}
		h := vm.nextFileHandle
		vm.nextFileHandle++
		vm.fileHandles[h] = f
		if mode == 0 {
			vm.fileReaders[h] = bufio.NewReader(f)
		}
		vm.push(h)

	case OpReadLine:
		if len(vm.stack) == 0 {
			return fmt.Errorf("stack underflow for ReadLine")
		}
		h := valueToInt(vm.pop())
		f, ok := vm.fileHandles[h]
		if !ok {
			return fmt.Errorf("invalid file handle: %d", h)
		}
		rd := vm.fileReaders[h]
		if rd == nil {
			rd = bufio.NewReader(f)
			vm.fileReaders[h] = rd
		}
		line, err := rd.ReadString('\n')
		if err != nil && line == "" {
			vm.push("")
		} else {
			vm.push(strings.TrimSuffix(line, "\n"))
		}

	case OpWriteLine:
		if len(vm.stack) < 2 {
			return fmt.Errorf("stack underflow for WriteLine")
		}
		text := valueToString(vm.pop())
		h := valueToInt(vm.pop())
		f, ok := vm.fileHandles[h]
		if !ok {
			return fmt.Errorf("invalid file handle: %d", h)
		}
		if _, err := fmt.Fprintln(f, text); err != nil {
			return fmt.Errorf("WriteLine: %w", err)
		}

	case OpCloseFile:
		if len(vm.stack) == 0 {
			return fmt.Errorf("stack underflow for CloseFile")
		}
		h := valueToInt(vm.pop())
		if f, ok := vm.fileHandles[h]; ok {
			f.Close()
			delete(vm.fileHandles, h)
			delete(vm.fileReaders, h)
		}

	case OpCreateArray:
		if vm.ip+1 >= len(vm.chunk.Code) {
			return fmt.Errorf("unexpected end of code for OpCreateArray")
		}
		nDims := int(vm.chunk.Code[vm.ip])
		vm.ip++
		if nDims < 1 || nDims > 8 {
			return fmt.Errorf("invalid array dimensions count: %d", nDims)
		}
		dims := make([]int, nDims)
		size := 1
		for i := 0; i < nDims; i++ {
			if vm.ip >= len(vm.chunk.Code) {
				return fmt.Errorf("unexpected end of code for OpCreateArray dims")
			}
			ci := int(vm.chunk.Code[vm.ip])
			vm.ip++
			if ci >= len(vm.chunk.Constants) {
				return fmt.Errorf("constant index out of bounds in OpCreateArray")
			}
			dims[i] = valueToInt(vm.chunk.Constants[ci])
			size *= dims[i]
		}
		if vm.ip >= len(vm.chunk.Code) {
			return fmt.Errorf("unexpected end of code for OpCreateArray varIndex")
		}
		varIndex := int(vm.chunk.Code[vm.ip])
		vm.ip++
		arr := make([]Value, size)
		for i := range arr {
			arr[i] = 0
		}
		for varIndex >= len(vm.stack) {
			vm.stack = append(vm.stack, nil)
		}
		vm.stack[varIndex] = arr

	case OpLoadArray:
		if vm.ip >= len(vm.chunk.Code) {
			return fmt.Errorf("unexpected end of code for OpLoadArray")
		}
		varIndex := int(vm.chunk.Code[vm.ip])
		vm.ip++
		var dims []int
		for name, idx := range vm.chunk.Variables {
			if idx == varIndex {
				dims = vm.chunk.VarDims[name]
				break
			}
		}
		if dims == nil {
			return fmt.Errorf("OpLoadArray: variable %d is not an array", varIndex)
		}
		if varIndex >= len(vm.stack) {
			return fmt.Errorf("OpLoadArray: variable %d not initialized", varIndex)
		}
		arr, ok := vm.stack[varIndex].([]Value)
		if !ok {
			return fmt.Errorf("OpLoadArray: slot %d is not an array", varIndex)
		}
		idx := 0
		stride := 1
		for d := len(dims) - 1; d >= 0; d-- {
			if len(vm.stack) == 0 {
				return fmt.Errorf("stack underflow for OpLoadArray indices")
			}
			i := valueToInt(vm.pop())
			idx += i * stride
			stride *= dims[d]
		}
		if idx < 0 || idx >= len(arr) {
			return fmt.Errorf("array index out of bounds: %d", idx)
		}
		vm.push(arr[idx])

	case OpStoreArray:
		if vm.ip >= len(vm.chunk.Code) {
			return fmt.Errorf("unexpected end of code for OpStoreArray")
		}
		varIndex := int(vm.chunk.Code[vm.ip])
		vm.ip++
		var dims []int
		for name, idx := range vm.chunk.Variables {
			if idx == varIndex {
				dims = vm.chunk.VarDims[name]
				break
			}
		}
		if dims == nil {
			return fmt.Errorf("OpStoreArray: variable %d is not an array", varIndex)
		}
		if varIndex >= len(vm.stack) {
			return fmt.Errorf("OpStoreArray: variable %d not initialized", varIndex)
		}
		arr, ok := vm.stack[varIndex].([]Value)
		if !ok {
			return fmt.Errorf("OpStoreArray: slot %d is not an array", varIndex)
		}
		// Stack is [..., value, index0, index1, ...] with index1 on top; pop indices first, then value
		idx := 0
		stride := 1
		for d := len(dims) - 1; d >= 0; d-- {
			if len(vm.stack) == 0 {
				return fmt.Errorf("stack underflow for OpStoreArray indices")
			}
			i := valueToInt(vm.pop())
			idx += i * stride
			stride *= dims[d]
		}
		if len(vm.stack) == 0 {
			return fmt.Errorf("stack underflow for OpStoreArray value")
		}
		value := vm.pop()
		if idx < 0 || idx >= len(arr) {
			return fmt.Errorf("array index out of bounds: %d", idx)
		}
		arr[idx] = value

	case OpQuit:
		vm.running = false

	case OpHalt:
		vm.running = false

	// Game-specific operations
	case OpLoadImage:
		if len(vm.stack) == 0 {
			return fmt.Errorf("stack underflow for LOADIMAGE")
		}
		filename := vm.pop()
		if vm.runtime != nil {
			if err := vm.runtime.LoadImage(valueToString(filename)); err != nil {
				return err
			}
		} else {
			fmt.Printf("LOADIMAGE: %v\n", filename)
		}

	case OpCreateSprite:
		if len(vm.stack) < 4 {
			return fmt.Errorf("stack underflow for CREATESPRITE")
		}
		y := vm.pop()
		x := vm.pop()
		image := vm.pop()
		id := vm.pop()
		if vm.runtime != nil {
			if err := vm.runtime.CreateSprite(valueToString(id), valueToString(image), valueToFloat64(x), valueToFloat64(y)); err != nil {
				return err
			}
		} else {
			fmt.Printf("CREATESPRITE: %v, %v, %v, %v\n", id, image, x, y)
		}

	case OpSetSpritePosition:
		if len(vm.stack) < 3 {
			return fmt.Errorf("stack underflow for SETSPRITEPOSITION")
		}
		y := vm.pop()
		x := vm.pop()
		id := vm.pop()
		if vm.runtime != nil {
			if err := vm.runtime.SetSpritePosition(valueToString(id), valueToFloat64(x), valueToFloat64(y)); err != nil {
				return err
			}
		} else {
			fmt.Printf("SETSPRITEPOSITION: %v, %v, %v\n", id, x, y)
		}

	case OpDrawSprite:
		if len(vm.stack) == 0 {
			return fmt.Errorf("stack underflow for DRAWSPRITE")
		}
		id := vm.pop()
		if vm.runtime != nil {
			if err := vm.runtime.DrawSprite(valueToString(id)); err != nil {
				return err
			}
		} else {
			fmt.Printf("DRAWSPRITE: %v\n", id)
		}

	case OpLoadModel:
		if len(vm.stack) == 0 {
			return fmt.Errorf("stack underflow for LOADMODEL")
		}
		filename := vm.pop()
		if vm.runtime != nil {
			if err := vm.runtime.LoadModel(valueToString(filename)); err != nil {
				return err
			}
		} else {
			fmt.Printf("LOADMODEL: %v\n", filename)
		}

	case OpCreateCamera:
		if len(vm.stack) < 4 {
			return fmt.Errorf("stack underflow for CREATECAMERA")
		}
		z := vm.pop()
		y := vm.pop()
		x := vm.pop()
		id := vm.pop()
		if vm.runtime != nil {
			if err := vm.runtime.CreateCamera(valueToString(id), valueToFloat64(x), valueToFloat64(y), valueToFloat64(z)); err != nil {
				return err
			}
		} else {
			fmt.Printf("CREATECAMERA: %v, %v, %v, %v\n", id, x, y, z)
		}

	case OpSetCameraPosition:
		if len(vm.stack) < 4 {
			return fmt.Errorf("stack underflow for SETCAMERAPOSITION")
		}
		z := vm.pop()
		y := vm.pop()
		x := vm.pop()
		id := vm.pop()
		if vm.runtime != nil {
			if err := vm.runtime.SetCameraPosition(valueToString(id), valueToFloat64(x), valueToFloat64(y), valueToFloat64(z)); err != nil {
				return err
			}
		} else {
			fmt.Printf("SETCAMERAPOSITION: %v, %v, %v, %v\n", id, x, y, z)
		}

	case OpDrawModel:
		if len(vm.stack) < 5 {
			return fmt.Errorf("stack underflow for DRAWMODEL")
		}
		scale := vm.pop()
		z := vm.pop()
		y := vm.pop()
		x := vm.pop()
		id := vm.pop()
		if vm.runtime != nil {
			if err := vm.runtime.DrawModel(valueToString(id), valueToFloat64(x), valueToFloat64(y), valueToFloat64(z), valueToFloat64(scale)); err != nil {
				return err
			}
		} else {
			fmt.Printf("DRAWMODEL: %v, %v, %v, %v, %v\n", id, x, y, z, scale)
		}

	case OpPlayMusic:
		if len(vm.stack) == 0 {
			return fmt.Errorf("stack underflow for PLAYMUSIC")
		}
		filename := vm.pop()
		if vm.runtime != nil {
			if err := vm.runtime.PlayMusic(valueToString(filename)); err != nil {
				return err
			}
		} else {
			fmt.Printf("PLAYMUSIC: %v\n", filename)
		}

	case OpPlaySound:
		if len(vm.stack) == 0 {
			return fmt.Errorf("stack underflow for PLAYSOUND")
		}
		filename := vm.pop()
		if vm.runtime != nil {
			if err := vm.runtime.PlaySound(valueToString(filename)); err != nil {
				return err
			}
		} else {
			fmt.Printf("PLAYSOUND: %v\n", filename)
		}

	case OpLoadSound:
		if len(vm.stack) == 0 {
			return fmt.Errorf("stack underflow for LOADSOUND")
		}
		filename := vm.pop()
		if vm.runtime != nil {
			if err := vm.runtime.LoadSound(valueToString(filename)); err != nil {
				return err
			}
		} else {
			fmt.Printf("LOADSOUND: %v\n", filename)
		}

	case OpCreatePhysicsBody:
		if len(vm.stack) < 6 {
			return fmt.Errorf("stack underflow for CREATEPHYSICSBODY")
		}
		mass := vm.pop()
		z := vm.pop()
		y := vm.pop()
		x := vm.pop()
		bodyType := vm.pop()
		id := vm.pop()
		if vm.runtime != nil {
			if err := vm.runtime.CreatePhysicsBody(valueToString(id), valueToString(bodyType), valueToFloat64(x), valueToFloat64(y), valueToFloat64(z), valueToFloat64(mass)); err != nil {
				return err
			}
		} else {
			fmt.Printf("CREATEPHYSICSBODY: %v, %v, %v, %v, %v, %v\n", id, bodyType, x, y, z, mass)
		}

	case OpSetVelocity:
		if len(vm.stack) < 4 {
			return fmt.Errorf("stack underflow for SETVELOCITY")
		}
		vz := vm.pop()
		vy := vm.pop()
		vx := vm.pop()
		id := vm.pop()
		if vm.runtime != nil {
			if err := vm.runtime.SetVelocity(valueToString(id), valueToFloat64(vx), valueToFloat64(vy), valueToFloat64(vz)); err != nil {
				return err
			}
		} else {
			fmt.Printf("SETVELOCITY: %v, %v, %v, %v\n", id, vx, vy, vz)
		}

	case OpApplyForce:
		if len(vm.stack) < 4 {
			return fmt.Errorf("stack underflow for APPLYFORCE")
		}
		fz := vm.pop()
		fy := vm.pop()
		fx := vm.pop()
		id := vm.pop()
		if vm.runtime != nil {
			if err := vm.runtime.ApplyForce(valueToString(id), valueToFloat64(fx), valueToFloat64(fy), valueToFloat64(fz)); err != nil {
				return err
			}
		} else {
			fmt.Printf("APPLYFORCE: %v, %v, %v, %v\n", id, fx, fy, fz)
		}

	case OpRayCast3D:
		if len(vm.stack) < 7 {
			return fmt.Errorf("stack underflow for RAYCAST3D")
		}
		maxDistance := vm.pop()
		dirz := vm.pop()
		diry := vm.pop()
		dirx := vm.pop()
		startz := vm.pop()
		starty := vm.pop()
		startx := vm.pop()
		if vm.runtime != nil {
			hit, outX, outY, outZ, err := vm.runtime.RayCast3D(valueToFloat64(startx), valueToFloat64(starty), valueToFloat64(startz), valueToFloat64(dirx), valueToFloat64(diry), valueToFloat64(dirz), valueToFloat64(maxDistance))
			if err != nil {
				return err
			}
			vm.push(hit)
			vm.push(outX)
			vm.push(outY)
			vm.push(outZ)
		} else {
			fmt.Printf("RAYCAST3D: %v, %v, %v, %v, %v, %v, %v\n", startx, starty, startz, dirx, diry, dirz, maxDistance)
			vm.push(false)
		}

	default:
		return fmt.Errorf("unknown opcode: %d", op)
	}

	return nil
}