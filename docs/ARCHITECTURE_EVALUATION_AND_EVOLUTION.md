# CyberBASIC2 — Architecture Evaluation & Evolution Plan

This document evaluates the current architecture of CyberBASIC2 (cloned from CharmingBlaze/cyberbasic2), identifies missing or incomplete systems, and outlines concrete steps to evolve it into a **clean, powerful, modular game-focused language** with strong compiler-engineering practices.

---

## 1. Current State Summary

### What’s in the repo

- **Compiler pipeline:** Lexer → Parser → (single-pass) Code generation → Bytecode Chunk.
- **VM:** Custom bytecode VM with stack-based execution, fibers/coroutines, foreign function calls.
- **Runtime:** Game loop (mainloop / SYNC), fixed timestep, implicit OnUpdate/OnDraw handlers, render queues (2D/3D/UI).
- **Bindings:** raylib (graphics, input, audio, raygui), Box2D, Bullet, ECS, net, std, scene, game, dbp, terrain, water, vegetation, navigation, indoor, objects, procedural, world, SQL, aseprite, model (OBJ/GLTF/FBX).
- **Docs:** Extensive guides (2D/3D, multiplayer, GUI, ECS, tutorials, DBP parity, troubleshooting).

### Build and run

```bash
cd cyberbasic2
go build -o cyberbasic .
./cyberbasic examples/first_game.bas
```

Build succeeds; the project is in active use with a single-binary, Go-only stack.

---

## 2. Architecture Evaluation

### Strengths

| Area | Assessment |
|------|------------|
| **Lexer** | Isolated in `compiler/lexer/` (lexer.go, token.go). Token-only; no parsing. Clear boundary. |
| **Parser** | Isolated in `compiler/parser/` (parser.go, ast.go, parser_*.go). Consumes tokens, produces AST. No codegen. |
| **VM** | Well-structured: bytecode.go (opcodes), vm.go (execution), vm_ops.go, vm_run.go, vm_foreign.go, vm_fibers.go, runtime_iface.go. Clear separation from compiler. |
| **Runtime** | Separate package `compiler/runtime/` (loop, renderer, camera, assets, animation, etc.). Bindings in `compiler/bindings/<pkg>`. |
| **Game focus** | DBP-style commands, 2D/3D APIs, physics, ECS, multiplayer, terrain, water, vegetation. Good fit for “BASIC for games.” |

### Gaps vs. desired “modular compiler” design

The document **MODULAR_COMPILER_ARCHITECTURE.md** describes a five-stage pipeline. Current reality:

| Stage | Doc plan | Current state |
|-------|----------|----------------|
| **1. Lexer** | `compiler/lexer/` | ✅ Implemented as specified. |
| **2. Parser** | `compiler/parser/` | ✅ Implemented as specified. |
| **3. Semantic analysis** | `compiler/semantic/` (semantic.go, symbol.go) | ❌ **Missing.** Type defs, entity names, user func/sub collection, and scope-like logic live inside `compiler.go` / `generateCode()` (first pass). No dedicated symbol table or semantic package. |
| **4. Code generation** | `compiler/codegen/` (codegen.go, stmt.go, expr.go) + slim `compiler.go` | ❌ **Not split.** All codegen lives in `compiler` package: `compiler.go`, `codegen_statements.go`, `codegen_expr.go`, `codegen_call.go`, `codegen_util.go`. Driver and codegen are one blob. |
| **5. Runtime / stdlib** | VM + bindings | ✅ Largely as specified (VM + bindings are separate from compilation). |

So the **compiler** is currently a **three-stage monolith**: Lexer → Parser → (semantics + codegen mixed in one pass in `compiler`). That blocks:

- Independent testing of semantic analysis (e.g. type/scope checks without generating code).
- Clear dependency rules (e.g. codegen depending only on AST + symbol table, not on ad-hoc maps in Compiler).
- Easier evolution of the language (add new semantic checks or optimizations in one place).

### Other architectural observations

- **main.go** is large: CLI, preprocessing (#include/IMPORT), compile, VM setup, and registration of many bindings in one place. A small `cmd/cyberbasic/main.go` plus a shared “run pipeline” could improve clarity.
- **Preprocessing** (include/import) is in main.go as a string pass before lexing; it’s simple but lives outside the compiler package.
- **Module path:** `go.mod` uses `module cyberbasic`; imports are `cyberbasic/compiler/...`. Repo name is cyberbasic2; consider aligning module name if you want consistency (optional).

---

## 3. Using Go libraries to support the fix

Yes — we can use a few Go libraries to make the refactor and error reporting better. The **modular split** (semantic + codegen packages) is still done by moving code into new packages; there is no “compiler framework” in Go for custom languages. These libraries help around the edges:

| Need | Library | Use |
|------|---------|-----|
| **Pretty errors with source snippets** | `github.com/rhysd/locerr` | After semantic/codegen have line/column (from AST or tokens), wrap errors in `locerr.ErrorAt(pos, msg)` or `locerr.ErrorIn(start, end, msg)`. Gives snippet + caret, optional notes. |
| **Collect multiple errors in one pass** | **stdlib** `errors.Join` (Go 1.20+) | Semantic pass can collect many errors (e.g. duplicate type, undefined ref) and return `errors.Join(err1, err2, ...)` so the user sees all at once. No new dependency. |
| **Test assertions** | `github.com/stretchr/testify` (optional) | Cleaner semantic/codegen tests with `require.NoError`, `assert.Equal`; not required but improves readability. |
| **Parser generator** | `github.com/mna/pigeon` (optional, later) | Only if you ever replace the hand-written parser with a PEG grammar. Not needed for the current fix; keep existing parser. |

**What not to use**

- **go/ast, go/parser, go/types** — For Go source only; we have a custom BASIC AST, so semantic analysis stays in our own `compiler/semantic` package.
- **A “symbol table” library** — Go doesn’t have a standard one for custom ASTs; a simple `map[string]T` and a slice of scopes (stack of maps) in `symbol.go` is enough.

**Suggested order**

1. Do Phase A (semantic + codegen packages, slim driver) with stdlib only.
2. Add **locerr** when attaching positions to semantic/codegen errors (Phase B).
3. Use **errors.Join** in the semantic pass as soon as you want “report all errors” instead of fail-fast.
4. Add **testify** in tests if you prefer it for compiler/semantic/codegen tests.

---

## 4. Missing or Incomplete Systems

### 4.1 Compiler / language (from MODULAR_COMPILER_ARCHITECTURE + ROADMAP)

- **Dedicated semantic analysis**
  - Symbol table (scopes, function signatures, type definitions, entity names).
  - Move “first pass” from `generateCode()` into `compiler/semantic/`: typeDefs, entityNames, userFuncs (and any future scope/type checks).
  - Output: annotated AST or a semantic result struct consumed by codegen only.

- **Dedicated codegen package**
  - `compiler/codegen/`: entry (e.g. `Emit(ast, semanticInfo) -> *vm.Chunk`), stmt.go, expr.go; optionally call.go and util in codegen.
  - `compiler.go` reduced to: Lex → Parse → Semantic.Analyze(ast) → Codegen.Emit(ast, semanticInfo) → return Chunk.

- **Error reporting**
  - Line/column in compiler errors, “did you mean?” where feasible, clearer type/signature messages (ROADMAP: “Better error messages”).

### 4.2 Runtime / engine (from ROADMAP_IMPLEMENTATION.md)

- **Shadows:** Point-light and spot-light shadows; cascaded shadow maps (directional exists).
- **Level lighting:** GLTF KHR_lights_punctual not wired to level loader.
- **Asset cache:** LoadLevel/LoadLevelWithHierarchy not using parsed-model cache; LoadPrefab not cache-backed; BuildModel texture reuse (use resource cache).
- **3D physics:** Native Bullet backend optional; 3D constraint joints (hinge, slider) stubbed; CreateStaticMesh3D placeholder; capsule/swept queries approximated.
- **Multiplayer:** Lockstep, rollback, prediction, matchmaking, interest management not implemented (design in MULTIPLAYER_DESIGN.md).

### 4.3 Game / DX (from ROADMAP.md and DBP_GAP.md)

- **UI:** Expand BeginUI/Label/Button beyond stubs; solid immediate-mode or raygui usage for menus and tools.
- **Debugger:** Breakpoints, step/next, watch, call stack, bytecode dump, ECS inspector.
- **Stdlib strings/math (DBP parity):** LEFT$, RIGHT$, MID$, LEN, CHR$, ASC, STR$, VAL, RND, INT; CopyFile, Dir/ListDir, ExecuteFile (see DBP_GAP.md).
- **Strings/arrays (ROADMAP):** Dynamic arrays, slicing, INSTR, string interpolation.

### 4.4 Tooling and quality

- **CI:** GitHub Actions for build and tests (ROADMAP).
- **Tests:** Parser and VM tests exist; compiler/codegen and semantic tests would benefit from the modular split.
- **Package/module distribution:** No package manager yet (longer-term).

---

## 5. Evolution Steps — Prioritized

### Phase A: Compiler modularity (strong compiler-engineering base)

1. **Add `compiler/semantic`**
   - Introduce `symbol.go`: symbol table, scopes, function signatures, type definitions, entity names.
   - Introduce `semantic.go`: `Analyze(ast) -> (*Result, error)`; move from compiler’s first pass: typeDefs, entityNames, userFuncs collection and any scope/type checks.
   - Keep AST in parser package; semantic package only depends on parser (AST types).

2. **Add `compiler/codegen`**
   - `codegen.go`: entry `Emit(program *parser.Program, semantic *semantic.Result) (*vm.Chunk, error)`; orchestrate compilation of main + decls + event handlers + patches.
   - Move from `compiler`: all `compileStatement`, `compileExpression`, `compileDecl`, `compileCall`, `compile*Statement`, `compile*Expr`, `emitFrameWrap`, `emitHybridLoopBody`, and helpers (e.g. `qualifiedName`, `bodyCallsUserSub`) into `codegen` (e.g. stmt.go, expr.go, call.go, util.go).
   - Codegen depends only on: parser (AST), vm (Chunk, opcodes), semantic (Result). No lexer.

3. **Slim `compiler.go` to driver only**
   - `Compile(source string)`: Lex → Parse → Semantic.Analyze(ast) → Codegen.Emit(ast, semanticResult) → return Chunk.
   - Remove all generateCode/compile* logic from compiler package; compiler depends on lexer, parser, semantic, codegen, vm.

4. **Tests and docs**
   - Update `compiler_test.go` (and imports) so tests still pass.
   - Add semantic tests (e.g. duplicate type/entity names, undefined refs if you add that check).
   - Optionally add codegen tests (e.g. golden bytecode for small programs).
   - Update MODULAR_COMPILER_ARCHITECTURE.md to “Implementation status: done” for semantic and codegen.

### Phase B: Language and errors

5. **Error messages**
   - Attach line/column (from tokens or AST) to compiler errors; surface in semantic and codegen.
   - Add simple “did you mean?” for unknown identifiers (e.g. nearest name from symbol table).

6. **Stdlib and DBP parity (high-value, low-risk)**
   - Implement missing std functions from DBP_GAP.md in `compiler/bindings/std/`: Left, Right, Mid, Len, Chr, Asc, Str, Val, Rnd, Int, CopyFile, Dir/ListDir, ExecuteFile (and document in DBP_COMPAT.md).

7. **Strings/arrays (incremental)**
   - Dynamic arrays and slicing as per ROADMAP; string interpolation and INSTR when the rest of the pipeline is stable.

### Phase C: Runtime and game features

8. **Asset cache and level loading**
   - Route LoadLevel / LoadLevelWithHierarchy through assets cache; refcount prefabs; use resource cache for BuildModel textures (see ROADMAP_IMPLEMENTATION.md).

9. **Shadows and lighting**
   - Point/spot shadow passes; cascaded shadow maps; wire GLTF KHR_lights_punctual in level loader.

10. **Physics completeness**
    - Box2D/Bullet: joints, callbacks, raycasts (fill stubs); 3D constraints or clear “unsupported” messages.

11. **UI and debugging**
    - Expand UI (BeginUI, Label, Button, etc.) for menus and tools.
    - First-cut debugger: breakpoints, step, stack trace, bytecode dump (even text-based).

### Phase D: Ecosystem and polish

12. **CI and examples**
    - GitHub Actions: build and test on push.
    - More examples: platformer, top-down shooter, 3D scene with physics, UI demo, ECS example.

13. **DX**
    - VSCode extension (syntax, run/compile); optional REPL; later, package/module distribution.

---

## 6. Dependency Rules (target state)

- **Lexer:** no internal compiler deps (stdlib only).
- **Parser:** depends only on lexer (tokens).
- **Semantic:** depends only on parser (AST). No VM, no codegen.
- **Codegen:** depends on parser (AST), vm (Chunk), semantic (Result). No lexer.
- **Driver (compiler.go):** depends on lexer, parser, semantic, codegen, vm; orchestrates only.
- **Runtime / bindings:** depend on VM for registration; no compiler internals.

---

## 7. Summary

| Goal | Current | Target |
|------|---------|--------|
| **Modular compiler** | Lexer + Parser + monolith (semantics+codegen in compiler) | Lexer → Parser → Semantic → Codegen → Chunk, with clear packages |
| **Game-focused** | Strong (2D/3D, physics, ECS, net, terrain, etc.) | Same; add UI, debugger, DBP stdlib, asset cache |
| **Simple** | BASIC-like syntax, single binary | Keep; improve errors and onboarding |
| **Fast** | Bytecode VM, single process | Keep; optional JIT/REPL later |
| **Modular** | Bindings are modular; compiler is not | Full pipeline modularity + optional package distribution later |

Implementing **Phase A** (semantic package + codegen package + slim driver) gives a clean, modular compiler foundation. Phases B–D then build on that for a modern, game-focused BASIC that is simple, fast, and built with strong compiler-engineering practices.

---

## References

- [MODULAR_COMPILER_ARCHITECTURE.md](MODULAR_COMPILER_ARCHITECTURE.md) — Target compiler layout.
- [ROADMAP_IMPLEMENTATION.md](ROADMAP_IMPLEMENTATION.md) — Current implementation status and gaps.
- [ROADMAP.md](../ROADMAP.md) — Forward-looking priorities.
- [DBP_GAP.md](DBP_GAP.md) — Missing DBP-style stdlib commands.
- [DOCUMENTATION_INDEX.md](DOCUMENTATION_INDEX.md) — Full doc index.
