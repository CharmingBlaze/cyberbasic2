// Package renderer: shadow mapping support.
//
// The first shipping implementation intentionally keeps scope tight:
// one shadow map, one active directional light, and low/medium/high presets
// so projects can scale from lower-end to higher-end machines.
package renderer

import (
	"math"
	"strings"
	"sync"

	"cyberbasic/compiler/bindings/raylib"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type ShadowLight struct {
	Type      int
	Position  rl.Vector3
	Direction rl.Vector3
	Range     float32
	Shadows   bool
}

const (
	shadowLightDirectional = 1
	shadowDarkness         = float32(0.55)
)

const shadowDepthVS = `#version 330
in vec3 vertexPosition;
uniform mat4 matProjection;
uniform mat4 matView;
uniform mat4 matModel;
void main() {
  gl_Position = matProjection * matView * matModel * vec4(vertexPosition, 1.0);
}
`

const shadowDepthFS = `#version 330
out vec4 finalColor;
void main() {
  float depth = gl_FragCoord.z;
  finalColor = vec4(depth, depth, depth, 1.0);
}
`

const shadowMainVS = `#version 330
in vec3 vertexPosition;
in vec2 vertexTexCoord;
in vec4 vertexColor;
out vec2 fragTexCoord;
out vec4 fragColor;
out vec4 fragShadowPos;
uniform mat4 matProjection;
uniform mat4 matView;
uniform mat4 matModel;
uniform mat4 shadowLightVP;
void main() {
  vec4 worldPos = matModel * vec4(vertexPosition, 1.0);
  fragTexCoord = vertexTexCoord;
  fragColor = vertexColor;
  fragShadowPos = shadowLightVP * worldPos;
  gl_Position = matProjection * matView * worldPos;
}
`

const shadowMainFS = `#version 330
in vec2 fragTexCoord;
in vec4 fragColor;
in vec4 fragShadowPos;
uniform sampler2D texture0;
uniform sampler2D shadowMap;
uniform float shadowBias;
uniform float shadowDarkness;
out vec4 finalColor;

float readShadowDepth(vec2 uv) {
  return texture(shadowMap, uv).r;
}

float computeShadow(vec4 shadowPos) {
  vec3 proj = shadowPos.xyz / max(shadowPos.w, 0.0001);
  proj = proj * 0.5 + 0.5;
  if (proj.z > 1.0 || proj.x < 0.0 || proj.x > 1.0 || proj.y < 0.0 || proj.y > 1.0) {
    return 1.0;
  }
  float closestDepth = readShadowDepth(proj.xy);
  float currentDepth = proj.z - shadowBias;
  return currentDepth <= closestDepth ? 1.0 : shadowDarkness;
}

void main() {
  vec4 texColor = texture(texture0, fragTexCoord);
  if (texColor.a < 0.01) discard;
  vec4 baseColor = texColor * fragColor;
  float shadowFactor = computeShadow(fragShadowPos);
  finalColor = vec4(baseColor.rgb * shadowFactor, baseColor.a);
}
`

var (
	shadowMapWidth  int32   = 1024
	shadowMapHeight int32   = 1024
	shadowBias      float32 = 0.005
	shadowQuality   string  = "medium"
	shadowMu        sync.RWMutex

	shadowLights      []ShadowLight
	activeShadowLight ShadowLight
	lightViewProj     rl.Matrix
	activeLightValid  bool
	shadowPassActive  bool

	shadowRT      rl.RenderTexture2D
	shadowRTReady bool

	shadowDepthShader      rl.Shader
	shadowDepthShaderReady bool
	shadowMainShader       rl.Shader
	shadowMainShaderReady  bool

	shadowLightVPLoc    int32 = -1
	shadowBiasLoc       int32 = -1
	shadowDarknessLoc   int32 = -1
	shadowMapSamplerLoc int32 = -1
)

func clampShadowSize(v int32) int32 {
	if v < 256 {
		return 256
	}
	if v > 4096 {
		return 4096
	}
	return v
}

func normalizeShadowDirection(dir rl.Vector3) rl.Vector3 {
	if rl.Vector3Length(dir) <= 1e-6 {
		return rl.Vector3{X: 0.45, Y: -0.8, Z: 0.4}
	}
	return rl.Vector3Normalize(dir)
}

func ensureShadowRenderTextureLocked() {
	width := clampShadowSize(shadowMapWidth)
	height := clampShadowSize(shadowMapHeight)
	if shadowRTReady && rl.IsRenderTextureValid(shadowRT) && shadowRT.Texture.Width == width && shadowRT.Texture.Height == height {
		return
	}
	if shadowRTReady && rl.IsRenderTextureValid(shadowRT) {
		rl.UnloadRenderTexture(shadowRT)
	}
	shadowRT = rl.LoadRenderTexture(width, height)
	shadowRTReady = rl.IsRenderTextureValid(shadowRT)
	shadowMapWidth = width
	shadowMapHeight = height
}

func ensureShadowShadersLocked() {
	if !shadowDepthShaderReady {
		sh := rl.LoadShaderFromMemory(shadowDepthVS, shadowDepthFS)
		if rl.IsShaderValid(sh) {
			shadowDepthShader = sh
			shadowDepthShaderReady = true
		}
	}
	if !shadowMainShaderReady {
		sh := rl.LoadShaderFromMemory(shadowMainVS, shadowMainFS)
		if rl.IsShaderValid(sh) {
			shadowMainShader = sh
			shadowLightVPLoc = rl.GetShaderLocation(sh, "shadowLightVP")
			shadowBiasLoc = rl.GetShaderLocation(sh, "shadowBias")
			shadowDarknessLoc = rl.GetShaderLocation(sh, "shadowDarkness")
			shadowMapSamplerLoc = rl.GetShaderLocation(sh, "shadowMap")
			shadowMainShader = sh
			shadowMainShaderReady = true
		}
	}
}

func pickActiveShadowLightLocked() bool {
	activeLightValid = false
	var fallback *ShadowLight
	for i := range shadowLights {
		light := shadowLights[i]
		if light.Type != shadowLightDirectional {
			continue
		}
		if fallback == nil {
			copyLight := light
			fallback = &copyLight
		}
		if light.Shadows {
			activeShadowLight = light
			activeShadowLight.Direction = normalizeShadowDirection(activeShadowLight.Direction)
			activeLightValid = true
			return true
		}
	}
	if fallback != nil {
		activeShadowLight = *fallback
		activeShadowLight.Direction = normalizeShadowDirection(activeShadowLight.Direction)
		activeLightValid = true
	}
	return activeLightValid
}

func buildShadowCamera(light ShadowLight, sceneCam rl.Camera3D) rl.Camera3D {
	dir := normalizeShadowDirection(light.Direction)
	target := sceneCam.Target
	if rl.Vector3Length(rl.Vector3Subtract(sceneCam.Target, sceneCam.Position)) <= 1e-6 {
		target = rl.Vector3{X: 0, Y: 0, Z: 0}
	}
	distance := light.Range
	if distance < 25 {
		distance = 25
	}
	position := rl.Vector3Subtract(target, rl.Vector3Scale(dir, distance*0.8))
	up := rl.Vector3{X: 0, Y: 1, Z: 0}
	if math.Abs(float64(rl.Vector3DotProduct(dir, up))) > 0.98 {
		up = rl.Vector3{X: 0, Y: 0, Z: 1}
	}
	return rl.Camera3D{
		Position:   position,
		Target:     target,
		Up:         up,
		Fovy:       distance,
		Projection: rl.CameraOrthographic,
	}
}

func computeLightViewProjection(lightCam rl.Camera3D) rl.Matrix {
	view := rl.MatrixLookAt(lightCam.Position, lightCam.Target, lightCam.Up)
	halfHeight := lightCam.Fovy * 0.5
	if halfHeight < 10 {
		halfHeight = 10
	}
	aspect := float32(1)
	if shadowMapHeight > 0 {
		aspect = float32(shadowMapWidth) / float32(shadowMapHeight)
		if aspect <= 0 {
			aspect = 1
		}
	}
	halfWidth := halfHeight * aspect
	proj := rl.MatrixOrtho(-halfWidth, halfWidth, -halfHeight, halfHeight, 0.1, halfHeight*8)
	return rl.MatrixMultiply(proj, view)
}

func applyShadowMainUniformsLocked() {
	if !shadowMainShaderReady || !shadowRTReady {
		return
	}
	if shadowLightVPLoc >= 0 {
		rl.SetShaderValueMatrix(shadowMainShader, shadowLightVPLoc, lightViewProj)
	}
	if shadowBiasLoc >= 0 {
		rl.SetShaderValue(shadowMainShader, shadowBiasLoc, []float32{shadowBias}, rl.ShaderUniformFloat)
	}
	if shadowDarknessLoc >= 0 {
		rl.SetShaderValue(shadowMainShader, shadowDarknessLoc, []float32{shadowDarkness}, rl.ShaderUniformFloat)
	}
	if shadowMapSamplerLoc >= 0 {
		rl.SetShaderValueTexture(shadowMainShader, shadowMapSamplerLoc, shadowRT.Texture)
	}
}

// SetShadowLights updates the light snapshot used by the renderer to choose a shadow caster.
func SetShadowLights(lights []ShadowLight) {
	shadowMu.Lock()
	shadowLights = append(shadowLights[:0], lights...)
	shadowMu.Unlock()
}

// SetShadowMapSize sets the shadow map resolution.
func SetShadowMapSize(width, height int32) {
	shadowMu.Lock()
	shadowMapWidth = clampShadowSize(width)
	shadowMapHeight = clampShadowSize(height)
	shadowMu.Unlock()
}

// SetShadowBias sets the depth bias to reduce acne and peter-panning.
func SetShadowBias(bias float32) {
	shadowMu.Lock()
	if bias < 0.0001 {
		bias = 0.0001
	}
	if bias > 0.05 {
		bias = 0.05
	}
	shadowBias = bias
	shadowMu.Unlock()
}

// SetShadowQuality applies performance-oriented presets for lower/mid/higher-end hardware.
func SetShadowQuality(name string) {
	quality := strings.ToLower(strings.TrimSpace(name))
	shadowMu.Lock()
	switch quality {
	case "low":
		shadowQuality = "low"
		shadowMapWidth, shadowMapHeight = 512, 512
		shadowBias = 0.008
	case "mid", "medium", "":
		shadowQuality = "medium"
		shadowMapWidth, shadowMapHeight = 1024, 1024
		shadowBias = 0.005
	case "high":
		shadowQuality = "high"
		shadowMapWidth, shadowMapHeight = 2048, 2048
		shadowBias = 0.0035
	default:
		shadowQuality = quality
		shadowMapWidth, shadowMapHeight = 1024, 1024
		shadowBias = 0.005
	}
	shadowMu.Unlock()
}

// ShadowQuality returns the current quality preset label.
func ShadowQuality() string {
	shadowMu.RLock()
	defer shadowMu.RUnlock()
	return shadowQuality
}

// IsShadowPassActive reports whether the engine is currently rendering the light-view shadow pass.
func IsShadowPassActive() bool {
	shadowMu.RLock()
	defer shadowMu.RUnlock()
	return shadowPassActive
}

// IsShadowLightingActive reports whether the main pass can sample a valid shadow map this frame.
func IsShadowLightingActive() bool {
	shadowMu.RLock()
	defer shadowMu.RUnlock()
	return raylib.ShadowsEnabled() && activeLightValid && shadowRTReady && shadowMainShaderReady && !shadowPassActive
}

// DepthShader returns the shader used for the depth-encoding shadow pass.
func DepthShader() (rl.Shader, bool) {
	shadowMu.RLock()
	defer shadowMu.RUnlock()
	return shadowDepthShader, shadowDepthShaderReady
}

// ShadowShader returns the shader used in the main pass to sample shadows.
func ShadowShader() (rl.Shader, bool) {
	shadowMu.RLock()
	defer shadowMu.RUnlock()
	return shadowMainShader, shadowMainShaderReady
}

// PrepareShadowShader updates per-frame uniforms for the main-pass shadow shader.
func PrepareShadowShader() {
	shadowMu.Lock()
	defer shadowMu.Unlock()
	applyShadowMainUniformsLocked()
}

// RenderShadowPass renders the scene from the active light's point of view into the shadow map.
func RenderShadowPass() {
	if !raylib.ShadowsEnabled() {
		return
	}
	var lightCam rl.Camera3D
	shadowMu.Lock()
	ensureShadowRenderTextureLocked()
	ensureShadowShadersLocked()
	if !shadowRTReady || !shadowDepthShaderReady || !pickActiveShadowLightLocked() {
		shadowMu.Unlock()
		return
	}
	sceneCam := raylib.GetCamera3D()
	lightCam = buildShadowCamera(activeShadowLight, sceneCam)
	lightViewProj = computeLightViewProjection(lightCam)
	shadowPassActive = true
	shadowMu.Unlock()

	rl.BeginTextureMode(shadowRT)
	rl.ClearBackground(rl.White)
	rl.BeginMode3D(lightCam)
	if draw3DFn != nil {
		draw3DFn()
	}
	rl.EndMode3D()
	rl.EndTextureMode()
	shadowMu.Lock()
	shadowPassActive = false
	applyShadowMainUniformsLocked()
	shadowMu.Unlock()
}

// ShadowMapTexture returns the encoded depth texture generated during the shadow pass.
func ShadowMapTexture() rl.Texture2D {
	shadowMu.RLock()
	defer shadowMu.RUnlock()
	if shadowRTReady {
		return shadowRT.Texture
	}
	return rl.Texture2D{}
}
