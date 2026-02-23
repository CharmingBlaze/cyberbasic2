// Package raylib: 3D camera, models, and primitives (rmodels).
package raylib

import (
	"cyberbasic/compiler/vm"
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func register3D(v *vm.VM) {
	v.RegisterForeign("SetCamera3D", func(args []interface{}) (interface{}, error) {
		if len(args) < 9 {
			return nil, fmt.Errorf("SetCamera3D requires (posX, posY, posZ, targetX, targetY, targetZ, upX, upY, upZ)")
		}
		camera3D.Position = rl.Vector3{X: toFloat32(args[0]), Y: toFloat32(args[1]), Z: toFloat32(args[2])}
		camera3D.Target = rl.Vector3{X: toFloat32(args[3]), Y: toFloat32(args[4]), Z: toFloat32(args[5])}
		camera3D.Up = rl.Vector3{X: toFloat32(args[6]), Y: toFloat32(args[7]), Z: toFloat32(args[8])}
		camera3D.Fovy = 60.0
		camera3D.Projection = rl.CameraPerspective
		return nil, nil
	})
	v.RegisterForeign("BeginMode3D", func(args []interface{}) (interface{}, error) {
		rl.BeginMode3D(camera3D)
		return nil, nil
	})
	v.RegisterForeign("EndMode3D", func(args []interface{}) (interface{}, error) {
		rl.EndMode3D()
		return nil, nil
	})
	v.RegisterForeign("LoadModel", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("LoadModel requires (fileName)")
		}
		path := toString(args[0])
		model := rl.LoadModel(path)
		modelMu.Lock()
		modelCounter++
		id := fmt.Sprintf("model_%d", modelCounter)
		models[id] = model
		modelMu.Unlock()
		return id, nil
	})
	v.RegisterForeign("LoadModelFromMesh", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("LoadModelFromMesh requires (meshId)")
		}
		meshId := toString(args[0])
		meshMu.Lock()
		mesh, ok := meshes[meshId]
		meshMu.Unlock()
		if !ok {
			return nil, fmt.Errorf("unknown mesh id: %s", meshId)
		}
		model := rl.LoadModelFromMesh(mesh)
		modelMu.Lock()
		modelCounter++
		id := fmt.Sprintf("model_%d", modelCounter)
		models[id] = model
		modelMu.Unlock()
		return id, nil
	})
	v.RegisterForeign("UnloadModel", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("UnloadModel requires (id)")
		}
		id := toString(args[0])
		modelMu.Lock()
		model, ok := models[id]
		delete(models, id)
		modelMu.Unlock()
		if ok {
			rl.UnloadModel(model)
		}
		return nil, nil
	})
	v.RegisterForeign("DrawModel", func(args []interface{}) (interface{}, error) {
		if len(args) < 5 {
			return nil, fmt.Errorf("DrawModel requires (id, posX, posY, posZ, scale) and optional tint")
		}
		id := toString(args[0])
		modelMu.Lock()
		model, ok := models[id]
		modelMu.Unlock()
		if !ok {
			return nil, fmt.Errorf("unknown model id: %s", id)
		}
		pos := rl.Vector3{X: toFloat32(args[1]), Y: toFloat32(args[2]), Z: toFloat32(args[3])}
		scale := toFloat32(args[4])
		c := rl.White
		if len(args) >= 9 {
			c = argsToColor(args, 5)
		}
		rl.DrawModel(model, pos, scale, c)
		return nil, nil
	})
	v.RegisterForeign("DrawCube", func(args []interface{}) (interface{}, error) {
		if len(args) < 7 {
			return nil, fmt.Errorf("DrawCube requires (posX, posY, posZ, width, height, length, color)")
		}
		pos := rl.Vector3{X: toFloat32(args[0]), Y: toFloat32(args[1]), Z: toFloat32(args[2])}
		w, h, l := toFloat32(args[3]), toFloat32(args[4]), toFloat32(args[5])
		c := rl.White
		if len(args) >= 11 {
			c = argsToColor(args, 6)
		}
		rl.DrawCube(pos, w, h, l, c)
		return nil, nil
	})
	v.RegisterForeign("DrawCubeWires", func(args []interface{}) (interface{}, error) {
		if len(args) < 7 {
			return nil, fmt.Errorf("DrawCubeWires requires (posX, posY, posZ, width, height, length, color)")
		}
		pos := rl.Vector3{X: toFloat32(args[0]), Y: toFloat32(args[1]), Z: toFloat32(args[2])}
		w, h, l := toFloat32(args[3]), toFloat32(args[4]), toFloat32(args[5])
		c := rl.White
		if len(args) >= 11 {
			c = argsToColor(args, 6)
		}
		rl.DrawCubeWires(pos, w, h, l, c)
		return nil, nil
	})
	v.RegisterForeign("DrawSphere", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("DrawSphere requires (posX, posY, posZ, radius, color)")
		}
		pos := rl.Vector3{X: toFloat32(args[0]), Y: toFloat32(args[1]), Z: toFloat32(args[2])}
		radius := toFloat32(args[3])
		c := rl.White
		if len(args) >= 8 {
			c = argsToColor(args, 4)
		}
		rl.DrawSphere(pos, radius, c)
		return nil, nil
	})
	v.RegisterForeign("DrawSphereWires", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("DrawSphereWires requires (posX, posY, posZ, radius, color)")
		}
		pos := rl.Vector3{X: toFloat32(args[0]), Y: toFloat32(args[1]), Z: toFloat32(args[2])}
		radius := toFloat32(args[3])
		rings, slices := int32(16), int32(16)
		if len(args) >= 6 {
			rings, slices = toInt32(args[4]), toInt32(args[5])
		}
		c := rl.White
		if len(args) >= 10 {
			c = argsToColor(args, 6)
		} else if len(args) >= 9 {
			c = argsToColor(args, 5)
		}
		rl.DrawSphereWires(pos, radius, rings, slices, c)
		return nil, nil
	})
	v.RegisterForeign("DrawPlane", func(args []interface{}) (interface{}, error) {
		if len(args) < 6 {
			return nil, fmt.Errorf("DrawPlane requires (centerX, centerY, centerZ, sizeX, sizeZ, color)")
		}
		center := rl.Vector3{X: toFloat32(args[0]), Y: toFloat32(args[1]), Z: toFloat32(args[2])}
		sizeX, sizeZ := toFloat32(args[3]), toFloat32(args[4])
		c := rl.White
		if len(args) >= 10 {
			c = argsToColor(args, 5)
		}
		rl.DrawPlane(center, rl.Vector2{X: sizeX, Y: sizeZ}, c)
		return nil, nil
	})
	v.RegisterForeign("DrawLine3D", func(args []interface{}) (interface{}, error) {
		if len(args) < 7 {
			return nil, fmt.Errorf("DrawLine3D requires (startX,startY,startZ, endX,endY,endZ, color)")
		}
		start := rl.Vector3{X: toFloat32(args[0]), Y: toFloat32(args[1]), Z: toFloat32(args[2])}
		end := rl.Vector3{X: toFloat32(args[3]), Y: toFloat32(args[4]), Z: toFloat32(args[5])}
		c := rl.White
		if len(args) >= 11 {
			c = argsToColor(args, 6)
		}
		rl.DrawLine3D(start, end, c)
		return nil, nil
	})
	v.RegisterForeign("DrawPoint3D", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("DrawPoint3D requires (x, y, z, color)")
		}
		pos := rl.Vector3{X: toFloat32(args[0]), Y: toFloat32(args[1]), Z: toFloat32(args[2])}
		c := rl.White
		if len(args) >= 8 {
			c = argsToColor(args, 3)
		}
		rl.DrawPoint3D(pos, c)
		return nil, nil
	})
	v.RegisterForeign("DrawCircle3D", func(args []interface{}) (interface{}, error) {
		if len(args) < 8 {
			return nil, fmt.Errorf("DrawCircle3D requires (centerX,Y,Z, radius, rotAxisX,Y,Z, rotAngle, color)")
		}
		center := rl.Vector3{X: toFloat32(args[0]), Y: toFloat32(args[1]), Z: toFloat32(args[2])}
		radius := toFloat32(args[3])
		rotAxis := rl.Vector3{X: toFloat32(args[4]), Y: toFloat32(args[5]), Z: toFloat32(args[6])}
		rotAngle := toFloat32(args[7])
		c := rl.White
		if len(args) >= 12 {
			c = argsToColor(args, 8)
		}
		rl.DrawCircle3D(center, radius, rotAxis, rotAngle, c)
		return nil, nil
	})
	v.RegisterForeign("DrawCubeV", func(args []interface{}) (interface{}, error) {
		if len(args) < 7 {
			return nil, fmt.Errorf("DrawCubeV requires (posX,posY,posZ, sizeX,sizeY,sizeZ, color)")
		}
		pos := rl.Vector3{X: toFloat32(args[0]), Y: toFloat32(args[1]), Z: toFloat32(args[2])}
		size := rl.Vector3{X: toFloat32(args[3]), Y: toFloat32(args[4]), Z: toFloat32(args[5])}
		c := rl.White
		if len(args) >= 11 {
			c = argsToColor(args, 6)
		}
		rl.DrawCubeV(pos, size, c)
		return nil, nil
	})
	v.RegisterForeign("DrawCylinder", func(args []interface{}) (interface{}, error) {
		if len(args) < 7 {
			return nil, fmt.Errorf("DrawCylinder requires (posX,posY,posZ, radiusTop, radiusBottom, height, slices, color)")
		}
		pos := rl.Vector3{X: toFloat32(args[0]), Y: toFloat32(args[1]), Z: toFloat32(args[2])}
		radiusTop := toFloat32(args[3])
		radiusBottom := toFloat32(args[4])
		height := toFloat32(args[5])
		slices := toInt32(args[6])
		c := rl.White
		if len(args) >= 11 {
			c = argsToColor(args, 7)
		}
		rl.DrawCylinder(pos, radiusTop, radiusBottom, height, slices, c)
		return nil, nil
	})
	v.RegisterForeign("DrawCylinderWires", func(args []interface{}) (interface{}, error) {
		if len(args) < 7 {
			return nil, fmt.Errorf("DrawCylinderWires requires (posX,posY,posZ, radiusTop, radiusBottom, height, slices, color)")
		}
		pos := rl.Vector3{X: toFloat32(args[0]), Y: toFloat32(args[1]), Z: toFloat32(args[2])}
		radiusTop := toFloat32(args[3])
		radiusBottom := toFloat32(args[4])
		height := toFloat32(args[5])
		slices := toInt32(args[6])
		c := rl.White
		if len(args) >= 11 {
			c = argsToColor(args, 7)
		}
		rl.DrawCylinderWires(pos, radiusTop, radiusBottom, height, slices, c)
		return nil, nil
	})
	v.RegisterForeign("DrawRay", func(args []interface{}) (interface{}, error) {
		if len(args) < 7 {
			return nil, fmt.Errorf("DrawRay requires (posX,posY,posZ, dirX,dirY,dirZ, color)")
		}
		ray := rl.Ray{
			Position:  rl.Vector3{X: toFloat32(args[0]), Y: toFloat32(args[1]), Z: toFloat32(args[2])},
			Direction: rl.Vector3{X: toFloat32(args[3]), Y: toFloat32(args[4]), Z: toFloat32(args[5])},
		}
		c := rl.White
		if len(args) >= 11 {
			c = argsToColor(args, 6)
		}
		rl.DrawRay(ray, c)
		return nil, nil
	})
	v.RegisterForeign("DrawTriangle3D", func(args []interface{}) (interface{}, error) {
		if len(args) < 10 {
			return nil, fmt.Errorf("DrawTriangle3D requires (v1x,v1y,v1z, v2x,v2y,v2z, v3x,v3y,v3z, color)")
		}
		v1 := rl.Vector3{X: toFloat32(args[0]), Y: toFloat32(args[1]), Z: toFloat32(args[2])}
		v2 := rl.Vector3{X: toFloat32(args[3]), Y: toFloat32(args[4]), Z: toFloat32(args[5])}
		v3 := rl.Vector3{X: toFloat32(args[6]), Y: toFloat32(args[7]), Z: toFloat32(args[8])}
		c := rl.White
		if len(args) >= 14 {
			c = argsToColor(args, 9)
		}
		rl.DrawTriangle3D(v1, v2, v3, c)
		return nil, nil
	})
	v.RegisterForeign("DrawTriangleStrip3D", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("DrawTriangleStrip3D requires (pointCount, x1,y1,z1, x2,y2,z2, ..., color)")
		}
		pointCount := int(toInt32(args[0]))
		if pointCount < 3 {
			return nil, fmt.Errorf("DrawTriangleStrip3D pointCount must be >= 3")
		}
		colorOffset := 1 + pointCount*3
		if len(args) < colorOffset+4 {
			return nil, fmt.Errorf("DrawTriangleStrip3D needs %d coords + color (4)", colorOffset)
		}
		points := make([]rl.Vector3, pointCount)
		for i := 0; i < pointCount; i++ {
			base := 1 + i*3
			points[i] = rl.Vector3{X: toFloat32(args[base]), Y: toFloat32(args[base+1]), Z: toFloat32(args[base+2])}
		}
		c := rl.White
		if len(args) >= colorOffset+4 {
			c = argsToColor(args, colorOffset)
		}
		for i := 0; i < pointCount-2; i++ {
			rl.DrawTriangle3D(points[i], points[i+1], points[i+2], c)
		}
		return nil, nil
	})
	v.RegisterForeign("DrawCubeWiresV", func(args []interface{}) (interface{}, error) {
		if len(args) < 7 {
			return nil, fmt.Errorf("DrawCubeWiresV requires (posX,posY,posZ, sizeX,sizeY,sizeZ, color)")
		}
		pos := rl.Vector3{X: toFloat32(args[0]), Y: toFloat32(args[1]), Z: toFloat32(args[2])}
		size := rl.Vector3{X: toFloat32(args[3]), Y: toFloat32(args[4]), Z: toFloat32(args[5])}
		c := rl.White
		if len(args) >= 11 {
			c = argsToColor(args, 6)
		}
		rl.DrawCubeWiresV(pos, size, c)
		return nil, nil
	})
	v.RegisterForeign("DrawSphereEx", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("DrawSphereEx requires (centerX,centerY,centerZ, radius, rings, slices, color)")
		}
		center := rl.Vector3{X: toFloat32(args[0]), Y: toFloat32(args[1]), Z: toFloat32(args[2])}
		radius := toFloat32(args[3])
		rings, slices := int32(16), int32(16)
		if len(args) >= 6 {
			rings, slices = toInt32(args[4]), toInt32(args[5])
		}
		c := rl.White
		if len(args) >= 10 {
			c = argsToColor(args, 6)
		} else if len(args) >= 9 {
			c = argsToColor(args, 5)
		}
		rl.DrawSphereEx(center, radius, rings, slices, c)
		return nil, nil
	})
	v.RegisterForeign("DrawCylinderEx", func(args []interface{}) (interface{}, error) {
		if len(args) < 8 {
			return nil, fmt.Errorf("DrawCylinderEx requires (startX,startY,startZ, endX,endY,endZ, startRadius, endRadius, sides, color)")
		}
		start := rl.Vector3{X: toFloat32(args[0]), Y: toFloat32(args[1]), Z: toFloat32(args[2])}
		end := rl.Vector3{X: toFloat32(args[3]), Y: toFloat32(args[4]), Z: toFloat32(args[5])}
		startRadius := toFloat32(args[6])
		endRadius := toFloat32(args[7])
		sides := int32(18)
		if len(args) >= 9 {
			sides = toInt32(args[8])
		}
		c := rl.White
		if len(args) >= 13 {
			c = argsToColor(args, 9)
		} else if len(args) >= 12 {
			c = argsToColor(args, 8)
		}
		rl.DrawCylinderEx(start, end, startRadius, endRadius, sides, c)
		return nil, nil
	})
	v.RegisterForeign("DrawCylinderWiresEx", func(args []interface{}) (interface{}, error) {
		if len(args) < 8 {
			return nil, fmt.Errorf("DrawCylinderWiresEx requires (startX,startY,startZ, endX,endY,endZ, startRadius, endRadius, sides, color)")
		}
		start := rl.Vector3{X: toFloat32(args[0]), Y: toFloat32(args[1]), Z: toFloat32(args[2])}
		end := rl.Vector3{X: toFloat32(args[3]), Y: toFloat32(args[4]), Z: toFloat32(args[5])}
		startRadius := toFloat32(args[6])
		endRadius := toFloat32(args[7])
		sides := int32(18)
		if len(args) >= 9 {
			sides = toInt32(args[8])
		}
		c := rl.White
		if len(args) >= 13 {
			c = argsToColor(args, 9)
		} else if len(args) >= 12 {
			c = argsToColor(args, 8)
		}
		rl.DrawCylinderWiresEx(start, end, startRadius, endRadius, sides, c)
		return nil, nil
	})
	v.RegisterForeign("DrawCapsule", func(args []interface{}) (interface{}, error) {
		if len(args) < 7 {
			return nil, fmt.Errorf("DrawCapsule requires (startX,startY,startZ, endX,endY,endZ, radius, slices, rings, color)")
		}
		start := rl.Vector3{X: toFloat32(args[0]), Y: toFloat32(args[1]), Z: toFloat32(args[2])}
		end := rl.Vector3{X: toFloat32(args[3]), Y: toFloat32(args[4]), Z: toFloat32(args[5])}
		radius := toFloat32(args[6])
		slices, rings := int32(16), int32(8)
		if len(args) >= 9 {
			slices, rings = toInt32(args[7]), toInt32(args[8])
		}
		c := rl.White
		if len(args) >= 13 {
			c = argsToColor(args, 9)
		} else if len(args) >= 12 {
			c = argsToColor(args, 8)
		}
		rl.DrawCapsule(start, end, radius, slices, rings, c)
		return nil, nil
	})
	v.RegisterForeign("DrawCapsuleWires", func(args []interface{}) (interface{}, error) {
		if len(args) < 7 {
			return nil, fmt.Errorf("DrawCapsuleWires requires (startX,startY,startZ, endX,endY,endZ, radius, slices, rings, color)")
		}
		start := rl.Vector3{X: toFloat32(args[0]), Y: toFloat32(args[1]), Z: toFloat32(args[2])}
		end := rl.Vector3{X: toFloat32(args[3]), Y: toFloat32(args[4]), Z: toFloat32(args[5])}
		radius := toFloat32(args[6])
		slices, rings := int32(16), int32(8)
		if len(args) >= 9 {
			slices, rings = toInt32(args[7]), toInt32(args[8])
		}
		c := rl.White
		if len(args) >= 13 {
			c = argsToColor(args, 9)
		} else if len(args) >= 12 {
			c = argsToColor(args, 8)
		}
		rl.DrawCapsuleWires(start, end, radius, slices, rings, c)
		return nil, nil
	})
	v.RegisterForeign("DrawModelEx", func(args []interface{}) (interface{}, error) {
		if len(args) < 9 {
			return nil, fmt.Errorf("DrawModelEx requires (id, posX,posY,posZ, rotAxisX,Y,Z, rotAngle, scaleX,scaleY,scaleZ, tint)")
		}
		id := toString(args[0])
		modelMu.Lock()
		model, ok := models[id]
		modelMu.Unlock()
		if !ok {
			return nil, fmt.Errorf("unknown model id: %s", id)
		}
		pos := rl.Vector3{X: toFloat32(args[1]), Y: toFloat32(args[2]), Z: toFloat32(args[3])}
		rotAxis := rl.Vector3{X: toFloat32(args[4]), Y: toFloat32(args[5]), Z: toFloat32(args[6])}
		rotAngle := toFloat32(args[7])
		scale := rl.Vector3{X: toFloat32(args[8]), Y: toFloat32(args[9]), Z: toFloat32(args[10])}
		c := rl.White
		if len(args) >= 15 {
			c = argsToColor(args, 11)
		}
		rl.DrawModelEx(model, pos, rotAxis, rotAngle, scale, c)
		return nil, nil
	})
	v.RegisterForeign("DrawModelWires", func(args []interface{}) (interface{}, error) {
		if len(args) < 5 {
			return nil, fmt.Errorf("DrawModelWires requires (id, posX, posY, posZ, scale, tint)")
		}
		id := toString(args[0])
		modelMu.Lock()
		model, ok := models[id]
		modelMu.Unlock()
		if !ok {
			return nil, fmt.Errorf("unknown model id: %s", id)
		}
		pos := rl.Vector3{X: toFloat32(args[1]), Y: toFloat32(args[2]), Z: toFloat32(args[3])}
		scale := toFloat32(args[4])
		c := rl.White
		if len(args) >= 9 {
			c = argsToColor(args, 5)
		}
		rl.DrawModelWires(model, pos, scale, c)
		return nil, nil
	})
	v.RegisterForeign("DrawBoundingBox", func(args []interface{}) (interface{}, error) {
		if len(args) < 7 {
			return nil, fmt.Errorf("DrawBoundingBox requires (minX,minY,minZ, maxX,maxY,maxZ, color)")
		}
		box := rl.BoundingBox{
			Min: rl.Vector3{X: toFloat32(args[0]), Y: toFloat32(args[1]), Z: toFloat32(args[2])},
			Max: rl.Vector3{X: toFloat32(args[3]), Y: toFloat32(args[4]), Z: toFloat32(args[5])},
		}
		c := rl.White
		if len(args) >= 11 {
			c = argsToColor(args, 6)
		}
		rl.DrawBoundingBox(box, c)
		return nil, nil
	})
	v.RegisterForeign("IsModelValid", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return false, nil
		}
		id := toString(args[0])
		modelMu.Lock()
		_, ok := models[id]
		modelMu.Unlock()
		return ok, nil
	})
	v.RegisterForeign("GetModelBoundingBox", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("GetModelBoundingBox requires (modelId)")
		}
		id := toString(args[0])
		modelMu.Lock()
		model, ok := models[id]
		modelMu.Unlock()
		if !ok {
			return nil, fmt.Errorf("unknown model id: %s", id)
		}
		box := rl.GetModelBoundingBox(model)
		return []interface{}{box.Min.X, box.Min.Y, box.Min.Z, box.Max.X, box.Max.Y, box.Max.Z}, nil
	})
	v.RegisterForeign("DrawModelWiresEx", func(args []interface{}) (interface{}, error) {
		if len(args) < 11 {
			return nil, fmt.Errorf("DrawModelWiresEx requires (id, posX,posY,posZ, rotAxisX,Y,Z, rotAngle, scaleX,scaleY,scaleZ, tint)")
		}
		id := toString(args[0])
		modelMu.Lock()
		model, ok := models[id]
		modelMu.Unlock()
		if !ok {
			return nil, fmt.Errorf("unknown model id: %s", id)
		}
		pos := rl.Vector3{X: toFloat32(args[1]), Y: toFloat32(args[2]), Z: toFloat32(args[3])}
		rotAxis := rl.Vector3{X: toFloat32(args[4]), Y: toFloat32(args[5]), Z: toFloat32(args[6])}
		rotAngle := toFloat32(args[7])
		scale := rl.Vector3{X: toFloat32(args[8]), Y: toFloat32(args[9]), Z: toFloat32(args[10])}
		c := rl.White
		if len(args) >= 15 {
			c = argsToColor(args, 11)
		}
		rl.DrawModelWiresEx(model, pos, rotAxis, rotAngle, scale, c)
		return nil, nil
	})
	v.RegisterForeign("DrawModelPoints", func(args []interface{}) (interface{}, error) {
		if len(args) < 5 {
			return nil, fmt.Errorf("DrawModelPoints requires (id, posX, posY, posZ, scale, tint)")
		}
		id := toString(args[0])
		modelMu.Lock()
		model, ok := models[id]
		modelMu.Unlock()
		if !ok {
			return nil, fmt.Errorf("unknown model id: %s", id)
		}
		pos := rl.Vector3{X: toFloat32(args[1]), Y: toFloat32(args[2]), Z: toFloat32(args[3])}
		scale := toFloat32(args[4])
		c := rl.White
		if len(args) >= 9 {
			c = argsToColor(args, 5)
		}
		rl.DrawModelPoints(model, pos, scale, c)
		return nil, nil
	})
	v.RegisterForeign("DrawModelPointsEx", func(args []interface{}) (interface{}, error) {
		if len(args) < 11 {
			return nil, fmt.Errorf("DrawModelPointsEx requires (id, posX,posY,posZ, rotAxisX,Y,Z, rotAngle, scaleX,scaleY,scaleZ, tint)")
		}
		id := toString(args[0])
		modelMu.Lock()
		model, ok := models[id]
		modelMu.Unlock()
		if !ok {
			return nil, fmt.Errorf("unknown model id: %s", id)
		}
		pos := rl.Vector3{X: toFloat32(args[1]), Y: toFloat32(args[2]), Z: toFloat32(args[3])}
		rotAxis := rl.Vector3{X: toFloat32(args[4]), Y: toFloat32(args[5]), Z: toFloat32(args[6])}
		rotAngle := toFloat32(args[7])
		scale := rl.Vector3{X: toFloat32(args[8]), Y: toFloat32(args[9]), Z: toFloat32(args[10])}
		c := rl.White
		if len(args) >= 15 {
			c = argsToColor(args, 11)
		}
		rl.DrawModelPointsEx(model, pos, rotAxis, rotAngle, scale, c)
		return nil, nil
	})
	v.RegisterForeign("DrawBillboard", func(args []interface{}) (interface{}, error) {
		if len(args) < 5 {
			return nil, fmt.Errorf("DrawBillboard requires (textureId, centerX,centerY,centerZ, scale, tint)")
		}
		texId := toString(args[0])
		texMu.Lock()
		tex, ok := textures[texId]
		texMu.Unlock()
		if !ok {
			return nil, fmt.Errorf("unknown texture id: %s", texId)
		}
		center := rl.Vector3{X: toFloat32(args[1]), Y: toFloat32(args[2]), Z: toFloat32(args[3])}
		scale := toFloat32(args[4])
		c := rl.White
		if len(args) >= 9 {
			c = argsToColor(args, 5)
		}
		rl.DrawBillboard(camera3D, tex, center, scale, c)
		return nil, nil
	})
	v.RegisterForeign("DrawBillboardRec", func(args []interface{}) (interface{}, error) {
		if len(args) < 9 {
			return nil, fmt.Errorf("DrawBillboardRec requires (textureId, srcX,srcY,srcW,srcH, centerX,centerY,centerZ, sizeX,sizeY, tint)")
		}
		texId := toString(args[0])
		texMu.Lock()
		tex, ok := textures[texId]
		texMu.Unlock()
		if !ok {
			return nil, fmt.Errorf("unknown texture id: %s", texId)
		}
		src := rl.Rectangle{X: toFloat32(args[1]), Y: toFloat32(args[2]), Width: toFloat32(args[3]), Height: toFloat32(args[4])}
		center := rl.Vector3{X: toFloat32(args[5]), Y: toFloat32(args[6]), Z: toFloat32(args[7])}
		size := rl.Vector2{X: toFloat32(args[8]), Y: toFloat32(args[9])}
		c := rl.White
		if len(args) >= 14 {
			c = argsToColor(args, 10)
		}
		rl.DrawBillboardRec(camera3D, tex, src, center, size, c)
		return nil, nil
	})
	v.RegisterForeign("SetModelMeshMaterial", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("SetModelMeshMaterial requires (modelId, meshIndex, materialIndex)")
		}
		id := toString(args[0])
		meshIndex := toInt32(args[1])
		materialIndex := toInt32(args[2])
		modelMu.Lock()
		model, ok := models[id]
		modelMu.Unlock()
		if !ok {
			return nil, fmt.Errorf("unknown model id: %s", id)
		}
		rl.SetModelMeshMaterial(&model, meshIndex, materialIndex)
		modelMu.Lock()
		models[id] = model
		modelMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("DrawBillboardPro", func(args []interface{}) (interface{}, error) {
		if len(args) < 13 {
			return nil, fmt.Errorf("DrawBillboardPro requires (textureId, srcX,srcY,srcW,srcH, posX,posY,posZ, upX,upY,upZ, sizeX,sizeY, originX,originY, rotation, tint)")
		}
		texId := toString(args[0])
		texMu.Lock()
		tex, ok := textures[texId]
		texMu.Unlock()
		if !ok {
			return nil, fmt.Errorf("unknown texture id: %s", texId)
		}
		src := rl.Rectangle{X: toFloat32(args[1]), Y: toFloat32(args[2]), Width: toFloat32(args[3]), Height: toFloat32(args[4])}
		pos := rl.Vector3{X: toFloat32(args[5]), Y: toFloat32(args[6]), Z: toFloat32(args[7])}
		up := rl.Vector3{X: toFloat32(args[8]), Y: toFloat32(args[9]), Z: toFloat32(args[10])}
		size := rl.Vector2{X: toFloat32(args[11]), Y: toFloat32(args[12])}
		origin := rl.Vector2{X: 0.5, Y: 0.5}
		rotation := float32(0)
		if len(args) >= 16 {
			origin.X, origin.Y = toFloat32(args[13]), toFloat32(args[14])
		}
		if len(args) >= 17 {
			rotation = toFloat32(args[15])
		}
		c := rl.White
		if len(args) >= 21 {
			c = argsToColor(args, 16)
		} else if len(args) >= 20 {
			c = argsToColor(args, 15)
		}
		rl.DrawBillboardPro(camera3D, tex, src, pos, up, size, origin, rotation, c)
		return nil, nil
	})

	// --- Model animations (storage in raylib.go: animations, lastLoadedAnimIds) ---
	v.RegisterForeign("LoadModelAnimations", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("LoadModelAnimations requires (fileName)")
		}
		path := toString(args[0])
		anims := rl.LoadModelAnimations(path)
		animMu.Lock()
		for _, a := range lastLoadedAnimIds {
			if anim, ok := animations[a]; ok {
				rl.UnloadModelAnimation(anim)
				delete(animations, a)
			}
		}
		lastLoadedAnimIds = make([]string, 0, len(anims))
		for _, a := range anims {
			animCounter++
			id := fmt.Sprintf("anim_%d", animCounter)
			animations[id] = a
			lastLoadedAnimIds = append(lastLoadedAnimIds, id)
		}
		animMu.Unlock()
		return int32(len(anims)), nil
	})

	v.RegisterForeign("GetModelAnimationId", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return "", nil
		}
		index := int(toInt32(args[0]))
		animMu.Lock()
		defer animMu.Unlock()
		if index < 0 || index >= len(lastLoadedAnimIds) {
			return "", nil
		}
		return lastLoadedAnimIds[index], nil
	})

	v.RegisterForeign("UpdateModelAnimation", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("UpdateModelAnimation requires (modelId, animId, frame)")
		}
		modelId := toString(args[0])
		animId := toString(args[1])
		frame := toInt32(args[2])
		modelMu.Lock()
		model, okModel := models[modelId]
		modelMu.Unlock()
		if !okModel {
			return nil, fmt.Errorf("unknown model id: %s", modelId)
		}
		animMu.Lock()
		anim, okAnim := animations[animId]
		animMu.Unlock()
		if !okAnim {
			return nil, fmt.Errorf("unknown animation id: %s", animId)
		}
		rl.UpdateModelAnimation(model, anim, frame)
		modelMu.Lock()
		models[modelId] = model
		modelMu.Unlock()
		return nil, nil
	})

	v.RegisterForeign("UpdateModelAnimationBones", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("UpdateModelAnimationBones requires (modelId, animId, frame)")
		}
		modelId := toString(args[0])
		animId := toString(args[1])
		frame := toInt32(args[2])
		modelMu.Lock()
		model, okModel := models[modelId]
		modelMu.Unlock()
		if !okModel {
			return nil, fmt.Errorf("unknown model id: %s", modelId)
		}
		animMu.Lock()
		anim, okAnim := animations[animId]
		animMu.Unlock()
		if !okAnim {
			return nil, fmt.Errorf("unknown animation id: %s", animId)
		}
		rl.UpdateModelAnimationBones(model, anim, frame)
		modelMu.Lock()
		models[modelId] = model
		modelMu.Unlock()
		return nil, nil
	})

	v.RegisterForeign("UnloadModelAnimation", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("UnloadModelAnimation requires (animId)")
		}
		animId := toString(args[0])
		animMu.Lock()
		anim, ok := animations[animId]
		delete(animations, animId)
		animMu.Unlock()
		if ok {
			rl.UnloadModelAnimation(anim)
		}
		return nil, nil
	})

	v.RegisterForeign("UnloadModelAnimations", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, nil
		}
		var list []rl.ModelAnimation
		animMu.Lock()
		for _, v := range args {
			id := toString(v)
			if anim, ok := animations[id]; ok {
				list = append(list, anim)
				delete(animations, id)
			}
		}
		// remove unloaded ids from lastLoadedAnimIds
		newLoaded := make([]string, 0, len(lastLoadedAnimIds))
		gone := make(map[string]bool)
		for _, v := range args {
			gone[toString(v)] = true
		}
		for _, id := range lastLoadedAnimIds {
			if !gone[id] {
				newLoaded = append(newLoaded, id)
			}
		}
		lastLoadedAnimIds = newLoaded
		animMu.Unlock()
		if len(list) > 0 {
			rl.UnloadModelAnimations(list)
		}
		return nil, nil
	})

	v.RegisterForeign("IsModelAnimationValid", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return false, nil
		}
		modelId := toString(args[0])
		animId := toString(args[1])
		modelMu.Lock()
		model, okModel := models[modelId]
		modelMu.Unlock()
		if !okModel {
			return false, nil
		}
		animMu.Lock()
		anim, okAnim := animations[animId]
		animMu.Unlock()
		if !okAnim {
			return false, nil
		}
		return rl.IsModelAnimationValid(model, anim), nil
	})
}
