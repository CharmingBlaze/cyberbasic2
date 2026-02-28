package vegetation

import (
	"fmt"
	"math"
	"sync"
)

// TreeType holds model and texture refs for a tree type.
type TreeType struct {
	ModelID         string
	TrunkTextureID  string
	LeafTextureID   string
	WindStrength    float32
	WindSpeed       float32
}

// TreeInstance holds position, scale, rotation for one placed tree.
type TreeInstance struct {
	TypeID          string
	X, Y, Z         float32
	Scale           float32
	Rotation        float32
	CollisionEnabled bool
	CollisionRadius  float32
}

var (
	treeTypes     = make(map[string]*TreeType)
	treeTypesSeq  int
	treeTypesMu   sync.Mutex
	treeSystems   = make(map[string][]string) // systemId -> list of treeInstanceIds
	treeInstances = make(map[string]*TreeInstance)
	treeInstancesSeq int
	treeInstancesMu  sync.Mutex
	treeSystemSeq int
	treeSystemMu  sync.Mutex
	lodDistances  = make(map[string][3]float32) // systemId -> near, mid, far
	instancingOn  = make(map[string]bool)
)

func getTreeInstance(id string) *TreeInstance {
	treeInstancesMu.Lock()
	t := treeInstances[id]
	treeInstancesMu.Unlock()
	return t
}

// TreeTypeCreate stores a tree type; returns type id.
func TreeTypeCreate(modelID, trunkTextureID, leafTextureID string) string {
	treeTypesMu.Lock()
	treeTypesSeq++
	id := fmt.Sprintf("treetype_%d", treeTypesSeq)
	treeTypes[id] = &TreeType{
		ModelID:        modelID,
		TrunkTextureID: trunkTextureID,
		LeafTextureID:  leafTextureID,
	}
	treeTypesMu.Unlock()
	return id
}

// TreeSystemCreate creates a new tree system; returns system id.
func TreeSystemCreate() string {
	treeSystemMu.Lock()
	treeSystemSeq++
	id := fmt.Sprintf("treesys_%d", treeSystemSeq)
	treeSystems[id] = nil
	treeSystemMu.Unlock()
	return id
}

// TreePlace adds a tree instance to a system. Returns tree instance id.
func TreePlace(systemID, typeID string, x, y, z, scale, rotation float32) (string, error) {
	treeTypesMu.Lock()
	_, okType := treeTypes[typeID]
	treeTypesMu.Unlock()
	if !okType {
		return "", fmt.Errorf("unknown tree type id: %s", typeID)
	}
	treeSystemMu.Lock()
	_, okSys := treeSystems[systemID]
	if !okSys {
		treeSystemMu.Unlock()
		return "", fmt.Errorf("unknown tree system id: %s", systemID)
	}
	treeInstancesMu.Lock()
	treeInstancesSeq++
	tid := fmt.Sprintf("tree_%d", treeInstancesSeq)
	treeInstances[tid] = &TreeInstance{TypeID: typeID, X: x, Y: y, Z: z, Scale: scale, Rotation: rotation}
	treeInstancesMu.Unlock()
	treeSystems[systemID] = append(treeSystems[systemID], tid)
	treeSystemMu.Unlock()
	return tid, nil
}

// TreeRemove removes a tree instance.
func TreeRemove(treeID string) error {
	treeInstancesMu.Lock()
	_, ok := treeInstances[treeID]
	delete(treeInstances, treeID)
	treeInstancesMu.Unlock()
	if !ok {
		return fmt.Errorf("unknown tree id: %s", treeID)
	}
	treeSystemMu.Lock()
	for sysID, list := range treeSystems {
		for i, id := range list {
			if id == treeID {
				treeSystems[sysID] = append(list[:i], list[i+1:]...)
				break
			}
		}
	}
	treeSystemMu.Unlock()
	return nil
}

// TreeSetPosition sets tree instance position.
func TreeSetPosition(treeID string, x, y, z float32) error {
	t := getTreeInstance(treeID)
	if t == nil {
		return fmt.Errorf("unknown tree id: %s", treeID)
	}
	t.X, t.Y, t.Z = x, y, z
	return nil
}

// TreeSetScale sets tree instance scale.
func TreeSetScale(treeID string, scale float32) error {
	t := getTreeInstance(treeID)
	if t == nil {
		return fmt.Errorf("unknown tree id: %s", treeID)
	}
	t.Scale = scale
	return nil
}

// TreeSetRotation sets tree instance rotation (radians).
func TreeSetRotation(treeID string, rotation float32) error {
	t := getTreeInstance(treeID)
	if t == nil {
		return fmt.Errorf("unknown tree id: %s", treeID)
	}
	t.Rotation = rotation
	return nil
}

// TreeSystemSetLOD sets LOD distances (near, mid, far) for a system.
func TreeSystemSetLOD(systemID string, near, mid, far float32) {
	treeSystemMu.Lock()
	lodDistances[systemID] = [3]float32{near, mid, far}
	treeSystemMu.Unlock()
}

// TreeSystemEnableInstancing enables or disables instancing for a system.
func TreeSystemEnableInstancing(systemID string, on bool) {
	treeSystemMu.Lock()
	instancingOn[systemID] = on
	treeSystemMu.Unlock()
}

// TreeGetAt returns the nearest tree id in the system to (x, z), or "" if none.
func TreeGetAt(systemID string, x, z float32) string {
	treeSystemMu.Lock()
	list := treeSystems[systemID]
	treeSystemMu.Unlock()
	if list == nil {
		return ""
	}
	var best string
	bestDistSq := float32(math.MaxFloat32)
	for _, tid := range list {
		t := getTreeInstance(tid)
		if t == nil {
			continue
		}
		dx := t.X - x
		dz := t.Z - z
		d2 := dx*dx + dz*dz
		if d2 < bestDistSq {
			bestDistSq = d2
			best = tid
		}
	}
	return best
}

// GetTreeSystemInstanceIds returns a copy of tree instance ids for a system (for DrawTrees).
func GetTreeSystemInstanceIds(systemID string) []string {
	treeSystemMu.Lock()
	list := treeSystems[systemID]
	if list == nil {
		treeSystemMu.Unlock()
		return nil
	}
	out := make([]string, len(list))
	copy(out, list)
	treeSystemMu.Unlock()
	return out
}

// GetTreeType returns the tree type by id.
func GetTreeType(typeID string) *TreeType {
	treeTypesMu.Lock()
	tt := treeTypes[typeID]
	treeTypesMu.Unlock()
	return tt
}

// TreeSetCollisionEnabled sets whether a tree instance has collision (capsule).
func TreeSetCollisionEnabled(treeInstanceID string, on bool) {
	treeInstancesMu.Lock()
	if t := treeInstances[treeInstanceID]; t != nil {
		t.CollisionEnabled = on
	}
	treeInstancesMu.Unlock()
}

// TreeSetCollisionRadius sets capsule radius for a tree instance.
func TreeSetCollisionRadius(treeInstanceID string, radius float32) {
	treeInstancesMu.Lock()
	if t := treeInstances[treeInstanceID]; t != nil {
		t.CollisionRadius = radius
	}
	treeInstancesMu.Unlock()
}

// TreeTypeSetWind sets wind strength and speed for a tree type (shader bending).
func TreeTypeSetWind(typeID string, strength, speed float32) {
	treeTypesMu.Lock()
	if tt := treeTypes[typeID]; tt != nil {
		tt.WindStrength = strength
		tt.WindSpeed = speed
	}
	treeTypesMu.Unlock()
}
