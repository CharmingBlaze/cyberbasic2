// Package std provides standard library bindings: file (ReadFile, WriteFile, DeleteFile, CopyFile, ListDir),
// string/math (Left, Right, Mid, Len, Chr, Asc, Str, Val, Rnd, Int), JSON, HTTP, HELP.
package std

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"sync"

	"cyberbasic/compiler/vm"
	"github.com/google/uuid"
)

func toString(v interface{}) string {
	if v == nil {
		return ""
	}
	return fmt.Sprint(v)
}

func isTruthy(v interface{}) bool {
	if v == nil {
		return false
	}
	switch x := v.(type) {
	case bool:
		return x
	case int:
		return x != 0
	case float64:
		return x != 0
	case string:
		return x != ""
	default:
		return true
	}
}

var (
	jsonStore   = make(map[string]interface{})
	jsonCounter int
	jsonMu      sync.Mutex

	lastDirEntries []string
	dirMu          sync.Mutex

	// TimerStart(name) / TimerElapsed(name)
	timerStarts = make(map[string]time.Time)
	timerMu     sync.Mutex

	// enumRegistry is set by RegisterEnums(chunk.Enums) before running; used by Enum.getValue/getName/hasValue
	enumRegistry map[string]map[string]int64
	enumRegMu    sync.RWMutex

	logLevel int // SetLogLevel(level): 0=off, 1=error, 2=warn, 3=info, 4=debug

	// Save/Load
	autosaveIntervalSec float64
	autosavePath        string
	autosaveMu          sync.Mutex

	// Localization: code -> key -> translated string
	locLanguages     = make(map[string]map[string]string)
	locCurrentCode   string
	locMu            sync.RWMutex
)

func toFloat64(v interface{}) float64 {
	if v == nil {
		return 0
	}
	switch x := v.(type) {
	case int:
		return float64(x)
	case float64:
		return x
	case string:
		f, _ := strconv.ParseFloat(x, 64)
		return f
	default:
		return 0
	}
}

func toInt(v interface{}) int {
	if v == nil {
		return 0
	}
	switch x := v.(type) {
	case int:
		return x
	case float64:
		return int(math.Trunc(x))
	case string:
		n, _ := strconv.Atoi(x)
		return n
	default:
		return 0
	}
}

// RegisterEnums stores the chunk's enum map so Enum.getValue/getName/hasValue can look up at runtime.
// Call this after Compile and before Run (e.g. right after LoadChunk in main).
func RegisterEnums(enums map[string]vm.EnumMembers) {
	enumRegMu.Lock()
	defer enumRegMu.Unlock()
	enumRegistry = make(map[string]map[string]int64)
	for name, members := range enums {
		m := make(map[string]int64)
		for k, v := range members {
			m[k] = v
		}
		enumRegistry[name] = m
	}
}

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
	// LoadText(path): read entire file as string (alias for ReadFile)
	v.RegisterForeign("LoadText", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("LoadText(path) requires 1 argument")
		}
		data, err := os.ReadFile(toString(args[0]))
		if err != nil {
			return nil, err
		}
		return string(data), nil
	})
	// SaveText(path, text): write string to file (alias for WriteFile)
	v.RegisterForeign("SaveText", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("SaveText(path, text) requires 2 arguments")
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
	v.RegisterForeign("CopyFile", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("CopyFile(src, dst) requires 2 arguments")
		}
		data, err := os.ReadFile(toString(args[0]))
		if err != nil {
			return nil, err
		}
		err = os.WriteFile(toString(args[1]), data, 0644)
		return err == nil, err
	})
	v.RegisterForeign("ListDir", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("ListDir(path) requires 1 argument")
		}
		path := toString(args[0])
		entries, err := os.ReadDir(path)
		if err != nil {
			return nil, err
		}
		dirMu.Lock()
		lastDirEntries = make([]string, 0, len(entries))
		for _, e := range entries {
			lastDirEntries = append(lastDirEntries, e.Name())
		}
		dirMu.Unlock()
		return len(lastDirEntries), nil
	})
	v.RegisterForeign("DirectoryList", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("DirectoryList(path) requires 1 argument")
		}
		path := toString(args[0])
		entries, err := os.ReadDir(path)
		if err != nil {
			return nil, err
		}
		dirMu.Lock()
		lastDirEntries = make([]string, 0, len(entries))
		for _, e := range entries {
			lastDirEntries = append(lastDirEntries, e.Name())
		}
		dirMu.Unlock()
		return len(lastDirEntries), nil
	})
	v.RegisterForeign("GetDirItem", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return "", nil
		}
		idx := toInt(args[0])
		dirMu.Lock()
		defer dirMu.Unlock()
		if idx < 0 || idx >= len(lastDirEntries) {
			return "", nil
		}
		return lastDirEntries[idx], nil
	})
	v.RegisterForeign("ExecuteFile", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("ExecuteFile(path) requires 1 argument")
		}
		path := toString(args[0])
		path, err := filepath.Abs(path)
		if err != nil {
			return nil, err
		}
		cmd := exec.Command(path)
		cmd.Stdin = nil
		cmd.Stdout = nil
		cmd.Stderr = nil
		err = cmd.Start()
		return err == nil, err
	})

	// --- Multi-window (env and SpawnWindow) ---
	v.RegisterForeign("GetEnv", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return "", nil
		}
		return os.Getenv(toString(args[0])), nil
	})
	v.RegisterForeign("IsWindowProcess", func(args []interface{}) (interface{}, error) {
		return os.Getenv("CYBERBASIC_WINDOW") == "1", nil
	})
	v.RegisterForeign("GetWindowTitle", func(args []interface{}) (interface{}, error) {
		s := os.Getenv("CYBERBASIC_WINDOW_TITLE")
		if s == "" {
			return "Window", nil
		}
		return s, nil
	})
	v.RegisterForeign("GetWindowWidth", func(args []interface{}) (interface{}, error) {
		s := os.Getenv("CYBERBASIC_WINDOW_WIDTH")
		if s == "" {
			return 400, nil
		}
		n, _ := strconv.Atoi(s)
		if n <= 0 {
			return 400, nil
		}
		return n, nil
	})
	v.RegisterForeign("GetWindowHeight", func(args []interface{}) (interface{}, error) {
		s := os.Getenv("CYBERBASIC_WINDOW_HEIGHT")
		if s == "" {
			return 300, nil
		}
		n, _ := strconv.Atoi(s)
		if n <= 0 {
			return 300, nil
		}
		return n, nil
	})
	v.RegisterForeign("SpawnWindow", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return 0, fmt.Errorf("SpawnWindow(port, title, width, height) requires 4 arguments")
		}
		scriptPath := os.Getenv("CYBERBASIC_SCRIPT")
		if scriptPath == "" {
			return 0, nil
		}
		exe, err := os.Executable()
		if err != nil {
			return 0, nil
		}
		port := toInt(args[0])
		title := toString(args[1])
		width := toInt(args[2])
		height := toInt(args[3])
		cmd := exec.Command(exe, scriptPath, "--window",
			"--parent=127.0.0.1:"+strconv.Itoa(port),
			"--title="+title,
			"--width="+strconv.Itoa(width),
			"--height="+strconv.Itoa(height))
		cmd.Env = append(os.Environ(), "CYBERBASIC_SCRIPT="+scriptPath)
		cmd.Stdin = nil
		cmd.Stdout = nil
		cmd.Stderr = nil
		if err := cmd.Start(); err != nil {
			return 0, nil
		}
		return 1, nil
	})

	// --- String (DBP-style: Left, Right, Mid, Len, Chr, Asc, Str, Val) ---
	v.RegisterForeign("Left", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return "", nil
		}
		s := toString(args[0])
		n := toInt(args[1])
		if n <= 0 {
			return "", nil
		}
		r := []rune(s)
		if n >= len(r) {
			return s, nil
		}
		return string(r[:n]), nil
	})
	v.RegisterForeign("Right", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return "", nil
		}
		s := toString(args[0])
		n := toInt(args[1])
		if n <= 0 {
			return "", nil
		}
		r := []rune(s)
		if n >= len(r) {
			return s, nil
		}
		return string(r[len(r)-n:]), nil
	})
	v.RegisterForeign("Mid", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return "", nil
		}
		s := toString(args[0])
		start1 := toInt(args[1])
		r := []rune(s)
		if start1 < 1 || len(r) == 0 {
			return "", nil
		}
		start0 := start1 - 1
		if start0 >= len(r) {
			return "", nil
		}
		if len(args) < 3 {
			return string(r[start0:]), nil
		}
		count := toInt(args[2])
		if count <= 0 {
			return "", nil
		}
		end := start0 + count
		if end > len(r) {
			end = len(r)
		}
		return string(r[start0:end]), nil
	})
	// SUBSTR(s, start, count): 0-based start; alias for substring. Same as Mid(s, start+1, count).
	v.RegisterForeign("Substr", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return "", nil
		}
		s := toString(args[0])
		start0 := toInt(args[1])
		r := []rune(s)
		if start0 < 0 || start0 >= len(r) {
			return "", nil
		}
		if len(args) < 3 {
			return string(r[start0:]), nil
		}
		count := toInt(args[2])
		if count <= 0 {
			return "", nil
		}
		end := start0 + count
		if end > len(r) {
			end = len(r)
		}
		return string(r[start0:end]), nil
	})
	v.RegisterForeign("Instr", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return 0, nil
		}
		s := toString(args[0])
		sub := toString(args[1])
		idx := strings.Index(s, sub)
		if idx < 0 {
			return 0, nil
		}
		return idx + 1, nil
	})
	v.RegisterForeign("Upper", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return "", nil
		}
		return strings.ToUpper(toString(args[0])), nil
	})
	v.RegisterForeign("Lower", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return "", nil
		}
		return strings.ToLower(toString(args[0])), nil
	})
	// Random(n): integer 0 to n-1. Random(min, max): integer in [min, max] inclusive.
	v.RegisterForeign("Random", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return rand.Int(), nil
		}
		if len(args) >= 2 {
			lo, hi := toInt(args[0]), toInt(args[1])
			if hi < lo {
				lo, hi = hi, lo
			}
			n := hi - lo + 1
			if n <= 0 {
				return lo, nil
			}
			return rand.Intn(n) + lo, nil
		}
		n := toInt(args[0])
		if n <= 0 {
			return 0, nil
		}
		return rand.Intn(n), nil
	})
	v.RegisterForeign("TimeNow", func(args []interface{}) (interface{}, error) {
		return float64(time.Now().UnixNano()) / 1e9, nil
	})
	v.RegisterForeign("TimerStart", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("TimerStart(name) requires 1 argument")
		}
		name := toString(args[0])
		timerMu.Lock()
		timerStarts[name] = time.Now()
		timerMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("TimerElapsed", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return 0.0, nil
		}
		name := toString(args[0])
		timerMu.Lock()
		start, ok := timerStarts[name]
		timerMu.Unlock()
		if !ok {
			return 0.0, nil
		}
		return time.Since(start).Seconds(), nil
	})
	// PrintDebug(value): print value to stderr for debugging (e.g. PrintDebug("x=" + Str(x)))
	v.RegisterForeign("PrintDebug", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			fmt.Fprintln(os.Stderr, "[debug]")
			return nil, nil
		}
		fmt.Fprintln(os.Stderr, "[debug]", fmt.Sprint(args[0]))
		return nil, nil
	})
	v.RegisterForeign("Len", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return 0, nil
		}
		return len([]rune(toString(args[0]))), nil
	})
	v.RegisterForeign("Chr", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return "", nil
		}
		code := toInt(args[0])
		return string(rune(code)), nil
	})
	v.RegisterForeign("Asc", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return 0, nil
		}
		s := toString(args[0])
		r := []rune(s)
		if len(r) == 0 {
			return 0, nil
		}
		return int(r[0]), nil
	})
	v.RegisterForeign("Str", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return "", nil
		}
		return fmt.Sprint(args[0]), nil
	})
	v.RegisterForeign("Val", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return 0.0, nil
		}
		s := strings.TrimSpace(toString(args[0]))
		f, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return 0.0, nil
		}
		return f, nil
	})
	// Assert(condition, message): if condition is falsy, return error so VM stops
	v.RegisterForeign("Assert", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, nil
		}
		if !isTruthy(args[0]) {
			msg := "assertion failed"
			if len(args) >= 2 {
				msg = toString(args[1])
			}
			return nil, fmt.Errorf("%s", msg)
		}
		return nil, nil
	})
	v.RegisterForeign("SetLogLevel", func(args []interface{}) (interface{}, error) {
		if len(args) >= 1 {
			logLevel = toInt(args[0])
		}
		return nil, nil
	})
	v.RegisterForeign("UUID", func(args []interface{}) (interface{}, error) {
		return uuid.New().String(), nil
	})

	// --- Math: Sin, Cos, Tan, Sqrt (radians for trig) ---
	v.RegisterForeign("Sin", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return 0.0, nil
		}
		return math.Sin(toFloat64(args[0])), nil
	})
	v.RegisterForeign("Cos", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return 1.0, nil
		}
		return math.Cos(toFloat64(args[0])), nil
	})
	v.RegisterForeign("Tan", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return 0.0, nil
		}
		return math.Tan(toFloat64(args[0])), nil
	})
	v.RegisterForeign("Sqrt", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return 0.0, nil
		}
		x := toFloat64(args[0])
		if x < 0 {
			return 0.0, nil
		}
		return math.Sqrt(x), nil
	})
	v.RegisterForeign("RandomFloat", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return 0.0, nil
		}
		lo, hi := toFloat64(args[0]), toFloat64(args[1])
		if hi < lo {
			lo, hi = hi, lo
		}
		return lo + rand.Float64()*(hi-lo), nil
	})
	v.RegisterForeign("RandomInt", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return 0, nil
		}
		lo, hi := toInt(args[0]), toInt(args[1])
		if hi < lo {
			lo, hi = hi, lo
		}
		n := hi - lo + 1
		if n <= 0 {
			return lo, nil
		}
		return rand.Intn(n) + lo, nil
	})
	v.RegisterForeign("Log", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			fmt.Fprintln(os.Stderr, "[log]")
			return nil, nil
		}
		fmt.Fprintln(os.Stderr, "[log]", fmt.Sprint(args[0]))
		return nil, nil
	})

	// --- Math (DBP-style: Rnd, Int) ---
	v.RegisterForeign("Rnd", func(args []interface{}) (interface{}, error) {
		if len(args) == 0 {
			return rand.Float64(), nil
		}
		n := toInt(args[0])
		if n <= 0 {
			return 1, nil
		}
		return rand.Intn(n) + 1, nil
	})
	v.RegisterForeign("Int", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return 0, nil
		}
		return int(math.Trunc(toFloat64(args[0]))), nil
	})
	// Radians(degrees): convert degrees to radians
	v.RegisterForeign("Radians", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return 0.0, nil
		}
		return toFloat64(args[0]) * math.Pi / 180, nil
	})
	// Degrees(radians): convert radians to degrees
	v.RegisterForeign("Degrees", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return 0.0, nil
		}
		return toFloat64(args[0]) * 180 / math.Pi, nil
	})
	// AngleWrap(angle): wrap angle in radians to [-PI, PI]
	v.RegisterForeign("AngleWrap", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return 0.0, nil
		}
		a := toFloat64(args[0])
		const twoPi = 2 * math.Pi
		a = math.Mod(a+math.Pi, twoPi)
		if a < 0 {
			a += twoPi
		}
		return a - math.Pi, nil
	})
	// WrapAngle(angle): alias for AngleWrap — wrap angle in radians to [-PI, PI]
	v.RegisterForeign("WrapAngle", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return 0.0, nil
		}
		a := toFloat64(args[0])
		const twoPi = 2 * math.Pi
		a = math.Mod(a+math.Pi, twoPi)
		if a < 0 {
			a += twoPi
		}
		return a - math.Pi, nil
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
			return nil, fmt.Errorf("GetJSONKey(handleOrDict, key) requires 2 arguments")
		}
		key := toString(args[1])
		// If first arg is a map (e.g. from dict literal or CreateDict), use it directly
		if m, ok := args[0].(map[string]interface{}); ok {
			val, _ := m[key]
			return val, nil
		}
		id := toString(args[0])
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
	// CreateDict returns a new empty map for use as a dictionary; use SetDictKey to add pairs.
	v.RegisterForeign("CreateDict", func(args []interface{}) (interface{}, error) {
		return make(map[string]interface{}), nil
	})
	// SetDictKey(dict, key, value) sets dict[key]=value and returns dict (for chaining or assignment).
	v.RegisterForeign("SetDictKey", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("SetDictKey(dict, key, value) requires 3 arguments")
		}
		m, ok := args[0].(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("SetDictKey: first argument must be a dictionary")
		}
		m[toString(args[1])] = args[2]
		return m, nil
	})
	// Dictionary.* helpers
	v.RegisterForeign("Dictionary.has", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("Dictionary.has(dict, key) requires 2 arguments")
		}
		m, ok := args[0].(map[string]interface{})
		if !ok {
			return false, nil
		}
		_, ok = m[toString(args[1])]
		return ok, nil
	})
	v.RegisterForeign("Dictionary.keys", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("Dictionary.keys(dict) requires 1 argument")
		}
		m, ok := args[0].(map[string]interface{})
		if !ok {
			return []interface{}{}, nil
		}
		keys := make([]interface{}, 0, len(m))
		for k := range m {
			keys = append(keys, k)
		}
		return keys, nil
	})
	v.RegisterForeign("Dictionary.values", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("Dictionary.values(dict) requires 1 argument")
		}
		m, ok := args[0].(map[string]interface{})
		if !ok {
			return []interface{}{}, nil
		}
		vals := make([]interface{}, 0, len(m))
		for _, v := range m {
			vals = append(vals, v)
		}
		return vals, nil
	})
	v.RegisterForeign("Dictionary.size", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("Dictionary.size(dict) requires 1 argument")
		}
		m, ok := args[0].(map[string]interface{})
		if !ok {
			return 0, nil
		}
		return len(m), nil
	})
	v.RegisterForeign("Dictionary.remove", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("Dictionary.remove(dict, key) requires 2 arguments")
		}
		m, ok := args[0].(map[string]interface{})
		if !ok {
			return args[0], nil
		}
		delete(m, toString(args[1]))
		return m, nil
	})
	v.RegisterForeign("Dictionary.clear", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("Dictionary.clear(dict) requires 1 argument")
		}
		m, ok := args[0].(map[string]interface{})
		if !ok {
			return args[0], nil
		}
		for k := range m {
			delete(m, k)
		}
		return m, nil
	})
	v.RegisterForeign("Dictionary.merge", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("Dictionary.merge(dict1, dict2) requires 2 arguments")
		}
		m1, ok := args[0].(map[string]interface{})
		if !ok {
			return args[0], nil
		}
		m2, ok := args[1].(map[string]interface{})
		if !ok {
			return args[0], nil
		}
		out := make(map[string]interface{})
		for k, v := range m1 {
			out[k] = v
		}
		for k, v := range m2 {
			out[k] = v
		}
		return out, nil
	})
	v.RegisterForeign("Dictionary.get", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("Dictionary.get(dict, key [, default]) requires 2 or 3 arguments")
		}
		m, ok := args[0].(map[string]interface{})
		if !ok {
			if len(args) >= 3 {
				return args[2], nil
			}
			return nil, nil
		}
		key := toString(args[1])
		if v, ok := m[key]; ok {
			return v, nil
		}
		if len(args) >= 3 {
			return args[2], nil
		}
		return nil, nil
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

	// --- Save / Load system ---
	v.RegisterForeign("SaveGame", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("SaveGame(path, data) requires 2 arguments")
		}
		path := toString(args[0])
		var data []byte
		var err error
		switch d := args[1].(type) {
		case string:
			data = []byte(d)
		case map[string]interface{}:
			data, err = json.MarshalIndent(d, "", "  ")
			if err != nil {
				return nil, err
			}
		default:
			// treat as JSON handle if it's a known id
			id := toString(args[1])
			jsonMu.Lock()
			obj, ok := jsonStore[id]
			jsonMu.Unlock()
			if ok {
				data, err = json.MarshalIndent(obj, "", "  ")
				if err != nil {
					return nil, err
				}
			} else {
				data, err = json.Marshal(args[1])
				if err != nil {
					return nil, err
				}
			}
		}
		if err := os.WriteFile(path, data, 0644); err != nil {
			return nil, err
		}
		return true, nil
	})
	v.RegisterForeign("LoadGame", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("LoadGame(path) requires 1 argument")
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
	v.RegisterForeign("Autosave", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("Autosave(intervalSeconds) requires 1 argument")
		}
		autosaveMu.Lock()
		autosaveIntervalSec = toFloat64(args[0])
		if len(args) >= 2 {
			autosavePath = toString(args[1])
		}
		autosaveMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("SaveExists", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return false, nil
		}
		path := toString(args[0])
		_, err := os.Stat(path)
		return err == nil, nil
	})

	// --- Localization ---
	v.RegisterForeign("LoadLanguage", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("LoadLanguage(path) requires 1 argument")
		}
		path := toString(args[0])
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}
		var raw map[string]interface{}
		if err := json.Unmarshal(data, &raw); err != nil {
			return nil, err
		}
		code := strings.TrimSuffix(filepath.Base(path), ".json")
		if code == filepath.Base(path) {
			code = "default"
		}
		locMu.Lock()
		if locLanguages[code] == nil {
			locLanguages[code] = make(map[string]string)
		}
		for k, v := range raw {
			if s, ok := v.(string); ok {
				locLanguages[code][k] = s
			} else {
				locLanguages[code][k] = fmt.Sprint(v)
			}
		}
		locMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("SetLanguage", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, nil
		}
		locMu.Lock()
		locCurrentCode = toString(args[0])
		locMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("Translate", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return "", nil
		}
		key := toString(args[0])
		locMu.RLock()
		code := locCurrentCode
		if code == "" {
			code = "default"
		}
		lang := locLanguages[code]
		var out string
		if lang != nil {
			out = lang[key]
		}
		locMu.RUnlock()
		if out == "" {
			return key, nil
		}
		return out, nil
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

	// --- Enum runtime API (requires RegisterEnums(chunk.Enums) before Run) ---
	v.RegisterForeign("Enum.getValue", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("Enum.getValue(enumName, valueName) requires 2 arguments")
		}
		enumRegMu.RLock()
		defer enumRegMu.RUnlock()
		if enumRegistry == nil {
			return nil, fmt.Errorf("enum registry not set (call RegisterEnums before Run)")
		}
		enumName := strings.ToLower(toString(args[0]))
		valueName := strings.ToLower(toString(args[1]))
		m, ok := enumRegistry[enumName]
		if !ok {
			return nil, fmt.Errorf("unknown enum: %s", enumName)
		}
		val, ok := m[valueName]
		if !ok {
			return nil, fmt.Errorf("enum %s has no value %s", enumName, valueName)
		}
		return int(val), nil
	})
	v.RegisterForeign("Enum.getName", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("Enum.getName(enumName, value) requires 2 arguments")
		}
		enumRegMu.RLock()
		defer enumRegMu.RUnlock()
		if enumRegistry == nil {
			return nil, fmt.Errorf("enum registry not set (call RegisterEnums before Run)")
		}
		enumName := strings.ToLower(toString(args[0]))
		var target int64
		switch x := args[1].(type) {
		case int:
			target = int64(x)
		case float64:
			target = int64(x)
		default:
			return nil, fmt.Errorf("Enum.getName value must be numeric")
		}
		m, ok := enumRegistry[enumName]
		if !ok {
			return nil, fmt.Errorf("unknown enum: %s", enumName)
		}
		for name, v := range m {
			if v == target {
				return name, nil
			}
		}
		return "", nil
	})
	v.RegisterForeign("Enum.hasValue", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("Enum.hasValue(enumName, valueName) requires 2 arguments")
		}
		enumRegMu.RLock()
		defer enumRegMu.RUnlock()
		if enumRegistry == nil {
			return false, nil
		}
		enumName := strings.ToLower(toString(args[0]))
		valueName := strings.ToLower(toString(args[1]))
		m, ok := enumRegistry[enumName]
		if !ok {
			return false, nil
		}
		_, ok = m[valueName]
		return ok, nil
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
		if len(args) == 1 {
			if s, ok := args[0].(string); ok {
				msg := helpCommandLine(strings.ToLower(strings.TrimSpace(s)))
				if msg != "" {
					fmt.Println(msg)
					return nil, nil
				}
			}
		}
		fmt.Println("CyberBasic: Full command list → docs/COMMAND_REFERENCE.md and API_REFERENCE.md in the project root.")
		fmt.Println("From the CLI run: cyberbasic --list-commands for a short list, or cyberbasic --help for options.")
		fmt.Println("For one command: Help(\"DrawRectangle\") or Help(\"rect\")")
		return nil, nil
	})
	v.RegisterForeign("?", func(args []interface{}) (interface{}, error) {
		if len(args) == 1 {
			if s, ok := args[0].(string); ok {
				msg := helpCommandLine(strings.ToLower(strings.TrimSpace(s)))
				if msg != "" {
					fmt.Println(msg)
					return nil, nil
				}
			}
		}
		fmt.Println("CyberBasic: Full command list → docs/COMMAND_REFERENCE.md and API_REFERENCE.md.")
		fmt.Println("CLI: cyberbasic --list-commands")
		return nil, nil
	})
}

// helpCommandLine returns a one-line help string for a command name, or "" if unknown.
func helpCommandLine(cmd string) string {
	helpMap := map[string]string{
		"initwindow":       "InitWindow(width, height, title) – open game window",
		"closewindow":      "CloseWindow() – close window and exit",
		"settargetfps":     "SetTargetFPS(fps) – target frame rate",
		"getframetime":     "GetFrameTime() – delta time since last frame (seconds)",
		"windowshouldclose": "WindowShouldClose() – true when user requested close",
		"drawrectangle":   "DrawRectangle(x, y, w, h, r, g, b, a) – filled rectangle. Alias: rect(...)",
		"rect":             "rect(x, y, w, h, color...) – alias of DrawRectangle",
		"drawcircle":       "DrawCircle(x, y, radius, r, g, b, a) – filled circle. Alias: circle(...)",
		"circle":           "circle(x, y, radius, color...) – alias of DrawCircle",
		"drawtext":         "DrawText(text, x, y, size, r, g, b, a) – draw text",
		"drawtexture":      "DrawTexture(id, x, y [, tint]) – draw texture. Alias: sprite(...)",
		"sprite":           "sprite(id, x, y [, tint]) – alias of DrawTexture",
		"clearbackground":  "ClearBackground(r, g, b, a) – clear screen",
		"drawcube":         "DrawCube(x, y, z, w, h, d, color) – 3D cube. Alias: cube(...)",
		"cube":             "cube(x, y, z, w, h, d, color) – alias of DrawCube",
		"guibutton":        "GuiButton(x, y, w, h, text) – button; returns 1 if clicked. Alias: button(...)",
		"button":           "button(x, y, w, h, text) – alias of GuiButton",
		"keydown":          "KeyDown(key) – true while key held (e.g. KEY_W, KEY_ESCAPE)",
		"keypressed":      "KeyPressed(key) – true once when key pressed",
		"getmousex":        "GetMouseX() – mouse X",
		"getmousey":        "GetMouseY() – mouse Y",
		"createworld2d":    "CreateWorld2D() – create 2D physics world; returns world id",
		"createbox2d":      "CreateBox2D(worldId, bodyId, x, y, w, h, density) – box body",
		"stepallphysics2d": "StepAllPhysics2D(dt) – step all 2D worlds (called automatically in hybrid loop)",
		"createworld3d":    "CreateWorld3D() – create 3D physics world",
		"createbox3d":      "CreateBox3D(worldId, bodyId, x, y, z, w, h, d, mass) – 3D box body",
		"stepallphysics3d": "StepAllPhysics3D(dt) – step all 3D worlds (called automatically in hybrid loop)",
		"print":            "Print(value) – print to console",
		"str":              "Str(value) – convert to string",
		"int":              "Int(value) – convert to integer",
	}
	return helpMap[cmd]
}
