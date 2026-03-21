// Command foreignaudit scans compiler/bindings for RegisterForeign names and optionally
// diffs them against top-level exported functions in github.com/gen2brain/raylib-go/raylib.
package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

var reRegister = regexp.MustCompile(`RegisterForeign\(\s*"([^"]+)"`)

type Command struct {
	Name      string   `json:"name"`
	Packages  []string `json:"packages"`
	Files     []string `json:"files"`
	PrimaryPkg string  `json:"primary_package"`
}

type Inventory struct {
	GeneratedAt string    `json:"generated_at"`
	Commands    []Command `json:"commands"`
	Total       int       `json:"total"`
}

type ParityReport struct {
	RaylibModule     string   `json:"raylib_go_module"`
	RaylibVersion    string   `json:"raylib_go_version"`
	RaylibFuncs      []string `json:"raylib_go_exported_funcs"`
	InRaylibNotBound []string `json:"in_raylib_not_in_bindings_raylib"`
	InBindingsRaylib []string `json:"bound_in_compiler_bindings_raylib"`
	Note             string   `json:"note"`
}

func main() {
	root := flag.String("root", ".", "repository root")
	bindingsRel := flag.String("bindings", "compiler/bindings", "bindings path relative to root")
	outJSON := flag.String("out-json", "", "write inventory JSON (e.g. docs/generated/foreign_commands.json)")
	outMD := flag.String("out-md", "", "write inventory markdown index")
	parityJSON := flag.String("parity-json", "", "write raylib-go vs bindings/raylib diff JSON")
	skipTests := flag.Bool("skip-tests", true, "ignore *_test.go when scanning bindings")
	flag.Parse()

	absRoot, err := filepath.Abs(*root)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	bindingsDir := filepath.Join(absRoot, filepath.FromSlash(*bindingsRel))

	cmds, err := scanBindings(bindingsDir, *skipTests)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	inv := Inventory{
		GeneratedAt: timeRFC3339(),
		Commands:    cmds,
		Total:       len(cmds),
	}

	if *outJSON != "" {
		if err := writeJSON(*outJSON, inv); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		fmt.Fprintf(os.Stderr, "wrote %d commands -> %s\n", inv.Total, *outJSON)
	}

	if *outMD != "" {
		if err := writeInventoryMD(*outMD, inv); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		fmt.Fprintf(os.Stderr, "wrote markdown -> %s\n", *outMD)
	}

	if *parityJSON != "" {
		rep, err := raylibParity(cmds)
		if err != nil {
			fmt.Fprintln(os.Stderr, "parity:", err)
			os.Exit(1)
		}
		if err := writeJSON(*parityJSON, rep); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		fmt.Fprintf(os.Stderr, "wrote parity (%d raylib funcs, %d unbound) -> %s\n",
			len(rep.RaylibFuncs), len(rep.InRaylibNotBound), *parityJSON)
	}
}

func timeRFC3339() string {
	return time.Now().UTC().Format(time.RFC3339)
}

func scanBindings(bindingsDir string, skipTests bool) ([]Command, error) {
	type acc struct {
		files map[string]struct{}
	}
	byName := make(map[string]*acc)

	err := filepath.Walk(bindingsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if !strings.HasSuffix(path, ".go") {
			return nil
		}
		if skipTests && strings.HasSuffix(path, "_test.go") {
			return nil
		}
		rel, _ := filepath.Rel(bindingsDir, path)
		relSlash := filepath.ToSlash(rel)

		f, err := os.Open(path)
		if err != nil {
			return err
		}
		sc := bufio.NewScanner(f)
		for sc.Scan() {
			line := sc.Text()
			for _, m := range reRegister.FindAllStringSubmatch(line, -1) {
				name := m[1]
				a := byName[name]
				if a == nil {
					a = &acc{files: make(map[string]struct{})}
					byName[name] = a
				}
				a.files[relSlash] = struct{}{}
			}
		}
		_ = f.Close()
		return sc.Err()
	})
	if err != nil {
		return nil, err
	}

	names := make([]string, 0, len(byName))
	for n := range byName {
		names = append(names, n)
	}
	sort.Strings(names)

	out := make([]Command, 0, len(names))
	for _, n := range names {
		a := byName[n]
		fs := make([]string, 0, len(a.files))
		for f := range a.files {
			fs = append(fs, f)
		}
		sort.Strings(fs)
		pkgSet := make(map[string]struct{})
		for _, f := range fs {
			parts := strings.SplitN(f, "/", 2)
			if len(parts) > 0 {
				pkgSet[parts[0]] = struct{}{}
			}
		}
		pkgs := make([]string, 0, len(pkgSet))
		for p := range pkgSet {
			pkgs = append(pkgs, p)
		}
		sort.Strings(pkgs)
		primary := pkgs[0]
		for _, p := range pkgs {
			if p == "raylib" {
				primary = "raylib"
				break
			}
		}
		out = append(out, Command{
			Name:        n,
			Packages:    pkgs,
			Files:       fs,
			PrimaryPkg:  primary,
		})
	}
	return out, nil
}

var reTopLevelFunc = regexp.MustCompile(`^\s*func\s+([A-Z][A-Za-z0-9]*)\s*\(`)

func raylibGoModule() (path, version, dir string, err error) {
	cmd := exec.Command("go", "list", "-m", "-json", "github.com/gen2brain/raylib-go/raylib")
	out, e := cmd.Output()
	if e != nil {
		err = fmt.Errorf("go list -m -json raylib-go/raylib: %w (run go mod download)", e)
		return
	}
	var meta struct {
		Path    string `json:"Path"`
		Version string `json:"Version"`
		Dir     string `json:"Dir"`
	}
	if e := json.Unmarshal(out, &meta); e != nil {
		err = e
		return
	}
	if meta.Dir == "" {
		err = fmt.Errorf("empty Dir in go list -m -json output")
		return
	}
	return meta.Path, meta.Version, meta.Dir, nil
}

func listRaylibExportedFuncs(dir string) ([]string, error) {
	set := make(map[string]struct{})
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return err
		}
		if !strings.HasSuffix(path, ".go") {
			return nil
		}
		if strings.HasSuffix(path, "_test.go") {
			return nil
		}
		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()
		sc := bufio.NewScanner(f)
		for sc.Scan() {
			m := reTopLevelFunc.FindStringSubmatch(sc.Text())
			if len(m) == 2 {
				set[m[1]] = struct{}{}
			}
		}
		return sc.Err()
	})
	if err != nil {
		return nil, err
	}
	names := make([]string, 0, len(set))
	for n := range set {
		names = append(names, n)
	}
	sort.Strings(names)
	return names, nil
}

func raylibParity(cmds []Command) (*ParityReport, error) {
	modPath, modVer, dir, err := raylibGoModule()
	if err != nil {
		return nil, err
	}
	_ = dir
	rlFuncs, err := listRaylibExportedFuncs(dir)
	if err != nil {
		return nil, err
	}

	boundRaylib := make(map[string]struct{})
	for _, c := range cmds {
		for _, f := range c.Files {
			if strings.HasPrefix(f, "raylib/") {
				boundRaylib[c.Name] = struct{}{}
				break
			}
		}
	}
	boundNames := make([]string, 0, len(boundRaylib))
	for n := range boundRaylib {
		boundNames = append(boundNames, n)
	}
	sort.Strings(boundNames)

	var missing []string
	for _, fn := range rlFuncs {
		if _, ok := boundRaylib[fn]; !ok {
			missing = append(missing, fn)
		}
	}

	return &ParityReport{
		RaylibModule:     modPath,
		RaylibVersion:    modVer,
		RaylibFuncs:      rlFuncs,
		InRaylibNotBound: missing,
		InBindingsRaylib: boundNames,
		Note: "Top-level exported funcs in raylib-go/raylib only; methods and non-Go wrappers are excluded. " +
			"CyberBasic may expose a different name or wrap via DBP/game.",
	}, nil
}

func writeJSON(path string, v interface{}) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(v)
}

func writeInventoryMD(path string, inv Inventory) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	var b strings.Builder
	b.WriteString("# Foreign command inventory\n\n")
	b.WriteString("Auto-generated by `go run ./internal/tools/foreignaudit`. Do not edit by hand.\n\n")
	b.WriteString(fmt.Sprintf("**Total:** %d unique `RegisterForeign` names under `compiler/bindings/`.\n\n", inv.Total))

	byPkg := make(map[string][]Command)
	for _, c := range inv.Commands {
		p := c.PrimaryPkg
		byPkg[p] = append(byPkg[p], c)
	}
	pkgs := make([]string, 0, len(byPkg))
	for p := range byPkg {
		pkgs = append(pkgs, p)
	}
	sort.Strings(pkgs)

	for _, p := range pkgs {
		list := byPkg[p]
		sort.Slice(list, func(i, j int) bool { return list[i].Name < list[j].Name })
		b.WriteString(fmt.Sprintf("## %s (%d)\n\n", p, len(list)))
		for _, c := range list {
			loc := c.Files[0]
			if len(c.Files) > 1 {
				loc = strings.Join(c.Files, ", ")
			}
			b.WriteString(fmt.Sprintf("- **%s** — `%s`", c.Name, loc))
			if len(c.Packages) > 1 {
				b.WriteString(fmt.Sprintf(" _(packages: %s)_", strings.Join(c.Packages, ", ")))
			}
			b.WriteString("\n")
		}
		b.WriteString("\n")
	}

	return os.WriteFile(path, []byte(b.String()), 0644)
}
