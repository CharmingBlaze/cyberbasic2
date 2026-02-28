package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"cyberbasic/compiler"
	"cyberbasic/compiler/bindings/box2d"
	"cyberbasic/compiler/bindings/bullet"
	"cyberbasic/compiler/bindings/ecs"
	"cyberbasic/compiler/bindings/game"
	"cyberbasic/compiler/bindings/indoor"
	"cyberbasic/compiler/bindings/navigation"
	"cyberbasic/compiler/bindings/net"
	"cyberbasic/compiler/bindings/raylib"
	"cyberbasic/compiler/bindings/scene"
	"cyberbasic/compiler/bindings/sql"
	"cyberbasic/compiler/bindings/std"
	"cyberbasic/compiler/bindings/objects"
	"cyberbasic/compiler/bindings/procedural"
	"cyberbasic/compiler/bindings/terrain"
	"cyberbasic/compiler/bindings/vegetation"
	"cyberbasic/compiler/bindings/water"
	"cyberbasic/compiler/bindings/world"
	"cyberbasic/compiler/gogen"
	"cyberbasic/compiler/lexer"
	"cyberbasic/compiler/parser"
	"cyberbasic/compiler/runtime"
	"cyberbasic/compiler/vm"
)

func main() {
	fmt.Println("CyberBasic starting...")

	// Check for --help first
	for _, arg := range os.Args {
		if arg == "--help" {
			fmt.Println("CyberBasic - A BASIC-like language with Raylib + Bullet physics")
			fmt.Println("Usage: cyberbasic <filename.bas> [options]")
			fmt.Println("Options:")
			fmt.Println("  --compile-only    Compile but don't run")
			fmt.Println("  --gen-go [file]   Generate Go source that calls raylib directly (default: <basename>_gen.go)")
			fmt.Println("  --debug           Enable debug output")
			fmt.Println("  --list-commands   Print built-in command names (2D, 3D, GUI, Physics, Std)")
			fmt.Println("  --lint            Check program (compile only, no run); same as --compile-only")
			fmt.Println("  --help            Show this help")
			fmt.Println("  (Multi-window: --window --parent=host:port --title=... --width=... --height=...)")
			fmt.Println("Exit codes: 0 = success, 1 = compile/file error, 2 = runtime error")
			os.Exit(0)
		}
		if arg == "--list-commands" {
			printCommandList()
			os.Exit(0)
		}
	}

	// Filename is the first argument that does not start with -
	var filename string
	for i := 1; i < len(os.Args); i++ {
		if !strings.HasPrefix(os.Args[i], "-") {
			filename = os.Args[i]
			break
		}
		if os.Args[i] == "--gen-go" && i+1 < len(os.Args) {
			i++ // skip gen-go output path
		}
	}
	if filename == "" {
		// No file: default to 3D physics demo or show usage
		exeDir := filepath.Dir(os.Args[0])
		defaultBas := filepath.Join(exeDir, "examples", "run_3d_physics_demo.bas")
		if _, err := os.Stat(defaultBas); err == nil {
			filename = defaultBas
			fmt.Println("No file specified, running default: examples/run_3d_physics_demo.bas")
		} else {
			defaultBas = filepath.Join("examples", "run_3d_physics_demo.bas")
			if _, err := os.Stat(defaultBas); err == nil {
				filename = defaultBas
				fmt.Println("No file specified, running default: examples/run_3d_physics_demo.bas")
			} else {
				fmt.Println("CyberBasic - A BASIC-like language with Raylib + Bullet physics")
				fmt.Println("Usage: cyberbasic <filename.bas> [options]  (or: cyberbasic [options] <filename.bas>)")
				fmt.Println("Options:")
				fmt.Println("  --compile-only    Compile but don't run")
				fmt.Println("  --gen-go [file]   Generate Go source that calls raylib directly")
				fmt.Println("  --debug           Enable debug output")
				os.Exit(1)
			}
		}
	}
	compileOnly := false
	debug := false
	genGo := false
	genGoOut := ""

	// Parse command line arguments (all args except exe name)
	for i := 1; i < len(os.Args); i++ {
		switch os.Args[i] {
		case "--compile-only":
			compileOnly = true
		case "--lint":
			compileOnly = true // lint = compile without running
		case "--debug":
			debug = true
		case "--gen-go":
			genGo = true
			if i+1 < len(os.Args) && len(os.Args[i+1]) > 0 && !strings.HasPrefix(os.Args[i+1], "-") {
				i++
				genGoOut = os.Args[i]
			}
		}
	}
	if genGo && genGoOut == "" {
		base := strings.TrimSuffix(filepath.Base(filename), filepath.Ext(filename))
		genGoOut = filepath.Join("generated", base+"_gen.go")
	}

	// Check if file exists
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		fmt.Printf("Error: File '%s' not found\n", filename)
		os.Exit(1)
	}

	// Set script path for SpawnWindow (child processes re-run same .bas)
	if absPath, err := filepath.Abs(filename); err == nil {
		os.Setenv("CYBERBASIC_SCRIPT", absPath)
	}
	// Parse window-mode args so the script can branch (IsWindowProcess) and get title/size
	for _, arg := range os.Args {
		if arg == "--window" {
			os.Setenv("CYBERBASIC_WINDOW", "1")
		} else if strings.HasPrefix(arg, "--parent=") {
			os.Setenv("CYBERBASIC_PARENT", strings.TrimPrefix(arg, "--parent="))
		} else if strings.HasPrefix(arg, "--title=") {
			os.Setenv("CYBERBASIC_WINDOW_TITLE", strings.TrimPrefix(arg, "--title="))
		} else if strings.HasPrefix(arg, "--width=") {
			os.Setenv("CYBERBASIC_WINDOW_WIDTH", strings.TrimPrefix(arg, "--width="))
		} else if strings.HasPrefix(arg, "--height=") {
			os.Setenv("CYBERBASIC_WINDOW_HEIGHT", strings.TrimPrefix(arg, "--height="))
		}
	}

	// Read source file
	source, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		os.Exit(1)
	}

	// Preprocess: expand #include "file.bas"
	baseDir := filepath.Dir(filename)
	source = preprocessIncludes(source, baseDir, nil)

	// --gen-go: parse and generate Go source (no bytecode)
	if genGo {
		l := lexer.New(string(source))
		tokens, err := l.Tokenize()
		if err != nil {
			fmt.Printf("Lex error: %v\n", err)
			os.Exit(1)
		}
		p := parser.New(tokens)
		program, err := p.Parse()
		if err != nil {
			fmt.Printf("Parse error: %v\n", err)
			os.Exit(1)
		}
		goCode, err := gogen.Generate(program)
		if err != nil {
			fmt.Printf("Go gen error: %v\n", err)
			os.Exit(1)
		}
		if dir := filepath.Dir(genGoOut); dir != "." {
			_ = os.MkdirAll(dir, 0755)
		}
		if err := os.WriteFile(genGoOut, []byte(goCode), 0644); err != nil {
			fmt.Printf("Write error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Generated %s\n", genGoOut)
		os.Exit(0)
	}

	fmt.Printf("Compiling %s...\n", filename)

	// Create compiler
	comp := compiler.New()

	// Compile source code
	chunk, err := comp.Compile(string(source))
	if err != nil {
		fmt.Printf("Compilation error: %v\n", err)
		os.Exit(1)
	}

	if debug {
		fmt.Printf("Compiled %d bytes of bytecode with %d constants\n", len(chunk.Code), len(chunk.Constants))
	}

	if compileOnly {
		fmt.Println("Compilation successful!")
		os.Exit(0)
	}

	// Create runtime
	rt := runtime.NewRuntime()

	// Load bytecode into VM and wire runtime so game opcodes call into it
	rt.GetVM().LoadChunk(chunk)
	std.RegisterEnums(chunk.Enums)
	rt.GetVM().SetRuntime(rt)
	// Expose raylib and Bullet as foreign API: RL.*, BULLET.*
	raylib.RegisterRaylib(rt.GetVM())
	bullet.RegisterBullet(rt.GetVM())
	box2d.RegisterBox2D(rt.GetVM())
	ecs.RegisterECS(rt.GetVM())
	net.RegisterNet(rt.GetVM())
	scene.RegisterScene(rt.GetVM())
	game.RegisterGame(rt.GetVM())
	sql.RegisterSQL(rt.GetVM())
	terrain.RegisterTerrain(rt.GetVM())
	objects.RegisterObjects(rt.GetVM())
	procedural.RegisterProcedural(rt.GetVM())
	water.RegisterWater(rt.GetVM())
	vegetation.RegisterVegetation(rt.GetVM())
	world.RegisterWorld(rt.GetVM())
	navigation.RegisterNavigation(rt.GetVM())
	indoor.RegisterIndoor(rt.GetVM())
	std.RegisterStd(rt.GetVM())

	fmt.Println("Running program...")

	// Run the program
	err = rt.GetVM().Run()
	if err != nil {
		fmt.Printf("Runtime error: %v\n", err)
		if debug {
			for i, f := range rt.GetVM().StackTrace() {
				fmt.Printf("  #%d line %d (ip %d)\n", i, f.Line, f.IP)
			}
		}
		rt.CloseWindow()
		os.Exit(2)
	}

	rt.CloseWindow()
	fmt.Println("Program completed successfully!")
	os.Exit(0)
}

// printCommandList prints built-in command names grouped by category (for --list-commands).
func printCommandList() {
	fmt.Println("Built-in commands (see docs/COMMAND_REFERENCE.md and API_REFERENCE.md for full reference)")
	fmt.Println()
	fmt.Println("Window & system: InitWindow, CloseWindow, SetTargetFPS, GetFrameTime, WindowShouldClose, DisableCursor, EnableCursor")
	fmt.Println("Game loop (hybrid): ClearRenderQueues, FlushRenderQueues, StepAllPhysics2D, StepAllPhysics3D")
	fmt.Println("2D (shapes/textures): DrawRectangle, rect, DrawCircle, circle, DrawLine, DrawText, DrawTexture, sprite, ClearBackground, Background")
	fmt.Println("3D: DrawCube, cube, DrawSphere, DrawModel, DrawModelSimple, DrawGrid, BeginMode3D, EndMode3D")
	fmt.Println("GUI: GuiButton, button, GuiLabel, GuiSlider, GuiCheckbox, GuiTextbox, GuiProgressBar")
	fmt.Println("Physics 2D: CreateWorld2D, CreateBox2D, CreateCircle2D, Step2D, StepAllPhysics2D, GetPositionX2D, GetPositionY2D, ApplyForce2D")
	fmt.Println("Physics 3D: CreateWorld3D, CreateBox3D, CreateSphere3D, Step3D, StepAllPhysics3D, GetPositionX3D, GetPositionY3D, GetPositionZ3D, ApplyForce3D")
	fmt.Println("Input: KeyDown, KeyPressed, GetMouseX, GetMouseY, GetMouseDeltaX, GetMouseDeltaY")
	fmt.Println("Math: Clamp, Lerp, Vec2, Vec3, Color")
	fmt.Println("Std: Print, Str, Int, Rnd, OpenFile, ReadLine, WriteLine, CloseFile, Left, Right, Mid, Len")
}

// Additional utility functions for debugging and development
func printTokens(tokens []lexer.Token) {
	fmt.Println("Tokens:")
	for _, token := range tokens {
		fmt.Printf("  %s: %s (line %d, col %d)\n", token.Type, token.Value, token.Line, token.Col)
	}
}

func printAST(program *parser.Program) {
	fmt.Println("Abstract Syntax Tree:")
	fmt.Println(program.String())
}

func printBytecode(chunk *vm.Chunk) {
	fmt.Println("Bytecode:")
	for i, instruction := range chunk.Code {
		fmt.Printf("  %04d: %d", i, instruction)
		if i < len(chunk.Code)-1 {
			fmt.Printf(" %d", chunk.Code[i+1])
			i++
		}
		fmt.Println()
	}

	fmt.Println("\nConstants:")
	for i, constant := range chunk.Constants {
		fmt.Printf("  %d: %v\n", i, constant)
	}
}

func runTests() {
	fmt.Println("Running CyberBasic tests...")

	// Test basic arithmetic
	testArithmetic()

	// Test control flow
	testControlFlow()

	// Test functions
	testFunctions()

	fmt.Println("All tests completed!")
}

func testArithmetic() {
	fmt.Println("Testing arithmetic...")

	source := `
DIM a AS INTEGER
DIM b AS INTEGER
a = 10
b = 20
DIM result AS INTEGER
result = a + b
`

	comp := compiler.New()
	chunk, err := comp.Compile(source)
	if err != nil {
		fmt.Printf("Arithmetic test failed: %v\n", err)
		return
	}

	v := vm.NewVM()
	v.LoadChunk(chunk)
	err = v.Run()
	if err != nil {
		fmt.Printf("Arithmetic test runtime error: %v\n", err)
		return
	}

	fmt.Println("Arithmetic test passed!")
}

func testControlFlow() {
	fmt.Println("Testing control flow...")

	source := `
DIM i AS INTEGER
DIM sum AS INTEGER
sum = 0
FOR i = 1 TO 10
    sum = sum + i
NEXT
`

	comp := compiler.New()
	chunk, err := comp.Compile(source)
	if err != nil {
		fmt.Printf("Control flow test failed: %v\n", err)
		return
	}

	v := vm.NewVM()
	v.LoadChunk(chunk)
	err = v.Run()
	if err != nil {
		fmt.Printf("Control flow test runtime error: %v\n", err)
		return
	}

	fmt.Println("Control flow test passed!")
}

func testFunctions() {
	fmt.Println("Testing functions...")

	source := `
FUNCTION AddNumbers(a AS INTEGER, b AS INTEGER) AS INTEGER
    RETURN a + b
END FUNCTION

DIM result AS INTEGER
result = AddNumbers(5, 7)
`

	comp := compiler.New()
	chunk, err := comp.Compile(source)
	if err != nil {
		fmt.Printf("Function test failed: %v\n", err)
		return
	}

	v := vm.NewVM()
	v.LoadChunk(chunk)
	err = v.Run()
	if err != nil {
		fmt.Printf("Function test runtime error: %v\n", err)
		return
	}

	fmt.Println("Function test passed!")
}

// Helper function to get examples directory
func getExamplesDir() string {
	dir, err := os.Getwd()
	if err != nil {
		return "examples"
	}
	return filepath.Join(dir, "examples")
}

// Function to list available examples
func listExamples() {
	examplesDir := getExamplesDir()

	fmt.Println("Available examples:")

	files, err := os.ReadDir(examplesDir)
	if err != nil {
		fmt.Printf("Error reading examples directory: %v\n", err)
		return
	}

	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".bas" {
			fmt.Printf("  %s\n", file.Name())
		}
	}
}

// Function to run an example
func runExample(name string) {
	examplesDir := getExamplesDir()
	filename := filepath.Join(examplesDir, name)

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		fmt.Printf("Example '%s' not found\n", name)
		listExamples()
		return
	}

	// Read and run the example
	source, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading example: %v\n", err)
		return
	}

	fmt.Printf("Running example: %s\n", name)

	comp := compiler.New()
	chunk, err := comp.Compile(string(source))
	if err != nil {
		fmt.Printf("Example compilation error: %v\n", err)
		return
	}

	rt := runtime.NewRuntime()
	rt.GetVM().LoadChunk(chunk)

	err = rt.GetVM().Run()
	if err != nil {
		fmt.Printf("Example runtime error: %v\n", err)
		return
	}

	fmt.Printf("Example '%s' completed successfully!\n", name)
}

// preprocessIncludes expands #include "file.bas" and IMPORT "file.bas" with file contents (relative to baseDir). seen prevents cycles.
func preprocessIncludes(source []byte, baseDir string, seen map[string]bool) []byte {
	if seen == nil {
		seen = make(map[string]bool)
	}
	includeRe := regexp.MustCompile(`^\s*#include\s*"([^"]+)"\s*$`)
	importRe := regexp.MustCompile(`(?i)^\s*IMPORT\s*"([^"]+)"\s*$`)
	var out strings.Builder
	sc := bufio.NewScanner(strings.NewReader(string(source)))
	for sc.Scan() {
		line := sc.Text()
		var filePath string
		if m := includeRe.FindStringSubmatch(line); m != nil {
			filePath = m[1]
		} else if m := importRe.FindStringSubmatch(line); m != nil {
			filePath = m[1]
		}
		if filePath != "" {
			path := filepath.Join(baseDir, filePath)
			abs, _ := filepath.Abs(path)
			if seen[abs] {
				continue
			}
			seen[abs] = true
			inc, err := os.ReadFile(path)
			if err != nil {
				out.WriteString(line)
				out.WriteByte('\n')
				continue
			}
			incDir := filepath.Dir(path)
			inc = preprocessIncludes(inc, incDir, seen)
			out.Write(inc)
			if len(inc) > 0 && inc[len(inc)-1] != '\n' {
				out.WriteByte('\n')
			}
			continue
		}
		out.WriteString(line)
		out.WriteByte('\n')
	}
	return []byte(out.String())
}
