// Package dbp provides DarkBASIC Pro-style high-level commands as thin wrappers over raylib.
//
// Uses integer IDs for images, objects, sounds, textures, and materials (DBP convention).
// Commands are split into modular files:
//   - dbp.go: Core 2D/3D, objects, scene, input, FPS camera
//   - dbp_textures.go, dbp_materials.go: Texture and material registries
//   - dbp_camera.go: CameraFollow, CameraOrbit, CameraShake, CameraSmooth
//   - dbp_world.go: SetSkybox, SetAmbientLight, SetFog
//   - dbp_groups.go, dbp_players.go: Object groups and player state
//   - dbp_audio.go, dbp_lighting.go: Music and light registries
//   - dbp_level.go: LoadLevel, DrawLevel, UnloadLevel
//   - dbp_model.go: BuildModel (unified importer pipeline)
//   - dbp_physics.go: Physics wrappers (PhysicsOn/Off, MakeRigidBody, SetGravity)
//   - dbp_collision.go: Raycast, Spherecast, ObjectCollides, PointInObject
//   - dbp_particles.go: MakeParticles, SetParticleColor/Size/Speed, EmitParticlesAt, DrawParticles
//   - dbp_net.go: Networking (NetConnect, NetSend, etc. from net package)
//   - dbp_file.go: SaveString, LoadString, SaveValue, LoadValue
//   - dbp_runtime.go: StopTask, PauseTask, ResumeTask (stubs)
//   - dbp_replication.go: Replication (ReplicatePosition, etc. from game package)
//
// See docs/DBP_EXTENDED.md for the full command reference.
package dbp

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"cyberbasic/compiler/bindings/model"
	"cyberbasic/compiler/vm"
	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
	images   = make(map[int]rl.Texture2D)
	imagesMu sync.Mutex

	objects   = make(map[int]*dbpObject)
	objectsMu sync.Mutex

	sounds   = make(map[int]rl.Sound)
	soundsMu sync.Mutex

	fonts   = make(map[int]rl.Font)
	fontsMu sync.Mutex
	currentFontID int = -1 // -1 = default font

	// Current draw color for Ink(r, g, b)
	inkR, inkG, inkB, inkA uint8 = 255, 255, 255, 255
	inkMu                   sync.Mutex

	// FPS camera state
	fpsCameraOn    bool
	fpsCamX        float32
	fpsCamY        float32
	fpsCamZ        float32
	fpsYaw         float32
	fpsPitch       float32
	fpsMoveSpeed   float32 = 5.0
	fpsLookSpeed   float32 = 0.002
	fpsCameraMu    sync.Mutex
)

type dbpObject struct {
	model       rl.Model
	x, y, z     float32
	pitch       float32
	yaw         float32
	roll        float32
	scaleX      float32
	scaleY      float32
	scaleZ      float32
	visible     bool
	colorR      uint8
	colorG      uint8
	colorB      uint8
	colorA      uint8
	textureId   int
	shaderId    int // For custom shaders; DrawModelEx doesn't take shader
	wireframe   bool
	collision   bool
	fixed       bool
	parentID    int   // -1 = no parent; for ParentObject/UnparentObject
	tag         string
	ownerID     int   // For multiplayer: player who owns this object
	syncMe      bool  // Mark for replication
}

func newDbpObject(model rl.Model) *dbpObject {
	return &dbpObject{
		model:    model,
		scaleX:   1, scaleY: 1, scaleZ: 1,
		visible:  true,
		colorR:   255, colorG: 255, colorB: 255, colorA: 255,
		parentID: -1,
	}
}

func toInt(v interface{}) int {
	switch x := v.(type) {
	case int:
		return x
	case float64:
		return int(math.Trunc(x))
	case string:
		n, _ := strconv.Atoi(x)
		return n
	default:
		return 0
	}
}

func toFloat32(v interface{}) float32 {
	switch x := v.(type) {
	case int:
		return float32(x)
	case float64:
		return float32(x)
	case string:
		f, _ := strconv.ParseFloat(x, 32)
		return float32(f)
	default:
		return 0
	}
}

func toString(v interface{}) string {
	if v == nil {
		return ""
	}
	return fmt.Sprint(v)
}

// getObjectWorldTransform returns effective world position, rotation, and scale for an object.
// When parentID >= 0, recursively composes parent transforms.
func getObjectWorldTransform(id int) (x, y, z, pitch, yaw, roll, scaleX, scaleY, scaleZ float32) {
	objectsMu.Lock()
	obj, ok := objects[id]
	objectsMu.Unlock()
	if !ok {
		return 0, 0, 0, 0, 0, 0, 1, 1, 1
	}
	x, y, z = obj.x, obj.y, obj.z
	pitch, yaw, roll = obj.pitch, obj.yaw, obj.roll
	scaleX, scaleY, scaleZ = obj.scaleX, obj.scaleY, obj.scaleZ
	if obj.parentID < 0 {
		return x, y, z, pitch, yaw, roll, scaleX, scaleY, scaleZ
	}
	// Compose with parent's world transform
	px, py, pz, pp, pyaw, pr, psx, psy, psz := getObjectWorldTransform(obj.parentID)
	// Rotate child local offset by parent yaw (Y-up), scale by parent scale, add to parent pos
	rad := pyaw * math.Pi / 180
	c, s := float32(math.Cos(float64(rad))), float32(math.Sin(float64(rad)))
	lx := (x*psx)*c - (z*psz)*s
	lz := (x*psx)*s + (z*psz)*c
	ly := y * psy
	x = px + lx
	y = py + ly
	z = pz + lz
	pitch = pp + pitch
	yaw = pyaw + yaw
	roll = pr + roll
	scaleX = psx * scaleX
	scaleY = psy * scaleY
	scaleZ = psz * scaleZ
	return x, y, z, pitch, yaw, roll, scaleX, scaleY, scaleZ
}

// RegisterDBP registers DBP-style commands with the VM.
func RegisterDBP(v *vm.VM) {
	registerTextures(v)
	registerMaterials(v)
	registerCameraExtras(v)
	registerWorld(v)
	registerGroups(v)
	registerPlayers(v)
	registerAudioExpanded(v)
	registerLighting(v)
	registerPhysics(v)
	registerCollision(v)
	registerParticles(v)
	registerNet(v)
	registerFile(v)
	registerRuntime(v)
	registerReplication(v)
	register3D(v)
	registerLevel(v)
	registerPrefab(v)
	registerIK(v)
	registerInstancing(v)
	registerNav(v)
	// register2D is called from main after game so SetTile/GetTile overwrite game's
	// RegisterWater and RegisterTerrain are called from main after water/terrain packages (integer-ID API)
	// --- 2D Graphics ---
	// LoadImage(path, id): DBP-style, load texture and store at integer id.
	// Overrides raylib LoadImage (which loads to Image, not Texture). Use LoadTexture for raylib-style.
	v.RegisterForeign("LoadImage", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("LoadImage(path, id) requires 2 arguments")
		}
		path := toString(args[0])
		id := toInt(args[1])
		tex := rl.LoadTexture(path)
		imagesMu.Lock()
		images[id] = tex
		imagesMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("Sprite", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("Sprite(id, x, y) requires 3 arguments")
		}
		id := toInt(args[0])
		x, y := toFloat32(args[1]), toFloat32(args[2])
		imagesMu.Lock()
		tex, ok := images[id]
		imagesMu.Unlock()
		if !ok {
			return nil, fmt.Errorf("LoadImage: unknown image id %d", id)
		}
		inkMu.Lock()
		c := rl.NewColor(inkR, inkG, inkB, inkA)
		inkMu.Unlock()
		rl.DrawTexture(tex, int32(x), int32(y), c)
		return nil, nil
	})
	v.RegisterForeign("Cls", func(args []interface{}) (interface{}, error) {
		inkMu.Lock()
		c := rl.NewColor(inkR, inkG, inkB, inkA)
		inkMu.Unlock()
		rl.ClearBackground(c)
		return nil, nil
	})
	v.RegisterForeign("CLS", func(args []interface{}) (interface{}, error) {
		rl.ClearBackground(rl.Black)
		return nil, nil
	})
	v.RegisterForeign("Ink", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("Ink(r, g, b) requires 3 arguments")
		}
		inkMu.Lock()
		inkR = uint8(toInt(args[0]) & 0xff)
		inkG = uint8(toInt(args[1]) & 0xff)
		inkB = uint8(toInt(args[2]) & 0xff)
		inkA = 255
		inkMu.Unlock()
		return nil, nil
	})

	// --- 3D ---
	v.RegisterForeign("LoadObject", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("LoadObject(path, id) requires 2 arguments")
		}
		path := toString(args[0])
		id := toInt(args[1])
		model := rl.LoadModel(path)
		objectsMu.Lock()
		objects[id] = newDbpObject(model)
		objectsMu.Unlock()
		return nil, nil
	})
	// LoadObject(id, path): DBP arg order - id first, then path. Uses unified importer pipeline.
	v.RegisterForeign("LoadObjectId", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("LoadObject(id, path) requires 2 arguments")
		}
		id := toInt(args[0])
		path := toString(args[1])
		m, err := model.Load(path)
		if err != nil {
			return nil, fmt.Errorf("LoadObject: %w", err)
		}
		basePath := filepath.Dir(path)
		res, err := BuildModel(m, id, basePath)
		if err != nil {
			return nil, fmt.Errorf("LoadObject build: %w", err)
		}
		// For LoadObject we use the first object only (DBP convention: one id = one object)
		if len(res.ObjectIDs) > 1 {
			// Remove extra objects - user expects single object at id
			objectsMu.Lock()
			for i := 1; i < len(res.ObjectIDs); i++ {
				if obj, ok := objects[res.ObjectIDs[i]]; ok {
					delete(objects, res.ObjectIDs[i])
					if obj.model.MeshCount > 0 {
						rl.UnloadModel(obj.model)
					}
				}
			}
			objectsMu.Unlock()
		}
		if len(res.ObjectIDs) > 0 && res.ObjectIDs[0] != id {
			objectsMu.Lock()
			if obj, ok := objects[res.ObjectIDs[0]]; ok {
				delete(objects, res.ObjectIDs[0])
				objects[id] = obj
			}
			objectsMu.Unlock()
		}
		return nil, nil
	})
	// LoadCube(id, size): create procedural cube, DBP-style. Use when no .obj file available.
	v.RegisterForeign("LoadCube", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("LoadCube(id, size) requires 2 arguments")
		}
		id := toInt(args[0])
		size := toFloat32(args[1])
		mesh := rl.GenMeshCube(size, size, size)
		model := rl.LoadModelFromMesh(mesh)
		objectsMu.Lock()
		objects[id] = newDbpObject(model)
		objectsMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("PositionObject", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("PositionObject(id, x, y, z) requires 4 arguments")
		}
		id := toInt(args[0])
		x, y, z := toFloat32(args[1]), toFloat32(args[2]), toFloat32(args[3])
		objectsMu.Lock()
		obj, ok := objects[id]
		if ok {
			obj.x, obj.y, obj.z = x, y, z
		}
		objectsMu.Unlock()
		if !ok {
			return nil, fmt.Errorf("LoadObject: unknown object id %d", id)
		}
		return nil, nil
	})
	v.RegisterForeign("RotateObject", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("RotateObject(id, pitch, yaw, roll) requires 4 arguments")
		}
		id := toInt(args[0])
		pitch, yaw, roll := toFloat32(args[1]), toFloat32(args[2]), toFloat32(args[3])
		objectsMu.Lock()
		obj, ok := objects[id]
		if ok {
			obj.pitch, obj.yaw, obj.roll = pitch, yaw, roll
		}
		objectsMu.Unlock()
		if !ok {
			return nil, fmt.Errorf("LoadObject: unknown object id %d", id)
		}
		return nil, nil
	})
	v.RegisterForeign("YRotateObject", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("YRotateObject(id, angle) requires 2 arguments")
		}
		id := toInt(args[0])
		angle := toFloat32(args[1])
		objectsMu.Lock()
		obj, ok := objects[id]
		if ok {
			obj.yaw += angle
		}
		objectsMu.Unlock()
		if !ok {
			return nil, fmt.Errorf("LoadObject: unknown object id %d", id)
		}
		return nil, nil
	})
	// DrawObject(id): draw 3D object at its stored position/rotation/scale. Call between BeginMode3D/EndMode3D.
	v.RegisterForeign("DrawObject", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("DrawObject(id) requires 1 argument")
		}
		id := toInt(args[0])
		objectsMu.Lock()
		obj, ok := objects[id]
		objectsMu.Unlock()
		if !ok {
			return nil, fmt.Errorf("LoadObject: unknown object id %d", id)
		}
		if !obj.visible {
			return nil, nil
		}
		UpdateObjectAnimation(id, obj)
		wx, wy, wz, _, wyaw, _, wsx, wsy, wsz := getObjectWorldTransform(id)
		pos := rl.Vector3{X: wx, Y: wy, Z: wz}
		rotAxis := rl.Vector3{X: 0, Y: 1, Z: 0}
		rotAngle := wyaw * math.Pi / 180
		scale := rl.Vector3{X: wsx, Y: wsy, Z: wsz}
		tint := rl.NewColor(obj.colorR, obj.colorG, obj.colorB, obj.colorA)
		if obj.wireframe {
			rl.DrawModelWiresEx(obj.model, pos, rotAxis, rotAngle, scale, tint)
		} else {
			rl.DrawModelEx(obj.model, pos, rotAxis, rotAngle, scale, tint)
		}
		return nil, nil
	})
	// PointCamera(x, y, z): set camera target - delegates to raylib SetCameraTarget
	v.RegisterForeign("PointCamera", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("PointCamera(x, y, z) requires 3 arguments")
		}
		_, err := v.CallForeign("SetCameraTarget", args)
		return nil, err
	})
	// MakeCube(id, size): alias for LoadCube
	v.RegisterForeign("MakeCube", func(args []interface{}) (interface{}, error) {
		return v.CallForeign("LoadCube", args)
	})
	// MakeSphere(id, radius): create procedural sphere
	v.RegisterForeign("MakeSphere", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("MakeSphere(id, radius) requires 2 arguments")
		}
		id := toInt(args[0])
		radius := toFloat32(args[1])
		mesh := rl.GenMeshSphere(radius, 16, 16)
		model := rl.LoadModelFromMesh(mesh)
		objectsMu.Lock()
		objects[id] = newDbpObject(model)
		objectsMu.Unlock()
		return nil, nil
	})
	// MakePlane(id, width, height): create procedural plane
	v.RegisterForeign("MakePlane", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("MakePlane(id, width, height) requires 3 arguments")
		}
		id := toInt(args[0])
		w, h := toFloat32(args[1]), toFloat32(args[2])
		mesh := rl.GenMeshPlane(w, h, 1, 1)
		model := rl.LoadModelFromMesh(mesh)
		objectsMu.Lock()
		objects[id] = newDbpObject(model)
		objectsMu.Unlock()
		return nil, nil
	})
	// DeleteObject(id): remove object and unload model
	v.RegisterForeign("DeleteObject", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("DeleteObject(id) requires 1 argument")
		}
		id := toInt(args[0])
		UnregisterObjectModel(id)
		objectsMu.Lock()
		obj, ok := objects[id]
		delete(objects, id)
		objectsMu.Unlock()
		if ok && obj.model.MeshCount > 0 {
			rl.UnloadModel(obj.model)
		}
		return nil, nil
	})
	// ScaleObject(id, sx, sy, sz): set object scale
	v.RegisterForeign("ScaleObject", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("ScaleObject(id, sx, sy, sz) requires 4 arguments")
		}
		id := toInt(args[0])
		sx, sy, sz := toFloat32(args[1]), toFloat32(args[2]), toFloat32(args[3])
		objectsMu.Lock()
		obj, ok := objects[id]
		if ok {
			obj.scaleX, obj.scaleY, obj.scaleZ = sx, sy, sz
		}
		objectsMu.Unlock()
		if !ok {
			return nil, fmt.Errorf("unknown object id %d", id)
		}
		return nil, nil
	})
	// MoveObject(id, x, y, z): add to position
	v.RegisterForeign("MoveObject", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("MoveObject(id, x, y, z) requires 4 arguments")
		}
		id := toInt(args[0])
		dx, dy, dz := toFloat32(args[1]), toFloat32(args[2]), toFloat32(args[3])
		objectsMu.Lock()
		obj, ok := objects[id]
		if ok {
			obj.x += dx
			obj.y += dy
			obj.z += dz
		}
		objectsMu.Unlock()
		if !ok {
			return nil, fmt.Errorf("unknown object id %d", id)
		}
		return nil, nil
	})
	// TurnObject(id, pitch, yaw, roll): add to rotation
	v.RegisterForeign("TurnObject", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("TurnObject(id, pitch, yaw, roll) requires 4 arguments")
		}
		id := toInt(args[0])
		dp, dy, dr := toFloat32(args[1]), toFloat32(args[2]), toFloat32(args[3])
		objectsMu.Lock()
		obj, ok := objects[id]
		if ok {
			obj.pitch += dp
			obj.yaw += dy
			obj.roll += dr
		}
		objectsMu.Unlock()
		if !ok {
			return nil, fmt.Errorf("unknown object id %d", id)
		}
		return nil, nil
	})
	// HideObject(id): set visible=false
	v.RegisterForeign("HideObject", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("HideObject(id) requires 1 argument")
		}
		id := toInt(args[0])
		objectsMu.Lock()
		obj, ok := objects[id]
		if ok {
			obj.visible = false
		}
		objectsMu.Unlock()
		if !ok {
			return nil, fmt.Errorf("unknown object id %d", id)
		}
		return nil, nil
	})
	// ShowObject(id): set visible=true
	v.RegisterForeign("ShowObject", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("ShowObject(id) requires 1 argument")
		}
		id := toInt(args[0])
		objectsMu.Lock()
		obj, ok := objects[id]
		if ok {
			obj.visible = true
		}
		objectsMu.Unlock()
		if !ok {
			return nil, fmt.Errorf("unknown object id %d", id)
		}
		return nil, nil
	})

	// --- Object extras ---
	v.RegisterForeign("CloneObject", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("CloneObject(newID, sourceID) requires 2 arguments")
		}
		newID := toInt(args[0])
		srcID := toInt(args[1])
		objectsMu.Lock()
		src, ok := objects[srcID]
		if !ok {
			objectsMu.Unlock()
			return nil, fmt.Errorf("unknown object id %d", srcID)
		}
		// Create procedural cube with same scale as source (raylib has no model clone)
		mesh := rl.GenMeshCube(src.scaleX*2, src.scaleY*2, src.scaleZ*2)
		newModel := rl.LoadModelFromMesh(mesh)
		clone := &dbpObject{
			model:     newModel,
			x:         src.x, y: src.y, z: src.z,
			pitch:     src.pitch, yaw: src.yaw, roll: src.roll,
			scaleX:    1, scaleY: 1, scaleZ: 1,
			visible:   src.visible,
			colorR:    src.colorR, colorG: src.colorG, colorB: src.colorB, colorA: src.colorA,
			textureId: src.textureId, wireframe: src.wireframe, collision: src.collision, fixed: src.fixed,
		}
		objects[newID] = clone
		objectsMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("FixObject", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("FixObject(id) requires 1 argument")
		}
		id := toInt(args[0])
		objectsMu.Lock()
		obj, ok := objects[id]
		if ok {
			obj.fixed = true
		}
		objectsMu.Unlock()
		if !ok {
			return nil, fmt.Errorf("unknown object id %d", id)
		}
		return nil, nil
	})
	v.RegisterForeign("UnfixObject", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("UnfixObject(id) requires 1 argument")
		}
		id := toInt(args[0])
		objectsMu.Lock()
		obj, ok := objects[id]
		if ok {
			obj.fixed = false
		}
		objectsMu.Unlock()
		if !ok {
			return nil, fmt.Errorf("unknown object id %d", id)
		}
		return nil, nil
	})
	v.RegisterForeign("SetObjectColor", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("SetObjectColor(id, r, g, b) requires 4 arguments")
		}
		id := toInt(args[0])
		r, g, b := toInt(args[1])&0xff, toInt(args[2])&0xff, toInt(args[3])&0xff
		objectsMu.Lock()
		obj, ok := objects[id]
		if ok {
			obj.colorR, obj.colorG, obj.colorB = uint8(r), uint8(g), uint8(b)
			obj.colorA = 255
		}
		objectsMu.Unlock()
		if !ok {
			return nil, fmt.Errorf("unknown object id %d", id)
		}
		return nil, nil
	})
	v.RegisterForeign("SetObjectAlpha", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("SetObjectAlpha(id, value) requires 2 arguments")
		}
		id := toInt(args[0])
		alpha := toInt(args[1]) & 0xff
		objectsMu.Lock()
		obj, ok := objects[id]
		if ok {
			obj.colorA = uint8(alpha)
		}
		objectsMu.Unlock()
		if !ok {
			return nil, fmt.Errorf("unknown object id %d", id)
		}
		return nil, nil
	})
	v.RegisterForeign("SetObjectTexture", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("SetObjectTexture(id, textureID) requires 2 arguments")
		}
		id := toInt(args[0])
		texID := toInt(args[1])
		texturesMu.Lock()
		tex, texOk := textures[texID]
		texturesMu.Unlock()
		if !texOk {
			return nil, fmt.Errorf("unknown texture id %d", texID)
		}
		objectsMu.Lock()
		obj, ok := objects[id]
		if ok {
			obj.textureId = texID
			if obj.model.Materials != nil && obj.model.Materials.Maps != nil {
				obj.model.Materials.Maps.Texture = tex
			}
		}
		objectsMu.Unlock()
		if !ok {
			return nil, fmt.Errorf("unknown object id %d", id)
		}
		return nil, nil
	})
	v.RegisterForeign("SetObjectShader", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("SetObjectShader(id, shaderID) requires 2 arguments")
		}
		id := toInt(args[0])
		shaderID := toInt(args[1])
		objectsMu.Lock()
		obj, ok := objects[id]
		if ok {
			obj.shaderId = shaderID
		}
		objectsMu.Unlock()
		if !ok {
			return nil, fmt.Errorf("unknown object id %d", id)
		}
		// raylib DrawModelEx doesn't take shader; custom draw would use shader
		return nil, nil
	})
	v.RegisterForeign("SetObjectWireframe", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("SetObjectWireframe(id, onOff) requires 2 arguments")
		}
		id := toInt(args[0])
		onOff := toInt(args[1]) != 0
		objectsMu.Lock()
		obj, ok := objects[id]
		if ok {
			obj.wireframe = onOff
		}
		objectsMu.Unlock()
		if !ok {
			return nil, fmt.Errorf("unknown object id %d", id)
		}
		return nil, nil
	})
	v.RegisterForeign("SetObjectCollision", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("SetObjectCollision(id, onOff) requires 2 arguments")
		}
		id := toInt(args[0])
		onOff := toInt(args[1]) != 0
		objectsMu.Lock()
		obj, ok := objects[id]
		if ok {
			obj.collision = onOff
		}
		objectsMu.Unlock()
		if !ok {
			return nil, fmt.Errorf("unknown object id %d", id)
		}
		return nil, nil
	})

	// --- Scene ---
	v.RegisterForeign("Start3D", func(args []interface{}) (interface{}, error) {
		_, err := v.CallForeign("BeginMode3D", nil)
		return nil, err
	})
	v.RegisterForeign("End3D", func(args []interface{}) (interface{}, error) {
		_, err := v.CallForeign("EndMode3D", nil)
		return nil, err
	})
	// DrawGrid(size, spacing): DBP-style - calls rl.DrawGrid (slices, spacing)
	v.RegisterForeign("DrawGrid", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("DrawGrid(size, spacing) requires 2 arguments")
		}
		slices := int32(toInt(args[0]))
		spacing := toFloat32(args[1])
		rl.DrawGrid(slices, spacing)
		return nil, nil
	})
	v.RegisterForeign("Clear", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			rl.ClearBackground(rl.Black)
			return nil, nil
		}
		r, g, b := toInt(args[0])&0xff, toInt(args[1])&0xff, toInt(args[2])&0xff
		rl.ClearBackground(rl.NewColor(uint8(r), uint8(g), uint8(b), 255))
		return nil, nil
	})
	v.RegisterForeign("BackgroundColor", func(args []interface{}) (interface{}, error) {
		return v.CallForeign("Clear", args)
	})

	// --- Input aliases ---
	v.RegisterForeign("KeyDown", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("KeyDown(key) requires 1 argument")
		}
		return v.CallForeign("IsKeyDown", args)
	})
	v.RegisterForeign("KeyHit", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("KeyHit(key) requires 1 argument")
		}
		return v.CallForeign("IsKeyPressed", args)
	})
	v.RegisterForeign("KeyUp", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("KeyUp(key) requires 1 argument")
		}
		return v.CallForeign("IsKeyUp", args)
	})
	v.RegisterForeign("MouseX", func(args []interface{}) (interface{}, error) {
		return v.CallForeign("GetMouseX", nil)
	})
	v.RegisterForeign("MouseY", func(args []interface{}) (interface{}, error) {
		return v.CallForeign("GetMouseY", nil)
	})
	v.RegisterForeign("MouseMoveX", func(args []interface{}) (interface{}, error) {
		return v.CallForeign("GetMouseDeltaX", nil)
	})
	v.RegisterForeign("MouseMoveY", func(args []interface{}) (interface{}, error) {
		return v.CallForeign("GetMouseDeltaY", nil)
	})
	v.RegisterForeign("MouseButtonDown", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("MouseButtonDown(button) requires 1 argument")
		}
		return v.CallForeign("IsMouseButtonDown", args)
	})
	v.RegisterForeign("MouseButtonHit", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("MouseButtonHit(button) requires 1 argument")
		}
		return v.CallForeign("IsMouseButtonPressed", args)
	})
	v.RegisterForeign("MouseButtonUp", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("MouseButtonUp(button) requires 1 argument")
		}
		return v.CallForeign("IsMouseButtonReleased", args)
	})
	v.RegisterForeign("HideMouse", func(args []interface{}) (interface{}, error) {
		_, err := v.CallForeign("DisableCursor", nil)
		return nil, err
	})
	v.RegisterForeign("ShowMouse", func(args []interface{}) (interface{}, error) {
		_, err := v.CallForeign("EnableCursor", nil)
		return nil, err
	})
	v.RegisterForeign("LockMouse", func(args []interface{}) (interface{}, error) {
		_, err := v.CallForeign("DisableCursor", nil)
		return nil, err
	})
	v.RegisterForeign("UnlockMouse", func(args []interface{}) (interface{}, error) {
		_, err := v.CallForeign("EnableCursor", nil)
		return nil, err
	})

	// --- Time ---
	v.RegisterForeign("DeltaTime", func(args []interface{}) (interface{}, error) {
		return v.CallForeign("GetFrameTime", nil)
	})
	v.RegisterForeign("FPS", func(args []interface{}) (interface{}, error) {
		return v.CallForeign("GetFPS", nil)
	})

	// --- Math ---
	v.RegisterForeign("Clamp", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("Clamp(value, min, max) requires 3 arguments")
		}
		return v.CallForeign("Clamp", args)
	})
	v.RegisterForeign("Lerp", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("Lerp(a, b, t) requires 3 arguments")
		}
		return v.CallForeign("Lerp", args)
	})
	v.RegisterForeign("RandomRange", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("RandomRange(min, max) requires 2 arguments")
		}
		return v.CallForeign("GetRandomValue", args)
	})

	// --- 2D Drawing (DBP-style aliases) ---
	v.RegisterForeign("DrawRect", func(args []interface{}) (interface{}, error) {
		if len(args) < 7 {
			return nil, fmt.Errorf("DrawRect(x, y, w, h, r, g, b) requires 7 arguments")
		}
		x, y, w, h := toInt(args[0]), toInt(args[1]), toInt(args[2]), toInt(args[3])
		r, g, b := toInt(args[4])&0xff, toInt(args[5])&0xff, toInt(args[6])&0xff
		rl.DrawRectangle(int32(x), int32(y), int32(w), int32(h), rl.NewColor(uint8(r), uint8(g), uint8(b), 255))
		return nil, nil
	})
	v.RegisterForeign("DrawCircle", func(args []interface{}) (interface{}, error) {
		if len(args) < 7 {
			return nil, fmt.Errorf("DrawCircle(x, y, radius, r, g, b) requires 7 arguments")
		}
		x, y := int32(toInt(args[0])), int32(toInt(args[1]))
		radius := toFloat32(args[2])
		r, g, b := toInt(args[3])&0xff, toInt(args[4])&0xff, toInt(args[5])&0xff
		rl.DrawCircle(x, y, radius, rl.NewColor(uint8(r), uint8(g), uint8(b), 255))
		return nil, nil
	})
	v.RegisterForeign("DrawLine", func(args []interface{}) (interface{}, error) {
		if len(args) < 7 {
			return nil, fmt.Errorf("DrawLine(x1, y1, x2, y2, r, g, b) requires 7 arguments")
		}
		x1, y1, x2, y2 := int32(toInt(args[0])), int32(toInt(args[1])), int32(toInt(args[2])), int32(toInt(args[3]))
		r, g, b := toInt(args[4])&0xff, toInt(args[5])&0xff, toInt(args[6])&0xff
		rl.DrawLine(x1, y1, x2, y2, rl.NewColor(uint8(r), uint8(g), uint8(b), 255))
		return nil, nil
	})
	v.RegisterForeign("DrawSprite", func(args []interface{}) (interface{}, error) {
		return v.CallForeign("Sprite", args)
	})
	v.RegisterForeign("LoadSprite", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("LoadSprite(id, path) requires 2 arguments")
		}
		return v.CallForeign("LoadImage", []interface{}{args[1], args[0]})
	})

	// --- Text (DBP-style) ---
	v.RegisterForeign("DrawText", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("DrawText(text$, x, y, size) or DrawText(text$, x, y, size, r, g, b)")
		}
		text := toString(args[0])
		x, y, size := int32(toInt(args[1])), int32(toInt(args[2])), int32(toInt(args[3]))
		c := rl.White
		if len(args) >= 7 {
			r, g, b := toInt(args[4])&0xff, toInt(args[5])&0xff, toInt(args[6])&0xff
			c = rl.NewColor(uint8(r), uint8(g), uint8(b), 255)
		}
		fontsMu.Lock()
		fid := currentFontID
		fontsMu.Unlock()
		if fid >= 0 {
			fontsMu.Lock()
			f, ok := fonts[fid]
			fontsMu.Unlock()
			if ok {
				rl.DrawTextEx(f, text, rl.Vector2{X: float32(x), Y: float32(y)}, float32(size), 1, c)
				return nil, nil
			}
		}
		rl.DrawText(text, x, y, size, c)
		return nil, nil
	})
	v.RegisterForeign("LoadFont", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("LoadFont(id, path) requires 2 arguments")
		}
		id := toInt(args[0])
		path := toString(args[1])
		f := rl.LoadFont(path)
		fontsMu.Lock()
		fonts[id] = f
		fontsMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("SetFont", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("SetFont(id) requires 1 argument")
		}
		fontsMu.Lock()
		currentFontID = toInt(args[0])
		fontsMu.Unlock()
		return nil, nil
	})

	// --- File I/O (DBP-style) ---
	v.RegisterForeign("FileExists", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("FileExists(path) requires 1 argument")
		}
		_, err := os.Stat(toString(args[0]))
		return err == nil, nil
	})
	v.RegisterForeign("AppendFile", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("AppendFile(path, text) requires 2 arguments")
		}
		f, err := os.OpenFile(toString(args[0]), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return nil, err
		}
		defer f.Close()
		_, err = f.WriteString(toString(args[1]))
		return nil, err
	})

	// --- UI (DBP-style aliases to Gui*) ---
	v.RegisterForeign("UIButton", func(args []interface{}) (interface{}, error) {
		if len(args) < 6 {
			return nil, fmt.Errorf("UIButton(id, x, y, w, h, text$) requires 6 arguments")
		}
		return v.CallForeign("GuiButton", []interface{}{args[1], args[2], args[3], args[4], args[5]})
	})
	v.RegisterForeign("UILabel", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("UILabel(x, y, text$) requires 3 arguments")
		}
		return v.CallForeign("GuiLabel", []interface{}{args[0], args[1], 200, 20, args[2]})
	})
	v.RegisterForeign("UICheckbox", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("UICheckbox(id, x, y, text$) requires 4 arguments")
		}
		return v.CallForeign("GuiCheckBox", []interface{}{args[1], args[2], 20, 20, args[3], false})
	})
	v.RegisterForeign("UISlider", func(args []interface{}) (interface{}, error) {
		if len(args) < 7 {
			return nil, fmt.Errorf("UISlider(id, x, y, w, min, max, value) requires 7 arguments")
		}
		return v.CallForeign("GuiSlider", []interface{}{args[1], args[2], args[3], args[4], args[5], args[6]})
	})

	// --- Debugging ---
	v.RegisterForeign("DebugLog", func(args []interface{}) (interface{}, error) {
		if len(args) >= 1 {
			log.Println(toString(args[0]))
		}
		return nil, nil
	})
	v.RegisterForeign("DebugDrawLine", func(args []interface{}) (interface{}, error) {
		if len(args) < 6 {
			return nil, fmt.Errorf("DebugDrawLine(x1,y1,z1, x2,y2,z2) requires 6 arguments")
		}
		start := rl.Vector3{X: toFloat32(args[0]), Y: toFloat32(args[1]), Z: toFloat32(args[2])}
		end := rl.Vector3{X: toFloat32(args[3]), Y: toFloat32(args[4]), Z: toFloat32(args[5])}
		rl.DrawLine3D(start, end, rl.Red)
		return nil, nil
	})
	v.RegisterForeign("DebugDrawBox", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("DebugDrawBox(x, y, z, size) requires 4 arguments")
		}
		pos := rl.Vector3{X: toFloat32(args[0]), Y: toFloat32(args[1]), Z: toFloat32(args[2])}
		size := toFloat32(args[3])
		rl.DrawCubeWires(pos, size, size, size, rl.Red)
		return nil, nil
	})

	// --- Random / Utility ---
	v.RegisterForeign("Randomize", func(args []interface{}) (interface{}, error) {
		seed := time.Now().UnixNano()
		if len(args) >= 1 {
			s := toInt(args[0])
			if s != 0 {
				seed = int64(s)
			}
		}
		rand.Seed(seed)
		return nil, nil
	})
	v.RegisterForeign("RandomMinMax", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("RandomMinMax(a, b) requires 2 arguments")
		}
		a, b := toFloat32(args[0]), toFloat32(args[1])
		return float64(a + (b-a)*float32(rand.Float64())), nil
	})
	v.RegisterForeign("Distance", func(args []interface{}) (interface{}, error) {
		if len(args) < 6 {
			return nil, fmt.Errorf("Distance(x1,y1,z1, x2,y2,z2) requires 6 arguments")
		}
		x1, y1, z1 := toFloat32(args[0]), toFloat32(args[1]), toFloat32(args[2])
		x2, y2, z2 := toFloat32(args[3]), toFloat32(args[4]), toFloat32(args[5])
		dx, dy, dz := x2-x1, y2-y1, z2-z1
		return float64(float32(math.Sqrt(float64(dx*dx + dy*dy + dz*dz)))), nil
	})
	v.RegisterForeign("AngleBetween", func(args []interface{}) (interface{}, error) {
		if len(args) < 6 {
			return nil, fmt.Errorf("AngleBetween(x1,y1,z1, x2,y2,z2) requires 6 arguments")
		}
		x1, y1, z1 := toFloat32(args[0]), toFloat32(args[1]), toFloat32(args[2])
		x2, y2, z2 := toFloat32(args[3]), toFloat32(args[4]), toFloat32(args[5])
		dot := x1*x2 + y1*y2 + z1*z2
		len1 := float32(math.Sqrt(float64(x1*x1 + y1*y1 + z1*z1)))
		len2 := float32(math.Sqrt(float64(x2*x2 + y2*y2 + z2*z2)))
		if len1 < 1e-6 || len2 < 1e-6 {
			return 0.0, nil
		}
		cos := dot / (len1 * len2)
		if cos > 1 {
			cos = 1
		}
		if cos < -1 {
			cos = -1
		}
		return float64(math.Acos(float64(cos))), nil
	})
	v.RegisterForeign("Dot", func(args []interface{}) (interface{}, error) {
		if len(args) < 6 {
			return nil, fmt.Errorf("Dot(x1,y1,z1, x2,y2,z2) requires 6 arguments")
		}
		x1, y1, z1 := toFloat32(args[0]), toFloat32(args[1]), toFloat32(args[2])
		x2, y2, z2 := toFloat32(args[3]), toFloat32(args[4]), toFloat32(args[5])
		return float64(x1*x2 + y1*y2 + z1*z2), nil
	})

	// --- Game loop ---
	v.RegisterForeign("StartDraw", func(args []interface{}) (interface{}, error) {
		_, err := v.CallForeign("BeginDrawing", nil)
		return nil, err
	})
	v.RegisterForeign("EndDraw", func(args []interface{}) (interface{}, error) {
		_, err := v.CallForeign("EndDrawing", nil)
		return nil, err
	})

	// --- Audio (DBP-style int IDs) ---
	v.RegisterForeign("LoadSound", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("LoadSound(id, path) requires 2 arguments")
		}
		id := toInt(args[0])
		path := toString(args[1])
		snd := rl.LoadSound(path)
		soundsMu.Lock()
		sounds[id] = snd
		soundsMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("PlaySound", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("PlaySound(id) requires 1 argument")
		}
		id := toInt(args[0])
		soundsMu.Lock()
		snd, ok := sounds[id]
		soundsMu.Unlock()
		if !ok {
			return nil, fmt.Errorf("unknown sound id %d", id)
		}
		rl.PlaySound(snd)
		return nil, nil
	})
	v.RegisterForeign("StopSound", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("StopSound(id) requires 1 argument")
		}
		id := toInt(args[0])
		soundsMu.Lock()
		snd, ok := sounds[id]
		soundsMu.Unlock()
		if !ok {
			return nil, fmt.Errorf("unknown sound id %d", id)
		}
		rl.StopSound(snd)
		return nil, nil
	})
	v.RegisterForeign("SetSoundVolume", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("SetSoundVolume(id, value) requires 2 arguments")
		}
		id := toInt(args[0])
		vol := toFloat32(args[1])
		soundsMu.Lock()
		snd, ok := sounds[id]
		soundsMu.Unlock()
		if !ok {
			return nil, fmt.Errorf("unknown sound id %d", id)
		}
		rl.SetSoundVolume(snd, vol)
		return nil, nil
	})

	// --- FPS Camera ---
	v.RegisterForeign("FpsCameraOn", func(args []interface{}) (interface{}, error) {
		fpsCameraMu.Lock()
		fpsCameraOn = true
		fpsCameraMu.Unlock()
		_, err := v.CallForeign("DisableCursor", nil)
		return nil, err
	})
	v.RegisterForeign("FpsCameraOff", func(args []interface{}) (interface{}, error) {
		fpsCameraMu.Lock()
		fpsCameraOn = false
		fpsCameraMu.Unlock()
		_, err := v.CallForeign("EnableCursor", nil)
		return nil, err
	})
	v.RegisterForeign("FpsCameraPosition", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("FpsCameraPosition(x, y, z) requires 3 arguments")
		}
		x, y, z := toFloat32(args[0]), toFloat32(args[1]), toFloat32(args[2])
		fpsCameraMu.Lock()
		fpsCamX, fpsCamY, fpsCamZ = x, y, z
		fpsCameraMu.Unlock()
		_, err := v.CallForeign("SetCameraPosition", []interface{}{x, y, z})
		return nil, err
	})
	v.RegisterForeign("FpsMoveSpeed", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("FpsMoveSpeed(value) requires 1 argument")
		}
		fpsCameraMu.Lock()
		fpsMoveSpeed = toFloat32(args[0])
		fpsCameraMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("FpsLookSpeed", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("FpsLookSpeed(value) requires 1 argument")
		}
		fpsCameraMu.Lock()
		fpsLookSpeed = toFloat32(args[0])
		fpsCameraMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("FpsUpdate", func(args []interface{}) (interface{}, error) {
		dx, _ := v.CallForeign("GetMouseDeltaX", nil)
		dy, _ := v.CallForeign("GetMouseDeltaY", nil)
		dt, _ := v.CallForeign("GetFrameTime", nil)
		dtf := float32(0.016)
		if dt != nil {
			if d, ok := dt.(float64); ok {
				dtf = float32(d)
			}
		}
		dxf, dyf := float32(0), float32(0)
		if dx != nil {
			if d, ok := dx.(float64); ok {
				dxf = float32(d)
			}
		}
		if dy != nil {
			if d, ok := dy.(float64); ok {
				dyf = float32(d)
			}
		}
		kw, _ := v.CallForeign("IsKeyDown", []interface{}{int(rl.KeyW)})
		ks, _ := v.CallForeign("IsKeyDown", []interface{}{int(rl.KeyS)})
		ka, _ := v.CallForeign("IsKeyDown", []interface{}{int(rl.KeyA)})
		kd, _ := v.CallForeign("IsKeyDown", []interface{}{int(rl.KeyD)})
		fpsCameraMu.Lock()
		fpsYaw -= dxf * fpsLookSpeed
		fpsPitch -= dyf * fpsLookSpeed
		if fpsPitch > 1.5 {
			fpsPitch = 1.5
		}
		if fpsPitch < -1.5 {
			fpsPitch = -1.5
		}
		moveSpeed := fpsMoveSpeed * dtf
		cosYaw := float32(math.Cos(float64(fpsYaw)))
		sinYaw := float32(math.Sin(float64(fpsYaw)))
		if kw == true {
			fpsCamX -= sinYaw * moveSpeed
			fpsCamZ -= cosYaw * moveSpeed
		}
		if ks == true {
			fpsCamX += sinYaw * moveSpeed
			fpsCamZ += cosYaw * moveSpeed
		}
		if ka == true {
			fpsCamX -= cosYaw * moveSpeed
			fpsCamZ += sinYaw * moveSpeed
		}
		if kd == true {
			fpsCamX += cosYaw * moveSpeed
			fpsCamZ -= sinYaw * moveSpeed
		}
		camX, camY, camZ := fpsCamX, fpsCamY, fpsCamZ
		yaw, pitch := fpsYaw, fpsPitch
		fpsCameraMu.Unlock()
		targetX := camX - float32(math.Sin(float64(yaw)))*float32(math.Cos(float64(pitch)))
		targetY := camY + float32(math.Sin(float64(pitch)))
		targetZ := camZ - float32(math.Cos(float64(yaw)))*float32(math.Cos(float64(pitch)))
		_, err := v.CallForeign("SetCameraPosition", []interface{}{camX, camY, camZ})
		if err != nil {
			return nil, err
		}
		_, err = v.CallForeign("SetCameraTarget", []interface{}{targetX, targetY, targetZ})
		return nil, err
	})

	// --- Frame sync and input ---
	// Sync: end frame and present (user calls BeginDrawing + draw calls before Sync)
	v.RegisterForeign("Sync", func(args []interface{}) (interface{}, error) {
		rl.EndDrawing()
		return nil, nil
	})
	v.RegisterForeign("SYNC", func(args []interface{}) (interface{}, error) {
		rl.EndDrawing()
		return nil, nil
	})
	v.RegisterForeign("EscapeKey", func(args []interface{}) (interface{}, error) {
		return rl.IsKeyDown(rl.KeyEscape), nil
	})
	v.RegisterForeign("ESCAPEKEY", func(args []interface{}) (interface{}, error) {
		return rl.IsKeyDown(rl.KeyEscape), nil
	})
}
