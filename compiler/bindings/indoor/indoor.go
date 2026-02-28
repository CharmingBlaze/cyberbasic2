// Package indoor provides Room, Portal, Door, Trigger, Interactable stubs for indoor gameplay.
package indoor

import (
	"fmt"

	"cyberbasic/compiler/vm"
)

func toString(v interface{}) string {
	if v == nil {
		return ""
	}
	return fmt.Sprint(v)
}

func toFloat64(v interface{}) float64 {
	switch x := v.(type) {
	case int:
		return float64(x)
	case float64:
		return x
	case float32:
		return float64(x)
	default:
		return 0
	}
}

// RegisterIndoor registers Room, Portal, Door, Trigger, Interactable commands (stubs).
func RegisterIndoor(v *vm.VM) {
	v.RegisterForeign("RoomCreate", func(args []interface{}) (interface{}, error) {
		return fmt.Sprintf("room_%d", 0), nil
	})
	v.RegisterForeign("RoomSetBounds", func(args []interface{}) (interface{}, error) { return nil, nil })
	v.RegisterForeign("RoomAddPortal", func(args []interface{}) (interface{}, error) { return nil, nil })
	v.RegisterForeign("PortalCreate", func(args []interface{}) (interface{}, error) {
		return fmt.Sprintf("portal_%d", 0), nil
	})
	v.RegisterForeign("PortalSetOpen", func(args []interface{}) (interface{}, error) { return nil, nil })
	v.RegisterForeign("DoorCreate", func(args []interface{}) (interface{}, error) {
		return fmt.Sprintf("door_%d", 0), nil
	})
	v.RegisterForeign("DoorSetOpen", func(args []interface{}) (interface{}, error) { return nil, nil })
	v.RegisterForeign("DoorToggle", func(args []interface{}) (interface{}, error) { return nil, nil })
	v.RegisterForeign("DoorSetLocked", func(args []interface{}) (interface{}, error) { return nil, nil })
	v.RegisterForeign("LeverCreate", func(args []interface{}) (interface{}, error) {
		return fmt.Sprintf("lever_%d", 0), nil
	})
	v.RegisterForeign("ButtonCreate", func(args []interface{}) (interface{}, error) {
		return fmt.Sprintf("button_%d", 0), nil
	})
	v.RegisterForeign("SwitchCreate", func(args []interface{}) (interface{}, error) {
		return fmt.Sprintf("switch_%d", 0), nil
	})
	v.RegisterForeign("TriggerCreate", func(args []interface{}) (interface{}, error) {
		return fmt.Sprintf("trigger_%d", 0), nil
	})
	v.RegisterForeign("InteractableCreate", func(args []interface{}) (interface{}, error) {
		return fmt.Sprintf("interact_%d", 0), nil
	})
	v.RegisterForeign("PickupCreate", func(args []interface{}) (interface{}, error) {
		return fmt.Sprintf("pickup_%d", 0), nil
	})
	v.RegisterForeign("LightZoneCreate", func(args []interface{}) (interface{}, error) {
		return fmt.Sprintf("lightzone_%d", 0), nil
	})
	v.RegisterForeign("WorldSaveInteractables", func(args []interface{}) (interface{}, error) { return nil, nil })
	v.RegisterForeign("WorldLoadInteractables", func(args []interface{}) (interface{}, error) { return nil, nil })
}
