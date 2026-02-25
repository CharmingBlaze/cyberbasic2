// Package game provides high-level bindings: particles, AI, animation, coroutine stubs.
package game

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"sync"
	"time"

	"cyberbasic/compiler/vm"
	rl "github.com/gen2brain/raylib-go/raylib"
)

func toString(v interface{}) string {
	if v == nil {
		return ""
	}
	return fmt.Sprint(v)
}

func toFloat64(v interface{}) float64 {
	switch x := v.(type) {
	case int:
		return float64(x)
	case float64:
		return x
	case float32:
		return float64(x)
	default:
		return 0
	}
}

// --- Particle system ---
type particle struct {
	X, Y, Z    float32
	VX, VY, VZ float32
	R, G, B, A uint8
	Life       float32
	MaxLife    float32
}

type particleSystem struct {
	Particles  []particle
	DefaultR   uint8
	DefaultG   uint8
	DefaultB   uint8
	DefaultA   uint8
	Lifetime   float32
	VelX, VelY, VelZ float32
}

var (
	particleSystems   = make(map[string]*particleSystem)
	particleSeq       int
	particleMu        sync.Mutex

	// AI state: entity id -> position, velocity, speed, target
	aiPos     = make(map[string][3]float64)
	aiVel     = make(map[string][3]float64)
	aiSpeed   = make(map[string]float64)
	aiTarget  = make(map[string][3]float64)
	aiWander  = make(map[string]float64) // radius
	aiMu      sync.RWMutex

	// Animation: key -> { startValue, targetValue, startTime, duration }
	animValue    = make(map[string]struct{ Start, Target float64; StartT time.Time; Dur time.Duration })
	animColor    = make(map[string]struct{ R0, G0, B0, A0, R1, G1, B1, A1 float64; StartT time.Time; Dur time.Duration })
	animPosition = make(map[string]struct{ X0, Y0, Z0, X1, Y1, Z1 float64; StartT time.Time; Dur time.Duration })
	animRotation = make(map[string]struct{ P0, Y0, R0, P1, Y1, R1 float64; StartT time.Time; Dur time.Duration })
	animMu       sync.Mutex

	// Tilemap: id -> grid and solid set
	tilemaps   = make(map[string]*tilemapData)
	tilemapSeq int
	tilemapMu  sync.RWMutex

	// Dialogue
	dialogueNodes    = make(map[string]map[string]interface{})
	dialogueVars     = make(map[string]interface{})
	dialogueCurrent  string
	dialogueMu       sync.RWMutex

	// Inventory: invId -> slots []{itemID, amount}, maxSlots
	inventories   = make(map[string]*invData)
	invSeq        int
	invMu         sync.RWMutex
	itemDefs      = make(map[string]*itemDef)
	itemDefMu     sync.RWMutex

	// Behavior trees
	aiTrees   = make(map[string]*btNode)
	aiTreeSeq int
	aiTreeMu  sync.Mutex

	// Replication (state to sync)
	replicateVars   = make(map[string]map[string]bool) // entity -> set of var names
	replicatePos    = make(map[string]bool)
	replicateRot    = make(map[string]bool)
	replicateMu     sync.Mutex

	// Shader graph (nodes + connections; compile = stub)
	shaderGraphNodes   = make(map[string]*sgNode)
	shaderGraphGraphs  = make(map[string]*sgGraph)
	shaderGraphSeq     int
	shaderGraphMu      sync.Mutex

	// Anim state machine
	animStates     = make(map[string]*animStateData)
	animTransitions = make(map[string][]*animTransition)
	animParams     = make(map[string]float64)   // param name -> value
	animEntityState = make(map[string]string)   // entityId -> current state name
	animStateSeq    int
	animStateMu     sync.Mutex
)

type invData struct {
	Slots    []struct{ ItemID string; Amount int }
	MaxSlots int
}

type itemDef struct {
	Name      string
	Icon      string
	StackSize int
	Props     map[string]interface{}
}

type btNode struct {
	Type     string
	Children []string
	FuncName string
}

type sgNode struct {
	ID   string
	Type string
	Args []interface{}
}

type sgGraph struct {
	Nodes []string
	Conns []struct{ Out, In string }
}

type animStateData struct {
	Name string
	Clip string
}

type animTransition struct {
	From, To string
	Condition string
}

type tilemapData struct {
	Tiles    [][]int
	TileSize int
	Solid    map[int]bool
}

func ensureTilemapSize(tm *tilemapData, w, h int) {
	if w <= 0 {
		w = 1
	}
	if h <= 0 {
		h = 1
	}
	for len(tm.Tiles) < h {
		tm.Tiles = append(tm.Tiles, make([]int, w))
	}
	for y := range tm.Tiles {
		for len(tm.Tiles[y]) < w {
			tm.Tiles[y] = append(tm.Tiles[y], 0)
		}
	}
}

// valueNoise2D returns deterministic noise in [0,1] for procedural gen.
func valueNoise2D(x, y float64) float64 {
	ix, iy := int(math.Floor(x))&0xff, int(math.Floor(y))&0xff
	fx := x - math.Floor(x)
	fy := y - math.Floor(y)
	fx = fx * fx * (3 - 2*fx)
	fy = fy * fy * (3 - 2*fy)
	h := func(a, b int) float64 {
		n := (a*313 + b*757) % 1024
		if n < 0 {
			n += 1024
		}
		return float64(n) / 1024
	}
	v00 := h(ix, iy)
	v10 := h(ix+1, iy)
	v01 := h(ix, iy+1)
	v11 := h(ix+1, iy+1)
	return v00*(1-fx)*(1-fy) + v10*fx*(1-fy) + v01*(1-fx)*fy + v11*fx*fy
}

func valueNoise3D(x, y, z float64) float64 {
	ix, iy, iz := int(math.Floor(x))&0x1f, int(math.Floor(y))&0x1f, int(math.Floor(z))&0x1f
	fx := x - math.Floor(x)
	fy := y - math.Floor(y)
	fz := z - math.Floor(z)
	fx = fx * fx * (3 - 2*fx)
	fy = fy * fy * (3 - 2*fy)
	fz = fz * fz * (3 - 2*fz)
	h := func(a, b, c int) float64 {
		n := (a*313 + b*757 + c*419) % 1024
		if n < 0 {
			n += 1024
		}
		return float64(n) / 1024
	}
	v000 := h(ix, iy, iz)
	v100 := h(ix+1, iy, iz)
	v010 := h(ix, iy+1, iz)
	v110 := h(ix+1, iy+1, iz)
	v001 := h(ix, iy, iz+1)
	v101 := h(ix+1, iy, iz+1)
	v011 := h(ix, iy+1, iz+1)
	v111 := h(ix+1, iy+1, iz+1)
	return v000*(1-fx)*(1-fy)*(1-fz) + v100*fx*(1-fy)*(1-fz) + v010*(1-fx)*fy*(1-fz) + v110*fx*fy*(1-fz) +
		v001*(1-fx)*(1-fy)*fz + v101*fx*(1-fy)*fz + v011*(1-fx)*fy*fz + v111*fx*fy*fz
}

// RegisterGame registers particle, AI, animation, and coroutine bindings.
func RegisterGame(v *vm.VM) {
	// --- Particles ---
	v.RegisterForeign("CreateParticleSystem", func(args []interface{}) (interface{}, error) {
		particleMu.Lock()
		particleSeq++
		id := fmt.Sprintf("ps_%d", particleSeq)
		particleSystems[id] = &particleSystem{
			Particles: make([]particle, 0, 256),
			DefaultR:  255, DefaultG: 255, DefaultB: 255, DefaultA: 255,
			Lifetime: 1.0, VelX: 0, VelY: 0, VelZ: 0,
		}
		particleMu.Unlock()
		return id, nil
	})
	v.RegisterForeign("EmitParticles", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("EmitParticles requires (systemId, count)")
		}
		id := toString(args[0])
		count := int(toFloat64(args[1]))
		if count <= 0 || count > 1000 {
			count = 10
		}
		particleMu.Lock()
		ps := particleSystems[id]
		if ps == nil {
			particleMu.Unlock()
			return nil, nil
		}
		for i := 0; i < count; i++ {
			ps.Particles = append(ps.Particles, particle{
				X: 0, Y: 0, Z: 0,
				VX: ps.VelX, VY: ps.VelY, VZ: ps.VelZ,
				R: ps.DefaultR, G: ps.DefaultG, B: ps.DefaultB, A: ps.DefaultA,
				Life: ps.Lifetime, MaxLife: ps.Lifetime,
			})
		}
		particleMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("SetParticleColor", func(args []interface{}) (interface{}, error) {
		if len(args) < 5 {
			return nil, fmt.Errorf("SetParticleColor requires (systemId, r, g, b, a)")
		}
		id := toString(args[0])
		particleMu.Lock()
		ps := particleSystems[id]
		if ps != nil {
			ps.DefaultR = uint8(toFloat64(args[1])) & 0xff
			ps.DefaultG = uint8(toFloat64(args[2])) & 0xff
			ps.DefaultB = uint8(toFloat64(args[3])) & 0xff
			ps.DefaultA = uint8(toFloat64(args[4])) & 0xff
		}
		particleMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("SetParticleLifetime", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("SetParticleLifetime requires (systemId, seconds)")
		}
		id := toString(args[0])
		particleMu.Lock()
		ps := particleSystems[id]
		if ps != nil {
			ps.Lifetime = float32(toFloat64(args[1]))
		}
		particleMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("SetParticleVelocity", func(args []interface{}) (interface{}, error) {
		if len(args) < 5 {
			return nil, fmt.Errorf("SetParticleVelocity requires (systemId, vx, vy, vz)")
		}
		id := toString(args[0])
		particleMu.Lock()
		ps := particleSystems[id]
		if ps != nil {
			ps.VelX = float32(toFloat64(args[1]))
			ps.VelY = float32(toFloat64(args[2]))
			ps.VelZ = float32(toFloat64(args[3]))
		}
		particleMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("DrawParticles", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("DrawParticles requires (systemId)")
		}
		id := toString(args[0])
		dt := rl.GetFrameTime()
		particleMu.Lock()
		ps := particleSystems[id]
		if ps == nil {
			particleMu.Unlock()
			return nil, nil
		}
		live := ps.Particles[:0]
		for _, p := range ps.Particles {
			p.X += p.VX * dt
			p.Y += p.VY * dt
			p.Z += p.VZ * dt
			p.Life -= dt
			if p.Life > 0 {
				live = append(live, p)
			}
		}
		ps.Particles = live
		particleMu.Unlock()
		for _, p := range live {
			c := rl.NewColor(p.R, p.G, p.B, uint8(float32(p.A)*p.Life/p.MaxLife))
			rl.DrawSphere(rl.Vector3{X: p.X, Y: p.Y, Z: p.Z}, 0.1, c)
		}
		return nil, nil
	})

	// --- AI ---
	v.RegisterForeign("AISetPosition", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("AISetPosition requires (entityId, x, y, z)")
		}
		eid := toString(args[0])
		aiMu.Lock()
		aiPos[eid] = [3]float64{toFloat64(args[1]), toFloat64(args[2]), toFloat64(args[3])}
		if _, ok := aiSpeed[eid]; !ok {
			aiSpeed[eid] = 1
		}
		aiMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("GetAIPosition", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return []interface{}{0.0, 0.0, 0.0}, nil
		}
		eid := toString(args[0])
		aiMu.RLock()
		p := aiPos[eid]
		aiMu.RUnlock()
		return []interface{}{p[0], p[1], p[2]}, nil
	})
	v.RegisterForeign("AIUpdate", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, nil
		}
		eid := toString(args[0])
		dt := rl.GetFrameTime()
		aiMu.Lock()
		pos, ok := aiPos[eid]
		if !ok {
			aiPos[eid] = [3]float64{0, 0, 0}
			aiSpeed[eid] = 1
			aiMu.Unlock()
			return nil, nil
		}
		targ := aiTarget[eid]
		sp := aiSpeed[eid]
		if sp <= 0 {
			sp = 1
		}
		dx := targ[0] - pos[0]
		dy := targ[1] - pos[1]
		dz := targ[2] - pos[2]
		dist := math.Sqrt(dx*dx + dy*dy + dz*dz)
		if dist > 1e-6 {
			move := sp * float64(dt)
			if move > dist {
				move = dist
			}
			pos[0] += dx / dist * move
			pos[1] += dy / dist * move
			pos[2] += dz / dist * move
		}
		aiPos[eid] = pos
		aiMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("AIMoveTo", func(args []interface{}) (interface{}, error) {
		if len(args) < 5 {
			return nil, fmt.Errorf("AIMoveTo requires (entityId, x, y, z)")
		}
		eid := toString(args[0])
		tx, ty, tz := toFloat64(args[1]), toFloat64(args[2]), toFloat64(args[3])
		aiMu.Lock()
		aiTarget[eid] = [3]float64{tx, ty, tz}
		if _, ok := aiPos[eid]; !ok {
			aiPos[eid] = [3]float64{0, 0, 0}
			aiSpeed[eid] = 1
		}
		aiMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("AISetSpeed", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("AISetSpeed requires (entityId, speed)")
		}
		eid := toString(args[0])
		aiMu.Lock()
		aiSpeed[eid] = toFloat64(args[1])
		aiMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("AIWander", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("AIWander requires (entityId, radius)")
		}
		eid := toString(args[0])
		aiMu.Lock()
		aiWander[eid] = toFloat64(args[1])
		aiMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("AIChase", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("AIChase requires (entityId, targetEntityId)")
		}
		eid := toString(args[0])
		tid := toString(args[1])
		aiMu.Lock()
		if tpos, ok := aiPos[tid]; ok {
			aiTarget[eid] = tpos
		}
		aiMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("AIFlee", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("AIFlee requires (entityId, targetEntityId)")
		}
		eid := toString(args[0])
		tid := toString(args[1])
		aiMu.Lock()
		pos := aiPos[eid]
		tpos := aiTarget[tid]
		if _, ok := aiPos[tid]; ok {
			tpos = aiPos[tid]
		}
		dx, dy, dz := pos[0]-tpos[0], pos[1]-tpos[1], pos[2]-tpos[2]
		n := math.Sqrt(dx*dx + dy*dy + dz*dz)
		if n > 1e-6 {
			dx, dy, dz = dx/n, dy/n, dz/n
			sp := aiSpeed[eid]
			if sp <= 0 {
				sp = 1
			}
			aiTarget[eid] = [3]float64{pos[0] + dx*sp, pos[1] + dy*sp, pos[2] + dz*sp}
		}
		aiMu.Unlock()
		return nil, nil
	})

	// --- Animation (lerp over time; script polls current value each frame) ---
	v.RegisterForeign("AnimateValue", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("AnimateValue requires (key, target, durationSeconds)")
		}
		key := toString(args[0])
		target := toFloat64(args[1])
		dur := time.Duration(toFloat64(args[2]) * 1e9)
		animMu.Lock()
		start := 0.0
		if a, ok := animValue[key]; ok {
			elapsed := time.Since(a.StartT).Seconds()
			if elapsed >= a.Dur.Seconds() {
				start = a.Target
			} else {
				t := elapsed / a.Dur.Seconds()
				start = a.Start + (a.Target-a.Start)*t
			}
		}
		animValue[key] = struct{ Start, Target float64; StartT time.Time; Dur time.Duration }{start, target, time.Now(), dur}
		animMu.Unlock()
		return start, nil
	})
	v.RegisterForeign("AnimateColor", func(args []interface{}) (interface{}, error) {
		if len(args) < 6 {
			return nil, fmt.Errorf("AnimateColor requires (key, r, g, b, a, durationSeconds)")
		}
		key := toString(args[0])
		r, g, b, a := toFloat64(args[1]), toFloat64(args[2]), toFloat64(args[3]), toFloat64(args[4])
		dur := time.Duration(toFloat64(args[5]) * 1e9)
		animMu.Lock()
		r0, g0, b0, a0 := r, g, b, a
		if prev, ok := animColor[key]; ok {
			elapsed := time.Since(prev.StartT).Seconds()
			t := elapsed / prev.Dur.Seconds()
			if t < 1 {
				r0 = prev.R0 + (prev.R1-prev.R0)*t
				g0 = prev.G0 + (prev.G1-prev.G0)*t
				b0 = prev.B0 + (prev.B1-prev.B0)*t
				a0 = prev.A0 + (prev.A1-prev.A0)*t
			}
		}
		animColor[key] = struct{ R0, G0, B0, A0, R1, G1, B1, A1 float64; StartT time.Time; Dur time.Duration }{r0, g0, b0, a0, r, g, b, a, time.Now(), dur}
		animMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("AnimatePosition", func(args []interface{}) (interface{}, error) {
		if len(args) < 5 {
			return nil, fmt.Errorf("AnimatePosition requires (key, x, y, z, durationSeconds)")
		}
		key := toString(args[0])
		x, y, z := toFloat64(args[1]), toFloat64(args[2]), toFloat64(args[3])
		dur := time.Duration(toFloat64(args[4]) * 1e9)
		animMu.Lock()
		x0, y0, z0 := x, y, z
		if prev, ok := animPosition[key]; ok {
			elapsed := time.Since(prev.StartT).Seconds()
			t := elapsed / prev.Dur.Seconds()
			if t < 1 {
				x0 = prev.X0 + (prev.X1-prev.X0)*t
				y0 = prev.Y0 + (prev.Y1-prev.Y0)*t
				z0 = prev.Z0 + (prev.Z1-prev.Z0)*t
			}
		}
		animPosition[key] = struct{ X0, Y0, Z0, X1, Y1, Z1 float64; StartT time.Time; Dur time.Duration }{x0, y0, z0, x, y, z, time.Now(), dur}
		animMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("AnimateRotation", func(args []interface{}) (interface{}, error) {
		if len(args) < 5 {
			return nil, fmt.Errorf("AnimateRotation requires (key, pitch, yaw, roll, durationSeconds)")
		}
		key := toString(args[0])
		p, y, r := toFloat64(args[1]), toFloat64(args[2]), toFloat64(args[3])
		dur := time.Duration(toFloat64(args[4]) * 1e9)
		animMu.Lock()
		p0, y0, r0 := p, y, r
		if prev, ok := animRotation[key]; ok {
			elapsed := time.Since(prev.StartT).Seconds()
			t := elapsed / prev.Dur.Seconds()
			if t < 1 {
				p0 = prev.P0 + (prev.P1-prev.P0)*t
				y0 = prev.Y0 + (prev.Y1-prev.Y0)*t
				r0 = prev.R0 + (prev.R1-prev.R0)*t
			}
		}
		animRotation[key] = struct{ P0, Y0, R0, P1, Y1, R1 float64; StartT time.Time; Dur time.Duration }{p0, y0, r0, p, y, r, time.Now(), dur}
		animMu.Unlock()
		return nil, nil
	})

	// --- Coroutine stubs (VM cannot yield; no-op) ---
	v.RegisterForeign("CoroutineStart", func(args []interface{}) (interface{}, error) { return "", nil })
	v.RegisterForeign("CoroutineYield", func(args []interface{}) (interface{}, error) { return nil, nil })
	v.RegisterForeign("CoroutineWait", func(args []interface{}) (interface{}, error) { return nil, nil })
	v.RegisterForeign("CoroutineStop", func(args []interface{}) (interface{}, error) { return nil, nil })

	// --- Tilemap (2D grid of tile IDs) ---
	v.RegisterForeign("LoadTilemap", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("LoadTilemap requires (path)")
		}
		path := toString(args[0])
		_ = path
		tilemapMu.Lock()
		tilemapSeq++
		id := fmt.Sprintf("tm_%d", tilemapSeq)
		tilemaps[id] = &tilemapData{Tiles: [][]int{}, TileSize: 32, Solid: make(map[int]bool)}
		tilemapMu.Unlock()
		return id, nil
	})
	v.RegisterForeign("DrawTilemap", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("DrawTilemap requires (mapId)")
		}
		id := toString(args[0])
		tilemapMu.RLock()
		tm := tilemaps[id]
		tilemapMu.RUnlock()
		if tm == nil || len(tm.Tiles) == 0 {
			return nil, nil
		}
		for y := 0; y < len(tm.Tiles); y++ {
			for x := 0; x < len(tm.Tiles[y]); x++ {
				tid := tm.Tiles[y][x]
				if tid != 0 {
					px := float32(x * tm.TileSize)
					py := float32(y * tm.TileSize)
					rl.DrawRectangle(int32(px), int32(py), int32(tm.TileSize), int32(tm.TileSize), rl.Gray)
				}
			}
		}
		return nil, nil
	})
	v.RegisterForeign("SetTile", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("SetTile requires (mapId, x, y, tileID)")
		}
		id := toString(args[0])
		x, y := int(toFloat64(args[1])), int(toFloat64(args[2]))
		tid := int(toFloat64(args[3]))
		tilemapMu.Lock()
		tm := tilemaps[id]
		if tm != nil {
			ensureTilemapSize(tm, x+1, y+1)
			tm.Tiles[y][x] = tid
		}
		tilemapMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("GetTile", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("GetTile requires (mapId, x, y)")
		}
		id := toString(args[0])
		x, y := int(toFloat64(args[1])), int(toFloat64(args[2]))
		tilemapMu.RLock()
		tm := tilemaps[id]
		tilemapMu.RUnlock()
		if tm == nil || y < 0 || y >= len(tm.Tiles) || x < 0 || x >= len(tm.Tiles[y]) {
			return 0, nil
		}
		return tm.Tiles[y][x], nil
	})
	v.RegisterForeign("TilemapCollision", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("TilemapCollision requires (mapId, x, y)")
		}
		id := toString(args[0])
		wx, wy := toFloat64(args[1]), toFloat64(args[2])
		tilemapMu.RLock()
		tm := tilemaps[id]
		tilemapMu.RUnlock()
		if tm == nil || tm.TileSize <= 0 {
			return false, nil
		}
		tx := int(wx) / tm.TileSize
		ty := int(wy) / tm.TileSize
		if ty < 0 || ty >= len(tm.Tiles) || tx < 0 || tx >= len(tm.Tiles[ty]) {
			return false, nil
		}
		tid := tm.Tiles[ty][tx]
		return tm.Solid[tid], nil
	})

	// --- Pathfinding (simple grid A* stub; returns empty path) ---
	v.RegisterForeign("PathfindGrid", func(args []interface{}) (interface{}, error) {
		if len(args) < 5 {
			return nil, fmt.Errorf("PathfindGrid requires (mapId, startX, startY, endX, endY)")
		}
		return []interface{}{}, nil
	})
	v.RegisterForeign("PathfindNavmesh", func(args []interface{}) (interface{}, error) {
		return []interface{}{}, nil
	})
	v.RegisterForeign("FollowPath", func(args []interface{}) (interface{}, error) {
		return nil, nil
	})

	// --- Scripting / event stubs (VM cannot pass function refs) ---
	v.RegisterForeign("OnKeyPress", func(args []interface{}) (interface{}, error) { return nil, nil })
	v.RegisterForeign("OnMouseClick", func(args []interface{}) (interface{}, error) { return nil, nil })
	v.RegisterForeign("OnUpdate", func(args []interface{}) (interface{}, error) { return nil, nil })
	v.RegisterForeign("OnDraw", func(args []interface{}) (interface{}, error) { return nil, nil })
	v.RegisterForeign("OnCollision", func(args []interface{}) (interface{}, error) { return nil, nil })

	// --- Debug ---
	v.RegisterForeign("DebugDrawGrid", func(args []interface{}) (interface{}, error) {
		slices := int32(10)
		spacing := float32(1)
		if len(args) >= 1 {
			slices = int32(toFloat64(args[0]))
		}
		if len(args) >= 2 {
			spacing = float32(toFloat64(args[1]))
		}
		rl.DrawGrid(slices, spacing)
		return nil, nil
	})
	v.RegisterForeign("DebugDrawBounds", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, nil
		}
		_ = toString(args[0])
		return nil, nil
	})
	v.RegisterForeign("DebugLog", func(args []interface{}) (interface{}, error) {
		if len(args) >= 1 {
			fmt.Fprintln(os.Stderr, "[debug]", fmt.Sprint(args[0]))
		}
		return nil, nil
	})
	v.RegisterForeign("DebugWatch", func(args []interface{}) (interface{}, error) {
		return nil, nil
	})

	// --- Procedural generation ---
	v.RegisterForeign("Noise2D", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return 0.0, nil
		}
		x, y := toFloat64(args[0]), toFloat64(args[1])
		return valueNoise2D(x, y), nil
	})
	v.RegisterForeign("Noise3D", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return 0.0, nil
		}
		x, y, z := toFloat64(args[0]), toFloat64(args[1]), toFloat64(args[2])
		return valueNoise3D(x, y, z), nil
	})
	v.RegisterForeign("GenerateDungeon", func(args []interface{}) (interface{}, error) {
		w, h := 32, 24
		if len(args) >= 1 {
			w = int(toFloat64(args[0]))
		}
		if len(args) >= 2 {
			h = int(toFloat64(args[1]))
		}
		if w <= 0 {
			w = 32
		}
		if h <= 0 {
			h = 24
		}
		tilemapMu.Lock()
		tilemapSeq++
		id := fmt.Sprintf("dungeon_%d", tilemapSeq)
		tm := &tilemapData{Tiles: make([][]int, h), TileSize: 32, Solid: map[int]bool{0: true}}
		for y := 0; y < h; y++ {
			tm.Tiles[y] = make([]int, w)
			for x := 0; x < w; x++ {
				tm.Tiles[y][x] = 0
			}
		}
		// Simple random rooms + corridors
		roomCount := (w * h) / 80
		if roomCount < 3 {
			roomCount = 3
		}
		seeded := time.Now().UnixNano()
		for i := 0; i < roomCount; i++ {
			rw := 3 + (int(seeded+int64(i*7)) % 5)
			rh := 3 + (int(seeded+int64(i*11)) % 4)
			rx := 1 + (int(seeded+int64(i*13)) % (w - rw - 1))
			ry := 1 + (int(seeded+int64(i*17)) % (h - rh - 1))
			for yy := ry; yy < ry+rh && yy < h; yy++ {
				for xx := rx; xx < rx+rw && xx < w; xx++ {
					tm.Tiles[yy][xx] = 1
				}
			}
		}
		tilemaps[id] = tm
		tilemapMu.Unlock()
		return id, nil
	})
	v.RegisterForeign("GenerateTree", func(args []interface{}) (interface{}, error) {
		seed := int64(0)
		if len(args) >= 1 {
			seed = int64(toFloat64(args[0]))
		}
		if seed == 0 {
			seed = time.Now().UnixNano()
		}
		// Return a deterministic tree id (script can use for drawing or placement)
		id := fmt.Sprintf("tree_%d", seed)
		return id, nil
	})
	v.RegisterForeign("GenerateCity", func(args []interface{}) (interface{}, error) {
		size := 16
		if len(args) >= 1 {
			size = int(toFloat64(args[0]))
		}
		if size <= 0 {
			size = 16
		}
		w, h := size*4, size*4
		tilemapMu.Lock()
		tilemapSeq++
		id := fmt.Sprintf("city_%d", tilemapSeq)
		tm := &tilemapData{Tiles: make([][]int, h), TileSize: 16, Solid: map[int]bool{1: true}}
		for y := 0; y < h; y++ {
			tm.Tiles[y] = make([]int, w)
			for x := 0; x < w; x++ {
				tm.Tiles[y][x] = 0
			}
		}
		// Grid of streets (0) and blocks; some cells = building (1)
		for y := 0; y < h; y++ {
			for x := 0; x < w; x++ {
				if x%4 == 0 || y%4 == 0 {
					tm.Tiles[y][x] = 0
				} else {
					if (x+y)%3 != 0 {
						tm.Tiles[y][x] = 1
					}
				}
			}
		}
		tilemaps[id] = tm
		tilemapMu.Unlock()
		return id, nil
	})

	// --- Dialogue system ---
	v.RegisterForeign("DialogueLoad", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("DialogueLoad requires (path)")
		}
		data, err := os.ReadFile(toString(args[0]))
		if err != nil {
			return nil, err
		}
		var nodes map[string]map[string]interface{}
		if err := json.Unmarshal(data, &nodes); err != nil {
			return nil, err
		}
		dialogueMu.Lock()
		for k, v := range nodes {
			dialogueNodes[k] = v
		}
		dialogueMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("DialogueStart", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("DialogueStart requires (id)")
		}
		dialogueMu.Lock()
		dialogueCurrent = toString(args[0])
		dialogueMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("DialogueNext", func(args []interface{}) (interface{}, error) {
		dialogueMu.Lock()
		cur := dialogueCurrent
		nodes := dialogueNodes
		dialogueMu.Unlock()
		if cur == "" {
			return nil, nil
		}
		node := nodes[cur]
		if node == nil {
			return nil, nil
		}
		if next, ok := node["next"].(string); ok && next != "" {
			dialogueMu.Lock()
			dialogueCurrent = next
			dialogueMu.Unlock()
		}
		return nil, nil
	})
	v.RegisterForeign("DialogueChoice", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, nil
		}
		idx := int(toFloat64(args[0]))
		dialogueMu.Lock()
		cur := dialogueCurrent
		node := dialogueNodes[cur]
		dialogueMu.Unlock()
		if node == nil {
			return nil, nil
		}
		if choices, ok := node["choices"].([]interface{}); ok && idx >= 0 && idx < len(choices) {
			if c, ok := choices[idx].(map[string]interface{}); ok {
				if next, ok := c["next"].(string); ok {
					dialogueMu.Lock()
					dialogueCurrent = next
					dialogueMu.Unlock()
				}
			}
		}
		return nil, nil
	})
	v.RegisterForeign("DialogueShowText", func(args []interface{}) (interface{}, error) { return nil, nil })
	v.RegisterForeign("DialogueShowChoices", func(args []interface{}) (interface{}, error) { return nil, nil })
	v.RegisterForeign("DialogueSetVar", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, nil
		}
		dialogueMu.Lock()
		dialogueVars[toString(args[0])] = args[1]
		dialogueMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("DialogueGetVar", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, nil
		}
		dialogueMu.RLock()
		v := dialogueVars[toString(args[0])]
		dialogueMu.RUnlock()
		return v, nil
	})

	// --- Inventory ---
	v.RegisterForeign("InventoryCreate", func(args []interface{}) (interface{}, error) {
		size := 20
		if len(args) >= 1 {
			size = int(toFloat64(args[0]))
		}
		if size <= 0 {
			size = 20
		}
		invMu.Lock()
		invSeq++
		id := fmt.Sprintf("inv_%d", invSeq)
		inventories[id] = &invData{Slots: make([]struct{ ItemID string; Amount int }, 0, size), MaxSlots: size}
		invMu.Unlock()
		return id, nil
	})
	v.RegisterForeign("InventoryAddItem", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("InventoryAddItem requires (invId, itemID, amount)")
		}
		invId := toString(args[0])
		itemID := toString(args[1])
		amount := int(toFloat64(args[2]))
		invMu.Lock()
		inv := inventories[invId]
		if inv != nil {
			for i := range inv.Slots {
				if inv.Slots[i].ItemID == itemID {
					inv.Slots[i].Amount += amount
					invMu.Unlock()
					return nil, nil
				}
			}
			if len(inv.Slots) < inv.MaxSlots {
				inv.Slots = append(inv.Slots, struct{ ItemID string; Amount int }{itemID, amount})
			}
		}
		invMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("InventoryRemoveItem", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, nil
		}
		invId := toString(args[0])
		itemID := toString(args[1])
		amount := int(toFloat64(args[2]))
		invMu.Lock()
		inv := inventories[invId]
		if inv != nil {
			for i := range inv.Slots {
				if inv.Slots[i].ItemID == itemID {
					inv.Slots[i].Amount -= amount
					if inv.Slots[i].Amount <= 0 {
						inv.Slots = append(inv.Slots[:i], inv.Slots[i+1:]...)
					}
					break
				}
			}
		}
		invMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("InventoryHasItem", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return false, nil
		}
		invId := toString(args[0])
		itemID := toString(args[1])
		invMu.RLock()
		inv := inventories[invId]
		has := false
		if inv != nil {
			for i := range inv.Slots {
				if inv.Slots[i].ItemID == itemID && inv.Slots[i].Amount > 0 {
					has = true
					break
				}
			}
		}
		invMu.RUnlock()
		return has, nil
	})
	v.RegisterForeign("ItemDefine", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("ItemDefine requires (id, name, icon, stackSize)")
		}
		id := toString(args[0])
		itemDefMu.Lock()
		itemDefs[id] = &itemDef{Name: toString(args[1]), Icon: toString(args[2]), StackSize: int(toFloat64(args[3])), Props: make(map[string]interface{})}
		itemDefMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("ItemSetProperty", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, nil
		}
		id := toString(args[0])
		itemDefMu.Lock()
		if d := itemDefs[id]; d != nil && d.Props != nil {
			d.Props[toString(args[1])] = args[2]
		}
		itemDefMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("InventoryDraw", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, nil
		}
		invId := toString(args[0])
		x, y := int32(toFloat64(args[1])), int32(toFloat64(args[2]))
		invMu.RLock()
		inv := inventories[invId]
		invMu.RUnlock()
		if inv == nil {
			return nil, nil
		}
		slotSize := int32(40)
		for i := range inv.Slots {
			px := x + int32(i%5)*slotSize
			py := y + int32(i/5)*slotSize
			rl.DrawRectangle(px, py, slotSize-2, slotSize-2, rl.DarkGray)
			rl.DrawRectangleLines(px, py, slotSize-2, slotSize-2, rl.White)
		}
		return nil, nil
	})

	// --- Physics joints (stubs; use BULLET.* for real joints) ---
	v.RegisterForeign("CreateHingeJoint", func(args []interface{}) (interface{}, error) { return "", nil })
	v.RegisterForeign("CreateBallJoint", func(args []interface{}) (interface{}, error) { return "", nil })
	v.RegisterForeign("CreateSliderJoint", func(args []interface{}) (interface{}, error) { return "", nil })
	v.RegisterForeign("CreateRagdoll", func(args []interface{}) (interface{}, error) { return "", nil })
	v.RegisterForeign("RagdollEnable", func(args []interface{}) (interface{}, error) { return nil, nil })
	v.RegisterForeign("RagdollDisable", func(args []interface{}) (interface{}, error) { return nil, nil })

	// --- AI behavior trees ---
	v.RegisterForeign("AIBehaviorTreeCreate", func(args []interface{}) (interface{}, error) {
		aiTreeMu.Lock()
		aiTreeSeq++
		id := fmt.Sprintf("bt_%d", aiTreeSeq)
		aiTrees[id] = &btNode{Type: "root", Children: nil}
		aiTreeMu.Unlock()
		return id, nil
	})
	v.RegisterForeign("AISelector", func(args []interface{}) (interface{}, error) {
		aiTreeMu.Lock()
		aiTreeSeq++
		id := fmt.Sprintf("bt_%d", aiTreeSeq)
		children := make([]string, 0)
		for i := 0; i < len(args); i++ {
			children = append(children, toString(args[i]))
		}
		aiTrees[id] = &btNode{Type: "selector", Children: children}
		aiTreeMu.Unlock()
		return id, nil
	})
	v.RegisterForeign("AISequence", func(args []interface{}) (interface{}, error) {
		aiTreeMu.Lock()
		aiTreeSeq++
		id := fmt.Sprintf("bt_%d", aiTreeSeq)
		children := make([]string, 0)
		for i := 0; i < len(args); i++ {
			children = append(children, toString(args[i]))
		}
		aiTrees[id] = &btNode{Type: "sequence", Children: children}
		aiTreeMu.Unlock()
		return id, nil
	})
	v.RegisterForeign("AIAction", func(args []interface{}) (interface{}, error) {
		aiTreeMu.Lock()
		aiTreeSeq++
		id := fmt.Sprintf("bt_%d", aiTreeSeq)
		fn := ""
		if len(args) >= 1 {
			fn = toString(args[0])
		}
		aiTrees[id] = &btNode{Type: "action", FuncName: fn}
		aiTreeMu.Unlock()
		return id, nil
	})
	v.RegisterForeign("AICondition", func(args []interface{}) (interface{}, error) {
		aiTreeMu.Lock()
		aiTreeSeq++
		id := fmt.Sprintf("bt_%d", aiTreeSeq)
		fn := ""
		if len(args) >= 1 {
			fn = toString(args[0])
		}
		aiTrees[id] = &btNode{Type: "condition", FuncName: fn}
		aiTreeMu.Unlock()
		return id, nil
	})
	v.RegisterForeign("AIRun", func(args []interface{}) (interface{}, error) { return nil, nil })

	// --- Multiplayer replication ---
	v.RegisterForeign("NetStartServer", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("NetStartServer requires (port)")
		}
		return nil, nil
	})
	v.RegisterForeign("NetStartClient", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("NetStartClient requires (ip, port)")
		}
		return nil, nil
	})
	v.RegisterForeign("ReplicateVariable", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, nil
		}
		eid := toString(args[0])
		replicateMu.Lock()
		if replicateVars[eid] == nil {
			replicateVars[eid] = make(map[string]bool)
		}
		replicateVars[eid][toString(args[1])] = true
		replicateMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("ReplicatePosition", func(args []interface{}) (interface{}, error) {
		if len(args) >= 1 {
			replicateMu.Lock()
			replicatePos[toString(args[0])] = true
			replicateMu.Unlock()
		}
		return nil, nil
	})
	v.RegisterForeign("ReplicateRotation", func(args []interface{}) (interface{}, error) {
		if len(args) >= 1 {
			replicateMu.Lock()
			replicateRot[toString(args[0])] = true
			replicateMu.Unlock()
		}
		return nil, nil
	})
	v.RegisterForeign("RPC", func(args []interface{}) (interface{}, error) { return nil, nil })

	// --- Shader graph (stub; compile returns empty) ---
	v.RegisterForeign("ShaderNodeTexture", func(args []interface{}) (interface{}, error) {
		shaderGraphMu.Lock()
		shaderGraphSeq++
		id := fmt.Sprintf("sg_%d", shaderGraphSeq)
		tex := ""
		if len(args) >= 1 {
			tex = toString(args[0])
		}
		shaderGraphNodes[id] = &sgNode{ID: id, Type: "texture", Args: []interface{}{tex}}
		shaderGraphMu.Unlock()
		return id, nil
	})
	v.RegisterForeign("ShaderNodeColor", func(args []interface{}) (interface{}, error) {
		shaderGraphMu.Lock()
		shaderGraphSeq++
		id := fmt.Sprintf("sg_%d", shaderGraphSeq)
		shaderGraphNodes[id] = &sgNode{ID: id, Type: "color", Args: args}
		shaderGraphMu.Unlock()
		return id, nil
	})
	v.RegisterForeign("ShaderNodeAdd", func(args []interface{}) (interface{}, error) {
		shaderGraphMu.Lock()
		shaderGraphSeq++
		id := fmt.Sprintf("sg_%d", shaderGraphSeq)
		shaderGraphNodes[id] = &sgNode{ID: id, Type: "add", Args: args}
		shaderGraphMu.Unlock()
		return id, nil
	})
	v.RegisterForeign("ShaderNodeMultiply", func(args []interface{}) (interface{}, error) {
		shaderGraphMu.Lock()
		shaderGraphSeq++
		id := fmt.Sprintf("sg_%d", shaderGraphSeq)
		shaderGraphNodes[id] = &sgNode{ID: id, Type: "multiply", Args: args}
		shaderGraphMu.Unlock()
		return id, nil
	})
	v.RegisterForeign("ShaderNodeTime", func(args []interface{}) (interface{}, error) {
		shaderGraphMu.Lock()
		shaderGraphSeq++
		id := fmt.Sprintf("sg_%d", shaderGraphSeq)
		shaderGraphNodes[id] = &sgNode{ID: id, Type: "time", Args: nil}
		shaderGraphMu.Unlock()
		return id, nil
	})
	v.RegisterForeign("ShaderGraphCreate", func(args []interface{}) (interface{}, error) {
		shaderGraphMu.Lock()
		shaderGraphSeq++
		id := fmt.Sprintf("graph_%d", shaderGraphSeq)
		shaderGraphGraphs[id] = &sgGraph{}
		shaderGraphMu.Unlock()
		return id, nil
	})
	v.RegisterForeign("ShaderGraphConnect", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("ShaderGraphConnect requires (graphId, outputNodeId, inputNodeId)")
		}
		gid := toString(args[0])
		shaderGraphMu.Lock()
		g := shaderGraphGraphs[gid]
		if g != nil {
			g.Conns = append(g.Conns, struct{ Out, In string }{toString(args[1]), toString(args[2])})
		}
		shaderGraphMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("ShaderGraphCompile", func(args []interface{}) (interface{}, error) {
		return "", nil
	})

	// --- Animation state machine ---
	v.RegisterForeign("AnimStateCreate", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("AnimStateCreate requires (name)")
		}
		animStateMu.Lock()
		animStateSeq++
		id := fmt.Sprintf("animstate_%d", animStateSeq)
		animStates[id] = &animStateData{Name: toString(args[0])}
		animStateMu.Unlock()
		return id, nil
	})
	v.RegisterForeign("AnimStateSetClip", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, nil
		}
		animStateMu.Lock()
		if s := animStates[toString(args[0])]; s != nil {
			s.Clip = toString(args[1])
		}
		animStateMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("AnimTransition", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, nil
		}
		from, to := toString(args[0]), toString(args[1])
		cond := toString(args[2])
		animStateMu.Lock()
		animTransitions[from] = append(animTransitions[from], &animTransition{From: from, To: to, Condition: cond})
		animStateMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("AnimSetParameter", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, nil
		}
		animStateMu.Lock()
		animParams[toString(args[0])] = toFloat64(args[1])
		animStateMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("AnimSetState", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, nil
		}
		animStateMu.Lock()
		animEntityState[toString(args[0])] = toString(args[1])
		animStateMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("AnimUpdate", func(args []interface{}) (interface{}, error) {
		return nil, nil
	})
}
