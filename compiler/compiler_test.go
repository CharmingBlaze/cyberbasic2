package compiler

import (
	"cyberbasic/compiler/bindings/std"
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

func TestCompileElseIf(t *testing.T) {
	src := `Function Choose(op, a, b)
  If op = 1 Then
    Return a + b
  ElseIf op = 2 Then
    Return a * b
  Else
    Return 0
  End If
EndFunction
VAR n = Choose(1, 3, 4)
`
	chunk := mustCompile(t, src)
	if !chunkContainsOp(chunk, vm.OpJumpIfFalse) {
		t.Error("expected OpJumpIfFalse for IF/ELSEIF")
	}
	if _, ok := chunk.Functions["choose"]; !ok {
		t.Error("expected chunk.Functions to contain 'choose'")
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

// TestCompileAndRunEntityCreateOnly verifies ENTITY creates a dict in globals (no assignment).
func TestCompileAndRunEntityCreateOnly(t *testing.T) {
	src := `ENTITY Player
  x = 100
  y = 200
END ENTITY
VAR a = Player.x
VAR b = Player.y
`
	c := New()
	chunk, err := c.Compile(src)
	if err != nil {
		t.Fatalf("compile: %v", err)
	}
	v := vm.NewVM()
	std.RegisterStd(v)
	v.LoadChunk(chunk)
	err = v.Run()
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	g, ok := v.Globals()["player"]
	if !ok {
		t.Fatal("expected global 'player'")
	}
	m, _ := g.(map[string]interface{})
	if m["x"] != float64(100) && m["x"] != 100 {
		t.Errorf("player.x: got %v", m["x"])
	}
	if m["y"] != float64(200) && m["y"] != 200 {
		t.Errorf("player.y: got %v", m["y"])
	}
}

// TestCompileAndRunEntity verifies ENTITY creates a dict in globals and property read/write works.
func TestCompileAndRunEntity(t *testing.T) {
	src := `ENTITY Player
  x = 100
  y = 200
END ENTITY
Player.x = 50
VAR a = Player.x
VAR b = Player.y
`
	c := New()
	chunk, err := c.Compile(src)
	if err != nil {
		t.Fatalf("compile: %v", err)
	}
	if !chunkContainsOp(chunk, vm.OpStoreGlobal) {
		t.Error("expected OpStoreGlobal for entity")
	}
	if !chunkContainsOp(chunk, vm.OpLoadEntityProp) {
		t.Error("expected OpLoadEntityProp for entity property read")
	}
	v := vm.NewVM()
	std.RegisterStd(v)
	v.LoadChunk(chunk)
	err = v.Run()
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	// Check globals: player entity with x=50 (updated by assignment), y=200
	g, ok := v.Globals()["player"]
	if !ok {
		t.Fatal("expected global 'player' after run")
	}
	m, ok := g.(map[string]interface{})
	if !ok {
		t.Fatalf("expected player to be map, got %T", g)
	}
	if x, ok := m["x"]; !ok {
		t.Error("expected player.x")
	} else if xi, ok := x.(float64); ok && int(xi) != 50 {
		t.Errorf("player.x: expected 50, got %v", x)
	} else if xi, ok := x.(int); ok && xi != 50 {
		t.Errorf("player.x: expected 50, got %v", x)
	}
	if y, ok := m["y"]; !ok {
		t.Error("expected player.y")
	} else if yi, ok := y.(float64); ok && int(yi) != 200 {
		t.Errorf("player.y: expected 200, got %v", y)
	} else if yi, ok := y.(int); ok && yi != 200 {
		t.Errorf("player.y: expected 200, got %v", y)
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

// chunkHasConstant returns true if chunk.Constants contains the given string (e.g. "begindrawing").
func chunkHasConstant(chunk *vm.Chunk, name string) bool {
	for _, c := range chunk.Constants {
		if s, ok := c.(string); ok && s == name {
			return true
		}
	}
	return false
}

// TestGameLoopFrameWrap asserts that WHILE NOT WindowShouldClose() ... WEND gets automatic
// BeginDrawing, EndDrawing, BeginMode2D, EndMode2D so the user doesn't need to call them.
func TestGameLoopFrameWrap(t *testing.T) {
	src := `InitWindow(800, 600, "test")
WHILE NOT WindowShouldClose()
  ClearBackground(0, 0, 0, 255)
WEND
CloseWindow()
`
	chunk := mustCompile(t, src)
	for _, name := range []string{"begindrawing", "enddrawing", "beginmode2d", "endmode2d"} {
		if !chunkHasConstant(chunk, name) {
			t.Errorf("expected chunk to contain frame-wrap constant %q", name)
		}
	}
	// 2D loop should not inject 3D mode
	if chunkHasConstant(chunk, "beginmode3d") {
		t.Errorf("2D loop should not contain beginmode3d")
	}
}

// TestGameLoopFrameWrap3D asserts that a WHILE NOT WindowShouldClose() loop with 3D drawing
// gets automatic BeginMode3D/EndMode3D instead of 2D.
func TestGameLoopFrameWrap3D(t *testing.T) {
	src := `InitWindow(800, 600, "test")
SetCamera3D(0, 10, 10, 0, 0, 0, 0, 1, 0)
WHILE NOT WindowShouldClose()
  ClearBackground(0, 0, 0, 255)
  DrawCube(0, 0, 0, 2, 2, 2, 255, 0, 0, 255)
WEND
CloseWindow()
`
	chunk := mustCompile(t, src)
	for _, name := range []string{"begindrawing", "enddrawing", "beginmode3d", "endmode3d"} {
		if !chunkHasConstant(chunk, name) {
			t.Errorf("expected 3D loop to contain frame-wrap constant %q", name)
		}
	}
}
