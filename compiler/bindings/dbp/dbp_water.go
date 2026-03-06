// Package dbp: DBP-style water commands with integer IDs.
//
// MakeWater(id, width, depth), SetWaterTexture(id, path), PositionWater(id, x, y, z),
// SetWaterScroll, SetWaterWaveStrength, SetWaterWaveSpeed, SetWaterReflection,
// SetWaterRefraction, SetWaterNormalmap, SetWaterFoamTexture, SetWaterDepthColor,
// SetWaterShallowColor, DrawWater.
// All texture loads are safe: missing files fall back to solid blue.
package dbp

import (
	"fmt"
	"sync"

	"cyberbasic/compiler/bindings/water"
	"cyberbasic/compiler/vm"
	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
	idToWater     = make(map[int]string)
	idToWaterMu   sync.Mutex
	waterTexs     = make(map[string]rl.Texture2D)
	waterTexSeq   int
	waterTexMu    sync.Mutex
)

// loadTextureSafe loads texture from path. Returns (textureID, true) on success.
// On failure (missing file, invalid), returns fallback - never panics.
// Fallback: creates 1x1 blue texture. Water package stores id for future shader use.
func loadTextureSafe(path string) (string, bool) {
	if path == "" {
		return createFallbackTextureBlue(), false
	}
	tex := rl.LoadTexture(path)
	if tex.ID == 0 {
		return createFallbackTextureBlue(), false
	}
	waterTexMu.Lock()
	waterTexSeq++
	id := fmt.Sprintf("water_tex_%d", waterTexSeq)
	waterTexs[id] = tex
	waterTexMu.Unlock()
	return id, true
}

// createFallbackTextureBlue creates 1x1 solid blue texture for water fallback.
func createFallbackTextureBlue() string {
	img := rl.GenImageColor(1, 1, rl.NewColor(0, 0, 255, 255))
	tex := rl.LoadTextureFromImage(img)
	rl.UnloadImage(img)
	if tex.ID == 0 {
		return ""
	}
	waterTexMu.Lock()
	waterTexSeq++
	id := fmt.Sprintf("water_tex_%d", waterTexSeq)
	waterTexs[id] = tex
	waterTexMu.Unlock()
	return id
}

func registerWater(v *vm.VM) {
	// MakeWater(id, width, depth): Create water plane; map integer id to internal string id.
	v.RegisterForeign("MakeWater", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("MakeWater(id, width, depth) requires 3 arguments")
		}
		id := toInt(args[0])
		width := toFloat32(args[1])
		depth := toFloat32(args[2])
		tileSize := float32(16)
		if width > 0 && depth > 0 {
			tileSize = width / 16
			if tileSize < 1 {
				tileSize = 1
			}
		}
		res, err := v.CallForeign("WaterCreate", []interface{}{width, depth, tileSize})
		if err != nil {
			return nil, err
		}
		internalID, ok := res.(string)
		if !ok || internalID == "" {
			return nil, fmt.Errorf("WaterCreate did not return water id")
		}
		idToWaterMu.Lock()
		idToWater[id] = internalID
		idToWaterMu.Unlock()
		return nil, nil
	})

	// SetWaterTexture(id, file$): Load texture from path; fallback to solid blue if missing.
	v.RegisterForeign("SetWaterTexture", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("SetWaterTexture(id, path) requires 2 arguments")
		}
		id := toInt(args[0])
		path := toString(args[1])
		idToWaterMu.Lock()
		internalID, ok := idToWater[id]
		idToWaterMu.Unlock()
		if !ok {
			return nil, nil // unknown water id - no-op, don't crash
		}
		texID, _ := loadTextureSafe(path)
		if texID != "" {
			_, _ = v.CallForeign("SetWaterTexture", []interface{}{internalID, texID})
		}
		return nil, nil
	})

	// PositionWater(id, x, y, z)
	v.RegisterForeign("PositionWater", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("PositionWater(id, x, y, z) requires 4 arguments")
		}
		id := toInt(args[0])
		idToWaterMu.Lock()
		internalID, ok := idToWater[id]
		idToWaterMu.Unlock()
		if !ok {
			return nil, nil
		}
		return v.CallForeign("SetWaterPosition", []interface{}{
			internalID, toFloat32(args[1]), toFloat32(args[2]), toFloat32(args[3]),
		})
	})

	// SetWaterLevel(id, height): Set Y position only.
	v.RegisterForeign("SetWaterLevel", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("SetWaterLevel(id, height) requires 2 arguments")
		}
		id := toInt(args[0])
		idToWaterMu.Lock()
		internalID, ok := idToWater[id]
		idToWaterMu.Unlock()
		if !ok {
			return nil, nil
		}
		// Get current position, update Y only
		w := water.GetWaterByID(internalID)
		if w == nil {
			return nil, nil
		}
		return v.CallForeign("SetWaterPosition", []interface{}{
			internalID, w.PosX, toFloat32(args[1]), w.PosZ,
		})
	})

	// SetWaterColor(id, r, g, b)
	v.RegisterForeign("SetWaterColor", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("SetWaterColor(id, r, g, b) requires 4 arguments")
		}
		id := toInt(args[0])
		idToWaterMu.Lock()
		internalID, ok := idToWater[id]
		idToWaterMu.Unlock()
		if !ok {
			return nil, nil
		}
		r, g, b := toFloat32(args[1]), toFloat32(args[2]), toFloat32(args[3])
		if r > 1 || g > 1 || b > 1 {
			r, g, b = r/255, g/255, b/255
		}
		return v.CallForeign("SetWaterColor", []interface{}{internalID, r, g, b, float32(0.8)})
	})

	// SetWaterScroll(id, uSpeed, vSpeed)
	v.RegisterForeign("SetWaterScroll", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("SetWaterScroll(id, uSpeed, vSpeed) requires 3 arguments")
		}
		id := toInt(args[0])
		idToWaterMu.Lock()
		internalID, ok := idToWater[id]
		idToWaterMu.Unlock()
		if !ok {
			return nil, nil
		}
		return v.CallForeign("SetWaterScroll", []interface{}{
			internalID, toFloat32(args[1]), toFloat32(args[2]),
		})
	})

	// SetWaterWaveStrength(id, value): Alias for SetWaterWaveHeight.
	v.RegisterForeign("SetWaterWaveStrength", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("SetWaterWaveStrength(id, value) requires 2 arguments")
		}
		id := toInt(args[0])
		idToWaterMu.Lock()
		internalID, ok := idToWater[id]
		idToWaterMu.Unlock()
		if !ok {
			return nil, nil
		}
		return v.CallForeign("SetWaterWaveHeight", []interface{}{internalID, toFloat32(args[1])})
	})

	// SetWaterWaveSpeed(id, value)
	v.RegisterForeign("SetWaterWaveSpeed", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("SetWaterWaveSpeed(id, value) requires 2 arguments")
		}
		id := toInt(args[0])
		idToWaterMu.Lock()
		internalID, ok := idToWater[id]
		idToWaterMu.Unlock()
		if !ok {
			return nil, nil
		}
		return v.CallForeign("SetWaterWaveSpeed", []interface{}{internalID, toFloat32(args[1])})
	})

	// SetWaterReflection(id, onOff)
	v.RegisterForeign("SetWaterReflection", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("SetWaterReflection(id, onOff) requires 2 arguments")
		}
		id := toInt(args[0])
		idToWaterMu.Lock()
		internalID, ok := idToWater[id]
		idToWaterMu.Unlock()
		if !ok {
			return nil, nil
		}
		return v.CallForeign("SetWaterReflection", []interface{}{internalID, args[1]})
	})

	// SetWaterRefraction(id, onOff)
	v.RegisterForeign("SetWaterRefraction", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("SetWaterRefraction(id, onOff) requires 2 arguments")
		}
		id := toInt(args[0])
		idToWaterMu.Lock()
		internalID, ok := idToWater[id]
		idToWaterMu.Unlock()
		if !ok {
			return nil, nil
		}
		return v.CallForeign("SetWaterRefraction", []interface{}{internalID, args[1]})
	})

	// SetWaterNormalmap(id, file$): Load normal map; skip layer if missing.
	v.RegisterForeign("SetWaterNormalmap", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("SetWaterNormalmap(id, path) requires 2 arguments")
		}
		id := toInt(args[0])
		path := toString(args[1])
		idToWaterMu.Lock()
		internalID, ok := idToWater[id]
		idToWaterMu.Unlock()
		if !ok {
			return nil, nil
		}
		texID, _ := loadTextureSafe(path)
		if texID != "" {
			_, _ = v.CallForeign("SetWaterNormalMap", []interface{}{internalID, texID})
		}
		return nil, nil
	})

	// SetWaterFoamTexture(id, file$): Load foam texture; skip if missing.
	v.RegisterForeign("SetWaterFoamTexture", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("SetWaterFoamTexture(id, path) requires 2 arguments")
		}
		id := toInt(args[0])
		path := toString(args[1])
		idToWaterMu.Lock()
		internalID, ok := idToWater[id]
		idToWaterMu.Unlock()
		if !ok {
			return nil, nil
		}
		texID, _ := loadTextureSafe(path)
		if texID != "" {
			_, _ = v.CallForeign("SetWaterFoamTexture", []interface{}{internalID, texID})
			_, _ = v.CallForeign("WaterEnableFoam", []interface{}{internalID, 1})
		}
		return nil, nil
	})

	// SetWaterDepthColor(id, r, g, b)
	v.RegisterForeign("SetWaterDepthColor", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("SetWaterDepthColor(id, r, g, b) requires 4 arguments")
		}
		id := toInt(args[0])
		idToWaterMu.Lock()
		internalID, ok := idToWater[id]
		idToWaterMu.Unlock()
		if !ok {
			return nil, nil
		}
		return v.CallForeign("SetWaterDepthColor", []interface{}{
			internalID, toFloat32(args[1]), toFloat32(args[2]), toFloat32(args[3]),
		})
	})

	// SetWaterShallowColor(id, r, g, b)
	v.RegisterForeign("SetWaterShallowColor", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("SetWaterShallowColor(id, r, g, b) requires 4 arguments")
		}
		id := toInt(args[0])
		idToWaterMu.Lock()
		internalID, ok := idToWater[id]
		idToWaterMu.Unlock()
		if !ok {
			return nil, nil
		}
		return v.CallForeign("SetWaterShallowColor", []interface{}{
			internalID, toFloat32(args[1]), toFloat32(args[2]), toFloat32(args[3]),
		})
	})

	// DrawWater(id): Draw at stored position.
	v.RegisterForeign("DrawWater", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("DrawWater(id) requires 1 argument")
		}
		id := toInt(args[0])
		idToWaterMu.Lock()
		internalID, ok := idToWater[id]
		idToWaterMu.Unlock()
		if !ok {
			return nil, nil // unknown id - no-op
		}
		w := water.GetWaterByID(internalID)
		if w == nil {
			return nil, nil
		}
		return v.CallForeign("DrawWater", []interface{}{
			internalID, w.PosX, w.PosY, w.PosZ,
		})
	})

	v.RegisterRenderType("drawwater", vm.Render3D)
}

// RegisterWater registers DBP-style water commands (MakeWater, SetWaterTexture, etc.) with integer IDs.
// Call after water.RegisterWater so DBP commands overwrite for the public API.
func RegisterWater(v *vm.VM) {
	registerWater(v)
}
