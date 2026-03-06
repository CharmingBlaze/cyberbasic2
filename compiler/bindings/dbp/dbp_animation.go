// Package dbp: Animation - LoadAnimation, PlayAnimation, SetAnimationFrame, GetAnimationFrame, GetAnimationLength, GetAnimationName.
package dbp

import (
	"fmt"
	"sync"

	"cyberbasic/compiler/vm"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type dbpAnimState struct {
	animId int
	frame  float32
	speed  float32
}

var (
	dbpAnims        = make(map[int]rl.ModelAnimation)
	dbpAnimsMu      sync.Mutex
	objectAnimState = make(map[int]*dbpAnimState)
	objectAnimMu    sync.Mutex
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
		dbpAnims[id] = anims[0]
		dbpAnimsMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("PlayAnimation", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("PlayAnimation(objectID, animID, speed) requires 3 arguments")
		}
		objID := toInt(args[0])
		animID := toInt(args[1])
		speed := toFloat32(args[2])
		dbpAnimsMu.Lock()
		_, hasAnim := dbpAnims[animID]
		dbpAnimsMu.Unlock()
		if !hasAnim {
			return nil, nil // no-op when no animation
		}
		objectAnimMu.Lock()
		objectAnimState[objID] = &dbpAnimState{animId: animID, frame: 0, speed: speed}
		objectAnimMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("SetAnimationFrame", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("SetAnimationFrame(objectID, frame) requires 2 arguments")
		}
		objID := toInt(args[0])
		frame := toFloat32(args[1])
		var animID int
		objectAnimMu.Lock()
		if st, ok := objectAnimState[objID]; ok {
			st.frame = frame
			animID = st.animId
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
		anim, ok := dbpAnims[animID]
		dbpAnimsMu.Unlock()
		if !ok {
			return nil, nil // graceful: no-op when no anim
		}
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
		dbpAnimsMu.Lock()
		anim, ok := dbpAnims[animID]
		dbpAnimsMu.Unlock()
		if !ok {
			return 0, nil // graceful: return 0 when no anim
		}
		return anim.FrameCount, nil
	})
	v.RegisterForeign("GetAnimationName", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return "", nil
		}
		animID := toInt(args[0])
		dbpAnimsMu.Lock()
		anim, ok := dbpAnims[animID]
		dbpAnimsMu.Unlock()
		if !ok {
			return "", nil
		}
		return anim.Name, nil
	})
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
	anim, ok := dbpAnims[st.animId]
	dbpAnimsMu.Unlock()
	if !ok || obj.model.MeshCount == 0 {
		return
	}
	dt := rl.GetFrameTime()
	objectAnimMu.Lock()
	st.frame += st.speed * dt * float32(anim.FrameCount) / 60.0
	fc := float32(anim.FrameCount)
	if fc > 0 {
		for st.frame >= fc {
			st.frame -= fc
		}
		for st.frame < 0 {
			st.frame += fc
		}
	}
	frame := int32(st.frame)
	objectAnimMu.Unlock()
	rl.UpdateModelAnimation(obj.model, anim, frame)
}
