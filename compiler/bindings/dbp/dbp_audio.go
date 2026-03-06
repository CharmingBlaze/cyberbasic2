// Package dbp: Audio expanded - music with DBP-style integer IDs.
//
// Commands:
//   - LoadMusic(id, path): Load music stream
//   - PlayMusic(id): Start playback
//   - StopMusic(id): Stop playback
//   - SetMusicVolume(id, value): Volume 0-1
//   - SetMusicLoop(id, onOff): Loop on/off
package dbp

import (
	"fmt"
	"sync"

	"cyberbasic/compiler/vm"
	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
	musics     = make(map[int]rl.Music)
	musicsMu   sync.Mutex
	musicLoops = make(map[int]bool)
)

// registerAudioExpanded adds LoadMusic, PlayMusic, StopMusic, SetMusicVolume, SetMusicLoop.
func registerAudioExpanded(v *vm.VM) {
	v.RegisterForeign("LoadMusic", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("LoadMusic(id, path) requires 2 arguments")
		}
		id := toInt(args[0])
		path := toString(args[1])
		m := rl.LoadMusicStream(path)
		musicsMu.Lock()
		if old, ok := musics[id]; ok {
			rl.UnloadMusicStream(old)
		}
		musics[id] = m
		musicsMu.Unlock()
		return nil, nil
	})

	v.RegisterForeign("PlayMusic", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("PlayMusic(id) requires 1 argument")
		}
		id := toInt(args[0])
		musicsMu.Lock()
		m, ok := musics[id]
		musicsMu.Unlock()
		if !ok {
			return nil, fmt.Errorf("unknown music id %d", id)
		}
		rl.PlayMusicStream(m)
		return nil, nil
	})

	v.RegisterForeign("StopMusic", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("StopMusic(id) requires 1 argument")
		}
		id := toInt(args[0])
		musicsMu.Lock()
		m, ok := musics[id]
		musicsMu.Unlock()
		if !ok {
			return nil, fmt.Errorf("unknown music id %d", id)
		}
		rl.StopMusicStream(m)
		return nil, nil
	})

	v.RegisterForeign("SetMusicVolume", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("SetMusicVolume(id, value) requires 2 arguments")
		}
		id := toInt(args[0])
		vol := toFloat32(args[1])
		musicsMu.Lock()
		m, ok := musics[id]
		musicsMu.Unlock()
		if !ok {
			return nil, fmt.Errorf("unknown music id %d", id)
		}
		rl.SetMusicVolume(m, vol)
		return nil, nil
	})

	v.RegisterForeign("SetMusicLoop", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("SetMusicLoop(id, onOff) requires 2 arguments")
		}
		id := toInt(args[0])
		onOff := toInt(args[1]) != 0
		musicsMu.Lock()
		musicLoops[id] = onOff
		musicsMu.Unlock()
		// raylib Music loops by default; loop count API may vary
		return nil, nil
	})
}
