// Package model provides a canonical internal representation for 3D models.
// All formats (GLTF, OBJ, FBX) are converted into this struct.
package model

import "math"

// Model is the unified internal representation. Level and Prefab wrap this.
type Model struct {
	Meshes     []Mesh
	Materials  []Material
	Textures   []Texture
	Nodes      []Node
	Skeleton   *Skeleton // from GLTF skin; nil for static meshes
	Animations []Animation
	Lights     []Light
	Colliders  []Collider
}

// Mesh holds vertex data for one mesh.
type Mesh struct {
	Vertices []float32 // x,y,z per vertex
	Normals  []float32 // x,y,z per vertex
	Texcoords []float32 // u,v per vertex
	Indices   []uint32
	MaterialIndex int // -1 if none
}

// Material holds PBR material properties.
type Material struct {
	BaseColorR, BaseColorG, BaseColorB, BaseColorA float32
	Metallic                                       float32
	Roughness                                      float32
	BaseColorTextureIndex                          int // -1 if none
	NormalTextureIndex                             int
	MetallicRoughnessTextureIndex                  int
}

// Texture references an image file path (resolved relative to model).
type Texture struct {
	Path string
}

// Node is a scene graph node with optional mesh.
type Node struct {
	Name       string
	Transform  Transform
	MeshIndex  int   // -1 if no mesh
	Children   []int
}

// Transform is position, rotation (pitch/yaw/roll in degrees), scale.
type Transform struct {
	X, Y, Z     float32
	Pitch, Yaw, Roll float32
	ScaleX, ScaleY, ScaleZ float32
}

// Bone holds joint data for skeletal animation.
type Bone struct {
	Name        string
	Parent      int   // -1 for root
	BindPose    [16]float32 // Matrix4 column-major
	InverseBind [16]float32
}

// Skeleton holds the bone hierarchy for skinned meshes.
type Skeleton struct {
	Bones []Bone
}

// SkinnedModel wraps Model with skeleton and animations (for reference).
type SkinnedModel struct {
	Model      *Model
	Skeleton   *Skeleton
	Animations []Animation
}

// Animation holds animation clip data (keyframes, etc.).
type Animation struct {
	Name     string
	Duration float32
	Channels []AnimationChannel
}

// AnimationChannel targets a node/bone and property.
type AnimationChannel struct {
	NodeIndex int
	BoneIndex int   // -1 if node animation; for skeletal
	Property  string // "translation", "rotation", "scale"
	Keyframes []Keyframe
}

// Keyframe is one keyframe (time, value).
type Keyframe struct {
	Time  float32
	Value []float32 // 3 for translation/scale, 4 for rotation quat
}

// Light types.
const (
	LightPoint       = 0
	LightDirectional = 1
	LightSpot        = 2
)

// Light holds light data.
type Light struct {
	Type       int
	X, Y, Z    float32
	DirX, DirY, DirZ float32
	R, G, B    float32
	Intensity  float32
	Range      float32
	InnerCone  float32
	OuterCone  float32
}

// Collider types.
const (
	ColliderMesh    = 0
	ColliderBox     = 1
	ColliderSphere  = 2
	ColliderCapsule = 3
)

// Collider holds collision shape data.
type Collider struct {
	Type       int
	Transform  Transform
	MeshIndex  int   // for mesh colliders
	SizeX, SizeY, SizeZ float32 // for box
	Radius     float32 // for sphere/capsule
	Height     float32 // for capsule
}

// DefaultTransform returns identity transform.
func DefaultTransform() Transform {
	return Transform{
		ScaleX: 1, ScaleY: 1, ScaleZ: 1,
	}
}

// Radians converts degrees to radians.
func Radians(deg float32) float32 {
	return deg * float32(math.Pi) / 180
}

// MeshBounds returns axis-aligned bounds (minX, minY, minZ, maxX, maxY, maxZ) for a mesh.
// Returns zeros if mesh is empty.
func MeshBounds(m *Mesh) (minX, minY, minZ, maxX, maxY, maxZ float32) {
	if len(m.Vertices) < 3 {
		return 0, 0, 0, 0, 0, 0
	}
	minX, minY, minZ = m.Vertices[0], m.Vertices[1], m.Vertices[2]
	maxX, maxY, maxZ = minX, minY, minZ
	for i := 3; i < len(m.Vertices); i += 3 {
		x, y, z := m.Vertices[i], m.Vertices[i+1], m.Vertices[i+2]
		if x < minX {
			minX = x
		}
		if x > maxX {
			maxX = x
		}
		if y < minY {
			minY = y
		}
		if y > maxY {
			maxY = y
		}
		if z < minZ {
			minZ = z
		}
		if z > maxZ {
			maxZ = z
		}
	}
	return minX, minY, minZ, maxX, maxY, maxZ
}

// ComputeFlatNormals fills mesh.Normals from triangle faces using cross product.
// Call when normals are missing or invalid. Overwrites existing normals.
func ComputeFlatNormals(m *Mesh) {
	vCount := len(m.Vertices) / 3
	if vCount == 0 {
		return
	}
	normals := make([]float32, vCount*3)
	if len(m.Indices) > 0 {
		for i := 0; i+2 < len(m.Indices); i += 3 {
			i0, i1, i2 := int(m.Indices[i]), int(m.Indices[i+1]), int(m.Indices[i+2])
			if i0 >= vCount || i1 >= vCount || i2 >= vCount {
				continue
			}
			v0 := [3]float32{m.Vertices[i0*3], m.Vertices[i0*3+1], m.Vertices[i0*3+2]}
			v1 := [3]float32{m.Vertices[i1*3], m.Vertices[i1*3+1], m.Vertices[i1*3+2]}
			v2 := [3]float32{m.Vertices[i2*3], m.Vertices[i2*3+1], m.Vertices[i2*3+2]}
			e1 := [3]float32{v1[0] - v0[0], v1[1] - v0[1], v1[2] - v0[2]}
			e2 := [3]float32{v2[0] - v0[0], v2[1] - v0[1], v2[2] - v0[2]}
			nx := e1[1]*e2[2] - e1[2]*e2[1]
			ny := e1[2]*e2[0] - e1[0]*e2[2]
			nz := e1[0]*e2[1] - e1[1]*e2[0]
			len2 := nx*nx + ny*ny + nz*nz
			if len2 > 1e-12 {
				inv := float32(1.0 / math.Sqrt(float64(len2)))
				nx, ny, nz = nx*inv, ny*inv, nz*inv
			} else {
				nx, ny, nz = 0, 1, 0
			}
			normals[i0*3], normals[i0*3+1], normals[i0*3+2] = nx, ny, nz
			normals[i1*3], normals[i1*3+1], normals[i1*3+2] = nx, ny, nz
			normals[i2*3], normals[i2*3+1], normals[i2*3+2] = nx, ny, nz
		}
	} else {
		for i := 0; i+2 < vCount; i += 3 {
			v0 := [3]float32{m.Vertices[i*3], m.Vertices[i*3+1], m.Vertices[i*3+2]}
			v1 := [3]float32{m.Vertices[(i+1)*3], m.Vertices[(i+1)*3+1], m.Vertices[(i+1)*3+2]}
			v2 := [3]float32{m.Vertices[(i+2)*3], m.Vertices[(i+2)*3+1], m.Vertices[(i+2)*3+2]}
			e1 := [3]float32{v1[0] - v0[0], v1[1] - v0[1], v1[2] - v0[2]}
			e2 := [3]float32{v2[0] - v0[0], v2[1] - v0[1], v2[2] - v0[2]}
			nx := e1[1]*e2[2] - e1[2]*e2[1]
			ny := e1[2]*e2[0] - e1[0]*e2[2]
			nz := e1[0]*e2[1] - e1[1]*e2[0]
			len2 := nx*nx + ny*ny + nz*nz
			if len2 > 1e-12 {
				inv := float32(1.0 / math.Sqrt(float64(len2)))
				nx, ny, nz = nx*inv, ny*inv, nz*inv
			} else {
				nx, ny, nz = 0, 1, 0
			}
			for j := 0; j < 3; j++ {
				normals[(i+j)*3], normals[(i+j)*3+1], normals[(i+j)*3+2] = nx, ny, nz
			}
		}
	}
	m.Normals = normals
}
