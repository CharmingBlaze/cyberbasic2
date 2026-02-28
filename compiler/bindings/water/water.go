// Package water provides a water plane mesh and draw/params API for CyberBasic.
package water

import (
	"fmt"
	"math"
	"sync"

	"cyberbasic/compiler/vm"
)

// WaterState holds mesh, material, and wave params for a water plane.
type WaterState struct {
	MeshID     string
	MaterialID string
	Width      float32
	Depth      float32
	PosX       float32
	PosY       float32
	PosZ       float32
	WaveSpeed  float32
	WaveHeight float32
	WaveFreq   float32
	Time       float32
	// Optional texture refs (for future shader use)
	ReflectionTexture string
	RefractionTexture string
	NormalMap         string
	ColorR            float32
	ColorG            float32
	ColorB            float32
	ColorA            float32
	Shininess         float32
	FoamEnabled       bool
	FoamIntensity     float32
	DepthFade         float32
	Transparency      float32
	Density           float32
	DragLinear        float32
	DragAngular       float32
}

var (
	waters   = make(map[string]*WaterState)
	waterSeq int
	waterMu  sync.Mutex
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

// RegisterWater registers water creation, drawing, and param bindings with the VM.
func RegisterWater(v *vm.VM) {
	v.RegisterForeign("WaterCreate", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("WaterCreate requires (width, depth, tileSize)")
		}
		width := toFloat32(args[0])
		depth := toFloat32(args[1])
		tileSize := toFloat32(args[2])
		if width <= 0 {
			width = 10
		}
		if depth <= 0 {
			depth = 10
		}
		resX, resZ := int32(16), int32(16)
		if tileSize > 0 {
			resX = int32(width/tileSize) + 1
			resZ = int32(depth/tileSize) + 1
			if resX < 2 {
				resX = 2
			}
			if resZ < 2 {
				resZ = 2
			}
		}
		meshRes, err := v.CallForeign("GenMeshPlane", []interface{}{width, depth, resX, resZ})
		if err != nil {
			return nil, err
		}
		meshID, ok := meshRes.(string)
		if !ok || meshID == "" {
			return nil, fmt.Errorf("GenMeshPlane did not return mesh id")
		}
		matRes, err := v.CallForeign("LoadMaterialDefault", nil)
		if err != nil {
			return nil, err
		}
		matID := ""
		if matRes != nil {
			matID, _ = matRes.(string)
		}
		waterMu.Lock()
		waterSeq++
		id := fmt.Sprintf("water_%d", waterSeq)
		waters[id] = &WaterState{
			MeshID:     meshID,
			MaterialID: matID,
			Width:      width,
			Depth:      depth,
			WaveSpeed:  1,
			WaveHeight: 0.2,
			WaveFreq:   0.5,
			ColorA:     0.8,
			Transparency: 0.8,
		}
		waterMu.Unlock()
		return id, nil
	})

	v.RegisterForeign("DrawWater", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("DrawWater requires (waterId)")
		}
		id := toString(args[0])
		waterMu.Lock()
		w, ok := waters[id]
		waterMu.Unlock()
		if !ok {
			return nil, fmt.Errorf("unknown water id: %s", id)
		}
		posX, posY, posZ := w.PosX, w.PosY, w.PosZ
		if len(args) >= 4 {
			posX = toFloat32(args[1])
			posY = toFloat32(args[2])
			posZ = toFloat32(args[3])
		}
		matID := w.MaterialID
		if matID == "" {
			res, _ := v.CallForeign("LoadMaterialDefault", nil)
			if res != nil {
				matID, _ = res.(string)
			}
		}
		_, err := v.CallForeign("DrawMesh", []interface{}{
			w.MeshID, matID,
			posX, posY, posZ,
			float32(1), float32(1), float32(1),
		})
		return nil, err
	})

	v.RegisterRenderType("drawwater", vm.Render3D)

	v.RegisterForeign("SetWaterPosition", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("SetWaterPosition requires (waterId, x, y, z)")
		}
		id := toString(args[0])
		waterMu.Lock()
		defer waterMu.Unlock()
		w, ok := waters[id]
		if !ok {
			return nil, fmt.Errorf("unknown water id: %s", id)
		}
		w.PosX = toFloat32(args[1])
		w.PosY = toFloat32(args[2])
		w.PosZ = toFloat32(args[3])
		return nil, nil
	})

	v.RegisterForeign("SetWaterWaveSpeed", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("SetWaterWaveSpeed requires (waterId, speed)")
		}
		waterMu.Lock()
		defer waterMu.Unlock()
		if w, ok := waters[toString(args[0])]; ok {
			w.WaveSpeed = toFloat32(args[1])
		}
		return nil, nil
	})
	v.RegisterForeign("SetWaterWaveHeight", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("SetWaterWaveHeight requires (waterId, height)")
		}
		waterMu.Lock()
		defer waterMu.Unlock()
		if w, ok := waters[toString(args[0])]; ok {
			w.WaveHeight = toFloat32(args[1])
		}
		return nil, nil
	})
	v.RegisterForeign("SetWaterWaveFrequency", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("SetWaterWaveFrequency requires (waterId, frequency)")
		}
		waterMu.Lock()
		defer waterMu.Unlock()
		if w, ok := waters[toString(args[0])]; ok {
			w.WaveFreq = toFloat32(args[1])
		}
		return nil, nil
	})
	v.RegisterForeign("SetWaterTime", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("SetWaterTime requires (waterId, time)")
		}
		waterMu.Lock()
		defer waterMu.Unlock()
		if w, ok := waters[toString(args[0])]; ok {
			w.Time = toFloat32(args[1])
		}
		return nil, nil
	})

	v.RegisterForeign("WaterGetHeight", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("WaterGetHeight requires (waterId, x, z)")
		}
		id := toString(args[0])
		x, z := toFloat32(args[1]), toFloat32(args[2])
		waterMu.Lock()
		w, ok := waters[id]
		waterMu.Unlock()
		if !ok {
			return float64(0), nil
		}
		// Simple wave formula: base (PosY) + wave
		base := float64(w.PosY)
		wave := float64(w.WaveHeight) * math.Sin(float64(x)*float64(w.WaveFreq)+float64(w.Time)*float64(w.WaveSpeed)) * math.Cos(float64(z)*float64(w.WaveFreq)+float64(w.Time)*float64(w.WaveSpeed))
		return base + wave, nil
	})

	v.RegisterForeign("SetWaterTexture", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("SetWaterTexture requires (waterId, textureId)")
		}
		waterMu.Lock()
		defer waterMu.Unlock()
		if w, ok := waters[toString(args[0])]; ok {
			w.RefractionTexture = toString(args[1])
		}
		return nil, nil
	})
	v.RegisterForeign("SetWaterReflectionTexture", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("SetWaterReflectionTexture requires (waterId, textureId)")
		}
		waterMu.Lock()
		defer waterMu.Unlock()
		if w, ok := waters[toString(args[0])]; ok {
			w.ReflectionTexture = toString(args[1])
		}
		return nil, nil
	})
	v.RegisterForeign("SetWaterRefractionTexture", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("SetWaterRefractionTexture requires (waterId, textureId)")
		}
		waterMu.Lock()
		defer waterMu.Unlock()
		if w, ok := waters[toString(args[0])]; ok {
			w.RefractionTexture = toString(args[1])
		}
		return nil, nil
	})
	v.RegisterForeign("SetWaterNormalMap", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("SetWaterNormalMap requires (waterId, textureId)")
		}
		waterMu.Lock()
		defer waterMu.Unlock()
		if w, ok := waters[toString(args[0])]; ok {
			w.NormalMap = toString(args[1])
		}
		return nil, nil
	})
	v.RegisterForeign("SetWaterColor", func(args []interface{}) (interface{}, error) {
		if len(args) < 5 {
			return nil, fmt.Errorf("SetWaterColor requires (waterId, r, g, b, a)")
		}
		waterMu.Lock()
		defer waterMu.Unlock()
		if w, ok := waters[toString(args[0])]; ok {
			w.ColorR = toFloat32(args[1])
			w.ColorG = toFloat32(args[2])
			w.ColorB = toFloat32(args[3])
			w.ColorA = toFloat32(args[4])
		}
		return nil, nil
	})
	v.RegisterForeign("SetWaterShininess", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("SetWaterShininess requires (waterId, shininess)")
		}
		waterMu.Lock()
		defer waterMu.Unlock()
		if w, ok := waters[toString(args[0])]; ok {
			w.Shininess = toFloat32(args[1])
		}
		return nil, nil
	})
	v.RegisterForeign("WaterEnableFoam", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("WaterEnableFoam requires (waterId, enabled)")
		}
		waterMu.Lock()
		defer waterMu.Unlock()
		if w, ok := waters[toString(args[0])]; ok {
			w.FoamEnabled = toFloat32(args[1]) != 0
		}
		return nil, nil
	})
	v.RegisterForeign("WaterSetFoamIntensity", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("WaterSetFoamIntensity requires (waterId, intensity)")
		}
		waterMu.Lock()
		defer waterMu.Unlock()
		if w, ok := waters[toString(args[0])]; ok {
			w.FoamIntensity = toFloat32(args[1])
		}
		return nil, nil
	})
	v.RegisterForeign("WaterSetDepthFade", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("WaterSetDepthFade requires (waterId, fade)")
		}
		waterMu.Lock()
		defer waterMu.Unlock()
		if w, ok := waters[toString(args[0])]; ok {
			w.DepthFade = toFloat32(args[1])
		}
		return nil, nil
	})
	v.RegisterForeign("WaterSetTransparency", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("WaterSetTransparency requires (waterId, alpha)")
		}
		waterMu.Lock()
		defer waterMu.Unlock()
		if w, ok := waters[toString(args[0])]; ok {
			w.Transparency = toFloat32(args[1])
		}
		return nil, nil
	})
	v.RegisterForeign("WaterSetDensity", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("WaterSetDensity requires (waterId, value)")
		}
		waterMu.Lock()
		defer waterMu.Unlock()
		if w, ok := waters[toString(args[0])]; ok {
			w.Density = toFloat32(args[1])
		}
		return nil, nil
	})
	v.RegisterForeign("WaterSetDrag", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("WaterSetDrag requires (waterId, linear, angular)")
		}
		waterMu.Lock()
		defer waterMu.Unlock()
		if w, ok := waters[toString(args[0])]; ok {
			w.DragLinear = toFloat32(args[1])
			w.DragAngular = toFloat32(args[2])
		}
		return nil, nil
	})
	v.RegisterForeign("WaterApplyBuoyancy", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("WaterApplyBuoyancy requires (bodyId, waterId)")
		}
		// Stub: sample water height at body position and apply upward force via physics; full impl would use Bullet body.
		return nil, nil
	})
}
