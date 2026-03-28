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

func (c *cameraRootDot) CallMethod(name string, args []vm.Value) (vm.Value, error) {
	ia := dotargs.From(args)
	switch strings.ToLower(name) {
	case "setcamera3d", "setcamera":
		return c.v.CallForeign("SetCamera3D", ia)
	case "beginmode3d", "begin3d":
		return c.v.CallForeign("BeginMode3D", ia)
	case "endmode3d", "end3d":
		return c.v.CallForeign("EndMode3D", ia)
	case "camera3d":
		return c.v.CallForeign("CAMERA3D", ia)
	default:
		return nil, fmt.Errorf("camera: use setcamera3d, beginmode3d, endmode3d, camera3d, or camera.fx.*")
	}
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
