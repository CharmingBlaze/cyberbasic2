package vm

// GameRuntime is the interface the VM uses to execute game/engine opcodes.
// The runtime package implements this interface so the VM can call into
// graphics, audio, and physics without importing runtime.
type GameRuntime interface {
	LoadImage(filename string) error
	CreateSprite(id, image string, x, y float64) error
	SetSpritePosition(id string, x, y float64) error
	DrawSprite(id string) error
	LoadModel(filename string) error
	CreateCamera(id string, x, y, z float64) error
	SetCameraPosition(id string, x, y, z float64) error
	DrawModel(id string, x, y, z, scale float64) error
	PlayMusic(filename string) error
	PlaySound(filename string) error
	LoadSound(filename string) error
	CreatePhysicsBody(id, bodyType string, x, y, z, mass float64) error
	SetVelocity(id string, vx, vy, vz float64) error
	ApplyForce(id string, fx, fy, fz float64) error
	RayCast3D(startX, startY, startZ, dirX, dirY, dirZ, maxDistance float64) (hit bool, outX, outY, outZ float64, err error)
	InitializeGraphics(width, height int, title string) error
	InitializePhysics() error
	// ShouldClose returns true when the user requested to close the window (e.g. close button)
	ShouldClose() bool
	// Sync completes one frame: draw, swap buffers, wait for target FPS
	Sync() error
	// IsKeyDown returns true if the key (e.g. "ESCAPE", "W") is currently held (for On KeyDown handlers)
	IsKeyDown(keyName string) bool
	// IsKeyPressed returns true if the key was pressed this frame (for On KeyPressed handlers)
	IsKeyPressed(keyName string) bool
}
