// Dotmapgen prints a commented Go map stub from docs/generated/foreign_commands.json
// to help author modfacade method→foreign tables. Does not modify the repo.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"
)

type doc struct {
	Commands []struct {
		Name string `json:"name"`
	} `json:"commands"`
}

func main() {
	jsonPath := flag.String("json", "docs/generated/foreign_commands.json", "path to foreign_commands.json")
	prefix := flag.String("prefix", "", "if set, only names starting with this prefix (case-insensitive)")
	flag.Parse()
	data, err := os.ReadFile(*jsonPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "read %s: %v\n", *jsonPath, err)
		os.Exit(1)
	}
	var d doc
	if err := json.Unmarshal(data, &d); err != nil {
		fmt.Fprintf(os.Stderr, "json: %v\n", err)
		os.Exit(1)
	}
	var names []string
	pfx := strings.ToLower(*prefix)
	for _, e := range d.Commands {
		if e.Name == "" {
			continue
		}
		if pfx != "" && !strings.HasPrefix(strings.ToLower(e.Name), pfx) {
			continue
		}
		names = append(names, e.Name)
	}
	sort.Strings(names)
	fmt.Println("// Generated stub — hand-map friendly method names to RegisterForeign names.")
	fmt.Println("var exampleV2 = map[string]string{")
	for _, n := range names {
		key := strings.ToLower(n)
		fmt.Printf("\t%q: %q,\n", key, n)
	}
	fmt.Println("}")
}
