package bullet

import (
	"math"
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

	// BulletJointsAvailable: 1 = PointToPoint and Fixed joints supported
	jointsAvail, err := v.CallForeign("BulletJointsAvailable", nil)
	if err != nil {
		t.Fatalf("BulletJointsAvailable failed: %v", err)
	}
	if jointsAvail != 1 {
		t.Fatalf("expected BulletJointsAvailable 1 (PointToPoint/Fixed), got %v", jointsAvail)
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
		"CreateHingeJoint3D",
		"CreateSliderJoint3D",
		"CreateConeTwistJoint3D",
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

func TestApplyTorqueImpulse3DChangesAngularVelocity(t *testing.T) {
	v := vm.NewVM()
	RegisterBullet(v)
	if _, err := v.CallForeign("CreateWorld3D", []interface{}{"w", 0.0, -9.81, 0.0}); err != nil {
		t.Fatalf("CreateWorld3D failed: %v", err)
	}
	if _, err := v.CallForeign("CreateBox3D", []interface{}{"w", "b", 0.0, 0.0, 0.0, 1.0, 1.0, 1.0, 1.0}); err != nil {
		t.Fatalf("CreateBox3D failed: %v", err)
	}
	if _, err := v.CallForeign("ApplyTorqueImpulse3D", []interface{}{"w", "b", 1.0, 0.0, 0.0}); err != nil {
		t.Fatalf("ApplyTorqueImpulse3D failed: %v", err)
	}
	avx, _ := v.CallForeign("GetAngularVelocityX3D", []interface{}{"w", "b"})
	avy, _ := v.CallForeign("GetAngularVelocityY3D", []interface{}{"w", "b"})
	avz, _ := v.CallForeign("GetAngularVelocityZ3D", []interface{}{"w", "b"})
	if toF(avx) == 0 && toF(avy) == 0 && toF(avz) == 0 {
		t.Fatalf("ApplyTorqueImpulse3D did not change angular velocity: got %v, %v, %v", avx, avy, avz)
	}
}

func toF(v interface{}) float64 {
	switch x := v.(type) {
	case float64:
		return x
	case int:
		return float64(x)
	default:
		return 0
	}
}

func TestPointToPointJointKeepsBodiesConnected(t *testing.T) {
	v := vm.NewVM()
	RegisterBullet(v)
	if _, err := v.CallForeign("CreateWorld3D", []interface{}{"w", 0.0, -9.81, 0.0}); err != nil {
		t.Fatalf("CreateWorld3D failed: %v", err)
	}
	if _, err := v.CallForeign("CreateSphere3D", []interface{}{"w", "a", 0.0, 2.0, 0.0, 0.5, 1.0}); err != nil {
		t.Fatalf("CreateSphere3D failed: %v", err)
	}
	if _, err := v.CallForeign("CreateSphere3D", []interface{}{"w", "b", 0.0, 4.0, 0.0, 0.5, 1.0}); err != nil {
		t.Fatalf("CreateSphere3D failed: %v", err)
	}
	// PointToPoint at body centers (0,0,0) in local space for both
	if _, err := v.CallForeign("CreatePointToPointJoint3D", []interface{}{"w", "j", "a", "b", 0.0, 0.0, 0.0, 0.0, 0.0, 0.0}); err != nil {
		t.Fatalf("CreatePointToPointJoint3D failed: %v", err)
	}
	// Step a few times
	for i := 0; i < 20; i++ {
		if _, err := v.CallForeign("Step3D", []interface{}{"w", 0.016}); err != nil {
			t.Fatalf("Step3D failed: %v", err)
		}
	}
	ax, _ := v.CallForeign("GetPositionX3D", []interface{}{"w", "a"})
	ay, _ := v.CallForeign("GetPositionY3D", []interface{}{"w", "a"})
	az, _ := v.CallForeign("GetPositionZ3D", []interface{}{"w", "a"})
	bx, _ := v.CallForeign("GetPositionX3D", []interface{}{"w", "b"})
	by, _ := v.CallForeign("GetPositionY3D", []interface{}{"w", "b"})
	bz, _ := v.CallForeign("GetPositionZ3D", []interface{}{"w", "b"})
	dist := math.Sqrt(math.Pow(toF(bx)-toF(ax), 2) + math.Pow(toF(by)-toF(ay), 2) + math.Pow(toF(bz)-toF(az), 2))
	// Joint constrains centers to coincide; distance should be near 0 (within solver tolerance)
	if dist > 0.5 {
		t.Fatalf("PointToPoint joint did not keep bodies connected: distance=%v (ax=%v ay=%v az=%v bx=%v by=%v bz=%v)", dist, ax, ay, az, bx, by, bz)
	}
}
