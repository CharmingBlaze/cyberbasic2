// Package animation provides animation types: skeletal, mesh frames, manual.
package animation

import (
	"sync"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// AnimationType distinguishes skeletal, mesh-frame, and manual animation.
type AnimationType int

const (
	AnimationSkeletal  AnimationType = iota
	AnimationMeshFrames
	AnimationManual
)

// AnimationComponent holds animation state for an object.
type AnimationComponent struct {
	Type        AnimationType
	Skeleton    *Skeleton
	Clips       []AnimationClip
	CurrentClip *AnimationClip
	FrameIndex  int
	Speed       float32
	MeshFrames  []rl.Model
}

// Skeleton holds bone hierarchy.
type Skeleton struct {
	Bones []Bone
}

// Bone represents a single bone with transform.
type Bone struct {
	Name     string
	Parent   int
	BindPose rl.Matrix
	Local    rl.Matrix
	Global   rl.Matrix
}

// AnimationClip holds keyframe data.
type AnimationClip struct {
	Name      string
	Duration  float32
	Keyframes []Keyframe
}

// Keyframe holds time and bone transforms.
type Keyframe struct {
	Time   float32
	Bones  []rl.Matrix
}

var (
	components   = make(map[int]*AnimationComponent)
	componentsMu sync.RWMutex
)

// SetComponent stores animation component for object id.
func SetComponent(id int, c *AnimationComponent) {
	componentsMu.Lock()
	defer componentsMu.Unlock()
	components[id] = c
}

// GetComponent returns animation component for object id.
func GetComponent(id int) *AnimationComponent {
	componentsMu.RLock()
	defer componentsMu.RUnlock()
	return components[id]
}

// SolveTwoBone solves two-bone IK. Placeholder for full implementation.
func SolveTwoBone(skel *Skeleton, upper, lower string, target rl.Vector3) {
	// Full implementation would compute joint angles for target.
	_ = skel
	_ = upper
	_ = lower
	_ = target
}
