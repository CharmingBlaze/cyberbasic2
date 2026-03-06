// Package model: FBX importer (stub - go-fbx/fbx returns 404, use GLTF or OBJ).
package model

import "fmt"

// importFBX loads an FBX file. Not implemented - use GLTF or OBJ.
func importFBX(path string) (*Model, error) {
	return nil, fmt.Errorf("FBX import not yet supported: use .gltf or .obj (file: %s)", path)
}
