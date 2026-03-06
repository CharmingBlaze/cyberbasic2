// Package dbp: DBP-style terrain commands with integer IDs.
//
// MakeTerrain(id, width, depth), LoadHeightmap(id, path, width, depth, heightScale),
// SetTerrainTexture(id, path), PositionTerrain(id, x, y, z), SetTerrainLayer,
// SetTerrainSplatmap, GenerateTerrainNoise, DrawTerrain.
// All loads are safe: missing texture -> green fallback, missing heightmap -> flat.
package dbp

import (
	"fmt"
	"sync"

	"cyberbasic/compiler/bindings/terrain"
	"cyberbasic/compiler/vm"
	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
	idToTerrain   = make(map[int]string)
	idToTerrainMu sync.Mutex
	terrainTexs   = make(map[string]rl.Texture2D)
	terrainTexSeq int
	terrainTexMu  sync.Mutex
)

// loadTerrainTextureSafe loads texture from path for terrain. Fallback: solid green.
func loadTerrainTextureSafe(path string) string {
	if path == "" {
		return createFallbackTextureGreen()
	}
	tex := rl.LoadTexture(path)
	if tex.ID == 0 {
		return createFallbackTextureGreen()
	}
	terrainTexMu.Lock()
	terrainTexSeq++
	id := fmt.Sprintf("terrain_tex_%d", terrainTexSeq)
	terrainTexs[id] = tex
	terrainTexMu.Unlock()
	return id
}

func createFallbackTextureGreen() string {
	img := rl.GenImageColor(1, 1, rl.NewColor(0, 128, 0, 255))
	tex := rl.LoadTextureFromImage(img)
	rl.UnloadImage(img)
	if tex.ID == 0 {
		return ""
	}
	terrainTexMu.Lock()
	terrainTexSeq++
	id := fmt.Sprintf("terrain_tex_%d", terrainTexSeq)
	terrainTexs[id] = tex
	terrainTexMu.Unlock()
	return id
}

func registerTerrain(v *vm.VM) {
	// MakeTerrain(id, width, depth): Create flat terrain; map integer id to internal string id.
	v.RegisterForeign("MakeTerrain", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("MakeTerrain(id, width, depth) requires 3 arguments")
		}
		id := toInt(args[0])
		width := toFloat32(args[1])
		depth := toFloat32(args[2])
		res, err := terrain.MakeTerrainFlat(v, width, depth)
		if err != nil {
			return nil, err
		}
		idToTerrainMu.Lock()
		idToTerrain[id] = res
		idToTerrainMu.Unlock()
		return nil, nil
	})

	// LoadHeightmap(id, file$, width, depth, heightScale): Load from file; fallback flat if missing.
	v.RegisterForeign("LoadHeightmap", func(args []interface{}) (interface{}, error) {
		if len(args) < 5 {
			return nil, fmt.Errorf("LoadHeightmap(id, path, width, depth, heightScale) requires 5 arguments")
		}
		id := toInt(args[0])
		path := toString(args[1])
		width := int(toInt(args[2]))
		depth := int(toInt(args[3]))
		heightScale := toFloat32(args[4])
		hmID, err := terrain.LoadHeightmapFromFile(path, width, depth, heightScale)
		if err != nil {
			return nil, err
		}
		// Create terrain from heightmap
		res, err := v.CallForeign("TerrainCreate", []interface{}{hmID, float32(width), float32(depth), heightScale})
		if err != nil {
			return nil, err
		}
		internalID, ok := res.(string)
		if !ok || internalID == "" {
			return nil, fmt.Errorf("TerrainCreate did not return terrain id")
		}
		idToTerrainMu.Lock()
		idToTerrain[id] = internalID
		idToTerrainMu.Unlock()
		return nil, nil
	})

	// SetTerrainTexture(id, file$): Load texture from path; fallback green if missing.
	v.RegisterForeign("SetTerrainTexture", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("SetTerrainTexture(id, path) requires 2 arguments")
		}
		id := toInt(args[0])
		path := toString(args[1])
		idToTerrainMu.Lock()
		internalID, ok := idToTerrain[id]
		idToTerrainMu.Unlock()
		if !ok {
			return nil, nil
		}
		texID := loadTerrainTextureSafe(path)
		if texID != "" {
			_ = terrain.SetTerrainTexture(internalID, texID)
		}
		return nil, nil
	})

	// PositionTerrain(id, x, y, z)
	v.RegisterForeign("PositionTerrain", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("PositionTerrain(id, x, y, z) requires 4 arguments")
		}
		id := toInt(args[0])
		idToTerrainMu.Lock()
		internalID, ok := idToTerrain[id]
		idToTerrainMu.Unlock()
		if !ok {
			return nil, nil
		}
		return nil, terrain.SetTerrainPosition(internalID, toFloat32(args[1]), toFloat32(args[2]), toFloat32(args[3]))
	})

	// SetTerrainLayer(id, layerIndex, texture$)
	v.RegisterForeign("SetTerrainLayer", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("SetTerrainLayer(id, layerIndex, path) requires 3 arguments")
		}
		id := toInt(args[0])
		layerIndex := toInt(args[1])
		path := toString(args[2])
		idToTerrainMu.Lock()
		internalID, ok := idToTerrain[id]
		idToTerrainMu.Unlock()
		if !ok {
			return nil, nil
		}
		texID := loadTerrainTextureSafe(path)
		if texID != "" {
			_ = terrain.SetTerrainLayer(internalID, layerIndex, texID)
		}
		return nil, nil
	})

	// SetTerrainSplatmap(id, file$): Fallback to layer 0 only if missing.
	v.RegisterForeign("SetTerrainSplatmap", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("SetTerrainSplatmap(id, path) requires 2 arguments")
		}
		id := toInt(args[0])
		path := toString(args[1])
		idToTerrainMu.Lock()
		internalID, ok := idToTerrain[id]
		idToTerrainMu.Unlock()
		if !ok {
			return nil, nil
		}
		_ = terrain.SetTerrainSplatmap(internalID, path)
		return nil, nil
	})

	// GenerateTerrainNoise(id, seed, octaves, scale): Procedural heightmap; deterministic.
	v.RegisterForeign("GenerateTerrainNoise", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("GenerateTerrainNoise(id, seed, octaves, scale) requires 4 arguments")
		}
		id := toInt(args[0])
		idToTerrainMu.Lock()
		internalID, ok := idToTerrain[id]
		idToTerrainMu.Unlock()
		if !ok {
			return nil, fmt.Errorf("unknown terrain id %d", id)
		}
		seed := int64(toFloat64(args[1]))
		octaves := toInt(args[2])
		scale := toFloat32(args[3])
		if scale <= 0 {
			scale = 0.01
		}
		return nil, terrain.GenerateTerrainNoiseForTerrain(v, internalID, seed, octaves, float64(scale))
	})

	// DrawTerrain(id): Draw at stored position.
	v.RegisterForeign("DrawTerrain", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("DrawTerrain(id) requires 1 argument")
		}
		id := toInt(args[0])
		idToTerrainMu.Lock()
		internalID, ok := idToTerrain[id]
		idToTerrainMu.Unlock()
		if !ok {
			return nil, nil
		}
		ts := terrain.GetTerrainState(internalID)
		if ts == nil {
			return nil, nil
		}
		return nil, terrain.DrawTerrain(v, internalID, ts.PosX, ts.PosY, ts.PosZ)
	})

	v.RegisterRenderType("drawterrain", vm.Render3D)
}

// RegisterTerrain registers DBP-style terrain commands (MakeTerrain, LoadHeightmap, etc.) with integer IDs.
// Call after terrain.RegisterTerrain so DBP commands overwrite for the public API.
func RegisterTerrain(v *vm.VM) {
	registerTerrain(v)
}
