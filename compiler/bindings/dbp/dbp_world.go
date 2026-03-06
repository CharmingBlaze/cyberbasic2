// Package dbp: World and scene extras - skybox, ambient light, fog.
//
// These commands control the 3D environment:
//   - SetSkybox(path): Load and display a skybox cubemap
//   - SetAmbientLight(r, g, b): Set ambient light color (stored for shaders)
//   - SetFog(onOff): Enable or disable distance fog
//   - SetFogColor(r, g, b): Fog color
//   - SetFogRange(near, far): Fog density range
package dbp

import (
	"fmt"
	"sync"

	"cyberbasic/compiler/vm"
	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
	skyboxTex     rl.Texture2D
	skyboxLoaded  bool
	skyboxMu      sync.Mutex
	ambientR      float32 = 0.4
	ambientG      float32 = 0.4
	ambientB      float32 = 0.5
	fogEnabled    bool
	fogColorR     float32 = 0.5
	fogColorG     float32 = 0.6
	fogColorB     float32 = 0.7
	fogDensity    float32 = 0.02
	fogNear       float32 = 10
	fogFar        float32 = 50
	worldMu       sync.Mutex

	// Clouds
	cloudsOn      bool
	cloudTexture  rl.Texture2D
	cloudTexLoaded bool
	cloudSpeed    float32 = 0.001
	cloudDensity  float32 = 0.5
	cloudHeight   float32 = 100
	cloudColorR   float32 = 1
	cloudColorG   float32 = 1
	cloudColorB   float32 = 1
	cloudMu       sync.Mutex

	// Sun
	sunDirX       float32 = 0.5
	sunDirY       float32 = -0.7
	sunDirZ       float32 = 0.5
	sunColorR     float32 = 1
	sunColorG     float32 = 0.95
	sunColorB     float32 = 0.9
	sunIntensity  float32 = 1
	sunMu         sync.Mutex

	// World time (0-24 hours)
	worldTime     float32 = 12
	worldTimeScale float32 = 1
	weatherPreset string = "Clear"
	worldTimeMu   sync.Mutex
)

// registerWorld adds SetSkybox, SetAmbientLight, SetFog, SetFogColor, SetFogRange.
func registerWorld(v *vm.VM) {
	// SetSkybox(path): Load skybox. If missing/invalid -> fallback to solid color sky (never crash).
	v.RegisterForeign("SetSkybox", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("SetSkybox(path) requires 1 argument")
		}
		path := toString(args[0])
		skyboxMu.Lock()
		if skyboxLoaded && skyboxTex.ID > 0 {
			rl.UnloadTexture(skyboxTex)
		}
		skyboxLoaded = false
		if path != "" {
			tex := rl.LoadTexture(path)
			if tex.ID > 0 {
				skyboxTex = tex
				skyboxLoaded = true
			}
		}
		skyboxMu.Unlock()
		return nil, nil
	})

	// SetAmbientLight(r, g, b): Set ambient light color (0-1 or 0-255).
	v.RegisterForeign("SetAmbientLight", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("SetAmbientLight(r, g, b) requires 3 arguments")
		}
		r, g, b := toFloat32(args[0]), toFloat32(args[1]), toFloat32(args[2])
		if r > 1 || g > 1 || b > 1 {
			r, g, b = r/255, g/255, b/255
		}
		worldMu.Lock()
		ambientR, ambientG, ambientB = r, g, b
		worldMu.Unlock()
		return nil, nil
	})

	// SetFog(onOff): Enable or disable fog. Uses stored fog color and density.
	v.RegisterForeign("SetFog", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("SetFog(onOff) requires 1 argument")
		}
		worldMu.Lock()
		fogEnabled = toInt(args[0]) != 0
		worldMu.Unlock()
		// Raylib SetFog expects (enable, density, r, g, b)
		_, _ = v.CallForeign("SetFog", []interface{}{fogEnabled, fogDensity, fogColorR * 255, fogColorG * 255, fogColorB * 255})
		return nil, nil
	})

	// SetFogColor(r, g, b): Set fog color (0-255).
	v.RegisterForeign("SetFogColor", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("SetFogColor(r, g, b) requires 3 arguments")
		}
		r, g, b := toFloat32(args[0]), toFloat32(args[1]), toFloat32(args[2])
		worldMu.Lock()
		fogColorR, fogColorG, fogColorB = r/255, g/255, b/255
		worldMu.Unlock()
		return nil, nil
	})

	// SetFogRange(near, far): Fog start and end distance.
	v.RegisterForeign("SetFogRange", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("SetFogRange(near, far) requires 2 arguments")
		}
		worldMu.Lock()
		fogNear = toFloat32(args[0])
		fogFar = toFloat32(args[1])
		worldMu.Unlock()
		return nil, nil
	})

	// Skybox and ambient light queries
	v.RegisterForeign("GetSkybox", func(args []interface{}) (interface{}, error) {
		skyboxMu.Lock()
		loaded := skyboxLoaded && skyboxTex.ID > 0
		skyboxMu.Unlock()
		if loaded {
			return 1, nil
		}
		return 0, nil
	})
	v.RegisterForeign("GetAmbientLightR", func(args []interface{}) (interface{}, error) {
		worldMu.Lock()
		r := ambientR
		worldMu.Unlock()
		return float64(r * 255), nil
	})
	v.RegisterForeign("GetAmbientLightG", func(args []interface{}) (interface{}, error) {
		worldMu.Lock()
		g := ambientG
		worldMu.Unlock()
		return float64(g * 255), nil
	})
	v.RegisterForeign("GetAmbientLightB", func(args []interface{}) (interface{}, error) {
		worldMu.Lock()
		b := ambientB
		worldMu.Unlock()
		return float64(b * 255), nil
	})

	// --- Clouds ---
	v.RegisterForeign("SetFogOff", func(args []interface{}) (interface{}, error) {
		worldMu.Lock()
		fogEnabled = false
		worldMu.Unlock()
		_, _ = v.CallForeign("SetFog", []interface{}{false, fogDensity, fogColorR * 255, fogColorG * 255, fogColorB * 255})
		return nil, nil
	})

	v.RegisterForeign("SetCloudsOn", func(args []interface{}) (interface{}, error) {
		cloudMu.Lock()
		cloudsOn = true
		cloudMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("SetCloudsOff", func(args []interface{}) (interface{}, error) {
		cloudMu.Lock()
		cloudsOn = false
		cloudMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("SetCloudTexture", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("SetCloudTexture(path) requires 1 argument")
		}
		path := toString(args[0])
		cloudMu.Lock()
		if cloudTexLoaded && cloudTexture.ID > 0 {
			rl.UnloadTexture(cloudTexture)
		}
		cloudTexLoaded = false
		if path != "" {
			tex := rl.LoadTexture(path)
			if tex.ID > 0 {
				cloudTexture = tex
				cloudTexLoaded = true
			}
		}
		cloudMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("SetCloudSpeed", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("SetCloudSpeed(value) requires 1 argument")
		}
		cloudMu.Lock()
		cloudSpeed = toFloat32(args[0])
		cloudMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("SetCloudDensity", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("SetCloudDensity(value) requires 1 argument")
		}
		cloudMu.Lock()
		cloudDensity = toFloat32(args[0])
		if cloudDensity < 0 {
			cloudDensity = 0
		}
		if cloudDensity > 1 {
			cloudDensity = 1
		}
		cloudMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("SetCloudHeight", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("SetCloudHeight(value) requires 1 argument")
		}
		cloudMu.Lock()
		cloudHeight = toFloat32(args[0])
		cloudMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("SetCloudColor", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("SetCloudColor(r, g, b) requires 3 arguments")
		}
		cloudMu.Lock()
		cloudColorR = toFloat32(args[0])
		cloudColorG = toFloat32(args[1])
		cloudColorB = toFloat32(args[2])
		if cloudColorR > 1 {
			cloudColorR /= 255
		}
		if cloudColorG > 1 {
			cloudColorG /= 255
		}
		if cloudColorB > 1 {
			cloudColorB /= 255
		}
		cloudMu.Unlock()
		return nil, nil
	})

	// --- Sun ---
	v.RegisterForeign("SetSunDirection", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("SetSunDirection(x, y, z) requires 3 arguments")
		}
		sunMu.Lock()
		sunDirX = toFloat32(args[0])
		sunDirY = toFloat32(args[1])
		sunDirZ = toFloat32(args[2])
		sunMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("SetSunColor", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("SetSunColor(r, g, b) requires 3 arguments")
		}
		sunMu.Lock()
		sunColorR = toFloat32(args[0])
		sunColorG = toFloat32(args[1])
		sunColorB = toFloat32(args[2])
		if sunColorR > 1 {
			sunColorR /= 255
		}
		if sunColorG > 1 {
			sunColorG /= 255
		}
		if sunColorB > 1 {
			sunColorB /= 255
		}
		sunMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("SetSunIntensity", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("SetSunIntensity(value) requires 1 argument")
		}
		sunMu.Lock()
		sunIntensity = toFloat32(args[0])
		if sunIntensity < 0 {
			sunIntensity = 0
		}
		sunMu.Unlock()
		return nil, nil
	})

	// --- World time ---
	v.RegisterForeign("SetWorldTime", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("SetWorldTime(hours) requires 1 argument")
		}
		worldTimeMu.Lock()
		worldTime = toFloat32(args[0])
		for worldTime < 0 {
			worldTime += 24
		}
		for worldTime >= 24 {
			worldTime -= 24
		}
		worldTimeMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("GetWorldTime", func(args []interface{}) (interface{}, error) {
		worldTimeMu.Lock()
		t := worldTime
		worldTimeMu.Unlock()
		return float64(t), nil
	})
	v.RegisterForeign("SetWorldTimeScale", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("SetWorldTimeScale(value) requires 1 argument")
		}
		worldTimeMu.Lock()
		worldTimeScale = toFloat32(args[0])
		if worldTimeScale < 0 {
			worldTimeScale = 0
		}
		worldTimeMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("SetWeatherPreset", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("SetWeatherPreset(name) requires 1 argument")
		}
		worldTimeMu.Lock()
		weatherPreset = toString(args[0])
		worldTimeMu.Unlock()
		return nil, nil
	})
}

// CloudsOn returns whether clouds are enabled.
func CloudsOn() bool {
	cloudMu.Lock()
	on := cloudsOn
	cloudMu.Unlock()
	return on
}

// CloudTex returns cloud texture if loaded.
func CloudTex() (rl.Texture2D, bool) {
	cloudMu.Lock()
	loaded := cloudTexLoaded && cloudTexture.ID > 0
	tex := cloudTexture
	cloudMu.Unlock()
	return tex, loaded
}

// SkyboxTex returns the loaded skybox texture for use with raylib DrawSkybox.
// Returns (texture, true) if loaded, (zero, false) otherwise.
func SkyboxTex() (rl.Texture2D, bool) {
	skyboxMu.Lock()
	defer skyboxMu.Unlock()
	if skyboxLoaded && skyboxTex.ID > 0 {
		return skyboxTex, true
	}
	return rl.Texture2D{}, false
}
