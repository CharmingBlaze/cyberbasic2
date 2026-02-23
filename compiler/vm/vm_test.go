package vm

import (
	"testing"
)

// TestSetRuntime ensures the VM accepts a GameRuntime and stores it
func TestSetRuntime(t *testing.T) {
	v := NewVM()
	if v.runtime != nil {
		t.Fatal("new VM should have nil runtime")
	}
	// SetRuntime with nil should not panic
	v.SetRuntime(nil)
	v.SetRuntime(&mockRuntime{})
	if v.runtime == nil {
		t.Fatal("SetRuntime should set runtime")
	}
}

// mockRuntime implements GameRuntime for tests
type mockRuntime struct {
	LoadImageCalled          bool
	CreateSpriteCalled       bool
	SetSpritePositionCalled  bool
	DrawSpriteCalled         bool
	CreatePhysicsBodyCalled  bool
	InitializeGraphicsCalled bool
}

func (m *mockRuntime) LoadImage(filename string) error {
	m.LoadImageCalled = true
	return nil
}
func (m *mockRuntime) CreateSprite(id, image string, x, y float64) error {
	m.CreateSpriteCalled = true
	return nil
}
func (m *mockRuntime) SetSpritePosition(id string, x, y float64) error {
	m.SetSpritePositionCalled = true
	return nil
}
func (m *mockRuntime) DrawSprite(id string) error {
	m.DrawSpriteCalled = true
	return nil
}
func (m *mockRuntime) LoadModel(filename string) error { return nil }
func (m *mockRuntime) CreateCamera(id string, x, y, z float64) error { return nil }
func (m *mockRuntime) SetCameraPosition(id string, x, y, z float64) error { return nil }
func (m *mockRuntime) DrawModel(id string, x, y, z, scale float64) error { return nil }
func (m *mockRuntime) PlayMusic(filename string) error { return nil }
func (m *mockRuntime) PlaySound(filename string) error { return nil }
func (m *mockRuntime) LoadSound(filename string) error { return nil }
func (m *mockRuntime) CreatePhysicsBody(id, bodyType string, x, y, z, mass float64) error {
	m.CreatePhysicsBodyCalled = true
	return nil
}
func (m *mockRuntime) SetVelocity(id string, vx, vy, vz float64) error { return nil }
func (m *mockRuntime) ApplyForce(id string, fx, fy, fz float64) error   { return nil }
func (m *mockRuntime) RayCast3D(startX, startY, startZ, dirX, dirY, dirZ, maxDist float64) (bool, float64, float64, float64, error) {
	return false, 0, 0, 0, nil
}
func (m *mockRuntime) InitializeGraphics(width, height int, title string) error {
	m.InitializeGraphicsCalled = true
	return nil
}
func (m *mockRuntime) InitializePhysics() error { return nil }
func (m *mockRuntime) ShouldClose() bool         { return false }
func (m *mockRuntime) Sync() error              { return nil }
func (m *mockRuntime) IsKeyDown(keyName string) bool    { return false }
func (m *mockRuntime) IsKeyPressed(keyName string) bool { return false }

// TestGameOpcodesInvokeRuntime verifies that when runtime is set, game opcodes call it
func TestGameOpcodesInvokeRuntime(t *testing.T) {
	// Build a minimal chunk that does LOADIMAGE, CREATESPRITE, SETSPRITEPOSITION, DRAWSPRITE
	chunk := NewChunk()
	// Push "file.png" and emit OpLoadImage
	chunk.Write(byte(OpLoadConst))
	ci := chunk.WriteConstant("file.png")
	chunk.Write(byte(ci))
	chunk.Write(byte(OpLoadImage))

	chunk.Write(byte(OpLoadConst))
	ci = chunk.WriteConstant("s1")
	chunk.Write(byte(ci))
	chunk.Write(byte(OpLoadConst))
	ci = chunk.WriteConstant("file.png")
	chunk.Write(byte(ci))
	chunk.Write(byte(OpLoadConst))
	ci = chunk.WriteConstant(100.0)
	chunk.Write(byte(ci))
	chunk.Write(byte(OpLoadConst))
	ci = chunk.WriteConstant(200.0)
	chunk.Write(byte(ci))
	chunk.Write(byte(OpCreateSprite))

	chunk.Write(byte(OpLoadConst))
	ci = chunk.WriteConstant("s1")
	chunk.Write(byte(ci))
	chunk.Write(byte(OpLoadConst))
	ci = chunk.WriteConstant(50.0)
	chunk.Write(byte(ci))
	chunk.Write(byte(OpLoadConst))
	ci = chunk.WriteConstant(60.0)
	chunk.Write(byte(ci))
	chunk.Write(byte(OpSetSpritePosition))

	chunk.Write(byte(OpLoadConst))
	ci = chunk.WriteConstant("s1")
	chunk.Write(byte(ci))
	chunk.Write(byte(OpDrawSprite))

	chunk.Write(byte(OpHalt))

	mock := &mockRuntime{}
	v := NewVM()
	v.SetRuntime(mock)
	v.LoadChunk(chunk)
	err := v.Run()
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if !mock.LoadImageCalled {
		t.Error("LoadImage was not called")
	}
	if !mock.CreateSpriteCalled {
		t.Error("CreateSprite was not called")
	}
	if !mock.SetSpritePositionCalled {
		t.Error("SetSpritePosition was not called")
	}
	if !mock.DrawSpriteCalled {
		t.Error("DrawSprite was not called")
	}
}

// TestUserFunctionCall verifies OpCallUser and OpReturnVal: call a user function and get return value
func TestUserFunctionCall(t *testing.T) {
	// Chunk layout: main (push 3, push 5, OpCallUser), OpHalt, then function "add" body
	chunk := NewChunk()
	// Main: push 3, push 5, OpCallUser add with 2 args
	chunk.Write(byte(OpLoadConst))
	ci := chunk.WriteConstant(3)
	chunk.Write(byte(ci))
	chunk.Write(byte(OpLoadConst))
	ci = chunk.WriteConstant(5)
	chunk.Write(byte(ci))
	chunk.Write(byte(OpCallUser))
	ci = chunk.WriteConstant("add")
	chunk.Write(byte(ci))
	chunk.Write(byte(2))
	chunk.Write(byte(OpHalt))
	// Function "add" body: stack is [3, 5]. LoadVar 0, LoadVar 1, Add, ReturnVal
	chunk.Functions["add"] = len(chunk.Code)
	chunk.Write(byte(OpLoadVar))
	chunk.Write(byte(0))
	chunk.Write(byte(OpLoadVar))
	chunk.Write(byte(1))
	chunk.Write(byte(OpAdd))
	chunk.Write(byte(OpReturnVal))

	v := NewVM()
	v.LoadChunk(chunk)
	err := v.Run()
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if len(v.stack) != 1 {
		t.Fatalf("expected 1 value on stack, got %d", len(v.stack))
	}
	got, ok := v.stack[0].(int)
	if !ok {
		// might be float64 from constant
		if f, ok := v.stack[0].(float64); ok {
			got = int(f)
		} else {
			t.Fatalf("stack top: got %T %v", v.stack[0], v.stack[0])
		}
	}
	if got != 8 {
		t.Errorf("expected 8, got %d", got)
	}
}

// TestEventRegistration verifies that OpRegisterEvent adds a handler to the VM's event list
func TestEventRegistration(t *testing.T) {
	chunk := NewChunk()
	etIdx := chunk.WriteConstant("keydown")
	keyIdx := chunk.WriteConstant("X")
	chunk.Write(byte(OpRegisterEvent))
	chunk.Write(byte(etIdx))
	chunk.Write(byte(keyIdx))
	handlerIP := len(chunk.Code) + 2
	chunk.Write(byte(handlerIP & 0xff))
	chunk.Write(byte(handlerIP >> 8))
	chunk.Write(byte(OpHalt))

	v := NewVM()
	v.LoadChunk(chunk)
	err := v.Run()
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if len(v.eventHandlers) != 1 {
		t.Fatalf("expected 1 event handler, got %d", len(v.eventHandlers))
	}
	h := v.eventHandlers[0]
	if h.eventType != "keydown" || h.key != "X" {
		t.Errorf("handler: eventType=%q key=%q", h.eventType, h.key)
	}
}

// TestCoroutineSwitch verifies StartCoroutine and Yield: main and fiber both run
func TestCoroutineSwitch(t *testing.T) {
	// Main: set main1, StartCoroutine worker, Yield, set main2, Halt.
	// Worker: set co1, Yield, set co2, Return (fiber exits).
	chunk := NewChunk()
	chunk.Write(byte(OpLoadConst))
	ci := chunk.WriteConstant(1)
	chunk.Write(byte(ci))
	chunk.Write(byte(OpStoreGlobal))
	gi := chunk.WriteConstant("main1")
	chunk.Write(byte(gi))
	chunk.Write(byte(OpStartCoroutine))
	chunk.Write(byte(0)) // worker IP low (patched below)
	chunk.Write(byte(0)) // worker IP high
	chunk.Write(byte(OpYield))
	chunk.Write(byte(OpLoadConst))
	ci = chunk.WriteConstant(2)
	chunk.Write(byte(ci))
	chunk.Write(byte(OpStoreGlobal))
	gi = chunk.WriteConstant("main2")
	chunk.Write(byte(gi))
	chunk.Write(byte(OpYield)) // yield again so worker can run and set co2
	chunk.Write(byte(OpHalt))
	// Patch StartCoroutine target: the two bytes after the opcode are at indices 5 and 6
	workerStart := len(chunk.Code)
	chunk.Code[5] = byte(workerStart & 0xff)
	chunk.Code[6] = byte(workerStart >> 8)
	// Worker: set co1, Yield, set co2, Return
	chunk.Write(byte(OpLoadConst))
	ci = chunk.WriteConstant(3)
	chunk.Write(byte(ci))
	chunk.Write(byte(OpStoreGlobal))
	gi = chunk.WriteConstant("co1")
	chunk.Write(byte(gi))
	chunk.Write(byte(OpYield))
	chunk.Write(byte(OpLoadConst))
	ci = chunk.WriteConstant(4)
	chunk.Write(byte(ci))
	chunk.Write(byte(OpStoreGlobal))
	gi = chunk.WriteConstant("co2")
	chunk.Write(byte(gi))
	chunk.Write(byte(OpReturn))

	v := NewVM()
	v.LoadChunk(chunk)
	err := v.Run()
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	// After Run: main1, main2, co1, co2 should be set. Order depends on scheduler (main runs first, then Yield to worker, etc.)
	if v.globals["main1"] != 1 {
		t.Errorf("main1: got %v", v.globals["main1"])
	}
	if v.globals["main2"] != 2 {
		t.Errorf("main2: got %v", v.globals["main2"])
	}
	if v.globals["co1"] != 3 {
		t.Errorf("co1: got %v", v.globals["co1"])
	}
	if v.globals["co2"] != 4 {
		t.Errorf("co2: got %v", v.globals["co2"])
	}
}
