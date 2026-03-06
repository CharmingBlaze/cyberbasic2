// Package dbp: Mesh and model queries - LoadMesh, GetModelBounds, GetMeshVertexCount, GetMeshTriangleCount.
package dbp

import (
	"fmt"

	"cyberbasic/compiler/vm"
	rl "github.com/gen2brain/raylib-go/raylib"
)

// register3DMesh adds LoadMesh, GetModelBounds, GetMeshVertexCount, GetMeshTriangleCount.
func register3DMesh(v *vm.VM) {
	v.RegisterForeign("LoadMesh", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("LoadMesh(id, path) requires 2 arguments")
		}
		id := toInt(args[0])
		path := toString(args[1])
		model := rl.LoadModel(path)
		objectsMu.Lock()
		objects[id] = newDbpObject(model)
		objectsMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("GetModelBounds", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return []interface{}{0.0, 0.0, 0.0, 0.0, 0.0, 0.0}, nil
		}
		id := toInt(args[0])
		objectsMu.Lock()
		obj, ok := objects[id]
		objectsMu.Unlock()
		if !ok || obj.model.MeshCount == 0 {
			return []interface{}{0.0, 0.0, 0.0, 0.0, 0.0, 0.0}, nil
		}
		meshes := obj.model.GetMeshes()
		if len(meshes) == 0 {
			return []interface{}{0.0, 0.0, 0.0, 0.0, 0.0, 0.0}, nil
		}
		box := rl.GetMeshBoundingBox(meshes[0])
		return []interface{}{box.Min.X, box.Min.Y, box.Min.Z, box.Max.X, box.Max.Y, box.Max.Z}, nil
	})
	v.RegisterForeign("GetMeshVertexCount", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return 0, nil
		}
		id := toInt(args[0])
		objectsMu.Lock()
		obj, ok := objects[id]
		objectsMu.Unlock()
		if !ok || obj.model.MeshCount == 0 {
			return 0, nil
		}
		meshes := obj.model.GetMeshes()
		if len(meshes) == 0 {
			return 0, nil
		}
		return meshes[0].VertexCount, nil
	})
	v.RegisterForeign("GetMeshTriangleCount", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return 0, nil
		}
		id := toInt(args[0])
		objectsMu.Lock()
		obj, ok := objects[id]
		objectsMu.Unlock()
		if !ok || obj.model.MeshCount == 0 {
			return 0, nil
		}
		meshes := obj.model.GetMeshes()
		if len(meshes) == 0 {
			return 0, nil
		}
		return meshes[0].TriangleCount, nil
	})
}
