// Package texturedot provides texture.load → TextureDot (draw, unload) over raylib texture ids.
package texturedot

import (
	"cyberbasic/compiler/vm"
	"fmt"
	"strings"
)

// RegisterTextureDot registers global "texture".
func RegisterTextureDot(v *vm.VM) {
	v.SetGlobal("texture", &textureModuleDot{v: v})
}

type textureModuleDot struct {
	v *vm.VM
}

func (t *textureModuleDot) GetProp([]string) (vm.Value, error) { return nil, nil }
func (t *textureModuleDot) SetProp([]string, vm.Value) error {
	return fmt.Errorf("texture: namespace is not assignable")
}

func (t *textureModuleDot) CallMethod(name string, args []vm.Value) (vm.Value, error) {
	ia := make([]interface{}, len(args))
	for i := range args {
		ia[i] = args[i]
	}
	switch strings.ToLower(name) {
	case "load":
		if len(args) < 1 {
			return nil, fmt.Errorf("texture.load requires (path$)")
		}
		r, err := t.v.CallForeign("LoadTexture", ia)
		if err != nil {
			return nil, err
		}
		id := fmt.Sprint(r)
		return &TextureDot{v: t.v, id: id}, nil
	default:
		return nil, fmt.Errorf("texture: unknown method %q (load)", name)
	}
}

// TextureDot is a vm.DotObject for a loaded texture id (string).
type TextureDot struct {
	v  *vm.VM
	id string
}

func (t *TextureDot) GetProp(path []string) (vm.Value, error) {
	if len(path) != 1 {
		return nil, fmt.Errorf("texture handle: single property")
	}
	if strings.ToLower(path[0]) == "id" {
		return t.id, nil
	}
	return nil, nil
}

func (t *TextureDot) SetProp([]string, vm.Value) error {
	return fmt.Errorf("texture handle: not assignable")
}

func (t *TextureDot) CallMethod(name string, args []vm.Value) (vm.Value, error) {
	switch strings.ToLower(name) {
	case "draw":
		ia := append([]interface{}{t.id}, valuesToIface(args)...)
		return t.v.CallForeign("DrawTexture", ia)
	case "drawex":
		ia := append([]interface{}{t.id}, valuesToIface(args)...)
		return t.v.CallForeign("DrawTextureEx", ia)
	case "unload":
		_, err := t.v.CallForeign("UnloadTexture", []interface{}{t.id})
		return nil, err
	default:
		return nil, fmt.Errorf("texture: use draw, drawex, unload")
	}
}

func valuesToIface(a []vm.Value) []interface{} {
	out := make([]interface{}, len(a))
	for i := range a {
		out[i] = a[i]
	}
	return out
}
