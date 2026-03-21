package shadersys

import (
	"testing"

	"cyberbasic/compiler/vm"
)

func TestEmbeddedGLSLNonEmpty(t *testing.T) {
	if len(presetVertexShader) < 80 {
		t.Fatal("preset vertex shader unexpectedly short")
	}
	for _, name := range []struct {
		n string
		s string
	}{
		{"pbr", presetPBRFragment},
		{"toon", presetToonFragment},
		{"dissolve", presetDissolveFragment},
	} {
		if len(name.s) < 40 {
			t.Fatalf("preset %s fragment unexpectedly short", name.n)
		}
	}
}

func TestRegisterShaderSysVersionForeign(t *testing.T) {
	v := vm.NewVM()
	RegisterShaderSys(v)
	out, err := v.CallForeign("ShaderSysVersion", []interface{}{})
	if err != nil {
		t.Fatal(err)
	}
	if out != "v2-raylib" {
		t.Fatalf("ShaderSysVersion = %v, want v2-raylib", out)
	}
	ai := v.Globals()["shader"]
	if ai == nil {
		t.Fatal("expected global shader")
	}
}

func TestShaderHandleSetSignature(t *testing.T) {
	v := vm.NewVM()
	RegisterShaderSys(v)
	_, err := v.CallForeign("ShaderHandleSet", []interface{}{})
	if err == nil {
		t.Fatal("expected error for too few args")
	}
}
