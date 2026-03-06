// Package dbp - Networking: DBP-style wrappers over the net package.
//
// The net package (registered in main) provides:
//   - NetConnect(ip$, port) - Connect to server, returns connectionId$
//   - NetSend(connectionId$, data$) - Send text
//   - NetReceive(connectionId$) - Receive next message (or empty)
//   - NetDisconnect(connectionId$) - Close connection
//   - NetIsConnected(connectionId$) - 1 if connected, 0 otherwise
//   - NetIsServer() - 1 if hosting, 0 otherwise
//   - NetPlayerID() - First connection id (or "" if none)
//   - NetPing(connectionId$) / NetLatency(connectionId$) - RTT in ms
//   - Host(port) / Accept(serverId$) - Server API
package dbp

import (
	"cyberbasic/compiler/vm"
)

// registerNet registers DBP-style networking. The net package provides NetConnect,
// NetSend, NetReceive, NetDisconnect, NetIsServer, NetPlayerID, NetPing, NetLatency.
// This module documents the API; net package is registered separately in main.
func registerNet(v *vm.VM) {
	_ = v // reserved for future DBP-specific net wrappers
}
