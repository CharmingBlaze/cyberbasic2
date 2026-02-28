// Package raylib: 3D camera, models, and primitives (rmodels).
package raylib

import (
	"cyberbasic/compiler/vm"
	"fmt"
	"math"
	"sync"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type modelAnimState struct {
	ModelId      string
	AnimId       string
	CurrentTime  float64
	FPS          float64
	Loop         bool
	CurrentFrame int32
	FrameCount   int32
}

type modelTransform struct {
	Position   rl.Vector3
	RotAxis    rl.Vector3
	RotAngle   float32
	Scale      rl.Vector3
}

var (
	modelAnimStates     = make(map[string]*modelAnimState)
	modelAnimStateCtr   int
	modelAnimStateMu    sync.Mutex
	modelTransformState = make(map[string]*modelTransform)
	modelTransformMu    sync.Mutex
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
	// CAMERA3D() returns a new camera id; use SetCameraPosition/Target/Up/Fovy/Projection then SetCurrentCamera(id) before BeginMode3D()
	v.RegisterForeign("CAMERA3D", func(args []interface{}) (interface{}, error) {
		camMu.Lock()
		camCounter++
		id := fmt.Sprintf("cam_%d", camCounter)
		cameras[id] = rl.Camera3D{
			Position:   rl.Vector3{X: 0, Y: 0, Z: 0},
			Target:     rl.Vector3{X: 0, Y: 0, Z: -1},
			Up:         rl.Vector3{X: 0, Y: 1, Z: 0},
			Fovy:       60,
			Projection: rl.CameraPerspective,
		}
		camMu.Unlock()
		return id, nil
	})
	v.RegisterForeign("SetCameraPosition", func(args []interface{}) (interface{}, error) {
		if len(args) == 3 {
			// SetCameraPosition(x, y, z): set global/default camera position
			camera3D.Position = rl.Vector3{X: toFloat32(args[0]), Y: toFloat32(args[1]), Z: toFloat32(args[2])}
			return nil, nil
		}
		if len(args) < 4 {
			return nil, fmt.Errorf("SetCameraPosition requires (x, y, z) or (cameraId, x, y, z)")
		}
		id := toString(args[0])
		camMu.Lock()
		defer camMu.Unlock()
		c, ok := cameras[id]
		if !ok {
			return nil, fmt.Errorf("unknown camera id: %s", id)
		}
		c.Position = rl.Vector3{X: toFloat32(args[1]), Y: toFloat32(args[2]), Z: toFloat32(args[3])}
		cameras[id] = c
		return nil, nil
	})
	v.RegisterForeign("SetCameraTarget", func(args []interface{}) (interface{}, error) {
		if len(args) == 3 {
			tx, ty, tz := toFloat32(args[0]), toFloat32(args[1]), toFloat32(args[2])
			camera3D.Target = rl.Vector3{X: tx, Y: ty, Z: tz}
			orbitStateMu.Lock()
			orbitTargetX, orbitTargetY, orbitTargetZ = tx, ty, tz
			orbitStateMu.Unlock()
			return nil, nil
		}
		if len(args) < 4 {
			return nil, fmt.Errorf("SetCameraTarget requires (x, y, z) or (cameraId, x, y, z)")
		}
		id := toString(args[0])
		camMu.Lock()
		defer camMu.Unlock()
		c, ok := cameras[id]
		if !ok {
			return nil, fmt.Errorf("unknown camera id: %s", id)
		}
		c.Target = rl.Vector3{X: toFloat32(args[1]), Y: toFloat32(args[2]), Z: toFloat32(args[3])}
		cameras[id] = c
		return nil, nil
	})
	v.RegisterForeign("SetCameraUp", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("SetCameraUp requires (cameraId, x, y, z)")
		}
		id := toString(args[0])
		camMu.Lock()
		defer camMu.Unlock()
		c, ok := cameras[id]
		if !ok {
			return nil, fmt.Errorf("unknown camera id: %s", id)
		}
		c.Up = rl.Vector3{X: toFloat32(args[1]), Y: toFloat32(args[2]), Z: toFloat32(args[3])}
		cameras[id] = c
		return nil, nil
	})
	v.RegisterForeign("SetCameraFovy", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("SetCameraFovy requires (cameraId, fovy)")
		}
		id := toString(args[0])
		camMu.Lock()
		defer camMu.Unlock()
		c, ok := cameras[id]
		if !ok {
			return nil, fmt.Errorf("unknown camera id: %s", id)
		}
		c.Fovy = toFloat32(args[1])
		cameras[id] = c
		return nil, nil
	})
	v.RegisterForeign("SetCameraProjection", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("SetCameraProjection requires (cameraId, projection)")
		}
		id := toString(args[0])
		proj := toInt32(args[1])
		camMu.Lock()
		defer camMu.Unlock()
		c, ok := cameras[id]
		if !ok {
			return nil, fmt.Errorf("unknown camera id: %s", id)
		}
		c.Projection = rl.CameraProjection(proj)
		cameras[id] = c
		return nil, nil
	})
	v.RegisterForeign("SetCurrentCamera", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("SetCurrentCamera requires (cameraId)")
		}
		id := toString(args[0])
		camMu.Lock()
		c, ok := cameras[id]
		camMu.Unlock()
		if !ok {
			return nil, fmt.Errorf("unknown camera id: %s", id)
		}
		camera3D = c
		return nil, nil
	})
	v.RegisterForeign("CAMERA_PERSPECTIVE", func(args []interface{}) (interface{}, error) {
		return int(rl.CameraPerspective), nil
	})
	v.RegisterForeign("CAMERA_ORTHOGRAPHIC", func(args []interface{}) (interface{}, error) {
		return int(rl.CameraOrthographic), nil
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
	v.RegisterForeign("LoadModelAnimated", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("LoadModelAnimated requires (path)")
		}
		path := toString(args[0])
		model := rl.LoadModel(path)
		modelMu.Lock()
		modelCounter++
		id := fmt.Sprintf("model_%d", modelCounter)
		models[id] = model
		modelMu.Unlock()
		anims := rl.LoadModelAnimations(path)
		animMu.Lock()
		for _, a := range anims {
			animCounter++
			animId := fmt.Sprintf("anim_%d", animCounter)
			animations[animId] = a
			lastLoadedAnimIds = append(lastLoadedAnimIds, animId)
		}
		animMu.Unlock()
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
		modelStateMu.Lock()
		delete(modelColors, id)
		delete(modelAngles, id)
		modelStateMu.Unlock()
		return nil, nil
	})
	// LoadCube(size): create a cube model (GenMeshCube + LoadModelFromMesh); returns model id
	v.RegisterForeign("LoadCube", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("LoadCube requires (size)")
		}
		s := toFloat32(args[0])
		mesh := rl.GenMeshCube(s, s, s)
		meshMu.Lock()
		meshCounter++
		meshId := fmt.Sprintf("mesh_%d", meshCounter)
		meshes[meshId] = mesh
		meshMu.Unlock()
		model := rl.LoadModelFromMesh(mesh)
		modelMu.Lock()
		modelCounter++
		id := fmt.Sprintf("model_%d", modelCounter)
		models[id] = model
		modelMu.Unlock()
		return id, nil
	})
	// SetModelColor(modelId, r, g, b, a): store tint for simplified DrawModel
	v.RegisterForeign("SetModelColor", func(args []interface{}) (interface{}, error) {
		if len(args) < 5 {
			return nil, fmt.Errorf("SetModelColor requires (modelId, r, g, b, a)")
		}
		id := toString(args[0])
		c := rl.NewColor(toUint8(args[1]), toUint8(args[2]), toUint8(args[3]), toUint8(args[4]))
		modelStateMu.Lock()
		modelColors[id] = c
		modelStateMu.Unlock()
		return nil, nil
	})
	// RotateModel(modelId, speedDegPerSec): add speed*GetFrameTime() to stored angle (radians) for this model
	v.RegisterForeign("RotateModel", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("RotateModel requires (modelId, speedDegPerSec)")
		}
		id := toString(args[0])
		speedDeg := toFloat32(args[1])
		dt := rl.GetFrameTime()
		radPerSec := speedDeg * (3.14159265 / 180)
		modelStateMu.Lock()
		modelAngles[id] += radPerSec * dt
		modelStateMu.Unlock()
		return nil, nil
	})
	// DrawModelSimple(id, x, y, z [, angle]): draw with scale 1, axis Y; uses SetModelColor tint and RotateModel angle when omitted
	v.RegisterForeign("DrawModelSimple", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("DrawModelSimple requires (id, x, y, z) or (id, x, y, z, angle)")
		}
		id := toString(args[0])
		modelMu.Lock()
		model, ok := models[id]
		modelMu.Unlock()
		if !ok {
			return nil, fmt.Errorf("unknown model id: %s", id)
		}
		pos := rl.Vector3{X: toFloat32(args[1]), Y: toFloat32(args[2]), Z: toFloat32(args[3])}
		var angle float32
		if len(args) >= 5 {
			angle = toFloat32(args[4])
		} else {
			modelStateMu.Lock()
			angle = modelAngles[id]
			modelStateMu.Unlock()
		}
		modelStateMu.Lock()
		c := modelColors[id]
		modelStateMu.Unlock()
		if c.A == 0 && c.R == 0 && c.G == 0 && c.B == 0 {
			c = rl.White
		}
		rotAxis := rl.Vector3{X: 0, Y: 1, Z: 0}
		scale := rl.Vector3{X: 1, Y: 1, Z: 1}
		rl.DrawModelEx(model, pos, rotAxis, angle, scale, c)
		return nil, nil
	})
	v.RegisterForeign("DrawModel", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("DrawModel requires (id, posX, posY, posZ, scale) or (id, VECTOR3, scale [, tint])")
		}
		id := toString(args[0])
		modelMu.Lock()
		model, ok := models[id]
		modelMu.Unlock()
		if !ok {
			return nil, fmt.Errorf("unknown model id: %s", id)
		}
		var pos rl.Vector3
		var scale float32
		var tintStart int
		if sl, ok := args[1].([]interface{}); ok && len(sl) >= 3 {
			pos = rl.Vector3{X: toFloat32(sl[0]), Y: toFloat32(sl[1]), Z: toFloat32(sl[2])}
			scale = toFloat32(args[2])
			tintStart = 3
		} else {
			if len(args) < 5 {
				return nil, fmt.Errorf("DrawModel requires (id, posX, posY, posZ, scale) and optional tint")
			}
			pos = rl.Vector3{X: toFloat32(args[1]), Y: toFloat32(args[2]), Z: toFloat32(args[3])}
			scale = toFloat32(args[4])
			tintStart = 5
		}
		c := rl.White
		if len(args) >= tintStart+4 {
			c = argsToColor(args, tintStart)
		}
		rl.DrawModel(model, pos, scale, c)
		return nil, nil
	})
	// SetModelPosition(modelId, x, y, z): store position for DrawModelWithState
	v.RegisterForeign("SetModelPosition", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("SetModelPosition requires (modelId, x, y, z)")
		}
		id := toString(args[0])
		modelTransformMu.Lock()
		if modelTransformState[id] == nil {
			modelTransformState[id] = &modelTransform{
				Scale: rl.Vector3{X: 1, Y: 1, Z: 1},
			}
		}
		modelTransformState[id].Position = rl.Vector3{X: toFloat32(args[1]), Y: toFloat32(args[2]), Z: toFloat32(args[3])}
		modelTransformMu.Unlock()
		return nil, nil
	})
	// SetModelRotation(modelId, axisX, axisY, axisZ, angleRad): store rotation for DrawModelWithState
	v.RegisterForeign("SetModelRotation", func(args []interface{}) (interface{}, error) {
		if len(args) < 5 {
			return nil, fmt.Errorf("SetModelRotation requires (modelId, axisX, axisY, axisZ, angleRad)")
		}
		id := toString(args[0])
		modelTransformMu.Lock()
		if modelTransformState[id] == nil {
			modelTransformState[id] = &modelTransform{
				Scale: rl.Vector3{X: 1, Y: 1, Z: 1},
			}
		}
		modelTransformState[id].RotAxis = rl.Vector3{X: toFloat32(args[1]), Y: toFloat32(args[2]), Z: toFloat32(args[3])}
		modelTransformState[id].RotAngle = toFloat32(args[4])
		modelTransformMu.Unlock()
		return nil, nil
	})
	// SetModelScale(modelId, sx, sy, sz): store scale for DrawModelWithState
	v.RegisterForeign("SetModelScale", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("SetModelScale requires (modelId, sx, sy, sz)")
		}
		id := toString(args[0])
		modelTransformMu.Lock()
		if modelTransformState[id] == nil {
			modelTransformState[id] = &modelTransform{
				Scale: rl.Vector3{X: 1, Y: 1, Z: 1},
			}
		}
		modelTransformState[id].Scale = rl.Vector3{X: toFloat32(args[1]), Y: toFloat32(args[2]), Z: toFloat32(args[3])}
		modelTransformMu.Unlock()
		return nil, nil
	})
	// DrawModelWithState(modelId [, tint]): draw using stored position, rotation, scale (defaults if never set)
	v.RegisterForeign("DrawModelWithState", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("DrawModelWithState requires (modelId) and optional tint")
		}
		id := toString(args[0])
		modelMu.Lock()
		model, ok := models[id]
		modelMu.Unlock()
		if !ok {
			return nil, fmt.Errorf("unknown model id: %s", id)
		}
		modelTransformMu.Lock()
		t := modelTransformState[id]
		modelTransformMu.Unlock()
		pos := rl.Vector3{X: 0, Y: 0, Z: 0}
		rotAxis := rl.Vector3{X: 0, Y: 1, Z: 0}
		rotAngle := float32(0)
		scale := rl.Vector3{X: 1, Y: 1, Z: 1}
		if t != nil {
			pos = t.Position
			rotAxis = t.RotAxis
			rotAngle = t.RotAngle
			scale = t.Scale
		}
		c := rl.White
		if len(args) >= 5 {
			c = argsToColor(args, 1)
		}
		rl.DrawModelEx(model, pos, rotAxis, rotAngle, scale, c)
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
	// DrawText3D(fontId, text, posX, posY, posZ, fontSize, spacing, r, g, b, a): draw text at 3D world position (projects to screen using current camera).
	v.RegisterForeign("DrawText3D", func(args []interface{}) (interface{}, error) {
		if len(args) < 7 {
			return nil, fmt.Errorf("DrawText3D requires (fontId, text, posX, posY, posZ, fontSize, spacing) and optional tint r,g,b,a")
		}
		fontId := toString(args[0])
		text := toString(args[1])
		pos := rl.Vector3{X: toFloat32(args[2]), Y: toFloat32(args[3]), Z: toFloat32(args[4])}
		fontSize := toFloat32(args[5])
		spacing := toFloat32(args[6])
		screenPos := rl.GetWorldToScreen(pos, camera3D)
		var font rl.Font
		if fontId != "" {
			fontMu.Lock()
			f, ok := fonts[fontId]
			fontMu.Unlock()
			if ok {
				font = f
			} else {
				font = rl.GetFontDefault()
			}
		} else {
			font = rl.GetFontDefault()
		}
		c := rl.White
		if len(args) >= 11 {
			c = argsToColor(args, 7)
		}
		rl.DrawTextEx(font, text, screenPos, fontSize, spacing, c)
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
	v.RegisterForeign("PlayModelAnimation", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("PlayModelAnimation requires (modelId, animId)")
		}
		modelId := toString(args[0])
		animId := toString(args[1])
		modelMu.Lock()
		model, okModel := models[modelId]
		modelMu.Unlock()
		if !okModel {
			return nil, nil
		}
		animMu.Lock()
		anim, okAnim := animations[animId]
		animMu.Unlock()
		if !okAnim {
			return nil, nil
		}
		rl.UpdateModelAnimation(model, anim, 0)
		modelMu.Lock()
		models[modelId] = model
		modelMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("SetModelTexture", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("SetModelTexture requires (modelId, textureId)")
		}
		modelId := toString(args[0])
		texId := toString(args[1])
		modelMu.Lock()
		model, okModel := models[modelId]
		modelMu.Unlock()
		if !okModel {
			return nil, fmt.Errorf("unknown model id: %s", modelId)
		}
		texMu.Lock()
		tex, okTex := textures[texId]
		texMu.Unlock()
		if !okTex {
			return nil, fmt.Errorf("unknown texture id: %s", texId)
		}
		if model.MaterialCount > 0 && model.Materials != nil {
			// raylib-go: Materials is *Material, Maps is *MaterialMap; set diffuse texture
			model.Materials.Maps.Texture = tex
		}
		modelMu.Lock()
		models[modelId] = model
		modelMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("SetMaterialTexture", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("SetMaterialTexture requires (modelId, textureId)")
		}
		modelId := toString(args[0])
		texId := toString(args[1])
		modelMu.Lock()
		model, okModel := models[modelId]
		modelMu.Unlock()
		if !okModel {
			return nil, fmt.Errorf("unknown model id: %s", modelId)
		}
		texMu.Lock()
		tex, okTex := textures[texId]
		texMu.Unlock()
		if !okTex {
			return nil, fmt.Errorf("unknown texture id: %s", texId)
		}
		if model.MaterialCount > 0 && model.Materials != nil {
			model.Materials.Maps.Texture = tex
		}
		modelMu.Lock()
		models[modelId] = model
		modelMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("SetMaterialColor", func(args []interface{}) (interface{}, error) {
		if len(args) < 5 {
			return nil, fmt.Errorf("SetMaterialColor requires (modelId, r, g, b, a)")
		}
		id := toString(args[0])
		c := rl.NewColor(toUint8(args[1]), toUint8(args[2]), toUint8(args[3]), toUint8(args[4]))
		modelStateMu.Lock()
		modelColors[id] = c
		modelStateMu.Unlock()
		return nil, nil
	})
	// SetModelShader(modelId, shaderId): set the model's first material shader for custom rendering.
	v.RegisterForeign("SetModelShader", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("SetModelShader requires (modelId, shaderId)")
		}
		modelId := toString(args[0])
		shaderId := toString(args[1])
		modelMu.Lock()
		model, okModel := models[modelId]
		modelMu.Unlock()
		if !okModel {
			return nil, fmt.Errorf("unknown model id: %s", modelId)
		}
		shaderMu.Lock()
		sh, okShader := shaders[shaderId]
		shaderMu.Unlock()
		if !okShader {
			return nil, fmt.Errorf("unknown shader id: %s", shaderId)
		}
		if model.MaterialCount > 0 && model.Materials != nil {
			model.Materials.Shader = sh
			modelMu.Lock()
			models[modelId] = model
			modelMu.Unlock()
		}
		return nil, nil
	})
	// SetMaterialFloat(modelId, paramName, value): set float uniform on the model's first material shader.
	v.RegisterForeign("SetMaterialFloat", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("SetMaterialFloat requires (modelId, paramName, value)")
		}
		modelId := toString(args[0])
		paramName := toString(args[1])
		value := toFloat32(args[2])
		modelMu.Lock()
		model, ok := models[modelId]
		modelMu.Unlock()
		if !ok {
			return nil, fmt.Errorf("unknown model id: %s", modelId)
		}
		if model.MaterialCount == 0 || model.Materials == nil {
			return nil, nil
		}
		sh := model.Materials.Shader
		loc := rl.GetShaderLocation(sh, paramName)
		if loc >= 0 {
			rl.SetShaderValue(sh, loc, []float32{value}, rl.ShaderUniformFloat)
		}
		return nil, nil
	})
	// SetMaterialVector(modelId, paramName, x, y, z): set vec3 uniform on the model's first material shader.
	v.RegisterForeign("SetMaterialVector", func(args []interface{}) (interface{}, error) {
		if len(args) < 6 {
			return nil, fmt.Errorf("SetMaterialVector requires (modelId, paramName, x, y, z)")
		}
		modelId := toString(args[0])
		paramName := toString(args[1])
		vec := []float32{toFloat32(args[2]), toFloat32(args[3]), toFloat32(args[4])}
		modelMu.Lock()
		model, ok := models[modelId]
		modelMu.Unlock()
		if !ok {
			return nil, fmt.Errorf("unknown model id: %s", modelId)
		}
		if model.MaterialCount == 0 || model.Materials == nil {
			return nil, nil
		}
		sh := model.Materials.Shader
		loc := rl.GetShaderLocation(sh, paramName)
		if loc >= 0 {
			rl.SetShaderValue(sh, loc, vec, rl.ShaderUniformVec3)
		}
		return nil, nil
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

	v.RegisterForeign("GetModelAnimationFrameCount", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return 0, nil
		}
		animId := toString(args[0])
		animMu.Lock()
		anim, ok := animations[animId]
		animMu.Unlock()
		if !ok {
			return 0, nil
		}
		return int(anim.FrameCount), nil
	})

	// --- Model animation state (time-based playback) ---
	v.RegisterForeign("CreateModelAnimState", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("CreateModelAnimState requires (modelId, animId, fps [, loop])")
		}
		modelId := toString(args[0])
		animId := toString(args[1])
		fps := toFloat64(args[2])
		if fps <= 0 {
			fps = 24
		}
		loop := true
		if len(args) >= 4 && args[3] != nil {
			switch x := args[3].(type) {
			case bool:
				loop = x
			case int:
				loop = x != 0
			case float64:
				loop = x != 0
			}
		}
		modelMu.Lock()
		_, okModel := models[modelId]
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
		frameCount := anim.FrameCount
		if frameCount <= 0 {
			frameCount = 1
		}
		modelAnimStateMu.Lock()
		modelAnimStateCtr++
		stateId := fmt.Sprintf("modelanim_%d", modelAnimStateCtr)
		modelAnimStates[stateId] = &modelAnimState{
			ModelId:      modelId,
			AnimId:       animId,
			FPS:          fps,
			Loop:         loop,
			FrameCount:   frameCount,
		}
		modelAnimStateMu.Unlock()
		return stateId, nil
	})

	v.RegisterForeign("UpdateModelAnimState", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("UpdateModelAnimState requires (stateId, deltaTime)")
		}
		stateId := toString(args[0])
		dt := toFloat64(args[1])
		modelAnimStateMu.Lock()
		st, ok := modelAnimStates[stateId]
		modelAnimStateMu.Unlock()
		if !ok {
			return nil, fmt.Errorf("unknown model anim state id: %s", stateId)
		}
		st.CurrentTime += dt
		frame := int32(st.CurrentTime * st.FPS)
		if st.Loop && st.FrameCount > 0 {
			for frame >= st.FrameCount {
				frame -= st.FrameCount
			}
			for frame < 0 {
				frame += st.FrameCount
			}
		} else {
			if frame >= st.FrameCount {
				frame = st.FrameCount - 1
			}
			if frame < 0 {
				frame = 0
			}
		}
		st.CurrentFrame = frame
		modelMu.Lock()
		model, okModel := models[st.ModelId]
		modelMu.Unlock()
		if !okModel {
			return nil, nil
		}
		animMu.Lock()
		anim, okAnim := animations[st.AnimId]
		animMu.Unlock()
		if !okAnim {
			return nil, nil
		}
		rl.UpdateModelAnimation(model, anim, frame)
		modelMu.Lock()
		models[st.ModelId] = model
		modelMu.Unlock()
		return nil, nil
	})

	v.RegisterForeign("SetModelAnimStateFrame", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("SetModelAnimStateFrame requires (stateId, frameIndex)")
		}
		stateId := toString(args[0])
		idx := toInt32(args[1])
		modelAnimStateMu.Lock()
		st, ok := modelAnimStates[stateId]
		modelAnimStateMu.Unlock()
		if !ok {
			return nil, fmt.Errorf("unknown model anim state id: %s", stateId)
		}
		if idx < 0 {
			idx = 0
		}
		if idx >= st.FrameCount {
			idx = st.FrameCount - 1
		}
		st.CurrentFrame = idx
		st.CurrentTime = float64(idx) / st.FPS
		modelMu.Lock()
		model, okModel := models[st.ModelId]
		modelMu.Unlock()
		if !okModel {
			return nil, nil
		}
		animMu.Lock()
		anim, okAnim := animations[st.AnimId]
		animMu.Unlock()
		if !okAnim {
			return nil, nil
		}
		rl.UpdateModelAnimation(model, anim, idx)
		modelMu.Lock()
		models[st.ModelId] = model
		modelMu.Unlock()
		return nil, nil
	})

	v.RegisterForeign("GetModelAnimStateFrame", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return 0, nil
		}
		stateId := toString(args[0])
		modelAnimStateMu.Lock()
		defer modelAnimStateMu.Unlock()
		st, ok := modelAnimStates[stateId]
		if !ok {
			return 0, nil
		}
		return int(st.CurrentFrame), nil
	})

	v.RegisterForeign("DestroyModelAnimState", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, nil
		}
		stateId := toString(args[0])
		modelAnimStateMu.Lock()
		delete(modelAnimStates, stateId)
		modelAnimStateMu.Unlock()
		return nil, nil
	})

	// --- Lighting (state stored for custom shaders; raylib has no built-in lighting) ---
	v.RegisterForeign("ENABLELIGHTING", func(args []interface{}) (interface{}, error) {
		lightingOn = true
		return nil, nil
	})
	v.RegisterForeign("LIGHT", func(args []interface{}) (interface{}, error) {
		lightMu.Lock()
		lightCtr++
		id := fmt.Sprintf("light_%d", lightCtr)
		lightIds[id] = true
		lightMu.Unlock()
		return id, nil
	})
	v.RegisterForeign("LIGHT_DIRECTIONAL", func(args []interface{}) (interface{}, error) {
		return 0, nil
	})
	v.RegisterForeign("CreateLight", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("CreateLight requires (type, x, y, z)")
		}
		lightType := toInt32(args[0])
		x, y, z := toFloat32(args[1]), toFloat32(args[2]), toFloat32(args[3])
		lightMu.Lock()
		lightCtr++
		id := fmt.Sprintf("light_%d", lightCtr)
		lightIds[id] = true
		lightMu.Unlock()
		lightDataMu.Lock()
		lightData[id] = &struct {
			Type      int
			X, Y, Z   float32
			R, G, B   uint8
			Intensity float32
			DirX, DirY, DirZ float32
		}{Type: int(lightType), X: x, Y: y, Z: z, R: 255, G: 255, B: 255, Intensity: 1}
		lightDataMu.Unlock()
		return id, nil
	})
	v.RegisterForeign("SetLightType", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("SetLightType requires (lightId, type)")
		}
		id := toString(args[0])
		lightDataMu.Lock()
		if d, ok := lightData[id]; ok {
			d.Type = int(toInt32(args[1]))
		}
		lightDataMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("SetLightPosition", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("SetLightPosition requires (lightId, x, y, z)")
		}
		id := toString(args[0])
		lightDataMu.Lock()
		if d, ok := lightData[id]; ok {
			d.X, d.Y, d.Z = toFloat32(args[1]), toFloat32(args[2]), toFloat32(args[3])
		}
		lightDataMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("SetLightTarget", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("SetLightTarget requires (lightId, x, y, z)")
		}
		id := toString(args[0])
		tx, ty, tz := toFloat32(args[1]), toFloat32(args[2]), toFloat32(args[3])
		lightDataMu.Lock()
		if d, ok := lightData[id]; ok {
			dx, dy, dz := tx-d.X, ty-d.Y, tz-d.Z
			len := float32(math.Sqrt(float64(dx*dx + dy*dy + dz*dz)))
			if len > 1e-6 {
				d.DirX, d.DirY, d.DirZ = dx/len, dy/len, dz/len
			}
		}
		lightDataMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("SetLightColor", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("SetLightColor requires (lightId, r, g, b)")
		}
		id := toString(args[0])
		lightDataMu.Lock()
		if d, ok := lightData[id]; ok {
			d.R = uint8(toFloat64(args[1])) & 0xff
			d.G = uint8(toFloat64(args[2])) & 0xff
			d.B = uint8(toFloat64(args[3])) & 0xff
		}
		lightDataMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("SetLightIntensity", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("SetLightIntensity requires (lightId, amount)")
		}
		id := toString(args[0])
		lightDataMu.Lock()
		if d, ok := lightData[id]; ok {
			d.Intensity = toFloat32(args[1])
		}
		lightDataMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("SetLightDirection", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("SetLightDirection requires (lightId, x, y, z)")
		}
		id := toString(args[0])
		dx, dy, dz := toFloat32(args[1]), toFloat32(args[2]), toFloat32(args[3])
		lightDataMu.Lock()
		if d, ok := lightData[id]; ok {
			len := float32(math.Sqrt(float64(dx*dx + dy*dy + dz*dz)))
			if len > 1e-6 {
				d.DirX, d.DirY, d.DirZ = dx/len, dy/len, dz/len
			}
		}
		lightDataMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("EnableShadows", func(args []interface{}) (interface{}, error) {
		shadowsEnabled = true
		return nil, nil
	})
	v.RegisterForeign("DisableShadows", func(args []interface{}) (interface{}, error) {
		shadowsEnabled = false
		return nil, nil
	})
	v.RegisterForeign("SETAMBIENTLIGHT", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("SETAMBIENTLIGHT requires (r, g, b)")
		}
		ambientR = toFloat32(args[0])
		ambientG = toFloat32(args[1])
		ambientB = toFloat32(args[2])
		return nil, nil
	})

	// Phase 8: Optimization (store flags for future use)
	v.RegisterForeign("SetCullingDistance", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("SetCullingDistance requires (distance)")
		}
		cullingDistance = toFloat32(args[0])
		return nil, nil
	})
	v.RegisterForeign("EnableFrustumCulling", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("EnableFrustumCulling requires (flag)")
		}
		frustumCulling = toFloat32(args[0]) != 0
		return nil, nil
	})
	v.RegisterForeign("Enable2DCulling", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("Enable2DCulling requires (flag)")
		}
		enable2DCulling = toFloat32(args[0]) != 0
		return nil, nil
	})
	v.RegisterForeign("SetCullingMargin", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("SetCullingMargin requires (pixels)")
		}
		cullingMargin = toFloat32(args[0])
		return nil, nil
	})
}
