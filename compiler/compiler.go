package compiler

import (
	"cyberbasic/compiler/lexer"
	"cyberbasic/compiler/parser"
	"cyberbasic/compiler/vm"
	"fmt"
	"strings"
)

// Compiler handles the compilation of BASIC source code to bytecode
type Compiler struct {
	lexer         *lexer.Lexer
	parser        *parser.Parser
	vm            *vm.VM
	constIndices  map[string]byte            // const name -> chunk constant index (set during generateCode)
	typeDefs      map[string]*parser.TypeDecl // UDT name (lowercase) -> TYPE definition (filled in first pass)
	loopExitStack     [][]int // for EXIT FOR / EXIT WHILE: positions of jump offsets to patch
	loopContinueStack [][]int // for CONTINUE FOR / CONTINUE WHILE: positions of jump offsets to patch (target = loop head)
	userFuncs        map[string]bool            // user Sub/Function names (lowercase) for call resolution during codegen
	funcParamIndices map[string]int             // when compiling a function body: param name -> stack index 0,1,2,...
	eventPatchList         []eventPatch               // (patchPos, OnEventStatement) for patching handler offsets
	startCoroutinePatchList []startCoroutinePatch     // (patchPos, subName) for patching after decls
}

type eventPatch struct {
	patchPos int
	stmt     *parser.OnEventStatement
}

type startCoroutinePatch struct {
	patchPos int
	subName  string
}

// New creates a new compiler instance
func New() *Compiler {
	return &Compiler{
		vm: vm.NewVM(),
	}
}

// Compile compiles BASIC source code to bytecode
func (c *Compiler) Compile(source string) (*vm.Chunk, error) {
	// Lexical analysis
	l := lexer.New(source)
	tokens, err := l.Tokenize()
	if err != nil {
		return nil, fmt.Errorf("lexical error: %w", err)
	}

	// Parsing
	p := parser.New(tokens)
	program, err := p.Parse()
	if err != nil {
		return nil, fmt.Errorf("parse error: %w", err)
	}

	// Code generation
	chunk := vm.NewChunk()
	err = c.generateCode(program, chunk)
	if err != nil {
		return nil, fmt.Errorf("code generation error: %w", err)
	}

	return chunk, nil
}

// generateCode converts AST to bytecode
func (c *Compiler) generateCode(program *parser.Program, chunk *vm.Chunk) error {
	c.constIndices = make(map[string]byte)
	c.typeDefs = make(map[string]*parser.TypeDecl)
	c.eventPatchList = nil
	c.startCoroutinePatchList = nil
	// First pass: collect TYPE definitions and user function/sub names
	userFuncs := make(map[string]bool)
	var mainStmts, decls []parser.Node
	for _, stmt := range program.Statements {
		if td, ok := stmt.(*parser.TypeDecl); ok {
			c.typeDefs[strings.ToLower(td.Name)] = td
		}
		switch s := stmt.(type) {
		case *parser.ModuleStatement:
			for _, node := range s.Body {
				q := qualifiedName(node)
				userFuncs[q] = true
				decls = append(decls, node)
			}
		case *parser.FunctionDecl, *parser.SubDecl:
			userFuncs[qualifiedName(s)] = true
			decls = append(decls, stmt)
		default:
			mainStmts = append(mainStmts, stmt)
		}
	}
	c.userFuncs = userFuncs
	// Compile main program (no Function/Sub bodies)
	for _, stmt := range mainStmts {
		if err := c.compileStatement(stmt, chunk); err != nil {
			return err
		}
	}
	// Jump over all function/sub bodies
	jumpPos := len(chunk.Code)
	chunk.Write(byte(vm.OpJump))
	chunk.WriteJumpOffset(0)
	// Compile each Sub/Function: record offset, param pops, body, optional OpReturn
	for _, stmt := range decls {
		if err := c.compileDecl(stmt, chunk); err != nil {
			return err
		}
	}
	// Patch StartCoroutine target offsets (subs are now in chunk.Functions)
	for _, p := range c.startCoroutinePatchList {
		target, ok := chunk.GetFunction(p.subName)
		if !ok {
			return fmt.Errorf("unknown sub for StartCoroutine: %s", p.subName)
		}
		chunk.PatchJumpOffset(p.patchPos, target)
	}
	// Compile event handlers (On KeyDown/KeyPressed) and patch registration offsets
	for _, ep := range c.eventPatchList {
		handlerStart := len(chunk.Code)
		for _, stmt := range ep.stmt.Body.Statements {
			if err := c.compileStatement(stmt, chunk); err != nil {
				return err
			}
		}
		chunk.Write(byte(vm.OpReturn))
		chunk.PatchJumpOffset(ep.patchPos, handlerStart)
	}
	endPos := len(chunk.Code)
	chunk.PatchJumpOffset(jumpPos+1, endPos-(jumpPos+3))
	chunk.Write(byte(vm.OpHalt))
	return nil
}

// qualifiedName returns the lowercase name for a user function/sub (Module.Name or Name).
func qualifiedName(node parser.Node) string {
	switch n := node.(type) {
	case *parser.FunctionDecl:
		if n.ModuleName != "" {
			return strings.ToLower(n.ModuleName) + "." + strings.ToLower(n.Name)
		}
		return strings.ToLower(n.Name)
	case *parser.SubDecl:
		if n.ModuleName != "" {
			return strings.ToLower(n.ModuleName) + "." + strings.ToLower(n.Name)
		}
		return strings.ToLower(n.Name)
	default:
		return ""
	}
}

// compileDecl compiles a single FunctionDecl or SubDecl (called after main code; records chunk.Functions and emits body).
func (c *Compiler) compileDecl(stmt parser.Node, chunk *vm.Chunk) error {
	switch node := stmt.(type) {
	case *parser.FunctionDecl:
		return c.compileFunctionDecl(node, chunk)
	case *parser.SubDecl:
		return c.compileSubDecl(node, chunk)
	default:
		return fmt.Errorf("compileDecl: expected FunctionDecl or SubDecl, got %T", stmt)
	}
}

// getSourceLine returns the source line for error reporting (0 if node has no location).
func getSourceLine(node parser.Node) int {
	if loc, ok := node.(parser.HasSourceLoc); ok {
		return loc.GetLine()
	}
	return 0
}

// compileStatement compiles a single statement
func (c *Compiler) compileStatement(stmt parser.Node, chunk *vm.Chunk) error {
	if line := getSourceLine(stmt); line > 0 {
		chunk.SetLine(line)
	}
	switch node := stmt.(type) {
	case *parser.Assignment:
		return c.compileAssignment(node, chunk)
	case *parser.CompoundAssign:
		return c.compileCompoundAssign(node, chunk)
	case *parser.Call:
		return c.compileCall(node, chunk)
	case *parser.IfStatement:
		return c.compileIfStatement(node, chunk)
	case *parser.ForStatement:
		return c.compileForStatement(node, chunk)
	case *parser.WhileStatement:
		return c.compileWhileStatement(node, chunk)
	case *parser.FunctionDecl:
		return c.compileFunctionDecl(node, chunk)
	case *parser.SubDecl:
		return c.compileSubDecl(node, chunk)
	case *parser.ReturnStatement:
		return c.compileReturnStatement(node, chunk)
	case *parser.DimStatement:
		return c.compileDimStatement(node, chunk)
	case *parser.ConstStatement:
		return c.compileConstStatement(node, chunk)
	case *parser.EnumStatement:
		return c.compileEnumStatement(node, chunk)
	case *parser.TypeDecl:
		// TYPE definitions already collected in first pass; no code to emit
		return nil
	case *parser.Identifier:
		// Bare identifier - ignore (could be a variable reference without assignment)
		return nil
	case *parser.GameCommand:
		return c.compileGameCommand(node, chunk)
	case *parser.SelectCaseStatement:
		return c.compileSelectCaseStatement(node, chunk)
	case *parser.RepeatStatement:
		return c.compileRepeatStatement(node, chunk)
	case *parser.ExitLoopStatement:
		return c.compileExitLoopStatement(node, chunk)
	case *parser.ContinueLoopStatement:
		return c.compileContinueLoopStatement(node, chunk)
	case *parser.AssertStatement:
		return c.compileAssertStatement(node, chunk)
	case *parser.OnEventStatement:
		return c.compileOnEventStatement(node, chunk)
	case *parser.StartCoroutineStatement:
		return c.compileStartCoroutineStatement(node, chunk)
	case *parser.YieldStatement:
		chunk.Write(byte(vm.OpYield))
		return nil
	case *parser.WaitSecondsStatement:
		if err := c.compileExpression(node.Seconds, chunk); err != nil {
			return err
		}
		chunk.Write(byte(vm.OpWaitSeconds))
		return nil
	default:
		return fmt.Errorf("unsupported statement type: %T", stmt)
	}
}

// compileExpression compiles an expression
func (c *Compiler) compileExpression(expr parser.Node, chunk *vm.Chunk) error {
	if line := getSourceLine(expr); line > 0 {
		chunk.SetLine(line)
	}
	switch node := expr.(type) {
	case *parser.Number:
		return c.compileNumber(node, chunk)
	case *parser.StringLiteral:
		return c.compileString(node, chunk)
	case *parser.Boolean:
		return c.compileBoolean(node, chunk)
	case *parser.NilLiteral:
		return c.compileNilLiteral(chunk)
	case *parser.Identifier:
		return c.compileIdentifier(node, chunk)
	case *parser.BinaryOp:
		return c.compileBinaryOp(node, chunk)
	case *parser.UnaryOp:
		return c.compileUnaryOp(node, chunk)
	case *parser.Call:
		return c.compileCall(node, chunk)
	case *parser.MemberAccess:
		return c.compileMemberAccess(node, chunk)
	case *parser.JSONIndexAccess:
		return c.compileJSONIndexAccess(node, chunk)
	case *parser.DictLiteral:
		return c.compileDictLiteral(node, chunk)
	default:
		return fmt.Errorf("unsupported expression type: %T", expr)
	}
}

// compileDictLiteral compiles { k: v, ... } as CreateDict then SetDictKey for each pair.
func (c *Compiler) compileDictLiteral(node *parser.DictLiteral, chunk *vm.Chunk) error {
	ci := chunk.WriteConstant("createdict")
	if ci > 255 {
		return fmt.Errorf("too many constants")
	}
	chunk.Write(byte(vm.OpCallForeign))
	chunk.Write(byte(ci))
	chunk.Write(byte(0))
	for i, p := range node.Pairs {
		if i > 0 {
			chunk.Write(byte(vm.OpDup))
		}
		keyIdx := chunk.WriteConstant(p.Key)
		if keyIdx > 255 {
			return fmt.Errorf("too many constants for dict key")
		}
		if err := c.compileExpression(p.Value, chunk); err != nil {
			return err
		}
		setIdx := chunk.WriteConstant("setdictkey")
		if setIdx > 255 {
			return fmt.Errorf("too many constants")
		}
		chunk.Write(byte(vm.OpCallForeign))
		chunk.Write(byte(setIdx))
		chunk.Write(byte(3))
		// After SetDictKey we have (map, map) when i>0 because we dup'd; pop the duplicate
		if i < len(node.Pairs)-1 {
			chunk.Write(byte(vm.OpPop))
		}
	}
	return nil
}

// compileJSONIndexAccess compiles obj["key"] as GetJSONKey(obj, "key")
func (c *Compiler) compileJSONIndexAccess(node *parser.JSONIndexAccess, chunk *vm.Chunk) error {
	if err := c.compileExpression(node.Object, chunk); err != nil {
		return err
	}
	keyIdx := chunk.WriteConstant(node.Key)
	if keyIdx > 255 {
		return fmt.Errorf("too many constants for JSON key")
	}
	chunk.Write(byte(vm.OpLoadConst))
	chunk.Write(byte(keyIdx))
	nameIdx := chunk.WriteConstant("getjsonkey")
	if nameIdx > 255 {
		return fmt.Errorf("too many constants")
	}
	chunk.Write(byte(vm.OpCallForeign))
	chunk.Write(byte(nameIdx))
	chunk.Write(byte(2))
	return nil
}

// compileMemberAccess compiles expr.member: UDT constant group, namespace constant (RL.*), or getter (pos.x).
func (c *Compiler) compileMemberAccess(m *parser.MemberAccess, chunk *vm.Chunk) error {
	mb := strings.ToLower(m.Member)
	if id, ok := m.Object.(*parser.Identifier); ok {
		objLower := strings.ToLower(id.Name)
		// UDT constant group: TypeName.Member (e.g. Color.Red, Key.W)
		if c.typeDefs != nil {
			if td, ok := c.typeDefs[objLower]; ok {
				val, err := c.resolveUDTConstantMember(td, mb)
				if err == nil {
					key := objLower + "." + mb
					if idx, has := c.constIndices[key]; has {
						chunk.Write(byte(vm.OpLoadConst))
						chunk.Write(byte(idx))
						return nil
					}
					ci := chunk.WriteConstant(val)
					if ci > 255 {
						return fmt.Errorf("too many constants")
					}
					c.constIndices[key] = byte(ci)
					chunk.Write(byte(vm.OpLoadConst))
					chunk.Write(byte(ci))
					return nil
				}
			}
		}
		// Namespace constant: RL.DarkGray, BOX2D.Something (0-arg foreign), no object on stack
		if (objLower == "rl" || objLower == "box2d" || objLower == "bullet" || objLower == "game") && mb != "x" && mb != "y" && mb != "z" {
			idx := chunk.WriteConstant(mb)
			if idx > 255 {
				return fmt.Errorf("too many constants for foreign constant")
			}
			chunk.Write(byte(vm.OpCallForeign))
			chunk.Write(byte(idx))
			chunk.Write(byte(0))
			return nil
		}
	}
	// Getter: compile object (e.g. pos or GetMousePosition()), then call getvector2x etc.
	if err := c.compileExpression(m.Object, chunk); err != nil {
		return err
	}
	var name string
	switch mb {
	case "x":
		name = "getvector2x"
	case "y":
		name = "getvector2y"
	case "z":
		name = "getvector3z"
	default:
		name = "getvector2" + mb
	}
	idx := chunk.WriteConstant(name)
	if idx > 255 {
		return fmt.Errorf("too many constants for member getter")
	}
	chunk.Write(byte(vm.OpCallForeign))
	chunk.Write(byte(idx))
	chunk.Write(byte(1))
	return nil
}

// compileNumber compiles a number literal
func (c *Compiler) compileNumber(num *parser.Number, chunk *vm.Chunk) error {
	// If it contains a decimal point, parse as float first (so 0.5 is not parsed as int 0)
	if strings.Contains(num.Value, ".") {
		if floatVal, err := parseFloat(num.Value); err == nil {
			constIndex := chunk.WriteConstant(floatVal)
			chunk.Write(byte(vm.OpLoadConst))
			chunk.Write(byte(constIndex))
			return nil
		}
	}
	// Try integer first
	if intVal, err := parseInt(num.Value); err == nil {
		constIndex := chunk.WriteConstant(intVal)
		chunk.Write(byte(vm.OpLoadConst))
		chunk.Write(byte(constIndex))
		return nil
	}
	// Parse as float
	if floatVal, err := parseFloat(num.Value); err == nil {
		constIndex := chunk.WriteConstant(floatVal)
		chunk.Write(byte(vm.OpLoadConst))
		chunk.Write(byte(constIndex))
		return nil
	}

	return fmt.Errorf("invalid number format: %s", num.Value)
}

// compileString compiles a string literal
func (c *Compiler) compileString(str *parser.StringLiteral, chunk *vm.Chunk) error {
	constIndex := chunk.WriteConstant(str.Value)
	chunk.Write(byte(vm.OpLoadString))
	chunk.Write(byte(constIndex))
	return nil
}

// compileBoolean compiles a boolean literal
func (c *Compiler) compileBoolean(bool *parser.Boolean, chunk *vm.Chunk) error {
	constIndex := chunk.WriteConstant(bool.Value)
	chunk.Write(byte(vm.OpLoadConst))
	chunk.Write(byte(constIndex))
	return nil
}

// compileNilLiteral compiles the null/nil literal (pushes nil onto the stack)
func (c *Compiler) compileNilLiteral(chunk *vm.Chunk) error {
	constIndex := chunk.WriteConstant(nil)
	chunk.Write(byte(vm.OpLoadConst))
	chunk.Write(byte(constIndex))
	return nil
}

// compileIdentifier compiles a variable reference, a CONST, or a qualified "constant" (e.g. RL.DarkGray)
func (c *Compiler) compileIdentifier(ident *parser.Identifier, chunk *vm.Chunk) error {
	// In a function body, params are stack indices 0, 1, ...
	if c.funcParamIndices != nil {
		if idx, ok := c.funcParamIndices[strings.ToLower(ident.Name)]; ok {
			chunk.Write(byte(vm.OpLoadVar))
			chunk.Write(byte(idx))
			return nil
		}
	}
	// Check if it's a global variable
	if varIndex, exists := chunk.GetVariable(ident.Name); exists {
		chunk.Write(byte(vm.OpLoadVar))
		chunk.Write(byte(varIndex))
		return nil
	}

	// Check if it's a CONST (compile-time constant)
	if c.constIndices != nil {
		if idx, ok := c.constIndices[strings.ToLower(ident.Name)]; ok {
			chunk.Write(byte(vm.OpLoadConst))
			chunk.Write(byte(idx))
			return nil
		}
	}

	// Qualified name starting with known namespace (RL., BOX2D., BULLET., GAME.) used as expression -> 0-arg foreign call
	nameLower := strings.ToLower(ident.Name)
	if strings.HasPrefix(nameLower, "rl.") || strings.HasPrefix(nameLower, "box2d.") || strings.HasPrefix(nameLower, "bullet.") || strings.HasPrefix(nameLower, "game.") {
		nameConst := nameLower
		if strings.HasPrefix(nameConst, "rl.") {
			nameConst = nameConst[3:]
		} else if strings.HasPrefix(nameConst, "box2d.") {
			nameConst = nameConst[6:]
		} else 		if strings.HasPrefix(nameConst, "bullet.") {
			nameConst = nameConst[7:]
		} else if strings.HasPrefix(nameConst, "game.") {
			nameConst = nameConst[5:]
		}
		idx := chunk.WriteConstant(nameConst)
		if idx > 255 {
			return fmt.Errorf("too many constants for foreign constant")
		}
		chunk.Write(byte(vm.OpCallForeign))
		chunk.Write(byte(idx))
		chunk.Write(byte(0))
		return nil
	}

	// Load as global
	constIndex := chunk.WriteConstant(ident.Name)
	chunk.Write(byte(vm.OpLoadGlobal))
	chunk.Write(byte(constIndex))
	return nil
}

// compileAssignment compiles an assignment statement (scalar or array element)
func (c *Compiler) compileAssignment(assign *parser.Assignment, chunk *vm.Chunk) error {
	err := c.compileExpression(assign.Value, chunk)
	if err != nil {
		return err
	}

	if len(assign.Indices) > 0 {
		for _, idx := range assign.Indices {
			if err := c.compileExpression(idx, chunk); err != nil {
				return err
			}
		}
		varIndex, exists := chunk.GetVariable(assign.Variable)
		if !exists {
			return fmt.Errorf("array variable not declared: %s", assign.Variable)
		}
		chunk.Write(byte(vm.OpStoreArray))
		chunk.Write(byte(varIndex))
		return nil
	}

	if c.funcParamIndices != nil {
		if idx, ok := c.funcParamIndices[strings.ToLower(assign.Variable)]; ok {
			chunk.Write(byte(vm.OpStoreVar))
			chunk.Write(byte(idx))
			return nil
		}
	}
	if varIndex, exists := chunk.GetVariable(assign.Variable); exists {
		chunk.Write(byte(vm.OpStoreVar))
		chunk.Write(byte(varIndex))
	} else {
		varIndex := chunk.AddVariable(assign.Variable)
		chunk.Write(byte(vm.OpStoreVar))
		chunk.Write(byte(varIndex))
	}
	return nil
}

// compileCompoundAssign compiles +=, -=, *=, /= (load var, load value, op, store var).
func (c *Compiler) compileCompoundAssign(ca *parser.CompoundAssign, chunk *vm.Chunk) error {
	var varIndex int
	if c.funcParamIndices != nil {
		if idx, ok := c.funcParamIndices[strings.ToLower(ca.Variable)]; ok {
			varIndex = idx
			goto emitCompound
		}
	}
	if idx, exists := chunk.GetVariable(ca.Variable); exists {
		varIndex = idx
	} else {
		varIndex = chunk.AddVariable(ca.Variable)
	}
emitCompound:
	// Load current value
	chunk.Write(byte(vm.OpLoadVar))
	chunk.Write(byte(varIndex))
	// Load RHS
	if err := c.compileExpression(ca.Value, chunk); err != nil {
		return err
	}
	switch ca.Op {
	case "+=":
		chunk.Write(byte(vm.OpAdd))
	case "-=":
		chunk.Write(byte(vm.OpSub))
	case "*=":
		chunk.Write(byte(vm.OpMul))
	case "/=":
		chunk.Write(byte(vm.OpDiv))
	default:
		return fmt.Errorf("unsupported compound assign op: %s", ca.Op)
	}
	chunk.Write(byte(vm.OpStoreVar))
	chunk.Write(byte(varIndex))
	return nil
}

// compileBinaryOp compiles a binary operation
func (c *Compiler) compileBinaryOp(op *parser.BinaryOp, chunk *vm.Chunk) error {
	// Compile left operand
	err := c.compileExpression(op.Left, chunk)
	if err != nil {
		return err
	}

	// Compile right operand
	err = c.compileExpression(op.Right, chunk)
	if err != nil {
		return err
	}

	// Emit operation instruction (case-insensitive for AND/OR)
	opCanon := strings.ToLower(op.Operator)
	switch opCanon {
	case "+":
		chunk.Write(byte(vm.OpAdd))
	case "-":
		chunk.Write(byte(vm.OpSub))
	case "*":
		chunk.Write(byte(vm.OpMul))
	case "/":
		chunk.Write(byte(vm.OpDiv))
	case "%":
		chunk.Write(byte(vm.OpMod))
	case "^":
		chunk.Write(byte(vm.OpPower))
	case "\\":
		chunk.Write(byte(vm.OpIntDiv))
	case "=", "==":
		chunk.Write(byte(vm.OpEqual))
	case "<>":
		chunk.Write(byte(vm.OpNotEqual))
	case "<":
		chunk.Write(byte(vm.OpLess))
	case "<=":
		chunk.Write(byte(vm.OpLessEqual))
	case ">":
		chunk.Write(byte(vm.OpGreater))
	case ">=":
		chunk.Write(byte(vm.OpGreaterEqual))
	case "and":
		chunk.Write(byte(vm.OpAnd))
	case "or":
		chunk.Write(byte(vm.OpOr))
	case "xor":
		chunk.Write(byte(vm.OpXor))
	default:
		return fmt.Errorf("unsupported binary operator: %s", op.Operator)
	}

	return nil
}

// compileUnaryOp compiles a unary operation
func (c *Compiler) compileUnaryOp(op *parser.UnaryOp, chunk *vm.Chunk) error {
	// Compile operand
	err := c.compileExpression(op.Operand, chunk)
	if err != nil {
		return err
	}

	// Emit operation instruction (case-insensitive for NOT)
	switch {
	case op.Operator == "-":
		chunk.Write(byte(vm.OpNeg))
	case strings.EqualFold(op.Operator, "NOT"):
		chunk.Write(byte(vm.OpNot))
	default:
		return fmt.Errorf("unsupported unary operator: %s", op.Operator)
	}

	return nil
}

// compileCall compiles a function/procedure call or array element read
func (c *Compiler) compileCall(call *parser.Call, chunk *vm.Chunk) error {
	// Built-in: ShouldClose() pushes runtime.ShouldClose() result (case-insensitive)
	if len(call.Arguments) == 0 && strings.EqualFold(call.Name, "shouldclose") {
		chunk.Write(byte(vm.OpShouldClose))
		return nil
	}

	// Array element read: a(i,j) when a is an array
	if !strings.Contains(call.Name, ".") {
		if dims, ok := chunk.GetVarDims(call.Name); ok && len(dims) == len(call.Arguments) {
			for _, arg := range call.Arguments {
				if err := c.compileExpression(arg, chunk); err != nil {
					return err
				}
			}
			varIndex, _ := chunk.GetVariable(call.Name)
			chunk.Write(byte(vm.OpLoadArray))
			chunk.Write(byte(varIndex))
			return nil
		}
	}

	// Dotted name: user module function (Math3D.Dot) or foreign (RL.InitWindow, BOX2D.Step)
	if strings.Contains(call.Name, ".") {
		for _, arg := range call.Arguments {
			if err := c.compileExpression(arg, chunk); err != nil {
				return err
			}
		}
		nameConst := strings.ToLower(call.Name)
		if c.userFuncs != nil && c.userFuncs[nameConst] {
			idx := chunk.WriteConstant(nameConst)
			if idx > 255 {
				return fmt.Errorf("too many constants for user call")
			}
			chunk.Write(byte(vm.OpCallUser))
			chunk.Write(byte(idx))
			chunk.Write(byte(len(call.Arguments)))
			return nil
		}
		if strings.HasPrefix(nameConst, "rl.") {
			nameConst = nameConst[3:] // "rl.initwindow" -> "initwindow"
		}
		idx := chunk.WriteConstant(nameConst)
		if idx > 255 {
			return fmt.Errorf("too many constants for foreign call")
		}
		chunk.Write(byte(vm.OpCallForeign))
		chunk.Write(byte(idx))
		chunk.Write(byte(len(call.Arguments)))
		return nil
	}

	// Built-in functions (case-insensitive) - compile args per call as needed
	name := strings.ToLower(call.Name)
	switch name {
	case "print":
		for _, arg := range call.Arguments {
			if err := c.compileExpression(arg, chunk); err != nil {
				return err
			}
			chunk.Write(byte(vm.OpPrint))
		}
		return nil
	case "matmul":
		if len(call.Arguments) != 3 {
			return fmt.Errorf("MatMul(resultName, aName, bName) expects 3 arguments")
		}
		getName := func(n parser.Node) (string, bool) {
			switch v := n.(type) {
			case *parser.Identifier:
				return v.Name, true
			case *parser.StringLiteral:
				return v.Value, true
			default:
				return "", false
			}
		}
		r, ok1 := getName(call.Arguments[0])
		a, ok2 := getName(call.Arguments[1])
		b, ok3 := getName(call.Arguments[2])
		if !ok1 || !ok2 || !ok3 {
			return fmt.Errorf("MatMul: all 3 args must be variable names or string literals")
		}
		ri := chunk.WriteConstant(strings.ToLower(r))
		ai := chunk.WriteConstant(strings.ToLower(a))
		bi := chunk.WriteConstant(strings.ToLower(b))
		if ri > 255 || ai > 255 || bi > 255 {
			return fmt.Errorf("too many constants")
		}
		chunk.Write(byte(vm.OpMatMul))
		chunk.Write(byte(ri))
		chunk.Write(byte(ai))
		chunk.Write(byte(bi))
		return nil
	}

	// Compile arguments for remaining calls
	for _, arg := range call.Arguments {
		err := c.compileExpression(arg, chunk)
		if err != nil {
			return err
		}
	}

	// User Sub/Function call (name resolved from AST; chunk.Functions filled when compiling decls)
	if c.userFuncs != nil && c.userFuncs[name] {
		idx := chunk.WriteConstant(name)
		if idx > 255 {
			return fmt.Errorf("too many constants for user call")
		}
		chunk.Write(byte(vm.OpCallUser))
		chunk.Write(byte(idx))
		chunk.Write(byte(len(call.Arguments)))
		return nil
	}

	switch name {
	case "str":
		if len(call.Arguments) != 1 {
			return fmt.Errorf("STR() expects 1 argument")
		}
		chunk.Write(byte(vm.OpStr))
		return nil
	case "random":
		if len(call.Arguments) == 0 {
			chunk.Write(byte(vm.OpRandom))
			return nil
		}
		if len(call.Arguments) == 1 {
			chunk.Write(byte(vm.OpRandomN))
			return nil
		}
		return fmt.Errorf("Random() expects 0 or 1 argument")
	case "sleep", "wait":
		if len(call.Arguments) != 1 {
			return fmt.Errorf("Sleep/Wait expect 1 argument (milliseconds)")
		}
		chunk.Write(byte(vm.OpSleep))
		return nil
	case "int":
		if len(call.Arguments) != 1 {
			return fmt.Errorf("Int() expects 1 argument")
		}
		chunk.Write(byte(vm.OpInt))
		return nil
	case "timer":
		if len(call.Arguments) != 0 {
			return fmt.Errorf("Timer() takes no arguments")
		}
		chunk.Write(byte(vm.OpTimer))
		return nil
	case "resettimer":
		if len(call.Arguments) != 0 {
			return fmt.Errorf("ResetTimer() takes no arguments")
		}
		chunk.Write(byte(vm.OpResetTimer))
		return nil
	case "quit":
		if len(call.Arguments) != 0 {
			return fmt.Errorf("Quit takes no arguments")
		}
		chunk.Write(byte(vm.OpQuit))
		return nil
	case "sin":
		if len(call.Arguments) != 1 {
			return fmt.Errorf("Sin() expects 1 argument")
		}
		chunk.Write(byte(vm.OpSin))
		return nil
	case "cos":
		if len(call.Arguments) != 1 {
			return fmt.Errorf("Cos() expects 1 argument")
		}
		chunk.Write(byte(vm.OpCos))
		return nil
	case "tan":
		if len(call.Arguments) != 1 {
			return fmt.Errorf("Tan() expects 1 argument")
		}
		chunk.Write(byte(vm.OpTan))
		return nil
	case "sqrt":
		if len(call.Arguments) != 1 {
			return fmt.Errorf("Sqrt() expects 1 argument")
		}
		chunk.Write(byte(vm.OpSqrt))
		return nil
	case "abs":
		if len(call.Arguments) != 1 {
			return fmt.Errorf("Abs() expects 1 argument")
		}
		chunk.Write(byte(vm.OpAbs))
		return nil
	case "lerp":
		if len(call.Arguments) != 3 {
			return fmt.Errorf("Lerp(a, b, t) expects 3 arguments")
		}
		chunk.Write(byte(vm.OpLerp))
		return nil
	case "noise", "noise2d", "perlin", "simplex":
		if len(call.Arguments) != 2 {
			return fmt.Errorf("Noise(x, y) expects 2 arguments")
		}
		chunk.Write(byte(vm.OpNoise2D))
		return nil
	case "openfile":
		if len(call.Arguments) != 2 {
			return fmt.Errorf("OpenFile(path, mode) expects 2 arguments; mode 0=read, 1=write, 2=append")
		}
		chunk.Write(byte(vm.OpOpenFile))
		return nil
	case "readline":
		if len(call.Arguments) != 1 {
			return fmt.Errorf("ReadLine(handle) expects 1 argument")
		}
		chunk.Write(byte(vm.OpReadLine))
		return nil
	case "writeline":
		if len(call.Arguments) != 2 {
			return fmt.Errorf("WriteLine(handle, text) expects 2 arguments")
		}
		chunk.Write(byte(vm.OpWriteLine))
		return nil
	case "closefile":
		if len(call.Arguments) != 1 {
			return fmt.Errorf("CloseFile(handle) expects 1 argument")
		}
		chunk.Write(byte(vm.OpCloseFile))
		return nil
	case "floor":
		if len(call.Arguments) != 1 {
			return fmt.Errorf("Floor() expects 1 argument")
		}
		chunk.Write(byte(vm.OpFloor))
		return nil
	case "ceil":
		if len(call.Arguments) != 1 {
			return fmt.Errorf("Ceil() expects 1 argument")
		}
		chunk.Write(byte(vm.OpCeil))
		return nil
	case "round":
		if len(call.Arguments) != 1 {
			return fmt.Errorf("Round() expects 1 argument")
		}
		chunk.Write(byte(vm.OpRound))
		return nil
	case "min":
		if len(call.Arguments) != 2 {
			return fmt.Errorf("Min(a, b) expects 2 arguments")
		}
		chunk.Write(byte(vm.OpMin))
		return nil
	case "max":
		if len(call.Arguments) != 2 {
			return fmt.Errorf("Max(a, b) expects 2 arguments")
		}
		chunk.Write(byte(vm.OpMax))
		return nil
	case "clamp":
		if len(call.Arguments) != 3 {
			return fmt.Errorf("Clamp(x, lo, hi) expects 3 arguments")
		}
		chunk.Write(byte(vm.OpClamp))
		return nil
	case "pow":
		if len(call.Arguments) != 2 {
			return fmt.Errorf("Pow(base, exp) expects 2 arguments")
		}
		chunk.Write(byte(vm.OpPow))
		return nil
	case "exp":
		if len(call.Arguments) != 1 {
			return fmt.Errorf("Exp() expects 1 argument")
		}
		chunk.Write(byte(vm.OpExp))
		return nil
	case "log":
		if len(call.Arguments) != 1 {
			return fmt.Errorf("Log() expects 1 argument")
		}
		chunk.Write(byte(vm.OpLog))
		return nil
	case "log10":
		if len(call.Arguments) != 1 {
			return fmt.Errorf("Log10() expects 1 argument")
		}
		chunk.Write(byte(vm.OpLog10))
		return nil
	case "atan2":
		if len(call.Arguments) != 2 {
			return fmt.Errorf("Atan2(y, x) expects 2 arguments")
		}
		chunk.Write(byte(vm.OpAtan2))
		return nil
	case "sign":
		if len(call.Arguments) != 1 {
			return fmt.Errorf("Sign() expects 1 argument")
		}
		chunk.Write(byte(vm.OpSign))
		return nil
	case "deg2rad":
		if len(call.Arguments) != 1 {
			return fmt.Errorf("Deg2Rad() expects 1 argument")
		}
		chunk.Write(byte(vm.OpDeg2Rad))
		return nil
	case "rad2deg":
		if len(call.Arguments) != 1 {
			return fmt.Errorf("Rad2Deg() expects 1 argument")
		}
		chunk.Write(byte(vm.OpRad2Deg))
		return nil
	case "distance2d":
		if len(call.Arguments) != 4 {
			return fmt.Errorf("Distance2D(x1, y1, x2, y2) expects 4 arguments")
		}
		chunk.Write(byte(vm.OpDistance2D))
		return nil
	case "distance3d":
		if len(call.Arguments) != 6 {
			return fmt.Errorf("Distance3D(x1, y1, z1, x2, y2, z2) expects 6 arguments")
		}
		chunk.Write(byte(vm.OpDistance3D))
		return nil
	case "distsq2d":
		if len(call.Arguments) != 4 {
			return fmt.Errorf("DistSq2D(x1, y1, x2, y2) expects 4 arguments")
		}
		chunk.Write(byte(vm.OpDistSq2D))
		return nil
	case "distsq3d":
		if len(call.Arguments) != 6 {
			return fmt.Errorf("DistSq3D(x1, y1, z1, x2, y2, z2) expects 6 arguments")
		}
		chunk.Write(byte(vm.OpDistSq3D))
		return nil
	case "inradius2d":
		if len(call.Arguments) != 5 {
			return fmt.Errorf("InRadius2D(x1, y1, x2, y2, radius) expects 5 arguments")
		}
		chunk.Write(byte(vm.OpInRadius2D))
		return nil
	case "inradius3d":
		if len(call.Arguments) != 7 {
			return fmt.Errorf("InRadius3D(x1, y1, z1, x2, y2, z2, radius) expects 7 arguments")
		}
		chunk.Write(byte(vm.OpInRadius3D))
		return nil
	case "angle2d":
		if len(call.Arguments) != 4 {
			return fmt.Errorf("Angle2D(x1, y1, x2, y2) expects 4 arguments")
		}
		chunk.Write(byte(vm.OpAngle2D))
		return nil
	case "left", "left$":
		if len(call.Arguments) != 2 {
			return fmt.Errorf("Left(s, n) expects 2 arguments")
		}
		chunk.Write(byte(vm.OpLeftStr))
		return nil
	case "right", "right$":
		if len(call.Arguments) != 2 {
			return fmt.Errorf("Right(s, n) expects 2 arguments")
		}
		chunk.Write(byte(vm.OpRightStr))
		return nil
	case "mid", "mid$":
		if len(call.Arguments) != 3 {
			return fmt.Errorf("Mid(s, start, n) expects 3 arguments; start is 1-based")
		}
		chunk.Write(byte(vm.OpMidStr))
		return nil
	case "len":
		if len(call.Arguments) != 1 {
			return fmt.Errorf("Len(s) expects 1 argument")
		}
		chunk.Write(byte(vm.OpLenStr))
		return nil
	case "eof":
		if len(call.Arguments) != 1 {
			return fmt.Errorf("EOF(handle) expects 1 argument")
		}
		chunk.Write(byte(vm.OpEOF))
		return nil
	}

	// Unrecognized name: treat as foreign API (InitWindow, BULLET.Step, etc.).
	// Names without a dot (e.g. InitWindow) are raylib; with a dot (e.g. BULLET.CreateWorld) we already handled above.
	nameConst := strings.ToLower(call.Name)
	if strings.HasPrefix(nameConst, "rl.") {
		nameConst = nameConst[3:]
	}
	idx := chunk.WriteConstant(nameConst)
	if idx > 255 {
		return fmt.Errorf("too many constants for foreign call")
	}
	chunk.Write(byte(vm.OpCallForeign))
	chunk.Write(byte(idx))
	chunk.Write(byte(len(call.Arguments)))
	return nil
}

// compileSelectCaseStatement compiles SELECT CASE expr ... CASE val: block ... END SELECT
func (c *Compiler) compileSelectCaseStatement(s *parser.SelectCaseStatement, chunk *vm.Chunk) error {
	if err := c.compileExpression(s.Expr, chunk); err != nil {
		return err
	}
	var endJumpPositions []int
	for _, k := range s.Cases {
		chunk.Write(byte(vm.OpDup))
		if err := c.compileExpression(k.Value, chunk); err != nil {
			return err
		}
		chunk.Write(byte(vm.OpEqual))
		chunk.Write(byte(vm.OpJumpIfFalse))
		chunk.Write(byte(0))
		chunk.Write(byte(0))
		skipPos := len(chunk.Code) - 2
		chunk.Write(byte(vm.OpPop))
		chunk.Write(byte(vm.OpPop))
		for _, stmt := range k.Block.Statements {
			if err := c.compileStatement(stmt, chunk); err != nil {
				return err
			}
		}
		chunk.Write(byte(vm.OpJump))
		chunk.Write(byte(0))
		chunk.Write(byte(0))
		endJumpPositions = append(endJumpPositions, len(chunk.Code)-2)
		// Patch skip to here
		chunk.PatchJumpOffset(skipPos, len(chunk.Code)-skipPos-2)
	}
	if s.ElseBlock != nil {
		chunk.Write(byte(vm.OpPop)) // drop the SELECT expr
		for _, stmt := range s.ElseBlock.Statements {
			if err := c.compileStatement(stmt, chunk); err != nil {
				return err
			}
		}
	} else {
		chunk.Write(byte(vm.OpPop))
	}
	endTarget := len(chunk.Code)
	for _, pos := range endJumpPositions {
		chunk.PatchJumpOffset(pos, endTarget-pos-2)
	}
	return nil
}

// compileExitLoopStatement compiles EXIT FOR or EXIT WHILE (or BREAK) (jump to end of innermost loop).
func (c *Compiler) compileExitLoopStatement(e *parser.ExitLoopStatement, chunk *vm.Chunk) error {
	if len(c.loopExitStack) == 0 {
		return fmt.Errorf("EXIT/BREAK %s outside loop", e.Kind)
	}
	chunk.Write(byte(vm.OpJump))
	chunk.Write(byte(0))
	chunk.Write(byte(0))
	c.loopExitStack[len(c.loopExitStack)-1] = append(c.loopExitStack[len(c.loopExitStack)-1], len(chunk.Code)-2)
	return nil
}

// compileContinueLoopStatement compiles CONTINUE FOR or CONTINUE WHILE (jump to loop head of innermost matching loop).
func (c *Compiler) compileContinueLoopStatement(cl *parser.ContinueLoopStatement, chunk *vm.Chunk) error {
	if len(c.loopContinueStack) == 0 {
		return fmt.Errorf("CONTINUE %s outside loop", cl.Kind)
	}
	chunk.Write(byte(vm.OpJump))
	chunk.Write(byte(0))
	chunk.Write(byte(0))
	c.loopContinueStack[len(c.loopContinueStack)-1] = append(c.loopContinueStack[len(c.loopContinueStack)-1], len(chunk.Code)-2)
	return nil
}

// compileAssertStatement compiles ASSERT condition [, message] to condition + message + CallForeign Assert(2).
func (c *Compiler) compileAssertStatement(a *parser.AssertStatement, chunk *vm.Chunk) error {
	if err := c.compileExpression(a.Condition, chunk); err != nil {
		return err
	}
	if a.Message != nil {
		if err := c.compileExpression(a.Message, chunk); err != nil {
			return err
		}
	} else {
		idx := chunk.WriteConstant("assertion failed")
		chunk.Write(byte(vm.OpLoadConst))
		chunk.Write(byte(idx))
	}
	idx := chunk.WriteConstant("assert")
	chunk.Write(byte(vm.OpCallForeign))
	chunk.Write(byte(idx))
	chunk.Write(byte(2))
	return nil
}

// compileRepeatStatement compiles REPEAT ... UNTIL condition (jump back when condition false)
func (c *Compiler) compileRepeatStatement(r *parser.RepeatStatement, chunk *vm.Chunk) error {
	loopStart := len(chunk.Code)
	wrapFrame := isGameLoopCondition(r.Condition, true)
	hybridMode := wrapFrame && (c.userFuncs["update"] || c.userFuncs["draw"])
	if hybridMode {
		wrapFrame = false
	} else if wrapFrame && r.Body != nil && c.bodyCallsUserSub(r.Body.Statements) {
		wrapFrame = false
	}
	use3D := wrapFrame && r.Body != nil && bodyContains3DDraw(r.Body.Statements)
	if hybridMode {
		c.emitHybridLoopBody(chunk)
	} else {
		if wrapFrame {
			c.emitFrameWrap(chunk, "BeginDrawing")
			if use3D {
				c.emitFrameWrap(chunk, "BeginMode3D")
			} else {
				c.emitFrameWrap(chunk, "BeginMode2D")
			}
		}
		for _, stmt := range r.Body.Statements {
			if err := c.compileStatement(stmt, chunk); err != nil {
				return err
			}
		}
		if wrapFrame {
			if use3D {
				c.emitFrameWrap(chunk, "EndMode3D")
			} else {
				c.emitFrameWrap(chunk, "EndMode2D")
			}
			c.emitFrameWrap(chunk, "EndDrawing")
		}
	}
	if err := c.compileExpression(r.Condition, chunk); err != nil {
		return err
	}
	chunk.Write(byte(vm.OpJumpIfFalse))
	chunk.Write(byte(0))
	chunk.Write(byte(0))
	offset := loopStart - len(chunk.Code)
	chunk.PatchJumpOffset(len(chunk.Code)-2, offset)
	return nil
}

// compileIfStatement compiles an IF statement (with optional ELSEIF and ELSE)
func (c *Compiler) compileIfStatement(ifStmt *parser.IfStatement, chunk *vm.Chunk) error {
	// Compile condition
	err := c.compileExpression(ifStmt.Condition, chunk)
	if err != nil {
		return err
	}

	// Emit jump if false (skip then block; go to first ELSEIF or ELSE or end)
	chunk.Write(byte(vm.OpJumpIfFalse))
	chunk.Write(byte(0))
	chunk.Write(byte(0))
	jumpPos := len(chunk.Code) - 2

	// Compile then block
	for _, stmt := range ifStmt.ThenBlock.Statements {
		err = c.compileStatement(stmt, chunk)
		if err != nil {
			return err
		}
	}

	// Jump over all ELSEIF/ELSE to end (patched later)
	chunk.Write(byte(vm.OpJump))
	chunk.Write(byte(0))
	chunk.Write(byte(0))
	endJumpPos := len(chunk.Code) - 2

	// Patch first "jump if false" to here (start of first ELSEIF or ELSE or after end)
	chunk.PatchJumpOffset(jumpPos, len(chunk.Code)-jumpPos-2)

	var endJumpPositions []int
	endJumpPositions = append(endJumpPositions, endJumpPos)

	// Compile each ELSEIF branch
	for _, branch := range ifStmt.ElseIfs {
		err = c.compileExpression(branch.Condition, chunk)
		if err != nil {
			return err
		}
		chunk.Write(byte(vm.OpJumpIfFalse))
		chunk.Write(byte(0))
		chunk.Write(byte(0))
		elseIfJumpPos := len(chunk.Code) - 2

		for _, stmt := range branch.Block.Statements {
			err = c.compileStatement(stmt, chunk)
			if err != nil {
				return err
			}
		}
		chunk.Write(byte(vm.OpJump))
		chunk.Write(byte(0))
		chunk.Write(byte(0))
		elseIfEndPos := len(chunk.Code) - 2
		endJumpPositions = append(endJumpPositions, elseIfEndPos)
		chunk.PatchJumpOffset(elseIfJumpPos, len(chunk.Code)-elseIfJumpPos-2)
	}

	// Optional ELSE block
	if ifStmt.ElseBlock != nil {
		for _, stmt := range ifStmt.ElseBlock.Statements {
			err = c.compileStatement(stmt, chunk)
			if err != nil {
				return err
			}
		}
	}

	// Patch all "jump to end" offsets (after then block and after each ELSEIF block)
	for _, pos := range endJumpPositions {
		chunk.PatchJumpOffset(pos, len(chunk.Code)-pos-2)
	}

	return nil
}

// compileForStatement compiles a FOR loop
func (c *Compiler) compileForStatement(forStmt *parser.ForStatement, chunk *vm.Chunk) error {
	c.loopExitStack = append(c.loopExitStack, nil)
	c.loopContinueStack = append(c.loopContinueStack, nil)
	// Initialize loop variable
	err := c.compileExpression(forStmt.Start, chunk)
	if err != nil {
		return err
	}

	varIndex := chunk.AddVariable(forStmt.Variable)
	chunk.Write(byte(vm.OpStoreVar))
	chunk.Write(byte(varIndex))

	// Loop start
	loopStart := len(chunk.Code)

	// Check condition
	err = c.compileIdentifier(&parser.Identifier{Name: forStmt.Variable}, chunk)
	if err != nil {
		return err
	}

	err = c.compileExpression(forStmt.End, chunk)
	if err != nil {
		return err
	}

	chunk.Write(byte(vm.OpGreater))
	chunk.Write(byte(vm.OpJumpIfTrue))
	chunk.Write(byte(0))
	chunk.Write(byte(0))
	exitJumpPos := len(chunk.Code) - 2

	// Compile loop body
	for _, stmt := range forStmt.Body.Statements {
		err = c.compileStatement(stmt, chunk)
		if err != nil {
			return err
		}
	}

	// CONTINUE FOR jumps here (increment then re-check condition)
	continueTargetIP := len(chunk.Code)
	// Increment loop variable
	err = c.compileIdentifier(&parser.Identifier{Name: forStmt.Variable}, chunk)
	if err != nil {
		return err
	}

	if forStmt.Step != nil {
		err = c.compileExpression(forStmt.Step, chunk)
	} else {
		// Default step is 1
		err = c.compileNumber(&parser.Number{Value: "1"}, chunk)
	}

	if err != nil {
		return err
	}

	chunk.Write(byte(vm.OpAdd))
	chunk.Write(byte(vm.OpStoreVar))
	chunk.Write(byte(varIndex))

	// Jump back to loop start (2-byte offset for large loop bodies)
	chunk.Write(byte(vm.OpJump))
	chunk.Write(byte(0))
	chunk.Write(byte(0))
	chunk.PatchJumpOffset(len(chunk.Code)-2, loopStart-len(chunk.Code))

	// Fix exit jump offset
	chunk.PatchJumpOffset(exitJumpPos, len(chunk.Code)-exitJumpPos-2)
	// Patch EXIT FOR jumps
	for _, pos := range c.loopExitStack[len(c.loopExitStack)-1] {
		chunk.PatchJumpOffset(pos, len(chunk.Code)-pos-2)
	}
	c.loopExitStack = c.loopExitStack[:len(c.loopExitStack)-1]
	// Patch CONTINUE FOR jumps to increment
	for _, pos := range c.loopContinueStack[len(c.loopContinueStack)-1] {
		chunk.PatchJumpOffset(pos, continueTargetIP-pos-2)
	}
	c.loopContinueStack = c.loopContinueStack[:len(c.loopContinueStack)-1]

	return nil
}

func normWindowShouldCloseName(name string) string {
	s := strings.ToLower(name)
	if strings.HasPrefix(s, "rl.") {
		s = s[3:]
	}
	return s
}

// isGameLoopCondition returns true if the condition is NOT WindowShouldClose() (for WHILE) or WindowShouldClose() (for REPEAT).
func isGameLoopCondition(condition parser.Node, forRepeat bool) bool {
	if forRepeat {
		call, ok := condition.(*parser.Call)
		return ok && normWindowShouldCloseName(call.Name) == "windowshouldclose" && len(call.Arguments) == 0
	}
	un, ok := condition.(*parser.UnaryOp)
	if !ok || strings.ToLower(un.Operator) != "not" {
		return false
	}
	call, ok := un.Operand.(*parser.Call)
	return ok && normWindowShouldCloseName(call.Name) == "windowshouldclose" && len(call.Arguments) == 0
}

// bodyCallsUserSub returns true if the given statements (or any nested block) contain a call to a user-defined SUB or FUNCTION.
// When true, we skip automatic BeginDrawing/EndDrawing so the user's own Draw() (or similar) is not double-wrapped and flicker is avoided.
func (c *Compiler) bodyCallsUserSub(statements []parser.Node) bool {
	for _, n := range statements {
		if c.nodeCallsUserSub(n) {
			return true
		}
	}
	return false
}

func (c *Compiler) nodeCallsUserSub(node parser.Node) bool {
	if c.userFuncs == nil {
		return false
	}
	switch n := node.(type) {
	case *parser.Call:
		name := strings.ToLower(n.Name)
		return c.userFuncs[name]
	case *parser.IfStatement:
		if n.ThenBlock != nil && c.bodyCallsUserSub(n.ThenBlock.Statements) {
			return true
		}
		for _, b := range n.ElseIfs {
			if b.Block != nil && c.bodyCallsUserSub(b.Block.Statements) {
				return true
			}
		}
		if n.ElseBlock != nil && c.bodyCallsUserSub(n.ElseBlock.Statements) {
			return true
		}
		return false
	case *parser.ForStatement:
		if n.Body != nil {
			return c.bodyCallsUserSub(n.Body.Statements)
		}
		return false
	case *parser.WhileStatement:
		if n.Body != nil {
			return c.bodyCallsUserSub(n.Body.Statements)
		}
		return false
	case *parser.RepeatStatement:
		if n.Body != nil {
			return c.bodyCallsUserSub(n.Body.Statements)
		}
		return false
	case *parser.SelectCaseStatement:
		for _, k := range n.Cases {
			if k.Block != nil && c.bodyCallsUserSub(k.Block.Statements) {
				return true
			}
		}
		if n.ElseBlock != nil && c.bodyCallsUserSub(n.ElseBlock.Statements) {
			return true
		}
		return false
	default:
		return false
	}
}

// bodyContains3DDraw returns true if the given statements (or any nested block) contain a 3D draw call (DrawCube, DrawSphere, DrawModel, etc.).
func bodyContains3DDraw(nodes []parser.Node) bool {
	for _, n := range nodes {
		if nodeContains3DDraw(n) {
			return true
		}
	}
	return false
}

var threeDDrawNames = map[string]bool{
	"drawcube": true, "drawcubewires": true, "drawsphere": true, "drawspherewires": true,
	"drawmodel": true, "drawmodelsimple": true, "drawmodelex": true, "drawmodelwires": true, "drawplane": true,
	"drawline3d": true, "drawpoint3d": true, "drawcircle3d": true, "drawgrid": true,
	"drawcylinder": true, "drawcylinderwires": true, "drawray": true, "drawtriangle3d": true,
	"beginmode3d": true,
}

func nodeContains3DDraw(node parser.Node) bool {
	switch n := node.(type) {
	case *parser.Call:
		name := normWindowShouldCloseName(n.Name)
		return threeDDrawNames[name]
	case *parser.Statement:
		return n.Value != nil && nodeContains3DDraw(n.Value)
	case *parser.IfStatement:
		if n.ThenBlock != nil && bodyContains3DDraw(n.ThenBlock.Statements) {
			return true
		}
		if n.ElseBlock != nil && bodyContains3DDraw(n.ElseBlock.Statements) {
			return true
		}
		return false
	case *parser.ForStatement:
		if n.Body != nil {
			return bodyContains3DDraw(n.Body.Statements)
		}
		return false
	case *parser.WhileStatement:
		if n.Body != nil {
			return bodyContains3DDraw(n.Body.Statements)
		}
		return false
	case *parser.RepeatStatement:
		if n.Body != nil {
			return bodyContains3DDraw(n.Body.Statements)
		}
		return false
	case *parser.SelectCaseStatement:
		for _, c := range n.Cases {
			if c.Block != nil && bodyContains3DDraw(c.Block.Statements) {
				return true
			}
		}
		if n.ElseBlock != nil && bodyContains3DDraw(n.ElseBlock.Statements) {
			return true
		}
		return false
	default:
		return false
	}
}

// emitFrameWrap emits a no-arg foreign call (BeginDrawing, EndDrawing, BeginMode2D, EndMode2D). Used for automatic frame wrapping.
func (c *Compiler) emitFrameWrap(chunk *vm.Chunk, name string) {
	idx := chunk.WriteConstant(strings.ToLower(name))
	if idx > 255 {
		return
	}
	chunk.Write(byte(vm.OpCallForeign))
	chunk.Write(byte(idx))
	chunk.Write(byte(0))
}

// emitHybridLoopBody emits the hybrid update/draw loop body: GetFrameTime, StepAllPhysics2D/3D, update(dt), ClearRenderQueues, draw(), FlushRenderQueues.
func (c *Compiler) emitHybridLoopBody(chunk *vm.Chunk) {
	emitForeign := func(name string, argCount int) {
		idx := chunk.WriteConstant(strings.ToLower(name))
		if idx > 255 {
			return
		}
		chunk.Write(byte(vm.OpCallForeign))
		chunk.Write(byte(idx))
		chunk.Write(byte(argCount))
	}
	emitCallUser := func(name string, argCount int) {
		idx := chunk.WriteConstant(strings.ToLower(name))
		if idx > 255 {
			return
		}
		chunk.Write(byte(vm.OpCallUser))
		chunk.Write(byte(idx))
		chunk.Write(byte(argCount))
	}
	// dt = GetFrameTime()
	emitForeign("getframetime", 0)
	// StepAllPhysics2D(dt): need dt on stack, then dup for StepAllPhysics3D
	chunk.Write(byte(vm.OpDup))
	emitForeign("stepallphysics2d", 1)
	chunk.Write(byte(vm.OpDup))
	emitForeign("stepallphysics3d", 1)
	if c.userFuncs["update"] {
		emitCallUser("update", 1)
	}
	emitForeign("clearrenderqueues", 0)
	if c.userFuncs["draw"] {
		emitCallUser("draw", 0)
	}
	emitForeign("flushrenderqueues", 0)
}

// compileWhileStatement compiles a WHILE loop
func (c *Compiler) compileWhileStatement(whileStmt *parser.WhileStatement, chunk *vm.Chunk) error {
	c.loopExitStack = append(c.loopExitStack, nil)
	c.loopContinueStack = append(c.loopContinueStack, nil)
	// Loop start (CONTINUE WHILE jumps here)
	loopStart := len(chunk.Code)
	wrapFrame := isGameLoopCondition(whileStmt.Condition, false)
	hybridMode := wrapFrame && (c.userFuncs["update"] || c.userFuncs["draw"])
	if hybridMode {
		wrapFrame = false
	} else if wrapFrame && whileStmt.Body != nil && c.bodyCallsUserSub(whileStmt.Body.Statements) {
		wrapFrame = false // user does their own BeginDrawing/EndDrawing (e.g. in Draw()), avoid double wrap and flicker
	}
	use3D := wrapFrame && whileStmt.Body != nil && bodyContains3DDraw(whileStmt.Body.Statements)

	// Compile condition
	err := c.compileExpression(whileStmt.Condition, chunk)
	if err != nil {
		return err
	}

	// Jump if false (2-byte offset)
	chunk.Write(byte(vm.OpJumpIfFalse))
	chunk.Write(byte(0))
	chunk.Write(byte(0))
	exitJumpPos := len(chunk.Code) - 2

	if hybridMode {
		c.emitHybridLoopBody(chunk)
	} else {
		if wrapFrame {
			c.emitFrameWrap(chunk, "BeginDrawing")
			if use3D {
				c.emitFrameWrap(chunk, "BeginMode3D")
			} else {
				c.emitFrameWrap(chunk, "BeginMode2D")
			}
		}
		for _, stmt := range whileStmt.Body.Statements {
			err = c.compileStatement(stmt, chunk)
			if err != nil {
				return err
			}
		}
		if wrapFrame {
			if use3D {
				c.emitFrameWrap(chunk, "EndMode3D")
			} else {
				c.emitFrameWrap(chunk, "EndMode2D")
			}
			c.emitFrameWrap(chunk, "EndDrawing")
		}
	}

	// Jump back to loop start (2-byte offset)
	chunk.Write(byte(vm.OpJump))
	chunk.Write(byte(0))
	chunk.Write(byte(0))
	chunk.PatchJumpOffset(len(chunk.Code)-2, loopStart-len(chunk.Code))

	// Fix exit jump offset
	chunk.PatchJumpOffset(exitJumpPos, len(chunk.Code)-exitJumpPos-2)
	// Patch EXIT WHILE jumps
	for _, pos := range c.loopExitStack[len(c.loopExitStack)-1] {
		chunk.PatchJumpOffset(pos, len(chunk.Code)-pos-2)
	}
	c.loopExitStack = c.loopExitStack[:len(c.loopExitStack)-1]
	// Patch CONTINUE WHILE jumps to condition
	for _, pos := range c.loopContinueStack[len(c.loopContinueStack)-1] {
		chunk.PatchJumpOffset(pos, loopStart-pos-2)
	}
	c.loopContinueStack = c.loopContinueStack[:len(c.loopContinueStack)-1]

	return nil
}

// compileFunctionDecl compiles a function declaration. VM replaces stack with [arg0, arg1, ...] on call, so no param stores.
func (c *Compiler) compileFunctionDecl(fn *parser.FunctionDecl, chunk *vm.Chunk) error {
	name := qualifiedName(fn)
	chunk.Functions[name] = len(chunk.Code)
	// Params map to stack indices 0, 1, ... so body sees a=0, b=1, etc.
	c.funcParamIndices = make(map[string]int)
	for i, p := range fn.Parameters {
		c.funcParamIndices[strings.ToLower(p)] = i
	}
	for _, stmt := range fn.Body.Statements {
		if err := c.compileStatement(stmt, chunk); err != nil {
			return err
		}
	}
	c.funcParamIndices = nil
	return nil
}

// compileSubDecl compiles a sub procedure declaration. VM replaces stack with [arg0, arg1, ...] on call.
func (c *Compiler) compileSubDecl(sub *parser.SubDecl, chunk *vm.Chunk) error {
	name := qualifiedName(sub)
	chunk.Functions[name] = len(chunk.Code)
	c.funcParamIndices = make(map[string]int)
	for i, p := range sub.Parameters {
		c.funcParamIndices[strings.ToLower(p)] = i
	}
	for _, stmt := range sub.Body.Statements {
		if err := c.compileStatement(stmt, chunk); err != nil {
			return err
		}
	}
	c.funcParamIndices = nil
	chunk.Write(byte(vm.OpReturn))
	return nil
}

// compileStartCoroutineStatement compiles StartCoroutine SubName(): emit OpStartCoroutine with 2-byte target offset (patched after decls)
func (c *Compiler) compileStartCoroutineStatement(stmt *parser.StartCoroutineStatement, chunk *vm.Chunk) error {
	name := strings.ToLower(stmt.SubName)
	if !c.userFuncs[name] {
		return fmt.Errorf("unknown sub for StartCoroutine: %s", stmt.SubName)
	}
	chunk.Write(byte(vm.OpStartCoroutine))
	chunk.WriteJumpOffset(0)
	c.startCoroutinePatchList = append(c.startCoroutinePatchList, startCoroutinePatch{patchPos: len(chunk.Code) - 2, subName: name})
	return nil
}

// compileOnEventStatement compiles On KeyDown("X") ... End On: emit OpRegisterEvent (handler offset patched later)
func (c *Compiler) compileOnEventStatement(on *parser.OnEventStatement, chunk *vm.Chunk) error {
	eventTypeConst := chunk.WriteConstant(strings.ToLower(on.EventType))
	keyConst := chunk.WriteConstant(on.Key)
	if eventTypeConst > 255 || keyConst > 255 {
		return fmt.Errorf("too many constants for On event")
	}
	chunk.Write(byte(vm.OpRegisterEvent))
	chunk.Write(byte(eventTypeConst))
	chunk.Write(byte(keyConst))
	chunk.WriteJumpOffset(0)
	c.eventPatchList = append(c.eventPatchList, eventPatch{patchPos: len(chunk.Code) - 2, stmt: on})
	return nil
}

// compileReturnStatement compiles a RETURN statement
func (c *Compiler) compileReturnStatement(ret *parser.ReturnStatement, chunk *vm.Chunk) error {
	if ret.Value != nil {
		err := c.compileExpression(ret.Value, chunk)
		if err != nil {
			return err
		}
		chunk.Write(byte(vm.OpReturnVal))
	} else {
		chunk.Write(byte(vm.OpReturn))
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
func (c *Compiler) compileDimStatement(dim *parser.DimStatement, chunk *vm.Chunk) error {
	for _, v := range dim.Variables {
		varIndex := chunk.AddVariable(v.Name)

		if len(v.Dimensions) > 0 {
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
				ci := chunk.WriteConstant(n)
				if ci > 255 {
					return fmt.Errorf("too many constants")
				}
				constIndices[i] = byte(ci)
			}
			chunk.SetVarDims(v.Name, dims)
			chunk.Write(byte(vm.OpCreateArray))
			chunk.Write(byte(len(dims)))
			for _, b := range constIndices {
				chunk.Write(b)
			}
			chunk.Write(byte(varIndex))
			continue
		}

		// Scalar: initialize with default value (dynamic type when v.Type == "")
		switch strings.ToLower(v.Type) {
		case "integer", "int", "":
			chunk.Write(byte(vm.OpLoadConst))
			chunk.Write(byte(chunk.WriteConstant(0)))
		case "string", "str":
			chunk.Write(byte(vm.OpLoadString))
			chunk.Write(byte(chunk.WriteConstant("")))
		case "float", "single", "double":
			chunk.Write(byte(vm.OpLoadConst))
			chunk.Write(byte(chunk.WriteConstant(0.0)))
		case "boolean", "bool":
			chunk.Write(byte(vm.OpLoadConst))
			chunk.Write(byte(chunk.WriteConstant(false)))
		case "vector2", "vector3", "body", "color":
			// Optional type hints (e.g. DIM pos AS Vector2); stored as dynamic 0, use .x/.y/.z or API later
			chunk.Write(byte(vm.OpLoadConst))
			chunk.Write(byte(chunk.WriteConstant(0)))
		default:
			chunk.Write(byte(vm.OpLoadConst))
			chunk.Write(byte(chunk.WriteConstant(0)))
		}
		chunk.Write(byte(vm.OpStoreVar))
		chunk.Write(byte(varIndex))
	}
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
func (c *Compiler) compileConstStatement(cs *parser.ConstStatement, chunk *vm.Chunk) error {
	for _, d := range cs.Decls {
		val, err := evalConstExpr(d.Value)
		if err != nil {
			return err
		}
		idx := chunk.WriteConstant(val)
		if idx > 255 {
			return fmt.Errorf("too many constants")
		}
		c.constIndices[strings.ToLower(d.Name)] = byte(idx)
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
func (c *Compiler) resolveUDTConstantMember(td *parser.TypeDecl, memberLower string) (interface{}, error) {
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
	return nil, fmt.Errorf("unknown member %s", memberLower)
}

// compileEnumStatement compiles ENUM Name : a, b = 2, c ... members as constants (auto-increment from 0 or explicit value).
// Also records enum name -> member -> value in chunk.Enums for Enum.getValue/getName/hasValue at runtime.
func (c *Compiler) compileEnumStatement(es *parser.EnumStatement, chunk *vm.Chunk) error {
	enumName := strings.ToLower(es.Name)
	if chunk.Enums == nil {
		chunk.Enums = make(map[string]vm.EnumMembers)
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
		idx := chunk.WriteConstant(nextVal)
		if idx > 255 {
			return fmt.Errorf("too many constants")
		}
		c.constIndices[memLower] = byte(idx)
		nextVal++
	}
	chunk.Enums[enumName] = members
	return nil
}

// compileGameCommand compiles game-specific commands
func (c *Compiler) compileGameCommand(cmd *parser.GameCommand, chunk *vm.Chunk) error {
	// Compile arguments
	for _, arg := range cmd.Arguments {
		err := c.compileExpression(arg, chunk)
		if err != nil {
			return err
		}
	}

	// Case-insensitive: canonical form is lowercase
	command := strings.ToLower(cmd.Command)

	// Emit game command instruction
	switch command {
	case "print":
		chunk.Write(byte(vm.OpPrint))
	case "str":
		chunk.Write(byte(vm.OpStr))
	case "loadimage":
		chunk.Write(byte(vm.OpLoadImage))
	case "createsprite":
		chunk.Write(byte(vm.OpCreateSprite))
	case "setspriteposition":
		chunk.Write(byte(vm.OpSetSpritePosition))
	case "drawsprite":
		chunk.Write(byte(vm.OpDrawSprite))
	case "loadmodel":
		chunk.Write(byte(vm.OpLoadModel))
	case "createcamera":
		chunk.Write(byte(vm.OpCreateCamera))
	case "setcameraposition":
		chunk.Write(byte(vm.OpSetCameraPosition))
	case "drawmodel":
		chunk.Write(byte(vm.OpDrawModel))
	case "playmusic":
		chunk.Write(byte(vm.OpPlayMusic))
	case "playsound":
		chunk.Write(byte(vm.OpPlaySound))
	case "loadsound":
		chunk.Write(byte(vm.OpLoadSound))
	case "createphysicsbody":
		chunk.Write(byte(vm.OpCreatePhysicsBody))
	case "setvelocity":
		chunk.Write(byte(vm.OpSetVelocity))
	case "applyforce":
		chunk.Write(byte(vm.OpApplyForce))
	case "raycast3d":
		chunk.Write(byte(vm.OpRayCast3D))
	case "sync":
		chunk.Write(byte(vm.OpSync))
	case "shouldclose":
		chunk.Write(byte(vm.OpShouldClose))
	default:
		return fmt.Errorf("unsupported game command: %s", cmd.Command)
	}

	return nil
}

// Helper functions for parsing numbers
func parseInt(s string) (int, error) {
	var result int
	_, err := fmt.Sscanf(s, "%d", &result)
	return result, err
}

func parseFloat(s string) (float64, error) {
	var result float64
	_, err := fmt.Sscanf(s, "%f", &result)
	return result, err
}
