// Package objectdot provides object.load → ObjectDot over DBP 3D object ids.
//
// ObjectDot methods call the same RegisterForeign handlers as flat DBP commands (id is prepended):
// DrawObject, DeleteObject, PositionObject, RotateObject, ScaleObject, MoveObject, TurnObject,
// YRotateObject, HideObject, ShowObject, CloneObject, CopyObject, ObjectExists, FixObject, UnfixObject,
// SetObjectColor, SetObjectAlpha, SetObjectTexture, SetObjectNormalmap, SetObjectRoughness,
// SetObjectMetallic, SetObjectEmissive, SetObjectShader, SetObjectWireframe, SetObjectCollision.
package objectdot

import (
	"cyberbasic/compiler/vm"
	"fmt"
	"strings"
	"sync"
)

var autoObjectMu sync.Mutex
var autoObjectID = 100000

func nextAutoObjectID() int {
	autoObjectMu.Lock()
	defer autoObjectMu.Unlock()
	autoObjectID++
	return autoObjectID
}

// RegisterObjectDot registers global "object".
func RegisterObjectDot(v *vm.VM) {
	v.SetGlobal("object", &objectModuleDot{v: v})
}

type objectModuleDot struct {
	v *vm.VM
}

func (o *objectModuleDot) GetProp([]string) (vm.Value, error) { return nil, nil }
func (o *objectModuleDot) SetProp([]string, vm.Value) error {
	return fmt.Errorf("object: namespace is not assignable")
}

func (o *objectModuleDot) CallMethod(name string, args []vm.Value) (vm.Value, error) {
	ia := make([]interface{}, len(args))
	for i := range args {
		ia[i] = args[i]
	}
	switch strings.ToLower(name) {
	case "load":
		if len(args) < 1 {
			return nil, fmt.Errorf("object.load requires (path$) or (id, path$)")
		}
		var id int
		var path string
		if len(args) >= 2 {
			id = toInt(args[0])
			path = fmt.Sprint(args[1])
			if _, err := o.v.CallForeign("LoadObjectId", []interface{}{id, path}); err != nil {
				return nil, err
			}
		} else {
			id = nextAutoObjectID()
			path = fmt.Sprint(args[0])
			if _, err := o.v.CallForeign("LoadObjectId", []interface{}{id, path}); err != nil {
				return nil, err
			}
		}
		return &ObjectDot{v: o.v, id: id}, nil
	default:
		return nil, fmt.Errorf("object: unknown method %q (load)", name)
	}
}

func toInt(v vm.Value) int {
	switch x := v.(type) {
	case int:
		return x
	case int32:
		return int(x)
	case int64:
		return int(x)
	case float64:
		return int(x)
	default:
		return 0
	}
}

// ObjectDot wraps a DBP integer object id.
type ObjectDot struct {
	v  *vm.VM
	id int
}

func (o *ObjectDot) prependArgs(tail []vm.Value) []interface{} {
	ia := make([]interface{}, 0, 1+len(tail))
	ia = append(ia, o.id)
	for _, a := range tail {
		ia = append(ia, a)
	}
	return ia
}

func (o *ObjectDot) GetProp(path []string) (vm.Value, error) {
	if len(path) == 0 {
		return nil, fmt.Errorf("empty path")
	}
	p := strings.ToLower(path[0])
	switch p {
	case "id":
		return float64(o.id), nil
	case "x":
		return callFloat(o.v, "GetObjectX", o.id)
	case "y":
		return callFloat(o.v, "GetObjectY", o.id)
	case "z":
		return callFloat(o.v, "GetObjectZ", o.id)
	case "pitch":
		return callFloat(o.v, "GetObjectPitch", o.id)
	case "yaw":
		return callFloat(o.v, "GetObjectYaw", o.id)
	case "roll":
		return callFloat(o.v, "GetObjectRoll", o.id)
	case "scalex":
		return callFloat(o.v, "GetObjectScaleX", o.id)
	case "scaley":
		return callFloat(o.v, "GetObjectScaleY", o.id)
	case "scalez":
		return callFloat(o.v, "GetObjectScaleZ", o.id)
	default:
		return nil, fmt.Errorf("object: unknown property %q (id, x, y, z, pitch, yaw, roll, scalex, scaley, scalez)", path[0])
	}
}

func callFloat(v *vm.VM, fn string, id int) (float64, error) {
	r, err := v.CallForeign(fn, []interface{}{id})
	if err != nil {
		return 0, err
	}
	switch x := r.(type) {
	case float64:
		return x, nil
	case float32:
		return float64(x), nil
	case int:
		return float64(x), nil
	default:
		return 0, nil
	}
}

func (o *ObjectDot) SetProp(path []string, val vm.Value) error {
	if len(path) != 1 {
		return fmt.Errorf("object: nested property set not supported")
	}
	p := strings.ToLower(path[0])
	fv := toFloat64(val)
	switch p {
	case "x", "y", "z":
		x, _ := callFloat(o.v, "GetObjectX", o.id)
		y, _ := callFloat(o.v, "GetObjectY", o.id)
		z, _ := callFloat(o.v, "GetObjectZ", o.id)
		switch p {
		case "x":
			x = fv
		case "y":
			y = fv
		case "z":
			z = fv
		}
		_, err := o.v.CallForeign("PositionObject", []interface{}{o.id, x, y, z})
		return err
	case "pitch", "yaw", "roll":
		pitch, _ := callFloat(o.v, "GetObjectPitch", o.id)
		yaw, _ := callFloat(o.v, "GetObjectYaw", o.id)
		roll, _ := callFloat(o.v, "GetObjectRoll", o.id)
		switch p {
		case "pitch":
			pitch = fv
		case "yaw":
			yaw = fv
		case "roll":
			roll = fv
		}
		_, err := o.v.CallForeign("RotateObject", []interface{}{o.id, pitch, yaw, roll})
		return err
	case "scalex", "scaley", "scalez":
		sx, _ := callFloat(o.v, "GetObjectScaleX", o.id)
		sy, _ := callFloat(o.v, "GetObjectScaleY", o.id)
		sz, _ := callFloat(o.v, "GetObjectScaleZ", o.id)
		switch p {
		case "scalex":
			sx = fv
		case "scaley":
			sy = fv
		case "scalez":
			sz = fv
		}
		_, err := o.v.CallForeign("ScaleObject", []interface{}{o.id, sx, sy, sz})
		return err
	default:
		return fmt.Errorf("object: unknown or read-only property %q", path[0])
	}
}

func toFloat64(v vm.Value) float64 {
	switch x := v.(type) {
	case float64:
		return x
	case int:
		return float64(x)
	case int32:
		return float64(x)
	default:
		return 0
	}
}

func (o *ObjectDot) CallMethod(name string, args []vm.Value) (vm.Value, error) {
	n := strings.ToLower(name)
	ia := o.prependArgs(args)

	switch n {
	case "draw":
		return o.v.CallForeign("DrawObject", []interface{}{o.id})
	case "delete":
		return o.v.CallForeign("DeleteObject", []interface{}{o.id})
	case "position":
		if len(args) < 3 {
			return nil, fmt.Errorf("position(x, y, z) requires 3 arguments; flat: PositionObject")
		}
		return o.v.CallForeign("PositionObject", []interface{}{o.id, toFloat64(args[0]), toFloat64(args[1]), toFloat64(args[2])})
	case "rotate":
		if len(args) < 3 {
			return nil, fmt.Errorf("rotate(pitch, yaw, roll) requires 3 arguments; flat: RotateObject")
		}
		return o.v.CallForeign("RotateObject", []interface{}{o.id, toFloat64(args[0]), toFloat64(args[1]), toFloat64(args[2])})
	case "scale":
		if len(args) < 3 {
			return nil, fmt.Errorf("scale(sx, sy, sz) requires 3 arguments; flat: ScaleObject")
		}
		return o.v.CallForeign("ScaleObject", []interface{}{o.id, toFloat64(args[0]), toFloat64(args[1]), toFloat64(args[2])})
	case "move", "moveobject":
		if len(args) < 3 {
			return nil, fmt.Errorf("move(dx, dy, dz) requires 3 arguments; flat: MoveObject")
		}
		return o.v.CallForeign("MoveObject", ia)
	case "turn", "turnobject":
		if len(args) < 3 {
			return nil, fmt.Errorf("turn(dpitch, dyaw, droll) requires 3 arguments; flat: TurnObject")
		}
		return o.v.CallForeign("TurnObject", ia)
	case "yrotate", "yrotateobject":
		if len(args) < 1 {
			return nil, fmt.Errorf("yrotate(angle) requires 1 argument; flat: YRotateObject")
		}
		return o.v.CallForeign("YRotateObject", ia)
	case "hide", "hideobject":
		return o.v.CallForeign("HideObject", []interface{}{o.id})
	case "show", "showobject":
		return o.v.CallForeign("ShowObject", []interface{}{o.id})
	case "clone", "cloneobject":
		if len(args) < 1 {
			return nil, fmt.Errorf("clone(newId) requires 1 argument; flat: CloneObject(newId, sourceId)")
		}
		return o.v.CallForeign("CloneObject", []interface{}{toInt(args[0]), o.id})
	case "copy", "copyobject":
		if len(args) < 1 {
			return nil, fmt.Errorf("copy(newId) requires 1 argument; flat: CopyObject (same as CloneObject)")
		}
		return o.v.CallForeign("CopyObject", []interface{}{toInt(args[0]), o.id})
	case "exists", "objectexists":
		return o.v.CallForeign("ObjectExists", []interface{}{o.id})
	case "fix", "fixobject":
		return o.v.CallForeign("FixObject", []interface{}{o.id})
	case "unfix", "unfixobject":
		return o.v.CallForeign("UnfixObject", []interface{}{o.id})
	case "setcolor", "setobjectcolor":
		if len(args) < 3 {
			return nil, fmt.Errorf("setcolor(r, g, b) requires 3 arguments; flat: SetObjectColor")
		}
		return o.v.CallForeign("SetObjectColor", ia)
	case "setalpha", "setobjectalpha":
		if len(args) < 1 {
			return nil, fmt.Errorf("setalpha(value) requires 1 argument; flat: SetObjectAlpha")
		}
		return o.v.CallForeign("SetObjectAlpha", ia)
	case "settexture", "setobjecttexture":
		if len(args) < 1 {
			return nil, fmt.Errorf("settexture(textureIdOrPath$) requires 1 argument; flat: SetObjectTexture")
		}
		return o.v.CallForeign("SetObjectTexture", ia)
	case "setnormalmap", "setobjectnormalmap":
		if len(args) < 1 {
			return nil, fmt.Errorf("setnormalmap(path$) requires 1 argument; flat: SetObjectNormalmap")
		}
		return o.v.CallForeign("SetObjectNormalmap", ia)
	case "setroughness", "setobjectroughness":
		if len(args) < 1 {
			return nil, fmt.Errorf("setroughness(value) requires 1 argument; flat: SetObjectRoughness")
		}
		return o.v.CallForeign("SetObjectRoughness", ia)
	case "setmetallic", "setobjectmetallic":
		if len(args) < 1 {
			return nil, fmt.Errorf("setmetallic(value) requires 1 argument; flat: SetObjectMetallic")
		}
		return o.v.CallForeign("SetObjectMetallic", ia)
	case "setemissive", "setobjectemissive":
		if len(args) < 3 {
			return nil, fmt.Errorf("setemissive(r, g, b) requires 3 arguments; flat: SetObjectEmissive")
		}
		return o.v.CallForeign("SetObjectEmissive", ia)
	case "setshader", "setobjectshader":
		if len(args) < 1 {
			return nil, fmt.Errorf("setshader(shaderId) requires 1 argument; flat: SetObjectShader")
		}
		return o.v.CallForeign("SetObjectShader", ia)
	case "setwireframe", "setobjectwireframe":
		if len(args) < 1 {
			return nil, fmt.Errorf("setwireframe(onOff) requires 1 argument; flat: SetObjectWireframe")
		}
		return o.v.CallForeign("SetObjectWireframe", ia)
	case "setcollision", "setobjectcollision":
		if len(args) < 1 {
			return nil, fmt.Errorf("setcollision(onOff) requires 1 argument; flat: SetObjectCollision")
		}
		return o.v.CallForeign("SetObjectCollision", ia)
	default:
		return nil, fmt.Errorf("object: unknown method %q (draw, delete, position, rotate, scale, move, turn, yrotate, hide, show, clone, copy, exists, fix, unfix, setcolor, setalpha, settexture, setnormalmap, setroughness, setmetallic, setemissive, setshader, setwireframe, setcollision)", name)
	}
}
