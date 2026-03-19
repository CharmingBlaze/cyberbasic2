package semantic

import "cyberbasic/compiler/parser"

// Result holds the read-only output of semantic analysis for use by codegen.
// Codegen must not mutate any of these fields.
type Result struct {
	// TypeDefs maps UDT name (lowercase) to TYPE definition.
	TypeDefs map[string]*parser.TypeDecl
	// EntityNames is the set of entity names (lowercase) from each ENTITY decl.
	EntityNames map[string]bool
	// UserFuncs is the set of qualified Sub/Function names (lowercase) for call resolution.
	UserFuncs map[string]bool
	// MainStmts are top-level statements that are not TYPE/ENTITY/MODULE/FUNCTION/SUB (executed first).
	MainStmts []parser.Node
	// Decls are FUNCTION/SUB and MODULE body members in order (compiled after main, jump over).
	Decls []parser.Node
}
