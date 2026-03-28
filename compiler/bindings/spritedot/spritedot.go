// Package spritedot exposes sprite.load using raylib LoadSprite (texture id).
package spritedot

import (
	"cyberbasic/compiler/vm"
	"fmt"
	"strings"
)

// RegisterSpriteDot registers global "sprite".
func RegisterSpriteDot(v *vm.VM) {
	v.SetGlobal("sprite", &spriteModuleDot{v: v})
}

type spriteModuleDot struct {
	v *vm.VM
}

func (s *spriteModuleDot) GetProp([]string) (vm.Value, error) { return nil, nil }
func (s *spriteModuleDot) SetProp([]string, vm.Value) error {
	return fmt.Errorf("sprite: namespace is not assignable")
}

func (s *spriteModuleDot) CallMethod(name string, args []vm.Value) (vm.Value, error) {
	ia := make([]interface{}, len(args))
	for i := range args {
		ia[i] = args[i]
	}
	switch strings.ToLower(name) {
	case "load":
		if len(args) < 1 {
			return nil, fmt.Errorf("sprite.load requires (path$)")
		}
		r, err := s.v.CallForeign("LoadSprite", ia)
		if err != nil {
			return nil, err
		}
		id := fmt.Sprint(r)
		return &SpriteDot{v: s.v, id: id}, nil
	default:
		return nil, fmt.Errorf("sprite: unknown method %q (load)", name)
	}
}

// SpriteDot wraps a texture id for 2D sprite-style drawing.
type SpriteDot struct {
	v  *vm.VM
	id string
}

func (s *SpriteDot) GetProp(path []string) (vm.Value, error) {
	if len(path) != 1 {
		return nil, fmt.Errorf("sprite handle: single property")
	}
	if strings.ToLower(path[0]) == "id" {
		return s.id, nil
	}
	return nil, nil
}

func (s *SpriteDot) SetProp([]string, vm.Value) error {
	return fmt.Errorf("sprite handle: use position via game/db commands or flat SetObject for scene sprites")
}

func (s *SpriteDot) CallMethod(name string, args []vm.Value) (vm.Value, error) {
	switch strings.ToLower(name) {
	case "draw":
		ia := append([]interface{}{s.id}, valuesToIface(args)...)
		return s.v.CallForeign("DrawTexture", ia)
	case "unload":
		_, err := s.v.CallForeign("UnloadTexture", []interface{}{s.id})
		return nil, err
	case "delete":
		_, err := s.v.CallForeign("UnloadTexture", []interface{}{s.id})
		return nil, err
	default:
		return nil, fmt.Errorf("sprite: use draw, unload, delete")
	}
}

func valuesToIface(a []vm.Value) []interface{} {
	out := make([]interface{}, len(a))
	for i := range a {
		out[i] = a[i]
	}
	return out
}
