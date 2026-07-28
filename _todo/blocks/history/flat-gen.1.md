---
id: flat-gen.1
title: Flat Gen Package (gen/routes.go instead of per-dir dree.go)
status: 13
phase: v0.0.2
requires:
  - transpiler.1
  - routing.1
created: 2026-07-26
changed: 2026-07-26
---

All route handlers in gen/routes.go (package gen). No more _ "import". Solves Go import path problem with [ and ( in directory names.
