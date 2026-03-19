// Package audiosys provides a high-level AUDIO.* API over raylib (additive).
package audiosys

import (
	"cyberbasic/compiler/vm"
	"fmt"
	"strings"
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
	case "load":
		return a.v.CallForeign("AudioLoadSound", ia)
	case "playsoundid", "playid":
		return a.v.CallForeign("AudioPlaySoundId", ia)
	default:
		return nil, fmt.Errorf("unknown audio method %q (use load, playsoundid)", name)
	}
}

// SoundDot stub for future DotObject sound handles.
type SoundDot struct {
	path string
	v    *vm.VM
}

func (s *SoundDot) GetProp(path []string) (vm.Value, error) {
	if len(path) == 1 && strings.ToLower(path[0]) == "path" {
		return s.path, nil
	}
	return nil, nil
}
func (s *SoundDot) SetProp(path []string, val vm.Value) error { return nil }
func (s *SoundDot) CallMethod(name string, args []vm.Value) (vm.Value, error) {
	if name == "play" {
		_, err := s.v.CallForeign("PlaySound", []interface{}{s.path})
		return nil, err
	}
	return nil, nil
}
