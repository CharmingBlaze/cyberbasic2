package semantic

import (
	"cyberbasic/compiler/parser"
	"errors"
	"fmt"
	"strings"
)

// QualifiedName returns the lowercase qualified name for a user function/sub (Module.Name or Name).
// It is exported for use by codegen or tests.
func QualifiedName(node parser.Node) string {
	switch n := node.(type) {
	case *parser.FunctionDecl:
		if n.ModuleName != "" {
			return strings.ToLower(n.ModuleName) + "." + strings.ToLower(n.Name)
		}
		return strings.ToLower(n.Name)
	case *parser.SubDecl:
		if n.ModuleName != "" {
			return strings.ToLower(n.ModuleName) + "." + strings.ToLower(n.Name)
		}
		return strings.ToLower(n.Name)
	default:
		return ""
	}
}

// Analyze walks the program AST and builds the symbol table and statement split.
// It collects all validation errors (e.g. duplicate TYPE, ENTITY, or function/sub) and returns
// them together via errors.Join so the user sees every issue in one run.
func Analyze(program *parser.Program) (*Result, error) {
	typeDefs := make(map[string]*parser.TypeDecl)
	entityNames := make(map[string]bool)
	userFuncs := make(map[string]bool)
	var mainStmts, decls []parser.Node
	var errs []error

	for _, stmt := range program.Statements {
		if td, ok := stmt.(*parser.TypeDecl); ok {
			key := strings.ToLower(td.Name)
			if _, exists := typeDefs[key]; exists {
				errs = append(errs, fmt.Errorf("duplicate type %q", td.Name))
			} else {
				typeDefs[key] = td
			}
		}
		if ed, ok := stmt.(*parser.EntityDecl); ok {
			key := strings.ToLower(ed.Name)
			if entityNames[key] {
				errs = append(errs, fmt.Errorf("duplicate entity %q", ed.Name))
			} else {
				entityNames[key] = true
			}
		}
		switch s := stmt.(type) {
		case *parser.ModuleStatement:
			for _, node := range s.Body {
				q := QualifiedName(node)
				if q != "" {
					if userFuncs[q] {
						errs = append(errs, fmt.Errorf("duplicate function or sub %q", q))
					} else {
						userFuncs[q] = true
					}
				}
				decls = append(decls, node)
			}
		case *parser.FunctionDecl, *parser.SubDecl:
			q := QualifiedName(s)
			if q != "" {
				if userFuncs[q] {
					errs = append(errs, fmt.Errorf("duplicate function or sub %q", q))
				} else {
					userFuncs[q] = true
				}
			}
			decls = append(decls, stmt)
		default:
			mainStmts = append(mainStmts, stmt)
		}
	}

	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}

	return &Result{
		TypeDefs:    typeDefs,
		EntityNames: entityNames,
		UserFuncs:   userFuncs,
		MainStmts:   mainStmts,
		Decls:       decls,
	}, nil
}
