// Package scene provides scene/window concepts: CreateScene, LoadScene, UnloadScene, SetCurrentScene, GetCurrentScene, SaveScene, LoadSceneFromFile.
package scene

import (
	"cyberbasic/compiler/vm"
	"encoding/json"
	"fmt"
	"os"
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
}
