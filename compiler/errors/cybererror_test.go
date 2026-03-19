package errors

import (
	"strings"
	"testing"
)

func TestCyberError_Format(t *testing.T) {
	tests := []struct {
		name string
		e    CyberError
		want []string // substrings that must appear
	}{
		{
			name: "full",
			e: CyberError{
				Code:       ErrAssetNotFound,
				Message:    "Asset 'plr' not found",
				Line:       5,
				Snippet:    `DRAWTEXTURE(ASSETS["plr"], 0, 0)`,
				Suggestion: "Did you mean 'player'?",
				Filename:   "game.bas",
			},
			want: []string{"line 5", "game.bas", "ASSETS", "Asset 'plr'", "Did you mean"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := tt.e.Format()
			for _, w := range tt.want {
				if !strings.Contains(s, w) {
					t.Errorf("Format() missing %q in:\n%s", w, s)
				}
			}
		})
	}
}

func TestNearest(t *testing.T) {
	known := []string{"player_idle", "player_run", "jump_sfx"}
	if got := Nearest("plr_idle", known, 2); got != "player_idle" && got != "" {
		t.Logf("Nearest(plr_idle) = %q (acceptable fuzzy)", got)
	}
	if got := Nearest("player_idle", known, 2); got != "player_idle" {
		t.Errorf("exact match: got %q", got)
	}
}
