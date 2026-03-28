package codegen

import (
	"cyberbasic/compiler/parser"
	"cyberbasic/compiler/vm"
	"fmt"
	"strings"
)

// dotObjectRoots are identifiers that always use OpGetProp (DotObject), not vector .x/.y/.z helpers.
// Keys are lowercase to match globals registered by bindings (same convention as WINDOW -> "window").
var dotObjectRoots = map[string]bool{
	"window":     true,
	"physics":    true,
	"audio":      true,
	"input":      true,
	"assets":     true,
	"shader":     true,
	"ai":         true,
	"scenes":     true,
	"ecs":        true,
	"net":        true,
	"sql":        true,
	"nakama":     true,
	"terrain":    true,
	"water":      true,
	"vegetation": true,
	"world":      true,
	"navigation": true,
	"nav":        true,
	"indoor":     true,
	"procedural": true,
	"objects":    true,
	"object":     true,
	"effect":     true,
	"tween":      true,
	"camera":     true,
	"engine":     true,
	"std":        true,
	"draw":       true,
	"texture":    true,
	"sprite":     true,
	"file":       true,
	"http":       true,
	"model":      true,
	"shapes3d":   true,
	"mesh":       true,
	"image":      true,
	"font":       true,
	"rlaudio":    true,
	"box2d":      true,
	"bullet":     true,
	"game":       true,
}

func collectMemberAccessChain(ma *parser.MemberAccess) ([]string, parser.Node) {
	var segs []string
	cur := parser.Node(ma)
	for {
		m, ok := cur.(*parser.MemberAccess)
		if !ok {
			break
		}
		segs = append([]string{m.Member}, segs...)
		cur = m.Object
	}
	return segs, cur
}

func (e *Emitter) emitOpGetProp(path []string) error {
	if len(path) < 1 || len(path) > 32 {
		return fmt.Errorf("invalid property path length")
	}
	e.chunk.Write(byte(vm.OpGetProp))
	e.chunk.Write(byte(len(path)))
	for _, seg := range path {
		low := strings.ToLower(seg)
		ci := e.chunk.WriteConstant(low)
		if err := checkConstIndex(ci, " for OpGetProp path"); err != nil {
			return err
		}
		e.chunk.Write(byte(ci))
	}
	return nil
}

// compileDotMethodCall compiles base.path.method(args) for DotObject (not rl./box2d./…).
func (e *Emitter) compileDotMethodCall(call *parser.Call, parts []string) error {
	if len(parts) < 2 {
		return fmt.Errorf("invalid method call name")
	}
	method := strings.ToLower(parts[len(parts)-1])
	baseName := parts[0]
	path := parts[1 : len(parts)-1]

	ident := &parser.Identifier{Name: baseName, Line: call.Line, Col: call.Col}
	if err := e.compileIdentifier(ident); err != nil {
		return err
	}
	if len(path) > 0 {
		if err := e.emitOpGetProp(path); err != nil {
			return err
		}
	}
	for _, arg := range call.Arguments {
		if err := e.compileExpression(arg); err != nil {
			return err
		}
	}
	mi := e.chunk.WriteConstant(method)
	if err := checkConstIndex(mi, " for OpCallMethod"); err != nil {
		return err
	}
	e.chunk.Write(byte(vm.OpCallMethod))
	e.chunk.Write(byte(mi))
	e.chunk.Write(byte(len(call.Arguments)))
	return nil
}

func (e *Emitter) emitOpSetProp(path []string) error {
	if len(path) < 1 || len(path) > 32 {
		return fmt.Errorf("invalid property path length")
	}
	e.chunk.Write(byte(vm.OpSetProp))
	e.chunk.Write(byte(len(path)))
	for _, seg := range path {
		low := strings.ToLower(seg)
		ci := e.chunk.WriteConstant(low)
		if err := checkConstIndex(ci, " for OpSetProp path"); err != nil {
			return err
		}
		e.chunk.Write(byte(ci))
	}
	return nil
}
