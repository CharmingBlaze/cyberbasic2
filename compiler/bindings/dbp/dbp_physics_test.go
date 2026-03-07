package dbp

import (
	"testing"

	"cyberbasic/compiler/vm"
)

func resetPhysicsWrapperState() {
	physicsBodyMapMu.Lock()
	physicsBodyMap = make(map[int]string)
	physicsBodyMapMu.Unlock()
}

func TestPhysicsWrappersAcceptDocumentedColliderArity(t *testing.T) {
	resetPhysicsWrapperState()
	v := vm.NewVM()

	physicsEnableCalls := 0
	sphereCalls := 0
	capsuleCalls := 0
	static2DCalls := 0
	velocityXCalls := 0
	velocityYCalls := 0
	var sphereArgs []interface{}
	var capsuleArgs []interface{}
	var static2DArgs []interface{}
	var velocityXArgs []interface{}
	var velocityYArgs []interface{}

	v.RegisterForeign("CreateSphere3D", func(args []interface{}) (interface{}, error) {
		sphereCalls++
		sphereArgs = append([]interface{}{}, args...)
		return nil, nil
	})
	v.RegisterForeign("CreateCapsule3D", func(args []interface{}) (interface{}, error) {
		capsuleCalls++
		capsuleArgs = append([]interface{}{}, args...)
		return nil, nil
	})
	v.RegisterForeign("CreateBox2D", func(args []interface{}) (interface{}, error) {
		static2DCalls++
		static2DArgs = append([]interface{}{}, args...)
		return nil, nil
	})
	v.RegisterForeign("GetVelocityX2DByBodyId", func(args []interface{}) (interface{}, error) {
		velocityXCalls++
		velocityXArgs = append([]interface{}{}, args...)
		return 12.5, nil
	})
	v.RegisterForeign("GetVelocityY2DByBodyId", func(args []interface{}) (interface{}, error) {
		velocityYCalls++
		velocityYArgs = append([]interface{}{}, args...)
		return -3.5, nil
	})
	v.RegisterForeign("PhysicsEnable", func(args []interface{}) (interface{}, error) {
		physicsEnableCalls++
		return nil, nil
	})

	registerPhysics(v)

	if _, err := v.CallForeign("MakeSphereCollider", []interface{}{7, 1.25}); err != nil {
		t.Fatalf("MakeSphereCollider failed: %v", err)
	}
	if sphereCalls != 1 {
		t.Fatalf("expected one sphere collider call, got %d", sphereCalls)
	}
	if len(sphereArgs) != 7 || sphereArgs[0] != defaultPhysicsWorld3D || sphereArgs[5] != 1.25 {
		t.Fatalf("unexpected sphere args: %#v", sphereArgs)
	}

	if _, err := v.CallForeign("MakeCapsuleCollider", []interface{}{8, 0.5, 2.0}); err != nil {
		t.Fatalf("MakeCapsuleCollider failed: %v", err)
	}
	if capsuleCalls != 1 {
		t.Fatalf("expected one capsule collider call, got %d", capsuleCalls)
	}
	if len(capsuleArgs) != 8 || capsuleArgs[0] != defaultPhysicsWorld3D || capsuleArgs[5] != 0.5 || capsuleArgs[6] != 2.0 {
		t.Fatalf("unexpected capsule args: %#v", capsuleArgs)
	}

	if _, err := v.CallForeign("MakeStaticBody2D", []interface{}{"wall", 1.0, 2.0, 3.0, 4.0}); err != nil {
		t.Fatalf("MakeStaticBody2D failed: %v", err)
	}
	if static2DCalls != 1 {
		t.Fatalf("expected one static 2D body call, got %d", static2DCalls)
	}
	if len(static2DArgs) != 8 || static2DArgs[0] != defaultPhysicsWorld2D || static2DArgs[1] != "wall" || static2DArgs[6] != 0 || static2DArgs[7] != 0 {
		t.Fatalf("unexpected static 2D args: %#v", static2DArgs)
	}

	gotVX, err := v.CallForeign("GetVelocityX2D", []interface{}{"player"})
	if err != nil {
		t.Fatalf("GetVelocityX2D default-world wrapper failed: %v", err)
	}
	if gotVX != 12.5 || velocityXCalls != 1 || len(velocityXArgs) != 2 || velocityXArgs[0] != defaultPhysicsWorld2D || velocityXArgs[1] != "player" {
		t.Fatalf("unexpected velocity X wrapper result=%v args=%#v calls=%d", gotVX, velocityXArgs, velocityXCalls)
	}

	gotVY, err := v.CallForeign("GetVelocityY2D", []interface{}{"arena", "enemy"})
	if err != nil {
		t.Fatalf("GetVelocityY2D explicit-world wrapper failed: %v", err)
	}
	if gotVY != -3.5 || velocityYCalls != 1 || len(velocityYArgs) != 2 || velocityYArgs[0] != "arena" || velocityYArgs[1] != "enemy" {
		t.Fatalf("unexpected velocity Y wrapper result=%v args=%#v calls=%d", gotVY, velocityYArgs, velocityYCalls)
	}

	if physicsEnableCalls != 2 {
		t.Fatalf("expected PhysicsEnable for 3D collider creation only, got %d", physicsEnableCalls)
	}
}

func Test3DBodyQueryAliasesUseDefaultWorld(t *testing.T) {
	v := vm.NewVM()

	var posXArgs []interface{}
	var posYArgs []interface{}
	var posZArgs []interface{}
	var velXArgs []interface{}
	var velYArgs []interface{}
	var velZArgs []interface{}

	v.RegisterForeign("GetPositionX3D", func(args []interface{}) (interface{}, error) {
		posXArgs = append([]interface{}{}, args...)
		return 10.0, nil
	})
	v.RegisterForeign("GetPositionY3D", func(args []interface{}) (interface{}, error) {
		posYArgs = append([]interface{}{}, args...)
		return 20.0, nil
	})
	v.RegisterForeign("GetPositionZ3D", func(args []interface{}) (interface{}, error) {
		posZArgs = append([]interface{}{}, args...)
		return 30.0, nil
	})
	v.RegisterForeign("GetVelocityX3D", func(args []interface{}) (interface{}, error) {
		velXArgs = append([]interface{}{}, args...)
		return 1.0, nil
	})
	v.RegisterForeign("GetVelocityY3D", func(args []interface{}) (interface{}, error) {
		velYArgs = append([]interface{}{}, args...)
		return 2.0, nil
	})
	v.RegisterForeign("GetVelocityZ3D", func(args []interface{}) (interface{}, error) {
		velZArgs = append([]interface{}{}, args...)
		return 3.0, nil
	})

	registerPhysics(v)

	got, err := v.CallForeign("GetBodyX", []interface{}{"player"})
	if err != nil || got != 10.0 || len(posXArgs) != 2 || posXArgs[0] != defaultPhysicsWorld3D || posXArgs[1] != "player" {
		t.Fatalf("unexpected GetBodyX result=%v err=%v args=%#v", got, err, posXArgs)
	}
	got, err = v.CallForeign("GetBodyY", []interface{}{"player"})
	if err != nil || got != 20.0 || len(posYArgs) != 2 || posYArgs[0] != defaultPhysicsWorld3D || posYArgs[1] != "player" {
		t.Fatalf("unexpected GetBodyY result=%v err=%v args=%#v", got, err, posYArgs)
	}
	got, err = v.CallForeign("GetBodyZ", []interface{}{"player"})
	if err != nil || got != 30.0 || len(posZArgs) != 2 || posZArgs[0] != defaultPhysicsWorld3D || posZArgs[1] != "player" {
		t.Fatalf("unexpected GetBodyZ result=%v err=%v args=%#v", got, err, posZArgs)
	}
	got, err = v.CallForeign("GetBodyVX", []interface{}{"player"})
	if err != nil || got != 1.0 || len(velXArgs) != 2 || velXArgs[0] != defaultPhysicsWorld3D || velXArgs[1] != "player" {
		t.Fatalf("unexpected GetBodyVX result=%v err=%v args=%#v", got, err, velXArgs)
	}
	got, err = v.CallForeign("GetBodyVY", []interface{}{"player"})
	if err != nil || got != 2.0 || len(velYArgs) != 2 || velYArgs[0] != defaultPhysicsWorld3D || velYArgs[1] != "player" {
		t.Fatalf("unexpected GetBodyVY result=%v err=%v args=%#v", got, err, velYArgs)
	}
	got, err = v.CallForeign("GetBodyVZ", []interface{}{"player"})
	if err != nil || got != 3.0 || len(velZArgs) != 2 || velZArgs[0] != defaultPhysicsWorld3D || velZArgs[1] != "player" {
		t.Fatalf("unexpected GetBodyVZ result=%v err=%v args=%#v", got, err, velZArgs)
	}
}

func TestMakeMeshColliderReturnsClearFallbackError(t *testing.T) {
	v := vm.NewVM()
	registerPhysics(v)

	if _, err := v.CallForeign("MakeMeshCollider", []interface{}{1, 2}); err == nil {
		t.Fatalf("expected MakeMeshCollider to return an unsupported-feature error")
	}
}

func TestDeleteBody2DAndDeleteBody3DUseDefaultWorld(t *testing.T) {
	v := vm.NewVM()

	destroy2DCalls := 0
	destroy3DCalls := 0
	var destroy2DArgs []interface{}
	var destroy3DArgs []interface{}

	v.RegisterForeign("DestroyBody2D", func(args []interface{}) (interface{}, error) {
		destroy2DCalls++
		destroy2DArgs = append([]interface{}{}, args...)
		return nil, nil
	})
	v.RegisterForeign("DestroyBody3D", func(args []interface{}) (interface{}, error) {
		destroy3DCalls++
		destroy3DArgs = append([]interface{}{}, args...)
		return nil, nil
	})

	registerPhysics(v)

	if _, err := v.CallForeign("DeleteBody2D", []interface{}{"player"}); err != nil {
		t.Fatalf("DeleteBody2D failed: %v", err)
	}
	if destroy2DCalls != 1 {
		t.Fatalf("expected one DestroyBody2D call, got %d", destroy2DCalls)
	}
	if len(destroy2DArgs) != 2 || destroy2DArgs[0] != defaultPhysicsWorld2D || destroy2DArgs[1] != "player" {
		t.Fatalf("unexpected DestroyBody2D args: %#v", destroy2DArgs)
	}

	if _, err := v.CallForeign("DeleteBody3D", []interface{}{"orb"}); err != nil {
		t.Fatalf("DeleteBody3D failed: %v", err)
	}
	if destroy3DCalls != 1 {
		t.Fatalf("expected one DestroyBody3D call, got %d", destroy3DCalls)
	}
	if len(destroy3DArgs) != 2 || destroy3DArgs[0] != defaultPhysicsWorld3D || destroy3DArgs[1] != "orb" {
		t.Fatalf("unexpected DestroyBody3D args: %#v", destroy3DArgs)
	}
}
