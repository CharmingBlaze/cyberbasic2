// Package raylibdot registers dot namespaces (model, shapes3d, mesh, image, font, rlaudio)
// that forward to existing raylib RegisterForeign names. Requires raylib bindings registered first.
package raylibdot

import (
	"strings"

	"cyberbasic/compiler/bindings/modfacade"
	"cyberbasic/compiler/vm"
)

func lowerMap(names []string) map[string]string {
	m := make(map[string]string, len(names))
	for _, n := range names {
		m[strings.ToLower(n)] = n
	}
	return m
}

// Register installs globals: model, shapes3d, mesh, image, font, rlaudio.
func Register(v *vm.VM) {
	v.SetGlobal("model", modfacade.New(v, lowerMap(modelNames)))
	v.SetGlobal("shapes3d", modfacade.New(v, lowerMap(shapes3dNames)))
	v.SetGlobal("mesh", modfacade.New(v, lowerMap(meshNames)))
	v.SetGlobal("image", modfacade.New(v, lowerMap(imageNames)))
	v.SetGlobal("font", modfacade.New(v, lowerMap(fontNames)))
	v.SetGlobal("rlaudio", modfacade.New(v, lowerMap(rlaudioNames)))
}
