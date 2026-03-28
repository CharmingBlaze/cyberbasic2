// Package box2ddot exposes global "box2d" as a modfacade over flat Box2D foreigns.
package box2ddot

import (
	"strings"

	"cyberbasic/compiler/bindings/modfacade"
	"cyberbasic/compiler/vm"
)

var box2dNames = []string{
	"Box2DBackendName", "Box2DBackendMode",
	"CreateWorld2D", "Physics2DCreateWorld", "DestroyWorld2D",
	"Step2D", "Physics2DStep", "Physics2DSetGravity", "Physics2DRaycast", "Physics2DSetLayerCollision", "StepAllPhysics2D",
	"CreateBody2D", "DestroyBody2D", "GetBodyCount2D", "GetBodyId2D", "CreateBodyAtScreen2D",
	"CreateBox2D", "CreateCircle2D", "CreatePolygon2D", "CreateEdge2D", "CreateChain2D",
	"SetSensor2D",
	"GetPositionX2D", "GetPositionY2D", "SetPosition2D", "GetAngle2D", "SetAngle2D",
	"GetVelocityX2D", "GetVelocityX2DByBodyId", "GetVelocityY2D", "GetVelocityY2DByBodyId", "SetVelocity2D",
	"ApplyForce2D", "ApplyForce2DByBodyId", "ApplyImpulse2D", "ApplyImpulse2DByBodyId",
	"ApplyTorque2D", "SetAngularVelocity2D", "GetAngularVelocity2D",
	"SetFriction2D", "SetRestitution2D", "SetDamping2D", "SetFixedRotation2D", "SetGravityScale2D", "SetMass2D", "SetBullet2D",
	"CreateDistanceJoint2D", "CreateRevoluteJoint2D", "CreatePrismaticJoint2D", "CreateWeldJoint2D", "CreateRopeJoint2D",
	"CreateWheelJoint2D", "CreatePulleyJoint2D", "CreateGearJoint2D",
	"SetJointLimits2D", "SetJointMotor2D", "DestroyJoint2D",
	"RayCast2D", "RayHitX2D", "RayHitY2D", "RayHitBody2D", "RayHitNormalX2D", "RayHitNormalY2D",
	"GetCollisionCount2D", "GetCollisionOther2D", "GetCollisionNormalX2D", "GetCollisionNormalY2D",
}

func lowerMap(names []string) map[string]string {
	m := make(map[string]string, len(names))
	for _, n := range names {
		m[strings.ToLower(n)] = n
	}
	return m
}

// Register installs global "box2d" after box2d.RegisterBox2D.
func Register(v *vm.VM) {
	v.SetGlobal("box2d", modfacade.New(v, lowerMap(box2dNames)))
}
