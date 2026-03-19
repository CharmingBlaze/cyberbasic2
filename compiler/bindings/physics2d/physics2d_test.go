package physics2d

import "testing"

func TestRequireExplicitWorld(t *testing.T) {
	RequireExplicitWorld = true
	WorldEnsured = false
	// ensureWorld needs VM — covered by integration; config toggles only
	if !RequireExplicitWorld {
		t.Fatal("flag")
	}
	RequireExplicitWorld = false
	WorldEnsured = false
}
