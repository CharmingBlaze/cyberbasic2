package parser

import (
	"cyberbasic/compiler/lexer"
	"strings"
	"testing"
)

func mustParse(t *testing.T, source string) *Program {
	t.Helper()
	l := lexer.New(source)
	tokens, err := l.Tokenize()
	if err != nil {
		t.Fatalf("tokenize: %v", err)
	}
	p := New(tokens)
	prog, err := p.Parse()
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	return prog
}

func TestParseModule(t *testing.T) {
	src := `Module Math3D
  Function Dot(x, y)
    Return 0
  End Function
End Module
`
	prog := mustParse(t, src)
	if len(prog.Statements) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(prog.Statements))
	}
	mod, ok := prog.Statements[0].(*ModuleStatement)
	if !ok {
		t.Fatalf("expected ModuleStatement, got %T", prog.Statements[0])
	}
	if strings.ToLower(mod.Name) != "math3d" {
		t.Errorf("module name: got %q", mod.Name)
	}
	if len(mod.Body) != 1 {
		t.Fatalf("expected 1 body node, got %d", len(mod.Body))
	}
	fn, ok := mod.Body[0].(*FunctionDecl)
	if !ok {
		t.Fatalf("expected FunctionDecl in module body, got %T", mod.Body[0])
	}
	if strings.ToLower(fn.ModuleName) != "math3d" {
		t.Errorf("FunctionDecl.ModuleName: got %q", fn.ModuleName)
	}
	if strings.ToLower(fn.Name) != "dot" {
		t.Errorf("function name: got %q", fn.Name)
	}
}

func TestParseSubAndFunction(t *testing.T) {
	src := `Sub Foo()
  Print("hi")
End Sub
Function Bar(a)
  Return a
End Function
`
	prog := mustParse(t, src)
	if len(prog.Statements) != 2 {
		t.Fatalf("expected 2 statements, got %d", len(prog.Statements))
	}
	sub, ok := prog.Statements[0].(*SubDecl)
	if !ok {
		t.Fatalf("expected SubDecl, got %T", prog.Statements[0])
	}
	if strings.ToLower(sub.Name) != "foo" {
		t.Errorf("sub name: got %q", sub.Name)
	}
	fn, ok := prog.Statements[1].(*FunctionDecl)
	if !ok {
		t.Fatalf("expected FunctionDecl, got %T", prog.Statements[1])
	}
	if strings.ToLower(fn.Name) != "bar" || len(fn.Parameters) != 1 || strings.ToLower(fn.Parameters[0]) != "a" {
		t.Errorf("function: got %q %v", fn.Name, fn.Parameters)
	}
}

func TestParseSelectCase(t *testing.T) {
	src := `Select Case x
  Case 1
    Print("one")
  Case 2
    Print("two")
End Select
`
	prog := mustParse(t, src)
	if len(prog.Statements) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(prog.Statements))
	}
	_, ok := prog.Statements[0].(*SelectCaseStatement)
	if !ok {
		t.Fatalf("expected SelectCaseStatement, got %T", prog.Statements[0])
	}
}

func TestParseTypeDecl(t *testing.T) {
	src := `Type Point
  x As Float
  y As Float
End Type
`
	prog := mustParse(t, src)
	if len(prog.Statements) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(prog.Statements))
	}
	td, ok := prog.Statements[0].(*TypeDecl)
	if !ok {
		t.Fatalf("expected TypeDecl, got %T", prog.Statements[0])
	}
	if strings.ToLower(td.Name) != "point" {
		t.Errorf("type name: got %q", td.Name)
	}
}

func TestParseEnum(t *testing.T) {
	src := `Enum Color: Red, Green, Blue
`
	prog := mustParse(t, src)
	if len(prog.Statements) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(prog.Statements))
	}
	en, ok := prog.Statements[0].(*EnumStatement)
	if !ok {
		t.Fatalf("expected EnumStatement, got %T", prog.Statements[0])
	}
	if strings.ToLower(en.Name) != "color" || len(en.Members) != 3 {
		t.Errorf("enum: got %q %d members", en.Name, len(en.Members))
	}
}

func TestParseEnumEndEnum(t *testing.T) {
	// "End Enum" at top level must be consumed by enum parser so we don't hang
	src := `Enum Color: Red, Green, Blue
End Enum
Sub Main()
End Sub
`
	prog := mustParse(t, src)
	if len(prog.Statements) != 2 {
		t.Fatalf("expected 2 statements (enum + sub), got %d", len(prog.Statements))
	}
	en, ok := prog.Statements[0].(*EnumStatement)
	if !ok {
		t.Fatalf("expected first statement EnumStatement, got %T", prog.Statements[0])
	}
	if strings.ToLower(en.Name) != "color" || len(en.Members) != 3 {
		t.Errorf("enum: got %q %d members", en.Name, len(en.Members))
	}
	_, ok = prog.Statements[1].(*SubDecl)
	if !ok {
		t.Fatalf("expected second statement SubDecl, got %T", prog.Statements[1])
	}
}

func TestParseOnKeyDownEndOn(t *testing.T) {
	src := `On KeyDown("ESCAPE")
  Print("quit")
End On
`
	prog := mustParse(t, src)
	if len(prog.Statements) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(prog.Statements))
	}
	on, ok := prog.Statements[0].(*OnEventStatement)
	if !ok {
		t.Fatalf("expected OnEventStatement, got %T", prog.Statements[0])
	}
	if on.EventType != "keydown" {
		t.Errorf("event type: got %q", on.EventType)
	}
	// Key may be stored with or without surrounding quotes depending on lexer
	if on.Key != "ESCAPE" && on.Key != "\"ESCAPE\"" && !strings.Contains(on.Key, "ESCAPE") {
		t.Errorf("event key: got %q", on.Key)
	}
	if on.Body == nil || len(on.Body.Statements) != 1 {
		n := 0
		if on.Body != nil {
			n = len(on.Body.Statements)
		}
		t.Errorf("expected 1 body statement, got %d", n)
	}
}

func TestParseStartCoroutineYieldWaitSeconds(t *testing.T) {
	src := `Sub Worker()
  Yield
  WaitSeconds(1.5)
End Sub
StartCoroutine Worker()
`
	prog := mustParse(t, src)
	if len(prog.Statements) != 2 {
		t.Fatalf("expected 2 statements, got %d", len(prog.Statements))
	}
	sub, ok := prog.Statements[0].(*SubDecl)
	if !ok {
		t.Fatalf("expected SubDecl, got %T", prog.Statements[0])
	}
	if strings.ToLower(sub.Name) != "worker" {
		t.Errorf("sub name: got %q", sub.Name)
	}
	// Second statement is StartCoroutine
	sc, ok := prog.Statements[1].(*StartCoroutineStatement)
	if !ok {
		t.Fatalf("expected StartCoroutineStatement, got %T", prog.Statements[1])
	}
	if strings.ToLower(sc.SubName) != "worker" {
		t.Errorf("StartCoroutine sub name: got %q", sc.SubName)
	}
	// Sub body: Yield and WaitSeconds
	if sub.Body == nil || len(sub.Body.Statements) != 2 {
		n := 0
		if sub.Body != nil {
			n = len(sub.Body.Statements)
		}
		t.Fatalf("expected 2 body statements (Yield, WaitSeconds), got %d", n)
	}
	_, ok = sub.Body.Statements[0].(*YieldStatement)
	if !ok {
		t.Fatalf("expected first body statement YieldStatement, got %T", sub.Body.Statements[0])
	}
	ws, ok := sub.Body.Statements[1].(*WaitSecondsStatement)
	if !ok {
		t.Fatalf("expected second body statement WaitSecondsStatement, got %T", sub.Body.Statements[1])
	}
	if ws.Seconds == nil {
		t.Error("WaitSecondsStatement.Seconds should be set")
	}
}

func TestParseJSONIndexSugar(t *testing.T) {
	src := `VAR x = cfg["key"]
VAR y = cfg["a"]["b"]
`
	prog := mustParse(t, src)
	if len(prog.Statements) != 2 {
		t.Fatalf("expected 2 statements, got %d", len(prog.Statements))
	}
	// First: VAR x = cfg["key"] -> Assignment{Variable: "x", Value: JSONIndexAccess}
	assign, ok := prog.Statements[0].(*Assignment)
	if !ok {
		t.Fatalf("expected Assignment, got %T", prog.Statements[0])
	}
	if strings.ToLower(assign.Variable) != "x" {
		t.Errorf("assign.Variable: got %q", assign.Variable)
	}
	j, ok := assign.Value.(*JSONIndexAccess)
	if !ok {
		t.Fatalf("expected JSONIndexAccess for cfg[\"key\"], got %T", assign.Value)
	}
	if id, ok := j.Object.(*Identifier); !ok || strings.ToLower(id.Name) != "cfg" {
		t.Errorf("JSONIndexAccess.Object: got %T %v", j.Object, j.Object)
	}
	if j.Key != "key" && !strings.Contains(j.Key, "key") {
		t.Errorf("JSONIndexAccess.Key: got %q", j.Key)
	}
	// Second: VAR y = cfg["a"]["b"]
	assign2, ok := prog.Statements[1].(*Assignment)
	if !ok {
		t.Fatalf("expected Assignment, got %T", prog.Statements[1])
	}
	j2, ok := assign2.Value.(*JSONIndexAccess)
	if !ok {
		t.Fatalf("expected JSONIndexAccess for cfg[\"a\"][\"b\"], got %T", assign2.Value)
	}
	if j2.Key != "b" && !strings.Contains(j2.Key, "b") {
		t.Errorf("outer JSONIndexAccess.Key: got %q", j2.Key)
	}
	inner, ok := j2.Object.(*JSONIndexAccess)
	if !ok {
		t.Fatalf("expected inner JSONIndexAccess, got %T", j2.Object)
	}
	if inner.Key != "a" && !strings.Contains(inner.Key, "a") {
		t.Errorf("inner JSONIndexAccess.Key: got %q", inner.Key)
	}
}
