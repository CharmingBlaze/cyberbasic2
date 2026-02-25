// Package raylib: distance fog for 3D (shader-based).
// Use SetFog or SetFogDensity/SetFogColor, then BeginFog() before 3D draw, EndFog() after.
package raylib

import (
	"cyberbasic/compiler/vm"
	"fmt"
	"sync"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	fogVS = `#version 330
in vec3 vertexPosition;
in vec2 vertexTexCoord;
in vec3 vertexNormal;
in vec4 vertexColor;
out vec2 fragTexCoord;
out vec4 fragColor;
out vec3 fragWorldPos;
uniform mat4 matProjection;
uniform mat4 matView;
uniform mat4 matModel;
void main() {
  fragTexCoord = vertexTexCoord;
  fragColor = vertexColor;
  vec4 w = matModel * vec4(vertexPosition, 1.0);
  fragWorldPos = w.xyz;
  gl_Position = matProjection * matView * w;
}
`
	fogFS = `#version 330
in vec2 fragTexCoord;
in vec4 fragColor;
in vec3 fragWorldPos;
uniform vec3 cameraPosition;
uniform float fogDensity;
uniform vec3 fogColor;
uniform sampler2D texture0;
out vec4 finalColor;
void main() {
  vec4 texColor = texture(texture0, fragTexCoord);
  if (texColor.a < 0.01) discard;
  vec4 baseColor = texColor * fragColor;
  float dist = length(fragWorldPos - cameraPosition);
  float fogAmount = 1.0 - exp(-dist * fogDensity);
  fogAmount = clamp(fogAmount, 0.0, 1.0);
  finalColor = mix(baseColor, vec4(fogColor, 1.0), fogAmount);
}
`
)

var (
	fogMu          sync.Mutex
	fogEnabled     bool
	fogDensity     float32 = 0.02
	fogColorR      float32 = 0.5
	fogColorG      float32 = 0.6
	fogColorB      float32 = 0.7
	fogShader      rl.Shader
	fogShaderReady bool
	fogCamLoc      int32 = -1
	fogDensLoc     int32 = -1
	fogColorLoc    int32 = -1
)

func initFogShader() {
	if fogShaderReady {
		return
	}
	fogMu.Lock()
	defer fogMu.Unlock()
	if fogShaderReady {
		return
	}
	sh := rl.LoadShaderFromMemory(fogVS, fogFS)
	if !rl.IsShaderValid(sh) {
		return
	}
	fogShader = sh
	fogCamLoc = rl.GetShaderLocation(fogShader, "cameraPosition")
	fogDensLoc = rl.GetShaderLocation(fogShader, "fogDensity")
	fogColorLoc = rl.GetShaderLocation(fogShader, "fogColor")
	fogShaderReady = true
}

func registerFog(v *vm.VM) {
	v.RegisterForeign("SetFog", func(args []interface{}) (interface{}, error) {
		if len(args) < 5 {
			return nil, fmt.Errorf("SetFog requires (enable, density, r, g, b)")
		}
		fogMu.Lock()
		fogEnabled = toInt32(args[0]) != 0
		fogDensity = toFloat32(args[1])
		fogColorR = toFloat32(args[2]) / 255.0
		fogColorG = toFloat32(args[3]) / 255.0
		fogColorB = toFloat32(args[4]) / 255.0
		fogMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("SetFogDensity", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("SetFogDensity requires (density)")
		}
		fogMu.Lock()
		fogDensity = toFloat32(args[0])
		fogMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("SetFogColor", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("SetFogColor requires (r, g, b)")
		}
		fogMu.Lock()
		fogColorR = toFloat32(args[0]) / 255.0
		fogColorG = toFloat32(args[1]) / 255.0
		fogColorB = toFloat32(args[2]) / 255.0
		fogMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("EnableFog", func(args []interface{}) (interface{}, error) {
		fogMu.Lock()
		fogEnabled = true
		fogMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("DisableFog", func(args []interface{}) (interface{}, error) {
		fogMu.Lock()
		fogEnabled = false
		fogMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("IsFogEnabled", func(args []interface{}) (interface{}, error) {
		fogMu.Lock()
		en := fogEnabled
		fogMu.Unlock()
		if en {
			return 1, nil
		}
		return 0, nil
	})
	v.RegisterForeign("BeginFog", func(args []interface{}) (interface{}, error) {
		fogMu.Lock()
		en := fogEnabled
		fogMu.Unlock()
		if !en {
			return nil, nil
		}
		initFogShader()
		if !fogShaderReady {
			return nil, nil
		}
		rl.BeginShaderMode(fogShader)
		camPos := []float32{camera3D.Position.X, camera3D.Position.Y, camera3D.Position.Z}
		rl.SetShaderValue(fogShader, fogCamLoc, camPos, rl.ShaderUniformVec3)
		rl.SetShaderValue(fogShader, fogDensLoc, []float32{fogDensity}, rl.ShaderUniformFloat)
		fogCol := []float32{fogColorR, fogColorG, fogColorB}
		rl.SetShaderValue(fogShader, fogColorLoc, fogCol, rl.ShaderUniformVec3)
		return nil, nil
	})
	v.RegisterForeign("EndFog", func(args []interface{}) (interface{}, error) {
		fogMu.Lock()
		en := fogEnabled
		fogMu.Unlock()
		if !en {
			return nil, nil
		}
		if fogShaderReady {
			rl.EndShaderMode()
		}
		return nil, nil
	})
}
