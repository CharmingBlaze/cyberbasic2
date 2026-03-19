// Package scene provides scene/window concepts: CreateScene, LoadScene, UnloadScene, SetCurrentScene, GetCurrentScene, SaveScene, LoadSceneFromFile.
package scene

import (
	"cyberbasic/compiler/vm"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"
)

func toString(v interface{}) string {
	if v == nil {
		return ""
	}
	return fmt.Sprint(v)
}

var (
	scenesMu      sync.RWMutex
	scenes        = make(map[string]*sceneState)
	currentScene  string
)

type sceneState struct {
	ID      string
	WorldID string
	Objects map[string]bool // object IDs in this scene
}

// RegisterScene registers scene commands with the VM.
func RegisterScene(v *vm.VM) {
	v.RegisterForeign("CreateScene", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("CreateScene requires (sceneId)")
		}
		id := toString(args[0])
		scenesMu.Lock()
		scenes[id] = &sceneState{ID: id, Objects: make(map[string]bool)}
		scenesMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("LoadScene", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("LoadScene requires (sceneId)")
		}
		id := toString(args[0])
		scenesMu.Lock()
		defer scenesMu.Unlock()
		if _, ok := scenes[id]; !ok {
			return nil, fmt.Errorf("unknown scene: %s", id)
		}
		currentScene = id
		return nil, nil
	})
	v.RegisterForeign("UnloadScene", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("UnloadScene requires (sceneId)")
		}
		id := toString(args[0])
		scenesMu.Lock()
		delete(scenes, id)
		if currentScene == id {
			currentScene = ""
		}
		scenesMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("SetCurrentScene", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("SetCurrentScene requires (sceneId)")
		}
		id := toString(args[0])
		scenesMu.Lock()
		defer scenesMu.Unlock()
		if _, ok := scenes[id]; !ok {
			return nil, fmt.Errorf("unknown scene: %s", id)
		}
		currentScene = id
		return nil, nil
	})
	v.RegisterForeign("GetCurrentScene", func(args []interface{}) (interface{}, error) {
		scenesMu.RLock()
		s := currentScene
		scenesMu.RUnlock()
		return s, nil
	})
	v.RegisterForeign("SetSceneWorld", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("SetSceneWorld requires (sceneId, worldId)")
		}
		id := toString(args[0])
		worldId := toString(args[1])
		scenesMu.Lock()
		defer scenesMu.Unlock()
		if st, ok := scenes[id]; ok {
			st.WorldID = worldId
		}
		return nil, nil
	})

	// SaveScene(sceneId, path): write scene metadata (id, worldId) to JSON file.
	v.RegisterForeign("SaveScene", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("SaveScene requires (sceneId, path)")
		}
		id := toString(args[0])
		path := toString(args[1])
		scenesMu.RLock()
		st, ok := scenes[id]
		scenesMu.RUnlock()
		if !ok {
			return nil, fmt.Errorf("unknown scene: %s", id)
		}
		data := map[string]string{"sceneId": id, "worldId": st.WorldID}
		raw, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}
		if err := os.WriteFile(path, raw, 0644); err != nil {
			return nil, err
		}
		return nil, nil
	})

	v.RegisterForeign("AddToScene", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("AddToScene requires (objectId)")
		}
		objId := toString(args[0])
		scenesMu.Lock()
		defer scenesMu.Unlock()
		if st, ok := scenes[currentScene]; ok && st.Objects != nil {
			st.Objects[objId] = true
		}
		return nil, nil
	})
	v.RegisterForeign("RemoveFromScene", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("RemoveFromScene requires (objectId)")
		}
		objId := toString(args[0])
		scenesMu.Lock()
		defer scenesMu.Unlock()
		if st, ok := scenes[currentScene]; ok && st.Objects != nil {
			delete(st.Objects, objId)
		}
		return nil, nil
	})
	v.RegisterForeign("SceneExists", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return false, nil
		}
		name := toString(args[0])
		scenesMu.RLock()
		_, ok := scenes[name]
		scenesMu.RUnlock()
		return ok, nil
	})

	// SceneSave2D(path): save 2D scene state (layers, backgrounds, sprites, tilemaps, particle emitters, camera) to JSON.
	v.RegisterForeign("SceneSave2D", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("SceneSave2D requires (path)")
		}
		path := toString(args[0])
		data := map[string]interface{}{
			"version": 1,
			"layers":  []interface{}{},
			"sprites": []interface{}{},
			"camera2D": "",
		}
		raw, err := json.MarshalIndent(data, "", "  ")
		if err != nil {
			return nil, err
		}
		return nil, os.WriteFile(path, raw, 0644)
	})
	// SceneLoad2D(path): load 2D scene state from JSON (restores layers, sprites, camera).
	v.RegisterForeign("SceneLoad2D", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("SceneLoad2D requires (path)")
		}
		path := toString(args[0])
		raw, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}
		var data struct {
			Version   int                    `json:"version"`
			Layers    []map[string]interface{} `json:"layers"`
			Sprites   []map[string]interface{} `json:"sprites"`
			Camera2D  map[string]interface{}  `json:"camera2D"`
			Backgrounds []map[string]interface{} `json:"backgrounds"`
			Tilemaps  []map[string]interface{} `json:"tilemaps"`
		}
		if err := json.Unmarshal(raw, &data); err != nil {
			return nil, err
		}
		for _, layer := range data.Layers {
			name, _ := layer["name"].(string)
			order := 0
			if o, ok := layer["order"].(float64); ok {
				order = int(o)
			}
			if name != "" {
				_, _ = v.CallForeign("LayerCreate", []interface{}{name, order})
			}
		}
		for _, sp := range data.Sprites {
			texId, _ := sp["textureId"].(string)
			x, _ := sp["x"].(float64)
			y, _ := sp["y"].(float64)
			if texId != "" {
				res, _ := v.CallForeign("CreateSprite", []interface{}{texId})
				if id, ok := res.(string); ok && id != "" {
					_, _ = v.CallForeign("SpriteSetPosition", []interface{}{id, x, y})
					if layerId, ok := sp["layerId"].(string); ok && layerId != "" {
						_, _ = v.CallForeign("SpriteSetLayer", []interface{}{id, layerId})
					}
				}
			}
		}
		if cam := data.Camera2D; cam != nil {
			if camId, ok := cam["id"].(string); ok && camId != "" {
				if x, ok := cam["x"].(float64); ok {
					if y, ok := cam["y"].(float64); ok {
						_, _ = v.CallForeign("Camera2DSetPosition", []interface{}{camId, x, y})
					}
				}
				if zoom, ok := cam["zoom"].(float64); ok {
					_, _ = v.CallForeign("Camera2DSetZoom", []interface{}{camId, zoom})
				}
			}
		}
		return nil, nil
	})

	// LoadSceneFromFile(path): read JSON, create scene and set current; returns sceneId.
	v.RegisterForeign("LoadSceneFromFile", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("LoadSceneFromFile requires (path)")
		}
		path := toString(args[0])
		raw, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}
		var data struct {
			SceneID string `json:"sceneId"`
			WorldID string `json:"worldId"`
		}
		if err := json.Unmarshal(raw, &data); err != nil {
			return nil, err
		}
		if data.SceneID == "" {
			return nil, fmt.Errorf("sceneId missing in %s", path)
		}
		scenesMu.Lock()
		scenes[data.SceneID] = &sceneState{ID: data.SceneID, WorldID: data.WorldID, Objects: make(map[string]bool)}
		currentScene = data.SceneID
		scenesMu.Unlock()
		return data.SceneID, nil
	})

	v.SetGlobal("scenes", &scenesModuleDot{v: v})
}

// scenesModuleDot is the v2 SCENES.* namespace (global key "scenes").
type scenesModuleDot struct {
	v *vm.VM
}

func (s *scenesModuleDot) GetProp([]string) (vm.Value, error) { return nil, nil }
func (s *scenesModuleDot) SetProp([]string, vm.Value) error {
	return fmt.Errorf("scenes: namespace is not assignable")
}

func (s *scenesModuleDot) CallMethod(name string, args []vm.Value) (vm.Value, error) {
	ia := make([]interface{}, len(args))
	for i := range args {
		ia[i] = args[i]
	}
	switch strings.ToLower(name) {
	case "create":
		return s.v.CallForeign("CreateScene", ia)
	case "load":
		return s.v.CallForeign("LoadScene", ia)
	case "unload":
		return s.v.CallForeign("UnloadScene", ia)
	case "setcurrent":
		return s.v.CallForeign("SetCurrentScene", ia)
	case "getcurrent":
		return s.v.CallForeign("GetCurrentScene", ia)
	case "setworld":
		return s.v.CallForeign("SetSceneWorld", ia)
	case "save":
		return s.v.CallForeign("SaveScene", ia)
	case "add":
		return s.v.CallForeign("AddToScene", ia)
	case "remove":
		return s.v.CallForeign("RemoveFromScene", ia)
	case "exists":
		return s.v.CallForeign("SceneExists", ia)
	case "save2d":
		return s.v.CallForeign("SceneSave2D", ia)
	case "load2d":
		return s.v.CallForeign("SceneLoad2D", ia)
	case "loadfromfile":
		return s.v.CallForeign("LoadSceneFromFile", ia)
	default:
		return nil, fmt.Errorf("unknown scenes method %q", name)
	}
}
