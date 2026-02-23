// Package std provides standard library bindings: file (ReadFile, WriteFile, DeleteFile),
// JSON (LoadJSON, GetJSONKey, SaveJSON), HTTP (HttpGet, HttpPost, DownloadFile), and HELP.
package std

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"

	"cyberbasic/compiler/vm"
)

func toString(v interface{}) string {
	if v == nil {
		return ""
	}
	return fmt.Sprint(v)
}

var (
	jsonStore   = make(map[string]interface{})
	jsonCounter int
	jsonMu      sync.Mutex
)

// RegisterStd registers file, JSON, HTTP, and HELP bindings with the VM.
func RegisterStd(v *vm.VM) {
	// --- File ---
	v.RegisterForeign("ReadFile", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("ReadFile(path) requires 1 argument")
		}
		data, err := os.ReadFile(toString(args[0]))
		if err != nil {
			return nil, err
		}
		return string(data), nil
	})
	v.RegisterForeign("WriteFile", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("WriteFile(path, contents) requires 2 arguments")
		}
		err := os.WriteFile(toString(args[0]), []byte(toString(args[1])), 0644)
		return err == nil, err
	})
	v.RegisterForeign("DeleteFile", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("DeleteFile(path) requires 1 argument")
		}
		err := os.Remove(toString(args[0]))
		return err == nil, err
	})

	// --- JSON ---
	v.RegisterForeign("LoadJSON", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("LoadJSON(path) requires 1 argument")
		}
		path := toString(args[0])
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}
		var obj interface{}
		if err := json.Unmarshal(data, &obj); err != nil {
			return nil, err
		}
		jsonMu.Lock()
		jsonCounter++
		id := fmt.Sprintf("json_%d", jsonCounter)
		jsonStore[id] = obj
		jsonMu.Unlock()
		return id, nil
	})
	v.RegisterForeign("LoadJSONFromString", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("LoadJSONFromString(str) requires 1 argument")
		}
		s := toString(args[0])
		var obj interface{}
		if err := json.Unmarshal([]byte(s), &obj); err != nil {
			return nil, err
		}
		jsonMu.Lock()
		jsonCounter++
		id := fmt.Sprintf("json_%d", jsonCounter)
		jsonStore[id] = obj
		jsonMu.Unlock()
		return id, nil
	})
	v.RegisterForeign("GetJSONKey", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("GetJSONKey(handle, key) requires 2 arguments")
		}
		id := toString(args[0])
		key := toString(args[1])
		jsonMu.Lock()
		obj, ok := jsonStore[id]
		jsonMu.Unlock()
		if !ok {
			return nil, fmt.Errorf("unknown JSON handle: %s", id)
		}
		m, ok := obj.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("JSON handle is not an object")
		}
		val, ok := m[key]
		if !ok {
			return nil, nil
		}
		return val, nil
	})
	v.RegisterForeign("SaveJSON", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("SaveJSON(path, handle) requires 2 arguments")
		}
		path := toString(args[0])
		id := toString(args[1])
		jsonMu.Lock()
		obj, ok := jsonStore[id]
		jsonMu.Unlock()
		if !ok {
			return nil, fmt.Errorf("unknown JSON handle: %s", id)
		}
		data, err := json.MarshalIndent(obj, "", "  ")
		if err != nil {
			return nil, err
		}
		if err := os.WriteFile(path, data, 0644); err != nil {
			return nil, err
		}
		return true, nil
	})

	// --- HTTP ---
	v.RegisterForeign("HttpGet", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("HttpGet(url) requires 1 argument")
		}
		resp, err := http.Get(toString(args[0]))
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		return string(body), nil
	})
	v.RegisterForeign("HttpPost", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("HttpPost(url, body) requires 2 arguments")
		}
		urlStr := toString(args[0])
		bodyStr := toString(args[1])
		resp, err := http.Post(urlStr, "text/plain", io.NopCloser(strings.NewReader(bodyStr)))
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		out, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		return string(out), nil
	})
	v.RegisterForeign("DownloadFile", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("DownloadFile(url, path) requires 2 arguments")
		}
		resp, err := http.Get(toString(args[0]))
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		f, err := os.Create(toString(args[1]))
		if err != nil {
			return nil, err
		}
		defer f.Close()
		_, err = io.Copy(f, resp.Body)
		return err == nil, err
	})

	// --- Null check ---
	v.RegisterForeign("IsNull", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("IsNull(value) requires 1 argument")
		}
		return args[0] == nil, nil
	})

	// --- HELP ---
	v.RegisterForeign("HELP", func(args []interface{}) (interface{}, error) {
		fmt.Println("CyberBasic API: See API_REFERENCE.md in the project root.")
		fmt.Println("Quick ref: InitWindow, BeginDrawing, ClearBackground, DrawCircle, EndDrawing, WindowShouldClose, CloseWindow, GetFrameTime, SetTargetFPS, IsKeyDown, KEY_W, KEY_ESCAPE, etc.")
		return nil, nil
	})
	v.RegisterForeign("?", func(args []interface{}) (interface{}, error) {
		fmt.Println("CyberBasic API: See API_REFERENCE.md in the project root.")
		return nil, nil
	})
}
