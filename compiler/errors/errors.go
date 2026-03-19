// Package errors provides pretty-printing of compiler errors with source context.
// It uses github.com/rhysd/locerr to show the offending line and a caret when
// the error message contains "line N: ..." (as produced by codegen.errWithLine).
package errors

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/rhysd/locerr"
)

// linePrefix matches "line 123: " at the start of a segment (we look for the last occurrence in the chain).
var lineRegex = regexp.MustCompile(`line (\d+): (.+)`)

// PrettyPrint writes the compilation error to w. If the error chain contains
// a "line N: message" segment (from codegen.errWithLine), it uses locerr to
// print a source snippet and the offending line; otherwise it prints the
// error as-is. source is the full file content; filename is used only in the
// locerr source name when non-empty.
func PrettyPrint(w io.Writer, source string, filename string, err error) {
	if err == nil {
		return
	}
	msg := err.Error()
	line, rest := parseLineFromError(msg)
	if line <= 0 || source == "" {
		fmt.Fprintf(w, "Compilation error: %v\n", err)
		return
	}
	_ = filename // locerr.NewDummySource does not take a name; snippet still shows line/col
	src := locerr.NewDummySource(source)
	offset := lineToOffset(source, line)
	pos := locerr.Pos{
		Offset: offset,
		Line:   line,
		Column: 1,
		File:   src,
	}
	locErr := locerr.ErrorAt(pos, strings.TrimSpace(rest))
	if f, ok := w.(*os.File); ok {
		locErr.PrintToFile(f)
	} else {
		fmt.Fprintf(w, "Compilation error: %v\n", locErr)
	}
}

// parseLineFromError finds the last "line N: message" in the error chain and returns (N, message).
// If none found, returns (0, "").
func parseLineFromError(errMsg string) (line int, message string) {
	// Unwrap chain: "file: phase: line 5: msg" -> we want the last "line N: msg"
	idx := strings.LastIndex(errMsg, "line ")
	if idx < 0 {
		return 0, ""
	}
	sub := errMsg[idx:]
	if !lineRegex.MatchString(sub) {
		return 0, ""
	}
	matches := lineRegex.FindStringSubmatch(sub)
	if len(matches) != 3 {
		return 0, ""
	}
	var n int
	if _, err := fmt.Sscanf(matches[1], "%d", &n); err != nil {
		return 0, ""
	}
	return n, matches[2]
}

func lineToOffset(source string, line int) int {
	currentLine := 1
	for i, r := range source {
		if currentLine == line {
			return i
		}
		if r == '\n' {
			currentLine++
		}
	}
	return 0
}
