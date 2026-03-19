package runtime

import "regexp"

// Window run mode for programs without mandatory InitWindow.
type WindowMode int

const (
	// ModeExplicit: source calls InitWindow — user drives mainloop; do not auto-run implicit loop.
	ModeExplicit WindowMode = iota
	// ModeImplicit: ON UPDATE / ON DRAW (or Sub OnUpdate/OnDraw) without InitWindow — runtime opens window and loops.
	ModeImplicit
	// ModeConsole: no InitWindow and no implicit handlers — stdout-only programs.
	ModeConsole
)

var (
	reInitWindow = regexp.MustCompile(`(?i)initwindow\s*\(`)
	reOnUpdate   = regexp.MustCompile(`(?i)on\s+update`)
	reOnDraw     = regexp.MustCompile(`(?i)on\s+draw`)
	reSubOnUp    = regexp.MustCompile(`(?i)sub\s+onupdate`)
	reSubOnDr    = regexp.MustCompile(`(?i)sub\s+ondraw`)
)

// DetectWindowMode classifies source for post-run implicit loop behaviour.
func DetectWindowMode(source string) WindowMode {
	if reInitWindow.MatchString(source) {
		return ModeExplicit
	}
	if reOnUpdate.MatchString(source) || reOnDraw.MatchString(source) ||
		reSubOnUp.MatchString(source) || reSubOnDr.MatchString(source) {
		return ModeImplicit
	}
	return ModeConsole
}
