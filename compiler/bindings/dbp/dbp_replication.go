// Package dbp - Replication: DBP-style wrappers over game package replication.
//
// The game package provides:
//   - ReplicatePosition(entityId$) - Mark entity position for sync
//   - ReplicateRotation(entityId$) - Mark entity rotation for sync
//   - ReplicateScale(entityId$) - Mark entity scale for sync
//   - ReplicateValue(entityId$, varName$) - Alias for ReplicateVariable
package dbp

import (
	"cyberbasic/compiler/vm"
)

// registerReplication documents the replication API. The game package registers
// ReplicatePosition, ReplicateRotation, ReplicateScale, ReplicateValue.
func registerReplication(v *vm.VM) {
	_ = v // game package registers these
}
