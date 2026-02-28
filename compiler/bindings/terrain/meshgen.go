package terrain

import (
	"fmt"
	"math"

	"cyberbasic/compiler/vm"
)

// GenTerrainMesh builds a mesh from a heightmap and registers it via the VM's MeshCreate.
// sizeX, sizeZ are world-space size; heightScale scales height values.
// lodLevel 0 = full resolution; each increment typically halves resolution.
func GenTerrainMesh(v *vm.VM, heightmapID string, sizeX, sizeZ, heightScale float32, lodLevel int) (string, error) {
	hm := GetHeightmap(heightmapID)
	if hm == nil {
		return "", fmt.Errorf("unknown heightmap id: %s", heightmapID)
	}
	w, d := hm.Width, hm.Depth
	if w < 2 || d < 2 {
		return "", fmt.Errorf("heightmap too small")
	}
	// LOD: skip vertices
	step := 1
	for i := 0; i < lodLevel && step*2 < w && step*2 < d; i++ {
		step *= 2
	}
	// Grid dimensions after LOD
	nx := (w-1)/step + 1
	nz := (d-1)/step + 1
	vertexCount := nx * nz
	// Allocate buffers
	vertices := make([]float32, vertexCount*3)
	normals := make([]float32, vertexCount*3)
	uvs := make([]float32, vertexCount*2)
	indices := make([]uint16, 0, (nx-1)*(nz-1)*6)
	stepX := sizeX / float32(nx-1)
	stepZ := sizeZ / float32(nz-1)
	for j := 0; j < nz; j++ {
		for i := 0; i < nx; i++ {
			gi, gj := i*step, j*step
			if gi >= w {
				gi = w - 1
			}
			if gj >= d {
				gj = d - 1
			}
			idx := gj*w + gi
			h := hm.Heights[idx]
			vidx := (j*nx + i) * 3
			vertices[vidx] = float32(i)*stepX - sizeX/2
			vertices[vidx+1] = h * heightScale
			vertices[vidx+2] = float32(j)*stepZ - sizeZ/2
			// UV
			uvidx := (j*nx + i) * 2
			uvs[uvidx] = float32(i) / float32(nx-1)
			uvs[uvidx+1] = float32(j) / float32(nz-1)
			// Normal from finite difference
			var dx, dz float32
			if gi+step < w {
				dx = (hm.Heights[gj*w+gi+step] - h) * heightScale / (float32(step) * stepX)
			} else if gi > 0 {
				dx = (h - hm.Heights[gj*w+gi-step]) * heightScale / (float32(step) * stepX)
			}
			if gj+step < d {
				dz = (hm.Heights[(gj+step)*w+gi] - h) * heightScale / (float32(step) * stepZ)
			} else if gj > 0 {
				dz = (h - hm.Heights[(gj-step)*w+gi]) * heightScale / (float32(step) * stepZ)
			}
			// tangentZ = (0, dz, 1), tangentX = (1, dx, 0) -> normal = cross(Z,X) = (-dx, 1, -dz) then normalize
			nxNorm := -dx
			nyNorm := float32(1)
			nzNorm := -dz
			lenNorm := float32(math.Sqrt(float64(nxNorm*nxNorm + nyNorm*nyNorm + nzNorm*nzNorm)))
			if lenNorm < 1e-6 {
				lenNorm = 1
			}
			normals[vidx] = nxNorm / lenNorm
			normals[vidx+1] = nyNorm / lenNorm
			normals[vidx+2] = nzNorm / lenNorm
		}
	}
	// Indices: two triangles per quad
	for j := 0; j < nz-1; j++ {
		for i := 0; i < nx-1; i++ {
			a := uint16(j*nx + i)
			b := uint16(j*nx + (i + 1))
			c := uint16((j+1)*nx + i)
			d := uint16((j+1)*nx + (i + 1))
			indices = append(indices, a, c, b, b, c, d)
		}
	}
	// Convert to []interface{} for VM
	vertsIf := make([]interface{}, len(vertices))
	for i, v := range vertices {
		vertsIf[i] = v
	}
	normsIf := make([]interface{}, len(normals))
	for i, n := range normals {
		normsIf[i] = n
	}
	uvsIf := make([]interface{}, len(uvs))
	for i, u := range uvs {
		uvsIf[i] = u
	}
	indicesIf := make([]interface{}, len(indices))
	for i, idx := range indices {
		indicesIf[i] = int(idx)
	}
	result, err := v.CallForeign("MeshCreate", []interface{}{vertsIf, normsIf, uvsIf, indicesIf})
	if err != nil {
		return "", err
	}
	if id, ok := result.(string); ok {
		return id, nil
	}
	return "", fmt.Errorf("MeshCreate did not return mesh id")
}
