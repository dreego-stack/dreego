---
id: codegen-errors.1
title: Replace Silent CodeGen Failures with Errors
status: 39
phase: v0.0.22
requires:
  - transpiler.1
created: 2026-07-31
changed: 2026-08-08
---

Audit CodeGen for places that return empty strings or ignore errors instead of propagating them. Examples:
- invalid else-if chains in `genTemplateNodeComp` returning ""
- ignored render errors in generated form post handlers
- typed block generation silently skipping unsupported types

All CodeGen paths should return meaningful errors so `dreego generate` fails fast with a clear message.

Done in v0.0.22: all `core/codegen*.go` template generators return `(string, error)`.
