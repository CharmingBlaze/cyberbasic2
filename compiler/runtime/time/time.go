// Package time provides time state for the game engine: delta time, fixed step, time scale.
package time

import (
	"sync"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
	mu sync.RWMutex

	// DeltaTime is the scaled frame delta (GetFrameTime * TimeScale).
	deltaTime float32 = 0

	// FixedDeltaTime is the fixed timestep for physics (e.g. 1/60).
	FixedDeltaTime float32 = 1.0 / 60.0

	// Accumulator accumulates time for fixed-step physics.
	Accumulator float32 = 0

	// TimeScale multiplies delta time (1 = normal, 0.5 = slow motion, 0 = pause).
	TimeScale float32 = 1.0

	// FrameCounter is the number of frames since start.
	FrameCounter uint64 = 0
)

const maxFrameDelta float32 = 0.25

// Update advances time state. Call once per frame before physics/update.
// dt is typically from rl.GetFrameTime().
func Update(dt float32) {
	mu.Lock()
	defer mu.Unlock()
	if dt < 0 {
		dt = 0
	}
	if dt > maxFrameDelta {
		dt = maxFrameDelta
	}
	deltaTime = dt * TimeScale
	Accumulator += deltaTime
	FrameCounter++
}

// DeltaTime returns the scaled frame delta for this frame.
func DeltaTime() float32 {
	mu.RLock()
	defer mu.RUnlock()
	return deltaTime
}

// GetFixedDeltaTime returns the fixed timestep for physics.
func GetFixedDeltaTime() float32 {
	mu.RLock()
	defer mu.RUnlock()
	return FixedDeltaTime
}

// SetFixedDeltaTime sets the fixed timestep used by the runtime loop.
func SetFixedDeltaTime(value float32) {
	mu.Lock()
	defer mu.Unlock()
	if value <= 0 {
		value = 1.0 / 60.0
	}
	FixedDeltaTime = value
}

// GetAccumulator returns the current accumulator value.
func GetAccumulator() float32 {
	mu.RLock()
	defer mu.RUnlock()
	return Accumulator
}

// ConsumeAccumulator subtracts fixedStep from Accumulator. Call after each physics step.
func ConsumeAccumulator(fixedStep float32) {
	mu.Lock()
	defer mu.Unlock()
	Accumulator -= fixedStep
	if Accumulator < 0 {
		Accumulator = 0
	}
}

// ClampAccumulator caps the stored accumulator to avoid runaway catch-up loops.
func ClampAccumulator(max float32) {
	mu.Lock()
	defer mu.Unlock()
	if max < 0 {
		max = 0
	}
	if Accumulator > max {
		Accumulator = max
	}
}

// SetTimeScale sets the time scale multiplier.
func SetTimeScale(value float32) {
	mu.Lock()
	defer mu.Unlock()
	if value < 0 {
		value = 0
	}
	TimeScale = value
}

// GetTimeScale returns the current time scale.
func GetTimeScale() float32 {
	mu.RLock()
	defer mu.RUnlock()
	return TimeScale
}

// GetFrameCounter returns the frame count.
func GetFrameCounter() uint64 {
	mu.RLock()
	defer mu.RUnlock()
	return FrameCounter
}

// GetRawDeltaTime returns raw frame time from raylib (unscaled).
func GetRawDeltaTime() float32 {
	return rl.GetFrameTime()
}
