package dbp

import (
	"testing"

	"cyberbasic/compiler/vm"
)

func resetTilemapWrapperState() {
	tilemapId2StrMu.Lock()
	tilemapId2Str = make(map[int]string)
	tilemapId2StrMu.Unlock()
	tilemapVisibleMu.Lock()
	tilemapVisible = make(map[int]bool)
	tilemapVisibleMu.Unlock()
}

func TestTilemapWrappersUseInternalAliases(t *testing.T) {
	resetTilemapWrapperState()
	v := vm.NewVM()

	loadCalls := 0
	drawCalls := 0
	tilesetCalls := 0
	deleteCalls := 0
	var drawArgs []interface{}
	var tilesetArgs []interface{}
	var deleteArgs []interface{}

	v.RegisterForeign("TilemapLoadByPath", func(args []interface{}) (interface{}, error) {
		loadCalls++
		return "tm_internal", nil
	})
	v.RegisterForeign("TilemapDrawByMapId", func(args []interface{}) (interface{}, error) {
		drawCalls++
		drawArgs = append([]interface{}{}, args...)
		return nil, nil
	})
	v.RegisterForeign("TilemapSetTilesetByMapId", func(args []interface{}) (interface{}, error) {
		tilesetCalls++
		tilesetArgs = append([]interface{}{}, args...)
		return nil, nil
	})
	v.RegisterForeign("SetTileByMapId", func(args []interface{}) (interface{}, error) { return nil, nil })
	v.RegisterForeign("GetTileByMapId", func(args []interface{}) (interface{}, error) { return 0, nil })
	v.RegisterForeign("TilemapDeleteByMapId", func(args []interface{}) (interface{}, error) {
		deleteCalls++
		deleteArgs = append([]interface{}{}, args...)
		return nil, nil
	})

	register2DTilemaps(v)

	if _, err := v.CallForeign("LoadTilemap", []interface{}{7, "levels/level1.json"}); err != nil {
		t.Fatalf("LoadTilemap failed: %v", err)
	}
	if loadCalls != 1 {
		t.Fatalf("expected one internal load call, got %d", loadCalls)
	}
	if _, err := v.CallForeign("DrawTilemap", []interface{}{7, 16, 32}); err != nil {
		t.Fatalf("DrawTilemap failed: %v", err)
	}
	if drawCalls != 1 {
		t.Fatalf("expected one internal draw call, got %d", drawCalls)
	}
	if len(drawArgs) != 3 || drawArgs[0] != "tm_internal" || drawArgs[1] != 16 || drawArgs[2] != 32 {
		t.Fatalf("unexpected draw args: %#v", drawArgs)
	}
	if _, err := v.CallForeign("TilemapSetTileset", []interface{}{7, "tiles/dungeon.png"}); err != nil {
		t.Fatalf("TilemapSetTileset failed: %v", err)
	}
	if tilesetCalls != 1 {
		t.Fatalf("expected one internal tileset call, got %d", tilesetCalls)
	}
	if len(tilesetArgs) != 2 || tilesetArgs[0] != "tm_internal" || tilesetArgs[1] != "tiles/dungeon.png" {
		t.Fatalf("unexpected tileset args: %#v", tilesetArgs)
	}
	if _, err := v.CallForeign("DeleteTilemap", []interface{}{7}); err != nil {
		t.Fatalf("DeleteTilemap failed: %v", err)
	}
	if deleteCalls != 1 {
		t.Fatalf("expected one internal delete call, got %d", deleteCalls)
	}
	if len(deleteArgs) != 1 || deleteArgs[0] != "tm_internal" {
		t.Fatalf("unexpected delete args: %#v", deleteArgs)
	}
}

func Test2DPhysicsWrappersUseInternalAliases(t *testing.T) {
	v := vm.NewVM()

	staticCalls := 0
	forceCalls := 0
	impulseCalls := 0
	var staticArgs []interface{}
	var forceArgs []interface{}
	var impulseArgs []interface{}

	v.RegisterForeign("MakeStaticBody2D", func(args []interface{}) (interface{}, error) {
		staticCalls++
		staticArgs = append([]interface{}{}, args...)
		return "player", nil
	})
	v.RegisterForeign("ApplyForce2DByBodyId", func(args []interface{}) (interface{}, error) {
		forceCalls++
		forceArgs = append([]interface{}{}, args...)
		return nil, nil
	})
	v.RegisterForeign("ApplyImpulse2DByBodyId", func(args []interface{}) (interface{}, error) {
		impulseCalls++
		impulseArgs = append([]interface{}{}, args...)
		return nil, nil
	})

	register2DPhysics(v)

	if _, err := v.CallForeign("MakeStatic2D", []interface{}{"player"}); err != nil {
		t.Fatalf("MakeStatic2D failed: %v", err)
	}
	if staticCalls != 1 {
		t.Fatalf("expected one static body call, got %d", staticCalls)
	}
	if len(staticArgs) != 5 || staticArgs[0] != "player" {
		t.Fatalf("unexpected static args: %#v", staticArgs)
	}

	if _, err := v.CallForeign("ApplyForce2D", []interface{}{"player", 3.0, -1.5}); err != nil {
		t.Fatalf("ApplyForce2D failed: %v", err)
	}
	if forceCalls != 1 {
		t.Fatalf("expected one force alias call, got %d", forceCalls)
	}
	if len(forceArgs) != 4 || forceArgs[0] != "default" || forceArgs[1] != "player" || forceArgs[2] != 3.0 || forceArgs[3] != -1.5 {
		t.Fatalf("unexpected force args: %#v", forceArgs)
	}

	if _, err := v.CallForeign("ApplyImpulse2D", []interface{}{"player", 2.0, 4.0}); err != nil {
		t.Fatalf("ApplyImpulse2D failed: %v", err)
	}
	if impulseCalls != 1 {
		t.Fatalf("expected one impulse alias call, got %d", impulseCalls)
	}
	if len(impulseArgs) != 4 || impulseArgs[0] != "default" || impulseArgs[1] != "player" || impulseArgs[2] != 2.0 || impulseArgs[3] != 4.0 {
		t.Fatalf("unexpected impulse args: %#v", impulseArgs)
	}
}
