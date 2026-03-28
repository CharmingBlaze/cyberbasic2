package codegen

import (
	"cyberbasic/compiler/parser"
	"cyberbasic/compiler/vm"
	"fmt"
	"strings"
)

// compileExpression compiles an expression
func (e *Emitter) compileExpression(expr parser.Node) error {
	if line := getSourceLine(expr); line > 0 {
		e.chunk.SetLine(line)
	}
	switch node := expr.(type) {
	case *parser.Number:
		return e.compileNumber(node)
	case *parser.StringLiteral:
		return e.compileString(node)
	case *parser.Boolean:
		return e.compileBoolean(node)
	case *parser.NilLiteral:
		return e.compileNilLiteral()
	case *parser.Identifier:
		return e.compileIdentifier(node)
	case *parser.BinaryOp:
		return e.compileBinaryOp(node)
	case *parser.UnaryOp:
		return e.compileUnaryOp(node)
	case *parser.Call:
		return e.compileCall(node)
	case *parser.MemberAccess:
		return e.compileMemberAccess(node)
	case *parser.JSONIndexAccess:
		return e.compileJSONIndexAccess(node)
	case *parser.SliceExpr:
		return e.compileSliceExpr(node)
	case *parser.InterpolatedString:
		return e.compileInterpolatedString(node)
	case *parser.DictLiteral:
		return e.compileDictLiteral(node)
	default:
		return errWithLine(expr, fmt.Errorf("unsupported expression type: %T", expr))
	}
}

// compileDictLiteral compiles { k: v, ... } as CreateDict then SetDictKey for each pair.
func (e *Emitter) compileDictLiteral(node *parser.DictLiteral) error {
	ci := e.chunk.WriteConstant("createdict")
	if err := checkConstIndex(ci, ""); err != nil {
		return err
	}
	e.chunk.Write(byte(vm.OpCallForeign))
	e.chunk.Write(byte(ci))
	e.chunk.Write(byte(0))
	for i, p := range node.Pairs {
		if i > 0 {
			e.chunk.Write(byte(vm.OpDup))
		}
		keyIdx := e.chunk.WriteConstant(p.Key)
		if err := checkConstIndex(keyIdx, " for dict key"); err != nil {
			return err
		}
		if err := e.compileExpression(p.Value); err != nil {
			return err
		}
		setIdx := e.chunk.WriteConstant("setdictkey")
		if err := checkConstIndex(setIdx, ""); err != nil {
			return err
		}
		e.chunk.Write(byte(vm.OpCallForeign))
		e.chunk.Write(byte(setIdx))
		e.chunk.Write(byte(3))
		if i < len(node.Pairs)-1 {
			e.chunk.Write(byte(vm.OpPop))
		}
	}
	return nil
}

// compileSliceExpr compiles s[start:end], s[i] (strings), or arr[i,j] (multi-dim arrays).
func (e *Emitter) compileSliceExpr(node *parser.SliceExpr) error {
	// Multi-dim array access: arr[i,j,k]
	if node.Indices != nil {
		if id, ok := node.Object.(*parser.Identifier); ok && !strings.Contains(id.Name, ".") {
			dims, hasDims := e.chunk.GetVarDims(id.Name)
			if hasDims && (len(dims) == len(node.Indices) || (len(dims) == 0 && len(node.Indices) == 1)) {
				// Push indices in natural order (first dim first; OpLoadArray pops last dim first)
				for _, idx := range node.Indices {
					if err := e.compileExpression(idx); err != nil {
						return err
					}
				}
				varIndex, _ := e.chunk.GetVariable(id.Name)
				e.chunk.Write(byte(vm.OpLoadArray))
				e.chunk.Write(byte(varIndex))
				return nil
			}
		}
		return fmt.Errorf("multi-index [i,j,...] requires an array variable with matching dimensions")
	}
	if err := e.compileExpression(node.Object); err != nil {
		return err
	}
	if !node.HasColon && node.Start != nil {
		// Single index: s[i] (string) or arr[i] (1D array)
		if id, ok := node.Object.(*parser.Identifier); ok && !strings.Contains(id.Name, ".") {
			dims, hasDims := e.chunk.GetVarDims(id.Name)
			if hasDims && (len(dims) == 1 || (len(dims) == 0)) {
				// 1D array access
				if err := e.compileExpression(node.Start); err != nil {
					return err
				}
				varIndex, _ := e.chunk.GetVariable(id.Name)
				e.chunk.Write(byte(vm.OpLoadArray))
				e.chunk.Write(byte(varIndex))
				return nil
			}
		}
		// Single char s[i]: push i, i+1 (0-based, end exclusive)
		if err := e.compileExpression(node.Start); err != nil {
			return err
		}
		if err := e.compileExpression(node.Start); err != nil {
			return err
		}
		ci := e.chunk.WriteConstant(1)
		if err := checkConstIndex(ci, ""); err != nil {
			return err
		}
		e.chunk.Write(byte(vm.OpLoadConst))
		e.chunk.Write(byte(ci))
		e.chunk.Write(byte(vm.OpAdd))
	} else {
		// Range: s[start:end], s[start:], s[:end], s[:]
		if node.End == nil && node.Start != nil {
			// s[start:]: use OpStrSliceFrom (obj, start) -> s[start:]
			if err := e.compileExpression(node.Start); err != nil {
				return err
			}
			e.chunk.Write(byte(vm.OpStrSliceFrom))
			return nil
		} else if node.End == nil && node.Start == nil {
			// s[:]: full string -> [obj, 0, len(obj)]
			e.chunk.Write(byte(vm.OpDup))
			e.chunk.Write(byte(vm.OpLenStr))
			ci := e.chunk.WriteConstant(0)
			if err := checkConstIndex(ci, ""); err != nil {
				return err
			}
			e.chunk.Write(byte(vm.OpLoadConst))
			e.chunk.Write(byte(ci))
			e.chunk.Write(byte(vm.OpSwap))
			e.chunk.Write(byte(vm.OpStrSlice))
		} else if node.End == nil {
			// s[:end] - start is nil, end is set. Shouldn't happen with our grammar.
			ci := e.chunk.WriteConstant(0)
			if err := checkConstIndex(ci, ""); err != nil {
				return err
			}
			e.chunk.Write(byte(vm.OpLoadConst))
			e.chunk.Write(byte(ci))
			if err := e.compileExpression(node.End); err != nil {
				return err
			}
			e.chunk.Write(byte(vm.OpStrSlice))
		} else {
			if node.Start != nil {
				if err := e.compileExpression(node.Start); err != nil {
					return err
				}
			} else {
				ci := e.chunk.WriteConstant(0)
				if err := checkConstIndex(ci, ""); err != nil {
					return err
				}
				e.chunk.Write(byte(vm.OpLoadConst))
				e.chunk.Write(byte(ci))
			}
			if err := e.compileExpression(node.End); err != nil {
				return err
			}
		}
	}
	e.chunk.Write(byte(vm.OpStrSlice))
	return nil
}

// compileInterpolatedString compiles "Hello {x}!" as "Hello " + Str(x) + "!"
func (e *Emitter) compileInterpolatedString(node *parser.InterpolatedString) error {
	for i, p := range node.Parts {
		if sl, ok := p.(*parser.StringLiteral); ok {
			if err := e.compileString(sl); err != nil {
				return err
			}
		} else {
			if err := e.compileExpression(p); err != nil {
				return err
			}
			e.chunk.Write(byte(vm.OpStr))
		}
		if i > 0 {
			e.chunk.Write(byte(vm.OpAdd))
		}
	}
	if len(node.Parts) == 0 {
		ci := e.chunk.WriteConstant("")
		if err := checkConstIndex(ci, ""); err != nil {
			return err
		}
		e.chunk.Write(byte(vm.OpLoadConst))
		e.chunk.Write(byte(ci))
	}
	return nil
}

// compileJSONIndexAccess compiles obj["key"] as GetJSONKey(obj, "key")
func (e *Emitter) compileJSONIndexAccess(node *parser.JSONIndexAccess) error {
	if err := e.compileExpression(node.Object); err != nil {
		return err
	}
	keyIdx := e.chunk.WriteConstant(node.Key)
	if err := checkConstIndex(keyIdx, " for JSON key"); err != nil {
		return err
	}
	e.chunk.Write(byte(vm.OpLoadConst))
	e.chunk.Write(byte(keyIdx))
	nameIdx := e.chunk.WriteConstant("getjsonkey")
	if err := checkConstIndex(nameIdx, ""); err != nil {
		return err
	}
	e.chunk.Write(byte(vm.OpCallForeign))
	e.chunk.Write(byte(nameIdx))
	e.chunk.Write(byte(2))
	return nil
}

// compileMemberAccess compiles expr.member: UDT constant group, entity property, namespace constant (RL.*), vector .x/.y/.z, or DotObject OpGetProp chain.
func (e *Emitter) compileMemberAccess(m *parser.MemberAccess) error {
	segs, base := collectMemberAccessChain(m)
	if len(segs) == 0 {
		return fmt.Errorf("empty member access chain")
	}
	mb := strings.ToLower(segs[len(segs)-1])
	if id, ok := base.(*parser.Identifier); ok {
		objLower := strings.ToLower(id.Name)
		if e.sem.EntityNames != nil && e.sem.EntityNames[objLower] && len(segs) == 1 {
			entityIdx := e.chunk.WriteConstant(objLower)
			propIdx := e.chunk.WriteConstant(mb)
			if err := checkConstIndex(entityIdx, " for entity prop"); err != nil {
				return err
			}
			if err := checkConstIndex(propIdx, " for entity prop"); err != nil {
				return err
			}
			e.chunk.Write(byte(vm.OpLoadEntityProp))
			e.chunk.Write(byte(entityIdx))
			e.chunk.Write(byte(propIdx))
			return nil
		}
		if e.sem.TypeDefs != nil {
			if td, ok := e.sem.TypeDefs[objLower]; ok && len(segs) == 1 {
				val, err := e.resolveUDTConstantMember(td, mb)
				if err == nil {
					key := objLower + "." + mb
					if idx, has := e.constIndices[key]; has {
						e.chunk.Write(byte(vm.OpLoadConst))
						e.chunk.Write(byte(idx))
						return nil
					}
					ci := e.chunk.WriteConstant(val)
					if err := checkConstIndex(ci, ""); err != nil {
						return err
					}
					e.constIndices[key] = byte(ci)
					e.chunk.Write(byte(vm.OpLoadConst))
					e.chunk.Write(byte(ci))
					return nil
				}
			}
		}
		if (objLower == "rl" || objLower == "box2d" || objLower == "bullet" || objLower == "game") && len(segs) == 1 && mb != "x" && mb != "y" && mb != "z" {
			idx := e.chunk.WriteConstant(mb)
			if err := checkConstIndex(idx, " for foreign constant"); err != nil {
				return err
			}
			e.chunk.Write(byte(vm.OpCallForeign))
			e.chunk.Write(byte(idx))
			e.chunk.Write(byte(0))
			return nil
		}
	}

	// DotObject roots (window, physics, …): never treat .x/.y/.z as vector swizzle on the namespace itself.
	if id, ok := base.(*parser.Identifier); ok {
		if dotObjectRoots[strings.ToLower(id.Name)] {
			if err := e.compileIdentifier(id); err != nil {
				return err
			}
			return e.emitOpGetProp(segs)
		}
	}

	// Vector components: single segment x/y/z on any expression (e.g. pos.x, GetMousePosition().x)
	if len(segs) == 1 && (mb == "x" || mb == "y" || mb == "z") {
		if err := e.compileExpression(base); err != nil {
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
		}
		idx := e.chunk.WriteConstant(name)
		if err := checkConstIndex(idx, " for member getter"); err != nil {
			return err
		}
		e.chunk.Write(byte(vm.OpCallForeign))
		e.chunk.Write(byte(idx))
		e.chunk.Write(byte(1))
		return nil
	}

	// DotObject property path (WINDOW.*, nested handles, VAR.prop...)
	if bid, ok := base.(*parser.Identifier); ok {
		if err := e.compileIdentifier(bid); err != nil {
			return err
		}
		return e.emitOpGetProp(segs)
	}
	if err := e.compileExpression(base); err != nil {
		return err
	}
	return e.emitOpGetProp(segs)
}

// compileNumber compiles a number literal
func (e *Emitter) compileNumber(num *parser.Number) error {
	if strings.Contains(num.Value, ".") {
		if floatVal, err := parseFloat(num.Value); err == nil {
			constIndex := e.chunk.WriteConstant(floatVal)
			e.chunk.Write(byte(vm.OpLoadConst))
			e.chunk.Write(byte(constIndex))
			return nil
		}
	}
	if intVal, err := parseInt(num.Value); err == nil {
		constIndex := e.chunk.WriteConstant(intVal)
		e.chunk.Write(byte(vm.OpLoadConst))
		e.chunk.Write(byte(constIndex))
		return nil
	}
	if floatVal, err := parseFloat(num.Value); err == nil {
		constIndex := e.chunk.WriteConstant(floatVal)
		e.chunk.Write(byte(vm.OpLoadConst))
		e.chunk.Write(byte(constIndex))
		return nil
	}
	return fmt.Errorf("invalid number format: %s", num.Value)
}

// compileString compiles a string literal
func (e *Emitter) compileString(str *parser.StringLiteral) error {
	constIndex := e.chunk.WriteConstant(str.Value)
	e.chunk.Write(byte(vm.OpLoadString))
	e.chunk.Write(byte(constIndex))
	return nil
}

// compileBoolean compiles a boolean literal
func (e *Emitter) compileBoolean(bool *parser.Boolean) error {
	constIndex := e.chunk.WriteConstant(bool.Value)
	e.chunk.Write(byte(vm.OpLoadConst))
	e.chunk.Write(byte(constIndex))
	return nil
}

// compileNilLiteral compiles the null/nil literal (pushes nil onto the stack)
func (e *Emitter) compileNilLiteral() error {
	constIndex := e.chunk.WriteConstant(nil)
	e.chunk.Write(byte(vm.OpLoadConst))
	e.chunk.Write(byte(constIndex))
	return nil
}

// compileIdentifier compiles a variable reference, a CONST, or a qualified "constant" (e.g. RL.DarkGray)
func (e *Emitter) compileIdentifier(ident *parser.Identifier) error {
	if e.funcParamIndices != nil {
		if idx, ok := e.funcParamIndices[strings.ToLower(ident.Name)]; ok {
			e.chunk.Write(byte(vm.OpLoadParam))
			e.chunk.Write(byte(idx))
			return nil
		}
	}
	if varIndex, exists := e.chunk.GetVariable(ident.Name); exists {
		e.chunk.Write(byte(vm.OpLoadVar))
		e.chunk.Write(byte(varIndex))
		return nil
	}
	if e.constIndices != nil {
		if idx, ok := e.constIndices[strings.ToLower(ident.Name)]; ok {
			e.chunk.Write(byte(vm.OpLoadConst))
			e.chunk.Write(byte(idx))
			return nil
		}
	}
	nameLower := strings.ToLower(ident.Name)
	if strings.HasPrefix(nameLower, "rl.") || strings.HasPrefix(nameLower, "box2d.") || strings.HasPrefix(nameLower, "bullet.") || strings.HasPrefix(nameLower, "game.") {
		nameConst := nameLower
		if flat := physicsNamespaceToFlat(nameConst); flat != "" {
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
		idx := e.chunk.WriteConstant(nameConst)
		if err := checkConstIndex(idx, " for foreign constant"); err != nil {
			return err
		}
		e.chunk.Write(byte(vm.OpCallForeign))
		e.chunk.Write(byte(idx))
		e.chunk.Write(byte(0))
		return nil
	}
	constIndex := e.chunk.WriteConstant(ident.Name)
	e.chunk.Write(byte(vm.OpLoadGlobal))
	e.chunk.Write(byte(constIndex))
	return nil
}

// compileBinaryOp compiles a binary operation
func (e *Emitter) compileBinaryOp(op *parser.BinaryOp) error {
	if err := e.compileExpression(op.Left); err != nil {
		return err
	}
	if err := e.compileExpression(op.Right); err != nil {
		return err
	}
	opCanon := strings.ToLower(op.Operator)
	switch opCanon {
	case "+":
		e.chunk.Write(byte(vm.OpAdd))
	case "-":
		e.chunk.Write(byte(vm.OpSub))
	case "*":
		e.chunk.Write(byte(vm.OpMul))
	case "/":
		e.chunk.Write(byte(vm.OpDiv))
	case "%":
		e.chunk.Write(byte(vm.OpMod))
	case "^":
		e.chunk.Write(byte(vm.OpPower))
	case "\\":
		e.chunk.Write(byte(vm.OpIntDiv))
	case "=", "==":
		e.chunk.Write(byte(vm.OpEqual))
	case "<>":
		e.chunk.Write(byte(vm.OpNotEqual))
	case "<":
		e.chunk.Write(byte(vm.OpLess))
	case "<=":
		e.chunk.Write(byte(vm.OpLessEqual))
	case ">":
		e.chunk.Write(byte(vm.OpGreater))
	case ">=":
		e.chunk.Write(byte(vm.OpGreaterEqual))
	case "and":
		e.chunk.Write(byte(vm.OpAnd))
	case "or":
		e.chunk.Write(byte(vm.OpOr))
	case "xor":
		e.chunk.Write(byte(vm.OpXor))
	default:
		return errWithLine(op, fmt.Errorf("unsupported binary operator: %s", op.Operator))
	}
	return nil
}

// compileUnaryOp compiles a unary operation
func (e *Emitter) compileUnaryOp(op *parser.UnaryOp) error {
	if err := e.compileExpression(op.Operand); err != nil {
		return err
	}
	switch {
	case op.Operator == "-":
		e.chunk.Write(byte(vm.OpNeg))
	case strings.EqualFold(op.Operator, "NOT"):
		e.chunk.Write(byte(vm.OpNot))
	default:
		return errWithLine(op, fmt.Errorf("unsupported unary operator: %s", op.Operator))
	}
	return nil
}
