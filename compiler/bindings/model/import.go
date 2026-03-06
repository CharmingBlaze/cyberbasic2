// Package model: format detection and unified Load entry point.
package model

import (
	"fmt"
	"path/filepath"
	"strings"
)

// Load parses the file and returns a canonical Model. Uses file extension to detect format.
func Load(path string) (*Model, error) {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".gltf", ".glb":
		return importGLTF(path)
	case ".obj":
		return importOBJ(path)
	case ".fbx":
		return importFBX(path)
	default:
		return nil, fmt.Errorf("unsupported format: %s (use .gltf, .glb, or .obj)", ext)
	}
}
