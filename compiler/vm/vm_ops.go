package vm

import (
	"cyberbasic/compiler/valueutil"
	"fmt"
	"math"
	"strconv"
)

// Stack operations
func (vm *VM) push(value Value) {
	vm.stack = append(vm.stack, value)
}

func (vm *VM) pop() Value {
	if len(vm.stack) == 0 {
		return nil
	}
	value := vm.stack[len(vm.stack)-1]
	vm.stack = vm.stack[:len(vm.stack)-1]
	return value
}

func (vm *VM) peek() Value {
	if len(vm.stack) == 0 {
		return nil
	}
	return vm.stack[len(vm.stack)-1]
}

// valueToString converts a VM value to string for runtime calls
func valueToString(v Value) string {
	if v == nil {
		return ""
	}
	switch x := v.(type) {
	case string:
		return x
	case int:
		return strconv.Itoa(x)
	case float64:
		return strconv.FormatFloat(x, 'f', -1, 64)
	case bool:
		if x {
			return "true"
		}
		return "false"
	default:
		return fmt.Sprintf("%v", v)
	}
}

// valueToFloat64 converts a VM value to float64 for runtime calls
func valueToFloat64(v Value) float64 {
	if v == nil {
		return 0
	}
	switch x := v.(type) {
	case int:
		return float64(x)
	case float64:
		return x
	case string:
		f, _ := strconv.ParseFloat(x, 64)
		return f
	default:
		return 0
	}
}

// valueToInt converts a VM value to int for runtime calls
func valueToInt(v Value) int {
	if v == nil {
		return 0
	}
	switch x := v.(type) {
	case int:
		return x
	case float64:
		return int(x)
	case string:
		i, _ := strconv.Atoi(x)
		return i
	default:
		return 0
	}
}

// Arithmetic operations
func (vm *VM) add(a, b Value) (Value, error) {
	switch a := a.(type) {
	case int:
		switch b := b.(type) {
		case int:
			return a + b, nil
		case float64:
			return float64(a) + b, nil
		case string:
			return strconv.Itoa(a) + b, nil
		}
	case float64:
		switch b := b.(type) {
		case int:
			return a + float64(b), nil
		case float64:
			return a + b, nil
		}
	case string:
		switch b := b.(type) {
		case int:
			return a + strconv.Itoa(b), nil
		case float64:
			return a + strconv.FormatFloat(b, 'f', -1, 64), nil
		case string:
			return a + b, nil
		}
	}
	return nil, fmt.Errorf("invalid operands for +")
}

func (vm *VM) subtract(a, b Value) (Value, error) {
	switch a := a.(type) {
	case int:
		switch b := b.(type) {
		case int:
			return a - b, nil
		case float64:
			return float64(a) - b, nil
		}
	case float64:
		switch b := b.(type) {
		case int:
			return a - float64(b), nil
		case float64:
			return a - b, nil
		}
	}
	return nil, fmt.Errorf("invalid operands for -")
}

func (vm *VM) multiply(a, b Value) (Value, error) {
	switch a := a.(type) {
	case int:
		switch b := b.(type) {
		case int:
			return a * b, nil
		case float64:
			return float64(a) * b, nil
		}
	case float64:
		switch b := b.(type) {
		case int:
			return a * float64(b), nil
		case float64:
			return a * b, nil
		}
	}
	return nil, fmt.Errorf("invalid operands for *")
}

func (vm *VM) divide(a, b Value) (Value, error) {
	switch a := a.(type) {
	case int:
		switch b := b.(type) {
		case int:
			if b == 0 {
				return nil, fmt.Errorf("division by zero")
			}
			return a / b, nil
		case float64:
			if b == 0 {
				return nil, fmt.Errorf("division by zero")
			}
			return float64(a) / b, nil
		}
	case float64:
		switch b := b.(type) {
		case int:
			if b == 0 {
				return nil, fmt.Errorf("division by zero")
			}
			return a / float64(b), nil
		case float64:
			if b == 0 {
				return nil, fmt.Errorf("division by zero")
			}
			return a / b, nil
		}
	}
	return nil, fmt.Errorf("invalid operands for /")
}

func (vm *VM) modulo(a, b Value) (Value, error) {
	switch a := a.(type) {
	case int:
		switch b := b.(type) {
		case int:
			if b == 0 {
				return nil, fmt.Errorf("division by zero")
			}
			return a % b, nil
		}
	}
	return nil, fmt.Errorf("invalid operands for modulo")
}

func (vm *VM) power(a, b Value) (Value, error) {
	af := valueToFloat64(a)
	bf := valueToFloat64(b)
	return math.Pow(af, bf), nil
}

func (vm *VM) intDiv(a, b Value) (Value, error) {
	af := valueToFloat64(a)
	bf := valueToFloat64(b)
	if bf == 0 {
		return nil, fmt.Errorf("division by zero")
	}
	return int(math.Trunc(af / bf)), nil
}

func (vm *VM) negate(a Value) (Value, error) {
	switch a := a.(type) {
	case int:
		return -a, nil
	case float64:
		return -a, nil
	}
	return nil, fmt.Errorf("invalid operand for unary -")
}

// Comparison operations
func (vm *VM) less(a, b Value) (bool, error) {
	switch a := a.(type) {
	case int:
		switch b := b.(type) {
		case int:
			return a < b, nil
		case float64:
			return float64(a) < b, nil
		}
	case float64:
		switch b := b.(type) {
		case int:
			return a < float64(b), nil
		case float64:
			return a < b, nil
		}
	case string:
		switch b := b.(type) {
		case string:
			return a < b, nil
		}
	}
	return false, fmt.Errorf("invalid operands for <")
}

func (vm *VM) lessEqual(a, b Value) (bool, error) {
	switch a := a.(type) {
	case int:
		switch b := b.(type) {
		case int:
			return a <= b, nil
		case float64:
			return float64(a) <= b, nil
		}
	case float64:
		switch b := b.(type) {
		case int:
			return a <= float64(b), nil
		case float64:
			return a <= b, nil
		}
	case string:
		switch b := b.(type) {
		case string:
			return a <= b, nil
		}
	}
	return false, fmt.Errorf("invalid operands for <=")
}

func (vm *VM) greater(a, b Value) (bool, error) {
	switch a := a.(type) {
	case int:
		switch b := b.(type) {
		case int:
			return a > b, nil
		case float64:
			return float64(a) > b, nil
		}
	case float64:
		switch b := b.(type) {
		case int:
			return a > float64(b), nil
		case float64:
			return a > b, nil
		}
	case string:
		switch b := b.(type) {
		case string:
			return a > b, nil
		}
	}
	return false, fmt.Errorf("invalid operands for >")
}

func (vm *VM) greaterEqual(a, b Value) (bool, error) {
	switch a := a.(type) {
	case int:
		switch b := b.(type) {
		case int:
			return a >= b, nil
		case float64:
			return float64(a) >= b, nil
		}
	case float64:
		switch b := b.(type) {
		case int:
			return a >= float64(b), nil
		case float64:
			return a >= b, nil
		}
	case string:
		switch b := b.(type) {
		case string:
			return a >= b, nil
		}
	}
	return false, fmt.Errorf("invalid operands for >=")
}

func (vm *VM) isTruthy(value Value) bool {
	return valueutil.IsTruthy(value)
}

// simplex2D returns 2D Simplex-style noise in [-1, 1] (simple hash-based implementation)
func simplex2D(x, y float64) float64 {
	const scale = 0.1
	xi := int(math.Floor(x*scale)) & 255
	yi := int(math.Floor(y*scale)) & 255
	h := (xi*37+yi*97)*971 + 1
	return (float64(h%1024)/512 - 1)
}
