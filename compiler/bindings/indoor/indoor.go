// Package indoor provides Room, Portal, Door, Trigger, Interactable for indoor gameplay.
package indoor

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"cyberbasic/compiler/bindings/modfacade"
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

type room struct {
	MinX, MinY, MinZ, MaxX, MaxY, MaxZ float64
	Portals                             []string
}

type portal struct {
	RoomA, RoomB string
	PlaneX, PlaneY, PlaneZ, PlaneD float64
	Open                          bool
}

type door struct {
	Open   bool
	Locked bool
}

type trigger struct {
	MinX, MinY, MinZ, MaxX, MaxY, MaxZ float64
}

type interactable struct {
	X, Y, Z float64
	Data    map[string]interface{}
}

type pickup struct {
	X, Y, Z float64
	ItemId  string
	Amount  int
}

type lightZone struct {
	MinX, MinY, MinZ, MaxX, MaxY, MaxZ float64
	R, G, B, Intensity                 float64
}

type lever struct {
	On bool
}

type button struct {
	Pressed bool
}

type switchState struct {
	On bool
}

var (
	rooms       = make(map[string]*room)
	portals     = make(map[string]*portal)
	doors       = make(map[string]*door)
	triggers    = make(map[string]*trigger)
	interactables = make(map[string]*interactable)
	pickups     = make(map[string]*pickup)
	lightZones  = make(map[string]*lightZone)
	levers      = make(map[string]*lever)
	buttons     = make(map[string]*button)
	switches    = make(map[string]*switchState)
	roomSeq     int
	portalSeq   int
	doorSeq     int
	triggerSeq  int
	interactSeq int
	pickupSeq   int
	lightSeq    int
	leverSeq    int
	buttonSeq   int
	switchSeq   int
	indoorMu    sync.RWMutex
)

// RegisterIndoor registers Room, Portal, Door, Trigger, Interactable commands.
func RegisterIndoor(v *vm.VM) {
	v.RegisterForeign("RoomCreate", func(args []interface{}) (interface{}, error) {
		indoorMu.Lock()
		roomSeq++
		id := fmt.Sprintf("room_%d", roomSeq)
		rooms[id] = &room{Portals: []string{}}
		indoorMu.Unlock()
		return id, nil
	})
	v.RegisterForeign("RoomSetBounds", func(args []interface{}) (interface{}, error) {
		if len(args) < 7 {
			return nil, nil
		}
		id := toString(args[0])
		indoorMu.Lock()
		if r := rooms[id]; r != nil {
			r.MinX, r.MinY, r.MinZ = toFloat64(args[1]), toFloat64(args[2]), toFloat64(args[3])
			r.MaxX, r.MaxY, r.MaxZ = toFloat64(args[4]), toFloat64(args[5]), toFloat64(args[6])
		}
		indoorMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("RoomAddPortal", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, nil
		}
		roomId := toString(args[0])
		portalId := toString(args[1])
		indoorMu.Lock()
		if r := rooms[roomId]; r != nil {
			r.Portals = append(r.Portals, portalId)
		}
		indoorMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("PortalCreate", func(args []interface{}) (interface{}, error) {
		roomA, roomB := "", ""
		if len(args) >= 2 {
			roomA, roomB = toString(args[0]), toString(args[1])
		}
		indoorMu.Lock()
		portalSeq++
		id := fmt.Sprintf("portal_%d", portalSeq)
		portals[id] = &portal{RoomA: roomA, RoomB: roomB, Open: true}
		indoorMu.Unlock()
		return id, nil
	})
	v.RegisterForeign("PortalSetOpen", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, nil
		}
		id := toString(args[0])
		indoorMu.Lock()
		if p := portals[id]; p != nil {
			p.Open = toFloat64(args[1]) != 0
		}
		indoorMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("DoorCreate", func(args []interface{}) (interface{}, error) {
		indoorMu.Lock()
		doorSeq++
		id := fmt.Sprintf("door_%d", doorSeq)
		doors[id] = &door{Open: false, Locked: false}
		indoorMu.Unlock()
		return id, nil
	})
	v.RegisterForeign("DoorSetOpen", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, nil
		}
		id := toString(args[0])
		indoorMu.Lock()
		if d := doors[id]; d != nil && !d.Locked {
			d.Open = toFloat64(args[1]) != 0
		}
		indoorMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("DoorToggle", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, nil
		}
		id := toString(args[0])
		indoorMu.Lock()
		if d := doors[id]; d != nil && !d.Locked {
			d.Open = !d.Open
		}
		indoorMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("DoorSetLocked", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, nil
		}
		id := toString(args[0])
		indoorMu.Lock()
		if d := doors[id]; d != nil {
			d.Locked = toFloat64(args[1]) != 0
		}
		indoorMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("LeverCreate", func(args []interface{}) (interface{}, error) {
		indoorMu.Lock()
		leverSeq++
		id := fmt.Sprintf("lever_%d", leverSeq)
		levers[id] = &lever{On: false}
		indoorMu.Unlock()
		return id, nil
	})
	v.RegisterForeign("ButtonCreate", func(args []interface{}) (interface{}, error) {
		indoorMu.Lock()
		buttonSeq++
		id := fmt.Sprintf("button_%d", buttonSeq)
		buttons[id] = &button{Pressed: false}
		indoorMu.Unlock()
		return id, nil
	})
	v.RegisterForeign("SwitchCreate", func(args []interface{}) (interface{}, error) {
		indoorMu.Lock()
		switchSeq++
		id := fmt.Sprintf("switch_%d", switchSeq)
		switches[id] = &switchState{On: false}
		indoorMu.Unlock()
		return id, nil
	})
	v.RegisterForeign("TriggerCreate", func(args []interface{}) (interface{}, error) {
		minX, minY, minZ, maxX, maxY, maxZ := 0.0, 0.0, 0.0, 1.0, 1.0, 1.0
		if len(args) >= 6 {
			minX, minY, minZ = toFloat64(args[0]), toFloat64(args[1]), toFloat64(args[2])
			maxX, maxY, maxZ = toFloat64(args[3]), toFloat64(args[4]), toFloat64(args[5])
		}
		indoorMu.Lock()
		triggerSeq++
		id := fmt.Sprintf("trigger_%d", triggerSeq)
		triggers[id] = &trigger{MinX: minX, MinY: minY, MinZ: minZ, MaxX: maxX, MaxY: maxY, MaxZ: maxZ}
		indoorMu.Unlock()
		return id, nil
	})
	v.RegisterForeign("TriggerSetBounds", func(args []interface{}) (interface{}, error) {
		if len(args) < 7 {
			return nil, nil
		}
		id := toString(args[0])
		indoorMu.Lock()
		if t := triggers[id]; t != nil {
			t.MinX, t.MinY, t.MinZ = toFloat64(args[1]), toFloat64(args[2]), toFloat64(args[3])
			t.MaxX, t.MaxY, t.MaxZ = toFloat64(args[4]), toFloat64(args[5]), toFloat64(args[6])
		}
		indoorMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("InteractableCreate", func(args []interface{}) (interface{}, error) {
		x, y, z := 0.0, 0.0, 0.0
		if len(args) >= 3 {
			x, y, z = toFloat64(args[0]), toFloat64(args[1]), toFloat64(args[2])
		}
		indoorMu.Lock()
		interactSeq++
		id := fmt.Sprintf("interact_%d", interactSeq)
		interactables[id] = &interactable{X: x, Y: y, Z: z, Data: make(map[string]interface{})}
		indoorMu.Unlock()
		return id, nil
	})
	v.RegisterForeign("PickupCreate", func(args []interface{}) (interface{}, error) {
		x, y, z := 0.0, 0.0, 0.0
		itemId := ""
		amount := 1
		if len(args) >= 3 {
			x, y, z = toFloat64(args[0]), toFloat64(args[1]), toFloat64(args[2])
		}
		if len(args) >= 4 {
			itemId = toString(args[3])
		}
		if len(args) >= 5 {
			amount = int(toFloat64(args[4]))
		}
		indoorMu.Lock()
		pickupSeq++
		id := fmt.Sprintf("pickup_%d", pickupSeq)
		pickups[id] = &pickup{X: x, Y: y, Z: z, ItemId: itemId, Amount: amount}
		indoorMu.Unlock()
		return id, nil
	})
	v.RegisterForeign("LightZoneCreate", func(args []interface{}) (interface{}, error) {
		minX, minY, minZ, maxX, maxY, maxZ := 0.0, 0.0, 0.0, 10.0, 10.0, 10.0
		if len(args) >= 6 {
			minX, minY, minZ = toFloat64(args[0]), toFloat64(args[1]), toFloat64(args[2])
			maxX, maxY, maxZ = toFloat64(args[3]), toFloat64(args[4]), toFloat64(args[5])
		}
		indoorMu.Lock()
		lightSeq++
		id := fmt.Sprintf("lightzone_%d", lightSeq)
		lightZones[id] = &lightZone{MinX: minX, MinY: minY, MinZ: minZ, MaxX: maxX, MaxY: maxY, MaxZ: maxZ, R: 1, G: 1, B: 1, Intensity: 1}
		indoorMu.Unlock()
		return id, nil
	})

	// WorldSaveInteractables / WorldLoadInteractables
	type interactExport struct {
		Rooms   map[string]*room
		Portals map[string]*portal
		Doors   map[string]*door
		Triggers map[string]*trigger
		Interactables map[string]*interactable
		Pickups map[string]*pickup
		LightZones map[string]*lightZone
		Levers  map[string]*lever
		Buttons map[string]*button
		Switches map[string]*switchState
	}
	v.RegisterForeign("WorldSaveInteractables", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("WorldSaveInteractables requires (path)")
		}
		path := toString(args[0])
		indoorMu.RLock()
		exp := interactExport{
			Rooms:   copyRooms(rooms),
			Portals: copyPortals(portals),
			Doors:   copyDoors(doors),
			Triggers: copyTriggers(triggers),
			Interactables: copyInteractables(interactables),
			Pickups: copyPickups(pickups),
			LightZones: copyLightZones(lightZones),
			Levers:  copyLevers(levers),
			Buttons: copyButtons(buttons),
			Switches: copySwitches(switches),
		}
		indoorMu.RUnlock()
		raw, err := json.MarshalIndent(exp, "", "  ")
		if err != nil {
			return nil, err
		}
		return nil, os.WriteFile(path, raw, 0644)
	})
	v.RegisterForeign("WorldLoadInteractables", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("WorldLoadInteractables requires (path)")
		}
		path := toString(args[0])
		raw, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}
		var exp interactExport
		if err := json.Unmarshal(raw, &exp); err != nil {
			return nil, err
		}
		indoorMu.Lock()
		if exp.Rooms != nil {
			rooms = exp.Rooms
		}
		if exp.Portals != nil {
			portals = exp.Portals
		}
		if exp.Doors != nil {
			doors = exp.Doors
		}
		if exp.Triggers != nil {
			triggers = exp.Triggers
		}
		if exp.Interactables != nil {
			interactables = exp.Interactables
		}
		if exp.Pickups != nil {
			pickups = exp.Pickups
		}
		if exp.LightZones != nil {
			lightZones = exp.LightZones
		}
		if exp.Levers != nil {
			levers = exp.Levers
		}
		if exp.Buttons != nil {
			buttons = exp.Buttons
		}
		if exp.Switches != nil {
			switches = exp.Switches
		}
		indoorMu.Unlock()
		return nil, nil
	})

	v.SetGlobal("indoor", modfacade.New(v, indoorV2))
}

func copyRooms(m map[string]*room) map[string]*room {
	out := make(map[string]*room)
	for k, v := range m {
		if v != nil {
			portalsCopy := make([]string, len(v.Portals))
			copy(portalsCopy, v.Portals)
			out[k] = &room{MinX: v.MinX, MinY: v.MinY, MinZ: v.MinZ, MaxX: v.MaxX, MaxY: v.MaxY, MaxZ: v.MaxZ, Portals: portalsCopy}
		}
	}
	return out
}

func copyPortals(m map[string]*portal) map[string]*portal {
	out := make(map[string]*portal)
	for k, v := range m {
		if v != nil {
			out[k] = &portal{RoomA: v.RoomA, RoomB: v.RoomB, Open: v.Open}
		}
	}
	return out
}

func copyDoors(m map[string]*door) map[string]*door {
	out := make(map[string]*door)
	for k, v := range m {
		if v != nil {
			out[k] = &door{Open: v.Open, Locked: v.Locked}
		}
	}
	return out
}

func copyTriggers(m map[string]*trigger) map[string]*trigger {
	out := make(map[string]*trigger)
	for k, v := range m {
		if v != nil {
			out[k] = &trigger{MinX: v.MinX, MinY: v.MinY, MinZ: v.MinZ, MaxX: v.MaxX, MaxY: v.MaxY, MaxZ: v.MaxZ}
		}
	}
	return out
}

func copyInteractables(m map[string]*interactable) map[string]*interactable {
	out := make(map[string]*interactable)
	for k, v := range m {
		if v != nil {
			data := make(map[string]interface{})
			for dk, dv := range v.Data {
				data[dk] = dv
			}
			out[k] = &interactable{X: v.X, Y: v.Y, Z: v.Z, Data: data}
		}
	}
	return out
}

func copyPickups(m map[string]*pickup) map[string]*pickup {
	out := make(map[string]*pickup)
	for k, v := range m {
		if v != nil {
			out[k] = &pickup{X: v.X, Y: v.Y, Z: v.Z, ItemId: v.ItemId, Amount: v.Amount}
		}
	}
	return out
}

func copyLightZones(m map[string]*lightZone) map[string]*lightZone {
	out := make(map[string]*lightZone)
	for k, v := range m {
		if v != nil {
			out[k] = &lightZone{MinX: v.MinX, MinY: v.MinY, MinZ: v.MinZ, MaxX: v.MaxX, MaxY: v.MaxY, MaxZ: v.MaxZ, R: v.R, G: v.G, B: v.B, Intensity: v.Intensity}
		}
	}
	return out
}

func copyLevers(m map[string]*lever) map[string]*lever {
	out := make(map[string]*lever)
	for k, v := range m {
		if v != nil {
			out[k] = &lever{On: v.On}
		}
	}
	return out
}

func copyButtons(m map[string]*button) map[string]*button {
	out := make(map[string]*button)
	for k, v := range m {
		if v != nil {
			out[k] = &button{Pressed: v.Pressed}
		}
	}
	return out
}

func copySwitches(m map[string]*switchState) map[string]*switchState {
	out := make(map[string]*switchState)
	for k, v := range m {
		if v != nil {
			out[k] = &switchState{On: v.On}
		}
	}
	return out
}
