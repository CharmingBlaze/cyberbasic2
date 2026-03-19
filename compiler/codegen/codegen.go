package codegen

import (
	"cyberbasic/compiler/parser"
	"cyberbasic/compiler/semantic"
	"cyberbasic/compiler/vm"
	"fmt"
)

// eventPatch records a position to patch with the OnEvent handler offset.
type eventPatch struct {
	patchPos int
	stmt     *parser.OnEventStatement
}

// startCoroutinePatch records a position to patch with the sub's bytecode offset after decls are compiled.
type startCoroutinePatch struct {
	patchPos int
	subName  string
}

// Emitter holds mutable codegen state and the read-only semantic result.
type Emitter struct {
	chunk                   *vm.Chunk
	sem                     *semantic.Result
	constIndices            map[string]byte
	loopExitStack           [][]int
	loopContinueStack       [][]int
	funcParamIndices        map[string]int
	eventPatchList          []eventPatch
	startCoroutinePatchList []startCoroutinePatch
}

// Emit compiles the program AST into bytecode using the semantic analysis result.
func Emit(program *parser.Program, sem *semantic.Result) (*vm.Chunk, error) {
	chunk := vm.NewChunk()
	e := &Emitter{
		chunk:                   chunk,
		sem:                     sem,
		constIndices:            make(map[string]byte),
		eventPatchList:          nil,
		startCoroutinePatchList: nil,
	}
	// Compile main program (no Function/Sub bodies)
	for _, stmt := range sem.MainStmts {
		if err := e.compileStatement(stmt); err != nil {
			return nil, err
		}
	}
	// Jump over all function/sub bodies
	jumpPos := len(chunk.Code)
	chunk.Write(byte(vm.OpJump))
	chunk.WriteJumpOffset(0)
	// Compile each Sub/Function
	for _, stmt := range sem.Decls {
		if err := e.compileDecl(stmt); err != nil {
			return nil, err
		}
	}
	// Patch StartCoroutine target offsets
	for _, p := range e.startCoroutinePatchList {
		target, ok := chunk.GetFunction(p.subName)
		if !ok {
			candidates := make([]string, 0, len(e.sem.UserFuncs))
			for q := range e.sem.UserFuncs {
				candidates = append(candidates, q)
			}
			msg := fmt.Sprintf("unknown sub for StartCoroutine: %s", p.subName)
			if sug := nearestName(p.subName, candidates, 3); sug != "" {
				msg += " (did you mean " + sug + "?)"
			}
			return nil, fmt.Errorf("%s", msg)
		}
		chunk.PatchJumpOffset(p.patchPos, target)
	}
	// Compile event handlers and patch registration offsets
	for _, ep := range e.eventPatchList {
		handlerStart := len(chunk.Code)
		for _, stmt := range ep.stmt.Body.Statements {
			if err := e.compileStatement(stmt); err != nil {
				return nil, err
			}
		}
		chunk.Write(byte(vm.OpReturn))
		chunk.PatchJumpOffset(ep.patchPos, handlerStart)
	}
	endPos := len(chunk.Code)
	chunk.PatchJumpOffset(jumpPos+1, endPos-(jumpPos+3))
	chunk.Write(byte(vm.OpHalt))
	return chunk, nil
}

// compileDecl compiles a single FunctionDecl or SubDecl.
func (e *Emitter) compileDecl(stmt parser.Node) error {
	switch node := stmt.(type) {
	case *parser.FunctionDecl:
		return e.compileFunctionDecl(node)
	case *parser.SubDecl:
		return e.compileSubDecl(node)
	default:
		return fmt.Errorf("compileDecl: expected FunctionDecl or SubDecl, got %T", stmt)
	}
}
