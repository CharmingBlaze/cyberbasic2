package terrain

import (
	"fmt"
	"math"
)

// worldToNormalized maps world (x,z) to heightmap normalized [0,1] using terrain size.
func worldToNormalized(ts *TerrainState, x, z float64) (nx, nz float32) {
	nx = float32((x + float64(ts.SizeX/2)) / float64(ts.SizeX))
	nz = float32((z + float64(ts.SizeZ/2)) / float64(ts.SizeZ))
	return nx, nz
}

// gridCoord maps normalized [0,1] to heightmap grid (i, j).
func gridCoord(hm *Heightmap, nx, nz float32) (i, j int) {
	i = int(nx * float32(hm.Width-1))
	j = int(nz * float32(hm.Depth-1))
	if i < 0 {
		i = 0
	}
	if j < 0 {
		j = 0
	}
	if i >= hm.Width {
		i = hm.Width - 1
	}
	if j >= hm.Depth {
		j = hm.Depth - 1
	}
	return i, j
}

// TerrainRaise raises heights within radius of (x,z) by amount.
func TerrainRaise(terrainID string, x, z, radius, amount float64) error {
	ts := GetTerrainState(terrainID)
	if ts == nil {
		return fmt.Errorf("unknown terrain id: %s", terrainID)
	}
	hm := GetHeightmap(ts.HeightmapID)
	if hm == nil {
		return fmt.Errorf("terrain has no heightmap")
	}
	radiusSq := radius * radius
	for j := 0; j < hm.Depth; j++ {
		for i := 0; i < hm.Width; i++ {
			// Grid cell center in world (approx)
			wx := float64(i)*float64(ts.SizeX)/float64(hm.Width-1) - float64(ts.SizeX/2)
			wz := float64(j)*float64(ts.SizeZ)/float64(hm.Depth-1) - float64(ts.SizeZ/2)
			distSq := (wx-x)*(wx-x) + (wz-z)*(wz-z)
			if distSq <= radiusSq {
				idx := j*hm.Width + i
				hm.Heights[idx] += float32(amount)
				if hm.Heights[idx] < 0 {
					hm.Heights[idx] = 0
				}
				if hm.Heights[idx] > 1 {
					hm.Heights[idx] = 1
				}
			}
		}
	}
	return nil
}

// TerrainLower lowers heights within radius.
func TerrainLower(terrainID string, x, z, radius, amount float64) error {
	return TerrainRaise(terrainID, x, z, radius, -amount)
}

// TerrainSmooth smooths heights within radius (local average).
func TerrainSmooth(terrainID string, x, z, radius float64) error {
	ts := GetTerrainState(terrainID)
	if ts == nil {
		return fmt.Errorf("unknown terrain id: %s", terrainID)
	}
	hm := GetHeightmap(ts.HeightmapID)
	if hm == nil {
		return nil
	}
	radiusSq := radius * radius
	// Copy heights, then write smoothed back
	tmp := make([]float32, len(hm.Heights))
	copy(tmp, hm.Heights)
	for j := 0; j < hm.Depth; j++ {
		for i := 0; i < hm.Width; i++ {
			wx := float64(i)*float64(ts.SizeX)/float64(hm.Width-1) - float64(ts.SizeX/2)
			wz := float64(j)*float64(ts.SizeZ)/float64(hm.Depth-1) - float64(ts.SizeZ/2)
			if (wx-x)*(wx-x)+(wz-z)*(wz-z) > radiusSq {
				continue
			}
			var sum float32
			var n int
			for jj := j - 1; jj <= j+1; jj++ {
				for ii := i - 1; ii <= i+1; ii++ {
					if ii >= 0 && ii < hm.Width && jj >= 0 && jj < hm.Depth {
						sum += tmp[jj*hm.Width+ii]
						n++
					}
				}
			}
			if n > 0 {
				hm.Heights[j*hm.Width+i] = sum / float32(n)
			}
		}
	}
	return nil
}

// TerrainFlatten sets heights within radius to targetHeight (0-1).
func TerrainFlatten(terrainID string, x, z, radius, targetHeight float64) error {
	ts := GetTerrainState(terrainID)
	if ts == nil {
		return fmt.Errorf("unknown terrain id: %s", terrainID)
	}
	hm := GetHeightmap(ts.HeightmapID)
	if hm == nil {
		return nil
	}
	radiusSq := radius * radius
	th := float32(targetHeight)
	for j := 0; j < hm.Depth; j++ {
		for i := 0; i < hm.Width; i++ {
			wx := float64(i)*float64(ts.SizeX)/float64(hm.Width-1) - float64(ts.SizeX/2)
			wz := float64(j)*float64(ts.SizeZ)/float64(hm.Depth-1) - float64(ts.SizeZ/2)
			if (wx-x)*(wx-x)+(wz-z)*(wz-z) <= radiusSq {
				hm.Heights[j*hm.Width+i] = th
			}
		}
	}
	return nil
}

// TerrainPaint blends a paint value (0-1) within radius. Stores in a separate layer if we had one; for now we blend with existing height.
func TerrainPaint(terrainID string, x, z, radius, paintValue, blend float64) error {
	ts := GetTerrainState(terrainID)
	if ts == nil {
		return fmt.Errorf("unknown terrain id: %s", terrainID)
	}
	hm := GetHeightmap(ts.HeightmapID)
	if hm == nil {
		return nil
	}
	radiusSq := radius * radius
	pv := float32(paintValue)
	blendF := float32(blend)
	if blendF <= 0 {
		blendF = 0.5
	}
	for j := 0; j < hm.Depth; j++ {
		for i := 0; i < hm.Width; i++ {
			wx := float64(i)*float64(ts.SizeX)/float64(hm.Width-1) - float64(ts.SizeX/2)
			wz := float64(j)*float64(ts.SizeZ)/float64(hm.Depth-1) - float64(ts.SizeZ/2)
			distSq := (wx-x)*(wx-x) + (wz-z)*(wz-z)
			if distSq > radiusSq {
				continue
			}
			// Falloff by distance
			dist := math.Sqrt(distSq)
			factor := blendF * float32(1-dist/radius)
			if factor > 1 {
				factor = 1
			}
			idx := j*hm.Width + i
			hm.Heights[idx] = hm.Heights[idx]*(1-factor) + pv*factor
		}
	}
	return nil
}
