# Stub and backlog (not part of the CI health gate)

CI verifies **build, tests, example compile (`--lint`), foreign-inventory parity, and golangci-lint**. The items below are **intentional follow-ups** — track them as separate issues or milestones when you prioritize “full featured” work.

| Area | Status | Notes |
|------|--------|--------|
| **effect / camera.fx** | Stub / queue | Real post-processing needs a render-target path or render graph; see [`COMMAND_COVERAGE.md`](COMMAND_COVERAGE.md). |
| **Raylib parity** | Partial | [`generated/raylib_parity.json`](generated/raylib_parity.json) lists unbound `raylib-go` symbols; add in small, game-relevant tranches. |
| **BTREE / Phase 12 AI** | Deferred | Parser/VM scope; see [`BTREE_PHASE12.md`](BTREE_PHASE12.md) and [`REFACTOR_PLAN.md`](../REFACTOR_PLAN.md). |
| **Shader presets** | Minimal | `shader.pbr` / `toon` / `dissolve` are teaching shaders, not full PBR; custom GLSL via `shader.load`. |

When opening GitHub issues, prefer one ticket per row (or per tranche) so PRs stay reviewable.
