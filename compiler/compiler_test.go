package compiler

import (
	"cyberbasic/compiler/vm"
	"testing"
)

func mustCompile(t *testing.T, source string) *vm.Chunk {
	t.Helper()
	c := New()
	chunk, err := c.Compile(source)
	if err != nil {
		t.Fatalf("compile: %v", err)
	}
	return chunk
}

func chunkContainsOp(chunk *vm.Chunk, op vm.OpCode) bool {
	b := byte(op)
	for _, c := range chunk.Code {
		if c == b {
			return true
		}
	}
	return false
}

func TestCompileForeignCall(t *testing.T) {
	chunk := mustCompile(t, `InitWindow(800, 600, "test")`)
	if !chunkContainsOp(chunk, vm.OpCallForeign) {
		t.Error("expected chunk to contain OpCallForeign")
	}
	if !chunkContainsOp(chunk, vm.OpHalt) {
		t.Error("expected chunk to contain OpHalt")
	}
}

func TestCompileUserFunctionCall(t *testing.T) {
	src := `Function Add(a, b)
  Return a + b
End Function
VAR n = Add(2, 3)
`
	chunk := mustCompile(t, src)
	if !chunkContainsOp(chunk, vm.OpCallUser) {
		t.Error("expected chunk to contain OpCallUser")
	}
	if _, ok := chunk.Functions["add"]; !ok {
		t.Error("expected chunk.Functions to contain 'add'")
	}
}

func TestCompileModuleFunctionCall(t *testing.T) {
	src := `Module M
  Function F(x)
    Return x
  End Function
End Module
VAR n = M.F(1)
`
	chunk := mustCompile(t, src)
	if !chunkContainsOp(chunk, vm.OpCallUser) {
		t.Error("expected chunk to contain OpCallUser")
	}
	if _, ok := chunk.Functions["m.f"]; !ok {
		t.Error("expected chunk.Functions to contain 'm.f'")
	}
}

func TestCompilePrint(t *testing.T) {
	chunk := mustCompile(t, `Print("hello")`)
	if !chunkContainsOp(chunk, vm.OpPrint) {
		t.Error("expected chunk to contain OpPrint")
	}
}

// TestE2ECompileAndRun runs a minimal script (user function only, no foreigns) and checks it completes
func TestE2ECompileAndRun(t *testing.T) {
	src := `Function F()
  Return 42
End Function
VAR x = F()
`
	c := New()
	chunk, err := c.Compile(src)
	if err != nil {
		t.Fatalf("compile: %v", err)
	}
	v := vm.NewVM()
	v.LoadChunk(chunk)
	err = v.Run()
	if err != nil {
		t.Fatalf("run: %v", err)
	}
}

func TestCompileOnEvent(t *testing.T) {
	src := `On KeyDown("X")
  Print("ok")
End On
`
	chunk := mustCompile(t, src)
	if !chunkContainsOp(chunk, vm.OpRegisterEvent) {
		t.Error("expected chunk to contain OpRegisterEvent")
	}
}

func TestCompileCoroutineOpcodes(t *testing.T) {
	src := `Sub Co()
  Yield
  WaitSeconds(1)
End Sub
StartCoroutine Co()
`
	chunk := mustCompile(t, src)
	if !chunkContainsOp(chunk, vm.OpStartCoroutine) {
		t.Error("expected chunk to contain OpStartCoroutine")
	}
	if !chunkContainsOp(chunk, vm.OpYield) {
		t.Error("expected chunk to contain OpYield")
	}
	if !chunkContainsOp(chunk, vm.OpWaitSeconds) {
		t.Error("expected chunk to contain OpWaitSeconds")
	}
}
