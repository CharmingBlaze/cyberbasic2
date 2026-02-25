# GUI Guide

CyberBasic provides two UI options:

1. **Full raygui** (recommended when using CGO): the real [raygui](https://github.com/raysan5/raygui) library via **Gui*** functions. You place controls by (x, y, w, h). Requires CGO (gen2brain/raylib-go/raygui).
2. **Pure-Go layout UI**: **BeginUI()** / **EndUI()** with a vertical cursor; widgets (Label, Button, Slider, etc.) advance the layout. No CGO; use when building with CGO_ENABLED=0.

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

## Widgets

| Widget | Call | Returns |
|--------|------|---------|
| **Label** | Label(text) | — |
| **Button** | Button(text) | true when clicked |
| **Slider** | Slider(text, value, min, max) | current value |
| **Checkbox** | Checkbox(text, checked) | 1 or 0 |
| **TextBox** | TextBox(id, text) | current text (edit with same id each frame) |
| **Dropdown** | Dropdown(id, itemsText, activeIndex) | new activeIndex (itemsText = "A;B;C") |
| **ProgressBar** | ProgressBar(text, value, min, max) | value (display only) |
| **WindowBox** | WindowBox(title) … EndWindowBox() | — |
| **GroupBox** | GroupBox(text) … EndGroupBox() | — |

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

See [API_REFERENCE.md](../API_REFERENCE.md) for the full list.
