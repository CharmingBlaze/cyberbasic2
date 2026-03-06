// Package model: GLTF/GLB importer using qmuntal/gltf.
package model

import (
	"fmt"
	"math"
	"path/filepath"
	"strings"

	"github.com/qmuntal/gltf"
	"github.com/qmuntal/gltf/modeler"
)

// importGLTF loads a GLTF or GLB file into the canonical Model struct.
func importGLTF(path string) (*Model, error) {
	doc, err := gltf.Open(path)
	if err != nil {
		return nil, fmt.Errorf("gltf open: %w", err)
	}
	baseDir := filepath.Dir(path)

	m := &Model{
		Meshes:    make([]Mesh, 0),
		Materials: make([]Material, 0),
		Textures:  make([]Texture, 0),
		Nodes:     make([]Node, 0),
		Lights:    make([]Light, 0),
		Colliders: make([]Collider, 0),
	}

	// Import materials - add default if none
	for i, gmat := range doc.Materials {
		mat := importMaterial(gmat, i)
		m.Materials = append(m.Materials, mat)
	}
	if len(m.Materials) == 0 {
		m.Materials = append(m.Materials, Material{
			BaseColorR: 0.8, BaseColorG: 0.8, BaseColorB: 0.8, BaseColorA: 1,
			Metallic: 0, Roughness: 1, BaseColorTextureIndex: -1,
		})
	}

	// Import textures (from images)
	for _, img := range doc.Images {
		texPath := img.URI
		if texPath == "" && img.BufferView != nil {
			// Embedded image - skip for now (would need to decode)
			continue
		}
		if texPath != "" && !filepath.IsAbs(texPath) {
			texPath = filepath.Join(baseDir, texPath)
		}
		m.Textures = append(m.Textures, Texture{Path: texPath})
	}

	// Import meshes (one per GLTF mesh, using first primitive)
	for _, gmesh := range doc.Meshes {
		if len(gmesh.Primitives) == 0 {
			continue
		}
		mesh, err := importPrimitive(doc, gmesh.Primitives[0])
		if err != nil {
			return nil, fmt.Errorf("mesh %s: %w", gmesh.Name, err)
		}
		m.Meshes = append(m.Meshes, mesh)
	}

	// Import nodes (scene graph)
	sceneIdx := 0
	if doc.Scene != nil {
		sceneIdx = *doc.Scene
	}
	if sceneIdx < len(doc.Scenes) {
		scene := doc.Scenes[sceneIdx]
		for _, nodeIdx := range scene.Nodes {
			importNodes(doc, doc.Nodes, nodeIdx, m)
		}
	}
	if len(m.Nodes) == 0 && len(m.Meshes) > 0 {
		// Flat: one node per mesh
		for i := range m.Meshes {
			m.Nodes = append(m.Nodes, Node{
				Name:      fmt.Sprintf("mesh_%d", i),
				Transform: DefaultTransform(),
				MeshIndex: i,
			})
		}
	}

	// Import skeleton from first skin (if any)
	if len(doc.Skins) > 0 {
		sk := doc.Skins[0]
		skeleton := &Skeleton{Bones: make([]Bone, 0, len(sk.Joints))}
		invBind := make([][16]float32, 0)
		if sk.InverseBindMatrices != nil {
			acr := doc.Accessors[*sk.InverseBindMatrices]
			matrices, err := modeler.ReadInverseBindMatrices(doc, acr, nil)
			if err == nil {
				for _, mat := range matrices {
					var m [16]float32
					for col := 0; col < 4; col++ {
						for row := 0; row < 4; row++ {
							m[col*4+row] = mat[row][col]
						}
					}
					invBind = append(invBind, m)
				}
			}
		}
		jointSet := make(map[int]int) // nodeIdx -> bone index
		for i, jointIdx := range sk.Joints {
			jointSet[jointIdx] = i
		}
		for i, jointIdx := range sk.Joints {
			bone := Bone{Parent: -1}
			if jointIdx >= 0 && jointIdx < len(doc.Nodes) {
				bone.Name = doc.Nodes[jointIdx].Name
				for p, pnode := range doc.Nodes {
					for _, c := range pnode.Children {
						if c == jointIdx {
							if bi, ok := jointSet[p]; ok {
								bone.Parent = bi
							}
							break
						}
					}
				}
			}
			if i < len(invBind) {
				bone.InverseBind = invBind[i]
			}
			skeleton.Bones = append(skeleton.Bones, bone)
		}
		if len(skeleton.Bones) > 0 {
			m.Skeleton = skeleton
		}
	}

	// Import animations
	jointSet := make(map[int]int)
	if m.Skeleton != nil && len(doc.Skins) > 0 {
		for i, jointIdx := range doc.Skins[0].Joints {
			jointSet[jointIdx] = i
		}
	}
	for _, ganim := range doc.Animations {
		anim := importAnimation(doc, ganim, jointSet)
		if len(anim.Channels) > 0 {
			m.Animations = append(m.Animations, anim)
		}
	}

	// Import lights from KHR_lights_punctual extension
	if ext, ok := doc.Extensions["KHR_lights_punctual"]; ok {
		if lights, ok := ext.(map[string]any); ok {
			_ = lights // TODO: parse light nodes
		}
	}

	return m, nil
}

func importAnimation(doc *gltf.Document, ganim *gltf.Animation, jointSet map[int]int) Animation {
	anim := Animation{
		Name:     ganim.Name,
		Duration: 0,
		Channels: make([]AnimationChannel, 0),
	}
	for _, ch := range ganim.Channels {
		if ch.Sampler >= len(ganim.Samplers) {
			continue
		}
		sampler := ganim.Samplers[ch.Sampler]
		if sampler.Input >= len(doc.Accessors) || sampler.Output >= len(doc.Accessors) {
			continue
		}
		inputAcr := doc.Accessors[sampler.Input]
		outputAcr := doc.Accessors[sampler.Output]
		times, err := readAccessorFloat32(doc, inputAcr)
		if err != nil || len(times) == 0 {
			continue
		}
		values, err := readAccessorVec(doc, outputAcr, ch.Target.Path)
		if err != nil || len(values) == 0 {
			continue
		}
		if len(times) != len(values) {
			continue
		}
		keyframes := make([]Keyframe, len(times))
		for i := range times {
			keyframes[i] = Keyframe{Time: times[i], Value: values[i]}
		}
		nodeIdx := -1
		if ch.Target.Node != nil {
			nodeIdx = *ch.Target.Node
		}
		boneIdx := -1
		if nodeIdx >= 0 {
			if bi, ok := jointSet[nodeIdx]; ok {
				boneIdx = bi
			}
		}
		prop := string(ch.Target.Path)
		if prop == "" {
			prop = "translation"
		}
		channel := AnimationChannel{
			NodeIndex:  nodeIdx,
			BoneIndex:  boneIdx,
			Property:   prop,
			Keyframes:  keyframes,
		}
		anim.Channels = append(anim.Channels, channel)
		if times[len(times)-1] > anim.Duration {
			anim.Duration = times[len(times)-1]
		}
	}
	return anim
}

func readAccessorFloat32(doc *gltf.Document, acr *gltf.Accessor) ([]float32, error) {
	data, err := modeler.ReadAccessor(doc, acr, nil)
	if err != nil {
		return nil, err
	}
	switch v := data.(type) {
	case []float32:
		return v, nil
	case []float64:
		out := make([]float32, len(v))
		for i, x := range v {
			out[i] = float32(x)
		}
		return out, nil
	default:
		return nil, fmt.Errorf("unexpected accessor type for time")
	}
}

func readAccessorVec(doc *gltf.Document, acr *gltf.Accessor, path gltf.TRSProperty) ([][]float32, error) {
	data, err := modeler.ReadAccessor(doc, acr, nil)
	if err != nil {
		return nil, err
	}
	switch path {
	case gltf.TRSTranslation, gltf.TRSScale:
		// VEC3
		switch v := data.(type) {
		case [][3]float32:
			out := make([][]float32, len(v))
			for i, x := range v {
				out[i] = []float32{x[0], x[1], x[2]}
			}
			return out, nil
		case [][3]float64:
			out := make([][]float32, len(v))
			for i, x := range v {
				out[i] = []float32{float32(x[0]), float32(x[1]), float32(x[2])}
			}
			return out, nil
		}
	case gltf.TRSRotation:
		// VEC4 quaternion
		switch v := data.(type) {
		case [][4]float32:
			out := make([][]float32, len(v))
			for i, x := range v {
				out[i] = []float32{x[0], x[1], x[2], x[3]}
			}
			return out, nil
		case [][4]float64:
			out := make([][]float32, len(v))
			for i, x := range v {
				out[i] = []float32{float32(x[0]), float32(x[1]), float32(x[2]), float32(x[3])}
			}
			return out, nil
		}
	}
	return nil, fmt.Errorf("unsupported path %s", path)
}

func importMaterial(gmat *gltf.Material, idx int) Material {
	mat := Material{
		BaseColorR: 1, BaseColorG: 1, BaseColorB: 1, BaseColorA: 1,
		Metallic:  0, Roughness: 1,
		BaseColorTextureIndex: -1, NormalTextureIndex: -1, MetallicRoughnessTextureIndex: -1,
		EmissiveTextureIndex: -1,
	}
	if gmat.PBRMetallicRoughness != nil {
		pbr := gmat.PBRMetallicRoughness
		if len(pbr.BaseColorFactor) >= 4 {
			mat.BaseColorR = float32(pbr.BaseColorFactor[0])
			mat.BaseColorG = float32(pbr.BaseColorFactor[1])
			mat.BaseColorB = float32(pbr.BaseColorFactor[2])
			mat.BaseColorA = float32(pbr.BaseColorFactor[3])
		}
		if pbr.MetallicFactor != nil {
			mat.Metallic = float32(*pbr.MetallicFactor)
		}
		if pbr.RoughnessFactor != nil {
			mat.Roughness = float32(*pbr.RoughnessFactor)
		}
		if pbr.BaseColorTexture != nil && pbr.BaseColorTexture.Index >= 0 {
			mat.BaseColorTextureIndex = pbr.BaseColorTexture.Index
		}
		if pbr.MetallicRoughnessTexture != nil && pbr.MetallicRoughnessTexture.Index >= 0 {
			mat.MetallicRoughnessTextureIndex = pbr.MetallicRoughnessTexture.Index
		}
	}
	if gmat.NormalTexture != nil && gmat.NormalTexture.Index != nil {
		mat.NormalTextureIndex = *gmat.NormalTexture.Index
	}
	if len(gmat.EmissiveFactor) >= 3 {
		mat.EmissiveFactorR = float32(gmat.EmissiveFactor[0])
		mat.EmissiveFactorG = float32(gmat.EmissiveFactor[1])
		mat.EmissiveFactorB = float32(gmat.EmissiveFactor[2])
	}
	if gmat.EmissiveTexture != nil && gmat.EmissiveTexture.Index >= 0 {
		mat.EmissiveTextureIndex = gmat.EmissiveTexture.Index
	}
	return mat
}

func importPrimitive(doc *gltf.Document, prim *gltf.Primitive) (Mesh, error) {
	mesh := Mesh{MaterialIndex: 0}
	if prim.Material != nil {
		idx := *prim.Material
		if idx >= 0 && idx < len(doc.Materials) {
			mesh.MaterialIndex = idx
		}
	}

	// Positions (required)
	posIdx, ok := prim.Attributes[gltf.POSITION]
	if !ok {
		return mesh, fmt.Errorf("primitive has no POSITION")
	}
	posAcr := doc.Accessors[posIdx]
	positions, err := modeler.ReadPosition(doc, posAcr, nil)
	if err != nil {
		return mesh, fmt.Errorf("read positions: %w", err)
	}
	mesh.Vertices = make([]float32, len(positions)*3)
	for i, p := range positions {
		mesh.Vertices[i*3] = p[0]
		mesh.Vertices[i*3+1] = p[1]
		mesh.Vertices[i*3+2] = p[2]
	}

	// Normals (optional) - compute flat normals when missing
	if normIdx, ok := prim.Attributes[gltf.NORMAL]; ok {
		normAcr := doc.Accessors[normIdx]
		normals, err := modeler.ReadNormal(doc, normAcr, nil)
		if err == nil && len(normals) > 0 {
			mesh.Normals = make([]float32, len(normals)*3)
			for i, n := range normals {
				mesh.Normals[i*3] = n[0]
				mesh.Normals[i*3+1] = n[1]
				mesh.Normals[i*3+2] = n[2]
			}
		}
	}
	if len(mesh.Normals) < len(mesh.Vertices) {
		ComputeFlatNormals(&mesh)
	}

	// Texcoords (optional)
	if tcIdx, ok := prim.Attributes[gltf.TEXCOORD_0]; ok {
		tcAcr := doc.Accessors[tcIdx]
		tcs, err := modeler.ReadTextureCoord(doc, tcAcr, nil)
		if err == nil {
			mesh.Texcoords = make([]float32, len(tcs)*2)
			for i, t := range tcs {
				mesh.Texcoords[i*2] = t[0]
				mesh.Texcoords[i*2+1] = t[1]
			}
		}
	}
	if len(mesh.Texcoords) == 0 {
		mesh.Texcoords = make([]float32, len(mesh.Vertices)/3*2)
	}

	// Indices (optional)
	if prim.Indices != nil {
		idxAcr := doc.Accessors[*prim.Indices]
		mesh.Indices, err = modeler.ReadIndices(doc, idxAcr, nil)
		if err != nil {
			return mesh, fmt.Errorf("read indices: %w", err)
		}
	}

	// Skinning (optional) - JOINTS_0 and WEIGHTS_0
	if jointsIdx, ok := prim.Attributes[gltf.JOINTS_0]; ok {
		jointsAcr := doc.Accessors[jointsIdx]
		joints, jerr := modeler.ReadJoints(doc, jointsAcr, nil)
		if jerr == nil && len(joints) > 0 {
			mesh.BoneIndices = make([]uint8, len(joints)*4)
			for i, j := range joints {
				mesh.BoneIndices[i*4] = uint8(j[0])
				mesh.BoneIndices[i*4+1] = uint8(j[1])
				mesh.BoneIndices[i*4+2] = uint8(j[2])
				mesh.BoneIndices[i*4+3] = uint8(j[3])
			}
		}
	}
	if weightsIdx, ok := prim.Attributes[gltf.WEIGHTS_0]; ok {
		weightsAcr := doc.Accessors[weightsIdx]
		weights, werr := modeler.ReadWeights(doc, weightsAcr, nil)
		if werr == nil && len(weights) > 0 {
			mesh.BoneWeights = make([]float32, len(weights)*4)
			for i, w := range weights {
				mesh.BoneWeights[i*4] = w[0]
				mesh.BoneWeights[i*4+1] = w[1]
				mesh.BoneWeights[i*4+2] = w[2]
				mesh.BoneWeights[i*4+3] = w[3]
			}
		}
	}

	return mesh, nil
}

// isCollisionNode returns true if the node name or extras indicate a collision volume.
func isCollisionNode(name string, extras any) bool {
	n := strings.TrimSpace(strings.ToLower(name))
	if strings.HasPrefix(n, "col_") || strings.HasPrefix(n, "collision_") || strings.HasPrefix(n, "collision") {
		return true
	}
	if extras != nil {
		if m, ok := extras.(map[string]any); ok {
			if v, ok := m["collision"]; ok {
				if b, ok := v.(bool); ok && b {
					return true
				}
			}
		}
	}
	return false
}

// colliderTypeFromName infers ColliderBox, ColliderSphere, or ColliderCapsule from name.
func colliderTypeFromName(name string) int {
	n := strings.ToLower(name)
	if strings.Contains(n, "sphere") {
		return ColliderSphere
	}
	if strings.Contains(n, "capsule") {
		return ColliderCapsule
	}
	return ColliderBox
}

func importNodes(doc *gltf.Document, nodes []*gltf.Node, nodeIdx int, m *Model) {
	if nodeIdx < 0 || nodeIdx >= len(nodes) {
		return
	}
	gn := nodes[nodeIdx]
	tr := nodeToTransform(gn)
	meshIdx := -1
	if gn.Mesh != nil {
		// Map gltf mesh index to our mesh index (we flatten primitives)
		meshIdx = *gn.Mesh
		if meshIdx >= len(m.Meshes) {
			meshIdx = -1
		}
	}
	// Collision detection: add to Colliders if node name/extras indicate collision
	if isCollisionNode(gn.Name, gn.Extras) {
		col := Collider{
			Type:      colliderTypeFromName(gn.Name),
			Transform: tr,
			MeshIndex: meshIdx,
		}
		if meshIdx >= 0 && meshIdx < len(m.Meshes) {
			minX, minY, minZ, maxX, maxY, maxZ := MeshBounds(&m.Meshes[meshIdx])
			sx := (maxX - minX) / 2
			sy := (maxY - minY) / 2
			sz := (maxZ - minZ) / 2
			if col.Type == ColliderBox {
				col.SizeX = sx * tr.ScaleX
				col.SizeY = sy * tr.ScaleY
				col.SizeZ = sz * tr.ScaleZ
				if col.SizeX < 0.01 {
					col.SizeX = 0.5
				}
				if col.SizeY < 0.01 {
					col.SizeY = 0.5
				}
				if col.SizeZ < 0.01 {
					col.SizeZ = 0.5
				}
			} else if col.Type == ColliderSphere {
				r := sx
				if sy > r {
					r = sy
				}
				if sz > r {
					r = sz
				}
				col.Radius = r * tr.ScaleX
				if col.Radius < 0.01 {
					col.Radius = 0.5
				}
			} else {
				col.Radius = sx * tr.ScaleX
				col.Height = (maxY - minY) * tr.ScaleY
				if col.Radius < 0.01 {
					col.Radius = 0.5
				}
				if col.Height < 0.01 {
					col.Height = 1
				}
			}
		} else {
			// No mesh: use scale as size
			col.SizeX = tr.ScaleX
			col.SizeY = tr.ScaleY
			col.SizeZ = tr.ScaleZ
			if col.SizeX < 0.01 {
				col.SizeX = 0.5
			}
			if col.SizeY < 0.01 {
				col.SizeY = 0.5
			}
			if col.SizeZ < 0.01 {
				col.SizeZ = 0.5
			}
			if col.Type == ColliderSphere {
				col.Radius = col.SizeX
			} else if col.Type == ColliderCapsule {
				col.Radius = col.SizeX
				col.Height = col.SizeY
			}
		}
		m.Colliders = append(m.Colliders, col)
	}
	node := Node{
		Name:       gn.Name,
		Transform:  tr,
		MeshIndex:  meshIdx,
		Children:   make([]int, 0),
	}
	ourIdx := len(m.Nodes)
	m.Nodes = append(m.Nodes, node)
	for _, childIdx := range gn.Children {
		childStart := len(m.Nodes)
		importNodes(doc, nodes, childIdx, m)
		if len(m.Nodes) > childStart {
			m.Nodes[ourIdx].Children = append(m.Nodes[ourIdx].Children, childStart)
		}
	}
}

func nodeToTransform(gn *gltf.Node) Transform {
	tr := DefaultTransform()
	t := gn.TranslationOrDefault()
	s := gn.ScaleOrDefault()
	tr.X = float32(t[0])
	tr.Y = float32(t[1])
	tr.Z = float32(t[2])
	tr.ScaleX = float32(s[0])
	tr.ScaleY = float32(s[1])
	tr.ScaleZ = float32(s[2])
	if tr.ScaleX == 0 {
		tr.ScaleX = 1
	}
	if tr.ScaleY == 0 {
		tr.ScaleY = 1
	}
	if tr.ScaleZ == 0 {
		tr.ScaleZ = 1
	}
	r := gn.RotationOrDefault()
	w, x, y, z := float32(r[3]), float32(r[0]), float32(r[1]), float32(r[2])
	tr.Pitch, tr.Yaw, tr.Roll = quatToEuler(w, x, y, z)
	return tr
}

func quatToEuler(w, x, y, z float32) (pitch, yaw, roll float32) {
	// Convert quaternion to euler angles (degrees)
	sinP := -2 * (x*z - w*y)
	if sinP > 1 {
		sinP = 1
	}
	if sinP < -1 {
		sinP = -1
	}
	pitch = float32(math.Asin(float64(sinP))) * 180 / float32(math.Pi)
	yaw = float32(math.Atan2(float64(2*(w*z+x*y)), float64(1-2*(y*y+z*z)))) * 180 / float32(math.Pi)
	roll = float32(math.Atan2(float64(2*(w*x+y*z)), float64(1-2*(x*x+y*y)))) * 180 / float32(math.Pi)
	return pitch, yaw, roll
}
