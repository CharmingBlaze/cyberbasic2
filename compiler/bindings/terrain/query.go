package terrain

import (
	"fmt"
	"math"
)

// TerrainGetHeight returns world Y at (x,z) using bilinear sampling. Returns 0 if terrain invalid.
func TerrainGetHeight(terrainID string, x, z float64) (float64, error) {
	ts := GetTerrainState(terrainID)
	if ts == nil {
		return 0, fmt.Errorf("unknown terrain id: %s", terrainID)
	}
	hm := GetHeightmap(ts.HeightmapID)
	if hm == nil {
		return 0, nil
	}
	nx, nz := worldToNormalized(ts, x, z)
	h := hm.SampleHeight(nx, nz)
	return float64(h) * float64(ts.HeightScale), nil
}

// TerrainGetNormal returns the normal vector at (x,z) as (nx, ny, nz). Approximated from gradient.
func TerrainGetNormal(terrainID string, x, z float64) (nx, ny, nz float64, err error) {
	ts := GetTerrainState(terrainID)
	if ts == nil {
		return 0, 1, 0, fmt.Errorf("unknown terrain id: %s", terrainID)
	}
	hm := GetHeightmap(ts.HeightmapID)
	if hm == nil {
		return 0, 1, 0, nil
	}
	delta := 0.01
	h, _ := TerrainGetHeight(terrainID, x, z)
	hx, _ := TerrainGetHeight(terrainID, x+delta, z)
	hz, _ := TerrainGetHeight(terrainID, x, z+delta)
	dx := (hx - h) / delta
	dz := (hz - h) / delta
	// Normal = cross(Z_tangent, X_tangent) = (-dx, 1, -dz) normalized
	nx = -dx
	ny = 1
	nz = -dz
	len := math.Sqrt(nx*nx + ny*ny + nz*nz)
	if len < 1e-10 {
		len = 1
	}
	return nx / len, ny / len, nz / len, nil
}

// TerrainRaycast performs a ray-vs-heightmap intersection. Origin and direction are 3D vectors.
// Returns hit (true/false), distance, and hit position (x,y,z).
func TerrainRaycast(terrainID string, ox, oy, oz, dx, dy, dz float64) (hit bool, dist float64, hx, hy, hz float64, err error) {
	ts := GetTerrainState(terrainID)
	if ts == nil {
		return false, 0, 0, 0, 0, fmt.Errorf("unknown terrain id: %s", terrainID)
	}
	hm := GetHeightmap(ts.HeightmapID)
	if hm == nil {
		return false, 0, 0, 0, 0, nil
	}
	// Step along ray and sample height
	step := 0.5
	maxDist := 10000.0
	lenDir := math.Sqrt(dx*dx + dy*dy + dz*dz)
	if lenDir < 1e-10 {
		return false, 0, 0, 0, 0, nil
	}
	dx /= lenDir
	dy /= lenDir
	dz /= lenDir
	for d := 0.0; d < maxDist; d += step {
		px := ox + d*dx
		py := oy + d*dy
		pz := oz + d*dz
		terrainY, _ := TerrainGetHeight(terrainID, px, pz)
		if py <= terrainY {
			// Binary search refinement
			lo, hi := d-step, d
			for i := 0; i < 10; i++ {
				mid := (lo + hi) / 2
				px := ox + mid*dx
				py := oy + mid*dy
				pz := oz + mid*dz
				ty, _ := TerrainGetHeight(terrainID, px, pz)
				if py <= ty {
					hi = mid
				} else {
					lo = mid
				}
			}
			dist = (lo + hi) / 2
			hx = ox + dist*dx
			hy = oy + dist*dy
			hz = oz + dist*dz
			return true, dist, hx, hy, hz, nil
		}
	}
	return false, 0, 0, 0, 0, nil
}
