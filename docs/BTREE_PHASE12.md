# BTREE (Phase 12) — deferred mini-plan

**Status:** Not implemented in the main codebase path. Roadmap C completes **navigation delegation** and a **patrol example** first.

## Why defer

- BTREE requires **parser / AST** extensions (or a data-driven runtime format), a **VM runner**, and tests — significantly larger than the `ai.*` → `navigation.*` facade.
- Patrol and agent flows should stay stable before layering behavior trees.

## Suggested next steps (when you pick this up)

1. Choose **syntax** (inline BTREE blocks vs. JSON/tables loaded at runtime).
2. Add a minimal **tick** runner that calls into existing `navigation` / agent updates.
3. Land `aisys_test.go` integration tests for tree evaluation without GL where possible.

See **[REFACTOR_PLAN.md](../REFACTOR_PLAN.md)** Phase 12 for the original scope.
