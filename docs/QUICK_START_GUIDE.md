# CyberBasic Quick Start Guide

Get started with CyberBasic in minutes! This guide will take you from zero to running your first game.

---

## What is CyberBasic?

CyberBasic is a modern BASIC-inspired programming language designed specifically for game development. It combines:
- **Simple BASIC syntax** - Easy to learn and read
- **Powerful graphics** - Built-in 2D/3D rendering via raylib
- **Physics engines** - Box2D for 2D, Bullet for 3D
- **Networking** - Built-in multiplayer support
- **Cross-platform** - Runs on Windows, Mac, and Linux

---

## Installation

### Option 1: Download Binary (Easiest)

1. Go to the [GitHub Releases](https://github.com/CharmingBlaze/cyberbasic2/releases)
2. Download the latest binary for your platform
3. Extract to a folder (e.g., `C:\CyberBasic` on Windows)
4. Add to PATH or run from the folder

### Option 2: Build from Source

```bash
# Clone the repository
git clone https://github.com/CharmingBlaze/cyberbasic2.git
cd cyberbasic2

# Build the binary
go build -o cyberbasic .

# Run your first program
./cyberbasic examples/first_game.bas
```

---

## Your First Program

Let's create a simple window that says "Hello, World!":

```basic
// hello.bas
PRINT "Hello, CyberBasic!"
```

Run it:
```bash
cyberbasic hello.bas
```

---

## Your First Graphics Program

Now let's create a window with graphics:

```basic
// first_graphics.bas
InitWindow(800, 600, "My First Game")
SetTargetFPS(60)

WHILE NOT WindowShouldClose()
    ClearBackground(20, 20, 40, 255)
    DrawText("Hello, Graphics!", 300, 280, 30, 255, 255, 255, 255)
WEND

CloseWindow()
```

Run it:
```bash
cyberbasic first_graphics.bas
```

You should see a window with white text on a dark background!

---

## Your First Interactive Game

Let's make something you can control:

```basic
// interactive_game.bas
InitWindow(800, 600, "Move the Circle")
SetTargetFPS(60)

VAR x = 400
VAR y = 300

WHILE NOT WindowShouldClose()
    // Move with arrow keys
    IF IsKeyDown(KEY_LEFT) THEN x = x - 5
    IF IsKeyDown(KEY_RIGHT) THEN x = x + 5
    IF IsKeyDown(KEY_UP) THEN y = y - 5
    IF IsKeyDown(KEY_DOWN) THEN y = y + 5
    
    // Keep on screen
    x = Clamp(x, 20, 780)
    y = Clamp(y, 20, 580)
    
    // Drawing
    ClearBackground(30, 30, 50, 255)
    DrawCircle(x, y, 20, 100, 200, 255, 255)
    DrawText("Use arrow keys to move", 10, 10, 20, 255, 255, 255, 255)
WEND

CloseWindow()
```

---

## Learning Path

Now that you have the basics, follow this structured learning path:

### 1. Language Fundamentals (30 minutes)
Learn the BASIC syntax, variables, functions, and control flow.

**Start with**: [LEARNING_PATH.md](LEARNING_PATH.md) - Module 1

### 2. 2D Game Development (2 hours)
Create complete 2D games with graphics, collision, and physics.

**Start with**: [TUTORIAL_2D_GAMES.md](TUTORIAL_2D_GAMES.md)

### 3. 3D Game Development (3 hours)
Master 3D graphics, cameras, and physics.

**Start with**: [TUTORIAL_3D_GAMES.md](TUTORIAL_3D_GAMES.md)

### 4. GUI Development (1 hour)
Create menus, HUDs, and user interfaces.

**Start with**: [TUTORIAL_GUI_DEVELOPMENT.md](TUTORIAL_GUI_DEVELOPMENT.md)

### 5. Multiplayer Games (2 hours)
Add networking and create multiplayer experiences.

**Start with**: [TUTORIAL_MULTIPLAYER.md](TUTORIAL_MULTIPLAYER.md)

---

## Quick Reference

### Essential Functions

**Window:**
```basic
InitWindow(width, height, title)    // Create window
WindowShouldClose()                // Check if user wants to quit
ClearBackground(r, g, b, a)        // Clear screen
CloseWindow()                      // Close window
```

**Input:**
```basic
IsKeyDown(KEY_W)                   // Key held down
IsKeyPressed(KEY_SPACE)            // Key pressed once
GetMouseX(), GetMouseY()           // Mouse position
IsMouseButtonPressed(0)             // Mouse click
```

**Drawing:**
```basic
DrawText(text, x, y, size, r, g, b, a)
DrawCircle(x, y, radius, r, g, b, a)
DrawRectangle(x, y, width, height, r, g, b, a)
DrawLine(x1, y1, x2, y2, r, g, b, a)
```

**Game Loop Pattern:**
```basic
InitWindow(800, 600, "Game")
SetTargetFPS(60)

WHILE NOT WindowShouldClose()
    // Handle input
    // Update game logic
    // Draw everything
WEND

CloseWindow()
```

---

## Common Key Codes

```
KEY_A, KEY_B, KEY_C, ..., KEY_Z     // Letter keys
KEY_0, KEY_1, ..., KEY_9            // Number keys
KEY_UP, KEY_DOWN, KEY_LEFT, KEY_RIGHT // Arrow keys
KEY_SPACE, KEY_ENTER, KEY_ESCAPE    // Special keys
KEY_LEFT_SHIFT, KEY_RIGHT_SHIFT     // Modifier keys
```

---

## Colors

Colors use RGBA format (Red, Green, Blue, Alpha), each value 0-255:

```basic
// Common colors
DrawText("Red", 10, 10, 20, 255, 0, 0, 255)
DrawText("Green", 10, 40, 20, 0, 255, 0, 255)
DrawText("Blue", 10, 70, 20, 0, 0, 255, 255)
DrawText("White", 10, 100, 20, 255, 255, 255, 255)
DrawText("Black", 10, 130, 20, 0, 0, 0, 255)
```

---

## Project Ideas for Practice

Once you're comfortable with the basics, try these projects:

### Beginner (1-2 hours)
- **Pong** - Classic paddle game
- **Snake** - Eat and grow game
- **Breakout** - Break bricks with ball
- **Tic-Tac-Toe** - 3x3 grid game

### Intermediate (3-5 hours)
- **Platformer** - Jump between platforms
- **Space Shooter** - Shoot asteroids
- **Puzzle Game** - Match-3 or sliding puzzle
- **Racing Game** - Top-down racing

### Advanced (1-2 weeks)
- **RPG** - Character stats and inventory
- **Strategy Game** - Resource management
- **3D Explorer** - First-person movement
- **Multiplayer Game** - Networked play

---

## Getting Help

### Documentation
- **[LEARNING_PATH.md](LEARNING_PATH.md)** - Complete structured learning
- **[LANGUAGE_SPEC.md](../LANGUAGE_SPEC.md)** - Full language reference
- **[API_REFERENCE.md](../API_REFERENCE.md)** - All functions listed

### Examples
The `examples/` folder contains working code for every feature:
```bash
# Try some examples
cyberbasic examples/first_game.bas
cyberbasic examples/box2d_demo.bas
cyberbasic examples/3d_graphics_demo.bas
cyberbasic examples/gui_demo.bas
```

### Troubleshooting
- **"Command not found"** - Make sure cyberbasic is in your PATH
- **"Parse error"** - Check syntax and quotes
- **"Window appears and closes"** - Add a game loop
- **"No sound"** - Initialize audio device first

### Community
- **GitHub Issues** - Report bugs and request features
- **Examples** - Learn from existing code
- **Documentation** - Comprehensive guides available

---

## Next Steps

1. **Complete Module 1** of the learning path
2. **Try the examples** in the examples folder
3. **Build a simple game** from the project ideas
4. **Explore advanced features** like physics and networking
5. **Share your creations** with the community

---

## Tips for Success

### Start Simple
Don't try to build a complex MMO as your first game. Start with Pong or Snake and work your way up.

### Use Examples
The examples folder is your best friend. Study them, modify them, and learn from them.

### Experiment
Try changing values, adding features, and breaking things. That's how you learn!

### Save Often
Keep backup copies of your work. Use version control if you know Git.

### Ask Questions
If you're stuck, look at the documentation or ask for help.

---

## You're Ready!

You now have everything you need to start making games with CyberBasic. The journey from beginner to game developer starts with a single line of code.

**Your first game is waiting to be created. Happy coding!**
