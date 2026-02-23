// +build ignore

// Run with: go run test_raylib_window.go
// If a window opens, raylib works and the issue is in the BASIC/VM path.
// If no window or crash, the issue is raylib/OpenGL/drivers on this machine.
package main

import (
	"fmt"
	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	fmt.Println("Opening raylib window...")
	rl.InitWindow(800, 600, "Raylib Test")
	defer rl.CloseWindow()
	if !rl.IsWindowReady() {
		fmt.Println("ERROR: IsWindowReady() returned false - window/OpenGL init failed")
		return
	}
	fmt.Println("Window ready. Close the window to exit.")
	rl.SetTargetFPS(60)
	for !rl.WindowShouldClose() {
		rl.BeginDrawing()
		rl.ClearBackground(rl.DarkGray)
		rl.DrawText("If you see this, raylib works!", 50, 50, 20, rl.White)
		rl.EndDrawing()
	}
	fmt.Println("Done.")
}
