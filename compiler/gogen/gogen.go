// Package gogen generates Go source code from CyberBasic AST.
// The generated code calls raylib (and other libraries) directly.
package gogen

import (
	"fmt"
	"strings"

	"cyberbasic/compiler/parser"
)

// useLibsResult holds which libraries are used in the program.
type useLibsResult struct {
	useRL, useBullet, useBox2D, useNet, useECS bool
}

// usesLibs scans the program and returns which libraries are used.
func usesLibs(program *parser.Program) useLibsResult {
	var r useLibsResult
	for _, stmt := range program.Statements {
		if node, ok := stmt.(*parser.Call); ok && strings.Contains(node.Name, ".") {
			lib := strings.ToLower(strings.SplitN(node.Name, ".", 2)[0])
			switch lib {
			case "rl":
				r.useRL = true
			case "bullet":
				r.useBullet = true
			case "box2d":
				r.useBox2D = true
			case "net":
				r.useNet = true
			case "ecs":
				r.useECS = true
			}
		}
		if _, ok := stmt.(*parser.GameCommand); ok {
			r.useRL = true
		}
	}
	return r
}

// Generate produces Go source code from a parsed BASIC program.
// Output is package main with imports for raylib and/or bullet, box2d, etc. as used.
func Generate(program *parser.Program) (string, error) {
	used := usesLibs(program)
	var b strings.Builder
	b.WriteString("//go:build ignore\n\npackage main\n\n")
	b.WriteString("import (\n")
	b.WriteString("\t\"fmt\"\n")
	if used.useBullet {
		b.WriteString("\t\"cyberbasic/compiler/bindings/bullet\"\n")
	}
	if used.useBox2D {
		b.WriteString("\t\"cyberbasic/compiler/bindings/box2d\"\n")
	}
	if used.useNet {
		b.WriteString("\t\"cyberbasic/compiler/bindings/net\"\n")
	}
	if used.useECS {
		b.WriteString("\t\"cyberbasic/compiler/bindings/ecs\"\n")
	}
	if used.useRL {
		b.WriteString("\trl \"github.com/gen2brain/raylib-go/raylib\"\n")
	}
	b.WriteString(")\n\n")
	b.WriteString("func main() {\n")

	// Track sprite positions for DRAWSPRITE (simplified: draw as rectangles)
	spritePos := make(map[string]string) // id -> "x, y" or "x, y, w, h"

	for _, stmt := range program.Statements {
		code, err := emitStatement(stmt, &spritePos, "")
		if err != nil {
			return "", err
		}
		if code != "" {
			b.WriteString("\t" + code + "\n")
		}
	}

	b.WriteString("}\n")
	return b.String(), nil
}

func emitStatement(stmt parser.Node, spritePos *map[string]string, indent string) (string, error) {
	switch node := stmt.(type) {
	case *parser.DimStatement:
		return emitDim(node)
	case *parser.Assignment:
		return emitAssignment(node)
	case *parser.Call:
		return emitCall(node)
	case *parser.GameCommand:
		return emitGameCommand(node, spritePos)
	case *parser.IfStatement:
		return emitIf(node, spritePos, indent)
	case *parser.ForStatement:
		return emitFor(node, spritePos, indent)
	case *parser.WhileStatement:
		return emitWhile(node, spritePos, indent)
	case *parser.ReturnStatement:
		if node.Value != nil {
			expr, _ := emitExpr(node.Value)
			return "return " + expr, nil
		}
		return "return", nil
	default:
		return "", nil
	}
}

func emitDim(d *parser.DimStatement) (string, error) {
	var parts []string
	for _, v := range d.Variables {
		goType := basicTypeToGo(v.Type)
		parts = append(parts, fmt.Sprintf("var %s %s", v.Name, goType))
	}
	return strings.Join(parts, "\n\t"), nil
}

func basicTypeToGo(t string) string {
	switch strings.ToLower(t) {
	case "integer", "int":
		return "int"
	case "float", "single", "double":
		return "float64"
	case "string", "str":
		return "string"
	case "boolean", "bool":
		return "bool"
	default:
		return "interface{}"
	}
}

func emitAssignment(a *parser.Assignment) (string, error) {
	rhs, err := emitExpr(a.Value)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s = %s", a.Variable, rhs), nil
}

// box2dGoName maps BASIC BOX2D.* method (lowercase) to flat Go export (no namespace in generated code).
var box2dGoName = map[string]string{
	"createworld": "CreateWorld2D", "destroyworld": "DestroyWorld2D", "step": "Step2D",
	"createbody": "CreateBody2D", "destroybody": "DestroyBody2D",
	"getbodycount": "GetBodyCount2D", "getbodyid": "GetBodyId2D", "createbodyatscreen": "CreateBodyAtScreen2D",
	"getpositionx": "GetPositionX2D", "getpositiony": "GetPositionY2D", "getangle": "GetAngle2D",
	"setlinearvelocity": "SetVelocity2D", "settransform": "SetTransform2D", "applyforce": "ApplyForce2D",
}

// bulletGoName maps BASIC BULLET.* method (lowercase) to flat Go export (no namespace in generated code).
var bulletGoName = map[string]string{
	"createworld": "CreateWorld3D", "destroyworld": "DestroyWorld3D", "setgravity": "SetWorldGravity3D", "step": "Step3D",
	"createbox": "CreateBox3D", "createsphere": "CreateSphere3D", "destroybody": "DestroyBody3D",
	"setposition": "SetPosition3D", "getpositionx": "GetPositionX3D", "getpositiony": "GetPositionY3D", "getpositionz": "GetPositionZ3D",
	"setvelocity": "SetVelocity3D", "getvelocityx": "GetVelocityX3D", "getvelocityy": "GetVelocityY3D", "getvelocityz": "GetVelocityZ3D",
	"getrotationx": "GetYaw3D", "getrotationy": "GetPitch3D", "getrotationz": "GetRoll3D", "setrotation": "SetRotation3D",
	"applyforce": "ApplyForce3D", "applycentralforce": "ApplyForce3D", "applyimpulse": "ApplyImpulse3D",
	"raycast": "RayCastFromDir3D",
	"getraycasthitx": "RayHitX3D", "getraycasthity": "RayHitY3D", "getraycasthitz": "RayHitZ3D",
	"getraycasthitbody": "RayHitBody3D",
	"getraycasthitnormalx": "RayHitNormalX3D", "getraycasthitnormaly": "RayHitNormalY3D", "getraycasthitnormalz": "RayHitNormalZ3D",
}

// rlGoName maps BASIC RL.* (canonical lowercase) to raylib-go export (PascalCase).
var rlGoName = map[string]string{
	"initwindow": "InitWindow", "settargetfps": "SetTargetFPS", "windowshouldclose": "WindowShouldClose",
	"closewindow": "CloseWindow", "setwindowposition": "SetWindowPosition",
	"begindrawing": "BeginDrawing", "enddrawing": "EndDrawing", "clearbackground": "ClearBackground",
	"drawrectangle": "DrawRectangle", "drawcircle": "DrawCircle", "drawtext": "DrawText",
	"drawline": "DrawLine", "getscreenwidth": "GetScreenWidth", "getscreenheight": "GetScreenHeight",
	"iskeypressed": "IsKeyPressed", "iskeydown": "IsKeyDown", "white": "White", "black": "Black",
}

func emitCall(c *parser.Call) (string, error) {
	// Foreign API: RL.InitWindow(...) -> rl.InitWindow(...), BULLET.* -> bullet.*
	if strings.Contains(c.Name, ".") {
		parts := strings.SplitN(c.Name, ".", 2)
		lib, fn := parts[0], parts[1]
		if strings.EqualFold(lib, "rl") {
			args, err := emitExprList(c.Arguments)
			if err != nil {
				return "", err
			}
			return fmt.Sprintf("rl.%s(%s)", fn, strings.Join(args, ", ")), nil
		}
		if strings.EqualFold(lib, "bullet") {
			args, err := emitExprList(c.Arguments)
			if err != nil {
				return "", err
			}
			goName := bulletGoName[strings.ToLower(fn)]
			if goName == "" {
				goName = fn
			}
			return fmt.Sprintf("bullet.%s(%s)", goName, strings.Join(args, ", ")), nil
		}
		if strings.EqualFold(lib, "box2d") {
			args, err := emitExprList(c.Arguments)
			if err != nil {
				return "", err
			}
			goName := box2dGoName[strings.ToLower(fn)]
			if goName == "" {
				goName = fn
			}
			return fmt.Sprintf("box2d.%s(%s)", goName, strings.Join(args, ", ")), nil
		}
		if strings.EqualFold(lib, "audio") {
			args, err := emitExprList(c.Arguments)
			if err != nil {
				return "", err
			}
			return fmt.Sprintf("audio.%s(%s)", fn, strings.Join(args, ", ")), nil
		}
		if strings.EqualFold(lib, "ui") {
			args, err := emitExprList(c.Arguments)
			if err != nil {
				return "", err
			}
			return fmt.Sprintf("ui.%s(%s)", fn, strings.Join(args, ", ")), nil
		}
		if strings.EqualFold(lib, "tile") {
			args, err := emitExprList(c.Arguments)
			if err != nil {
				return "", err
			}
			return fmt.Sprintf("tile.%s(%s)", fn, strings.Join(args, ", ")), nil
		}
		if strings.EqualFold(lib, "net") {
			args, err := emitExprList(c.Arguments)
			if err != nil {
				return "", err
			}
			return fmt.Sprintf("net.%s(%s)", fn, strings.Join(args, ", ")), nil
		}
		if strings.EqualFold(lib, "assimp") {
			args, err := emitExprList(c.Arguments)
			if err != nil {
				return "", err
			}
			return fmt.Sprintf("assimp.%s(%s)", fn, strings.Join(args, ", ")), nil
		}
		if strings.EqualFold(lib, "path") {
			args, err := emitExprList(c.Arguments)
			if err != nil {
				return "", err
			}
			return fmt.Sprintf("path.%s(%s)", fn, strings.Join(args, ", ")), nil
		}
		if strings.EqualFold(lib, "ecs") {
			args, err := emitExprList(c.Arguments)
			if err != nil {
				return "", err
			}
			return fmt.Sprintf("ecs.%s(%s)", fn, strings.Join(args, ", ")), nil
		}
		return "", fmt.Errorf("unknown library: %s", lib)
	}
	// Built-in: Print -> fmt.Println (case-insensitive)
	if strings.EqualFold(c.Name, "print") {
		if len(c.Arguments) == 0 {
			return "fmt.Println()", nil
		}
		args, err := emitExprList(c.Arguments)
		if err != nil {
			return "", err
		}
		return "fmt.Println(" + strings.Join(args, ", ") + ")", nil
	}
	// Multi-window API: not supported in --gen-go (VM-only binding)
	if isMultiWindowCall(c.Name) {
		return "", fmt.Errorf("multi-window API %s is not supported in --gen-go; run without --gen-go", c.Name)
	}
	return "", fmt.Errorf("unsupported call: %s", c.Name)
}

// isMultiWindowCall returns true for Window*, Channel*, State*, Dock*, OnWindow* APIs.
func isMultiWindowCall(name string) bool {
	n := strings.ToLower(name)
	return strings.HasPrefix(n, "window") || strings.HasPrefix(n, "channel") ||
		strings.HasPrefix(n, "state") || strings.HasPrefix(n, "dock") ||
		strings.HasPrefix(n, "onwindow")
}

func emitGameCommand(g *parser.GameCommand, spritePos *map[string]string) (string, error) {
	cmd := strings.ToLower(g.Command)
	args, err := emitExprList(g.Arguments)
	if err != nil {
		return "", err
	}

	switch cmd {
	case "initgraphics3d", "initgraphics":
		if len(args) < 3 {
			return "", fmt.Errorf("%s requires (width, height, title)", cmd)
		}
		return fmt.Sprintf("rl.InitWindow(%s, %s, %s)\n\trl.SetTargetFPS(60)\n\trl.SetWindowPosition(120, 80)", args[0], args[1], args[2]), nil
	case "createsprite":
		if len(args) >= 3 {
			(*spritePos)[args[0]] = args[1] + ", " + args[2] + ", 64, 64"
		}
		return "", nil
	case "setspriteposition":
		if len(args) >= 3 {
			(*spritePos)[args[0]] = args[1] + ", " + args[2] + ", 64, 64"
		}
		return "", nil
	case "drawsprite":
		if len(args) < 1 {
			return "", fmt.Errorf("DRAWSPRITE requires sprite id")
		}
		id := strings.Trim(args[0], "\"")
		pos, ok := (*spritePos)[args[0]]
		if !ok {
			pos = "0, 0, 64, 64"
		}
		_ = id
		return fmt.Sprintf("rl.DrawRectangle(%s, rl.White)", pos), nil
	case "sync":
		// Emit only EndDrawing; caller ensures BeginDrawing/ClearBackground at loop start
		return "rl.EndDrawing()", nil
	case "shouldclose":
		return "", fmt.Errorf("SHOULDCLOSE is used in condition only")
	default:
		return "", fmt.Errorf("unsupported game command in Go gen: %s", cmd)
	}
}

func emitIf(i *parser.IfStatement, spritePos *map[string]string, indent string) (string, error) {
	cond, err := emitExpr(i.Condition)
	if err != nil {
		return "", err
	}
	// Go requires bool; BASIC uses int (0=false). Add " != 0" unless condition is clearly boolean.
	if !strings.Contains(cond, "WindowShouldClose") && !strings.Contains(cond, "true") && !strings.Contains(cond, "false") && !strings.Contains(cond, "==") && !strings.Contains(cond, "!=") {
		cond = cond + " != 0"
	}
	tab := indent + "\t"
	var b strings.Builder
	b.WriteString(fmt.Sprintf("if %s {\n", cond))
	for _, s := range i.ThenBlock.Statements {
		line, err := emitStatement(s, spritePos, tab)
		if err != nil {
			return "", err
		}
		if line != "" {
			b.WriteString(tab + line + "\n")
		}
	}
	if i.ElseBlock != nil {
		b.WriteString(indent + "} else {\n")
		for _, s := range i.ElseBlock.Statements {
			line, err := emitStatement(s, spritePos, tab)
			if err != nil {
				return "", err
			}
			if line != "" {
				b.WriteString(tab + line + "\n")
			}
		}
	}
	b.WriteString(indent + "}")
	return b.String(), nil
}

func emitFor(f *parser.ForStatement, spritePos *map[string]string, indent string) (string, error) {
	start, _ := emitExpr(f.Start)
	end, _ := emitExpr(f.End)
	step := "1"
	if f.Step != nil {
		step, _ = emitExpr(f.Step)
	}
	tab := indent + "\t"
	var b strings.Builder
	b.WriteString(fmt.Sprintf("for %s := %s; %s <= %s; %s += %s {\n", f.Variable, start, f.Variable, end, f.Variable, step))
	for _, s := range f.Body.Statements {
		line, err := emitStatement(s, spritePos, tab)
		if err != nil {
			return "", err
		}
		if line != "" {
			b.WriteString(tab + line + "\n")
		}
	}
	b.WriteString(indent + "}")
	return b.String(), nil
}

func blockContainsDrawOrSync(statements []parser.Node) bool {
	for _, s := range statements {
		if g, ok := s.(*parser.GameCommand); ok && strings.EqualFold(g.Command, "sync") {
			return true
		}
		// Any draw-like call (ClearBackground, DrawCircle, etc.) - treat as needing frame
		if st, ok := s.(*parser.Statement); ok && st.Value != nil {
			if c, ok := st.Value.(*parser.Call); ok {
				nm := strings.ToLower(c.Name)
				if strings.HasSuffix(nm, "background") || strings.HasPrefix(nm, "draw") {
					return true
				}
			}
		}
	}
	return false
}

func emitWhile(w *parser.WhileStatement, spritePos *map[string]string, indent string) (string, error) {
	cond, err := emitExpr(w.Condition)
	if err != nil {
		return "", err
	}
	// BASIC: WHILE NOT ShouldClose() -> Go: for !rl.WindowShouldClose()
	if cond == "!rl.WindowShouldClose()" || strings.Contains(cond, "WindowShouldClose") {
		cond = "!rl.WindowShouldClose()"
	}
	tab := indent + "\t"
	var b strings.Builder
	b.WriteString(fmt.Sprintf("for %s {\n", cond))
	injectFrame := blockContainsDrawOrSync(w.Body.Statements)
	if injectFrame {
		b.WriteString(tab + "rl.BeginDrawing()\n")
		b.WriteString(tab + "rl.ClearBackground(rl.NewColor(20, 20, 30, 255))\n")
	}
	for _, s := range w.Body.Statements {
		line, err := emitStatement(s, spritePos, tab)
		if err != nil {
			return "", err
		}
		if line != "" {
			b.WriteString(tab + line + "\n")
		}
	}
	if injectFrame {
		b.WriteString(tab + "rl.EndDrawing()\n")
	}
	b.WriteString(indent + "}")
	return b.String(), nil
}

func emitExprList(nodes []parser.Node) ([]string, error) {
	var out []string
	for _, n := range nodes {
		s, err := emitExpr(n)
		if err != nil {
			return nil, err
		}
		out = append(out, s)
	}
	return out, nil
}

func emitExpr(expr parser.Node) (string, error) {
	switch node := expr.(type) {
	case *parser.Number:
		return node.Value, nil
	case *parser.StringLiteral:
		return `"` + strings.Trim(node.Value, `"`) + `"`, nil
	case *parser.Boolean:
		if node.Value {
			return "true", nil
		}
		return "false", nil
	case *parser.Identifier:
		return node.Name, nil
	case *parser.Call:
		if strings.EqualFold(node.Name, "shouldclose") {
			return "rl.WindowShouldClose()", nil
		}
		if strings.Contains(node.Name, ".") {
			parts := strings.SplitN(node.Name, ".", 2)
			if strings.EqualFold(parts[0], "rl") {
				args, _ := emitExprList(node.Arguments)
				return "rl." + parts[1] + "(" + strings.Join(args, ", ") + ")", nil
			}
			if strings.EqualFold(parts[0], "bullet") {
				args, _ := emitExprList(node.Arguments)
				goName := bulletGoName[strings.ToLower(parts[1])]
				if goName == "" {
					goName = parts[1]
				}
				return "bullet." + goName + "(" + strings.Join(args, ", ") + ")", nil
			}
			if strings.EqualFold(parts[0], "box2d") {
				args, _ := emitExprList(node.Arguments)
				goName := box2dGoName[strings.ToLower(parts[1])]
				if goName == "" {
					goName = parts[1]
				}
				return "box2d." + goName + "(" + strings.Join(args, ", ") + ")", nil
			}
		}
		args, _ := emitExprList(node.Arguments)
		return node.Name + "(" + strings.Join(args, ", ") + ")", nil
	case *parser.BinaryOp:
		left, _ := emitExpr(node.Left)
		right, _ := emitExpr(node.Right)
		op := node.Operator
		if strings.EqualFold(op, "and") {
			op = "&&"
		}
		if strings.EqualFold(op, "or") {
			op = "||"
		}
		return fmt.Sprintf("(%s %s %s)", left, op, right), nil
	case *parser.UnaryOp:
		operand, _ := emitExpr(node.Operand)
		if strings.EqualFold(node.Operator, "NOT") {
			return "!" + operand, nil
		}
		return node.Operator + operand, nil
	default:
		return "", fmt.Errorf("unsupported expression: %T", expr)
	}
}
