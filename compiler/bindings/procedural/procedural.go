// Package procedural provides noise and scatter helpers for world generation.
package procedural

import (
	"fmt"
	"math"
	"math/rand"
	"sync"

	"cyberbasic/compiler/vm"
)

func toFloat32(v interface{}) float32 {
	switch x := v.(type) {
	case int:
		return float32(x)
	case int32:
		return float32(x)
	case float64:
		return float32(x)
	case float32:
		return x
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

func toInt32(v interface{}) int32 {
	switch x := v.(type) {
	case int:
		return int32(x)
	case int32:
		return x
	case float64:
		return int32(x)
	default:
		return 0
	}
}

// valueNoise2D returns deterministic noise in [0,1].
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

var procMu sync.Mutex

// RegisterProcedural registers noise and scatter with the VM.
func RegisterProcedural(v *vm.VM) {
	v.RegisterForeign("NoisePerlin2D", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("NoisePerlin2D requires (x, y, scale)")
		}
		x := toFloat32(args[0]) * toFloat32(args[2])
		y := toFloat32(args[1]) * toFloat32(args[2])
		return valueNoise2D(float64(x), float64(y)), nil
	})

	v.RegisterForeign("NoiseFractal2D", func(args []interface{}) (interface{}, error) {
		if len(args) < 5 {
			return nil, fmt.Errorf("NoiseFractal2D requires (x, y, octaves, persistence, lacunarity)")
		}
		x, y := float64(toFloat32(args[0])), float64(toFloat32(args[1]))
		octaves := int(toInt32(args[2]))
		persistence := toFloat32(args[3])
		lacunarity := toFloat32(args[4])
		if octaves <= 0 {
			octaves = 4
		}
		var sum, amp, freq float64 = 0, 1, 1
		maxVal := 0.0
		for i := 0; i < octaves; i++ {
			sum += amp * valueNoise2D(x*freq, y*freq)
			maxVal += amp
			amp *= float64(persistence)
			freq *= float64(lacunarity)
		}
		return sum / maxVal, nil
	})

	v.RegisterForeign("NoiseSimplex2D", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("NoiseSimplex2D requires (x, y, scale)")
		}
		x := toFloat32(args[0]) * toFloat32(args[2])
		y := toFloat32(args[1]) * toFloat32(args[2])
		return valueNoise2D(float64(x), float64(y)), nil
	})

	v.RegisterForeign("ScatterTrees", func(args []interface{}) (interface{}, error) {
		if len(args) < 5 {
			return nil, fmt.Errorf("ScatterTrees requires (treeSystemId, treeTypeId, areaX, areaZ, density)")
		}
		sysID := toString(args[0])
		typeID := toString(args[1])
		areaX, areaZ := toFloat32(args[2]), toFloat32(args[3])
		density := toFloat32(args[4])
		n := int(density)
		if n <= 0 {
			n = 20
		}
		procMu.Lock()
		for i := 0; i < n; i++ {
			x := (float32(rand.Float64()) - 0.5) * 2 * areaX
			z := (float32(rand.Float64()) - 0.5) * 2 * areaZ
			scale := 0.8 + float32(rand.Float64())*0.4
			rot := float32(rand.Float64() * 2 * math.Pi)
			_, _ = v.CallForeign("TreePlace", []interface{}{sysID, typeID, x, 0, z, scale, rot})
		}
		procMu.Unlock()
		return nil, nil
	})

	v.RegisterForeign("ScatterGrass", func(args []interface{}) (interface{}, error) {
		if len(args) < 5 {
			return nil, fmt.Errorf("ScatterGrass requires (grassId, centerX, centerZ, radius, density)")
		}
		grassID := toString(args[0])
		cx, cz := toFloat32(args[1]), toFloat32(args[2])
		radius, density := toFloat32(args[3]), toFloat32(args[4])
		_, _ = v.CallForeign("GrassPaint", []interface{}{grassID, cx, cz, radius, density})
		return nil, nil
	})

	v.RegisterForeign("ScatterObjects", func(args []interface{}) (interface{}, error) {
		if len(args) < 5 {
			return nil, fmt.Errorf("ScatterObjects requires (modelId, areaX, areaZ, count)")
		}
		modelID := toString(args[0])
		areaX, areaZ := toFloat32(args[1]), toFloat32(args[2])
		count := int(toInt32(args[3]))
		if count <= 0 {
			count = 10
		}
		minS, maxS := float32(0.8), float32(1.2)
		if len(args) >= 6 {
			minS, maxS = toFloat32(args[4]), toFloat32(args[5])
		}
		_, err := v.CallForeign("ObjectRandomScatter", []interface{}{modelID, areaX, areaZ, count, minS, maxS})
		return nil, err
	})
}
