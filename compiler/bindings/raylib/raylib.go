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

	camera3D rl.Camera3D
	camera2D rl.Camera2D

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
	registerUI(v)
}
