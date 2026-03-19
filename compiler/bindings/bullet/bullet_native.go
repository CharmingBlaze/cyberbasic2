//go:build bullet

// Package bullet: native Bullet backend stub. When built with -tags bullet, BulletNativeAvailable returns 1.
// A full native implementation requires CGO and Bullet C libraries (-lBulletDynamics -lBulletCollision -lLinearMath).
// This stub allows the build to succeed; replace with actual CGO bindings to Bullet for full fidelity.
package bullet

func init() {
	bulletNativeAvailable = true
}
