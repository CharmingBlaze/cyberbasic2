// Package raylib: texture atlas (load JSON + texture, get region by name).
package raylib

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"cyberbasic/compiler/vm"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type atlasState struct {
	TextureID string
	Regions   map[string]rl.Rectangle
}

var (
	atlases   = make(map[string]*atlasState)
	atlasSeq  int
	atlasMu   sync.Mutex
)

type atlasJSON struct {
	Texture string              `json:"texture"`
	Regions map[string][]float64 `json:"regions"`
}

func registerAtlas(v *vm.VM) {
	v.RegisterForeign("AtlasLoad", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("AtlasLoad requires (path)")
		}
		path := toString(args[0])
		atlasMu.Lock()
		atlasSeq++
		id := fmt.Sprintf("atlas_%d", atlasSeq)
		atlasMu.Unlock()
		dir := filepath.Dir(path)
		raw, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}
		var data atlasJSON
		if err := json.Unmarshal(raw, &data); err != nil {
			return nil, err
		}
		texPath := data.Texture
		if !filepath.IsAbs(texPath) {
			texPath = filepath.Join(dir, texPath)
		}
		texMu.Lock()
		texCounter++
		texID := fmt.Sprintf("atlas_tex_%d", texCounter)
		tex := rl.LoadTexture(texPath)
		textures[texID] = tex
		texMu.Unlock()
		regions := make(map[string]rl.Rectangle)
		for name, r := range data.Regions {
			if len(r) >= 4 {
				regions[name] = rl.Rectangle{
					X: float32(r[0]), Y: float32(r[1]),
					Width: float32(r[2]), Height: float32(r[3]),
				}
			}
		}
		atlasMu.Lock()
		atlases[id] = &atlasState{TextureID: texID, Regions: regions}
		atlasMu.Unlock()
		return id, nil
	})
	v.RegisterForeign("AtlasGetRegion", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("AtlasGetRegion requires (atlasId, name)")
		}
		id := toString(args[0])
		name := toString(args[1])
		atlasMu.Lock()
		a := atlases[id]
		atlasMu.Unlock()
		if a == nil {
			return nil, fmt.Errorf("unknown atlas: %s", id)
		}
		r, ok := a.Regions[name]
		if !ok {
			return []interface{}{float64(0), float64(0), float64(0), float64(0)}, nil
		}
		return []interface{}{float64(r.X), float64(r.Y), float64(r.Width), float64(r.Height)}, nil
	})
	v.RegisterForeign("AtlasGetTextureId", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("AtlasGetTextureId requires (atlasId)")
		}
		id := toString(args[0])
		atlasMu.Lock()
		a := atlases[id]
		atlasMu.Unlock()
		if a == nil {
			return nil, fmt.Errorf("unknown atlas: %s", id)
		}
		return a.TextureID, nil
	})
}
