package assets

import (
	"cyberbasic/compiler/errors"
	"testing"
)

func TestNearestInError(t *testing.T) {
	s := errors.Nearest("plr", []string{"player", "enemy"}, 2)
	if s == "" {
		t.Log("no fuzzy match (ok)")
	}
}
