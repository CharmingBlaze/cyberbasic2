// Package terrain provides heightmap and terrain mesh generation for CyberBasic.
package terrain

import (
	"fmt"
	"math"
	"sync"

	"cyberbasic/compiler/bindings/raylib"
)

// Heightmap is a 2D grid of heights (width × depth).
type Heightmap struct {
	Width  int
	Depth  int
	Heights []float32
}

var (
	heightmaps   = make(map[string]*Heightmap)
	heightmapSeq int
	heightmapMu  sync.Mutex
)

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

// LoadHeightmapFromImage fills a new heightmap from raylib image (grayscale 0-1). Returns heightmap id.
func LoadHeightmapFromImage(imageID string) (string, error) {
	width, height, heights, ok := raylib.GetImageDataForHeightmap(imageID)
	if !ok {
		return "", fmt.Errorf("unknown or invalid image id: %s", imageID)
	}
	if width <= 0 || height <= 0 {
		return "", fmt.Errorf("image has invalid dimensions")
	}
	hm := &Heightmap{Width: width, Depth: height, Heights: heights}
	heightmapMu.Lock()
	heightmapSeq++
	id := fmt.Sprintf("heightmap_%d", heightmapSeq)
	heightmaps[id] = hm
	heightmapMu.Unlock()
	return id, nil
}

// GenHeightmap creates a heightmap using 2D value noise. Returns heightmap id.
func GenHeightmap(width, depth int, noiseScale float64) (string, error) {
	if width <= 0 {
		width = 32
	}
	if depth <= 0 {
		depth = 32
	}
	if noiseScale <= 0 {
		noiseScale = 1
	}
	heights := make([]float32, width*depth)
	for z := 0; z < depth; z++ {
		for x := 0; x < width; x++ {
			n := valueNoise2D(float64(x)*noiseScale, float64(z)*noiseScale)
			heights[z*width+x] = float32(n)
		}
	}
	hm := &Heightmap{Width: width, Depth: depth, Heights: heights}
	heightmapMu.Lock()
	heightmapSeq++
	id := fmt.Sprintf("heightmap_%d", heightmapSeq)
	heightmaps[id] = hm
	heightmapMu.Unlock()
	return id, nil
}

// GetHeightmap returns the heightmap by id (caller must not modify).
func GetHeightmap(id string) *Heightmap {
	heightmapMu.Lock()
	hm := heightmaps[id]
	heightmapMu.Unlock()
	return hm
}

// SampleHeight performs bilinear sampling of the heightmap at normalized [0,1] x [0,1].
func (h *Heightmap) SampleHeight(nx, nz float32) float32 {
	if h == nil || len(h.Heights) == 0 {
		return 0
	}
	w, d := float32(h.Width-1), float32(h.Depth-1)
	if w <= 0 || d <= 0 {
		return h.Heights[0]
	}
	x := nx * w
	z := nz * d
	ix := int(x)
	iz := int(z)
	if ix < 0 {
		ix = 0
	}
	if iz < 0 {
		iz = 0
	}
	if ix >= h.Width-1 {
		ix = h.Width - 2
	}
	if iz >= h.Depth-1 {
		iz = h.Depth - 2
	}
	fx := x - float32(ix)
	fz := z - float32(iz)
	fx = fx * fx * (3 - 2*fx)
	fz = fz * fz * (3 - 2*fz)
	i00 := iz*h.Width + ix
	i10 := iz*h.Width + (ix + 1)
	i01 := (iz+1)*h.Width + ix
	i11 := (iz+1)*h.Width + (ix + 1)
	h00 := h.Heights[i00]
	h10 := h.Heights[i10]
	h01 := h.Heights[i01]
	h11 := h.Heights[i11]
	return h00*(1-fx)*(1-fz) + h10*fx*(1-fz) + h01*(1-fx)*fz + h11*fx*fz
}
