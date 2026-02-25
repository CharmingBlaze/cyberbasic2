// Package raylib: audio device, sounds, and music (raudio).
package raylib

import (
	"fmt"
	"os"
	"strings"
	"unsafe"

	"cyberbasic/compiler/vm"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func registerAudio(v *vm.VM) {
	v.RegisterForeign("InitAudioDevice", func(args []interface{}) (interface{}, error) {
		rl.InitAudioDevice()
		return nil, nil
	})
	v.RegisterForeign("CloseAudioDevice", func(args []interface{}) (interface{}, error) {
		rl.CloseAudioDevice()
		return nil, nil
	})
	v.RegisterForeign("IsAudioDeviceReady", func(args []interface{}) (interface{}, error) {
		return rl.IsAudioDeviceReady(), nil
	})
	v.RegisterForeign("LoadSound", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("LoadSound requires (fileName)")
		}
		path := toString(args[0])
		sound := rl.LoadSound(path)
		soundMu.Lock()
		soundCounter++
		id := fmt.Sprintf("sound_%d", soundCounter)
		sounds[id] = sound
		soundMu.Unlock()
		return id, nil
	})
	v.RegisterForeign("PlaySound", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("PlaySound requires (id)")
		}
		id := toString(args[0])
		soundMu.Lock()
		sound, ok := sounds[id]
		soundMu.Unlock()
		if !ok {
			return nil, fmt.Errorf("unknown sound id: %s", id)
		}
		rl.PlaySound(sound)
		return nil, nil
	})
	v.RegisterForeign("StopSound", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("StopSound requires (id)")
		}
		id := toString(args[0])
		soundMu.Lock()
		sound, ok := sounds[id]
		soundMu.Unlock()
		if !ok {
			return nil, fmt.Errorf("unknown sound id: %s", id)
		}
		rl.StopSound(sound)
		return nil, nil
	})
	v.RegisterForeign("SetSoundVolume", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("SetSoundVolume requires (id, volume)")
		}
		id := toString(args[0])
		soundMu.Lock()
		sound, ok := sounds[id]
		soundMu.Unlock()
		if !ok {
			return nil, fmt.Errorf("unknown sound id: %s", id)
		}
		rl.SetSoundVolume(sound, toFloat32(args[1]))
		return nil, nil
	})
	v.RegisterForeign("UnloadSound", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("UnloadSound requires (id)")
		}
		id := toString(args[0])
		soundMu.Lock()
		sound, ok := sounds[id]
		delete(sounds, id)
		soundMu.Unlock()
		if ok {
			rl.UnloadSound(sound)
		}
		return nil, nil
	})
	v.RegisterForeign("LoadMusicStream", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("LoadMusicStream requires (fileName)")
		}
		path := toString(args[0])
		stream := rl.LoadMusicStream(path)
		musicMu.Lock()
		musicCounter++
		id := fmt.Sprintf("music_%d", musicCounter)
		music[id] = stream
		musicMu.Unlock()
		return id, nil
	})
	v.RegisterForeign("PlayMusicStream", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("PlayMusicStream requires (id)")
		}
		id := toString(args[0])
		musicMu.Lock()
		m, ok := music[id]
		musicMu.Unlock()
		if !ok {
			return nil, fmt.Errorf("unknown music id: %s", id)
		}
		rl.PlayMusicStream(m)
		return nil, nil
	})
	v.RegisterForeign("UpdateMusicStream", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("UpdateMusicStream requires (id)")
		}
		id := toString(args[0])
		musicMu.Lock()
		m, ok := music[id]
		musicMu.Unlock()
		if !ok {
			return nil, fmt.Errorf("unknown music id: %s", id)
		}
		rl.UpdateMusicStream(m)
		return nil, nil
	})
	v.RegisterForeign("StopMusicStream", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("StopMusicStream requires (id)")
		}
		id := toString(args[0])
		musicMu.Lock()
		m, ok := music[id]
		musicMu.Unlock()
		if !ok {
			return nil, fmt.Errorf("unknown music id: %s", id)
		}
		rl.StopMusicStream(m)
		return nil, nil
	})
	v.RegisterForeign("SetMusicVolume", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("SetMusicVolume requires (id, volume)")
		}
		id := toString(args[0])
		musicMu.Lock()
		m, ok := music[id]
		musicMu.Unlock()
		if !ok {
			return nil, fmt.Errorf("unknown music id: %s", id)
		}
		rl.SetMusicVolume(m, toFloat32(args[1]))
		return nil, nil
	})
	v.RegisterForeign("UnloadMusicStream", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("UnloadMusicStream requires (id)")
		}
		id := toString(args[0])
		musicMu.Lock()
		m, ok := music[id]
		delete(music, id)
		musicMu.Unlock()
		if ok {
			rl.UnloadMusicStream(m)
		}
		return nil, nil
	})
	v.RegisterForeign("SetMasterVolume", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("SetMasterVolume requires (volume)")
		}
		rl.SetMasterVolume(toFloat32(args[0]))
		return nil, nil
	})
	v.RegisterForeign("GetMasterVolume", func(args []interface{}) (interface{}, error) {
		return float64(rl.GetMasterVolume()), nil
	})
	v.RegisterForeign("PauseSound", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("PauseSound requires (id)")
		}
		id := toString(args[0])
		soundMu.Lock()
		sound, ok := sounds[id]
		soundMu.Unlock()
		if !ok {
			return nil, fmt.Errorf("unknown sound id: %s", id)
		}
		rl.PauseSound(sound)
		return nil, nil
	})
	v.RegisterForeign("ResumeSound", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("ResumeSound requires (id)")
		}
		id := toString(args[0])
		soundMu.Lock()
		sound, ok := sounds[id]
		soundMu.Unlock()
		if !ok {
			return nil, fmt.Errorf("unknown sound id: %s", id)
		}
		rl.ResumeSound(sound)
		return nil, nil
	})
	v.RegisterForeign("IsSoundPlaying", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("IsSoundPlaying requires (id)")
		}
		id := toString(args[0])
		soundMu.Lock()
		sound, ok := sounds[id]
		soundMu.Unlock()
		if !ok {
			return nil, fmt.Errorf("unknown sound id: %s", id)
		}
		return rl.IsSoundPlaying(sound), nil
	})
	v.RegisterForeign("SetSoundPitch", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("SetSoundPitch requires (id, pitch)")
		}
		id := toString(args[0])
		soundMu.Lock()
		sound, ok := sounds[id]
		soundMu.Unlock()
		if !ok {
			return nil, fmt.Errorf("unknown sound id: %s", id)
		}
		rl.SetSoundPitch(sound, toFloat32(args[1]))
		return nil, nil
	})
	v.RegisterForeign("SetSoundPan", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("SetSoundPan requires (id, pan)")
		}
		id := toString(args[0])
		soundMu.Lock()
		sound, ok := sounds[id]
		soundMu.Unlock()
		if !ok {
			return nil, fmt.Errorf("unknown sound id: %s", id)
		}
		rl.SetSoundPan(sound, toFloat32(args[1]))
		return nil, nil
	})
	v.RegisterForeign("PauseMusicStream", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("PauseMusicStream requires (id)")
		}
		id := toString(args[0])
		musicMu.Lock()
		m, ok := music[id]
		musicMu.Unlock()
		if !ok {
			return nil, fmt.Errorf("unknown music id: %s", id)
		}
		rl.PauseMusicStream(m)
		return nil, nil
	})
	v.RegisterForeign("ResumeMusicStream", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("ResumeMusicStream requires (id)")
		}
		id := toString(args[0])
		musicMu.Lock()
		m, ok := music[id]
		musicMu.Unlock()
		if !ok {
			return nil, fmt.Errorf("unknown music id: %s", id)
		}
		rl.ResumeMusicStream(m)
		return nil, nil
	})
	v.RegisterForeign("IsMusicStreamPlaying", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("IsMusicStreamPlaying requires (id)")
		}
		id := toString(args[0])
		musicMu.Lock()
		m, ok := music[id]
		musicMu.Unlock()
		if !ok {
			return nil, fmt.Errorf("unknown music id: %s", id)
		}
		return rl.IsMusicStreamPlaying(m), nil
	})
	// Aliases: LoadMusic, PlayMusic, PauseMusic, ResumeMusic, IsMusicPlaying
	v.RegisterForeign("LoadMusic", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("LoadMusic requires (path)")
		}
		path := toString(args[0])
		stream := rl.LoadMusicStream(path)
		musicMu.Lock()
		musicCounter++
		id := fmt.Sprintf("music_%d", musicCounter)
		music[id] = stream
		musicMu.Unlock()
		return id, nil
	})
	v.RegisterForeign("PlayMusic", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("PlayMusic requires (musicId)")
		}
		id := toString(args[0])
		musicMu.Lock()
		m, ok := music[id]
		musicMu.Unlock()
		if !ok {
			return nil, fmt.Errorf("unknown music id: %s", id)
		}
		rl.PlayMusicStream(m)
		return nil, nil
	})
	v.RegisterForeign("PauseMusic", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("PauseMusic requires (musicId)")
		}
		id := toString(args[0])
		musicMu.Lock()
		m, ok := music[id]
		musicMu.Unlock()
		if !ok {
			return nil, fmt.Errorf("unknown music id: %s", id)
		}
		rl.PauseMusicStream(m)
		return nil, nil
	})
	v.RegisterForeign("ResumeMusic", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("ResumeMusic requires (musicId)")
		}
		id := toString(args[0])
		musicMu.Lock()
		m, ok := music[id]
		musicMu.Unlock()
		if !ok {
			return nil, fmt.Errorf("unknown music id: %s", id)
		}
		rl.ResumeMusicStream(m)
		return nil, nil
	})
	v.RegisterForeign("IsMusicPlaying", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("IsMusicPlaying requires (musicId)")
		}
		id := toString(args[0])
		musicMu.Lock()
		m, ok := music[id]
		musicMu.Unlock()
		if !ok {
			return nil, fmt.Errorf("unknown music id: %s", id)
		}
		return rl.IsMusicStreamPlaying(m), nil
	})
	v.RegisterForeign("SeekMusicStream", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("SeekMusicStream requires (id, position)")
		}
		id := toString(args[0])
		musicMu.Lock()
		m, ok := music[id]
		musicMu.Unlock()
		if !ok {
			return nil, fmt.Errorf("unknown music id: %s", id)
		}
		rl.SeekMusicStream(m, toFloat32(args[1]))
		return nil, nil
	})
	v.RegisterForeign("SetMusicPitch", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("SetMusicPitch requires (id, pitch)")
		}
		id := toString(args[0])
		musicMu.Lock()
		m, ok := music[id]
		musicMu.Unlock()
		if !ok {
			return nil, fmt.Errorf("unknown music id: %s", id)
		}
		rl.SetMusicPitch(m, toFloat32(args[1]))
		return nil, nil
	})
	v.RegisterForeign("SetMusicPan", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("SetMusicPan requires (id, pan)")
		}
		id := toString(args[0])
		musicMu.Lock()
		m, ok := music[id]
		musicMu.Unlock()
		if !ok {
			return nil, fmt.Errorf("unknown music id: %s", id)
		}
		rl.SetMusicPan(m, toFloat32(args[1]))
		return nil, nil
	})
	v.RegisterForeign("GetMusicTimeLength", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("GetMusicTimeLength requires (id)")
		}
		id := toString(args[0])
		musicMu.Lock()
		m, ok := music[id]
		musicMu.Unlock()
		if !ok {
			return nil, fmt.Errorf("unknown music id: %s", id)
		}
		return float64(rl.GetMusicTimeLength(m)), nil
	})
	v.RegisterForeign("GetMusicTimePlayed", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("GetMusicTimePlayed requires (id)")
		}
		id := toString(args[0])
		musicMu.Lock()
		m, ok := music[id]
		musicMu.Unlock()
		if !ok {
			return nil, fmt.Errorf("unknown music id: %s", id)
		}
		return float64(rl.GetMusicTimePlayed(m)), nil
	})
	v.RegisterForeign("IsMusicValid", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return false, nil
		}
		musicMu.Lock()
		m, ok := music[toString(args[0])]
		musicMu.Unlock()
		return ok && rl.IsMusicValid(m), nil
	})
	v.RegisterForeign("LoadMusicStreamFromMemory", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("LoadMusicStreamFromMemory requires (fileType, data, dataSize)")
		}
		var data []byte
		switch d := args[1].(type) {
		case string:
			data = []byte(d)
		case []byte:
			data = d
		default:
			return nil, fmt.Errorf("data must be string or []byte")
		}
		dataSize := toInt32(args[2])
		if int(dataSize) < len(data) {
			data = data[:dataSize]
		}
		m := rl.LoadMusicStreamFromMemory(toString(args[0]), data, dataSize)
		musicMu.Lock()
		musicCounter++
		id := fmt.Sprintf("music_%d", musicCounter)
		music[id] = m
		musicMu.Unlock()
		return id, nil
	})

	// Wave
	v.RegisterForeign("LoadWave", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("LoadWave requires (fileName)")
		}
		w := rl.LoadWave(toString(args[0]))
		waveMu.Lock()
		waveCounter++
		id := fmt.Sprintf("wave_%d", waveCounter)
		waves[id] = w
		waveMu.Unlock()
		return id, nil
	})
	v.RegisterForeign("LoadWaveFromMemory", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("LoadWaveFromMemory requires (fileType, data, dataSize)")
		}
		var data []byte
		switch d := args[1].(type) {
		case string:
			data = []byte(d)
		case []byte:
			data = d
		default:
			return nil, fmt.Errorf("data must be string or []byte")
		}
		dataSize := toInt32(args[2])
		if int(dataSize) < len(data) {
			data = data[:dataSize]
		}
		w := rl.LoadWaveFromMemory(toString(args[0]), data, dataSize)
		waveMu.Lock()
		waveCounter++
		id := fmt.Sprintf("wave_%d", waveCounter)
		waves[id] = w
		waveMu.Unlock()
		return id, nil
	})
	v.RegisterForeign("IsWaveValid", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return false, nil
		}
		waveMu.Lock()
		w, ok := waves[toString(args[0])]
		waveMu.Unlock()
		return ok && rl.IsWaveValid(w), nil
	})
	v.RegisterForeign("UnloadWave", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("UnloadWave requires (waveId)")
		}
		id := toString(args[0])
		waveMu.Lock()
		w, ok := waves[id]
		delete(waves, id)
		waveMu.Unlock()
		if ok {
			rl.UnloadWave(w)
		}
		return nil, nil
	})
	v.RegisterForeign("ExportWave", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("ExportWave requires (waveId, fileName)")
		}
		waveMu.Lock()
		w, ok := waves[toString(args[0])]
		waveMu.Unlock()
		if !ok {
			return nil, fmt.Errorf("unknown wave id: %s", toString(args[0]))
		}
		rl.ExportWave(w, toString(args[1]))
		return nil, nil
	})
	v.RegisterForeign("WaveCopy", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("WaveCopy requires (waveId)")
		}
		waveMu.Lock()
		w, ok := waves[toString(args[0])]
		waveMu.Unlock()
		if !ok {
			return nil, fmt.Errorf("unknown wave id: %s", toString(args[0]))
		}
		copyW := rl.WaveCopy(w)
		waveMu.Lock()
		waveCounter++
		id := fmt.Sprintf("wave_%d", waveCounter)
		waves[id] = copyW
		waveMu.Unlock()
		return id, nil
	})
	v.RegisterForeign("WaveCrop", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("WaveCrop requires (waveId, initFrame, finalFrame)")
		}
		waveMu.Lock()
		w, ok := waves[toString(args[0])]
		if !ok {
			waveMu.Unlock()
			return nil, fmt.Errorf("unknown wave id: %s", toString(args[0]))
		}
		rl.WaveCrop(&w, toInt32(args[1]), toInt32(args[2]))
		waves[toString(args[0])] = w
		waveMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("WaveFormat", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("WaveFormat requires (waveId, sampleRate, sampleSize, channels)")
		}
		waveMu.Lock()
		w, ok := waves[toString(args[0])]
		if !ok {
			waveMu.Unlock()
			return nil, fmt.Errorf("unknown wave id: %s", toString(args[0]))
		}
		rl.WaveFormat(&w, toInt32(args[1]), toInt32(args[2]), toInt32(args[3]))
		waves[toString(args[0])] = w
		waveMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("LoadWaveSamples", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("LoadWaveSamples requires (waveId)")
		}
		waveMu.Lock()
		w, ok := waves[toString(args[0])]
		waveMu.Unlock()
		if !ok {
			return nil, fmt.Errorf("unknown wave id: %s", toString(args[0]))
		}
		samples := rl.LoadWaveSamples(w)
		lastWaveSamplesMu.Lock()
		lastWaveSamples = samples
		lastWaveSamplesMu.Unlock()
		return len(samples), nil
	})
	v.RegisterForeign("UnloadWaveSamples", func(args []interface{}) (interface{}, error) {
		lastWaveSamplesMu.Lock()
		if len(lastWaveSamples) > 0 {
			rl.UnloadWaveSamples(lastWaveSamples)
			lastWaveSamples = nil
		}
		lastWaveSamplesMu.Unlock()
		return nil, nil
	})
	// ExportWaveAsCode: export wave PCM data as C header (raylib-go has no native API; we generate .h from Wave).
	v.RegisterForeign("ExportWaveAsCode", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("ExportWaveAsCode requires (waveId, fileName)")
		}
		waveMu.Lock()
		w, ok := waves[toString(args[0])]
		waveMu.Unlock()
		if !ok {
			return false, nil
		}
		bytesPerSample := int(w.SampleSize / 8)
		if bytesPerSample == 0 {
			bytesPerSample = 1
		}
		size := int(w.FrameCount) * bytesPerSample * int(w.Channels)
		if size <= 0 || w.Data == nil {
			return false, nil
		}
		data := unsafe.Slice((*byte)(w.Data), size)
		var b strings.Builder
		b.WriteString("// Exported by CyberBasic ExportWaveAsCode\n")
		b.WriteString("#ifndef WAVE_EXPORT_H\n#define WAVE_EXPORT_H\n\n")
		fmt.Fprintf(&b, "static const unsigned int WAVE_FRAME_COUNT = %d;\n", w.FrameCount)
		fmt.Fprintf(&b, "static const unsigned int WAVE_SAMPLE_RATE = %d;\n", w.SampleRate)
		fmt.Fprintf(&b, "static const unsigned int WAVE_SAMPLE_SIZE = %d;\n", w.SampleSize)
		fmt.Fprintf(&b, "static const unsigned int WAVE_CHANNELS = %d;\n", w.Channels)
		b.WriteString("static const unsigned char WAVE_DATA[] = {\n")
		for i := 0; i < len(data); i++ {
			if i > 0 {
				b.WriteByte(',')
			}
			if i%16 == 0 {
				b.WriteString("\n    ")
			}
			fmt.Fprintf(&b, "%d", data[i])
		}
		b.WriteString("\n};\n\n#endif\n")
		if err := os.WriteFile(toString(args[1]), []byte(b.String()), 0644); err != nil {
			return false, err
		}
		return true, nil
	})

	// Sound from wave, alias, update, valid
	v.RegisterForeign("LoadSoundFromWave", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("LoadSoundFromWave requires (waveId)")
		}
		waveMu.Lock()
		w, ok := waves[toString(args[0])]
		waveMu.Unlock()
		if !ok {
			return nil, fmt.Errorf("unknown wave id: %s", toString(args[0]))
		}
		sound := rl.LoadSoundFromWave(w)
		soundMu.Lock()
		soundCounter++
		id := fmt.Sprintf("sound_%d", soundCounter)
		sounds[id] = sound
		soundMu.Unlock()
		return id, nil
	})
	v.RegisterForeign("LoadSoundAlias", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("LoadSoundAlias requires (sourceSoundId)")
		}
		soundMu.Lock()
		src, ok := sounds[toString(args[0])]
		soundMu.Unlock()
		if !ok {
			return nil, fmt.Errorf("unknown sound id: %s", toString(args[0]))
		}
		alias := rl.LoadSoundAlias(src)
		soundMu.Lock()
		soundCounter++
		id := fmt.Sprintf("sound_%d", soundCounter)
		sounds[id] = alias
		soundMu.Unlock()
		return id, nil
	})
	v.RegisterForeign("IsSoundValid", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return false, nil
		}
		soundMu.Lock()
		s, ok := sounds[toString(args[0])]
		soundMu.Unlock()
		return ok && rl.IsSoundValid(s), nil
	})
	v.RegisterForeign("UpdateSound", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("UpdateSound requires (soundId, data, sampleCount)")
		}
		soundMu.Lock()
		s, ok := sounds[toString(args[0])]
		soundMu.Unlock()
		if !ok {
			return nil, fmt.Errorf("unknown sound id: %s", toString(args[0]))
		}
		var data []byte
		switch d := args[1].(type) {
		case string:
			data = []byte(d)
		case []byte:
			data = d
		default:
			return nil, fmt.Errorf("data must be string or []byte")
		}
		sampleCount := toInt32(args[2])
		rl.UpdateSound(s, data, sampleCount)
		return nil, nil
	})
	v.RegisterForeign("UnloadSoundAlias", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("UnloadSoundAlias requires (aliasSoundId)")
		}
		id := toString(args[0])
		soundMu.Lock()
		s, ok := sounds[id]
		delete(sounds, id)
		soundMu.Unlock()
		if ok {
			rl.UnloadSoundAlias(s)
		}
		return nil, nil
	})

	// AudioStream
	v.RegisterForeign("LoadAudioStream", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("LoadAudioStream requires (sampleRate, sampleSize, channels)")
		}
		stream := rl.LoadAudioStream(uint32(toInt32(args[0])), uint32(toInt32(args[1])), uint32(toInt32(args[2])))
		audioStreamMu.Lock()
		audioStreamCounter++
		id := fmt.Sprintf("stream_%d", audioStreamCounter)
		audioStreams[id] = stream
		audioStreamMu.Unlock()
		return id, nil
	})
	v.RegisterForeign("IsAudioStreamValid", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return false, nil
		}
		audioStreamMu.Lock()
		st, ok := audioStreams[toString(args[0])]
		audioStreamMu.Unlock()
		return ok && rl.IsAudioStreamValid(st), nil
	})
	v.RegisterForeign("UnloadAudioStream", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("UnloadAudioStream requires (streamId)")
		}
		id := toString(args[0])
		audioStreamMu.Lock()
		st, ok := audioStreams[id]
		delete(audioStreams, id)
		audioStreamMu.Unlock()
		if ok {
			rl.UnloadAudioStream(st)
		}
		return nil, nil
	})
	v.RegisterForeign("UpdateAudioStream", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("UpdateAudioStream requires (streamId, ...floatSamples)")
		}
		audioStreamMu.Lock()
		st, ok := audioStreams[toString(args[0])]
		audioStreamMu.Unlock()
		if !ok {
			return nil, fmt.Errorf("unknown stream id: %s", toString(args[0]))
		}
		n := len(args) - 1
		data := make([]float32, n)
		for i := 0; i < n; i++ {
			data[i] = toFloat32(args[1+i])
		}
		rl.UpdateAudioStream(st, data)
		return nil, nil
	})
	v.RegisterForeign("IsAudioStreamProcessed", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return false, nil
		}
		audioStreamMu.Lock()
		st, ok := audioStreams[toString(args[0])]
		audioStreamMu.Unlock()
		return ok && rl.IsAudioStreamProcessed(st), nil
	})
	v.RegisterForeign("PlayAudioStream", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("PlayAudioStream requires (streamId)")
		}
		audioStreamMu.Lock()
		st, ok := audioStreams[toString(args[0])]
		audioStreamMu.Unlock()
		if !ok {
			return nil, fmt.Errorf("unknown stream id: %s", toString(args[0]))
		}
		rl.PlayAudioStream(st)
		return nil, nil
	})
	v.RegisterForeign("PauseAudioStream", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("PauseAudioStream requires (streamId)")
		}
		audioStreamMu.Lock()
		st, ok := audioStreams[toString(args[0])]
		audioStreamMu.Unlock()
		if !ok {
			return nil, fmt.Errorf("unknown stream id: %s", toString(args[0]))
		}
		rl.PauseAudioStream(st)
		return nil, nil
	})
	v.RegisterForeign("ResumeAudioStream", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("ResumeAudioStream requires (streamId)")
		}
		audioStreamMu.Lock()
		st, ok := audioStreams[toString(args[0])]
		audioStreamMu.Unlock()
		if !ok {
			return nil, fmt.Errorf("unknown stream id: %s", toString(args[0]))
		}
		rl.ResumeAudioStream(st)
		return nil, nil
	})
	v.RegisterForeign("IsAudioStreamPlaying", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return false, nil
		}
		audioStreamMu.Lock()
		st, ok := audioStreams[toString(args[0])]
		audioStreamMu.Unlock()
		return ok && rl.IsAudioStreamPlaying(st), nil
	})
	v.RegisterForeign("StopAudioStream", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("StopAudioStream requires (streamId)")
		}
		audioStreamMu.Lock()
		st, ok := audioStreams[toString(args[0])]
		audioStreamMu.Unlock()
		if !ok {
			return nil, fmt.Errorf("unknown stream id: %s", toString(args[0]))
		}
		rl.StopAudioStream(st)
		return nil, nil
	})
	v.RegisterForeign("SetAudioStreamVolume", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("SetAudioStreamVolume requires (streamId, volume)")
		}
		audioStreamMu.Lock()
		st, ok := audioStreams[toString(args[0])]
		audioStreamMu.Unlock()
		if !ok {
			return nil, fmt.Errorf("unknown stream id: %s", toString(args[0]))
		}
		rl.SetAudioStreamVolume(st, toFloat32(args[1]))
		return nil, nil
	})
	v.RegisterForeign("SetAudioStreamPitch", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("SetAudioStreamPitch requires (streamId, pitch)")
		}
		audioStreamMu.Lock()
		st, ok := audioStreams[toString(args[0])]
		audioStreamMu.Unlock()
		if !ok {
			return nil, fmt.Errorf("unknown stream id: %s", toString(args[0]))
		}
		rl.SetAudioStreamPitch(st, toFloat32(args[1]))
		return nil, nil
	})
	v.RegisterForeign("SetAudioStreamPan", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("SetAudioStreamPan requires (streamId, pan)")
		}
		audioStreamMu.Lock()
		st, ok := audioStreams[toString(args[0])]
		audioStreamMu.Unlock()
		if !ok {
			return nil, fmt.Errorf("unknown stream id: %s", toString(args[0]))
		}
		rl.SetAudioStreamPan(st, toFloat32(args[1]))
		return nil, nil
	})
	v.RegisterForeign("SetAudioStreamBufferSizeDefault", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("SetAudioStreamBufferSizeDefault requires (size)")
		}
		rl.SetAudioStreamBufferSizeDefault(toInt32(args[0]))
		return nil, nil
	})

	// Audio callbacks require a C/Go function pointer; BASIC cannot pass functions. These bindings exist for API completeness but return an error when called.
	v.RegisterForeign("SetAudioStreamCallback", func(args []interface{}) (interface{}, error) {
		return nil, fmt.Errorf("SetAudioStreamCallback requires a callback function; not supported from BASIC (use UpdateAudioStream to push samples instead)")
	})
	v.RegisterForeign("AttachAudioStreamProcessor", func(args []interface{}) (interface{}, error) {
		return nil, fmt.Errorf("AttachAudioStreamProcessor requires a callback function; not supported from BASIC")
	})
	v.RegisterForeign("DetachAudioStreamProcessor", func(args []interface{}) (interface{}, error) {
		return nil, fmt.Errorf("DetachAudioStreamProcessor requires a callback function; not supported from BASIC")
	})
	v.RegisterForeign("AttachAudioMixedProcessor", func(args []interface{}) (interface{}, error) {
		return nil, fmt.Errorf("AttachAudioMixedProcessor requires a callback function; not supported from BASIC")
	})
	v.RegisterForeign("DetachAudioMixedProcessor", func(args []interface{}) (interface{}, error) {
		return nil, fmt.Errorf("DetachAudioMixedProcessor requires a callback function; not supported from BASIC")
	})
}
