package compiler

import (
	"cyberbasic/compiler/codegen"
	"cyberbasic/compiler/lexer"
	"cyberbasic/compiler/parser"
	"cyberbasic/compiler/semantic"
	"cyberbasic/compiler/vm"
	"fmt"
)

// Compiler is the thin driver for the compilation pipeline: Lex → Parse → Semantic → Codegen.
// Set Filename before Compile to have errors prefixed with the source file name (e.g. "game.bas: line 5: ...").
type Compiler struct {
	Filename string // optional; used when wrapping errors for clearer diagnostics
}

// New creates a new compiler instance.
func New() *Compiler {
	return &Compiler{}
}

func (c *Compiler) wrapErr(phase string, err error) error {
	if err == nil {
		return nil
	}
	e := fmt.Errorf("%s: %w", phase, err)
	if c.Filename != "" {
		e = fmt.Errorf("%s: %w", c.Filename, e)
	}
	return e
}

// Compile compiles BASIC source code to bytecode.
func (c *Compiler) Compile(source string) (*vm.Chunk, error) {
	l := lexer.New(source)
	tokens, err := l.Tokenize()
	if err != nil {
		return nil, c.wrapErr("lexical error", err)
	}
	p := parser.New(tokens)
	program, err := p.Parse()
	if err != nil {
		return nil, c.wrapErr("parse error", err)
	}
	semResult, err := semantic.Analyze(program)
	if err != nil {
		return nil, c.wrapErr("semantic error", err)
	}
	chunk, err := codegen.Emit(program, semResult)
	if err != nil {
		return nil, c.wrapErr("code generation error", err)
	}
	return chunk, nil
}
