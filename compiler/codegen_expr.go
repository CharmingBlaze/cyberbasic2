package compiler

import (
	"cyberbasic/compiler/parser"
	"cyberbasic/compiler/vm"
	"fmt"
	"strings"
)

// compileExpression compiles an expression
func (c *Compiler) compileExpression(expr parser.Node, chunk *vm.Chunk) error {
	if line := getSourceLine(expr); line > 0 {
		chunk.SetLine(line)
	}
	switch node := expr.(type) {
	case *parser.Number:
		return c.compileNumber(node, chunk)
	case *parser.StringLiteral:
		return c.compileString(node, chunk)
	case *parser.Boolean:
		return c.compileBoolean(node, chunk)
	case *parser.NilLiteral:
		return c.compileNilLiteral(chunk)
	case *parser.Identifier:
		return c.compileIdentifier(node, chunk)
	case *parser.BinaryOp:
		return c.compileBinaryOp(node, chunk)
	case *parser.UnaryOp:
		return c.compileUnaryOp(node, chunk)
	case *parser.Call:
		return c.compileCall(node, chunk)
	case *parser.MemberAccess:
		return c.compileMemberAccess(node, chunk)
	case *parser.JSONIndexAccess:
		return c.compileJSONIndexAccess(node, chunk)
	case *parser.DictLiteral:
		return c.compileDictLiteral(node, chunk)
	default:
		return errWithLine(expr, fmt.Errorf("unsupported expression type: %T", expr))
	}
}

// compileDictLiteral compiles { k: v, ... } as CreateDict then SetDictKey for each pair.
func (c *Compiler) compileDictLiteral(node *parser.DictLiteral, chunk *vm.Chunk) error {
	ci := chunk.WriteConstant("createdict")
	if err := checkConstIndex(ci, ""); err != nil {
		return err
	}
	chunk.Write(byte(vm.OpCallForeign))
	chunk.Write(byte(ci))
	chunk.Write(byte(0))
	for i, p := range node.Pairs {
		if i > 0 {
			chunk.Write(byte(vm.OpDup))
		}
		keyIdx := chunk.WriteConstant(p.Key)
		if err := checkConstIndex(keyIdx, " for dict key"); err != nil {
			return err
		}
		if err := c.compileExpression(p.Value, chunk); err != nil {
			return err
		}
		setIdx := chunk.WriteConstant("setdictkey")
		if err := checkConstIndex(setIdx, ""); err != nil {
			return err
		}
		chunk.Write(byte(vm.OpCallForeign))
		chunk.Write(byte(setIdx))
		chunk.Write(byte(3))
		if i < len(node.Pairs)-1 {
			chunk.Write(byte(vm.OpPop))
		}
	}
	return nil
}

// compileJSONIndexAccess compiles obj["key"] as GetJSONKey(obj, "key")
func (c *Compiler) compileJSONIndexAccess(node *parser.JSONIndexAccess, chunk *vm.Chunk) error {
	if err := c.compileExpression(node.Object, chunk); err != nil {
		return err
	}
	keyIdx := chunk.WriteConstant(node.Key)
	if err := checkConstIndex(keyIdx, " for JSON key"); err != nil {
		return err
	}
	chunk.Write(byte(vm.OpLoadConst))
	chunk.Write(byte(keyIdx))
	nameIdx := chunk.WriteConstant("getjsonkey")
	if err := checkConstIndex(nameIdx, ""); err != nil {
		return err
	}
	chunk.Write(byte(vm.OpCallForeign))
	chunk.Write(byte(nameIdx))
	chunk.Write(byte(2))
	return nil
}

// compileMemberAccess compiles expr.member: UDT constant group, entity property, namespace constant (RL.*), or getter (pos.x).
func (c *Compiler) compileMemberAccess(m *parser.MemberAccess, chunk *vm.Chunk) error {
	mb := strings.ToLower(m.Member)
	if id, ok := m.Object.(*parser.Identifier); ok {
		objLower := strings.ToLower(id.Name)
		if c.entityNames != nil && c.entityNames[objLower] {
			entityIdx := chunk.WriteConstant(objLower)
			propIdx := chunk.WriteConstant(mb)
			if err := checkConstIndex(entityIdx, " for entity prop"); err != nil {
				return err
			}
			if err := checkConstIndex(propIdx, " for entity prop"); err != nil {
				return err
			}
			chunk.Write(byte(vm.OpLoadEntityProp))
			chunk.Write(byte(entityIdx))
			chunk.Write(byte(propIdx))
			return nil
		}
		if c.typeDefs != nil {
			if td, ok := c.typeDefs[objLower]; ok {
				val, err := c.resolveUDTConstantMember(td, mb)
				if err == nil {
					key := objLower + "." + mb
					if idx, has := c.constIndices[key]; has {
						chunk.Write(byte(vm.OpLoadConst))
						chunk.Write(byte(idx))
						return nil
					}
					ci := chunk.WriteConstant(val)
					if err := checkConstIndex(ci, ""); err != nil {
						return err
					}
					c.constIndices[key] = byte(ci)
					chunk.Write(byte(vm.OpLoadConst))
					chunk.Write(byte(ci))
					return nil
				}
			}
		}
		if (objLower == "rl" || objLower == "box2d" || objLower == "bullet" || objLower == "game") && mb != "x" && mb != "y" && mb != "z" {
			idx := chunk.WriteConstant(mb)
			if err := checkConstIndex(idx, " for foreign constant"); err != nil {
				return err
			}
			chunk.Write(byte(vm.OpCallForeign))
			chunk.Write(byte(idx))
			chunk.Write(byte(0))
			return nil
		}
	}
	if err := c.compileExpression(m.Object, chunk); err != nil {
		return err
	}
	var name string
	switch mb {
	case "x":
		name = "getvector2x"
	case "y":
		name = "getvector2y"
	case "z":
		name = "getvector3z"
	default:
		name = "getvector2" + mb
	}
	idx := chunk.WriteConstant(name)
	if err := checkConstIndex(idx, " for member getter"); err != nil {
		return err
	}
	chunk.Write(byte(vm.OpCallForeign))
	chunk.Write(byte(idx))
	chunk.Write(byte(1))
	return nil
}

// compileNumber compiles a number literal
func (c *Compiler) compileNumber(num *parser.Number, chunk *vm.Chunk) error {
	if strings.Contains(num.Value, ".") {
		if floatVal, err := parseFloat(num.Value); err == nil {
			constIndex := chunk.WriteConstant(floatVal)
			chunk.Write(byte(vm.OpLoadConst))
			chunk.Write(byte(constIndex))
			return nil
		}
	}
	if intVal, err := parseInt(num.Value); err == nil {
		constIndex := chunk.WriteConstant(intVal)
		chunk.Write(byte(vm.OpLoadConst))
		chunk.Write(byte(constIndex))
		return nil
	}
	if floatVal, err := parseFloat(num.Value); err == nil {
		constIndex := chunk.WriteConstant(floatVal)
		chunk.Write(byte(vm.OpLoadConst))
		chunk.Write(byte(constIndex))
		return nil
	}
	return fmt.Errorf("invalid number format: %s", num.Value)
}

// compileString compiles a string literal
func (c *Compiler) compileString(str *parser.StringLiteral, chunk *vm.Chunk) error {
	constIndex := chunk.WriteConstant(str.Value)
	chunk.Write(byte(vm.OpLoadString))
	chunk.Write(byte(constIndex))
	return nil
}

// compileBoolean compiles a boolean literal
func (c *Compiler) compileBoolean(bool *parser.Boolean, chunk *vm.Chunk) error {
	constIndex := chunk.WriteConstant(bool.Value)
	chunk.Write(byte(vm.OpLoadConst))
	chunk.Write(byte(constIndex))
	return nil
}

// compileNilLiteral compiles the null/nil literal (pushes nil onto the stack)
func (c *Compiler) compileNilLiteral(chunk *vm.Chunk) error {
	constIndex := chunk.WriteConstant(nil)
	chunk.Write(byte(vm.OpLoadConst))
	chunk.Write(byte(constIndex))
	return nil
}

// compileIdentifier compiles a variable reference, a CONST, or a qualified "constant" (e.g. RL.DarkGray)
func (c *Compiler) compileIdentifier(ident *parser.Identifier, chunk *vm.Chunk) error {
	if c.funcParamIndices != nil {
		if idx, ok := c.funcParamIndices[strings.ToLower(ident.Name)]; ok {
			chunk.Write(byte(vm.OpLoadVar))
			chunk.Write(byte(idx))
			return nil
		}
	}
	if varIndex, exists := chunk.GetVariable(ident.Name); exists {
		chunk.Write(byte(vm.OpLoadVar))
		chunk.Write(byte(varIndex))
		return nil
	}
	if c.constIndices != nil {
		if idx, ok := c.constIndices[strings.ToLower(ident.Name)]; ok {
			chunk.Write(byte(vm.OpLoadConst))
			chunk.Write(byte(idx))
			return nil
		}
	}
	nameLower := strings.ToLower(ident.Name)
	if strings.HasPrefix(nameLower, "rl.") || strings.HasPrefix(nameLower, "box2d.") || strings.HasPrefix(nameLower, "bullet.") || strings.HasPrefix(nameLower, "game.") {
		nameConst := nameLower
		if flat := PhysicsNamespaceToFlat(nameConst); flat != "" {
			nameConst = flat
		} else if strings.HasPrefix(nameConst, "rl.") {
			nameConst = nameConst[3:]
		} else if strings.HasPrefix(nameConst, "box2d.") {
			nameConst = nameConst[6:]
		} else if strings.HasPrefix(nameConst, "bullet.") {
			nameConst = nameConst[7:]
		} else if strings.HasPrefix(nameConst, "game.") {
			nameConst = nameConst[5:]
		}
		idx := chunk.WriteConstant(nameConst)
		if err := checkConstIndex(idx, " for foreign constant"); err != nil {
			return err
		}
		chunk.Write(byte(vm.OpCallForeign))
		chunk.Write(byte(idx))
		chunk.Write(byte(0))
		return nil
	}
	constIndex := chunk.WriteConstant(ident.Name)
	chunk.Write(byte(vm.OpLoadGlobal))
	chunk.Write(byte(constIndex))
	return nil
}

// compileBinaryOp compiles a binary operation
func (c *Compiler) compileBinaryOp(op *parser.BinaryOp, chunk *vm.Chunk) error {
	if err := c.compileExpression(op.Left, chunk); err != nil {
		return err
	}
	if err := c.compileExpression(op.Right, chunk); err != nil {
		return err
	}
	opCanon := strings.ToLower(op.Operator)
	switch opCanon {
	case "+":
		chunk.Write(byte(vm.OpAdd))
	case "-":
		chunk.Write(byte(vm.OpSub))
	case "*":
		chunk.Write(byte(vm.OpMul))
	case "/":
		chunk.Write(byte(vm.OpDiv))
	case "%":
		chunk.Write(byte(vm.OpMod))
	case "^":
		chunk.Write(byte(vm.OpPower))
	case "\\":
		chunk.Write(byte(vm.OpIntDiv))
	case "=", "==":
		chunk.Write(byte(vm.OpEqual))
	case "<>":
		chunk.Write(byte(vm.OpNotEqual))
	case "<":
		chunk.Write(byte(vm.OpLess))
	case "<=":
		chunk.Write(byte(vm.OpLessEqual))
	case ">":
		chunk.Write(byte(vm.OpGreater))
	case ">=":
		chunk.Write(byte(vm.OpGreaterEqual))
	case "and":
		chunk.Write(byte(vm.OpAnd))
	case "or":
		chunk.Write(byte(vm.OpOr))
	case "xor":
		chunk.Write(byte(vm.OpXor))
	default:
		return errWithLine(op, fmt.Errorf("unsupported binary operator: %s", op.Operator))
	}
	return nil
}

// compileUnaryOp compiles a unary operation
func (c *Compiler) compileUnaryOp(op *parser.UnaryOp, chunk *vm.Chunk) error {
	if err := c.compileExpression(op.Operand, chunk); err != nil {
		return err
	}
	switch {
	case op.Operator == "-":
		chunk.Write(byte(vm.OpNeg))
	case strings.EqualFold(op.Operator, "NOT"):
		chunk.Write(byte(vm.OpNot))
	default:
		return errWithLine(op, fmt.Errorf("unsupported unary operator: %s", op.Operator))
	}
	return nil
}
