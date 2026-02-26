# GUI Guide

CyberBasic provides two UI options:

1. **Full raygui** (recommended when using CGO): the real [raygui](https://github.com/raysan5/raygui) library via **Gui*** functions. You place controls by (x, y, w, h). Requires CGO (gen2brain/raylib-go/raygui).
2. **Pure-Go layout UI**: **BeginUI()** / **EndUI()** with a vertical cursor; widgets (Label, Button, Slider, etc.) advance the layout. No CGO; use when building with CGO_ENABLED=0.

## Table of Contents

1. [Full raygui (Gui* functions)](#full-raygui-gui-functions)
2. [Pure-Go layout UI (BeginUI / EndUI)](#pure-go-layout-ui-beginui--endui)
3. [Using GUI in the hybrid loop](#using-gui-in-the-hybrid-loop)
4. [Full widget reference](#full-widget-reference)
5. [Example: options menu](#example-options-menu)
6. [Notes](#notes)
7. [See also](#see-also)

---

## Full raygui (Gui* functions)

Use **Gui*** when building with CGO. Place controls by rectangle (x, y, width, height) in pixels:

```basic
VAR volume = 0.5
VAR enabled = 1
VAR name = "Player"

WHILE NOT WindowShouldClose()
  ClearBackground(40, 40, 50, 255)
  GuiWindowBox(24, 24, 280, 180, "Options")
  GuiLabel(32, 52, 80, 24, "Name:")
  name = GuiTextBox("tb1", 120, 48, 160, 28, name)
  volume = GuiSlider(32, 88, 240, 24, "0", "1", volume, 0, 1)
  enabled = GuiCheckBox(32, 120, 24, 24, "Enable sound", enabled)
  VAR clicked = GuiButton(32, 152, 120, 28, "OK")
  IF clicked THEN PRINT "OK pressed"
WEND
```

- **GuiTextBox(id, x, y, w, h, text)** returns the current text; use the same string **id** each frame so the buffer is preserved.
- **GuiDropdownBox(id, x, y, w, h, itemsText, active)** uses itemsText like `"Low;Medium;High"` and returns the new 0-based index.

---

## Pure-Go layout UI (BeginUI / EndUI)

Use inside your game loop, after **ClearBackground** and before other drawing:

```basic
WHILE NOT WindowShouldClose()
  ClearBackground(40, 40, 50, 255)
  BeginUI()
  Label("Hello")
  VAR clicked = Button("Click me")
  IF clicked THEN PRINT "Clicked!"
  VAR v = Slider("Volume", 0.5, 0, 1)
  VAR on = Checkbox("Enable", 1)
  EndUI()
  // ... rest of game draw
WEND
```

### Using GUI in the hybrid loop

When you define **update(dt)** and **draw()**, call **BeginUI**/ **EndUI** or **Gui*** inside your **draw()** so GUI commands are queued with other draw calls and rendered in order (2D, then 3D, then GUI). See [Program Structure](PROGRAM_STRUCTURE.md#hybrid-updatedraw-loop).

---

## Full widget reference

### Pure-Go layout (BeginUI / EndUI)

| Command | Arguments | Returns | Description |
|---------|-----------|---------|-------------|
| **BeginUI** | () | — | Start UI layout; resets cursor |
| **EndUI** | () | — | End UI layout |
| **Label** | (text) | — | Label widget |
| **Button** | (text) | 1 or 0 | 1 if clicked |
| **Slider** | (text, value, min, max) | float | Slider; returns current value |
| **Checkbox** | (text, checked) | 1 or 0 | Checkbox state |
| **TextBox** | (id, text) | string | Editable; use same id each frame |
| **Dropdown** | (id, itemsText, activeIndex) | int | itemsText = "A;B;C"; returns new index |
| **ProgressBar** | (text, value, min, max) | float | Progress bar (display) |
| **WindowBox** | (title) | — | Start window box |
| **EndWindowBox** | () | — | End window box |
| **GroupBox** | (text) | — | Start group box |
| **EndGroupBox** | () | — | End group box |

### Full raygui (Gui*; requires CGO)

All coordinates and sizes in pixels (x, y, width, height).

| Command | Arguments | Returns | Description |
|---------|-----------|---------|-------------|
| **GuiLabel** | (x, y, w, h, text) | — | Label |
| **GuiButton** | (x, y, w, h, text) | 1 or 0 | 1 if clicked |
| **GuiCheckBox** | (x, y, w, h, text, checked) | 1 or 0 | Checkbox |
| **GuiSlider** | (x, y, w, h, textLeft, textRight, value, min, max) | float | Slider |
| **GuiProgressBar** | (x, y, w, h, textLeft, textRight, value, min, max) | float | Progress bar |
| **GuiTextBox** | (id, x, y, w, h, text) | string | Editable (id = cache key) |
| **GuiDropdownBox** | (id, x, y, w, h, itemsText, active) | int | itemsText e.g. "One;Two;Three" |
| **GuiWindowBox** | (x, y, w, h, title) | 1 or 0 | 1 if close clicked |
| **GuiGroupBox** | (x, y, w, h, text) | — | Group box |
| **GuiLine** | (x, y, w, h, text) | — | Line |
| **GuiPanel** | (x, y, w, h, text) | — | Panel |

## Example: options menu

```basic
VAR volume = 0.7
VAR fullscreen = 0
VAR name = "Player"

WHILE NOT WindowShouldClose()
  ClearBackground(30, 30, 40, 255)
  BeginUI()
  WindowBox("Options")
  Label("Name:")
  name = TextBox("name", name)
  volume = Slider("Volume", volume, 0, 1)
  fullscreen = Checkbox("Fullscreen", fullscreen)
  VAR sel = Dropdown("quality", "Low;Medium;High", 1)
  EndWindowBox()
  IF Button("Apply") THEN
    // apply options...
  END IF
  EndUI()
WEND
```

## Notes

- **Full raygui:** Requires CGO (e.g. `go build` with CGO_ENABLED=1 and a C compiler). If you build with CGO_ENABLED=0 (e.g. purego raylib on Windows), use the pure-Go **BeginUI** / **EndUI** widgets instead.
- **TextBox / GuiTextBox:** Use a stable string id (e.g. "tb1") so the same buffer is used each frame. Click to focus; type to edit; click outside to unfocus.
- **Dropdown / GuiDropdownBox:** itemsText uses semicolons: "Item1;Item2;Item3". activeIndex is 0-based.
- **Slider / ProgressBar:** value, min, max are numbers; Slider returns the new value after dragging.

## See also

- [API Reference](../API_REFERENCE.md) (section 20) — full list of GUI and raygui functions
- [Command Reference](COMMAND_REFERENCE.md) — commands by feature
- [Program Structure](PROGRAM_STRUCTURE.md) — hybrid update/draw loop
- [Game Development Guide](GAME_DEVELOPMENT_GUIDE.md)
- [Documentation Index](DOCUMENTATION_INDEX.md)
