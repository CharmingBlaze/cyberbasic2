// Package app implements CLI orchestration: flags, compile, RegisterAll, run, REPL.
package app

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"cyberbasic/compiler"
	"cyberbasic/compiler/bindings"
	"cyberbasic/compiler/bindings/std"
	"cyberbasic/compiler/errors"
	"cyberbasic/compiler/gogen"
	"cyberbasic/compiler/runtime"
	"cyberbasic/compiler/vm"
)

// Main is the application entry (called from package main with build Version).
func Main(version string) {
	fmt.Println("CyberBasic starting...")

	// Check for --help and --version first
	for _, arg := range os.Args {
		if arg == "--version" {
			if version == "" {
				version = "dev"
			}
			fmt.Println("CyberBasic", version)
			os.Exit(0)
		}
		if arg == "--help" {
			printHelp()
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
	replMode := false
	for _, arg := range os.Args {
		if arg == "--repl" {
			replMode = true
			break
		}
	}
	if filename == "" && !replMode {
		fmt.Println("CyberBasic 2 — no script given, starting REPL (type EXIT to quit).")
		runREPL()
		os.Exit(0)
	}
	if replMode {
		runREPL()
		os.Exit(0)
	}

	compileOnly := false
	debug := false
	genGo := false
	genGoOut := ""
	var debuggerBreakpoints []int

	for i := 1; i < len(os.Args); i++ {
		switch os.Args[i] {
		case "--compile-only":
			compileOnly = true
		case "--lint":
			compileOnly = true
		case "--debug":
			debug = true
			_ = os.Setenv("CYBERBASIC_DEBUG", "1")
		case "--debugger":
			debug = true
		case "--break":
			arg := ""
			if strings.Contains(os.Args[i], "=") {
				parts := strings.SplitN(os.Args[i], "=", 2)
				arg = parts[1]
			} else if i+1 < len(os.Args) && !strings.HasPrefix(os.Args[i+1], "-") {
				i++
				arg = os.Args[i]
			}
			for _, s := range strings.Split(arg, ",") {
				var line int
				if _, err := fmt.Sscanf(strings.TrimSpace(s), "%d", &line); err == nil && line > 0 {
					debuggerBreakpoints = append(debuggerBreakpoints, line)
				}
			}
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

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		fmt.Printf("Error: File '%s' not found\n", filename)
		os.Exit(1)
	}

	if absPath, err := filepath.Abs(filename); err == nil {
		_ = os.Setenv("CYBERBASIC_SCRIPT", absPath)
	}
	for _, arg := range os.Args {
		if arg == "--window" {
			_ = os.Setenv("CYBERBASIC_WINDOW", "1")
		} else if strings.HasPrefix(arg, "--parent=") {
			_ = os.Setenv("CYBERBASIC_PARENT", strings.TrimPrefix(arg, "--parent="))
		} else if strings.HasPrefix(arg, "--title=") {
			_ = os.Setenv("CYBERBASIC_WINDOW_TITLE", strings.TrimPrefix(arg, "--title="))
		} else if strings.HasPrefix(arg, "--width=") {
			_ = os.Setenv("CYBERBASIC_WINDOW_WIDTH", strings.TrimPrefix(arg, "--width="))
		} else if strings.HasPrefix(arg, "--height=") {
			_ = os.Setenv("CYBERBASIC_WINDOW_HEIGHT", strings.TrimPrefix(arg, "--height="))
		}
	}

	source, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		os.Exit(1)
	}

	baseDir := filepath.Dir(filename)
	source = PreprocessIncludes(source, baseDir, nil)

	if genGo {
		runGenGo(string(source), genGoOut, filename)
		os.Exit(0)
	}

	fmt.Printf("Compiling %s...\n", filename)

	comp := compiler.New()
	comp.Filename = filename
	sourceStr := string(source)
	chunk, err := comp.Compile(sourceStr)
	if err != nil {
		errors.PrettyPrint(os.Stdout, sourceStr, filename, err)
		os.Exit(1)
	}

	if debug {
		fmt.Printf("Compiled %d bytes of bytecode with %d constants\n", len(chunk.Code), len(chunk.Constants))
	}

	if compileOnly {
		fmt.Println("Compilation successful!")
		os.Exit(0)
	}

	rt := runtime.NewRuntime()
	rt.GetVM().LoadChunk(chunk)
	stdRegisterEnumsAndRuntime(rt, chunk)
	if err := bindings.RegisterAll(rt.GetVM(), bindings.RegisterOptions{Source: sourceStr}); err != nil {
		fmt.Printf("Register bindings: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Running program...")

	v := rt.GetVM()
	if len(debuggerBreakpoints) > 0 {
		bpMap := make(map[int]bool)
		for _, l := range debuggerBreakpoints {
			bpMap[l] = true
		}
		v.SetBreakpoints(bpMap)
		v.SetDebugMode(true)
	}

	err = v.Run()
	if err != nil {
		if bp, ok := err.(*vm.ErrBreakpoint); ok {
			fmt.Printf("Breakpoint hit at line %d\n", bp.Line)
			for i, f := range v.StackTrace() {
				fmt.Printf("  #%d line %d (ip %d)\n", i, f.Line, f.IP)
			}
			rt.CloseWindow()
			os.Exit(0)
		}
		fmt.Printf("Runtime error: %v\n", err)
		if debug {
			for i, f := range v.StackTrace() {
				fmt.Printf("  #%d line %d (ip %d)\n", i, f.Line, f.IP)
			}
		}
		rt.CloseWindow()
		os.Exit(2)
	}

	if rt.HasImplicitHandlers() && runtime.DetectWindowMode(sourceStr) != runtime.ModeExplicit {
		err = rt.RunImplicitLoop()
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
	}

	rt.CloseWindow()
	fmt.Println("Program completed successfully!")
	os.Exit(0)
}

func stdRegisterEnumsAndRuntime(rt *runtime.Runtime, chunk *vm.Chunk) {
	std.RegisterEnums(chunk.Enums)
	rt.GetVM().SetRuntime(rt)
}

func runGenGo(source, genGoOut, basFilename string) {
	comp := compiler.New()
	comp.Filename = basFilename
	program, err := comp.Parse(source)
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
}

func printHelp() {
	fmt.Println("CyberBasic - A BASIC-like language with Raylib + Bullet physics")
	fmt.Println("Usage: cyberbasic <filename.bas> [options]")
	fmt.Println("Options:")
	fmt.Println("  --compile-only    Compile but don't run")
	fmt.Println("  --gen-go [file]   Generate Go source that calls raylib directly (default: <basename>_gen.go)")
	fmt.Println("  --debug           Enable debug output")
	fmt.Println("  --list-commands   Print built-in command names (2D, 3D, GUI, Physics, Std)")
	fmt.Println("  --lint            Check program (compile only, no run); same as --compile-only")
	fmt.Println("  --repl            Interactive REPL (read-eval-print loop)")
	fmt.Println("  --dev             Live reload (experimental; not fully implemented)")
	fmt.Println("  --debugger        Enable debugger (breakpoints, stack trace)")
	fmt.Println("  --break=5,10      Set breakpoints at lines 5 and 10")
	fmt.Println("  --help            Show this help")
	fmt.Println("  --version         Print version and exit")
	fmt.Println("  (Multi-window: --window --parent=host:port --title=... --width=... --height=...)")
	fmt.Println("Exit codes: 0 = success, 1 = compile/file error, 2 = runtime error")
}
