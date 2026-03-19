package codegen

import (
	"cyberbasic/compiler/parser"
	"fmt"
	"strings"
)

// physicsNamespaceToFlat maps BOX2D.* and BULLET.* dotted names (lowercase) to flat VM names.
func physicsNamespaceToFlat(name string) string {
	name = strings.ToLower(name)
	if flat, ok := physicsNamespaceFlatMap[name]; ok {
		return flat
	}
	return ""
}

var physicsNamespaceFlatMap = map[string]string{
	"box2d.createworld": "createworld2d", "box2d.destroyworld": "destroyworld2d", "box2d.step": "step2d",
	"box2d.createbody": "createbody2d", "box2d.destroybody": "destroybody2d", "box2d.getbodycount": "getbodycount2d",
	"box2d.getbodyid": "getbodyid2d", "box2d.createbodyatscreen": "createbodyatscreen2d",
	"box2d.getpositionx": "getpositionx2d", "box2d.getpositiony": "getpositiony2d", "box2d.getangle": "getangle2d",
	"box2d.setlinearvelocity": "setvelocity2d", "box2d.settransform": "settransform2d", "box2d.applyforce": "applyforce2d",
	"bullet.createworld": "createworld3d", "bullet.destroyworld": "destroyworld3d", "bullet.setgravity": "setworldgravity3d",
	"bullet.step": "step3d", "bullet.createbox": "createbox3d", "bullet.createsphere": "createsphere3d",
	"bullet.destroybody": "destroybody3d", "bullet.setposition": "setposition3d",
	"bullet.getpositionx": "getpositionx3d", "bullet.getpositiony": "getpositiony3d", "bullet.getpositionz": "getpositionz3d",
	"bullet.setvelocity": "setvelocity3d", "bullet.getvelocityx": "getvelocityx3d", "bullet.getvelocityy": "getvelocityy3d", "bullet.getvelocityz": "getvelocityz3d",
	"bullet.getrotationx": "getyaw3d", "bullet.getrotationy": "getpitch3d", "bullet.getrotationz": "getroll3d",
	"bullet.setrotation": "setrotation3d", "bullet.applyforce": "applyforce3d", "bullet.applycentralforce": "applyforce3d",
	"bullet.applyimpulse": "applyimpulse3d", "bullet.raycast": "raycastfromdir3d",
	"bullet.getraycasthitx": "rayhitx3d", "bullet.getraycasthity": "rayhity3d", "bullet.getraycasthitz": "rayhitz3d",
	"bullet.getraycasthitbody": "rayhitbody3d", "bullet.getraycasthitnormalx": "rayhitnormalx3d",
	"bullet.getraycasthitnormaly": "rayhitnormaly3d", "bullet.getraycasthitnormalz": "rayhitnormalz3d",
}

// MaxConstIndex is the maximum constant index (0-255) that fits in a single byte in the chunk.
const MaxConstIndex = 255

func checkConstIndex(idx int, kind string) error {
	if idx > MaxConstIndex {
		return fmt.Errorf("%s", "too many constants"+kind)
	}
	return nil
}

func parseInt(s string) (int, error) {
	var result int
	_, err := fmt.Sscanf(s, "%d", &result)
	return result, err
}

func parseFloat(s string) (float64, error) {
	var result float64
	_, err := fmt.Sscanf(s, "%f", &result)
	return result, err
}

func unwrapStatement(n parser.Node) parser.Node {
	if s, ok := n.(*parser.Statement); ok && s.Value != nil {
		return s.Value
	}
	return n
}

// WalkStatements visits every statement in nodes and recursively enters If/For/While/Repeat/Select blocks.
func WalkStatements(nodes []parser.Node, pred func(parser.Node) bool) bool {
	for _, n := range nodes {
		n2 := unwrapStatement(n)
		if pred(n2) {
			return true
		}
		switch v := n2.(type) {
		case *parser.IfStatement:
			if v.ThenBlock != nil && WalkStatements(v.ThenBlock.Statements, pred) {
				return true
			}
			for _, b := range v.ElseIfs {
				if b.Block != nil && WalkStatements(b.Block.Statements, pred) {
					return true
				}
			}
			if v.ElseBlock != nil && WalkStatements(v.ElseBlock.Statements, pred) {
				return true
			}
		case *parser.ForStatement:
			if v.Body != nil && WalkStatements(v.Body.Statements, pred) {
				return true
			}
		case *parser.WhileStatement:
			if v.Body != nil && WalkStatements(v.Body.Statements, pred) {
				return true
			}
		case *parser.RepeatStatement:
			if v.Body != nil && WalkStatements(v.Body.Statements, pred) {
				return true
			}
		case *parser.SelectCaseStatement:
			for _, c := range v.Cases {
				if c.Block != nil && WalkStatements(c.Block.Statements, pred) {
					return true
				}
			}
			if v.ElseBlock != nil && WalkStatements(v.ElseBlock.Statements, pred) {
				return true
			}
		}
	}
	return false
}

var threeDDrawNames = map[string]bool{
	"drawcube": true, "drawcubewires": true, "drawsphere": true, "drawspherewires": true,
	"drawmodel": true, "drawmodelsimple": true, "drawmodelex": true, "drawmodelwires": true, "drawplane": true,
	"drawline3d": true, "drawpoint3d": true, "drawcircle3d": true, "drawgrid": true,
	"drawcylinder": true, "drawcylinderwires": true, "drawray": true, "drawtriangle3d": true,
	"drawobject": true, "beginmode3d": true,
}

func predicate3DDraw(n parser.Node) bool {
	call, ok := n.(*parser.Call)
	if !ok {
		return false
	}
	name := normWindowShouldCloseName(call.Name)
	return threeDDrawNames[name]
}

func normWindowShouldCloseName(name string) string {
	s := strings.ToLower(name)
	if strings.HasPrefix(s, "rl.") {
		s = s[3:]
	}
	return s
}

func isGameLoopCondition(condition parser.Node, forRepeat bool) bool {
	if forRepeat {
		call, ok := condition.(*parser.Call)
		return ok && normWindowShouldCloseName(call.Name) == "windowshouldclose" && len(call.Arguments) == 0
	}
	un, ok := condition.(*parser.UnaryOp)
	if !ok || strings.ToLower(un.Operator) != "not" {
		return false
	}
	call, ok := un.Operand.(*parser.Call)
	return ok && normWindowShouldCloseName(call.Name) == "windowshouldclose" && len(call.Arguments) == 0
}

var frameBoundaryNames = map[string]bool{
	"begindrawing": true, "enddrawing": true, "beginframe": true, "endframe": true,
}

func bodyContainsFrameBoundaries(nodes []parser.Node) bool {
	return WalkStatements(nodes, func(n parser.Node) bool {
		call, ok := n.(*parser.Call)
		if !ok {
			return false
		}
		name := strings.ToLower(call.Name)
		if strings.HasPrefix(name, "rl.") {
			name = name[3:]
		}
		return frameBoundaryNames[name]
	})
}

func bodyContainsSync(nodes []parser.Node) bool {
	return WalkStatements(nodes, func(n parser.Node) bool {
		gc, ok := n.(*parser.GameCommand)
		return ok && strings.ToLower(gc.Command) == "sync"
	})
}

func bodyContains3DDraw(nodes []parser.Node) bool {
	return WalkStatements(nodes, predicate3DDraw)
}

func getSourceLine(node parser.Node) int {
	if loc, ok := node.(parser.HasSourceLoc); ok {
		return loc.GetLine()
	}
	return 0
}

func errWithLine(node parser.Node, err error) error {
	if err == nil {
		return nil
	}
	if line := getSourceLine(node); line > 0 {
		return fmt.Errorf("line %d: %w", line, err)
	}
	return err
}

// nearestName returns the candidate closest to bad by edit distance (Levenshtein), or "" if none within maxDist.
// Used for "did you mean?" suggestions for unknown subs, members, etc.
func nearestName(bad string, candidates []string, maxDist int) string {
	bad = strings.ToLower(bad)
	best := ""
	bestDist := maxDist + 1
	for _, c := range candidates {
		d := levenshtein(bad, strings.ToLower(c))
		if d < bestDist && d > 0 {
			bestDist = d
			best = c
		}
	}
	return best
}

func levenshtein(a, b string) int {
	if len(a) == 0 {
		return len(b)
	}
	if len(b) == 0 {
		return len(a)
	}
	prev := make([]int, len(b)+1)
	curr := make([]int, len(b)+1)
	for j := 0; j <= len(b); j++ {
		prev[j] = j
	}
	for i := 1; i <= len(a); i++ {
		curr[0] = i
		for j := 1; j <= len(b); j++ {
			del := prev[j] + 1
			ins := curr[j-1] + 1
			sub := prev[j-1]
			if a[i-1] != b[j-1] {
				sub++
			}
			curr[j] = min3(del, ins, sub)
		}
		prev, curr = curr, prev
	}
	return prev[len(b)]
}

func min3(a, b, c int) int {
	if a < b {
		if a < c {
			return a
		}
		return c
	}
	if b < c {
		return b
	}
	return c
}
