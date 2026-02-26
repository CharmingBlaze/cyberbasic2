// Package raylib binds raylib-go to the CyberBasic VM as a foreign API.
// BASIC can call InitWindow(800, 450, "Title"), DrawRectangle(...), etc. (no RL. namespace)
package raylib

import (
	"cyberbasic/compiler/vm"
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
	// seedableRand is used by SetRandomSeed/GetRandomValue when we want reproducible randomness.
	seedableRand   *rand.Rand
	seedableRandMu sync.Mutex
	textures   = make(map[string]rl.Texture2D)
	texCounter int
	texMu      sync.Mutex

	models       = make(map[string]rl.Model)
	modelCounter int
	modelMu      sync.Mutex

	camera3D   rl.Camera3D
	camera2D   rl.Camera2D
	cameras    = make(map[string]rl.Camera3D)
	camCounter int
	camMu      sync.Mutex

	lightIds   = make(map[string]bool)
	lightCtr   int
	lightMu    sync.Mutex
	ambientR   float32
	ambientG   float32
	ambientB   float32
	lightingOn bool

	renderTextures   = make(map[string]rl.RenderTexture2D)
	renderTexCounter int
	renderTexMu      sync.Mutex

	shaders       = make(map[string]rl.Shader)
	shaderCounter int
	shaderMu      sync.Mutex

	sounds       = make(map[string]rl.Sound)
	soundCounter int
	soundMu      sync.Mutex

	waves       = make(map[string]rl.Wave)
	waveCounter int
	waveMu      sync.Mutex

	music        = make(map[string]rl.Music)
	musicCounter int
	musicMu      sync.Mutex

	audioStreams       = make(map[string]rl.AudioStream)
	audioStreamCounter int
	audioStreamMu      sync.Mutex

	lastWaveSamples   []float32
	lastWaveSamplesMu sync.Mutex

	fonts       = make(map[string]rl.Font)
	fontCounter int
	fontMu      sync.Mutex

	meshes       = make(map[string]rl.Mesh)
	meshCounter  int
	meshMu       sync.Mutex
	materials       = make(map[string]rl.Material)
	materialCounter int
	materialMu      sync.Mutex

	animations            = make(map[string]rl.ModelAnimation)
	animCounter           int
	animMu                sync.Mutex
	lastLoadedAnimIds     []string
	lastLoadMaterialsCount int

	images       = make(map[string]*rl.Image)
	imageCounter int
	imageMu      sync.Mutex

	lastLoadImageAnimFrames int32
	lastImageColors         []rl.Color
	lastImageColorsMu       sync.Mutex
	lastImagePalette        []rl.Color
	lastImagePaletteMu      sync.Mutex

	lastFontData   []rl.GlyphInfo
	lastFontDataMu sync.Mutex

	lastCodepoints   []rune
	lastCodepointsMu sync.Mutex
	lastTextSplit   []string
	lastTextSplitMu sync.Mutex

	// Orbit camera state (used by CameraZoom, CameraRotate(dx,dy), UpdateCamera, MouseOrbitCamera)
	orbitTargetX   float32
	orbitTargetY   float32
	orbitTargetZ   float32
	orbitAngle     float32
	orbitPitch     float32
	orbitDistance  float32
	orbitStateMu   sync.Mutex

	// FPS MouseLook state (yaw, pitch in radians)
	mouseLookYaw   float32
	mouseLookPitch float32
	mouseLookMu    sync.Mutex

	// Per-model state for simplified DrawModel / RotateModel / SetModelColor
	modelColors   = make(map[string]rl.Color)
	modelAngles   = make(map[string]float32) // radians
	modelStateMu  sync.Mutex

	// Camera clip plane (stored for custom projection if needed)
	cameraNearZ float32
	cameraFarZ  float32

	// Camera shake state
	cameraShakeAmount    float32
	cameraShakeDuration  float32
	cameraShakeMu        sync.Mutex

	// CollisionBox(id): center (x,y,z) + half-extents (hw,hh,hd)
	collisionBoxes   = make(map[string]struct{ Cx, Cy, Cz, Hw, Hh, Hd float64 })
	collisionBoxSeq  int
	collisionBoxMu   sync.Mutex

	// CreateLight / SetLight* state
	lightData   = make(map[string]*struct {
		Type       int
		X, Y, Z    float32
		R, G, B    uint8
		Intensity  float32
		DirX, DirY, DirZ float32
	})
	lightDataMu sync.Mutex
	shadowsEnabled bool

	// RemoveShader() ends current shader mode
	currentShaderId string
	currentShaderMu sync.Mutex

	// Skybox / sky
	skyboxEnabled bool
	skyColorR, skyColorG, skyColorB uint8 = 135, 206, 235

	// Post-processing (state only; actual effects need RT + shaders)
	bloomEnabled      bool
	bloomIntensity    float32
	motionBlurEnabled bool
	crtFilterEnabled  bool
	pixelateSize      int32

	// Terrain: id -> height grid and optional model
	terrainHeights   = make(map[string][]float32)
	terrainWidth     = make(map[string]int)
	terrainDepth     = make(map[string]int)
	terrainScale     = make(map[string]float32)
	terrainTexId     = make(map[string]string)
	terrainMaterial  = make(map[string]string)
	terrainSeq       int
	terrainMu        sync.Mutex
	terrainBrushSize     float32 = 3
	terrainBrushStrength float32 = 0.1
	terrainUndoStack     = make(map[string][]float32) // terrainId -> previous heights copy

	// Skybox: cubemap texture id (optional); drawing uses SetSkyColor clear or cubemap
	skyboxTexId string
)

// getRand returns the seedable RNG, creating it with a time-based seed if never set.
func getRand() *rand.Rand {
	seedableRandMu.Lock()
	defer seedableRandMu.Unlock()
	if seedableRand == nil {
		seedableRand = rand.New(rand.NewSource(time.Now().UnixNano()))
	}
	return seedableRand
}

// setRandSeed seeds (or reseeds) our custom RNG so GetRandomValue is reproducible.
func setRandSeed(seed int64) {
	seedableRandMu.Lock()
	defer seedableRandMu.Unlock()
	if seedableRand == nil {
		seedableRand = rand.New(rand.NewSource(seed))
	} else {
		seedableRand.Seed(seed)
	}
}

func toInt32(v interface{}) int32 {
	switch x := v.(type) {
	case int:
		return int32(x)
	case float64:
		return int32(x)
	case string:
		n, _ := strconv.Atoi(x)
		return int32(n)
	default:
		return 0
	}
}

func toFloat32(v interface{}) float32 {
	switch x := v.(type) {
	case int:
		return float32(x)
	case float64:
		return float32(x)
	case string:
		f, _ := strconv.ParseFloat(x, 32)
		return float32(f)
	default:
		return 0
	}
}

func toString(v interface{}) string {
	if v == nil {
		return ""
	}
	return fmt.Sprint(v)
}

func toFloat64(v interface{}) float64 {
	f, _ := toFloat64Safe(v)
	return f
}

func toFloat64Safe(v interface{}) (float64, bool) {
	switch x := v.(type) {
	case int:
		return float64(x), true
	case float64:
		return x, true
	case string:
		f, err := strconv.ParseFloat(x, 64)
		return f, err == nil
	default:
		return 0, false
	}
}

func toUint8(v interface{}) uint8 {
	switch x := v.(type) {
	case int:
		return uint8(x)
	case float64:
		return uint8(x)
	default:
		return 0
	}
}

// colorToPacked returns the format draw functions expect for a single int: R<<16|G<<8|B (A=255).
func colorToPacked(c rl.Color) int {
	return int(c.R)<<16 | int(c.G)<<8 | int(c.B)
}

// argsToColor reads r,g,b,a from args starting at offset; returns White if not enough args.
func argsToColor(args []interface{}, offset int) rl.Color {
	if len(args) < offset+4 {
		return rl.White
	}
	return rl.NewColor(toUint8(args[offset]), toUint8(args[offset+1]), toUint8(args[offset+2]), toUint8(args[offset+3]))
}

// argsToVector2Slice reads count Vector2s from args starting at startIndex (flat x1,y1, x2,y2, ...). Returns nil if not enough args.
func argsToVector2Slice(args []interface{}, startIndex, count int) []rl.Vector2 {
	if count <= 0 || len(args) < startIndex+count*2 {
		return nil
	}
	out := make([]rl.Vector2, count)
	for i := 0; i < count; i++ {
		o := startIndex + i*2
		out[i] = rl.Vector2{X: toFloat32(args[o]), Y: toFloat32(args[o+1])}
	}
	return out
}

// RegisterRaylib registers all raylib-go functions with the VM (modular: core, input, shapes, text, textures, 3d, audio, fonts, misc, game).
func RegisterRaylib(v *vm.VM) {
	registerFlags(v)
	registerCore(v)
	registerInput(v)
	registerShapes(v)
	registerText(v)
	registerTextures(v)
	registerImages(v)
	register3D(v)
	registerMesh(v)
	registerAudio(v)
	registerFonts(v)
	registerMisc(v)
	registerMath(v)
	registerGame(v)
	registerAnim2D(v)
	registerUI(v)
	registerRaygui(v)
	registerFog(v)
	registerViews(v)
	registerEditor(v)
	registerAdvanced(v)
	registerMultiWindow(v)
	registerHybrid(v)
}
