// Package math provides vector/matrix helpers for animation, physics, rendering.
package math

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

// Vec3 creates a Vector3.
func Vec3(x, y, z float32) rl.Vector3 {
	return rl.Vector3{X: x, Y: y, Z: z}
}

// Dot returns dot product of a and b.
func Dot(a, b rl.Vector3) float32 {
	return rl.Vector3DotProduct(a, b)
}

// Cross returns cross product of a and b.
func Cross(a, b rl.Vector3) rl.Vector3 {
	return rl.Vector3CrossProduct(a, b)
}

// Normalize returns normalized vector.
func Normalize(v rl.Vector3) rl.Vector3 {
	return rl.Vector3Normalize(v)
}
