// Package raylib: 2D particle emitter (texture-based quads).
package raylib

import (
	"fmt"
	"math"
	"math/rand"
	"sync"

	"cyberbasic/compiler/vm"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type particle2D struct {
	X, Y    float32
	VX, VY  float32
	Life    float32
	MaxLife float32
	R, G, B, A uint8
}

type particleEmitter2D struct {
	TextureID   string
	Rate        float32
	LifetimeMin float32
	LifetimeMax float32
	VelX, VelY  float32
	Spread      float32 // radians
	ColorR, ColorG, ColorB, ColorA float32
	Particles   []particle2D
	Accum       float32
	LayerID     string
	ZIndex      int
}

var (
	particleEmitters2D   = make(map[string]*particleEmitter2D)
	particleEmitter2DSeq int
	particleEmitter2DMu  sync.Mutex
)

func registerParticles2D(v *vm.VM) {
	v.RegisterForeign("ParticleEmitterCreate", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("ParticleEmitterCreate requires (textureId)")
		}
		texID := toString(args[0])
		texMu.Lock()
		_, ok := textures[texID]
		texMu.Unlock()
		if !ok {
			return nil, fmt.Errorf("unknown texture id: %s", texID)
		}
		particleEmitter2DMu.Lock()
		particleEmitter2DSeq++
		id := fmt.Sprintf("pe2d_%d", particleEmitter2DSeq)
		particleEmitters2D[id] = &particleEmitter2D{
			TextureID:   texID,
			Rate:        10,
			LifetimeMin: 0.5,
			LifetimeMax: 1.5,
			Spread:      float32(math.Pi / 4),
			ColorR: 1, ColorG: 1, ColorB: 1, ColorA: 1,
		}
		particleEmitter2DMu.Unlock()
		return id, nil
	})
	v.RegisterForeign("ParticleEmitterSetRate", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("ParticleEmitterSetRate requires (emitterId, rate)")
		}
		id := toString(args[0])
		particleEmitter2DMu.Lock()
		defer particleEmitter2DMu.Unlock()
		e := particleEmitters2D[id]
		if e == nil {
			return nil, fmt.Errorf("unknown emitter: %s", id)
		}
		e.Rate = toFloat32(args[1])
		return nil, nil
	})
	v.RegisterForeign("ParticleEmitterSetLifetime", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("ParticleEmitterSetLifetime requires (emitterId, min, max)")
		}
		id := toString(args[0])
		particleEmitter2DMu.Lock()
		defer particleEmitter2DMu.Unlock()
		e := particleEmitters2D[id]
		if e == nil {
			return nil, fmt.Errorf("unknown emitter: %s", id)
		}
		e.LifetimeMin = toFloat32(args[1])
		e.LifetimeMax = toFloat32(args[2])
		return nil, nil
	})
	v.RegisterForeign("ParticleEmitterSetVelocity", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("ParticleEmitterSetVelocity requires (emitterId, vx, vy)")
		}
		id := toString(args[0])
		particleEmitter2DMu.Lock()
		defer particleEmitter2DMu.Unlock()
		e := particleEmitters2D[id]
		if e == nil {
			return nil, fmt.Errorf("unknown emitter: %s", id)
		}
		e.VelX = toFloat32(args[1])
		e.VelY = toFloat32(args[2])
		return nil, nil
	})
	v.RegisterForeign("ParticleEmitterSetSpread", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("ParticleEmitterSetSpread requires (emitterId, angleRad)")
		}
		id := toString(args[0])
		particleEmitter2DMu.Lock()
		defer particleEmitter2DMu.Unlock()
		e := particleEmitters2D[id]
		if e == nil {
			return nil, fmt.Errorf("unknown emitter: %s", id)
		}
		e.Spread = toFloat32(args[1])
		return nil, nil
	})
	v.RegisterForeign("ParticleEmitterSetColor", func(args []interface{}) (interface{}, error) {
		if len(args) < 5 {
			return nil, fmt.Errorf("ParticleEmitterSetColor requires (emitterId, r, g, b, a)")
		}
		id := toString(args[0])
		particleEmitter2DMu.Lock()
		defer particleEmitter2DMu.Unlock()
		e := particleEmitters2D[id]
		if e == nil {
			return nil, fmt.Errorf("unknown emitter: %s", id)
		}
		e.ColorR = toFloat32(args[1])
		e.ColorG = toFloat32(args[2])
		e.ColorB = toFloat32(args[3])
		e.ColorA = toFloat32(args[4])
		return nil, nil
	})
	v.RegisterForeign("ParticleEmitterSetLayer", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("ParticleEmitterSetLayer requires (emitterId, layerId)")
		}
		id := toString(args[0])
		particleEmitter2DMu.Lock()
		defer particleEmitter2DMu.Unlock()
		e := particleEmitters2D[id]
		if e == nil {
			return nil, fmt.Errorf("unknown emitter: %s", id)
		}
		e.LayerID = toString(args[1])
		return nil, nil
	})
	v.RegisterForeign("DrawParticleEmitter", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("DrawParticleEmitter requires (emitterId)")
		}
		id := toString(args[0])
		drawParticleEmitter2D(id)
		return nil, nil
	})
}

func drawParticleEmitter2D(id string) {
	particleEmitter2DMu.Lock()
	e := particleEmitters2D[id]
	if e == nil {
		particleEmitter2DMu.Unlock()
		return
	}
	dt := rl.GetFrameTime()
	e.Accum += e.Rate * dt
	for e.Accum >= 1 {
		e.Accum -= 1
		life := e.LifetimeMin
		if e.LifetimeMax > e.LifetimeMin {
			life = e.LifetimeMin + float32(randFloat())*((e.LifetimeMax-e.LifetimeMin))
		}
		angle := -e.Spread/2 + float32(randFloat())*e.Spread
		vx := e.VelX*float32(math.Cos(float64(angle))) - e.VelY*float32(math.Sin(float64(angle)))
		vy := e.VelX*float32(math.Sin(float64(angle))) + e.VelY*float32(math.Cos(float64(angle)))
		e.Particles = append(e.Particles, particle2D{
			X: 0, Y: 0, VX: vx, VY: vy,
			Life: life, MaxLife: life,
			R: uint8(e.ColorR * 255), G: uint8(e.ColorG * 255), B: uint8(e.ColorB * 255), A: uint8(e.ColorA * 255),
		})
	}
	live := e.Particles[:0]
	for i := range e.Particles {
		p := &e.Particles[i]
		p.X += p.VX * dt
		p.Y += p.VY * dt
		p.Life -= dt
		if p.Life > 0 {
			live = append(live, *p)
		}
	}
	e.Particles = live
	toDraw := make([]particle2D, len(e.Particles))
	copy(toDraw, e.Particles)
	particleEmitter2DMu.Unlock()
	texMu.Lock()
	tex := textures[e.TextureID]
	texMu.Unlock()
	if tex.ID == 0 {
		return
	}
	for _, p := range toDraw {
		alpha := uint8(float32(p.A) * (p.Life / p.MaxLife))
		tint := rl.NewColor(p.R, p.G, p.B, alpha)
		rl.DrawTextureEx(tex, rl.Vector2{X: p.X, Y: p.Y}, 0, 1, tint)
	}
}

func randFloat() float64 {
	seedableRandMu.Lock()
	defer seedableRandMu.Unlock()
	if seedableRand != nil {
		return seedableRand.Float64()
	}
	return float64(rand.Int31()) / (1 << 31)
}

// GetParticleEmitter2DLayerAndZ returns layer and z for 2D emitter (for flush sorting).
func GetParticleEmitter2DLayerAndZ(emitterID string) (layerID string, zIndex int) {
	particleEmitter2DMu.Lock()
	defer particleEmitter2DMu.Unlock()
	e := particleEmitters2D[emitterID]
	if e == nil {
		return "", 0
	}
	return e.LayerID, e.ZIndex
}
