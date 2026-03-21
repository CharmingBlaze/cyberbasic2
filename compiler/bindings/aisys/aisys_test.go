package aisys

import (
	"testing"

	"cyberbasic/compiler/bindings/navigation"
	"cyberbasic/compiler/vm"
)

func TestAisysVersion(t *testing.T) {
	v := vm.NewVM()
	navigation.RegisterNavigation(v)
	RegisterAisys(v)
	out, err := v.CallForeign("AisysVersion", []interface{}{})
	if err != nil {
		t.Fatal(err)
	}
	if out != "v1-nav-delegate" {
		t.Fatalf("AisysVersion = %v", out)
	}
}

func TestAINavGridDelegatesToNavigation(t *testing.T) {
	v := vm.NewVM()
	navigation.RegisterNavigation(v)
	RegisterAisys(v)
	aiMod := v.Globals()["ai"].(vm.DotObject)
	gid, err := aiMod.CallMethod("navgridcreate", []vm.Value{5, 4})
	if err != nil {
		t.Fatal(err)
	}
	s := gid.(string)
	if s == "" {
		t.Fatal("empty grid id")
	}
	// Same id whether created via navigation or ai
	navMod := v.Globals()["navigation"].(vm.DotObject)
	gid2, err := navMod.CallMethod("navgridcreate", []vm.Value{5, 4})
	if err != nil {
		t.Fatal(err)
	}
	if gid2.(string) == s {
		t.Fatal("expected distinct grid ids")
	}
}

func TestAIAgentDotDelegates(t *testing.T) {
	v := vm.NewVM()
	navigation.RegisterNavigation(v)
	RegisterAisys(v)
	aiMod := v.Globals()["ai"].(vm.DotObject)
	gid, _ := aiMod.CallMethod("navgridcreate", []vm.Value{3, 3})
	aid, err := aiMod.CallMethod("navagentcreate", []vm.Value{"", gid})
	if err != nil {
		t.Fatal(err)
	}
	agWrap, err := aiMod.CallMethod("agent", []vm.Value{aid})
	if err != nil {
		t.Fatal(err)
	}
	ag := agWrap.(vm.DotObject)
	idProp, err := ag.GetProp([]string{"id"})
	if err != nil || idProp != aid {
		t.Fatalf("agent.id = %v, err %v", idProp, err)
	}
	_, err = ag.CallMethod("setdestination", []vm.Value{2.0, 0.0, 2.0})
	if err != nil {
		t.Fatal(err)
	}
}
