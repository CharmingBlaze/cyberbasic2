// Package assets provides a unified asset cache for models and textures.
// Use LoadAsset/PreloadAsset to load by path; UnloadAsset when done.
// DBP LoadObject and LoadLevel can use this cache to avoid re-parsing.
package assets

import (
	"fmt"
	"path/filepath"
	"strings"
	"sync"

	"cyberbasic/compiler/bindings/model"
	"cyberbasic/compiler/runtime/resources"
)

var (
	modelCache   = make(map[string]*modelEntry)
	modelCacheMu sync.RWMutex
)

type modelEntry struct {
	m    *model.Model
	refs int
}

// LoadAsset loads a model or texture by path. Path is the cache key.
// For .gltf, .glb, .obj: parses and caches model.Model.
// For .png, .jpg, .bmp, etc.: uses resources.LoadTexture.
// Returns the path on success. Reference counted.
func LoadAsset(path string) (string, error) {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".gltf", ".glb", ".obj":
		return loadModelAsset(path)
	case ".png", ".jpg", ".jpeg", ".bmp", ".tga", ".gif":
		return resources.LoadTexture(path)
	default:
		return "", fmt.Errorf("LoadAsset: unsupported extension %s", ext)
	}
}

func loadModelAsset(path string) (string, error) {
	modelCacheMu.Lock()
	defer modelCacheMu.Unlock()
	if e, ok := modelCache[path]; ok {
		e.refs++
		return path, nil
	}
	m, err := model.Load(path)
	if err != nil {
		return "", fmt.Errorf("LoadAsset: %w", err)
	}
	modelCache[path] = &modelEntry{m: m, refs: 1}
	return path, nil
}

// UnloadAsset decrements ref count; unloads when refs reach 0.
func UnloadAsset(path string) {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".gltf", ".glb", ".obj":
		unloadModelAsset(path)
	case ".png", ".jpg", ".jpeg", ".bmp", ".tga", ".gif":
		resources.UnloadTexture(path)
	}
}

func unloadModelAsset(path string) {
	modelCacheMu.Lock()
	defer modelCacheMu.Unlock()
	e, ok := modelCache[path]
	if !ok {
		return
	}
	e.refs--
	if e.refs <= 0 {
		delete(modelCache, path)
	}
}

// AssetExists returns true if the asset is cached.
func AssetExists(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".gltf", ".glb", ".obj":
		modelCacheMu.RLock()
		_, ok := modelCache[path]
		modelCacheMu.RUnlock()
		return ok
	case ".png", ".jpg", ".jpeg", ".bmp", ".tga", ".gif":
		return resources.TextureExists(path)
	default:
		return false
	}
}

// GetModel returns the cached model for path, or nil if not loaded.
// Used by BuildModel when integrating with the asset pipeline.
func GetModel(path string) *model.Model {
	modelCacheMu.RLock()
	defer modelCacheMu.RUnlock()
	if e, ok := modelCache[path]; ok {
		return e.m
	}
	return nil
}

// PreloadAsset is equivalent to LoadAsset (sync preload).
// Call at startup to load assets before first use.
func PreloadAsset(path string) (string, error) {
	return LoadAsset(path)
}

// LoadModelForBuild loads or returns cached model. Used by DBP when cache integration is enabled.
func LoadModelForBuild(path string) (*model.Model, error) {
	modelCacheMu.Lock()
	defer modelCacheMu.Unlock()
	if e, ok := modelCache[path]; ok {
		e.refs++
		return e.m, nil
	}
	m, err := model.Load(path)
	if err != nil {
		return nil, err
	}
	modelCache[path] = &modelEntry{m: m, refs: 1}
	return m, nil
}

// UnloadModelForBuild decrements ref after BuildModel is done with the model.
func UnloadModelForBuild(path string) {
	unloadModelAsset(path)
}
