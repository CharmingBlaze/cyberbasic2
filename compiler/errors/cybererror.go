package errors

import (
	"fmt"
	"strings"
)

// ErrorCode identifies a class of CyberBasic error.
type ErrorCode int

const (
	ErrUnknown ErrorCode = iota
	ErrPhysicsBodyBeforeWorld
	ErrAssetNotFound
	ErrBeginDrawOutsideLoop
	ErrUndefinedVariable
	ErrTypeMismatch
	ErrMissingEndKeyword
	ErrDotAccess
)

// CyberError is a structured runtime/compile error with user-facing context.
type CyberError struct {
	Code       ErrorCode
	Message    string
	Line       int
	Column     int
	Snippet    string
	Suggestion string
	Filename   string
}

func (e *CyberError) Error() string {
	return e.Format()
}

// Format renders the error in the standard multi-line form.
func (e *CyberError) Format() string {
	var b strings.Builder
	fn := e.Filename
	if fn == "" {
		fn = "(source)"
	}
	if e.Line > 0 {
		fmt.Fprintf(&b, "Error on line %d in %s:\n", e.Line, fn)
		if e.Snippet != "" {
			fmt.Fprintf(&b, "  %s\n\n", e.Snippet)
		} else {
			b.WriteString("\n")
		}
	} else {
		fmt.Fprintf(&b, "Error in %s:\n\n", fn)
	}
	b.WriteString(e.Message)
	if e.Suggestion != "" {
		b.WriteString("\nFix: ")
		b.WriteString(e.Suggestion)
	}
	b.WriteString("\n")
	return b.String()
}

// Nearest returns the closest string in known within maxLevenshtein inclusive, or "" if none.
func Nearest(key string, known []string, maxLevenshtein int) string {
	if key == "" || len(known) == 0 {
		return ""
	}
	best := ""
	bestD := maxLevenshtein + 1
	kl := strings.ToLower(strings.TrimSpace(key))
	for _, k := range known {
		if k == "" {
			continue
		}
		d := levenshtein(kl, strings.ToLower(strings.TrimSpace(k)))
		if d < bestD {
			bestD = d
			best = k
		}
	}
	if bestD <= maxLevenshtein {
		return best
	}
	return ""
}

func levenshtein(a, b string) int {
	if len(a) == 0 {
		return len(b)
	}
	if len(b) == 0 {
		return len(a)
	}
	row := make([]int, len(b)+1)
	for j := 0; j <= len(b); j++ {
		row[j] = j
	}
	for i := 1; i <= len(a); i++ {
		prev := row[0]
		row[0] = i
		for j := 1; j <= len(b); j++ {
			cur := row[j]
			cost := 0
			if a[i-1] != b[j-1] {
				cost = 1
			}
			row[j] = min3(row[j]+1, row[j-1]+1, prev+cost)
			prev = cur
		}
	}
	return row[len(b)]
}

func min3(a, b, c int) int {
	if a <= b && a <= c {
		return a
	}
	if b <= c {
		return b
	}
	return c
}
