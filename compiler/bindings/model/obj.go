// Package model: OBJ importer using flywave/go-obj.
package model

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/flywave/go-obj"
)

// importOBJ loads an OBJ file into the canonical Model struct.
func importOBJ(path string) (*Model, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open obj: %w", err)
	}
	defer f.Close()

	reader := &obj.ObjReader{}
	if err := reader.Read(f); err != nil {
		return nil, fmt.Errorf("read obj: %w", err)
	}

	buf := &reader.ObjBuffer
	baseDir := filepath.Dir(path)

	m := &Model{
		Meshes:    make([]Mesh, 0),
		Materials: make([]Material, 0),
		Textures:  make([]Texture, 0),
		Nodes:     make([]Node, 0),
		Lights:    make([]Light, 0),
		Colliders: make([]Collider, 0),
	}

	// Load MTL if present
	var mtlMap map[string]*obj.Material
	if buf.MTL != "" {
		mtlPath := filepath.Join(baseDir, buf.MTL)
		mtlMap, _ = obj.ReadMaterials(mtlPath)
	}

	// Convert faces to indexed mesh
	// OBJ uses 1-based indices; we need to build vertex/normal/uv arrays
	type vertexKey struct {
		v, vn, vt int
	}
	vertexMap := make(map[vertexKey]int)
	var vertices []float32
	var normals []float32
	var texcoords []float32
	var indices []uint32

	for _, face := range buf.F {
		tris := face.Triangulate(buf.V)
		for _, tri := range tris {
			for _, corner := range tri {
				vIdx := corner.VertexIndex
				vnIdx := -1
				vtIdx := -1
				if corner.NormalIndex >= 0 {
					vnIdx = corner.NormalIndex
				}
				if corner.TexCoordIndex >= 0 {
					vtIdx = corner.TexCoordIndex
				}
				key := vertexKey{v: vIdx, vn: vnIdx, vt: vtIdx}
				if idx, ok := vertexMap[key]; ok {
					indices = append(indices, uint32(idx))
				} else {
					idx := len(vertexMap)
					vertexMap[key] = idx
					indices = append(indices, uint32(idx))
					if vIdx >= 0 && vIdx < len(buf.V) {
						v := buf.V[vIdx]
						vertices = append(vertices, float32(v[0]), float32(v[1]), float32(v[2]))
					}
					if vnIdx >= 0 && vnIdx < len(buf.VN) {
						vn := buf.VN[vnIdx]
						normals = append(normals, float32(vn[0]), float32(vn[1]), float32(vn[2]))
					} else {
						normals = append(normals, 0, 0, 0)
					}
					if vtIdx >= 0 && vtIdx < len(buf.VT) {
						vt := buf.VT[vtIdx]
						texcoords = append(texcoords, float32(vt[0]), float32(vt[1]))
					} else {
						texcoords = append(texcoords, 0, 0)
					}
				}
			}
		}
	}

	// Ensure normals/texcoords match vertex count
	vCount := len(vertices) / 3
	for len(normals) < vCount*3 {
		normals = append(normals, 0, 0, 0)
	}
	for len(texcoords) < vCount*2 {
		texcoords = append(texcoords, 0, 0)
	}

	if len(vertices) == 0 {
		return nil, fmt.Errorf("obj has no geometry")
	}

	mesh := Mesh{
		Vertices:     vertices,
		Normals:      normals,
		Texcoords:    texcoords,
		Indices:      indices,
		MaterialIndex: -1,
	}
	// Compute flat normals when OBJ had no normals (all zeros)
	hasRealNormals := false
	for i := 0; i < len(normals) && !hasRealNormals; i++ {
		if normals[i] != 0 {
			hasRealNormals = true
		}
	}
	if !hasRealNormals {
		ComputeFlatNormals(&mesh)
	}
	m.Meshes = append(m.Meshes, mesh)

	// Add default material if we have MTL
	if len(mtlMap) > 0 {
		for _, gmat := range mtlMap {
			mat := Material{
				BaseColorR: 1, BaseColorG: 1, BaseColorB: 1, BaseColorA: 1,
				Metallic: 0, Roughness: 1,
				BaseColorTextureIndex: -1,
			}
			if len(gmat.Diffuse) >= 3 {
				mat.BaseColorR = gmat.Diffuse[0]
				mat.BaseColorG = gmat.Diffuse[1]
				mat.BaseColorB = gmat.Diffuse[2]
			}
			if gmat.DiffuseTexture != "" {
				m.Textures = append(m.Textures, Texture{
					Path: filepath.Join(baseDir, gmat.DiffuseTexture),
				})
				mat.BaseColorTextureIndex = len(m.Textures) - 1
			}
			m.Materials = append(m.Materials, mat)
			break
		}
	} else {
		m.Materials = append(m.Materials, Material{
			BaseColorR: 0.8, BaseColorG: 0.8, BaseColorB: 0.8, BaseColorA: 1,
			Metallic: 0, Roughness: 1, BaseColorTextureIndex: -1,
		})
	}
	// Ensure mesh uses valid material index
	if m.Meshes[0].MaterialIndex < 0 {
		m.Meshes[0].MaterialIndex = 0
	}

	// One root node
	m.Nodes = append(m.Nodes, Node{
		Name:       "root",
		Transform:  DefaultTransform(),
		MeshIndex:  0,
		Children:   nil,
	})

	return m, nil
}
