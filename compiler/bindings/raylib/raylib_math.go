// Package raylib: math utils (Clamp, Lerp, etc.), Vector2/3, Matrix, Quaternion.
package raylib

import (
	"cyberbasic/compiler/vm"
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// argsToMatrix reads 16 floats from args starting at startIndex (row-major: M0,M4,M8,M12, M1,M5,M9,M13, ...). Returns zero matrix if not enough args.
func argsToMatrix(args []interface{}, startIndex int) rl.Matrix {
	if len(args) < startIndex+16 {
		return rl.Matrix{}
	}
	return rl.NewMatrix(
		toFloat32(args[startIndex+0]), toFloat32(args[startIndex+1]), toFloat32(args[startIndex+2]), toFloat32(args[startIndex+3]),
		toFloat32(args[startIndex+4]), toFloat32(args[startIndex+5]), toFloat32(args[startIndex+6]), toFloat32(args[startIndex+7]),
		toFloat32(args[startIndex+8]), toFloat32(args[startIndex+9]), toFloat32(args[startIndex+10]), toFloat32(args[startIndex+11]),
		toFloat32(args[startIndex+12]), toFloat32(args[startIndex+13]), toFloat32(args[startIndex+14]), toFloat32(args[startIndex+15]),
	)
}

func matrixToSlice(m rl.Matrix) []interface{} {
	return []interface{}{float64(m.M0), float64(m.M4), float64(m.M8), float64(m.M12), float64(m.M1), float64(m.M5), float64(m.M9), float64(m.M13), float64(m.M2), float64(m.M6), float64(m.M10), float64(m.M14), float64(m.M3), float64(m.M7), float64(m.M11), float64(m.M15)}
}

func vec2ToSlice(v rl.Vector2) []interface{} { return []interface{}{float64(v.X), float64(v.Y)} }
func vec3ToSlice(v rl.Vector3) []interface{} { return []interface{}{float64(v.X), float64(v.Y), float64(v.Z)} }
func quatToSlice(q rl.Quaternion) []interface{} {
	return []interface{}{float64(q.X), float64(q.Y), float64(q.Z), float64(q.W)}
}

func registerMath(v *vm.VM) {
	// Utils math
	v.RegisterForeign("Clamp", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("Clamp requires (value, min, max)")
		}
		return float64(rl.Clamp(toFloat32(args[0]), toFloat32(args[1]), toFloat32(args[2]))), nil
	})
	v.RegisterForeign("Lerp", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("Lerp requires (start, end, amount)")
		}
		return float64(rl.Lerp(toFloat32(args[0]), toFloat32(args[1]), toFloat32(args[2]))), nil
	})
	v.RegisterForeign("Normalize", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("Normalize requires (value, start, end)")
		}
		return float64(rl.Normalize(toFloat32(args[0]), toFloat32(args[1]), toFloat32(args[2]))), nil
	})
	v.RegisterForeign("Remap", func(args []interface{}) (interface{}, error) {
		if len(args) < 5 {
			return nil, fmt.Errorf("Remap requires (value, inputStart, inputEnd, outputStart, outputEnd)")
		}
		return float64(rl.Remap(toFloat32(args[0]), toFloat32(args[1]), toFloat32(args[2]), toFloat32(args[3]), toFloat32(args[4]))), nil
	})
	v.RegisterForeign("Wrap", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("Wrap requires (value, min, max)")
		}
		return float64(rl.Wrap(toFloat32(args[0]), toFloat32(args[1]), toFloat32(args[2]))), nil
	})
	v.RegisterForeign("FloatEquals", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return false, nil
		}
		return rl.FloatEquals(toFloat32(args[0]), toFloat32(args[1])), nil
	})

	// Vector2 math
	v.RegisterForeign("Vector2Zero", func(args []interface{}) (interface{}, error) {
		return vec2ToSlice(rl.Vector2Zero()), nil
	})
	v.RegisterForeign("Vector2One", func(args []interface{}) (interface{}, error) {
		return vec2ToSlice(rl.Vector2One()), nil
	})
	v.RegisterForeign("Vector2Add", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("Vector2Add requires (x1,y1, x2,y2)")
		}
		v1 := rl.Vector2{X: toFloat32(args[0]), Y: toFloat32(args[1])}
		v2 := rl.Vector2{X: toFloat32(args[2]), Y: toFloat32(args[3])}
		return vec2ToSlice(rl.Vector2Add(v1, v2)), nil
	})
	v.RegisterForeign("Vector2AddValue", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("Vector2AddValue requires (x, y, add)")
		}
		vec := rl.Vector2{X: toFloat32(args[0]), Y: toFloat32(args[1])}
		return vec2ToSlice(rl.Vector2AddValue(vec, toFloat32(args[2]))), nil
	})
	v.RegisterForeign("Vector2Subtract", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("Vector2Subtract requires (x1,y1, x2,y2)")
		}
		v1 := rl.Vector2{X: toFloat32(args[0]), Y: toFloat32(args[1])}
		v2 := rl.Vector2{X: toFloat32(args[2]), Y: toFloat32(args[3])}
		return vec2ToSlice(rl.Vector2Subtract(v1, v2)), nil
	})
	v.RegisterForeign("Vector2SubtractValue", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("Vector2SubtractValue requires (x, y, sub)")
		}
		vec := rl.Vector2{X: toFloat32(args[0]), Y: toFloat32(args[1])}
		return vec2ToSlice(rl.Vector2SubtractValue(vec, toFloat32(args[2]))), nil
	})
	v.RegisterForeign("Vector2Length", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("Vector2Length requires (x, y)")
		}
		vec := rl.Vector2{X: toFloat32(args[0]), Y: toFloat32(args[1])}
		return float64(rl.Vector2Length(vec)), nil
	})
	v.RegisterForeign("Vector2LengthSqr", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("Vector2LengthSqr requires (x, y)")
		}
		vec := rl.Vector2{X: toFloat32(args[0]), Y: toFloat32(args[1])}
		return float64(rl.Vector2LengthSqr(vec)), nil
	})
	v.RegisterForeign("Vector2DotProduct", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("Vector2DotProduct requires (x1,y1, x2,y2)")
		}
		v1 := rl.Vector2{X: toFloat32(args[0]), Y: toFloat32(args[1])}
		v2 := rl.Vector2{X: toFloat32(args[2]), Y: toFloat32(args[3])}
		return float64(rl.Vector2DotProduct(v1, v2)), nil
	})
	v.RegisterForeign("Vector2Distance", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("Vector2Distance requires (x1,y1, x2,y2)")
		}
		v1 := rl.Vector2{X: toFloat32(args[0]), Y: toFloat32(args[1])}
		v2 := rl.Vector2{X: toFloat32(args[2]), Y: toFloat32(args[3])}
		return float64(rl.Vector2Distance(v1, v2)), nil
	})
	v.RegisterForeign("Vector2DistanceSqr", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("Vector2DistanceSqr requires (x1,y1, x2,y2)")
		}
		v1 := rl.Vector2{X: toFloat32(args[0]), Y: toFloat32(args[1])}
		v2 := rl.Vector2{X: toFloat32(args[2]), Y: toFloat32(args[3])}
		return float64(rl.Vector2DistanceSqr(v1, v2)), nil
	})
	v.RegisterForeign("Vector2Angle", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("Vector2Angle requires (x1,y1, x2,y2)")
		}
		v1 := rl.Vector2{X: toFloat32(args[0]), Y: toFloat32(args[1])}
		v2 := rl.Vector2{X: toFloat32(args[2]), Y: toFloat32(args[3])}
		return float64(rl.Vector2Angle(v1, v2)), nil
	})
	v.RegisterForeign("Vector2Scale", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("Vector2Scale requires (x, y, scale)")
		}
		vec := rl.Vector2{X: toFloat32(args[0]), Y: toFloat32(args[1])}
		return vec2ToSlice(rl.Vector2Scale(vec, toFloat32(args[2]))), nil
	})
	v.RegisterForeign("Vector2Multiply", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("Vector2Multiply requires (x1,y1, x2,y2)")
		}
		v1 := rl.Vector2{X: toFloat32(args[0]), Y: toFloat32(args[1])}
		v2 := rl.Vector2{X: toFloat32(args[2]), Y: toFloat32(args[3])}
		return vec2ToSlice(rl.Vector2Multiply(v1, v2)), nil
	})
	v.RegisterForeign("Vector2Negate", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("Vector2Negate requires (x, y)")
		}
		vec := rl.Vector2{X: toFloat32(args[0]), Y: toFloat32(args[1])}
		return vec2ToSlice(rl.Vector2Negate(vec)), nil
	})
	v.RegisterForeign("Vector2Divide", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("Vector2Divide requires (x1,y1, x2,y2)")
		}
		v1 := rl.Vector2{X: toFloat32(args[0]), Y: toFloat32(args[1])}
		v2 := rl.Vector2{X: toFloat32(args[2]), Y: toFloat32(args[3])}
		return vec2ToSlice(rl.Vector2Divide(v1, v2)), nil
	})
	v.RegisterForeign("Vector2Normalize", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("Vector2Normalize requires (x, y)")
		}
		vec := rl.Vector2{X: toFloat32(args[0]), Y: toFloat32(args[1])}
		return vec2ToSlice(rl.Vector2Normalize(vec)), nil
	})
	v.RegisterForeign("Vector2Transform", func(args []interface{}) (interface{}, error) {
		if len(args) < 2+16 {
			return nil, fmt.Errorf("Vector2Transform requires (x, y, 16 matrix floats)")
		}
		vec := rl.Vector2{X: toFloat32(args[0]), Y: toFloat32(args[1])}
		mat := argsToMatrix(args, 2)
		return vec2ToSlice(rl.Vector2Transform(vec, mat)), nil
	})
	v.RegisterForeign("Vector2Lerp", func(args []interface{}) (interface{}, error) {
		if len(args) < 5 {
			return nil, fmt.Errorf("Vector2Lerp requires (x1,y1, x2,y2, amount)")
		}
		v1 := rl.Vector2{X: toFloat32(args[0]), Y: toFloat32(args[1])}
		v2 := rl.Vector2{X: toFloat32(args[2]), Y: toFloat32(args[3])}
		return vec2ToSlice(rl.Vector2Lerp(v1, v2, toFloat32(args[4]))), nil
	})
	v.RegisterForeign("Vector2Reflect", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("Vector2Reflect requires (vx,vy, normalX, normalY)")
		}
		vec := rl.Vector2{X: toFloat32(args[0]), Y: toFloat32(args[1])}
		normal := rl.Vector2{X: toFloat32(args[2]), Y: toFloat32(args[3])}
		return vec2ToSlice(rl.Vector2Reflect(vec, normal)), nil
	})
	v.RegisterForeign("Vector2Rotate", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("Vector2Rotate requires (x, y, angle)")
		}
		vec := rl.Vector2{X: toFloat32(args[0]), Y: toFloat32(args[1])}
		return vec2ToSlice(rl.Vector2Rotate(vec, toFloat32(args[2]))), nil
	})
	v.RegisterForeign("Vector2MoveTowards", func(args []interface{}) (interface{}, error) {
		if len(args) < 5 {
			return nil, fmt.Errorf("Vector2MoveTowards requires (x,y, targetX,targetY, maxDistance)")
		}
		vec := rl.Vector2{X: toFloat32(args[0]), Y: toFloat32(args[1])}
		target := rl.Vector2{X: toFloat32(args[2]), Y: toFloat32(args[3])}
		return vec2ToSlice(rl.Vector2MoveTowards(vec, target, toFloat32(args[4]))), nil
	})
	v.RegisterForeign("Vector2Invert", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("Vector2Invert requires (x, y)")
		}
		vec := rl.Vector2{X: toFloat32(args[0]), Y: toFloat32(args[1])}
		return vec2ToSlice(rl.Vector2Invert(vec)), nil
	})
	v.RegisterForeign("Vector2Clamp", func(args []interface{}) (interface{}, error) {
		if len(args) < 6 {
			return nil, fmt.Errorf("Vector2Clamp requires (vx,vy, minX,minY, maxX,maxY)")
		}
		vec := rl.Vector2{X: toFloat32(args[0]), Y: toFloat32(args[1])}
		minV := rl.Vector2{X: toFloat32(args[2]), Y: toFloat32(args[3])}
		maxV := rl.Vector2{X: toFloat32(args[4]), Y: toFloat32(args[5])}
		return vec2ToSlice(rl.Vector2Clamp(vec, minV, maxV)), nil
	})
	v.RegisterForeign("Vector2ClampValue", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("Vector2ClampValue requires (x, y, min, max)")
		}
		vec := rl.Vector2{X: toFloat32(args[0]), Y: toFloat32(args[1])}
		return vec2ToSlice(rl.Vector2ClampValue(vec, toFloat32(args[2]), toFloat32(args[3]))), nil
	})
	v.RegisterForeign("Vector2Equals", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return false, nil
		}
		p := rl.Vector2{X: toFloat32(args[0]), Y: toFloat32(args[1])}
		q := rl.Vector2{X: toFloat32(args[2]), Y: toFloat32(args[3])}
		return rl.Vector2Equals(p, q), nil
	})

	// Vector3 math
	v.RegisterForeign("Vector3Zero", func(args []interface{}) (interface{}, error) {
		return vec3ToSlice(rl.Vector3Zero()), nil
	})
	v.RegisterForeign("Vector3One", func(args []interface{}) (interface{}, error) {
		return vec3ToSlice(rl.Vector3One()), nil
	})
	v.RegisterForeign("Vector3Add", func(args []interface{}) (interface{}, error) {
		if len(args) < 6 {
			return nil, fmt.Errorf("Vector3Add requires (x1,y1,z1, x2,y2,z2)")
		}
		v1 := rl.Vector3{X: toFloat32(args[0]), Y: toFloat32(args[1]), Z: toFloat32(args[2])}
		v2 := rl.Vector3{X: toFloat32(args[3]), Y: toFloat32(args[4]), Z: toFloat32(args[5])}
		return vec3ToSlice(rl.Vector3Add(v1, v2)), nil
	})
	v.RegisterForeign("Vector3AddValue", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("Vector3AddValue requires (x, y, z, add)")
		}
		vec := rl.Vector3{X: toFloat32(args[0]), Y: toFloat32(args[1]), Z: toFloat32(args[2])}
		return vec3ToSlice(rl.Vector3AddValue(vec, toFloat32(args[3]))), nil
	})
	v.RegisterForeign("Vector3Subtract", func(args []interface{}) (interface{}, error) {
		if len(args) < 6 {
			return nil, fmt.Errorf("Vector3Subtract requires (x1,y1,z1, x2,y2,z2)")
		}
		v1 := rl.Vector3{X: toFloat32(args[0]), Y: toFloat32(args[1]), Z: toFloat32(args[2])}
		v2 := rl.Vector3{X: toFloat32(args[3]), Y: toFloat32(args[4]), Z: toFloat32(args[5])}
		return vec3ToSlice(rl.Vector3Subtract(v1, v2)), nil
	})
	v.RegisterForeign("Vector3SubtractValue", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("Vector3SubtractValue requires (x, y, z, sub)")
		}
		vec := rl.Vector3{X: toFloat32(args[0]), Y: toFloat32(args[1]), Z: toFloat32(args[2])}
		return vec3ToSlice(rl.Vector3SubtractValue(vec, toFloat32(args[3]))), nil
	})
	v.RegisterForeign("Vector3Scale", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("Vector3Scale requires (x, y, z, scalar)")
		}
		vec := rl.Vector3{X: toFloat32(args[0]), Y: toFloat32(args[1]), Z: toFloat32(args[2])}
		return vec3ToSlice(rl.Vector3Scale(vec, toFloat32(args[3]))), nil
	})
	v.RegisterForeign("Vector3Multiply", func(args []interface{}) (interface{}, error) {
		if len(args) < 6 {
			return nil, fmt.Errorf("Vector3Multiply requires (x1,y1,z1, x2,y2,z2)")
		}
		v1 := rl.Vector3{X: toFloat32(args[0]), Y: toFloat32(args[1]), Z: toFloat32(args[2])}
		v2 := rl.Vector3{X: toFloat32(args[3]), Y: toFloat32(args[4]), Z: toFloat32(args[5])}
		return vec3ToSlice(rl.Vector3Multiply(v1, v2)), nil
	})
	v.RegisterForeign("Vector3CrossProduct", func(args []interface{}) (interface{}, error) {
		if len(args) < 6 {
			return nil, fmt.Errorf("Vector3CrossProduct requires (x1,y1,z1, x2,y2,z2)")
		}
		v1 := rl.Vector3{X: toFloat32(args[0]), Y: toFloat32(args[1]), Z: toFloat32(args[2])}
		v2 := rl.Vector3{X: toFloat32(args[3]), Y: toFloat32(args[4]), Z: toFloat32(args[5])}
		return vec3ToSlice(rl.Vector3CrossProduct(v1, v2)), nil
	})
	v.RegisterForeign("Vector3Perpendicular", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("Vector3Perpendicular requires (x, y, z)")
		}
		vec := rl.Vector3{X: toFloat32(args[0]), Y: toFloat32(args[1]), Z: toFloat32(args[2])}
		return vec3ToSlice(rl.Vector3Perpendicular(vec)), nil
	})
	v.RegisterForeign("Vector3Length", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("Vector3Length requires (x, y, z)")
		}
		vec := rl.Vector3{X: toFloat32(args[0]), Y: toFloat32(args[1]), Z: toFloat32(args[2])}
		return float64(rl.Vector3Length(vec)), nil
	})
	v.RegisterForeign("Vector3LengthSqr", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("Vector3LengthSqr requires (x, y, z)")
		}
		vec := rl.Vector3{X: toFloat32(args[0]), Y: toFloat32(args[1]), Z: toFloat32(args[2])}
		return float64(rl.Vector3LengthSqr(vec)), nil
	})
	v.RegisterForeign("Vector3DotProduct", func(args []interface{}) (interface{}, error) {
		if len(args) < 6 {
			return nil, fmt.Errorf("Vector3DotProduct requires (x1,y1,z1, x2,y2,z2)")
		}
		v1 := rl.Vector3{X: toFloat32(args[0]), Y: toFloat32(args[1]), Z: toFloat32(args[2])}
		v2 := rl.Vector3{X: toFloat32(args[3]), Y: toFloat32(args[4]), Z: toFloat32(args[5])}
		return float64(rl.Vector3DotProduct(v1, v2)), nil
	})
	v.RegisterForeign("Vector3Distance", func(args []interface{}) (interface{}, error) {
		if len(args) < 6 {
			return nil, fmt.Errorf("Vector3Distance requires (x1,y1,z1, x2,y2,z2)")
		}
		v1 := rl.Vector3{X: toFloat32(args[0]), Y: toFloat32(args[1]), Z: toFloat32(args[2])}
		v2 := rl.Vector3{X: toFloat32(args[3]), Y: toFloat32(args[4]), Z: toFloat32(args[5])}
		return float64(rl.Vector3Distance(v1, v2)), nil
	})
	v.RegisterForeign("Vector3DistanceSqr", func(args []interface{}) (interface{}, error) {
		if len(args) < 6 {
			return nil, fmt.Errorf("Vector3DistanceSqr requires (x1,y1,z1, x2,y2,z2)")
		}
		v1 := rl.Vector3{X: toFloat32(args[0]), Y: toFloat32(args[1]), Z: toFloat32(args[2])}
		v2 := rl.Vector3{X: toFloat32(args[3]), Y: toFloat32(args[4]), Z: toFloat32(args[5])}
		return float64(rl.Vector3DistanceSqr(v1, v2)), nil
	})
	v.RegisterForeign("Vector3Angle", func(args []interface{}) (interface{}, error) {
		if len(args) < 6 {
			return nil, fmt.Errorf("Vector3Angle requires (x1,y1,z1, x2,y2,z2)")
		}
		v1 := rl.Vector3{X: toFloat32(args[0]), Y: toFloat32(args[1]), Z: toFloat32(args[2])}
		v2 := rl.Vector3{X: toFloat32(args[3]), Y: toFloat32(args[4]), Z: toFloat32(args[5])}
		return float64(rl.Vector3Angle(v1, v2)), nil
	})
	v.RegisterForeign("Vector3Negate", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("Vector3Negate requires (x, y, z)")
		}
		vec := rl.Vector3{X: toFloat32(args[0]), Y: toFloat32(args[1]), Z: toFloat32(args[2])}
		return vec3ToSlice(rl.Vector3Negate(vec)), nil
	})
	v.RegisterForeign("Vector3Divide", func(args []interface{}) (interface{}, error) {
		if len(args) < 6 {
			return nil, fmt.Errorf("Vector3Divide requires (x1,y1,z1, x2,y2,z2)")
		}
		v1 := rl.Vector3{X: toFloat32(args[0]), Y: toFloat32(args[1]), Z: toFloat32(args[2])}
		v2 := rl.Vector3{X: toFloat32(args[3]), Y: toFloat32(args[4]), Z: toFloat32(args[5])}
		return vec3ToSlice(rl.Vector3Divide(v1, v2)), nil
	})
	v.RegisterForeign("Vector3Normalize", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("Vector3Normalize requires (x, y, z)")
		}
		vec := rl.Vector3{X: toFloat32(args[0]), Y: toFloat32(args[1]), Z: toFloat32(args[2])}
		return vec3ToSlice(rl.Vector3Normalize(vec)), nil
	})
	v.RegisterForeign("Vector3OrthoNormalize", func(args []interface{}) (interface{}, error) {
		if len(args) < 6 {
			return nil, fmt.Errorf("Vector3OrthoNormalize requires (v1x,v1y,v1z, v2x,v2y,v2z)")
		}
		v1 := rl.Vector3{X: toFloat32(args[0]), Y: toFloat32(args[1]), Z: toFloat32(args[2])}
		v2 := rl.Vector3{X: toFloat32(args[3]), Y: toFloat32(args[4]), Z: toFloat32(args[5])}
		rl.Vector3OrthoNormalize(&v1, &v2)
		return []interface{}{float64(v1.X), float64(v1.Y), float64(v1.Z), float64(v2.X), float64(v2.Y), float64(v2.Z)}, nil
	})
	v.RegisterForeign("Vector3Transform", func(args []interface{}) (interface{}, error) {
		if len(args) < 3+16 {
			return nil, fmt.Errorf("Vector3Transform requires (x,y,z, 16 matrix floats)")
		}
		vec := rl.Vector3{X: toFloat32(args[0]), Y: toFloat32(args[1]), Z: toFloat32(args[2])}
		mat := argsToMatrix(args, 3)
		return vec3ToSlice(rl.Vector3Transform(vec, mat)), nil
	})
	v.RegisterForeign("Vector3RotateByQuaternion", func(args []interface{}) (interface{}, error) {
		if len(args) < 7 {
			return nil, fmt.Errorf("Vector3RotateByQuaternion requires (x,y,z, qx,qy,qz,qw)")
		}
		vec := rl.Vector3{X: toFloat32(args[0]), Y: toFloat32(args[1]), Z: toFloat32(args[2])}
		q := rl.Quaternion{X: toFloat32(args[3]), Y: toFloat32(args[4]), Z: toFloat32(args[5]), W: toFloat32(args[6])}
		return vec3ToSlice(rl.Vector3RotateByQuaternion(vec, q)), nil
	})
	v.RegisterForeign("Vector3RotateByAxisAngle", func(args []interface{}) (interface{}, error) {
		if len(args) < 7 {
			return nil, fmt.Errorf("Vector3RotateByAxisAngle requires (vx,vy,vz, axisX,axisY,axisZ, angle)")
		}
		vec := rl.Vector3{X: toFloat32(args[0]), Y: toFloat32(args[1]), Z: toFloat32(args[2])}
		axis := rl.Vector3{X: toFloat32(args[3]), Y: toFloat32(args[4]), Z: toFloat32(args[5])}
		return vec3ToSlice(rl.Vector3RotateByAxisAngle(vec, axis, toFloat32(args[6]))), nil
	})
	v.RegisterForeign("Vector3Lerp", func(args []interface{}) (interface{}, error) {
		if len(args) < 7 {
			return nil, fmt.Errorf("Vector3Lerp requires (x1,y1,z1, x2,y2,z2, amount)")
		}
		v1 := rl.Vector3{X: toFloat32(args[0]), Y: toFloat32(args[1]), Z: toFloat32(args[2])}
		v2 := rl.Vector3{X: toFloat32(args[3]), Y: toFloat32(args[4]), Z: toFloat32(args[5])}
		return vec3ToSlice(rl.Vector3Lerp(v1, v2, toFloat32(args[6]))), nil
	})
	v.RegisterForeign("Vector3Reflect", func(args []interface{}) (interface{}, error) {
		if len(args) < 6 {
			return nil, fmt.Errorf("Vector3Reflect requires (vx,vy,vz, normalX,normalY,normalZ)")
		}
		vec := rl.Vector3{X: toFloat32(args[0]), Y: toFloat32(args[1]), Z: toFloat32(args[2])}
		normal := rl.Vector3{X: toFloat32(args[3]), Y: toFloat32(args[4]), Z: toFloat32(args[5])}
		return vec3ToSlice(rl.Vector3Reflect(vec, normal)), nil
	})
	v.RegisterForeign("Vector3Min", func(args []interface{}) (interface{}, error) {
		if len(args) < 6 {
			return nil, fmt.Errorf("Vector3Min requires (x1,y1,z1, x2,y2,z2)")
		}
		v1 := rl.Vector3{X: toFloat32(args[0]), Y: toFloat32(args[1]), Z: toFloat32(args[2])}
		v2 := rl.Vector3{X: toFloat32(args[3]), Y: toFloat32(args[4]), Z: toFloat32(args[5])}
		return vec3ToSlice(rl.Vector3Min(v1, v2)), nil
	})
	v.RegisterForeign("Vector3Max", func(args []interface{}) (interface{}, error) {
		if len(args) < 6 {
			return nil, fmt.Errorf("Vector3Max requires (x1,y1,z1, x2,y2,z2)")
		}
		v1 := rl.Vector3{X: toFloat32(args[0]), Y: toFloat32(args[1]), Z: toFloat32(args[2])}
		v2 := rl.Vector3{X: toFloat32(args[3]), Y: toFloat32(args[4]), Z: toFloat32(args[5])}
		return vec3ToSlice(rl.Vector3Max(v1, v2)), nil
	})
	v.RegisterForeign("Vector3Barycenter", func(args []interface{}) (interface{}, error) {
		if len(args) < 12 {
			return nil, fmt.Errorf("Vector3Barycenter requires (px,py,pz, ax,ay,az, bx,by,bz, cx,cy,cz)")
		}
		p := rl.Vector3{X: toFloat32(args[0]), Y: toFloat32(args[1]), Z: toFloat32(args[2])}
		a := rl.Vector3{X: toFloat32(args[3]), Y: toFloat32(args[4]), Z: toFloat32(args[5])}
		b := rl.Vector3{X: toFloat32(args[6]), Y: toFloat32(args[7]), Z: toFloat32(args[8])}
		c := rl.Vector3{X: toFloat32(args[9]), Y: toFloat32(args[10]), Z: toFloat32(args[11])}
		return vec3ToSlice(rl.Vector3Barycenter(p, a, b, c)), nil
	})
	v.RegisterForeign("Vector3Unproject", func(args []interface{}) (interface{}, error) {
		if len(args) < 3+16+16 {
			return nil, fmt.Errorf("Vector3Unproject requires (x,y,z, projectionMat16, viewMat16)")
		}
		source := rl.Vector3{X: toFloat32(args[0]), Y: toFloat32(args[1]), Z: toFloat32(args[2])}
		proj := argsToMatrix(args, 3)
		view := argsToMatrix(args, 3+16)
		return vec3ToSlice(rl.Vector3Unproject(source, proj, view)), nil
	})
	v.RegisterForeign("Vector3Invert", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("Vector3Invert requires (x, y, z)")
		}
		vec := rl.Vector3{X: toFloat32(args[0]), Y: toFloat32(args[1]), Z: toFloat32(args[2])}
		return vec3ToSlice(rl.Vector3Invert(vec)), nil
	})
	v.RegisterForeign("Vector3Clamp", func(args []interface{}) (interface{}, error) {
		if len(args) < 9 {
			return nil, fmt.Errorf("Vector3Clamp requires (vx,vy,vz, minX,minY,minZ, maxX,maxY,maxZ)")
		}
		vec := rl.Vector3{X: toFloat32(args[0]), Y: toFloat32(args[1]), Z: toFloat32(args[2])}
		minV := rl.Vector3{X: toFloat32(args[3]), Y: toFloat32(args[4]), Z: toFloat32(args[5])}
		maxV := rl.Vector3{X: toFloat32(args[6]), Y: toFloat32(args[7]), Z: toFloat32(args[8])}
		return vec3ToSlice(rl.Vector3Clamp(vec, minV, maxV)), nil
	})
	v.RegisterForeign("Vector3ClampValue", func(args []interface{}) (interface{}, error) {
		if len(args) < 5 {
			return nil, fmt.Errorf("Vector3ClampValue requires (x, y, z, min, max)")
		}
		vec := rl.Vector3{X: toFloat32(args[0]), Y: toFloat32(args[1]), Z: toFloat32(args[2])}
		return vec3ToSlice(rl.Vector3ClampValue(vec, toFloat32(args[3]), toFloat32(args[4]))), nil
	})
	v.RegisterForeign("Vector3Equals", func(args []interface{}) (interface{}, error) {
		if len(args) < 6 {
			return false, nil
		}
		p := rl.Vector3{X: toFloat32(args[0]), Y: toFloat32(args[1]), Z: toFloat32(args[2])}
		q := rl.Vector3{X: toFloat32(args[3]), Y: toFloat32(args[4]), Z: toFloat32(args[5])}
		return rl.Vector3Equals(p, q), nil
	})
	v.RegisterForeign("Vector3Refract", func(args []interface{}) (interface{}, error) {
		if len(args) < 7 {
			return nil, fmt.Errorf("Vector3Refract requires (vx,vy,vz, nx,ny,nz, r)")
		}
		vec := rl.Vector3{X: toFloat32(args[0]), Y: toFloat32(args[1]), Z: toFloat32(args[2])}
		n := rl.Vector3{X: toFloat32(args[3]), Y: toFloat32(args[4]), Z: toFloat32(args[5])}
		return vec3ToSlice(rl.Vector3Refract(vec, n, toFloat32(args[6]))), nil
	})
	v.RegisterForeign("Vector3ToFloatV", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("Vector3ToFloatV requires (x, y, z)")
		}
		vec := rl.Vector3{X: toFloat32(args[0]), Y: toFloat32(args[1]), Z: toFloat32(args[2])}
		arr := rl.Vector3ToFloatV(vec)
		return []interface{}{float64(arr[0]), float64(arr[1]), float64(arr[2])}, nil
	})

	// Matrix math
	v.RegisterForeign("MatrixDeterminant", func(args []interface{}) (interface{}, error) {
		if len(args) < 16 {
			return nil, fmt.Errorf("MatrixDeterminant requires (16 matrix floats)")
		}
		return float64(rl.MatrixDeterminant(argsToMatrix(args, 0))), nil
	})
	v.RegisterForeign("MatrixTrace", func(args []interface{}) (interface{}, error) {
		if len(args) < 16 {
			return nil, fmt.Errorf("MatrixTrace requires (16 matrix floats)")
		}
		return float64(rl.MatrixTrace(argsToMatrix(args, 0))), nil
	})
	v.RegisterForeign("MatrixTranspose", func(args []interface{}) (interface{}, error) {
		if len(args) < 16 {
			return nil, fmt.Errorf("MatrixTranspose requires (16 matrix floats)")
		}
		return matrixToSlice(rl.MatrixTranspose(argsToMatrix(args, 0))), nil
	})
	v.RegisterForeign("MatrixInvert", func(args []interface{}) (interface{}, error) {
		if len(args) < 16 {
			return nil, fmt.Errorf("MatrixInvert requires (16 matrix floats)")
		}
		return matrixToSlice(rl.MatrixInvert(argsToMatrix(args, 0))), nil
	})
	v.RegisterForeign("MatrixIdentity", func(args []interface{}) (interface{}, error) {
		return matrixToSlice(rl.MatrixIdentity()), nil
	})
	v.RegisterForeign("MatrixAdd", func(args []interface{}) (interface{}, error) {
		if len(args) < 32 {
			return nil, fmt.Errorf("MatrixAdd requires (left16, right16)")
		}
		left := argsToMatrix(args, 0)
		right := argsToMatrix(args, 16)
		return matrixToSlice(rl.MatrixAdd(left, right)), nil
	})
	v.RegisterForeign("MatrixSubtract", func(args []interface{}) (interface{}, error) {
		if len(args) < 32 {
			return nil, fmt.Errorf("MatrixSubtract requires (left16, right16)")
		}
		left := argsToMatrix(args, 0)
		right := argsToMatrix(args, 16)
		return matrixToSlice(rl.MatrixSubtract(left, right)), nil
	})
	v.RegisterForeign("MatrixMultiply", func(args []interface{}) (interface{}, error) {
		if len(args) < 32 {
			return nil, fmt.Errorf("MatrixMultiply requires (left16, right16)")
		}
		left := argsToMatrix(args, 0)
		right := argsToMatrix(args, 16)
		return matrixToSlice(rl.MatrixMultiply(left, right)), nil
	})
	v.RegisterForeign("MatrixTranslate", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("MatrixTranslate requires (x, y, z)")
		}
		return matrixToSlice(rl.MatrixTranslate(toFloat32(args[0]), toFloat32(args[1]), toFloat32(args[2]))), nil
	})
	v.RegisterForeign("MatrixRotate", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("MatrixRotate requires (axisX, axisY, axisZ, angle)")
		}
		axis := rl.Vector3{X: toFloat32(args[0]), Y: toFloat32(args[1]), Z: toFloat32(args[2])}
		return matrixToSlice(rl.MatrixRotate(axis, toFloat32(args[3]))), nil
	})
	v.RegisterForeign("MatrixRotateX", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("MatrixRotateX requires (angle)")
		}
		return matrixToSlice(rl.MatrixRotateX(toFloat32(args[0]))), nil
	})
	v.RegisterForeign("MatrixRotateY", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("MatrixRotateY requires (angle)")
		}
		return matrixToSlice(rl.MatrixRotateY(toFloat32(args[0]))), nil
	})
	v.RegisterForeign("MatrixRotateZ", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("MatrixRotateZ requires (angle)")
		}
		return matrixToSlice(rl.MatrixRotateZ(toFloat32(args[0]))), nil
	})
	v.RegisterForeign("MatrixRotateXYZ", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("MatrixRotateXYZ requires (angleX, angleY, angleZ)")
		}
		ang := rl.Vector3{X: toFloat32(args[0]), Y: toFloat32(args[1]), Z: toFloat32(args[2])}
		return matrixToSlice(rl.MatrixRotateXYZ(ang)), nil
	})
	v.RegisterForeign("MatrixRotateZYX", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("MatrixRotateZYX requires (angleX, angleY, angleZ)")
		}
		ang := rl.Vector3{X: toFloat32(args[0]), Y: toFloat32(args[1]), Z: toFloat32(args[2])}
		return matrixToSlice(rl.MatrixRotateZYX(ang)), nil
	})
	v.RegisterForeign("MatrixScale", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("MatrixScale requires (x, y, z)")
		}
		return matrixToSlice(rl.MatrixScale(toFloat32(args[0]), toFloat32(args[1]), toFloat32(args[2]))), nil
	})
	v.RegisterForeign("MatrixFrustum", func(args []interface{}) (interface{}, error) {
		if len(args) < 6 {
			return nil, fmt.Errorf("MatrixFrustum requires (left, right, bottom, top, near, far)")
		}
		return matrixToSlice(rl.MatrixFrustum(toFloat32(args[0]), toFloat32(args[1]), toFloat32(args[2]), toFloat32(args[3]), toFloat32(args[4]), toFloat32(args[5]))), nil
	})
	v.RegisterForeign("MatrixPerspective", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("MatrixPerspective requires (fovy, aspect, near, far)")
		}
		return matrixToSlice(rl.MatrixPerspective(toFloat32(args[0]), toFloat32(args[1]), toFloat32(args[2]), toFloat32(args[3]))), nil
	})
	v.RegisterForeign("MatrixOrtho", func(args []interface{}) (interface{}, error) {
		if len(args) < 6 {
			return nil, fmt.Errorf("MatrixOrtho requires (left, right, bottom, top, near, far)")
		}
		return matrixToSlice(rl.MatrixOrtho(toFloat32(args[0]), toFloat32(args[1]), toFloat32(args[2]), toFloat32(args[3]), toFloat32(args[4]), toFloat32(args[5]))), nil
	})
	v.RegisterForeign("MatrixLookAt", func(args []interface{}) (interface{}, error) {
		if len(args) < 9 {
			return nil, fmt.Errorf("MatrixLookAt requires (eyeX,eyeY,eyeZ, targetX,targetY,targetZ, upX,upY,upZ)")
		}
		eye := rl.Vector3{X: toFloat32(args[0]), Y: toFloat32(args[1]), Z: toFloat32(args[2])}
		target := rl.Vector3{X: toFloat32(args[3]), Y: toFloat32(args[4]), Z: toFloat32(args[5])}
		up := rl.Vector3{X: toFloat32(args[6]), Y: toFloat32(args[7]), Z: toFloat32(args[8])}
		return matrixToSlice(rl.MatrixLookAt(eye, target, up)), nil
	})
	v.RegisterForeign("MatrixToFloatV", func(args []interface{}) (interface{}, error) {
		if len(args) < 16 {
			return nil, fmt.Errorf("MatrixToFloatV requires (16 matrix floats)")
		}
		arr := rl.MatrixToFloatV(argsToMatrix(args, 0))
		out := make([]interface{}, 16)
		for i := range arr {
			out[i] = float64(arr[i])
		}
		return out, nil
	})

	// Quaternion math
	v.RegisterForeign("QuaternionAdd", func(args []interface{}) (interface{}, error) {
		if len(args) < 8 {
			return nil, fmt.Errorf("QuaternionAdd requires (q1x,y,z,w, q2x,y,z,w)")
		}
		q1 := rl.Quaternion{X: toFloat32(args[0]), Y: toFloat32(args[1]), Z: toFloat32(args[2]), W: toFloat32(args[3])}
		q2 := rl.Quaternion{X: toFloat32(args[4]), Y: toFloat32(args[5]), Z: toFloat32(args[6]), W: toFloat32(args[7])}
		return quatToSlice(rl.QuaternionAdd(q1, q2)), nil
	})
	v.RegisterForeign("QuaternionAddValue", func(args []interface{}) (interface{}, error) {
		if len(args) < 5 {
			return nil, fmt.Errorf("QuaternionAddValue requires (qx,qy,qz,qw, add)")
		}
		q := rl.Quaternion{X: toFloat32(args[0]), Y: toFloat32(args[1]), Z: toFloat32(args[2]), W: toFloat32(args[3])}
		return quatToSlice(rl.QuaternionAddValue(q, toFloat32(args[4]))), nil
	})
	v.RegisterForeign("QuaternionSubtract", func(args []interface{}) (interface{}, error) {
		if len(args) < 8 {
			return nil, fmt.Errorf("QuaternionSubtract requires (q1x,y,z,w, q2x,y,z,w)")
		}
		q1 := rl.Quaternion{X: toFloat32(args[0]), Y: toFloat32(args[1]), Z: toFloat32(args[2]), W: toFloat32(args[3])}
		q2 := rl.Quaternion{X: toFloat32(args[4]), Y: toFloat32(args[5]), Z: toFloat32(args[6]), W: toFloat32(args[7])}
		return quatToSlice(rl.QuaternionSubtract(q1, q2)), nil
	})
	v.RegisterForeign("QuaternionSubtractValue", func(args []interface{}) (interface{}, error) {
		if len(args) < 5 {
			return nil, fmt.Errorf("QuaternionSubtractValue requires (qx,qy,qz,qw, sub)")
		}
		q := rl.Quaternion{X: toFloat32(args[0]), Y: toFloat32(args[1]), Z: toFloat32(args[2]), W: toFloat32(args[3])}
		return quatToSlice(rl.QuaternionSubtractValue(q, toFloat32(args[4]))), nil
	})
	v.RegisterForeign("QuaternionIdentity", func(args []interface{}) (interface{}, error) {
		return quatToSlice(rl.QuaternionIdentity()), nil
	})
	v.RegisterForeign("QuaternionLength", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("QuaternionLength requires (x, y, z, w)")
		}
		q := rl.Quaternion{X: toFloat32(args[0]), Y: toFloat32(args[1]), Z: toFloat32(args[2]), W: toFloat32(args[3])}
		return float64(rl.QuaternionLength(q)), nil
	})
	v.RegisterForeign("QuaternionNormalize", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("QuaternionNormalize requires (x, y, z, w)")
		}
		q := rl.Quaternion{X: toFloat32(args[0]), Y: toFloat32(args[1]), Z: toFloat32(args[2]), W: toFloat32(args[3])}
		return quatToSlice(rl.QuaternionNormalize(q)), nil
	})
	v.RegisterForeign("QuaternionInvert", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("QuaternionInvert requires (x, y, z, w)")
		}
		q := rl.Quaternion{X: toFloat32(args[0]), Y: toFloat32(args[1]), Z: toFloat32(args[2]), W: toFloat32(args[3])}
		return quatToSlice(rl.QuaternionInvert(q)), nil
	})
	v.RegisterForeign("QuaternionMultiply", func(args []interface{}) (interface{}, error) {
		if len(args) < 8 {
			return nil, fmt.Errorf("QuaternionMultiply requires (q1x,y,z,w, q2x,y,z,w)")
		}
		q1 := rl.Quaternion{X: toFloat32(args[0]), Y: toFloat32(args[1]), Z: toFloat32(args[2]), W: toFloat32(args[3])}
		q2 := rl.Quaternion{X: toFloat32(args[4]), Y: toFloat32(args[5]), Z: toFloat32(args[6]), W: toFloat32(args[7])}
		return quatToSlice(rl.QuaternionMultiply(q1, q2)), nil
	})
	v.RegisterForeign("QuaternionScale", func(args []interface{}) (interface{}, error) {
		if len(args) < 5 {
			return nil, fmt.Errorf("QuaternionScale requires (x, y, z, w, mul)")
		}
		q := rl.Quaternion{X: toFloat32(args[0]), Y: toFloat32(args[1]), Z: toFloat32(args[2]), W: toFloat32(args[3])}
		return quatToSlice(rl.QuaternionScale(q, toFloat32(args[4]))), nil
	})
	v.RegisterForeign("QuaternionDivide", func(args []interface{}) (interface{}, error) {
		if len(args) < 8 {
			return nil, fmt.Errorf("QuaternionDivide requires (q1x,y,z,w, q2x,y,z,w)")
		}
		q1 := rl.Quaternion{X: toFloat32(args[0]), Y: toFloat32(args[1]), Z: toFloat32(args[2]), W: toFloat32(args[3])}
		q2 := rl.Quaternion{X: toFloat32(args[4]), Y: toFloat32(args[5]), Z: toFloat32(args[6]), W: toFloat32(args[7])}
		return quatToSlice(rl.QuaternionDivide(q1, q2)), nil
	})
	v.RegisterForeign("QuaternionLerp", func(args []interface{}) (interface{}, error) {
		if len(args) < 9 {
			return nil, fmt.Errorf("QuaternionLerp requires (q1x,y,z,w, q2x,y,z,w, amount)")
		}
		q1 := rl.Quaternion{X: toFloat32(args[0]), Y: toFloat32(args[1]), Z: toFloat32(args[2]), W: toFloat32(args[3])}
		q2 := rl.Quaternion{X: toFloat32(args[4]), Y: toFloat32(args[5]), Z: toFloat32(args[6]), W: toFloat32(args[7])}
		return quatToSlice(rl.QuaternionLerp(q1, q2, toFloat32(args[8]))), nil
	})
	v.RegisterForeign("QuaternionNlerp", func(args []interface{}) (interface{}, error) {
		if len(args) < 9 {
			return nil, fmt.Errorf("QuaternionNlerp requires (q1x,y,z,w, q2x,y,z,w, amount)")
		}
		q1 := rl.Quaternion{X: toFloat32(args[0]), Y: toFloat32(args[1]), Z: toFloat32(args[2]), W: toFloat32(args[3])}
		q2 := rl.Quaternion{X: toFloat32(args[4]), Y: toFloat32(args[5]), Z: toFloat32(args[6]), W: toFloat32(args[7])}
		return quatToSlice(rl.QuaternionNlerp(q1, q2, toFloat32(args[8]))), nil
	})
	v.RegisterForeign("QuaternionSlerp", func(args []interface{}) (interface{}, error) {
		if len(args) < 9 {
			return nil, fmt.Errorf("QuaternionSlerp requires (q1x,y,z,w, q2x,y,z,w, amount)")
		}
		q1 := rl.Quaternion{X: toFloat32(args[0]), Y: toFloat32(args[1]), Z: toFloat32(args[2]), W: toFloat32(args[3])}
		q2 := rl.Quaternion{X: toFloat32(args[4]), Y: toFloat32(args[5]), Z: toFloat32(args[6]), W: toFloat32(args[7])}
		return quatToSlice(rl.QuaternionSlerp(q1, q2, toFloat32(args[8]))), nil
	})
	v.RegisterForeign("QuaternionFromVector3ToVector3", func(args []interface{}) (interface{}, error) {
		if len(args) < 6 {
			return nil, fmt.Errorf("QuaternionFromVector3ToVector3 requires (fromX,Y,Z, toX,Y,Z)")
		}
		from := rl.Vector3{X: toFloat32(args[0]), Y: toFloat32(args[1]), Z: toFloat32(args[2])}
		to := rl.Vector3{X: toFloat32(args[3]), Y: toFloat32(args[4]), Z: toFloat32(args[5])}
		return quatToSlice(rl.QuaternionFromVector3ToVector3(from, to)), nil
	})
	v.RegisterForeign("QuaternionFromMatrix", func(args []interface{}) (interface{}, error) {
		if len(args) < 16 {
			return nil, fmt.Errorf("QuaternionFromMatrix requires (16 matrix floats)")
		}
		return quatToSlice(rl.QuaternionFromMatrix(argsToMatrix(args, 0))), nil
	})
	v.RegisterForeign("QuaternionToMatrix", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("QuaternionToMatrix requires (x, y, z, w)")
		}
		q := rl.Quaternion{X: toFloat32(args[0]), Y: toFloat32(args[1]), Z: toFloat32(args[2]), W: toFloat32(args[3])}
		return matrixToSlice(rl.QuaternionToMatrix(q)), nil
	})
	v.RegisterForeign("QuaternionFromAxisAngle", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("QuaternionFromAxisAngle requires (axisX, axisY, axisZ, angle)")
		}
		axis := rl.Vector3{X: toFloat32(args[0]), Y: toFloat32(args[1]), Z: toFloat32(args[2])}
		return quatToSlice(rl.QuaternionFromAxisAngle(axis, toFloat32(args[3]))), nil
	})
	v.RegisterForeign("QuaternionToAxisAngle", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("QuaternionToAxisAngle requires (qx, qy, qz, qw)")
		}
		q := rl.Quaternion{X: toFloat32(args[0]), Y: toFloat32(args[1]), Z: toFloat32(args[2]), W: toFloat32(args[3])}
		var outAxis rl.Vector3
		var outAngle float32
		rl.QuaternionToAxisAngle(q, &outAxis, &outAngle)
		return []interface{}{float64(outAxis.X), float64(outAxis.Y), float64(outAxis.Z), float64(outAngle)}, nil
	})
	v.RegisterForeign("QuaternionFromEuler", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("QuaternionFromEuler requires (pitch, yaw, roll)")
		}
		return quatToSlice(rl.QuaternionFromEuler(toFloat32(args[0]), toFloat32(args[1]), toFloat32(args[2]))), nil
	})
	v.RegisterForeign("QuaternionToEuler", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("QuaternionToEuler requires (x, y, z, w)")
		}
		q := rl.Quaternion{X: toFloat32(args[0]), Y: toFloat32(args[1]), Z: toFloat32(args[2]), W: toFloat32(args[3])}
		return vec3ToSlice(rl.QuaternionToEuler(q)), nil
	})
	v.RegisterForeign("QuaternionTransform", func(args []interface{}) (interface{}, error) {
		if len(args) < 4+16 {
			return nil, fmt.Errorf("QuaternionTransform requires (qx,qy,qz,qw, 16 matrix floats)")
		}
		q := rl.Quaternion{X: toFloat32(args[0]), Y: toFloat32(args[1]), Z: toFloat32(args[2]), W: toFloat32(args[3])}
		mat := argsToMatrix(args, 4)
		return quatToSlice(rl.QuaternionTransform(q, mat)), nil
	})
	v.RegisterForeign("QuaternionEquals", func(args []interface{}) (interface{}, error) {
		if len(args) < 8 {
			return false, nil
		}
		p := rl.Quaternion{X: toFloat32(args[0]), Y: toFloat32(args[1]), Z: toFloat32(args[2]), W: toFloat32(args[3])}
		q := rl.Quaternion{X: toFloat32(args[4]), Y: toFloat32(args[5]), Z: toFloat32(args[6]), W: toFloat32(args[7])}
		return rl.QuaternionEquals(p, q), nil
	})
}
