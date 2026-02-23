// Package raylib: mesh generation, upload, draw; materials; model animations.
package raylib

import (
	"cyberbasic/compiler/vm"
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func registerMesh(v *vm.VM) {
	// Mesh generation - each returns meshId
	v.RegisterForeign("GenMeshPoly", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("GenMeshPoly requires (sides, radius)")
		}
		sides := toInt32(args[0])
		radius := toFloat32(args[1])
		mesh := rl.GenMeshPoly(int(sides), radius)
		meshMu.Lock()
		meshCounter++
		id := fmt.Sprintf("mesh_%d", meshCounter)
		meshes[id] = mesh
		meshMu.Unlock()
		return id, nil
	})
	v.RegisterForeign("GenMeshPlane", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("GenMeshPlane requires (width, length, resX, resZ)")
		}
		width := toFloat32(args[0])
		length := toFloat32(args[1])
		resX, resZ := int32(1), int32(1)
		if len(args) >= 4 {
			resX, resZ = toInt32(args[2]), toInt32(args[3])
		}
		mesh := rl.GenMeshPlane(width, length, int(resX), int(resZ))
		meshMu.Lock()
		meshCounter++
		id := fmt.Sprintf("mesh_%d", meshCounter)
		meshes[id] = mesh
		meshMu.Unlock()
		return id, nil
	})
	v.RegisterForeign("GenMeshCube", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("GenMeshCube requires (width, height, length)")
		}
		w, h, l := toFloat32(args[0]), toFloat32(args[1]), toFloat32(args[2])
		mesh := rl.GenMeshCube(w, h, l)
		meshMu.Lock()
		meshCounter++
		id := fmt.Sprintf("mesh_%d", meshCounter)
		meshes[id] = mesh
		meshMu.Unlock()
		return id, nil
	})
	v.RegisterForeign("GenMeshSphere", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("GenMeshSphere requires (radius, rings, slices)")
		}
		radius := toFloat32(args[0])
		rings, slices := toInt32(args[1]), toInt32(args[2])
		mesh := rl.GenMeshSphere(radius, int(rings), int(slices))
		meshMu.Lock()
		meshCounter++
		id := fmt.Sprintf("mesh_%d", meshCounter)
		meshes[id] = mesh
		meshMu.Unlock()
		return id, nil
	})
	v.RegisterForeign("GenMeshHemiSphere", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("GenMeshHemiSphere requires (radius, rings, slices)")
		}
		radius := toFloat32(args[0])
		rings, slices := toInt32(args[1]), toInt32(args[2])
		mesh := rl.GenMeshHemiSphere(radius, int(rings), int(slices))
		meshMu.Lock()
		meshCounter++
		id := fmt.Sprintf("mesh_%d", meshCounter)
		meshes[id] = mesh
		meshMu.Unlock()
		return id, nil
	})
	v.RegisterForeign("GenMeshCylinder", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("GenMeshCylinder requires (radius, height, slices)")
		}
		radius := toFloat32(args[0])
		height := toFloat32(args[1])
		slices := toInt32(args[2])
		mesh := rl.GenMeshCylinder(radius, height, int(slices))
		meshMu.Lock()
		meshCounter++
		id := fmt.Sprintf("mesh_%d", meshCounter)
		meshes[id] = mesh
		meshMu.Unlock()
		return id, nil
	})
	v.RegisterForeign("GenMeshCone", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("GenMeshCone requires (radius, height, slices)")
		}
		radius := toFloat32(args[0])
		height := toFloat32(args[1])
		slices := toInt32(args[2])
		mesh := rl.GenMeshCone(radius, height, int(slices))
		meshMu.Lock()
		meshCounter++
		id := fmt.Sprintf("mesh_%d", meshCounter)
		meshes[id] = mesh
		meshMu.Unlock()
		return id, nil
	})
	v.RegisterForeign("GenMeshTorus", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("GenMeshTorus requires (radius, size, radSeg, sides)")
		}
		radius := toFloat32(args[0])
		size := toFloat32(args[1])
		radSeg, sides := toInt32(args[2]), toInt32(args[3])
		mesh := rl.GenMeshTorus(radius, size, int(radSeg), int(sides))
		meshMu.Lock()
		meshCounter++
		id := fmt.Sprintf("mesh_%d", meshCounter)
		meshes[id] = mesh
		meshMu.Unlock()
		return id, nil
	})
	v.RegisterForeign("GenMeshKnot", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("GenMeshKnot requires (radius, size, radSeg, sides)")
		}
		radius := toFloat32(args[0])
		size := toFloat32(args[1])
		radSeg, sides := toInt32(args[2]), toInt32(args[3])
		mesh := rl.GenMeshKnot(radius, size, int(radSeg), int(sides))
		meshMu.Lock()
		meshCounter++
		id := fmt.Sprintf("mesh_%d", meshCounter)
		meshes[id] = mesh
		meshMu.Unlock()
		return id, nil
	})
	v.RegisterForeign("GenMeshHeightmap", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("GenMeshHeightmap requires (fileName, sizeX, sizeY, sizeZ)")
		}
		path := toString(args[0])
		size := rl.Vector3{X: toFloat32(args[1]), Y: toFloat32(args[2]), Z: toFloat32(args[3])}
		img := rl.LoadImage(path)
		mesh := rl.GenMeshHeightmap(*img, size)
		rl.UnloadImage(img)
		meshMu.Lock()
		meshCounter++
		id := fmt.Sprintf("mesh_%d", meshCounter)
		meshes[id] = mesh
		meshMu.Unlock()
		return id, nil
	})
	v.RegisterForeign("GenMeshCubicmap", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("GenMeshCubicmap requires (fileName, cubeSizeX, cubeSizeY, cubeSizeZ)")
		}
		path := toString(args[0])
		cubeSize := rl.Vector3{X: toFloat32(args[1]), Y: toFloat32(args[2]), Z: toFloat32(args[3])}
		img := rl.LoadImage(path)
		mesh := rl.GenMeshCubicmap(*img, cubeSize)
		rl.UnloadImage(img)
		meshMu.Lock()
		meshCounter++
		id := fmt.Sprintf("mesh_%d", meshCounter)
		meshes[id] = mesh
		meshMu.Unlock()
		return id, nil
	})

	v.RegisterForeign("UploadMesh", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("UploadMesh requires (meshId, dynamic)")
		}
		meshId := toString(args[0])
		dynamic := len(args) > 1 && toFloat32(args[1]) != 0
		meshMu.Lock()
		mesh, ok := meshes[meshId]
		meshMu.Unlock()
		if !ok {
			return nil, fmt.Errorf("unknown mesh id: %s", meshId)
		}
		rl.UploadMesh(&mesh, dynamic)
		meshMu.Lock()
		meshes[meshId] = mesh
		meshMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("UnloadMesh", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("UnloadMesh requires (meshId)")
		}
		meshId := toString(args[0])
		meshMu.Lock()
		mesh, ok := meshes[meshId]
		delete(meshes, meshId)
		meshMu.Unlock()
		if ok {
			rl.UnloadMesh(&mesh)
		}
		return nil, nil
	})
	v.RegisterForeign("GetMeshBoundingBox", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("GetMeshBoundingBox requires (meshId)")
		}
		meshId := toString(args[0])
		meshMu.Lock()
		mesh, ok := meshes[meshId]
		meshMu.Unlock()
		if !ok {
			return nil, fmt.Errorf("unknown mesh id: %s", meshId)
		}
		box := rl.GetMeshBoundingBox(mesh)
		return []interface{}{box.Min.X, box.Min.Y, box.Min.Z, box.Max.X, box.Max.Y, box.Max.Z}, nil
	})
	v.RegisterForeign("ExportMesh", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("ExportMesh requires (meshId, fileName)")
		}
		meshId := toString(args[0])
		path := toString(args[1])
		meshMu.Lock()
		mesh, ok := meshes[meshId]
		meshMu.Unlock()
		if !ok {
			return nil, fmt.Errorf("unknown mesh id: %s", meshId)
		}
		rl.ExportMesh(mesh, path)
		return nil, nil
	})
	v.RegisterForeign("DrawMesh", func(args []interface{}) (interface{}, error) {
		if len(args) < 7 {
			return nil, fmt.Errorf("DrawMesh requires (meshId, materialId, posX,posY,posZ, scaleX,scaleY,scaleZ)")
		}
		meshId := toString(args[0])
		matId := toString(args[1])
		pos := rl.Vector3{X: toFloat32(args[2]), Y: toFloat32(args[3]), Z: toFloat32(args[4])}
		scale := rl.Vector3{X: toFloat32(args[5]), Y: toFloat32(args[6]), Z: toFloat32(args[7])}
		meshMu.Lock()
		mesh, ok := meshes[meshId]
		meshMu.Unlock()
		if !ok {
			return nil, fmt.Errorf("unknown mesh id: %s", meshId)
		}
		materialMu.Lock()
		mat, okMat := materials[matId]
		materialMu.Unlock()
		if !okMat {
			return nil, fmt.Errorf("unknown material id: %s", matId)
		}
		transform := rl.MatrixMultiply(rl.MatrixScale(scale.X, scale.Y, scale.Z), rl.MatrixTranslate(pos.X, pos.Y, pos.Z))
		rl.DrawMesh(mesh, mat, transform)
		return nil, nil
	})

	v.RegisterForeign("UpdateMeshBuffer", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("UpdateMeshBuffer requires (meshId, bufferIndex, data, offset)")
		}
		meshId := toString(args[0])
		bufferIndex := toInt32(args[1])
		offset := toInt32(args[3])
		var data []byte
		switch v := args[2].(type) {
		case string:
			data = []byte(v)
		case []byte:
			data = v
		default:
			return nil, fmt.Errorf("UpdateMeshBuffer data must be string or bytes")
		}
		meshMu.Lock()
		mesh, ok := meshes[meshId]
		meshMu.Unlock()
		if !ok {
			return nil, fmt.Errorf("unknown mesh id: %s", meshId)
		}
		rl.UpdateMeshBuffer(mesh, int(bufferIndex), data, int(offset))
		meshMu.Lock()
		meshes[meshId] = mesh
		meshMu.Unlock()
		return nil, nil
	})

	v.RegisterForeign("DrawMeshInstanced", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("DrawMeshInstanced requires (meshId, materialId, instanceCount, ...16*count matrix floats)")
		}
		meshId := toString(args[0])
		matId := toString(args[1])
		instanceCount := toInt32(args[2])
		if instanceCount <= 0 {
			return nil, nil
		}
		required := 3 + int(instanceCount)*16
		if len(args) < required {
			return nil, fmt.Errorf("DrawMeshInstanced needs %d args (meshId, materialId, count + %d floats)", required, instanceCount*16)
		}
		meshMu.Lock()
		mesh, ok := meshes[meshId]
		meshMu.Unlock()
		if !ok {
			return nil, fmt.Errorf("unknown mesh id: %s", meshId)
		}
		materialMu.Lock()
		mat, okMat := materials[matId]
		materialMu.Unlock()
		if !okMat {
			return nil, fmt.Errorf("unknown material id: %s", matId)
		}
		transforms := make([]rl.Matrix, instanceCount)
		for i := int32(0); i < instanceCount; i++ {
			base := 3 + int(i)*16
			transforms[i] = rl.NewMatrix(
				toFloat32(args[base+0]), toFloat32(args[base+4]), toFloat32(args[base+8]), toFloat32(args[base+12]),
				toFloat32(args[base+1]), toFloat32(args[base+5]), toFloat32(args[base+9]), toFloat32(args[base+13]),
				toFloat32(args[base+2]), toFloat32(args[base+6]), toFloat32(args[base+10]), toFloat32(args[base+14]),
				toFloat32(args[base+3]), toFloat32(args[base+7]), toFloat32(args[base+11]), toFloat32(args[base+15]),
			)
		}
		rl.DrawMeshInstanced(mesh, mat, transforms, int(instanceCount))
		return nil, nil
	})

	v.RegisterForeign("LoadMaterialDefault", func(args []interface{}) (interface{}, error) {
		mat := rl.LoadMaterialDefault()
		materialMu.Lock()
		materialCounter++
		id := fmt.Sprintf("mat_%d", materialCounter)
		materials[id] = mat
		materialMu.Unlock()
		return id, nil
	})
	v.RegisterForeign("IsMaterialValid", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return false, nil
		}
		materialMu.Lock()
		_, ok := materials[toString(args[0])]
		materialMu.Unlock()
		return ok, nil
	})
	v.RegisterForeign("UnloadMaterial", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("UnloadMaterial requires (materialId)")
		}
		id := toString(args[0])
		materialMu.Lock()
		mat, ok := materials[id]
		delete(materials, id)
		materialMu.Unlock()
		if ok {
			rl.UnloadMaterial(mat)
		}
		return nil, nil
	})
	v.RegisterForeign("SetMaterialTexture", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("SetMaterialTexture requires (materialId, mapType, textureId)")
		}
		matId := toString(args[0])
		mapType := toInt32(args[1])
		texId := toString(args[2])
		materialMu.Lock()
		mat, ok := materials[matId]
		materialMu.Unlock()
		if !ok {
			return nil, fmt.Errorf("unknown material id: %s", matId)
		}
		texMu.Lock()
		tex, okTex := textures[texId]
		texMu.Unlock()
		if !okTex {
			return nil, fmt.Errorf("unknown texture id: %s", texId)
		}
		rl.SetMaterialTexture(&mat, mapType, tex)
		materialMu.Lock()
		materials[matId] = mat
		materialMu.Unlock()
		return nil, nil
	})

	v.RegisterForeign("LoadMaterials", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("LoadMaterials requires (fileName)")
		}
		path := toString(args[0])
		materialMu.Lock()
		for i := 0; i < lastLoadMaterialsCount; i++ {
			id := fmt.Sprintf("mat_loaded_%d", i)
			if mat, ok := materials[id]; ok {
				rl.UnloadMaterial(mat)
				delete(materials, id)
			}
		}
		materialMu.Unlock()
		mats := rl.LoadMaterials(path)
		materialMu.Lock()
		lastLoadMaterialsCount = len(mats)
		for i, m := range mats {
			materials[fmt.Sprintf("mat_loaded_%d", i)] = m
		}
		materialMu.Unlock()
		return int32(len(mats)), nil
	})

	v.RegisterForeign("GetMaterialIdFromLoad", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return "", nil
		}
		index := int(toInt32(args[0]))
		if index < 0 || index >= lastLoadMaterialsCount {
			return "", nil
		}
		return fmt.Sprintf("mat_loaded_%d", index), nil
	})

	// GetRayCollisionMesh: ray (6), meshId, transform (pos 3 + scale 3)
	v.RegisterForeign("GetRayCollisionMesh", func(args []interface{}) (interface{}, error) {
		if len(args) < 13 {
			return nil, fmt.Errorf("GetRayCollisionMesh requires (rayPosX,Y,Z, rayDirX,Y,Z, meshId, posX,posY,posZ, scaleX,scaleY,scaleZ)")
		}
		ray := rl.Ray{
			Position:  rl.Vector3{X: toFloat32(args[0]), Y: toFloat32(args[1]), Z: toFloat32(args[2])},
			Direction: rl.Vector3{X: toFloat32(args[3]), Y: toFloat32(args[4]), Z: toFloat32(args[5])},
		}
		meshId := toString(args[6])
		pos := rl.Vector3{X: toFloat32(args[7]), Y: toFloat32(args[8]), Z: toFloat32(args[9])}
		scale := rl.Vector3{X: toFloat32(args[10]), Y: toFloat32(args[11]), Z: toFloat32(args[12])}
		meshMu.Lock()
		mesh, ok := meshes[meshId]
		meshMu.Unlock()
		if !ok {
			return nil, fmt.Errorf("unknown mesh id: %s", meshId)
		}
		transform := rl.MatrixMultiply(rl.MatrixScale(scale.X, scale.Y, scale.Z), rl.MatrixTranslate(pos.X, pos.Y, pos.Z))
		coll := rl.GetRayCollisionMesh(ray, mesh, transform)
		lastRayCollisionMu.Lock()
		lastRayCollision = coll
		lastRayCollisionMu.Unlock()
		if coll.Hit {
			return 1, nil
		}
		return 0, nil
	})
}
