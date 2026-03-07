package box2d

import (
	"testing"

	"cyberbasic/compiler/vm"
)

func TestBox2DBackendMetadata(t *testing.T) {
	v := vm.NewVM()
	RegisterBox2D(v)

	name, err := v.CallForeign("Box2DBackendName", nil)
	if err != nil {
		t.Fatalf("Box2DBackendName failed: %v", err)
	}
	if name != "bytearena-box2d" {
		t.Fatalf("unexpected backend name: %v", name)
	}

	mode, err := v.CallForeign("Box2DBackendMode", nil)
	if err != nil {
		t.Fatalf("Box2DBackendMode failed: %v", err)
	}
	if mode != "authoritative" {
		t.Fatalf("unexpected backend mode: %v", mode)
	}
}
