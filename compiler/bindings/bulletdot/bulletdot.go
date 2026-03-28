// Package bulletdot exposes global "bullet" as a modfacade over flat 3D physics foreigns.
package bulletdot

import (
	"strings"

	"cyberbasic/compiler/bindings/modfacade"
	"cyberbasic/compiler/vm"
)

var bulletNames = []string{
	"BulletBackendName", "BulletBackendMode", "BulletNativeAvailable", "BulletFeatureAvailable",
	"CreateWorld3D", "DestroyWorld3D", "SetWorldGravity3D", "DestroyBody3D",
	"Step3D", "StepAllPhysics3D",
	"CreateSphere3D", "CreateBox3D", "CreateCapsule3D", "CreateStaticMesh3D", "CreateCylinder3D", "CreateCone3D",
	"CreateHeightmap3D", "CreateCompound3D", "AddShapeToCompound3D",
	"GetPositionX3D", "GetPositionY3D", "GetPositionZ3D", "SetPosition3D",
	"GetYaw3D", "GetPitch3D", "GetRoll3D", "SetRotation3D", "SetScale3D",
	"GetVelocityX3D", "GetVelocityY3D", "GetVelocityZ3D", "SetVelocity3D",
	"SetAngularVelocity3D", "GetAngularVelocityX3D", "GetAngularVelocityY3D", "GetAngularVelocityZ3D",
	"ApplyForce3D", "ApplyImpulse3D", "ApplyTorque3D", "ApplyTorqueImpulse3D",
	"SetFriction3D", "SetRestitution3D", "SetDamping3D", "SetKinematic3D", "SetGravity3D", "SetMass3D", "GetMass3D",
	"SetLinearFactor3D", "SetAngularFactor3D", "SetCCD3D",
	"BulletJointsAvailable", "CreateHingeJoint3D", "CreateSliderJoint3D", "CreateConeTwistJoint3D",
	"CreatePointToPointJoint3D", "CreateFixedJoint3D", "SetJointLimits3D", "SetJointMotor3D",
	"RayCast3D", "RayCastFromDir3D",
	"RayHitX3D", "RayHitY3D", "RayHitZ3D", "RayHitBody3D", "RayHitNormalX3D", "RayHitNormalY3D", "RayHitNormalZ3D",
	"GetCollisionCount3D", "GetCollisionOther3D", "GetCollisionNormalX3D", "GetCollisionNormalY3D", "GetCollisionNormalZ3D",
	"PhysicsEnable", "PhysicsDisable", "PhysicsSetGravity",
	"CreateRigidBody", "ApplyForce", "ApplyImpulse", "SetBodyPosition", "GetBodyPosition", "SetBodyVelocity", "GetBodyVelocity",
	"CheckCollision3D",
}

func lowerMap(names []string) map[string]string {
	m := make(map[string]string, len(names))
	for _, n := range names {
		m[strings.ToLower(n)] = n
	}
	return m
}

// Register installs global "bullet" after bullet.RegisterBullet.
func Register(v *vm.VM) {
	v.SetGlobal("bullet", modfacade.New(v, lowerMap(bulletNames)))
}
