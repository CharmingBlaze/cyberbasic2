// Package dbp - Particles: MakeParticles, SetParticleColor/Size/Speed, EmitParticles.
package dbp

import (
	"fmt"
	"math/rand"
	"strconv"
	"sync"

	"cyberbasic/compiler/vm"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type dbpParticle struct {
	x, y, z     float32
	vx, vy, vz float32
	r, g, b, a uint8
	life        float32
	maxLife     float32
	size        float32
}

type dbpParticleSystem struct {
	particles []dbpParticle
	maxCount  int
	defaultR  uint8
	defaultG  uint8
	defaultB  uint8
	defaultA  uint8
	lifetime   float32
	velX       float32
	velY       float32
	velZ       float32
	size       float32
}

var (
	dbpParticleSystems   = make(map[int]*dbpParticleSystem)
	dbpParticleSystemsMu sync.Mutex
)

func toFloat32Particle(v interface{}) float32 {
	switch x := v.(type) {
	case int:
		return float32(x)
	case float64:
		return float32(x)
	case string:
		f, _ := strconv.ParseFloat(x, 32)
		return float32(f)
	default:
		return 0
	}
}

// registerParticles registers DBP-style particle commands.
func registerParticles(v *vm.VM) {
	v.RegisterForeign("MakeParticles", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("MakeParticles(id) requires 1 argument")
		}
		id := toInt(args[0])
		dbpParticleSystemsMu.Lock()
		dbpParticleSystems[id] = &dbpParticleSystem{
			particles: make([]dbpParticle, 0, 256),
			maxCount:  256,
			defaultR: 255, defaultG: 255, defaultB: 255, defaultA: 255,
			lifetime: 2.0,
			size:     0.1,
		}
		dbpParticleSystemsMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("SetParticleColor", func(args []interface{}) (interface{}, error) {
		if len(args) < 5 {
			return nil, fmt.Errorf("SetParticleColor(id, r, g, b, a) requires 5 arguments")
		}
		id := toInt(args[0])
		r, g, b, a := toUint8(args[1]), toUint8(args[2]), toUint8(args[3]), toUint8(args[4])
		dbpParticleSystemsMu.Lock()
		if ps, ok := dbpParticleSystems[id]; ok {
			ps.defaultR, ps.defaultG, ps.defaultB, ps.defaultA = r, g, b, a
		}
		dbpParticleSystemsMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("SetParticleSize", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("SetParticleSize(id, size) requires 2 arguments")
		}
		id := toInt(args[0])
		size := toFloat32Particle(args[1])
		if size <= 0 {
			size = 0.1
		}
		dbpParticleSystemsMu.Lock()
		if ps, ok := dbpParticleSystems[id]; ok {
			ps.size = size
		}
		dbpParticleSystemsMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("SetParticleSpeed", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("SetParticleSpeed(id, vx, vy, vz) requires 4 arguments")
		}
		id := toInt(args[0])
		vx, vy, vz := toFloat32Particle(args[1]), toFloat32Particle(args[2]), toFloat32Particle(args[3])
		dbpParticleSystemsMu.Lock()
		if ps, ok := dbpParticleSystems[id]; ok {
			ps.velX, ps.velY, ps.velZ = vx, vy, vz
		}
		dbpParticleSystemsMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("SetParticleLifetime", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("SetParticleLifetime(id, seconds) requires 2 arguments")
		}
		id := toInt(args[0])
		life := toFloat32Particle(args[1])
		if life <= 0 {
			life = 1.0
		}
		dbpParticleSystemsMu.Lock()
		if ps, ok := dbpParticleSystems[id]; ok {
			ps.lifetime = life
		}
		dbpParticleSystemsMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("EmitParticlesAt", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("EmitParticlesAt(id, x, y, z [, count]) requires 3+ arguments")
		}
		id := toInt(args[0])
		x, y, z := toFloat32Particle(args[1]), toFloat32Particle(args[2]), toFloat32Particle(args[3])
		count := 10
		if len(args) >= 5 {
			count = toInt(args[4])
		}
		if count < 1 {
			count = 1
		}
		if count > 100 {
			count = 100
		}
		dbpParticleSystemsMu.Lock()
		ps, ok := dbpParticleSystems[id]
		if !ok {
			dbpParticleSystemsMu.Unlock()
			return nil, fmt.Errorf("unknown particle system id %d", id)
		}
		maxCap := ps.maxCount
		if maxCap <= 0 {
			maxCap = 10000
		}
		for i := 0; i < count && len(ps.particles) < maxCap; i++ {
			vx := ps.velX + float32(rand.Float64()*0.2-0.1)
			vy := ps.velY + float32(rand.Float64()*0.2-0.1)
			vz := ps.velZ + float32(rand.Float64()*0.2-0.1)
			life := ps.lifetime * (0.8 + float32(rand.Float64())*0.4)
			ps.particles = append(ps.particles, dbpParticle{
				x: x, y: y, z: z,
				vx: vx, vy: vy, vz: vz,
				r: ps.defaultR, g: ps.defaultG, b: ps.defaultB, a: ps.defaultA,
				life: life, maxLife: life,
				size: ps.size,
			})
		}
		dbpParticleSystemsMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("DrawParticles", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("DrawParticles(id) requires 1 argument")
		}
		id := toInt(args[0])
		drawDbpParticles(id)
		return nil, nil
	})
	v.RegisterForeign("DeleteParticles", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("DeleteParticles(id) requires 1 argument")
		}
		id := toInt(args[0])
		dbpParticleSystemsMu.Lock()
		delete(dbpParticleSystems, id)
		dbpParticleSystemsMu.Unlock()
		return nil, nil
	})

	// 3D particle aliases with maxCount
	v.RegisterForeign("MakeParticles3D", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("MakeParticles3D(id, maxCount) requires 2 arguments")
		}
		id := toInt(args[0])
		maxCount := toInt(args[1])
		if maxCount < 1 {
			maxCount = 256
		}
		if maxCount > 10000 {
			maxCount = 10000
		}
		dbpParticleSystemsMu.Lock()
		ps := &dbpParticleSystem{
			particles: make([]dbpParticle, 0, maxCount),
			maxCount:  maxCount,
			defaultR:  255, defaultG: 255, defaultB: 255, defaultA: 255,
			lifetime: 2.0,
			size:     0.1,
		}
		dbpParticleSystems[id] = ps
		dbpParticleSystemsMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("SetParticles3DColor", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("SetParticles3DColor(id, r, g, b) requires 4 arguments")
		}
		return v.CallForeign("SetParticleColor", []interface{}{args[0], args[1], args[2], args[3], 255})
	})
	v.RegisterForeign("SetParticles3DSize", func(args []interface{}) (interface{}, error) {
		return v.CallForeign("SetParticleSize", args)
	})
	v.RegisterForeign("SetParticles3DSpeed", func(args []interface{}) (interface{}, error) {
		return v.CallForeign("SetParticleSpeed", args)
	})
	v.RegisterForeign("EmitParticles3D", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("EmitParticles3D(id, count [, x, y, z]) requires 2+ arguments")
		}
		id := toInt(args[0])
		count := toInt(args[1])
		x, y, z := float32(0), float32(0), float32(0)
		if len(args) >= 5 {
			x = toFloat32Particle(args[2])
			y = toFloat32Particle(args[3])
			z = toFloat32Particle(args[4])
		}
		return v.CallForeign("EmitParticlesAt", []interface{}{id, x, y, z, count})
	})
	v.RegisterForeign("GetParticles3DCount", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return 0, nil
		}
		id := toInt(args[0])
		dbpParticleSystemsMu.Lock()
		ps, ok := dbpParticleSystems[id]
		dbpParticleSystemsMu.Unlock()
		if !ok {
			return 0, nil
		}
		return len(ps.particles), nil
	})
	v.RegisterForeign("GetParticles3DMax", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return 0, nil
		}
		id := toInt(args[0])
		dbpParticleSystemsMu.Lock()
		ps, ok := dbpParticleSystems[id]
		dbpParticleSystemsMu.Unlock()
		if !ok {
			return 0, nil
		}
		return ps.maxCount, nil
	})
	v.RegisterForeign("DrawParticles3D", func(args []interface{}) (interface{}, error) {
		return v.CallForeign("DrawParticles", args)
	})
}

func toUint8(v interface{}) uint8 {
	switch x := v.(type) {
	case int:
		if x < 0 {
			return 0
		}
		if x > 255 {
			return 255
		}
		return uint8(x)
	case float64:
		if x < 0 {
			return 0
		}
		if x > 255 {
			return 255
		}
		return uint8(x)
	default:
		return 255
	}
}

func drawDbpParticles(id int) {
	dt := rl.GetFrameTime()
	dbpParticleSystemsMu.Lock()
	ps, ok := dbpParticleSystems[id]
	if !ok {
		dbpParticleSystemsMu.Unlock()
		return
	}
	live := ps.particles[:0]
	for _, p := range ps.particles {
		p.x += p.vx * dt
		p.y += p.vy * dt
		p.z += p.vz * dt
		p.life -= dt
		if p.life > 0 {
			live = append(live, p)
		}
	}
	ps.particles = live
	toDraw := make([]dbpParticle, len(live))
	copy(toDraw, live)
	dbpParticleSystemsMu.Unlock()
	for _, p := range toDraw {
		alpha := uint8(float32(p.a) * p.life / p.maxLife)
		c := rl.NewColor(p.r, p.g, p.b, alpha)
		rl.DrawSphere(rl.Vector3{X: p.x, Y: p.y, Z: p.z}, p.size, c)
	}
}
