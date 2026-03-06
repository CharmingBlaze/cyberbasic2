package runtime

import (
	"strings"
	"sync"

	gametime "cyberbasic/compiler/runtime/time"
)

var (
	fixedUpdateMu    sync.RWMutex
	fixedUpdateRate  float64 = 60
	fixedUpdateLabel string
)

// SetFixedUpdateRate configures the fixed update rate and keeps the time package in sync.
func SetFixedUpdateRate(rate float64) {
	if rate <= 0 {
		rate = 60
	}
	fixedUpdateMu.Lock()
	fixedUpdateRate = rate
	fixedUpdateMu.Unlock()
	gametime.SetFixedDeltaTime(1.0 / float32(rate))
}

// FixedUpdateRate returns the configured fixed update rate.
func FixedUpdateRate() float64 {
	fixedUpdateMu.RLock()
	defer fixedUpdateMu.RUnlock()
	return fixedUpdateRate
}

// SetFixedUpdateLabel sets the user callback invoked on each fixed step.
func SetFixedUpdateLabel(label string) {
	fixedUpdateMu.Lock()
	fixedUpdateLabel = strings.TrimSpace(label)
	fixedUpdateMu.Unlock()
}

// FixedUpdateLabel returns the user callback invoked on each fixed step.
func FixedUpdateLabel() string {
	fixedUpdateMu.RLock()
	defer fixedUpdateMu.RUnlock()
	return fixedUpdateLabel
}
