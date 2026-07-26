---
id: flat-gen.1
title: Flat Gen-Package (gen/routes.go statt per-dir dree.go)
status: 13
phase: v0.0.2
requires:
  - transpiler.1
  - routing.1
created: 2026-07-26
changed: 2026-07-26
---

Alle Route-Handler in gen/routes.go (package gen). Kein _ "import" mehr. Löst Go-Import-Path-Problem mit [ und ( in Verzeichnisnamen.
