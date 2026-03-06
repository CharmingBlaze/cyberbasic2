// Package dbp: DBP-style material registry.
//
// MakeMaterial(id) creates materials that can be applied to objects
// via ApplyMaterial(id, objectID). Uses raylib Material internally.
package dbp

import (
	"fmt"
	"sync"

	"cyberbasic/compiler/vm"
	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
	materials   = make(map[int]rl.Material)
	materialsMu sync.Mutex
)

func registerMaterials(v *vm.VM) {
	v.RegisterForeign("CreatePBRMaterial", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("CreatePBRMaterial(id) requires 1 argument")
		}
		id := toInt(args[0])
		mat := rl.LoadMaterialDefault()
		if mat.Maps != nil {
			if metalMap := mat.GetMap(rl.MapMetalness); metalMap != nil {
				metalMap.Value = 0
			}
			if roughMap := mat.GetMap(rl.MapRoughness); roughMap != nil {
				roughMap.Value = 0.5
			}
		}
		materialsMu.Lock()
		materials[id] = mat
		materialsMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("SetMaterialPBR", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("SetMaterialPBR(id, metallic, roughness, normalTextureID) requires 4 arguments")
		}
		id := toInt(args[0])
		metallic := toFloat32(args[1])
		roughness := toFloat32(args[2])
		normTexID := toInt(args[3])
		materialsMu.Lock()
		mat, ok := materials[id]
		materialsMu.Unlock()
		if !ok {
			return nil, fmt.Errorf("unknown material id %d", id)
		}
		if mat.Maps != nil {
			if metalMap := mat.GetMap(rl.MapMetalness); metalMap != nil {
				metalMap.Value = metallic
			}
			if roughMap := mat.GetMap(rl.MapRoughness); roughMap != nil {
				roughMap.Value = roughness
			}
			if normTexID > 0 {
				texturesMu.Lock()
				tex, texOk := textures[normTexID]
				texturesMu.Unlock()
				if texOk {
					rl.SetMaterialTexture(&mat, rl.MapNormal, tex)
				}
			}
			materialsMu.Lock()
			materials[id] = mat
			materialsMu.Unlock()
		}
		return nil, nil
	})
	v.RegisterForeign("MakeMaterial", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("MakeMaterial(id) requires 1 argument")
		}
		id := toInt(args[0])
		mat := rl.LoadMaterialDefault()
		materialsMu.Lock()
		materials[id] = mat
		materialsMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("SetMaterialColor", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("SetMaterialColor(id, r, g, b) requires 4 arguments")
		}
		id := toInt(args[0])
		r, g, b := toInt(args[1])&0xff, toInt(args[2])&0xff, toInt(args[3])&0xff
		materialsMu.Lock()
		mat, ok := materials[id]
		if ok && mat.Maps != nil {
			mat.Maps.Color = rl.NewColor(uint8(r), uint8(g), uint8(b), 255)
			materials[id] = mat
		}
		materialsMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("SetMaterialMetalness", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("SetMaterialMetalness(id, value) requires 2 arguments")
		}
		// raylib Material may not have metalness; store for future use
		return nil, nil
	})
	v.RegisterForeign("SetMaterialRoughness", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("SetMaterialRoughness(id, value) requires 2 arguments")
		}
		return nil, nil
	})
	v.RegisterForeign("SetMaterialTexture", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("SetMaterialTexture(id, textureID) requires 2 arguments")
		}
		id := toInt(args[0])
		texID := toInt(args[1])
		materialsMu.Lock()
		mat, ok := materials[id]
		materialsMu.Unlock()
		if !ok {
			return nil, fmt.Errorf("unknown material id %d", id)
		}
		texturesMu.Lock()
		tex, texOk := textures[texID]
		texturesMu.Unlock()
		if !texOk {
			return nil, fmt.Errorf("unknown texture id %d", texID)
		}
		if mat.Maps != nil {
			materialsMu.Lock()
			mat.Maps.Texture = tex
			materials[id] = mat
			materialsMu.Unlock()
		}
		return nil, nil
	})
	v.RegisterForeign("ApplyMaterial", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("ApplyMaterial(id, objectID) requires 2 arguments")
		}
		matID := toInt(args[0])
		objID := toInt(args[1])
		materialsMu.Lock()
		mat, matOk := materials[matID]
		materialsMu.Unlock()
		if !matOk {
			return nil, fmt.Errorf("unknown material id %d", matID)
		}
		objectsMu.Lock()
		obj, objOk := objects[objID]
		if objOk && obj.model.MeshCount > 0 && obj.model.Materials != nil {
			*obj.model.Materials = mat
		}
		objectsMu.Unlock()
		if !objOk {
			return nil, fmt.Errorf("unknown object id %d", objID)
		}
		return nil, nil
	})
	v.RegisterForeign("DeleteMaterial", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("DeleteMaterial(id) requires 1 argument")
		}
		id := toInt(args[0])
		materialsMu.Lock()
		mat, ok := materials[id]
		if ok {
			delete(materials, id)
		}
		materialsMu.Unlock()
		if ok && mat.Maps != nil {
			rl.UnloadMaterial(mat)
		}
		return nil, nil
	})
	v.RegisterForeign("MaterialExists", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("MaterialExists(id) requires 1 argument")
		}
		id := toInt(args[0])
		materialsMu.Lock()
		_, ok := materials[id]
		materialsMu.Unlock()
		if ok {
			return 1, nil
		}
		return 0, nil
	})
}
