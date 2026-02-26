// Package box2d binds github.com/bytearena/box2d to the CyberBasic VM as BOX2D.*.
// BASIC can call BOX2D.CreateWorld, BOX2D.CreateBody, BOX2D.Step, etc.
package box2d

import (
	"fmt"
	"strconv"
	"strings"
	"sync"

	"cyberbasic/compiler/vm"
	"github.com/bytearena/box2d"
)

func toFloat64(v interface{}) float64 {
	switch x := v.(type) {
	case int:
		return float64(x)
	case float64:
		return x
	case string:
		f, _ := strconv.ParseFloat(x, 64)
		return f
	default:
		return 0
	}
}

func toString(v interface{}) string {
	if v == nil {
		return ""
	}
	return fmt.Sprint(v)
}

func toInt(v interface{}) int {
	switch x := v.(type) {
	case int:
		return x
	case float64:
		return int(x)
	case string:
		n, _ := strconv.Atoi(x)
		return n
	default:
		return 0
	}
}

type collisionHit2D struct {
	otherId string
	nx, ny  float64
}

var (
	worlds             = make(map[string]*box2d.B2World)
	worldIdByPtr       = make(map[*box2d.B2World]string)
	worldMu            sync.RWMutex
	bodies             = make(map[string]*box2d.B2Body) // key: worldId.bodyId
	bodiesMu           sync.RWMutex
	bodyOrder          = make(map[string][]string) // worldId -> ordered body IDs for iteration
	bodyOrderMu        sync.RWMutex
	collisionBuffer2D  = make(map[string][]collisionHit2D)
	collisionBuffer2DMu sync.RWMutex
	lastRay2D          struct {
		hit    bool
		x, y   float64
		bodyId string
		nx, ny float64
	}
	lastRay2DMu sync.Mutex
)

type contactListener struct{}

func (contactListener) BeginContact(contact box2d.B2ContactInterface) {
	fa := contact.GetFixtureA()
	fb := contact.GetFixtureB()
	if fa == nil || fb == nil {
		return
	}
	bodyA := fa.GetBody()
	bodyB := fb.GetBody()
	if bodyA == nil || bodyB == nil {
		return
	}
	world := bodyA.GetWorld()
	if world == nil {
		return
	}
	worldMu.RLock()
	worldId := worldIdByPtr[world]
	worldMu.RUnlock()
	if worldId == "" {
		return
	}
	idA := bodyIdFromBody(worldId, bodyA)
	idB := bodyIdFromBody(worldId, bodyB)
	if idA == "" || idB == "" {
		return
	}
	var wm box2d.B2WorldManifold
	contact.GetWorldManifold(&wm)
	nx, ny := wm.Normal.X, wm.Normal.Y
	collisionBuffer2DMu.Lock()
	keyA := bodyKey(worldId, idA)
	keyB := bodyKey(worldId, idB)
	collisionBuffer2D[keyA] = append(collisionBuffer2D[keyA], collisionHit2D{idB, nx, ny})
	collisionBuffer2D[keyB] = append(collisionBuffer2D[keyB], collisionHit2D{idA, -nx, -ny})
	collisionBuffer2DMu.Unlock()
}

func (contactListener) EndContact(box2d.B2ContactInterface)   {}
func (contactListener) PreSolve(box2d.B2ContactInterface, box2d.B2Manifold) {}
func (contactListener) PostSolve(box2d.B2ContactInterface, *box2d.B2ContactImpulse) {}

func bodyKey(worldId, bodyId string) string { return worldId + "\x00" + bodyId }

// bodyIdFromBody returns the bodyId for the given world and body, or "" if not found.
func bodyIdFromBody(worldId string, body *box2d.B2Body) string {
	if body == nil {
		return ""
	}
	bodiesMu.RLock()
	defer bodiesMu.RUnlock()
	prefix := worldId + "\x00"
	for k, b := range bodies {
		if len(k) > len(prefix) && k[:len(prefix)] == prefix && b == body {
			return k[len(prefix):]
		}
	}
	return ""
}

// RegisterBox2D registers Box2D 2D physics functions with the VM (BOX2D.*).
func RegisterBox2D(v *vm.VM) {
	// World
	v.RegisterForeign("BOX2D.CreateWorld", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("CreateWorld requires (worldId, gravityX, gravityY)")
		}
		worldId := toString(args[0])
		gx, gy := toFloat64(args[1]), toFloat64(args[2])
		gravity := box2d.MakeB2Vec2(gx, gy)
		w := new(box2d.B2World)
		*w = box2d.MakeB2World(gravity)
		w.SetContactListener(contactListener{})
		worldMu.Lock()
		worlds[worldId] = w
		worldIdByPtr[w] = worldId
		worldMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("BOX2D.Step", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("Step requires (worldId, timeStep, velocityIters, positionIters)")
		}
		worldId := toString(args[0])
		dt := toFloat64(args[1])
		velIters := toInt(args[2])
		posIters := toInt(args[3])
		if velIters <= 0 {
			velIters = 8
		}
		if posIters <= 0 {
			posIters = 3
		}
		worldMu.RLock()
		w := worlds[worldId]
		worldMu.RUnlock()
		if w == nil {
			return nil, fmt.Errorf("world not found: %s", worldId)
		}
		collisionBuffer2DMu.Lock()
		for k := range collisionBuffer2D {
			if strings.HasPrefix(k, worldId+"\x00") {
				delete(collisionBuffer2D, k)
			}
		}
		collisionBuffer2DMu.Unlock()
		w.Step(dt, velIters, posIters)
		return nil, nil
	})
	v.RegisterForeign("BOX2D.DestroyWorld", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("DestroyWorld requires (worldId)")
		}
		worldId := toString(args[0])
		worldMu.Lock()
		if w := worlds[worldId]; w != nil {
			delete(worldIdByPtr, w)
		}
		delete(worlds, worldId)
		worldMu.Unlock()
		bodyOrderMu.Lock()
		delete(bodyOrder, worldId)
		bodyOrderMu.Unlock()
		bodiesMu.Lock()
		for k := range bodies {
			if len(k) > len(worldId) && k[:len(worldId)] == worldId && k[len(worldId)] == '\x00' {
				delete(bodies, k)
			}
		}
		bodiesMu.Unlock()
		return nil, nil
	})

	// Body: type 0=static, 1=kinematic, 2=dynamic; shape 0=box, 1=circle
	v.RegisterForeign("BOX2D.CreateBody", func(args []interface{}) (interface{}, error) {
		if len(args) < 7 {
			return nil, fmt.Errorf("CreateBody requires (worldId, bodyId, type, shape, x, y, ...)")
		}
		worldId := toString(args[0])
		bodyId := toString(args[1])
		bodyType := toInt(args[2]) // 0 static, 1 kinematic, 2 dynamic
		shapeKind := toInt(args[3])
		x, y := toFloat64(args[4]), toFloat64(args[5])
		density := toFloat64(args[6])
		if density <= 0 {
			density = 1
		}
		worldMu.RLock()
		w := worlds[worldId]
		worldMu.RUnlock()
		if w == nil {
			return nil, fmt.Errorf("world not found: %s", worldId)
		}
		def := box2d.NewB2BodyDef()
		def.Position = box2d.MakeB2Vec2(x, y)
		switch bodyType {
		case 0:
			def.Type = box2d.B2BodyType.B2_staticBody
		case 1:
			def.Type = box2d.B2BodyType.B2_kinematicBody
		default:
			def.Type = box2d.B2BodyType.B2_dynamicBody
		}
		body := w.CreateBody(def)
		if body == nil {
			return nil, fmt.Errorf("CreateBody failed")
		}
		// Static bodies must use density 0 so mass is zero
		fixtureDensity := density
		if bodyType == 0 {
			fixtureDensity = 0
		}
		var shape box2d.B2ShapeInterface
		if shapeKind == 1 {
			// circle: args can be radius as 7th or we use 1
			radius := 1.0
			if len(args) >= 8 {
				radius = toFloat64(args[7])
			}
			circle := box2d.NewB2CircleShape()
			circle.SetRadius(radius)
			shape = circle
		} else {
			// box: half-width, half-height
			hx, hy := 0.5, 0.5
			if len(args) >= 9 {
				hx = toFloat64(args[7])
				hy = toFloat64(args[8])
			}
			poly := box2d.NewB2PolygonShape()
			poly.SetAsBox(hx, hy)
			shape = poly
		}
		body.CreateFixture(shape, fixtureDensity)
		bodiesMu.Lock()
		bodies[bodyKey(worldId, bodyId)] = body
		bodiesMu.Unlock()
		bodyOrderMu.Lock()
		bodyOrder[worldId] = append(bodyOrder[worldId], bodyId)
		bodyOrderMu.Unlock()
		return bodyId, nil
	})
	v.RegisterForeign("BOX2D.DestroyBody", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("DestroyBody requires (worldId, bodyId)")
		}
		worldId := toString(args[0])
		bodyId := toString(args[1])
		worldMu.RLock()
		w := worlds[worldId]
		worldMu.RUnlock()
		if w == nil {
			return nil, nil
		}
		bodiesMu.Lock()
		b := bodies[bodyKey(worldId, bodyId)]
		delete(bodies, bodyKey(worldId, bodyId))
		bodiesMu.Unlock()
		bodyOrderMu.Lock()
		for i, id := range bodyOrder[worldId] {
			if id == bodyId {
				bodyOrder[worldId] = append(bodyOrder[worldId][:i], bodyOrder[worldId][i+1:]...)
				break
			}
		}
		bodyOrderMu.Unlock()
		if b != nil {
			w.DestroyBody(b)
		}
		return nil, nil
	})

	v.RegisterForeign("BOX2D.GetBodyCount", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("GetBodyCount requires (worldId)")
		}
		worldId := toString(args[0])
		bodyOrderMu.RLock()
		n := len(bodyOrder[worldId])
		bodyOrderMu.RUnlock()
		return n, nil
	})
	v.RegisterForeign("BOX2D.GetBodyId", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("GetBodyId requires (worldId, index)")
		}
		worldId := toString(args[0])
		idx := toInt(args[1])
		bodyOrderMu.RLock()
		order := bodyOrder[worldId]
		bodyOrderMu.RUnlock()
		if idx < 0 || idx >= len(order) {
			return "", nil
		}
		return order[idx], nil
	})

	// CreateBodyAtScreen: create dynamic box at screen (pixel) position; body ID is auto-generated
	v.RegisterForeign("BOX2D.CreateBodyAtScreen", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("CreateBodyAtScreen requires (worldId, screenX, screenY, scale)")
		}
		worldId := toString(args[0])
		screenX := toFloat64(args[1])
		screenY := toFloat64(args[2])
		scale := toFloat64(args[3])
		if scale <= 0 {
			scale = 50
		}
		ox, oy := 400.0, 350.0
		if len(args) >= 6 {
			ox, oy = toFloat64(args[4]), toFloat64(args[5])
		}
		wx := (screenX - ox) / scale
		wy := (oy - screenY) / scale
		bodyOrderMu.Lock()
		n := len(bodyOrder[worldId])
		bodyOrderMu.Unlock()
		bodyId := fmt.Sprintf("box%d", n)
		worldMu.RLock()
		w := worlds[worldId]
		worldMu.RUnlock()
		if w == nil {
			return nil, fmt.Errorf("world not found: %s", worldId)
		}
		def := box2d.NewB2BodyDef()
		def.Position = box2d.MakeB2Vec2(wx, wy)
		def.Type = box2d.B2BodyType.B2_dynamicBody
		body := w.CreateBody(def)
		if body == nil {
			return nil, fmt.Errorf("CreateBody failed")
		}
		poly := box2d.NewB2PolygonShape()
		poly.SetAsBox(0.5, 0.5)
		body.CreateFixture(poly, 1)
		bodiesMu.Lock()
		bodies[bodyKey(worldId, bodyId)] = body
		bodiesMu.Unlock()
		bodyOrderMu.Lock()
		bodyOrder[worldId] = append(bodyOrder[worldId], bodyId)
		bodyOrderMu.Unlock()
		return nil, nil
	})

	// Position / velocity
	v.RegisterForeign("BOX2D.GetPosition", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("GetPosition requires (worldId, bodyId)")
		}
		bodiesMu.RLock()
		b := bodies[bodyKey(toString(args[0]), toString(args[1]))]
		bodiesMu.RUnlock()
		if b == nil {
			return nil, fmt.Errorf("body not found")
		}
		p := b.GetPosition()
		return []interface{}{p.X, p.Y}, nil
	})
	v.RegisterForeign("BOX2D.GetPositionX", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("GetPositionX requires (worldId, bodyId)")
		}
		bodiesMu.RLock()
		b := bodies[bodyKey(toString(args[0]), toString(args[1]))]
		bodiesMu.RUnlock()
		if b == nil {
			return nil, fmt.Errorf("body not found")
		}
		return b.GetPosition().X, nil
	})
	v.RegisterForeign("BOX2D.GetPositionY", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("GetPositionY requires (worldId, bodyId)")
		}
		bodiesMu.RLock()
		b := bodies[bodyKey(toString(args[0]), toString(args[1]))]
		bodiesMu.RUnlock()
		if b == nil {
			return nil, fmt.Errorf("body not found")
		}
		return b.GetPosition().Y, nil
	})
	v.RegisterForeign("BOX2D.GetAngle", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("GetAngle requires (worldId, bodyId)")
		}
		bodiesMu.RLock()
		b := bodies[bodyKey(toString(args[0]), toString(args[1]))]
		bodiesMu.RUnlock()
		if b == nil {
			return nil, fmt.Errorf("body not found")
		}
		return b.GetAngle(), nil
	})
	v.RegisterForeign("BOX2D.SetLinearVelocity", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("SetLinearVelocity requires (worldId, bodyId, vx, vy)")
		}
		bodiesMu.RLock()
		b := bodies[bodyKey(toString(args[0]), toString(args[1]))]
		bodiesMu.RUnlock()
		if b == nil {
			return nil, fmt.Errorf("body not found")
		}
		vx, vy := toFloat64(args[2]), toFloat64(args[3])
		b.SetLinearVelocity(box2d.MakeB2Vec2(vx, vy))
		return nil, nil
	})
	v.RegisterForeign("BOX2D.GetLinearVelocity", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("GetLinearVelocity requires (worldId, bodyId)")
		}
		bodiesMu.RLock()
		b := bodies[bodyKey(toString(args[0]), toString(args[1]))]
		bodiesMu.RUnlock()
		if b == nil {
			return nil, fmt.Errorf("body not found")
		}
		v := b.GetLinearVelocity()
		return []interface{}{v.X, v.Y}, nil
	})
	v.RegisterForeign("BOX2D.SetTransform", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("SetTransform requires (worldId, bodyId, x, y, angle)")
		}
		bodiesMu.RLock()
		b := bodies[bodyKey(toString(args[0]), toString(args[1]))]
		bodiesMu.RUnlock()
		if b == nil {
			return nil, fmt.Errorf("body not found")
		}
		x, y := toFloat64(args[2]), toFloat64(args[3])
		angle := 0.0
		if len(args) >= 5 {
			angle = toFloat64(args[4])
		}
		b.SetTransform(box2d.MakeB2Vec2(x, y), angle)
		return nil, nil
	})
	v.RegisterForeign("BOX2D.ApplyForce", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("ApplyForce requires (worldId, bodyId, fx, fy)")
		}
		bodiesMu.RLock()
		b := bodies[bodyKey(toString(args[0]), toString(args[1]))]
		bodiesMu.RUnlock()
		if b == nil {
			return nil, fmt.Errorf("body not found")
		}
		fx, fy := toFloat64(args[2]), toFloat64(args[3])
		b.ApplyForceToCenter(box2d.MakeB2Vec2(fx, fy), true)
		return nil, nil
	})

	// --- Flat 2D commands (no namespace, case-insensitive via VM) ---
	registerFlat2D(v)
}

// registerFlat2D registers flat CreateWorld2D, Step2D, CreateBox2D, etc. (no BOX2D. prefix).
func registerFlat2D(v *vm.VM) {
	// World
	v.RegisterForeign("CreateWorld2D", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("CreateWorld2D requires (worldName$, gravityX, gravityY)")
		}
		worldId := toString(args[0])
		gx, gy := toFloat64(args[1]), toFloat64(args[2])
		gravity := box2d.MakeB2Vec2(gx, gy)
		w := new(box2d.B2World)
		*w = box2d.MakeB2World(gravity)
		w.SetContactListener(contactListener{})
		worldMu.Lock()
		worlds[worldId] = w
		worldIdByPtr[w] = worldId
		worldMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("DestroyWorld2D", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("DestroyWorld2D requires (worldName$)")
		}
		worldId := toString(args[0])
		worldMu.Lock()
		if w := worlds[worldId]; w != nil {
			delete(worldIdByPtr, w)
		}
		delete(worlds, worldId)
		worldMu.Unlock()
		bodiesMu.Lock()
		for k := range bodies {
			if len(k) > len(worldId) && k[:len(worldId)] == worldId && k[len(worldId)] == '\x00' {
				delete(bodies, k)
			}
		}
		bodiesMu.Unlock()
		bodyOrderMu.Lock()
		delete(bodyOrder, worldId)
		bodyOrderMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("Step2D", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("Step2D requires (worldName$, dt)")
		}
		worldId := toString(args[0])
		dt := toFloat64(args[1])
		worldMu.RLock()
		w := worlds[worldId]
		worldMu.RUnlock()
		if w == nil {
			return nil, fmt.Errorf("world not found: %s", worldId)
		}
		collisionBuffer2DMu.Lock()
		for k := range collisionBuffer2D {
			if strings.HasPrefix(k, worldId+"\x00") {
				delete(collisionBuffer2D, k)
			}
		}
		collisionBuffer2DMu.Unlock()
		w.Step(dt, 8, 3)
		return nil, nil
	})
	v.RegisterForeign("StepAllPhysics2D", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("StepAllPhysics2D requires (dt)")
		}
		dt := toFloat64(args[0])
		worldMu.RLock()
		ids := make([]string, 0, len(worlds))
		for id := range worlds {
			ids = append(ids, id)
		}
		worldMu.RUnlock()
		for _, worldId := range ids {
			worldMu.RLock()
			w := worlds[worldId]
			worldMu.RUnlock()
			if w != nil {
				w.Step(dt, 8, 3)
			}
		}
		return nil, nil
	})

	// Bodies
	v.RegisterForeign("CreateBox2D", func(args []interface{}) (interface{}, error) {
		if len(args) < 8 {
			return nil, fmt.Errorf("CreateBox2D requires (world$, body$, x, y, width, height, mass, isDynamic)")
		}
		worldId := toString(args[0])
		bodyId := toString(args[1])
		x, y := toFloat64(args[2]), toFloat64(args[3])
		width, height := toFloat64(args[4]), toFloat64(args[5])
		mass := toFloat64(args[6])
		isDynamic := toFloat64(args[7]) != 0
		worldMu.RLock()
		w := worlds[worldId]
		worldMu.RUnlock()
		if w == nil {
			return nil, fmt.Errorf("world not found: %s", worldId)
		}
		def := box2d.NewB2BodyDef()
		def.Position = box2d.MakeB2Vec2(x, y)
		if isDynamic {
			def.Type = box2d.B2BodyType.B2_dynamicBody
		} else {
			def.Type = box2d.B2BodyType.B2_staticBody
		}
		body := w.CreateBody(def)
		if body == nil {
			return nil, fmt.Errorf("CreateBody failed")
		}
		density := mass
		if !isDynamic {
			density = 0
		}
		poly := box2d.NewB2PolygonShape()
		poly.SetAsBox(width/2, height/2)
		body.CreateFixture(poly, density)
		bodiesMu.Lock()
		bodies[bodyKey(worldId, bodyId)] = body
		bodiesMu.Unlock()
		bodyOrderMu.Lock()
		bodyOrder[worldId] = append(bodyOrder[worldId], bodyId)
		bodyOrderMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("CreateCircle2D", func(args []interface{}) (interface{}, error) {
		if len(args) < 7 {
			return nil, fmt.Errorf("CreateCircle2D requires (world$, body$, x, y, radius, mass, isDynamic)")
		}
		worldId := toString(args[0])
		bodyId := toString(args[1])
		x, y := toFloat64(args[2]), toFloat64(args[3])
		radius := toFloat64(args[4])
		mass := toFloat64(args[5])
		isDynamic := toFloat64(args[6]) != 0
		worldMu.RLock()
		w := worlds[worldId]
		worldMu.RUnlock()
		if w == nil {
			return nil, fmt.Errorf("world not found: %s", worldId)
		}
		def := box2d.NewB2BodyDef()
		def.Position = box2d.MakeB2Vec2(x, y)
		if isDynamic {
			def.Type = box2d.B2BodyType.B2_dynamicBody
		} else {
			def.Type = box2d.B2BodyType.B2_staticBody
		}
		body := w.CreateBody(def)
		if body == nil {
			return nil, fmt.Errorf("CreateBody failed")
		}
		density := mass
		if !isDynamic {
			density = 0
		}
		circle := box2d.NewB2CircleShape()
		circle.SetRadius(radius)
		body.CreateFixture(circle, density)
		bodiesMu.Lock()
		bodies[bodyKey(worldId, bodyId)] = body
		bodiesMu.Unlock()
		bodyOrderMu.Lock()
		bodyOrder[worldId] = append(bodyOrder[worldId], bodyId)
		bodyOrderMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("CreatePolygon2D", func(args []interface{}) (interface{}, error) {
		if len(args) < 12 {
			return nil, fmt.Errorf("CreatePolygon2D requires (world$, body$, x, y, mass, isDynamic, v1x,v1y, v2x,v2y, v3x,v3y, ...)")
		}
		worldId := toString(args[0])
		bodyId := toString(args[1])
		x, y := toFloat64(args[2]), toFloat64(args[3])
		mass := toFloat64(args[4])
		isDynamic := toFloat64(args[5]) != 0
		n := (len(args) - 6) / 2
		if n < 3 {
			return nil, fmt.Errorf("CreatePolygon2D needs at least 3 vertices")
		}
		verts := make([]box2d.B2Vec2, n)
		for i := 0; i < n; i++ {
			verts[i] = box2d.MakeB2Vec2(toFloat64(args[6+i*2]), toFloat64(args[7+i*2]))
		}
		worldMu.RLock()
		w := worlds[worldId]
		worldMu.RUnlock()
		if w == nil {
			return nil, fmt.Errorf("world not found: %s", worldId)
		}
		def := box2d.NewB2BodyDef()
		def.Position = box2d.MakeB2Vec2(x, y)
		if isDynamic {
			def.Type = box2d.B2BodyType.B2_dynamicBody
		} else {
			def.Type = box2d.B2BodyType.B2_staticBody
		}
		body := w.CreateBody(def)
		if body == nil {
			return nil, fmt.Errorf("CreateBody failed")
		}
		density := mass
		if !isDynamic {
			density = 0
		}
		poly := box2d.NewB2PolygonShape()
		poly.Set(verts, n)
		body.CreateFixture(poly, density)
		bodiesMu.Lock()
		bodies[bodyKey(worldId, bodyId)] = body
		bodiesMu.Unlock()
		bodyOrderMu.Lock()
		bodyOrder[worldId] = append(bodyOrder[worldId], bodyId)
		bodyOrderMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("CreateEdge2D", func(args []interface{}) (interface{}, error) {
		if len(args) < 6 {
			return nil, fmt.Errorf("CreateEdge2D requires (world$, body$, x1, y1, x2, y2)")
		}
		worldId := toString(args[0])
		bodyId := toString(args[1])
		v1 := box2d.MakeB2Vec2(toFloat64(args[2]), toFloat64(args[3]))
		v2 := box2d.MakeB2Vec2(toFloat64(args[4]), toFloat64(args[5]))
		worldMu.RLock()
		w := worlds[worldId]
		worldMu.RUnlock()
		if w == nil {
			return nil, fmt.Errorf("world not found: %s", worldId)
		}
		def := box2d.NewB2BodyDef()
		def.Position = box2d.MakeB2Vec2(0, 0)
		def.Type = box2d.B2BodyType.B2_staticBody
		body := w.CreateBody(def)
		if body == nil {
			return nil, fmt.Errorf("CreateBody failed")
		}
		edge := box2d.NewB2EdgeShape()
		edge.Set(v1, v2)
		body.CreateFixture(edge, 0)
		bodiesMu.Lock()
		bodies[bodyKey(worldId, bodyId)] = body
		bodiesMu.Unlock()
		bodyOrderMu.Lock()
		bodyOrder[worldId] = append(bodyOrder[worldId], bodyId)
		bodyOrderMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("CreateChain2D", func(args []interface{}) (interface{}, error) {
		return nil, nil
	})
	v.RegisterForeign("SetSensor2D", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("SetSensor2D requires (world$, body$, sensor)")
		}
		bodiesMu.RLock()
		b := bodies[bodyKey(toString(args[0]), toString(args[1]))]
		bodiesMu.RUnlock()
		if b == nil {
			return nil, nil
		}
		sensor := toFloat64(args[2]) != 0
		for f := b.GetFixtureList(); f != nil; f = f.GetNext() {
			f.SetSensor(sensor)
		}
		return nil, nil
	})

	// Position
	v.RegisterForeign("GetPositionX2D", func(args []interface{}) (interface{}, error) {
		bodiesMu.RLock()
		b := bodies[bodyKey(toString(args[0]), toString(args[1]))]
		bodiesMu.RUnlock()
		if b == nil {
			return 0.0, nil
		}
		return b.GetPosition().X, nil
	})
	v.RegisterForeign("GetPositionY2D", func(args []interface{}) (interface{}, error) {
		bodiesMu.RLock()
		b := bodies[bodyKey(toString(args[0]), toString(args[1]))]
		bodiesMu.RUnlock()
		if b == nil {
			return 0.0, nil
		}
		return b.GetPosition().Y, nil
	})
	v.RegisterForeign("SetPosition2D", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("SetPosition2D requires (world$, body$, x, y)")
		}
		bodiesMu.RLock()
		b := bodies[bodyKey(toString(args[0]), toString(args[1]))]
		bodiesMu.RUnlock()
		if b == nil {
			return nil, fmt.Errorf("body not found")
		}
		b.SetTransform(box2d.MakeB2Vec2(toFloat64(args[2]), toFloat64(args[3])), b.GetAngle())
		return nil, nil
	})

	// Rotation
	v.RegisterForeign("GetAngle2D", func(args []interface{}) (interface{}, error) {
		bodiesMu.RLock()
		b := bodies[bodyKey(toString(args[0]), toString(args[1]))]
		bodiesMu.RUnlock()
		if b == nil {
			return 0.0, nil
		}
		return b.GetAngle(), nil
	})
	v.RegisterForeign("SetAngle2D", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("SetAngle2D requires (world$, body$, angle)")
		}
		bodiesMu.RLock()
		b := bodies[bodyKey(toString(args[0]), toString(args[1]))]
		bodiesMu.RUnlock()
		if b == nil {
			return nil, fmt.Errorf("body not found")
		}
		b.SetTransform(b.GetPosition(), toFloat64(args[2]))
		return nil, nil
	})

	// Velocity
	v.RegisterForeign("GetVelocityX2D", func(args []interface{}) (interface{}, error) {
		bodiesMu.RLock()
		b := bodies[bodyKey(toString(args[0]), toString(args[1]))]
		bodiesMu.RUnlock()
		if b == nil {
			return 0.0, nil
		}
		return b.GetLinearVelocity().X, nil
	})
	v.RegisterForeign("GetVelocityY2D", func(args []interface{}) (interface{}, error) {
		bodiesMu.RLock()
		b := bodies[bodyKey(toString(args[0]), toString(args[1]))]
		bodiesMu.RUnlock()
		if b == nil {
			return 0.0, nil
		}
		return b.GetLinearVelocity().Y, nil
	})
	v.RegisterForeign("SetVelocity2D", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("SetVelocity2D requires (world$, body$, vx, vy)")
		}
		bodiesMu.RLock()
		b := bodies[bodyKey(toString(args[0]), toString(args[1]))]
		bodiesMu.RUnlock()
		if b == nil {
			return nil, fmt.Errorf("body not found")
		}
		b.SetLinearVelocity(box2d.MakeB2Vec2(toFloat64(args[2]), toFloat64(args[3])))
		return nil, nil
	})

	// Forces
	v.RegisterForeign("ApplyForce2D", func(args []interface{}) (interface{}, error) {
		if len(args) < 5 {
			return nil, fmt.Errorf("ApplyForce2D requires (world$, body$, fx, fy)")
		}
		bodiesMu.RLock()
		b := bodies[bodyKey(toString(args[0]), toString(args[1]))]
		bodiesMu.RUnlock()
		if b == nil {
			return nil, fmt.Errorf("body not found")
		}
		b.ApplyForceToCenter(box2d.MakeB2Vec2(toFloat64(args[2]), toFloat64(args[3])), true)
		return nil, nil
	})
	v.RegisterForeign("ApplyImpulse2D", func(args []interface{}) (interface{}, error) {
		if len(args) < 5 {
			return nil, fmt.Errorf("ApplyImpulse2D requires (world$, body$, ix, iy)")
		}
		bodiesMu.RLock()
		b := bodies[bodyKey(toString(args[0]), toString(args[1]))]
		bodiesMu.RUnlock()
		if b == nil {
			return nil, fmt.Errorf("body not found")
		}
		b.ApplyLinearImpulseToCenter(box2d.MakeB2Vec2(toFloat64(args[2]), toFloat64(args[3])), true)
		return nil, nil
	})
	v.RegisterForeign("ApplyTorque2D", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("ApplyTorque2D requires (world$, body$, torque)")
		}
		bodiesMu.RLock()
		b := bodies[bodyKey(toString(args[0]), toString(args[1]))]
		bodiesMu.RUnlock()
		if b == nil {
			return nil, nil
		}
		b.ApplyTorque(toFloat64(args[2]), true)
		return nil, nil
	})
	v.RegisterForeign("SetAngularVelocity2D", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("SetAngularVelocity2D requires (world$, body$, omega)")
		}
		bodiesMu.RLock()
		b := bodies[bodyKey(toString(args[0]), toString(args[1]))]
		bodiesMu.RUnlock()
		if b == nil {
			return nil, nil
		}
		b.SetAngularVelocity(toFloat64(args[2]))
		return nil, nil
	})
	v.RegisterForeign("GetAngularVelocity2D", func(args []interface{}) (interface{}, error) {
		bodiesMu.RLock()
		b := bodies[bodyKey(toString(args[0]), toString(args[1]))]
		bodiesMu.RUnlock()
		if b == nil {
			return 0.0, nil
		}
		return b.GetAngularVelocity(), nil
	})

	// Body properties
	v.RegisterForeign("SetFriction2D", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("SetFriction2D requires (world$, body$, friction)")
		}
		bodiesMu.RLock()
		b := bodies[bodyKey(toString(args[0]), toString(args[1]))]
		bodiesMu.RUnlock()
		if b == nil {
			return nil, nil
		}
		for f := b.GetFixtureList(); f != nil; f = f.GetNext() {
			f.SetFriction(toFloat64(args[2]))
		}
		return nil, nil
	})
	v.RegisterForeign("SetRestitution2D", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("SetRestitution2D requires (world$, body$, bounce)")
		}
		bodiesMu.RLock()
		b := bodies[bodyKey(toString(args[0]), toString(args[1]))]
		bodiesMu.RUnlock()
		if b == nil {
			return nil, nil
		}
		for f := b.GetFixtureList(); f != nil; f = f.GetNext() {
			f.SetRestitution(toFloat64(args[2]))
		}
		return nil, nil
	})
	v.RegisterForeign("SetDamping2D", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("SetDamping2D requires (world$, body$, linearDamp, angularDamp)")
		}
		bodiesMu.RLock()
		b := bodies[bodyKey(toString(args[0]), toString(args[1]))]
		bodiesMu.RUnlock()
		if b == nil {
			return nil, nil
		}
		b.SetLinearDamping(toFloat64(args[2]))
		b.SetAngularDamping(toFloat64(args[3]))
		return nil, nil
	})
	v.RegisterForeign("SetFixedRotation2D", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("SetFixedRotation2D requires (world$, body$, fixed)")
		}
		bodiesMu.RLock()
		b := bodies[bodyKey(toString(args[0]), toString(args[1]))]
		bodiesMu.RUnlock()
		if b == nil {
			return nil, nil
		}
		b.SetFixedRotation(toFloat64(args[2]) != 0)
		return nil, nil
	})
	v.RegisterForeign("SetGravityScale2D", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("SetGravityScale2D requires (world$, body$, scale)")
		}
		bodiesMu.RLock()
		b := bodies[bodyKey(toString(args[0]), toString(args[1]))]
		bodiesMu.RUnlock()
		if b == nil {
			return nil, nil
		}
		b.SetGravityScale(toFloat64(args[2]))
		return nil, nil
	})
	v.RegisterForeign("SetMass2D", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("SetMass2D requires (world$, body$, mass)")
		}
		bodiesMu.RLock()
		b := bodies[bodyKey(toString(args[0]), toString(args[1]))]
		bodiesMu.RUnlock()
		if b == nil {
			return nil, nil
		}
		for f := b.GetFixtureList(); f != nil; f = f.GetNext() {
			f.SetDensity(toFloat64(args[2]))
		}
		b.ResetMassData()
		return nil, nil
	})
	v.RegisterForeign("SetBullet2D", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("SetBullet2D requires (world$, body$, bullet)")
		}
		bodiesMu.RLock()
		b := bodies[bodyKey(toString(args[0]), toString(args[1]))]
		bodiesMu.RUnlock()
		if b == nil {
			return nil, nil
		}
		b.SetBullet(toFloat64(args[2]) != 0)
		return nil, nil
	})

	// Joints: CreateDistanceJoint2D implemented; others stubbed
	v.RegisterForeign("CreateDistanceJoint2D", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("CreateDistanceJoint2D requires (worldId, bodyAId, bodyBId, length)")
		}
		worldId := toString(args[0])
		bodyAId := toString(args[1])
		bodyBId := toString(args[2])
		length := toFloat64(args[3])
		if length <= 0 {
			length = 1
		}
		worldMu.RLock()
		w := worlds[worldId]
		worldMu.RUnlock()
		if w == nil {
			return nil, fmt.Errorf("world not found: %s", worldId)
		}
		bodiesMu.RLock()
		bodyA := bodies[bodyKey(worldId, bodyAId)]
		bodyB := bodies[bodyKey(worldId, bodyBId)]
		bodiesMu.RUnlock()
		if bodyA == nil {
			return nil, fmt.Errorf("body not found: %s", bodyAId)
		}
		if bodyB == nil {
			return nil, fmt.Errorf("body not found: %s", bodyBId)
		}
		anchorA := bodyA.GetPosition()
		anchorB := bodyB.GetPosition()
		def := box2d.MakeB2DistanceJointDef()
		def.Initialize(bodyA, bodyB, anchorA, anchorB)
		def.Length = length
		w.CreateJoint(&def)
		return nil, nil
	})
	v.RegisterForeign("CreateRevoluteJoint2D", func(args []interface{}) (interface{}, error) { return nil, nil })
	v.RegisterForeign("CreatePrismaticJoint2D", func(args []interface{}) (interface{}, error) { return nil, nil })
	v.RegisterForeign("CreatePulleyJoint2D", func(args []interface{}) (interface{}, error) { return nil, nil })
	v.RegisterForeign("CreateGearJoint2D", func(args []interface{}) (interface{}, error) { return nil, nil })
	v.RegisterForeign("CreateWeldJoint2D", func(args []interface{}) (interface{}, error) { return nil, nil })
	v.RegisterForeign("CreateRopeJoint2D", func(args []interface{}) (interface{}, error) { return nil, nil })
	v.RegisterForeign("CreateWheelJoint2D", func(args []interface{}) (interface{}, error) { return nil, nil })
	v.RegisterForeign("SetJointLimits2D", func(args []interface{}) (interface{}, error) { return nil, nil })
	v.RegisterForeign("SetJointMotor2D", func(args []interface{}) (interface{}, error) { return nil, nil })

	// Raycast
	v.RegisterForeign("RayCast2D", func(args []interface{}) (interface{}, error) {
		if len(args) < 5 {
			return nil, fmt.Errorf("RayCast2D requires (world$, fromX, fromY, toX, toY)")
		}
		worldId := toString(args[0])
		fx, fy := toFloat64(args[1]), toFloat64(args[2])
		tx, ty := toFloat64(args[3]), toFloat64(args[4])
		worldMu.RLock()
		w := worlds[worldId]
		worldMu.RUnlock()
		if w == nil {
			lastRay2DMu.Lock()
			lastRay2D.hit = false
			lastRay2DMu.Unlock()
			return 0, nil
		}
		p1 := box2d.MakeB2Vec2(fx, fy)
		p2 := box2d.MakeB2Vec2(tx, ty)
		var hitX, hitY, hitNx, hitNy float64
		var hitBodyId string
		hit := false
		callback := func(fixture *box2d.B2Fixture, point box2d.B2Vec2, normal box2d.B2Vec2, fraction float64) float64 {
			hitX, hitY = point.X, point.Y
			hitNx, hitNy = normal.X, normal.Y
			hitBodyId = bodyIdFromBody(worldId, fixture.GetBody())
			hit = true
			return fraction
		}
		w.RayCast(box2d.B2RaycastCallback(callback), p1, p2)
		lastRay2DMu.Lock()
		lastRay2D.hit = hit
		lastRay2D.x, lastRay2D.y = hitX, hitY
		lastRay2D.bodyId = hitBodyId
		lastRay2D.nx, lastRay2D.ny = hitNx, hitNy
		lastRay2DMu.Unlock()
		if hit {
			return 1, nil
		}
		return 0, nil
	})
	v.RegisterForeign("RayHitX2D", func(args []interface{}) (interface{}, error) {
		lastRay2DMu.Lock()
		defer lastRay2DMu.Unlock()
		return lastRay2D.x, nil
	})
	v.RegisterForeign("RayHitY2D", func(args []interface{}) (interface{}, error) {
		lastRay2DMu.Lock()
		defer lastRay2DMu.Unlock()
		return lastRay2D.y, nil
	})
	v.RegisterForeign("RayHitBody2D", func(args []interface{}) (interface{}, error) {
		lastRay2DMu.Lock()
		defer lastRay2DMu.Unlock()
		return lastRay2D.bodyId, nil
	})
	v.RegisterForeign("RayHitNormalX2D", func(args []interface{}) (interface{}, error) {
		lastRay2DMu.Lock()
		defer lastRay2DMu.Unlock()
		return lastRay2D.nx, nil
	})
	v.RegisterForeign("RayHitNormalY2D", func(args []interface{}) (interface{}, error) {
		lastRay2DMu.Lock()
		defer lastRay2DMu.Unlock()
		return lastRay2D.ny, nil
	})

	// Collision events
	v.RegisterForeign("GetCollisionCount2D", func(args []interface{}) (interface{}, error) {
		key := bodyKey(toString(args[0]), toString(args[1]))
		collisionBuffer2DMu.RLock()
		n := len(collisionBuffer2D[key])
		collisionBuffer2DMu.RUnlock()
		return n, nil
	})
	v.RegisterForeign("GetCollisionOther2D", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return "", nil
		}
		key := bodyKey(toString(args[0]), toString(args[1]))
		idx := toInt(args[2])
		collisionBuffer2DMu.RLock()
		hits := collisionBuffer2D[key]
		collisionBuffer2DMu.RUnlock()
		if idx < 0 || idx >= len(hits) {
			return "", nil
		}
		return hits[idx].otherId, nil
	})
	v.RegisterForeign("GetCollisionNormalX2D", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return 0.0, nil
		}
		key := bodyKey(toString(args[0]), toString(args[1]))
		idx := toInt(args[2])
		collisionBuffer2DMu.RLock()
		hits := collisionBuffer2D[key]
		collisionBuffer2DMu.RUnlock()
		if idx < 0 || idx >= len(hits) {
			return 0.0, nil
		}
		return hits[idx].nx, nil
	})
	v.RegisterForeign("GetCollisionNormalY2D", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return 0.0, nil
		}
		key := bodyKey(toString(args[0]), toString(args[1]))
		idx := toInt(args[2])
		collisionBuffer2DMu.RLock()
		hits := collisionBuffer2D[key]
		collisionBuffer2DMu.RUnlock()
		if idx < 0 || idx >= len(hits) {
			return 0.0, nil
		}
		return hits[idx].ny, nil
	})
}

// Exported API for --gen-go (same names as VM: box2d.CreateWorld, etc.)

func CreateWorld(worldId string, gx, gy float64) {
	gravity := box2d.MakeB2Vec2(gx, gy)
	w := new(box2d.B2World)
	*w = box2d.MakeB2World(gravity)
	w.SetContactListener(contactListener{})
	worldMu.Lock()
	worlds[worldId] = w
	worldIdByPtr[w] = worldId
	worldMu.Unlock()
}

func Step(worldId string, dt float64, velocityIters, positionIters int) {
	if velocityIters <= 0 {
		velocityIters = 8
	}
	if positionIters <= 0 {
		positionIters = 3
	}
	worldMu.RLock()
	w := worlds[worldId]
	worldMu.RUnlock()
	if w != nil {
		w.Step(dt, velocityIters, positionIters)
	}
}

func DestroyWorld(worldId string) {
	worldMu.Lock()
	delete(worlds, worldId)
	worldMu.Unlock()
	bodiesMu.Lock()
	for k := range bodies {
		if len(k) > len(worldId) && k[:len(worldId)] == worldId && k[len(worldId)] == '\x00' {
			delete(bodies, k)
		}
	}
	bodiesMu.Unlock()
}

func CreateBody(worldId, bodyId string, bodyType, shapeKind int, x, y, density float64, extra ...float64) {
	worldMu.RLock()
	w := worlds[worldId]
	worldMu.RUnlock()
	if w == nil {
		return
	}
	def := box2d.NewB2BodyDef()
	def.Position = box2d.MakeB2Vec2(x, y)
	switch bodyType {
	case 0:
		def.Type = box2d.B2BodyType.B2_staticBody
	case 1:
		def.Type = box2d.B2BodyType.B2_kinematicBody
	default:
		def.Type = box2d.B2BodyType.B2_dynamicBody
	}
	body := w.CreateBody(def)
	if body == nil {
		return
	}
	fixtureDensity := density
	if bodyType == 0 {
		fixtureDensity = 0
	}
	if shapeKind == 1 {
		radius := 1.0
		if len(extra) >= 1 {
			radius = extra[0]
		}
		circle := box2d.NewB2CircleShape()
		circle.SetRadius(radius)
		body.CreateFixture(circle, fixtureDensity)
	} else {
		hx, hy := 0.5, 0.5
		if len(extra) >= 2 {
			hx, hy = extra[0], extra[1]
		}
		poly := box2d.NewB2PolygonShape()
		poly.SetAsBox(hx, hy)
		body.CreateFixture(poly, fixtureDensity)
	}
	bodiesMu.Lock()
	bodies[bodyKey(worldId, bodyId)] = body
	bodiesMu.Unlock()
}

func DestroyBody(worldId, bodyId string) {
	worldMu.RLock()
	w := worlds[worldId]
	worldMu.RUnlock()
	if w == nil {
		return
	}
	bodiesMu.Lock()
	b := bodies[bodyKey(worldId, bodyId)]
	delete(bodies, bodyKey(worldId, bodyId))
	bodiesMu.Unlock()
	if b != nil {
		w.DestroyBody(b)
	}
}

func GetPosition(worldId, bodyId string) (float64, float64) {
	bodiesMu.RLock()
	b := bodies[bodyKey(worldId, bodyId)]
	bodiesMu.RUnlock()
	if b == nil {
		return 0, 0
	}
	p := b.GetPosition()
	return p.X, p.Y
}

// GetPositionX returns the X position of the body (world coordinates).
func GetPositionX(worldId, bodyId string) float64 {
	x, _ := GetPosition(worldId, bodyId)
	return x
}

// GetPositionY returns the Y position of the body (world coordinates).
func GetPositionY(worldId, bodyId string) float64 {
	_, y := GetPosition(worldId, bodyId)
	return y
}

func GetAngle(worldId, bodyId string) float64 {
	bodiesMu.RLock()
	b := bodies[bodyKey(worldId, bodyId)]
	bodiesMu.RUnlock()
	if b == nil {
		return 0
	}
	return b.GetAngle()
}

func SetLinearVelocity(worldId, bodyId string, vx, vy float64) {
	bodiesMu.RLock()
	b := bodies[bodyKey(worldId, bodyId)]
	bodiesMu.RUnlock()
	if b != nil {
		b.SetLinearVelocity(box2d.MakeB2Vec2(vx, vy))
	}
}

func GetLinearVelocity(worldId, bodyId string) (float64, float64) {
	bodiesMu.RLock()
	b := bodies[bodyKey(worldId, bodyId)]
	bodiesMu.RUnlock()
	if b == nil {
		return 0, 0
	}
	v := b.GetLinearVelocity()
	return v.X, v.Y
}

func SetTransform(worldId, bodyId string, x, y, angle float64) {
	bodiesMu.RLock()
	b := bodies[bodyKey(worldId, bodyId)]
	bodiesMu.RUnlock()
	if b != nil {
		b.SetTransform(box2d.MakeB2Vec2(x, y), angle)
	}
}

// RayCastQuery performs a raycast in the world; returns true if hit, and the hit point/body/normal.
// Used by IsOnGround2D and other game helpers.
func RayCastQuery(worldId string, fromX, fromY, toX, toY float64) (hit bool, hitX, hitY float64, hitBodyId string, normX, normY float64) {
	worldMu.RLock()
	w := worlds[worldId]
	worldMu.RUnlock()
	if w == nil {
		return false, 0, 0, "", 0, 0
	}
	p1 := box2d.MakeB2Vec2(fromX, fromY)
	p2 := box2d.MakeB2Vec2(toX, toY)
	callback := func(fixture *box2d.B2Fixture, point box2d.B2Vec2, normal box2d.B2Vec2, fraction float64) float64 {
		hitX, hitY = point.X, point.Y
		normX, normY = normal.X, normal.Y
		hitBodyId = bodyIdFromBody(worldId, fixture.GetBody())
		hit = true
		return fraction
	}
	w.RayCast(box2d.B2RaycastCallback(callback), p1, p2)
	return hit, hitX, hitY, hitBodyId, normX, normY
}

// GetCollisionCountForBody returns the number of current collisions for the body (for use by GAME.ProcessCollisions2D).
func GetCollisionCountForBody(worldId, bodyId string) int {
	key := bodyKey(worldId, bodyId)
	collisionBuffer2DMu.RLock()
	n := len(collisionBuffer2D[key])
	collisionBuffer2DMu.RUnlock()
	return n
}

// GetCollisionOtherForBody returns the other body id for the collision at index.
func GetCollisionOtherForBody(worldId, bodyId string, index int) string {
	key := bodyKey(worldId, bodyId)
	collisionBuffer2DMu.RLock()
	hits := collisionBuffer2D[key]
	collisionBuffer2DMu.RUnlock()
	if index < 0 || index >= len(hits) {
		return ""
	}
	return hits[index].otherId
}
