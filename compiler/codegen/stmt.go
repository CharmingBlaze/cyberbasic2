package codegen

import (
	"cyberbasic/compiler/parser"
	"cyberbasic/compiler/semantic"
	"cyberbasic/compiler/vm"
	"fmt"
	"strconv"
	"strings"
)

// compileStatement compiles a single statement
func (e *Emitter) compileStatement(stmt parser.Node) error {
	if line := getSourceLine(stmt); line > 0 {
		e.chunk.SetLine(line)
	}
	switch node := stmt.(type) {
	case *parser.Assignment:
		return e.compileAssignment(node)
	case *parser.CompoundAssign:
		return e.compileCompoundAssign(node)
	case *parser.Call:
		return e.compileCall(node)
	case *parser.IfStatement:
		return e.compileIfStatement(node)
	case *parser.ForStatement:
		return e.compileForStatement(node)
	case *parser.WhileStatement:
		return e.compileWhileStatement(node)
	case *parser.MainLoopStatement:
		return e.compileMainLoopStatement(node)
	case *parser.FunctionDecl:
		return e.compileFunctionDecl(node)
	case *parser.SubDecl:
		return e.compileSubDecl(node)
	case *parser.ReturnStatement:
		return e.compileReturnStatement(node)
	case *parser.DimStatement:
		return e.compileDimStatement(node)
	case *parser.RedimStatement:
		return e.compileRedimStatement(node)
	case *parser.AppendStatement:
		return e.compileAppendStatement(node)
	case *parser.ConstStatement:
		return e.compileConstStatement(node)
	case *parser.EnumStatement:
		return e.compileEnumStatement(node)
	case *parser.TypeDecl:
		return nil
	case *parser.EntityDecl:
		return e.compileEntityDecl(node)
	case *parser.Identifier:
		return nil
	case *parser.GameCommand:
		return e.compileGameCommand(node)
	case *parser.SelectCaseStatement:
		return e.compileSelectCaseStatement(node)
	case *parser.RepeatStatement:
		return e.compileRepeatStatement(node)
	case *parser.ExitLoopStatement:
		return e.compileExitLoopStatement(node)
	case *parser.ContinueLoopStatement:
		return e.compileContinueLoopStatement(node)
	case *parser.AssertStatement:
		return e.compileAssertStatement(node)
	case *parser.OnEventStatement:
		return e.compileOnEventStatement(node)
	case *parser.StartCoroutineStatement:
		return e.compileStartCoroutineStatement(node)
	case *parser.YieldStatement:
		e.chunk.Write(byte(vm.OpYield))
		return nil
	case *parser.WaitSecondsStatement:
		if err := e.compileExpression(node.Seconds); err != nil {
			return err
		}
		e.chunk.Write(byte(vm.OpWaitSeconds))
		return nil
	case *parser.WaitFramesStatement:
		if err := e.compileExpression(node.Frames); err != nil {
			return err
		}
		e.chunk.Write(byte(vm.OpLoadConst))
		ci := e.chunk.WriteConstant(60.0)
		if err := checkConstIndex(ci, " for WaitFrames divisor"); err != nil {
			return err
		}
		e.chunk.Write(byte(ci))
		e.chunk.Write(byte(vm.OpDiv))
		e.chunk.Write(byte(vm.OpWaitSeconds))
		return nil
	case *parser.DataStatement:
		return e.compileDataStatement(node)
	case *parser.ReadStatement:
		return e.compileReadStatement(node)
	case *parser.RestoreStatement:
		e.chunk.Write(byte(vm.OpRestore))
		return nil
	case *parser.GosubStatement:
		return e.compileGosubStatement(node)
	default:
		return errWithLine(stmt, fmt.Errorf("unsupported statement type: %T", stmt))
	}
}

// compileEntityDecl emits CreateDict, SetDictKey for each property, then StoreGlobal(entityName).
func (e *Emitter) compileEntityDecl(ed *parser.EntityDecl) error {
	ci := e.chunk.WriteConstant("createdict")
	if err := checkConstIndex(ci, ""); err != nil {
		return err
	}
	e.chunk.Write(byte(vm.OpCallForeign))
	e.chunk.Write(byte(ci))
	e.chunk.Write(byte(0))
	entityLower := strings.ToLower(ed.Name)
	for i, p := range ed.Properties {
		if i > 0 {
			e.chunk.Write(byte(vm.OpDup))
		}
		keyIdx := e.chunk.WriteConstant(p.Name)
		if err := checkConstIndex(keyIdx, " for entity property key"); err != nil {
			return err
		}
		e.chunk.Write(byte(vm.OpLoadConst))
		e.chunk.Write(byte(keyIdx))
		if err := e.compileExpression(p.Value); err != nil {
			return err
		}
		setIdx := e.chunk.WriteConstant("setdictkey")
		if err := checkConstIndex(setIdx, ""); err != nil {
			return err
		}
		e.chunk.Write(byte(vm.OpCallForeign))
		e.chunk.Write(byte(setIdx))
		e.chunk.Write(byte(3))
		// SetDictKey returns the dict; keep it on stack for next property or StoreGlobal
	}
	globalIdx := e.chunk.WriteConstant(entityLower)
	if err := checkConstIndex(globalIdx, " for entity global"); err != nil {
		return err
	}
	e.chunk.Write(byte(vm.OpStoreGlobal))
	e.chunk.Write(byte(globalIdx))
	return nil
}

// compileAssignment compiles an assignment statement (scalar, array element, or entity property).
func (e *Emitter) compileAssignment(assign *parser.Assignment) error {
	// Entity property write: "EntityName.prop = value" -> value, OpStoreEntityProp(entityIdx, propIdx)
	if len(assign.Indices) == 0 && e.sem.EntityNames != nil && strings.Contains(assign.Variable, ".") {
		parts := strings.SplitN(assign.Variable, ".", 2)
		if len(parts) == 2 && e.sem.EntityNames[strings.ToLower(parts[0])] {
			entityLower := strings.ToLower(parts[0])
			propName := parts[1]
			if err := e.compileExpression(assign.Value); err != nil {
				return err
			}
			entityIdx := e.chunk.WriteConstant(entityLower)
			propIdx := e.chunk.WriteConstant(propName)
			if err := checkConstIndex(entityIdx, " for entity prop"); err != nil {
				return err
			}
			if err := checkConstIndex(propIdx, " for entity prop"); err != nil {
				return err
			}
			e.chunk.Write(byte(vm.OpStoreEntityProp))
			e.chunk.Write(byte(entityIdx))
			e.chunk.Write(byte(propIdx))
			return nil
		}
	}

	// DotObject property write: "window.title = value" or "a.b.c = value" (qualified name from parser)
	if len(assign.Indices) == 0 && strings.Contains(assign.Variable, ".") {
		parts := strings.Split(assign.Variable, ".")
		if len(parts) >= 2 {
			baseName := parts[0]
			path := parts[1:]
			lowerBase := strings.ToLower(baseName)
			if lowerBase != "rl" && lowerBase != "box2d" && lowerBase != "bullet" && lowerBase != "game" {
				if err := e.compileExpression(assign.Value); err != nil {
					return err
				}
				ident := &parser.Identifier{Name: baseName, Line: assign.Line, Col: assign.Col}
				if err := e.compileIdentifier(ident); err != nil {
					return err
				}
				return e.emitOpSetProp(path)
			}
		}
	}

	err := e.compileExpression(assign.Value)
	if err != nil {
		return err
	}

	if len(assign.Indices) > 0 {
		for _, idx := range assign.Indices {
			if err := e.compileExpression(idx); err != nil {
				return err
			}
		}
		varIndex, exists := e.chunk.GetVariable(assign.Variable)
		if !exists {
			return errWithLine(assign, fmt.Errorf("array variable not declared: %s", assign.Variable))
		}
		e.chunk.Write(byte(vm.OpStoreArray))
		e.chunk.Write(byte(varIndex))
		return nil
	}

	if e.funcParamIndices != nil {
		if idx, ok := e.funcParamIndices[strings.ToLower(assign.Variable)]; ok {
			e.chunk.Write(byte(vm.OpStoreParam))
			e.chunk.Write(byte(idx))
			return nil
		}
	}
	if varIndex, exists := e.chunk.GetVariable(assign.Variable); exists {
		e.chunk.Write(byte(vm.OpStoreVar))
		e.chunk.Write(byte(varIndex))
	} else {
		varIndex := e.chunk.AddVariable(assign.Variable)
		e.chunk.Write(byte(vm.OpStoreVar))
		e.chunk.Write(byte(varIndex))
	}
	return nil
}

// compileCompoundAssign compiles +=, -=, *=, /= (load var, load value, op, store var).
func (e *Emitter) compileCompoundAssign(ca *parser.CompoundAssign) error {
	var varIndex int
	if e.funcParamIndices != nil {
		if idx, ok := e.funcParamIndices[strings.ToLower(ca.Variable)]; ok {
			e.chunk.Write(byte(vm.OpLoadParam))
			e.chunk.Write(byte(idx))
			if err := e.compileExpression(ca.Value); err != nil {
				return err
			}
			switch ca.Op {
			case "+=":
				e.chunk.Write(byte(vm.OpAdd))
			case "-=":
				e.chunk.Write(byte(vm.OpSub))
			case "*=":
				e.chunk.Write(byte(vm.OpMul))
			case "/=":
				e.chunk.Write(byte(vm.OpDiv))
			default:
				return errWithLine(ca, fmt.Errorf("unsupported compound assign op: %s", ca.Op))
			}
			e.chunk.Write(byte(vm.OpStoreParam))
			e.chunk.Write(byte(idx))
			return nil
		}
	}
	if idx, exists := e.chunk.GetVariable(ca.Variable); exists {
		varIndex = idx
	} else {
		varIndex = e.chunk.AddVariable(ca.Variable)
	}
	// Load current value
	e.chunk.Write(byte(vm.OpLoadVar))
	e.chunk.Write(byte(varIndex))
	// Load RHS
	if err := e.compileExpression(ca.Value); err != nil {
		return err
	}
	switch ca.Op {
	case "+=":
		e.chunk.Write(byte(vm.OpAdd))
	case "-=":
		e.chunk.Write(byte(vm.OpSub))
	case "*=":
		e.chunk.Write(byte(vm.OpMul))
	case "/=":
		e.chunk.Write(byte(vm.OpDiv))
	default:
		return errWithLine(ca, fmt.Errorf("unsupported compound assign op: %s", ca.Op))
	}
	e.chunk.Write(byte(vm.OpStoreVar))
	e.chunk.Write(byte(varIndex))
	return nil
}

// compileSelectCaseStatement compiles SELECT CASE expr ... CASE val: block ... END SELECT
func (e *Emitter) compileSelectCaseStatement(s *parser.SelectCaseStatement) error {
	if err := e.compileExpression(s.Expr); err != nil {
		return err
	}
	var endJumpPositions []int
	for _, k := range s.Cases {
		e.chunk.Write(byte(vm.OpDup))
		if err := e.compileExpression(k.Value); err != nil {
			return err
		}
		e.chunk.Write(byte(vm.OpEqual))
		e.chunk.Write(byte(vm.OpJumpIfFalse))
		e.chunk.Write(byte(0))
		e.chunk.Write(byte(0))
		skipPos := len(e.chunk.Code) - 2
		e.chunk.Write(byte(vm.OpPop))
		e.chunk.Write(byte(vm.OpPop))
		for _, stmt := range k.Block.Statements {
			if err := e.compileStatement(stmt); err != nil {
				return err
			}
		}
		e.chunk.Write(byte(vm.OpJump))
		e.chunk.Write(byte(0))
		e.chunk.Write(byte(0))
		endJumpPositions = append(endJumpPositions, len(e.chunk.Code)-2)
		// Patch skip to here
		e.chunk.PatchJumpOffset(skipPos, len(e.chunk.Code)-skipPos-2)
	}
	if s.ElseBlock != nil {
		e.chunk.Write(byte(vm.OpPop)) // drop the SELECT expr
		for _, stmt := range s.ElseBlock.Statements {
			if err := e.compileStatement(stmt); err != nil {
				return err
			}
		}
	} else {
		e.chunk.Write(byte(vm.OpPop))
	}
	endTarget := len(e.chunk.Code)
	for _, pos := range endJumpPositions {
		e.chunk.PatchJumpOffset(pos, endTarget-pos-2)
	}
	return nil
}

// compileExitLoopStatement compiles EXIT FOR or EXIT WHILE (or BREAK) (jump to end of innermost loop).
func (e *Emitter) compileExitLoopStatement(ex *parser.ExitLoopStatement) error {
	if len(e.loopExitStack) == 0 {
		return errWithLine(ex, fmt.Errorf("EXIT/BREAK %s outside loop", ex.Kind))
	}
	e.chunk.Write(byte(vm.OpJump))
	e.chunk.Write(byte(0))
	e.chunk.Write(byte(0))
	e.loopExitStack[len(e.loopExitStack)-1] = append(e.loopExitStack[len(e.loopExitStack)-1], len(e.chunk.Code)-2)
	return nil
}

// compileContinueLoopStatement compiles CONTINUE FOR or CONTINUE WHILE (jump to loop head of innermost matching loop).
func (e *Emitter) compileContinueLoopStatement(cl *parser.ContinueLoopStatement) error {
	if len(e.loopContinueStack) == 0 {
		return errWithLine(cl, fmt.Errorf("CONTINUE %s outside loop", cl.Kind))
	}
	e.chunk.Write(byte(vm.OpJump))
	e.chunk.Write(byte(0))
	e.chunk.Write(byte(0))
	e.loopContinueStack[len(e.loopContinueStack)-1] = append(e.loopContinueStack[len(e.loopContinueStack)-1], len(e.chunk.Code)-2)
	return nil
}

// compileAssertStatement compiles ASSERT condition [, message] to condition + message + CallForeign Assert(2).
func (e *Emitter) compileAssertStatement(a *parser.AssertStatement) error {
	if err := e.compileExpression(a.Condition); err != nil {
		return err
	}
	if a.Message != nil {
		if err := e.compileExpression(a.Message); err != nil {
			return err
		}
	} else {
		idx := e.chunk.WriteConstant("assertion failed")
		e.chunk.Write(byte(vm.OpLoadConst))
		e.chunk.Write(byte(idx))
	}
	idx := e.chunk.WriteConstant("assert")
	e.chunk.Write(byte(vm.OpCallForeign))
	e.chunk.Write(byte(idx))
	e.chunk.Write(byte(2))
	return nil
}

// compileRepeatStatement compiles REPEAT ... UNTIL condition (jump back when condition false)
func (e *Emitter) compileRepeatStatement(r *parser.RepeatStatement) error {
	loopStart := len(e.chunk.Code)
	wrapFrame := isGameLoopCondition(r.Condition, true)
	hybridMode := wrapFrame && (e.sem.UserFuncs["update"] || e.sem.UserFuncs["draw"])
	if hybridMode {
		wrapFrame = false
	} else if wrapFrame && r.Body != nil && (e.bodyCallsUserSub(r.Body.Statements) || bodyContainsFrameBoundaries(r.Body.Statements)) {
		wrapFrame = false
	}
	use3D := wrapFrame && r.Body != nil && bodyContains3DDraw(r.Body.Statements)
	if hybridMode {
		e.emitHybridLoopBody()
	} else {
		if wrapFrame {
			e.emitFrameWrap( "BeginDrawing")
			if use3D {
				e.emitFrameWrap( "BeginMode3D")
			} else {
				e.emitFrameWrap( "BeginMode2D")
			}
		}
		for _, stmt := range r.Body.Statements {
			// Emit EndMode2D/EndMode3D right before SYNC so 2D/3D content is flushed before EndDrawing.
			if wrapFrame && bodyContainsSync(r.Body.Statements) {
				n := unwrapStatement(stmt)
				if gc, ok := n.(*parser.GameCommand); ok && strings.ToLower(gc.Command) == "sync" {
					if use3D {
						e.emitFrameWrap( "EndMode3D")
					} else {
						e.emitFrameWrap( "EndMode2D")
					}
				}
			}
			if err := e.compileStatement(stmt); err != nil {
				return err
			}
		}
		if wrapFrame && !bodyContainsSync(r.Body.Statements) {
			// When body contains SYNC, omit EndDrawing; SYNC does it.
			if use3D {
				e.emitFrameWrap( "EndMode3D")
			} else {
				e.emitFrameWrap( "EndMode2D")
			}
			e.emitFrameWrap( "EndDrawing")
		}
	}
	if err := e.compileExpression(r.Condition); err != nil {
		return err
	}
	e.chunk.Write(byte(vm.OpJumpIfFalse))
	e.chunk.Write(byte(0))
	e.chunk.Write(byte(0))
	offset := loopStart - len(e.chunk.Code)
	e.chunk.PatchJumpOffset(len(e.chunk.Code)-2, offset)
	return nil
}

// compileIfStatement compiles an IF statement (with optional ELSEIF and ELSE)
func (e *Emitter) compileIfStatement(ifStmt *parser.IfStatement) error {
	// Compile condition
	err := e.compileExpression(ifStmt.Condition)
	if err != nil {
		return err
	}

	// Emit jump if false (skip then block; go to first ELSEIF or ELSE or end)
	e.chunk.Write(byte(vm.OpJumpIfFalse))
	e.chunk.Write(byte(0))
	e.chunk.Write(byte(0))
	jumpPos := len(e.chunk.Code) - 2

	// Compile then block
	for _, stmt := range ifStmt.ThenBlock.Statements {
		err = e.compileStatement(stmt)
		if err != nil {
			return err
		}
	}

	// Jump over all ELSEIF/ELSE to end (patched later)
	e.chunk.Write(byte(vm.OpJump))
	e.chunk.Write(byte(0))
	e.chunk.Write(byte(0))
	endJumpPos := len(e.chunk.Code) - 2

	// Patch first "jump if false" to here (start of first ELSEIF or ELSE or after end)
	e.chunk.PatchJumpOffset(jumpPos, len(e.chunk.Code)-jumpPos-2)

	var endJumpPositions []int
	endJumpPositions = append(endJumpPositions, endJumpPos)

	// Compile each ELSEIF branch
	for _, branch := range ifStmt.ElseIfs {
		err = e.compileExpression(branch.Condition)
		if err != nil {
			return err
		}
		e.chunk.Write(byte(vm.OpJumpIfFalse))
		e.chunk.Write(byte(0))
		e.chunk.Write(byte(0))
		elseIfJumpPos := len(e.chunk.Code) - 2

		for _, stmt := range branch.Block.Statements {
			err = e.compileStatement(stmt)
			if err != nil {
				return err
			}
		}
		e.chunk.Write(byte(vm.OpJump))
		e.chunk.Write(byte(0))
		e.chunk.Write(byte(0))
		elseIfEndPos := len(e.chunk.Code) - 2
		endJumpPositions = append(endJumpPositions, elseIfEndPos)
		e.chunk.PatchJumpOffset(elseIfJumpPos, len(e.chunk.Code)-elseIfJumpPos-2)
	}

	// Optional ELSE block
	if ifStmt.ElseBlock != nil {
		for _, stmt := range ifStmt.ElseBlock.Statements {
			err = e.compileStatement(stmt)
			if err != nil {
				return err
			}
		}
	}

	// Patch all "jump to end" offsets (after then block and after each ELSEIF block)
	for _, pos := range endJumpPositions {
		e.chunk.PatchJumpOffset(pos, len(e.chunk.Code)-pos-2)
	}

	return nil
}

// compileForStatement compiles a FOR loop
func (e *Emitter) compileForStatement(forStmt *parser.ForStatement) error {
	e.loopExitStack = append(e.loopExitStack, nil)
	e.loopContinueStack = append(e.loopContinueStack, nil)
	// Initialize loop variable
	err := e.compileExpression(forStmt.Start)
	if err != nil {
		return err
	}

	varIndex := e.chunk.AddVariable(forStmt.Variable)
	e.chunk.Write(byte(vm.OpStoreVar))
	e.chunk.Write(byte(varIndex))

	// Loop start
	loopStart := len(e.chunk.Code)

	// Check condition
	err = e.compileIdentifier(&parser.Identifier{Name: forStmt.Variable})
	if err != nil {
		return err
	}

	err = e.compileExpression(forStmt.End)
	if err != nil {
		return err
	}

	e.chunk.Write(byte(vm.OpGreater))
	e.chunk.Write(byte(vm.OpJumpIfTrue))
	e.chunk.Write(byte(0))
	e.chunk.Write(byte(0))
	exitJumpPos := len(e.chunk.Code) - 2

	// Compile loop body
	for _, stmt := range forStmt.Body.Statements {
		err = e.compileStatement(stmt)
		if err != nil {
			return err
		}
	}

	// CONTINUE FOR jumps here (increment then re-check condition)
	continueTargetIP := len(e.chunk.Code)
	// Increment loop variable
	err = e.compileIdentifier(&parser.Identifier{Name: forStmt.Variable})
	if err != nil {
		return err
	}

	if forStmt.Step != nil {
		err = e.compileExpression(forStmt.Step)
	} else {
		// Default step is 1
		err = e.compileNumber(&parser.Number{Value: "1"})
	}

	if err != nil {
		return err
	}

	e.chunk.Write(byte(vm.OpAdd))
	e.chunk.Write(byte(vm.OpStoreVar))
	e.chunk.Write(byte(varIndex))

	// Jump back to loop start (2-byte offset for large loop bodies)
	e.chunk.Write(byte(vm.OpJump))
	e.chunk.Write(byte(0))
	e.chunk.Write(byte(0))
	e.chunk.PatchJumpOffset(len(e.chunk.Code)-2, loopStart-len(e.chunk.Code))

	// Fix exit jump offset
	e.chunk.PatchJumpOffset(exitJumpPos, len(e.chunk.Code)-exitJumpPos-2)
	// Patch EXIT FOR jumps
	for _, pos := range e.loopExitStack[len(e.loopExitStack)-1] {
		e.chunk.PatchJumpOffset(pos, len(e.chunk.Code)-pos-2)
	}
	e.loopExitStack = e.loopExitStack[:len(e.loopExitStack)-1]
	// Patch CONTINUE FOR jumps to increment
	for _, pos := range e.loopContinueStack[len(e.loopContinueStack)-1] {
		e.chunk.PatchJumpOffset(pos, continueTargetIP-pos-2)
	}
	e.loopContinueStack = e.loopContinueStack[:len(e.loopContinueStack)-1]

	return nil
}

// bodyCallsUserSub returns true if the given statements (or any nested block) contain a call to a user-defined SUB or FUNCTION.
// When true, we skip automatic BeginDrawing/EndDrawing so the user's own Draw() (or similar) is not double-wrapped and flicker is avoided.
func (e *Emitter) bodyCallsUserSub(statements []parser.Node) bool {
	return WalkStatements(statements, func(n parser.Node) bool {
		call, ok := n.(*parser.Call)
		return ok && e.sem.UserFuncs != nil && e.sem.UserFuncs[strings.ToLower(call.Name)]
	})
}

// emitFrameWrap emits a no-arg foreign call (BeginDrawing, EndDrawing, BeginMode2D, EndMode2D). Used for automatic frame wrapping.
func (e *Emitter) emitFrameWrap(name string) {
	idx := e.chunk.WriteConstant(strings.ToLower(name))
	if idx > MaxConstIndex {
		return
	}
	e.chunk.Write(byte(vm.OpCallForeign))
	e.chunk.Write(byte(idx))
	e.chunk.Write(byte(0))
}

// emitHybridLoopBody emits a single runtime StepFrame call so all hybrid entry points share the same fixed-step behavior.
func (e *Emitter) emitHybridLoopBody() {
	idx := e.chunk.WriteConstant("stepframe")
	if idx > MaxConstIndex {
		return
	}
	e.chunk.Write(byte(vm.OpCallForeign))
	e.chunk.Write(byte(idx))
	e.chunk.Write(byte(0))
}

// compileWhileStatement compiles a WHILE loop
func (e *Emitter) compileWhileStatement(whileStmt *parser.WhileStatement) error {
	e.loopExitStack = append(e.loopExitStack, nil)
	e.loopContinueStack = append(e.loopContinueStack, nil)
	// Loop start (CONTINUE WHILE jumps here)
	loopStart := len(e.chunk.Code)
	wrapFrame := isGameLoopCondition(whileStmt.Condition, false)
	hybridMode := wrapFrame && (e.sem.UserFuncs["update"] || e.sem.UserFuncs["draw"])
	if hybridMode {
		wrapFrame = false
	} else if wrapFrame && whileStmt.Body != nil && (e.bodyCallsUserSub(whileStmt.Body.Statements) || bodyContainsFrameBoundaries(whileStmt.Body.Statements)) {
		wrapFrame = false // user does their own BeginDrawing/EndDrawing (or calls Draw()), avoid double wrap and flicker
	}
	use3D := wrapFrame && whileStmt.Body != nil && bodyContains3DDraw(whileStmt.Body.Statements)

	// Compile condition
	err := e.compileExpression(whileStmt.Condition)
	if err != nil {
		return err
	}

	// Jump if false (2-byte offset)
	e.chunk.Write(byte(vm.OpJumpIfFalse))
	e.chunk.Write(byte(0))
	e.chunk.Write(byte(0))
	exitJumpPos := len(e.chunk.Code) - 2

	if hybridMode {
		e.emitHybridLoopBody()
	} else {
		if wrapFrame {
			e.emitFrameWrap( "BeginDrawing")
			if use3D {
				e.emitFrameWrap( "BeginMode3D")
			} else {
				e.emitFrameWrap( "BeginMode2D")
			}
		}
		for _, stmt := range whileStmt.Body.Statements {
			// Emit EndMode2D/EndMode3D right before SYNC so 2D/3D content is flushed before EndDrawing.
			if wrapFrame && bodyContainsSync(whileStmt.Body.Statements) {
				n := unwrapStatement(stmt)
				if gc, ok := n.(*parser.GameCommand); ok && strings.ToLower(gc.Command) == "sync" {
					if use3D {
						e.emitFrameWrap( "EndMode3D")
					} else {
						e.emitFrameWrap( "EndMode2D")
					}
				}
			}
			err = e.compileStatement(stmt)
			if err != nil {
				return err
			}
		}
		if wrapFrame && !bodyContainsSync(whileStmt.Body.Statements) {
			// When body contains SYNC, omit EndDrawing; SYNC does it.
			if use3D {
				e.emitFrameWrap( "EndMode3D")
			} else {
				e.emitFrameWrap( "EndMode2D")
			}
			e.emitFrameWrap( "EndDrawing")
		}
	}

	// Jump back to loop start (2-byte offset)
	e.chunk.Write(byte(vm.OpJump))
	e.chunk.Write(byte(0))
	e.chunk.Write(byte(0))
	e.chunk.PatchJumpOffset(len(e.chunk.Code)-2, loopStart-len(e.chunk.Code))

	// Fix exit jump offset
	e.chunk.PatchJumpOffset(exitJumpPos, len(e.chunk.Code)-exitJumpPos-2)
	// Patch EXIT WHILE jumps
	for _, pos := range e.loopExitStack[len(e.loopExitStack)-1] {
		e.chunk.PatchJumpOffset(pos, len(e.chunk.Code)-pos-2)
	}
	e.loopExitStack = e.loopExitStack[:len(e.loopExitStack)-1]
	// Patch CONTINUE WHILE jumps to condition
	for _, pos := range e.loopContinueStack[len(e.loopContinueStack)-1] {
		e.chunk.PatchJumpOffset(pos, loopStart-pos-2)
	}
	e.loopContinueStack = e.loopContinueStack[:len(e.loopContinueStack)-1]

	return nil
}

// compileMainLoopStatement compiles MAINLOOP...ENDMAIN (equivalent to WHILE NOT WindowShouldClose()...WEND with frame wrap).
func (e *Emitter) compileMainLoopStatement(m *parser.MainLoopStatement) error {
	e.loopExitStack = append(e.loopExitStack, nil)
	e.loopContinueStack = append(e.loopContinueStack, nil)
	loopStart := len(e.chunk.Code)

	wrapFrame := true
	if m.Body != nil && (e.bodyCallsUserSub(m.Body.Statements) || bodyContainsFrameBoundaries(m.Body.Statements)) {
		wrapFrame = false
	}
	use3D := wrapFrame && m.Body != nil && bodyContains3DDraw(m.Body.Statements)
	body := m.Body
	if body == nil {
		body = &parser.Block{}
	}

	// Condition: NOT WindowShouldClose()
	idx := e.chunk.WriteConstant("windowshouldclose")
	if err := checkConstIndex(idx, ""); err != nil {
		return err
	}
	e.chunk.Write(byte(vm.OpCallForeign))
	e.chunk.Write(byte(idx))
	e.chunk.Write(byte(0))
	e.chunk.Write(byte(vm.OpNot))

	e.chunk.Write(byte(vm.OpJumpIfFalse))
	e.chunk.Write(byte(0))
	e.chunk.Write(byte(0))
	exitJumpPos := len(e.chunk.Code) - 2

	if wrapFrame {
		e.emitFrameWrap( "BeginDrawing")
		if use3D {
			e.emitFrameWrap( "BeginMode3D")
		} else {
			e.emitFrameWrap( "BeginMode2D")
		}
	}
	for _, stmt := range body.Statements {
		if wrapFrame && bodyContainsSync(body.Statements) {
			n := unwrapStatement(stmt)
			if gc, ok := n.(*parser.GameCommand); ok && strings.ToLower(gc.Command) == "sync" {
				if use3D {
					e.emitFrameWrap( "EndMode3D")
				} else {
					e.emitFrameWrap( "EndMode2D")
				}
			}
		}
		err := e.compileStatement(stmt)
		if err != nil {
			return err
		}
	}
	if wrapFrame && !bodyContainsSync(body.Statements) {
		if use3D {
			e.emitFrameWrap( "EndMode3D")
		} else {
			e.emitFrameWrap( "EndMode2D")
		}
		e.emitFrameWrap( "EndDrawing")
	}

	e.chunk.Write(byte(vm.OpJump))
	e.chunk.Write(byte(0))
	e.chunk.Write(byte(0))
	e.chunk.PatchJumpOffset(len(e.chunk.Code)-2, loopStart-len(e.chunk.Code))

	e.chunk.PatchJumpOffset(exitJumpPos, len(e.chunk.Code)-exitJumpPos-2)
	for _, pos := range e.loopExitStack[len(e.loopExitStack)-1] {
		e.chunk.PatchJumpOffset(pos, len(e.chunk.Code)-pos-2)
	}
	e.loopExitStack = e.loopExitStack[:len(e.loopExitStack)-1]
	for _, pos := range e.loopContinueStack[len(e.loopContinueStack)-1] {
		e.chunk.PatchJumpOffset(pos, loopStart-pos-2)
	}
	e.loopContinueStack = e.loopContinueStack[:len(e.loopContinueStack)-1]

	return nil
}

// compileFunctionDecl compiles a function declaration. VM replaces stack with [arg0, arg1, ...] on call, so no param stores.
func (e *Emitter) compileFunctionDecl(fn *parser.FunctionDecl) error {
	name := semantic.QualifiedName(fn)
	e.chunk.Functions[name] = len(e.chunk.Code)
	// Params map to stack indices 0, 1, ... so body sees a=0, b=1, etc.
	e.funcParamIndices = make(map[string]int)
	for i, p := range fn.Parameters {
		e.funcParamIndices[strings.ToLower(p)] = i
	}
	for _, stmt := range fn.Body.Statements {
		if err := e.compileStatement(stmt); err != nil {
			return err
		}
	}
	e.funcParamIndices = nil
	return nil
}

// compileSubDecl compiles a sub procedure declaration. Parameters use OpLoadParam/OpStoreParam; other vars use global stack slots; caller stack is preserved across OpCallUser.
func (e *Emitter) compileSubDecl(sub *parser.SubDecl) error {
	name := semantic.QualifiedName(sub)
	e.chunk.Functions[name] = len(e.chunk.Code)
	e.funcParamIndices = make(map[string]int)
	for i, p := range sub.Parameters {
		e.funcParamIndices[strings.ToLower(p)] = i
	}
	for _, stmt := range sub.Body.Statements {
		if err := e.compileStatement(stmt); err != nil {
			return err
		}
	}
	e.funcParamIndices = nil
	e.chunk.Write(byte(vm.OpReturn))
	return nil
}

func (e *Emitter) compileDataStatement(d *parser.DataStatement) error {
	for _, v := range d.Values {
		var val vm.Value
		switch n := v.(type) {
		case *parser.Number:
			f, _ := strconv.ParseFloat(n.Value, 64)
			val = f
		case *parser.StringLiteral:
			val = n.Value
		case *parser.Boolean:
			val = n.Value
		case *parser.NilLiteral:
			val = nil
		case *parser.Identifier:
			val = float64(0) // identifier in DATA: use 0 (could resolve consts if added)
		default:
			val = float64(0)
		}
		e.chunk.DataValues = append(e.chunk.DataValues, val)
	}
	return nil
}

func (e *Emitter) compileReadStatement(r *parser.ReadStatement) error {
	for _, v := range r.Variables {
		name := ""
		switch n := v.(type) {
		case *parser.Identifier:
			name = n.Name
		case *parser.Call:
			if len(n.Arguments) > 0 {
				name = n.Name // READ a(i) - use base name for now (scalar)
			}
		}
		if name == "" {
			return errWithLine(v, fmt.Errorf("READ requires variable name"))
		}
		idx := e.chunk.AddVariable(name)
		if idx > 255 {
			return fmt.Errorf("too many variables for READ")
		}
		e.chunk.Write(byte(vm.OpRead))
		e.chunk.Write(byte(idx))
	}
	return nil
}

func (e *Emitter) compileGosubStatement(g *parser.GosubStatement) error {
	name := strings.ToLower(g.SubName)
	if !e.sem.UserFuncs[name] {
		return errWithLine(g, fmt.Errorf("unknown sub for GOSUB: %s", g.SubName))
	}
	idx := e.chunk.WriteConstant(name)
	if err := checkConstIndex(idx, " for GOSUB"); err != nil {
		return err
	}
	e.chunk.Write(byte(vm.OpGosub))
	e.chunk.Write(byte(idx))
	e.chunk.Write(byte(0)) // arg count
	return nil
}

// compileStartCoroutineStatement compiles StartCoroutine SubName(): emit OpStartCoroutine with 2-byte target offset (patched after decls)
func (e *Emitter) compileStartCoroutineStatement(stmt *parser.StartCoroutineStatement) error {
	name := strings.ToLower(stmt.SubName)
	if !e.sem.UserFuncs[name] {
		candidates := make([]string, 0, len(e.sem.UserFuncs))
		for q := range e.sem.UserFuncs {
			candidates = append(candidates, q)
		}
		msg := fmt.Sprintf("unknown sub for StartCoroutine: %s", stmt.SubName)
		if sug := nearestName(name, candidates, 3); sug != "" {
			msg += " (did you mean " + sug + "?)"
		}
		return errWithLine(stmt, fmt.Errorf("%s", msg))
	}
	e.chunk.Write(byte(vm.OpStartCoroutine))
	e.chunk.WriteJumpOffset(0)
	nameIdx := e.chunk.WriteConstant(name)
	if err := checkConstIndex(nameIdx, " for StartCoroutine name"); err != nil {
		return err
	}
	e.chunk.Write(byte(nameIdx))
	e.startCoroutinePatchList = append(e.startCoroutinePatchList, startCoroutinePatch{patchPos: len(e.chunk.Code) - 3, subName: name})
	return nil
}

// compileOnEventStatement compiles On KeyDown("X") ... End On: emit OpRegisterEvent (handler offset patched later)
func (e *Emitter) compileOnEventStatement(on *parser.OnEventStatement) error {
	eventTypeConst := e.chunk.WriteConstant(strings.ToLower(on.EventType))
	keyConst := e.chunk.WriteConstant(on.Key)
	if err := checkConstIndex(eventTypeConst, " for On event"); err != nil {
		return err
	}
	if err := checkConstIndex(keyConst, " for On event"); err != nil {
		return err
	}
	e.chunk.Write(byte(vm.OpRegisterEvent))
	e.chunk.Write(byte(eventTypeConst))
	e.chunk.Write(byte(keyConst))
	e.chunk.WriteJumpOffset(0)
	e.eventPatchList = append(e.eventPatchList, eventPatch{patchPos: len(e.chunk.Code) - 2, stmt: on})
	return nil
}

// compileReturnStatement compiles a RETURN statement
func (e *Emitter) compileReturnStatement(ret *parser.ReturnStatement) error {
	if ret.Value != nil {
		err := e.compileExpression(ret.Value)
		if err != nil {
			return err
		}
		e.chunk.Write(byte(vm.OpReturnVal))
	} else {
		e.chunk.Write(byte(vm.OpReturn))
	}
	return nil
}

// constIntFromNode returns the integer value of a constant dimension (Number node). Returns error if not constant.
func constIntFromNode(n parser.Node) (int, error) {
	switch node := n.(type) {
	case *parser.Number:
		return parseInt(node.Value)
	default:
		return 0, fmt.Errorf("array dimension must be a constant number")
	}
}

// compileDimStatement compiles a DIM statement (scalar or array)
func (e *Emitter) compileDimStatement(dim *parser.DimStatement) error {
	for _, v := range dim.Variables {
		varIndex := e.chunk.AddVariable(v.Name)

		if len(v.Dimensions) > 0 {
			// Fixed-size array: DIM a(10, 20)
			// Array: DIM a(10, 20) AS Integer
			dims := make([]int, len(v.Dimensions))
			constIndices := make([]byte, len(v.Dimensions))
			for i, d := range v.Dimensions {
				n, err := constIntFromNode(d)
				if err != nil {
					return fmt.Errorf("dimension %d for %s: %w", i+1, v.Name, err)
				}
				if n < 1 {
					return fmt.Errorf("dimension %d for %s must be >= 1", i+1, v.Name)
				}
				dims[i] = n
				ci := e.chunk.WriteConstant(n)
				if err := checkConstIndex(ci, ""); err != nil {
					return err
				}
				constIndices[i] = byte(ci)
			}
			e.chunk.SetVarDims(v.Name, dims)
			e.chunk.Write(byte(vm.OpCreateArray))
			e.chunk.Write(byte(len(dims)))
			for _, b := range constIndices {
				e.chunk.Write(b)
			}
			e.chunk.Write(byte(varIndex))
			continue
		}

		// DIM a() - empty dynamic array (dimensions is [] from parser when we saw ())
		if v.Dimensions != nil && len(v.Dimensions) == 0 {
			e.chunk.SetVarDims(v.Name, []int{})
			e.chunk.Write(byte(vm.OpCreateArray))
			e.chunk.Write(byte(1))
			ci := e.chunk.WriteConstant(0)
			if err := checkConstIndex(ci, ""); err != nil {
				return err
			}
			e.chunk.Write(byte(ci))
			e.chunk.Write(byte(varIndex))
			continue
		}

		// Scalar: initialize with default value (dynamic type when v.Type == "")
		switch strings.ToLower(v.Type) {
		case "integer", "int", "":
			e.chunk.Write(byte(vm.OpLoadConst))
			e.chunk.Write(byte(e.chunk.WriteConstant(0)))
		case "string", "str":
			e.chunk.Write(byte(vm.OpLoadString))
			e.chunk.Write(byte(e.chunk.WriteConstant("")))
		case "float", "single", "double":
			e.chunk.Write(byte(vm.OpLoadConst))
			e.chunk.Write(byte(e.chunk.WriteConstant(0.0)))
		case "boolean", "bool":
			e.chunk.Write(byte(vm.OpLoadConst))
			e.chunk.Write(byte(e.chunk.WriteConstant(false)))
		case "vector2", "vector3", "body", "color":
			// Optional type hints (e.g. DIM pos AS Vector2); stored as dynamic 0, use .x/.y/.z or API later
			e.chunk.Write(byte(vm.OpLoadConst))
			e.chunk.Write(byte(e.chunk.WriteConstant(0)))
		default:
			e.chunk.Write(byte(vm.OpLoadConst))
			e.chunk.Write(byte(e.chunk.WriteConstant(0)))
		}
		e.chunk.Write(byte(vm.OpStoreVar))
		e.chunk.Write(byte(varIndex))
	}
	return nil
}

// compileRedimStatement compiles REDIM a(n) or REDIM a(n, m)
func (e *Emitter) compileRedimStatement(r *parser.RedimStatement) error {
	varIndex, exists := e.chunk.GetVariable(r.Variable)
	if !exists {
		return errWithLine(r, fmt.Errorf("variable %s not declared for REDIM", r.Variable))
	}
	for i := len(r.Dimensions) - 1; i >= 0; i-- {
		if err := e.compileExpression(r.Dimensions[i]); err != nil {
			return err
		}
	}
	e.chunk.Write(byte(vm.OpResizeArray))
	e.chunk.Write(byte(len(r.Dimensions)))
	e.chunk.Write(byte(varIndex))
	return nil
}

// compileAppendStatement compiles APPEND a, value
func (e *Emitter) compileAppendStatement(a *parser.AppendStatement) error {
	varIndex, exists := e.chunk.GetVariable(a.Variable)
	if !exists {
		return errWithLine(a, fmt.Errorf("variable %s not declared for APPEND", a.Variable))
	}
	if err := e.compileExpression(a.Value); err != nil {
		return err
	}
	e.chunk.Write(byte(vm.OpAppendArray))
	e.chunk.Write(byte(varIndex))
	return nil
}

// evalConstExpr evaluates a constant expression (number, string, boolean, or -number). Returns error if not constant.
func evalConstExpr(n parser.Node) (interface{}, error) {
	switch node := n.(type) {
	case *parser.Number:
		if strings.Contains(node.Value, ".") {
			f, err := parseFloat(node.Value)
			if err != nil {
				return nil, fmt.Errorf("CONST value not a valid number: %s", node.Value)
			}
			return f, nil
		}
		i, err := parseInt(node.Value)
		if err != nil {
			return nil, fmt.Errorf("CONST value not a valid integer: %s", node.Value)
		}
		return i, nil
	case *parser.StringLiteral:
		return node.Value, nil
	case *parser.Boolean:
		return node.Value, nil
	case *parser.UnaryOp:
		if node.Operator != "-" {
			return nil, fmt.Errorf("CONST value must be a literal or -number")
		}
		num, ok := node.Operand.(*parser.Number)
		if !ok {
			return nil, fmt.Errorf("CONST value must be a literal or -number")
		}
		if strings.Contains(num.Value, ".") {
			f, err := parseFloat(num.Value)
			if err != nil {
				return nil, fmt.Errorf("CONST value not a valid number: %s", num.Value)
			}
			return -f, nil
		}
		i, err := parseInt(num.Value)
		if err != nil {
			return nil, fmt.Errorf("CONST value not a valid integer: %s", num.Value)
		}
		return -i, nil
	default:
		return nil, fmt.Errorf("CONST value must be a literal (number, string, true/false) or -number")
	}
}

// compileConstStatement compiles CONST name = value (, name = value)*
func (e *Emitter) compileConstStatement(cs *parser.ConstStatement) error {
	for _, d := range cs.Decls {
		val, err := evalConstExpr(d.Value)
		if err != nil {
			return err
		}
		idx := e.chunk.WriteConstant(val)
		if err := checkConstIndex(idx, ""); err != nil {
			return err
		}
		e.constIndices[strings.ToLower(d.Name)] = byte(idx)
	}
	return nil
}

// toInt64 converts a constant value to int64 for enum members.
func toInt64(v interface{}) (int64, error) {
	switch x := v.(type) {
	case int:
		return int64(x), nil
	case int64:
		return x, nil
	case float64:
		return int64(x), nil
	default:
		return 0, fmt.Errorf("enum value must be numeric, got %T", v)
	}
}

// resolveUDTConstantMember returns the value for TypeName.Member when the type is used as a constant group (eval or auto-increment).
func (e *Emitter) resolveUDTConstantMember(td *parser.TypeDecl, memberLower string) (interface{}, error) {
	nextVal := int64(0)
	for _, f := range td.Fields {
		if strings.ToLower(f.Name) == memberLower {
			if f.ConstValue != nil {
				return evalConstExpr(f.ConstValue)
			}
			return nextVal, nil
		}
		if f.ConstValue != nil {
			val, err := evalConstExpr(f.ConstValue)
			if err != nil {
				return nil, err
			}
			nextVal, err = toInt64(val)
			if err != nil {
				return nil, err
			}
		}
		nextVal++
	}
	fieldNames := make([]string, 0, len(td.Fields))
	for _, f := range td.Fields {
		fieldNames = append(fieldNames, f.Name)
	}
	msg := fmt.Sprintf("unknown member %s", memberLower)
	if sug := nearestName(memberLower, fieldNames, 3); sug != "" {
		msg += " (did you mean " + sug + "?)"
	}
	return nil, fmt.Errorf("%s", msg)
}

// compileEnumStatement compiles ENUM Name : a, b = 2, c ... members as constants (auto-increment from 0 or explicit value).
// Also records enum name -> member -> value in e.chunk.Enums for Enum.getValue/getName/hasValue at runtime.
func (e *Emitter) compileEnumStatement(es *parser.EnumStatement) error {
	enumName := strings.ToLower(es.Name)
	if e.chunk.Enums == nil {
		e.chunk.Enums = make(map[string]vm.EnumMembers)
	}
	members := make(vm.EnumMembers)
	nextVal := int64(0)
	for _, m := range es.Members {
		if m.Value != nil {
			val, err := evalConstExpr(m.Value)
			if err != nil {
				return fmt.Errorf("enum member %s: %w", m.Name, err)
			}
			nextVal, err = toInt64(val)
			if err != nil {
				return fmt.Errorf("enum member %s: %w", m.Name, err)
			}
		}
		memLower := strings.ToLower(m.Name)
		members[memLower] = nextVal
		idx := e.chunk.WriteConstant(nextVal)
		if err := checkConstIndex(idx, ""); err != nil {
			return err
		}
		e.constIndices[memLower] = byte(idx)
		nextVal++
	}
	e.chunk.Enums[enumName] = members
	return nil
}

// compileGameCommand compiles game-specific commands
func (e *Emitter) compileGameCommand(cmd *parser.GameCommand) error {
	// Compile arguments
	for _, arg := range cmd.Arguments {
		err := e.compileExpression(arg)
		if err != nil {
			return err
		}
	}

	// Case-insensitive: canonical form is lowercase
	command := strings.ToLower(cmd.Command)

	// Emit game command instruction
	switch command {
	case "print":
		e.chunk.Write(byte(vm.OpPrint))
	case "str":
		e.chunk.Write(byte(vm.OpStr))
	case "loadimage":
		e.chunk.Write(byte(vm.OpLoadImage))
	case "createsprite":
		e.chunk.Write(byte(vm.OpCreateSprite))
	case "setspriteposition":
		e.chunk.Write(byte(vm.OpSetSpritePosition))
	case "drawsprite":
		e.chunk.Write(byte(vm.OpDrawSprite))
	case "loadmodel":
		e.chunk.Write(byte(vm.OpLoadModel))
	case "createcamera":
		e.chunk.Write(byte(vm.OpCreateCamera))
	case "setcameraposition":
		e.chunk.Write(byte(vm.OpSetCameraPosition))
	case "drawmodel":
		e.chunk.Write(byte(vm.OpDrawModel))
	case "playmusic":
		e.chunk.Write(byte(vm.OpPlayMusic))
	case "playsound":
		e.chunk.Write(byte(vm.OpPlaySound))
	case "loadsound":
		e.chunk.Write(byte(vm.OpLoadSound))
	case "createphysicsbody":
		e.chunk.Write(byte(vm.OpCreatePhysicsBody))
	case "setvelocity":
		e.chunk.Write(byte(vm.OpSetVelocity))
	case "applyforce":
		e.chunk.Write(byte(vm.OpApplyForce))
	case "raycast3d":
		e.chunk.Write(byte(vm.OpRayCast3D))
	case "sync":
		e.chunk.Write(byte(vm.OpSync))
	case "shouldclose":
		e.chunk.Write(byte(vm.OpShouldClose))
	default:
		return fmt.Errorf("unsupported game command: %s", cmd.Command)
	}

	return nil
}
