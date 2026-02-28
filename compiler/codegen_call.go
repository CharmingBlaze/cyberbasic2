package compiler

import (
	"cyberbasic/compiler/parser"
	"cyberbasic/compiler/vm"
	"fmt"
	"strings"
)

// compileCall compiles a function/procedure call or array element read
func (c *Compiler) compileCall(call *parser.Call, chunk *vm.Chunk) error {
	if len(call.Arguments) == 0 && strings.EqualFold(call.Name, "shouldclose") {
		chunk.Write(byte(vm.OpShouldClose))
		return nil
	}

	if !strings.Contains(call.Name, ".") {
		if dims, ok := chunk.GetVarDims(call.Name); ok && len(dims) == len(call.Arguments) {
			for _, arg := range call.Arguments {
				if err := c.compileExpression(arg, chunk); err != nil {
					return err
				}
			}
			varIndex, _ := chunk.GetVariable(call.Name)
			chunk.Write(byte(vm.OpLoadArray))
			chunk.Write(byte(varIndex))
			return nil
		}
	}

	if strings.Contains(call.Name, ".") {
		for _, arg := range call.Arguments {
			if err := c.compileExpression(arg, chunk); err != nil {
				return err
			}
		}
		nameConst := strings.ToLower(call.Name)
		if c.userFuncs != nil && c.userFuncs[nameConst] {
			idx := chunk.WriteConstant(nameConst)
			if err := checkConstIndex(idx, " for user call"); err != nil {
				return err
			}
			chunk.Write(byte(vm.OpCallUser))
			chunk.Write(byte(idx))
			chunk.Write(byte(len(call.Arguments)))
			return nil
		}
		if flat := PhysicsNamespaceToFlat(nameConst); flat != "" {
			nameConst = flat
		} else if strings.HasPrefix(nameConst, "rl.") {
			nameConst = nameConst[3:]
		}
		idx := chunk.WriteConstant(nameConst)
		if err := checkConstIndex(idx, " for foreign call"); err != nil {
			return err
		}
		chunk.Write(byte(vm.OpCallForeign))
		chunk.Write(byte(idx))
		chunk.Write(byte(len(call.Arguments)))
		return nil
	}

	name := strings.ToLower(call.Name)
	switch name {
	case "print":
		for _, arg := range call.Arguments {
			if err := c.compileExpression(arg, chunk); err != nil {
				return err
			}
			chunk.Write(byte(vm.OpPrint))
		}
		return nil
	case "matmul":
		if len(call.Arguments) != 3 {
			return fmt.Errorf("MatMul(resultName, aName, bName) expects 3 arguments")
		}
		getName := func(n parser.Node) (string, bool) {
			switch v := n.(type) {
			case *parser.Identifier:
				return v.Name, true
			case *parser.StringLiteral:
				return v.Value, true
			default:
				return "", false
			}
		}
		r, ok1 := getName(call.Arguments[0])
		a, ok2 := getName(call.Arguments[1])
		b, ok3 := getName(call.Arguments[2])
		if !ok1 || !ok2 || !ok3 {
			return fmt.Errorf("MatMul: all 3 args must be variable names or string literals")
		}
		ri := chunk.WriteConstant(strings.ToLower(r))
		ai := chunk.WriteConstant(strings.ToLower(a))
		bi := chunk.WriteConstant(strings.ToLower(b))
		if err := checkConstIndex(ri, ""); err != nil {
			return err
		}
		if err := checkConstIndex(ai, ""); err != nil {
			return err
		}
		if err := checkConstIndex(bi, ""); err != nil {
			return err
		}
		chunk.Write(byte(vm.OpMatMul))
		chunk.Write(byte(ri))
		chunk.Write(byte(ai))
		chunk.Write(byte(bi))
		return nil
	}

	for _, arg := range call.Arguments {
		if err := c.compileExpression(arg, chunk); err != nil {
			return err
		}
	}

	if c.userFuncs != nil && c.userFuncs[name] {
		idx := chunk.WriteConstant(name)
		if err := checkConstIndex(idx, " for user call"); err != nil {
			return err
		}
		chunk.Write(byte(vm.OpCallUser))
		chunk.Write(byte(idx))
		chunk.Write(byte(len(call.Arguments)))
		return nil
	}

	switch name {
	case "str":
		if len(call.Arguments) != 1 {
			return fmt.Errorf("STR() expects 1 argument")
		}
		chunk.Write(byte(vm.OpStr))
		return nil
	case "random":
		if len(call.Arguments) == 0 {
			chunk.Write(byte(vm.OpRandom))
			return nil
		}
		if len(call.Arguments) == 1 {
			chunk.Write(byte(vm.OpRandomN))
			return nil
		}
		return fmt.Errorf("Random() expects 0 or 1 argument")
	case "sleep", "wait":
		if len(call.Arguments) != 1 {
			return fmt.Errorf("Sleep/Wait expect 1 argument (milliseconds)")
		}
		chunk.Write(byte(vm.OpSleep))
		return nil
	case "int":
		if len(call.Arguments) != 1 {
			return fmt.Errorf("Int() expects 1 argument")
		}
		chunk.Write(byte(vm.OpInt))
		return nil
	case "timer":
		if len(call.Arguments) != 0 {
			return fmt.Errorf("Timer() takes no arguments")
		}
		chunk.Write(byte(vm.OpTimer))
		return nil
	case "resettimer":
		if len(call.Arguments) != 0 {
			return fmt.Errorf("ResetTimer() takes no arguments")
		}
		chunk.Write(byte(vm.OpResetTimer))
		return nil
	case "quit":
		if len(call.Arguments) != 0 {
			return fmt.Errorf("Quit takes no arguments")
		}
		chunk.Write(byte(vm.OpQuit))
		return nil
	case "sin":
		if len(call.Arguments) != 1 {
			return fmt.Errorf("Sin() expects 1 argument")
		}
		chunk.Write(byte(vm.OpSin))
		return nil
	case "cos":
		if len(call.Arguments) != 1 {
			return fmt.Errorf("Cos() expects 1 argument")
		}
		chunk.Write(byte(vm.OpCos))
		return nil
	case "tan":
		if len(call.Arguments) != 1 {
			return fmt.Errorf("Tan() expects 1 argument")
		}
		chunk.Write(byte(vm.OpTan))
		return nil
	case "sqrt":
		if len(call.Arguments) != 1 {
			return fmt.Errorf("Sqrt() expects 1 argument")
		}
		chunk.Write(byte(vm.OpSqrt))
		return nil
	case "abs":
		if len(call.Arguments) != 1 {
			return fmt.Errorf("Abs() expects 1 argument")
		}
		chunk.Write(byte(vm.OpAbs))
		return nil
	case "lerp":
		if len(call.Arguments) != 3 {
			return fmt.Errorf("Lerp(a, b, t) expects 3 arguments")
		}
		chunk.Write(byte(vm.OpLerp))
		return nil
	case "noise", "noise2d", "perlin", "simplex":
		if len(call.Arguments) != 2 {
			return fmt.Errorf("Noise(x, y) expects 2 arguments")
		}
		chunk.Write(byte(vm.OpNoise2D))
		return nil
	case "openfile":
		if len(call.Arguments) != 2 {
			return fmt.Errorf("OpenFile(path, mode) expects 2 arguments; mode 0=read, 1=write, 2=append")
		}
		chunk.Write(byte(vm.OpOpenFile))
		return nil
	case "readline":
		if len(call.Arguments) != 1 {
			return fmt.Errorf("ReadLine(handle) expects 1 argument")
		}
		chunk.Write(byte(vm.OpReadLine))
		return nil
	case "writeline":
		if len(call.Arguments) != 2 {
			return fmt.Errorf("WriteLine(handle, text) expects 2 arguments")
		}
		chunk.Write(byte(vm.OpWriteLine))
		return nil
	case "closefile":
		if len(call.Arguments) != 1 {
			return fmt.Errorf("CloseFile(handle) expects 1 argument")
		}
		chunk.Write(byte(vm.OpCloseFile))
		return nil
	case "floor":
		if len(call.Arguments) != 1 {
			return fmt.Errorf("Floor() expects 1 argument")
		}
		chunk.Write(byte(vm.OpFloor))
		return nil
	case "ceil":
		if len(call.Arguments) != 1 {
			return fmt.Errorf("Ceil() expects 1 argument")
		}
		chunk.Write(byte(vm.OpCeil))
		return nil
	case "round":
		if len(call.Arguments) != 1 {
			return fmt.Errorf("Round() expects 1 argument")
		}
		chunk.Write(byte(vm.OpRound))
		return nil
	case "min":
		if len(call.Arguments) != 2 {
			return fmt.Errorf("Min(a, b) expects 2 arguments")
		}
		chunk.Write(byte(vm.OpMin))
		return nil
	case "max":
		if len(call.Arguments) != 2 {
			return fmt.Errorf("Max(a, b) expects 2 arguments")
		}
		chunk.Write(byte(vm.OpMax))
		return nil
	case "clamp":
		if len(call.Arguments) != 3 {
			return fmt.Errorf("Clamp(x, lo, hi) expects 3 arguments")
		}
		chunk.Write(byte(vm.OpClamp))
		return nil
	case "pow":
		if len(call.Arguments) != 2 {
			return fmt.Errorf("Pow(base, exp) expects 2 arguments")
		}
		chunk.Write(byte(vm.OpPow))
		return nil
	case "exp":
		if len(call.Arguments) != 1 {
			return fmt.Errorf("Exp() expects 1 argument")
		}
		chunk.Write(byte(vm.OpExp))
		return nil
	case "log":
		if len(call.Arguments) != 1 {
			return fmt.Errorf("Log() expects 1 argument")
		}
		chunk.Write(byte(vm.OpLog))
		return nil
	case "log10":
		if len(call.Arguments) != 1 {
			return fmt.Errorf("Log10() expects 1 argument")
		}
		chunk.Write(byte(vm.OpLog10))
		return nil
	case "atan2":
		if len(call.Arguments) != 2 {
			return fmt.Errorf("Atan2(y, x) expects 2 arguments")
		}
		chunk.Write(byte(vm.OpAtan2))
		return nil
	case "sign":
		if len(call.Arguments) != 1 {
			return fmt.Errorf("Sign() expects 1 argument")
		}
		chunk.Write(byte(vm.OpSign))
		return nil
	case "deg2rad":
		if len(call.Arguments) != 1 {
			return fmt.Errorf("Deg2Rad() expects 1 argument")
		}
		chunk.Write(byte(vm.OpDeg2Rad))
		return nil
	case "rad2deg":
		if len(call.Arguments) != 1 {
			return fmt.Errorf("Rad2Deg() expects 1 argument")
		}
		chunk.Write(byte(vm.OpRad2Deg))
		return nil
	case "distance2d":
		if len(call.Arguments) != 4 {
			return fmt.Errorf("Distance2D(x1, y1, x2, y2) expects 4 arguments")
		}
		chunk.Write(byte(vm.OpDistance2D))
		return nil
	case "distance3d":
		if len(call.Arguments) != 6 {
			return fmt.Errorf("Distance3D(x1, y1, z1, x2, y2, z2) expects 6 arguments")
		}
		chunk.Write(byte(vm.OpDistance3D))
		return nil
	case "distsq2d":
		if len(call.Arguments) != 4 {
			return fmt.Errorf("DistSq2D(x1, y1, x2, y2) expects 4 arguments")
		}
		chunk.Write(byte(vm.OpDistSq2D))
		return nil
	case "distsq3d":
		if len(call.Arguments) != 6 {
			return fmt.Errorf("DistSq3D(x1, y1, z1, x2, y2, z2) expects 6 arguments")
		}
		chunk.Write(byte(vm.OpDistSq3D))
		return nil
	case "inradius2d":
		if len(call.Arguments) != 5 {
			return fmt.Errorf("InRadius2D(x1, y1, x2, y2, radius) expects 5 arguments")
		}
		chunk.Write(byte(vm.OpInRadius2D))
		return nil
	case "inradius3d":
		if len(call.Arguments) != 7 {
			return fmt.Errorf("InRadius3D(x1, y1, z1, x2, y2, z2, radius) expects 7 arguments")
		}
		chunk.Write(byte(vm.OpInRadius3D))
		return nil
	case "angle2d":
		if len(call.Arguments) != 4 {
			return fmt.Errorf("Angle2D(x1, y1, x2, y2) expects 4 arguments")
		}
		chunk.Write(byte(vm.OpAngle2D))
		return nil
	case "left", "left$":
		if len(call.Arguments) != 2 {
			return fmt.Errorf("Left(s, n) expects 2 arguments")
		}
		chunk.Write(byte(vm.OpLeftStr))
		return nil
	case "right", "right$":
		if len(call.Arguments) != 2 {
			return fmt.Errorf("Right(s, n) expects 2 arguments")
		}
		chunk.Write(byte(vm.OpRightStr))
		return nil
	case "mid", "mid$":
		if len(call.Arguments) != 3 {
			return fmt.Errorf("Mid(s, start, n) expects 3 arguments; start is 1-based")
		}
		chunk.Write(byte(vm.OpMidStr))
		return nil
	case "len":
		if len(call.Arguments) != 1 {
			return fmt.Errorf("Len(s) expects 1 argument")
		}
		chunk.Write(byte(vm.OpLenStr))
		return nil
	case "eof":
		if len(call.Arguments) != 1 {
			return fmt.Errorf("EOF(handle) expects 1 argument")
		}
		chunk.Write(byte(vm.OpEOF))
		return nil
	}

	nameConst := strings.ToLower(call.Name)
	if strings.HasPrefix(nameConst, "rl.") {
		nameConst = nameConst[3:]
	}
	idx := chunk.WriteConstant(nameConst)
	if err := checkConstIndex(idx, " for foreign call"); err != nil {
		return err
	}
	chunk.Write(byte(vm.OpCallForeign))
	chunk.Write(byte(idx))
	chunk.Write(byte(len(call.Arguments)))
	return nil
}
