// Package dbp: Shared pipeline to convert model.Model -> GPU resources + DBP objects.
package dbp

import (
	"fmt"
	"path/filepath"
	"sync"

	"cyberbasic/compiler/bindings/model"
	rl "github.com/gen2brain/raylib-go/raylib"
)

// BuildResult holds the result of building a model into runtime resources.
type BuildResult struct {
	ObjectIDs   []int
	TextureIDs  []int
	MaterialIDs []int
	LightIDs    []int
}

// nextTexID, nextMatID for level-scoped resources
var (
	levelTexCounter     int
	levelMatCounter     int
	levelLightCounter   int
	levelBuildMu        sync.Mutex
	defaultTexture      rl.Texture2D
	defaultTextureOnce  sync.Once
	defaultMaterial     rl.Material
	defaultMaterialOnce sync.Once
)

// getDefaultTexture returns a 1x1 white pixel texture for fallback when texture load fails.
func getDefaultTexture() rl.Texture2D {
	defaultTextureOnce.Do(func() {
		img := rl.GenImageColor(1, 1, rl.NewColor(255, 255, 255, 255))
		defaultTexture = rl.LoadTextureFromImage(img)
		rl.UnloadImage(img)
	})
	return defaultTexture
}

// getDefaultMaterial returns white (0.8,0.8,0.8) material for fallback.
func getDefaultMaterial() rl.Material {
	defaultMaterialOnce.Do(func() {
		defaultMaterial = rl.LoadMaterialDefault()
		if defaultMaterial.Maps != nil {
			defaultMaterial.Maps.Texture = getDefaultTexture()
			defaultMaterial.Maps.Color = rl.NewColor(204, 204, 204, 255) // 0.8 * 255
		}
	})
	return defaultMaterial
}

// BuildModel converts a model.Model into GPU resources and DBP objects.
// objectIDBase: for levels use levelID*levelObjectIDBase; for single object use user id.
// basePath: directory of the source file (for resolving texture paths).
func BuildModel(m *model.Model, objectIDBase int, basePath string) (*BuildResult, error) {
	return buildModelInternal(m, objectIDBase, basePath, false)
}

// BuildModelWithHierarchy is like BuildModel but creates parent-child links from the GLTF node tree.
// Mesh nodes are parented to their nearest mesh ancestor. Use for LoadLevel when hierarchy matters.
func BuildModelWithHierarchy(m *model.Model, objectIDBase int, basePath string) (*BuildResult, error) {
	return buildModelInternal(m, objectIDBase, basePath, true)
}

func buildModelInternal(m *model.Model, objectIDBase int, basePath string, withHierarchy bool) (*BuildResult, error) {
	res := &BuildResult{}
	levelBuildMu.Lock()
	texBase := levelTexCounter
	matBase := levelMatCounter
	lightBase := levelLightCounter
	levelBuildMu.Unlock()

	// Load textures - use default on failure or empty path
	for i, tex := range m.Textures {
		t := getDefaultTexture()
		if tex.Path != "" {
			path := tex.Path
			if !filepath.IsAbs(path) && basePath != "" {
				path = filepath.Join(basePath, path)
			}
			loaded := rl.LoadTexture(path)
			if loaded.ID != 0 {
				t = loaded
			}
		}
		tid := objectIDBase + 10000 + i
		if tid <= 0 {
			tid = 1
		}
		texturesMu.Lock()
		textures[tid] = t
		texturesMu.Unlock()
		res.TextureIDs = append(res.TextureIDs, tid)
	}

	// Create materials - ensure at least one default
	materialsList := m.Materials
	if len(materialsList) == 0 {
		materialsList = []model.Material{{
			BaseColorR: 0.8, BaseColorG: 0.8, BaseColorB: 0.8, BaseColorA: 1,
			Metallic: 0, Roughness: 1, BaseColorTextureIndex: -1,
		}}
	}
	for range materialsList {
		mat := rl.LoadMaterialDefault()
		mid := objectIDBase + 20000 + len(res.MaterialIDs)
		if mid <= 0 {
			mid = 1
		}
		materialsMu.Lock()
		materials[mid] = mat
		materialsMu.Unlock()
		res.MaterialIDs = append(res.MaterialIDs, mid)
	}

	// Apply textures to materials - use default texture when index invalid
	for i, mat := range materialsList {
		if i >= len(res.MaterialIDs) {
			break
		}
		materialsMu.Lock()
		m := materials[res.MaterialIDs[i]]
		materialsMu.Unlock()
		if m.Maps != nil {
			t := getDefaultTexture()
			if mat.BaseColorTextureIndex >= 0 && mat.BaseColorTextureIndex < len(res.TextureIDs) {
				texturesMu.Lock()
				t = textures[res.TextureIDs[mat.BaseColorTextureIndex]]
				texturesMu.Unlock()
			}
			rl.SetMaterialTexture(&m, rl.MapAlbedo, t)
			if albedoMap := m.GetMap(rl.MapAlbedo); albedoMap != nil {
				albedoMap.Color = rl.NewColor(
					uint8(mat.BaseColorR*255), uint8(mat.BaseColorG*255),
					uint8(mat.BaseColorB*255), uint8(mat.BaseColorA*255),
				)
			}
			// Normal map
			if mat.NormalTextureIndex >= 0 && mat.NormalTextureIndex < len(res.TextureIDs) {
				texturesMu.Lock()
				normTex := textures[res.TextureIDs[mat.NormalTextureIndex]]
				texturesMu.Unlock()
				rl.SetMaterialTexture(&m, rl.MapNormal, normTex)
			}
			// Metallic-roughness (GLTF uses combined texture; set on both slots)
			if mat.MetallicRoughnessTextureIndex >= 0 && mat.MetallicRoughnessTextureIndex < len(res.TextureIDs) {
				texturesMu.Lock()
				mrTex := textures[res.TextureIDs[mat.MetallicRoughnessTextureIndex]]
				texturesMu.Unlock()
				rl.SetMaterialTexture(&m, rl.MapMetalness, mrTex)
				rl.SetMaterialTexture(&m, rl.MapRoughness, mrTex)
			} else {
				if metalMap := m.GetMap(rl.MapMetalness); metalMap != nil {
					metalMap.Value = mat.Metallic
				}
				if roughMap := m.GetMap(rl.MapRoughness); roughMap != nil {
					roughMap.Value = mat.Roughness
				}
			}
			materialsMu.Lock()
			materials[res.MaterialIDs[i]] = m
			materialsMu.Unlock()
		}
	}

	// Build meshes and create objects for each mesh node - skip empty meshes
	meshModels := make([]rl.Model, len(m.Meshes))
	for i := range m.Meshes {
		rlMesh, indicesKeep, err := meshToRaylib(&m.Meshes[i])
		if err != nil {
			continue // skip empty meshes, do not fail
		}
		rl.UploadMesh(&rlMesh, false)
		mod := rl.LoadModelFromMesh(rlMesh)
		meshModels[i] = mod
		_ = indicesKeep // keep alive until UploadMesh done
	}

	parentOf := make([]int, len(m.Nodes))
	for i := range parentOf {
		parentOf[i] = -1
	}
	for i, node := range m.Nodes {
		for _, childIdx := range node.Children {
			if childIdx >= 0 && childIdx < len(m.Nodes) {
				parentOf[childIdx] = i
			}
		}
	}

	nearestMeshAncestor := make([]int, len(m.Nodes))
	for i := range nearestMeshAncestor {
		nearestMeshAncestor[i] = -2
	}
	var findMeshAncestor func(int) int
	findMeshAncestor = func(idx int) int {
		if idx < 0 || idx >= len(m.Nodes) {
			return -1
		}
		if nearestMeshAncestor[idx] != -2 {
			return nearestMeshAncestor[idx]
		}
		p := parentOf[idx]
		if p < 0 {
			nearestMeshAncestor[idx] = -1
			return -1
		}
		if m.Nodes[p].MeshIndex >= 0 && m.Nodes[p].MeshIndex < len(meshModels) {
			nearestMeshAncestor[idx] = p
			return p
		}
		nearestMeshAncestor[idx] = findMeshAncestor(p)
		return nearestMeshAncestor[idx]
	}

	buildNodeTransform := func(tr model.Transform) objectTransform {
		return makeObjectTransform(tr.X, tr.Y, tr.Z, tr.Pitch, tr.Yaw, tr.Roll, tr.ScaleX, tr.ScaleY, tr.ScaleZ)
	}
	composeNodePath := func(nodeIdx, stopParent int) objectTransform {
		path := make([]int, 0, 4)
		for cur := nodeIdx; cur >= 0 && cur != stopParent; cur = parentOf[cur] {
			path = append(path, cur)
		}
		acc := identityObjectTransform()
		for i := len(path) - 1; i >= 0; i-- {
			acc = composeObjectTransform(acc, buildNodeTransform(m.Nodes[path[i]].Transform))
		}
		return acc
	}

	// Create DBP objects for nodes that have meshes
	for i, node := range m.Nodes {
		if node.MeshIndex < 0 || node.MeshIndex >= len(meshModels) {
			continue
		}
		objID := objectIDBase + i
		obj := newDbpObject(meshModels[node.MeshIndex])
		effective := buildNodeTransform(node.Transform)
		if withHierarchy {
			ancestor := findMeshAncestor(i)
			effective = composeNodePath(i, ancestor)
			obj.parentID = -1
			if ancestor >= 0 && ancestor != i {
				obj.parentID = objectIDBase + ancestor
			}
		}
		obj.x = effective.position.X
		obj.y = effective.position.Y
		obj.z = effective.position.Z
		obj.pitch, obj.yaw, obj.roll = quaternionToDegrees(effective.rotation)
		obj.scaleX = effective.scale.X
		obj.scaleY = effective.scale.Y
		obj.scaleZ = effective.scale.Z
		matIdx := m.Meshes[node.MeshIndex].MaterialIndex
		if matIdx < 0 || matIdx >= len(res.MaterialIDs) {
			matIdx = 0
		}
		if matIdx < len(res.MaterialIDs) {
			matID := res.MaterialIDs[matIdx]
			materialsMu.Lock()
			mat, ok := materials[matID]
			materialsMu.Unlock()
			if ok && mat.Maps != nil {
				obj.model.Materials = &mat
			} else {
				_ = getDefaultMaterial() // ensure init
				obj.model.Materials = &defaultMaterial
			}
			if matIdx >= 0 && matIdx < len(materialsList) {
				srcMat := materialsList[matIdx]
				obj.roughness = srcMat.Roughness
				obj.metallic = srcMat.Metallic
				obj.emissiveR = uint8(srcMat.EmissiveFactorR * 255)
				obj.emissiveG = uint8(srcMat.EmissiveFactorG * 255)
				obj.emissiveB = uint8(srcMat.EmissiveFactorB * 255)
			}
		}
		objectsMu.Lock()
		objects[objID] = obj
		objectsMu.Unlock()
		res.ObjectIDs = append(res.ObjectIDs, objID)
	}

	// Create lights
	for i, light := range m.Lights {
		lightsMu.Lock()
		lid := objectIDBase + 30000 + lightBase + i
		if lid <= 0 {
			lid = 1
		}
		lights[lid] = &dbpLight{
			lightType: light.Type,
			x:         light.X, y: light.Y, z: light.Z,
			r: light.R, g: light.G, b: light.B,
			intensity: light.Intensity,
			range_:    light.Range,
		}
		lightsMu.Unlock()
		res.LightIDs = append(res.LightIDs, lid)
	}
	if len(m.Lights) > 0 {
		syncRendererShadowLights()
	}

	_ = texBase
	_ = matBase
	return res, nil
}

func meshToRaylib(m *model.Mesh) (rl.Mesh, []uint16, error) {
	vCount := len(m.Vertices) / 3
	if vCount == 0 {
		return rl.Mesh{}, nil, fmt.Errorf("mesh has no vertices")
	}
	if len(m.Normals) < vCount*3 {
		model.ComputeFlatNormals(m)
	}
	normals := m.Normals
	texcoords := m.Texcoords
	if len(texcoords) < vCount*2 {
		texcoords = make([]float32, vCount*2)
	}
	triCount := vCount / 3
	var indicesKeep []uint16
	mesh := rl.Mesh{
		VertexCount:   int32(vCount),
		TriangleCount: int32(triCount),
		Vertices:      &m.Vertices[0],
		Normals:       &normals[0],
		Texcoords:     &texcoords[0],
	}
	if len(m.Indices) > 0 {
		triCount = len(m.Indices) / 3
		mesh.TriangleCount = int32(triCount)
		if len(m.Indices) <= 65535 {
			indicesKeep = make([]uint16, len(m.Indices))
			for i, idx := range m.Indices {
				indicesKeep[i] = uint16(idx)
			}
			mesh.Indices = &indicesKeep[0]
		}
	}
	return mesh, indicesKeep, nil
}
