package app

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"cyberbasic/compiler"
	"cyberbasic/compiler/bindings"
	"cyberbasic/compiler/errors"
	"cyberbasic/compiler/runtime"
)

func runREPL() {
	fmt.Println("CyberBasic REPL - type statements and press Enter. Empty line or QUIT to exit.")
	comp := compiler.New()
	comp.Filename = "<repl>"
	var session strings.Builder
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		if strings.EqualFold(line, "QUIT") || strings.EqualFold(line, "EXIT") {
			fmt.Println("Goodbye!")
			break
		}
		session.WriteString(line)
		session.WriteByte('\n')
		source := session.String()
		chunk, err := comp.Compile(source)
		if err != nil {
			errors.PrettyPrint(os.Stdout, source, "<repl>", err)
			revert := len(source) - len(line) - 1
			if revert > 0 {
				session.Reset()
				session.WriteString(source[:revert])
			} else {
				session.Reset()
			}
			continue
		}
		rt := runtime.NewRuntime()
		rt.GetVM().LoadChunk(chunk)
		stdRegisterEnumsAndRuntime(rt, chunk)
		if err := bindings.RegisterAll(rt.GetVM(), bindings.RegisterOptions{Source: source}); err != nil {
			fmt.Printf("Register bindings: %v\n", err)
			revert := len(source) - len(line) - 1
			if revert > 0 {
				session.Reset()
				session.WriteString(source[:revert])
			} else {
				session.Reset()
			}
			continue
		}
		err = rt.GetVM().Run()
		if err != nil {
			fmt.Printf("Runtime error: %v\n", err)
			revert := len(source) - len(line) - 1
			if revert > 0 {
				session.Reset()
				session.WriteString(source[:revert])
			} else {
				session.Reset()
			}
			continue
		}
		if rt.HasImplicitHandlers() {
			fmt.Println("(Program has OnUpdate/OnDraw - use file mode for graphics)")
		}
	}
}

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
