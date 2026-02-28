package compiler

import (
	"cyberbasic/compiler/lexer"
	"cyberbasic/compiler/parser"
	"cyberbasic/compiler/vm"
	"fmt"
	"strings"
)

// Compiler handles the compilation of BASIC source code to bytecode.
// Codegen uses compile* for "compile this AST to bytecode" and emit* for "emit a fixed opcode sequence" (e.g. emitFrameWrap).
type Compiler struct {
	lexer         *lexer.Lexer
	parser        *parser.Parser
	vm            *vm.VM
	constIndices  map[string]byte            // const name -> chunk constant index (set during generateCode)
	typeDefs      map[string]*parser.TypeDecl // UDT name (lowercase) -> TYPE definition (filled in first pass)
	entityNames   map[string]bool             // entity name (lowercase) from each ENTITY decl (first pass)
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
	c.entityNames = make(map[string]bool)
	c.eventPatchList = nil
	c.startCoroutinePatchList = nil
	// First pass: collect TYPE definitions, entity names, and user function/sub names
	userFuncs := make(map[string]bool)
	var mainStmts, decls []parser.Node
	for _, stmt := range program.Statements {
		if td, ok := stmt.(*parser.TypeDecl); ok {
			c.typeDefs[strings.ToLower(td.Name)] = td
		}
		if ed, ok := stmt.(*parser.EntityDecl); ok {
			c.entityNames[strings.ToLower(ed.Name)] = true
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

// errWithLine wraps err with "line N: " when node has a source line, for clearer compiler errors.
func errWithLine(node parser.Node, err error) error {
	if err == nil {
		return nil
	}
	if line := getSourceLine(node); line > 0 {
		return fmt.Errorf("line %d: %w", line, err)
	}
	return err
}
