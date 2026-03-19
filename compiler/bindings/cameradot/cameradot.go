// Package cameradot exposes camera.fx.add / clear for post-FX queue (stubs wire to effect package foreigns).
package cameradot

import (
	"cyberbasic/compiler/vm"
	"fmt"
	"strings"

	"cyberbasic/compiler/bindings/dotargs"
)

// RegisterCameraDot installs global camera DotObject (nested camera.fx).
func RegisterCameraDot(v *vm.VM) {
	v.SetGlobal("camera", &cameraRootDot{v: v})
}

type cameraRootDot struct {
	v *vm.VM
}

func (c *cameraRootDot) GetProp(path []string) (vm.Value, error) {
	if len(path) == 1 && strings.ToLower(path[0]) == "fx" {
		return &cameraFXDot{v: c.v}, nil
	}
	return nil, fmt.Errorf("camera: unknown property %v (use fx)", path)
}

func (c *cameraRootDot) SetProp([]string, vm.Value) error {
	return fmt.Errorf("camera: namespace is not assignable")
}

func (c *cameraRootDot) CallMethod(string, []vm.Value) (vm.Value, error) {
	return nil, fmt.Errorf("camera: use camera.fx.add / camera.fx.clear")
}

type cameraFXDot struct {
	v *vm.VM
}

func (f *cameraFXDot) GetProp([]string) (vm.Value, error) { return nil, nil }
func (f *cameraFXDot) SetProp([]string, vm.Value) error {
	return fmt.Errorf("camera.fx: not assignable")
}

func (f *cameraFXDot) CallMethod(name string, args []vm.Value) (vm.Value, error) {
	ia := dotargs.From(args)
	switch strings.ToLower(name) {
	case "add":
		return f.v.CallForeign("CameraFXAddStub", ia)
	case "clear":
		return f.v.CallForeign("CameraFXClearStub", ia)
	default:
		return nil, fmt.Errorf("camera.fx: unknown method %q (add, clear)", name)
	}
}
