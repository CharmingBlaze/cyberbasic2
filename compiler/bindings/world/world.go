// Package world provides WorldSave, WorldLoad, and JSON export/import for CyberBasic.
package world

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"sync"

	"cyberbasic/compiler/bindings/modfacade"
	"cyberbasic/compiler/bindings/objects"
	"cyberbasic/compiler/vm"
)

const saveVersion = 1

// WorldData is the top-level structure saved to disk.
type WorldData struct {
	Version int                         `json:"version"`
	Objects map[string]objects.ObjectExport `json:"objects,omitempty"`
}

var worldMu sync.Mutex

// World streaming state
type chunkKey struct {
	X, Z int
}
var (
	worldStreamEnabled bool
	worldStreamRadius  float64
	worldStreamCenterX, worldStreamCenterY, worldStreamCenterZ float64
	worldLoadedChunks  = make(map[chunkKey]bool)
	worldStreamMu      sync.RWMutex
)

func exportObjects() map[string]objects.ObjectExport {
	return objects.ExportForSave()
}

func importObjects(data map[string]objects.ObjectExport) {
	objects.ImportFromLoad(data)
}

func toInt(v interface{}) int {
	switch x := v.(type) {
	case int:
		return x
	case int32:
		return int(x)
	case float64:
		return int(x)
	default:
		return 0
	}
}

func toFloat64(v interface{}) float64 {
	switch x := v.(type) {
	case int:
		return float64(x)
	case float64:
		return x
	case float32:
		return float64(x)
	default:
		return 0
	}
}

// RegisterWorld registers WorldSave, WorldLoad, WorldExportJSON, WorldImportJSON with the VM.
func RegisterWorld(v *vm.VM) {
	v.RegisterForeign("WorldSave", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("WorldSave requires (path)")
		}
		path := fmt.Sprint(args[0])
		worldMu.Lock()
		defer worldMu.Unlock()
		data := WorldData{Version: saveVersion, Objects: exportObjects()}
		raw, err := json.MarshalIndent(data, "", "  ")
		if err != nil {
			return nil, err
		}
		return nil, os.WriteFile(path, raw, 0644)
	})

	v.RegisterForeign("WorldLoad", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("WorldLoad requires (path)")
		}
		path := fmt.Sprint(args[0])
		raw, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}
		var data WorldData
		if err := json.Unmarshal(raw, &data); err != nil {
			return nil, err
		}
		worldMu.Lock()
		defer worldMu.Unlock()
		if data.Objects != nil {
			importObjects(data.Objects)
		}
		return nil, nil
	})

	v.RegisterForeign("WorldExportJSON", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("WorldExportJSON requires (path)")
		}
		path := fmt.Sprint(args[0])
		worldMu.Lock()
		defer worldMu.Unlock()
		data := WorldData{Version: saveVersion, Objects: exportObjects()}
		raw, err := json.MarshalIndent(data, "", "  ")
		if err != nil {
			return nil, err
		}
		return nil, os.WriteFile(path, raw, 0644)
	})

	v.RegisterForeign("WorldStreamEnable", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("WorldStreamEnable requires (flag)")
		}
		worldStreamMu.Lock()
		worldStreamEnabled = toInt(args[0]) != 0
		worldStreamMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("WorldStreamSetRadius", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, nil
		}
		worldStreamMu.Lock()
		worldStreamRadius = toFloat64(args[0])
		worldStreamMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("WorldStreamSetCenter", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, nil
		}
		worldStreamMu.Lock()
		worldStreamCenterX = toFloat64(args[0])
		worldStreamCenterY = toFloat64(args[1])
		worldStreamCenterZ = toFloat64(args[2])
		worldStreamMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("WorldLoadChunk", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("WorldLoadChunk requires (chunkX, chunkZ)")
		}
		cx, cz := toInt(args[0]), toInt(args[1])
		worldStreamMu.Lock()
		worldLoadedChunks[chunkKey{X: cx, Z: cz}] = true
		worldStreamMu.Unlock()
		if len(args) >= 3 {
			path := fmt.Sprint(args[2])
			if path != "" {
				if _, err := os.Stat(path); err == nil {
					_, _ = v.CallForeign("LoadLevel", []interface{}{path})
				}
			}
		} else {
			chunkPath := filepath.Join("chunks", "chunk_"+strconv.Itoa(cx)+"_"+strconv.Itoa(cz)+".json")
			if _, err := os.Stat(chunkPath); err == nil {
				_, _ = v.CallForeign("LoadLevel", []interface{}{chunkPath})
			}
		}
		return nil, nil
	})
	v.RegisterForeign("WorldUnloadChunk", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, nil
		}
		cx, cz := toInt(args[0]), toInt(args[1])
		worldStreamMu.Lock()
		delete(worldLoadedChunks, chunkKey{X: cx, Z: cz})
		worldStreamMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("WorldIsChunkLoaded", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return false, nil
		}
		cx, cz := toInt(args[0]), toInt(args[1])
		worldStreamMu.RLock()
		loaded := worldLoadedChunks[chunkKey{X: cx, Z: cz}]
		worldStreamMu.RUnlock()
		return loaded, nil
	})
	v.RegisterForeign("WorldGetLoadedChunks", func(args []interface{}) (interface{}, error) {
		worldStreamMu.RLock()
		out := make([]interface{}, 0, len(worldLoadedChunks)*2)
		for k := range worldLoadedChunks {
			out = append(out, k.X, k.Z)
		}
		worldStreamMu.RUnlock()
		return out, nil
	})

	v.RegisterForeign("WorldImportJSON", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("WorldImportJSON requires (path)")
		}
		path := fmt.Sprint(args[0])
		raw, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}
		var data WorldData
		if err := json.Unmarshal(raw, &data); err != nil {
			return nil, err
		}
		worldMu.Lock()
		defer worldMu.Unlock()
		if data.Objects != nil {
			importObjects(data.Objects)
		}
		return nil, nil
	})

	v.SetGlobal("world", modfacade.New(v, worldV2))
}
