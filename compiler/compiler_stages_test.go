package compiler

import (
	"cyberbasic/compiler/codegen"
	"cyberbasic/compiler/vm"
	"strings"
	"testing"
)

func TestTokenizeNonEmpty(t *testing.T) {
	c := New()
	toks, err := c.Tokenize(`VAR x = 1`)
	if err != nil {
		t.Fatal(err)
	}
	if len(toks) == 0 {
		t.Fatal("expected tokens")
	}
}

func TestParseTokensRoundTrip(t *testing.T) {
	c := New()
	src := `Print(42)`
	toks, err := c.Tokenize(src)
	if err != nil {
		t.Fatal(err)
	}
	p1, err := c.ParseTokens(toks)
	if err != nil {
		t.Fatal(err)
	}
	p2, err := c.Parse(src)
	if err != nil {
		t.Fatal(err)
	}
	if len(p1.Statements) != len(p2.Statements) {
		t.Fatalf("statement count mismatch: %d vs %d", len(p1.Statements), len(p2.Statements))
	}
}

func TestAnalyzeCollectsUserFunc(t *testing.T) {
	c := New()
	src := `Function F()
  Return 1
End Function
`
	p, err := c.Parse(src)
	if err != nil {
		t.Fatal(err)
	}
	sem, err := c.Analyze(p)
	if err != nil {
		t.Fatal(err)
	}
	if !sem.UserFuncs["f"] {
		t.Fatal("expected user func 'f' in UserFuncs")
	}
}

func TestStagedEmitMatchesCompile(t *testing.T) {
	src := `Print(1)`
	c := New()
	p, err := c.Parse(src)
	if err != nil {
		t.Fatal(err)
	}
	sem, err := c.Analyze(p)
	if err != nil {
		t.Fatal(err)
	}
	chunk1, err := codegen.Emit(p, sem)
	if err != nil {
		t.Fatal(err)
	}
	chunk2, err := c.Compile(src)
	if err != nil {
		t.Fatal(err)
	}
	if len(chunk1.Code) != len(chunk2.Code) {
		t.Errorf("bytecode length: staged %d vs Compile %d", len(chunk1.Code), len(chunk2.Code))
	}
}

func TestCompileWithOptionsFilenameInError(t *testing.T) {
	c := New()
	// Unterminated string → lexical error
	_, err := c.CompileWithOptions(`Print("oops`, CompileOptions{Filename: "bad.bas"})
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "bad.bas") {
		t.Errorf("error should mention filename: %v", err)
	}
}

func TestCompileRegressionStillEmitsHalt(t *testing.T) {
	c := New()
	chunk, err := c.Compile(`InitWindow(800, 600, "t")`)
	if err != nil {
		t.Fatal(err)
	}
	if !chunkContainsOp(chunk, vm.OpHalt) {
		t.Error("expected OpHalt in chunk")
	}
}
