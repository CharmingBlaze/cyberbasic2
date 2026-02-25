// Package raylib: terrain, skybox, post-processing, advanced GUI aliases, terrain sculpting.
package raylib

import (
	"fmt"
	"math"

	"cyberbasic/compiler/vm"
	rl "github.com/gen2brain/raylib-go/raylib"
)

func registerAdvanced(v *vm.VM) {
	// --- Terrain (state + mesh from heightmap) ---
	v.RegisterForeign("GenerateTerrain", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("GenerateTerrain requires (width, depth, scale)")
		}
		w, d := int(toFloat64(args[0])), int(toFloat64(args[1]))
		if w <= 0 {
			w = 32
		}
		if d <= 0 {
			d = 32
		}
		scale := toFloat32(args[2])
		if scale <= 0 {
			scale = 1
		}
		terrainMu.Lock()
		terrainSeq++
		id := fmt.Sprintf("terrain_%d", terrainSeq)
		heights := make([]float32, w*d)
		terrainHeights[id] = heights
		terrainWidth[id] = w
		terrainDepth[id] = d
		terrainScale[id] = scale
		terrainMu.Unlock()
		return id, nil
	})
	v.RegisterForeign("LoadHeightmap", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("LoadHeightmap requires (imageId)")
		}
		imgId := toString(args[0])
		imageMu.Lock()
		img, ok := images[imgId]
		imageMu.Unlock()
		if !ok || img == nil {
			return nil, fmt.Errorf("unknown image id: %s", imgId)
		}
		w, h := int(img.Width), int(img.Height)
		heights := make([]float32, w*h)
		// Sample gray from image (simple: use R as height 0-1)
		for y := 0; y < h; y++ {
			for x := 0; x < w; x++ {
				c := rl.GetImageColor(*img, int32(x), int32(y))
				heights[y*w+x] = float32(c.R) / 255
			}
		}
		terrainMu.Lock()
		terrainSeq++
		id := fmt.Sprintf("terrain_%d", terrainSeq)
		terrainHeights[id] = heights
		terrainWidth[id] = w
		terrainDepth[id] = h
		terrainScale[id] = 1
		terrainMu.Unlock()
		return id, nil
	})
	v.RegisterForeign("SetTerrainTexture", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("SetTerrainTexture requires (terrainId, textureId)")
		}
		tid := toString(args[0])
		texId := toString(args[1])
		terrainMu.Lock()
		terrainTexId[tid] = texId
		terrainMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("GetTerrainHeight", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("GetTerrainHeight requires (terrainId, x, z)")
		}
		tid := toString(args[0])
		x, z := toFloat64(args[1]), toFloat64(args[2])
		terrainMu.Lock()
		heights := terrainHeights[tid]
		w := terrainWidth[tid]
		d := terrainDepth[tid]
		scale := terrainScale[tid]
		terrainMu.Unlock()
		if len(heights) == 0 || w <= 0 || d <= 0 {
			return 0.0, nil
		}
		// Map world x,z to grid; bilinear sample
		gx := (x / float64(scale)) + float64(w)/2
		gz := (z / float64(scale)) + float64(d)/2
		ix, iy := int(gx), int(gz)
		if ix < 0 {
			ix = 0
		}
		if iy < 0 {
			iy = 0
		}
		if ix >= w {
			ix = w - 1
		}
		if iy >= d {
			iy = d - 1
		}
		return float64(heights[iy*w+ix]) * float64(scale), nil
	})

	// --- Terrain sculpting ---
	v.RegisterForeign("TerrainRaise", func(args []interface{}) (interface{}, error) {
		if len(args) < 5 {
			return nil, fmt.Errorf("TerrainRaise requires (terrainId, x, z, radius, amount)")
		}
		tid := toString(args[0])
		x, z := toFloat64(args[1]), toFloat64(args[2])
		radius := toFloat32(args[3])
		amount := toFloat32(args[4])
		terrainSculpt(tid, x, z, radius, amount, 1)
		return nil, nil
	})
	v.RegisterForeign("TerrainLower", func(args []interface{}) (interface{}, error) {
		if len(args) < 5 {
			return nil, fmt.Errorf("TerrainLower requires (terrainId, x, z, radius, amount)")
		}
		tid := toString(args[0])
		x, z := toFloat64(args[1]), toFloat64(args[2])
		radius := toFloat32(args[3])
		amount := toFloat32(args[4])
		terrainSculpt(tid, x, z, radius, -amount, 1)
		return nil, nil
	})
	v.RegisterForeign("TerrainSmooth", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("TerrainSmooth requires (terrainId, x, z, radius)")
		}
		tid := toString(args[0])
		x, z := toFloat64(args[1]), toFloat64(args[2])
		radius := toFloat32(args[3])
		terrainSmooth(tid, x, z, radius)
		return nil, nil
	})
	v.RegisterForeign("TerrainFlatten", func(args []interface{}) (interface{}, error) {
		if len(args) < 5 {
			return nil, fmt.Errorf("TerrainFlatten requires (terrainId, x, z, radius, height)")
		}
		tid := toString(args[0])
		x, z := toFloat64(args[1]), toFloat64(args[2])
		radius := toFloat32(args[3])
		targetH := toFloat32(args[4])
		terrainFlatten(tid, x, z, radius, targetH)
		return nil, nil
	})
	v.RegisterForeign("TerrainPaint", func(args []interface{}) (interface{}, error) {
		if len(args) < 5 {
			return nil, fmt.Errorf("TerrainPaint requires (terrainId, x, z, radius, textureID)")
		}
		tid := toString(args[0])
		texId := toString(args[4])
		terrainMu.Lock()
		terrainTexId[tid] = texId
		terrainMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("TerrainSetMaterial", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("TerrainSetMaterial requires (terrainId, material)")
		}
		tid := toString(args[0])
		mat := toString(args[1])
		terrainMu.Lock()
		terrainMaterial[tid] = mat
		terrainMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("TerrainBrushSetSize", func(args []interface{}) (interface{}, error) {
		if len(args) >= 1 {
			terrainBrushSize = toFloat32(args[0])
			if terrainBrushSize < 0.5 {
				terrainBrushSize = 0.5
			}
		}
		return nil, nil
	})
	v.RegisterForeign("TerrainBrushSetStrength", func(args []interface{}) (interface{}, error) {
		if len(args) >= 1 {
			terrainBrushStrength = toFloat32(args[0])
		}
		return nil, nil
	})
	v.RegisterForeign("TerrainUndo", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, nil
		}
		tid := toString(args[0])
		terrainMu.Lock()
		prev := terrainUndoStack[tid]
		delete(terrainUndoStack, tid)
		if len(prev) > 0 {
			if cur, ok := terrainHeights[tid]; ok && len(cur) == len(prev) {
				copy(cur, prev)
			}
		}
		terrainMu.Unlock()
		return nil, nil
	})

	// --- Skybox ---
	v.RegisterForeign("LoadSkybox", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("LoadSkybox requires (folderPath)")
		}
		_ = toString(args[0])
		// raylib LoadTextureCubemap; store id in skyboxTexId
		skyboxTexId = ""
		return nil, nil
	})
	v.RegisterForeign("SetSkyColor", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("SetSkyColor requires (r, g, b)")
		}
		skyColorR = uint8(toFloat64(args[0])) & 0xff
		skyColorG = uint8(toFloat64(args[1])) & 0xff
		skyColorB = uint8(toFloat64(args[2])) & 0xff
		return nil, nil
	})
	v.RegisterForeign("EnableSkybox", func(args []interface{}) (interface{}, error) {
		skyboxEnabled = true
		return nil, nil
	})
	v.RegisterForeign("DisableSkybox", func(args []interface{}) (interface{}, error) {
		skyboxEnabled = false
		return nil, nil
	})

	// --- Post-processing (state only) ---
	v.RegisterForeign("EnableBloom", func(args []interface{}) (interface{}, error) {
		bloomEnabled = true
		return nil, nil
	})
	v.RegisterForeign("SetBloomIntensity", func(args []interface{}) (interface{}, error) {
		if len(args) >= 1 {
			bloomIntensity = toFloat32(args[0])
		}
		return nil, nil
	})
	v.RegisterForeign("EnableMotionBlur", func(args []interface{}) (interface{}, error) {
		motionBlurEnabled = true
		return nil, nil
	})
	v.RegisterForeign("EnableCRTFilter", func(args []interface{}) (interface{}, error) {
		crtFilterEnabled = true
		return nil, nil
	})
	v.RegisterForeign("EnablePixelate", func(args []interface{}) (interface{}, error) {
		if len(args) >= 1 {
			pixelateSize = toInt32(args[0])
			if pixelateSize <= 0 {
				pixelateSize = 4
			}
		} else {
			pixelateSize = 4
		}
		return nil, nil
	})
}

func terrainPushUndo(tid string) {
	terrainMu.Lock()
	defer terrainMu.Unlock()
	heights := terrainHeights[tid]
	if len(heights) == 0 {
		return
	}
	backup := make([]float32, len(heights))
	copy(backup, heights)
	terrainUndoStack[tid] = backup
}

func terrainSculpt(tid string, worldX, worldZ float64, radius, amount float32, sign float32) {
	terrainMu.Lock()
	heights := terrainHeights[tid]
	w := terrainWidth[tid]
	d := terrainDepth[tid]
	scale := terrainScale[tid]
	terrainMu.Unlock()
	if len(heights) == 0 || w <= 0 || d <= 0 {
		return
	}
	terrainPushUndo(tid)
	cx := (float32(worldX)/scale + float32(w)/2)
	cz := (float32(worldZ)/scale + float32(d)/2)
	r := radius
	if r < 0.5 {
		r = 0.5
	}
	str := terrainBrushStrength * float32(math.Abs(float64(amount)))
	if str > 1 {
		str = 1
	}
	terrainMu.Lock()
	defer terrainMu.Unlock()
	for iz := 0; iz < d; iz++ {
		for ix := 0; ix < w; ix++ {
			dx := float32(ix) - cx
			dz := float32(iz) - cz
			dist := float32(math.Sqrt(float64(dx*dx + dz*dz)))
			if dist > r {
				continue
			}
			falloff := float32(1) - dist/r*0.7
			if falloff < 0 {
				falloff = 0
			}
			idx := iz*w + ix
			heights[idx] += sign * str * falloff
			if heights[idx] < 0 {
				heights[idx] = 0
			}
			if heights[idx] > 1 {
				heights[idx] = 1
			}
		}
	}
}

func terrainSmooth(tid string, worldX, worldZ float64, radius float32) {
	terrainMu.Lock()
	heights := terrainHeights[tid]
	w := terrainWidth[tid]
	d := terrainDepth[tid]
	scale := terrainScale[tid]
	terrainMu.Unlock()
	if len(heights) == 0 || w <= 0 || d <= 0 {
		return
	}
	terrainPushUndo(tid)
	cx := int((float32(worldX)/scale + float32(w)/2))
	cz := int((float32(worldZ)/scale + float32(d)/2))
	r := int(radius)
	if r < 1 {
		r = 1
	}
	tmp := make([]float32, len(heights))
	copy(tmp, heights)
	terrainMu.Lock()
	defer terrainMu.Unlock()
	for iz := cz - r; iz <= cz+r; iz++ {
		for ix := cx - r; ix <= cx+r; ix++ {
			if ix < 0 || ix >= w || iz < 0 || iz >= d {
				continue
			}
			var sum float32
			var n int
			for dz := -1; dz <= 1; dz++ {
				for dx := -1; dx <= 1; dx++ {
					nx, nz := ix+dx, iz+dz
					if nx >= 0 && nx < w && nz >= 0 && nz < d {
						sum += tmp[nz*w+nx]
						n++
					}
				}
			}
			if n > 0 {
				heights[iz*w+ix] = sum / float32(n)
			}
		}
	}
}

func terrainFlatten(tid string, worldX, worldZ float64, radius, targetHeight float32) {
	terrainMu.Lock()
	heights := terrainHeights[tid]
	w := terrainWidth[tid]
	d := terrainDepth[tid]
	scale := terrainScale[tid]
	terrainMu.Unlock()
	if len(heights) == 0 || w <= 0 || d <= 0 {
		return
	}
	terrainPushUndo(tid)
	cx := (float32(worldX)/scale + float32(w)/2)
	cz := (float32(worldZ)/scale + float32(d)/2)
	r := radius
	if r < 0.5 {
		r = 0.5
	}
	str := terrainBrushStrength
	terrainMu.Lock()
	defer terrainMu.Unlock()
	for iz := 0; iz < d; iz++ {
		for ix := 0; ix < w; ix++ {
			dx := float32(ix) - cx
			dz := float32(iz) - cz
			if dx*dx+dz*dz > r*r {
				continue
			}
			idx := iz*w + ix
			heights[idx] += (targetHeight - heights[idx]) * str
			if heights[idx] < 0 {
				heights[idx] = 0
			}
			if heights[idx] > 1 {
				heights[idx] = 1
			}
		}
	}
}
