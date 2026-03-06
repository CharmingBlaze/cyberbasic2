// Package aseprite parses Aseprite JSON export format for sprite sheets.
package aseprite

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
)

// Sheet represents a parsed Aseprite sprite sheet.
type Sheet struct {
	Frames []Frame
	Tags   map[string]FrameTag
	Slices map[string]Slice
}

// Frame holds one frame's rect and duration.
type Frame struct {
	X, Y, W, H int
	DurationMs int
}

// FrameTag defines an animation (name, from-to range, direction).
type FrameTag struct {
	Name      string
	From      int
	To        int
	Direction string // "forward", "reverse", "pingpong"
}

// Slice holds per-frame bounds (e.g. hitboxes).
type Slice struct {
	Keys map[int]SliceKey // frame index -> bounds
}

// SliceKey holds bounds for one frame.
type SliceKey struct {
	X, Y, W, H int
}

type rawSheet struct {
	Frames interface{} `json:"frames"`
	Meta   struct {
		FrameTags []struct {
			Name      string `json:"name"`
			From      int    `json:"from"`
			To        int    `json:"to"`
			Direction string `json:"direction"`
		} `json:"frameTags"`
		Slices []struct {
			Name string `json:"name"`
			Keys []struct {
				Frame  int `json:"frame"`
				Bounds struct {
					X int `json:"x"`
					Y int `json:"y"`
					W int `json:"w"`
					H int `json:"h"`
				} `json:"bounds"`
			} `json:"keys"`
		} `json:"slices"`
	} `json:"meta"`
}

// Load parses an Aseprite JSON file and returns a Sheet.
func Load(path string) (*Sheet, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("aseprite load: %w", err)
	}
	var raw rawSheet
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("aseprite parse: %w", err)
	}
	s := &Sheet{
		Tags:   make(map[string]FrameTag),
		Slices: make(map[string]Slice),
	}
	// Parse frames - can be hash (map) or array
	switch v := raw.Frames.(type) {
	case map[string]interface{}:
		// Hash format: keys are frame names, need to sort by name for consistent order
		names := make([]string, 0, len(v))
		for k := range v {
			names = append(names, k)
		}
		sort.Strings(names)
		for _, name := range names {
			fr, err := parseFrame(v[name])
			if err != nil {
				continue
			}
			s.Frames = append(s.Frames, fr)
		}
	case []interface{}:
		for _, item := range v {
			fr, err := parseFrame(item)
			if err != nil {
				continue
			}
			s.Frames = append(s.Frames, fr)
		}
	default:
		return nil, fmt.Errorf("aseprite: unknown frames format")
	}
	for _, ft := range raw.Meta.FrameTags {
		s.Tags[ft.Name] = FrameTag{
			Name:      ft.Name,
			From:      ft.From,
			To:        ft.To,
			Direction: ft.Direction,
		}
	}
	for _, sl := range raw.Meta.Slices {
		keys := make(map[int]SliceKey)
		for _, k := range sl.Keys {
			keys[k.Frame] = SliceKey{
				X: k.Bounds.X, Y: k.Bounds.Y,
				W: k.Bounds.W, H: k.Bounds.H,
			}
		}
		s.Slices[sl.Name] = Slice{Keys: keys}
	}
	return s, nil
}

func parseFrame(v interface{}) (Frame, error) {
	m, ok := v.(map[string]interface{})
	if !ok {
		return Frame{}, fmt.Errorf("frame not object")
	}
	fr := Frame{}
	if f, ok := m["frame"].(map[string]interface{}); ok {
		fr.X = intNum(f["x"])
		fr.Y = intNum(f["y"])
		fr.W = intNum(f["w"])
		fr.H = intNum(f["h"])
	}
	fr.DurationMs = intNum(m["duration"])
	if fr.DurationMs <= 0 {
		fr.DurationMs = 100
	}
	return fr, nil
}

func intNum(v interface{}) int {
	switch x := v.(type) {
	case float64:
		return int(x)
	case int:
		return x
	default:
		return 0
	}
}

// GetTagFrameRange returns (from, to) for a tag, clamped to frame count.
func (s *Sheet) GetTagFrameRange(tagName string, frameCount int) (from, to int, ok bool) {
	t, ok := s.Tags[tagName]
	if !ok {
		return 0, 0, false
	}
	from = t.From
	to = t.To
	if from < 0 {
		from = 0
	}
	if to >= frameCount {
		to = frameCount - 1
	}
	if from > to {
		from, to = to, from
	}
	return from, to, true
}

// GetSliceBounds returns the slice bounds for the given frame, or (0,0,0,0) if not found.
func (s *Sheet) GetSliceBounds(sliceName string, frameIndex int) (x, y, w, h int) {
	sl, ok := s.Slices[sliceName]
	if !ok {
		return 0, 0, 0, 0
	}
	if k, ok := sl.Keys[frameIndex]; ok {
		return k.X, k.Y, k.W, k.H
	}
	// Fallback: use first available key
	for _, k := range sl.Keys {
		return k.X, k.Y, k.W, k.H
	}
	return 0, 0, 0, 0
}
