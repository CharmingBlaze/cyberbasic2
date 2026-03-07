package runtime

import (
	"cyberbasic/compiler/vm"
	"fmt"
	"strconv"
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Runtime provides the BASIC runtime environment with game API functions
type Runtime struct {
	vm       *vm.VM
	sprites  map[string]*Sprite
	models   map[string]*Model
	cameras  map[string]*Camera
	sounds   map[string]*Sound
	physics  *PhysicsEngine
	graphics *GraphicsEngine
	textures map[string]rl.Texture2D // filename -> texture for LOADIMAGE/CREATESPRITE/DRAWSPRITE
	models3D map[string]rl.Model     // filename -> model for LOADMODEL/DRAWMODEL
}

// Sprite represents a 2D sprite
type Sprite struct {
	ID          string
	Image       string
	X, Y        float64
	Width       float64
	Height      float64
	Visible     bool
	PhysicsBody *PhysicsBody
}

// Model represents a 3D model
type Model struct {
	ID          string
	Mesh        string
	X, Y, Z     float64
	Rotation    [3]float64
	Scale       [3]float64
	Visible     bool
	PhysicsBody *PhysicsBody
}

// Camera represents a 3D camera
type Camera struct {
	ID      string
	X, Y, Z float64
	TargetX float64
	TargetY float64
	TargetZ float64
	UpX     float64
	UpY     float64
	UpZ     float64
	FOV     float64
	Active  bool
}

// Sound represents an audio resource
type Sound struct {
	ID     string
	File   string
	Loaded bool
	Volume float64
}

// PhysicsBody represents a physics object
type PhysicsBody struct {
	ID         string
	Type       string // "box", "sphere", "plane", etc.
	X, Y, Z    float64
	VX, VY, VZ float64
	Mass       float64
	Active     bool
}

// GraphicsEngine handles rendering operations
type GraphicsEngine struct {
	screenWidth  int
	screenHeight int
	fps          int
	title        string
	windowOpen   bool // true after raylib InitWindow
}

// PhysicsEngine handles physics simulation
type PhysicsEngine struct {
	gravity [3]float64
	bodies  map[string]*PhysicsBody
	running bool
}

// NewRuntime creates a new runtime instance
func NewRuntime() *Runtime {
	return &Runtime{
		vm:       vm.NewVM(),
		sprites:  make(map[string]*Sprite),
		models:   make(map[string]*Model),
		cameras:  make(map[string]*Camera),
		sounds:   make(map[string]*Sound),
		textures: make(map[string]rl.Texture2D),
		models3D: make(map[string]rl.Model),
		graphics: &GraphicsEngine{
			screenWidth:  800,
			screenHeight: 600,
			fps:          60,
			title:        "CyberBasic Game",
		},
		physics: &PhysicsEngine{
			gravity: [3]float64{0, -9.81, 0},
			bodies:  make(map[string]*PhysicsBody),
		},
	}
}

// LoadImage loads an image file as a texture (for use with CreateSprite/DrawSprite)
func (r *Runtime) LoadImage(filename string) error {
	tex := rl.LoadTexture(filename)
	r.textures[filename] = tex
	return nil
}

// CreateSprite creates a new sprite (image = filename from LoadImage)
func (r *Runtime) CreateSprite(id, image string, x, y float64) error {
	tex, ok := r.textures[image]
	w, h := float64(64), float64(64)
	if ok && tex.ID != 0 {
		w, h = float64(tex.Width), float64(tex.Height)
	}
	sprite := &Sprite{
		ID:      id,
		Image:   image,
		X:       x,
		Y:       y,
		Width:   w,
		Height:  h,
		Visible: true,
	}
	r.sprites[id] = sprite
	return nil
}

// SetSpritePosition sets the position of a sprite
func (r *Runtime) SetSpritePosition(id string, x, y float64) error {
	sprite, exists := r.sprites[id]
	if !exists {
		return fmt.Errorf("sprite '%s' not found", id)
	}
	sprite.X = x
	sprite.Y = y
	return nil
}

// DrawSprite draws a sprite (called during render loop)
func (r *Runtime) DrawSprite(id string) error {
	sprite, exists := r.sprites[id]
	if !exists {
		return fmt.Errorf("sprite '%s' not found", id)
	}
	if !sprite.Visible {
		return nil
	}
	tex, ok := r.textures[sprite.Image]
	if ok && tex.ID != 0 {
		rl.DrawTexture(tex, int32(sprite.X), int32(sprite.Y), rl.White)
	} else {
		rl.DrawRectangle(int32(sprite.X), int32(sprite.Y), int32(sprite.Width), int32(sprite.Height), rl.White)
	}
	return nil
}

// LoadModel loads a 3D model file and registers it by filename for later DrawModel calls
func (r *Runtime) LoadModel(filename string) error {
	model := rl.LoadModel(filename)
	r.models3D[filename] = model
	r.models[filename] = &Model{
		ID:      filename,
		Mesh:    filename,
		Visible: true,
		Scale:   [3]float64{1, 1, 1},
	}
	return nil
}

// DrawModel updates model position/scale and marks it for drawing on next Render
func (r *Runtime) DrawModel(id string, x, y, z, scale float64) error {
	model, exists := r.models[id]
	if !exists {
		return fmt.Errorf("model '%s' not found", id)
	}
	model.X, model.Y, model.Z = x, y, z
	model.Scale[0], model.Scale[1], model.Scale[2] = scale, scale, scale
	model.Visible = true
	return nil
}

// CreateCamera creates a new 3D camera
func (r *Runtime) CreateCamera(id string, x, y, z float64) error {
	camera := &Camera{
		ID:      id,
		X:       x,
		Y:       y,
		Z:       z,
		TargetX: 0,
		TargetY: 0,
		TargetZ: 0,
		UpX:     0,
		UpY:     1,
		UpZ:     0,
		FOV:     45.0,
		Active:  false,
	}
	r.cameras[id] = camera
	fmt.Printf("Created camera '%s' at (%.1f, %.1f, %.1f)\n", id, x, y, z)
	return nil
}

// SetCameraPosition sets the position of a camera
func (r *Runtime) SetCameraPosition(id string, x, y, z float64) error {
	camera, exists := r.cameras[id]
	if !exists {
		return fmt.Errorf("camera '%s' not found", id)
	}
	camera.X = x
	camera.Y = y
	camera.Z = z
	fmt.Printf("Set camera '%s' position to (%.1f, %.1f, %.1f)\n", id, x, y, z)
	return nil
}

// LoadSound loads a sound file
func (r *Runtime) LoadSound(filename string) error {
	sound := &Sound{
		ID:     filename,
		File:   filename,
		Loaded: true,
		Volume: 1.0,
	}
	r.sounds[filename] = sound
	fmt.Printf("Loaded sound: %s\n", filename)
	return nil
}

// PlaySound plays a sound effect
func (r *Runtime) PlaySound(filename string) error {
	sound, exists := r.sounds[filename]
	if !exists {
		return fmt.Errorf("sound '%s' not loaded", filename)
	}
	fmt.Printf("Playing sound: %s (volume: %.1f)\n", filename, sound.Volume)
	// In a real implementation, this would use Raylib audio
	return nil
}

// PlayMusic plays background music
func (r *Runtime) PlayMusic(filename string) error {
	fmt.Printf("Playing music: %s\n", filename)
	// In a real implementation, this would use Raylib audio
	return nil
}

// CreatePhysicsBody creates a physics body
func (r *Runtime) CreatePhysicsBody(id, bodyType string, x, y, z float64, mass float64) error {
	body := &PhysicsBody{
		ID:     id,
		Type:   bodyType,
		X:      x,
		Y:      y,
		Z:      z,
		VX:     0,
		VY:     0,
		VZ:     0,
		Mass:   mass,
		Active: true,
	}
	r.physics.bodies[id] = body
	fmt.Printf("Created physics body '%s' (%s) at (%.1f, %.1f, %.1f) with mass %.1f\n",
		id, bodyType, x, y, z, mass)
	return nil
}

// SetVelocity sets the velocity of a physics body
func (r *Runtime) SetVelocity(id string, vx, vy, vz float64) error {
	body, exists := r.physics.bodies[id]
	if !exists {
		return fmt.Errorf("physics body '%s' not found", id)
	}
	body.VX = vx
	body.VY = vy
	body.VZ = vz
	fmt.Printf("Set velocity of body '%s' to (%.1f, %.1f, %.1f)\n", id, vx, vy, vz)
	return nil
}

// ApplyForce applies a force to a physics body
func (r *Runtime) ApplyForce(id string, fx, fy, fz float64) error {
	body, exists := r.physics.bodies[id]
	if !exists {
		return fmt.Errorf("physics body '%s' not found", id)
	}

	// In a real implementation, this would apply force using Bullet physics
	// For now, just log the force application
	fmt.Printf("Applied force (%.1f, %.1f, %.1f) to body '%s' (type: %s)\n",
		fx, fy, fz, id, body.Type)

	return nil
}

// RayCast3D performs a 3D raycast
func (r *Runtime) RayCast3D(startX, startY, startZ, dirX, dirY, dirZ, maxDistance float64) (bool, float64, float64, float64, error) {
	fmt.Printf("Raycasting from (%.1f, %.1f, %.1f) in direction (%.1f, %.1f, %.1f) with max distance %.1f\n",
		startX, startY, startZ, dirX, dirY, dirZ, maxDistance)

	// In a real implementation, this would use Bullet physics raycast
	// For now, return a dummy result
	return false, 0, 0, 0, nil
}

// UpdatePhysics updates the physics simulation
func (r *Runtime) UpdatePhysics(deltaTime float64) error {
	if !r.physics.running {
		return nil
	}

	// Simple physics update (in real implementation, use Bullet)
	for _, body := range r.physics.bodies {
		if body.Active {
			// Apply gravity
			body.VY += r.physics.gravity[1] * deltaTime

			// Update position
			body.X += body.VX * deltaTime
			body.Y += body.VY * deltaTime
			body.Z += body.VZ * deltaTime

			// Simple ground collision
			if body.Y < 0 {
				body.Y = 0
				body.VY = 0
			}
		}
	}

	return nil
}

// Render renders the current frame
func (r *Runtime) Render() error {
	// Draw all sprites
	for _, sprite := range r.sprites {
		if sprite.Visible {
			r.DrawSprite(sprite.ID)
		}
	}

	// Draw 3D models (requires BeginMode3D/EndMode3D - caller must ensure)
	for _, model := range r.models {
		if model.Visible {
			rlModel, has3D := r.models3D[model.Mesh]
			if has3D && rlModel.MeshCount > 0 {
				pos := rl.Vector3{X: float32(model.X), Y: float32(model.Y), Z: float32(model.Z)}
				scale := float32(1)
				if model.Scale[0] > 0 {
					scale = float32(model.Scale[0])
				}
				rl.DrawModel(rlModel, pos, scale, rl.White)
			}
		}
	}

	return nil
}

// InitializeGraphics initializes the graphics system and opens the window (raylib)
func (r *Runtime) InitializeGraphics(width, height int, title string) error {
	r.graphics.screenWidth = width
	r.graphics.screenHeight = height
	r.graphics.title = title
	if !r.graphics.windowOpen {
		rl.InitWindow(int32(width), int32(height), title)
		rl.SetTargetFPS(int32(r.graphics.fps))
		r.graphics.windowOpen = true
		// Position window so it's visible when launched from a terminal (not hidden behind it)
		rl.SetWindowPosition(120, 80)
		fmt.Println("Window opened. If you don't see it, check behind the terminal or the taskbar. Close the window to exit.")
	}
	return nil
}

// ShouldClose returns true when the user requested to close the window.
// Uses rl.IsWindowReady() so it works when window was opened via InitWindow or runtime.OpenWindow.
func (r *Runtime) ShouldClose() bool {
	if !rl.IsWindowReady() {
		return false
	}
	return rl.WindowShouldClose()
}

// keyNameToRaylib maps BASIC key name (e.g. "ESCAPE", "W") to raylib key code
func keyNameToRaylib(name string) (int32, bool) {
	n := strings.ToUpper(strings.TrimSpace(name))
	if k, ok := keyNameMap[n]; ok {
		return k, true
	}
	// Strip KEY_ prefix if present
	if strings.HasPrefix(n, "KEY_") {
		if k, ok := keyNameMap[n[4:]]; ok {
			return k, true
		}
	}
	// Try parsing as number (KEY_ESCAPE constant value)
	if i, err := strconv.Atoi(name); err == nil && i >= 0 {
		return int32(i), true
	}
	return 0, false
}

var keyNameMap = map[string]int32{
	"ESCAPE": int32(rl.KeyEscape), "ENTER": int32(rl.KeyEnter), "TAB": int32(rl.KeyTab),
	"BACKSPACE": int32(rl.KeyBackspace), "INSERT": int32(rl.KeyInsert), "DELETE": int32(rl.KeyDelete),
	"RIGHT": int32(rl.KeyRight), "LEFT": int32(rl.KeyLeft), "DOWN": int32(rl.KeyDown), "UP": int32(rl.KeyUp),
	"SPACE": int32(rl.KeySpace), "PAGEUP": int32(rl.KeyPageUp), "PAGEDOWN": int32(rl.KeyPageDown),
	"HOME": int32(rl.KeyHome), "END": int32(rl.KeyEnd),
	"F1": int32(rl.KeyF1), "F2": int32(rl.KeyF2), "F3": int32(rl.KeyF3), "F4": int32(rl.KeyF4),
	"F5": int32(rl.KeyF5), "F6": int32(rl.KeyF6), "F7": int32(rl.KeyF7), "F8": int32(rl.KeyF8),
	"F9": int32(rl.KeyF9), "F10": int32(rl.KeyF10), "F11": int32(rl.KeyF11), "F12": int32(rl.KeyF12),
	"A": int32(rl.KeyA), "B": int32(rl.KeyB), "C": int32(rl.KeyC), "D": int32(rl.KeyD),
	"E": int32(rl.KeyE), "F": int32(rl.KeyF), "G": int32(rl.KeyG), "H": int32(rl.KeyH),
	"I": int32(rl.KeyI), "J": int32(rl.KeyJ), "K": int32(rl.KeyK), "L": int32(rl.KeyL),
	"M": int32(rl.KeyM), "N": int32(rl.KeyN), "O": int32(rl.KeyO), "P": int32(rl.KeyP),
	"Q": int32(rl.KeyQ), "R": int32(rl.KeyR), "S": int32(rl.KeyS), "T": int32(rl.KeyT),
	"U": int32(rl.KeyU), "V": int32(rl.KeyV), "W": int32(rl.KeyW), "X": int32(rl.KeyX),
	"Y": int32(rl.KeyY), "Z": int32(rl.KeyZ),
}

// IsKeyDown returns true if the key (e.g. "ESCAPE", "W") is currently held (for On KeyDown handlers)
func (r *Runtime) IsKeyDown(keyName string) bool {
	if !r.graphics.windowOpen {
		return false
	}
	k, ok := keyNameToRaylib(keyName)
	if !ok {
		return false
	}
	return rl.IsKeyDown(k)
}

// IsKeyPressed returns true if the key was pressed this frame (for On KeyPressed handlers)
func (r *Runtime) IsKeyPressed(keyName string) bool {
	if !r.graphics.windowOpen {
		return false
	}
	k, ok := keyNameToRaylib(keyName)
	if !ok {
		return false
	}
	return rl.IsKeyPressed(k)
}

// Sync runs one frame. Delegates to SyncFrame so SYNC (statement) and Sync() (foreign) behave identically.
// Uses rl.IsWindowReady() so SYNC works when window was opened via InitWindow (raylib) or runtime.OpenWindow.
func (r *Runtime) Sync() error {
	if !rl.IsWindowReady() {
		return nil
	}
	SyncFrame()
	return nil
}

// CloseWindow closes the raylib window if it was opened (call when script ends)
func (r *Runtime) CloseWindow() {
	if r.graphics.windowOpen {
		rl.CloseWindow()
		r.graphics.windowOpen = false
	}
}

// InitializePhysics initializes the physics system
func (r *Runtime) InitializePhysics() error {
	r.physics.running = true
	fmt.Printf("Initialized physics with gravity (%.1f, %.1f, %.1f)\n",
		r.physics.gravity[0], r.physics.gravity[1], r.physics.gravity[2])
	// In a real implementation, this would initialize Bullet physics
	return nil
}

// GetVM returns the virtual machine instance
func (r *Runtime) GetVM() *vm.VM {
	return r.vm
}

// MainLoop runs the main game loop
func (r *Runtime) MainLoop() error {
	fmt.Println("Starting main game loop...")

	// Initialize systems
	err := r.InitializeGraphics(r.graphics.screenWidth, r.graphics.screenHeight, r.graphics.title)
	if err != nil {
		return err
	}

	err = r.InitializePhysics()
	if err != nil {
		return err
	}

	// Main loop (simplified)
	frameCount := 0
	maxFrames := 1000 // Limit for demo purposes

	for frameCount < maxFrames {
		// Update physics
		err = r.UpdatePhysics(1.0 / 60.0) // 60 FPS
		if err != nil {
			return err
		}

		// Render frame
		err = r.Render()
		if err != nil {
			return err
		}

		frameCount++
		if frameCount%60 == 0 {
			fmt.Printf("Frame %d\n", frameCount)
		}
	}

	fmt.Println("Main loop completed")
	return nil
}

// HasImplicitHandlers returns true if the loaded chunk has OnUpdate or OnDraw subs (DBP-style implicit loop).
func (r *Runtime) HasImplicitHandlers() bool {
	chunk := r.vm.Chunk()
	if chunk == nil {
		return false
	}
	_, hasUpdate := chunk.GetFunction("onupdate")
	_, hasDraw := chunk.GetFunction("ondraw")
	return hasUpdate || hasDraw
}

// RunImplicitLoop runs the DBP-style implicit loop: InitWindow, OnStart once, then loop with OnUpdate/OnDraw.
// Call this when HasImplicitHandlers() is true instead of vm.Run().
func (r *Runtime) RunImplicitLoop() error {
	chunk := r.vm.Chunk()
	if chunk == nil {
		return fmt.Errorf("no chunk loaded")
	}

	// Init window and FPS
	if err := r.InitializeGraphics(1280, 720, "CyberBASIC 2"); err != nil {
		return err
	}
	rl.SetTargetFPS(60)

	// OnStart once
	if _, ok := chunk.GetFunction("onstart"); ok {
		if err := r.vm.InvokeSub("OnStart", nil); err != nil {
			return err
		}
	}

	// Main loop
	for !r.ShouldClose() {
		if err := StepImplicitFrame(r.vm); err != nil {
			return err
		}
	}

	return nil
}
