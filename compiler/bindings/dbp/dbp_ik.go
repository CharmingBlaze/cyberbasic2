// Package dbp: IK (Inverse Kinematics) - IKSolveTwoBone, IKEnable.
package dbp

import (
	"fmt"
	"math"
	"strings"
	"sync"

	"cyberbasic/compiler/bindings/model"
	"cyberbasic/compiler/vm"
)

var (
	objectModelMap = make(map[int]*model.Model)
	objectModelMu  sync.RWMutex
	ikEnabled      = make(map[int]bool)
	ikEnabledMu    sync.RWMutex
)

// RegisterObjectModel links an object ID to its source model (for IK/skeleton lookup).
// Called by LoadLevel and SpawnPrefab when creating objects.
func RegisterObjectModel(objectID int, m *model.Model) {
	objectModelMu.Lock()
	objectModelMap[objectID] = m
	objectModelMu.Unlock()
}

// UnregisterObjectModel removes the object-model link (e.g. on DeleteObject, UnloadLevel).
func UnregisterObjectModel(objectID int) {
	objectModelMu.Lock()
	delete(objectModelMap, objectID)
	objectModelMu.Unlock()
	ikEnabledMu.Lock()
	delete(ikEnabled, objectID)
	ikEnabledMu.Unlock()
}

func boneIndexByName(skel *model.Skeleton, name string) int {
	if skel == nil {
		return -1
	}
	n := strings.ToLower(strings.TrimSpace(name))
	for i, b := range skel.Bones {
		if strings.ToLower(b.Name) == n {
			return i
		}
	}
	return -1
}

// solveTwoBoneIK computes elbow/knee position so end effector reaches target.
// root, mid, end are joint positions; target is desired end position; pole is hint for elbow.
func solveTwoBoneIK(root, mid, end, target [3]float32, pole [3]float32) (newMid [3]float32, ok bool) {
	len1 := float32(math.Sqrt(float64((mid[0]-root[0])*(mid[0]-root[0]) + (mid[1]-root[1])*(mid[1]-root[1]) + (mid[2]-root[2])*(mid[2]-root[2]))))
	len2 := float32(math.Sqrt(float64((end[0]-mid[0])*(end[0]-mid[0]) + (end[1]-mid[1])*(end[1]-mid[1]) + (end[2]-mid[2])*(end[2]-mid[2]))))
	if len1 < 1e-6 || len2 < 1e-6 {
		return mid, false
	}
	dir := [3]float32{target[0] - root[0], target[1] - root[1], target[2] - root[2]}
	dist := float32(math.Sqrt(float64(dir[0]*dir[0] + dir[1]*dir[1] + dir[2]*dir[2])))
	if dist < 1e-6 {
		return mid, false
	}
	dir[0] /= dist
	dir[1] /= dist
	dir[2] /= dist
	if dist >= len1+len2 {
		newMid[0] = root[0] + dir[0]*len1
		newMid[1] = root[1] + dir[1]*len1
		newMid[2] = root[2] + dir[2]*len1
		return newMid, true
	}
	if dist <= len1-len2 || dist <= len2-len1 {
		return mid, false
	}
	cosAngle := (len1*len1 + dist*dist - len2*len2) / (2 * len1 * dist)
	if cosAngle > 1 {
		cosAngle = 1
	}
	if cosAngle < -1 {
		cosAngle = -1
	}
	angle := float32(math.Acos(float64(cosAngle)))
	midAlong := [3]float32{root[0] + dir[0]*len1, root[1] + dir[1]*len1, root[2] + dir[2]*len1}
	perp := [3]float32{
		midAlong[0] - root[0],
		midAlong[1] - root[1],
		midAlong[2] - root[2],
	}
	perpLen := float32(math.Sqrt(float64(perp[0]*perp[0] + perp[1]*perp[1] + perp[2]*perp[2])))
	if perpLen < 1e-6 {
		perp = pole
	} else {
		perp[0] /= perpLen
		perp[1] /= perpLen
		perp[2] /= perpLen
	}
	sinAngle := float32(math.Sin(float64(angle)))
	newMid[0] = midAlong[0] + perp[0]*len1*sinAngle
	newMid[1] = midAlong[1] + perp[1]*len1*sinAngle
	newMid[2] = midAlong[2] + perp[2]*len1*sinAngle
	return newMid, true
}

func registerIK(v *vm.VM) {
	v.RegisterForeign("IKEnable", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("IKEnable(objectID, onOff) requires 2 arguments")
		}
		objID := toInt(args[0])
		onOff := toInt(args[1]) != 0
		ikEnabledMu.Lock()
		ikEnabled[objID] = onOff
		ikEnabledMu.Unlock()
		return nil, nil
	})

	v.RegisterForeign("IKSolveTwoBone", func(args []interface{}) (interface{}, error) {
		if len(args) < 6 {
			return nil, fmt.Errorf("IKSolveTwoBone(objectID, boneA$, boneB$, targetX, targetY, targetZ) requires 6 arguments")
		}
		objID := toInt(args[0])
		boneA := toString(args[1])
		boneB := toString(args[2])
		targetX := toFloat32(args[3])
		targetY := toFloat32(args[4])
		targetZ := toFloat32(args[5])
		ikEnabledMu.RLock()
		enabled := ikEnabled[objID]
		ikEnabledMu.RUnlock()
		if !enabled {
			return nil, nil
		}
		objectModelMu.RLock()
		m := objectModelMap[objID]
		objectModelMu.RUnlock()
		if m == nil || m.Skeleton == nil {
			return nil, nil
		}
		skel := m.Skeleton
		idxA := boneIndexByName(skel, boneA)
		idxB := boneIndexByName(skel, boneB)
		if idxA < 0 || idxB < 0 {
			return nil, nil
		}
		// Use bone hierarchy for positions: root at bone A, mid at bone B, end at bone B's child or estimate
		root := [3]float32{float32(idxA) * 0.1, 1, 0}
		mid := [3]float32{float32(idxB) * 0.1, 0.5, 0}
		end := [3]float32{float32(idxB)*0.1 + 0.1, 0, 0}
		if idxB+1 < len(skel.Bones) {
			end = [3]float32{float32(idxB+1) * 0.1, -0.5, 0}
		}
		target := [3]float32{targetX, targetY, targetZ}
		pole := [3]float32{0, 1, 0}
		_, ok := solveTwoBoneIK(root, mid, end, target, pole)
		if !ok {
			return nil, nil
		}
		// Apply: would update bone transforms for skinned models. Current BuildModel
		// produces per-mesh models; apply is no-op until skinned rendering is integrated.
		return nil, nil
	})
}
