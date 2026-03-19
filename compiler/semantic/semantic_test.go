package semantic

import (
	"cyberbasic/compiler/lexer"
	"cyberbasic/compiler/parser"
	"testing"
)

func mustParse(t *testing.T, source string) *parser.Program {
	t.Helper()
	l := lexer.New(source)
	tokens, err := l.Tokenize()
	if err != nil {
		t.Fatalf("lex: %v", err)
	}
	p := parser.New(tokens)
	program, err := p.Parse()
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	return program
}

func TestAnalyze_MinimalProgram(t *testing.T) {
	src := `DIM x AS INTEGER
x = 1
`
	program := mustParse(t, src)
	result, err := Analyze(program)
	if err != nil {
		t.Fatalf("Analyze: %v", err)
	}
	if result.TypeDefs == nil || len(result.TypeDefs) != 0 {
		t.Error("expected TypeDefs non-nil and empty")
	}
	if result.EntityNames == nil || len(result.EntityNames) != 0 {
		t.Error("expected EntityNames non-nil and empty")
	}
	if result.UserFuncs == nil || len(result.UserFuncs) != 0 {
		t.Error("expected UserFuncs non-nil and empty")
	}
	if len(result.MainStmts) != 2 {
		t.Errorf("expected 2 main stmts, got %d", len(result.MainStmts))
	}
	if len(result.Decls) != 0 {
		t.Errorf("expected 0 decls, got %d", len(result.Decls))
	}
}

func TestAnalyze_TypeAndEntity(t *testing.T) {
	src := `TYPE Player
  x AS INTEGER
  y AS INTEGER
END TYPE
ENTITY Hero
  n = 0
END ENTITY
DIM p AS Player
`
	program := mustParse(t, src)
	result, err := Analyze(program)
	if err != nil {
		t.Fatalf("Analyze: %v", err)
	}
	if len(result.TypeDefs) != 1 {
		t.Errorf("expected 1 type def, got %d", len(result.TypeDefs))
	}
	if _, ok := result.TypeDefs["player"]; !ok {
		t.Error("expected TypeDefs[\"player\"]")
	}
	if !result.EntityNames["hero"] {
		t.Error("expected EntityNames[\"hero\"]")
	}
	// TYPE and ENTITY are not FUNCTION/SUB/MODULE so they appear in MainStmts; DIM is also main.
	if len(result.MainStmts) != 3 {
		t.Errorf("expected 3 main stmts (TYPE, ENTITY, DIM), got %d", len(result.MainStmts))
	}
	if len(result.Decls) != 0 {
		t.Errorf("expected 0 decls, got %d", len(result.Decls))
	}
}

func TestAnalyze_UserFunctionAndSub(t *testing.T) {
	src := `SUB DoIt()
END SUB
FUNCTION Add(a, b)
  RETURN a + b
END FUNCTION
VAR n = Add(1, 2)
`
	program := mustParse(t, src)
	result, err := Analyze(program)
	if err != nil {
		t.Fatalf("Analyze: %v", err)
	}
	if !result.UserFuncs["doit"] {
		t.Error("expected UserFuncs[\"doit\"]")
	}
	if !result.UserFuncs["add"] {
		t.Error("expected UserFuncs[\"add\"]")
	}
	if len(result.Decls) != 2 {
		t.Errorf("expected 2 decls, got %d", len(result.Decls))
	}
	if len(result.MainStmts) != 1 {
		t.Errorf("expected 1 main stmt, got %d", len(result.MainStmts))
	}
}

func TestAnalyze_ModuleQualifiedNames(t *testing.T) {
	src := `MODULE M
  FUNCTION F(x)
    RETURN x
  END FUNCTION
  SUB S()
  END SUB
END MODULE
`
	program := mustParse(t, src)
	result, err := Analyze(program)
	if err != nil {
		t.Fatalf("Analyze: %v", err)
	}
	if !result.UserFuncs["m.f"] {
		t.Error("expected UserFuncs[\"m.f\"]")
	}
	if !result.UserFuncs["m.s"] {
		t.Error("expected UserFuncs[\"m.s\"]")
	}
	if len(result.Decls) != 2 {
		t.Errorf("expected 2 decls, got %d", len(result.Decls))
	}
}

func TestQualifiedName(t *testing.T) {
	tests := []struct {
		name     string
		node     parser.Node
		expected string
	}{
		{"plain function", &parser.FunctionDecl{Name: "Add", ModuleName: ""}, "add"},
		{"plain sub", &parser.SubDecl{Name: "DoIt", ModuleName: ""}, "doit"},
		{"module function", &parser.FunctionDecl{Name: "F", ModuleName: "M"}, "m.f"},
		{"module sub", &parser.SubDecl{Name: "S", ModuleName: "Util"}, "util.s"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := QualifiedName(tt.node)
			if got != tt.expected {
				t.Errorf("QualifiedName() = %q, want %q", got, tt.expected)
			}
		})
	}
}
