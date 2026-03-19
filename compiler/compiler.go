// Package compiler is the single front-end for CyberBasic: one pipeline (lex → parse → semantic → codegen)
// with modular subpackages (lexer, parser, semantic, codegen). Use Compiler.Compile or CompileWithOptions
// for end-to-end builds; use Tokenize, Parse, Analyze for tooling that stops early. All full compilations
// share the same internal pipeline so behavior cannot drift between callers.
package compiler

import (
	"cyberbasic/compiler/codegen"
	"cyberbasic/compiler/lexer"
	"cyberbasic/compiler/parser"
	"cyberbasic/compiler/semantic"
	"cyberbasic/compiler/vm"
	"fmt"
)

// Compiler is the driver for the compilation pipeline: Lex → Parse → Semantic → Codegen.
// Set Filename before staged calls or use CompileOptions for per-invocation source paths in errors.
type Compiler struct {
	Filename string // optional; used when wrapping errors (e.g. "game.bas: line 5: ...")
}

// CompileOptions configures a single CompileWithOptions invocation.
type CompileOptions struct {
	// Filename is prepended to phase errors (lexical, parse, semantic, codegen). Empty uses Compiler.Filename.
	Filename string
}

// New creates a new compiler instance.
func New() *Compiler {
	return &Compiler{}
}

// wrapFilenameErr prefixes an error with phase and optional source file name.
func wrapFilenameErr(filename, phase string, err error) error {
	if err == nil {
		return nil
	}
	e := fmt.Errorf("%s: %w", phase, err)
	if filename != "" {
		e = fmt.Errorf("%s: %w", filename, e)
	}
	return e
}

func (c *Compiler) effectiveFilename(opts *CompileOptions) string {
	if opts != nil && opts.Filename != "" {
		return opts.Filename
	}
	return c.Filename
}

func (c *Compiler) wrapErr(phase string, err error) error {
	return wrapFilenameErr(c.Filename, phase, err)
}

// Tokenize runs the lexer only. Errors are wrapped with Compiler.Filename when set.
func (c *Compiler) Tokenize(source string) ([]lexer.Token, error) {
	l := lexer.New(source)
	tokens, err := l.Tokenize()
	if err != nil {
		return nil, c.wrapErr("lexical error", err)
	}
	return tokens, nil
}

// ParseTokens builds an AST from an existing token stream (tests, tooling, incremental workflows).
func (c *Compiler) ParseTokens(tokens []lexer.Token) (*parser.Program, error) {
	p := parser.New(tokens)
	program, err := p.Parse()
	if err != nil {
		return nil, c.wrapErr("parse error", err)
	}
	return program, nil
}

// Parse runs Lex → Parse. Equivalent to Tokenize followed by ParseTokens.
func (c *Compiler) Parse(source string) (*parser.Program, error) {
	tokens, err := c.Tokenize(source)
	if err != nil {
		return nil, err
	}
	return c.ParseTokens(tokens)
}

// Analyze runs semantic analysis on an AST. Does not modify the AST.
func (c *Compiler) Analyze(program *parser.Program) (*semantic.Result, error) {
	semResult, err := semantic.Analyze(program)
	if err != nil {
		return nil, c.wrapErr("semantic error", err)
	}
	return semResult, nil
}

// Compile compiles BASIC source code to bytecode using Compiler.Filename for diagnostics.
func (c *Compiler) Compile(source string) (*vm.Chunk, error) {
	return c.fullPipeline(source, c.Filename)
}

// CompileWithOptions compiles source to bytecode. opts.Filename overrides Compiler.Filename for error prefixes when non-empty.
func (c *Compiler) CompileWithOptions(source string, opts CompileOptions) (*vm.Chunk, error) {
	return c.fullPipeline(source, c.effectiveFilename(&opts))
}

// fullPipeline is the only place that chains lexer → parser → semantic → codegen for a complete build.
// filename is used only for error messages (may be empty).
func (c *Compiler) fullPipeline(source, filename string) (*vm.Chunk, error) {
	tokens, err := lexer.New(source).Tokenize()
	if err != nil {
		return nil, wrapFilenameErr(filename, "lexical error", err)
	}
	p := parser.New(tokens)
	program, err := p.Parse()
	if err != nil {
		return nil, wrapFilenameErr(filename, "parse error", err)
	}
	semResult, err := semantic.Analyze(program)
	if err != nil {
		return nil, wrapFilenameErr(filename, "semantic error", err)
	}
	chunk, err := codegen.Emit(program, semResult)
	if err != nil {
		return nil, wrapFilenameErr(filename, "code generation error", err)
	}
	return chunk, nil
}
