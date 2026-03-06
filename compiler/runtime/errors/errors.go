// Package errors provides global error state for LASTERROR$ and safe fallbacks.
package errors

import (
	"sync"
)

var (
	mu    sync.RWMutex
	last  string
)

// SetLastError sets the last error message.
func SetLastError(msg string) {
	mu.Lock()
	defer mu.Unlock()
	last = msg
}

// LastError returns the last error message.
func LastError() string {
	mu.RLock()
	defer mu.RUnlock()
	return last
}

// ClearError clears the last error.
func ClearError() {
	mu.Lock()
	defer mu.Unlock()
	last = ""
}
