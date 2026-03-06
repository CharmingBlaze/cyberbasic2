// Package materials provides unified material format for PBR and simple materials.
package materials

import (
	"sync"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Material holds color, textures, and shader refs.
type Material struct {
	Color     rl.Color
	Texture   rl.Texture2D
	NormalMap rl.Texture2D
	Roughness float32
	Metallic  float32
	Emission  rl.Color
	ShaderID  int
}

var (
	materials   = make(map[int]*Material)
	materialsMu sync.RWMutex
)

// Set creates or updates a material by id.
func Set(id int, m *Material) {
	materialsMu.Lock()
	defer materialsMu.Unlock()
	materials[id] = m
}

// Get returns the material for id.
func Get(id int) *Material {
	materialsMu.RLock()
	defer materialsMu.RUnlock()
	return materials[id]
}

// Exists returns true if material id exists.
func Exists(id int) bool {
	materialsMu.RLock()
	defer materialsMu.RUnlock()
	_, ok := materials[id]
	return ok
}
