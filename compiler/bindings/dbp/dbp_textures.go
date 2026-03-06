// Package dbp: DBP-style texture registry (id-based).
//
// LoadTexture(id, path) stores textures by integer id for use with
// SetObjectTexture, SetMaterialTexture, etc.
package dbp

import (
	"fmt"
	"sync"

	"cyberbasic/compiler/vm"
	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
	textures   = make(map[int]rl.Texture2D)
	texturesMu sync.Mutex
)

func registerTextures(v *vm.VM) {
	v.RegisterForeign("LoadTexture", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("LoadTexture(id, path) requires 2 arguments")
		}
		id := toInt(args[0])
		path := toString(args[1])
		tex := rl.LoadTexture(path)
		texturesMu.Lock()
		textures[id] = tex
		texturesMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("DeleteTexture", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("DeleteTexture(id) requires 1 argument")
		}
		id := toInt(args[0])
		texturesMu.Lock()
		tex, ok := textures[id]
		delete(textures, id)
		texturesMu.Unlock()
		if ok && tex.ID > 0 {
			rl.UnloadTexture(tex)
		}
		return nil, nil
	})
	v.RegisterForeign("SetTextureFilter", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("SetTextureFilter(id, mode) requires 2 arguments")
		}
		id := toInt(args[0])
		mode := toInt(args[1])
		texturesMu.Lock()
		tex, ok := textures[id]
		texturesMu.Unlock()
		if !ok {
			return nil, fmt.Errorf("unknown texture id %d", id)
		}
		rl.SetTextureFilter(tex, rl.TextureFilterMode(mode))
		return nil, nil
	})
	v.RegisterForeign("SetTextureWrap", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("SetTextureWrap(id, mode) requires 2 arguments")
		}
		id := toInt(args[0])
		mode := toInt(args[1])
		texturesMu.Lock()
		tex, ok := textures[id]
		texturesMu.Unlock()
		if !ok {
			return nil, fmt.Errorf("unknown texture id %d", id)
		}
		rl.SetTextureWrap(tex, rl.TextureWrapMode(mode))
		return nil, nil
	})
}
