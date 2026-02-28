package vegetation

import (
	"fmt"
	"math"
	"math/rand"
	"sync"
)

// GrassInstance is one grass blade/patch position.
type GrassInstance struct {
	X, Y, Z   float32
	Scale     float32
	Rotation  float32
}

// GrassState holds texture and instances for a grass system.
type GrassState struct {
	TextureID  string
	Density    float32
	PatchSize  float32
	Instances  []GrassInstance
	WindSpeed  float32
	WindStrength float32
	Height     float32
	ColorR     float32
	ColorG     float32
	ColorB     float32
	ColorA     float32
	LODDist       float32
	Instancing    bool
	BendAmount   float32
	Interaction  bool
}

var (
	grassSystems   = make(map[string]*GrassState)
	grassSeq       int
	grassMu        sync.Mutex
)

// GrassCreate creates a grass system. Returns grass id.
func GrassCreate(textureID string, density, patchSize float32) string {
	grassMu.Lock()
	grassSeq++
	id := fmt.Sprintf("grass_%d", grassSeq)
	grassSystems[id] = &GrassState{
		TextureID:      textureID,
		Density:        density,
		PatchSize:      patchSize,
		Instances:      nil,
		WindSpeed:      1,
		WindStrength:   0.1,
		Height:         1,
		ColorR:         1,
		ColorG:         1,
		ColorB:         1,
		ColorA:         1,
	}
	grassMu.Unlock()
	return id
}

func getGrass(id string) *GrassState {
	grassMu.Lock()
	g := grassSystems[id]
	grassMu.Unlock()
	return g
}

// GrassSetWind sets wind speed and strength.
func GrassSetWind(grassID string, speed, strength float32) {
	if g := getGrass(grassID); g != nil {
		g.WindSpeed = speed
		g.WindStrength = strength
	}
}

// GrassSetHeight sets default grass blade height.
func GrassSetHeight(grassID string, height float32) {
	if g := getGrass(grassID); g != nil {
		g.Height = height
	}
}

// GrassSetColor sets grass color (r,g,b,a 0-1).
func GrassSetColor(grassID string, r, grn, b, a float32) {
	if g := getGrass(grassID); g != nil {
		g.ColorR, g.ColorG, g.ColorB, g.ColorA = r, grn, b, a
	}
}

// GrassPaint adds grass instances in a disk at (x,z) with radius and density (count).
func GrassPaint(grassID string, x, z, radius, density float32) {
	g := getGrass(grassID)
	if g == nil {
		return
	}
	n := int(density)
	if n <= 0 {
		n = 10
	}
	grassMu.Lock()
	for i := 0; i < n; i++ {
		angle := float32(rand.Float64() * 2 * math.Pi)
		r := float32(rand.Float64()) * radius
		gx := x + r*float32(math.Cos(float64(angle)))
		gz := z + r*float32(math.Sin(float64(angle)))
		g.Instances = append(g.Instances, GrassInstance{X: gx, Y: 0, Z: gz, Scale: g.Height, Rotation: float32(i) * 0.1})
	}
	grassMu.Unlock()
}

// GrassErase removes instances within radius of (x,z).
func GrassErase(grassID string, x, z, radius float32) {
	g := getGrass(grassID)
	if g == nil {
		return
	}
	radiusSq := radius * radius
	grassMu.Lock()
	filtered := g.Instances[:0]
	for _, inst := range g.Instances {
		dx := inst.X - x
		dz := inst.Z - z
		if dx*dx+dz*dz > radiusSq {
			filtered = append(filtered, inst)
		}
	}
	g.Instances = filtered
	grassMu.Unlock()
}

// GrassSetDensity sets the default density for new paints.
func GrassSetDensity(grassID string, density float32) {
	if g := getGrass(grassID); g != nil {
		g.Density = density
	}
}

// GrassSetLOD sets LOD distance for the grass system.
func GrassSetLOD(grassID string, dist float32) {
	if g := getGrass(grassID); g != nil {
		g.LODDist = dist
	}
}

// GrassEnableInstancing enables or disables instancing.
func GrassEnableInstancing(grassID string, on bool) {
	if g := getGrass(grassID); g != nil {
		g.Instancing = on
	}
}

// GrassSetBendAmount sets bend amount for wind (shader/displacement).
func GrassSetBendAmount(grassID string, value float32) {
	if g := getGrass(grassID); g != nil {
		g.BendAmount = value
	}
}

// GrassSetInteraction enables/disables interaction (e.g. player displacement).
func GrassSetInteraction(grassID string, flag bool) {
	if g := getGrass(grassID); g != nil {
		g.Interaction = flag
	}
}