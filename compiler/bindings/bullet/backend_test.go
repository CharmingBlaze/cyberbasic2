package bullet

import (
	"testing"

	"cyberbasic/compiler/vm"
)

func TestBulletBackendMetadata(t *testing.T) {
	v := vm.NewVM()
	RegisterBullet(v)

	name, err := v.CallForeign("BulletBackendName", nil)
	if err != nil {
		t.Fatalf("BulletBackendName failed: %v", err)
	}
	if name != "purego-fallback" {
		t.Fatalf("unexpected backend name: %v", name)
	}

	mode, err := v.CallForeign("BulletBackendMode", nil)
	if err != nil {
		t.Fatalf("BulletBackendMode failed: %v", err)
	}
	if mode != "fallback" {
		t.Fatalf("unexpected backend mode: %v", mode)
	}

	native, err := v.CallForeign("BulletNativeAvailable", nil)
	if err != nil {
		t.Fatalf("BulletNativeAvailable failed: %v", err)
	}
	if native != 0 {
		t.Fatalf("unexpected native flag: %v", native)
	}

	joints, err := v.CallForeign("BulletFeatureAvailable", []interface{}{"joints"})
	if err != nil {
		t.Fatalf("BulletFeatureAvailable(joints) failed: %v", err)
	}
	if joints != 0 {
		t.Fatalf("expected joints to be unavailable, got %v", joints)
	}

	sphere, err := v.CallForeign("BulletFeatureAvailable", []interface{}{"sphere"})
	if err != nil {
		t.Fatalf("BulletFeatureAvailable(sphere) failed: %v", err)
	}
	if sphere != 1 {
		t.Fatalf("expected sphere to be available, got %v", sphere)
	}
}

func TestBulletUnsupportedFeaturesReturnErrors(t *testing.T) {
	v := vm.NewVM()
	RegisterBullet(v)

	for _, name := range []string{
		"CreateHeightmap3D",
		"CreateCompound3D",
		"AddShapeToCompound3D",
		"ApplyTorque3D",
		"ApplyTorqueImpulse3D",
		"CreateHingeJoint3D",
		"CreateSliderJoint3D",
		"CreateConeTwistJoint3D",
		"CreatePointToPointJoint3D",
		"CreateFixedJoint3D",
		"SetJointLimits3D",
		"SetJointMotor3D",
	} {
		if _, err := v.CallForeign(name, []interface{}{}); err == nil {
			t.Fatalf("expected %s to return an unsupported-feature error", name)
		}
	}
}

func TestBulletSimpleBodyHelpers(t *testing.T) {
	v := vm.NewVM()
	RegisterBullet(v)

	if _, err := v.CallForeign("CreateWorld3D", []interface{}{"default", 0.0, -9.81, 0.0}); err != nil {
		t.Fatalf("CreateWorld3D failed: %v", err)
	}
	if _, err := v.CallForeign("CreateSphere3D", []interface{}{"default", "player", 1.0, 2.0, 3.0, 0.5, 1.0}); err != nil {
		t.Fatalf("CreateSphere3D failed: %v", err)
	}
	if _, err := v.CallForeign("SetBodyPosition", []interface{}{"player", 4.0, 5.0, 6.0}); err != nil {
		t.Fatalf("SetBodyPosition failed: %v", err)
	}
	pos, err := v.CallForeign("GetBodyPosition", []interface{}{"player"})
	if err != nil {
		t.Fatalf("GetBodyPosition failed: %v", err)
	}
	gotPos, ok := pos.([]interface{})
	if !ok || len(gotPos) != 3 || gotPos[0] != 4.0 || gotPos[1] != 5.0 || gotPos[2] != 6.0 {
		t.Fatalf("unexpected body position: %#v", pos)
	}
	if _, err := v.CallForeign("SetBodyVelocity", []interface{}{"player", 7.0, 8.0, 9.0}); err != nil {
		t.Fatalf("SetBodyVelocity failed: %v", err)
	}
	vel, err := v.CallForeign("GetBodyVelocity", []interface{}{"player"})
	if err != nil {
		t.Fatalf("GetBodyVelocity failed: %v", err)
	}
	gotVel, ok := vel.([]interface{})
	if !ok || len(gotVel) != 3 || gotVel[0] != 7.0 || gotVel[1] != 8.0 || gotVel[2] != 9.0 {
		t.Fatalf("unexpected body velocity: %#v", vel)
	}
}
