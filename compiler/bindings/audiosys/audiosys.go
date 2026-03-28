// Package audiosys provides a high-level AUDIO.* API over raylib (additive).
package audiosys

import (
	"cyberbasic/compiler/vm"
	"fmt"
	"strings"
	"sync"
)

// RegisterAudiosys registers AudioLoadSound (returns sound id), AudioPlaySoundId.
func RegisterAudiosys(v *vm.VM) {
	v.RegisterForeign("AudioLoadSound", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("AudioLoadSound requires (path$)")
		}
		path := fmt.Sprint(args[0])
		return v.CallForeign("LoadSound", []interface{}{path})
	})
	v.RegisterForeign("AudioPlaySoundId", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("AudioPlaySoundId requires (soundId$)")
		}
		return v.CallForeign("PlaySound", []interface{}{fmt.Sprint(args[0])})
	})

	v.SetGlobal("audio", &audioModuleDot{v: v})
}

type audioModuleDot struct {
	v *vm.VM
}

func (a *audioModuleDot) GetProp([]string) (vm.Value, error) { return nil, nil }
func (a *audioModuleDot) SetProp([]string, vm.Value) error {
	return fmt.Errorf("audio: namespace is not assignable")
}

func (a *audioModuleDot) CallMethod(name string, args []vm.Value) (vm.Value, error) {
	ia := make([]interface{}, len(args))
	for i := range args {
		ia[i] = args[i]
	}
	switch strings.ToLower(name) {
	case "load", "sound":
		r, err := a.v.CallForeign("AudioLoadSound", ia)
		if err != nil {
			return nil, err
		}
		id := fmt.Sprint(r)
		return &SoundDot{v: a.v, id: id, vol: 1.0}, nil
	case "playsoundid", "playid":
		return a.v.CallForeign("AudioPlaySoundId", ia)
	default:
		return nil, fmt.Errorf("unknown audio method %q (use load, sound, playsoundid)", name)
	}
}

// SoundDot is a vm.DotObject for a loaded sound id (string).
type SoundDot struct {
	v   *vm.VM
	id  string
	vol float64
	mu  sync.RWMutex
}

func (s *SoundDot) GetProp(path []string) (vm.Value, error) {
	if len(path) != 1 {
		return nil, fmt.Errorf("sound: single property")
	}
	switch strings.ToLower(path[0]) {
	case "id":
		return s.id, nil
	case "volume":
		s.mu.RLock()
		v := s.vol
		s.mu.RUnlock()
		return v, nil
	default:
		return nil, nil
	}
}

func (s *SoundDot) SetProp(path []string, val vm.Value) error {
	if len(path) != 1 {
		return fmt.Errorf("sound: single property")
	}
	if strings.ToLower(path[0]) != "volume" {
		return fmt.Errorf("sound: only volume is writable")
	}
	vol := toF64(val)
	s.mu.Lock()
	s.vol = vol
	s.mu.Unlock()
	_, err := s.v.CallForeign("SetSoundVolume", []interface{}{s.id, vol})
	return err
}

func (s *SoundDot) CallMethod(name string, args []vm.Value) (vm.Value, error) {
	switch strings.ToLower(name) {
	case "play":
		_, err := s.v.CallForeign("PlaySound", []interface{}{s.id})
		return nil, err
	case "stop":
		_, err := s.v.CallForeign("StopSound", []interface{}{s.id})
		return nil, err
	case "unload":
		_, err := s.v.CallForeign("UnloadSound", []interface{}{s.id})
		return nil, err
	default:
		return nil, fmt.Errorf("sound: use play, stop, unload")
	}
}

func toF64(v vm.Value) float64 {
	switch x := v.(type) {
	case float64:
		return x
	case int:
		return float64(x)
	case int32:
		return float64(x)
	default:
		return 0
	}
}
