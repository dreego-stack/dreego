---
version: patch
---

- Chore: remove the duplicate `exprKind` and source-position helpers from `internal/transpiler` and use the canonical `internal/transpiler/ir` implementations
- Chore: drop the pass-through helper wall in `internal/transpiler/codegen_html.go` and call the `ir`, `output`, and `head` packages directly
