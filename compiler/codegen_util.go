package compiler

import (
	"cyberbasic/compiler/parser"
	"fmt"
)

// MaxConstIndex is the maximum constant index (0-255) that fits in a single byte in the chunk.
// Chunk constant indices are 1-byte; programs with more than 256 constants will fail at codegen
// with "too many constants". See checkConstIndex for the shared error reporting.
const MaxConstIndex = 255

// checkConstIndex returns an error if idx exceeds MaxConstIndex. kind is appended to the message, e.g. " for dict key".
func checkConstIndex(idx int, kind string) error {
	if idx > MaxConstIndex {
		return fmt.Errorf("too many constants" + kind)
	}
	return nil
}

// parseInt parses a decimal integer string (used for constant folding in codegen).
func parseInt(s string) (int, error) {
	var result int
	_, err := fmt.Sscanf(s, "%d", &result)
	return result, err
}

// parseFloat parses a float string (used for constant folding in codegen).
func parseFloat(s string) (float64, error) {
	var result float64
	_, err := fmt.Sscanf(s, "%f", &result)
	return result, err
}

// unwrapStatement returns the inner node for Statement wrappers, otherwise the node itself.
func unwrapStatement(n parser.Node) parser.Node {
	if s, ok := n.(*parser.Statement); ok && s.Value != nil {
		return s.Value
	}
	return n
}

// WalkStatements visits every statement in nodes and recursively enters If/For/While/Repeat/Select blocks.
// It returns true if pred returns true for any node (after unwrapping Statement). Used to implement
// "does this body contain X?" checks (e.g. user sub calls, 3D draw calls) without duplicating traversal.
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

// threeDDrawNames lists RL.* call names that are 3D draw commands (used for hybrid loop wrapping).
var threeDDrawNames = map[string]bool{
	"drawcube": true, "drawcubewires": true, "drawsphere": true, "drawspherewires": true,
	"drawmodel": true, "drawmodelsimple": true, "drawmodelex": true, "drawmodelwires": true, "drawplane": true,
	"drawline3d": true, "drawpoint3d": true, "drawcircle3d": true, "drawgrid": true,
	"drawcylinder": true, "drawcylinderwires": true, "drawray": true, "drawtriangle3d": true,
	"beginmode3d": true,
}

// predicate3DDraw returns true if n is a Call whose normalized name is in threeDDrawNames.
func predicate3DDraw(n parser.Node) bool {
	call, ok := n.(*parser.Call)
	if !ok {
		return false
	}
	name := normWindowShouldCloseName(call.Name)
	return threeDDrawNames[name]
}
