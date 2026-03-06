// Package resources provides central asset caching with reference counting.
package resources

import (
	"fmt"
	"sync"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
	modelsMu   sync.RWMutex
	models     = make(map[string]*modelEntry)
	texturesMu sync.RWMutex
	textures   = make(map[string]*textureEntry)
)

type modelEntry struct {
	model rl.Model
	refs  int
}

type textureEntry struct {
	tex  rl.Texture2D
	refs int
}

// LoadModel loads a model from path. Returns id for UnloadModel. Reference counted.
func LoadModel(path string) (string, error) {
	modelsMu.Lock()
	defer modelsMu.Unlock()
	if e, ok := models[path]; ok {
		e.refs++
		return path, nil
	}
	model := rl.LoadModel(path)
	if model.MeshCount == 0 {
		return "", fmt.Errorf("LoadModel: failed to load %s", path)
	}
	models[path] = &modelEntry{model: model, refs: 1}
	return path, nil
}

// UnloadModel decrements ref count; unloads when refs reach 0.
func UnloadModel(path string) {
	modelsMu.Lock()
	defer modelsMu.Unlock()
	e, ok := models[path]
	if !ok {
		return
	}
	e.refs--
	if e.refs <= 0 {
		rl.UnloadModel(e.model)
		delete(models, path)
	}
}

// GetModel returns the model for path, or zero if not loaded.
func GetModel(path string) rl.Model {
	modelsMu.RLock()
	defer modelsMu.RUnlock()
	if e, ok := models[path]; ok {
		return e.model
	}
	return rl.Model{}
}

// ModelExists returns true if the model is cached.
func ModelExists(path string) bool {
	modelsMu.RLock()
	defer modelsMu.RUnlock()
	_, ok := models[path]
	return ok
}

// LoadTexture loads a texture from path. Reference counted.
func LoadTexture(path string) (string, error) {
	texturesMu.Lock()
	defer texturesMu.Unlock()
	if e, ok := textures[path]; ok {
		e.refs++
		return path, nil
	}
	tex := rl.LoadTexture(path)
	if tex.ID == 0 {
		return "", fmt.Errorf("LoadTexture: failed to load %s", path)
	}
	textures[path] = &textureEntry{tex: tex, refs: 1}
	return path, nil
}

// UnloadTexture decrements ref count; unloads when refs reach 0.
func UnloadTexture(path string) {
	texturesMu.Lock()
	defer texturesMu.Unlock()
	e, ok := textures[path]
	if !ok {
		return
	}
	e.refs--
	if e.refs <= 0 {
		rl.UnloadTexture(e.tex)
		delete(textures, path)
	}
}

// GetTexture returns the texture for path, or zero if not loaded.
func GetTexture(path string) rl.Texture2D {
	texturesMu.RLock()
	defer texturesMu.RUnlock()
	if e, ok := textures[path]; ok {
		return e.tex
	}
	return rl.Texture2D{}
}

// DefaultTexture returns a 1x1 white texture for fallback.
func DefaultTexture() rl.Texture2D {
	img := rl.GenImageColor(1, 1, rl.White)
	tex := rl.LoadTextureFromImage(img)
	rl.UnloadImage(img)
	return tex
}
