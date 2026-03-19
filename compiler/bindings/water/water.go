// Package water provides a water plane mesh and draw/params API for CyberBasic.
package water

import (
	"fmt"
	"math"
	"sync"

	"cyberbasic/compiler/bindings/modfacade"
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
	// UV scroll for animated water
	UScroll float32
	VScroll float32
	// Reflection/refraction (silently disabled if unsupported)
	ReflectionOn bool
	RefractionOn bool
	// Optional texture refs (for future shader use)
	ReflectionTexture string
	RefractionTexture string
	NormalMap         string
	FoamTexture      string
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
	// Depth-based color blending (deep vs shallow)
	DepthColorR   float32
	DepthColorG   float32
	DepthColorB   float32
	ShallowColorR float32
	ShallowColorG float32
	ShallowColorB float32
	// Visibility for Hide/Show
	Visible bool
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
		if width > 10000 {
			width = 10000
		}
		if depth > 10000 {
			depth = 10000
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
			MeshID:       meshID,
			MaterialID:   matID,
			Width:        width,
			Depth:        depth,
			WaveSpeed:    1,
			WaveHeight:   0.2,
			WaveFreq:     0.5,
			ColorA:       0.8,
			Transparency: 0.8,
			Visible:      true,
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
		if !w.Visible {
			return nil, nil
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
		bodyId := toString(args[0])
		waterId := toString(args[1])
		worldId := "default"
		if len(args) >= 3 {
			worldId = toString(args[2])
		}
		px, _ := v.CallForeign("GetPositionX3D", []interface{}{worldId, bodyId})
		py, _ := v.CallForeign("GetPositionY3D", []interface{}{worldId, bodyId})
		pz, _ := v.CallForeign("GetPositionZ3D", []interface{}{worldId, bodyId})
		bx := toFloat32(px)
		by := toFloat32(py)
		bz := toFloat32(pz)
		whRes, _ := v.CallForeign("WaterGetHeight", []interface{}{waterId, bx, bz})
		waterY := toFloat32(whRes)
		if by >= waterY {
			return nil, nil
		}
		waterMu.Lock()
		w, ok := waters[waterId]
		waterMu.Unlock()
		if !ok {
			return nil, nil
		}
		density := w.Density
		if density <= 0 {
			density = 1
		}
		submerged := waterY - by
		if submerged > 2 {
			submerged = 2
		}
		buoyancy := float64(density) * 9.81 * float64(submerged) * 0.5
		_, _ = v.CallForeign("ApplyForce3D", []interface{}{worldId, bodyId, 0.0, buoyancy, 0.0})
		return nil, nil
	})

	// SetWaterScroll(waterId, uSpeed, vSpeed): UV scroll for animated water.
	v.RegisterForeign("SetWaterScroll", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("SetWaterScroll requires (waterId, uSpeed, vSpeed)")
		}
		waterMu.Lock()
		defer waterMu.Unlock()
		if w, ok := waters[toString(args[0])]; ok {
			w.UScroll = toFloat32(args[1])
			w.VScroll = toFloat32(args[2])
		}
		return nil, nil
	})
	// SetWaterReflection(waterId, onOff): Enable reflection; silently disabled if unsupported.
	v.RegisterForeign("SetWaterReflection", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("SetWaterReflection requires (waterId, onOff)")
		}
		waterMu.Lock()
		defer waterMu.Unlock()
		if w, ok := waters[toString(args[0])]; ok {
			w.ReflectionOn = toFloat32(args[1]) != 0
		}
		return nil, nil
	})
	// SetWaterRefraction(waterId, onOff): Enable refraction; silently disabled if unsupported.
	v.RegisterForeign("SetWaterRefraction", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("SetWaterRefraction requires (waterId, onOff)")
		}
		waterMu.Lock()
		defer waterMu.Unlock()
		if w, ok := waters[toString(args[0])]; ok {
			w.RefractionOn = toFloat32(args[1]) != 0
		}
		return nil, nil
	})
	// SetWaterDepthColor(waterId, r, g, b): Deep water tint.
	v.RegisterForeign("SetWaterDepthColor", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("SetWaterDepthColor requires (waterId, r, g, b)")
		}
		waterMu.Lock()
		defer waterMu.Unlock()
		if w, ok := waters[toString(args[0])]; ok {
			w.DepthColorR = toFloat32(args[1])
			w.DepthColorG = toFloat32(args[2])
			w.DepthColorB = toFloat32(args[3])
		}
		return nil, nil
	})
	// SetWaterShallowColor(waterId, r, g, b): Shallow water tint.
	v.RegisterForeign("SetWaterShallowColor", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("SetWaterShallowColor requires (waterId, r, g, b)")
		}
		waterMu.Lock()
		defer waterMu.Unlock()
		if w, ok := waters[toString(args[0])]; ok {
			w.ShallowColorR = toFloat32(args[1])
			w.ShallowColorG = toFloat32(args[2])
			w.ShallowColorB = toFloat32(args[3])
		}
		return nil, nil
	})
	// SetWaterFoamTexture(waterId, textureId): Foam mask; skip layer if missing.
	v.RegisterForeign("SetWaterFoamTexture", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("SetWaterFoamTexture requires (waterId, textureId)")
		}
		waterMu.Lock()
		defer waterMu.Unlock()
		if w, ok := waters[toString(args[0])]; ok {
			w.FoamTexture = toString(args[1])
		}
		return nil, nil
	})
	// SetWaterVisible(waterId, onOff): Hide/Show water.
	v.RegisterForeign("SetWaterVisible", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("SetWaterVisible requires (waterId, onOff)")
		}
		waterMu.Lock()
		defer waterMu.Unlock()
		if w, ok := waters[toString(args[0])]; ok {
			w.Visible = toFloat32(args[1]) != 0
		}
		return nil, nil
	})

	v.SetGlobal("water", modfacade.New(v, waterV2))
}

// GetWaterByID returns WaterState for internal string id (for DBP id mapping).
func GetWaterByID(internalID string) *WaterState {
	waterMu.Lock()
	w := waters[internalID]
	waterMu.Unlock()
	return w
}

// WaterDelete removes water from the registry and unloads its mesh.
func WaterDelete(v *vm.VM, waterID string) error {
	waterMu.Lock()
	w, ok := waters[waterID]
	if !ok {
		waterMu.Unlock()
		return nil // already gone
	}
	meshID := w.MeshID
	delete(waters, waterID)
	waterMu.Unlock()
	if meshID != "" {
		_, _ = v.CallForeign("UnloadMesh", []interface{}{meshID})
	}
	return nil
}

// WaterClone creates a new water with the same state as the source. Returns new internal ID.
func WaterClone(v *vm.VM, srcID string) (string, error) {
	waterMu.Lock()
	src, ok := waters[srcID]
	if !ok {
		waterMu.Unlock()
		return "", fmt.Errorf("unknown water id: %s", srcID)
	}
	tileSize := src.Width / 16
	if tileSize < 1 {
		tileSize = 1
	}
	waterMu.Unlock()
	res, err := v.CallForeign("WaterCreate", []interface{}{src.Width, src.Depth, tileSize})
	if err != nil {
		return "", err
	}
	newID, ok := res.(string)
	if !ok || newID == "" {
		return "", fmt.Errorf("WaterCreate did not return water id")
	}
	waterMu.Lock()
	dst, ok := waters[newID]
	if !ok {
		waterMu.Unlock()
		return "", fmt.Errorf("new water not found")
	}
	dst.PosX, dst.PosY, dst.PosZ = src.PosX, src.PosY, src.PosZ
	dst.WaveSpeed, dst.WaveHeight, dst.WaveFreq = src.WaveSpeed, src.WaveHeight, src.WaveFreq
	dst.Time = src.Time
	dst.UScroll, dst.VScroll = src.UScroll, src.VScroll
	dst.ReflectionOn, dst.RefractionOn = src.ReflectionOn, src.RefractionOn
	dst.ReflectionTexture, dst.RefractionTexture = src.ReflectionTexture, src.RefractionTexture
	dst.NormalMap, dst.FoamTexture = src.NormalMap, src.FoamTexture
	dst.ColorR, dst.ColorG, dst.ColorB, dst.ColorA = src.ColorR, src.ColorG, src.ColorB, src.ColorA
	dst.Shininess, dst.FoamEnabled, dst.FoamIntensity = src.Shininess, src.FoamEnabled, src.FoamIntensity
	dst.DepthFade, dst.Transparency, dst.Density = src.DepthFade, src.Transparency, src.Density
	dst.DragLinear, dst.DragAngular = src.DragLinear, src.DragAngular
	dst.DepthColorR, dst.DepthColorG, dst.DepthColorB = src.DepthColorR, src.DepthColorG, src.DepthColorB
	dst.ShallowColorR, dst.ShallowColorG, dst.ShallowColorB = src.ShallowColorR, src.ShallowColorG, src.ShallowColorB
	dst.Visible = src.Visible
	waterMu.Unlock()
	return newID, nil
}
