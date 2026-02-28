// Package world provides WorldSave, WorldLoad, and JSON export/import for CyberBasic.
package world

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

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

func exportObjects() map[string]objects.ObjectExport {
	return objects.ExportForSave()
}

func importObjects(data map[string]objects.ObjectExport) {
	objects.ImportFromLoad(data)
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
		return nil, nil
	})
	v.RegisterForeign("WorldStreamSetRadius", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, nil
		}
		return nil, nil
	})
	v.RegisterForeign("WorldStreamSetCenter", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, nil
		}
		return nil, nil
	})
	v.RegisterForeign("WorldLoadChunk", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("WorldLoadChunk requires (chunkX, chunkZ)")
		}
		return nil, nil
	})
	v.RegisterForeign("WorldUnloadChunk", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, nil
		}
		return nil, nil
	})
	v.RegisterForeign("WorldIsChunkLoaded", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return false, nil
		}
		return false, nil
	})
	v.RegisterForeign("WorldGetLoadedChunks", func(args []interface{}) (interface{}, error) {
		return []interface{}{}, nil
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
}
