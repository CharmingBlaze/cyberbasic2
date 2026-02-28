package terrain

import (
	"fmt"

	"cyberbasic/compiler/vm"
)

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

func toFloat64(v interface{}) float64 {
	switch x := v.(type) {
	case int:
		return float64(x)
	case int32:
		return float64(x)
	case float64:
		return x
	case float32:
		return float64(x)
	default:
		return 0
	}
}

// RegisterTerrain registers heightmap and terrain mesh generation with the VM.
func RegisterTerrain(v *vm.VM) {
	// LoadHeightmap(imageID): create heightmap from raylib image (grayscale). Returns heightmap id.
	v.RegisterForeign("LoadHeightmap", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("LoadHeightmap requires (imageId)")
		}
		id, err := LoadHeightmapFromImage(toString(args[0]))
		if err != nil {
			return nil, err
		}
		return id, nil
	})

	// GenHeightmap(width, depth, noiseScale): procedural heightmap using value noise. Returns heightmap id.
	v.RegisterForeign("GenHeightmap", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("GenHeightmap requires (width, depth, noiseScale)")
		}
		w, d := int(toInt32(args[0])), int(toInt32(args[1]))
		scale := toFloat32(args[2])
		id, err := GenHeightmap(w, d, float64(scale))
		if err != nil {
			return nil, err
		}
		return id, nil
	})

	// GenHeightmapPerlin(width, depth, offsetX, offsetY, scale): create image via GenImagePerlinNoise then load as heightmap. Returns heightmap id.
	v.RegisterForeign("GenHeightmapPerlin", func(args []interface{}) (interface{}, error) {
		if len(args) < 5 {
			return nil, fmt.Errorf("GenHeightmapPerlin requires (width, depth, offsetX, offsetY, scale)")
		}
		w := toInt32(args[0])
		d := toInt32(args[1])
		ox := toInt32(args[2])
		oy := toInt32(args[3])
		scale := toFloat32(args[4])
		result, err := v.CallForeign("GenImagePerlinNoise", []interface{}{w, d, ox, oy, scale})
		if err != nil {
			return nil, err
		}
		imageID, ok := result.(string)
		if !ok || imageID == "" {
			return nil, fmt.Errorf("GenImagePerlinNoise did not return image id")
		}
		id, err := LoadHeightmapFromImage(imageID)
		if err != nil {
			return nil, err
		}
		return id, nil
	})

	// GenTerrainMesh(heightmapId, sizeX, sizeZ, heightScale [, lodLevel]): build mesh from heightmap, return mesh id.
	v.RegisterForeign("GenTerrainMesh", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("GenTerrainMesh requires (heightmapId, sizeX, sizeZ, heightScale)")
		}
		heightmapID := toString(args[0])
		sizeX := toFloat32(args[1])
		sizeZ := toFloat32(args[2])
		heightScale := toFloat32(args[3])
		lod := 0
		if len(args) >= 5 {
			lod = int(toInt32(args[4]))
		}
		id, err := GenTerrainMesh(v, heightmapID, sizeX, sizeZ, heightScale, lod)
		if err != nil {
			return nil, err
		}
		return id, nil
	})

	// --- Phase 3: High-level terrain API ---
	v.RegisterForeign("TerrainCreate", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("TerrainCreate requires (heightmapId, sizeX, sizeZ, heightScale)")
		}
		return TerrainCreate(v, toString(args[0]), toFloat32(args[1]), toFloat32(args[2]), toFloat32(args[3]))
	})
	v.RegisterForeign("TerrainUpdate", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("TerrainUpdate requires (terrainId)")
		}
		return nil, TerrainUpdate(v, toString(args[0]))
	})
	v.RegisterForeign("DrawTerrain", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("DrawTerrain requires (terrainId, posX, posY, posZ)")
		}
		return nil, DrawTerrain(v, toString(args[0]), toFloat32(args[1]), toFloat32(args[2]), toFloat32(args[3]))
	})
	v.RegisterRenderType("drawterrain", vm.Render3D)
	v.RegisterForeign("SetTerrainTexture", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("SetTerrainTexture requires (terrainId, textureId)")
		}
		return nil, SetTerrainTexture(toString(args[0]), toString(args[1]))
	})
	v.RegisterForeign("SetTerrainMaterial", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("SetTerrainMaterial requires (terrainId, materialId)")
		}
		return nil, SetTerrainMaterial(toString(args[0]), toString(args[1]))
	})
	v.RegisterForeign("SetTerrainLOD", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("SetTerrainLOD requires (terrainId, lodLevel)")
		}
		return nil, SetTerrainLOD(toString(args[0]), int(toInt32(args[1])))
	})
	v.RegisterForeign("TerrainRaise", func(args []interface{}) (interface{}, error) {
		if len(args) < 5 {
			return nil, fmt.Errorf("TerrainRaise requires (terrainId, x, z, radius, amount)")
		}
		return nil, TerrainRaise(toString(args[0]), toFloat64(args[1]), toFloat64(args[2]), toFloat64(args[3]), toFloat64(args[4]))
	})
	v.RegisterForeign("TerrainLower", func(args []interface{}) (interface{}, error) {
		if len(args) < 5 {
			return nil, fmt.Errorf("TerrainLower requires (terrainId, x, z, radius, amount)")
		}
		return nil, TerrainLower(toString(args[0]), toFloat64(args[1]), toFloat64(args[2]), toFloat64(args[3]), toFloat64(args[4]))
	})
	v.RegisterForeign("TerrainSmooth", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("TerrainSmooth requires (terrainId, x, z, radius)")
		}
		return nil, TerrainSmooth(toString(args[0]), toFloat64(args[1]), toFloat64(args[2]), toFloat64(args[3]))
	})
	v.RegisterForeign("TerrainFlatten", func(args []interface{}) (interface{}, error) {
		if len(args) < 5 {
			return nil, fmt.Errorf("TerrainFlatten requires (terrainId, x, z, radius, targetHeight)")
		}
		return nil, TerrainFlatten(toString(args[0]), toFloat64(args[1]), toFloat64(args[2]), toFloat64(args[3]), toFloat64(args[4]))
	})
	v.RegisterForeign("TerrainPaint", func(args []interface{}) (interface{}, error) {
		if len(args) < 5 {
			return nil, fmt.Errorf("TerrainPaint requires (terrainId, x, z, radius, paintValue [, blend])")
		}
		blend := 0.5
		if len(args) >= 6 {
			blend = toFloat64(args[5])
		}
		return nil, TerrainPaint(toString(args[0]), toFloat64(args[1]), toFloat64(args[2]), toFloat64(args[3]), toFloat64(args[4]), blend)
	})
	v.RegisterForeign("TerrainGetHeight", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("TerrainGetHeight requires (terrainId, x, z)")
		}
		return TerrainGetHeight(toString(args[0]), toFloat64(args[1]), toFloat64(args[2]))
	})
	v.RegisterForeign("TerrainGetNormal", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("TerrainGetNormal requires (terrainId, x, z)")
		}
		nx, ny, nz, err := TerrainGetNormal(toString(args[0]), toFloat64(args[1]), toFloat64(args[2]))
		if err != nil {
			return nil, err
		}
		return []interface{}{nx, ny, nz}, nil
	})
	v.RegisterForeign("TerrainRaycast", func(args []interface{}) (interface{}, error) {
		if len(args) < 6 {
			return nil, fmt.Errorf("TerrainRaycast requires (terrainId, ox, oy, oz, dx, dy, dz)")
		}
		hit, dist, hx, hy, hz, err := TerrainRaycast(toString(args[0]),
			toFloat64(args[1]), toFloat64(args[2]), toFloat64(args[3]),
			toFloat64(args[4]), toFloat64(args[5]), toFloat64(args[6]))
		if err != nil {
			return nil, err
		}
		if hit {
			return []interface{}{1, dist, hx, hy, hz}, nil
		}
		return []interface{}{0, 0.0, 0.0, 0.0, 0.0}, nil
	})
	v.RegisterForeign("TerrainEnableCollision", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("TerrainEnableCollision requires (terrainId, flag)")
		}
		id := toString(args[0])
		flag := toFloat32(args[1]) != 0
		terrainMu.Lock()
		defer terrainMu.Unlock()
		ts, ok := terrains[id]
		if !ok {
			return nil, fmt.Errorf("unknown terrain: %s", id)
		}
		ts.CollisionEnabled = flag
		return nil, nil
	})
	v.RegisterForeign("TerrainSetFriction", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("TerrainSetFriction requires (terrainId, value)")
		}
		id := toString(args[0])
		val := toFloat32(args[1])
		terrainMu.Lock()
		defer terrainMu.Unlock()
		ts, ok := terrains[id]
		if !ok {
			return nil, fmt.Errorf("unknown terrain: %s", id)
		}
		ts.Friction = val
		return nil, nil
	})
	v.RegisterForeign("TerrainSetBounce", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("TerrainSetBounce requires (terrainId, value)")
		}
		id := toString(args[0])
		val := toFloat32(args[1])
		terrainMu.Lock()
		defer terrainMu.Unlock()
		ts, ok := terrains[id]
		if !ok {
			return nil, fmt.Errorf("unknown terrain: %s", id)
		}
		ts.Bounce = val
		return nil, nil
	})
}
