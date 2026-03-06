// Package ui provides UIManager for layout and hit testing. Wraps raygui.
package ui

import (
	"sync"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Element represents a UI element for layout.
type Element struct {
	ID       string
	X        float32
	Y        float32
	Width    float32
	Height   float32
	Visible  bool
	HitTest  func(x, y float32) bool
}

var (
	elements   = make(map[string]*Element)
	elementsMu sync.RWMutex
)

// Register adds an element to the UI manager.
func Register(e *Element) {
	elementsMu.Lock()
	defer elementsMu.Unlock()
	elements[e.ID] = e
}

// Unregister removes an element.
func Unregister(id string) {
	elementsMu.Lock()
	defer elementsMu.Unlock()
	delete(elements, id)
}

// HitTestAt returns the element at (x, y), or "" if none.
func HitTestAt(x, y float32) string {
	elementsMu.RLock()
	defer elementsMu.RUnlock()
	for _, e := range elements {
		if e.Visible && e.HitTest != nil && e.HitTest(x, y) {
			return e.ID
		}
		if e.Visible && e.HitTest == nil {
			if x >= e.X && x <= e.X+e.Width && y >= e.Y && y <= e.Y+e.Height {
				return e.ID
			}
		}
	}
	return ""
}

// GetMouseElement returns the element under the mouse.
func GetMouseElement() string {
	m := rl.GetMousePosition()
	return HitTestAt(m.X, m.Y)
}
