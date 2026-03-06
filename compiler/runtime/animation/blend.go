// Package animation provides blend state for crossfading between animation clips.
package animation

// BlendState holds the state for a crossfade between two clips.
type BlendState struct {
	FromClipIndex int
	ToClipIndex   int
	FromFrame     float32
	ToFrame       float32
	BlendWeight   float32 // 0 = from, 1 = to
	BlendDuration float32
	BlendElapsed  float32
	Active        bool
}

// Advance advances the blend by dt. Returns true when blend is complete.
func (b *BlendState) Advance(dt float32) bool {
	if !b.Active || b.BlendDuration <= 0 {
		b.Active = false
		return true
	}
	b.BlendElapsed += dt
	if b.BlendElapsed >= b.BlendDuration {
		b.BlendWeight = 1
		b.Active = false
		return true
	}
	b.BlendWeight = b.BlendElapsed / b.BlendDuration
	return false
}
