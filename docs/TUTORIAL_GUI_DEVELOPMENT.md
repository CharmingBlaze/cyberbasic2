# GUI Development Tutorial - Complete Guide

Welcome to the complete GUI development tutorial! This guide will teach you how to create user interfaces, menus, and interactive elements for your CyberBasic games.

## What You'll Build

By the end of this tutorial, you'll have created:
- Complete menu systems
- Interactive UI elements
- Game settings screens
- HUD (Heads-Up Display) elements
- Custom UI components
- Responsive layouts

---

## Prerequisites

Before starting, make sure you've completed:
- **Module 1**: BASIC Programming Fundamentals (from LEARNING_PATH.md)
- **Module 2**: 2D Game Development (recommended for context)

---

## Lesson 1: Basic GUI Elements

### Understanding GUI Mode

CyberBasic provides two main approaches to GUI:
1. **Immediate Mode GUI** - Simple, stateless UI elements
2. **Custom Drawing** - Full control over appearance

```basic
// Basic GUI Elements Demo
InitWindow(800, 600, "GUI Elements Demo")
SetTargetFPS(60)

// GUI state variables
VAR sliderValue = 0.5
VAR checkboxState = 0
VAR buttonClicked = 0
VAR textInput = "Hello World"
VAR dropdownSelected = 0
VAR progressBarValue = 0.3

WHILE NOT WindowShouldClose()
    ClearBackground(60, 60, 80, 255)
    
    // Begin GUI mode
    BeginUI()
    
    // Labels
    Label("=== GUI Elements Demo ===")
    Label("")
    
    // Button
    buttonClicked = Button("Click Me!")
    IF buttonClicked THEN
        PRINT "Button was clicked!"
    ENDIF
    
    // Slider
    Label("Slider Value:")
    sliderValue = Slider("Adjust Me", sliderValue, 0, 1)
    Label("Current value: " + STR(Int(sliderValue * 100)) + "%")
    
    // Checkbox
    checkboxState = Checkbox("Enable Feature", checkboxState)
    IF checkboxState = 1 THEN
        Label("Feature is ENABLED")
    ELSE
        Label("Feature is DISABLED")
    ENDIF
    
    // Progress Bar
    Label("Progress:")
    progressBarValue = progressBarValue + 0.01
    IF progressBarValue > 1 THEN progressBarValue = 0
    ProgressBar("Loading", progressBarValue)
    
    // Text Input (basic simulation)
    Label("Text Input:")
    Label("Current: " + textInput)
    IF Button("Change Text") THEN
        textInput = "Updated at " + STR(Int(GetTime()))
    ENDIF
    
    // Dropdown (simulated with buttons)
    Label("Dropdown:")
    IF Button("Option 1") THEN dropdownSelected = 0
    IF Button("Option 2") THEN dropdownSelected = 1
    IF Button("Option 3") THEN dropdownSelected = 2
    Label("Selected: Option " + STR(dropdownSelected + 1))
    
    // End GUI mode
    EndUI()
    
    // Additional drawing outside GUI
    DrawText("ESC to exit", 10, 570, 16, 200, 200, 200, 255)
WEND

CloseWindow()
```

### Custom GUI Drawing

For more control, you can draw GUI elements manually:

```basic
// Custom GUI Elements
InitWindow(800, 600, "Custom GUI")
SetTargetFPS(60)

// Custom button state
VAR customButtonX = 300
VAR customButtonY = 200
VAR customButtonWidth = 200
VAR customButtonHeight = 50
VAR customButtonHover = 0
VAR customButtonClicks = 0

FUNCTION IsMouseOverButton(x, y, width, height)
    VAR mx = GetMouseX()
    VAR my = GetMouseY()
    RETURN mx >= x AND mx <= x + width AND my >= y AND my <= y + height
END FUNCTION

FUNCTION CustomButton(x, y, width, height, text)
    VAR isHover = IsMouseOverButton(x, y, width, height)
    VAR isClicked = 0
    
    // Draw button background
    IF isHover THEN
        DrawRectangle(x, y, width, height, 100, 150, 200, 255)
        DrawRectangleLines(x, y, width, height, 150, 200, 250, 255)
    ELSE
        DrawRectangle(x, y, width, height, 70, 100, 150, 255)
        DrawRectangleLines(x, y, width, height, 100, 130, 180, 255)
    ENDIF
    
    // Draw text centered
    VAR textWidth = MeasureText(text, 20)
    VAR textX = x + (width - textWidth) / 2
    VAR textY = y + (height - 20) / 2
    DrawText(text, textX, textY, 20, 255, 255, 255, 255)
    
    // Check for click
    IF isHover AND IsMouseButtonPressed(0) THEN
        isClicked = 1
    ENDIF
    
    RETURN isClicked
END FUNCTION

WHILE NOT WindowShouldClose()
    ClearBackground(40, 40, 60, 255)
    
    // Title
    DrawText("Custom GUI Elements", 250, 50, 30, 255, 255, 255, 255)
    
    // Custom button
    IF CustomButton(customButtonX, customButtonY, customButtonWidth, customButtonHeight, "Custom Button") THEN
        customButtonClicks = customButtonClicks + 1
    ENDIF
    
    // Show click count
    DrawText("Button clicked: " + STR(customButtonClicks) + " times", 300, 270, 20, 255, 255, 255, 255)
    
    // Custom slider
    DrawText("Custom Slider:", 300, 320, 18, 255, 255, 255, 255)
    DrawRectangle(300, 350, 200, 4, 100, 100, 100, 255)
    
    VAR sliderX = 300 + Int(GetMouseX() / 800.0 * 200)
    sliderX = Clamp(sliderX, 300, 500)
    DrawCircle(sliderX, 352, 8, 200, 200, 255, 255)
    
    // Custom checkbox
    DrawText("Custom Checkbox:", 300, 400, 18, 255, 255, 255, 255)
    DrawRectangleLines(300, 430, 20, 20, 200, 200, 200, 255)
    IF IsMouseOverButton(300, 430, 20, 20) AND IsMouseButtonPressed(0) THEN
        DrawRectangle(302, 432, 16, 16, 100, 255, 100, 255)
    ENDIF
    
    DrawText("Click to toggle", 330, 432, 16, 200, 200, 200, 255)
    
WEND

CloseWindow()
```

---

## Lesson 2: Menu Systems

### Main Menu Structure

```basic
// Complete Menu System
InitWindow(800, 600, "Menu System")
SetTargetFPS(60)

// Menu states
VAR currentMenu = "main"  // "main", "options", "credits", "game"

// Menu options
VAR mainMenuSelected = 0
VAR optionsMenuSelected = 0
VAR musicVolume = 0.7
VAR soundVolume = 0.8
VAR fullscreen = 0
VAR difficulty = 1

FUNCTION DrawMenuHeader(title)
    DrawText(title, 250, 100, 40, 255, 255, 255, 255)
    DrawText("================================", 200, 150, 20, 150, 150, 150, 255)
END FUNCTION

FUNCTION DrawMenuItem(text, y, selected)
    IF selected THEN
        DrawText("> " + text + " <", 250, y, 25, 255, 215, 0, 255)
    ELSE
        DrawText("  " + text + "  ", 250, y, 25, 255, 255, 255, 255)
    ENDIF
END FUNCTION

WHILE NOT WindowShouldClose()
    ClearBackground(20, 20, 40, 255)
    
    // Handle input
    IF IsKeyPressed(KEY_UP) THEN
        SELECT CASE currentMenu
            CASE "main": mainMenuSelected = (mainMenuSelected - 1 + 4) % 4
            CASE "options": optionsMenuSelected = (optionsMenuSelected - 1 + 5) % 5
        END SELECT
    ENDIF
    
    IF IsKeyPressed(KEY_DOWN) THEN
        SELECT CASE currentMenu
            CASE "main": mainMenuSelected = (mainMenuSelected + 1) % 4
            CASE "options": optionsMenuSelected = (optionsMenuSelected + 1) % 5
        END SELECT
    ENDIF
    
    IF IsKeyPressed(KEY_ENTER) OR IsKeyPressed(KEY_SPACE) THEN
        SELECT CASE currentMenu
            CASE "main"
                SELECT CASE mainMenuSelected
                    CASE 0: currentMenu = "game"
                    CASE 1: currentMenu = "options"
                    CASE 2: currentMenu = "credits"
                    CASE 3: EXIT WHILE
                END SELECT
            CASE "options"
                SELECT CASE optionsMenuSelected
                    CASE 0: musicVolume = 1.0 - musicVolume
                    CASE 1: soundVolume = 1.0 - soundVolume
                    CASE 2: fullscreen = 1 - fullscreen
                    CASE 3: difficulty = (difficulty + 1) % 3
                    CASE 4: currentMenu = "main"
                END SELECT
            CASE "credits"
                currentMenu = "main"
            CASE "game"
                currentMenu = "main"
        END SELECT
    ENDIF
    
    IF IsKeyPressed(KEY_ESCAPE) THEN
        IF currentMenu = "main" THEN
            EXIT WHILE
        ELSE
            currentMenu = "main"
        ENDIF
    ENDIF
    
    // Draw appropriate menu
    SELECT CASE currentMenu
        CASE "main"
            DrawMenuHeader("CYBERBASIC GAME")
            DrawMenuItem("Start Game", 220, mainMenuSelected = 0)
            DrawMenuItem("Options", 270, mainMenuSelected = 1)
            DrawMenuItem("Credits", 320, mainMenuSelected = 2)
            DrawMenuItem("Exit", 370, mainMenuSelected = 3)
            
            DrawText("Use UP/DOWN arrows to navigate", 220, 450, 16, 200, 200, 200, 255)
            DrawText("Press ENTER or SPACE to select", 220, 470, 16, 200, 200, 200, 255)
            DrawText("Press ESC to exit", 220, 490, 16, 200, 200, 200, 255)
            
        CASE "options"
            DrawMenuHeader("OPTIONS")
            
            // Music volume
            VAR musicText = "Music Volume: " + IIF(musicVolume > 0.5, "ON", "OFF")
            DrawMenuItem(musicText, 220, optionsMenuSelected = 0)
            
            // Sound volume
            VAR soundText = "Sound Volume: " + IIF(soundVolume > 0.5, "ON", "OFF")
            DrawMenuItem(soundText, 270, optionsMenuSelected = 1)
            
            // Fullscreen
            VAR fullscreenText = "Fullscreen: " + IIF(fullscreen = 1, "ON", "OFF")
            DrawMenuItem(fullscreenText, 320, optionsMenuSelected = 2)
            
            // Difficulty
            VAR diffText = "Difficulty: " + IIF(difficulty = 0, "Easy", IIF(difficulty = 1, "Normal", "Hard"))
            DrawMenuItem(diffText, 370, optionsMenuSelected = 3)
            
            DrawMenuItem("Back", 420, optionsMenuSelected = 4)
            
        CASE "credits"
            DrawMenuHeader("CREDITS")
            DrawText("Game developed with CyberBasic", 200, 220, 20, 255, 255, 255, 255)
            DrawText("Programming: Your Name", 250, 260, 18, 200, 200, 200, 255)
            DrawText("Graphics: Raylib", 280, 290, 18, 200, 200, 200, 255)
            DrawText("Physics: Box2D & Bullet", 260, 320, 18, 200, 200, 200, 255)
            DrawText("Thank you for playing!", 240, 380, 20, 255, 215, 0, 255)
            DrawText("Press any key to return", 250, 450, 16, 200, 200, 200, 255)
            
            IF IsKeyPressed(KEY_ANY) THEN
                currentMenu = "main"
            ENDIF
            
        CASE "game"
            DrawMenuHeader("GAME PLAYING")
            DrawText("This is where your game would be!", 200, 250, 20, 255, 255, 255, 255)
            DrawText("Press ESC to return to menu", 250, 350, 16, 200, 200, 200, 255)
            
            IF IsKeyPressed(KEY_ESCAPE) THEN
                currentMenu = "main"
            ENDIF
    END SELECT
    
WEND

CloseWindow()
```

---

## Lesson 3: HUD (Heads-Up Display)

### Game HUD Elements

```basic
// Game HUD System
InitWindow(1024, 768, "Game HUD")
SetTargetFPS(60)

// Game state
VAR playerHealth = 100
VAR playerMaxHealth = 100
VAR playerScore = 0
VAR playerLevel = 1
VAR playerAmmo = 30
VAR playerMaxAmmo = 30
VAR gameTime = 0

// Simulated game updates
VAR lastHealthUpdate = 0
VAR lastScoreUpdate = 0

WHILE NOT WindowShouldClose()
    VAR dt = GetFrameTime()
    gameTime = gameTime + dt
    
    // Simulate game events
    gameTime = gameTime + dt
    
    // Random health changes
    IF gameTime - lastHealthUpdate > 2.0 THEN
        VAR healthChange = GetRandomValue(-10, 5)
        playerHealth = Clamp(playerHealth + healthChange, 0, playerMaxHealth)
        lastHealthUpdate = gameTime
    ENDIF
    
    // Random score increases
    IF gameTime - lastScoreUpdate > 1.5 THEN
        playerScore = playerScore + GetRandomValue(10, 50)
        lastScoreUpdate = gameTime
    ENDIF
    
    // Ammo consumption
    IF IsMouseButtonPressed(0) AND playerAmmo > 0 THEN
        playerAmmo = playerAmmo - 1
    ENDIF
    IF IsKeyPressed(KEY_R) THEN
        playerAmmo = playerMaxAmmo
    ENDIF
    
    // Clear screen (game background)
    ClearBackground(30, 30, 50, 255)
    
    // Simulate game world
    DrawCircle(512, 384, 50, 100, 200, 255, 255)
    DrawText("GAME WORLD", 450, 500, 20, 255, 255, 255, 255)
    
    // Draw HUD elements
    // Health bar
    DrawText("Health", 20, 20, 18, 255, 255, 255, 255)
    DrawRectangleLines(20, 45, 200, 25, 255, 255, 255, 255)
    VAR healthWidth = Int((playerHealth / playerMaxHealth) * 198)
    IF playerHealth > 60 THEN
        DrawRectangle(21, 46, healthWidth, 23, 0, 255, 0, 255)
    ELSEIF playerHealth > 30 THEN
        DrawRectangle(21, 46, healthWidth, 23, 255, 255, 0, 255)
    ELSE
        DrawRectangle(21, 46, healthWidth, 23, 255, 0, 0, 255)
    ENDIF
    DrawText(STR(playerHealth) + "/" + STR(playerMaxHealth), 25, 48, 14, 255, 255, 255, 255)
    
    // Score
    DrawText("Score", 20, 80, 18, 255, 255, 255, 255)
    DrawText(STR(playerScore), 20, 105, 24, 255, 215, 0, 255)
    
    // Level
    DrawText("Level", 20, 140, 18, 255, 255, 255, 255)
    DrawText(STR(playerLevel), 20, 165, 24, 255, 255, 255, 255)
    
    // Ammo
    DrawText("Ammo", 20, 200, 18, 255, 255, 255, 255)
    DrawRectangleLines(20, 225, 150, 20, 255, 255, 255, 255)
    FOR i = 0 TO playerAmmo - 1
        IF i < 15 THEN  // Show max 15 bullets
            VAR bulletX = 25 + (i * 9)
            DrawRectangle(bulletX, 229, 6, 12, 255, 255, 100, 255)
        ENDIF
    NEXT i
    DrawText(STR(playerAmmo) + "/" + STR(playerMaxAmmo), 25, 250, 14, 200, 200, 200, 255)
    
    // Minimap (top-right corner)
    DrawRectangleLines(850, 20, 150, 150, 255, 255, 255, 255)
    DrawRectangle(851, 21, 148, 148, 40, 40, 60, 255)
    
    // Player position on minimap
    VAR minimapPlayerX = 850 + 75
    VAR minimapPlayerY = 21 + 75
    DrawCircle(minimapPlayerX, minimapPlayerY, 3, 0, 255, 0, 255)
    
    // Enemies on minimap
    DrawCircle(900, 50, 2, 255, 0, 0, 255)
    DrawCircle(870, 120, 2, 255, 0, 0, 255)
    DrawCircle(920, 90, 2, 255, 0, 0, 255)
    
    DrawText("MINIMAP", 880, 175, 12, 200, 200, 200, 255)
    
    // Crosshair (center screen)
    DrawLine(512 - 10, 384, 512 + 10, 384, 255, 255, 255, 255)
    DrawLine(512, 384 - 10, 512, 384 + 10, 255, 255, 255, 255)
    DrawCircle(512, 384, 15, 255, 255, 255, 100)
    
    // Game time
    VAR minutes = Int(gameTime / 60)
    VAR seconds = Int(gameTime % 60)
    VAR timeText = "Time: " + STR(minutes) + ":" + IIF(seconds < 10, "0", "") + STR(seconds)
    DrawText(timeText, 450, 20, 18, 255, 255, 255, 255)
    
    // Controls help
    DrawText("Click: Shoot | R: Reload | ESC: Menu", 350, 730, 16, 200, 200, 200, 255)
    
    // Warning messages
    IF playerHealth < 30 THEN
        DrawText("LOW HEALTH!", 400, 100, 30, 255, 0, 0, 255)
    ENDIF
    
    IF playerAmmo < 5 THEN
        DrawText("LOW AMMO! Press R to reload", 350, 130, 20, 255, 255, 0, 255)
    ENDIF
    
WEND

CloseWindow()
```

---

## Lesson 4: Advanced UI Components

### Tab System

```basic
// Tab Interface System
InitWindow(800, 600, "Tab System")
SetTargetFPS(60)

// Tab system
VAR activeTab = 0
VAR tabCount = 4
VAR tabNames[4] = ["Inventory", "Stats", "Skills", "Map"]
VAR tabWidth = 150
VAR tabHeight = 40
VAR tabX = 50
VAR tabY = 50

// Tab content data
VAR inventoryItems[5] = ["Sword", "Shield", "Potion", "Key", "Map"]
VAR playerStats[4] = [15, 10, 8, 12]  // STR, DEX, INT, CON
VAR skillNames[3] = ["Fireball", "Heal", "Lightning"]
VAR skillLevels[3] = [3, 2, 1]

FUNCTION DrawTab(x, y, width, height, text, isActive)
    IF isActive THEN
        // Active tab
        DrawRectangle(x, y, width, height, 80, 100, 120, 255)
        DrawRectangleLines(x, y, width, height, 150, 170, 190, 255)
        DrawText(text, x + 20, y + 10, 16, 255, 255, 255, 255)
    ELSE
        // Inactive tab
        DrawRectangle(x, y, width, height, 50, 60, 70, 255)
        DrawRectangleLines(x, y, width, height, 100, 110, 120, 255)
        DrawText(text, x + 20, y + 10, 16, 200, 200, 200, 255)
    ENDIF
END FUNCTION

FUNCTION IsTabClicked(x, y, width, height)
    VAR mx = GetMouseX()
    VAR my = GetMouseY()
    RETURN mx >= x AND mx <= x + width AND my >= y AND my <= y + height
END FUNCTION

WHILE NOT WindowShouldClose()
    ClearBackground(30, 30, 50, 255)
    
    // Handle tab clicks
    IF IsMouseButtonPressed(0) THEN
        FOR i = 0 TO tabCount - 1
            VAR currentTabX = tabX + (i * (tabWidth + 10))
            IF IsTabClicked(currentTabX, tabY, tabWidth, tabHeight) THEN
                activeTab = i
            ENDIF
        NEXT i
    ENDIF
    
    // Draw tabs
    FOR i = 0 TO tabCount - 1
        VAR currentTabX = tabX + (i * (tabWidth + 10))
        DrawTab(currentTabX, tabY, tabWidth, tabHeight, tabNames[i], activeTab = i)
    NEXT i
    
    // Draw tab content area
    DrawRectangle(40, 100, 720, 450, 40, 50, 60, 255)
    DrawRectangleLines(40, 100, 720, 450, 100, 110, 120, 255)
    
    // Draw content based on active tab
    SELECT CASE activeTab
        CASE 0  // Inventory
            DrawText("INVENTORY", 60, 120, 24, 255, 255, 255, 255)
            DrawText("Items:", 60, 160, 18, 255, 255, 255, 255)
            
            FOR i = 0 TO 4
                DrawRectangle(60, 190 + (i * 40), 300, 35, 60, 70, 80, 255)
                DrawRectangleLines(60, 190 + (i * 40), 300, 35, 100, 110, 120, 255)
                DrawText(inventoryItems[i], 70, 200 + (i * 40), 16, 255, 255, 255, 255)
            NEXT i
            
        CASE 1  // Stats
            DrawText("CHARACTER STATS", 60, 120, 24, 255, 255, 255, 255)
            
            DrawText("Strength: " + STR(playerStats[0]), 60, 180, 18, 255, 100, 100, 255)
            DrawRectangle(250, 180, Int(playerStats[0] * 10), 15, 255, 100, 100, 255)
            
            DrawText("Dexterity: " + STR(playerStats[1]), 60, 220, 18, 100, 255, 100, 255)
            DrawRectangle(250, 220, Int(playerStats[1] * 10), 15, 100, 255, 100, 255)
            
            DrawText("Intelligence: " + STR(playerStats[2]), 60, 260, 18, 100, 100, 255, 255)
            DrawRectangle(250, 260, Int(playerStats[2] * 10), 15, 100, 100, 255, 255)
            
            DrawText("Constitution: " + STR(playerStats[3]), 60, 300, 18, 255, 255, 100, 255)
            DrawRectangle(250, 300, Int(playerStats[3] * 10), 15, 255, 255, 100, 255)
            
        CASE 2  // Skills
            DrawText("SKILLS", 60, 120, 24, 255, 255, 255, 255)
            
            FOR i = 0 TO 2
                DrawText(skillNames[i], 60, 180 + (i * 60), 18, 255, 255, 255, 255)
                
                // Skill level stars
                FOR j = 0 TO 4
                    VAR starX = 250 + (j * 25)
                    VAR starY = 180 + (i * 60)
                    IF j < skillLevels[i] THEN
                        DrawText("*", starX, starY, 20, 255, 215, 0, 255)
                    ELSE
                        DrawText("*", starX, starY, 20, 100, 100, 100, 255)
                    ENDIF
                NEXT j
            NEXT i
            
        CASE 3  // Map
            DrawText("WORLD MAP", 60, 120, 24, 255, 255, 255, 255)
            
            // Draw simple map
            DrawRectangle(60, 160, 400, 300, 80, 120, 80, 255)
            DrawRectangleLines(60, 160, 400, 300, 150, 190, 150, 255)
            
            // Locations on map
            DrawCircle(150, 250, 8, 255, 255, 255, 255)  // Current location
            DrawText("You are here", 170, 240, 14, 255, 255, 255, 255)
            
            DrawCircle(300, 200, 6, 255, 100, 100, 255)   // Enemy location
            DrawText("Danger", 310, 190, 12, 255, 100, 100, 255)
            
            DrawCircle(380, 350, 6, 100, 255, 100, 255)   // Safe location
            DrawText("Town", 390, 340, 12, 100, 255, 100, 255)
            
            DrawCircle(200, 400, 6, 255, 255, 100, 255)   // Treasure
            DrawText("Treasure", 210, 390, 12, 255, 255, 100, 255)
    END SELECT
    
    // Instructions
    DrawText("Click on tabs to switch between sections", 60, 570, 16, 200, 200, 200, 255)
    
WEND

CloseWindow()
```

---

## Lesson 5: Complete GUI Game

### RPG Inventory System

```basic
// Complete RPG GUI Game
InitWindow(1024, 768, "RPG Inventory System")
SetTargetFPS(60)

// Game state
VAR gameState = "menu"  // "menu", "playing", "inventory", "dialogue"

// Player stats
VAR playerHealth = 100
VAR playerMaxHealth = 100
VAR playerMana = 50
VAR playerMaxMana = 50
VAR playerLevel = 1
VAR playerExp = 0
VAR playerExpToNext = 100
VAR playerGold = 100

// Inventory system
VAR inventorySize = 20
VAR inventoryNames[20] = ["", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""]
VAR inventoryTypes[20] = ["", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""]
VAR inventoryQuantities[20] = [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0]

// Initialize with some items
inventoryNames[0] = "Health Potion"
inventoryTypes[0] = "consumable"
inventoryQuantities[0] = 5

inventoryNames[1] = "Iron Sword"
inventoryTypes[1] = "weapon"
inventoryQuantities[1] = 1

inventoryNames[2] = "Leather Armor"
inventoryTypes[2] = "armor"
inventoryQuantities[2] = 1

inventoryNames[3] = "Mana Potion"
inventoryTypes[3] = "consumable"
inventoryQuantities[3] = 3

// UI state
VAR selectedSlot = 0
VAR inventoryOpen = 0
VAR messageText = ""
VAR messageTimer = 0

FUNCTION AddItem(name, type, quantity)
    // Find empty slot
    FOR i = 0 TO inventorySize - 1
        IF inventoryNames[i] = "" THEN
            inventoryNames[i] = name
            inventoryTypes[i] = type
            inventoryQuantities[i] = quantity
            RETURN 1
        ENDIF
    NEXT i
    RETURN 0
END FUNCTION

FUNCTION UseItem(slot)
    IF inventoryNames[slot] <> "" THEN
        SELECT CASE inventoryNames[slot]
            CASE "Health Potion"
                playerHealth = Clamp(playerHealth + 30, 0, playerMaxHealth)
                messageText = "Used Health Potion! +30 HP"
                inventoryQuantities[slot] = inventoryQuantities[slot] - 1
                IF inventoryQuantities[slot] <= 0 THEN
                    inventoryNames[slot] = ""
                    inventoryTypes[slot] = ""
                ENDIF
            CASE "Mana Potion"
                playerMana = Clamp(playerMana + 20, 0, playerMaxMana)
                messageText = "Used Mana Potion! +20 MP"
                inventoryQuantities[slot] = inventoryQuantities[slot] - 1
                IF inventoryQuantities[slot] <= 0 THEN
                    inventoryNames[slot] = ""
                    inventoryTypes[slot] = ""
                ENDIF
            CASE ELSE
                messageText = "Cannot use " + inventoryNames[slot]
        END SELECT
        messageTimer = 3.0
    ENDIF
END FUNCTION

FUNCTION DrawInventory()
    // Inventory background
    DrawRectangle(150, 100, 724, 500, 40, 40, 60, 230)
    DrawRectangleLines(150, 100, 724, 500, 150, 150, 170, 255)
    
    // Title
    DrawText("INVENTORY", 450, 120, 30, 255, 255, 255, 255)
    
    // Draw inventory grid
    VAR slotSize = 60
    VAR slotSpacing = 5
    VAR gridStartX = 200
    VAR gridStartY = 180
    VAR slotsPerRow = 10
    
    FOR i = 0 TO inventorySize - 1
        VAR row = Int(i / slotsPerRow)
        VAR col = i % slotsPerRow
        VAR slotX = gridStartX + (col * (slotSize + slotSpacing))
        VAR slotY = gridStartY + (row * (slotSize + slotSpacing))
        
        // Draw slot
        IF i = selectedSlot THEN
            DrawRectangle(slotX, slotY, slotSize, slotSize, 80, 100, 120, 255)
            DrawRectangleLines(slotX, slotY, slotSize, slotSize, 255, 215, 0, 255)
        ELSE
            DrawRectangle(slotX, slotY, slotSize, slotSize, 60, 60, 80, 255)
            DrawRectangleLines(slotX, slotY, slotSize, slotSize, 100, 100, 120, 255)
        ENDIF
        
        // Draw item
        IF inventoryNames[i] <> "" THEN
            // Item icon (colored rectangle based on type)
            SELECT CASE inventoryTypes[i]
                CASE "weapon": DrawRectangle(slotX + 10, slotY + 10, 40, 40, 200, 100, 100, 255)
                CASE "armor": DrawRectangle(slotX + 10, slotY + 10, 40, 40, 100, 100, 200, 255)
                CASE "consumable": DrawRectangle(slotX + 10, slotY + 10, 40, 40, 100, 200, 100, 255)
                CASE ELSE: DrawRectangle(slotX + 10, slotY + 10, 40, 40, 150, 150, 150, 255)
            END SELECT
            
            // Quantity
            IF inventoryQuantities[i] > 1 THEN
                DrawText(STR(inventoryQuantities[i]), slotX + 35, slotY + 35, 12, 255, 255, 255, 255)
            ENDIF
        ENDIF
        
        // Slot number
        DrawText(STR(i + 1), slotX + 2, slotY + 2, 10, 150, 150, 150, 255)
    NEXT i
    
    // Item details
    IF inventoryNames[selectedSlot] <> "" THEN
        DrawText("Item: " + inventoryNames[selectedSlot], 200, 450, 18, 255, 255, 255, 255)
        DrawText("Type: " + inventoryTypes[selectedSlot], 200, 475, 16, 200, 200, 200, 255)
        DrawText("Quantity: " + STR(inventoryQuantities[selectedSlot]), 200, 500, 16, 200, 200, 200, 255)
        
        IF inventoryTypes[selectedSlot] = "consumable" THEN
            DrawText("Press ENTER to use", 200, 530, 14, 255, 215, 0, 255)
        ENDIF
    ENDIF
    
    // Instructions
    DrawText("Arrow keys: Navigate | ENTER: Use | I: Close", 200, 570, 14, 200, 200, 200, 255)
END FUNCTION

WHILE NOT WindowShouldClose()
    VAR dt = GetFrameTime()
    
    // Update message timer
    IF messageTimer > 0 THEN
        messageTimer = messageTimer - dt
    ENDIF
    
    SELECT CASE gameState
        CASE "menu"
            ClearBackground(20, 20, 40, 255)
            
            DrawText("RPG INVENTORY SYSTEM", 300, 150, 40, 255, 255, 255, 255)
            DrawText("Press SPACE to start", 380, 300, 20, 255, 255, 255, 255)
            DrawText("Press I to open inventory", 360, 350, 20, 255, 255, 255, 255)
            
            IF IsKeyPressed(KEY_SPACE) THEN
                gameState = "playing"
            ENDIF
            IF IsKeyPressed(KEY_I) THEN
                gameState = "inventory"
            ENDIF
            
        CASE "playing"
            ClearBackground(30, 50, 30, 255)
            
            // Draw game world (simple representation)
            DrawCircle(512, 384, 50, 100, 200, 100, 255)
            DrawText("GAME WORLD", 450, 500, 20, 255, 255, 255, 255)
            
            // Draw HUD
            DrawText("Health: " + STR(playerHealth) + "/" + STR(playerMaxHealth), 20, 20, 18, 255, 100, 100, 255)
            DrawText("Mana: " + STR(playerMana) + "/" + STR(playerMaxMana), 20, 45, 18, 100, 100, 255, 255)
            DrawText("Level: " + STR(playerLevel), 20, 70, 18, 255, 255, 255, 255)
            DrawText("Gold: " + STR(playerGold), 20, 95, 18, 255, 215, 0, 255)
            
            DrawText("Press I for inventory | ESC for menu", 350, 700, 16, 200, 200, 200, 255)
            
            // Show message
            IF messageTimer > 0 THEN
                DrawText(messageText, 300, 250, 20, 255, 215, 0, 255)
            ENDIF
            
            IF IsKeyPressed(KEY_I) THEN
                gameState = "inventory"
            ENDIF
            IF IsKeyPressed(KEY_ESCAPE) THEN
                gameState = "menu"
            ENDIF
            
        CASE "inventory"
            ClearBackground(20, 20, 40, 255)
            
            // Handle input
            IF IsKeyPressed(KEY_ESCAPE) OR IsKeyPressed(KEY_I) THEN
                gameState = "playing"
            ENDIF
            
            IF IsKeyPressed(KEY_LEFT) THEN
                selectedSlot = (selectedSlot - 1 + inventorySize) % inventorySize
            ENDIF
            IF IsKeyPressed(KEY_RIGHT) THEN
                selectedSlot = (selectedSlot + 1) % inventorySize
            ENDIF
            IF IsKeyPressed(KEY_UP) THEN
                selectedSlot = (selectedSlot - 10 + inventorySize) % inventorySize
            ENDIF
            IF IsKeyPressed(KEY_DOWN) THEN
                selectedSlot = (selectedSlot + 10) % inventorySize
            ENDIF
            
            IF IsKeyPressed(KEY_ENTER) THEN
                UseItem(selectedSlot)
            ENDIF
            
            // Add test items with number keys
            IF IsKeyPressed(KEY_1) THEN AddItem("Health Potion", "consumable", 3)
            IF IsKeyPressed(KEY_2) THEN AddItem("Mana Potion", "consumable", 2)
            IF IsKeyPressed(KEY_3) THEN AddItem("Iron Sword", "weapon", 1)
            IF IsKeyPressed(KEY_4) THEN AddItem("Leather Armor", "armor", 1)
            
            DrawInventory()
    END SELECT
    
WEND

CloseWindow()
```

---

## Conclusion

Congratulations! You've now learned:

- **Basic GUI elements** - buttons, sliders, checkboxes
- **Custom GUI drawing** for full control over appearance
- **Complete menu systems** with navigation
- **HUD elements** for game information display
- **Advanced UI components** like tab systems
- **Complete RPG interface** with inventory management

### Next Steps

1. **Create custom themes**: Design your own visual style
2. **Add animations**: Smooth transitions and hover effects
3. **Implement localization**: Support multiple languages
4. **Create UI editors**: Build tools to design interfaces visually
5. **Optimize performance**: Efficient rendering for complex UIs

### Common GUI Patterns

- **RPG Interfaces**: Inventory, stats, skill trees
- **Strategy Games**: Resource panels, unit controls, minimaps
- **Action Games**: Health bars, ammo counters, objective markers
- **Puzzle Games**: Level select, settings, achievement screens
- **Simulation Games**: Dashboard panels, control interfaces, data displays

GUI development is crucial for player experience. A well-designed interface can make the difference between a good game and a great one. Keep experimenting and refining your UI skills!

**Happy coding!**
