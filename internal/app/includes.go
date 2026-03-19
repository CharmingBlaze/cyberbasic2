package app

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// PreprocessIncludes expands #include "file.bas" and IMPORT "file.bas" with file contents (relative to baseDir). seen prevents cycles.
func PreprocessIncludes(source []byte, baseDir string, seen map[string]bool) []byte {
	if seen == nil {
		seen = make(map[string]bool)
	}
	includeRe := regexp.MustCompile(`^\s*#include\s*"([^"]+)"\s*$`)
	importRe := regexp.MustCompile(`(?i)^\s*IMPORT\s*"([^"]+)"\s*$`)
	var out strings.Builder
	sc := bufio.NewScanner(strings.NewReader(string(source)))
	for sc.Scan() {
		line := sc.Text()
		var filePath string
		if m := includeRe.FindStringSubmatch(line); m != nil {
			filePath = m[1]
		} else if m := importRe.FindStringSubmatch(line); m != nil {
			filePath = m[1]
		}
		if filePath != "" {
			path := filepath.Join(baseDir, filePath)
			abs, _ := filepath.Abs(path)
			if seen[abs] {
				continue
			}
			seen[abs] = true
			inc, err := os.ReadFile(path)
			if err != nil {
				out.WriteString(line)
				out.WriteByte('\n')
				continue
			}
			incDir := filepath.Dir(path)
			inc = PreprocessIncludes(inc, incDir, seen)
			out.Write(inc)
			if len(inc) > 0 && inc[len(inc)-1] != '\n' {
				out.WriteByte('\n')
			}
			continue
		}
		out.WriteString(line)
		out.WriteByte('\n')
	}
	return []byte(out.String())
}
