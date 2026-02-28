// Package raylib: 2D layer registry and commands (LayerCreate, LayerSetOrder, LayerSetVisible, etc.).
package raylib

import (
	"fmt"
	"sync"

	"cyberbasic/compiler/bindings/game"
	"cyberbasic/compiler/vm"
)

type layerState struct {
	Order     int
	Visible   bool
	ParallaxX float32
	ParallaxY float32
	ScrollX   float32
	ScrollY   float32
	// Sprites/tilemaps/particles reference this layer by ID; we don't store drawables here.
}

var (
	layers     = make(map[string]*layerState)
	layerSeq   int
	layerMu    sync.RWMutex
	defaultLayerID = "" // empty string = default layer, order 0
)

func init() {
	layerMu.Lock()
	layers[defaultLayerID] = &layerState{Order: 0, Visible: true, ParallaxX: 1, ParallaxY: 1}
	layerMu.Unlock()
}

// GetLayerOrder returns the draw order for a layer (0 for unknown or default).
func GetLayerOrder(layerID string) int {
	layerMu.RLock()
	defer layerMu.RUnlock()
	if layerID == "" {
		return 0
	}
	l, ok := layers[layerID]
	if !ok {
		return 0
	}
	return l.Order
}

// GetLayerVisible returns whether the layer is visible.
func GetLayerVisible(layerID string) bool {
	layerMu.RLock()
	defer layerMu.RUnlock()
	if layerID == "" {
		return true
	}
	l, ok := layers[layerID]
	if !ok {
		return true
	}
	return l.Visible
}

// GetLayerParallax returns parallax factors (px, py) for the layer.
func GetLayerParallax(layerID string) (float32, float32) {
	layerMu.RLock()
	defer layerMu.RUnlock()
	if layerID == "" {
		return 1, 1
	}
	l, ok := layers[layerID]
	if !ok {
		return 1, 1
	}
	return l.ParallaxX, l.ParallaxY
}

// GetLayerScroll returns scroll offset (sx, sy) for the layer.
func GetLayerScroll(layerID string) (float32, float32) {
	layerMu.RLock()
	defer layerMu.RUnlock()
	if layerID == "" {
		return 0, 0
	}
	l, ok := layers[layerID]
	if !ok {
		return 0, 0
	}
	return l.ScrollX, l.ScrollY
}

func registerLayers(v *vm.VM) {
	v.RegisterForeign("LayerCreate", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("LayerCreate requires (name, order)")
		}
		name := toString(args[0])
		order := int(toFloat32(args[1]))
		layerMu.Lock()
		layerSeq++
		id := name
		if id == "" {
			id = fmt.Sprintf("layer_%d", layerSeq)
		}
		if _, exists := layers[id]; exists {
			layerMu.Unlock()
			return nil, fmt.Errorf("layer already exists: %s", id)
		}
		layers[id] = &layerState{Order: order, Visible: true, ParallaxX: 1, ParallaxY: 1}
		layerMu.Unlock()
		return id, nil
	})
	v.RegisterForeign("LayerSetOrder", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("LayerSetOrder requires (layerID, order)")
		}
		id := toString(args[0])
		order := int(toFloat32(args[1]))
		layerMu.Lock()
		defer layerMu.Unlock()
		l, ok := layers[id]
		if !ok {
			return nil, fmt.Errorf("unknown layer: %s", id)
		}
		l.Order = order
		return nil, nil
	})
	v.RegisterForeign("LayerSetVisible", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("LayerSetVisible requires (layerID, flag)")
		}
		id := toString(args[0])
		visible := toFloat32(args[1]) != 0
		layerMu.Lock()
		defer layerMu.Unlock()
		l, ok := layers[id]
		if !ok {
			return nil, fmt.Errorf("unknown layer: %s", id)
		}
		l.Visible = visible
		return nil, nil
	})
	v.RegisterForeign("LayerSetParallax", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("LayerSetParallax requires (layerID, parallaxX, parallaxY)")
		}
		id := toString(args[0])
		px := toFloat32(args[1])
		py := toFloat32(args[2])
		layerMu.Lock()
		defer layerMu.Unlock()
		l, ok := layers[id]
		if !ok {
			return nil, fmt.Errorf("unknown layer: %s", id)
		}
		l.ParallaxX = px
		l.ParallaxY = py
		return nil, nil
	})
	v.RegisterForeign("LayerSetScroll", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("LayerSetScroll requires (layerID, scrollX, scrollY)")
		}
		id := toString(args[0])
		sx := toFloat32(args[1])
		sy := toFloat32(args[2])
		layerMu.Lock()
		defer layerMu.Unlock()
		l, ok := layers[id]
		if !ok {
			return nil, fmt.Errorf("unknown layer: %s", id)
		}
		l.ScrollX = sx
		l.ScrollY = sy
		return nil, nil
	})
	v.RegisterForeign("LayerClear", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("LayerClear requires (layerID)")
		}
		id := toString(args[0])
		clearLayerAssignments(id)
		return nil, nil
	})
	v.RegisterForeign("LayerSortSprites", func(args []interface{}) (interface{}, error) {
		// No-op: flush already sorts by (layerOrder, zIndex). Kept for API compatibility.
		return nil, nil
	})
}

// clearLayerAssignments removes layerID from all sprites, tilemaps, and particle systems that use it.
func clearLayerAssignments(layerID string) {
	clearSpriteLayer(layerID)
	game.ClearLayerAssignments(layerID)
}
