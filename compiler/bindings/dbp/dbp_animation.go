// Package dbp: Animation - LoadAnimation, PlayAnimation, SetAnimationFrame, GetAnimationFrame, GetAnimationLength, GetAnimationName.
package dbp

import (
	"fmt"
	"path/filepath"
	"sync"

	blendanim "cyberbasic/compiler/runtime/animation"
	"cyberbasic/compiler/vm"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type dbpAnimState struct {
	animId    int
	clipIndex int
	frame     float32
	speed     float32
	loop      bool
	blend     *blendanim.BlendState
}

type dbpMeshAnimState struct {
	frames  []rl.Model
	frame   float32
	speed   float32
	playing bool
	loop    bool
}

var (
	dbpAnims        = make(map[int][]rl.ModelAnimation) // animID -> all clips
	dbpAnimsMu      sync.Mutex
	objectAnimState = make(map[int]*dbpAnimState)
	objectAnimMu    sync.Mutex
	meshAnimFrames  = make(map[int][]rl.Model)
	meshAnimState   = make(map[int]*dbpMeshAnimState)
	meshAnimMu      sync.Mutex
)

// register3DAnimation adds LoadAnimation, PlayAnimation, SetAnimationFrame, GetAnimationFrame, GetAnimationLength, GetAnimationName.
func register3DAnimation(v *vm.VM) {
	v.RegisterForeign("LoadAnimation", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("LoadAnimation(id, path) requires 2 arguments")
		}
		id := toInt(args[0])
		path := toString(args[1])
		anims := rl.LoadModelAnimations(path)
		if len(anims) == 0 {
			return nil, nil // graceful: succeed with no animations; GetAnimationLength returns 0
		}
		dbpAnimsMu.Lock()
		dbpAnims[id] = anims
		dbpAnimsMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("PlayAnimation", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("PlayAnimation(objectID, animID, speed) or (objectID, animID, clipIndex, speed) requires 3-4 arguments")
		}
		objID := toInt(args[0])
		animID := toInt(args[1])
		clipIndex := 0
		speed := toFloat32(args[2])
		if len(args) >= 4 {
			clipIndex = toInt(args[2])
			speed = toFloat32(args[3])
		}
		dbpAnimsMu.Lock()
		clips, hasAnim := dbpAnims[animID]
		dbpAnimsMu.Unlock()
		if !hasAnim || len(clips) == 0 {
			return nil, nil // no-op when no animation
		}
		if clipIndex < 0 || clipIndex >= len(clips) {
			clipIndex = 0
		}
		objectAnimMu.Lock()
		objectAnimState[objID] = &dbpAnimState{animId: animID, clipIndex: clipIndex, frame: 0, speed: speed, loop: true}
		objectAnimMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("CrossfadeAnimation", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("CrossfadeAnimation(objectID, animID, clipIndex, duration) requires 4 arguments")
		}
		objID := toInt(args[0])
		animID := toInt(args[1])
		clipIndex := toInt(args[2])
		duration := toFloat32(args[3])
		if duration <= 0 {
			duration = 0.2
		}
		dbpAnimsMu.Lock()
		clips, hasAnim := dbpAnims[animID]
		dbpAnimsMu.Unlock()
		if !hasAnim || len(clips) == 0 {
			return nil, nil
		}
		if clipIndex < 0 || clipIndex >= len(clips) {
			clipIndex = 0
		}
		objectAnimMu.Lock()
		st, ok := objectAnimState[objID]
		if !ok || st == nil {
			objectAnimState[objID] = &dbpAnimState{animId: animID, clipIndex: clipIndex, frame: 0, speed: 1, loop: true}
			objectAnimMu.Unlock()
			return nil, nil
		}
		if st.animId != animID {
			objectAnimState[objID] = &dbpAnimState{animId: animID, clipIndex: clipIndex, frame: 0, speed: st.speed, loop: st.loop}
			objectAnimMu.Unlock()
			return nil, nil
		}
		fromClip := st.clipIndex
		fromFrame := st.frame
		if st.blend != nil {
			fromClip = st.blend.ToClipIndex
			fromFrame = st.blend.ToFrame
		}
		if fromClip == clipIndex {
			objectAnimMu.Unlock()
			return nil, nil
		}
		objectAnimState[objID] = &dbpAnimState{
			animId: animID, clipIndex: fromClip, frame: fromFrame, speed: st.speed, loop: st.loop,
			blend: &blendanim.BlendState{
				FromClipIndex: fromClip,
				ToClipIndex:   clipIndex,
				FromFrame:     fromFrame,
				ToFrame:       0,
				BlendDuration: duration,
				Active:        true,
			},
		}
		objectAnimMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("SetAnimationBlend", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("SetAnimationBlend(objectID, weight) requires 2 arguments")
		}
		objID := toInt(args[0])
		weight := toFloat32(args[1])
		if weight < 0 {
			weight = 0
		}
		if weight > 1 {
			weight = 1
		}
		objectAnimMu.Lock()
		if st, ok := objectAnimState[objID]; ok && st.blend != nil {
			st.blend.BlendElapsed = weight * st.blend.BlendDuration
			st.blend.BlendWeight = weight
			if weight >= 1 {
				st.clipIndex = st.blend.ToClipIndex
				st.frame = st.blend.ToFrame
				st.blend = nil
			}
		}
		objectAnimMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("StopAnimation", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("StopAnimation(objectID) requires 1 argument")
		}
		objID := toInt(args[0])
		objectAnimMu.Lock()
		delete(objectAnimState, objID)
		objectAnimMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("SetAnimationSpeed", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("SetAnimationSpeed(objectID, speed) requires 2 arguments")
		}
		objID := toInt(args[0])
		speed := toFloat32(args[1])
		objectAnimMu.Lock()
		if st, ok := objectAnimState[objID]; ok {
			st.speed = speed
		}
		objectAnimMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("SetAnimationLoop", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("SetAnimationLoop(objectID, onOff) requires 2 arguments")
		}
		objID := toInt(args[0])
		onOff := toInt(args[1]) != 0
		objectAnimMu.Lock()
		if st, ok := objectAnimState[objID]; ok {
			st.loop = onOff
		}
		objectAnimMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("ResetBones", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("ResetBones(objectID) requires 1 argument")
		}
		objID := toInt(args[0])
		objectAnimMu.Lock()
		st, ok := objectAnimState[objID]
		objectAnimMu.Unlock()
		if !ok || st.animId == 0 {
			return nil, nil
		}
		dbpAnimsMu.Lock()
		clips, ok := dbpAnims[st.animId]
		dbpAnimsMu.Unlock()
		if !ok || len(clips) == 0 {
			return nil, nil
		}
		ci := st.clipIndex
		if ci < 0 || ci >= len(clips) {
			ci = 0
		}
		anim := clips[ci]
		objectsMu.Lock()
		obj, ok := objects[objID]
		objectsMu.Unlock()
		if !ok || obj.model.MeshCount == 0 {
			return nil, nil
		}
		rl.UpdateModelAnimation(obj.model, anim, 0)
		return nil, nil
	})
	v.RegisterForeign("SetAnimationFrame", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("SetAnimationFrame(objectID, frame) requires 2 arguments")
		}
		objID := toInt(args[0])
		frame := toFloat32(args[1])
		var animID, clipIndex int
		objectAnimMu.Lock()
		if st, ok := objectAnimState[objID]; ok {
			st.frame = frame
			animID = st.animId
			clipIndex = st.clipIndex
		} else {
			objectAnimState[objID] = &dbpAnimState{frame: frame}
		}
		objectAnimMu.Unlock()
		if animID == 0 {
			return nil, nil
		}
		objectsMu.Lock()
		obj, ok := objects[objID]
		objectsMu.Unlock()
		if !ok || obj.model.MeshCount == 0 {
			return nil, nil
		}
		dbpAnimsMu.Lock()
		clips, ok := dbpAnims[animID]
		dbpAnimsMu.Unlock()
		if !ok || len(clips) == 0 {
			return nil, nil // graceful: no-op when no anim
		}
		if clipIndex < 0 || clipIndex >= len(clips) {
			clipIndex = 0
		}
		anim := clips[clipIndex]
		rl.UpdateModelAnimation(obj.model, anim, int32(frame))
		return nil, nil
	})
	v.RegisterForeign("GetAnimationFrame", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return 0.0, nil
		}
		objID := toInt(args[0])
		objectAnimMu.Lock()
		st, ok := objectAnimState[objID]
		objectAnimMu.Unlock()
		if !ok {
			return 0.0, nil
		}
		return float64(st.frame), nil
	})
	v.RegisterForeign("GetAnimationLength", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return 0, nil
		}
		animID := toInt(args[0])
		clipIndex := 0
		if len(args) >= 2 {
			clipIndex = toInt(args[1])
		}
		dbpAnimsMu.Lock()
		clips, ok := dbpAnims[animID]
		dbpAnimsMu.Unlock()
		if !ok || len(clips) == 0 {
			return 0, nil // graceful: return 0 when no anim
		}
		if clipIndex < 0 || clipIndex >= len(clips) {
			clipIndex = 0
		}
		return clips[clipIndex].FrameCount, nil
	})
	v.RegisterForeign("GetAnimationName", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return "", nil
		}
		animID := toInt(args[0])
		clipIndex := 0
		if len(args) >= 2 {
			clipIndex = toInt(args[1])
		}
		dbpAnimsMu.Lock()
		clips, ok := dbpAnims[animID]
		dbpAnimsMu.Unlock()
		if !ok || len(clips) == 0 {
			return "", nil
		}
		if clipIndex < 0 || clipIndex >= len(clips) {
			clipIndex = 0
		}
		return clips[clipIndex].Name, nil
	})
	// SetBoneRotation(id, boneName, pitch, yaw, roll): Manual bone control. Stored for future IK/skeletal use.
	v.RegisterForeign("SetBoneRotation", func(args []interface{}) (interface{}, error) {
		if len(args) < 5 {
			return nil, fmt.Errorf("SetBoneRotation(id, boneName, pitch, yaw, roll) requires 5 arguments")
		}
		// Placeholder: store for manual animation; full implementation in runtime/animation
		return nil, nil
	})
	// SetBonePosition(id, boneName, x, y, z): Manual bone control.
	v.RegisterForeign("SetBonePosition", func(args []interface{}) (interface{}, error) {
		if len(args) < 5 {
			return nil, fmt.Errorf("SetBonePosition(id, boneName, x, y, z) requires 5 arguments")
		}
		return nil, nil
	})
	// LoadMeshAnimation(id, folder$, frameCount): Load frame-by-frame mesh animation from folder.
	v.RegisterForeign("LoadMeshAnimation", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("LoadMeshAnimation(id, folder, frameCount) requires 3 arguments")
		}
		id := toInt(args[0])
		folder := toString(args[1])
		frameCount := toInt(args[2])
		if frameCount <= 0 {
			return nil, nil
		}
		var frames []rl.Model
		for i := 0; i < frameCount; i++ {
			path := filepath.Join(folder, fmt.Sprintf("%03d.obj", i+1))
			model := rl.LoadModel(path)
			if model.MeshCount == 0 {
				path = filepath.Join(folder, fmt.Sprintf("frame_%03d.obj", i+1))
				model = rl.LoadModel(path)
			}
			if model.MeshCount == 0 {
				path = filepath.Join(folder, fmt.Sprintf("%d.obj", i+1))
				model = rl.LoadModel(path)
			}
			frames = append(frames, model)
		}
		meshAnimMu.Lock()
		meshAnimFrames[id] = frames
		meshAnimState[id] = &dbpMeshAnimState{frames: frames, frame: 0, speed: 1, loop: true}
		meshAnimMu.Unlock()
		return nil, nil
	})
	// PlayMeshAnimation(id, speed): Play frame-by-frame animation.
	v.RegisterForeign("PlayMeshAnimation", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("PlayMeshAnimation(id, speed) requires 2 arguments")
		}
		id := toInt(args[0])
		speed := toFloat32(args[1])
		meshAnimMu.Lock()
		if st, ok := meshAnimState[id]; ok {
			st.playing = true
			st.speed = speed
		}
		meshAnimMu.Unlock()
		return nil, nil
	})
	// SetMeshAnimationFrame(id, frame): Set current mesh animation frame.
	v.RegisterForeign("SetMeshAnimationFrame", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("SetMeshAnimationFrame(id, frame) requires 2 arguments")
		}
		id := toInt(args[0])
		frame := toFloat32(args[1])
		meshAnimMu.Lock()
		if st, ok := meshAnimState[id]; ok {
			st.frame = frame
		}
		meshAnimMu.Unlock()
		return nil, nil
	})
}

// GetMeshAnimationModel returns the model to draw for object with mesh animation, or nil if none.
func GetMeshAnimationModel(objID int) *rl.Model {
	meshAnimMu.Lock()
	st, ok := meshAnimState[objID]
	meshAnimMu.Unlock()
	if !ok || st == nil || len(st.frames) == 0 {
		return nil
	}
	frameIdx := int(st.frame)
	if frameIdx < 0 {
		frameIdx = 0
	}
	if frameIdx >= len(st.frames) {
		frameIdx = len(st.frames) - 1
	}
	return &st.frames[frameIdx]
}

// UpdateMeshAnimation advances mesh animation for an object. Call from draw loop.
func UpdateMeshAnimation(objID int) {
	meshAnimMu.Lock()
	st, ok := meshAnimState[objID]
	if !ok || st == nil || !st.playing || len(st.frames) == 0 {
		meshAnimMu.Unlock()
		return
	}
	dt := rl.GetFrameTime()
	st.frame += st.speed * dt * 30
	fc := float32(len(st.frames))
	if st.loop {
		for st.frame >= fc {
			st.frame -= fc
		}
		for st.frame < 0 {
			st.frame += fc
		}
	} else {
		if st.frame >= fc {
			st.frame = fc - 1
			st.playing = false
		}
		if st.frame < 0 {
			st.frame = 0
			st.playing = false
		}
	}
	meshAnimMu.Unlock()
}

// UpdateObjectAnimation advances animation for an object and applies to model. Call from DrawObject.
func UpdateObjectAnimation(objID int, obj *dbpObject) {
	objectAnimMu.Lock()
	st, ok := objectAnimState[objID]
	objectAnimMu.Unlock()
	if !ok || st == nil || st.animId == 0 {
		return
	}
	dbpAnimsMu.Lock()
	clips, ok := dbpAnims[st.animId]
	dbpAnimsMu.Unlock()
	if !ok || len(clips) == 0 || obj.model.MeshCount == 0 {
		return
	}
	ci := st.clipIndex
	if ci < 0 || ci >= len(clips) {
		ci = 0
	}
	dt := rl.GetFrameTime()
	objectAnimMu.Lock()
	st = objectAnimState[objID]
	if st == nil {
		objectAnimMu.Unlock()
		return
	}
	if st.blend != nil {
		fromClip := st.blend.FromClipIndex
		toClip := st.blend.ToClipIndex
		if fromClip < 0 || fromClip >= len(clips) {
			fromClip = ci
		}
		if toClip < 0 || toClip >= len(clips) {
			toClip = ci
		}
		fromAnim := clips[fromClip]
		toAnim := clips[toClip]
		st.blend.FromFrame = advanceAnimationFrame(st.blend.FromFrame, fromAnim.FrameCount, dt, st.speed, st.loop)
		st.blend.ToFrame = advanceAnimationFrame(st.blend.ToFrame, toAnim.FrameCount, dt, st.speed, st.loop)
		if st.blend.Advance(dt) {
			st.clipIndex = st.blend.ToClipIndex
			st.frame = st.blend.ToFrame
			st.blend = nil
			objectAnimMu.Unlock()
			rl.UpdateModelAnimation(obj.model, toAnim, int32(st.frame))
			return
		}
		fromFrame := int32(st.blend.FromFrame)
		toFrame := int32(st.blend.ToFrame)
		weight := st.blend.BlendWeight
		objectAnimMu.Unlock()
		if obj.model.BoneCount > 0 && obj.model.BindPose != nil {
			applyBlendedAnimationPose(&obj.model, fromAnim, fromFrame, toAnim, toFrame, weight)
		} else {
			if weight < 0.5 {
				rl.UpdateModelAnimation(obj.model, fromAnim, fromFrame)
			} else {
				rl.UpdateModelAnimation(obj.model, toAnim, toFrame)
			}
		}
		return
	}
	anim := clips[ci]
	st.frame = advanceAnimationFrame(st.frame, anim.FrameCount, dt, st.speed, st.loop)
	frame := int32(st.frame)
	objectAnimMu.Unlock()
	rl.UpdateModelAnimation(obj.model, anim, frame)
}

func advanceAnimationFrame(frame float32, frameCount int32, dt, speed float32, loop bool) float32 {
	if frameCount <= 0 {
		return 0
	}
	frame += speed * dt * float32(frameCount) / 60.0
	fc := float32(frameCount)
	if loop {
		for frame >= fc {
			frame -= fc
		}
		for frame < 0 {
			frame += fc
		}
		return frame
	}
	if frame >= fc {
		return fc - 1
	}
	if frame < 0 {
		return 0
	}
	return frame
}

func applyBlendedAnimationPose(model *rl.Model, fromAnim rl.ModelAnimation, fromFrame int32, toAnim rl.ModelAnimation, toFrame int32, weight float32) {
	if model == nil || model.BindPose == nil || model.BoneCount <= 0 {
		return
	}
	if weight < 0 {
		weight = 0
	}
	if weight > 1 {
		weight = 1
	}
	poses := model.GetBindPose()
	for i := 0; i < int(model.BoneCount) && i < len(poses); i++ {
		fromPose := fromAnim.GetFramePose(int(fromFrame), i)
		toPose := toAnim.GetFramePose(int(toFrame), i)
		q1 := rl.Quaternion{X: fromPose.Rotation.X, Y: fromPose.Rotation.Y, Z: fromPose.Rotation.Z, W: fromPose.Rotation.W}
		q2 := rl.Quaternion{X: toPose.Rotation.X, Y: toPose.Rotation.Y, Z: toPose.Rotation.Z, W: toPose.Rotation.W}
		q := rl.QuaternionNormalize(rl.QuaternionSlerp(q1, q2, weight))
		poses[i] = rl.Transform{
			Translation: rl.Vector3Lerp(fromPose.Translation, toPose.Translation, weight),
			Rotation:    rl.Vector4{X: q.X, Y: q.Y, Z: q.Z, W: q.W},
			Scale:       rl.Vector3Lerp(fromPose.Scale, toPose.Scale, weight),
		}
	}
}
