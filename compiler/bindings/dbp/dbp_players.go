// Package dbp: Player commands - explicit player state for multiplayer.
//
// Players have position and rotation stored explicitly. Use for characters
// that need deterministic, sync-friendly state.
//
// Commands:
//   - MakePlayer(id): Create player with default position
//   - SetPlayerPosition(id, x, y, z): Set position
//   - SetPlayerAngle(id, pitch, yaw, roll): Set rotation
//   - MovePlayer(id, x, y, z): Add to position
//   - TurnPlayer(id, pitch, yaw, roll): Add to rotation
//   - SyncPlayer(id): Sync player state to network (placeholder)
package dbp

import (
	"fmt"
	"sync"

	"cyberbasic/compiler/vm"
)

type playerState struct {
	x, y, z     float32
	pitch, yaw  float32
	roll        float32
}

var (
	players   = make(map[int]*playerState)
	playersMu sync.Mutex
)

// registerPlayers adds MakePlayer, SetPlayerPosition, SetPlayerAngle, MovePlayer, TurnPlayer, SyncPlayer.
func registerPlayers(v *vm.VM) {
	v.RegisterForeign("MakePlayer", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("MakePlayer(id) requires 1 argument")
		}
		id := toInt(args[0])
		playersMu.Lock()
		players[id] = &playerState{}
		playersMu.Unlock()
		return nil, nil
	})

	v.RegisterForeign("SetPlayerPosition", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("SetPlayerPosition(id, x, y, z) requires 4 arguments")
		}
		id := toInt(args[0])
		x, y, z := toFloat32(args[1]), toFloat32(args[2]), toFloat32(args[3])
		playersMu.Lock()
		if p, ok := players[id]; ok {
			p.x, p.y, p.z = x, y, z
		}
		playersMu.Unlock()
		return nil, nil
	})

	v.RegisterForeign("SetPlayerAngle", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("SetPlayerAngle(id, pitch, yaw, roll) requires 4 arguments")
		}
		id := toInt(args[0])
		p, y, r := toFloat32(args[1]), toFloat32(args[2]), toFloat32(args[3])
		playersMu.Lock()
		if pl, ok := players[id]; ok {
			pl.pitch, pl.yaw, pl.roll = p, y, r
		}
		playersMu.Unlock()
		return nil, nil
	})

	v.RegisterForeign("MovePlayer", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("MovePlayer(id, x, y, z) requires 4 arguments")
		}
		id := toInt(args[0])
		dx, dy, dz := toFloat32(args[1]), toFloat32(args[2]), toFloat32(args[3])
		playersMu.Lock()
		if p, ok := players[id]; ok {
			p.x += dx
			p.y += dy
			p.z += dz
		}
		playersMu.Unlock()
		return nil, nil
	})

	v.RegisterForeign("TurnPlayer", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("TurnPlayer(id, pitch, yaw, roll) requires 4 arguments")
		}
		id := toInt(args[0])
		dp, dy, dr := toFloat32(args[1]), toFloat32(args[2]), toFloat32(args[3])
		playersMu.Lock()
		if p, ok := players[id]; ok {
			p.pitch += dp
			p.yaw += dy
			p.roll += dr
		}
		playersMu.Unlock()
		return nil, nil
	})

	v.RegisterForeign("SyncPlayer", func(args []interface{}) (interface{}, error) {
		// Placeholder for net sync - would push player state to network
		return nil, nil
	})
}
