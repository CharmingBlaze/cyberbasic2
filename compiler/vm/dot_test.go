package vm

import (
	"testing"
)

type testDot struct {
	x float64
}

func (t *testDot) GetProp(path []string) (Value, error) {
	if len(path) == 1 && path[0] == "x" {
		return t.x, nil
	}
	return nil, nil
}

func (t *testDot) SetProp(path []string, val Value) error {
	if len(path) == 1 && path[0] == "x" {
		if f, ok := val.(float64); ok {
			t.x = f
		}
	}
	return nil
}

func (t *testDot) CallMethod(name string, args []Value) (Value, error) {
	return nil, nil
}

func TestOpGetPropSetProp(t *testing.T) {
	v := NewVM()
	ch := NewChunk()
	td := &testDot{x: 3}
	idx := ch.WriteConstant(td)
	cx := ch.WriteConstant("x")
	ch.Write(byte(OpLoadConst))
	ch.Write(byte(idx))
	ch.Write(byte(OpGetProp))
	ch.Write(byte(1))
	ch.Write(byte(cx))
	ch.Write(byte(OpHalt))
	v.LoadChunk(ch)
	v.ip = 0
	if err := v.Run(); err != nil {
		t.Fatal(err)
	}
	if len(v.stack) != 1 || v.stack[0].(float64) != 3 {
		t.Fatalf("getprop: stack %v", v.stack)
	}

	ch2 := NewChunk()
	td2 := &testDot{x: 0}
	iObj := ch2.WriteConstant(td2)
	ch2.Write(byte(OpLoadConst))
	ch2.Write(byte(iObj))
	ci := ch2.WriteConstant(42.0)
	ch2.Write(byte(OpLoadConst))
	ch2.Write(byte(ci))
	cx2 := ch2.WriteConstant("x")
	ch2.Write(byte(OpSetProp))
	ch2.Write(byte(1))
	ch2.Write(byte(cx2))
	ch2.Write(byte(OpHalt))
	v2 := NewVM()
	v2.LoadChunk(ch2)
	v2.ip = 0
	if err := v2.Run(); err != nil {
		t.Fatal(err)
	}
	if td2.x != 42 {
		t.Fatalf("setprop: x=%v", td2.x)
	}
}
