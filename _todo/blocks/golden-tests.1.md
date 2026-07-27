---
id: golden-tests.1
title: Golden-File Tests für Generator
status: planned
phase: v0.0.5
requires:
  - transpiler.1
  - dreegotest.1
created: 2026-07-27
changed: 2026-07-27
---

Templ-inspiriert: ~50 Test-Subdirectories mit `.dreego` Input + erwartetem `_dreego.go` Output. Golden-File Pattern: Parse → Generate → Compare. Jedes Template-Konstrukt ({#if}, {#each}, section, layout) bekommt eigene Test-Fixtures. Ersetzt manuelle CI-Checks.
